package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/GannettDigital/jstransform/generate"
)

func main() {
	var useMessagePack bool

	flag.BoolVar(&useMessagePack, "msgp", false, "generate MessagePack serialization methods")

	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		fmt.Printf("Usage: %s [-msgp] <JSON Schema Path> [output directory]\n", path.Base(os.Args[0]))
		os.Exit(1)
	}

	inputPath, err := filepath.Abs(args[0])
	if err != nil {
		fmt.Printf("Input directory \"%s\" error: %v", args[0], err)
		os.Exit(2)
	}

	var outputPath string
	if len(args) > 1 {
		outputPath, err = filepath.Abs(args[1])
		if err != nil {
			fmt.Printf("Output directory \"%s\" error: %v", args[1], err)
			os.Exit(3)
		}
	}

	if err = generate.BuildStructs(inputPath, outputPath, useMessagePack); err != nil {
		fmt.Printf("Golang Struct generation failed: %v", err)
		os.Exit(4)
	}
}
