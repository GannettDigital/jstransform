package generate

import (
	"bytes"
	"os"
	"reflect"
	"testing"

	"github.com/GannettDigital/jstransform/jsonschema"
)

func TestAddField(t *testing.T) {
	tests := []struct {
		description string
		fields      extractedFields
		tree        []string
		instance    jsonschema.Instance
		want        extractedFields
	}{
		{
			description: "Simple scalar field",
			fields:      make(map[string]*extractedField),
			tree:        []string{"field"},
			instance:    jsonschema.Instance{Type: []string{"string"}},
			want: extractedFields{
				"field": &extractedField{
					name:     "Field",
					jsonType: "string",
					jsonName: "field",
				},
			},
		},
		{
			description: "Array field",
			fields:      make(map[string]*extractedField),
			tree:        []string{"arrayfield"},
			instance:    jsonschema.Instance{Type: []string{"array"}, Items: []byte(`{ "type": "string" }`)},
			want: extractedFields{
				"arrayfield": &extractedField{
					name:     "Arrayfield",
					jsonType: "",
					jsonName: "arrayfield",
					array:    true,
				},
			},
		},
		{
			description: "Array field 2nd call",
			fields: extractedFields{
				"arrayfield": &extractedField{
					name:     "Arrayfield",
					jsonType: "",
					jsonName: "arrayfield",
					array:    true,
				},
			},
			tree:     []string{"arrayfield"},
			instance: jsonschema.Instance{Type: []string{"string"}},
			want: extractedFields{
				"arrayfield": &extractedField{
					name:     "Arrayfield",
					jsonType: "string",
					jsonName: "arrayfield",
					array:    true,
				},
			},
		},
		{
			description: "Struct field",
			fields:      make(map[string]*extractedField),
			tree:        []string{"structfield"},
			instance:    jsonschema.Instance{Type: []string{"object"}},
			want: extractedFields{
				"structfield": &extractedField{
					name:           "Structfield",
					jsonType:       "object",
					jsonName:       "structfield",
					fields:         make(map[string]*extractedField),
					requiredFields: make(map[string]bool),
				},
			},
		},
		{
			description: "Field in an existing child struct",
			fields: extractedFields{
				"structfield": &extractedField{
					name:           "Structfield",
					jsonType:       "object",
					jsonName:       "structfield",
					fields:         make(map[string]*extractedField),
					requiredFields: make(map[string]bool),
				},
			},
			tree:     []string{"structfield", "child"},
			instance: jsonschema.Instance{Type: []string{"string"}},
			want: extractedFields{
				"structfield": &extractedField{
					name:     "Structfield",
					jsonType: "object",
					jsonName: "structfield",
					fields: map[string]*extractedField{
						"child": {
							name:     "child",
							jsonType: "string",
							jsonName: "Child",
						},
					},
					requiredFields: make(map[string]bool),
				},
			},
		},
	}

	for _, test := range tests {
		if err := addField(test.fields, test.tree, test.instance, nil); err != nil {
			t.Fatalf("Test %q - failed to add fields: %v", test.description, err)
		}

		if len(test.fields) != len(test.want) {
			t.Errorf("Test %q - got %d fields want %d", test.description, len(test.fields), len(test.want))
		}
		for key, field := range test.fields {
			got, want := field, test.want[key]
			if want == nil {
				t.Errorf("Test %q - got value for key %q, want nil", test.description, key)
				continue
			}

			if got.fields != nil && want.fields != nil {
				// Note this is not recursive so I need a better solution if I ever have tests going more than two layers
				if gotf, wantf := len(got.fields), len(want.fields); gotf != wantf {
					t.Errorf("Test %q - got %d child fields, want %d", test.description, gotf, wantf)
				}
				got.fields, want.fields = nil, nil
			}
			if !reflect.DeepEqual(*got, *want) {
				t.Errorf("Test %q - got %#v, want %#v", test.description, *got, *want)
			}
		}
	}
}

func TestExtractedFields_Sorted(t *testing.T) {
	a := &extractedField{name: "A"}
	b := &extractedField{name: "B"}
	c := &extractedField{name: "C"}
	tests := []struct {
		description string
		efs         extractedFields
		want        []*extractedField
	}{
		{
			description: "already sorted",
			efs:         extractedFields{"a": a, "b": b},
			want:        []*extractedField{a, b},
		},
		{
			description: "needs sorting",
			efs:         extractedFields{"C": c, "a": a, "b": b},
			want:        []*extractedField{a, b, c},
		},
	}

	for _, test := range tests {
		got := test.efs.Sorted()
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Test %q - got %v, want %v", test.description, got, test.want)
		}
	}
}

