package transform

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	frontSchema, _           = jsonschema.SchemaFromFile("./test_data/front.json", "")

	transformerTests = []struct {
		description         string
		schema              *jsonschema.Schema
		transformIdentifier string
		in                  json.RawMessage
		want                json.RawMessage
		wantErr             bool
	}{
		{
			description:         "Use basic transforms, copy from input and default to build result",
			schema:              imageSchema,
			transformIdentifier: "cumulo",
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
			description:         "Input too simple, fails validation",
			schema:              imageSchema,
			transformIdentifier: "cumulo",
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
			description:         "Array transforms, tests arrays with string type and with a single object type",
			schema:              arrayTransformsSchema,
			transformIdentifier: "cumulo",
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
			description:         "Test all operations",
			schema:              operationsSchema,
			transformIdentifier: "cumulo",
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
							"url": "http://foo.com/blah",
							"startTime": "2019-05-16T21:00:00-04:00"
						}`),
			want: json.RawMessage(`{"caseSplit":["a","b","c","d"],"contributor":"two","duration":13,"startTime":"09:00","url":"http://gannettdigital.com/blah","valid":true}`),
		},
		{
			description:         "Test empty non-required object",
			schema:              imageSchema,
			transformIdentifier: "cumulo",
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
			description:         "Test empty non-required array",
			schema:              arrayTransformsSchema,
			transformIdentifier: "cumulo",
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
			description:         "Test nested arrays",
			schema:              doublearraySchema,
			transformIdentifier: "cumulo",
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
			description:         "Test format: date-time strings",
			schema:              dateTimesSchema,
			transformIdentifier: "cumulo",
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
		{
			description:         "Test special characters",
			schema:              frontSchema,
			transformIdentifier: "frontInput",
			in: json.RawMessage(`
						{
							"attributes": [
								{
									"canonicalurl": "canURL",
									"front-list-module-position": "frontlistmoduleposition"
								}
							],
							"og:image": "testOGIMAGE"
						}`),
			want: json.RawMessage(`{"attributes":[{"canonicalURL":"canURL","frontListModulePosition":"frontlistmoduleposition"}],"ogImage":"testOGIMAGE"}`),
		},
	}

	saveValueTests = []struct {
		description string
		tree        map[string]interface{}
		jsonPath    string
		value       interface{}
		want        map[string]interface{}
		wantErr     bool
	}{
		{
			description: "Simple string value at empty root",
			tree:        make(map[string]interface{}),
			jsonPath:    "$.test",
			value:       "string",
			want:        map[string]interface{}{"test": "string"},
		},
		{
			description: "nil value",
			tree:        make(map[string]interface{}),
			jsonPath:    "$.test",
			value:       nil,
			want:        map[string]interface{}{},
		},
		{
			description: "Simple string value at existing root",
			tree:        map[string]interface{}{"test1": 1},
			jsonPath:    "$.test",
			value:       "string",
			want:        map[string]interface{}{"test": "string", "test1": 1},
		},
		{
			description: "Simple string value overwriting existing value",
			tree:        map[string]interface{}{"test1": 1},
			jsonPath:    "$.test1",
			value:       "string",
			want:        map[string]interface{}{"test1": "string"},
		},
		{
			description: "Simple string value non-existent parent",
			tree:        map[string]interface{}{},
			jsonPath:    "$.test1.test2",
			value:       "string",
			want:        map[string]interface{}{"test1": map[string]interface{}{"test2": "string"}},
		},
		{
			description: "Simple int value at empty root",
			tree:        make(map[string]interface{}),
			jsonPath:    "$.test",
			value:       1,
			want:        map[string]interface{}{"test": 1},
		},
		{
			description: "New Map at empty root",
			tree:        map[string]interface{}{},
			jsonPath:    "$.test1",
			value:       map[string]interface{}{},
			want:        map[string]interface{}{"test1": map[string]interface{}{}},
		},
		{
			description: "Map with values at empty root",
			tree:        map[string]interface{}{},
			jsonPath:    "$.test1",
			value:       map[string]interface{}{"testA": "a"},
			want:        map[string]interface{}{"test1": map[string]interface{}{"testA": "a"}},
		},
		{
			description: "Save new value in existing Map",
			tree:        map[string]interface{}{"test1": map[string]interface{}{"testA": "a"}},
			jsonPath:    "$.test1.testB",
			value:       "B",
			want:        map[string]interface{}{"test1": map[string]interface{}{"testA": "a", "testB": "B"}},
		},
		{
			description: "Array at empty root",
			tree:        map[string]interface{}{},
			jsonPath:    "$.test1",
			value:       []interface{}{"a", "b"},
			want:        map[string]interface{}{"test1": []interface{}{"a", "b"}},
		},
		{
			description: "Save value in existing Array",
			tree:        map[string]interface{}{"test1": []interface{}{"a", "b"}},
			jsonPath:    "$.test1[2]",
			value:       "c",
			want:        map[string]interface{}{"test1": []interface{}{"a", "b", "c"}},
		},
		{
			description: "Save value in new Array of objects",
			tree:        map[string]interface{}{},
			jsonPath:    "$.test1[0].a",
			value:       "aValue",
			want:        map[string]interface{}{"test1": []interface{}{map[string]interface{}{"a": "aValue"}}},
		},
		{
			description: "Save value in existing Array of objects",
			tree:        map[string]interface{}{"test1": []interface{}{map[string]interface{}{"a": "aValue"}}},
			jsonPath:    "$.test1[0].b",
			value:       "bValue",
			want:        map[string]interface{}{"test1": []interface{}{map[string]interface{}{"a": "aValue", "b": "bValue"}}},
		},
		{
			description: "Save new array item value in existing Array of objects",
			tree:        map[string]interface{}{"test1": []interface{}{map[string]interface{}{"a": "aValue"}}},
			jsonPath:    "$.test1[1].a",
			value:       "a2ndValue",
			want:        map[string]interface{}{"test1": []interface{}{map[string]interface{}{"a": "aValue"}, map[string]interface{}{"a": "a2ndValue"}}},
		},
		{
			description: "Save new array item simple nested Array",
			tree:        map[string]interface{}{},
			jsonPath:    "$.test1[0][1]",
			value:       "nestedValue",
			want:        map[string]interface{}{"test1": []interface{}{[]interface{}{nil, "nestedValue"}}},
		},
		{
			description: "Save new array item new nested Array",
			tree:        map[string]interface{}{},
			jsonPath:    "$.test1[0].a[1]",
			value:       "nestedValue",
			want:        map[string]interface{}{"test1": []interface{}{map[string]interface{}{"a": []interface{}{nil, "nestedValue"}}}},
		},
		{
			description: "Save object field in new nested Array",
			tree:        map[string]interface{}{},
			jsonPath:    "$.test1[0].a[1].name",
			value:       "nestedName",
			want:        map[string]interface{}{"test1": []interface{}{map[string]interface{}{"a": []interface{}{nil, map[string]interface{}{"name": "nestedName"}}}}},
		},
		{
			description: "Save new array item existing nested Array",
			tree:        map[string]interface{}{"test1": []interface{}{map[string]interface{}{"a": []interface{}{"existingValue"}}}},
			jsonPath:    "$.test1[0].a[1]",
			value:       "nestedValue",
			want:        map[string]interface{}{"test1": []interface{}{map[string]interface{}{"a": []interface{}{"existingValue", "nestedValue"}}}},
		},
	}
)

func TestSaveValue(t *testing.T) {
	for _, test := range saveValueTests {
		err := saveInTree(test.tree, test.jsonPath, test.value)

		switch {
		case test.wantErr && err != nil:
			continue
		case test.wantErr && err == nil:
			t.Errorf("Test %q - got nil, want error", test.description)
		case !test.wantErr && err != nil:
			t.Errorf("Test %q - got error, want nil: %v", test.description, err)
		case !reflect.DeepEqual(test.tree, test.want):
			t.Errorf("Test %q - got %v, want %v", test.description, test.tree, test.want)
		}
	}
}

func TestTransformer(t *testing.T) {
	const parallelRuns = 4
	for _, test := range transformerTests {
		tr, err := NewTransformer(test.schema, test.transformIdentifier)
		if err != nil {
			t.Fatalf("Test %q - failed to initialize transformer: %v", test.description, err)
		}
		testFunc := func(description string, in json.RawMessage, wantErr bool, want json.RawMessage) func(t *testing.T) {
			return func(t *testing.T) {
				t.Parallel()
				got, err := tr.Transform(in)

				switch {
				case wantErr && err != nil:
					return
				case wantErr && err == nil:
					t.Errorf("Test %q - got nil, want error", description)
				case !wantErr && err != nil:
					t.Errorf("Test %q - got error, want nil: %v", description, err)
				case !reflect.DeepEqual(got, want):
					t.Errorf("Test %q - got\n%s\nwant\n%s", description, got, want)
				}
			}
		}
		for i := 0; i < parallelRuns; i++ {
			t.Run(fmt.Sprintf("%s-%d", test.description, i), testFunc(test.description, test.in, test.wantErr, test.want))
		}
	}
}

func TestNewXMLTransformer(t *testing.T) {
	tests := []struct {
		description         string
		transformIdentifier string
		schemaFilePath      string
		xmlFilePath         string
		wantFilePath        string
	}{
		{
			description:         "teams NBA",
			transformIdentifier: "sport",
			schemaFilePath:      "./test_data/xml/sports/teams/teams.json",
			xmlFilePath:         "./test_data/xml/sports/teams/teams_NBA.xml",
			wantFilePath:        "./test_data/xml/sports/teams/teamsNBA.out.json",
		},
		{
			description:         "teams MLB",
			transformIdentifier: "sport",
			schemaFilePath:      "./test_data/xml/sports/teams/teams.json",
			xmlFilePath:         "./test_data/xml/sports/teams/teams_MLB.xml",
			wantFilePath:        "./test_data/xml/sports/teams/teamsMLB.out.json",
		},
		{
			description:         "array-transforms",
			transformIdentifier: "sport",
			schemaFilePath:      "./test_data/xml/array-transforms.json",
			xmlFilePath:         "./test_data/xml/array-transforms.xml",
			wantFilePath:        "./test_data/xml/array-transforms.out.json",
		},
		{
			description:         "multiple-array-transforms",
			transformIdentifier: "sport",
			schemaFilePath:      "./test_data/xml/multiple-arrays.json",
			xmlFilePath:         "./test_data/xml/multiple-arrays.xml",
			wantFilePath:        "./test_data/xml/multiple-arrays.out.json",
		},
		{
			description:         "conversion-transforms",
			transformIdentifier: "sport",
			schemaFilePath:      "./test_data/xml/conversion-transforms.json",
			xmlFilePath:         "./test_data/xml/conversion-transforms.xml",
			wantFilePath:        "./test_data/xml/conversion-transforms.out.json",
		},
		{
			description:         "attribute-selection",
			transformIdentifier: "sport",
			schemaFilePath:      "./test_data/xml/attribute-selection.json",
			xmlFilePath:         "./test_data/xml/attribute-selection.xml",
			wantFilePath:        "./test_data/xml/attribute-selection.out.json",
		},
		{
			description:         "operations",
			transformIdentifier: "sport",
			schemaFilePath:      "./test_data/xml/operations.json",
			xmlFilePath:         "./test_data/xml/operations.xml",
			wantFilePath:        "./test_data/xml/operations.out.json",
		},
		{
			description:         "singleArrayElement",
			transformIdentifier: "sport",
			schemaFilePath:      "./test_data/xml/singleArrayElement.json",
			xmlFilePath:         "./test_data/xml/singleArrayElement.xml",
			wantFilePath:        "./test_data/xml/singleArrayElement.out.json",
		},
		{
			description:         "xmlRefsTest",
			transformIdentifier: "sport",
			schemaFilePath:      "./test_data/xml/xmlRefsTest/baseballBoxscores.json",
			xmlFilePath:         "./test_data/xml/xmlRefsTest/boxscoreBaseball.xml",
			wantFilePath:        "./test_data/xml/xmlRefsTest/boxscores.out.json",
		},
	}

	for _, test := range tests {
		schema, err := jsonschema.SchemaFromFile(test.schemaFilePath, "")
		if err != nil {
			t.Fatal(err)
		}

		tr, err := NewXMLTransformer(schema, test.transformIdentifier)
		if err != nil {
			t.Fatal(err)
		}

		rawXMLBytes, err := ioutil.ReadFile(test.xmlFilePath)
		if err != nil {
			t.Fatal(err)
		}

		output, err := tr.Transform(rawXMLBytes)
		if err != nil {
			t.Fatal(err)
		}

		want, err := ioutil.ReadFile(test.wantFilePath)

		var (
			outputMap map[string]interface{}
			wantMap   map[string]interface{}
		)

		if err := json.Unmarshal(output, &outputMap); err != nil {
			t.Fatal(err)
		}

		if err := json.Unmarshal(want, &wantMap); err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(outputMap, wantMap) {
			t.Fatalf("test %s failed \n got:\n %s \n want:\n %s", test.description, output, want)
		}
	}

}

func BenchmarkTransformer(b *testing.B) {
	for _, test := range transformerTests {
		if test.wantErr {
			continue
		}

		tr, err := NewTransformer(test.schema, test.transformIdentifier)
		if err != nil {
			b.Fatalf("Test %q - failed to initialize transformer: %v", test.description, err)
		}

		b.Run(test.description, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := tr.Transform(test.in)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkSaveInTree(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, test := range saveValueTests {
			err := saveInTree(test.tree, test.jsonPath, test.value)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}
