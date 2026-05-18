package main

import (
	"log"
	"os"
)

const specPath = "openapi.generated.yaml"

func main() {
	enumSites, compSites, reqSites, err := discoverAll(specPath)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	src, err := generateAll(enumSites, compSites, reqSites)
	if err != nil {
		log.Fatalf("error generating: %v", err)
	}

	if _, err := os.Stdout.Write(src); err != nil {
		log.Fatalf("error writing output: %v", err)
	}
}
