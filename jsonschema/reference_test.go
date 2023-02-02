package jsonschema

import (
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
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
		{
			description: "Sports Matchup Example",
			schemaPath:  "./test_data/reference/jsref_matchup-mlb.json",
		},
		{
			description: "Keep Source Transform when Referencing",
			schemaPath:  "./test_data/reference/jsref_keep-transform.json",
		},
	}
	for _, test := range tests {
		var want interface{}
		if !test.wantErr {
			wantPath := strings.Replace(test.schemaPath, "/jsref_", "/deref_", 1)
			wantJson, err := os.ReadFile(wantPath)
			if err != nil {
				t.Errorf("Test %q - failed to read json want file %q: %v", test.description, wantPath, err)
			}
			json.Unmarshal(wantJson, &want)
		}

		var got interface{}
		gotJson, err := os.ReadFile(test.schemaPath)
		if err != nil {
			t.Errorf("Test %q - failed to read json got file %q: %v", test.description, test.schemaPath, err)
		}

		gotJson, err = dereference(test.schemaPath, gotJson, "")
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
