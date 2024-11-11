package unzip

import (
	"archive/zip"
	"fmt"
	"gtools/internal/utils"
	"io"
	"os"
	"path/filepath"
)

type UnzipParameters struct {
	ZipFile     string
	Directory   string
	ExcludePath []string
	IncludePath []string
	Verbose     bool
}

func UnzipCommand(param UnzipParameters) error {

	err := unzip(param)

	return err
}

func unzip(param UnzipParameters) error {
	archive, err := zip.OpenReader(param.ZipFile)
	if err != nil {
		return fmt.Errorf("erreur pour lire le fichier %s : %w", param.ZipFile, err)
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(param.Directory, f.Name)
		toScan, err := fileToScan(filePath, param, true)
		if err != nil {
			return err
		}
		if toScan {
			if f.FileInfo().IsDir() {

				if param.Verbose {
					fmt.Println("unzipping file ", filePath)
				}
				if param.Verbose {
					fmt.Println("creating directory...")
				}
				err = os.MkdirAll(filePath, os.ModePerm)
				if err != nil {
					return fmt.Errorf("impossible de creer le répertoire %s : %w", filePath, err)
				}

			} else {
				toScan, err := fileToScan(filePath, param, false)
				if err != nil {
					return err
				}
				if toScan {
					if param.Verbose {
						fmt.Println("unzipping file ", filePath)
					}
					if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
						return fmt.Errorf("impossible de creer le répertoire %s : %w", filepath.Dir(filePath), err)
					}

					err = copieFichier(filePath, f)
				}
			}
		}
	}

	return nil
}

func copieFichier(filePath string, f *zip.File) error {
	dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return fmt.Errorf("impossible de creer le fichier %s : %w", filePath, err)
	}
	defer dstFile.Close()

	fileInArchive, err := f.Open()
	if err != nil {
		return fmt.Errorf("impossible d'ouvrir' le fichier %s dans le zip : %w", f.Name, err)
	}
	defer fileInArchive.Close()

	if _, err := io.Copy(dstFile, fileInArchive); err != nil {
		return fmt.Errorf("impossible de copier le fichier %s : %w", f.Name, err)
	}

	return nil
}

func fileToScan(file string, param UnzipParameters, exclude bool) (bool, error) {
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
