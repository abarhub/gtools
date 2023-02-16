package utils

import (
	"bufio"
	"fmt"
	"os"
)

type outputType int

const (
	fileOutput outputType = iota
	stdOutput
)

type OutputParameter struct {
	kind     outputType
	fileDest string
	file     *os.File
}

func FileOutputParameter(file string) (*OutputParameter, error) {
	if file == "" {
		return nil, fmt.Errorf("filename emtpy")
	}
	return &OutputParameter{fileOutput, file, nil}, nil
}

func StdOutputParameter() (*OutputParameter, error) {
	return &OutputParameter{stdOutput, "", nil}, nil
}

func (output *OutputParameter) Open() (*bufio.Writer, error) {
	if output.kind == stdOutput {
		return bufio.NewWriter(os.Stdout), nil
	} else if output.kind == fileOutput {
		file, err := os.Create(output.fileDest)
		output.file = file
		if err != nil {
			return nil, err
		}
		out := bufio.NewWriter(file)
		return out, nil
	} else {
		return nil, fmt.Errorf("invalide type of file")
	}
}

func (output OutputParameter) Close() error {
	if output.file != nil {
		return output.file.Close()
	}
	return nil
}
