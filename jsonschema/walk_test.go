package jsonschema

import (
	"encoding/json"
	"reflect"
	"testing"
)

type testWalker struct {
	calls map[string]json.RawMessage
}

func (tw *testWalker) walkFn(path string, value json.RawMessage) error {
	tw.calls[path] = value
	return nil
}

func TestWalkJSONSchema(t *testing.T) {
	tests := []struct {
		description string
		oneOfType   string
		schemaPath  string
		want        map[string]json.RawMessage
		wantErr     bool
	}{
		{
			description: "Basic walk, no allOf, no oneOf",
			schemaPath:  "./test_data/image.json",
			want: map[string]json.RawMessage{
				"$.type": []byte(`{
      "type": "string",
      "enum": [
        "image"
      ]
    }`),
				"$.crops": []byte(`{
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "name": {
            "type": "string",
            "default": "name"
          },
          "width": {
            "type": "number"
          },
          "height": {
            "type": "number"
          },
          "path": {
            "type": "string"
          },
          "relativePath": {
            "type": "string"
          }
        },
        "required":[
          "name",
          "width",
          "height",
          "path",
          "relativePath"
        ]
      }
    }`),
				"$.crops[*]": []byte(`{
        "type": "object",
        "properties": {
          "name": {
            "type": "string",
            "default": "name"
          },
          "width": {
            "type": "number"
          },
          "height": {
            "type": "number"
          },
          "path": {
            "type": "string"
          },
          "relativePath": {
            "type": "string"
          }
        },
        "required":[
          "name",
          "width",
          "height",
          "path",
          "relativePath"
        ]
      }`),
				"$.crops[*].name": []byte(`{
            "type": "string",
            "default": "name"
          }`),
				"$.crops[*].width": []byte(`{
            "type": "number"
          }`),
				"$.crops[*].height": []byte(`{
            "type": "number"
          }`),
				"$.crops[*].path": []byte(`{
            "type": "string"
          }`),
				"$.crops[*].relativePath": []byte(`{
            "type": "string"
          }`),
				"$.URL": []byte(`{
      "type": "object",
      "properties": {
        "publish": {
          "type": "string",
          "transform": {
            "cumulo": {
              "from" : [
                {
                  "jsonPath": "$.publishUrl"
                }
              ]
            }
          }
        },
        "absolute": {
          "type": "string",
          "transform": {
            "cumulo": {
              "from" : [
                {
                  "jsonPath": "$.absoluteUrl"
                }
              ]
            }
          }
        }
      },
      "required":[
        "publish",
        "absolute"
      ]
    }`),
				"$.URL.publish": []byte(`{
          "type": "string",
          "transform": {
            "cumulo": {
              "from" : [
                {
                  "jsonPath": "$.publishUrl"
                }
              ]
            }
          }
        }`),
				"$.URL.absolute": []byte(`{
          "type": "string",
          "transform": {
            "cumulo": {
              "from" : [
                {
                  "jsonPath": "$.absoluteUrl"
                }
              ]
            }
          }
        }`),
			},
		},
		{
			description: "Walk with allOf, no oneOf",
			schemaPath:  "./test_data/embed_parent.json",
			want: map[string]json.RawMessage{"$.type": []byte(`{
      "type": "string",
      "enum": [
        "embed"
      ]
    }`),
			},
		},
		{
			description: "Walk with oneOf, no allOf",
			oneOfType:   "image",
			schemaPath:  "./test_data/image_parent.json",
			want: map[string]json.RawMessage{
				"$.type": []byte(`{
      "type": "string",
      "enum": [
        "image"
      ]
    }`),
				"$.crops": []byte(`{
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "name": {
            "type": "string",
            "default": "name"
          },
          "width": {
            "type": "number"
          },
          "height": {
            "type": "number"
          },
          "path": {
            "type": "string"
          },
          "relativePath": {
            "type": "string"
          }
        },
        "required":[
          "name",
          "width",
          "height",
          "path",
          "relativePath"
        ]
      }
    }`),
				"$.crops[*]": []byte(`{
        "type": "object",
        "properties": {
          "name": {
            "type": "string",
            "default": "name"
          },
          "width": {
            "type": "number"
          },
          "height": {
            "type": "number"
          },
          "path": {
            "type": "string"
          },
          "relativePath": {
            "type": "string"
          }
        },
        "required":[
          "name",
          "width",
          "height",
          "path",
          "relativePath"
        ]
      }`),
				"$.crops[*].name": []byte(`{
            "type": "string",
            "default": "name"
          }`),
				"$.crops[*].width": []byte(`{
            "type": "number"
          }`),
				"$.crops[*].height": []byte(`{
            "type": "number"
          }`),
				"$.crops[*].path": []byte(`{
            "type": "string"
          }`),
				"$.crops[*].relativePath": []byte(`{
            "type": "string"
          }`),
				"$.URL": []byte(`{
      "type": "object",
      "properties": {
        "publish": {
          "type": "string",
          "transform": {
            "cumulo": {
              "from" : [
                {
                  "jsonPath": "$.publishUrl"
                }
              ]
            }
          }
        },
        "absolute": {
          "type": "string",
          "transform": {
            "cumulo": {
              "from" : [
                {
                  "jsonPath": "$.absoluteUrl"
                }
              ]
            }
          }
        }
      },
      "required":[
        "publish",
        "absolute"
      ]
    }`),
				"$.URL.publish": []byte(`{
          "type": "string",
          "transform": {
            "cumulo": {
              "from" : [
                {
                  "jsonPath": "$.publishUrl"
                }
              ]
            }
          }
        }`),
				"$.URL.absolute": []byte(`{
          "type": "string",
          "transform": {
            "cumulo": {
              "from" : [
                {
                  "jsonPath": "$.absoluteUrl"
                }
              ]
            }
          }
        }`),
			},
		},
		{
			description: "Advanced walk does it all",
			oneOfType:   "array-of-array",
			schemaPath:  "./test_data/parent.json",
			want: map[string]json.RawMessage{
				"$.type": []byte(`{
      "type": "string",
      "enum": [
        "array-of-array"
      ]
    }`),
				"$.crops": []byte(`{
      "type": "array",
      "items": {
        "type": "array",
        "items": {
          "type": "object",
          "properties": {
            "name": {
              "type": "string"
            }
          }
        }
      }
    }`),
				"$.crops[*]": []byte(`{
        "type": "array",
        "items": {
          "type": "object",
          "properties": {
            "name": {
              "type": "string"
            }
          }
        }
      }`),
				"$.crops[*][*]": []byte(`{
          "type": "object",
          "properties": {
            "name": {
              "type": "string"
            }
          }
        }`),
				"$.crops[*][*].name": []byte(`{
              "type": "string"
            }`),
			},
		},
		{
			description: "Object with missing properties",
			schemaPath:  "./test_data/bad-object.json",
			wantErr:     true,
		},
		{
			description: "Array with missing Items",
			schemaPath:  "./test_data/bad-array.json",
			wantErr:     true,
		},
	}

	for _, test := range tests {
		walker := testWalker{calls: make(map[string]json.RawMessage)}
		schema, err := SchemaFromFile(test.schemaPath, test.oneOfType)
		if err != nil {
			t.Fatal(err)
		}
		err = Walk(schema, walker.walkFn)

		switch {
		case test.wantErr && err != nil:
			continue
		case test.wantErr && err == nil:
			t.Errorf("Test %q - got nil error want error", test.description)
		case !test.wantErr && err != nil:
			t.Errorf("Test %q - got error: %v", test.description, err)
			continue
		}
		if got, want := len(walker.calls), len(test.want); got != want {
			t.Errorf("Test %q - got %d calls, want %d", test.description, got, want)
		}
		for key, call := range walker.calls {
			if !reflect.DeepEqual(call, test.want[key]) {
				t.Errorf("Test %q - at got key %q got call\n%s\n\twant\n%s", test.description, key, call, test.want[key])
			}
		}
		for key, call := range test.want {
			if !reflect.DeepEqual(call, walker.calls[key]) {
				t.Errorf("Test %q - at want key %q got call\n%s\n\twant\n%s", test.description, key, walker.calls[key], call)
			}
		}
	}
}
