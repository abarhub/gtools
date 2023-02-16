package base64

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

func TestEncodeDecodeBase64(t *testing.T) {
	type args struct {
		param Base64Parameters
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := EncodeDecodeBase64(tt.args.param); (err != nil) != tt.wantErr {
				t.Errorf("EncodeDecodeBase64() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_decode(t *testing.T) {
	type args struct {
		in  string
		out string
		nb  int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test1", args{"YWJj", "abc", 100}, false},
		{"test2", args{"YWJjYWFhYWRnNGY1ZGJ2YzQxNWI2NGZkZ3M1NmI0NWZnNDVidjRjNWI0eGN2YjR4YzViNGN2eDQxMzI=",
			"abcaaaadg4f5dbvc415b64fdgs56b45fg45bv4c5b4xcvb4xc5b4cvx4132",
			100}, false},
		{"test3", args{
			"YWJjYWFhYWRnNGY1ZGJ2YzQxNWI2NGZkZ3M1NmI0NWZnNDVidjRjNWI0eGN2YjR4YzViNGN2eDQxMzI=",
			"abcaaaadg4f5dbvc415b64fdgs56b45fg45bv4c5b4xcvb4xc5b4cvx4132", 16}, false},
		{"test4", args{"YWJjYWFhYWRnNGY1ZGJ2YzQxNWI2NGZkZ3M1NmI0NWZnNDVidjRjNWI0eGN2YjR4YzViNGN2eDQxMzI=",
			"abcaaaadg4f5dbvc415b64fdgs56b45fg45bv4c5b4xcvb4xc5b4cvx4132",
			1000}, false},
		{"test5", args{"YWJjYWJjYWJj", "abcabcabc", 4}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tt.args.in))
			buf := new(bytes.Buffer)
			writer := bufio.NewWriter(buf)
			err := decode(reader, writer, tt.args.nb)
			err2 := writer.Flush()
			if err2 != nil {
				t.Errorf("decode() error for flush: %v", err2)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("decode() error = %v, wantErr %v", err, tt.wantErr)
			} else if buf.String() != tt.args.out {
				t.Errorf("decode() out = %v, want %v", buf.String(), tt.args.out)
			}
			//if err := decode(tt.args.in, tt.args.out, tt.args.nb); (err != nil) != tt.wantErr {
			//	t.Errorf("decode() error = %v, wantErr %v", err, tt.wantErr)
			//}
		})
	}
}

func Test_encode(t *testing.T) {
	type args struct {
		in  string
		out string
		nb  int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test1", args{"abc", "YWJj", 100}, false},
		{"test2", args{"abcaaaadg4f5dbvc415b64fdgs56b45fg45bv4c5b4xcvb4xc5b4cvx4132",
			"YWJjYWFhYWRnNGY1ZGJ2YzQxNWI2NGZkZ3M1NmI0NWZnNDVidjRjNWI0eGN2YjR4YzViNGN2eDQxMzI=", 100}, false},
		{"test3", args{"abcaaaadg4f5dbvc415b64fdgs56b45fg45bv4c5b4xcvb4xc5b4cvx4132",
			"YWJjYWFhYWRnNGY1ZGJ2YzQxNWI2NGZkZ3M1NmI0NWZnNDVidjRjNWI0eGN2YjR4YzViNGN2eDQxMzI=", 16}, false},
		{"test4", args{"abcabcabc", "YWJjYWJjYWJj", 4}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tt.args.in))
			buf := new(bytes.Buffer)
			writer := bufio.NewWriter(buf)
			err := encode(reader, writer, tt.args.nb)
			err2 := writer.Flush()
			if err2 != nil {
				t.Errorf("encode() error for flush: %v", err2)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("encode() error = %v, wantErr %v", err, tt.wantErr)
			} else if buf.String() != tt.args.out {
				t.Errorf("encode() out = %v, want %v", buf.String(), tt.args.out)
			}
		})
	}
}
