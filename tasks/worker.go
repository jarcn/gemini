package tasks

import (
	"bytes"
	"context"
	"fmt"
	"gemini/tlp"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"log"
	"text/template"
)

func CallGemini(ocrCv, profileCv, key string) string {
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
	}
	candidates := resp.Candidates
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()
	if resp == nil || len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return ""
	}
	part := candidates[0].Content.Parts[0]
	return fmt.Sprintf("%s", part)
}

func parseContent(ocrCv, profileCv string) string {
	tpl, err := template.New("gemini").Parse(tlp.STEP1)
	if err != nil {
		panic(err)
	}
	data := Data{OcrCV: ocrCv, ProfileCV: profileCv}
	var buf bytes.Buffer
	err = tpl.Execute(&buf, data)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

type Data struct {
	OcrCV     string
	ProfileCV string
}
