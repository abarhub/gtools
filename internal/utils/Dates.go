package utils

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type Dates struct {
	Dates []DateCompare
}

type DateDirection int

const (
	DateGT DateDirection = iota
	DateLT
	DateGE
	DateLE
	DateEQ
)

type DateCompare struct {
	Date          time.Time
	DateDirection DateDirection
	WithTime      bool
}

func ParseDates(s []string) (Dates, error) {
	var dates Dates
	if len(s) > 0 {

		tab := s
		for _, t := range tab {
			direction := DateEQ
			if strings.HasPrefix(t, ">=") {
				direction = DateGE
				t = t[2:]
			} else if strings.HasPrefix(t, "<=") {
				direction = DateLE
				t = t[2:]
			} else if strings.HasPrefix(t, ">") {
				direction = DateGT
				t = t[1:]
			} else if strings.HasPrefix(t, "<") {
				direction = DateLT
				t = t[1:]
			} else if strings.HasPrefix(t, "=") {
				direction = DateEQ
				t = t[1:]
			}
			match, _ := regexp.MatchString(`\d{4}-\d{2}-\d{2}(T\d{2}:\d{2}:\d{2})?`, t)

			if match {

				if len(t) == 10 {
					time1, err := time.Parse("2006-01-02", t)
					if err != nil {
						return Dates{}, fmt.Errorf("invalid Date '%s' : %v", t, err)
					}
					dates.Dates = append(dates.Dates, DateCompare{Date: time1, DateDirection: direction, WithTime: false})
				} else if len(t) == 19 {
					time1, err := time.Parse("2006-01-02T15:04:05", t)
					if err != nil {
						return Dates{}, fmt.Errorf("invalid Date '%s' : %v", t, err)
					}
					dates.Dates = append(dates.Dates, DateCompare{Date: time1, DateDirection: direction, WithTime: true})
				} else {
					return Dates{}, fmt.Errorf("invalid Date '%s' : format is 'YYYY-MM-DD(THH:MM:SS)?'", t)
				}

			} else {
				return Dates{}, fmt.Errorf("invalid Date '%s' : format is 'YYYY-MM-DD(THH:MM:SS)?'", t)
			}
		}

	}

	return dates, nil
}

func (dates Dates) IsDateOk(time2 time.Time) bool {

	for _, date := range dates.Dates {
		var t, t2 time.Time
		if date.WithTime {
			t = time2.Truncate(1 * time.Second)
			t2 = date.Date.Truncate(1 * time.Second)
		} else {
			t = time2.Truncate(24 * time.Hour)
			t2 = date.Date.Truncate(24 * time.Hour)
		}

		if date.DateDirection == DateEQ {
			if !(t.Equal(t2)) {
				return false
			}
		} else if date.DateDirection == DateGT {
			if !t.After(t2) {
				return false
			}
		} else if date.DateDirection == DateLT {
			if !t.Before(t2) {
				return false
			}
		} else if date.DateDirection == DateGE {
			if !(t.After(t2) || t.Equal(t2)) {
				return false
			}
		} else if date.DateDirection == DateLE {
			if !(t.Before(t2) || t.Equal(t2)) {
				return false
			}
		}

	}

	return true
}
