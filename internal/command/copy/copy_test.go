package copy

import (
	testutils "gtools/internal/testutils"
	"os"
	"path"
	"runtime"
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

var defaultFilesUnique = map[string][]byte{
	"test1.txt": {1, 2, 3},
}

func TestCopyDirIncludeExclude(t *testing.T) {
	// TODO fix TU in gitub action
	isWindowsOs := isWindows()
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
		{"test2", args{"src", "dest", []string{"*/test1.txt"}, []string{}, false, false, testutils.Remove(defaultFiles, []string{"test1.txt"})}, false},
		{"test3", args{"src", "dest", []string{"*/test1.txt", "*/test2.txt"}, []string{}, false, false, testutils.Remove(defaultFiles, []string{"test1.txt", "test2.txt"})}, false},
		{"test4", args{"src", "dest2", []string{}, []string{}, false, false, defaultFiles}, true},
		{"test5", args{"src", "dest2", []string{}, []string{}, false, true, defaultFiles}, false},
		{"test6", args{"src", "dest", []string{"*/*.log"}, []string{}, false, false, testutils.Remove(defaultFiles, []string{"test5.log", "dir2/test02_4.log"})}, false},
		{"test7", args{"src", "dest", []string{}, []string{"*/*.txt"}, false, false, testutils.Remove(defaultFiles, []string{"test3.csv", "test4.csv", "test5.log", "dir1/test03.csv", "dir2/test02_2.csv", "dir2/test02_4.log"})}, false},
		{"test8", args{"src", "dest", []string{"*/dir1"}, []string{}, false, false, testutils.Remove(defaultFiles, []string{"dir1/test01.txt", "dir1/test02.txt", "dir1/test03.csv"})}, false},
		{"test9", args{"src", "dest", []string{"*/dir1"}, []string{"*/*.txt"}, false, false, testutils.Remove(defaultFiles, []string{"test3.csv", "test4.csv", "test5.log", "dir1/test01.txt", "dir1/test02.txt", "dir1/test03.csv", "dir2/test02_2.csv", "dir2/test02_4.log"})}, false},
		{"test10", args{"src/test1.txt", "dest/test1.txt", []string{}, []string{}, false, false, defaultFilesUnique}, false},
		{"test11", args{"src/test1.txt", "dest", []string{}, []string{}, false, false, defaultFilesUnique}, false},
	}
	for _, test := range tests {
		if !isWindowsOs {
			switch test.name {
			case "test2", "test3", "test6", "test7", "test8", "test9":
				//t.Logf("ignore test: %v", test.name)
				//t.Skipf("ignore test: %v", test.name)
			}
		}
		t.Run(test.name, func(t *testing.T) {
			rootDir := t.TempDir()
			createTestDirectory(t, rootDir)
			t.Logf("exclude: %v", test.args.exclude)
			t.Logf("include: %v", test.args.include)
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
				testutils.CheckFs(t, test.args.dirDest, path.Join(rootDir, test.args.dest))
			}
		})
	}
}

var defaultFilesDryRun1 = map[string][]byte{
	"test1.txt": {1, 2, 3, 4},
}

var defaultFilesDryRunRes = testutils.Add(defaultFiles, defaultFilesDryRun1)

func TestCopyDirDryRun(t *testing.T) {
	type args struct {
		dryRun           bool
		verbose          bool
		copyIfFileExists FileExists
		dirDest          map[string][]byte
		updateDest       map[string][]byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test1", args{false, false, CopyFileExists, defaultFiles, map[string][]byte{}}, false},
		{"test2", args{false, true, CopyFileExists, defaultFiles, map[string][]byte{}}, false},
		{"test3", args{true, false, CopyFileExists, map[string][]byte{}, map[string][]byte{}}, false},
		{"test4", args{false, true, CopyFileExists, defaultFiles, defaultFilesDryRun1}, false},
		{"test5", args{false, true, NoCopyFileExists, defaultFilesDryRunRes, defaultFilesDryRun1}, false},
		{"test6", args{false, true, NoCopyFileExisteSizeFile, defaultFiles, defaultFilesDryRun1}, false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rootDir := t.TempDir()
			createTestDirectory(t, rootDir)
			testutils.AddFs(t, test.args.updateDest, path.Join(rootDir, "dest"))
			t.Logf("dryRun: %v", test.args.dryRun)
			t.Logf("verbose: %v", test.args.verbose)
			t.Logf("copyIfFileExists: %v", test.args.copyIfFileExists)
			param := CopyParameters{
				PathSrc:          path.Join(rootDir, "src"),
				PathDest:         path.Join(rootDir, "dest"),
				DryRun:           test.args.dryRun,
				Verbose:          test.args.verbose,
				CopyIfFileExists: test.args.copyIfFileExists,
			}
			err := CopyDir(param)
			if (err != nil) != test.wantErr {
				t.Errorf("CopyDir() error = %v, wantErr %v", err, test.wantErr)
			} else if err != nil {
				//  pass next test
			} else {
				testutils.CheckFs(t, test.args.dirDest, path.Join(rootDir, "dest"))
			}
		})
	}
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

