package main

import (
	"bytes"
	"context"
	"encoding/json"
	"gemini/db"
	"gemini/profile"
	"gemini/store"
	"github.com/buger/jsonparser"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"log"
	"strconv"
)

func init() {
	db.MustInitMySQL("kp_user_local:Kupu123!@#@tcp(10.131.0.206:3306)/qiyee_job_data") //预发环境
	//db.MustInitMySQL("sc_kupu:Sc_kupu_1234@tcp(10.128.0.28:3306)/qiyee_job_data") //生产环境
}

func main() {
	result := store.GeminiResult{}
	all, _ := result.FindAll(db.Client())
	var resumeArr []profile.Resume
	for _, d := range all {
		step1 := d.GeminiStep1
		data := []byte(step1)
		resume := parseJsonData(data)
		resume.BasicInformation.ProfileUrl = d.CVURL
		resume.ID = d.ID
		resumeArr = append(resumeArr, resume)
	}
	insert2ES(resumeArr)
}

func insert2ES(resumeArr []profile.Resume) {
	// Elasticsearch 连接配置
	cfg := elasticsearch.Config{
		//Addresses: []string{"http://10.128.0.165:9200", "http://10.128.0.72:9200"}, //生产环境
		Addresses: []string{"http://10.129.0.251:9200", "http://10.129.0.217:9200", "http://10.129.0.146:9200"}, //预发环境
		Password:  "Qiyi123!@#",
		Username:  "elastic",
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// 执行批量插入请求
	// 批量插入数据
	for _, doc := range resumeArr {
		// 将文档转换为 JSON 格式
		docJSON, err := json.Marshal(doc)
		if err != nil {
			log.Printf("Error marshalling document: %s", err)
			continue
		}

		// 准备批量插入请求
		req := esapi.IndexRequest{
			//Index:      "gemini_hra", // 替换为你的索引名称
			Index:      "hra_cv", // 替换为你的索引名称
			DocumentID: strconv.Itoa(int(doc.ID)),
			Body:       bytes.NewReader(docJSON),
			Refresh:    "true",
		}

		// 执行批量插入请求
		res, err := req.Do(context.Background(), es)
		if err != nil {
			log.Printf("Error indexing document: %s", err)
			continue
		}
		defer res.Body.Close()

		if res.IsError() {
			log.Printf("[%s] Error indexing document: %s", res.Status(), res.String())
		} else {
			log.Printf("[%s] Indexed document with ID: %d", res.Status(), doc.ID)
		}
	}
}

func parseJsonData(jsonData []byte) profile.Resume {
	basicInfo := parseBasicInfo(jsonData)
	contactInfo := parseContactInformation(jsonData)
	introduction := parseSefIntroduction(jsonData)
	exp := parseWorkExp(jsonData)
	education := parseEducation(jsonData)
	languages := parseLanguage(jsonData)
	skills := parseSkills(jsonData)
	certs := parseCertifications(jsonData)
	resume := profile.Resume{
		BasicInformation:   basicInfo,
		ContactInformation: contactInfo,
		SefIntroduction:    introduction,
		WorkExperience:     exp,
		Education:          education,
		Language:           languages,
		Skills:             skills,
		Certifications:     certs,
	}
	return resume
}

func parseCertifications(jsonData []byte) []string {
	var certs []string
	jsonparser.ArrayEach(jsonData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		name, _ := jsonparser.GetString(value, "certifications")
		certs = append(certs, name)
	}, "certifications")
	return certs
}

func parseSkills(jsonData []byte) []profile.Skill {
	var skills []profile.Skill
	jsonparser.ArrayEach(jsonData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		language := profile.Skill{}
		name, _ := jsonparser.GetString(value, "skill")
		proficiency, _ := jsonparser.GetString(value, "proficiency")
		language.Skill = name
		language.Proficiency = proficiency
		skills = append(skills, language)
	}, "skills")
	return skills
}

func parseLanguage(jsonData []byte) []profile.Language {
	var languages []profile.Language
	jsonparser.ArrayEach(jsonData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		language := profile.Language{}
		name, _ := jsonparser.GetString(value, "language")
		proficiency, _ := jsonparser.GetString(value, "proficiency")
		language.Language = name
		language.Proficiency = proficiency
		languages = append(languages, language)
	}, "language")
	return languages
}

