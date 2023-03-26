package mv

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
	"test01.txt": {7, 8, 9},
	"test02.txt": {7, 8, 9},
	"test03.csv": {7, 8, 9},
}

var result2 = map[string][]byte{
	"test02_1.txt": {4, 5, 6},
	"test02_2.csv": {4, 5, 6},
	"test02_3.txt": {4, 5, 6},
	"test02_4.log": {4, 5, 6},
}

func Test_mvCommandWriter(t *testing.T) {
	type args struct {
		param MvParameters
	}
	tests := []struct {
		name      string
		args      args
		files     map[string][]byte
		filesDest map[string][]byte
		out       []string
		wantErr   bool
	}{
		{"test1", args{param: MvParameters{PathSrc: "src", PathDest: "dest", Verbose: false, CopyAndDelete: false}}, defaultFiles, defaultFiles, []string{}, false},
		{"test2", args{param: MvParameters{PathSrc: "src/dir1", PathDest: "dest", Verbose: false, CopyAndDelete: false}}, defaultFiles, result, []string{}, false},
		{"test3", args{param: MvParameters{PathSrc: "src/dir2", PathDest: "dest", Verbose: false, CopyAndDelete: false}}, defaultFiles, result2, []string{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rootDir := t.TempDir()
			dir := rootDir
			src := path.Join(rootDir, "src")
			err := os.Mkdir(src, 0700)
			if err != nil {
				t.Errorf("mvCommandWriter() error = %v", err)
				return
			}
			testutils.AddFs(t, tt.files, src)
			if !t.Failed() {
				testutils.ChangeDir(t, dir, func() {
					out := &bytes.Buffer{}
					err := mvCommandWriter(tt.args.param, out)
					if (err != nil) != tt.wantErr {
						t.Errorf("mvCommandWriter() error = %v, wantErr %v", err, tt.wantErr)
						return
					}
					testutils.CheckFs(t, tt.filesDest, path.Join(rootDir, tt.args.param.PathDest))
					res, err := splitString(out.String())
					if err != nil {
						t.Errorf("error for split : %v", err)
					} else if !reflect.DeepEqual(tt.out, res) {
						t.Errorf("error for result : %v != %v", tt.out, res)
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

func compareSliceString(slice, slice2 []string) bool {
	if len(slice) != len(slice2) {
		return false
	}
	for i := 0; i < len(slice); i++ {
		if slice[i] != slice2[i] {
			return false
		}
	}
	return true
}
