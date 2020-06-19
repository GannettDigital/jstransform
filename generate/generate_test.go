package generate

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestBuildStructs(t *testing.T) {
	testdir := "test_data"
	outDir := filepath.Join(testdir, "rename")
	os.Mkdir(outDir, os.ModePerm|os.ModePerm)
	defer os.RemoveAll(outDir)
	outDir2 := "nonest"
	os.Mkdir(outDir2, os.ModePerm|os.ModePerm)
	defer os.RemoveAll(outDir2)

	tests := []struct {
		description   string
		buildArgs     BuildArgs
		expectedFiles []string
		wantFiles     []string
	}{
		{
			description: "without oneOfTypes",
			buildArgs: BuildArgs{
				SchemaPath:             filepath.Join(testdir, "complex.json"),
				OutputDir:              outDir,
				GenerateMessagePack:    false,
				StructNameMap:          nil,
				DescriptionAsStructTag: true,
			},
			expectedFiles: []string{"complex.go"},
			wantFiles:     []string{"complex.go.out2"},
		},
		{
			description: "without oneOfTypes, with no nested structs and descriptions as comments",
			buildArgs: BuildArgs{
				SchemaPath:             filepath.Join(testdir, "test_schema2.json"),
				OutputDir:              outDir2,
				GenerateMessagePack:    false,
				StructNameMap:          nil,
				DescriptionAsStructTag: false,
				NoNestedStructs:        true,
			},
			expectedFiles: []string{"simple.go", "complex.go"},
			wantFiles:     []string{"nonest/simple.go", "nonest/complex.go"},
		},
		{
			description: "with oneOfType",
			buildArgs: BuildArgs{
				SchemaPath:             filepath.Join(testdir, "test_schema.json"),
				OutputDir:              outDir,
				GenerateMessagePack:    true,
				StructNameMap:          nil,
				DescriptionAsStructTag: true,
			},
			expectedFiles: []string{"simple.go", "complex.go", "rename_msgp.go", "rename_msgp_test.go"},
			wantFiles:     []string{"simple.go.out2", "complex.go.out2", "rename_msgp.go.out", "rename_msgp_test.go.out"},
		},
		{
			description: "without oneOfTypes, renamed",
			buildArgs: BuildArgs{
				SchemaPath:          filepath.Join(testdir, "complex.json"),
				OutputDir:           outDir,
				GenerateMessagePack: false,
				StructNameMap: map[string]string{
					"complex": "ReallyComplex",
				},
				DescriptionAsStructTag: true,
			},
			expectedFiles: []string{"complex.go"},
			wantFiles:     []string{"complex.go.out-rename"},
		},
		{
			description: "with oneOfType, renamed",
			buildArgs: BuildArgs{
				SchemaPath:          filepath.Join(testdir, "test_schema.json"),
				OutputDir:           outDir,
				GenerateMessagePack: true,
				StructNameMap: map[string]string{
					"simple":  "TotallySimple",
					"complex": "ReallyComplex",
					"height":  "Not-Renamed",
					"Height":  "Not-Either",
				},
				DescriptionAsStructTag: true,
			},
			expectedFiles: []string{"simple.go", "complex.go", "rename_msgp.go", "rename_msgp_test.go"},
			wantFiles:     []string{"simple.go.out-rename", "complex.go.out-rename", "rename_msgp.go.out-rename", "rename_msgp_test.go.out-rename"},
		},
		{
			description: "one allOf with additional properties at the top level",
			buildArgs: BuildArgs{
				SchemaPath:             filepath.Join(testdir, "all_of_with_properties.json"),
				OutputDir:              outDir,
				GenerateMessagePack:    false,
				StructNameMap:          nil,
				DescriptionAsStructTag: true,
			},
			expectedFiles: []string{"all_of_with_properties.go", "simple.go"},
			wantFiles:     []string{"all_of_with_properties.go.out", "simple.go.out2"},
		},
		{
			description: "test formatting of times",
			buildArgs: BuildArgs{
				SchemaPath:             filepath.Join(testdir, "times.json"),
				OutputDir:              outDir,
				GenerateMessagePack:    false,
				StructNameMap:          nil,
				DescriptionAsStructTag: true,
			},
			expectedFiles: []string{"times.go"},
			wantFiles:     []string{"times.go.out"},
		},
	}

	for _, test := range tests {
		if err := BuildStructsWithArgs(test.buildArgs); err != nil {
			t.Fatalf("Test %q - BuildStructsRename failed: %v", test.description, err)
		}

		for i := range test.expectedFiles {
			got, err := ioutil.ReadFile(filepath.Join(test.buildArgs.OutputDir, test.expectedFiles[i]))
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
		}
	}
}

func TestGoType(t *testing.T) {
	tests := []struct {
		description string
		jsonType    string
		array       bool
		required    bool
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
			required:    true,
			want:        "time.Time",
		},
		{
			description: "JSON string date-time array",
			jsonType:    "date-time",
			array:       true,
			required:    true,
			want:        "[]time.Time",
		},
		{
			description: "JSON string date-time, omitempty",
			jsonType:    "date-time",
			array:       false,
			want:        "*time.Time",
		},
		{
			description: "JSON string date-time array, omitempty",
			jsonType:    "date-time",
			array:       true,
			want:        "[]*time.Time",
		},
	}

	for _, test := range tests {
		got := goType(test.jsonType, test.array, test.required)
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
