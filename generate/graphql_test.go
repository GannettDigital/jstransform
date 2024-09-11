package generate

// For now a direct copy of 'struct_test.go' tests.
import (
	"bytes"
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/GannettDigital/jstransform/jsonschema"
)

func TestGraphQLAddField(t *testing.T) {
	tests := []struct {
		description string
		fields      gqlExtractedFields
		tree        []string
		instance    jsonschema.Instance
		want        gqlExtractedFields
	}{
		{
			description: "Simple scalar field",
			fields:      make(map[string]*gqlExtractedField),
			tree:        []string{"field"},
			instance:    jsonschema.Instance{Type: []string{"string"}},
			want: gqlExtractedFields{
				"field": &gqlExtractedField{
					name:     "Field",
					jsonType: "string",
					jsonName: "field",
				},
			},
		},
		{
			description: "Array field",
			fields:      make(map[string]*gqlExtractedField),
			tree:        []string{"arrayfield"},
			instance:    jsonschema.Instance{Type: []string{"array"}, Items: []byte(`{ "type": "string" }`)},
			want: gqlExtractedFields{
				"arrayfield": &gqlExtractedField{
					name:     "Arrayfield",
					jsonType: "",
					jsonName: "arrayfield",
					array:    true,
				},
				"totalArrayfield": &gqlExtractedField{
					description: "The total length of the arrayfield list at this same level in the data, this number is unaffected by filtering.",
					name:        "TotalArrayfield",
					jsonType:    "integer",
					jsonName:    "totalArrayfield",
				},
			},
		},
		{
			description: "Array field 2nd call",
			fields: gqlExtractedFields{
				"arrayfield": &gqlExtractedField{
					name:     "Arrayfield",
					jsonType: "",
					jsonName: "arrayfield",
					array:    true,
				},
			},
			tree:     []string{"arrayfield"},
			instance: jsonschema.Instance{Type: []string{"string"}},
			want: gqlExtractedFields{
				"arrayfield": &gqlExtractedField{
					name:     "Arrayfield",
					jsonType: "string",
					jsonName: "arrayfield",
					array:    true,
				},
			},
		},
		{
			description: "Struct field",
			fields:      make(map[string]*gqlExtractedField),
			tree:        []string{"structfield"},
			instance: jsonschema.Instance{
				Type:       []string{"object"},
				Properties: map[string]json.RawMessage{"test": nil},
			},
			want: gqlExtractedFields{
				"structfield": &gqlExtractedField{
					name:           "Structfield",
					jsonType:       "object",
					jsonName:       "structfield",
					fields:         make(map[string]*gqlExtractedField),
					requiredFields: make(map[string]bool),
				},
			},
		},
		{
			description: "Field in an existing child struct",
			fields: gqlExtractedFields{
				"structfield": &gqlExtractedField{
					name:           "Structfield",
					jsonType:       "object",
					jsonName:       "structfield",
					fields:         make(map[string]*gqlExtractedField),
					requiredFields: make(map[string]bool),
				},
			},
			tree:     []string{"structfield", "child"},
			instance: jsonschema.Instance{Type: []string{"string"}},
			want: gqlExtractedFields{
				"structfield": &gqlExtractedField{
					name:       "Structfield",
					jsonType:   "object",
					jsonName:   "structfield",
					fieldOrder: []string{"child"},
					fields: map[string]*gqlExtractedField{
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
		ef := gqlExtractedField{
			fields:         test.fields,
			requiredFields: make(map[string]bool),
		}
		if err := ef.addField(test.tree, nil, test.instance); err != nil {
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
				t.Errorf("Test %q\nwant: %#v\ngot:  %#v", test.description, *want, *got)
			}
		}
	}
}

func TestGraphQLExtractedFields_Sorted(t *testing.T) {
	a := &gqlExtractedField{name: "A"}
	b := &gqlExtractedField{name: "B"}
	c := &gqlExtractedField{name: "C"}
	tests := []struct {
		description string
		efs         gqlExtractedFields
		want        []*gqlExtractedField
	}{
		{
			description: "already sorted",
			efs:         gqlExtractedFields{"a": a, "b": b},
			want:        []*gqlExtractedField{a, b},
		},
		{
			description: "needs sorting",
			efs:         gqlExtractedFields{"C": c, "a": a, "b": b},
			want:        []*gqlExtractedField{a, b, c},
		},
	}

	for _, test := range tests {
		got := test.efs.Sorted()
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Test %q - got %v, want %v", test.description, got, test.want)
		}
	}
}

func TestGraphQLExtractedField_Write(t *testing.T) {
	/*
		a := &gqlExtractedField{name: "A", jsonName: "a", jsonType: "string", array: true}
		b := &gqlExtractedField{name: "B", jsonName: "b", jsonType: "boolean", array: false}
	*/
	tests := []struct {
		description            string
		ef                     *gqlExtractedField
		prefix                 string
		descriptionAsStructTag bool
		required               bool
		want                   string
	}{
		{
			description: "Write scalr with no prefix",
			ef:          &gqlExtractedField{name: "Field", jsonName: "field", jsonType: "number"},
			want:        "field: Float!\n",
		},
		{
			description: "Write scalr with descrition as comment",
			ef:          &gqlExtractedField{name: "Field", jsonName: "field", jsonType: "number", description: "I expect a better description"},
			want:        "\"I expect a better description\"\nfield: Float!\n",
		},
		{
			description:            "Write scalr with descrition as struct tag",
			ef:                     &gqlExtractedField{name: "Field", jsonName: "field", jsonType: "number", description: "I expect a better description"},
			descriptionAsStructTag: true,
			want:                   "\"I expect a better description\"\nfield: Float!\n",
		},
		{
			description: "Write array scalr with prefix, required",
			ef:          &gqlExtractedField{name: "Field", jsonName: "field", jsonType: "number", array: true},
			prefix:      "\t",
			required:    true,
			want:        "\tfield(\n    \"A List Filter expression such as '{Field: \\\"position\\\", Operation: \\\"<=\\\", Argument: {Value: 10}}'\"\n    filter: ListFilter\n\n    \"Sort the list, ie '{Field: \\\"position\\\", Order: \\\"ASC\\\"}'\"\n    sort: ListSortParams\n  ): [Float!]! @goField(forceResolver: true)\n",
		},
		/* GraphQL generation never nests structures.
		   		{
		   			description: "Write struct",
		   			ef: &gqlExtractedField{
		                   name: "Field",
		                   jsonName: "field",
		                   jsonType: "object",
		                   fields: gqlExtractedFields{"a": a, "b": b},
		               },
		   			want:        "",
		   		},
		   		{
		   			description: "Write struct, required children",
		   			ef: &gqlExtractedField{
		   				name:           "Field",
		   				jsonName:       "field",
		   				jsonType:       "object",
		   				fields:         gqlExtractedFields{"a": a, "b": b},
		   				requiredFields: map[string]bool{"a": true, "b": true},
		   			},
		   			want: "",
		   		},
		*/
	}

	for _, test := range tests {
		buf := &bytes.Buffer{}
		if err := test.ef.write(buf, test.prefix, test.required, test.descriptionAsStructTag, false); err != nil {
			t.Fatalf("Test %q - failed write: %v", test.description, err)
		}
		if got, want := string(buf.Bytes()), test.want; got != want {
			t.Errorf("Test %q\nwant: %s\ngot:  %s", test.description, want, got)
		}
	}
}

// TODO: This isn't a full test of the GraphQL generation at the moment because some logic is in `buildGraphQLFile()`.
func TestGraphQLGeneratedStruct(t *testing.T) {
	tests := []struct {
		description            string
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
			schemaPath:             "graphql_test_data/test_schema.json",
			packageName:            "nonest",
			oneOfType:              "simple",
			descriptionAsStructTag: true,
			noNestedStruct:         true,
			wantFilePath:           "graphql_test_data/nonest/simple.graphqls",
		},
		{
			description:    "Complex schema - no nest",
			schemaPath:     "graphql_test_data/test_schema.json",
			packageName:    "nonest",
			noNestedStruct: true,
			oneOfType:      "complex",
			wantFilePath:   "graphql_test_data/nonest/complex.graphqls",
		},
		{
			description:            "Simple schema",
			schemaPath:             "graphql_test_data/test_schema.json",
			packageName:            "test_data",
			oneOfType:              "simple",
			descriptionAsStructTag: true,
			wantFilePath:           "graphql_test_data/simple.graphqls",
		},
		{
			description: "Simple schema with field rename",
			schemaPath:  "graphql_test_data/test_schema.json",
			packageName: "test_data",
			oneOfType:   "simple",
			renameFieldMap: map[string]string{
				"type":  "typeRenamed",
				"dates": "times",
			},
			wantFilePath: "graphql_test_data/simple.graphqls.out-rename-fields",
		},
		{
			description:            "valid schema with numeric graphql field names, renamed to words",
			schemaPath:             "graphql_test_data/needs_field_rename.json",
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
			wantFilePath: "graphql_test_data/needs_field_rename.graphqls.out",
		},
		{
			description:  "Complex schema",
			schemaPath:   "graphql_test_data/test_schema.json",
			packageName:  "test_data",
			oneOfType:    "complex",
			wantFilePath: "graphql_test_data/complex.graphqls",
		},
		{
			description:            "nested array structs schema - no nest",
			schemaPath:             "graphql_test_data/nested.json",
			packageName:            "nonest",
			oneOfType:              "nested",
			descriptionAsStructTag: false,
			noNestedStruct:         true,
			wantFilePath:           "graphql_test_data/nonest/nested.graphqls",
		},
		{
			description:  "nested array structs schema",
			schemaPath:   "graphql_test_data/nested.json",
			packageName:  "test_data",
			oneOfType:    "nested",
			wantFilePath: "graphql_test_data/nested.graphqls",
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
		g, err := newGeneratedGraphQLFile(schema.Instance, test.oneOfType, test.packageName, false, bArgs)
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
			_ = os.WriteFile(test.wantFilePath+".got", got, 0600)
			t.Errorf("Test %q\nwant: %s\ngot:  %s", test.description, want, got)
			t.Errorf("Test %q\nwant: %v\ngot:  %v", test.description, want, got)
			lwant := strings.Split(string(want), "\n")
			lgot := strings.Split(string(got), "\n")
			for idx := range lwant {
				if idx <= len(lgot)-1 && lwant[idx] != lgot[idx] {
					t.Logf("line %d\nwant: %s\ngot:  %s", idx, lwant[idx], lgot[idx])
					t.Logf("line %d\nwant: %v\ngot:  %v", idx, []byte(lwant[idx]), []byte(lgot[idx]))
				}
			}
		} else {
			_ = os.Remove(test.wantFilePath + ".got")
		}
	}
}

func TestGraphQLType(t *testing.T) {
	tests := []struct {
		description string
		jsonType    string
		array       bool
		required    bool
		pointers    bool
		wantArgs    string
		wantType    string
	}{
		{
			description: "JSON boolean",
			jsonType:    "boolean",
			wantType:    "Boolean!",
		},
		{
			description: "JSON boolean array",
			jsonType:    "boolean",
			array:       true,
			wantArgs:    "(\n    \"A List Filter expression such as '{Field: \\\"position\\\", Operation: \\\"<=\\\", Argument: {Value: 10}}'\"\n    filter: ListFilter\n\n    \"Sort the list, ie '{Field: \\\"position\\\", Order: \\\"ASC\\\"}'\"\n    sort: ListSortParams\n  )",
			wantType:    "[Boolean!]! @goField(forceResolver: true)",
		},
		{
			description: "JSON integer",
			jsonType:    "integer",
			wantType:    "Int!",
		},
		{
			description: "JSON integer array",
			jsonType:    "integer",
			array:       true,
			wantArgs:    "(\n    \"A List Filter expression such as '{Field: \\\"position\\\", Operation: \\\"<=\\\", Argument: {Value: 10}}'\"\n    filter: ListFilter\n\n    \"Sort the list, ie '{Field: \\\"position\\\", Order: \\\"ASC\\\"}'\"\n    sort: ListSortParams\n  )",
			wantType:    "[Int!]! @goField(forceResolver: true)",
		},
		{
			description: "JSON number",
			jsonType:    "number",
			wantType:    "Float!",
		},
		{
			description: "JSON number array",
			jsonType:    "number",
			array:       true,
			wantArgs:    "(\n    \"A List Filter expression such as '{Field: \\\"position\\\", Operation: \\\"<=\\\", Argument: {Value: 10}}'\"\n    filter: ListFilter\n\n    \"Sort the list, ie '{Field: \\\"position\\\", Order: \\\"ASC\\\"}'\"\n    sort: ListSortParams\n  )",
			wantType:    "[Float!]! @goField(forceResolver: true)",
		},
		{
			description: "JSON string",
			jsonType:    "string",
			wantType:    "String!",
		},
		{
			description: "JSON string array",
			jsonType:    "string",
			array:       true,
			wantArgs:    "(\n    \"A List Filter expression such as '{Field: \\\"position\\\", Operation: \\\"<=\\\", Argument: {Value: 10}}'\"\n    filter: ListFilter\n\n    \"Sort the list, ie '{Field: \\\"position\\\", Order: \\\"ASC\\\"}'\"\n    sort: ListSortParams\n  )",
			wantType:    "[String!]! @goField(forceResolver: true)",
		},
		{
			description: "JSON object",
			jsonType:    "object",
			wantType:    "(!", // Case not expected to occur.
		},
		{
			description: "JSON object",
			jsonType:    "object",
			array:       true,
			wantArgs:    "(\n    \"A List Filter expression such as '{Field: \\\"position\\\", Operation: \\\"<=\\\", Argument: {Value: 10}}'\"\n    filter: ListFilter\n\n    \"Sort the list, ie '{Field: \\\"position\\\", Order: \\\"ASC\\\"}'\"\n    sort: ListSortParams\n  )",
			wantType:    "[(!]! @goField(forceResolver: true)", // Case not expected to occur.
		},
		{
			description: "JSON string date-time",
			jsonType:    "date-time",
			array:       false,
			required:    true,
			wantType:    "DateTime",
		},
		{
			description: "JSON string date-time array",
			jsonType:    "date-time",
			array:       true,
			required:    true,
			wantArgs:    "(\n    \"A List Filter expression such as '{Field: \\\"position\\\", Operation: \\\"<=\\\", Argument: {Value: 10}}'\"\n    filter: ListFilter\n\n    \"Sort the list, ie '{Field: \\\"position\\\", Order: \\\"ASC\\\"}'\"\n    sort: ListSortParams\n  )",
			wantType:    "[DateTime]! @goField(forceResolver: true)",
		},
		{
			description: "JSON string date-time, omitempty",
			jsonType:    "date-time",
			array:       false,
			wantType:    "DateTime",
		},
		{
			description: "JSON string date-time array, omitempty",
			jsonType:    "date-time",
			array:       true,
			wantArgs:    "(\n    \"A List Filter expression such as '{Field: \\\"position\\\", Operation: \\\"<=\\\", Argument: {Value: 10}}'\"\n    filter: ListFilter\n\n    \"Sort the list, ie '{Field: \\\"position\\\", Order: \\\"ASC\\\"}'\"\n    sort: ListSortParams\n  )",
			wantType:    "[DateTime]! @goField(forceResolver: true)",
		},
		{
			description: "JSON string date-time, omitempty, pointers",
			jsonType:    "date-time",
			array:       false,
			pointers:    true,
			wantType:    "DateTime",
		},
		{
			description: "JSON string date-time array, omitempty, pointers",
			jsonType:    "date-time",
			array:       true,
			pointers:    true,
			wantArgs:    "(\n    \"A List Filter expression such as '{Field: \\\"position\\\", Operation: \\\"<=\\\", Argument: {Value: 10}}'\"\n    filter: ListFilter\n\n    \"Sort the list, ie '{Field: \\\"position\\\", Order: \\\"ASC\\\"}'\"\n    sort: ListSortParams\n  )",
			wantType:    "[DateTime] @goField(forceResolver: true)",
		},
	}

	for _, test := range tests {
		ef := gqlExtractedField{
			array:    test.array,
			jsonType: test.jsonType,
		}
		gotArgs, gotType := ef.graphqlType(test.required, test.pointers)
		if gotArgs != test.wantArgs {
			t.Errorf("Test %q arguments\nwant: %q\ngot:  %q", test.description, test.wantArgs, gotArgs)
		}
		if gotType != test.wantType {
			t.Errorf("Test %q type\nwant: %q\ngot:  %q", test.description, test.wantType, gotType)
		}
	}
}
