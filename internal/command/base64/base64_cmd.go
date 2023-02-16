package base64

import (
	"bufio"
	b64 "encoding/base64"
	"gtools/internal/utils"
	"io"
)

type Base64Parameters struct {
	Input    *utils.InputParameter
	Output   *utils.OutputParameter
	FileSrc  string
	FileDest string
	Encode   bool
}

func EncodeDecodeBase64(param Base64Parameters) error {
	if param.Encode {
		var buf = []byte{}
		in, err := param.Input.Open()
		if err != nil {
			return err
		}
		defer param.Input.Close()
		buf, err = readBytes(in)
		encodedBase64 := b64.StdEncoding.EncodeToString(buf)

		out, err := param.Output.Open()
		if err != nil {
			return err
		}
		defer param.Output.Close()
		_, err = out.WriteString(encodedBase64)
		if err != nil {
			return err
		}

		err = out.Flush()
		if err != nil {
			return err
		}
		return nil
	} else {
		return nil
	}
}

func readBytes(in *bufio.Reader) ([]byte, error) {
	var buf = []byte{}
	for {
		c, err := in.ReadByte()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		buf = append(buf, c)
	}
	return buf, nil
}
