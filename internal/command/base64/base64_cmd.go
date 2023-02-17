package base64

import (
	"bufio"
	b64 "encoding/base64"
	"fmt"
	"gtools/internal/utils"
	"io"
)

// encoding by buffer. multiple of 4
const bufEncoding = 100

type Base64Parameters struct {
	Input  *utils.InputParameter
	Output *utils.OutputParameter
	Encode bool
}

func EncodeDecodeBase64(param Base64Parameters) error {

	// open input
	in, err := param.Input.Open()
	if err != nil {
		return err
	}
	defer param.Input.Close()

	// open output
	out, err := param.Output.Open()
	if err != nil {
		return err
	}
	defer param.Output.Close()

	if param.Encode {
		err = encode(in, out, bufEncoding)
		if err != nil {
			return err
		}
	} else {
		err = decode(in, out, bufEncoding)
		if err != nil {
			return err
		}
	}

	err = out.Flush()
	if err != nil {
		return err
	}

	return nil
}

func encode(in *bufio.Reader, out *bufio.Writer, nb int) error {
	var buf = make([]byte, nb)
	var buf3 = []byte{}
	for {
		nb2, err := in.Read(buf)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		buf3 = append(buf3, buf[:nb2]...)
		lenBuf := len(buf3)
		if lenBuf >= nb && b64.StdEncoding.EncodedLen(lenBuf) < b64.StdEncoding.EncodedLen(lenBuf+1) {
			n := lenBuf
			for n > 0 {
				if b64.StdEncoding.EncodedLen(n) < b64.StdEncoding.EncodedLen(n+1) {
					break
				}
			}
			if n > 0 {
				var buf2 = make([]byte, b64.StdEncoding.EncodedLen(n))
				b64.StdEncoding.Encode(buf2, buf3[:n])
				_, err = out.Write(buf2)
				if err != nil {
					return err
				}
				var buf4 = buf3
				buf3 = []byte{}
				if n < lenBuf {
					buf3 = append(buf3, buf4[n+1:]...)
				}
			}
		}
	}
	if len(buf3) > 0 {
		var err error
		var buf2 = make([]byte, b64.StdEncoding.EncodedLen(len(buf3)))
		b64.StdEncoding.Encode(buf2, buf3)
		_, err = out.Write(buf2)
		if err != nil {
			return err
		}
	}

	return nil
}

func decode(in *bufio.Reader, out *bufio.Writer, nb int) error {
	var buf = []byte{}
	for {
		c, err := in.ReadByte()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		buf = append(buf, c)
		if len(buf) >= nb {
			var buf2 = make([]byte, b64.StdEncoding.DecodedLen(len(buf)))
			var n int
			n, err = b64.StdEncoding.Decode(buf2, buf)
			if err != nil {
				return fmt.Errorf("error decoding base64: %v", err)
			}
			_, err = out.Write(buf2[:n])
			if err != nil {
				return err
			}
			buf = []byte{}
		}
	}
	if len(buf) > 0 {
		var buf2 = make([]byte, b64.StdEncoding.DecodedLen(len(buf)))
		n, err := b64.StdEncoding.Decode(buf2, buf)
		if err != nil {
			return fmt.Errorf("error decoding base64: %v", err)
		}
		_, err = out.Write(buf2[:n])
		if err != nil {
			return err
		}
	}

	return nil
}
