package ls

import (
	"fmt"
	"gtools/internal/utils"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type LsParameters struct {
	Path             string
	LongFormat       bool
	Recurvive        bool
	ExcludePath      []string
	IncludePath      []string
	DisplayDirectory bool
	DisplayAll       bool
}

func LsCommand(param LsParameters) error {

	return lsCommandWriter(param, os.Stdout)
}

func lsCommandWriter(param LsParameters, out io.Writer) error {

	var dir, repInit, directory string
	directory = param.Path

	if strings.Contains(directory, "**") {
		param.Recurvive = true
		i := strings.Index(directory, "**")
		rep := param.Path[0:i]
		include := param.Path[i:]
		directory = rep
		param.IncludePath = append(param.IncludePath, include)
	} else if strings.Contains(directory, "*") || strings.Contains(directory, "?") {
		i := strings.Index(directory, "*")
		j := strings.Index(directory, "?")
		var k, l int
		if i == -1 {
			k = j
		} else if j == -1 {
			k = i
		} else if i > j {
			k = j
		} else {
			k = i
		}
		if k == 0 {
			l = k
		} else {
			n := strings.LastIndexAny(param.Path[0:k], "/\\")
			if n == -1 {
				l = 0
			} else {
				l = n + 1
			}
		}
		rep := param.Path[0:l]
		include := param.Path[l:]
		directory = rep
		param.IncludePath = append(param.IncludePath, include)
	}

	if len(directory) > 0 {
		dir = directory
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

		skip := false
		if !param.DisplayAll && strings.HasPrefix(file.Name(), ".") {
			skip = true
		}

		if !skip {
			toScan, err := fileToScan(filename, param, true)
			if err != nil {
				return err
			} else if toScan {
				if file.IsDir() {
					if param.DisplayDirectory {
						err = display(param, file, filename, out)
						if err != nil {
							return err
						}
					}
					if param.Recurvive {
						err := listFiles(filename, param, filepath.Join(rep, file.Name()), out)
						if err != nil {
							return err
						}
					}
				} else {
					toScan, err := fileToScan(filename, param, false)
					if err != nil {
						return err
					} else if toScan {
						err = display(param, file, filename, out)
						if err != nil {
							return err
						}
					}
				}
			}
		}
	}
	return nil
}

func display(param LsParameters, file os.DirEntry, filename string, out io.Writer) error {
	if param.LongFormat {
		s := ""
		if file.IsDir() {
			s += "D"
		} else {
			s += "F"
		}
		s += " " + filename
		_, err := fmt.Fprintln(out, s)
		if err != nil {
			return err
		}
	} else {
		_, err := fmt.Fprintln(out, filename)
		if err != nil {
			return err
		}
	}
	return nil
}

func fileToScan(file string, param LsParameters, exclude bool) (bool, error) {
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
