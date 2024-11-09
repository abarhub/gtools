package time

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

type TimeParameters struct {
	Cmd        string
	Parameters []string
}

func TimeCommand(param TimeParameters) error {

	start := time.Now()

	cmd := exec.Command(param.Cmd, param.Parameters...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("could not run command: %v", err)
	}

	elapsed := time.Since(start)
	log.Printf("Exec took %s", elapsed)

	return nil
}
