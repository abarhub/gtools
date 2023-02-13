package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

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

func CopyDir(src string, dest string) error {

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

	if _, err := os.Stat(dest); err != nil {
		return err
	}

	files, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, f := range files {

		if f.IsDir() {

			err = CopyDir(src+"/"+f.Name(), dest+"/"+f.Name())
			if err != nil {
				return err
			}

		}

		if !f.IsDir() {

			content, err := os.ReadFile(src + "/" + f.Name())
			if err != nil {
				return err

			}

			err = os.WriteFile(dest+"/"+f.Name(), content, 0755)
			if err != nil {
				return err

			}

		}

	}

	return nil
}
