package transform

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	jsonpath "github.com/GannettDigital/PaesslerAG_jsonpath"
	"github.com/GannettDigital/jsonparser"

	"github.com/antchfx/xmlquery"
)

// pathModifier is used to modify the JSON path of an instance to indicate.
type pathModifier func(string) string

func pathReplace(old, new string, modifier pathModifier) pathModifier {
	return func(path string) string {
		if modifier != nil {
			path = modifier(path)
		}
		path = strings.Replace(path, old, new, 1)
		return path
	}
}

// instanceTransformer represents a JSON schema instance with transform details in it.
// The primary function it performs is to transform (or not) data given as input.
// In some cases instances contain other instances, in such a case the children transformers are called as needed to
// build up each value depth first.
type instanceTransformer interface {
	addChild(instanceTransformer) error
	child() instanceTransformer // Arrays return a child object all others nil
	path() string
	selectChild(string) instanceTransformer // This returns nil for everything except objects
	transform(interface{}, pathModifier) (interface{}, error)
}

// arrayTransformer represents a JSON instance type array in the case of a JSON transform or an array of xmlquery.Node in the case of an XML transform.
// in both cases the output will be JSON.
type arrayTransformer struct {
	childTransformer instanceTransformer
	defaultValue     []interface{}
	jsonPath         string
	format           inputFormat
	transforms       *transformInstructions
}

func newArrayTransformer(path, transformIdentifier string, raw json.RawMessage, format inputFormat) (*arrayTransformer, error) {
	at := &arrayTransformer{
		jsonPath: path,
		format:   format,
	}

	var err error
	at.transforms, err = extractTransformInstructions(raw, transformIdentifier, path, "array")
	if err != nil {
		return nil, err
	}

	rawDefault, err := schemaDefault(raw)
	if err != nil {
		return nil, err
	}
	if rawDefault != nil {
		var ok bool
		at.defaultValue, ok = rawDefault.([]interface{})
		if !ok {
			return nil, fmt.Errorf("default value for path %q is not an array", path)
		}
	}

	return at, nil
}

func (at *arrayTransformer) addChild(child instanceTransformer) error {
	at.childTransformer = child
	return nil
}

func (at *arrayTransformer) baseValueJSON(in interface{}, path string, modifier pathModifier) ([]interface{}, bool, error) {
	// 1. Use a transform if it exists
	if at.transforms != nil {
		rawValue, err := at.transforms.transform(in, "array", modifier, at.format)
		if err != nil {
			return nil, false, err
		}
		if rawValue != nil {
			newValue, ok := rawValue.([]interface{})
			if !ok {
				newValue = []interface{}{rawValue}
			}
			return newValue, true, nil
		}
	}

	// 2. Look for the same jsonPath in the input and use directly if possible.
	rawValue, err := jsonpath.Get(path, in)
	if err == nil && rawValue != nil {
		newValue, ok := rawValue.([]interface{})
		if !ok {
			newValue = []interface{}{rawValue}
		}
		return newValue, false, nil
	}

	// 3. Fall back to the JSON Schema default value.
	if at.defaultValue != nil {
		return at.defaultValue, true, nil
	}
	return nil, false, nil
}

func (at *arrayTransformer) baseValueXML(in interface{}, path string, modifier pathModifier) ([]interface{}, bool, error) {
	// 1. Use a transform if it exists
	if at.transforms != nil {
		rawValue, err := at.transforms.transform(in, "array", modifier, at.format)
		if err != nil {
			return nil, false, err
		}

		// if rawValue is an array of xml nodes we need to append them to newValue for return as []interface{}
		xmlNodeArray, ok := rawValue.([]*xmlquery.Node)
		if ok {
			newValue := make([]interface{}, len(xmlNodeArray))
			for i, item := range xmlNodeArray {
				newValue[i] = item
			}
			return newValue, false, nil
		}

		if rawValue != nil {
			newValue, ok := rawValue.([]interface{})
			if !ok {
				newValue = []interface{}{rawValue}
			}
			return newValue, true, nil
		}
	}

	// 2. Fall back to the JSON Schema default value.
	if at.defaultValue != nil {
		return at.defaultValue, true, nil
	}
	return nil, false, nil
}

