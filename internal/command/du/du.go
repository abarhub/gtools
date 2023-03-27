package du

import (
	"fmt"
	"gtools/internal/utils"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

type DuParameters struct {
	Path          string
	HumanReadable bool
	ThresholdStr  string
	MaxDepth      int
	ExcludePath   []string
	IncludePath   []string
}

func diskUsage(currPath string, depth int, maxDepth int, humanReadable bool, threshold int64, out io.Writer, param DuParameters) (sizeResult int64, errResult error) {
	var size int64

	dir, err := os.Open(currPath)
	if err != nil {
		fmt.Println(err)
		return size, nil
	}
	defer func() {
		if tempErr := dir.Close(); tempErr != nil {
			tempErr = fmt.Errorf("error for close %v : %v", currPath, tempErr.Error())
			if errResult == nil {
				errResult = tempErr
			} else {
				errResult = fmt.Errorf("%v. %v", errResult, tempErr.Error())
			}
		}
	}()

	files, err := dir.Readdir(-1)
	if err != nil {
		return 0, err
	}

	for _, file := range files {
		filePath := filepath.Join(currPath, file.Name())
		toScan, err := fileToScan(filePath, param, true)
		if err != nil {
			return 0, err
		} else if toScan {
			if file.IsDir() {
				sizeDir, err := diskUsage(fmt.Sprintf("%s/%s", currPath, file.Name()), depth+1, maxDepth, humanReadable,
					threshold, out, param)
				if err != nil {
					return 0, err
				}
				size += sizeDir
			} else {
				filePath := filepath.Join(currPath, file.Name())
				toScan, err := fileToScan(filePath, param, true)
				if err != nil {
					return 0, err
				} else if toScan {
					size += file.Size()
				}
			}
		}
	}

	if (maxDepth) <= 0 || (maxDepth) >= depth {
		if threshold <= 0 || size >= threshold {
			err := prettyPrintSize(size, humanReadable, out)
			if err != nil {
				return 0, err
			}
			_, err = fmt.Fprintf(out, "\t %s%c\n", currPath, filepath.Separator)
			if err != nil {
				return 0, err
			}
		}
	}

	return size, nil
}

func prettyPrintSize(size int64, humanReadable bool, out io.Writer) error {
	if humanReadable {
		switch {
		case size > 1024*1024*1024:
			_, err := fmt.Fprintf(out, "%.1fG", float64(size)/(1024*1024*1024))
			if err != nil {
				return err
			}
		case size > 1024*1024:
			_, err := fmt.Fprintf(out, "%.1fM", float64(size)/(1024*1024))
			if err != nil {
				return err
			}
		case size > 1024:
			_, err := fmt.Fprintf(out, "%.1fK", float64(size)/1024)
			if err != nil {
				return err
			}
		default:
			_, err := fmt.Fprintf(out, "%d", size)
			if err != nil {
				return err
			}
		}
	} else {
		_, err := fmt.Fprintf(out, "%d", size)
		if err != nil {
			return err
		}
	}
	return nil
}

func DiskUsage(param DuParameters) error {
	return DiskUsageWriter(param, os.Stdout)
}

func DiskUsageWriter(param DuParameters, out io.Writer) error {
	var threshold int64

	var path = param.Path
	var humanReadable = param.HumanReadable
	var thresholdStr = param.ThresholdStr
	var maxDepth = param.MaxDepth

	l := len(thresholdStr)
	if l > 0 {
		t, err := strconv.Atoi(thresholdStr)
		if err != nil { // threshold string ends with a unit char
			i, err := strconv.Atoi((thresholdStr)[0:(l - 1)])
			if err != nil {
				return err
			}

			switch (thresholdStr)[l-1:] {
			case "G":
				t = i * 1024 * 1024 * 1024
			case "M":
				t = i * 1024 * 1024
			case "K":
				t = i * 1024
			}
		}
		threshold = int64(t)
	}

	var dir string

	if path == "" {
		dir = "."
	} else {
		dir = path

		_, err := os.Lstat(dir)
		if err != nil {
			return err
		}
	}

	_, err := diskUsage(dir, 0, maxDepth, humanReadable, threshold, out, param)

	return err
}

func fileToScan(file string, param DuParameters, exclude bool) (bool, error) {
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
