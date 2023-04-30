package rename

import (
	"bytes"
	"gtools/internal/testutils"
	"os"
	"path"
	"reflect"
	"sort"
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
	"test1.zip":         {1, 2, 3},
	"test2.zip":         {1, 2, 3},
	"test3.csv":         {1, 2, 3},
	"test4.csv":         {4, 5, 6},
	"test5.log":         {3, 2, 1},
	"dir1/test01.zip":   {7, 8, 9},
	"dir1/test02.zip":   {7, 8, 9},
	"dir1/test03.csv":   {7, 8, 9},
	"dir2/test02_1.zip": {4, 5, 6},
	"dir2/test02_2.csv": {4, 5, 6},
	"dir2/test02_3.zip": {4, 5, 6},
	"dir2/test02_4.log": {4, 5, 6},
}

var result2 = map[string][]byte{
	"test1.doc":         {1, 2, 3},
	"test2.doc":         {1, 2, 3},
	"test3.csv":         {1, 2, 3},
	"test4.csv":         {4, 5, 6},
	"test5.log":         {3, 2, 1},
	"dir1/test01.doc":   {7, 8, 9},
	"dir1/test02.doc":   {7, 8, 9},
	"dir1/test03.csv":   {7, 8, 9},
	"dir2/test02_1.doc": {4, 5, 6},
	"dir2/test02_2.csv": {4, 5, 6},
	"dir2/test02_3.doc": {4, 5, 6},
	"dir2/test02_4.log": {4, 5, 6},
}

var result3 = map[string][]byte{
	"test1.zip":         {1, 2, 3},
	"test2.zip":         {1, 2, 3},
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

var empty = []string{}

var out = []string{
	"rename src/dir1/test01.txt -> src/dir1/test01.zip",
	"rename src/dir1/test02.txt -> src/dir1/test02.zip",
	"rename src/dir2/test02_1.txt -> src/dir2/test02_1.zip",
	"rename src/dir2/test02_3.txt -> src/dir2/test02_3.zip",
	"rename src/test1.txt -> src/test1.zip",
	"rename src/test2.txt -> src/test2.zip",
}

func Test_renameCommandWriter(t *testing.T) {
	type args struct {
		param RenameParameters
	}
	tests := []struct {
		name    string
		args    args
		init    map[string][]byte
		res     map[string][]byte
		wantOut []string
		wantErr bool
	}{
		{"test1", args{param: RenameParameters{Files: ".txt", FilesRenamed: ".zip", Recursive: true,
			Verbose: false, DryRun: false, Directory: ""}}, defaultFiles, result, empty, false},
		{"test2", args{param: RenameParameters{Files: ".txt", FilesRenamed: ".doc", Recursive: true,
			Verbose: false, DryRun: false, Directory: ""}}, defaultFiles, result2, empty, false},
		{"test3", args{param: RenameParameters{Files: ".txt", FilesRenamed: ".zip", Recursive: true,
			Verbose: true, DryRun: false, Directory: ""}}, defaultFiles, result, out, false},
		{"test4_no_recursive", args{param: RenameParameters{Files: ".txt", FilesRenamed: ".zip", Recursive: false,
			Verbose: false, DryRun: false, Directory: "src"}}, defaultFiles, result3, empty, false},
		{"test5_", args{param: RenameParameters{Files: ".txt", FilesRenamed: ".zip", Recursive: true,
			Verbose: false, DryRun: false, Directory: "src"}}, defaultFiles, result, empty, false},
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
			testutils.AddFs(t, tt.init, src)
			if !t.Failed() {
				testutils.ChangeDir(t, dir, func() {
					out := &bytes.Buffer{}
					err := renameCommandWriter(tt.args.param, out)
					if (err != nil) != tt.wantErr {
						t.Errorf("renameCommandWriter() error = %v, wantErr %v", err, tt.wantErr)
						return
					}
					testutils.CheckFs(t, tt.res, src)
					if !t.Failed() {
						res, err := splitString(out.String())
						if err != nil {
							t.Errorf("error for split : %v", err)
						} else {
							normalize := normalize(res)
							sort.Strings(normalize)
							tmp := []string{}
							tmp = append(tt.wantOut)
							sort.Strings(tmp)
							if !reflect.DeepEqual(tmp, normalize) {
								t.Errorf("renameCommandWriter() error for result : %v != %v", tmp, normalize)
							}
						}

						//if gotOut := out.String(); gotOut != tt.wantOut {
						//	t.Errorf("renameCommandWriter() gotOut = %v, want %v", gotOut, tt.wantOut)
						//}
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

func normalize(list []string) []string {
	res := []string{}
	for _, s := range list {
		res = append(res, strings.ReplaceAll(s, "\\", "/"))
	}
	return res
}
