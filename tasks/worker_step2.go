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

func DoDeduce(msg []byte, key string) bool {
	data := make(map[string]int64)
	err := json.Unmarshal(msg, &data)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return true
	}
	id, _ := data["id"]
	result := &store.GeminiResult{
		GeminiKey: key,
	}
	result, err = result.QueryById(db.Client(), id)
	if err != nil {
		fmt.Printf("query step1 result error:%s\r\n", err)
	}
	if result.GeminiStep1 == "" {
		fmt.Printf("id: %d not have step1 result\r\n", id)
		return true
	}
	step2 := geminiStep2Deduce(result.GeminiStep1, key)
	if step2 == "" {
		return false
	}
	jsonResult := step2ResultToJson(step2)
	result.GeminiStep1 = jsonResult
	result.ID = id
	err = result.Step2Update(db.Client())
	fmt.Println("update gemini step2 result", jsonResult)
	return err == nil
}

func step2ResultToJson(s string) string {
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start == -1 || end == -1 || start >= end {
		return ""
	}
	return s[start : end+1]
}

func geminiStep2Deduce(step1Result, key string) string {
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
	content := step2ContentBuilder(step1Result)
	resp, err := model.GenerateContent(ctx, genai.Text(content))
	if err != nil {
		fmt.Println("call gemini error:", err)
		return ""
	}
	candidates := resp.Candidates
	defer func() string {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
		errorMsg, _ := json.Marshal(resp)
		fmt.Println("step2 call gemini response:", string(errorMsg))
		return string(errorMsg)
	}()
	if resp == nil || len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return ""
	}
	part := candidates[0].Content.Parts[0]
	return fmt.Sprintf("%s", part)
}

func step2ContentBuilder(step1Result string) string {
	temple, err := template.New("gemini-step2").Parse(tpl.STEP2)
	if err != nil {
		panic(err)
	}
	data := Step1Result{Step1Result: step1Result}
	var buf bytes.Buffer
	err = temple.Execute(&buf, data)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

type Step1Result struct {
	Step1Result string
}
