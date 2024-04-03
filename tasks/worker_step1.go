package tasks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gemini/db"
	"gemini/store"
	"gemini/tpl"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"log"
	"strings"
	"text/template"
)

func DoMerge(msg []byte, key string) bool {
	data := make(map[string]interface{})
	err := json.Unmarshal(msg, &data)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return true
	}
	profile, _ := data["profile"].(string)
	url, _ := data["url"].(string)
	result := store.GeminiResult{
		GeminiKey:   key,
		ProfileData: "",
		CVURL:       url,
		CVData:      profile,
	}
	exists, err := result.CvExists(db.Client(), url)
	if exists {
		fmt.Println("cv url has exists")
		return true
	}
	id, err := result.Create(db.Client())
	if err != nil {
		fmt.Println("insert data error:", err)
		return false
	}
	step1 := GeminiStep1Merge(profile, "", key)
	if step1 == "" {
		return false
	}
	jsonResult := GetJSON(step1)
	result.GeminiStep1 = jsonResult
	result.ID = id
	err = result.Update(db.Client())
	fmt.Println("update gemini result", jsonResult)
	if err != nil {
		// todo 送到任务二的队列进行处理
	}
	return err == nil
}

func GetJSON(s string) string {
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start == -1 || end == -1 || start >= end {
		return ""
	}
	return s[start : end+1]
}

func GeminiStep1Merge(ocrCv, profileCv, key string) string {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(key))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	model := client.GenerativeModel("gemini-1.0-pro")
	model.SetTemperature(0.9)
	model.SetTopK(1)
	model.SetTopP(1)
	model.SetMaxOutputTokens(2048)
	model.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockMediumAndAbove,
		},
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockMediumAndAbove,
		},
		{
			Category:  genai.HarmCategorySexuallyExplicit,
			Threshold: genai.HarmBlockMediumAndAbove,
		},
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockMediumAndAbove,
		},
	}
	content := parseContent(ocrCv, profileCv)
	resp, err := model.GenerateContent(ctx, genai.Text(content))
	if err != nil {
		fmt.Println("call gemini error:", err)
		return ""
	}
	candidates := resp.Candidates
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
		errorMsg, _ := json.Marshal(resp)
		fmt.Println("step1 call gemini response:", string(errorMsg))
	}()
	if resp == nil || len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return ""
	}
	part := candidates[0].Content.Parts[0]
	return fmt.Sprintf("%s", part)
}

func parseContent(ocrCv, profileCv string) string {
	temple, err := template.New("gemini-step1").Parse(tpl.STEP1)
	if err != nil {
		panic(err)
	}
	data := Data{OcrCV: ocrCv, ProfileCV: profileCv}
	var buf bytes.Buffer
	err = temple.Execute(&buf, data)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

type Data struct {
	OcrCV     string
	ProfileCV string
}
