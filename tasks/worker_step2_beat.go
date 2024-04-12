package tasks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gemini/db"
	"gemini/store"
	"gemini/tpl"
	"gemini/utils"
	"log"
	"text/template"
)

func ReDoDeduceBeat(result store.GeminiResult, key string) bool {
	client := db.Client()
	if result.GeminiStep1 == "" {
		log.Printf("id:%d,url:%s not hava step1 result,start step1 task", result.ID, result.CVURL)
		step1 := GeminiStep1MergeBeat(result.ProfileData, result.CVData, key)
		if step1 == "error" || step1 == "" {
			log.Println("redo step1 and gemini response is null")
			return true
		}
		jsonResult := getJSONBeat(step1)
		result.GeminiStep1 = jsonResult
		result.Update(client)
		log.Printf("update step1 gemini result length:%d\r\n", len(jsonResult))
	}
	step2ByStep1Data := GeminiStep2DeduceBeat(result.GeminiStep1, key) //使用任务一的结构跑任务二
	step2Result := getJSONBeat(step2ByStep1Data)
	result.GeminiStep2 = step2Result
	result.GeminiKey = key
	result.Step2Update(client)
	log.Printf("update id: %d gemini step2 result length:%d \r\n", result.ID, len(result.GeminiStep2))
	return true
}

func GeminiStep2DeduceBeat(step1Result, key string) string {
	content := step2ContentBuilderBeta(step1Result)
	log.Printf("call step2 request para length:%d\r\n", len(content))
	resp, err := utils.InvokeGemini(content, key)
	if resp == nil || err != nil {
		log.Println("step2 gemini response data is null", err)
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

func step2ContentBuilderBeta(step1Result string) string {
	temple, err := template.New("gemini-step2").Parse(tpl.STEP2)
	if err != nil {
		panic(err)
	}
	var out bytes.Buffer
	err = json.Compact(&out, []byte(step1Result))
	var data Step1ResultBeat
	if err != nil {
		fmt.Println("Error compacting JSON:", err)
		data.Step1Result = step1Result
	} else {
		data.Step1Result = out.String()
	}
	var buf bytes.Buffer
	err = temple.Execute(&buf, data)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

type Step1ResultBeat struct {
	Step1Result string
}
