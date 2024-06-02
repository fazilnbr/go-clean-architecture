package main

import (
	"log"

	"github.com/fazilnbr/go-clean-architecture/internal/server"
)

func main() {
	srv := server.New()
	if err := srv.Start(); err != nil {
		log.Fatalf("failed to start the server, err: %s", err.Error())
	}
}
