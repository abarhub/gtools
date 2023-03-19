package mv

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type MvParameters struct {
	PathSrc       string
	PathDest      string
	CopyAndDelete bool
	Verbose       bool
}

func MvCommand(param MvParameters) error {

	if len(param.PathSrc) == 0 {
		return fmt.Errorf("source is empty")
	}
	if len(param.PathDest) == 0 {
		return fmt.Errorf("destination is empty")
	}
	info, err := os.Stat(param.PathSrc)
	if os.IsNotExist(err) {
		return fmt.Errorf("file source %v not exists", param.PathSrc)
	} else if err != nil {
		return fmt.Errorf("error for source : %v", err)
	}
	dest := ""
	if info.IsDir() { // source is directory
		if param.CopyAndDelete {
			return fmt.Errorf("copy and delete not implemented for directory source")
		}
		dest = param.PathDest
	} else { // source is file
		info, err = os.Stat(param.PathDest)
		if os.IsNotExist(err) {
			dest = param.PathDest
		} else if err != nil {
			return fmt.Errorf("error for destination : %v", err)
		} else if info.IsDir() {
			dest = filepath.Join(param.PathDest, filepath.Base(param.PathSrc))
			info, err = os.Stat(dest)
			if os.IsNotExist(err) { // not exists
			} else if err != nil {
				return fmt.Errorf("error for destination : %v", err)
			} else {
				return fmt.Errorf("destination %v exists", dest)
			}
		} else {
			return fmt.Errorf("destination %v exists", param.PathDest)
		}
	}

	if param.PathSrc == dest { // do nothing
		return nil
	} else {
		if param.CopyAndDelete {
			err2 := MoveFile(param.PathSrc, dest, param)
			if err2 != nil {
				s := ""
				for _, err3 := range err2 {
					s = s + err3.Error()
					if len(s) > 0 && s[len(s)-1] != '.' {
						s += "."
					}
				}
				return fmt.Errorf("error for move : %v", s)
			}
		} else {
			if param.Verbose {
				fmt.Printf("%v -> %v\n", param.PathSrc, dest)
			}
			err = os.Rename(param.PathSrc, dest)
			if err != nil {
				return fmt.Errorf("error for move : %v", err)
			}
		}
	}
	return nil
}

func MoveFile(source, destination string, param MvParameters) (result []error) {
	srcClosed := false
	src, err := os.Open(source)
	if err != nil {
		return []error{err}
	}
	defer func(src *os.File) {
		if !srcClosed {
			res := src.Close()
			if res != nil {
				result = append(result, res)
			}
		}
	}(src)
	fi, err := src.Stat()
	if err != nil {
		return []error{err}
	}
	flag := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	perm := fi.Mode() & os.ModePerm
	dst, err := os.OpenFile(destination, flag, perm)
	if err != nil {
		return []error{err}
	}
	defer func(dst *os.File) {
		res := dst.Close()
		if res != nil {
			result = append(result, res)
		}
	}(dst)
	if param.Verbose {
		fmt.Printf("copy %v -> %v\n", source, destination)
	}
	_, err = io.Copy(dst, src)
	if err != nil {
		res := []error{err}
		res2 := os.Remove(destination)
		if res2 != nil {
			res = append(res, res2)
		}
		return res
	}
	srcClosed = true
	err = src.Close()
	if err != nil {
		return []error{err}
	}
	if param.Verbose {
		fmt.Printf("rm %v\n", source)
	}
	err = os.Remove(source)
	if err != nil {
		return []error{err}
	}
	return nil
}
