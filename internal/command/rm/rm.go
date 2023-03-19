package rm

import (
	"fmt"
	"os"
	"path"
)

type RmParameters struct {
	Path         string
	Confirmation bool
	Recursive    bool
	Verbose      bool
}

func RmCommand(param RmParameters) error {

	if len(param.Path) == 0 {
		return fmt.Errorf("file is empty")
	}

	info, err := os.Stat(param.Path)
	if os.IsNotExist(err) {
		return fmt.Errorf("file %v not exists", param.Path)
	} else if err != nil {
		return fmt.Errorf("error for source : %v", err)
	}
	if info.IsDir() {
		err = deleteDirectory(param.Path, param)
		if err != nil {
			return err
		}
	} else {
		if param.Verbose {
			fmt.Printf("%v\n", param.Path)
		}
		err := os.Remove(param.Path)
		if err != nil {
			return err
		}
	}

	return nil
}

func deleteDirectory(pathSrc string, param RmParameters) error {
	files, err := os.ReadDir(pathSrc)
	if err != nil {
		return err
	}

	for _, f := range files {
		srcFile := path.Join(pathSrc, f.Name())

		if f.IsDir() {
			err = deleteDirectory(srcFile, param)
			if err != nil {
				return err
			}
			if param.Verbose {
				fmt.Printf("%v\n", srcFile)
			}
			err := os.Remove(srcFile)
			if err != nil {
				return err
			}
		} else {
			if param.Verbose {
				fmt.Printf("%v\n", srcFile)
			}
			err := os.Remove(srcFile)
			if err != nil {
				return err
			}
		}
	}
	if param.Verbose {
		fmt.Printf("%v\n", pathSrc)
	}
	err = os.Remove(pathSrc)
	if err != nil {
		return err
	}
	return nil
}
