package testutils

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

func CheckFs(t *testing.T, dest map[string][]byte, rootDir string) {
	listeFileRef := []string{}
	for filename, content := range dest {
		listeFileRef = append(listeFileRef, normalizePath(filename))
		listeFileRef = addParent(listeFileRef, filename)
		var filePath string
		if path.Ext(rootDir) == ".txt" || path.Ext(rootDir) == ".csv" || path.Ext(rootDir) == ".log" {
			// dest is a file
			filePath = rootDir
		} else {
			filePath = path.Join(rootDir, filename)
		}
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("file %v not copied", filename)
		} else if err != nil {
			t.Errorf("error for read file %v : %v", filename, err)
		} else {
			dat, err := os.ReadFile(filePath)
			if err != nil {
				t.Errorf("error for read file %v : %v", filename, err)
			} else if bytes.Compare(dat, content) != 0 {
				t.Errorf("error for content of file %v : %v != %v", filename, content, dat)
			}
		}
	}
	t.Logf("list files: %v", listeFileRef)

	listeFileFs := []string{}
	err := filepath.Walk(rootDir,
		func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			f, err := filepath.Rel(rootDir, filePath)
			if err != nil {
				return fmt.Errorf("error for relativize %v", filePath)
			}
			t.Logf("file %v", f)
			if f != "." {
				f2 := normalizePath(f)
				listeFileFs = append(listeFileFs, f2)
				if !contains(listeFileRef, f2) {
					t.Errorf("file %v must not be is on FS", f2)
				}
			}
			return nil
		})
	if err != nil {
		t.Errorf("error for read directory %v : %v", rootDir, err)
	} else {

	}

}

func AddFs(t *testing.T, dest map[string][]byte, rootDir string) {
	listeFileRef := []string{}
	for filename, content := range dest {
		listeFileRef = append(listeFileRef, normalizePath(filename))
		listeFileRef = addParent(listeFileRef, filename)
		var filePath string
		if path.Ext(rootDir) == ".txt" || path.Ext(rootDir) == ".csv" || path.Ext(rootDir) == ".log" {
			// dest is a file
			filePath = rootDir
		} else {
			filePath = path.Join(rootDir, filename)
		}
		parent := filepath.Dir(filePath)
		if parent != "" && parent != "." {
			err := os.MkdirAll(parent, 0600)
			if err != nil {
				t.Errorf("error for create dir %v : %v", parent, err)
				return
			}
		}
		err := os.WriteFile(filePath, content, 0600)
		if err != nil {
			t.Errorf("error for create file %v : %v", filename, err)
			return
		}
	}
	t.Logf("list files: %v", listeFileRef)
}

func addParent(listPath []string, filename string) []string {
	filename = normalizePath(filename)
	s := filename
	for i := 0; i < 10; i++ {
		s := filepath.Dir(s)
		if s != "." && len(s) > 0 {
			if contains(listPath, s) {
				break
			} else {
				listPath = append(listPath, s)
			}
		} else {
			break
		}
	}

	return listPath
}

func normalizePath(file string) string {
	file = path.Clean(file)
	return strings.ReplaceAll(file, "\\", "/")
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func Remove(fileList map[string][]byte, filesToRemove []string) map[string][]byte {
	fileListModified := make(map[string][]byte)
	for filename, content := range fileList {
		if !contains(filesToRemove, filename) {
			fileListModified[filename] = content
		}
	}
	return fileListModified
}

func Add(fileList map[string][]byte, filesToAdd map[string][]byte) map[string][]byte {
	fileListModified := make(map[string][]byte)
	for filename, content := range fileList {
		fileListModified[filename] = content
	}
	for filename, content := range filesToAdd {
		fileListModified[filename] = content
	}
	return fileListModified
}
