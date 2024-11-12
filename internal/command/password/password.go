package password

import (
	"fmt"
	"io"
	"math/rand"
	"os"
)

type PasswordParameters struct {
	Size     int
	Number   bool
	Letter   bool
	Punctuer bool
	NewLine  bool
}

const (
	letterBytes  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specialBytes = "!@#$%^&*()_+-=[]{}\\|;':\",.<>/?`~"
	numBytes     = "0123456789"
)

func PasswordCommand(param PasswordParameters) error {
	return password(param, os.Stdout)
}

func password(param PasswordParameters, out io.Writer) error {

	var size int
	if param.Size > 0 {
		size = param.Size
	} else {
		size = 20
	}
	password := generatePassword(size, param.Number, param.Punctuer, param.Letter)
	newLine := ""
	if param.NewLine {
		newLine = "\n"
	}
	_, err := fmt.Fprintf(out, "%s%s", password, newLine)
	if err != nil {
		return err
	}

	return nil
}

func generatePassword(length int, useLetters bool, useSpecial bool, useNum bool) string {
	b := make([]byte, length)
	for i := range b {
		if useLetters {
			b[i] = letterBytes[rand.Intn(len(letterBytes))]
		} else if useSpecial {
			b[i] = specialBytes[rand.Intn(len(specialBytes))]
		} else if useNum {
			b[i] = numBytes[rand.Intn(len(numBytes))]
		}
	}
	return string(b)
}
