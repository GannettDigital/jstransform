package jsonschema

import (
	"encoding/json"
	"reflect"
	"testing"
)

var imageProperties = map[string]json.RawMessage{
	"type": []byte(`{
      "type": "string",
      "enum": [
        "image"
      ]
    }`),
	"crops": []byte(`{
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
	"URL": []byte(`{
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
}

func TestSchemaFromFile(t *testing.T) {
	tests := []struct {
		description string
		oneOfType   string
		schemaPath  string
		want        *Schema
		wantErr     bool
	}{
		{
			description: "Basic Schema Load",
			schemaPath:  "./test_data/image.json",
			want: &Schema{
				Instance{
					Properties: imageProperties,
					Required: []string{
						"type",
						"crops",
						"orientation",
						"credit",
						"URL",
						"caption",
						"originalSize",
						"datePhotoTaken",
					},
				},
			},
		},
		{
			description: "Load allOf, no oneOf",
			schemaPath:  "./test_data/embed_parent.json",
			want: &Schema{
				Instance{
					Properties: map[string]json.RawMessage{
						"type": json.RawMessage(`{
      "type": "string",
      "enum": [
        "embed"
      ]
    }`),
					},
				},
			},
		},
		{
			description: "Load oneOf, no allOf",
			oneOfType:   "image",
			schemaPath:  "./test_data/image_parent.json",
			want: &Schema{
				Instance{
					Properties: imageProperties,
					Required: []string{
						"type",
						"crops",
						"orientation",
						"credit",
						"URL",
						"caption",
						"originalSize",
						"datePhotoTaken",
					},
				},
			},
		},
		{
			description: "Choose oneOf from options",
			oneOfType:   "image",
			schemaPath:  "./test_data/parent.json",
			want: &Schema{
				Instance{
					Properties: imageProperties,
					Required: []string{
						"type",
						"crops",
						"orientation",
						"credit",
						"URL",
						"caption",
						"originalSize",
						"datePhotoTaken",
					},
				},
			},
		},
		{
			description: "Missing file",
			oneOfType:   "image",
			schemaPath:  "./test_data/does-not-exist.json",
			wantErr:     true,
		},
		{
			description: "Ref to missing file",
			oneOfType:   "image",
			schemaPath:  "./test_data/bad-parent.json",
			wantErr:     true,
		},
	}
	for _, test := range tests {
		got, err := SchemaFromFile(test.schemaPath, test.oneOfType)

		switch {
		case test.wantErr && err != nil:
			continue
		case test.wantErr && err == nil:
			t.Errorf("Test %q - got nil error want error", test.description)
		case !test.wantErr && err != nil:
			t.Errorf("Test %q - got error: %v", test.description, err)
		case !reflect.DeepEqual(got.Properties, test.want.Properties):
			t.Errorf("Test %q - got Properties\n%s\nwant\n%s", test.description, got.Properties, test.want.Properties)
		case !reflect.DeepEqual(got.Items, test.want.Items):
			t.Errorf("Test %q - got Items\n%s\nwant\n%s", test.description, got.Items, test.want.Items)
		case !reflect.DeepEqual(got.Required, test.want.Required):
			t.Errorf("Test %q - got Required\n%s\nwant\n%s", test.description, got.Required, test.want.Required)
		}
	}
}
