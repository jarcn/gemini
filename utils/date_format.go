package utils

func DateFormat(date string) string {
	if date == "Unknown" {
		return "2006-01-02"
	}
	return date
}
