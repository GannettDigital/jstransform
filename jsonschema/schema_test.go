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
				Instance: Instance{
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
				Instance: Instance{
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
				Instance: Instance{
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
				Instance: Instance{
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
		{
			description: "Nested AllOf",
			oneOfType:   "",
			schemaPath:  "./test_data/parent3.json",
			want: &Schema{
				Instance: Instance{
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
			description: "Nested AllOf with oneOf",
			oneOfType:   "image",
			schemaPath:  "./test_data/parent4.json",
			want: &Schema{
				Instance: Instance{
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
			description: "AllOf in referenced type",
			oneOfType:   "",
			schemaPath:  "./test_data/embed_embed.json",
			want: &Schema{
				Instance: Instance{
					Properties: map[string]json.RawMessage{
						"embed": json.RawMessage(`{
      "type": "array",
      "items": {
        "additionalProperties": true,
		"allOf": [
          {
            "additionalProperties": true,
            "fromRef": "./embed.json",
            "properties": {
              "type": {
      		    "type": "string",
      		    "enum": [
				  "embed"
      		    ]
              }
            },
		    "$schema": "http://json-schema.org/draft-04/schema#",
		    "type": [
				"object"
			]
          }
		],
        "properties": {
		  "type":{
      		"type": "string",
      		"enum": [
				"embed"
      		]
          }
        },
		"$schema": "http://json-schema.org/draft-04/schema#",
		"type": [
			"object"
		],
	  	"fromRef": "./embed_parent.json"
      }
  }`),
					},
				},
			},
		},
	}
	for _, test := range tests {
		got, err := SchemaFromFile(test.schemaPath, test.oneOfType)

		var gotProperties, wantProperties []byte
		if !test.wantErr && got != nil {
			gotProperties, err = json.MarshalIndent(got.Properties, "", "  ")
			if err != nil {
				t.Errorf("Test %q - failed to marshal got.Properties: %v", test.description, err)
			}
			wantProperties, err = json.MarshalIndent(test.want.Properties, "", "  ")
			if err != nil {
				t.Errorf("Test %q - failed to marshal test.want.Properties: %v", test.description, err)
			}
		}

		switch {
		case test.wantErr && err != nil:
			continue
		case test.wantErr && err == nil:
			t.Errorf("Test %q - got nil error want error", test.description)
		case !test.wantErr && err != nil:
			t.Errorf("Test %q - got error: %v", test.description, err)
		case !reflect.DeepEqual(gotProperties, wantProperties):
			t.Errorf("Test %q - got Properties\n%s\nwant\n%s", test.description, gotProperties, wantProperties)
		case !reflect.DeepEqual(got.Items, test.want.Items):
			t.Errorf("Test %q - got Items\n%s\nwant\n%s", test.description, got.Items, test.want.Items)
		case !reflect.DeepEqual(got.Required, test.want.Required):
			t.Errorf("Test %q - got Required\n%s\nwant\n%s", test.description, got.Required, test.want.Required)
		}
	}
}

func TestMappings(t *testing.T) {
	tests := []struct {
		description string
		oneOfType   string
		schemaPath  string
		want        Instance
	}{
		{
			description: "Simple map of top level fields",
			schemaPath:  "./test_data/simple.json",
			want: Instance{
				AdditionalProperties: false,
				Description:          "missing",
				Type:                 []string{"object"},
				Required:             []string{"required"},
			},
		},
		{
			description: "Missing additional properties should default appropriately",
			schemaPath:  "./test_data/missing-additional-properties.json",
			want: Instance{
				AdditionalProperties: true,
				Description:          "missing",
				Type:                 []string{"object"},
				Required:             []string{"required"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			got, err := SchemaFromFile(test.schemaPath, test.oneOfType)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(got.Instance, test.want) {
				t.Errorf("Test %q - got \n%v\n\twant\n%v", test.description, got.Instance, test.want)
			}
		})
	}
}

func TestSchemaTypes(t *testing.T) {
	tests := []struct {
		description    string
		path           string
		wantAllOf      []string
		wantOneOf      []string
		wantProperties []string
	}{
		{
			description: "Single AllOf",
			path:        "test_data/embed_parent.json",
			wantAllOf:   []string{"./embed.json"},
		},
		{
			description: "Single oneOf",
			path:        "test_data/image_parent.json",
			wantOneOf:   []string{"./image.json"},
		},
		{
			description: "Multiple oneOf",
			path:        "test_data/parent.json",
			wantOneOf:   []string{"./image.json", "./array-of-array.json"},
		},
		{
			description: "Multiple allOf and oneOf",
			path:        "test_data/parent2.json",
			wantAllOf:   []string{"./embed.json", "./operations.json"},
			wantOneOf:   []string{"./image.json", "./array-of-array.json"},
		},
		{
			description:    "No oneOf or allOf",
			path:           "test_data/image.json",
			wantProperties: []string{"URL", "crops", "type"},
		},
	}

	for _, test := range tests {
		gotAllOf, gotOneOf, gotProperties, err := SchemaTypes(test.path)
		if err != nil {
			t.Fatalf("Test %q - failed: %v", test.description, err)
		}

		if !reflect.DeepEqual(gotAllOf, test.wantAllOf) {
			t.Errorf("Test %q - got AllOf %v, want %v", test.description, gotAllOf, test.wantAllOf)
		}

		if !reflect.DeepEqual(gotOneOf, test.wantOneOf) {
			t.Errorf("Test %q - got OneOf %v, want %v", test.description, gotOneOf, test.wantOneOf)
		}

		if !reflect.DeepEqual(gotProperties, test.wantProperties) {
			t.Errorf("Test %q - got properties %v, want %v", test.description, gotProperties, test.wantProperties)
		}
	}
}
