package jsonschema

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/PaesslerAG/jsonpath"
)

// Transformer uses a JSON schema and the transform sections within it to take a set of JSON and transform it to
// matching the schema.
// More details on the transform section of the schema are found at
// https://github.com/GannettDigital/content-schema/blob/master/v1/assets/README.md
//
// Note: Root level arrays in the JSON output are not currently supported
type Transformer struct {
	in                  interface{}
	processedArrays     map[string][]interface{}
	schema              *Schema
	transformIdentifier string // Used to select the proper transform Instructions
	transformed         map[string]interface{}
}

// NewTransformer returns a Transformer using the schema given.
// The transformIdentifier is used to select the appropriate transform section from the schema.
func NewTransformer(schema *Schema, tranformIdentifier string) (*Transformer, error) {
	return &Transformer{schema: schema, transformIdentifier: tranformIdentifier}, nil
}

// Transform takes the provided JSON and converts the JSON to match the pre-defined JSON Schema using the transform
// sections in the schema.
// By default fields with no Transform section but with matching path and type are copied verbatim into the new
// JSON structure. Fields which are missing from the input are set to a default value in the output.
//
// Errors are returned for failures to perform operations but are not returned for empty fields which are either
// omitted from the output or set to an empty value.
//
// Validation of the output against the schema is not done and is the responsibility of the caller.
func (tr *Transformer) Transform(in json.RawMessage) (json.RawMessage, error) {
	// reset transformed and processed so that each time this is called it repeats the operation
	tr.transformed = make(map[string]interface{})
	tr.processedArrays = make(map[string][]interface{})
	if err := json.Unmarshal(in, &tr.in); err != nil {
		return nil, fmt.Errorf("failed to parse input JSON: %v", err)
	}

	if err := Walk(tr.schema, tr.walker); err != nil {
		return nil, err
	}

	out, err := json.Marshal(tr.transformed)
	if err != nil {
		return nil, fmt.Errorf("failed to JSON marsal transformed data: %v", err)
	}

	return out, nil
}

// walker is a WalkFunc for the Transformer which does the bulk of the work instance by instance.
// It includes the logic to handle arrays properly.
func (tr *Transformer) walker(path string, value json.RawMessage) error {
	ifields := struct {
		Type      string    `json:"type"`
		Transform transform `json:"transform"`
	}{}
	if err := json.Unmarshal(value, &ifields); err != nil {
		return fmt.Errorf("failed to extract transform: %v", err)
	}

	// For array items process every item in the array at the same time
	for key, arraySrc := range tr.processedArrays {
		if strings.Contains(path, key) {
			relativePath := strings.Replace(path, key+"[*]", "", 1)
			rawArray, err := jsonpath.Get(key, tr.transformed)
			if err != nil {
				return err
			}
			array, ok := rawArray.([]interface{})
			if !ok {
				return fmt.Errorf("expected array in transformed object at path %q", key)
			}
			if err := tr.processArrayItems(relativePath, arraySrc, array, ifields.Type, ifields.Transform, value); err != nil {
				return fmt.Errorf("failed processing array items at path %q: %v", path, err)
			}
			return nil
		}
	}

	newValue, err := tr.getInstanceValue(ifields.Transform, tr.in, ifields.Type, path, value)
	if err != nil {
		return err
	}

	// Arrays are processed as a whole independent, save the to be used when processing the subfields and write an
	// empty array to the output which will be filled in by the subfield calls
	if ifields.Type == "array" {
		newArray, ok := newValue.([]interface{})
		if !ok {
			newArray = []interface{}{}
			if newValue != nil {
				newArray = []interface{}{newValue}
			} else {
				newArray = []interface{}{}
			}
		}
		tr.processedArrays[path] = newArray
		return tr.saveValue(path, make([]interface{}, len(newArray)))
	}

	return tr.saveValue(path, newValue)
}

