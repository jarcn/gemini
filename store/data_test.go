package store

import (
	"encoding/json"
	"fmt"
	"gemini/db"
	"gemini/profile"
	"testing"
)

func TestQuery(t *testing.T) {
	db.MustInitMySQL(db.MysqlAddr)
	cvData := CvData{}
	cv, err := cvData.GetCvByUrl(db.Client(), "http://10.128.0.250/18ebc71cac2e84a3c9fd38bd3f680d67.pdf")
	if err != nil {
		fmt.Println(err)
	}
	msg := cv.ResumeMsg
	fmt.Println(msg)
}

func TestQueryById(t *testing.T) {
	db.MustInitMySQL(db.MysqlAddr)
	cvData := GeminiResult{}
	cv, err := cvData.QueryById(db.Client(), 1)
	if err != nil {
		fmt.Println(err)
	}
	msg := cv.Type
	fmt.Println(msg)
}

func TestName(t *testing.T) {
	var js = `{"basic_information":{"name":"Shane Soetaniman","gender":"Male","age":28,"birthday":"1993-12-02","marital_status":"Unknown","religion":"Unknown","residential_city":"Jakarta","job_search_status":"I_AM_LOOKING_FOR_JOB","perfer_work_city":[],"perfer_position":["Product Manager"],"expect_salary":{"min":null,"max":18000000}},"contact_information":{"phone":["817788239"],"email":["soetanimans@gmail.com"],"social_network":[{"name":"LinkedIn","url":"https://www.linkedin.com/in/shanesoetaniman/"}]},"sef_introduction":{"desc":"An aspiring Product Manager with a background in startup strategy and financial modeling fueled by an Economics and Business Management degree. Passionate in learning about Product Management, UX Research, Product Design, and Growth Hacking. Enjoy working with full stack teams who are banded around a common vision and love solving problems."},"work_experience":[{"company_name":"Grasshopper Pte Ltd & ChatQ","position_info":[{"job_title":"Investor Relations Manager","duration":"2019-07 to 2020-01","responsibilities":[{"1.":"Assisted with company's efforts to find external investors for trading capital and retail investor app ChatQ – conducting investor profile analysis and generating strategies of approach","2.":"Managed all variations of IR products for Grasshopper and ChatQ – Investment Deck, DCF Model, Investment Memo (company model and overview, market analysis, competitive analysis, expansion plans, financials)"}],"job_title":"Regional Sales Coordinator","duration":"2015-01-02 to 2016-12-02","responsibilities":[{"1.":"Organized teams to cater all clientele of two independent national markets for a Harvard-based educational tutoring venture"}]}]},{"company_name":"Travelio.com","position_info":[{"job_title":"Product Manager","duration":"2018-06-02 to 2019-10-02","responsibilities":[{"1.":"Acted as sole PM and scrum master for all projects reporting directly to management","2.":"Restructured developer team’s departmental SOP – implementing agile methodology cycle to a scrum framework and PRDs, streamlined Sprint meetings, built internal request and ticketing system","3.":"Lead and designed all UI/UX research, projects and protocols in optimizing lower funnel conversion & user retention - compiled both hard and soft data, conducted A/B Testing, Persona Testing and Guerilla Testing","4.":"Set engineers’ day-to-day tasks and weekly objectives - implement a soft KPI system and automate weekly reports for management","5.":"Managed development of operation’s handover system and feedback loop master data in orchestrating cross-functional teams","6.":"Assisted in closing company’s Series B funding round, building and updating financial projections and maintaining investor relations due diligence","7.":"Streamlined Sales/Business Development team’s B2C and B2B feasibility models, automating it based on client and building criteria using Excel VBA (sidelining the need for CS involvement)","8.":"Conducted department audits to find most efficient work flow protocols and practices"}]}]},{"company_name":"LaunchByte.io","position_info":[{"job_title":"Project Analyst","duration":"2017-01-02 to 2018-01-02","responsibilities":[{"1.":"Consulted directly with clients in designing a valid pricing structure and revenue model for products","2.":"Corroborated models with existing internal client data and product, market and competitive analysis","3.":"Created all financial models, generating EBITDA and valuations for clients with consumer-face mobile and web applications, consolidating units for investor memo and deck and reporting directly to the respective Project Manager","4.":"Assisted with client financial and intelligence due diligence with prospective investors"}]}]},{"company_name":"Search Fund Accelerator","position_info":[{"job_title":"Summer Analyst","duration":"2016-05-02 to 2016-08-02","responsibilities":[{"1.":"Conducted analytical research in finding potential private companies in recommending strategies for high ROI","2.":"Assisted with product and industry analysis for newly appointed or potential entrepreneurs in a search fund"}]}]},{"company_name":"Grasshopper Asia Pte Ltd","position_info":[{"job_title":"IT Intern","duration":"2015-05-02 to 2016-02-02","responsibilities":[{"1.":"Shadowed and mentored by software engineers and quant developers in creating methodical features and tools for inhouse traders in the derivative market","2.":"Created company culture book ‘The Green Book’ for new recruits","3.":"Architect for company’s Confluence internal database deploying strategy and exchange information, code guides and data (framework adopted by SGX)"},{"job_title":"Co-founder & Project Director","duration":"2018-01-02 to 2019-05-02","responsibilities":[{"1.":"Created eSports marketing analytics platform for gamers, influencers, teams, and businesses (no longer involved with operations)","2.":"Conducted UX research via Guerilla testing and designed web platform’s UI managing team of graphic designers","3.":"Managed design team to deliver content and branding for young upcoming gaming influencers"}]}]},{"company_name":"Park Hyatt, Hyatt Hotel Corp.","position_info":[{"job_title":"Business Administration Intern","duration":"2015-01-02 to 2015-05-02","responsibilities":[]}]}],"education":[{"school":"Emmanuel College","degree":"BA","major":"Economics & Business Management","duration":"2013-09 to 2017-05"}],"language":[{"language":"Indonesian","proficiency":"Fluent"},{"language":"Malay","proficiency":"Basic"},{"language":"Mandarin Chinese","proficiency":"Basic"}],"skills":[{"skill":"A/B Testing","proficiency":"Intermediate"},{"skill":"English","proficiency":"Fluent"},{"skill":"HTML","proficiency":"Basic"},{"skill":"It Product Management","proficiency":"Basic"},{"skill":"JavaScript","proficiency":"Basic"},{"skill":"Leadership","proficiency":"Intermediate"},{"skill":"Management","proficiency":"Intermediate"},{"skill":"Microsoft Excel","proficiency":"Intermediate"},{"skill":"Microsoft Office","proficiency":"Intermediate"},{"skill":"Microsoft Word","proficiency":"Intermediate"},{"skill":"PostgreSQL","proficiency":"Intermediate"},{"skill":"Product Development","proficiency":"Intermediate"},{"skill":"Project Management","proficiency":"Intermediate"},{"skill":"Python","proficiency":"Basic"},{"skill":"Waterfall Methodologies","proficiency":"Intermediate"},{"skill":"Agile Project Management","proficiency":"Intermediate"},{"skill":"Analytical Skills","proficiency":"Intermediate"},{"skill":"CSS","proficiency":"Basic"},{"skill":"Adobe Photoshop","proficiency":"Intermediate"},{"skill":"Good Communication Skills","proficiency":"Intermediate"}],"certifications":[]}`
	data := profile.ParseStep1JsonData([]byte(js))
	marshal, _ := json.Marshal(data)
	fmt.Println(string(marshal))
}
