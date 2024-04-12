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
	"strings"
	"text/template"
)

func DoMergeBeat(msg []byte, key string) bool {
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

func ReDoMergeBeat(data store.GeminiResult, key string) *store.GeminiResult {
	switch data.Type {
	case "ZH":
		return reDoZh(data, key)
	case "HRA":
		return reDoHra(data, key)
	}
	return nil
}

func reDoZh(data store.GeminiResult, key string) *store.GeminiResult {
	url := data.CVURL
	if strings.TrimSpace(data.CVData) == "" {
		cvData := store.CvData{}
		cv, err := cvData.GetCvByUrl(db.Client(), url)
		if err != nil || cv == nil {
			log.Printf("url:%s not have data \r", url)
			return nil
		}
		data.CVData = cv.ResumeMsg
	}
	step1 := GeminiStep1MergeBeat(data.ProfileData, data.CVData, key)
	jsonResult := getJSONBeat(step1)
	data.GeminiStep1 = jsonResult
	err := data.Update(db.Client())
	if err != nil {
		log.Println("", err)
	}
	log.Printf("update step1 gemini result length:%d\r\n", len(jsonResult))
	return &data
}
func reDoHra(data store.GeminiResult, key string) *store.GeminiResult {
	step1 := GeminiStep1MergeBeat(data.ProfileData, data.CVData, key)
	jsonResult := getJSONBeat(step1)
	data.GeminiStep1 = jsonResult
	data.Update(db.Client())
	log.Printf("update step1 gemini result length:%d\r\n", len(jsonResult))
	return &data
}

func getJSONBeat(s string) string {
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start == -1 || end == -1 || start >= end {
		return ""
	}
	return s[start : end+1]
}

func GeminiStep1MergeBeat(profileCv, ocrCv, key string) string {
	content := parseContentBeat(ocrCv, profileCv)
	log.Printf("call step1 para length:%d\r\n", len(content))
	resp, err := utils.InvokeGemini(content, key)
	if resp == nil || err != nil {
		log.Println("step1 gemini response data is null", err)
		return "error"
	}
	candidates := resp.Candidates
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from panic:", r)
		}
	}()
	if resp == nil || len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		log.Println("cell step1 not response data")
		return ""
	}
	part := candidates[0].Content.Parts[0]
	return fmt.Sprintf("%s", part)
}

func parseContentBeat(ocrCv, profileCv string) string {
	temple, err := template.New("gemini-step1").Parse(tpl.STEP1)
	if err != nil {
		panic(err)
	}
	var out bytes.Buffer
	err = json.Compact(&out, []byte(profileCv))
	data := Data{}
	if err != nil {
		data.ProfileCV = profileCv
		data.OcrCV = ocrCv
	} else {
		data.ProfileCV = out.String()
		data.OcrCV = ocrCv
	}
	var buf bytes.Buffer
	err = temple.Execute(&buf, data)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

type DataBeat struct {
	OcrCV     string
	ProfileCV string
}
