package rm

import (
	"fmt"
	"gtools/internal/utils"
	"io"
	"os"
	"path"
)

type RmParameters struct {
	Path         string
	Confirmation bool
	Recursive    bool
	Verbose      bool
	ExcludePath  []string
	IncludePath  []string
	DryRun       bool
}

func RmCommand(param RmParameters) error {
	return rmCommandWriter(param, os.Stdout)
}

func rmCommandWriter(param RmParameters, out io.Writer) error {

	if len(param.Path) == 0 {
		return fmt.Errorf("file is empty")
	}

	if param.DryRun {
		param.Verbose = true
	}

	info, err := os.Stat(param.Path)
	if os.IsNotExist(err) {
		return fmt.Errorf("file %v not exists", param.Path)
	} else if err != nil {
		return fmt.Errorf("error for source : %v", err)
	}
	if info.IsDir() && param.Recursive {
		err = deleteDirectory(param.Path, param, out)
		if err != nil {
			return err
		}
	} else {
		err = deleteFile(param.Path, param, out)
		if err != nil {
			return err
		}
	}

	return nil
}

func deleteDirectory(pathSrc string, param RmParameters, out io.Writer) error {
	files, err := os.ReadDir(pathSrc)
	if err != nil {
		return err
	}

	for _, f := range files {
		srcFile := path.Join(pathSrc, f.Name())

		toCopy, err := fileToRemove(srcFile, param, true)
		if err != nil {
			return err
		} else if toCopy {
			if f.IsDir() {
				err = deleteDirectory(srcFile, param, out)
				if err != nil {
					return err
				}
			} else {
				err = deleteFile(srcFile, param, out)
				if err != nil {
					return err
				}
			}
		}
	}
	err = deleteFile(pathSrc, param, out)
	if err != nil {
		return err
	}
	return nil
}

func deleteFile(file string, param RmParameters, out io.Writer) error {
	toCopy, err := fileToRemove(file, param, true)
	if err != nil {
		return err
	} else if toCopy {
		toCopy, err := fileToRemove(file, param, false)
		if err != nil {
			return err
		} else if toCopy {
			if param.Verbose {
				_, err = fmt.Fprintf(out, "%v\n", file)
				if err != nil {
					return err
				}
			}
			if !param.DryRun {
				err = os.Remove(file)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func fileToRemove(file string, param RmParameters, exclude bool) (bool, error) {
	if exclude && len(param.ExcludePath) > 0 {
		for _, s := range param.ExcludePath {
			match, err := matchGlob(file, s)
			if err != nil {
				return false, err
			} else if match {
				return false, nil
			}
		}
	}
	if !exclude && len(param.IncludePath) > 0 {
		for _, s := range param.IncludePath {
			match, err := matchGlob(file, s)
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

func matchGlob(file, pattern string) (bool, error) {
	return utils.MatchGlob(file, pattern)
}
