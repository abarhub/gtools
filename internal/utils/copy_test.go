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
		{"test1", args{"test1", "test20"}, false},
		{"test2", args{"test1", "test2"}, false},
		{"test3", args{"test1", "test1"}, true},
		{"test4", args{"test1", "test1_bis"}, false},
		{"test5", args{"test1", "test1/toto"}, true},
		{"test6", args{"test1", "test1/toto/tata"}, true},
		{"test7", args{"test1", "test1/toto/tata/titi"}, true},
		{"test8", args{"test1", "test2/toto/../tata"}, false},
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
						} else {
							if !tt.wantErr {
								content, err5 := os.ReadFile(dest + "/toto.txt")
								if err5 != nil {
									t.Errorf("CopyDir() error = %v", err5)
								} else {
									if len(content) != 3 || content[0] != 1 || content[1] != 2 || content[2] != 3 {
										t.Errorf("file is not valide")
									}
								}
							}
						}
					}
				}
			}
		})
	}
}
