package utils

import (
	"fmt"
	"path/filepath"
	"strings"
)

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
