package copy

import (
	"fmt"
	"gtools/internal/utils"
	"os"
)

type CopyParameters struct {
	PathSrc     string
	PathDest    string
	ExcludePath []string
	IncludePath []string
}

/*
copy src to dest
*/
func CopyDir(param CopyParameters) error {
	return copyDirectory(param.PathSrc, param.PathDest, param)
}

func copyDirectory(src string, dest string, param CopyParameters) error {

	isSubFolder, err := utils.IsSubfolder(dest, src)
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
