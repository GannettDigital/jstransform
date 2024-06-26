package generate

import (
	"go/ast"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestBuildAvroSchemaFile(t *testing.T) {
	tests := []struct {
		description string
		name        string
		goPath      string
		wantPath    string
	}{
		{
			description: "simple go struct",
			name:        "Simple",
			goPath:      "./avro_test_data/simple.go",
			wantPath:    "./avro_test_data/simple.avsc.out",
		},
		{
			description: "simple go struct - no nest",
			name:        "Simple",
			goPath:      "./avro_test_data/nonest/simple.go",
			wantPath:    "./avro_test_data/simple.avsc.out",
		},
		{
			description: "fields with repeated field names",
			name:        "Repeats",
			goPath:      "./avro_test_data/repeats.go",
			wantPath:    "./avro_test_data/repeats.avsc.out",
		},
		{
			description: "with embedded and nested structs, fields with descriptions",
			name:        "Complex",
			goPath:      "./avro_test_data/complex.go",
			wantPath:    "./avro_test_data/complex.avsc.out",
		},
		{
			description: "with embedded and named nested structs, fields with descriptions",
			name:        "Complex",
			goPath:      "./avro_test_data/nonest/complex.go",
			wantPath:    "./avro_test_data/complex.avsc.out",
		},
		{
			description: "Variations on Arrays",
			name:        "Arrays",
			goPath:      "./avro_test_data/arrays.go",
			wantPath:    "./avro_test_data/arrays.avsc.out",
		},
		{
			description: "Nested",
			name:        "Nested",
			goPath:      "./avro_test_data/nested.go",
			wantPath:    "./avro_test_data/nested.avsc.out",
		},
		{
			description: "Nested - nonest",
			name:        "Nested",
			goPath:      "./avro_test_data/nonest/nested.go",
			wantPath:    "./avro_test_data/nested.avsc.out",
		},
		{
			description: "Time fields",
			name:        "Times",
			goPath:      "./avro_test_data/times.go",
			wantPath:    "./avro_test_data/times.avsc.out",
		},
		{
			description: "Nested duplicate struct",
			name:        "Duplicate",
			goPath:      "./avro_test_data/duplicate.go",
			wantPath:    "./avro_test_data/duplicate.avsc.out",
		},
	}

	for _, test := range tests {
		outpath, err := buildAvroSchemaFile(test.name, test.goPath, true)
		if err != nil {
			t.Errorf("Test %q - failed to build Avro schema: %v", test.description, err)
		}

		got, err := os.ReadFile(outpath)
		if err != nil {
			t.Errorf("Test %q - failed to read Avro schema file: %v", test.description, err)
		}

		want, err := os.ReadFile(test.wantPath)
		if err != nil {
			t.Errorf("Test %q - failed to read want Avro schema file: %v", test.description, err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Test %q - got\n%s\nwant\n%s\n", test.description, got, want)
		}
		if err := os.Remove(outpath); err != nil {
			t.Errorf("Test %q - failed to remove generated Avro schema file at %q: %v", test.description, outpath, err)
		}
	}
}

func TestBuildAvroSerializationFunctions(t *testing.T) {
	tests := []struct {
		description string
		path        string
	}{
		{
			description: "simple avro file",
			path:        "./avro_test_data/simple.avsc.out",
		},
		{
			description: "fields with repeated field names",
			path:        "./avro_test_data/repeats.avsc.out",
		},
		{
			description: "complex avro file",
			path:        "./avro_test_data/complex.avsc.out",
		},
		{
			description: "Variations on Arrays",
			path:        "./avro_test_data/arrays.avsc.out",
		},
		{
			description: "Variations on time",
			path:        "./avro_test_data/times.avsc.out",
		},
	}

	for _, test := range tests {
		if err := buildAvroSerializationFunctions(test.path, BuildArgs{}); err != nil {
			t.Errorf("Test %q - failed: %v", test.description, err)
		}

		git := exec.Command("git", "diff", "--quiet", "*.go")
		schemaName := strings.Split(filepath.Base(test.path), ".")[0]
		git.Dir = filepath.Join("./avro_test_data/avro", schemaName)
		if err := git.Run(); err != nil {
			t.Errorf("Test %q - Differences in generated files found", test.description)
		}
	}
}

func TestParseGoStruct(t *testing.T) {
	tests := []struct {
		description string
		name        string
		path        string
	}{
		{
			description: "simple struct",
			name:        "Simple",
			path:        "./avro_test_data/simple.go",
		},
		{
			description: "a struct among many others",
			name:        "BuildArgs",
			path:        "./generate.go",
		},
		{
			description: "a struct among many files",
			name:        "BuildArgs",
			path:        ".",
		},
	}

	for _, test := range tests {
		got, err := parseGoStruct(test.name, test.path)
		if err != nil {
			t.Errorf("Test %q - got err: %v", test.description, err)
			continue
		}

		if got == nil {
			t.Errorf("Test %q - got nil", test.description)
		}

		if got.Name.Name != test.name {
			t.Errorf("Test %q - got unexpected name %q", test.description, got.Name.Name)
		}
	}
}

func TestParseStructTag(t *testing.T) {
	tests := []struct {
		description     string
		input           *ast.BasicLit
		wantName        string
		wantDescription string
		wantOmitEmpty   bool
	}{
		{
			description: "empty",
			input:       &ast.BasicLit{Value: ""},
		},
		{
			description:     "description only",
			input:           &ast.BasicLit{Value: `description:"blah"`},
			wantDescription: "blah",
		},
		{
			description: "json name only",
			input:       &ast.BasicLit{Value: `json:"name"`},
			wantName:    "name",
		},
		{
			description:   "json omitempty only",
			input:         &ast.BasicLit{Value: `json:",omitempty"`},
			wantOmitEmpty: true,
		},
		{
			description:   "json omitempty and name",
			input:         &ast.BasicLit{Value: `json:"name,omitempty"`},
			wantName:      "name",
			wantOmitEmpty: true,
		},
		{
			description:     "everything",
			input:           &ast.BasicLit{Value: `json:"name,omitempty" description:"blah"`},
			wantName:        "name",
			wantDescription: "blah",
			wantOmitEmpty:   true,
		},
		{
			description:     "everything, reverse order",
			input:           &ast.BasicLit{Value: `description:"blah" json:"name,omitempty"`},
			wantName:        "name",
			wantDescription: "blah",
			wantOmitEmpty:   true,
		},
	}

	for _, test := range tests {
		gotName, gotDescription, gotOmitEmpty := parseStructTag(test.input)

		if got, want := gotName, test.wantName; got != want {
			t.Errorf("Test %q - got name %q, want %q", test.description, got, want)
		}
		if got, want := gotDescription, test.wantDescription; got != want {
			t.Errorf("Test %q - got description %q, want %q", test.description, got, want)
		}
		if got, want := gotOmitEmpty, test.wantOmitEmpty; got != want {
			t.Errorf("Test %q - got omit empty %t, want %t", test.description, got, want)
		}
	}
}
