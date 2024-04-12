package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Request 结构体用于解析请求参数
type Request struct {
	Contents []Content `json:"contents"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

// Response 结构体用于解析响应参数
type Response struct {
	Candidates []Candidate `json:"candidates"`
}

type Candidate struct {
	Content       Content        `json:"content"`
	FinishReason  string         `json:"finishReason"`
	Index         int            `json:"index"`
	SafetyRatings []SafetyRating `json:"safetyRatings"`
}

type SafetyRating struct {
	Category    string `json:"category"`
	Probability string `json:"probability"`
}

var geminiUrl = "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-pro-latest:generateContent?key="

func InvokeGemini(content, key string) (*Response, error) {
	requestBody := Request{Contents: []Content{{Parts: []Part{{Text: content}}}}}
	requestBodyBytes, _ := json.Marshal(requestBody)
	req, err := http.NewRequest("POST", geminiUrl+key, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	respStr, _ := json.Marshal(resp)
	log.Printf("gemini resp:%s", respStr)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error parsing response:", err)
		return nil, err
	}

	return &response, nil
}
