package commandline

import (
	"github.com/spf13/cobra"
	"gtools/internal/command/base64"
	"gtools/internal/utils"
)

var (
	fileSrc  string
	fileDest string
	encode   bool
)

var base64Cmd = &cobra.Command{
	Use:   "base64",
	Short: "encode/decode in base64",
	RunE: func(cmd *cobra.Command, args []string) error {
		var input *utils.InputParameter
		var output *utils.OutputParameter
		var err error
		if fileSrc != "" {
			input, err = utils.FileInputParameter(fileSrc)
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
			output, err = utils.FileOutputParameter(fileDest)
			if err != nil {
				return err
			}
		} else {
			output, err = utils.StdOutputParameter()
		}
		param := base64.Base64Parameters{input, output, encode}
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
	rootCmd.AddCommand(base64Cmd)

}
