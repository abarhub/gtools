package ls

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type LsParameters struct {
	Path       string
	LongFormat bool
	Recurvive  bool
}

func LsCommand(param LsParameters) error {

	return lsCommandWriter(param, os.Stdout)
}

func lsCommandWriter(param LsParameters, out io.Writer) error {

	var dir, repInit string
	if len(param.Path) > 0 {
		dir = param.Path
		repInit = dir
	} else {
		dir = "."
		repInit = ""
	}

	return listFiles(dir, param, repInit, out)
}

func listFiles(path string, param LsParameters, rep string, out io.Writer) error {
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
			_, err = fmt.Fprintln(out, s)
			if err != nil {
				return err
			}
		} else {
			_, err = fmt.Fprintln(out, filename)
			if err != nil {
				return err
			}
		}
		if file.IsDir() && param.Recurvive {
			err := listFiles(filepath.Join(path, file.Name()), param, filepath.Join(rep, file.Name()), out)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
