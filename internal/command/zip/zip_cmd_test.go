package zip

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

func TestSimple1(t *testing.T) {
	rootDir := t.TempDir()
	rootDir2 := t.TempDir()
	dir := rootDir
	testutils.AddFs(t, defaultFiles, dir)
	if !t.Failed() {

		param := ZipParameters{}
		param.Directory = []string{dir}
		zipfile := path.Join(rootDir2, "test1.zip")
		param.ZipFile = zipfile
		param.Verbose = false
		param.Recurvive = true

		err := ZipCommand(param)
		if err != nil {
			t.Fatal(err)
		} else if _, err := os.Stat(zipfile); errors.Is(err, os.ErrNotExist) {
			t.Fatal(err)
		} else {
			resultat, err := readZip(zipfile)
			if err != nil {
				t.Fatal(err)
			} else {
				if !egal(defaultFiles, resultat) {
					t.Errorf("error for result : %v != %v", defaultFiles, resultat)
				}
			}
		}
	}
}

func TestNonRecursif2(t *testing.T) {
	rootDir := t.TempDir()
	rootDir2 := t.TempDir()
	dir := rootDir
	testutils.AddFs(t, defaultFiles, dir)

	var resAttendu = map[string][]byte{
		"test1.txt": {1, 2, 3},
		"test2.txt": {1, 2, 3},
		"test3.csv": {1, 2, 3},
		"test4.csv": {4, 5, 6},
		"test5.log": {3, 2, 1},
	}
	if !t.Failed() {

		param := ZipParameters{}
		param.Directory = []string{dir}
		zipfile := path.Join(rootDir2, "test1.zip")
		param.ZipFile = zipfile
		param.Verbose = false
		param.Recurvive = false

		err := ZipCommand(param)
		if err != nil {
			t.Fatal(err)
		} else if _, err := os.Stat(zipfile); errors.Is(err, os.ErrNotExist) {
			t.Fatal(err)
		} else {
			resultat, err := readZip(zipfile)
			if err != nil {
				t.Fatal(err)
			} else {

				if !egal(resAttendu, resultat) {
					t.Errorf("error for result : %v != %v", resAttendu, resultat)
				}
			}
		}
	}
}

func TestExclusion1(t *testing.T) {
	rootDir := t.TempDir()
	rootDir2 := t.TempDir()
	dir := rootDir
	testutils.AddFs(t, defaultFiles, dir)

	var resAttendu = map[string][]byte{
		"test1.txt":         {1, 2, 3},
		"test2.txt":         {1, 2, 3},
		"test5.log":         {3, 2, 1},
		"dir1/test01.txt":   {7, 8, 9},
		"dir1/test02.txt":   {7, 8, 9},
		"dir2/test02_1.txt": {4, 5, 6},
		"dir2/test02_3.txt": {4, 5, 6},
		"dir2/test02_4.log": {4, 5, 6},
	}
	if !t.Failed() {

		param := ZipParameters{}
		param.Directory = []string{dir}
		zipfile := path.Join(rootDir2, "test1.zip")
		param.ZipFile = zipfile
		param.Verbose = false
		param.Recurvive = true
		param.ExcludePath = []string{"*.csv"}

		err := ZipCommand(param)
		if err != nil {
			t.Fatal(err)
		} else if _, err := os.Stat(zipfile); errors.Is(err, os.ErrNotExist) {
			t.Fatal(err)
		} else {
			resultat, err := readZip(zipfile)
			if err != nil {
				t.Fatal(err)
			} else {
				if !egal(resAttendu, resultat) {
					t.Errorf("error for result : %v != %v", resAttendu, resultat)
				}
			}
		}
	}
}

