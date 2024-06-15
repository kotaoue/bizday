package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Holidays map[string]string

var years = map[string]Holidays{}

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
	}
	// second argument is a integer
	return printLaterBusinessDays(st, i)
}

func printLaterBusinessDays(st time.Time, i int) error {
	t, err := addBusinessDays(st, i)
	if err != nil {
		return err
	}

	fmt.Println(t.Format("2006-01-02"))
	return nil
}

func printNumberOfBusinessDays(st, et time.Time) error {
	i, err := numberOfBusinessDays(st, et)
	if err != nil {
		return err
	}
	fmt.Println(i)
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

func isHoliday(t time.Time) (bool, error) {
	_, exists := years[t.Format("2006")]
	if !exists {
		if err := setHolidays(t); err != nil {
			return false, err
		}
	}

	holidays := years[t.Format("2006")]
	_, exists = holidays[t.Format("2006-01-02")]
	return exists, nil
}

func isBusinessDay(date time.Time) (bool, error) {
	weekday := date.Weekday()
	if weekday == time.Saturday || weekday == time.Sunday {
		return false, nil
	}
	b, err := isHoliday(date)
	if err != nil {
		return false, err
	}
	return !b, nil
}

func numberOfBusinessDays(startTime, endTime time.Time) (int, error) {
	if startTime.After(endTime) {
		return 0, nil
	}

	i := 0
	for t := startTime; !t.After(endTime); t = t.AddDate(0, 0, 1) {
		b, err := isBusinessDay(t)
		if err != nil {
			return 0, err
		}
		if b {
			i++
		}
	}

	return i, nil
}

func addBusinessDays(st time.Time, i int) (time.Time, error) {
	t := st
	for i > 0 {
		fmt.Println(i, t)
		b, err := isBusinessDay(t)
		if err != nil {
			return time.Time{}, err
		}
		if b {
			i--
		}
		if i > 0 {
			t = t.AddDate(0, 0, 1)
		}
	}
	return t, nil
}

func getJson(t time.Time) error {
	apiURL := fmt.Sprintf("https://holidays-jp.github.io/api/v1/%s/date.json", t.Format("2006"))

	resp, err := http.Get(apiURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http error occurred. status: %s", resp.Status)
	}

	f, err := os.Create(fmt.Sprint(t.Format("2006"), ".json"))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func setHolidays(t time.Time) error {
	// 指定した日付のjsonをローカル読み込む
	name := fmt.Sprint(t.Format("2006"), ".json")
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		// なかったらAPIから取得
		getJson(t)
	}

	_, err = os.Open(name)
	if err != nil {
		return err
	}

	years[t.Format("2006")] = Holidays{
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

	return nil
}
