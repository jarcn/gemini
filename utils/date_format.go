package utils

import (
	"regexp"
	"strings"
	"time"
)

var monthMap = map[string]string{
	"Januari":   "January",
	"Februari":  "February",
	"Maret":     "March",
	"April":     "April",
	"Mei":       "May",
	"Juni":      "June",
	"Juli":      "July",
	"Agustus":   "August",
	"September": "September",
	"Oktober":   "October",
	"November":  "November",
	"Desember":  "December",
}

var monthEnMap = map[string]string{
	"Jan": "01 January",
	"Feb": "02 February",
	"Mar": "03 March",
	"Apr": "04 April",
	"May": "05 May",
	"Jun": "06 June",
	"Jul": "07 July",
	"Aug": "08 August",
	"Sep": "09 September",
	"Oct": "10 October",
	"Nov": "11 November",
	"Dec": "12 December",
}

func DateFormat(date string) string {
	if strings.EqualFold(date, "Not provided") || strings.EqualFold(date, "Unknown") || date == "" ||
		strings.Contains(date, "Unknown") || strings.Contains(date, "unknown") {
		return "2006-01-02"
	}
	year, month := extractDateInfo(date)
	if year != "" && month != "" {
		date = year + "-" + month
	}
	if year != "" && month == "" {
		date = year + "-01"
	}
	if year == "" && month != "" {
		date = "2006-" + month
	}
	return parseAndFormatTime(date)
}

func parseAndFormatTime(input string) string {
	// 替换印度尼西亚月份名称为英文月份名称
	for indoMonth, engMonth := range monthMap {
		input = strings.ReplaceAll(input, indoMonth, engMonth)
	}
	// 解析时间字符串
	var t time.Time
	var err error
	formats := []string{"01-2006", "02/01/2006", "2006", "02.01.2006", "2006-01-02", "2 January 2006", "2 January 2006", "2006-Jan", "2-Jan-2006", "2006-1", "02-01-2006", "2 January 2006", "01/January/2006", "01 January,06", "January 02 2006", "02 January 2006", "January,2006", "02-January-2006"}
	for _, format := range formats {
		t, err = time.Parse(format, input)
		if err == nil {
			break
		}
	}
	if err != nil {
		return "2006-01-02"
	}
	// 格式化为标准时间格式
	return t.Format("2006-01-02")
}

func extractDateInfo(dateStr string) (string, string) {
	// 定义正则表达式来匹配日期中的年份和月份
	yearRegex := regexp.MustCompile(`(\d{4})`)
	monthRegex := regexp.MustCompile(`(Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec)`)

	// 查找年份和月份
	// 查找年份和月份
	yearMatches := yearRegex.FindStringSubmatch(dateStr)
	monthMatches := monthRegex.FindStringSubmatch(dateStr)

	var year, month string
	if len(yearMatches) > 1 {
		year = yearMatches[1]
	}
	if len(monthMatches) > 0 {
		month = monthEnMap[monthMatches[0]]
	}
	return year, month
}
