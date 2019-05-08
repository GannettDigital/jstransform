// Package transform implements code which can use a JSON schema with transform sections to convert a JSON file to
// match the schema format.
package transform

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/antchfx/xmlquery"

	"github.com/GannettDigital/jstransform/jsonschema"

	"github.com/buger/jsonparser"
)

// inputFormat denotes the type of transform to perfrom, the options are 'JSON' or 'XML'
type inputFormat string

const (
	jsonInput = inputFormat("JSON")
	xmlInput  = inputFormat("XML")
)

// JSONTransformer - a type implemented by the jstransform.Transformer
type JSONTransformer interface {
	Transform(raw json.RawMessage) (json.RawMessage, error)
}

// Transformer uses a JSON schema and the transform sections within it to take a set of JSON and transform it to
// matching the schema.
// More details on the transform section of the schema are found at
// https://github.com/GannettDigital/jstransform/blob/master/transform.adoc
type Transformer struct {
	schema              *jsonschema.Schema
	transformIdentifier string // Used to select the proper transform Instructions
	root                instanceTransformer
	format              inputFormat
}

// NewTransformer returns a Transformer using the schema given.
// The transformIdentifier is used to select the appropriate transform section from the schema.
// It expects the transforms to be performed on JSON data
func NewTransformer(schema *jsonschema.Schema, tranformIdentifier string) (*Transformer, error) {
	return newTransformer(schema, tranformIdentifier, jsonInput)
}

// NewXMLTransformer returns a Transformer using the schema given.
// The transformIdentifier is used to select the appropriate transform section from the schema.
// It expects the transforms to be performed on XML data
func NewXMLTransformer(schema *jsonschema.Schema, tranformIdentifier string) (*Transformer, error) {
	return newTransformer(schema, tranformIdentifier, xmlInput)
}

func newTransformer(schema *jsonschema.Schema, tranformIdentifier string, format inputFormat) (*Transformer, error) {
	tr := &Transformer{schema: schema, transformIdentifier: tranformIdentifier, format: format}
	emptyJSON := []byte(`{}`)
	var err error
	if schema.Properties != nil {
		tr.root, err = newObjectTransformer("$", tranformIdentifier, emptyJSON, format)
	} else if schema.Items != nil {
		tr.root, err = newArrayTransformer("$", tranformIdentifier, emptyJSON, format)
	} else {
		return nil, errors.New("no Properties nor Items found for schema")
	}
	if err != nil {
		return nil, fmt.Errorf("failed initializing root transformer: %v", err)
	}

	if err := jsonschema.WalkRaw(schema, tr.walker); err != nil {
		return nil, err
	}

	return tr, nil
}

// Transform takes the provided JSON and converts the JSON to match the pre-defined JSON Schema using the transform
// sections in the schema.
//
// By default fields with no Transform section but with matching path and type are copied verbatim into the new
// JSON structure. Fields which are missing from the input are set to a default value in the output.
//
// Errors are returned for failures to perform operations but are not returned for empty fields which are either
// omitted from the output or set to an empty value.
//
// Validation of the output against the schema is the final step in the process.
func (tr *Transformer) Transform(raw json.RawMessage) (json.RawMessage, error) {
	if tr.format == jsonInput {
		return tr.jsonTransform(raw)
	}
	if tr.format == xmlInput {
		return tr.xmlTransform(raw)
	}
	return nil, fmt.Errorf("unknown transform type %s, must be 'JSON' or 'XML'", tr.format)
}

