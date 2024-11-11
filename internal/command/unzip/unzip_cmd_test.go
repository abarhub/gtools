package unzip

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"gtools/internal/testutils"
	"io"
	"os"
	"path"
	"reflect"
	"sort"
	"strings"
	"testing"
)

var defaultFiles = map[string][]byte{
	"test1.txt":         {1, 2, 3},
	"test2.txt":         {1, 2, 3},
	"test3.csv":         {1, 2, 3},
	"test4.csv":         {4, 5, 6},
	"test5.log":         {3, 2, 1},
	"dir1/test01.txt":   {7, 8, 9},
	"dir1/test02.txt":   {7, 8, 9},
	"dir1/test03.csv":   {7, 8, 9},
	"dir2/test02_1.txt": {4, 5, 6},
	"dir2/test02_2.csv": {4, 5, 6},
	"dir2/test02_3.txt": {4, 5, 6},
	"dir2/test02_4.log": {4, 5, 6},
}

func Test1(t *testing.T) {
	rootDir := t.TempDir()
	rootDir2 := t.TempDir()
	fichierZip := path.Join(rootDir, "test01.zip")
	err := createZip(fichierZip, defaultFiles)
	if err != nil {
		t.Fatal(err)
	}
	if !t.Failed() {

		param := UnzipParameters{}
		param.Directory = rootDir2
		param.ZipFile = fichierZip
		param.Verbose = false

		err := UnzipCommand(param)
		if err != nil {
			t.Fatal(err)
		} else if _, err := os.Stat(rootDir2); errors.Is(err, os.ErrNotExist) {
			t.Fatal(err)
		} else {
			testutils.CheckFs(t, defaultFiles, rootDir2)
		}
	}
}

func TestExclusion1(t *testing.T) {
	rootDir := t.TempDir()
	rootDir2 := t.TempDir()
	fichierZip := path.Join(rootDir, "test01.zip")
	err := createZip(fichierZip, defaultFiles)
	if err != nil {
		t.Fatal(err)
	}
	if !t.Failed() {

		var resultat = map[string][]byte{
			"test1.txt":         {1, 2, 3},
			"test2.txt":         {1, 2, 3},
			"test5.log":         {3, 2, 1},
			"dir1/test01.txt":   {7, 8, 9},
			"dir1/test02.txt":   {7, 8, 9},
			"dir2/test02_1.txt": {4, 5, 6},
			"dir2/test02_3.txt": {4, 5, 6},
			"dir2/test02_4.log": {4, 5, 6},
		}

		param := UnzipParameters{}
		param.Directory = rootDir2
		param.ZipFile = fichierZip
		param.Verbose = false
		param.ExcludePath = []string{"*.csv"}

		err := UnzipCommand(param)
		if err != nil {
			t.Fatal(err)
		} else if _, err := os.Stat(rootDir2); errors.Is(err, os.ErrNotExist) {
			t.Fatal(err)
		} else {
			testutils.CheckFs(t, resultat, rootDir2)
		}
	}
}

func TestInclusion1(t *testing.T) {
	rootDir := t.TempDir()
	rootDir2 := t.TempDir()
	fichierZip := path.Join(rootDir, "test01.zip")
	err := createZip(fichierZip, defaultFiles)
	if err != nil {
		t.Fatal(err)
	}
	if !t.Failed() {

		var resultat = map[string][]byte{
			"test5.log":         {3, 2, 1},
			"dir2/test02_4.log": {4, 5, 6},
		}

		param := UnzipParameters{}
		param.Directory = rootDir2
		param.ZipFile = fichierZip
		param.Verbose = false
		param.IncludePath = []string{"*.log"}

		err := UnzipCommand(param)
		if err != nil {
			t.Fatal(err)
		} else if _, err := os.Stat(rootDir2); errors.Is(err, os.ErrNotExist) {
			t.Fatal(err)
		} else {
			testutils.CheckFs(t, resultat, rootDir2)
		}
	}
}

