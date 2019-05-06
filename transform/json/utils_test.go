package json

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestSchemaDefault(t *testing.T) {
	tests := []struct {
		description string
		schema      json.RawMessage
		want        interface{}
	}{
		{
			description: "String with default",
			schema:      json.RawMessage(`{"type":"string","default":"defaultString"}`),
			want:        "defaultString",
		},
		{
			description: "Number with default",
			schema:      json.RawMessage(`{"type":"number","default":3.14}`),
			want:        3.14,
		},
		{
			description: "String without default",
			schema:      json.RawMessage(`{"type":"string"}`),
			want:        nil,
		},
		{
			description: "Number without default",
			schema:      json.RawMessage(`{"type":"number"}`),
			want:        nil,
		},
		{
			description: "object without default",
			schema:      json.RawMessage(`{"type":"object"}`),
			want:        nil,
		},
		{
			description: "array without default",
			schema:      json.RawMessage(`{"type":"array"}`),
			want:        nil,
		},
	}

	for _, test := range tests {
		got, err := schemaDefault(test.schema)
		if err != nil {
			t.Errorf("Test %q - got error: %v", test.description, err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Test %q - got %v, want %v", test.description, got, test.want)
		}
	}
}

func TestReplaceIndex(t *testing.T) {
	tests := []struct {
		description string
		path        string
		want        string
	}{
		{
			description: "Simple",
			path:        "a[0]",
			want:        "a[*]",
		},
		{
			description: "Array of objects",
			path:        "a[2].b",
			want:        "a[*].b",
		},
		{
			description: "Two Arrays",
			path:        "a[3].b[4]",
			want:        "a[*].b[*]",
		},
		{
			description: "No Arrays",
			path:        "a.b",
			want:        "a.b",
		},
		{
			description: "Already changed",
			path:        "a[*].b[*]",
			want:        "a[*].b[*]",
		},
		{
			description: "Complicated",
			path:        "a[3].ab.b[4].c.d[8]",
			want:        "a[*].ab.b[*].c.d[*]",
		},
		{
			description: "Large Index",
			path:        "a[13].b[4454]",
			want:        "a[*].b[*]",
		},
	}

	for _, test := range tests {
		if got, want := replaceIndex(test.path), test.want; got != want {
			t.Errorf("Test %q - got %q, want %q", test.description, got, want)
		}
	}
}
