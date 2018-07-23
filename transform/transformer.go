// Package transform implements code which can use a JSON schema with transform sections to convert a JSON file to
// match the schema format.
package transform

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/GannettDigital/jstransform/jsonschema"

	"github.com/PaesslerAG/jsonpath"
	"github.com/buger/jsonparser"
)

// Transformer uses a JSON schema and the transform sections within it to take a set of JSON and transform it to
// matching the schema.
// More details on the transform section of the schema are found at
// https://github.com/GannettDigital/jstransform/blob/master/transform.adoc
type Transformer struct {
	in                  interface{}
	relativePath        string
	schema              *jsonschema.Schema
	skipPrefix          string
	transformIdentifier string // Used to select the proper transform Instructions
	transformed         map[string]interface{}
}

// NewTransformer returns a Transformer using the schema given.
// The transformIdentifier is used to select the appropriate transform section from the schema.
func NewTransformer(schema *jsonschema.Schema, tranformIdentifier string) (*Transformer, error) {
	return &Transformer{schema: schema, transformIdentifier: tranformIdentifier}, nil
}

// Transform takes the provided JSON and converts the JSON to match the pre-defined JSON Schema using the transform
// sections in the schema.
//
// The Transform operation is not concurrency safe, only one Transform at a time should be performed for any given transformer.
//
// By default fields with no Transform section but with matching path and type are copied verbatim into the new
// JSON structure. Fields which are missing from the input are set to a default value in the output.
//
// Errors are returned for failures to perform operations but are not returned for empty fields which are either
// omitted from the output or set to an empty value.
//
// Validation of the output against the schema is the final step in the process.
func (tr *Transformer) Transform(in json.RawMessage) (json.RawMessage, error) {
	// reset transformed and processed so that each time this is called it repeats the operation
	tr.transformed = make(map[string]interface{})
	if err := json.Unmarshal(in, &tr.in); err != nil {
		return nil, fmt.Errorf("failed to parse input JSON: %v", err)
	}

	if err := jsonschema.WalkRaw(tr.schema, tr.walker); err != nil {
		return nil, err
	}

	out, err := json.Marshal(tr.transformed)
	if err != nil {
		return nil, fmt.Errorf("failed to JSON marsal transformed data: %v", err)
	}

	valid, err := tr.schema.Validate(out)
	if err != nil {
		return nil, fmt.Errorf("transformed result validation error: %v", err)
	}
	if !valid {
		return nil, errors.New("schema validation of the transformed result reports invalid")
	}

	return out, nil
}

// walker is a WalkFunc for the Transformer which does the bulk of the work instance by instance.
// It includes the logic to handle arrays properly.
func (tr *Transformer) walker(path string, value json.RawMessage) error {
	// arrays are processed as a group when encountered as part of the parent item
	if tr.skipPrefix != "" {
		if strings.HasPrefix(path, tr.skipPrefix) {
			path = strings.Replace(path, tr.skipPrefix, tr.relativePath, 1)
		} else {
			return nil
		}
	}
	if strings.Contains(path, "[*]") {
		return nil
	}

	ifields := struct {
		Type      string    `json:"type"`
		Format    string    `json:"format"`
		Transform transform `json:"transform"`
	}{}
	if err := json.Unmarshal(value, &ifields); err != nil {
		return fmt.Errorf("failed to extract transform: %v", err)
	}

	jsonType := ifields.Type
	if jsonType == "string" && ifields.Format == "date-time" {
		jsonType = "date-time"
	}

	if tr.relativePath != "" { // TODO this should also only apply the jsonPath's starting with @ not any
		ifields.Transform = ifields.Transform.replaceJSONPath("@", tr.relativePath)
	}

	newValue, err := tr.getInstanceValue(ifields.Transform, tr.in, jsonType, path, value)
	if err != nil {
		return err
	}

	// Arrays items are processed as a group when the parent is encountered
	if jsonType == "array" {
		if newValue == nil {
			return nil
		}
		newArray, ok := newValue.([]interface{})
		if !ok {
			newArray = []interface{}{newValue}
		}
		items, _, _, err := jsonparser.Get(value, "items")
		if err != nil {
			return fmt.Errorf("failed parsing array items at path %q: %v", path, err)
		}
		newValue, err = tr.processArrayItems(path, newArray, items, value)
		if err != nil {
			return fmt.Errorf("failed processing array items at path %q: %v", path, err)
		}
	}

	return tr.saveValue(path, newValue)
}

