package rm

import (
	"fmt"
	"io"
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
	return rmCommandWriter(param, os.Stdout)
}

func rmCommandWriter(param RmParameters, out io.Writer) error {

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
		if param.Recursive {
			err = deleteDirectory(param.Path, param, out)
			if err != nil {
				return err
			}
		} else {
			if param.Verbose {
				_, err := fmt.Fprintf(out, "%v\n", param.Path)
				if err != nil {
					return err
				}
			}
			err := os.Remove(param.Path)
			if err != nil {
				return err
			}
		}
	} else {
		if param.Verbose {
			_, err := fmt.Fprintf(out, "%v\n", param.Path)
			if err != nil {
				return err
			}
		}
		err := os.Remove(param.Path)
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

		if f.IsDir() {
			err = deleteDirectory(srcFile, param, out)
			if err != nil {
				return err
			}
			//if param.Verbose {
			//	_, err = fmt.Fprintf(out, "%v\n", srcFile)
			//	if err != nil {
			//		return err
			//	}
			//}
			//err := os.Remove(srcFile)
			//if err != nil {
			//	return err
			//}
		} else {
			if param.Verbose {
				_, err = fmt.Fprintf(out, "%v\n", srcFile)
				if err != nil {
					return err
				}
			}
			err := os.Remove(srcFile)
			if err != nil {
				return err
			}
		}
	}
	if param.Verbose {
		_, err = fmt.Fprintf(out, "%v\n", pathSrc)
		if err != nil {
			return err
		}
	}
	err = os.Remove(pathSrc)
	if err != nil {
		return err
	}
	return nil
}
