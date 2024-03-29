package copy

import (
	"fmt"
	"gtools/internal/utils"
	"io"
	"os"
	"path"
	"path/filepath"
)

type FileExists int

const (
	CopyFileExists FileExists = iota
	NoCopyFileExists
	NoCopyFileExisteSizeFile
)

type CopyParameters struct {
	PathSrc          string
	PathDest         string
	ExcludePath      []string
	IncludePath      []string
	CreateDestDir    bool
	GlobDoubleStar   bool
	Verbose          bool
	DryRun           bool
	CopyIfFileExists FileExists
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

	file, err := os.Stat(src)
	if err != nil {
		return err
	}

	dest2 := filepath.Clean(dest)

	if dest2 == "." {
		return fmt.Errorf("Destination " + dest + " is invalid !")
	}

	if !file.IsDir() {
		err = copyFile2(src, dest2, param)
		return err
	} else {
		err = copyDir2(src, dest, param)
		return err
	}

}

func copyDir2(src string, dest string, param CopyParameters) error {

	_, err := os.Stat(src)
	if err != nil {
		return err
	}

	dest2 := filepath.Clean(dest)

	if dest2 == "." {
		return fmt.Errorf("Destination %v is invalid !", dest)
	}

	if param.CreateDestDir && !param.DryRun {
		// create dest if not exists
		if _, err := os.Stat(dest2); os.IsNotExist(err) {
			err = os.Mkdir(dest2, 0755)
			if err != nil {
				return err
			}
		}
	}

	walk, err := utils.CreateWalktree(src, param.ExcludePath, param.IncludePath)
	if err != nil {
		return err
	}

	walk.SetFileParse(func(srcFile string, destFile string) error {
		if !param.DryRun {
			err = createDirIfNeeded(destFile)
			if err != nil {
				return err
			}
		}
		errors := copyFile(srcFile, destFile, param)
		err = convertErrorArryToError(errors)
		if err != nil {
			return err
		}
		return nil
	})

	walk.SetDir2(dest)

	err = walk.Parse()

	return err
}

func convertErrorArryToError(errors []error) error {
	if errors != nil && len(errors) > 0 {
		s := ""
		var e error
		for _, e = range errors {
			s = s + e.Error()
		}
		return fmt.Errorf(s)
	} else {
		return nil
	}
}

func copyFile2(src, dest string, param CopyParameters) error {
	f2, err := os.Stat(dest)
	if err == nil { // dest exists
		var srcFile, destFile string
		if f2.IsDir() {
			filename := path.Base(src)
			srcFile = src
			destFile = path.Join(dest, filename)
		} else {
			srcFile = src
			destFile = dest
		}
		errors := copyFile(srcFile, destFile, param)
		return convertErrorArryToError(errors)
	} else if os.IsNotExist(err) { // dest not exists
		if param.CreateDestDir {
			parent := path.Dir(dest)
			if _, err := os.Stat(parent); os.IsNotExist(err) {
				err = os.Mkdir(parent, 0755)
				if err != nil {
					return err
				}
			}
		}
		errors := copyFile(src, dest, param)
		return convertErrorArryToError(errors)

	} else {
		return err
	}
}

func createDirIfNeeded(file string) error {
	parent := filepath.Dir(file)
	if parent != "" && parent != "." && parent != "/" {
		if _, err := os.Stat(parent); os.IsNotExist(err) {
			err = os.Mkdir(parent, 0755)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func copyFile(srcFile, destFile string, param CopyParameters) (errorResult []error) {
	ok, err := copyFileToDest(destFile, param, srcFile)
	if err != nil {
		return []error{fmt.Errorf("error for check if dest must be create file %v : %v", destFile, err.Error())}
	} else if ok { // copy of file
		if param.DryRun {
			fmt.Printf("%v -> %v\n", srcFile, destFile)
			return nil
		} else {
			source, err := os.Open(srcFile)
			if err != nil {
				return []error{fmt.Errorf("error for open src file %v : %v", srcFile, err.Error())}
			}
			defer func() {
				if tempErr := source.Close(); tempErr != nil {
					tempErr = fmt.Errorf("error for close source file %v : %v", srcFile, tempErr.Error())
					errorResult = append(errorResult, tempErr)
				}
			}()

			destination, err := os.Create(destFile)
			if err != nil {
				return []error{fmt.Errorf("error for create dest file %v : %v", destFile, err.Error())}
			}
			defer func() {
				if tempErr := destination.Close(); tempErr != nil {
					tempErr = fmt.Errorf("error for close dest file %v : %v", destFile, tempErr.Error())
					errorResult = append(errorResult, tempErr)
				}
			}()
			if param.Verbose {
				fmt.Printf("%v -> %v\n", srcFile, destFile)
			}
			_, err = io.Copy(destination, source)
			if err != nil {
				return []error{fmt.Errorf("error for copy from %v to %v: %v", srcFile, destFile, err.Error())}
			} else {
				return nil
			}
		}
	} else { // no copy
		return nil
	}
}

func copyFileToDest(file string, param CopyParameters, srcFile string) (bool, error) {
	switch param.CopyIfFileExists {
	case CopyFileExists:
		return true, nil
	case NoCopyFileExists:
		_, err := os.Stat(file)
		if os.IsNotExist(err) {
			// file not exists => copy
			return true, nil
		} else {
			// file exists => no copy
			return false, nil
		}
	case NoCopyFileExisteSizeFile:
		fDest, err := os.Stat(file)
		if os.IsNotExist(err) {
			// file not exists => copy
			return true, nil
		} else {
			// file exists => check size
			fSrc, err2 := os.Stat(srcFile)
			if err2 != nil {
				return false, err2
			}
			return fDest.Size() != fSrc.Size(), nil
		}
	default:
		return false, fmt.Errorf("Invalide param copy if exists : %v !", param.CopyIfFileExists)
	}
}
