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

	// Note that 'go generate' only runs lines that start with '//go:generate' (no indentation)
	// but having each 'go:generate' comment with its respective test case will get indented
	// by 'go fmt'.  So to run the generation commands, temporarily remove the whitespace before
	// the '//go:generate' comments.
	tests := []struct {
		description string
		// BuildArgs.OutputDir is the directory name that holds the go files
		// go files created by these tests exist at {outDir} (this directory is cleaned up after each test)
		// expected go files that are used to compare against the test created files exist at generate_test_data/{outDir}
		buildArgs BuildArgs
		files     []string
	}{
		//go:generate go run .. --descriptionAsStructTag=true --embedAllOf=false --nestedStructs=true generate_test_data/complex.json generate_test_data/generated
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
		//go:generate go run .. --descriptionAsStructTag=true --embedAllOf=false --nestedStructs=true generate_test_data/all_of_with_properties.json generate_test_data/generated
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
		//go:generate go run .. --descriptionAsStructTag=true --embedAllOf=false --nestedStructs=true generate_test_data/times.json generate_test_data/generated
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
		//go:generate go run .. --descriptionAsStructTag=false --embedAllOf=false --nestedStructs=true generate_test_data/nested.json generate_test_data/generated
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
		//go:generate go run .. --descriptionAsStructTag=true --embedAllOf=false --nestedStructs=true generate_test_data/nested_to_primitive.json generate_test_data/generated
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
		//go:generate go run .. --descriptionAsStructTag=true --embedAllOf=false --nestedStructs=true --msgp generate_test_data/test_schema.json generate_test_data/msgp
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
		//go:generate go run .. --descriptionAsStructTag=true --embedAllOf=false --nestedStructs=false --rename complex=ReallyComplex generate_test_data/complex.json generate_test_data/rename
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
		//go:generate go run .. --descriptionAsStructTag=true --embedAllOf=false --nestedStructs=true --msgp --rename simple=TotallySimple,complex=ReallyComplex,height=Not-Renamed,Height=Not-Either generate_test_data/test_schema.json generate_test_data/rename
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
		//go:generate go run .. --descriptionAsStructTag=false --embedAllOf=false --nestedStructs=false generate_test_data/test_schema2.json generate_test_data/nonest
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
		//go:generate go run .. --descriptionAsStructTag=false --embedAllOf=false --nestedStructs=false generate_test_data/nested.json generate_test_data/nonest
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
		//go:generate go run .. --descriptionAsStructTag=false --embedAllOf=false --nestedStructs=false generate_test_data/nested_to_primitive.json generate_test_data/nonest
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
		//go:generate go run .. --descriptionAsStructTag=false --embedAllOf=false --nestedStructs=false --msgp --pointers generate_test_data/test_schema2.json generate_test_data/pointers
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
		//go:generate go run .. --descriptionAsStructTag=false --embedAllOf=false --nestedStructs=false --pointers generate_test_data/nested.json generate_test_data/pointers
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
		//go:generate go run .. --descriptionAsStructTag=false --embedAllOf=false --nestedStructs=false --pointers generate_test_data/nested_to_primitive.json generate_test_data/pointers
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
		//go:generate go run .. --descriptionAsStructTag=false --embedAllOf=false --nestedStructs=false --pointers generate_test_data/times.json generate_test_data/pointers
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
		//go:generate go run .. --descriptionAsStructTag=false --embedAllOf=true --nestedStructs=false generate_test_data/base.json generate_test_data/generated
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
		//go:generate go run .. --descriptionAsStructTag=false --embedAllOf=false --nestedStructs=false generate_test_data/simple_map.json generate_test_data/generated
		{
			description: "simple map",
			buildArgs: BuildArgs{
				SchemaPath:             filepath.Join(testdir, "simple_map.json"),
				OutputDir:              "generated",
				GenerateMessagePack:    false,
				StructNameMap:          nil,
				DescriptionAsStructTag: false,
				NoNestedStructs:        true,
			},
			files: []string{"simple_map.go"},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			outDir := test.buildArgs.OutputDir
			if err := os.Mkdir(outDir, 0o750); err != nil {
				t.Fatalf("Test %q - failed to create outDir %q: %v", test.description, outDir, err)
			}
			defer func() {
				if err := os.RemoveAll(outDir); err != nil {
					t.Errorf("Test %q - failed to cleanup output dir %s for test generated files", test.description, outDir)
				}
			}()

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
		})
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
