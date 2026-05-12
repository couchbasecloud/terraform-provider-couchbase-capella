package main

import (
	"log"
	"os"
)

const specPath = "openapi.generated.yaml"

func main() {
	sites, err := discover(specPath)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	src, err := generate(sites)
	if err != nil {
		log.Fatalf("error generating: %v", err)
	}

	if _, err := os.Stdout.Write(src); err != nil {
		log.Fatalf("error writing output: %v", err)
	}
}
