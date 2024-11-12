package merge

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestSimple1(t *testing.T) {
	rootDir := t.TempDir()
	file := filepath.Join(rootDir, "test01.txt.001")
	d1 := []byte("aaaaaaaa")
	err := os.WriteFile(file, d1, 0644)
	if err != nil {
		t.Fatal(err)
	} else {
		file2 := filepath.Join(rootDir, "test01.txt.002")
		d1 := []byte("bbbbbbb")
		err := os.WriteFile(file2, d1, 0644)
		if err != nil {
			t.Fatal(err)
		} else {
			fileresultat := filepath.Join(rootDir, "test01.txt")
			const RESULTAT = "aaaaaaaabbbbbbb"
			param := MergeParameters{File: file}
			err := MergeCommand(param)
			if err != nil {
				t.Fatal(err)
			} else if _, err := os.Stat(fileresultat); errors.Is(err, os.ErrNotExist) {
				t.Fatal(err)
			} else {
				dat, err := os.ReadFile(fileresultat)
				if err != nil {
					t.Fatal(err)
				} else {
					s := string(dat)
					if s != RESULTAT {
						t.Errorf("error for result : %v != %v", s, RESULTAT)
					}
				}
			}
		}
	}

}

func TestSimple2(t *testing.T) {
	rootDir := t.TempDir()
	file := filepath.Join(rootDir, "test01.txt.001")
	d1 := []byte("aaaaaaaa")
	err := os.WriteFile(file, d1, 0644)
	if err != nil {
		t.Fatal(err)
	} else {
		file2 := filepath.Join(rootDir, "test01.txt.002")
		d1 := []byte("bbbbbbb")
		err := os.WriteFile(file2, d1, 0644)
		if err != nil {
			t.Fatal(err)
		} else {
			fileresultat := filepath.Join(rootDir, "test01_aaa.txt")
			const RESULTAT = "aaaaaaaabbbbbbb"
			param := MergeParameters{File: file, OutputFile: fileresultat}
			err := MergeCommand(param)
			if err != nil {
				t.Fatal(err)
			} else if _, err := os.Stat(fileresultat); errors.Is(err, os.ErrNotExist) {
				t.Fatal(err)
			} else {
				dat, err := os.ReadFile(fileresultat)
				if err != nil {
					t.Fatal(err)
				} else {
					s := string(dat)
					if s != RESULTAT {
						t.Errorf("error for result : %v != %v", s, RESULTAT)
					}
				}
			}
		}
	}

}
