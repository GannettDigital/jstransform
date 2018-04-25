package jsonschema

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// Instance represents a JSON Schema instance.
type Instance struct {
	Type       string                     `json:"type"`
	Format     string                     `json:"format"`
	Items      json.RawMessage            `json:"items,omitempty"`
	Properties map[string]json.RawMessage `json:"properties,omitempty"`
	Required   []string                   `json:"Required,omitempty"`
}

// schemaRef represents a JSON Schema reference.
type schemaRef struct {
	Ref string `json:"$ref"`
}

// schemaJSON represents the basic supported structure of a JSON Schema file
type schemaJSON struct {
	Instance

	Schema      string      `json:"$schema"`
	Description string      `json:"description,omitempty"`
	AllOf       []schemaRef `json:"allOf,omitempty"`
	OneOf       []schemaRef `json:"oneOf,omitempty"`
	Required    []string    `json:"required,omitempty"`
}

// Schema represents a JSON Schema with the AllOf and OneOf references parsed and squashed into a single representation.
// This is not a fully spec compatible representation but a basic representation useful for walking through the schema
// instances within a schema.
//
// A fully spec compatible version of the schema is kept for validation purposes.
type Schema struct {
	Instance

	validator Validator
}

// Validate will check that the given json is validate according the schema.
func (s *Schema) Validate(raw json.RawMessage) (bool, error) {
	return s.validator.Validate(raw)
}

// SchemaFromFile parses a file at the given path and returns a schema based on its contents.
// The function traverses top level allOf fields within the schema. For oneOf fields the reference base
// name minus any extension is compared to the value of the oneOfType argument and if they match that file is also
// traversed.
//
// Only file references are supported.
//
// Note: A top level array will use the items from the first file to define them.
func SchemaFromFile(schemaPath string, oneOfType string) (*Schema, error) {
	data, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema file %q: %v", schemaPath, err)
	}
	v, err := NewValidator(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize schema validator: %v", err)
	}

	var sj schemaJSON
	if err := json.Unmarshal(data, &sj); err != nil {
		return nil, fmt.Errorf("failed to Unmarshal Schema: %v", err)
	}

	s := Schema{
		Instance: Instance{
			Items:      sj.Items,
			Properties: sj.Properties,
			Required:   sj.Required,
		},
		validator: v,
	}

	for _, all := range sj.AllOf {
		path := refPath(schemaPath, all.Ref)
		child, err := SchemaFromFile(path, oneOfType)
		if err != nil {
			return nil, fmt.Errorf("failed parsing allOf file %q: %v", path, err)
		}
		s.Properties = mergeProperties(s.Properties, child.Properties)
		s.Required = append(s.Required, child.Required...)
		if s.Items == nil {
			s.Items = child.Items
		}
	}

	for _, one := range sj.OneOf {
		subName := strings.Split(filepath.Base(one.Ref), ".")[0]
		if subName != oneOfType {
			continue
		}
		path := refPath(schemaPath, one.Ref)
		child, err := SchemaFromFile(path, oneOfType)
		if err != nil {
			return nil, fmt.Errorf("failed parsing oneOf file %q: %v", path, err)
		}
		s.Properties = mergeProperties(s.Properties, child.Properties)
		s.Required = append(s.Required, child.Required...)
		if s.Items == nil {
			s.Items = child.Items
		}
	}

	return &s, nil
}

// SchemaTypes will parse the given file and report which top level allOfTypes and oneOfTypes are found in that schema.
func SchemaTypes(schemaPath string) ([]string, []string, error) {
	data, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read schema file %q: %v", schemaPath, err)
	}

	var sj schemaJSON
	if err := json.Unmarshal(data, &sj); err != nil {
		return nil, nil, fmt.Errorf("failed to Unmarshal Schema: %v", err)
	}

	var allOfTypes []string
	for _, a := range sj.AllOf {
		allOfTypes = append(allOfTypes, a.Ref)
	}
	var oneOfTypes []string
	for _, one := range sj.OneOf {
		oneOfTypes = append(oneOfTypes, one.Ref)
	}

	return allOfTypes, oneOfTypes, nil
}

func mergeProperties(parent, child map[string]json.RawMessage) map[string]json.RawMessage {
	newProperties := make(map[string]json.RawMessage)
	if parent != nil {
		newProperties = parent
	}
	for key, value := range child {
		if _, ok := parent[key]; !ok {
			newProperties[key] = value
		}
	}
	return newProperties
}

// refPath returns the reference path if absolute or the combo if it and the parent if not.
func refPath(parent, ref string) string {
	if filepath.IsAbs(ref) {
		return ref
	}

	return filepath.Join(filepath.Dir(parent), ref)
}
