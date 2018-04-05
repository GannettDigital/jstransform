package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/GannettDigital/jstransform/generate"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <JSON Schema Path> [output directory]\n", path.Base(os.Args[0]))
		os.Exit(0)
	}

	var outputPath string
	if len(os.Args) == 3 {
		outputPath, _ = filepath.Abs(os.Args[2])
	}

	inputPath, _ := filepath.Abs(os.Args[1])
	err := generate.BuildStructs(inputPath, outputPath)
	if err == nil {
		fmt.Printf("Finished, check output directory %s\n", outputPath)
	} else {
		panic(err)
	}
}
