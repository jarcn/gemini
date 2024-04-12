package main

import (
	"fmt"
	"gemini/db"
	"gemini/profile"
	"gemini/store"
)

func init() {
	db.MustInitMySQL("kp_user_local:Kupu123!@#@tcp(10.131.0.206:3306)/qiyee_job_data")
}

//var WorkExperienceTotal = 0
//var WorkExperience_IndustryAttribute = 0
//var WorkExperience_CompanyIntroduction = 0
//var WorkExperience_PositionCategory = 0
//var WorkExperience_PositionLevel = 0
//var WorkExperience_ManagementScope = 0
//
//var EduTotal = 0
//var Edu_DegreeAttribute = 0
//var Edu_MajorAttribute = 0

//var Edu_IsTopUniversity = 0
//var Edu_IsChineseSchool = 0

//职位名称:职类
//公司名称:行业
//专业
//学历

func main() {
	var data = store.GeminiResult{}
	result, _ := data.SelectAllStep2Result(db.Client())
	//fmt.Printf("%s,%s,%s,%s\r\n", "公司名称", "行业分类", "职位名称", "职位类别")
	fmt.Printf("%s,%s\r\n", "degree", "major")
	for _, info := range result {
		step2 := info.GeminiStep2
		_, eduArrays := profile.ParseStep2JsonData([]byte(step2))
		//printExp(expArrays)
		printEdu(eduArrays)
	}
	//fmt.Printf("WorkExperienceTotal:%d,EduTotal:%d\r\n", WorkExperienceTotal, EduTotal)
	//fmt.Printf("WorkExperience:IndustryAttribute:%d,CompanyIntroduction:%d,PositionCategory:%d,PositionLevel:%d,ManagementScope:%d\r\n",
	//	WorkExperience_IndustryAttribute, WorkExperience_CompanyIntroduction, WorkExperience_PositionCategory, WorkExperience_PositionLevel, WorkExperience_ManagementScope)
	//fmt.Printf("Edu:DegreeAttribute:%d,Edu:MajorAttribute:%d\r\n", Edu_DegreeAttribute, Edu_MajorAttribute)
}

func printExp(expArr []profile.WorkExperienceArray) {
	//WorkExperienceTotal += len(expArr)
	for _, data := range expArr {
		title := data.JobTitle
		name := data.CompanyName
		companyIndustry := data.CompanyAdditionalInfo.IndustryAttribute.Value //行业
		positionCategory := data.CompanyAdditionalInfo.PositionCategory.Value //职类
		fmt.Printf("%s,%s,%s,%s\r\n", name, companyIndustry, title, positionCategory)
		//industry := data.CompanyAdditionalInfo.IndustryAttribute.Value
		//if strings.EqualFold(industry, "Not provided") || strings.EqualFold(industry, "Unknown") || industry == "" ||
		//	strings.Contains(industry, "Unknown") || strings.Contains(industry, "unknown") {
		//	WorkExperience_IndustryAttribute += 1
		//}
		//companyIntroduction := data.CompanyAdditionalInfo.CompanyIntroduction.Value
		//if strings.EqualFold(companyIntroduction, "Not provided") || strings.EqualFold(companyIntroduction, "Unknown") || companyIntroduction == "" ||
		//	strings.Contains(companyIntroduction, "Unknown") || strings.Contains(companyIntroduction, "unknown") {
		//	WorkExperience_CompanyIntroduction += 1
		//}
		//positionCategory := data.CompanyAdditionalInfo.PositionCategory.Value
		//if strings.EqualFold(positionCategory, "Not provided") || strings.EqualFold(positionCategory, "Unknown") || positionCategory == "" ||
		//	strings.Contains(positionCategory, "Unknown") || strings.Contains(positionCategory, "unknown") {
		//	WorkExperience_PositionCategory += 1
		//}
		//managementScope := data.CompanyAdditionalInfo.ManagementScope.Value
		//if strings.EqualFold(managementScope, "Not provided") || strings.EqualFold(managementScope, "Unknown") || managementScope == "" ||
		//	strings.Contains(managementScope, "Unknown") || strings.Contains(managementScope, "unknown") {
		//	WorkExperience_ManagementScope += 1
		//}
		//positionLevel := data.CompanyAdditionalInfo.PositionLevel.Value
		//if strings.EqualFold(positionLevel, "Not provided") || strings.EqualFold(positionLevel, "Unknown") || positionLevel == "" ||
		//	strings.Contains(positionLevel, "Unknown") || strings.Contains(positionLevel, "unknown") {
		//	WorkExperience_PositionLevel += 1
		//}
	}
}

func printEdu(eduArr []profile.EduInfoArray) {
	//EduTotal += len(eduArr)
	for _, data := range eduArr {
		degree := data.EducationAdditionalInfo.DegreeAttribute.Value
		major := data.EducationAdditionalInfo.MajorAttribute.Value
		fmt.Printf("%s,%s\r\n", degree, major)
		//degreeAttribute := data.EducationAdditionalInfo.DegreeAttribute.Value
		//if strings.EqualFold(degreeAttribute, "Not provided") || strings.EqualFold(degreeAttribute, "Unknown") || degreeAttribute == "" ||
		//	strings.Contains(degreeAttribute, "Unknown") || strings.Contains(degreeAttribute, "unknown") {
		//	Edu_DegreeAttribute += 1
		//}
		//majorAttribute := data.EducationAdditionalInfo.MajorAttribute.Value
		//if strings.EqualFold(majorAttribute, "Not provided") || strings.EqualFold(majorAttribute, "Unknown") || majorAttribute == "" ||
		//	strings.Contains(majorAttribute, "Unknown") || strings.Contains(majorAttribute, "unknown") {
		//	Edu_MajorAttribute += 1
		//}
	}
}
