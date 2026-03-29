package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// NextDate возвращает следующую дату после now по правилу repeat
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("правило повторения не указано")
	}

	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("некорректная дата: %w", err)
	}

	parts := strings.SplitN(repeat, " ", 2)

	switch parts[0] {
	case "y":
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				return date.Format(dateFormat), nil
			}
		}

	case "d":
		if len(parts) < 2 {
			return "", fmt.Errorf("правило d требует указания интервала")
		}
		n, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil || n < 1 || n > 400 {
			return "", fmt.Errorf("неверный интервал для правила d (допустимо 1–400)")
		}
		for {
			date = date.AddDate(0, 0, n)
			if afterNow(date, now) {
				return date.Format(dateFormat), nil
			}
		}

	case "w":
		if len(parts) < 2 {
			return "", fmt.Errorf("правило w требует указания дней недели")
		}
		wdays, err := parseWeekdays(parts[1])
		if err != nil {
			return "", err
		}
		nowDate := truncateToDate(now)
		if !date.After(nowDate) {
			date = nowDate
		}
		for i := 0; i < 8; i++ {
			if afterNow(date, now) {
				wd := int(date.Weekday())
				if wd == 0 {
					wd = 7
				}
				if wdays[wd] {
					return date.Format(dateFormat), nil
				}
			}
			date = date.AddDate(0, 0, 1)
		}
		return "", fmt.Errorf("не удалось найти следующую дату для правила w")

	case "m":
		if len(parts) < 2 {
			return "", fmt.Errorf("правило m требует указания дней")
		}
		mdays, mmonths, err := parseMonthRule(parts[1])
		if err != nil {
			return "", err
		}
		nowDate := truncateToDate(now)
		if !date.After(nowDate) {
			date = nowDate
		}
		for i := 0; i < 400; i++ {
			if afterNow(date, now) {
				month := int(date.Month())
				day := date.Day()
				lastDay := lastDayOfMonth(date)

				monthOk := len(mmonths) == 0 || mmonths[month]
				dayOk := mdays[day] ||
					(mdays[-1] && day == lastDay) ||
					(mdays[-2] && day == lastDay-1)

				if monthOk && dayOk {
					return date.Format(dateFormat), nil
				}
			}
			date = date.AddDate(0, 0, 1)
		}
		return "", fmt.Errorf("не удалось найти следующую дату для правила m")

	default:
		return "", fmt.Errorf("неподдерживаемое правило повторения: %q", parts[0])
	}
}

func afterNow(date, ref time.Time) bool {
	return date.Format(dateFormat) > ref.Format(dateFormat)
}

func truncateToDate(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func lastDayOfMonth(t time.Time) int {
	return time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func parseWeekdays(s string) (map[int]bool, error) {
	result := make(map[int]bool)
	for _, part := range strings.Split(s, ",") {
		n, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil || n < 1 || n > 7 {
			return nil, fmt.Errorf("недопустимое значение дня недели: %q", part)
		}
		result[n] = true
	}
	return result, nil
}

func parseMonthRule(s string) (map[int]bool, map[int]bool, error) {
	sparts := strings.SplitN(strings.TrimSpace(s), " ", 2)

	days := make(map[int]bool)
	for _, part := range strings.Split(sparts[0], ",") {
		n, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil || n == 0 || n < -2 || n > 31 {
			return nil, nil, fmt.Errorf("недопустимый день месяца: %q", part)
		}
		days[n] = true
	}

	months := make(map[int]bool)
	if len(sparts) > 1 && strings.TrimSpace(sparts[1]) != "" {
		for _, part := range strings.Split(sparts[1], ",") {
			n, err := strconv.Atoi(strings.TrimSpace(part))
			if err != nil || n < 1 || n > 12 {
				return nil, nil, fmt.Errorf("недопустимый месяц: %q", part)
			}
			months[n] = true
		}
	}

	return days, months, nil
}

// nextDateHandler godoc
// @Summary     Вычислить следующую дату
// @Tags        utils
// @Produce     plain
// @Param       now    query string true  "Текущая дата (20060102)"
// @Param       date   query string true  "Дата задачи (20060102)"
// @Param       repeat query string true  "Правило повторения (d N / y / w 1,3 / m 1,15)"
// @Success     200 {string} string "следующая дата в формате 20060102"
// @Failure     400 {string} string "ошибка"
// @Router      /api/nextdate [get]
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse(dateFormat, nowStr)
		if err != nil {
			http.Error(w, "некорректный параметр now", http.StatusBadRequest)
			return
		}
	}

	next, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := w.Write([]byte(next)); err != nil {
		log.Printf("nextdate write error: %v", err)
	}
}