func Test_fileToCopy(t *testing.T) {
	// filepath.Match don't work same between Windows and Linux
	var isWindows = isWindows()
	type args struct {
		file    string
		param   CopyParameters
		exclude bool
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"exclude_test1", args{"test1.txt", CopyParameters{ExcludePath: []string{"*.txt"}, IncludePath: []string{}}, true}, false, false},
		{"exclude_test2", args{"test2.txt", CopyParameters{ExcludePath: []string{"*.txt"}, IncludePath: []string{}}, true}, false, false},
		{"exclude_test3", args{"test1.csv", CopyParameters{ExcludePath: []string{"*.txt"}, IncludePath: []string{}}, true}, true, false},
		{"exclude_test4", args{"dir1/test1.txt", CopyParameters{ExcludePath: []string{"*/*.txt"}, IncludePath: []string{}}, true}, false, false},
		{"exclude_test5", args{"dir1/test2.txt", CopyParameters{ExcludePath: []string{"*/*.txt"}, IncludePath: []string{}}, true}, false, false},
		{"exclude_test6", args{"dir1/test3.csv", CopyParameters{ExcludePath: []string{"*/*.txt"}, IncludePath: []string{}}, true}, true, false},
		{"exclude_test7", args{"dir1/test1.txt", CopyParameters{ExcludePath: []string{"*.txt"}, IncludePath: []string{}}, true}, !isWindows, false},
		{"exclude_test8", args{"dir1/test2.txt", CopyParameters{ExcludePath: []string{"*.txt"}, IncludePath: []string{}}, true}, !isWindows, false},
		{"exclude_test9", args{"dir1/test3.csv", CopyParameters{ExcludePath: []string{"*.txt"}, IncludePath: []string{}}, true}, true, false},
		{"exclude_test10", args{"test1.txt", CopyParameters{ExcludePath: []string{"test1.txt"}, IncludePath: []string{}}, true}, false, false},
		{"exclude_test11", args{"test2.txt", CopyParameters{ExcludePath: []string{"test1.txt"}, IncludePath: []string{}}, true}, true, false},
		{"exclude_test12", args{"test3.csv", CopyParameters{ExcludePath: []string{"test1.txt"}, IncludePath: []string{}}, true}, true, false},
		{"exclude_test13", args{"test1.txt", CopyParameters{ExcludePath: []string{"*.txt", "*.csv"}, IncludePath: []string{}}, true}, false, false},
		{"exclude_test14", args{"test2.csv", CopyParameters{ExcludePath: []string{"*.txt", "*.csv"}, IncludePath: []string{}}, true}, false, false},
		{"exclude_test15", args{"test4.log", CopyParameters{ExcludePath: []string{"*.txt", "*.csv"}, IncludePath: []string{}}, true}, true, false},

		{"include_test1", args{"test1.txt", CopyParameters{ExcludePath: []string{}, IncludePath: []string{"*.txt"}}, false}, true, false},
		{"include_test2", args{"test2.txt", CopyParameters{ExcludePath: []string{}, IncludePath: []string{"*.txt"}}, false}, true, false},
		{"include_test3", args{"test3.csv", CopyParameters{ExcludePath: []string{}, IncludePath: []string{"*.txt"}}, false}, false, false},
		{"include_test4", args{"test1.txt", CopyParameters{ExcludePath: []string{}, IncludePath: []string{"*/*.txt"}}, false}, false, false},
		{"include_test5", args{"dir/test1.txt", CopyParameters{ExcludePath: []string{}, IncludePath: []string{"*.txt"}}, false}, isWindows, false},
		{"include_test6", args{"dir/test2.csv", CopyParameters{ExcludePath: []string{}, IncludePath: []string{"*.txt"}}, false}, false, false},
		{"include_test7", args{"dir/test1.txt", CopyParameters{ExcludePath: []string{}, IncludePath: []string{"*/*.txt"}}, false}, true, false},
		{"include_test8", args{"dir/test2.csv", CopyParameters{ExcludePath: []string{}, IncludePath: []string{"*/*.txt"}}, false}, false, false},
		{"include_test9", args{"test1.txt", CopyParameters{ExcludePath: []string{}, IncludePath: []string{}}, false}, true, false},
		{"include_test10", args{"test1.txt", CopyParameters{ExcludePath: []string{}, IncludePath: []string{"*.txt", "*.csv"}}, false}, true, false},
		{"include_test11", args{"test2.csv", CopyParameters{ExcludePath: []string{}, IncludePath: []string{"*.txt", "*.csv"}}, false}, true, false},
		{"include_test12", args{"test3.log", CopyParameters{ExcludePath: []string{}, IncludePath: []string{"*.txt", "*.csv"}}, false}, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fileToCopy(tt.args.file, tt.args.param, tt.args.exclude)
			if (err != nil) != tt.wantErr {
				t.Errorf("fileToCopy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("fileToCopy() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_matchGlob(t *testing.T) {
	// filepath.Match don't work same between Windows and Linux
	var isWindows = isWindows()
	type args struct {
		file    string
		pattern string
		param   CopyParameters
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"test1", args{"test1.txt", "*.txt", CopyParameters{GlobDoubleStar: false}}, true, false},
		{"test2", args{"test2.txt", "*.txt", CopyParameters{GlobDoubleStar: false}}, true, false},
		{"test3", args{"test3.csv", "*.txt", CopyParameters{GlobDoubleStar: false}}, false, false},
		{"test4", args{"rep/test4.txt", "*.txt", CopyParameters{GlobDoubleStar: false}}, isWindows, false},
		{"test5", args{"rep1/rep2/test5.txt", "*.txt", CopyParameters{GlobDoubleStar: false}}, isWindows, false},
		{"test6", args{"rep1/rep2/test6.csv", "*.txt", CopyParameters{GlobDoubleStar: false}}, false, false},
		{"test7", args{"test1.txt", "test1.txt", CopyParameters{GlobDoubleStar: false}}, true, false},
		{"test8", args{"test2.txt", "test1.txt", CopyParameters{GlobDoubleStar: false}}, false, false},
		{"test9", args{"rep1/rep2/test1.txt", "*/*.txt", CopyParameters{GlobDoubleStar: false}}, isWindows, false},
		{"test10", args{"rep1/rep2/test2.csv", "*/*.txt", CopyParameters{GlobDoubleStar: false}}, false, false},

		{"double_star_test1", args{"rep/test1.txt", "**/*.txt", CopyParameters{GlobDoubleStar: true}}, true, false},
		{"double_star_test2", args{"rep/test2.txt", "**/*.txt", CopyParameters{GlobDoubleStar: true}}, true, false},
		{"double_star_test3", args{"rep/test3.csv", "**/*.txt", CopyParameters{GlobDoubleStar: true}}, false, false},
		{"double_star_test4", args{"test1.txt", "*.txt", CopyParameters{GlobDoubleStar: true}}, true, false},
		{"double_star_test5", args{"test2.txt", "*.txt", CopyParameters{GlobDoubleStar: true}}, true, false},
		{"double_star_test6", args{"test3.csv", "*.txt", CopyParameters{GlobDoubleStar: true}}, false, false},
		{"double_star_test7", args{"test1.txt", "test1.txt", CopyParameters{GlobDoubleStar: true}}, true, false},
		{"double_star_test8", args{"test2.txt", "test1.txt", CopyParameters{GlobDoubleStar: true}}, false, false},
		{"double_star_test9", args{"test1.txt", "**/*.txt", CopyParameters{GlobDoubleStar: true}}, true, false},
		{"double_star_test10", args{"rep/rep2/test1.txt", "**/*.txt", CopyParameters{GlobDoubleStar: true}}, true, false},
		{"double_star_test11", args{"rep/rep2/test2.csv", "**/*.txt", CopyParameters{GlobDoubleStar: true}}, false, false},
		{"double_star_test12", args{"rep/rep2/rep3/test1.txt", "**/rep2/**/*.txt", CopyParameters{GlobDoubleStar: true}}, true, false},
		{"double_star_test13", args{"rep/rep4/rep3/test1.txt", "**/rep2/**/*.txt", CopyParameters{GlobDoubleStar: true}}, false, false},
		{"double_star_test14", args{"rep/rep2/rep3/test1.txt", "rep/**/*.txt", CopyParameters{GlobDoubleStar: true}}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := matchGlob(tt.args.file, tt.args.pattern, tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("matchGlob() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("matchGlob() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func isWindows() bool {
	return runtime.GOOS == "windows"
}