// baseValue routes to the correct arrayTransformer.baseValue format.
func (at *arrayTransformer) baseValue(in interface{}, path string, modifier pathModifier) ([]interface{}, bool, error) {
	if at.format == jsonInput {
		return at.baseValueJSON(in, path, modifier)
	}
	if at.format == xmlInput {
		return at.baseValueXML(in, path, modifier)
	}
	return nil, false, errors.New("unknown transform type in arrayTransformer baseValue")
}

func (at *arrayTransformer) child() instanceTransformer                 { return at.childTransformer }
func (at *arrayTransformer) path() string                               { return at.jsonPath }
func (at *arrayTransformer) selectChild(key string) instanceTransformer { return nil }

// arrayTransformJSON retrieves the value for this object by building the value for the base object and then adding in any
// transforms for all defined child fields.
func (at *arrayTransformer) arrayTransformJSON(in interface{}, modifier pathModifier) (interface{}, error) {
	path := at.jsonPath
	if modifier != nil {
		path = modifier(path)
	}
	base, changed, err := at.baseValue(in, path, modifier)
	if err != nil {
		return nil, err
	}

	if changed {
		// save the array base to in as children will use the value from this for their transforms
		if path == "$" {
			in = base
		} else {
			inMap, ok := in.(map[string]interface{})
			if !ok {
				return nil, errors.New("input is neither a JSON array nor object")
			}
			if err := saveInTree(inMap, path, base); err != nil {
				return nil, fmt.Errorf("failed to save array transform to input data: %v", err)
			}
		}
	}

	if at.childTransformer == nil {
		return base, nil
	}

	oldPath := path + "[*]"
	newArray := make([]interface{}, 0, len(base))

	for i := range base {
		currentPath := path + fmt.Sprintf("[%d]", i)

		childValue, err := at.childTransformer.transform(in, pathReplace(oldPath, currentPath, modifier))
		if err != nil {
			return nil, err
		}
		if childValue != nil {
			newArray = append(newArray, childValue)
		}
	}

	if len(newArray) == 0 {
		return nil, nil
	}
	return newArray, nil
}

// arrayTransformXML retrieves the value for this object by building the value for the base object and then adding in any
// transforms for all defined child fields.
func (at *arrayTransformer) arrayTransformXML(in interface{}, modifier pathModifier) (interface{}, error) {
	path := at.jsonPath
	if modifier != nil {
		path = modifier(path)
	}
	base, _, err := at.baseValue(in, path, modifier)
	if err != nil {
		return nil, err
	}

	if at.childTransformer == nil {
		return base, nil
	}

	oldPath := path + "[*]"
	newArray := make([]interface{}, 0, len(base))

	for i := range base {
		currentPath := path + fmt.Sprintf("[%d]", i)
		childValue := base[i]
		if _, ok := childValue.(*xmlquery.Node); ok {
			childValue, err = at.childTransformer.transform(childValue, pathReplace(oldPath, currentPath, modifier))
			if err != nil {
				return nil, err
			}
		}
		if childValue != nil {
			newArray = append(newArray, childValue)
		}
	}

	if len(newArray) == 0 {
		return nil, nil
	}
	return newArray, nil
}

// transform routes to the correct array transform type.
func (at *arrayTransformer) transform(in interface{}, modifier pathModifier) (interface{}, error) {
	if at.format == jsonInput {
		return at.arrayTransformJSON(in, modifier)
	}
	if at.format == xmlInput {
		return at.arrayTransformXML(in, modifier)
	}
	return nil, fmt.Errorf("Unrecognized transform type %s in arraytransformer transform, must be 'JSON' or 'XML' ", at.format)
}

// objectTransformer represents a JSON instance of type object and associated transforms.
type objectTransformer struct {
	children     map[string]instanceTransformer
	defaultValue map[string]interface{}
	jsonPath     string
	format       inputFormat
	transforms   *transformInstructions
}

func newObjectTransformer(path, transformIdentifier string, raw json.RawMessage, format inputFormat) (*objectTransformer, error) {
	ot := &objectTransformer{
		children: make(map[string]instanceTransformer),
		jsonPath: path,
		format:   format,
	}

	var err error
	ot.transforms, err = extractTransformInstructions(raw, transformIdentifier, path, "object")
	if err != nil {
		return nil, err
	}

	rawDefault, err := schemaDefault(raw)
	if err != nil {
		return nil, err
	}
	if rawDefault != nil {
		var ok bool
		ot.defaultValue, ok = rawDefault.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid default for an object `%v`", rawDefault)
		}
	}

	return ot, nil
}

