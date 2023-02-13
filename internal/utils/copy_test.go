package utils

import (
	"os"
	"testing"
)

func TestCopyDir(t *testing.T) {
	type args struct {
		src  string
		dest string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"test1", args{"test1", "test20"}, false},
		{"test1", args{"test1", "test2"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			t.Logf("src param=%v", tt.args.src)
			t.Logf("dest param=%v", tt.args.dest)
			rootDir := t.TempDir()
			temp1 := rootDir + "/test1"
			temp2 := rootDir + "/test2"
			t.Logf("temp1=%v", temp1)
			t.Logf("temp2=%v", temp2)
			err2 := os.Mkdir(temp1, 0755)
			if err2 != nil {
				t.Errorf("CopyDir() error = %v", err2)

			} else {
				err3 := os.WriteFile(temp1+"/toto.txt", []byte{1, 2, 3}, 0755)
				if err3 != nil {
					t.Errorf("CopyDir() error = %v", err3)

				} else {
					err4 := os.Mkdir(temp2, 0755)
					if err4 != nil {
						t.Errorf("CopyDir() error = %v", err4)
					} else {
						src := rootDir + "/" + tt.args.src
						dest := rootDir + "/" + tt.args.dest
						t.Logf("src=%v", src)
						t.Logf("dest=%v", dest)
						if err := CopyDir(src, dest); (err != nil) != tt.wantErr {
							t.Errorf("CopyDir() error = %v, wantErr %v", err, tt.wantErr)
						}
					}
				}
			}
		})
	}
}
