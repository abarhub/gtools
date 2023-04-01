package ls

import (
	"bytes"
	"gtools/internal/testutils"
	"reflect"
	"strings"
	"testing"
)

var defaultFiles = map[string][]byte{
	"test1.txt":         {1, 2, 3},
	"test2.txt":         {1, 2, 3},
	"test3.csv":         {1, 2, 3},
	"test4.csv":         {4, 5, 6},
	"test5.log":         {3, 2, 1},
	"dir1/test01.txt":   {7, 8, 9},
	"dir1/test02.txt":   {7, 8, 9},
	"dir1/test03.csv":   {7, 8, 9},
	"dir2/test02_1.txt": {4, 5, 6},
	"dir2/test02_2.csv": {4, 5, 6},
	"dir2/test02_3.txt": {4, 5, 6},
	"dir2/test02_4.log": {4, 5, 6},
}

var result = []string{"dir1", "dir2", "test1.txt", "test2.txt", "test3.csv", "test4.csv", "test5.log", ""}
var result2 = []string{"dir1/test01.txt", "dir1/test02.txt", "dir1/test03.csv", ""}
var result3 = []string{"dir1", "dir1/test01.txt", "dir1/test02.txt", "dir1/test03.csv",
	"dir2", "dir2/test02_1.txt", "dir2/test02_2.csv", "dir2/test02_3.txt", "dir2/test02_4.log",
	"test1.txt", "test2.txt", "test3.csv", "test4.csv", "test5.log", ""}
var result4 = []string{"test1.txt", "test2.txt", "test3.csv", "test4.csv", "test5.log", ""}
var result5 = []string{"dir1", "dir1/test03.csv",
	"dir2", "dir2/test02_2.csv", "dir2/test02_4.log",
	"test3.csv", "test4.csv", "test5.log", ""}
var result6 = []string{"dir1", "dir1/test01.txt", "dir1/test02.txt",
	"dir2", "dir2/test02_1.txt", "dir2/test02_3.txt",
	"test1.txt", "test2.txt", ""}
var result7 = []string{"dir1", "dir2", "test3.csv", "test4.csv", "test5.log", ""}
var result8 = []string{"dir1", "dir2", "test1.txt", "test2.txt", ""}

func Test_lsCommandWriter(t *testing.T) {
	empty := []string{}
	type args struct {
		param LsParameters
	}
	tests := []struct {
		name    string
		args    args
		files   map[string][]byte
		outList []string
		wantErr bool
	}{
		{"test1", args{LsParameters{"", false, false,
			empty, empty, true, true}}, defaultFiles,
			result, false},
		{"test2", args{LsParameters{"dir1", false, false,
			empty, empty, true, true}}, defaultFiles,
			result2, false},
		{"test3", args{LsParameters{"", false, true,
			empty, empty, true, true}}, defaultFiles,
			result3, false},
		{"test4_hidde_directory", args{LsParameters{"", false, false,
			empty, empty, false, true}}, defaultFiles,
			result4, false},
		{"test5_exclude_recursive", args{LsParameters{"", false, true,
			[]string{"*.txt"}, empty, true, true}}, defaultFiles,
			result5, false},
		{"test6_include_recursive", args{LsParameters{"", false, true,
			empty, []string{"*.txt"}, true, true}}, defaultFiles,
			result6, false},
		{"test7_exclude_no_recursive", args{LsParameters{"", false, false,
			[]string{"*.txt"}, empty, true, true}}, defaultFiles,
			result7, false},
		{"test8_include_no_recursive", args{LsParameters{"", false, false,
			empty, []string{"*.txt"}, true, true}}, defaultFiles,
			result8, false},
		{"test9_path", args{LsParameters{"*.txt", false, false,
			empty, empty, true, true}}, defaultFiles,
			result8, false},
		{"test10_path", args{LsParameters{"**/*.txt", false, false,
			empty, empty, true, true}}, defaultFiles,
			result6, false},
		{"test11_path", args{LsParameters{"test?.txt", false, false,
			empty, empty, true, true}}, defaultFiles,
			result8, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rootDir := t.TempDir()
			dir := rootDir
			testutils.AddFs(t, tt.files, dir)
			if !t.Failed() {
				testutils.ChangeDir(t, dir, func() {
					output := new(bytes.Buffer)
					t.Logf("lsCommandWriter ...")
					err := lsCommandWriter(tt.args.param, output)
					t.Logf("lsCommandWriter ok")
					if (err != nil) != tt.wantErr {
						t.Errorf("lsCommandWriter() error = %v, wantErr %v", err, tt.wantErr)
						return
					} else {
						res, err := splitString(output.String())
						if err != nil {
							t.Errorf("lsCommandWriter() split error = %v", err)
							return
						}
						if !reflect.DeepEqual(tt.outList, res) {
							t.Errorf("error for result : %v != %v", tt.outList, res)
						}
					}
				})
			}
		})
	}
}

func splitString(s string) ([]string, error) {
	lines := strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n")
	res := []string{}
	for _, s = range lines {
		s = strings.Trim(s, " ")
		s = strings.ReplaceAll(s, "\\", "/")
		res = append(res, s)
	}
	return res, nil
}
