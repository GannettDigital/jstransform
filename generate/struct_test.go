package generate

import (
	"bytes"
	"io/ioutil"
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
			instance:    jsonschema.Instance{Type: "string"},
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
			instance:    jsonschema.Instance{Type: "array", Items: []byte(`{ "type": "string" }`)},
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
			instance: jsonschema.Instance{Type: "string"},
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
			instance:    jsonschema.Instance{Type: "object"},
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
			instance: jsonschema.Instance{Type: "string"},
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
		if err := test.ef.write(buf, test.prefix, test.required, test.descriptionAsStructTag); err != nil {
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
			schemaPath:             "test_data/test_schema.json",
			packageName:            "nonestTest",
			oneOfType:              "simple",
			descriptionAsStructTag: true,
			noNestedStruct:         true,
			wantFilePath:           "test_data/nonest/simple.go",
		},
		{
			description:    "Complex schema - no nest",
			embeds:         []string{"Simple"},
			schemaPath:     "test_data/test_schema.json",
			packageName:    "nonestTest",
			noNestedStruct: true,
			oneOfType:      "complex",
			wantFilePath:   "test_data/nonest/complex.go",
		},
		{
			description:            "Simple schema",
			schemaPath:             "test_data/test_schema.json",
			packageName:            "test_data",
			oneOfType:              "simple",
			descriptionAsStructTag: true,
			wantFilePath:           "test_data/simple.go",
		},
		{
			description: "Simple schema with field rename",
			schemaPath:  "test_data/test_schema.json",
			packageName: "test_data",
			oneOfType:   "simple",
			renameFieldMap: map[string]string{
				"type":  "typeRenamed",
				"dates": "times",
			},
			wantFilePath: "test_data/simple.go.out-rename-fields",
		},
		{
			description:            "valid schema with invalid go field names should fail",
			schemaPath:             "test_data/needs_field_rename.json",
			packageName:            "test_data",
			oneOfType:              "rename",
			descriptionAsStructTag: true,
			wantWriteError:         true,
		},
		{
			description:            "valid schema with invalid go field names, renamed to work",
			schemaPath:             "test_data/needs_field_rename.json",
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
			wantFilePath: "test_data/needs_field_rename.go.out",
		},
		{
			description:  "Complex schema",
			embeds:       []string{"Simple"},
			schemaPath:   "test_data/test_schema.json",
			packageName:  "test_data",
			oneOfType:    "complex",
			wantFilePath: "test_data/complex.go",
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
			t.Fatalf("Test %q - expected failure but succeded to write", test.description)
		} else if test.wantWriteError {
			continue
		}
		got := buf.Bytes()

		want, err := ioutil.ReadFile(test.wantFilePath)
		if err != nil {
			t.Fatalf("Test %q - failed to read result file: %v", test.description, err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Test %q - got\n%s\nwant\n%s", test.description, got, want)
		}
	}
}