func (ot *objectTransformer) addChild(child instanceTransformer) error {
	pathSplits := strings.Split(child.path(), ".")
	name := pathSplits[len(pathSplits)-1]
	parentPath := strings.Join(pathSplits[:len(pathSplits)-1], ".")
	if parentPath != ot.jsonPath {
		return fmt.Errorf("path %q is not a child of %q", child.path(), ot.jsonPath)
	}

	ot.children[name] = child
	return nil
}

func (ot *objectTransformer) child() instanceTransformer                 { return nil }
func (ot *objectTransformer) path() string                               { return ot.jsonPath }
func (ot *objectTransformer) selectChild(key string) instanceTransformer { return ot.children[key] }

// objectTransformJSON retrieves the value for this object by building the value for the base object and then adding in any
// transforms for all defined child fields.
func (ot *objectTransformer) objectTransformJSON(in interface{}, modifier pathModifier) (interface{}, error) {
	path := ot.jsonPath
	if modifier != nil {
		path = modifier(path)
	}
	var newValue map[string]interface{}

	// For the object use a transform if it exists or the default or an empty map
	if ot.transforms != nil {
		rawValue, err := ot.transforms.transform(in, "object", modifier, ot.format)
		if err != nil {
			return nil, err
		}
		if rawValue != nil {
			var ok bool
			newValue, ok = rawValue.(map[string]interface{})
			if !ok {
				return nil, errors.New("transform returned non-object value")
			}
		}
	}
	if newValue == nil {
		if ot.defaultValue == nil {
			newValue = make(map[string]interface{})
		} else {
			newValue = ot.defaultValue
		}
	}

	// Add each child value to the paren
	for _, child := range ot.children {
		childValue, err := child.transform(in, modifier)
		if err != nil {
			return nil, err
		}

		savePath := strings.Replace(child.path(), ot.jsonPath, "$", 1)
		if err := saveInTree(newValue, savePath, childValue); err != nil {
			return nil, fmt.Errorf("path %q failed save: %v", child.path(), err)
		}
	}

	if len(newValue) == 0 {
		return nil, nil
	}

	return newValue, nil
}

// objectTransformXML retrieves the value for this object by building the value for the base object and then adding in any
// transforms for all defined child fields. If a transform is provided it transforms the children relative to the
// passed in node. If a transform is provided and not found the children of the object are skipped.
func (ot *objectTransformer) objectTransformXML(in interface{}, modifier pathModifier) (interface{}, error) {
	path := ot.jsonPath
	if modifier != nil {
		path = modifier(path)
	}

	// For the object use a transform if it exists, if the transform does not find a node it will return nil unless a
	// default value is specified in which case the default value will be returned
	if ot.transforms != nil {
		rawValue, err := ot.transforms.transform(in, "object", modifier, ot.format)
		if err != nil {
			return nil, err
		}

		if rawValue == nil {
			if ot.defaultValue == nil {
				return nil, nil
			} else {
				return ot.defaultValue, nil
			}
		} else if val, ok := rawValue.(string); ok {
			// If the XML node is returned as an empty string, then it likely indicates that the transformer encountered an empty XML tag, e.g. <tag /> or <tag></tag>.
			// While not particularly useful, it is also not an error. The end result is that the field won't show up in the output file.
			if val == "" {
				return nil, nil
			}
		}

		switch v := rawValue.(type) {
		case *xmlquery.Node:
			in = v
		case []*xmlquery.Node:
			if len(v) > 0 {
				in = v[0]
			}
		default:
			return nil, errors.New("non xml node returned from object transform")
		}
	}

	var newValue map[string]interface{}
	if ot.defaultValue == nil {
		newValue = make(map[string]interface{})
	} else {
		newValue = ot.defaultValue
	}

	// Add each child value to the parent if there is no object transform or if the object transform node is found
	for _, child := range ot.children {
		childValue, err := child.transform(in, modifier)
		if err != nil {
			return nil, err
		}

		savePath := strings.Replace(child.path(), ot.jsonPath, "$", 1)
		if err := saveInTree(newValue, savePath, childValue); err != nil {
			return nil, fmt.Errorf("path %q failed save: %v", child.path(), err)
		}
	}

	if len(newValue) == 0 {
		return nil, nil
	}

	return newValue, nil
}

