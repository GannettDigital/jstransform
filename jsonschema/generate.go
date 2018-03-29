package jsonschema

import "errors"

// GenerateStructs takes a JSON Schema and generates Golang structs that match the schema.
// The structs include struct tags for generating JSON.
// The JSON schema can specify more information than the structs enforce (like field size) and so validation of
// any JSON generated from the structs should still be done.
func GenerateStructs(schemaPath string) error {
	return errors.New("Not Implemented")
}
