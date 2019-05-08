package transform

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"
)

var (
	testIn = map[string]interface{}{
		"type": "image",
		"date": testTimeStr,
		"crops": []interface{}{
			map[string]interface{}{
				"height":       0,
				"path":         "path",
				"relativePath": "",
				"width":        1,
			},
			map[string]interface{}{
				"name":         "aname",
				"height":       0,
				"path":         "empty",
				"relativePath": "empty",
				"width":        0,
			},
		},
		"otherCrops": []interface{}{
			[]interface{}{
				map[string]interface{}{
					"height":       10,
					"path":         "otherpath",
					"relativePath": "other",
					"width":        11,
				},
				map[string]interface{}{
					"name":         "otheraname",
					"height":       10,
					"path":         "otherempty",
					"relativePath": "otherempty",
					"width":        10,
				},
			},
			[]interface{}{
				map[string]interface{}{
					"name":         "otheraname2",
					"height":       102,
					"path":         "otherempty2",
					"relativePath": "otherempty2",
					"width":        102,
				},
			},
		},
		"published":   true,
		"publishUrl":  "publishURL",
		"absoluteUrl": "absoluteURL",
	}

	testInBadTime = map[string]interface{}{
		"date": testBadTimeStr,
	}

	testTimeStr = "2018-01-01T01:01:00Z"

	testBadTimeStr = "2000-10-15"
)

