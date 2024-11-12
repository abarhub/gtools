package password

import (
	"bytes"
	"regexp"
	"strings"
	"testing"
)

func TestGeneratePasswordSimpleLetter1(t *testing.T) {
	s := generatePassword(15, true, false, false)
	if len(s) != 15 {
		t.Errorf("TestGeneratePasswordSimpleLetter1() error = %v, wantErr %v", len(s), 15)
	}
	if !t.Failed() {
		match, _ := regexp.MatchString("^[a-zA-Z]+$", s)
		if !match {
			t.Errorf("TestGeneratePasswordSimpleLetter1() letter error = %v, wantErr %v; s=%s", match, true, s)
		}
		match, _ = regexp.MatchString("^[0-9]+$", s)
		if match {
			t.Errorf("TestGeneratePasswordSimpleLetter1() number error = %v, wantErr %v; s=%s", match, false, s)
		}
		regex := "^[^a-zA-Z0-9]+$"
		match, _ = regexp.MatchString(regex, s)
		if match {
			t.Errorf("TestGeneratePasswordSimpleLetter1() punctuation error = %v, wantErr %v; s=%s", match, false, s)
		}
	}
}

func TestGeneratePasswordSimpleNumber1(t *testing.T) {
	s := generatePassword(15, false, false, true)
	if len(s) != 15 {
		t.Errorf("TestGeneratePasswordSimpleNumber1() error = %v, wantErr %v", len(s), 15)
	}
	if !t.Failed() {
		match, _ := regexp.MatchString("^[a-zA-Z]+$", s)
		if match {
			t.Errorf("TestGeneratePasswordSimpleNumber1() letter error = %v, wantErr %v; s=%s", match, false, s)
		}
		match, _ = regexp.MatchString("^[0-9]+$", s)
		if !match {
			t.Errorf("TestGeneratePasswordSimpleNumber1() number error = %v, wantErr %v; s=%s", match, true, s)
		}
		regex := "^[^a-zA-Z0-9]+$"
		match, _ = regexp.MatchString(regex, s)
		if match {
			t.Errorf("TestGeneratePasswordSimpleNumber1() punctuation error = %v, wantErr %v; s=%s", match, false, s)
		}
	}
}

func TestGeneratePasswordSimplePunctuation1(t *testing.T) {
	s := generatePassword(15, false, true, false)
	if len(s) != 15 {
		t.Errorf("TestGeneratePasswordSimplePunctuation1() error = %v, wantErr %v", len(s), 15)
	}
	if !t.Failed() {
		match, _ := regexp.MatchString("^[a-zA-Z]+$", s)
		if match {
			t.Errorf("TestGeneratePasswordSimplePunctuation1() letter error = %v, wantErr %v; s=%s", match, false, s)
		}
		match, _ = regexp.MatchString("^[0-9]+$", s)
		if match {
			t.Errorf("TestGeneratePasswordSimplePunctuation1() number error = %v, wantErr %v; s=%s", match, false, s)
		}
		regex := "^[^a-zA-Z0-9]+$"
		match, _ = regexp.MatchString(regex, s)
		if !match {
			t.Errorf("TestGeneratePasswordSimplePunctuation1() punctuation error = %v, wantErr %v; s=%s", match, true, s)
		}
	}
}

func TestGeneratePasswordSimpleLetter2(t *testing.T) {
	s := generatePassword(30, true, false, false)
	if len(s) != 30 {
		t.Errorf("TestGeneratePasswordSimpleLetter1() error = %v, wantErr %v", len(s), 15)
	}
	if !t.Failed() {
		match, _ := regexp.MatchString("^[a-zA-Z]+$", s)
		if !match {
			t.Errorf("TestGeneratePasswordSimpleLetter1() letter error = %v, wantErr %v; s=%s", match, true, s)
		}
		match, _ = regexp.MatchString("^[0-9]+$", s)
		if match {
			t.Errorf("TestGeneratePasswordSimpleLetter1() number error = %v, wantErr %v; s=%s", match, false, s)
		}
		regex := "^[^a-zA-Z0-9]+$"
		match, _ = regexp.MatchString(regex, s)
		if match {
			t.Errorf("TestGeneratePasswordSimpleLetter1() punctuation error = %v, wantErr %v; s=%s", match, false, s)
		}
	}
}

