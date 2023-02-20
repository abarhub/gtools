package utils

import (
	"fmt"
	"testing"
)

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
		{"test12", args{"dir1/test1.txt", "**/test1.txt"}, true, false},
		{"test13", args{"dir1/test2.txt", "**/test1.txt"}, false, false},
		{"test14", args{"dir1/dir2/dir3/test1.txt", "**/dir2/**/test1.txt"}, true, false},
		{"test15", args{"dir1/dir5/dir3/test1.txt", "**/dir2/**/test1.txt"}, false, false},
		{"test16", args{"dir1/dir5/dir3/test1.txt", "**/dir1/**/test1.txt"}, true, false},
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
		{"test1", args{"abc"}, "abc"},
		{"test2", args{"abc.toto"}, "abc\\.toto"},
		{"test3", args{"abc*toto"}, "abc[^\\\\/]*toto"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convert(tt.args.glob); got != tt.want {
				t.Errorf("convert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleMatchGlob() {
	fmt.Println(MatchGlob("abc.txt", "*.txt"))
	fmt.Println(MatchGlob("abc.doc", "*.txt"))
	fmt.Println(MatchGlob("toto/tata/test.txt", "**/*.txt"))
	fmt.Println(MatchGlob("toto/tata/test.doc", "**/*.txt"))
	// Output:
	// true <nil>
	// false <nil>
	// true <nil>
	// false <nil>
}
