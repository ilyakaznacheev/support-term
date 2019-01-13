package main

import (
	"log"

	"github.com/ilyakaznacheev/support-term/internal/support"
)

func main() {
	err := support.Run()
	if err != nil {
		log.Fatal(err)
	}
}
