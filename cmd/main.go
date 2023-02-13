package main

import (
	"gtools/internal/utils"
	"log"
	"os"
)

func main() {
	err := utils.Run(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
}
