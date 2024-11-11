package time

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

type TimeParameters struct {
	Cmd        string
	Parameters []string
}

func TimeCommand(param TimeParameters) error {

	return timeCommande(param, os.Stdout)
}

func timeCommande(param TimeParameters, out io.Writer) error {
	start := time.Now()

	cmd := exec.Command(param.Cmd, param.Parameters...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("could not run command: %v", err)
	}

	elapsed := time.Since(start)
	_, err := fmt.Fprintf(out, "Exec took %s", elapsed)
	if err != nil {
		return fmt.Errorf("could not run print: %v", err)
	}

	return nil
}
