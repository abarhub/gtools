package copy

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

func TestCopyDirectory(t *testing.T) {
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
			temp1 := path.Join(rootDir, "test1")
			temp2 := path.Join(rootDir, "test2")
			t.Logf("temp1=%v", temp1)
			t.Logf("temp2=%v", temp2)
			err2 := os.Mkdir(temp1, 0755)
			if err2 != nil {
				t.Errorf("CopyDir() error = %v", err2)

			} else {
				err3 := os.WriteFile(path.Join(temp1, "toto.txt"), []byte{1, 2, 3}, 0755)
				if err3 != nil {
					t.Errorf("CopyDir() error = %v", err3)

				} else {
					err4 := os.Mkdir(temp2, 0755)
					if err4 != nil {
						t.Errorf("CopyDir() error = %v", err4)
					} else {
						src := path.Join(rootDir, tt.args.src)
						dest := path.Join(rootDir, tt.args.dest)
						t.Logf("src=%v", src)
						t.Logf("dest=%v", dest)
						param := CopyParameters{}
						param.CreateDestDir = true
						if err := copyDirectory(src, dest, param); (err != nil) != tt.wantErr {
							t.Errorf("CopyDir() error = %v, wantErr %v", err, tt.wantErr)
						} else {
							if !tt.wantErr {
								file := path.Join(dest, "toto.txt")
								content, err5 := os.ReadFile(file)
								if err5 != nil {
									t.Logf("file=%v", file)
									t.Logf("file clean=%v", path.Clean(file))
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

func TestCopyDir(t *testing.T) {
	type args struct {
		src            string
		dest           string
		exclude        []string
		include        []string
		globDoubleStar bool
		createDestDir  bool
		dirDest        map[string][]byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test1", args{"src", "dest", []string{}, []string{}, false, false, defaultFiles}, false},
		{"test2", args{"src", "dest", []string{"*/test1.txt"}, []string{}, false, false, remove(defaultFiles, []string{"test1.txt"})}, false},
		{"test3", args{"src", "dest", []string{"*/test1.txt", "*/test2.txt"}, []string{}, false, false, remove(defaultFiles, []string{"test1.txt", "test2.txt"})}, false},
		{"test4", args{"src", "dest2", []string{}, []string{}, false, false, defaultFiles}, true},
		{"test5", args{"src", "dest2", []string{}, []string{}, false, true, defaultFiles}, false},
		{"test6", args{"src", "dest", []string{"*/*.log"}, []string{}, false, false, remove(defaultFiles, []string{"test5.log", "dir2/test02_4.log"})}, false},
		{"test7", args{"src", "dest", []string{}, []string{"*/*.txt"}, false, false, remove(defaultFiles, []string{"test3.csv", "test4.csv", "test5.log", "dir1/test03.csv", "dir2/test02_2.csv", "dir2/test02_4.log"})}, false},
		{"test8", args{"src", "dest", []string{"*/dir1"}, []string{}, false, false, remove(defaultFiles, []string{"dir1/test01.txt", "dir1/test02.txt", "dir1/test03.csv"})}, false},
		{"test9", args{"src", "dest", []string{"*/dir1"}, []string{"*/*.txt"}, false, false, remove(defaultFiles, []string{"test3.csv", "test4.csv", "test5.log", "dir1/test01.txt", "dir1/test02.txt", "dir1/test03.csv", "dir2/test02_2.csv", "dir2/test02_4.log"})}, false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rootDir := t.TempDir()
			createTestDirectory(t, rootDir)
			param := CopyParameters{
				PathSrc:        path.Join(rootDir, test.args.src),
				PathDest:       path.Join(rootDir, test.args.dest),
				ExcludePath:    test.args.exclude,
				IncludePath:    test.args.include,
				GlobDoubleStar: test.args.globDoubleStar,
				CreateDestDir:  test.args.createDestDir,
			}
			err := CopyDir(param)
			if (err != nil) != test.wantErr {
				t.Errorf("CopyDir() error = %v, wantErr %v", err, test.wantErr)
			} else if err != nil {
				//  pass next test
			} else {
				checkFs(t, test.args.dirDest, path.Join(rootDir, test.args.dest))
			}
		})
	}
}

func checkFs(t *testing.T, dest map[string][]byte, rootDir string) {
	listeFileRef := []string{}
	for filename, content := range dest {
		listeFileRef = append(listeFileRef, normalizePath(filename))
		listeFileRef = addParent(listeFileRef, filename)
		filePath := path.Join(rootDir, filename)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("file %v not copied", filename)
		} else if err != nil {
			t.Errorf("error for read file %v : %v", filename, err)
		} else {
			dat, err := os.ReadFile(filePath)
			if err != nil {
				t.Errorf("error for read file %v : %v", filename, err)
			} else if bytes.Compare(dat, content) != 0 {
				t.Errorf("error for content of file %v : %v != %v", filename, content, dat)
			}
		}
	}
	t.Logf("list files: %v", listeFileRef)

	listeFileFs := []string{}
	err := filepath.Walk(rootDir,
		func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			f, err := filepath.Rel(rootDir, filePath)
			if err != nil {
				return fmt.Errorf("error for relativize %v", filePath)
			}
			t.Logf("file %v", f)
			if f != "." {
				f2 := normalizePath(f)
				listeFileFs = append(listeFileFs, f2)
				if !contains(listeFileRef, f2) {
					t.Errorf("file %v must not be is on FS", f2)
				}
			}
			return nil
		})
	if err != nil {
		t.Errorf("error for read directory %v : %v", rootDir, err)
	} else {

	}

}

func addParent(listPath []string, filename string) []string {
	filename = normalizePath(filename)
	s := filename
	for i := 0; i < 10; i++ {
		s := filepath.Dir(s)
		if s != "." && len(s) > 0 {
			if contains(listPath, s) {
				break
			} else {
				listPath = append(listPath, s)
			}
		} else {
			break
		}
	}

	return listPath
}

func normalizePath(file string) string {
	file = path.Clean(file)
	return strings.ReplaceAll(file, "\\", "/")
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func remove(fileList map[string][]byte, filesToRemove []string) map[string][]byte {
	fileListModified := make(map[string][]byte)
	for filename, content := range fileList {
		if !contains(filesToRemove, filename) {
			fileListModified[filename] = content
		}
	}
	return fileListModified
}

func createTestDirectory(t *testing.T, rootDir string) bool {
	temp1 := path.Join(rootDir, "src")
	temp2 := path.Join(rootDir, "dest")
	t.Logf("temp1=%v", temp1)
	t.Logf("temp2=%v", temp2)
	err := os.Mkdir(temp1, 0755)
	if err != nil {
		t.Errorf("CopyDir() error = %v", err)
		return true
	}
	err3 := os.WriteFile(path.Join(temp1, "test1.txt"), []byte{1, 2, 3}, 0755)
	if err3 != nil {
		t.Errorf("CopyDir() error = %v", err3)
		return true
	}
	err3 = os.WriteFile(path.Join(temp1, "test2.txt"), []byte{1, 2, 3}, 0755)
	if err3 != nil {
		t.Errorf("CopyDir() error = %v", err3)
		return true
	}
	err3 = os.WriteFile(path.Join(temp1, "test3.csv"), []byte{1, 2, 3}, 0755)
	if err3 != nil {
		t.Errorf("CopyDir() error = %v", err3)
		return true
	}
	err3 = os.WriteFile(path.Join(temp1, "test4.csv"), []byte{4, 5, 6}, 0755)
	if err3 != nil {
		t.Errorf("CopyDir() error = %v", err3)
		return true
	}
	err3 = os.WriteFile(path.Join(temp1, "test5.log"), []byte{3, 2, 1}, 0755)
	if err3 != nil {
		t.Errorf("CopyDir() error = %v", err3)
		return true
	}
	temp1dir1 := path.Join(temp1, "dir1")
	err = os.Mkdir(temp1dir1, 0755)
	if err != nil {
		t.Errorf("CopyDir() error = %v", err)
		return true
	}
	err = os.WriteFile(path.Join(temp1dir1, "test01.txt"), []byte{7, 8, 9}, 0755)
	if err != nil {
		t.Errorf("CopyDir() error = %v", err)
		return true
	}
	err = os.WriteFile(path.Join(temp1dir1, "test02.txt"), []byte{7, 8, 9}, 0755)
	if err != nil {
		t.Errorf("CopyDir() error = %v", err)
		return true
	}
	err = os.WriteFile(path.Join(temp1dir1, "test03.csv"), []byte{7, 8, 9}, 0755)
	if err != nil {
		t.Errorf("CopyDir() error = %v", err)
		return true
	}
	temp1dir2 := path.Join(temp1, "dir2")
	err = os.Mkdir(temp1dir2, 0755)
	if err != nil {
		t.Errorf("CopyDir() error = %v", err)
		return true
	}
	err = os.WriteFile(path.Join(temp1dir2, "test02_1.txt"), []byte{4, 5, 6}, 0755)
	if err != nil {
		t.Errorf("CopyDir() error = %v", err)
		return true
	}
	err = os.WriteFile(path.Join(temp1dir2, "test02_2.csv"), []byte{4, 5, 6}, 0755)
	if err != nil {
		t.Errorf("CopyDir() error = %v", err)
		return true
	}
	err = os.WriteFile(path.Join(temp1dir2, "test02_3.txt"), []byte{4, 5, 6}, 0755)
	if err != nil {
		t.Errorf("CopyDir() error = %v", err)
		return true
	}
	err = os.WriteFile(path.Join(temp1dir2, "test02_4.log"), []byte{4, 5, 6}, 0755)
	if err != nil {
		t.Errorf("CopyDir() error = %v", err)
		return true
	}

	err = os.Mkdir(temp2, 0755)
	if err != nil {
		t.Errorf("CopyDir() error = %v", err)
		return true
	}

	return false
}
