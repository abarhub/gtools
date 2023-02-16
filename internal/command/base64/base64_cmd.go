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
	if param.Encode {
		//var buf = []byte{}

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

		//if false {
		//	// read input
		//	buf, err = readBytes(in)
		//
		//	// encode
		//	encodedBase64 := b64.StdEncoding.EncodeToString(buf)
		//
		//	// write output
		//	_, err = out.WriteString(encodedBase64)
		//	if err != nil {
		//		return err
		//	}
		//} else {
		err = encode(in, out, bufEncoding)
		if err != nil {
			return err
		}
		//}

		err = out.Flush()
		if err != nil {
			return err
		}
		return nil
	} else {

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

		err = decode(in, out, bufEncoding)
		if err != nil {
			return err
		}

		err = out.Flush()
		if err != nil {
			return err
		}

		return nil
	}
}

//func readBytes(in *bufio.Reader) ([]byte, error) {
//	var buf = []byte{}
//	for {
//		c, err := in.ReadByte()
//		if err == io.EOF {
//			break
//		} else if err != nil {
//			return nil, err
//		}
//		buf = append(buf, c)
//	}
//	return buf, nil
//}

func encode(in *bufio.Reader, out *bufio.Writer, nb int) error {
	var buf = []byte{}
	methode1 := true
	no := 0
	for {
		c, err := in.ReadByte()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		buf = append(buf, c)
		no++
		lenBuf := len(buf)
		if lenBuf >= nb && b64.StdEncoding.EncodedLen(lenBuf) < b64.StdEncoding.EncodedLen(lenBuf+1) {
			if methode1 {
				var buf2 = make([]byte, b64.StdEncoding.EncodedLen(len(buf)))
				b64.StdEncoding.Encode(buf2, buf)
				//s := b64.StdEncoding.EncodeToString(buf)
				//_, err = out.WriteString(s)
				_, err = out.Write(buf2)
			} else {
				s := b64.StdEncoding.EncodeToString(buf)
				_, err = out.WriteString(s)
			}
			if err != nil {
				return err
			}
			buf = []byte{}
			no = 0
		}
	}
	if len(buf) > 0 {
		var err error
		if methode1 {
			var buf2 = make([]byte, b64.StdEncoding.EncodedLen(len(buf)))
			b64.StdEncoding.Encode(buf2, buf)
			_, err = out.Write(buf2)
		} else {
			s := b64.StdEncoding.EncodeToString(buf)
			_, err = out.WriteString(s)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func decode(in *bufio.Reader, out *bufio.Writer, nb int) error {
	var buf = []byte{}
	no := 0
	for {
		c, err := in.ReadByte()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		buf = append(buf, c)
		no++
		//lenBuf:=len(buf)
		//if no+1 > nb {
		//if lenBuf>=nb &&b64.StdEncoding.DecodedLen(lenBuf)<b64.StdEncoding.DecodedLen(lenBuf+1){
		if len(buf) >= nb {
			if true {
				var buf2 = make([]byte, b64.StdEncoding.DecodedLen(len(buf)))
				_, err = b64.StdEncoding.Decode(buf2, buf)
				if err != nil {
					return fmt.Errorf("error decoding base64: %v", err)
				}
				//s := b64.StdEncoding.EncodeToString(buf)
				//_, err = out.WriteString(s)
				_, err = out.Write(buf2)
			} else {
				s := b64.StdEncoding.EncodeToString(buf)
				_, err = out.WriteString(s)
			}
			if err != nil {
				return err
			}
			buf = []byte{}
			no = 0
		}
	}
	if len(buf) > 0 {
		var buf2 = make([]byte, b64.StdEncoding.DecodedLen(len(buf)))
		_, err := b64.StdEncoding.Decode(buf2, buf)
		if err != nil {
			return fmt.Errorf("error decoding base64: %v", err)
		}
		_, err = out.Write(buf2)
		if err != nil {
			return err
		}
	}

	return nil
}
