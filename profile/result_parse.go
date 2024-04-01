package profile

import (
	"encoding/json"
	"github.com/buger/jsonparser"
	"strconv"
)

// ParseStep2WorkExperienceArray 解析任务二工作经历Json结果
func ParseStep2WorkExperienceArray(jsonData []byte) []CompanyAdditionalInfo {
	var deduceWorkExpArr []CompanyAdditionalInfo

	return deduceWorkExpArr
}

// ParseStep2EduInfoArray 解析任务二学历Json结果
func ParseStep2EduInfoArray(jsonData []byte) []EducationAdditionalInfo {
	return nil
}

// ParseStep1JsonData 解析任务一Json结果
func ParseStep1JsonData(jsonData []byte) Resume {
	basicInfo := parseBasicInfo(jsonData)
	contactInfo := parseContactInformation(jsonData)
	introduction := parseSefIntroduction(jsonData)
	exp := parseWorkExp(jsonData)
	education := parseEducation(jsonData)
	languages := parseLanguage(jsonData)
	skills := parseSkills(jsonData)
	certs := parseCertifications(jsonData)
	resume := Resume{
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

func parseSkills(jsonData []byte) []Skill {
	var skills []Skill
	jsonparser.ArrayEach(jsonData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		language := Skill{}
		name, _ := jsonparser.GetString(value, "skill")
		proficiency, _ := jsonparser.GetString(value, "proficiency")
		language.Skill = name
		language.Proficiency = proficiency
		skills = append(skills, language)
	}, "skills")
	return skills
}

func parseLanguage(jsonData []byte) []Language {
	var languages []Language
	jsonparser.ArrayEach(jsonData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		language := Language{}
		name, _ := jsonparser.GetString(value, "language")
		proficiency, _ := jsonparser.GetString(value, "proficiency")
		language.Language = name
		language.Proficiency = proficiency
		languages = append(languages, language)
	}, "language")
	return languages
}

func parseEducation(jsonData []byte) []Education {
	var educations []Education
	index := 1
	jsonparser.ArrayEach(jsonData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		education := Education{}
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

func parseWorkExp(jsonData []byte) []WorkExperience {
	var experienceArr []WorkExperience
	index := 1
	jsonparser.ArrayEach(jsonData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		var positionArr []PositionInfo
		companyName, _ := jsonparser.GetString(value, "company_name"+strconv.Itoa(index))
		jsonparser.ArrayEach(value, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			var positionInfo = PositionInfo{}
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
		var workExp = WorkExperience{
			CompanyName:  companyName,
			PositionInfo: positionArr,
		}
		experienceArr = append(experienceArr, workExp)
		index += 1
	}, "work_experience")
	return experienceArr
}

func parseSefIntroduction(jsonData []byte) SefIntroduction {
	desc, _ := jsonparser.GetString(jsonData, "sef_introduction", "desc")
	introduction := SefIntroduction{
		Desc: desc,
	}
	return introduction
}

func parseContactInformation(jsonData []byte) ContactInformation {
	var phoneArr []string
	var emailArr []string
	var socialArr []SocialNetwork
	jsonparser.ArrayEach(jsonData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		phoneArr = append(phoneArr, string(value))
	}, "contact_information", "phone")
	jsonparser.ArrayEach(jsonData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		emailArr = append(emailArr, string(value))
	}, "contact_information", "email")
	jsonparser.ArrayEach(jsonData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		var data = SocialNetwork{}
		json.Unmarshal(value, &data)
		socialArr = append(socialArr, data)
	}, "contact_information", "social_network")

	information := ContactInformation{
		Phone:         phoneArr,
		Email:         emailArr,
		SocialNetwork: socialArr,
	}
	return information
}

func parseBasicInfo(jsonData []byte) BasicInformation {
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
	expectSalary := ExpectSalary{
		Min: min,
		Max: max,
	}
	b := BasicInformation{
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
