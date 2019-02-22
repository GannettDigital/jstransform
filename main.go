package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/GannettDigital/jstransform/generate"
)

// mapFlags allows for "-opt key=value" flags.
type mapFlags struct {
	kv map[string]string
}
// String converts the provided options into the input format.
func (mf *mapFlags) String() string {
	kvs := make([]string, 0, len(mf.kv))
	for k, v := range mf.kv {
		kvs = append(kvs, fmt.Sprintf("-%s=%s", k, v))
	}
	return strings.Join(kvs, " ")
}
// Set stores the provided options into the data structure.
func (mf *mapFlags) Set(value string) error {
	kv := strings.SplitN(value, "=", 2)
	if len(kv) != 2 {
		return fmt.Errorf("Value in 'key=value' format: %v", value)
	}
	mf.kv[kv[0]] = kv[1]
	return nil
}

func main() {
	renameStructs := mapFlags{kv: make(map[string]string)}
	var useMessagePack bool

	flag.Var(&renameStructs, "rename", "Override generated name of structure; use '-rename old=new'.")
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

	if err = generate.BuildStructsRename(inputPath, outputPath, useMessagePack, renameStructs.kv); err != nil {
		fmt.Printf("Golang Struct generation failed: %v", err)
		os.Exit(4)
	}
}
