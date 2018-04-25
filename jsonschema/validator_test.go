package jsonschema

import (
	"io/ioutil"
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		description string
		schemaPath  string
		jsonPath    string
		wantValid   bool
		wantErr     bool
	}{
		{
			description: "valid image JSON",
			schemaPath:  "test_data/image.json",
			jsonPath:    "test_data/imageRaw.json",
			wantValid:   true,
		},
		{
			description: "valid image JSON using parent schema",
			schemaPath:  "test_data/parent.json",
			jsonPath:    "test_data/imageRaw.json",
			wantValid:   true,
		},
		{
			description: "invalid image JSON",
			schemaPath:  "test_data/image.json",
			jsonPath:    "test_data/embed.json",
			wantValid:   false,
			wantErr:     true,
		},
	}
	for _, test := range tests {
		v, err := NewValidator(test.schemaPath)
		if err != nil {
			t.Fatalf("Test %q - failed to load schema: %v", test.description, err)
		}
		raw, err := ioutil.ReadFile(test.jsonPath)
		if err != nil {
			t.Fatalf("Test %q - failed to load test JSON: %v", test.description, err)
		}

		got, err := v.Validate(raw)

		if (err != nil) != test.wantErr {
			t.Errorf("Test %q - got error %v, want error %t", test.description, err, test.wantErr)
		}
		if got != test.wantValid {
			t.Errorf("Test %q - got valid %t, want %t", test.description, got, test.wantValid)
		}
	}
}
