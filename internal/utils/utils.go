package utils

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

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

// track func execution time
func callerName(skip int) string {
	const unknown = "unknown"
	pcs := make([]uintptr, 1)
	n := runtime.Callers(skip+2, pcs)
	if n < 1 {
		return unknown
	}
	frame, _ := runtime.CallersFrames(pcs).Next()
	if frame.Function == "" {
		return unknown
	}
	return frame.Function
}

func FuncExecutionTime() func() {
	name := callerName(1)
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandString(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
