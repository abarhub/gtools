package copy

import (
	"fmt"
	"gtools/internal/utils"
	"io"
	"os"
	"path"
	"path/filepath"
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

	isSubFolder, err := utils.IsSubfolder(src, dest)
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

	dest2 := filepath.Clean(dest)

	// create dest if not exists
	if _, err := os.Stat(dest2); os.IsNotExist(err) {
		err = os.Mkdir(dest2, 0755)
		if err != nil {
			return err
		}
	}

	// check dest exists
	if _, err := os.Stat(dest2); err != nil {
		return err
	}

	files, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, f := range files {
		srcFile := path.Join(src, f.Name())
		destFile := path.Join(dest2, f.Name())
		toCopy, err := fileToCopy(srcFile, param)
		if err != nil {
			return err
		} else if toCopy {
			if f.IsDir() {
				err = copyDirectory(srcFile, destFile, param)
				if err != nil {
					return err
				}
			} else {
				err = copyFile(srcFile, destFile)
				if err != nil {
					return err
				}

			}
		}
	}

	return nil
}

func copyFile(srcFile, destFile string) error {
	source, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
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