func TestArrayTransform(t *testing.T) {
	var nilSlice []interface{}

	tests := []struct {
		description string
		format      inputFormat
		child       instanceTransformer
		path        string
		raw         json.RawMessage
		want        interface{}
	}{
		{
			description: "empty",
			format:      jsonInput,
			path:        "$.empty",
			raw:         json.RawMessage(`{"type":"array"}`),
			want:        nilSlice,
		},
		{
			description: "default with no child",
			format:      jsonInput,
			path:        "$.crops",
			raw:         json.RawMessage(`{"type":"array"}`),
			want: []interface{}{
				map[string]interface{}{
					"height":       0,
					"path":         "path",
					"relativePath": "",
					"width":        1,
				},
				map[string]interface{}{
					"name":         "aname",
					"height":       0,
					"path":         "empty",
					"relativePath": "empty",
					"width":        0,
				},
			},
		},
		{
			description: "transform with no child",
			format:      jsonInput,
			path:        "$.crops",
			raw:         json.RawMessage(`{"type":"object","transform":{"test":{"from":[{"jsonPath":"$.otherCrops[0]"}]}}}`),
			want: []interface{}{
				map[string]interface{}{
					"height":       10,
					"path":         "otherpath",
					"relativePath": "other",
					"width":        11,
				},
				map[string]interface{}{
					"name":         "otheraname",
					"height":       10,
					"path":         "otherempty",
					"relativePath": "otherempty",
					"width":        10,
				},
			},
		},
		{
			description: "scalar child",
			format:      jsonInput,
			child: &scalarTransformer{
				defaultValue: "name",
				format:       jsonInput,
				jsonType:     "string",
				jsonPath:     "$.crops[*]",
				transforms: &transformInstructions{
					From:   []*transformInstruction{{jsonPath: "$.crops[*].name"}},
					Method: 0,
				},
			},
			path: "$.crops",
			raw:  json.RawMessage(`{"type":"array"}`),
			want: []interface{}{"name", "aname"},
		},
		{
			description: "object child",
			format:      jsonInput,
			child: &objectTransformer{
				jsonPath: "$.crops[*]",
				format:   jsonInput,
				children: map[string]instanceTransformer{
					"name": &scalarTransformer{
						defaultValue: "name",
						jsonType:     "string",
						jsonPath:     "$.crops[*].name",
						format:       jsonInput,
						transforms: &transformInstructions{
							From:   []*transformInstruction{{jsonPath: "$.crops[*].name"}},
							Method: 0,
						},
					},
					"height": &scalarTransformer{
						jsonType: "number",
						jsonPath: "$.crops[*].height",
						format:   jsonInput,
						transforms: &transformInstructions{
							From:   []*transformInstruction{{jsonPath: "$.crops[*].height"}},
							Method: 0,
						},
					},
					"path": &scalarTransformer{
						jsonType: "string",
						jsonPath: "$.crops[*].path",
						format:   jsonInput,
						transforms: &transformInstructions{
							From:   []*transformInstruction{{jsonPath: "$.crops[*].path"}},
							Method: 0,
						},
					},
					"relativePath": &scalarTransformer{
						jsonType: "string",
						jsonPath: "$.crops[*].relativePath",
						format:   jsonInput,
						transforms: &transformInstructions{
							From:   []*transformInstruction{{jsonPath: "$.crops[*].relativePath"}},
							Method: 0,
						},
					},
					"width": &scalarTransformer{
						jsonType: "number",
						jsonPath: "$.crops[*].width",
						format:   jsonInput,
						transforms: &transformInstructions{
							From:   []*transformInstruction{{jsonPath: "$.crops[*].width"}},
							Method: 0,
						},
					},
				},
			},
			path: "$.crops",
			raw:  json.RawMessage(`{"type":"array"}`),
			want: []interface{}{
				map[string]interface{}{
					"name":         "name",
					"height":       0,
					"path":         "path",
					"relativePath": "",
					"width":        1,
				},
				map[string]interface{}{
					"name":         "aname",
					"height":       0,
					"path":         "empty",
					"relativePath": "empty",
					"width":        0,
				},
			},
		},
		{
			description: "nested array",
			format:      jsonInput,
			child: &arrayTransformer{
				jsonPath: "$.otherCrops[*]",
				format:   jsonInput,
				childTransformer: &objectTransformer{
					jsonPath: "$.otherCrops[*][*]",
					children: map[string]instanceTransformer{
						"name": &scalarTransformer{
							defaultValue: "name",
							jsonType:     "string",
							jsonPath:     "$.otherCrops[*][*].name",
							format:       jsonInput,
						},
						"height": &scalarTransformer{
							jsonType: "number",
							jsonPath: "$.otherCrops[*][*].height",
							format:   jsonInput,
						},
						"path": &scalarTransformer{
							jsonType: "string",
							jsonPath: "$.otherCrops[*][*].path",
							format:   jsonInput,
						},
						"relativePath": &scalarTransformer{
							jsonType: "string",
							jsonPath: "$.otherCrops[*][*].relativePath",
							format:   jsonInput,
						},
						"width": &scalarTransformer{
							jsonType: "number",
							jsonPath: "$.otherCrops[*][*].width",
							format:   jsonInput,
						},
					},
				},
			},
			path: "$.otherCrops",
			raw:  json.RawMessage(`{"type":"array"}`),
			want: []interface{}{
				[]interface{}{
					map[string]interface{}{
						"name":         "name",
						"height":       10,
						"path":         "otherpath",
						"relativePath": "other",
						"width":        11,
					},
					map[string]interface{}{
						"name":         "otheraname",
						"height":       10,
						"path":         "otherempty",
						"relativePath": "otherempty",
						"width":        10,
					},
				},
				[]interface{}{
					map[string]interface{}{
						"name":         "otheraname2",
						"height":       102,
						"path":         "otherempty2",
						"relativePath": "otherempty2",
						"width":        102,
					},
				},
			},
		},
	}

	for _, test := range tests {
		at, err := newArrayTransformer(test.path, "test", test.raw, test.format)
		if err != nil {
			t.Fatalf("Test %q - failed to initialize array transformer: %v", test.description, err)
		}

		at.childTransformer = test.child

		testInCopy := make(map[string]interface{})
		for k, v := range testIn {
			testInCopy[k] = v
		}
		got, err := at.transform(testInCopy, nil)
		if err != nil {
			t.Errorf("Test %q - failed transform: %v", test.description, err)
		}

		if !reflect.DeepEqual(test.want, got) {
			t.Errorf("Test %q - got %#v, want %#v", test.description, got, test.want)
		}
	}
}

