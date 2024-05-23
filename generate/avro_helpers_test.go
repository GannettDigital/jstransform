package generate

import (
	"os/exec"
	"path/filepath"
	"testing"
)

func TestBuildAvroHelperFunctions(t *testing.T) {
	tests := []struct {
		description string
		name        string
		filename    string
		importPath  string
		wantPath    string
	}{
		{
			description: "simple avro file",
			name:        "Simple",
			filename:    "simple.go",
			wantPath:    "simple_avro.go",
		},
		{
			description: "simple avro file with named nested structs",
			name:        "Simple",
			filename:    "nonest/simple.go",
			importPath:  "github.com/GannettDigital/jstransform/generate/avro_test_data/nonest",
			wantPath:    "nonest/simple_avro.go",
		},
		{
			description: "fields with repeated field names",
			name:        "Repeats",
			filename:    "repeats.go",
			wantPath:    "repeats_avro.go",
		},
		{
			description: "complex avro file",
			name:        "Complex",
			filename:    "complex.go",
			wantPath:    "complex_avro.go",
		},
		{
			description: "complex avro file with named nested structs",
			name:        "Complex",
			filename:    "nonest/complex.go",
			importPath:  "github.com/GannettDigital/jstransform/generate/avro_test_data/nonest",
			wantPath:    "nonest/complex_avro.go",
		},
		{
			description: "variations on Arrays",
			name:        "Arrays",
			filename:    "arrays.go",
			wantPath:    "arrays_avro.go",
		},
		{
			description: "nested array structs",
			name:        "Nested",
			filename:    "nested.go",
			wantPath:    "nested_avro.go",
		},
		{
			description: "nested array structs with named nested structs",
			name:        "Nested",
			filename:    "nonest/nested.go",
			importPath:  "github.com/GannettDigital/jstransform/generate/avro_test_data/nonest",
			wantPath:    "nonest/nested_avro.go",
		},
		{
			description: "variations on Times",
			name:        "Times",
			filename:    "times.go",
			wantPath:    "times_avro.go",
		},
	}

	testPath := "./avro_test_data"
	defaultImportPath := "github.com/GannettDigital/jstransform/generate/avro_test_data"
	for _, test := range tests {
		testImportPath := defaultImportPath
		if test.importPath != "" {
			testImportPath = test.importPath
		}
		if err := buildAvroHelperFunctions(test.name, filepath.Join(testPath, test.filename), testImportPath); err != nil {
			t.Errorf("Test %q - failed: %v", test.description, err)
		}

		git := exec.Command("git", "diff", "--quiet", test.wantPath)
		git.Dir = testPath
		if err := git.Run(); err != nil {
			t.Errorf("Test %q - failed git diff of generated file: %v", test.description, err)
		}
	}
}

func TestStructToFilename(t *testing.T) {
	tests := []struct {
		description string
		structName  string
		wantPath    string
	}{
		{
			description: "Simple",
			structName:  "Simple",
			wantPath:    "simple.go",
		},
		{
			description: "CamelCase",
			structName:  "CamelCase",
			wantPath:    "camelCase.go",
		},
		{
			description: "Snake_Case",
			structName:  "Snake_Case",
			wantPath:    "snake_Case.go",
		},
	}

	for _, test := range tests {
		if got := structToFilename(test.structName); got != test.wantPath {
			t.Errorf("Test %q - got unexpected path %q", test.description, got)
		}
	}
}
