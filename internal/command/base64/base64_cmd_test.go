package base64

import (
	"bufio"
	"bytes"
	"gtools/internal/utils"
	"os"
	"path"
	"strings"
	"testing"
)

func TestEncodeDecodeBase64(t *testing.T) {
	type args struct {
		input, output string
		encode        bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test1", args{
			"abcdef", "YWJjZGVm", true,
		}, false},
		{"test2", args{
			"YWJjZGVm", "abcdef", false,
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rootDir := t.TempDir()
			fileInputPath := path.Join(rootDir, "input.txt")
			fileOutputPath := path.Join(rootDir, "output.txt")

			err3 := os.WriteFile(fileInputPath, []byte(tt.args.input), 0755)
			if err3 != nil {
				t.Errorf("EncodeDecodeBase64() error = %v", err3)
			} else {
				input, err := utils.FileInputParameter(fileInputPath, 0)
				if err != nil {
					t.Errorf("EncodeDecodeBase64() error = %v", err)
				} else {
					output, err := utils.FileOutputParameter(fileOutputPath, 0)
					if err != nil {
						t.Errorf("EncodeDecodeBase64() error = %v", err)
					} else {
						param := Base64Parameters{input, output, tt.args.encode, 0}
						if err := EncodeDecodeBase64(param); (err != nil) != tt.wantErr {
							t.Errorf("EncodeDecodeBase64() error = %v, wantErr %v", err, tt.wantErr)
						} else {
							buf, err := os.ReadFile(fileOutputPath)
							if err != nil {
								t.Errorf("EncodeDecodeBase64() error = %v", err)
							} else if string(buf[:]) != tt.args.output {
								t.Errorf("EncodeDecodeBase64() output = %v, want %v", string(buf[:]), tt.args.output)
							}
						}
					}
				}
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
		{"test6", args{"YWJjYWJjYWJj", "abcabcabc", 5}, false},
		{"test7", args{"YWJjYWJjYWJj", "abcabcabc", 6}, false},
		{"test8", args{"YWJjYWJjYWJj", "abcabcabc", 7}, false},
		{"test9", args{"YWJjYWJjYWJj", "abcabcabc", 8}, false},
		{"test10", args{"YWJjYWJjYWJj", "abcabcabc", 9}, false},
		{"test11", args{"YWJjYWJjYWJj", "abcabcabc", 1}, false},
		{"test12", args{"YWJjYWJjYWJj", "abcabcabc", 2}, false},
		{"test13", args{"YWJjYWJjYWJj", "abcabcabc", 3}, false},
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
		{"test5", args{"abcabcabc", "YWJjYWJjYWJj", 5}, false},
		{"test6", args{"abcabcabc", "YWJjYWJjYWJj", 6}, false},
		{"test7", args{"abcabcabc", "YWJjYWJjYWJj", 7}, false},
		{"test8", args{"abcabcabc", "YWJjYWJjYWJj", 8}, false},
		{"test9", args{"abcabcabc", "YWJjYWJjYWJj", 9}, false},
		{"test10", args{"abcabcabc", "YWJjYWJjYWJj", 1}, false},
		{"test11", args{"abcabcabc", "YWJjYWJjYWJj", 2}, false},
		{"test12", args{"abcabcabc", "YWJjYWJjYWJj", 3}, false},
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

func TestEncodeDecodeBase64EncodeDecode(t *testing.T) {
	type args struct {
		input, output string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test1", args{"abcdef", "YWJjZGVm"}, false},
		{"test2", args{"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "YWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFh"}, false},
		{"test3", args{"1234567890", "MTIzNDU2Nzg5MA=="}, false},
		{"test4", args{"abcdefghijklmnopqrstuvwxyz", "YWJjZGVmZ2hpamtsbW5vcHFyc3R1dnd4eXo="}, false},
		{"test5", args{"ABCDEFGHIJKLMNOPQRSTUVWXYZ", "QUJDREVGR0hJSktMTU5PUFFSU1RVVldYWVo="}, false},
		{"test6", args{",;:!*$=)@_-('\"'&~#{[|`\\^@]}<>?./§%µ¤ ", "LDs6ISokPSlAXy0oJyInJn4je1t8YFxeQF19PD4/Li/CpyXCtcKkIA=="}, false},
		{"test7", args{"àèéïöüÿûôñç_ÀÈÉÏÖÜŸÛÔÑÇ", "w6DDqMOpw6/DtsO8w7/Du8O0w7HDp1/DgMOIw4nDj8OWw5zFuMObw5TDkcOH"}, false},
		{"test8", args{"\u0000\u0001\u0002\u0003", "AAECAw=="}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rootDir := t.TempDir()
			fileInputPath := path.Join(rootDir, "input.txt")
			fileOutputPath := path.Join(rootDir, "output.txt")

			err3 := os.WriteFile(fileInputPath, []byte(tt.args.input), 0755)
			if err3 != nil {
				t.Errorf("EncodeDecodeBase64() error = %v", err3)
			} else {
				input, err := utils.FileInputParameter(fileInputPath, 0)
				if err != nil {
					t.Errorf("EncodeDecodeBase64() error = %v", err)
				} else {
					output, err := utils.FileOutputParameter(fileOutputPath, 0)
					if err != nil {
						t.Errorf("EncodeDecodeBase64() error = %v", err)
					} else {
						param := Base64Parameters{input, output, true, 0}
						if err := EncodeDecodeBase64(param); (err != nil) != tt.wantErr {
							t.Errorf("EncodeDecodeBase64() error = %v, wantErr %v", err, tt.wantErr)
						} else {
							buf, err := os.ReadFile(fileOutputPath)
							if err != nil {
								t.Errorf("EncodeDecodeBase64() error = %v", err)
							} else if string(buf[:]) != tt.args.output {
								t.Errorf("EncodeDecodeBase64() output = %v, want %v", string(buf[:]), tt.args.output)
							} else {
								fileInputPath2 := path.Join(rootDir, "input2.txt")
								fileOutputPath2 := path.Join(rootDir, "output2.txt")

								err3 := os.WriteFile(fileInputPath2, buf, 0755)
								if err3 != nil {
									t.Errorf("EncodeDecodeBase64() error = %v", err3)
								} else {
									input, err := utils.FileInputParameter(fileInputPath2, 0)
									if err != nil {
										t.Errorf("EncodeDecodeBase64() error = %v", err)
									} else {
										output, err := utils.FileOutputParameter(fileOutputPath2, 0)
										if err != nil {
											t.Errorf("EncodeDecodeBase64() error = %v", err)
										} else {
											param := Base64Parameters{input, output, false, 0}
											if err := EncodeDecodeBase64(param); (err != nil) != tt.wantErr {
												t.Errorf("EncodeDecodeBase64() error = %v, wantErr %v", err, tt.wantErr)
											} else {
												buf, err := os.ReadFile(fileOutputPath2)
												if err != nil {
													t.Errorf("EncodeDecodeBase64() error = %v", err)
												} else if string(buf[:]) != tt.args.input {
													t.Errorf("encode() out = %v, want %v", string(buf), tt.args.input)
												}
											}
										}
									}
								}

							}
						}
					}
				}
			}
		})
	}
}

func benchmarkEncodeDecodeBase64(b *testing.B, nb int, bufferSize int) {
	rootDir := b.TempDir()
	fileInputPath := path.Join(rootDir, "input.txt")
	fileOutputPath := path.Join(rootDir, "output.txt")

	var input []byte
	input = []byte{}
	for i := 0; i < nb; i++ {
		input = append(input, 'a')
	}

	err3 := os.WriteFile(fileInputPath, input, 0755)
	if err3 != nil {
		b.Errorf("EncodeDecodeBase64() error = %v", err3)
	} else {
		input, err := utils.FileInputParameter(fileInputPath, bufferSize)
		if err != nil {
			b.Errorf("EncodeDecodeBase64() error = %v", err)
		} else {
			output, err := utils.FileOutputParameter(fileOutputPath, bufferSize)
			if err != nil {
				b.Errorf("EncodeDecodeBase64() error = %v", err)
			} else {
				param := Base64Parameters{input, output, true, bufferSize}
				if err := EncodeDecodeBase64(param); (err != nil) != false {
					b.Errorf("EncodeDecodeBase64() error = %v, wantErr %v", err, false)
				} else {
					_, err := os.ReadFile(fileOutputPath)
					if err != nil {
						b.Errorf("EncodeDecodeBase64() error = %v", err)
					}
				}
			}
		}
	}
}

func BenchmarkEncodeDecodeBase64_100(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkEncodeDecodeBase64(b, 10_000_000, 100)
	}
}

func BenchmarkEncodeDecodeBase64_1000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkEncodeDecodeBase64(b, 10_000_000, 1000)
	}
}

func BenchmarkEncodeDecodeBase64_10000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkEncodeDecodeBase64(b, 10_000_000, 10000)
	}
}

func BenchmarkEncodeDecodeBase64_100000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkEncodeDecodeBase64(b, 10_000_000, 100000)
	}
}
