package workers

import (
	"context"
	"gemini/queue"
	"github.com/google/generative-ai-go/genai"
	"log"
)

var template = `您好，请扮演一个专业的猎头,排除对该求职者的不重要信息,分析<附件简历>信息,帮我将所有信息中的关键信息提取出来:
                                     如:
                                     name:""(提取<附件简历>中的姓名)
                                     work_description:[
                                     {
                                        "organization":""(提取<附件简历>中的这一段工作经历的公司)
                                        "work_time":""(提取<附件简历>中的工作开始时间-工作结束时间 格式参考yyyy/MM/dd-yyyy/MM/dd,默认"")
                                        "work_year":""(提取<附件简历>中的这一段工作时间有多少年)
                                        "work_description":""(提取<附件简历>中的工作描述)
                                        "work_title":""(根据工作描述重新生成的标题)
                                        "skill":[""](提取工作描述中的技能 skill指的是人员曾经学习过或者有经验的,可以量化的一些技能标签,例如excel/驾驶汽车/教育/管理等)
                                        "language_skills":[""](提取工作描述中的语言技能)
                                        "tool":[""](提取工作描述中的工具,例如:[Mobil, sepeda motor, komputer, telepon genggam])
                                        "occupation_class":""(提取工作描述中的职级 例如:assistant、General staff、director、manager、director、executive)
                                        "work_industry_classification":[""](将work_description归类到所属行业,例如：汽车,it等,最多三个,默认[""])
                                        "work_functional_classification":[""](将work_description归类到所属职能,例如:销售,程序开发等,最多三个,默认[""])
                                     },(按时间排序取最近的前五条)...
                                     ]
                                     self_introduction:[""](提取<附件简历>中的自我介绍)
                                     gender:""(提取<附件简历>中的性别 "1"-男 "2"-女 "0"-未知)
                                     birthday:""(提取<附件简历>中的出生日期 格式参照:yyyy/MM/dd 默认"")
                                     email:[""](提取<附件简历>中提取的邮箱 默认"")
                                     phone:""(提取<附件简历>中的电话号码,如有多个，只取可能性最大的一个 默认"")
                                     height:""(提取<附件简历>中的身高 默认"")
                                     weight:""(提取<附件简历>中的体重 默认"")
                                     is_online:""(是否在职 0-不在职 1-在职 默认"")
                                     position_industry:[""](提取<附件简历>中提取的期望行业 例如:保险/IT/汽车等)
                                     position_functional:[""](提取<附件简历>中的期望职能 例如:销售/开发等)
                                     professional:[""](教育专业 例如[mechanical,Computer major])
                                     education:[{
                                        "education_time":""(提取<附件简历>中的教育经历的开始时间-教育经历的结束时间 格式参考yyyy/MM/dd-yyyy/MM/dd,默认"")
                                        "school":""(提取<附件简历>中本段教育经历的学校)
                                        "education_professional":""(提取<附件简历>中本段教育经历的专业)
                                        "education_level":""(提取<附件简历>中本段教育经历的学历 例如:SMA,SMK)
                                     },...]
                                     skill:[""](提取<附件简历>中的技能 skill指的是人员曾经学习过或者有经验的,可以量化的一些技能标签,例如excel/驾驶汽车/教育/管理等)
                                     language_skills:[""]示例:[English,Japanese],默认[""]
                                     tool:[""]示例:[Mobil, sepeda motor, komputer, telepon genggam],默认[""]
                                     position_salary_and_benefits:[""](期望薪资福利)
                                     management_team_experience:"默认0(1-Experience in managing a team,0-No requirements)"
                                     management_team_experience_year:默认""
                                     management_team_size:"默认0(0-No requirements,1-Manage 1-10 people,2-Manage 11-20 people,3-Manage 21-50 people,4-Manage 51-100 people,5-Manage more than 100 people)"
                                     position_city:[""](期望城市)
                                     religion:默认[""]
                                     address:""(提取<附件简历>中的家庭住址)
                                     tripartite_url:[""](第三方url链接)
                                     请以规范的JSON格式输出分析结果,判断<附件简历>采用的什么语种,以该语种返回,格式参考如下:
                                     {
                                    name:""
                                     work_description:[""]
                                     self_introduction:[""]
                                     gender:""
                                     birthday:""
                                     email:[""]
                                     phone:""
                                     height:""
                                     weight:""
                                     is_online:""
                                     position_industry:[""]
                                     position_functional:[""]
                                     professional:[""]
                                     education:[""]
                                     skill:[""]
                                     language_skills:[""]
                                     tool:[""]
                                     position_salary_and_benefits:[""]
                                     management_team_experience:""
                                     management_team_experience_year:""
                                     management_team_size:""
                                     position_city:[""]
                                     religion:默认[""]
                                     address:默认""
                                     tripartite_url:[""]`

func CallGemini(data string, pool *queue.GeminiQueue) {
	c := pool.Dequeue()
	defer pool.Enqueue(c)
	ctx := context.Background()
	generateContent, err := c.GenerateContent(ctx, genai.Text("Write a story about a magic backpack."))
	if err != nil {
		log.Println(err)
	}
	log.Println(generateContent)
}
