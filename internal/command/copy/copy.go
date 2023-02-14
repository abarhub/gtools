package copy

import (
	"fmt"
	"gtools/internal/utils"
	"os"
	"path/filepath"
	"strings"
)

type CopyParameters struct {
	PathSrc     string
	PathDest    string
	ExcludePath []string
	IncludePath []string
}

/*
*
return true iif dest if subdirectory of src
*/
func isSubfolder(src string, dest string) (bool, error) {
	if src == dest {
		return true, nil
	}
	path, err := filepath.Rel(src, dest)
	if err != nil {
		return false, fmt.Errorf("Invalid Path : %v", err)
	}
	if strings.Contains(path, "../") || path == ".." {
		return true, nil
	} else {
		return false, nil
	}
}

/*
copy src to dest
*/
func CopyDir(param CopyParameters) error {
	return copyDirectory(param.PathSrc, param.PathDest, param)
}

func copyDirectory(src string, dest string, param CopyParameters) error {

	isSubFolder, err := isSubfolder(dest, src)
	if err != nil {
		return err
	} else if isSubFolder {
		return fmt.Errorf("Cannot copy a folder into the folder itself!")
	}

	f, err := os.Open(src)
	if err != nil {
		return err
	}

	file, err := f.Stat()
	if err != nil {
		return err
	}
	if !file.IsDir() {
		return fmt.Errorf("Source " + file.Name() + " is not a directory!")
	}

	// create dest if not exists
	if _, err := os.Stat(dest); os.IsNotExist(err) {
		err = os.Mkdir(dest, 0755)
		if err != nil {
			return err
		}
	}

	// check dest exists
	if _, err := os.Stat(dest); err != nil {
		return err
	}

	files, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, f := range files {
		srcFile := src + "/" + f.Name()
		toCopy, err := fileToCopy(srcFile, param)
		if err != nil {
			return err
		} else if toCopy {
			if f.IsDir() {
				err = copyDirectory(srcFile, dest+"/"+f.Name(), param)
				if err != nil {
					return err
				}
			} else {
				content, err := os.ReadFile(srcFile)
				if err != nil {
					return err
				}

				err = os.WriteFile(dest+"/"+f.Name(), content, 0755)
				if err != nil {
					return err
				}

			}
		}
	}

	return nil
}

func fileToCopy(file string, param CopyParameters) (bool, error) {
	if len(param.ExcludePath) > 0 {
		for _, s := range param.ExcludePath {
			match, err := utils.MatchGlob(file, s)
			if err != nil {
				return false, err
			} else if match {
				return false, nil
			}
		}
	}
	if len(param.IncludePath) > 0 {
		for _, s := range param.IncludePath {
			match, err := utils.MatchGlob(file, s)
			if err != nil {
				return false, err
			} else if match {
				return true, nil
			}
		}
		return false, nil
	}
	return true, nil
}
