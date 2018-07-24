// Package jsonschema includes tools for walking JSON schema files and running a custom function for each instance.
package jsonschema

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/buger/jsonparser"
)

// WalkFunc processes a single Instance within a JSON schema file returning an error on any problems.
// The path corresponds to the JSONPath (http://goessner.net/articles/JsonPath/) of the instance within the JSON
// format described by the JSON Schema.
type WalkInstanceFunc func(path string, i Instance) error

// WalkRawFunc works similar to WalkFunc but rather than accepting an instance it accepts the raw JSON for each
// level of the schema.
type WalkRawFunc func(path string, value json.RawMessage) error

// Walk runs each instance in the JSON schema through the defined walk function, keeping track of the JSONPath.
// It assumes the JSON schema is valid and does not check for many errors such as bad property names.
// For instances that are objects or arrays the WalkFunc will be called for each child instance.
func Walk(s *Schema, walkFn WalkInstanceFunc) error {
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
func walkInstance(raw json.RawMessage, path string, walkFn WalkInstanceFunc) error {
	var i Instance
	if err := json.Unmarshal(raw, &i); err != nil {
		return fmt.Errorf("failed to unmarshal Instance at path %q: %v", path, err)
	}

	if err := walkFn(path, i); err != nil {
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

// WalkRaw works nearly identical to the Walk function but rather than calling a WalkFunc calls a WalkRawFunc.
// By skipping the JSON unmarshaling of the Instance it runs nearly 10 times faster then WalkFunc.
func WalkRaw(s *Schema, walkFn WalkRawFunc) error {
	rootPath := "$"

	for key, value := range s.Properties {
		if err := walkRaw(value, prependJSONPath(rootPath, key), walkFn); err != nil {
			return err
		}
	}

	if s.Items != nil {
		if err := walkRaw(s.Items, rootPath+"[*]", walkFn); err != nil {
			return err
		}
	}
	return nil
}

// walkRaw is similar to walkInstance, it is the recursive file which drives walkRaw as walkInstance is the recursive
// function with drives Walk.
func walkRaw(raw json.RawMessage, path string, walkFn WalkRawFunc) error {
	if err := walkFn(path, raw); err != nil {
		return fmt.Errorf("walkFn failed at path %q: %v", path, err)
	}

	iType, err := jsonparser.GetUnsafeString(raw, "type")
	if err != nil {
		return fmt.Errorf("failed to determine instance type at path %q: %v", path, err)
	}

	switch iType {
	case "object":
		if err := jsonparser.ObjectEach(raw, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			return walkRaw(value, prependJSONPath(path, string(key)), walkFn)
		}, "properties"); err != nil {
			return fmt.Errorf("failed processing properties at path %q: %v", path, err)
		}
	case "array":
		items, _, _, err := jsonparser.Get(raw, "items")
		if err != nil {
			return fmt.Errorf("failed extracting items at path %q: %v", path, err)
		}
		if err := walkRaw(items, path+"[*]", walkFn); err != nil {
			return fmt.Errorf("failed processing items at path %q: %v", path, err)
		}
	}
	return nil
}
