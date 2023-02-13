package main

import (
	"gtools/internal/utils"
	"log"
)

func main() {
	err := utils.Run()
	if err != nil {
		log.Fatal(err)
	}
}
