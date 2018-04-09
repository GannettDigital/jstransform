package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var schemas = map[string]string{
	"content-schema": "../content-api/schema/v1/assets",
	"test-schema":    "generate/test_data/buildstructs",
}

func main() {
	for title, source := range schemas {
		schemaPath, err := filepath.Abs(source)
		if err != nil {
			fmt.Printf("Path for %s not present!\n", title)
		} else {
			fmt.Printf("\nExamples using %s:\n", title)
			if stat, err := os.Stat(schemaPath); err == nil && stat.IsDir() {
				filepath.Walk(schemaPath, func(jsonPath string, f os.FileInfo, err error) error {
					if strings.HasSuffix(jsonPath, ".json") {
						fmt.Printf("\tgo run main.go %s %s\n", path.Base(jsonPath), source)
					}
					return nil
				})
			}
		}
	}
}
