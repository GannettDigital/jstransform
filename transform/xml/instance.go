package xmlTransform

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"github.com/antchfx/xmlquery"
	"github.com/buger/jsonparser"
)

// pathModifier is used to modify the JSON path of an instance to indicate
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

// arrayTransformer represents a JSON instance of type array.
type arrayTransformer struct {
	childTransformer instanceTransformer
	defaultValue     []interface{}
	jsonPath         string
	transforms       *transformInstructions
}

func newArrayTransformer(path, transformIdentifier string, raw json.RawMessage) (*arrayTransformer, error) {
	at := &arrayTransformer{
		jsonPath: path,
	}

	var err error
	at.transforms, err = extractTransformInstructions(raw, transformIdentifier, path)
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

func (at *arrayTransformer) baseValue(in interface{}, path string, modifier pathModifier) ([]interface{}, bool, error) {
	// 1. Use a transform if it exists
	if at.transforms != nil {
		rawValue, err := at.transforms.transform(in, "array", modifier)
		if err != nil {
			return nil, false, err
		}
		castedRaw := rawValue.([]*xmlquery.Node)
		var newValue []interface{}
		for _, item := range castedRaw {
			newValue = append(newValue, item)
		}
		return newValue, true, nil

		//if rawValue != nil {
		//	newValue, ok := rawValue.([]interface{})
		//	if !ok {
		//		newValue = []interface{}{rawValue}
		//	}
		//	return newValue, true, nil
		//}
	}

	//// 2. Look for the same JSONPath in the input and use directly if possible.
	//rawValue, err := jsonpath.Get(path, in)
	//if err == nil && rawValue != nil {
	//	newValue, ok := rawValue.([]interface{})
	//	if !ok {
	//		newValue = []interface{}{rawValue}
	//	}
	//	return newValue, false, nil
	//}

	//// 3. Fall back to the JSON Schema default value.
	//if at.defaultValue != nil {
	//	return at.defaultValue, true, nil
	//}
	return nil, false, nil
}

func (at *arrayTransformer) child() instanceTransformer                 { return at.childTransformer }
func (at *arrayTransformer) path() string                               { return at.jsonPath }
func (at *arrayTransformer) selectChild(key string) instanceTransformer { return nil }

// transform retrieves the value for this object by building the value for the base object and then adding in any
// transforms for all defined child fields.
func (at *arrayTransformer) transform(in interface{}, modifier pathModifier) (interface{}, error) {
	path := at.jsonPath
	if modifier != nil {
		path = modifier(path)
	}
	base, _, err := at.baseValue(in, path, modifier)
	if err != nil {
		return nil, err
	}

	//if changed {
	//	// save the array base to in as children will use the value from this for their transforms
	//	if path == "$" {
	//		in = base
	//	} else {
	//		inMap, ok := in.(map[string]interface{})
	//		if !ok {
	//			return nil, errors.New("input is neither a JSON array nor object")
	//		}
	//		if err := saveInTree(inMap, path, base); err != nil {
	//			return nil, fmt.Errorf("failed to save array transform to input data: %v", err)
	//		}
	//	}
	//}

	if at.childTransformer == nil {
		return base, nil
	}

	oldPath := path + "[*]"
	newArray := make([]interface{}, 0, len(base))

	for i := range base {
		currentPath := path + fmt.Sprintf("[%d]", i)
		in = base[i]
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

// objectTransformer represents a JSON instance of type object and associated transforms.
type objectTransformer struct {
	children     map[string]instanceTransformer
	defaultValue map[string]interface{}
	jsonPath     string
	transforms   *transformInstructions
}

func newObjectTransformer(path, transformIdentifier string, raw json.RawMessage) (*objectTransformer, error) {
	ot := &objectTransformer{
		children: make(map[string]instanceTransformer),
		jsonPath: path,
	}

	var err error
	ot.transforms, err = extractTransformInstructions(raw, transformIdentifier, path)
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

// transform retrieves the value for this object by building the value for the base object and then adding in any
// transforms for all defined child fields.
func (ot *objectTransformer) transform(in interface{}, modifier pathModifier) (interface{}, error) {
	path := ot.jsonPath
	if modifier != nil {
		path = modifier(path)
	}
	var newValue map[string]interface{}

	// For the object use a transform if it exists or the default or an empty map
	if ot.transforms != nil {
		rawValue, err := ot.transforms.transform(in, "object", modifier)
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

// scalarTransformer represents a JSON instance for a scalar type.
type scalarTransformer struct {
	defaultValue interface{}
	jsonType     string
	jsonPath     string
	transforms   *transformInstructions
}

func newScalarTransformer(path, transformIdentifier string, raw json.RawMessage, instanceType string) (*scalarTransformer, error) {
	st := &scalarTransformer{
		jsonType: instanceType,
		jsonPath: path,
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
	st.transforms, err = extractTransformInstructions(raw, transformIdentifier, path)
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

// transform retrieves the value for a scalar instance following this process:
//
// 1. Use a Transform if it exists.
//
// 2. Look for the same JSONPath in the input and use directly if possible.
//
// 3. Fall back to the JSON Schema default value.
func (st *scalarTransformer) transform(in interface{}, modifier pathModifier) (interface{}, error) {
	path := st.jsonPath
	if modifier != nil {
		path = modifier(path)
	}
	// 1. Use a transform if it exists
	if st.transforms != nil {
		newValue, err := st.transforms.transform(in, st.jsonType, modifier)
		if err != nil {
			return nil, err
		}
		if newValue != nil {
			return newValue, nil
		}
	}

	// 2. Look for the same JSONPath in the input and use directly if possible.
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
