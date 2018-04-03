package transform

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestConcat(t *testing.T) {
	tests := []struct {
		description string
		a           interface{}
		b           interface{}
		want        interface{}
		wantErr     bool
	}{
		{
			description: "type mismatch",
			a:           0,
			b:           "",
			wantErr:     true,
		},
		{
			description: "string concat",
			a:           "con",
			b:           "cat",
			want:        "concat",
		},
		{
			description: "int concat",
			a:           0,
			b:           1,
			wantErr:     true,
		},
		{
			description: "both nil",
			a:           nil,
			b:           nil,
			want:        nil,
		},
		{
			description: "a nil",
			a:           nil,
			b:           1,
			want:        1,
		},
		{
			description: "b nil",
			a:           0,
			b:           nil,
			want:        0,
		},
	}

	for _, test := range tests {
		got, err := concat(test.a, test.b)

		switch {
		case test.wantErr && err != nil:
			continue
		case test.wantErr && err == nil:
			t.Errorf("Test %q - got nil, want error", test.description)
		case !test.wantErr && err != nil:
			t.Errorf("Test %q - got error, want nil: %v", test.description, err)
		case !reflect.DeepEqual(got, test.want):
			t.Errorf("Test %q - got %v, want %v", test.description, got, test.want)
		}
	}
}

func TestConvert(t *testing.T) {
	tests := []struct {
		description string
		raw         interface{}
		jsonType    string
		want        interface{}
		wantErr     bool
	}{
		{
			description: "bool -> bool",
			raw:         true,
			jsonType:    "boolean",
			want:        true,
		},
		{
			description: "nil -> bool",
			raw:         nil,
			jsonType:    "boolean",
			want:        nil,
		},
		{
			description: "bool -> number, true",
			raw:         true,
			jsonType:    "number",
			want:        1,
		},
		{
			description: "bool -> number, false",
			raw:         false,
			jsonType:    "number",
			want:        0,
		},
		{
			description: "bool -> string",
			raw:         true,
			jsonType:    "string",
			want:        "true",
		},
		{
			description: "float -> bool, true",
			raw:         3.14,
			jsonType:    "boolean",
			want:        true,
		},
		{
			description: "float -> bool, false",
			raw:         -3.14,
			jsonType:    "boolean",
			want:        false,
		},
		{
			description: "float -> number",
			raw:         3.14,
			jsonType:    "number",
			want:        3.14,
		},
		{
			description: "float -> string",
			raw:         3.14,
			jsonType:    "string",
			want:        "3.14",
		},
		{
			description: "int -> bool, true",
			raw:         2,
			jsonType:    "boolean",
			want:        true,
		},
		{
			description: "int -> bool, false",
			raw:         0,
			jsonType:    "boolean",
			want:        false,
		},
		{
			description: "int -> number",
			raw:         2,
			jsonType:    "number",
			want:        2,
		},
		{
			description: "int -> string",
			raw:         2,
			jsonType:    "string",
			want:        "2",
		},
		{
			description: "string -> bool",
			raw:         "true",
			jsonType:    "boolean",
			want:        true,
		},
		{
			description: "empty string -> bool",
			raw:         "",
			jsonType:    "boolean",
			want:        false,
		},
		{
			description: "string -> bool - should error",
			raw:         "random",
			jsonType:    "boolean",
			wantErr:     true,
		},
		{
			description: "string -> number - int",
			raw:         "4",
			jsonType:    "number",
			want:        4,
		},
		{
			description: "string -> number - float",
			raw:         "4.3",
			jsonType:    "number",
			want:        4.3,
		},
		{
			description: "string -> number - should error",
			raw:         "hello",
			jsonType:    "number",
			wantErr:     true,
		},
		{
			description: "string -> string",
			raw:         "hello",
			jsonType:    "string",
			want:        "hello",
		},
		{
			description: "nil -> string",
			raw:         nil,
			jsonType:    "string",
			want:        nil,
		},
		{
			description: "[]interface{} with a single string in it -> string",
			raw:         []interface{}{"hello"},
			jsonType:    "string",
			want:        "hello",
		},
		{
			description: "[]interface{} with a single bool in it -> bool",
			raw:         []interface{}{true},
			jsonType:    "boolean",
			want:        true,
		},
		{
			description: "[]interface{} with a single number in it -> number",
			raw:         []interface{}{4.3},
			jsonType:    "number",
			want:        4.3,
		},
		{
			description: "[]interface{} returns nil",
			raw:         []interface{}{},
			jsonType:    "string",
			want:        nil,
		},
	}

	for _, test := range tests {
		got, err := convert(test.raw, test.jsonType)
		switch {
		case test.wantErr && err != nil:
			continue
		case test.wantErr && err == nil:
			t.Errorf("Test %q - got nil, want error", test.description)
		case !test.wantErr && err != nil:
			t.Errorf("Test %q - got error, want nil: %v", test.description, err)
		case !reflect.DeepEqual(got, test.want):
			t.Errorf("Test %q - got %v, want %v", test.description, got, test.want)
		}
	}
}

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
			want:        map[string]interface{}{},
		},
		{
			description: "array without default",
			schema:      json.RawMessage(`{"type":"array"}`),
			want:        []interface{}{},
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