func TestInclusion1(t *testing.T) {
	rootDir := t.TempDir()
	rootDir2 := t.TempDir()
	dir := rootDir
	testutils.AddFs(t, defaultFiles, dir)

	var resAttendu = map[string][]byte{
		"test1.txt":         {1, 2, 3},
		"test2.txt":         {1, 2, 3},
		"dir1/test01.txt":   {7, 8, 9},
		"dir1/test02.txt":   {7, 8, 9},
		"dir2/test02_1.txt": {4, 5, 6},
		"dir2/test02_3.txt": {4, 5, 6},
	}
	if !t.Failed() {

		param := ZipParameters{}
		param.Directory = []string{dir}
		zipfile := path.Join(rootDir2, "test1.zip")
		param.ZipFile = zipfile
		param.Verbose = false
		param.Recurvive = true
		param.IncludePath = []string{"*.txt"}

		err := ZipCommand(param)
		if err != nil {
			t.Fatal(err)
		} else if _, err := os.Stat(zipfile); errors.Is(err, os.ErrNotExist) {
			t.Fatal(err)
		} else {
			resultat, err := readZip(zipfile)
			if err != nil {
				t.Fatal(err)
			} else {
				if !egal(resAttendu, resultat) {
					t.Errorf("error for result : %v != %v", resAttendu, resultat)
				}
			}
		}
	}
}

func TestVerbose1(t *testing.T) {
	rootDir := t.TempDir()
	rootDir2 := t.TempDir()
	dir := rootDir
	testutils.AddFs(t, defaultFiles, dir)
	if !t.Failed() {

		param := ZipParameters{}
		param.Directory = []string{dir}
		zipfile := path.Join(rootDir2, "test1.zip")
		param.ZipFile = zipfile
		param.Verbose = true
		param.Recurvive = true

		out := new(bytes.Buffer)
		err := createZip(param, out)
		if err != nil {
			t.Fatal(err)
		} else if _, err := os.Stat(zipfile); errors.Is(err, os.ErrNotExist) {
			t.Fatal(err)
		} else {
			resultat, err := readZip(zipfile)
			if err != nil {
				t.Fatal(err)
			} else {
				if !egal(defaultFiles, resultat) {
					t.Errorf("error for result : %v != %v", defaultFiles, resultat)
				}
				listeRef := []string{"create dir1/test01.txt", "create dir1/test02.txt", "create dir1/test03.csv",
					"create dir2/test02_1.txt", "create dir2/test02_2.csv", "create dir2/test02_3.txt",
					"create dir2/test02_4.log", "create test1.txt", "create test2.txt", "create test3.csv",
					"create test4.csv", "create test5.log"}
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
	dir := rootDir
	testutils.AddFs(t, defaultFiles, dir)
	if !t.Failed() {

		param := ZipParameters{}
		param.Directory = []string{dir}
		zipfile := path.Join(rootDir2, "test1.zip")
		param.ZipFile = zipfile
		param.Verbose = false
		param.Recurvive = true

		out := new(bytes.Buffer)
		err := createZip(param, out)
		if err != nil {
			t.Fatal(err)
		} else if _, err := os.Stat(zipfile); errors.Is(err, os.ErrNotExist) {
			t.Fatal(err)
		} else {
			resultat, err := readZip(zipfile)
			if err != nil {
				t.Fatal(err)
			} else {
				if !egal(defaultFiles, resultat) {
					t.Errorf("error for result : %v != %v", defaultFiles, resultat)
				}
				liste := strings.Split(out.String(), "\n")
				liste = normaliseListe(liste)
				if len(liste) > 0 {
					t.Errorf("error for result : %v len>0", liste)
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

func egal(map1 map[string][]byte, map2 map[string][]byte) bool {
	var map01, map02 map[string][]byte
	map01 = map[string][]byte{}
	map02 = map[string][]byte{}

	for k, v := range map1 {
		k = strings.Replace(k, "\\", "/", -1)
		map01[k] = v
	}

	for k, v := range map2 {
		k = strings.Replace(k, "\\", "/", -1)
		map02[k] = v
	}
	return reflect.DeepEqual(map01, map02)
}

func readZip(fileZip string) (map[string][]byte, error) {
	res := map[string][]byte{}
	r, err := zip.OpenReader(fileZip)

	if err != nil {
		return nil, fmt.Errorf("error for open file %s : %w", fileZip, err)
	}
	defer r.Close()

	for _, f := range r.File {

		rc, err := f.Open()

		if err != nil {
			return nil, err
		}
		var b bytes.Buffer
		writer := io.Writer(&b)
		_, err = io.Copy(writer, rc)
		if err != nil {
			return nil, fmt.Errorf("error for read file %s in zip : %w", f.Name, err)
		}
		err = rc.Close()
		if err != nil {
			return nil, fmt.Errorf("error for close file %s in zip : %w", f.Name, err)
		}
		res[f.Name] = b.Bytes()
	}

	return res, nil
}