// processArrayItems handles the walker processing of Array items. These are different because the new array items
// are build based on the transformed data from the array instance and for each field in an array item processing of
// field for all array items happens in one step. This function can recursively handle nested arrays.
func (tr *Transformer) processArrayItems(path string, arraySrc []interface{}, rawSchema json.RawMessage, value json.RawMessage) ([]interface{}, error) {
	atrIn, ok := tr.in.(map[string]interface{})
	if !ok {
		atrIn = make(map[string]interface{})
	}
	if err := saveInTree(atrIn, path[2:], arraySrc); err != nil {
		return nil, fmt.Errorf("failed to initialize array walker: %v", err)
	}

	var newArray []interface{}

	for i := range arraySrc {
		atr := &Transformer{
			in:                  atrIn,
			relativePath:        fmt.Sprintf("%s[%d]", path, i),
			schema:              tr.schema,
			skipPrefix:          fmt.Sprintf("%s[*]", replaceIndex(path)),
			transformIdentifier: tr.transformIdentifier,
			transformed:         make(map[string]interface{}),
		}

		if err := jsonschema.WalkRaw(tr.schema, atr.walker); err != nil {
			return nil, err
		}
		if len(atr.transformed) != 0 {
			arrayValue, err := jsonpath.Get(fmt.Sprintf("%s[%d]", path, i), atr.transformed)
			if err != nil {
				continue
			}
			newArray = append(newArray, arrayValue)
		}
	}
	return newArray, nil

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
		return saveLeaf(tree, splits[0], value)
	}

	arraySplits := strings.SplitN(splits[0], "[", 2)
	if len(arraySplits) != 1 {
		index, err := strconv.Atoi(strings.Trim(arraySplits[1], "]"))
		if err != nil {
			return fmt.Errorf("failed ot determine index of %q", splits[0])
		}
		var slices []interface{}
		rawSlice, ok := tree[arraySplits[0]]
		if !ok {
			slices = make([]interface{}, 0)
		} else {
			slices, ok = rawSlice.([]interface{})
			if !ok {
				return fmt.Errorf("value at %q is not a slice of interface{}", arraySplits[0])
			}
		}

		for i := len(slices); i <= index; i++ {
			slices = append(slices, make(map[string]interface{}))
		}

		tree[arraySplits[0]] = slices
		sliceMap, ok := slices[index].(map[string]interface{})
		if !ok {
			return fmt.Errorf("value at %q is not a slice of map[string]interface{}", arraySplits[0])
		}

		return saveInTree(sliceMap, strings.Join(splits[1:], "."), value)
	}
	newTree, ok := tree[splits[0]]
	if ok {
		var newTreeMap map[string]interface{}
		if newTree == nil {
			newTreeMap = make(map[string]interface{})
			tree[splits[0]] = newTreeMap
		} else {
			var ok bool
			newTreeMap, ok = newTree.(map[string]interface{})
			if !ok {
				return fmt.Errorf("value at %q is not a map[string]interface{}", splits[0])
			}
		}
		return saveInTree(newTreeMap, strings.Join(splits[1:], "."), value)
	}
	newTreeMap := make(map[string]interface{})
	tree[splits[0]] = newTreeMap
	return saveInTree(newTreeMap, strings.Join(splits[1:], "."), value)
}

func saveLeaf(tree map[string]interface{}, path string, value interface{}) error {
	arraySplits := strings.Split(path, "[")
	if len(arraySplits) == 1 {
		tree[path] = value
		return nil
	}

	rawSlice, ok := tree[arraySplits[0]]
	if !ok {
		rawSlice = make([]interface{}, 0)
	}
	sValue, ok := rawSlice.([]interface{})
	if !ok {
		return fmt.Errorf("value at %q is not a []interface{}", arraySplits[0])
	}

	for i := len(arraySplits) - 1; i > 0; i-- {
		nestedArray := arraySplits[i]
		index, err := strconv.Atoi(strings.Trim(nestedArray, "]"))
		if err != nil {
			return fmt.Errorf("failed to determine index of %q", path)
		}

		intermediateValue := make([]interface{}, 0)
		if i == 1 {
			intermediateValue = sValue
		}
		for j := len(intermediateValue); j <= index; j++ {
			intermediateValue = append(intermediateValue, nil)
		}

		intermediateValue[index] = value
		value = intermediateValue

		if i == 1 {
			sValue = intermediateValue
		}
		// TODO this handles if the first level of the nested arrays has other values in it, but doesn't handle if
		// some of the nested arrays have other values
	}

	tree[arraySplits[0]] = sValue
	return nil
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
