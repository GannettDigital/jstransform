package generate

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// TestBuildStructs generates go files and compares them to their corresponding files located at generate_test_data/{outDir}
// To add new test cases that fit into the existing cases:
// 1. Generate go files in the generate_test_data/{outDir} that matches your test case
// To add new test cases that do NOT fit into the existing cases:
// 1. Add a new output directory with the name of your test case to generate_test_data/{outDir} (use one word all lowercase so the go files package name is simple)
// 2. Generate go files in the new output directory at generate_test_data/{outDir}.
func TestBuildStructs(t *testing.T) {
	testdir := "generate_test_data"

	tests := []struct {
		description string
		// BuildArgs.OutputDir is the directory name that holds the go files
		// go files created by these tests exist at {outDir} (this directory is cleaned up after each test)
		// expected go files that are used to compare against the test created files exist at generate_test_data/{outDir}
		buildArgs BuildArgs
		files     []string
	}{
		{
			description: "without oneOfTypes",
			buildArgs: BuildArgs{
				SchemaPath:             filepath.Join(testdir, "complex.json"),
				OutputDir:              "generated",
				GenerateMessagePack:    false,
				StructNameMap:          nil,
				DescriptionAsStructTag: true,
			},
			files: []string{"complex.go"},
		},
		{
			description: "one allOf with additional properties at the top level",
			buildArgs: BuildArgs{
				SchemaPath:             filepath.Join(testdir, "all_of_with_properties.json"),
				OutputDir:              "generated",
				GenerateMessagePack:    false,
				StructNameMap:          nil,
				DescriptionAsStructTag: true,
			},
			files: []string{"all_of_with_properties.go", "simple.go"},
		},
		{
			description: "test formatting of times",
			buildArgs: BuildArgs{
				SchemaPath:             filepath.Join(testdir, "times.json"),
				OutputDir:              "generated",
				GenerateMessagePack:    false,
				StructNameMap:          nil,
				DescriptionAsStructTag: true,
			},
			files: []string{"times.go"},
		},
		{
			description: "nested array structs",
			buildArgs: BuildArgs{
				SchemaPath:             filepath.Join(testdir, "nested.json"),
				OutputDir:              "generated",
				GenerateMessagePack:    false,
				StructNameMap:          nil,
				DescriptionAsStructTag: false,
				NoNestedStructs:        false,
			},
			files: []string{"nested.go"},
		},
		{
			description: "nested to primitive array structs",
			buildArgs: BuildArgs{
				SchemaPath:             filepath.Join(testdir, "nested_to_primitive.json"),
				OutputDir:              "generated",
				GenerateMessagePack:    false,
				StructNameMap:          nil,
				DescriptionAsStructTag: false,
				NoNestedStructs:        false,
			},
			files: []string{"nested_to_primitive.go"},
		},
		{
			description: "with oneOfType",
			buildArgs: BuildArgs{
				SchemaPath:             filepath.Join(testdir, "test_schema.json"),
				OutputDir:              "msgp",
				GenerateMessagePack:    true,
				StructNameMap:          nil,
				DescriptionAsStructTag: true,
			},
			files: []string{"simple.go", "complex.go", "msgp_msgp.go", "msgp_msgp_test.go"},
		},
		{
			description: "without oneOfTypes, renamed",
			buildArgs: BuildArgs{
				SchemaPath:          filepath.Join(testdir, "complex.json"),
				OutputDir:           "rename",
				GenerateMessagePack: false,
				StructNameMap: map[string]string{
					"complex": "ReallyComplex",
				},
				DescriptionAsStructTag: true,
			},
			files: []string{"complex.go"},
		},
		{
			description: "with oneOfType, renamed",
			buildArgs: BuildArgs{
				SchemaPath:          filepath.Join(testdir, "test_schema.json"),
				OutputDir:           "rename",
				GenerateMessagePack: true,
				StructNameMap: map[string]string{
					"simple":  "TotallySimple",
					"complex": "ReallyComplex",
					"height":  "Not-Renamed",
					"Height":  "Not-Either",
				},
				DescriptionAsStructTag: true,
			},
			files: []string{"simple.go", "complex.go", "rename_msgp.go", "rename_msgp_test.go"},
		},
		{
			description: "without oneOfTypes, with no nested structs and descriptions as comments",
			buildArgs: BuildArgs{
				SchemaPath:             filepath.Join(testdir, "test_schema2.json"),
				OutputDir:              "nonest",
				GenerateMessagePack:    false,
				StructNameMap:          nil,
				DescriptionAsStructTag: false,
				NoNestedStructs:        true,
			},
			files: []string{"simple.go", "complex.go"},
		},
		{
			description: "nested array structs - nonest",
			buildArgs: BuildArgs{
				SchemaPath:             filepath.Join(testdir, "nested.json"),
				OutputDir:              "nonest",
				GenerateMessagePack:    false,
				StructNameMap:          nil,
				DescriptionAsStructTag: false,
				NoNestedStructs:        true,
			},
			files: []string{"nested.go"},
		},
		{
			description: "nested to primitive array structs - nonest",
			buildArgs: BuildArgs{
				SchemaPath:             filepath.Join(testdir, "nested_to_primitive.json"),
				OutputDir:              "nonest",
				GenerateMessagePack:    false,
				StructNameMap:          nil,
				DescriptionAsStructTag: false,
				NoNestedStructs:        true,
			},
			files: []string{"nested_to_primitive.go"},
		},
		{
			description: "without oneOfTypes, with no nested structs and descriptions as comments - pointers",
			buildArgs: BuildArgs{
				SchemaPath:             filepath.Join(testdir, "test_schema2.json"),
				OutputDir:              "pointers",
				GenerateMessagePack:    true,
				StructNameMap:          nil,
				DescriptionAsStructTag: false,
				NoNestedStructs:        true,
				Pointers:               true,
			},
			files: []string{"simple.go", "complex.go"},
		},
		{
			description: "nested array structs - pointers",
			buildArgs: BuildArgs{
				SchemaPath:             filepath.Join(testdir, "nested.json"),
				OutputDir:              "pointers",
				GenerateMessagePack:    false,
				StructNameMap:          nil,
				DescriptionAsStructTag: false,
				NoNestedStructs:        true,
				Pointers:               true,
			},
			files: []string{"nested.go"},
		},
		{
			description: "nested to primitive array structs - pointers",
			buildArgs: BuildArgs{
				SchemaPath:             filepath.Join(testdir, "nested_to_primitive.json"),
				OutputDir:              "pointers",
				GenerateMessagePack:    false,
				StructNameMap:          nil,
				DescriptionAsStructTag: false,
				NoNestedStructs:        true,
				Pointers:               true,
			},
			files: []string{"nested_to_primitive.go"},
		},
		{
			description: "test formatting of times - pointers",
			buildArgs: BuildArgs{
				SchemaPath:             filepath.Join(testdir, "times.json"),
				OutputDir:              "pointers",
				GenerateMessagePack:    false,
				StructNameMap:          nil,
				DescriptionAsStructTag: false,
				NoNestedStructs:        true,
				Pointers:               true,
			},
			files: []string{"times.go"},
		},
		{
			description: "embedded allOf",
			buildArgs: BuildArgs{
				SchemaPath:             filepath.Join(testdir, "base.json"),
				OutputDir:              "generated",
				GenerateMessagePack:    false,
				StructNameMap:          nil,
				DescriptionAsStructTag: false,
				NoNestedStructs:        true,
				EmbedAllOf:             true,
			},
			files: []string{"simple_no_nested.go", "embedded.go"},
		},
	}

	for _, test := range tests {
		outDir := test.buildArgs.OutputDir
		os.Mkdir(outDir, os.ModePerm|os.ModePerm)

		if err := BuildStructsWithArgs(test.buildArgs); err != nil {
			t.Fatalf("Test %q - BuildStructsRename failed: %v", test.description, err)
		}

		for i := range test.files {
			got, err := os.ReadFile(filepath.Join(outDir, test.files[i]))
			if err != nil {
				t.Errorf("Test %q - failed to read expected file %q: %v", test.description, test.files[i], err)
			}

			want, err := os.ReadFile(filepath.Join(testdir, outDir, test.files[i]))
			if err != nil {
				t.Errorf("Test %q - failed to read want file %q: %v", test.description, test.files[i], err)
			}

			if string(got) != string(want) {
				t.Errorf("Test %q - file %q got\n%s\n!= want\n%s", test.description, test.files[i], got, want)
			}
		}

		if err := os.RemoveAll(outDir); err != nil {
			t.Errorf("Test %q - failed to cleanup output dir %s for test generated files", test.description, outDir)
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
