package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// todo: 休日を取得するAPIを使って動的に取得して、jsonで保存しておく
var holidays = map[string]string{
	"2024-01-01": "元日",
	"2024-01-08": "成人の日",
	"2024-02-11": "建国記念の日",
	"2024-02-12": "建国記念の日 振替休日",
	"2024-02-23": "天皇誕生日",
	"2024-03-20": "春分の日",
	"2024-04-29": "昭和の日",
	"2024-05-03": "憲法記念日",
	"2024-05-04": "みどりの日",
	"2024-05-05": "こどもの日",
	"2024-05-06": "こどもの日 振替休日",
	"2024-07-15": "海の日",
	"2024-08-11": "山の日",
	"2024-08-12": "休日 山の日",
	"2024-09-16": "敬老の日",
	"2024-09-22": "秋分の日",
	"2024-09-23": "秋分の日 振替休日",
	"2024-10-14": "スポーツの日",
	"2024-11-03": "文化の日",
	"2024-11-04": "文化の日 振替休日",
	"2024-11-23": "勤労感謝の日",
}

func main() {
	if err := Main(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func Main() error {
	args := os.Args[1:]

	if len(args) != 2 {
		return fmt.Errorf("arguments should be two dates or a date and integer.")
	}

	st, err := toTime(args[0])
	if err != nil {
		return err
	}

	i, err := strconv.Atoi(args[1])
	if err != nil {
		et, err := toTime(args[1])
		if err != nil {
			return err
		}
		// second argument is a date
		return printNumberOfBusinessDays(st, et)
	} else {
		// second argument is a integer
		return printLaterBusinessDays(st, i)
	}

	return fmt.Errorf("arguments should be two dates or a date and number.")
}

func printLaterBusinessDays(st time.Time, i int) error {
	t := addBusinessDays(st, i)
	fmt.Println(t.Format("2006-01-02"))
	return nil
}

func printNumberOfBusinessDays(st, et time.Time) error {
	fmt.Println(numberOfBusinessDays(st, et))
	return nil
}

func toTime(s string) (time.Time, error) {
	formats := []string{"2006-01-02", "2006/01/02", "2006.01.02"}

	for _, format := range formats {
		t, err := time.Parse(format, s)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid date format: %s\nsupported formats are %s", s, strings.Join(formats, ", "))
}

func isHoliday(date time.Time) bool {
	_, exists := holidays[date.Format("2006-01-02")]
	return exists
}

func isBusinessDay(date time.Time) bool {
	weekday := date.Weekday()
	if weekday == time.Saturday || weekday == time.Sunday {
		return false
	}
	if isHoliday(date) {
		return false
	}
	return true
}

func numberOfBusinessDays(startTime, endTime time.Time) int {
	if startTime.After(endTime) {
		return 0
	}

	i := 0
	for t := startTime; !t.After(endTime); t = t.AddDate(0, 0, 1) {
		if isBusinessDay(t) {
			i++
		}
	}

	return i
}

func addBusinessDays(st time.Time, i int) time.Time {
	t := st
	for i > 0 {
		if isBusinessDay(t) {
			i--
		}
		if i > 0 {
			t = t.AddDate(0, 0, 1)
		}
	}
	return t
}
