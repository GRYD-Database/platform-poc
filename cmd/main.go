package main

import (
	"log"

	"github.com/gryd-database/platform-poc/cmd/server"
)

func main() {
	if err := server.Init(); err != nil {
		log.Fatal(err)
	}
}
