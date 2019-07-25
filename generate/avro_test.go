package generate

import (
	"go/ast"
	"io/ioutil"
	"os"
	"reflect"
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
			goPath:      "./test_data/simple.go.out",
			wantPath:    "./test_data/simple.avsc.out",
		},
		{
			description: "fields with repeated field names",
			name:        "Repeats",
			goPath:      "./test_data/repeats.go.out",
			wantPath:    "./test_data/repeats.avsc.out",
		},
		{
			description: "with embedded and nested structs, fields with descriptions",
			name:        "Complex",
			goPath:      "./test_data/complex.go.out",
			wantPath:    "./test_data/complex.avsc.out",
		},
	}

	// For the complex.go to find the embedded struct simple.go must exist, other tests use that filename so create
	// and destroy it here
	if err := os.Link("./test_data/simple.go.out", "./test_data/simple.go"); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("./test_data/simple.go")

	for _, test := range tests {
		spec, err := parseGoStruct(test.name, test.goPath)
		if err != nil {
			t.Fatalf("Test %q - failed parsing go struct: %v", test.description, err)
		}

		outpath, err := buildAvroSchemaFile(test.name, "./test_data", spec, true)
		if err != nil {
			t.Errorf("Test %q - failed to build Avro schema: %v", test.description, err)
		}

		got, err := ioutil.ReadFile(outpath)
		if err != nil {
			t.Errorf("Test %q - failed to read Avro schema file: %v", test.description, err)
		}

		want, err := ioutil.ReadFile(test.wantPath)
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

func TestParseGoStruct(t *testing.T) {
	tests := []struct {
		description string
		name        string
		path        string
	}{
		{
			description: "simple struct",
			name:        "Simple",
			path:        "./test_data/simple.go.out",
		},
		{
			description: "a struct among many others",
			name:        "BuildArgs",
			path:        "./generate.go",
		},
	}

	for _, test := range tests {
		got, err := parseGoStruct(test.name, test.path)
		if err != nil {
			t.Errorf("Test %q - got err: %v", test.description, err)
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
