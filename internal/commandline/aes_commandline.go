package commandline

import (
	"github.com/spf13/cobra"
	"gtools/internal/command/aescrypt"
)

var (
	inputFileAes   string
	outpoutFileAes string
	decryptAes     bool
	passwordAes    string
)

var aesCmd = &cobra.Command{
	Use:   "aes",
	Short: "aes encrypt/decrypt",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		param := aescrypt.AesParameters{InputFile: inputFileAes, OutpoutFile: outpoutFileAes, Encrypt: !decryptAes,
			Password: passwordAes}
		commandError = aescrypt.AesCommand(param)
		return nil
	},
	DisableFlagsInUseLine: true,
}

func ConfigureAESCommandLine(rootCmd *cobra.Command) {
	aesCmd.Flags().StringVarP(&inputFileAes, "input", "i", "", "File input")
	aesCmd.Flags().StringVarP(&outpoutFileAes, "output", "o", "", "File output")
	aesCmd.Flags().BoolVarP(&decryptAes, "decrypt", "d", false, "decrypt. encrypt by default")
	aesCmd.Flags().StringVarP(&passwordAes, "password", "p", "", "Password in base64. if not present for crypt, the password was generated in base 64 and display in stdout")
	rootCmd.AddCommand(aesCmd)
}
