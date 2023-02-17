package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type inputType int

const (
	fileInput inputType = iota
	stringInput
	stdInput
)

type InputParameter struct {
	kind       inputType
	fileSrc    string
	str        string
	file       *os.File
	bufferSize int
}

func FileInputParameter(file string, bufferSize int) (*InputParameter, error) {
	if file == "" {
		return nil, fmt.Errorf("filename emtpy")
	}
	return &InputParameter{fileInput, file, "", nil, bufferSize}, nil
}

func StringInputParameter(str string) (*InputParameter, error) {
	return &InputParameter{stringInput, "", str, nil, 0}, nil
}

func StdInputParameter() (*InputParameter, error) {
	return &InputParameter{stdInput, "", "", nil, 0}, nil
}

func (input *InputParameter) Open() (*bufio.Reader, error) {
	if input.kind == stdInput {
		return bufio.NewReader(os.Stdin), nil
	} else if input.kind == fileInput {
		file, err := os.Open(input.fileSrc)
		if err != nil {
			return nil, err
		}
		input.file = file
		var in *bufio.Reader
		if input.bufferSize > 0 {
			in = bufio.NewReaderSize(file, input.bufferSize)
		} else {
			in = bufio.NewReader(file)
		}
		return in, nil
	} else if input.kind == stringInput {
		in := strings.NewReader(input.str)
		return bufio.NewReader(in), nil
	} else {
		return nil, fmt.Errorf("invalide type of file")
	}
}

func (input InputParameter) Close() error {
	if input.file != nil {
		return input.file.Close()
	}
	return nil
}
