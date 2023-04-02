package utils

import (
	"fmt"
	"path/filepath"
	"strings"
)

type SplitDir struct {
	Recusive bool
	Path     string
	GlobPath string
}

/*
*
return true iif subdirectory if subdirectory of directory
*/
func IsSubfolder(directory string, subdirectory string) (bool, error) {
	directory2 := filepath.Clean(directory)
	directory2 = filepath.FromSlash(directory2)
	subdirectory2 := filepath.Clean(subdirectory)
	subdirectory2 = filepath.FromSlash(subdirectory2)
	if directory2 == subdirectory2 {
		return true, nil
	}
	path, err := filepath.Rel(directory2, subdirectory2)
	if err != nil {
		return false, fmt.Errorf("Invalid Path : %v", err)
	}
	if !strings.Contains(path, "..") {
		return true, nil
	} else {
		return false, nil
	}
}

/*
split directory with path and regex
*/
func SplitDirGlob(directory string) SplitDir {
	if strings.Contains(directory, "**") {
		i := strings.Index(directory, "**")
		rep := directory[0:i]
		include := directory[i:]
		result := SplitDir{Recusive: true, Path: rep, GlobPath: include}
		return result
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
			n := strings.LastIndexAny(directory[0:k], "/\\")
			if n == -1 {
				l = 0
			} else {
				l = n + 1
			}
		}
		rep := directory[0:l]
		include := directory[l:]
		result := SplitDir{Recusive: false, Path: rep, GlobPath: include}
		return result
	} else {
		result := SplitDir{Recusive: false, Path: directory, GlobPath: ""}
		return result
	}
}
