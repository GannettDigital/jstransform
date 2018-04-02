// Package jsonschema includes tools for walking JSON schema files and running a custom function for each instance.
package jsonschema

import (
	"encoding/json"
	"fmt"
	"strings"
)

// WalkFunc processes a single Instance within a JSON schema file returning an error on any problems.
// The path corresponds to the JSONPath (http://goessner.net/articles/JsonPath/) of the instance within the JSON
// format described by the JSON Schema.
type WalkFunc func(path string, i Instance, value json.RawMessage) error

// Walk runs each instance in the JSON schema through the defined walk function, keeping track of the JSONPath.
// It assumes the JSON schema is valid and does not check for many errors such as bad property names.
// For instances that are objects or arrays the WalkFunc will be called for each child instance.
func Walk(s *Schema, walkFn WalkFunc) error {
	rootPath := "$"

	for key, value := range s.Properties {
		if err := walkInstance(value, prependJSONPath(rootPath, key), walkFn); err != nil {
			return err
		}
	}

	if s.Items != nil {
		if err := walkInstance(s.Items, rootPath+"[*]", walkFn); err != nil {
			return err
		}
	}
	return nil
}

// prependJSONPath a parent JSONPath to the beginning of a JSONPath.
// This allows for incrementally building up the full JSONPath.
func prependJSONPath(parent string, child string) string {
	newPath := parent
	if parent == "" {
		newPath = child
	}
	if parent != "" && child != "" {
		newPath = strings.Join([]string{parent, child}, ".")
	}
	return newPath
}

// walkInstance will recursively walk JSON Schema Instance calling defined walk functions for the root and for each
// property within an object and each item in an array.
// For each layer of depth the the JSONPath is added to so the walkFn received the full path of the JSON instance.
// Any error will halt the progress.
func walkInstance(raw json.RawMessage, path string, walkFn WalkFunc) error {
	var i Instance
	if err := json.Unmarshal(raw, &i); err != nil {
		return fmt.Errorf("failed to unmarshal Instance at path %q: %v", path, err)
	}

	if err := walkFn(path, i, raw); err != nil {
		return fmt.Errorf("walkFn failed at path %q: %v", path, err)
	}

	switch i.Type {
	case "object":
		if i.Properties == nil {
			return fmt.Errorf("object at path %q missing Properties", path)
		}
		for key, value := range i.Properties {
			if err := walkInstance(value, prependJSONPath(path, key), walkFn); err != nil {
				return err
			}
		}
	case "array":
		if i.Items == nil {
			return fmt.Errorf("array at path %q missing Items", path)
		}
		if i.Items != nil {
			if err := walkInstance(i.Items, path+"[*]", walkFn); err != nil {
				return err
			}
		}
	}
	return nil
}
