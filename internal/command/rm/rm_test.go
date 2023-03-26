package rm

import (
	"bytes"
	"gtools/internal/testutils"
	"os"
	"path"
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

var result = map[string][]byte{
	//"test01.txt": {7, 8, 9},
	//"test02.txt": {7, 8, 9},
	//"test03.csv": {7, 8, 9},
}

var outEmpty = []string{}

func Test_rmCommandWriter(t *testing.T) {
	type args struct {
		param RmParameters
	}
	tests := []struct {
		name      string
		args      args
		files     map[string][]byte
		filesDest map[string][]byte
		out       []string
		wantErr   bool
	}{
		// TODO: Add test cases.
		{"test1", args{param: RmParameters{Path: "src", Verbose: false, Confirmation: false, Recursive: false}},
			defaultFiles, result, outEmpty, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rootDir := t.TempDir()
			dir := rootDir
			src := path.Join(rootDir, "src")
			err := os.Mkdir(src, 0700)
			if err != nil {
				t.Errorf("rmCommandWriter() error = %v", err)
				return
			}
			testutils.AddFs(t, tt.files, src)
			if !t.Failed() {
				testutils.ChangeDir(t, dir, func() {
					out := &bytes.Buffer{}
					err := rmCommandWriter(tt.args.param, out)
					if (err != nil) != tt.wantErr {
						t.Errorf("RmCommandWriter() error = %v, wantErr %v", err, tt.wantErr)
						return
					}
					if _, err := os.Stat(src); !os.IsNotExist(err) {
						t.Errorf("RmCommandWriter() file exists : %v (%v)", src, err)
						return
					}
					if !t.Failed() {
						res, err := splitString(out.String())
						if err != nil {
							t.Errorf("error for split : %v", err)
						} else if !reflect.DeepEqual(tt.out, res) {
							t.Errorf("error for result : %v != %v", tt.out, res)
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
		res = append(res, s)
	}
	if len(res) > 0 && res[len(res)-1] == "" {
		res = res[0 : len(res)-1]
	}
	return res, nil
}
