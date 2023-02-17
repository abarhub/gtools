package commandline

import (
	"fmt"
	"github.com/spf13/cobra"
)

const VersionGTools = "1.0.0"

var commandError error = nil

var rootCmd = &cobra.Command{
	Use:     "gtools",
	Version: VersionGTools,
	Short:   "gtools - a simple CLI tools",
	Long: `gtools is a super simple CLI tools
   
simple in CLI`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("no command")
	},
}

func Run(args []string) error {

	ConfigureCopyCommandLine(rootCmd)
	ConfigureDuCommandLine(rootCmd)
	ConfigureBase64CommandLine(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		return fmt.Errorf("Whoops. There was an error while executing your CLI '%s'", err)
	}

	return commandError
}