func TestPasswordSimple1(t *testing.T) {
	param := PasswordParameters{20, true, false, false, true}
	out := new(bytes.Buffer)
	err := password(param, out)
	if err != nil {
		t.Errorf("TestPasswordSimple1() error = %v, wantErr %v", err, false)
	}
	if !t.Failed() {
		s := out.String()
		s = strings.TrimSuffix(s, "\n")
		if len(s) != 20 {
			t.Errorf("TestPasswordSimple1() error = %v, wantErr %v; s=%s", len(s), 20, s)
		}
		if !t.Failed() {
			match, _ := regexp.MatchString("^[a-zA-Z]+$", s)
			if !match {
				t.Errorf("TestPasswordSimple1() letter error = %v, wantErr %v; s=%s", match, true, s)
			}
			match, _ = regexp.MatchString("^[0-9]+$", s)
			if match {
				t.Errorf("TestPasswordSimple1() number error = %v, wantErr %v; s=%s", match, false, s)
			}
			regex := "^[^a-zA-Z0-9]+$"
			match, _ = regexp.MatchString(regex, s)
			if match {
				t.Errorf("TestPasswordSimple1() punctuation error = %v, wantErr %v; s=%s", match, false, s)
			}
		}
	}
}

func TestPasswordSimple2(t *testing.T) {
	param := PasswordParameters{25, false, true, false, true}
	out := new(bytes.Buffer)
	err := password(param, out)
	if err != nil {
		t.Errorf("TestPasswordSimple1() error = %v, wantErr %v", err, false)
	}
	if !t.Failed() {
		s := out.String()
		s = strings.TrimSuffix(s, "\n")
		if len(s) != 25 {
			t.Errorf("TestPasswordSimple1() error = %v, wantErr %v; s=%s", len(s), 25, s)
		}
		if !t.Failed() {
			match, _ := regexp.MatchString("^[a-zA-Z]+$", s)
			if match {
				t.Errorf("TestPasswordSimple1() letter error = %v, wantErr %v; s=%s", match, false, s)
			}
			match, _ = regexp.MatchString("^[0-9]+$", s)
			if !match {
				t.Errorf("TestPasswordSimple1() number error = %v, wantErr %v; s=%s", match, true, s)
			}
			regex := "^[^a-zA-Z0-9]+$"
			match, _ = regexp.MatchString(regex, s)
			if match {
				t.Errorf("TestPasswordSimple1() punctuation error = %v, wantErr %v; s=%s", match, false, s)
			}
		}
	}
}

func TestPasswordSimple3(t *testing.T) {
	param := PasswordParameters{-1, false, false, true, true}
	out := new(bytes.Buffer)
	err := password(param, out)
	if err != nil {
		t.Errorf("TestPasswordSimple1() error = %v, wantErr %v", err, false)
	}
	if !t.Failed() {
		s := out.String()
		s = strings.TrimSuffix(s, "\n")
		if len(s) != 20 {
			t.Errorf("TestPasswordSimple1() error = %v, wantErr %v; s=%s", len(s), 20, s)
		}
		if !t.Failed() {
			match, _ := regexp.MatchString("^[a-zA-Z]+$", s)
			if match {
				t.Errorf("TestPasswordSimple1() letter error = %v, wantErr %v; s=%s", match, false, s)
			}
			match, _ = regexp.MatchString("^[0-9]+$", s)
			if match {
				t.Errorf("TestPasswordSimple1() number error = %v, wantErr %v; s=%s", match, false, s)
			}
			regex := "^[^a-zA-Z0-9]+$"
			match, _ = regexp.MatchString(regex, s)
			if !match {
				t.Errorf("TestPasswordSimple1() punctuation error = %v, wantErr %v; s=%s", match, true, s)
			}
		}
	}
}

func TestPasswordSimple4(t *testing.T) {
	param := PasswordParameters{18, false, false, true, false}
	out := new(bytes.Buffer)
	err := password(param, out)
	if err != nil {
		t.Errorf("TestPasswordSimple1() error = %v, wantErr %v", err, false)
	}
	if !t.Failed() {
		s := out.String()
		if len(s) != 18 {
			t.Errorf("TestPasswordSimple1() error = %v, wantErr %v; s=%s", len(s), 18, s)
		}
		if !t.Failed() {
			match, _ := regexp.MatchString("^[a-zA-Z]+$", s)
			if match {
				t.Errorf("TestPasswordSimple1() letter error = %v, wantErr %v; s=%s", match, false, s)
			}
			match, _ = regexp.MatchString("^[0-9]+$", s)
			if match {
				t.Errorf("TestPasswordSimple1() number error = %v, wantErr %v; s=%s", match, false, s)
			}
			regex := "^[^a-zA-Z0-9]+$"
			match, _ = regexp.MatchString(regex, s)
			if !match {
				t.Errorf("TestPasswordSimple1() punctuation error = %v, wantErr %v; s=%s", match, true, s)
			}
		}
	}
}
