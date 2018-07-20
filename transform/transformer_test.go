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
	doublearraySchema, _     = jsonschema.SchemaFromFile("./test_data/double-array.json", "")
	operationsSchema, _      = jsonschema.SchemaFromFile("./test_data/operations.json", "")
	dateTimesSchema, _       = jsonschema.SchemaFromFile("./test_data/date-times.json", "")

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
										"height": 0,
										"path": "path",
										"relativePath": "",
										"width": 0
									},
									{
										"name": "aname",
										"height": 0,
										"path": "empty",
										"relativePath": "empty",
										"width": 0
									}
								],
								"publishUrl": "publishURL",
								"absoluteUrl": "absoluteURL"
							}`),
			want: json.RawMessage(`{"URL":{"absolute":"absoluteURL","publish":"publishURL"},"crops":[{"height":0,"name":"name","path":"path","relativePath":"","width":0},{"height":0,"name":"aname","path":"empty","relativePath":"empty","width":0}],"type":"image"}`),
		},
		{
			description: "Input too simple, fails validation",
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
			want:    json.RawMessage(`{"URL":{"absolute":"absoluteURL","publish":"publishURL"},"crops":[{"name":"name","path":"path"},{"name":"aname","relativePath":"empty"}],"type":"image"}`),
			wantErr: true,
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
		{
			description: "Test empty non-required object",
			transformer: Transformer{schema: imageSchema, transformIdentifier: "cumulo"},
			in: json.RawMessage(`
						{
							"type": "image",
							"crops": [
								{
									"height": 0,
									"path": "path",
									"relativePath": "",
									"width": 0
								},
								{
									"name": "aname",
									"height": 0,
									"path": "empty",
									"relativePath": "empty",
									"width": 0
								}
							]
						}`),
			want: json.RawMessage(`{"crops":[{"height":0,"name":"name","path":"path","relativePath":"","width":0},{"height":0,"name":"aname","path":"empty","relativePath":"empty","width":0}],"type":"image"}`),
		},
		{
			description: "Test empty non-required array",
			transformer: Transformer{schema: arrayTransformsSchema, transformIdentifier: "cumulo"},
			in: json.RawMessage(`
						{
							"type": "image",
							"data": {
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
			want: json.RawMessage(`{"lines":["line1","line2"],"wasSingleObject":[{"id":"1","name":"test1"}]}`),
		},
		{
			description: "Test nested arrays",
			transformer: Transformer{schema: doublearraySchema, transformIdentifier: "cumulo"},
			in: json.RawMessage(`
				{
					"data" : {
						"double": [
							["1-1", "1-2"],
							["2-1", "2-2"]
						]
					},
					"array1": [
						{
							"name": "array1-1",
							"array2": [
								{
									"name": "array1-1-1"
								},
								{
									"name": "array1-1-2"
								}
							]
						},
						{
							"name": "array1-2",
							"array2": [
								{
									"name": "array1-2-1"
								}
							]
						}
					]
				}`),
			want: json.RawMessage(`{"array1":[{"array2":[{"level2Name":"array1-1-1"},{"level2Name":"array1-1-2"}],"level1Name":"array1-1"},{"array2":[{"level2Name":"array1-2-1"}],"level1Name":"array1-2"}],"double":[["1-1","1-2"],["2-1","2-2"]]}`),
		},
		{
			description: "Test format: date-time strings",
			transformer: Transformer{schema: dateTimesSchema, transformIdentifier: "cumulo"},
			in: json.RawMessage(`
				{
					"dates": [
						1529958073,
						"2018-06-25T20:21:13Z"
					],
					"requiredDate": "2018-06-25T20:21:13Z",
					"optionalDate": ""
				}`),
			want: json.RawMessage(`{"dates":["2018-06-25T20:21:13Z","2018-06-25T20:21:13Z"],"requiredDate":"2018-06-25T20:21:13Z"}`),
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
			jsonPath:        "$.test1[2]",
			value:           "c",
			wantTransformed: map[string]interface{}{"test1": []interface{}{"a", "b", "c"}},
		},
		{
			description:     "Save value in new Array of objects",
			tr:              Transformer{transformed: map[string]interface{}{}},
			jsonPath:        "$.test1[0].a",
			value:           "aValue",
			wantTransformed: map[string]interface{}{"test1": []interface{}{map[string]interface{}{"a": "aValue"}}},
		},
		{
			description:     "Save value in existing Array of objects",
			tr:              Transformer{transformed: map[string]interface{}{"test1": []interface{}{map[string]interface{}{"a": "aValue"}}}},
			jsonPath:        "$.test1[0].b",
			value:           "bValue",
			wantTransformed: map[string]interface{}{"test1": []interface{}{map[string]interface{}{"a": "aValue", "b": "bValue"}}},
		},
		{
			description:     "Save new array item value in existing Array of objects",
			tr:              Transformer{transformed: map[string]interface{}{"test1": []interface{}{map[string]interface{}{"a": "aValue"}}}},
			jsonPath:        "$.test1[1].a",
			value:           "a2ndValue",
			wantTransformed: map[string]interface{}{"test1": []interface{}{map[string]interface{}{"a": "aValue"}, map[string]interface{}{"a": "a2ndValue"}}},
		},
		{
			description:     "Save new array item simple nested Array",
			tr:              Transformer{transformed: map[string]interface{}{}},
			jsonPath:        "$.test1[0][1]",
			value:           "nestedValue",
			wantTransformed: map[string]interface{}{"test1": []interface{}{[]interface{}{nil, "nestedValue"}}},
		},
		{
			description:     "Save new array item new nested Array",
			tr:              Transformer{transformed: map[string]interface{}{}},
			jsonPath:        "$.test1[0].a[1]",
			value:           "nestedValue",
			wantTransformed: map[string]interface{}{"test1": []interface{}{map[string]interface{}{"a": []interface{}{nil, "nestedValue"}}}},
		},
		{
			description:     "Save object field in new nested Array",
			tr:              Transformer{transformed: map[string]interface{}{}},
			jsonPath:        "$.test1[0].a[1].name",
			value:           "nestedName",
			wantTransformed: map[string]interface{}{"test1": []interface{}{map[string]interface{}{"a": []interface{}{map[string]interface{}{}, map[string]interface{}{"name": "nestedName"}}}}},
		},
		{
			description:     "Save new array item existing nested Array",
			tr:              Transformer{transformed: map[string]interface{}{"test1": []interface{}{map[string]interface{}{"a": []interface{}{"existingValue"}}}}},
			jsonPath:        "$.test1[0].a[1]",
			value:           "nestedValue",
			wantTransformed: map[string]interface{}{"test1": []interface{}{map[string]interface{}{"a": []interface{}{"existingValue", "nestedValue"}}}},
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
		if test.wantErr {
			continue
		}
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