func TestObjectTransform(t *testing.T) {
	tests := []struct {
		description string
		in          interface{}
		format      inputFormat
		children    map[string]instanceTransformer
		path        string
		raw         json.RawMessage
		want        interface{}
	}{
		{
			description: "empty",
			in:          testIn,
			format:      jsonInput,
			path:        "$",
			raw:         json.RawMessage(`{"type":"object"}`),
			want:        nil,
		},
		{
			description: "with scalar children",
			in:          testIn,
			format:      jsonInput,
			children: map[string]instanceTransformer{
				"name": &scalarTransformer{
					defaultValue: "name",
					jsonType:     "string",
					jsonPath:     "$.firstCrop.name",
					format:       jsonInput,
					transforms: &transformInstructions{
						From:   []*transformInstruction{{jsonPath: "$.crops[0].name"}},
						Method: 0,
					},
				},
				"height": &scalarTransformer{
					jsonType: "number",
					jsonPath: "$.firstCrop.height",
					format:   jsonInput,
					transforms: &transformInstructions{
						From:   []*transformInstruction{{jsonPath: "$.crops[0].height"}},
						Method: 0,
					},
				},
				"path": &scalarTransformer{
					jsonType: "string",
					jsonPath: "$.firstCrop.path",
					format:   jsonInput,
					transforms: &transformInstructions{
						From:   []*transformInstruction{{jsonPath: "$.crops[0].path"}},
						Method: 0,
					},
				},
				"relativePath": &scalarTransformer{
					jsonType: "string",
					jsonPath: "$.firstCrop.relativePath",
					format:   jsonInput,
					transforms: &transformInstructions{
						From:   []*transformInstruction{{jsonPath: "$.crops[0].relativePath"}},
						Method: 0,
					},
				},
				"width": &scalarTransformer{
					jsonType: "number",
					jsonPath: "$.firstCrop.width",
					format:   jsonInput,
					transforms: &transformInstructions{
						From:   []*transformInstruction{{jsonPath: "$.crops[0].width"}},
						Method: 0,
					},
				},
			},
			path: "$.firstCrop",
			raw:  json.RawMessage(`{"type":"object"}`),
			want: map[string]interface{}{
				"name":         "name",
				"height":       0,
				"path":         "path",
				"relativePath": "",
				"width":        1,
			},
		},
		{
			description: "with nil value after transform",
			in:          testIn,
			format:      jsonInput,
			children: map[string]instanceTransformer{
				"name": &scalarTransformer{
					defaultValue: "name",
					jsonType:     "string",
					jsonPath:     "$.firstCrop.name",
					format:       jsonInput,
					transforms: &transformInstructions{
						From:   []*transformInstruction{{jsonPath: "$.crops[0].name"}},
						Method: 0,
					},
				},
			},
			path: "$.firstCrop",
			raw:  json.RawMessage(`{"type":"object","transform":{"test":{"from":[{"jsonPath":"$.notFound"}]}}}`),
			want: map[string]interface{}{
				"name": "name",
			},
		},
	}

	for _, test := range tests {
		ot, err := newObjectTransformer(test.path, "test", test.raw, test.format)
		if err != nil {
			t.Fatalf("Test %q - failed to initialize object transformer: %v", test.description, err)
		}

		ot.children = test.children

		got, err := ot.transform(test.in, nil)
		if err != nil {
			t.Errorf("Test %q - failed transform: %v", test.description, err)
		}

		if !reflect.DeepEqual(test.want, got) {
			t.Errorf("Test %q - got %#v, want %#v", test.description, got, test.want)
		}
	}
}

