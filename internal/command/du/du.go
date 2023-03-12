package du

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

type DuParameters struct {
	Path          string
	HumanReadable bool
	ThresholdStr  string
	MaxDepth      int
}

func diskUsage(currPath string, depth int, maxDepth int, humanReadable bool, threshold int64) int64 {
	var size int64

	dir, err := os.Open(currPath)
	if err != nil {
		fmt.Println(err)
		return size
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, file := range files {
		if file.IsDir() {
			size += diskUsage(fmt.Sprintf("%s/%s", currPath, file.Name()), depth+1, maxDepth, humanReadable, threshold)
		} else {
			size += file.Size()
		}
	}

	if (maxDepth) <= 0 || (maxDepth) >= depth {
		if threshold == 0 || size >= threshold {
			prettyPrintSize(size, humanReadable)
			fmt.Printf("\t %s%c\n", currPath, filepath.Separator)
		}
	}

	return size
}

func prettyPrintSize(size int64, humanReadable bool) {
	if humanReadable {
		switch {
		case size > 1024*1024*1024:
			fmt.Printf("%.1fG", float64(size)/(1024*1024*1024))
		case size > 1024*1024:
			fmt.Printf("%.1fM", float64(size)/(1024*1024))
		case size > 1024:
			fmt.Printf("%.1fK", float64(size)/1024)
		default:
			fmt.Printf("%d", size)
		}
	} else {
		fmt.Printf("%d", size)
	}
}

func DiskUsage(param DuParameters) error {

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
		var err error
		dir, err = os.Getwd()
		if err != nil {
			return err
		}
	} else {
		dir = path
	}

	_, err := os.Lstat(dir)
	if err != nil {
		return err
	}

	diskUsage(dir, 0, maxDepth, humanReadable, threshold)

	return nil
}