func TestVerbose1(t *testing.T) {
	rootDir := t.TempDir()
	rootDir2 := t.TempDir()
	fichierZip := path.Join(rootDir, "test01.zip")
	err := createZip(fichierZip, defaultFiles)
	if err != nil {
		t.Fatal(err)
	}
	if rootDir == rootDir2 {
		t.Fatal(fmt.Errorf("les répertoires sont identiques : %s == %s", rootDir, rootDir2))
	}
	if !t.Failed() {

		param := UnzipParameters{}
		param.Directory = rootDir2
		param.ZipFile = fichierZip
		param.Verbose = true

		out := new(bytes.Buffer)
		err := unzip(param, out)
		if err != nil {
			t.Fatal(err)
		} else if _, err := os.Stat(rootDir2); errors.Is(err, os.ErrNotExist) {
			t.Fatal(err)
		} else {
			testutils.CheckFs(t, defaultFiles, rootDir2)
			if !t.Failed() {
				listeRef := []string{"unzipping file  dir1/test01.txt", "unzipping file  dir1/test02.txt", "unzipping file  dir1/test03.csv", "unzipping file  dir2/test02_1.txt",
					"unzipping file  dir2/test02_2.csv", "unzipping file  dir2/test02_3.txt", "unzipping file  dir2/test02_4.log", "unzipping file  test1.txt", "unzipping file  test2.txt",
					"unzipping file  test3.csv", "unzipping file  test4.csv", "unzipping file  test5.log"}
				liste := strings.Split(out.String(), "\n")
				liste = normaliseListe(liste)
				if !reflect.DeepEqual(liste, listeRef) {
					t.Errorf("error for result : %v != %v", liste, listeRef)
				}
			}
		}
	}
}

func TestVerbose2(t *testing.T) {
	rootDir := t.TempDir()
	rootDir2 := t.TempDir()
	fichierZip := path.Join(rootDir, "test01.zip")
	err := createZip(fichierZip, defaultFiles)
	if err != nil {
		t.Fatal(err)
	}
	if rootDir == rootDir2 {
		t.Fatal(fmt.Errorf("les répertoires sont identiques : %s == %s", rootDir, rootDir2))
	}
	if !t.Failed() {

		param := UnzipParameters{}
		param.Directory = rootDir2
		param.ZipFile = fichierZip
		param.Verbose = false

		out := new(bytes.Buffer)
		err := unzip(param, out)
		if err != nil {
			t.Fatal(err)
		} else if _, err := os.Stat(rootDir2); errors.Is(err, os.ErrNotExist) {
			t.Fatal(err)
		} else {
			testutils.CheckFs(t, defaultFiles, rootDir2)
			if !t.Failed() {
				listeRef := []string{}
				liste := strings.Split(out.String(), "\n")
				liste = normaliseListe(liste)
				if !reflect.DeepEqual(liste, listeRef) {
					t.Errorf("error for result : %v != %v", liste, listeRef)
				}
			}
		}
	}
}

func normaliseListe(liste []string) []string {
	listeResultat := make([]string, 0)
	for _, element := range liste {
		s := strings.Replace(element, "\\", "/", -1)
		if len(s) > 0 {
			listeResultat = append(listeResultat, s)
		}
	}
	sort.Strings(listeResultat)
	return listeResultat
}
func normaliseChemin(s string) string {
	return strings.Replace(s, "\\", "/", -1)
}

func createZip(file string, listeFile map[string][]byte) error {
	archive, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("erreur pour créer le fichier %s : %w", file, err)
	}
	defer archive.Close()
	zipWriter := zip.NewWriter(archive)

	for f2, data := range listeFile {

		w1, err := zipWriter.Create(f2)
		if err != nil {
			return fmt.Errorf("erreur pour créer le fichier %s dans le zip : %w", f2, err)
		}
		reader := bytes.NewReader(data)
		if _, err := io.Copy(w1, reader); err != nil {
			return fmt.Errorf("erreur pour ecrire le fichier %s dans le zip : %w", f2, err)
		}

	}

	err = zipWriter.Close()
	if err != nil {
		return err
	}

	return nil
}
