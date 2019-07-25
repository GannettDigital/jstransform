package generate

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestBuildStructsRename(t *testing.T) {
	testdir := "test_data"
	tests := []struct {
		description    string
		file           string
		expectedFiles  []string
		wantFiles      []string
		useMessagePack bool
		renameStructs  map[string]string
	}{
		{
			file:           "complex.json",
			description:    "without oneOfTypes",
			expectedFiles:  []string{"complex.go"},
			wantFiles:      []string{"complex.go.out2"},
			useMessagePack: false,
			renameStructs:  nil,
		},
		{
			file:           "test_schema.json",
			description:    "with oneOfType",
			expectedFiles:  []string{"simple.go", "complex.go", "test_data_msgp.go", "test_data_msgp_test.go"},
			wantFiles:      []string{"simple.go.out2", "complex.go.out2", "test_data_msgp.go.out", "test_data_msgp_test.go.out"},
			useMessagePack: true,
			renameStructs:  nil,
		},
		{
			file:           "complex.json",
			description:    "without oneOfTypes, renamed",
			expectedFiles:  []string{"complex.go"},
			wantFiles:      []string{"complex.go.out-rename"},
			useMessagePack: false,
			renameStructs: map[string]string{
				"complex": "ReallyComplex",
			},
		},
		{
			file:           "test_schema.json",
			description:    "with oneOfType, renamed",
			expectedFiles:  []string{"simple.go", "complex.go", "test_data_msgp.go", "test_data_msgp_test.go"},
			wantFiles:      []string{"simple.go.out-rename", "complex.go.out-rename", "test_data_msgp.go.out-rename", "test_data_msgp_test.go.out-rename"},
			useMessagePack: true,
			renameStructs: map[string]string{
				"simple":  "TotallySimple",
				"complex": "ReallyComplex",
				"height":  "Not-Renamed",
				"Height":  "Not-Either",
			},
		},
	}

	for _, test := range tests {
		if err := BuildStructsRename(filepath.Join(testdir, test.file), testdir, test.useMessagePack, test.renameStructs); err != nil {
			t.Fatalf("Test %q - BuildStructsRename failed: %v", test.description, err)
		}

		for i := range test.expectedFiles {
			got, err := ioutil.ReadFile(filepath.Join(testdir, test.expectedFiles[i]))
			if err != nil {
				t.Errorf("Test %q - failed to read expected file %q: %v", test.description, test.expectedFiles[i], err)
			}

			want, err := ioutil.ReadFile(filepath.Join(testdir, test.wantFiles[i]))
			if err != nil {
				t.Errorf("Test %q - failed to read want file %q: %v", test.description, test.wantFiles[i], err)
			}

			if string(got) != string(want) {
				t.Errorf("Test %q - file %q got\n%s\n!= want\n%s", test.description, test.expectedFiles[i], got, want)
			}

			if err := os.Remove(filepath.Join(testdir, test.expectedFiles[i])); err != nil {
				t.Errorf("Test %q - failed to remove file %q: %v", test.description, test.expectedFiles[i], err)
			}
		}
	}
}

func TestGoType(t *testing.T) {
	tests := []struct {
		description string
		jsonType    string
		array       bool
		want        string
	}{
		{
			description: "JSON boolean",
			jsonType:    "boolean",
			want:        "bool",
		},
		{
			description: "JSON boolean",
			jsonType:    "boolean",
			array:       true,
			want:        "[]bool",
		},
		{
			description: "JSON integer",
			jsonType:    "integer",
			want:        "int64",
		},
		{
			description: "JSON integer",
			jsonType:    "integer",
			array:       true,
			want:        "[]int64",
		},
		{
			description: "JSON number",
			jsonType:    "number",
			want:        "float64",
		},
		{
			description: "JSON number",
			jsonType:    "number",
			array:       true,
			want:        "[]float64",
		},
		{
			description: "JSON string",
			jsonType:    "string",
			want:        "string",
		},
		{
			description: "JSON string",
			jsonType:    "string",
			array:       true,
			want:        "[]string",
		},
		{
			description: "JSON object",
			jsonType:    "object",
			want:        "struct",
		},
		{
			description: "JSON object",
			jsonType:    "object",
			array:       true,
			want:        "[]struct",
		},
		{
			description: "JSON string date-time",
			jsonType:    "date-time",
			array:       false,
			want:        "time.Time",
		},
		{
			description: "JSON string date-time array",
			jsonType:    "date-time",
			array:       true,
			want:        "[]time.Time",
		},
	}

	for _, test := range tests {
		got := goType(test.jsonType, test.array)
		if got != test.want {
			t.Errorf("Test %q - got %q, want %q", test.description, got, test.want)
		}
	}
}

func TestSplitJSONPath(t *testing.T) {
	tests := []struct {
		description string
		path        string
		want        []string
	}{
		{
			description: "Top level field",
			path:        "$.field1",
			want:        []string{"field1"},
		},
		{
			description: "deep field",
			path:        "$.field1.field2.field3",
			want:        []string{"field1", "field2", "field3"},
		},
		{
			description: "deep field with array",
			path:        "$.field1.field2.field3[*]",
			want:        []string{"field1", "field2", "field3"},
		},
	}

	for _, test := range tests {
		got := splitJSONPath(test.path)

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Test %q - got %v, want %v", test.description, got, test.want)
		}
	}
}
