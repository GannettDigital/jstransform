package jsonschema

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/GannettDigital/jsonparser"
)

// Instance represents a JSON Schema instance.
type Instance struct {
	AdditionalProperties bool                       `json:"additionalProperties,omitempty"`
	AllOf                []Instance                 `json:"allOf,omitempty"`
	AnyOf                []Instance                 `json:"anyOf,omitempty"` // TODO unsupported
	Description          string                     `json:"description,omitempty"`
	Definitions          json.RawMessage            `json:"definitions,omitempty"`
	Format               string                     `json:"format,omitempty"`
	FromRef              string                     `json:"fromRef,omitempty"`           // Added as a way of tracking the ref which was already expanded
	GoModel              string                     `json:"goModel,omitempty"`           // For specifying the Go model which to reference.
	GraphQLArguments     []string                   `json:"graphql-arguments,omitempty"` // For type="graphql-hydration" to also require query arguments.
	Items                json.RawMessage            `json:"items,omitempty"`
	OneOf                []Instance                 `json:"oneOf,omitempty"`
	Properties           map[string]json.RawMessage `json:"properties,omitempty"`
	Ref                  string                     `json:"$ref,omitempty"`
	Required             []string                   `json:"required,omitempty"`
	Schema               string                     `json:"$schema,omitempty"`
	Target               string                     `json:"target,omitempty"` // For type="graphql-hydration" or GraphQL schema type overrides.
	Type                 []string                   `json:"type"`
}
type instanceUnmarshal Instance

func (i *Instance) UnmarshalJSON(data []byte) error {
	// The "type" field in JSON Schema is either a scalar or an array.
	// This code normalizes it to an array so that unmarshaling succeeds.
	value, dataType, _, err := jsonparser.Get(data, "type")
	if err != nil && !errors.Is(err, jsonparser.KeyPathNotFoundError) {
		return fmt.Errorf("error reading schema type: %w %s", err, data)
	} else if dataType == jsonparser.String {
		data, err = jsonparser.Set(data, []byte(fmt.Sprintf("[%q]", value)), "type")
		if err != nil {
			return fmt.Errorf("error changing schema type to array: %w", err)
		}
	}
	return json.Unmarshal(data, (*instanceUnmarshal)(i))
}

// Schema represents a JSON Schema with the AllOf and OneOf references parsed and squashed into a single representation.
// This is not a fully spec compatible representation but a basic representation useful for walking through the schema
// instances within a schema. Also note AnyOf fields are not supported at this time.
//
// A fully spec compatible version of the schema is kept for validation purposes.
type Schema struct {
	Instance

	Validator Validator `json:"-"`
}

// Validate will check that the given json is validate according the schema.
func (s *Schema) Validate(raw json.RawMessage) (bool, error) {
	return s.Validator.Validate(raw)
}

// SchemaFromFile parses a file at the given path and returns a schema based on its contents.
// The function traverses allOf fields within the schema. For oneOf fields the reference base
// name minus any extension is compared to the value of the oneOfType argument and if they match that file is also
// traversed. AnyOf fields are currently ignored.
//
// Referenced files are recursively processed. At this time only definition and file references are supported.
func SchemaFromFile(schemaPath string, oneOfType string) (*Schema, error) {
	data, err := os.ReadFile(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema file %q: %v", schemaPath, err)
	}
	v, err := NewValidator(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize schema validator: %v", err)
	}

	// dereferencing during walking is more efficient but more complicated so all dereferencing for a file is done immediately
	data, err = dereference(schemaPath, data, oneOfType)
	if err != nil {
		return nil, fmt.Errorf("failed to Dereference Schema: %v", err)
	}

	// json schema's default behavior is additionalProperties: true if the field is missing so mimic that behavior here
	var sj = Instance{
		AdditionalProperties: true,
	}
	if err := json.Unmarshal(data, &sj); err != nil {
		return nil, fmt.Errorf("failed to Unmarshal Schema: %v", err)
	}

	s := Schema{
		Instance:  sj,
		Validator: v,
	}

	// TODO this behavior is not spec compatible, according to the spec it is possible to have multiple allOf instances
	// that conflict. The legit use case for that is rare but it is in the spec. Rather than merge these the walk
	// should go through each set of files but this makes the raw walking much more complicated.
	for _, all := range sj.AllOf {
		s.Properties = mergeProperties(s.Properties, all.Properties)
		s.Required = append(s.Required, all.Required...)
		if s.Items == nil {
			s.Items = all.Items
		}
	}

	for _, one := range sj.OneOf {
		subName := strings.Split(filepath.Base(one.FromRef), ".")[0]
		if subName != oneOfType {
			continue
		}

		s.Properties = mergeProperties(s.Properties, one.Properties)
		s.Required = append(s.Required, one.Required...)
		if s.Items == nil {
			s.Items = one.Items
		}
	}

	return &s, nil
}

// SchemaTypes will parse the given file and report which top level allOfTypes, oneOfTypes, and properties are found in the schema.
func SchemaTypes(schemaPath string) ([]string, []string, []string, error) {
	data, err := os.ReadFile(schemaPath)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to read schema file %q: %v", schemaPath, err)
	}

	var sj Instance
	if err := json.Unmarshal(data, &sj); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to Unmarshal Schema: %v", err)
	}

	var allOfTypes []string
	for _, a := range sj.AllOf {
		allOfTypes = append(allOfTypes, a.Ref)
	}
	var oneOfTypes []string
	for _, one := range sj.OneOf {
		oneOfTypes = append(oneOfTypes, one.Ref)
	}
	var properties []string
	for prop := range sj.Properties {
		properties = append(properties, prop)
	}
	sort.Strings(properties)

	return allOfTypes, oneOfTypes, properties, nil
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

// FieldType returns the type name of a field and whether the field is nullable.
func FieldType(data []byte) (string, bool, error) {
	// The "type" field in JSON Schema is either a scalar or an array.
	value, dataType, _, err := jsonparser.Get(data, "type")
	if err != nil && !errors.Is(err, jsonparser.KeyPathNotFoundError) {
		return "", false, fmt.Errorf("error reading schema type: %w", err)
	} else if dataType == jsonparser.String {
		return string(value), false, nil
	} else if dataType == jsonparser.Array {
		var nullable bool
		var jsonType string
		_, aErr := jsonparser.ArrayEach(value, func(avalue []byte, adataType jsonparser.ValueType, aoffset int, aerr error) {
			if adataType != jsonparser.String {
				return
			}
			switch string(avalue) {
			case "null":
				nullable = true
			default:
				if jsonType != "" {
					err = fmt.Errorf("union types are not supported: %s", value)
				}
				jsonType = string(avalue)
			}
		})
		if aErr != nil && err == nil {
			err = aErr
		}
		if err != nil {
			return "", false, fmt.Errorf("error iterating over type array: %w", err)
		}
		return jsonType, nullable, nil
	} else {
		return "", false, fmt.Errorf("unknown schema type: %s", dataType)
	}
}
