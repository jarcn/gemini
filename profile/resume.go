package profile

type ExpectSalary struct {
	Min int64 `json:"min"`
	Max int64 `json:"max"`
}

type BasicInformation struct {
	Name            string       `json:"name"`
	Gender          string       `json:"gender"`
	Age             int64        `json:"age"`
	Birthday        string       `json:"birthday"`
	MaritalStatus   string       `json:"marital_status"`
	Religion        string       `json:"religion"`
	ResidentialCity string       `json:"residential_city"`
	JobSearchStatus string       `json:"job_search_status"`
	PerferWorkCity  string       `json:"perfer_work_city"`
	PerferPosition  string       `json:"perfer_position"`
	ProfileUrl      string       `json:"profile_url"`
	ExpectSalary    ExpectSalary `json:"expect_salary"`
}

type ContactInformation struct {
	Phone         []string        `json:"phone"`
	Email         []string        `json:"email"`
	SocialNetwork []SocialNetwork `json:"social_network"`
}

type SocialNetwork struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type PositionInfo struct {
	JobTitle              string                `json:"job_title"`
	Duration              string                `json:"duration"`
	Responsibilities      []string              `json:"responsibilities"`
	CompanyAdditionalInfo CompanyAdditionalInfo `json:"company_additional_info"`
}

type CompanyAdditionalInfo struct {
	IndustryAttribute    string     `json:"industry_attribute"`
	CompanyIntroduction  string     `json:"company_introduction"`
	PositionCategory     string     `json:"position_category"`
	PositionLevel        string     `json:"position_level"`
	IsManagementPosition string     `json:"is_management_position"`
	ManagementScope      string     `json:"management_scope"`
	WorkPeriod           WorkPeriod `json:"work_period"`
	Startup              Startup    `json:"startup"`
	KeySkillsExperience  []string   `json:"key_skills_experience"`
}

type Startup struct {
	IsStartup   string `json:"is_startup"`
	CreatedTime string `json:"created_time"`
}

type WorkPeriod struct {
	StartDate        string `json:"start_date"`
	EndDate          string `json:"end_date"`
	DurationInMonths int    `json:"duration_in_months"`
}

type WorkExperience struct {
	CompanyName  string         `json:"company_name"`
	PositionInfo []PositionInfo `json:"position_info"`
}

type Education struct {
	School                  string                  `json:"school"`
	Degree                  string                  `json:"degree"`
	Major                   string                  `json:"major"`
	Duration                string                  `json:"duration"`
	EducationAdditionalInfo EducationAdditionalInfo `json:"edu_info_array"`
}

type EducationAdditionalInfo struct {
	DegreeAttribute string `json:"degree_attribute"`
	MajorAttribute  string `json:"major_attribute"`
	IsTopUniversity string `json:"is_top_university"`
	IsChineseSchool string `json:"is_chinese_school"`
}

type Language struct {
	Language    string `json:"language"`
	Proficiency string `json:"proficiency"`
}

type Skill struct {
	Skill       string `json:"skill"`
	Proficiency string `json:"proficiency"`
}

type SefIntroduction struct {
	Desc string `json:"desc"`
}

type Resume struct {
	ID                 int64              `json:"id"`
	BasicInformation   BasicInformation   `json:"basic_information"`
	ContactInformation ContactInformation `json:"contact_information"`
	SefIntroduction    SefIntroduction    `json:"sef_introduction"`
	WorkExperience     []WorkExperience   `json:"work_experience"`
	Education          []Education        `json:"education"`
	Language           []Language         `json:"language"`
	Skills             []Skill            `json:"skills"`
	Certifications     []string           `json:"certifications"`
}
