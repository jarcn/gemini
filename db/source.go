package db

// KafkaBrokers 开发环境
var KafkaBrokers = []string{"10.129.0.78:9092", "10.129.0.180:9092", "10.129.0.85:9092"}
var EsBrokers = []string{"http://10.129.0.251:9200", "http://10.129.0.217:9200", "http://10.129.0.146:9200"}
var MysqlAddr = "kp_user_local:Kupu123!@#@tcp(10.131.0.206:3306)/qiyee_job_data"

// KafkaBrokers 生产环境
//var KafkaBrokers = []string{""10.128.0.94:9092", "10.128.0.156:9092", "10.128.0.124:9092"}
//var EsBrokers = []string{"http://10.128.0.224:9200", "http://10.128.0.237:9200"}
//var MysqlAddr = "sc_kupu:Sc_kupu_1234@tcp(10.128.0.28:3306)/qiyee_job_data"