func TestScalarTransform(t *testing.T) {
	testTime, err := time.Parse(time.RFC3339, testTimeStr)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		description  string
		in           interface{}
		path         string
		instanceType string
		format       inputFormat
		raw          json.RawMessage
		want         interface{}
		wantError    string
	}{
		{
			description:  "unchanged string",
			in:           testIn,
			path:         "$.type",
			instanceType: "string",
			format:       jsonInput,
			raw:          json.RawMessage(`{"type":"string","enum":["image"]}`),
			want:         "image",
		},
		{
			description:  "unchanged number",
			in:           testIn,
			path:         "$.crops[0].height",
			instanceType: "number",
			format:       jsonInput,
			raw:          json.RawMessage(`{ "type": "number" }`),
			want:         0,
		},
		{
			description:  "unchanged time",
			in:           testIn,
			path:         "$.date",
			instanceType: "string",
			format:       jsonInput,
			raw:          json.RawMessage(`{ "type": "string", "format": "date-time" }`),
			want:         testTime,
		},
		{
			description:  "unchanged bool",
			in:           testIn,
			path:         "$.published",
			instanceType: "boolean",
			format:       jsonInput,
			raw:          json.RawMessage(`{ "type": "boolean"}`),
			want:         true,
		},
		{
			description:  "string with default",
			in:           testIn,
			path:         "$.type2",
			instanceType: "string",
			format:       jsonInput,
			raw:          json.RawMessage(`{"type":"string","enum":["image"],"default":"type2"}`),
			want:         "type2",
		},
		{
			description:  "number with default",
			in:           testIn,
			path:         "$.crops[0].multiplier",
			instanceType: "number",
			format:       jsonInput,
			raw:          json.RawMessage(`{ "type": "number", "default": 10 }`),
			want:         float64(10),
		},
		{
			description:  "time with default",
			in:           testIn,
			path:         "$.date2",
			instanceType: "string",
			format:       jsonInput,
			raw:          json.RawMessage(fmt.Sprintf(`{ "type": "string", "format": "date-time", "default": "%s" }`, testTimeStr)),
			want:         testTimeStr,
		},
		{
			description:  "bool with default",
			in:           testIn,
			path:         "$.deleted",
			instanceType: "boolean",
			format:       jsonInput,
			raw:          json.RawMessage(`{ "type": "boolean", "default":true}`),
			want:         true,
		},
		{
			description:  "string transform",
			in:           testIn,
			path:         "$.type",
			instanceType: "string",
			format:       jsonInput,
			raw:          json.RawMessage(`{"type":"string","enum":["image"],"transform":{"test":{"from":[{"jsonPath":"$.publishUrl"}]}}}`),
			want:         "publishURL",
		},
		{
			description:  "number transform",
			in:           testIn,
			path:         "$.crops[0].height",
			instanceType: "number",
			format:       jsonInput,
			raw:          json.RawMessage(`{ "type": "number" ,"transform":{"test":{"from":[{"jsonPath":"$.crops[0].width"}]}}}`),
			want:         1,
		},
		{
			description:  "time transform",
			in:           testIn,
			path:         "$.anotherDate",
			instanceType: "string",
			format:       jsonInput,
			raw:          json.RawMessage(`{ "type": "string", "format": "date-time" ,"transform":{"test":{"from":[{"jsonPath":"$.date"}]}}}`),
			want:         testTime,
		},
		{
			description:  "bool transform",
			in:           testIn,
			path:         "$.WantToPublish",
			instanceType: "boolean",
			format:       jsonInput,
			raw:          json.RawMessage(`{ "type": "boolean","transform":{"test":{"from":[{"jsonPath":"$.published"}]}}}`),
			want:         true,
		},
		{
			description:  "time transform with bad formatting",
			in:           testInBadTime,
			path:         "$.date",
			instanceType: "string",
			format:       jsonInput,
			raw:          json.RawMessage(`{ "type": "string", "format": "date-time"}`),
			wantError:    "parsing time \"2000-10-15\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"\" as \"T\"",
		},
	}

	for _, test := range tests {
		st, err := newScalarTransformer(test.path, "test", test.raw, test.instanceType, test.format)
		if err != nil {
			t.Fatalf("Test %q - failed to initialize scalar transformer: %v", test.description, err)
		}

		got, err := st.transform(test.in, nil)

		if err != nil {
			if err.Error() == test.wantError {
				continue // pass
			}
			if test.wantError != "" {
				t.Errorf("Test %q - failed to produce error: %v, instead got: %v", test.description, test.wantError, err)
				continue
			}
			t.Errorf("Test %q - failed transform: %v", test.description, err)
		}

		if !reflect.DeepEqual(test.want, got) {
			t.Errorf("Test %q - got %#v, want %#v", test.description, got, test.want)
		}
	}
}
