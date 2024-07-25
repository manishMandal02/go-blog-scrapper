package utils

var MONTHS = map[string]string{
	"January":   "01",
	"February":  "02",
	"March":     "03",
	"April":     "04",
	"May":       "05",
	"June":      "06",
	"July":      "07",
	"August":    "08",
	"September": "09",
	"October":   "10",
	"November":  "11",
	"December":  "12",
}

// format date to UTC string YYYY-MM-DD
func FormatDateStringUTC(year, month, day string) string {

	return string(year) + "-" + MONTHS[month] + "-" + formatToTwoDigit(day)

}

func SafeMaxLimit(num int, limit int) int {

	if num > limit || num < 1 {
		num = limit
	}

	return num
}

func formatToTwoDigit(num string) string {
	if len(num) < 2 {
		return "0" + num
	}
	return num
}
