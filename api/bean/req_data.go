package bean

type SourceData struct {
	RequestId  string `json:"request_id"`  // 请求ID
	ResumeData string `json:"resume_data"` // CV原始数据
	ResumeId   string `json:"resume_id"`   // CV唯一标识
}
