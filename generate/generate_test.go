package generate

import (
	"bytes"
	"io/ioutil"
	"reflect"
	"testing"

	"fmt"
	"os/exec"

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
		if err := addField(test.fields, test.tree, test.instance); err != nil {
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

func TestBuildStructs(t *testing.T) {
	testdir := "test_data/buildstructs"

	if err := BuildStructs(fmt.Sprintf("%s/schema.json", testdir), testdir); err != nil {
		t.Fatalf("failed BuildStructs: %v", err)
	}

	cmd := exec.Command("git", "diff", "--quiet", testdir)
	if err := cmd.Run(); err != nil {
		t.Errorf("Build Structs run found differences, this left %q in a modified state: %v", testdir, err)
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
		description string
		ef          *extractedField
		prefix      string
		required    bool
		want        string
	}{
		{
			description: "Write scalr with no prefix",
			ef:          &extractedField{name: "Field", jsonName: "field", jsonType: "number"},
			want:        "Field\tfloat64\t`json:\"field,omitempty\"`\n",
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
		if err := test.ef.write(buf, test.prefix, test.required); err != nil {
			t.Fatalf("Test %q - failed write: %v", test.description, err)
		}
		if got, want := string(buf.Bytes()), test.want; got != want {
			t.Errorf("Test %q - got\n%s\nwant\n%s", test.description, got, want)
		}
	}
}

func TestGenerateGoStruct(t *testing.T) {
	tests := []struct {
		description  string
		schemaPath   string
		packageName  string
		oneOfType    string
		wantFilePath string
	}{
		{
			description:  "Simple schema",
			schemaPath:   "test_data/test_schema.json",
			packageName:  "test",
			oneOfType:    "simple",
			wantFilePath: "test_data/simple.go",
		},
		{
			description:  "Complex schema",
			schemaPath:   "test_data/test_schema.json",
			packageName:  "test",
			oneOfType:    "complex",
			wantFilePath: "test_data/complex.go",
		},
	}

	for _, test := range tests {
		buf := &bytes.Buffer{}
		if err := generateGoStruct(test.schemaPath, test.packageName, test.oneOfType, buf); err != nil {
			t.Fatalf("Test %q - failed: %v", test.description, err)
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
