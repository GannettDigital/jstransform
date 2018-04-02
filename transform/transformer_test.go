package transform

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/GannettDigital/jstransform/jsonschema"
)

// used for the Transformer test and benchmark
var (
	imageSchema, _           = jsonschema.SchemaFromFile("./test_data/image.json", "")
	arrayTransformsSchema, _ = jsonschema.SchemaFromFile("./test_data/array-transforms.json", "")
	operationsSchema, _      = jsonschema.SchemaFromFile("./test_data/operations.json", "")

	transformerTests = []struct {
		description string
		transformer Transformer
		in          json.RawMessage
		want        json.RawMessage
		wantErr     bool
	}{
		{
			description: "Use basic transforms, copy from input and default to build result",
			transformer: Transformer{schema: imageSchema, transformIdentifier: "cumulo"},
			in: json.RawMessage(`
			{
				"type": "image",
				"crops": [
					{
						"path": "path"
					},
					{
						"name": "aname",
						"relativePath": "empty"
					}
				],
				"publishUrl": "publishURL",
				"absoluteUrl": "absoluteURL"
			}`),
			want: json.RawMessage(`{"URL":{"absolute":"absoluteURL","publish":"publishURL"},"crops":[{"name":"name","path":"path"},{"name":"aname","relativePath":"empty"}],"type":"image"}`),
		},
		{
			description: "Array transforms, tests arrays with string type and with a single object type",
			transformer: Transformer{schema: arrayTransformsSchema, transformIdentifier: "cumulo"},
			in: json.RawMessage(`
				{
					"type": "image",
					"data": {
						"contributors": [
							{"id": 1, "fullname": "one"},
							{"id": 2, "fullname": "two"}
						],
						"lines": [
							"line1",
							"line2"
						]
					},
					"aSingleObject": [
						{
							"id": 1,
							"name": "test1"
						}
					]
				}`),
			want: json.RawMessage(`{"contributors":[{"id":"1","name":"one"},{"id":"2","name":"two"}],"lines":["line1","line2"],"wasSingleObject":[{"id":"1","name":"test1"}]}`),
		},
		{
			description: "Test all operations",
			transformer: Transformer{schema: operationsSchema, transformIdentifier: "cumulo"},
			in: json.RawMessage(`
				{
					"type": "image",
					"data": {
						"attributes": [
							{
								"name": "length",
								"value": "00:13"
							}
						],
						"contributors": [
							{"id": 1, "fullname": "one"},
							{"id": 2, "fullname": "two"}
						]
					},
					"mixedCase": "a|B|c|D",
					"invalid": false,
					"url": "http://foo.com/blah"
				}`),
			want: json.RawMessage(`{"caseSplit":["a","b","c","d"],"contributor":"two","duration":13,"url":"http://gannettdigital.com/blah","valid":true}`),
		},
	}
)

func TestSaveValue(t *testing.T) {
	tests := []struct {
		description     string
		tr              Transformer
		jsonPath        string
		value           interface{}
		wantTransformed map[string]interface{}
		wantErr         bool
	}{
		{
			description:     "Simple string value at empty root",
			tr:              Transformer{transformed: make(map[string]interface{})},
			jsonPath:        "$.test",
			value:           "string",
			wantTransformed: map[string]interface{}{"test": "string"},
		},
		{
			description:     "nil value",
			tr:              Transformer{transformed: make(map[string]interface{})},
			jsonPath:        "$.test",
			value:           nil,
			wantTransformed: map[string]interface{}{},
		},
		{
			description:     "Simple string value at existing root",
			tr:              Transformer{transformed: map[string]interface{}{"test1": 1}},
			jsonPath:        "$.test",
			value:           "string",
			wantTransformed: map[string]interface{}{"test": "string", "test1": 1},
		},
		{
			description:     "Simple string value overwriting existing value",
			tr:              Transformer{transformed: map[string]interface{}{"test1": 1}},
			jsonPath:        "$.test1",
			value:           "string",
			wantTransformed: map[string]interface{}{"test1": "string"},
		},
		{
			description:     "Simple string value non-existent parent",
			tr:              Transformer{transformed: map[string]interface{}{}},
			jsonPath:        "$.test1.test2",
			value:           "string",
			wantTransformed: map[string]interface{}{"test1": map[string]interface{}{"test2": "string"}},
		},
		{
			description:     "Simple int value at empty root",
			tr:              Transformer{transformed: make(map[string]interface{})},
			jsonPath:        "$.test",
			value:           1,
			wantTransformed: map[string]interface{}{"test": 1},
		},
		{
			description:     "New Map at empty root",
			tr:              Transformer{transformed: map[string]interface{}{}},
			jsonPath:        "$.test1",
			value:           map[string]interface{}{},
			wantTransformed: map[string]interface{}{"test1": map[string]interface{}{}},
		},
		{
			description:     "Map with values at empty root",
			tr:              Transformer{transformed: map[string]interface{}{}},
			jsonPath:        "$.test1",
			value:           map[string]interface{}{"testA": "a"},
			wantTransformed: map[string]interface{}{"test1": map[string]interface{}{"testA": "a"}},
		},
		{
			description:     "Save new value in existing Map",
			tr:              Transformer{transformed: map[string]interface{}{"test1": map[string]interface{}{"testA": "a"}}},
			jsonPath:        "$.test1.testB",
			value:           "B",
			wantTransformed: map[string]interface{}{"test1": map[string]interface{}{"testA": "a", "testB": "B"}},
		},
		{
			description:     "Array at empty root",
			tr:              Transformer{transformed: map[string]interface{}{}},
			jsonPath:        "$.test1",
			value:           []interface{}{"a", "b"},
			wantTransformed: map[string]interface{}{"test1": []interface{}{"a", "b"}},
		},
		{
			description:     "Save value in existing Array",
			tr:              Transformer{transformed: map[string]interface{}{"test1": []interface{}{"a", "b"}}},
			jsonPath:        "$.test1[3]",
			value:           "c",
			wantTransformed: map[string]interface{}{"test1": []interface{}{"a", "b", "c"}},
		},
	}

	for _, test := range tests {
		err := test.tr.saveValue(test.jsonPath, test.value)

		switch {
		case test.wantErr && err != nil:
			continue
		case test.wantErr && err == nil:
			t.Errorf("Test %q - got nil, want error", test.description)
		case !test.wantErr && err != nil:
			t.Errorf("Test %q - got error, want nil: %v", test.description, err)
		case !reflect.DeepEqual(test.tr.transformed, test.wantTransformed):
			t.Errorf("Test %q - got %v, want %v", test.description, test.tr.transformed, test.wantTransformed)
		}
	}
}

func TestTransformer(t *testing.T) {
	for _, test := range transformerTests {
		got, err := test.transformer.Transform(test.in)

		switch {
		case test.wantErr && err != nil:
			continue
		case test.wantErr && err == nil:
			t.Errorf("Test %q - got nil, want error", test.description)
		case !test.wantErr && err != nil:
			t.Errorf("Test %q - got error, want nil: %v", test.description, err)
		case !reflect.DeepEqual(got, test.want):
			t.Errorf("Test %q - got\n%s\nwant\n%s", test.description, got, test.want)
		}
	}
}

func BenchmarkTransformer(b *testing.B) {
	for _, test := range transformerTests {
		b.Run(test.description, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := test.transformer.Transform(test.in)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
