package commandline

import (
	"fmt"
	"github.com/spf13/cobra"
	"gtools/internal/command/base64"
	"gtools/internal/utils"
	"strconv"
)

var (
	fileSrc    string
	fileDest   string
	encode     bool
	bufferSize string
)

var base64Cmd = &cobra.Command{
	Use:   "base64",
	Short: "encode/decode in base64",
	RunE: func(cmd *cobra.Command, args []string) error {
		var input *utils.InputParameter
		var output *utils.OutputParameter
		var err error
		var bufferSizeInt int
		if len(bufferSize) > 0 {
			if n, err := strconv.Atoi(bufferSize); err == nil {
				bufferSizeInt = n
			} else if err != nil {
				return fmt.Errorf("bufferSize is not a number (%v)", bufferSize)
			}
		}
		if fileSrc != "" {
			input, err = utils.FileInputParameter(fileSrc, bufferSizeInt)
			if err != nil {
				return err
			}
		} else {
			input, err = utils.StdInputParameter()
			if err != nil {
				return err
			}
		}
		if fileDest != "" {
			output, err = utils.FileOutputParameter(fileDest, bufferSizeInt)
			if err != nil {
				return err
			}
		} else {
			output, err = utils.StdOutputParameter()
		}
		param := base64.Base64Parameters{input, output, encode, bufferSizeInt}
		commandError = cmdBase64(param)
		return nil
	},
}

func cmdBase64(param base64.Base64Parameters) error {
	return base64.EncodeDecodeBase64(param)
}

func ConfigureBase64CommandLine(rootCmd *cobra.Command) {

	base64Cmd.Flags().StringVarP(&fileSrc, "input", "i", "", "File input to encode/decode")
	base64Cmd.Flags().StringVarP(&fileDest, "output", "o", "", "File output to encode/decode")
	base64Cmd.Flags().BoolVarP(&encode, "encode", "e", true, "Encode or decode in base64")
	base64Cmd.Flags().StringVarP(&bufferSize, "bufferSize", "b", "", "Size of buffer. Must be a number")
	rootCmd.AddCommand(base64Cmd)

}
