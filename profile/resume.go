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
	JobTitle         string   `json:"job_title"`
	Duration         string   `json:"duration"`
	Responsibilities []string `json:"responsibilities"`
}

type WorkExperience struct {
	CompanyName  string         `json:"company_name"`
	PositionInfo []PositionInfo `json:"position_info"`
}

type Education struct {
	School   string `json:"school"`
	Degree   string `json:"degree"`
	Major    string `json:"major"`
	Duration string `json:"duration"`
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
