package du

import (
	"bytes"
	"fmt"
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

var result = map[string]string{
	"./dir1/": "9",
	"./dir2/": "12",
	"./":      "36",
}

var result2 = map[string]string{
	"dir1/": "9",
}

var result3 = map[string]string{
	"./dir1/": "3",
	"./dir2/": "6",
	"./":      "18",
}

var result4 = map[string]string{
	"./dir1/": "6",
	"./dir2/": "6",
	"./":      "18",
}

func TestDiskUsage(t *testing.T) {
	empty := []string{}
	type args struct {
		param  DuParameters
		files  map[string][]byte
		result map[string]string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test1", args{DuParameters{"", false, "", -1, empty, empty}, defaultFiles, result}, false},
		{"test2", args{DuParameters{".", false, "", -1, empty, empty}, defaultFiles, result}, false},
		{"test3", args{DuParameters{"dir1", false, "", -1, empty, empty}, defaultFiles, result2}, false},
		{"test4", args{DuParameters{"", false, "", -1, []string{"*.txt"}, empty}, defaultFiles, result3}, false},
		{"test5", args{DuParameters{"", false, "", -1, empty, []string{"*.txt"}}, defaultFiles, result4}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rootDir := t.TempDir()
			dir := rootDir
			testutils.AddFs(t, tt.args.files, dir)
			if !t.Failed() {

				testutils.ChangeDir(t, dir, func() {
					output := new(bytes.Buffer)
					t.Logf("DiskUsage ...")
					err := DiskUsageWriter(tt.args.param, output)
					t.Logf("DiskUsage end")
					if (err != nil) != tt.wantErr {
						t.Errorf("DiskUsage() error = %v, wantErr %v", err, tt.wantErr)
					} else {
						resultWanted := tt.args.result
						map0, err := splitString(output.String())
						if err != nil {
							t.Errorf("DiskUsage() error for split output : %v", err)
						} else if !reflect.DeepEqual(resultWanted, map0) {
							t.Errorf("error for result : %v != %v", resultWanted, map0)
						} else {
							t.Logf("ok for %v and %v", resultWanted, map0)
						}
					}
				})
			}
		})
	}
}

func splitString(s string) (map[string]string, error) {
	lines := strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n")
	res := make(map[string]string)
	for _, s = range lines {
		if len(s) > 0 {
			tab := strings.Split(s, "\t ")
			if len(tab) > 0 {
				if len(tab) == 2 {
					s2 := tab[1]
					s2 = strings.ReplaceAll(s2, "\\", "/")
					res[s2] = tab[0]
				} else {
					return nil, fmt.Errorf("Invalide size for line : %v", s)
				}
			}
		}
	}
	return res, nil
}
