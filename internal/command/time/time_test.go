package time

import (
	"bytes"
	"runtime"
	"strings"
	"testing"
)

func Test(t *testing.T) {
	var param TimeParameters
	os := runtime.GOOS
	if os == "windows" {
		param = TimeParameters{"cmd", []string{"/c", "echo", "Test"}}
	} else {
		param = TimeParameters{"/bin/sh", []string{"-c", "echo Test"}}
	}

	out := &bytes.Buffer{}
	err := timeCommande(param, out)
	if (err != nil) != false {
		t.Errorf("timeCommande() error = %v, wantErr %v", err, false)
		return
	}
	s := "Exec took "
	if gotOut := out.String(); !strings.HasPrefix(gotOut, s) {
		t.Errorf("timeCommande() gotOut = %v, want %v", gotOut, s)
	}
}
