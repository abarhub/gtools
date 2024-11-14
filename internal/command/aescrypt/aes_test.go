package aescrypt

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestSimple1(t *testing.T) {
	rootDir := t.TempDir()
	fileInput := filepath.Join(rootDir, "test01.txt")
	d1 := []byte("aaaaaaaa")
	err := os.WriteFile(fileInput, d1, 0644)
	if err != nil {
		t.Fatal(err)
	} else {

		key, err := createPassword()
		if err != nil {
			t.Fatal(err)
		} else if len(key) != 32 {
			t.Fatalf("wrong key length: %d", len(key))
		} else {
			fileOutput := filepath.Join(rootDir, "test01.txt.aes")

			err = encrypt(key, fileInput, fileOutput)
			if err != nil {
				t.Fatal(err)
			} else if _, err := os.Stat(fileOutput); errors.Is(err, os.ErrNotExist) {
				t.Fatal(err)
			} else {
				fileOutput2 := filepath.Join(rootDir, "test02.txt")
				err = decrypt(key, fileOutput, fileOutput2)
				if err != nil {
					t.Fatal(err)
				} else {
					buf, err := os.ReadFile(fileOutput2)
					if err != nil {
						t.Fatal(err)
					} else if !reflect.DeepEqual(buf, d1) {
						t.Fatalf("TestSimple1 : %s != %s", string(buf), string(d1))
					}
					buf2, err := os.ReadFile(fileOutput)
					if err != nil {
						t.Fatal(err)
					} else if reflect.DeepEqual(buf2, d1) {
						t.Fatalf("TestSimple1 : %s == %s", string(buf2), string(d1))
					}
				}
			}
		}
	}
}

func TestSimple2(t *testing.T) {
	rootDir := t.TempDir()
	fileInput := filepath.Join(rootDir, "test01.txt")
	d1 := []byte("aaaaaaaa")
	for i := range 5000 {
		d1 = append(d1, byte(i))
	}
	err := os.WriteFile(fileInput, d1, 0644)
	if err != nil {
		t.Fatal(err)
	} else {

		key, err := createPassword()
		if err != nil {
			t.Fatal(err)
		} else if len(key) != 32 {
			t.Fatalf("wrong key length: %d", len(key))
		} else {
			fileOutput := filepath.Join(rootDir, "test01.txt.aes")

			err = encrypt(key, fileInput, fileOutput)
			if err != nil {
				t.Fatal(err)
			} else if _, err := os.Stat(fileOutput); errors.Is(err, os.ErrNotExist) {
				t.Fatal(err)
			} else {
				fileOutput2 := filepath.Join(rootDir, "test02.txt")
				err = decrypt(key, fileOutput, fileOutput2)
				if err != nil {
					t.Fatal(err)
				} else {
					buf, err := os.ReadFile(fileOutput2)
					if err != nil {
						t.Fatal(err)
					} else if !reflect.DeepEqual(buf, d1) {
						t.Fatalf("TestSimple1 : %s != %s", string(buf), string(d1))
					}
					buf2, err := os.ReadFile(fileOutput)
					if err != nil {
						t.Fatal(err)
					} else if reflect.DeepEqual(buf2, d1) {
						t.Fatalf("TestSimple1 : %s == %s", string(buf2), string(d1))
					}
				}
			}
		}
	}
}
