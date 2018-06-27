package jsonschema

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestDereference(t *testing.T) {
	// create listener with desired port
	custom := "127.0.0.1:12345"
	tl, err := net.Listen("tcp", custom)
	if err != nil {
		t.Errorf("Test failed to create listener on %s %v", custom, err)
	}

	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs := http.FileServer(http.Dir("./test_data/reference"))
		fs.ServeHTTP(w, r)
	}))

	// Close listener, replace and start
	ts.Listener.Close()
	ts.Listener = tl
	ts.Start()
	defer ts.Close()

	tests := []struct {
		description string
		schemaPath  string
		wantErr     bool
	}{
		{
			description: "Self Reference",
			schemaPath:  "./test_data/reference/jsref_image-self.json",
		},
		{
			description: "File Reference",
			schemaPath:  "./test_data/reference/jsref_image-file.json",
		},
		{
			description: "Http Reference",
			schemaPath:  "./test_data/reference/jsref_image-http.json",
		},
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
			description: "Infinite Loop Reference",
			schemaPath:  "./test_data/reference/jsref_image-loop.json",
			wantErr:     true,
		},
		{
			description: "Ingestion Example",
			schemaPath:  "./test_data/reference/jsref_asset-event-facebook.json",
		},
	}
	for _, test := range tests {

		var want interface{}
		if !test.wantErr {
			wantPath := strings.Replace(test.schemaPath, "/jsref_", "/deref_", 1)
			wantJson, err := ioutil.ReadFile(wantPath)
			if err != nil {
				t.Errorf("Test %q - failed to read json want file %q: %v", test.description, wantPath, err)
			}
			json.Unmarshal(wantJson, &want)
		}

		var got interface{}
		gotJson, err := ioutil.ReadFile(test.schemaPath)
		if err != nil {
			t.Errorf("Test %q - failed to read json got file %q: %v", test.description, test.schemaPath, err)
		}

		gotJson, err = dereference(test.schemaPath, gotJson, true)
		if err != nil && !test.wantErr {
			t.Errorf("Test %q - failed to dereference json file %q: %v", test.description, test.schemaPath, err)
		}
		json.Unmarshal(gotJson, &got)

		switch {
		case test.wantErr && err != nil:
			continue
		case !reflect.DeepEqual(got, want):
			t.Errorf("Test %q - got\n%s\nwant\n%s", test.description, got, want)
		}
	}
}

func TestParseReference(t *testing.T) {
	var source interface{}
	defPath := "./test_data/reference/jsref_image-defs.json"
	defJson, err := ioutil.ReadFile(defPath)
	if err != nil {
		t.Errorf("Failed to read json file %q: %v", defPath, err)
	}
	json.Unmarshal(defJson, &source)

	tests := []struct {
		description string
		refPaths    []string
		want        json.RawMessage
	}{
		{
			description: "Simple Type Reference",
			refPaths:    []string{"definitions", "arraytype"},
			want:        json.RawMessage(`{"type":"array"}`),
		},
		{
			description: "Deeply Nested Reference",
			refPaths:    []string{"definitions", "deeply", "nested", "objecttype"},
			want:        json.RawMessage(`{"type":"object"}`),
		},
		{
			description: "Incorrect Reference",
			refPaths:    []string{"definitions", "deeply", "nested", "bad"},
			want:        json.RawMessage(`null`),
		},
		{
			description: "Image Url Reference",
			refPaths:    []string{"definitions", "imageurl"},
			want:        json.RawMessage(`{"properties":{"absolute":{"transform":{"cumulo":{"from":[{"jsonPath":"$.absoluteUrl"}]}},"type":"string"},"publish":{"transform":{"cumulo":{"from":[{"jsonPath":"$.publishUrl"}]}},"type":"string"}},"type":"object"}`),
		},
	}
	for _, test := range tests {
		var got json.RawMessage
		var gotJson interface{}
		gotJson = parseReference(source, test.refPaths)
		got, err := json.Marshal(gotJson)

		if err != nil {
			t.Errorf("Test %q - error marshal got\n%s: %v", test.description, gotJson, err)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Test %q - got\n%s\nwant\n%s", test.description, got, test.want)
		}
	}
}