func TestExtractedField_Write(t *testing.T) {
	a := &extractedField{name: "A", jsonName: "a", jsonType: "string", array: true}
	b := &extractedField{name: "B", jsonName: "b", jsonType: "boolean", array: false}
	tests := []struct {
		description            string
		ef                     *extractedField
		prefix                 string
		descriptionAsStructTag bool
		required               bool
		want                   string
	}{
		{
			description: "Write scalr with no prefix",
			ef:          &extractedField{name: "Field", jsonName: "field", jsonType: "number"},
			want:        "Field\tfloat64\t`json:\"field,omitempty\"`\n",
		},
		{
			description: "Write scalr with descrition as comment",
			ef:          &extractedField{name: "Field", jsonName: "field", jsonType: "number", description: "I expect a better description"},
			want:        "// I expect a better description\nField\tfloat64\t`json:\"field,omitempty\"`\n",
		},
		{
			description:            "Write scalr with descrition as struct tag",
			ef:                     &extractedField{name: "Field", jsonName: "field", jsonType: "number", description: "I expect a better description"},
			descriptionAsStructTag: true,
			want:                   "Field\tfloat64\t`json:\"field,omitempty\" description:\"I expect a better description\"`\n",
		},
		{
			description: "Write array scalr with prefix, required",
			ef:          &extractedField{name: "Field", jsonName: "field", jsonType: "number", array: true},
			prefix:      "\t",
			required:    true,
			want:        "\tField\t[]float64\t`json:\"field\"`\n",
		},
		{
			description: "Write struct",
			ef:          &extractedField{name: "Field", jsonName: "field", jsonType: "object", fields: extractedFields{"a": a, "b": b}},
			want:        "Field\tstruct {\n\tA\t[]string\t`json:\"a,omitempty\"`\n\tB\tbool\t`json:\"b,omitempty\"`\n\t}\t`json:\"field,omitempty\"`\n",
		},
		{
			description: "Write struct, required children",
			ef: &extractedField{
				name:           "Field",
				jsonName:       "field",
				jsonType:       "object",
				fields:         extractedFields{"a": a, "b": b},
				requiredFields: map[string]bool{"a": true, "b": true},
			},
			want: "Field\tstruct {\n\tA\t[]string\t`json:\"a\"`\n\tB\tbool\t`json:\"b\"`\n\t}\t`json:\"field,omitempty\"`\n",
		},
	}

	for _, test := range tests {
		buf := &bytes.Buffer{}
		if err := test.ef.write(buf, test.prefix, test.required, test.descriptionAsStructTag, false, nil); err != nil {
			t.Fatalf("Test %q - failed write: %v", test.description, err)
		}
		if got, want := string(buf.Bytes()), test.want; got != want {
			t.Errorf("Test %q - got\n%s\nwant\n%s", test.description, got, want)
		}
	}
}

