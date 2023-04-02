package utils

import (
	"reflect"
	"testing"
)

func TestIsSubfolder(t *testing.T) {
	type args struct {
		src  string
		dest string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"test1", args{"test1", "test2"}, false, false},
		{"test2", args{"test1/test2", "test2"}, false, false},
		{"test3", args{"test1", "test1/test2"}, true, false},
		{"test4", args{"/tmp/test1", "/usr/test2"}, false, false},
		{"test5", args{"/tmp/test1", "/tmp/test1/test2"}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsSubfolder(tt.args.src, tt.args.dest)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsSubfolder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsSubfolder() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSplitDirGlob(t *testing.T) {
	type args struct {
		directory string
	}
	tests := []struct {
		name string
		args args
		want SplitDir
	}{
		{"test1", args{"abc"}, SplitDir{Recusive: false, Path: "abc", GlobPath: ""}},
		{"test2", args{"temp/toto/tata"}, SplitDir{Recusive: false, Path: "temp/toto/tata", GlobPath: ""}},
		{"test3", args{"temp/toto/tata/*.txt"}, SplitDir{Recusive: false, Path: "temp/toto/tata/", GlobPath: "*.txt"}},
		{"test4", args{"myDir/test*.tgz"}, SplitDir{Recusive: false, Path: "myDir/", GlobPath: "test*.tgz"}},
		{"test5", args{"myDir/*"}, SplitDir{Recusive: false, Path: "myDir/", GlobPath: "*"}},
		{"test6", args{"*.txt"}, SplitDir{Recusive: false, Path: "", GlobPath: "*.txt"}},
		{"test7", args{"myDir/**/*.zip"}, SplitDir{Recusive: true, Path: "myDir/", GlobPath: "**/*.zip"}},
		{"test8", args{"myDir/**"}, SplitDir{Recusive: true, Path: "myDir/", GlobPath: "**"}},
		{"test9", args{"**"}, SplitDir{Recusive: true, Path: "", GlobPath: "**"}},
		{"test10", args{"myDir/test?.tgz"}, SplitDir{Recusive: false, Path: "myDir/", GlobPath: "test?.tgz"}},
		{"test11", args{"test?.tgz"}, SplitDir{Recusive: false, Path: "", GlobPath: "test?.tgz"}},
		{"test12", args{"?"}, SplitDir{Recusive: false, Path: "", GlobPath: "?"}},
		{"test13", args{"myDir/tot*/test?.tgz"}, SplitDir{Recusive: false, Path: "myDir/", GlobPath: "tot*/test?.tgz"}},
		{"test14", args{"myDir/test?/tit*.tgz"}, SplitDir{Recusive: false, Path: "myDir/", GlobPath: "test?/tit*.tgz"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SplitDirGlob(tt.args.directory); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SplitDirGlob() = %v, want %v", got, tt.want)
			}
		})
	}
}
