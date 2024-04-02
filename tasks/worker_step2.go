package tasks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gemini/db"
	"gemini/profile"
	"gemini/store"
	"gemini/tpl"
	deepcopier "gemini/utils"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"log"
	"strings"
	"text/template"
)

func DoDeduce(msg []byte, key string, isCvData bool) bool {
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
		return false
	}
	if isCvData {
		if result.CVData == "" {
			fmt.Printf("id: %d not have cv data \r\n", id)
			return true
		}
		if result.GeminiStep2 == "" {
			fmt.Printf("id:%d not have step2 result deduce\r\n", id)
			return true
		}
		step2ByCVData := geminiStep2Deduce(result.CVData, key) //使用原始CV数据跑任务二
		if step2ByCVData == "" {
			return false
		}
		step2Result := step2ResultToJson(step2ByCVData)
		result.GeminiStep4 = step2Result
		result.GeminiKey = key
		err = result.Step4Update(db.Client())
		fmt.Printf("update id: %d gemini step2 result:%s \r\n", id, step2ByCVData)
	} else {
		if result.GeminiStep1 == "" {
			fmt.Printf("id: %d not have step1 result\r\n", id)
			return true
		}
		//if result.GeminiStep2 != "" {
		//	fmt.Printf("id:%d step2 already deduce\r\n", id)
		//	return true
		//}
		step2ByStep1Data := geminiStep2Deduce(result.GeminiStep1, key) //使用任务一的结构跑任务二
		step2Result := step2ResultToJson(step2ByStep1Data)
		result.GeminiStep2 = step2Result
		result.GeminiKey = key
		err = result.Step2Update(db.Client())
		fmt.Printf("update id: %d gemini step2 result:%s \r\n", id, step2ByStep1Data)
	}
	step1AndStep2MergeJson := mergeStep1AndStep2([]byte(result.GeminiStep1), []byte(result.GeminiStep2))
	jsonData, _ := json.Marshal(step1AndStep2MergeJson)
	fmt.Println(jsonData)
	//TODO 保存合并后的数据
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
	fmt.Printf("call step2 request para:%s\r\n", content)
	resp, err := model.GenerateContent(ctx, genai.Text(content))
	if err != nil {
		fmt.Println("call gemini error:", err)
		return ""
	}
	candidates := resp.Candidates
	defer func() string {
		errorMsg, _ := json.Marshal(resp)
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
			fmt.Println("step2 call gemini response:", string(errorMsg))
		}
		return string(errorMsg)
	}()
	if resp == nil || len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return ""
	}
	part := candidates[0].Content.Parts[0]
	fmt.Printf("call step2 response data:%s\r\n", part)
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

func mergeStep1AndStep2(step1Json, step2Json []byte) *profile.Resume {
	resumeData := profile.ParseStep1JsonData(step1Json)
	workExpData, eduInfoData := profile.ParseStep2JsonData(step2Json)
	for i := 0; i < len(resumeData.WorkExperience); i++ {
		experience := resumeData.WorkExperience[i]
		positionInfoArr := resumeData.WorkExperience[i].PositionInfo
		for j := 0; j < len(positionInfoArr); j++ {
			info := &positionInfoArr[j]
			title := info.JobTitle
			for _, work := range workExpData {
				if experience.CompanyName == work.CompanyName || title == work.JobTitle {
					deepcopier.Copy(work.CompanyAdditionalInfo).To(&info.CompanyAdditionalInfo)
				}
			}
		}
	}
	for i := 0; i < len(resumeData.Educations); i++ {
		education := &resumeData.Educations[i]
		for _, info := range eduInfoData {
			if education.School == info.School || education.Degree == info.Degree {
				deepcopier.Copy(info.EducationAdditionalInfo).To(&education.EducationAdditionalInfo)
			}
		}
	}
	return resumeData
}

type Step1Result struct {
	Step1Result string
}
