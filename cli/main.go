package main

import (
	"log"
	"os"

	"github.com/iplcm/access/tools"
)

func main() {
	arg := os.Args
	if arg[len(arg)-1] == "uuid" {
		log.Println(tools.NewUUID().String())
	}
}
