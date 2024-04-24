package profile

// Talent struct representing the entire JSON object
type Talent struct {
	Id                   int            `json:"id"`
	Address              string         `json:"address"`
	AllowJobRec          bool           `json:"allow_job_rec"`
	Availability         Availability   `json:"availability"`
	BirthDate            string         `json:"birth_date"`
	Educations           []EktEducation `json:"educations"`
	Email                string         `json:"email"`
	ExperienceDuration   string         `json:"experience_duration"`
	Experiences          []Experience   `json:"experiences"`
	Followers            interface{}    `json:"followers"`
	FullName             string         `json:"fullname"`
	Gender               string         `json:"gender"`
	ImageURL             ImageURL       `json:"image_url"`
	IsFreelancer         bool           `json:"is_freelancer"`
	Link                 Link           `json:"link"`
	Onboarded            bool           `json:"onboarded"`
	OpenProject          bool           `json:"open_project"`
	Preference           Preference     `json:"preference"`
	TalentLanguages      interface{}    `json:"talent_languages"`
	Phone                string         `json:"phone"`
	TalentSkills         []TalentSkill  `json:"talent_skills"`
	TalentCertifications interface{}    `json:"talent_certifications"`
	Type                 string         `json:"type"`
}

// Availability struct
type Availability struct {
	Name string `json:"name"`
}

// EktEducation struct
type EktEducation struct {
	Current    interface{} `json:"current"`
	Degree     Degree      `json:"degree"`
	EndDate    string      `json:"end_date"`
	Major      Major       `json:"major"`
	StartDate  string      `json:"start_date"`
	University University  `json:"university"`
}

// Degree struct
type Degree struct {
	Name string `json:"name"`
}

// Major struct
type Major struct {
	Name string `json:"name"`
}

// University struct
type University struct {
	Name string `json:"name"`
}

// Experience struct
type Experience struct {
	Company              Company          `json:"company"`
	Description          string           `json:"description"`
	EndDate              string           `json:"end_date"`
	StartDate            string           `json:"start_date"`
	ExperienceBusinesses interface{}      `json:"experience_businesses"`
	ExperienceSkillSets  interface{}      `json:"experience_skill_sets"`
	Location             string           `json:"location"`
	Position             Position         `json:"position"`
	PositionFunction     PositionFunction `json:"position_function"`
}

// Company struct
type Company struct {
	LinkedinURL interface{} `json:"linkedin_url"`
	Name        string      `json:"name"`
	Pt          string      `json:"pt"`
}

// Position struct
type Position struct {
	Name                string              `json:"name"`
	PositionLevel       EktPositionLevel    `json:"position_level"`
	PositionSubCategory PositionSubCategory `json:"position_sub_category"`
}

// EktPositionLevel struct
type EktPositionLevel struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}

// PositionSubCategory struct
type PositionSubCategory struct {
	Name             string              `json:"name"`
	PositionCategory EktPositionCategory `json:"position_category"`
}

// EktPositionCategory struct
type EktPositionCategory struct {
	Name string `json:"name"`
}

// PositionFunction struct
type PositionFunction struct {
	Name string `json:"name"`
}

// ImageURL struct
type ImageURL struct {
	Medium   string `json:"medium"`
	Original string `json:"original"`
	Thumb    string `json:"thumb"`
}

// Link struct
type Link struct {
	Behance       interface{} `json:"behance"`
	Dribbble      interface{} `json:"dribbble"`
	Facebook      string      `json:"facebook"`
	Github        string      `json:"github"`
	Google        string      `json:"google"`
	Kaggle        interface{} `json:"kaggle"`
	Linkedin      string      `json:"linkedin"`
	Stackoverflow string      `json:"stackoverflow"`
	Web           string      `json:"web"`
}

// Preference struct
type Preference struct {
	CurMinSal int `json:"cur_min_sal"`
	ExpMinSal int `json:"exp_min_sal"`
}

// TalentSkill struct
type TalentSkill struct {
	Duration interface{} `json:"duration"`
	SkillSet SkillSet    `json:"skill_set"`
}

// SkillSet struct
type SkillSet struct {
	Name             string              `json:"name"`
	PositionCategory EktPositionCategory `json:"position_category"`
	SkillCategory    struct {
		Name interface{} `json:"name"`
	} `json:"skill_category"`
}
