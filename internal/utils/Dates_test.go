package utils

import (
	"reflect"
	"testing"
	"time"
)

func TestDates_IsDateOk(t *testing.T) {
	type fields struct {
		Dates []DateCompare
	}
	type args struct {
		time2 time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"test1", fields{Dates: []DateCompare{createDateCompare("2025-06-14T13:54:36", DateEQ, t)}},
			args{time2: time.Date(2025, time.June, 14, 13, 54, 36, 0, time.UTC)}, true},
		{"test2", fields{Dates: []DateCompare{createDateCompare("2025-06-14T13:54:36", DateEQ, t)}},
			args{time2: time.Date(2025, time.June, 14, 13, 54, 35, 0, time.UTC)}, false},
		{"test3", fields{Dates: []DateCompare{createDateCompare("2025-06-14T13:54:36", DateEQ, t)}},
			args{time2: time.Date(2025, time.June, 14, 13, 54, 37, 0, time.UTC)}, false},
		{"test4", fields{Dates: []DateCompare{createDateCompare("2025-06-14", DateEQ, t)}},
			args{time2: time.Date(2025, time.June, 14, 0, 0, 0, 0, time.UTC)}, true},
		{"test5", fields{Dates: []DateCompare{createDateCompare("2025-06-14", DateEQ, t)}},
			args{time2: time.Date(2025, time.June, 14, 14, 35, 10, 0, time.UTC)}, true},
		{"test6", fields{Dates: []DateCompare{createDateCompare("2025-06-14", DateEQ, t)}},
			args{time2: time.Date(2025, time.June, 13, 0, 0, 0, 0, time.UTC)}, false},
		{"test7", fields{Dates: []DateCompare{createDateCompare("2025-06-14", DateEQ, t)}},
			args{time2: time.Date(2025, time.June, 15, 0, 0, 0, 0, time.UTC)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dates := Dates{
				Dates: tt.fields.Dates,
			}
			if got := dates.IsDateOk(tt.args.time2); got != tt.want {
				t.Errorf("IsDateOk() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseDates(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name    string
		args    args
		want    Dates
		wantErr bool
	}{
		{"test1", args{s: []string{"2025-07-01T15:23:40"}}, Dates{Dates: []DateCompare{createDateCompare("2025-07-01T15:23:40", DateEQ, t)}}, false},
		{"test2", args{s: []string{"toto"}}, Dates{}, true},
		{"test3", args{s: []string{"2023-10-15"}}, Dates{Dates: []DateCompare{createDateCompare("2023-10-15", DateEQ, t)}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDates(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseDates() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func createDateCompare(s string, direction DateDirection, t *testing.T) DateCompare {
	if len(s) == 10 {
		time1, err := time.Parse("2006-01-02", s)
		if err != nil {
			t.Errorf("ParseDates() invalid date %s : %v", s, err)
		}
		return DateCompare{Date: time1, DateDirection: direction, WithTime: false}
	} else if len(s) == 19 {
		time1, err := time.Parse("2006-01-02T15:04:05", s)
		if err != nil {
			t.Errorf("ParseDates() invalid date %s : %v", s, err)
		}
		return DateCompare{Date: time1, DateDirection: direction, WithTime: true}
	} else {
		t.Errorf("ParseDates() invalid date %s", s)
		return DateCompare{}
	}
}

//func createDate(s string, t *testing.T)  {
//	if len(s) == 10 {
//		time1, err := time.Parse("2006-01-02", s)
//		if err != nil {
//			t.Errorf("ParseDates() invalid date %s : %v", s, err)
//		}
//		return DateCompare{Date: time1, DateDirection: direction, WithTime: false}
//	} else if len(s) == 19 {
//		time1, err := time.Parse("2006-01-02T15:04:05", s)
//		if err != nil {
//			t.Errorf("ParseDates() invalid date %s : %v", s, err)
//		}
//		return DateCompare{Date: time1, DateDirection: direction, WithTime: true}
//	} else {
//		t.Errorf("ParseDates() invalid date %s", s)
//		return DateCompare{}
//	}
//}
