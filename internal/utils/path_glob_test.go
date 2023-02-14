package utils

import "testing"

func TestMatchGlob(t *testing.T) {
	type args struct {
		file string
		glob string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"test1", args{"dir1/test1.txt", "*.txt"}, true, false},
		{"test2", args{"dir1/test1.txt", "*.csv"}, false, false},
		{"test3", args{"dir1/test1.txt", "test1.txt"}, true, false},
		{"test4", args{"dir1/test2.txt", "test1.txt"}, false, false},
		{"test5", args{"test1.txt", "test1.txt"}, true, false},
		{"test6", args{"test2.txt", "test1.txt"}, false, false},
		{"test7", args{"test1.txt", "test?.txt"}, true, false},
		{"test8", args{"test11.txt", "test?.txt"}, false, false},
		{"test9", args{"dir1/test1.txt", "**/*.txt"}, true, false},
		{"test10", args{"dir1/test1.txt", "**/*.doc"}, false, false},
		{"test11", args{"dir1/test1.txt", "**/test*.txt"}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MatchGlob(tt.args.file, tt.args.glob)
			if (err != nil) != tt.wantErr {
				t.Errorf("MatchGlob() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MatchGlob() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convert(t *testing.T) {
	type args struct {
		glob string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convert(tt.args.glob); got != tt.want {
				t.Errorf("convert() = %v, want %v", got, tt.want)
			}
		})
	}
}
