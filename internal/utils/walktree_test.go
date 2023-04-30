package utils

import (
	"testing"
)

func Test_fileToCopy(t *testing.T) {
	// filepath.Match don't work same between Windows and Linux
	var isWindows = isWindows()
	type args struct {
		file      string
		exclude   []string
		include   []string
		isExclude bool
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"exclude_test1", args{"test1.txt", []string{"*.txt"}, []string{}, true}, false, false},
		{"exclude_test2", args{"test2.txt", []string{"*.txt"}, []string{}, true}, false, false},
		{"exclude_test3", args{"test1.csv", []string{"*.txt"}, []string{}, true}, true, false},
		{"exclude_test4", args{"dir1/test1.txt", []string{"*/*.txt"}, []string{}, true}, false, false},
		{"exclude_test5", args{"dir1/test2.txt", []string{"*/*.txt"}, []string{}, true}, false, false},
		{"exclude_test6", args{"dir1/test3.csv", []string{"*/*.txt"}, []string{}, true}, true, false},
		{"exclude_test7", args{"dir1/test1.txt", []string{"*.txt"}, []string{}, true}, !isWindows, false},
		{"exclude_test8", args{"dir1/test2.txt", []string{"*.txt"}, []string{}, true}, !isWindows, false},
		{"exclude_test9", args{"dir1/test3.csv", []string{"*.txt"}, []string{}, true}, true, false},
		{"exclude_test10", args{"test1.txt", []string{"test1.txt"}, []string{}, true}, false, false},
		{"exclude_test11", args{"test2.txt", []string{"test1.txt"}, []string{}, true}, true, false},
		{"exclude_test12", args{"test3.csv", []string{"test1.txt"}, []string{}, true}, true, false},
		{"exclude_test13", args{"test1.txt", []string{"*.txt", "*.csv"}, []string{}, true}, false, false},
		{"exclude_test14", args{"test2.csv", []string{"*.txt", "*.csv"}, []string{}, true}, false, false},
		{"exclude_test15", args{"test4.log", []string{"*.txt", "*.csv"}, []string{}, true}, true, false},

		{"include_test1", args{"test1.txt", []string{}, []string{"*.txt"}, false}, true, false},
		{"include_test2", args{"test2.txt", []string{}, []string{"*.txt"}, false}, true, false},
		{"include_test3", args{"test3.csv", []string{}, []string{"*.txt"}, false}, false, false},
		{"include_test4", args{"test1.txt", []string{}, []string{"*/*.txt"}, false}, false, false},
		{"include_test5", args{"dir/test1.txt", []string{}, []string{"*.txt"}, false}, isWindows, false},
		{"include_test6", args{"dir/test2.csv", []string{}, []string{"*.txt"}, false}, false, false},
		{"include_test7", args{"dir/test1.txt", []string{}, []string{"*/*.txt"}, false}, true, false},
		{"include_test8", args{"dir/test2.csv", []string{}, []string{"*/*.txt"}, false}, false, false},
		{"include_test9", args{"test1.txt", []string{}, []string{}, false}, true, false},
		{"include_test10", args{"test1.txt", []string{}, []string{"*.txt", "*.csv"}, false}, true, false},
		{"include_test11", args{"test2.csv", []string{}, []string{"*.txt", "*.csv"}, false}, true, false},
		{"include_test12", args{"test3.log", []string{}, []string{"*.txt", "*.csv"}, false}, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := Walktree{exclude: tt.args.exclude, include: tt.args.include}
			got, err := obj.fileToCopy(tt.args.file, tt.args.isExclude)
			if (err != nil) != tt.wantErr {
				t.Errorf("fileToCopy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("fileToCopy() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_matchGlob(t *testing.T) {
	// filepath.Match don't work same between Windows and Linux
	var isWindows = isWindows()
	type args struct {
		file    string
		pattern string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"test1", args{"test1.txt", "*.txt"}, true, false},
		{"test2", args{"test2.txt", "*.txt"}, true, false},
		{"test3", args{"test3.csv", "*.txt"}, false, false},
		{"test4", args{"rep/test4.txt", "*.txt"}, isWindows, false},
		{"test5", args{"rep1/rep2/test5.txt", "*.txt"}, isWindows, false},
		{"test6", args{"rep1/rep2/test6.csv", "*.txt"}, false, false},
		{"test7", args{"test1.txt", "test1.txt"}, true, false},
		{"test8", args{"test2.txt", "test1.txt"}, false, false},
		{"test9", args{"rep1/rep2/test1.txt", "*/*.txt"}, isWindows, false},
		{"test10", args{"rep1/rep2/test2.csv", "*/*.txt"}, false, false},

		{"double_star_test1", args{"rep/test1.txt", "**/*.txt"}, true, false},
		{"double_star_test2", args{"rep/test2.txt", "**/*.txt"}, true, false},
		{"double_star_test3", args{"rep/test3.csv", "**/*.txt"}, false, false},
		{"double_star_test4", args{"test1.txt", "*.txt"}, true, false},
		{"double_star_test5", args{"test2.txt", "*.txt"}, true, false},
		{"double_star_test6", args{"test3.csv", "*.txt"}, false, false},
		{"double_star_test7", args{"test1.txt", "test1.txt"}, true, false},
		{"double_star_test8", args{"test2.txt", "test1.txt"}, false, false},
		{"double_star_test9", args{"test1.txt", "**/*.txt"}, true, false},
		{"double_star_test10", args{"rep/rep2/test1.txt", "**/*.txt"}, true, false},
		{"double_star_test11", args{"rep/rep2/test2.csv", "**/*.txt"}, false, false},
		{"double_star_test12", args{"rep/rep2/rep3/test1.txt", "**/rep2/**/*.txt"}, true, false},
		{"double_star_test13", args{"rep/rep4/rep3/test1.txt", "**/rep2/**/*.txt"}, false, false},
		{"double_star_test14", args{"rep/rep2/rep3/test1.txt", "rep/**/*.txt"}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := matchGlob(tt.args.file, tt.args.pattern)
			if (err != nil) != tt.wantErr {
				t.Errorf("matchGlob() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("matchGlob() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func isWindows() bool {
	//return runtime.GOOS == "windows"
	return true
}