// processArrayItems handles the walker processing of Array items. These are different because the new array items
// are build based on the transformed data from the array instance and for each field in an array item processing of
// field for all array items happens in one step.
func (tr *Transformer) processArrayItems(relativePath string, arraySrc []interface{}, array []interface{}, jsonType string, instanceTransform transform, value json.RawMessage) error {
	for i := range array {
		itemTransform := instanceTransform.replaceJSONPath("@", fmt.Sprintf("$[%d]", i))

		newValue, err := tr.getInstanceValue(itemTransform, arraySrc, jsonType, fmt.Sprintf("$[%d]%s", i, relativePath), value)
		if err != nil {
			return err
		}

		if trimmed := strings.TrimLeft(relativePath, "."); len(trimmed) > 0 {
			newMap, ok := array[i].(map[string]interface{})
			if !ok {
				return errors.New("expected map[string]interface{} array items")
			}
			if err := saveInTree(newMap, trimmed, newValue); err != nil {
				return err
			}
		} else {
			array[i] = newValue
		}
	}
	return nil

}

// processTransform determines the value for a given instance using a transform, returning nil if there is no value
// determined.
func (tr *Transformer) processTransform(t transform, in interface{}, jsonType string) (interface{}, error) {
	if t == nil {
		return nil, nil
	}

	instructions, found := t[tr.transformIdentifier]
	if !found {
		return nil, nil
	}

	newValue, err := instructions.transform(in, jsonType)
	if err != nil {
		return nil, err
	}
	return newValue, nil
}

// saveValue adds the given value to the tr.transformed object at the place specified by jsonPath.
// It does not support writing a value to an array which is inside another array, rather it assumes all array
// members are saved as a complete whole.
func (tr *Transformer) saveValue(jsonPath string, value interface{}) error {
	splits := strings.SplitN(jsonPath, ".", 2)
	if splits[0] != "$" {
		return errors.New("all JSONPaths are required to start at '$'")
	}
	return saveInTree(tr.transformed, splits[1], value)
}

// saveInTree is used recursively to accomplish the work of saveValue.
func saveInTree(tree map[string]interface{}, path string, value interface{}) error {
	if value == nil {
		return nil
	}
	splits := strings.Split(path, ".")
	if len(splits) == 1 {
		arraySplits := strings.SplitN(splits[0], "[", 2)
		if len(arraySplits) == 1 { // leaf value save it and return
			tree[splits[0]] = value
			return nil
		}

		rawSlice, ok := tree[arraySplits[0]]
		if !ok {
			tree[arraySplits[0]] = []interface{}{value}
			return nil
		}
		sValue, ok := rawSlice.([]interface{})
		if !ok {
			return fmt.Errorf("value at %q is not a []interface{}", arraySplits[0])
		}
		tree[arraySplits[0]] = append(sValue, value)
		return nil
	}
	newTree, ok := tree[splits[0]]
	if ok {
		newTreeMap, ok := newTree.(map[string]interface{})
		if !ok {
			return fmt.Errorf("value at %q is not a map[string]interface{}", splits[0])
		}
		return saveInTree(newTreeMap, strings.Join(splits[1:], "."), value)
	}
	newTreeMap := make(map[string]interface{})
	tree[splits[0]] = newTreeMap
	return saveInTree(newTreeMap, strings.Join(splits[1:], "."), value)
}

// getInstanceValue retrieves the value for a JSONSchema instance following this process:
//
// 1. Use a Transform if it exists.
//
// 2. Look for the same JSONPath in the input and use directly if possible.
//
// 3. Fall back to the JSON Schema default value.
func (tr *Transformer) getInstanceValue(t transform, in interface{}, jsonType string, path string, value json.RawMessage) (interface{}, error) {
	// 1. Use a transform if it exists
	newValue, err := tr.processTransform(t, in, jsonType)
	if err != nil {
		return nil, err
	}

	// 2. Look for the same JSONPath in the input and use directly if possible.
	if newValue == nil && jsonType != "object" {
		rawValue, err := jsonpath.Get(path, in)
		if err == nil {
			newValue, err = convert(rawValue, jsonType)
			if err != nil {
				return nil, err
			}
		}
	}

	// 3. Fall back to the JSON Schema default value.
	if newValue == nil {
		newValue, err = schemaDefault(value)
		if err != nil {
			return nil, err
		}
	}
	return newValue, nil
}
