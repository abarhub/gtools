package main

import (
	"gtools/internal/commandline"
	"log"
	"os"
)

func main() {
	err := commandline.Run(os.Args[1:])
	if err != nil {
		if err.Error() != "" {
			log.Fatal(err)
		} else {
			os.Exit(1)
		}
	}
}
