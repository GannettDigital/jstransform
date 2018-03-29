package jsonschema

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

// Concat will combine any two arbitrary values, though only strings are supported for non-trivial concatenation.
func concat(a, b interface{}) (interface{}, error) {
	switch {
	case a == nil && b == nil:
		return nil, nil
	case a == nil:
		return b, nil
	case b == nil:
		return a, nil
	}

	atype := reflect.TypeOf(a).String()
	btype := reflect.TypeOf(b).String()
	if atype != btype {
		return nil, fmt.Errorf("can't concat types %q and %q", atype, btype)
	}

	switch a.(type) {
	case string:
		return a.(string) + b.(string), nil
	default:
		return nil, fmt.Errorf("concatenation of types %q not supported", atype)
	}
}

// convert takes the raw value and checks to see if it matches the jsonType, if not it will attempt to convert it
// to the correct type. The function does not set defaults so a nil value will be returned as nil not as the desired
// types empty type.
func convert(raw interface{}, jsonType string) (interface{}, error) {
	if raw == nil {
		return nil, nil
	}
	if rawArray, ok := raw.([]interface{}); ok && len(rawArray) == 1 {
		raw = rawArray[0]
	} else if ok && len(rawArray) == 0 {
		return nil, nil
	}
	switch jsonType {
	case "boolean":
		return convertBoolean(raw)
	case "number":
		return convertNumber(raw)
	case "string":
		return convertString(raw)
	}
	return raw, nil
}

func convertBoolean(raw interface{}) (interface{}, error) {
	switch t := raw.(type) {
	case bool:
		return raw, nil
	case string:
		if t == "" {
			return false, nil
		}
		return strconv.ParseBool(t)
	case int:
		return t > 0, nil
	case float32:
		return t > 0, nil
	case float64:
		return t > 0, nil
	default:
		return nil, fmt.Errorf("unable to convert type %q to boolean", reflect.TypeOf(raw))
	}
}

func convertNumber(raw interface{}) (interface{}, error) {
	switch t := raw.(type) {
	case bool:
		if t {
			return 1, nil
		}
		return 0, nil
	case string:
		if t == "" {
			return nil, nil
		}
		if value, err := strconv.Atoi(t); err == nil {
			return value, nil
		}
		if value, err := strconv.ParseFloat(t, 64); err == nil {
			return value, nil
		}
		return nil, fmt.Errorf("failed to convert string %q to number", t)
	case int, float32, float64:
		return raw, nil
	default:
		return nil, fmt.Errorf("unable to convert type %q to a number", reflect.TypeOf(raw))
	}
}

func convertString(raw interface{}) (interface{}, error) {
	switch t := raw.(type) {
	case bool:
		return strconv.FormatBool(t), nil
	case string:
		return raw, nil
	case int:
		return strconv.Itoa(t), nil
	case uint64:
		return strconv.FormatUint(t, 10), nil
	case float32:
		t64 := float64(t)
		return strconv.FormatFloat(t64, 'f', -1, 32), nil
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64), nil
	default:
		return nil, fmt.Errorf("unable to convert type %q to a string", reflect.TypeOf(raw))
	}
}

// schemaDefault determines the default for an instance based on the JSONSchema.
// If no default is defined for arrays and objects it will return an empty value, for all other types nil.
func schemaDefault(schema json.RawMessage) (interface{}, error) {
	ifields := struct {
		Type    string      `json:"type"`
		Default interface{} `json:"default"`
	}{}
	if err := json.Unmarshal(schema, &ifields); err != nil {
		return nil, fmt.Errorf("failed to extract schema default: %v", err)
	}

	// Try the default
	if ifields.Default != nil {
		return ifields.Default, nil
	}

	switch ifields.Type {
	case "object":
		return map[string]interface{}{}, nil
	case "array":
		return []interface{}{}, nil
	default:
		return nil, nil
	}
}
