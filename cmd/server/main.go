package main

import (
	"log"

	"github.com/codecrafters-io/http-server-starter-go/internal/http"
)

func main() {
	srv := http.Server{}
	err := srv.Run("0.0.0.0:4221")
	if err != nil {
		log.Fatalf("failure during run: %v", err)
	}
}