func TestGeneratedStruct(t *testing.T) {
	tests := []struct {
		description            string
		embeds                 []string
		schemaPath             string
		packageName            string
		oneOfType              string
		descriptionAsStructTag bool
		noNestedStruct         bool
		renameFieldMap         map[string]string
		wantFilePath           string
		wantWriteError         bool
	}{
		{
			description:            "Simple schema - no nest",
			schemaPath:             "struct_test_data/test_schema.json",
			packageName:            "nonest",
			oneOfType:              "simple",
			descriptionAsStructTag: true,
			noNestedStruct:         true,
			wantFilePath:           "struct_test_data/nonest/simple.go",
		},
		{
			description:    "Complex schema - no nest",
			embeds:         []string{"Simple"},
			schemaPath:     "struct_test_data/test_schema.json",
			packageName:    "nonest",
			noNestedStruct: true,
			oneOfType:      "complex",
			wantFilePath:   "struct_test_data/nonest/complex.go",
		},
		{
			description:            "Simple schema",
			schemaPath:             "struct_test_data/test_schema.json",
			packageName:            "test_data",
			oneOfType:              "simple",
			descriptionAsStructTag: true,
			wantFilePath:           "struct_test_data/simple.go",
		},
		{
			description: "Simple schema with field rename",
			schemaPath:  "struct_test_data/test_schema.json",
			packageName: "test_data",
			oneOfType:   "simple",
			renameFieldMap: map[string]string{
				"type":  "typeRenamed",
				"dates": "times",
			},
			wantFilePath: "struct_test_data/simple.go.out-rename-fields",
		},
		{
			description:            "valid schema with invalid go field names should fail",
			schemaPath:             "struct_test_data/needs_field_rename.json",
			packageName:            "test_data",
			oneOfType:              "rename",
			descriptionAsStructTag: true,
			wantWriteError:         true,
		},
		{
			description:            "valid schema with invalid go field names, renamed to work",
			schemaPath:             "struct_test_data/needs_field_rename.json",
			packageName:            "test_data",
			oneOfType:              "rename",
			descriptionAsStructTag: true,
			renameFieldMap: map[string]string{
				"1_1":  "OneToOne",
				"3_4":  "ThreeToFour",
				"4_3":  "FourToThree",
				"16_9": "SixteenToNine",
				"9_16": "NineToSixteen",
			},
			wantFilePath: "struct_test_data/needs_field_rename.go.out",
		},
		{
			description:  "Complex schema",
			embeds:       []string{"Simple"},
			schemaPath:   "struct_test_data/test_schema.json",
			packageName:  "test_data",
			oneOfType:    "complex",
			wantFilePath: "struct_test_data/complex.go",
		},
		{
			description:            "nested array structs schema - no nest",
			schemaPath:             "struct_test_data/nested.json",
			packageName:            "nonest",
			oneOfType:              "nested",
			descriptionAsStructTag: false,
			noNestedStruct:         true,
			wantFilePath:           "struct_test_data/nonest/nested.go",
		},
		{
			description:  "nested array structs schema",
			schemaPath:   "struct_test_data/nested.json",
			packageName:  "test_data",
			oneOfType:    "nested",
			wantFilePath: "struct_test_data/nested.go",
		},
	}

	for _, test := range tests {
		schema, err := jsonschema.SchemaFromFile(test.schemaPath, test.oneOfType)
		if err != nil {
			t.Fatalf("Test %q - SchemaFromFile failed: %v", test.description, err)
		}
		bArgs := BuildArgs{
			DescriptionAsStructTag: test.descriptionAsStructTag,
			FieldNameMap:           test.renameFieldMap,
			NoNestedStructs:        test.noNestedStruct,
		}
		g, err := newGeneratedGoFile(schema, test.oneOfType, test.packageName, test.embeds, bArgs)
		if err != nil {
			t.Fatalf("Test %q - failed: %v", test.description, err)
		}

		buf := &bytes.Buffer{}
		err = g.write(buf)
		if !test.wantWriteError && err != nil {
			t.Fatalf("Test %q - failed write: %v", test.description, err)
		} else if test.wantWriteError && err == nil {
			t.Fatalf("Test %q - expected failure but succeeded to write", test.description)
		} else if test.wantWriteError {
			continue
		}
		got := buf.Bytes()

		want, err := os.ReadFile(test.wantFilePath)
		if err != nil {
			t.Fatalf("Test %q - failed to read result file: %v", test.description, err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Test %q - got\n%s\nwant\n%s", test.description, got, want)
		}
	}
}

func TestGoType(t *testing.T) {
	tests := []struct {
		description string
		jsonType    string
		array       bool
		required    bool
		pointers    bool
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
			want:        "time.Time",
		},
		{
			description: "JSON string date-time array, omitempty",
			jsonType:    "date-time",
			array:       true,
			want:        "[]time.Time",
		},
		{
			description: "JSON string date-time, omitempty, pointers",
			jsonType:    "date-time",
			array:       false,
			pointers:    true,
			want:        "*time.Time",
		},
		{
			description: "JSON string date-time array, omitempty, pointers",
			jsonType:    "date-time",
			array:       true,
			pointers:    true,
			want:        "[]*time.Time",
		},
	}

	for _, test := range tests {
		ef := extractedField{
			array:    test.array,
			jsonType: test.jsonType,
		}
		got := ef.goType(test.required, test.pointers)
		if got != test.want {
			t.Errorf("Test %q - got %q, want %q", test.description, got, test.want)
		}
	}
}
