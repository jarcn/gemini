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
	"text/template"
)

func DoDeduce(msg []byte, key string, isCvData bool) bool {
	client := db.Client()
	data := make(map[string]int64)
	err := json.Unmarshal(msg, &data)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		return true
	}
	id, _ := data["id"]
	result := &store.GeminiResult{
		GeminiKey: key,
		ID:        id,
	}
	result, err = result.QueryById(client, id)
	if err != nil {
		log.Printf("query step1 result error:%s\r\n", err)
		return false
	}
	if result.GeminiStep1 == "" {
		log.Printf("id:%d,url:%s not hava step1 result,start step1 task", result.ID, result.CVURL)
		step1 := GeminiStep1Merge(result.ProfileData, result.CVData, key)
		if step1 == "error" {
			return true
		}
		if step1 == "" {
			return false
		}
		jsonResult := GetJSON(step1)
		result.GeminiStep1 = jsonResult
		result.ID = id
		err = result.Update(client)
		log.Printf("update step1 gemini result length:%d\r\n", len(jsonResult))
	}
	if isCvData {
		if result.CVData == "" {
			log.Printf("id: %d not have cv data \r\n", id)
			return true
		}
		step2ByCVData := GeminiStep2Deduce(result.CVData, key) //使用原始CV数据跑任务二
		if step2ByCVData == "" {
			return false
		}
		step2Result := GetJSON(step2ByCVData)
		result.GeminiStep4 = step2Result
		result.GeminiKey = key
		err = result.Step2Update(client)
		log.Printf("update id: %d gemini step2 result:%s \r\n", id, step2ByCVData)
	} else {
		if result.GeminiStep1 == "" {
			log.Println("not have GeminiStep1 gemini result")
			return true
		}
		if result.GeminiStep2 != "" {
			log.Printf("id:%d gemini step2 already done \r\n", id)
			return true
		}
		step2ByStep1Data := GeminiStep2Deduce(result.GeminiStep1, key) //使用任务一的结构跑任务二
		step2Result := GetJSON(step2ByStep1Data)
		result.GeminiStep2 = step2Result
		result.GeminiKey = key
		err = result.Step2Update(client)
		if err != nil {
			log.Println(err)
		}
		log.Printf("update id: %d gemini step2 result length:%d \r\n", id, len(result.GeminiStep2))
	}
	return err == nil
}

func GeminiStep2Deduce(step1Result, key string) string {
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
	log.Printf("call step2 request para length:%d\r\n", len(content))
	text := genai.Text(content)
	resp, err := model.GenerateContent(ctx, text)
	if resp == nil || err != nil {
		log.Println("step2 gemini response data is null", err)
		return "error"
	}
	errorMsg, _ := json.Marshal(resp)
	reason := resp.Candidates[0].FinishReason
	if reason == 0 || reason == 2 || reason == 3 || reason == 4 || reason == 5 {
		log.Println("step2 call gemini response:", string(errorMsg))
		return "error"
	}
	candidates := resp.Candidates
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from panic:", r)
		}
	}()
	if resp == nil || len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "error"
	}
	part := candidates[0].Content.Parts[0]
	return fmt.Sprintf("%s", part)
}

func step2ContentBuilder(step1Result string) string {
	temple, err := template.New("gemini-step2").Parse(tpl.STEP2)
	if err != nil {
		panic(err)
	}
	var out bytes.Buffer
	err = json.Compact(&out, []byte(step1Result))
	if err != nil {
		fmt.Println("Error compacting JSON:", err)
		return step1Result
	}
	data := Step1Result{Step1Result: out.String()}
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
