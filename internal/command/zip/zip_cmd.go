package zip

import (
	"archive/zip"
	"fmt"
	"gtools/internal/utils"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type ZipParameters struct {
	ZipFile     string
	Directory   []string
	Recurvive   bool
	ExcludePath []string
	IncludePath []string
	Verbose     bool
}

func ZipCommand(param ZipParameters) error {

	err := createZip(param)

	return err
}

func createZip(param ZipParameters) error {
	archive, err := os.Create(param.ZipFile)
	if err != nil {
		return fmt.Errorf("error for create file %s : %w", param.ZipFile, err)
	}
	defer archive.Close()

	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()

	for _, dir := range param.Directory {
		err = listFiles(zipWriter, dir, param, "", os.Stdout)
		if err != nil {
			return err
		}
	}

	return nil
}

func listFiles(archive *zip.Writer, path string, param ZipParameters, rep string, out io.Writer) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, file := range files {
		filename := filepath.Join(path, file.Name())

		skip := false
		if strings.HasPrefix(file.Name(), ".") {
			skip = true
		}

		if !skip {
			toScan, err := fileToScan(filename, param, true)
			if err != nil {
				return err
			} else if toScan {
				if file.IsDir() {
					if param.Recurvive {
						err := listFiles(archive, filename, param, filepath.Join(rep, file.Name()), out)
						if err != nil {
							return err
						}
					}
				} else {
					toScan, err := fileToScan(filename, param, false)
					if err != nil {
						return err
					} else if toScan {
						filename := filepath.Join(rep, file.Name())
						if param.Verbose {
							fmt.Fprintf(out, "create %s\n", filename)
						}
						err = zipFile(archive, filepath.Join(path, file.Name()), filename)
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

func zipFile(archive *zip.Writer, file string, pathDest string) error {
	f1, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("error for open file %s : %w", file, err)
	}
	defer f1.Close()

	w1, err := archive.Create(pathDest)
	if err != nil {
		return fmt.Errorf("error for create file %s in zip : %w", pathDest, err)
	}
	if _, err := io.Copy(w1, f1); err != nil {
		return fmt.Errorf("error for read file %s : %w", file, err)
	}
	return nil
}

func fileToScan(file string, param ZipParameters, exclude bool) (bool, error) {
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