func (tr *Transformer) jsonTransform(raw json.RawMessage) (json.RawMessage, error) {
	var in interface{}
	if err := json.Unmarshal(raw, &in); err != nil {
		return nil, fmt.Errorf("failed to parse input JSON: %v", err)
	}

	transformed, err := tr.root.transform(in, nil)
	if err != nil {
		return nil, fmt.Errorf("failed transformation: %v", err)
	}

	out, err := json.Marshal(transformed)
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

func (tr *Transformer) xmlTransform(raw []byte) ([]byte, error) {
	xmlDoc, err := xmlquery.Parse(bytes.NewReader(raw))
	if err != nil {
		return nil, fmt.Errorf("failed to parse input XML: %v", err)
	}

	transformed, err := tr.root.transform(xmlDoc, nil)
	if err != nil {
		return nil, fmt.Errorf("failed transformation: %v", err)
	}

	out, err := json.Marshal(transformed)
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

// findParent walks the instanceTransformer tree to find the parent of the given path
func (tr *Transformer) findParent(path string) (instanceTransformer, error) {
	path = strings.Replace(path, "[", ".[", -1)
	splits := strings.Split(path, ".")
	if splits[0] != "$" {
		// TODO this will probably choke on a root level array
		return nil, errors.New("paths must start with '$'")
	}
	parentSplits := splits[1 : len(splits)-1]

	parent := tr.root
	for _, sp := range parentSplits {
		if sp == "[*]" {
			parent = parent.child()
			continue
		}

		parent = parent.selectChild(sp)
	}

	return parent, nil
}

// walker is a WalkFunc for the Transformer which builds an representation of the fields and transforms in the schema.
// This is later used to do the actual transform for incoming data
func (tr *Transformer) walker(path string, value json.RawMessage) error {
	instanceType, err := jsonparser.GetString(value, "type")
	if err != nil {
		return fmt.Errorf("failed to extract instance type: %v", err)
	}

	var iTransformer instanceTransformer
	switch instanceType {
	case "object":
		iTransformer, err = newObjectTransformer(path, tr.transformIdentifier, value, tr.format)
	case "array":
		iTransformer, err = newArrayTransformer(path, tr.transformIdentifier, value, tr.format)
	default:
		iTransformer, err = newScalarTransformer(path, tr.transformIdentifier, value, instanceType, tr.format)
	}
	if err != nil {
		return fmt.Errorf("failed to initialize transformer: %v", err)
	}

	parent, err := tr.findParent(path)
	if err != nil {
		return err
	}
	if err := parent.addChild(iTransformer); err != nil {
		return err
	}

	return nil
}

// saveInTree is used recursively to add values the tree based on the path even if the parents are nil.
func saveInTree(tree map[string]interface{}, path string, value interface{}) error {
	if value == nil {
		return nil
	}

	splits := strings.Split(path, ".")
	if splits[0] == "$" {
		path = path[2:]
		splits = splits[1:]
	}

	if len(splits) == 1 {
		return saveLeaf(tree, splits[0], value)
	}

	arraySplits := strings.Split(splits[0], "[")
	if len(arraySplits) != 1 { // the case of an array or nested arrays with an object in them
		var sValue []interface{}
		if rawSlice, ok := tree[arraySplits[0]]; ok {
			sValue = rawSlice.([]interface{})
		}

		newTreeMap := make(map[string]interface{})
		newValue, err := saveInSlice(sValue, arraySplits[1:], newTreeMap)
		if err != nil {
			return err
		}

		tree[arraySplits[0]] = newValue
		return saveInTree(newTreeMap, strings.Join(splits[1:], "."), value)
	}

	var newTreeMap map[string]interface{}
	newTree, ok := tree[splits[0]]
	if !ok || newTree == nil {
		newTreeMap = make(map[string]interface{})
	} else {
		newTreeMap, ok = newTree.(map[string]interface{})
		if !ok {
			return fmt.Errorf("value at %q is not a map[string]interface{}", splits[0])
		}
	}
	tree[splits[0]] = newTreeMap
	return saveInTree(newTreeMap, strings.Join(splits[1:], "."), value)
}

// saveLeaf will save a leaf value in the tree at the given path. If the path specifies an array or set of nested
// arrays it will build the array items as needed to reach the specified index. New array items are created as nil.
// Any nested array items will be recursively treated the same way.
func saveLeaf(tree map[string]interface{}, path string, value interface{}) error {
	arraySplits := strings.Split(path, "[")
	if len(arraySplits) == 1 {
		tree[path] = value
		return nil
	}

	var sValue []interface{}
	if rawSlice, ok := tree[arraySplits[0]]; ok {
		sValue = rawSlice.([]interface{})
	}

	newValue, err := saveInSlice(sValue, arraySplits[1:], value)
	if err != nil {
		return err
	}
	tree[arraySplits[0]] = newValue
	return nil
}

func saveInSlice(current []interface{}, arraySplits []string, value interface{}) ([]interface{}, error) {
	index, err := strconv.Atoi(strings.Trim(arraySplits[0], "]"))
	if err != nil {
		return nil, fmt.Errorf("failed to determine index of %q", arraySplits[0])
	}

	if current == nil {
		current = make([]interface{}, 0, index)
	}

	// fill up the slice slots with nil if the slice isn't the right size
	for j := len(current); j <= index; j++ {
		current = append(current, nil)
	}

	if len(arraySplits) == 1 {
		// if this is the last array split save the value and break
		if newValue, ok := value.(map[string]interface{}); ok { // special case combine existing values into new value if a map
			if oldValue, ok := current[index].(map[string]interface{}); ok {
				for k, v := range oldValue {
					if _, ok := newValue[k]; !ok {
						newValue[k] = v
					}
				}
				value = newValue
			}
		}
		current[index] = value
		return current, nil
	}

	// recurse as needed
	nested, ok := current[index].([]interface{})
	if !ok {
		nested = nil
	}

	newValue, err := saveInSlice(nested, arraySplits[1:], value)
	current[index] = newValue
	return current, nil
}