func parseEducation(jsonData []byte) []profile.Education {
	var educations []profile.Education
	index := 1
	jsonparser.ArrayEach(jsonData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		education := profile.Education{}
		school, _ := jsonparser.GetString(value, "school"+strconv.Itoa(index))
		degree, _ := jsonparser.GetString(value, "degree")
		major, _ := jsonparser.GetString(value, "major")
		duration, _ := jsonparser.GetString(value, "duration")
		education.School = school
		education.Degree = degree
		education.Major = major
		education.Duration = duration
		educations = append(educations, education)
		index += 1
	}, "education")
	return educations
}

func parseWorkExp(jsonData []byte) []profile.WorkExperience {
	var experienceArr []profile.WorkExperience
	index := 1
	jsonparser.ArrayEach(jsonData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		var positionArr []profile.PositionInfo
		companyName, _ := jsonparser.GetString(value, "company_name"+strconv.Itoa(index))
		jsonparser.ArrayEach(value, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			var positionInfo = profile.PositionInfo{}
			var responsibilities []string
			jobTitle, _ := jsonparser.GetString(value, "job_title")
			duration, _ := jsonparser.GetString(value, "duration")
			j := 1
			jsonparser.ArrayEach(value, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				desc, _ := jsonparser.GetString(value, strconv.Itoa(j)+".")
				if desc == "" {
					desc = string(value)
				}
				responsibilities = append(responsibilities, desc)
				j += 1
			}, "responsibilities")
			positionInfo.JobTitle = jobTitle
			positionInfo.Duration = duration
			positionInfo.Responsibilities = responsibilities
			positionArr = append(positionArr, positionInfo)
		}, "position_info")
		var workExp = profile.WorkExperience{
			CompanyName:  companyName,
			PositionInfo: positionArr,
		}
		experienceArr = append(experienceArr, workExp)
		index += 1
	}, "work_experience")
	return experienceArr
}

func parseSefIntroduction(jsonData []byte) profile.SefIntroduction {
	desc, _ := jsonparser.GetString(jsonData, "sef_introduction", "desc")
	introduction := profile.SefIntroduction{
		Desc: desc,
	}
	return introduction
}

func parseContactInformation(jsonData []byte) profile.ContactInformation {
	var phoneArr []string
	var emailArr []string
	var socialArr []profile.SocialNetwork
	jsonparser.ArrayEach(jsonData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		phoneArr = append(phoneArr, string(value))
	}, "contact_information", "phone")
	jsonparser.ArrayEach(jsonData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		emailArr = append(emailArr, string(value))
	}, "contact_information", "email")
	jsonparser.ArrayEach(jsonData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		var data = profile.SocialNetwork{}
		json.Unmarshal(value, &data)
		socialArr = append(socialArr, data)
	}, "contact_information", "social_network")

	information := profile.ContactInformation{
		Phone:         phoneArr,
		Email:         emailArr,
		SocialNetwork: socialArr,
	}
	return information
}

func parseBasicInfo(jsonData []byte) profile.BasicInformation {
	name, _ := jsonparser.GetString(jsonData, "basic_information", "name")
	gender, _ := jsonparser.GetString(jsonData, "basic_information", "gender")
	age, _ := jsonparser.GetInt(jsonData, "basic_information", "age")
	birthday, _ := jsonparser.GetString(jsonData, "basic_information", "birthday")
	maritalStatus, _ := jsonparser.GetString(jsonData, "basic_information", "marital_status")
	religion, _ := jsonparser.GetString(jsonData, "basic_information", "religion")
	residentialCity, _ := jsonparser.GetString(jsonData, "basic_information", "residential_city")
	jobSearchStatus, _ := jsonparser.GetString(jsonData, "basic_information", "job_search_status")
	perferWorkCity, _ := jsonparser.GetString(jsonData, "basic_information", "perfer_work_city")
	perferPosition, _ := jsonparser.GetString(jsonData, "basic_information", "perfer_position")
	min, _ := jsonparser.GetInt(jsonData, "basic_information", "expect_salary", "min")
	max, _ := jsonparser.GetInt(jsonData, "basic_information", "expect_salary", "max")
	expectSalary := profile.ExpectSalary{
		Min: min,
		Max: max,
	}
	b := profile.BasicInformation{
		Name:            name,
		Gender:          gender,
		Birthday:        birthday,
		Age:             age,
		MaritalStatus:   maritalStatus,
		Religion:        religion,
		ResidentialCity: residentialCity,
		JobSearchStatus: jobSearchStatus,
		PerferPosition:  perferPosition,
		PerferWorkCity:  perferWorkCity,
		ExpectSalary:    expectSalary,
	}
	return b
}
