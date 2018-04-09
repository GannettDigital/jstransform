package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/GannettDigital/jstransform/generate"
)

//go:generate go run examples/dynamic.go
func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <JSON Schema Path> [output directory]\n", path.Base(os.Args[0]))
		os.Exit(1)
	}

	inputPath, err := filepath.Abs(os.Args[1])
	if err != nil {
		fmt.Printf("Input directory %s error: %v", os.Args[1], err)
		os.Exit(2)
	}

	var outputPath string
	if len(os.Args) == 3 {
		outputPath, err = filepath.Abs(os.Args[2])
		if err != nil {
			fmt.Printf("Output directory %s error: %v", os.Args[2], err)
			os.Exit(3)
		}
	}

	if err = generate.BuildStructs(inputPath, outputPath); err != nil {
		fmt.Printf("Golang Struct generation failed: %v", err)
		os.Exit(4)
	}
}
