package main

import (
	"log"
	"os"

	"github.com/GannettDigital/jstransform/generate"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <JSON Schema Path> [output directory]", os.Args[0])
	}

	var outputPath string
	if len(os.Args) == 3 {
		outputPath = os.Args[2]
	}
	err := generate.BuildStructs(os.Args[1], outputPath)
	if err == nil {
		log.Print("Finished")
	} else {
		log.Fatal(err)
	}
}