// transform routes to the correct object transform type.
func (ot *objectTransformer) transform(in interface{}, modifier pathModifier) (interface{}, error) {
	if ot.format == jsonInput {
		return ot.objectTransformJSON(in, modifier)
	}
	if ot.format == xmlInput {
		return ot.objectTransformXML(in, modifier)
	}
	return nil, fmt.Errorf("Unrecognized transform type %s in objecttransformer transform, must be 'JSON' or 'XML' ", ot.format)
}

// scalarTransformer represents a JSON instance for a scalar type.
type scalarTransformer struct {
	defaultValue interface{}
	jsonType     string
	jsonPath     string
	format       inputFormat
	transforms   *transformInstructions
}

func newScalarTransformer(path, transformIdentifier string, raw json.RawMessage, instanceType string, format inputFormat) (*scalarTransformer, error) {
	st := &scalarTransformer{
		jsonType: instanceType,
		jsonPath: path,
		format:   format,
	}

	if instanceType == "string" {
		instanceFormat, err := jsonparser.GetString(raw, "format")
		if err != nil && err != jsonparser.KeyPathNotFoundError {
			return nil, fmt.Errorf("failed to extract instance format: %v", err)
		}
		if instanceFormat == "date-time" {
			st.jsonType = "date-time"
		}
	}

	var err error
	st.transforms, err = extractTransformInstructions(raw, transformIdentifier, path, "scalar")
	if err != nil {
		return nil, err
	}

	st.defaultValue, err = schemaDefault(raw)
	if err != nil {
		return nil, err
	}

	return st, nil
}

func (st *scalarTransformer) addChild(instanceTransformer) error     { return nil }
func (st *scalarTransformer) child() instanceTransformer             { return nil }
func (st *scalarTransformer) path() string                           { return st.jsonPath }
func (st *scalarTransformer) selectChild(string) instanceTransformer { return nil }

// transformScalarJSON retrieves the value for a scalar instance following this process:
//
// 1. Use a Transform if it exists.
//
// 2. Look for the same jsonPath in the input and use directly if possible.
//
// 3. Fall back to the JSON Schema default value.
func (st *scalarTransformer) transformScalarJSON(in interface{}, modifier pathModifier) (interface{}, error) {
	path := st.jsonPath
	if modifier != nil {
		path = modifier(path)
	}
	// 1. Use a transform if it exists
	if st.transforms != nil {
		newValue, err := st.transforms.transform(in, st.jsonType, modifier, st.format)
		if err != nil {
			return nil, err
		}
		if newValue != nil {
			return newValue, nil
		}
	}

	// 2. Look for the same jsonPath in the input and use directly if possible.
	rawValue, err := jsonpath.Get(path, in)
	if err == nil {
		newValue, err := convert(rawValue, st.jsonType)
		// if there is a conversion error fall through to the default
		if newValue != nil {
			return newValue, err
		}
	}

	// 3. Fall back to the JSON Schema default value.
	return st.defaultValue, nil
}

// transformScalarXML retrieves the value for a scalar instance following this process:
//
// 1. Use a Transform if it exists.
//
// 2. If transform does not exist or returns no value send back default.
func (st *scalarTransformer) transformScalarXML(in interface{}, modifier pathModifier) (interface{}, error) {
	path := st.jsonPath
	if modifier != nil {
		path = modifier(path)
	}

	// 1. Use a Transform if it exists.
	if st.transforms != nil {
		newValue, err := st.transforms.transform(in, st.jsonType, modifier, st.format)
		if err != nil {
			return nil, err
		}
		if newValue != nil {
			return newValue, nil
		}
	}

	// 2. If transform does not exist or returns no value send back default
	return st.defaultValue, nil
}

// transform routes to the correct scalar transform type.
func (st *scalarTransformer) transform(in interface{}, modifier pathModifier) (interface{}, error) {
	if st.format == jsonInput {
		return st.transformScalarJSON(in, modifier)
	}
	if st.format == xmlInput {
		return st.transformScalarXML(in, modifier)
	}
	return nil, fmt.Errorf("Unrecognized transform type %s in scalartransformer transform, must be 'JSON' or 'XML' ", st.format)
}
