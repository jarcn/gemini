package tasks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gemini/db"
	"gemini/starter/kfk_product"
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
		log.Println("Error decoding JSON:", err)
		return true
	}
	profile, _ := data["profile"].(string)
	url, _ := data["url"].(string)
	cvData := store.CvData{}
	cv, err := cvData.GetCvByUrl(db.Client(), url)
	if err != nil || cv == nil {
		log.Printf("url:%s not have data \r", url)
		return false
	}
	result := store.GeminiResult{
		GeminiKey:   key,
		ProfileData: profile,
		CVURL:       url,
		CVData:      cv.ResumeMsg,
		Type:        "ZH",
	}
	exists, err := result.CvExists(db.Client(), url)
	if exists != nil {
		log.Println("cv url has exists")
		doStep2(exists.ID)
		return true
	}
	id, err := result.Create(db.Client())
	if err != nil {
		log.Println("insert data error:", err)
		return false
	}
	step1 := GeminiStep1Merge(profile, cv.ResumeMsg, key)
	if step1 == "error" {
		return true
	}
	if step1 == "" {
		return false
	}
	jsonResult := GetJSON(step1)
	result.GeminiStep1 = jsonResult
	result.ID = id
	err = result.Update(db.Client())
	log.Println("update gemini result", jsonResult)
	doStep2(id)
	return err == nil
}

func SyncDoMerge(msg []byte, key string) *store.GeminiResult {
	data := make(map[string]interface{})
	err := json.Unmarshal(msg, &data)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		return nil
	}
	profile, _ := data["profile"].(string)
	url, _ := data["url"].(string)
	cvData := store.CvData{}
	cv, err := cvData.GetCvByUrl(db.Client(), url)
	if err != nil || cv == nil {
		log.Printf("url:%s not have data \r", url)
		return nil
	}
	result := store.GeminiResult{
		GeminiKey:   key,
		ProfileData: profile,
		CVURL:       url,
		CVData:      cv.ResumeMsg,
		Type:        "ZH",
	}
	exists, err := result.CvExistsN(db.Client(), url)
	if exists != nil {
		return exists
	}
	id, err := result.Create(db.Client())
	if err != nil {
		log.Println("insert data error:", err)
		return nil
	}
	step1 := GeminiStep1Merge(profile, cv.ResumeMsg, key)
	if step1 == "error" {
		return nil
	}
	if step1 == "" {
		return nil
	}
	jsonResult := GetJSON(step1)
	result.GeminiStep1 = jsonResult
	result.ID = id
	err = result.Update(db.Client())
	log.Println("update gemini result", jsonResult)
	return &result
}

func doStep2(id int64) {
	m := make(map[string]int64)
	m["id"] = id
	msgData, _ := json.Marshal(m)
	kfk_product.SendMsg(db.KafkaBrokers, db.Step2Topic, msgData) //送到任务二的队列进行处理
}

func GetJSON(s string) string {
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start == -1 || end == -1 || start >= end {
		return ""
	}
	return s[start : end+1]
}

func GeminiStep1Merge(profileCv, ocrCv, key string) string {
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
	if resp == nil || err != nil {
		log.Println("step1 gemini response data is null", err)
		return "error"
	}
	errorMsg, _ := json.Marshal(resp)
	reason := resp.Candidates[0].FinishReason
	if reason == 0 || reason == 2 || reason == 3 || reason == 4 || reason == 5 {
		log.Println("step1 call gemini response:", string(errorMsg))
		return "error"
	}
	candidates := resp.Candidates
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from panic:", r)
		}
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
	var out bytes.Buffer
	err = json.Compact(&out, []byte(profileCv))
	if err != nil {
		fmt.Println("Error compacting JSON:", err)
		return profileCv
	}
	data := Data{OcrCV: ocrCv, ProfileCV: out.String()}
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
