package jsonschema

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
)

func TestDereference(t *testing.T) {
	tests := []struct {
		description string
		schemaPath  string
	}{
		{
			description: "Self Reference",
			schemaPath:  "./test_data/reference/jsref_image-self.json",
		},
		{
			description: "File Reference",
			schemaPath:  "./test_data/reference/jsref_image-file.json",
		}, /*
			// TODO: figure out mock http client to test http json references
				{
					description: "Http Reference",
					schemaPath:  "./test_data/reference/jsref_image-http.json",
				},
		*/
		{
			description: "List Reference",
			schemaPath:  "./test_data/reference/jsref_image-list.json",
		},
		{
			description: "Shifted List Reference",
			schemaPath:  "./test_data/reference/jsref_image-list1.json",
		},
		{
			description: "Multiple List Reference",
			schemaPath:  "./test_data/reference/jsref_image-list2.json",
		},
		{
			description: "Nested List Reference",
			schemaPath:  "./test_data/reference/jsref_image-nest.json",
		},
		{
			description: "Ingestion Example",
			schemaPath:  "./test_data/reference/jsref_asset-event-facebook.json",
		},
	}
	for _, test := range tests {

		var want interface{}
		wantPath := strings.Replace(test.schemaPath, "/jsref_", "/deref_", 1)
		wantJson, err := ioutil.ReadFile(wantPath)
		if err != nil {
			t.Errorf("Test %q - failed to read json file %q: %v", test.description, wantPath, err)
		}
		json.Unmarshal(wantJson, &want)

		var got interface{}
		gotJson, err := ioutil.ReadFile(test.schemaPath)
		if err != nil {
			t.Errorf("Test %q - failed to read json file %q: %v", test.description, test.schemaPath, err)
		}

		var gj schemaJSON
		gj, gotJson, err = Dereference(test.schemaPath, gotJson)
		if err != nil {
			t.Errorf("Test %q - failed to dereference json file %q (%q): %v", test.description, test.schemaPath, gj, err)
		}
		json.Unmarshal(gotJson, &got)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Test %q - got\n%s\nwant\n%s", test.description, got, want)
		}
	}
}
