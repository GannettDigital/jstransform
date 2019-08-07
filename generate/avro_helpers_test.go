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
		wantPath    string
	}{
		{
			description: "simple avro file",
			name:        "Simple",
			filename:    "simple.go",
			wantPath:    "simple_avro.go",
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
			description: "variations on Arrays",
			name:        "Arrays",
			filename:    "arrays.go",
			wantPath:    "arrays_avro.go",
		},
	}

	testPath := "./test_data"
	testImportPath := "github.com/GannettDigital/jstransform/generate/test_data"
	for _, test := range tests {
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
