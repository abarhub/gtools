package split

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"path"
	"reflect"
	"testing"
)

func TestSimple1(t *testing.T) {
	rootDir := t.TempDir()

	fichier := path.Join(rootDir, "fichier.txt")
	err := creerFichier(fichier, 8000)
	if err != nil {
		t.Fatal(err)
	}

	param := SplitParameters{}
	param.File = fichier
	param.SizeStr = "4000"

	err = SplitCommand(param)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(fichier + ".001"); errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
	} else if _, err := os.Stat(fichier + ".002"); errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
	} else if _, err := os.Stat(fichier + ".003"); !errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
	} else {
		f1, err := lectureFichier(fichier)
		if err != nil {
			t.Fatal(err)
		}
		var f2 = []byte{}
		if !t.Failed() {
			f3, err := lectureFichier(fichier + ".001")
			if err != nil {
				t.Fatal(err)
			}
			f2 = append(f2, f3...)
		}
		if !t.Failed() {
			f3, err := lectureFichier(fichier + ".002")
			if err != nil {
				t.Fatal(err)
			}
			f2 = append(f2, f3...)
		}
		if !t.Failed() {
			if !reflect.DeepEqual(f1, f2) {
				t.Errorf("error for result : %v != %v", f1, f2)
			}
		}
	}
}

func TestSimple2(t *testing.T) {
	rootDir := t.TempDir()

	fichier := path.Join(rootDir, "fichier.txt")
	err := creerFichier(fichier, 6000)
	if err != nil {
		t.Fatal(err)
	}

	param := SplitParameters{}
	param.File = fichier
	param.SizeStr = "4000"

	err = SplitCommand(param)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(fichier + ".001"); errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
	} else if _, err := os.Stat(fichier + ".002"); errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
	} else if _, err := os.Stat(fichier + ".003"); !errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
	} else {
		f1, err := lectureFichier(fichier)
		if err != nil {
			t.Fatal(err)
		}
		var f2 = []byte{}
		if !t.Failed() {
			f3, err := lectureFichier(fichier + ".001")
			if err != nil {
				t.Fatal(err)
			}
			f2 = append(f2, f3...)
		}
		if !t.Failed() {
			f3, err := lectureFichier(fichier + ".002")
			if err != nil {
				t.Fatal(err)
			}
			f2 = append(f2, f3...)
		}
		if !t.Failed() {
			if !reflect.DeepEqual(f1, f2) {
				t.Errorf("error for result : %v != %v", f1, f2)
			}
		}
	}
}

func TestSimple3(t *testing.T) {
	rootDir := t.TempDir()

	fichier := path.Join(rootDir, "fichier.txt")
	err := creerFichier(fichier, 8000)
	if err != nil {
		t.Fatal(err)
	}

	param := SplitParameters{}
	param.File = fichier
	param.SizeStr = "10000"

	err = SplitCommand(param)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(fichier + ".001"); errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
	} else if _, err := os.Stat(fichier + ".002"); !errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
	} else if _, err := os.Stat(fichier + ".003"); !errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
	} else {
		f1, err := lectureFichier(fichier)
		if err != nil {
			t.Fatal(err)
		}
		var f2 = []byte{}
		if !t.Failed() {
			f3, err := lectureFichier(fichier + ".001")
			if err != nil {
				t.Fatal(err)
			}
			f2 = append(f2, f3...)
		}
		if !t.Failed() {
			if !reflect.DeepEqual(f1, f2) {
				t.Errorf("error for result : %v != %v", f1, f2)
			}
		}
	}
}

func creerFichier(fichier string, size int64) error {
	destination, err := os.Create(fichier)
	if err != nil {
		return err
	}
	defer destination.Close()

	var c byte = 0
	var data = make([]byte, size)
	for n := range size {

		data[n] = c
		var c3 = uint16(c + 1)
		var c2 = c3 % 256
		c = byte(c2)
	}
	reader := bytes.NewReader(data)
	_, err = io.Copy(destination, reader)
	if err != nil {
		return err
	}
	err = destination.Close()
	if err != nil {
		return err
	}
	return nil
}

func lectureFichier(nomFichier string) ([]byte, error) {
	file, err := os.Open(nomFichier)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Get the file size
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	// Read the file into a byte slice
	bs := make([]byte, stat.Size())
	_, err = bufio.NewReader(file).Read(bs)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return bs, nil
}
