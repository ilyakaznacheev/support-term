package main

import (
	"log"

	generator "github.com/ilyakaznacheev/support-term/internal/id-generator"
)

func main() {
	err := generator.Run()
	if err != nil {
		log.Fatal(err)
	}
}
