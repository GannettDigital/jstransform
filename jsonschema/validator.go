package jsonschema

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"runtime"
	"sort"

	"github.com/xeipuuv/gojsonschema"
)

// Validator defines an interface for validating JSON matches a JSON schema.
type Validator interface {
	Validate(raw json.RawMessage) (bool, error)
}

type validator struct {
	schema *gojsonschema.Schema
}

// NewValidator returns a Validator for the schema at the given file path.
// The Validate method on the Validator allows for verifying a JSON file matches the JSON schema.
func NewValidator(schemaPath string) (*validator, error) {
	schema, err := loadSchema(schemaPath)
	if err != nil {
		return nil, err
	}
	return &validator{
		schema: schema,
	}, nil
}

func loadSchema(schemasDir string) (*gojsonschema.Schema, error) {
	path, err := filepath.Abs(schemasDir)
	if err != nil {
		return nil, err
	}
	filePrefix := "file://"
	if runtime.GOOS == "windows" {
		// Convert all path separators (specifically Windows' \ ) to /
		path = filepath.ToSlash(path)
		// http://blogs.msdn.com/b/ie/archive/2006/12/06/file-uris-in-windows.aspx
		filePrefix += "/"
	}
	l := gojsonschema.NewReferenceLoader(filePrefix + path)
	schema, err := gojsonschema.NewSchema(l)
	if err != nil {
		return nil, err
	}
	return schema, nil
}

// Validate will check that the given json is validate according the schema loaded by the Validator.
func (v *validator) Validate(raw json.RawMessage) (bool, error) {
	l := gojsonschema.NewBytesLoader(raw)

	result, err := v.schema.Validate(l)
	if err != nil {
		return false, err
	}
	if len(result.Errors()) > 0 {
		sort.Slice(result.Errors(), func(i, j int) bool {
			return result.Errors()[i].String() < result.Errors()[j].String()
		})
		return false, fmt.Errorf("invalid schema: %v", result.Errors())
	}

	return result.Valid(), nil
}
