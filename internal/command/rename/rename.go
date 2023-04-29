package rename

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type RenameParameters struct {
	Files        string
	FilesRenamed string
	Recursive    bool
	Verbose      bool
	DryRun       bool
	Directory    string
}

func RenameCommand(param RenameParameters) error {
	return renameCommandWriter(param, os.Stdout)
}

func renameCommandWriter(param RenameParameters, out io.Writer) error {
	if len(param.Files) == 0 {
		return fmt.Errorf("source is empty")
	}
	if len(param.FilesRenamed) == 0 {
		return fmt.Errorf("destination is empty")
	}

	if param.DryRun {
		param.Verbose = true
	}

	var dir = "."
	if len(param.Directory) > 0 {
		dir = param.Directory
	}

	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return fmt.Errorf("directory '%v' not exists", dir)
	} else if err != nil {
		return fmt.Errorf("error for directory : %v", err)
	}
	if info.IsDir() { // dir is directory
		return rename(dir, param, out)
	} else { // dir is file
		return fmt.Errorf("%v is not a directory", dir)
	}
}

func rename(dir string, param RenameParameters, out io.Writer) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		filename := filepath.Join(dir, file.Name())
		if file.IsDir() {
			err = rename(filename, param, out)
			if err != nil {
				return err
			}
		} else {
			name := file.Name()
			if strings.Contains(name, param.Files) {
				s := strings.ReplaceAll(name, param.Files, param.FilesRenamed)
				filename2 := filepath.Join(dir, s)
				if param.Verbose {
					_, err := fmt.Fprintf(out, "rename %v -> %v\n", filename, filename2)
					if err != nil {
						return err
					}
				}
				if !param.DryRun {
					err = os.Rename(filename, filename2)
					if err != nil {
						return fmt.Errorf("error for rename of %v to %v : %v", filename, filename2, err)
					}
				}
			}
		}
	}
	return nil
}
