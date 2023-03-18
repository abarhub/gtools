package ls

import (
	"fmt"
	"os"
	"path/filepath"
)

type LsParameters struct {
	Path       string
	LongFormat bool
	Recurvive  bool
}

func LsCommand(param LsParameters) error {

	var dir string
	if len(param.Path) > 0 {
		dir = param.Path
	} else {
		dir = "."
	}

	return listFiles(dir, param, "")
}

func listFiles(path string, param LsParameters, rep string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, file := range files {
		filename := filepath.Join(rep, file.Name())
		if param.LongFormat {
			s := ""
			if file.IsDir() {
				s += "D"
			} else {
				s += "F"
			}
			s += " " + filename
			fmt.Println(s)
		} else {
			fmt.Println(filename)
		}
		if file.IsDir() && param.Recurvive {
			err := listFiles(filepath.Join(path, file.Name()), param, filepath.Join(rep, file.Name()))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
