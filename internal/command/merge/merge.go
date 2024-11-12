package merge

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type MergeParameters struct {
	File       string
	OutputFile string
}

func MergeCommand(param MergeParameters) error {
	return merge(param)
}

func merge(param MergeParameters) error {

	if param.File == "" {
		return fmt.Errorf("file is required")
	}
	if _, err := os.Stat(param.File); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("file not exists")
	}

	filename := filepath.Base(param.File)
	directory := filepath.Dir(param.File)

	regex := "^\\.[0-9]{3}$"
	extension := filepath.Ext(filename)
	match, _ := regexp.MatchString(regex, extension)
	if !match {
		return fmt.Errorf("filename must end with \\.[0-9]{3}$")
	}

	filenameFinal := filename[:len(filename)-len(extension)]

	var fileFinal string
	if len(param.OutputFile) > 0 {
		fileFinal = param.OutputFile
	} else {
		fileFinal = filepath.Join(directory, filenameFinal)
	}

	files, err := os.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}

	liste := []string{filename}
	for _, file := range files {
		if !file.IsDir() {
			name := file.Name()
			if name != filename && strings.HasPrefix(name, filenameFinal) && len(name) == len(filename) {
				match2, _ := regexp.MatchString(regex, filepath.Ext(filename))
				if match2 {
					liste = append(liste, name)
				}
			}
		}

	}
	sort.Strings(liste)

	for i, name := range liste {
		err := copyFile(directory, name, fileFinal, i)
		if err != nil {
			return err
		}
	}

	return nil
}

func copyFile(directory, name, fileFinal string, i int) (err error) {
	var flag int
	if i == 0 {
		flag = os.O_TRUNC | os.O_CREATE | os.O_WRONLY
	} else {
		flag = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	}
	f, err := os.OpenFile(fileFinal, flag, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	err = writeFile(directory, name, f)
	if err != nil {
		return err
	}

	return nil
}

func writeFile(directory, name string, f *os.File) error {
	fi, err := os.Open(filepath.Join(directory, name))
	if err != nil {
		return err
	}
	defer fi.Close()

	buf := make([]byte, 1024)
	for {
		// read a chunk
		n, err := fi.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		// write a chunk
		if _, err := f.Write(buf[:n]); err != nil {
			return err
		}
	}

	return nil
}
