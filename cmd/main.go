package main

import (
	"gtools/internal/utils"
	"log"
)

func main() {
	println("hello")
	err := utils.CopyDir("./testdir/test1", "./testdir/test_out")
	if err != nil {
		log.Fatal(err)
	}
}
