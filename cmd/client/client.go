package main

import (
	"log"

	"github.com/ilyakaznacheev/support-term/internal/client"
)

func main() {
	err := client.Run()
	if err != nil {
		log.Fatal(err)
	}
}
