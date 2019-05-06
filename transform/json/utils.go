package json

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/buger/jsonparser"
)

var indexRe = regexp.MustCompile(`\[([\d]+)\]`)

func extractTransformInstructions(raw json.RawMessage, transformIdentifier, path string) (*transformInstructions, error) {
	rawTransformInstruction, _, _, err := jsonparser.Get(raw, "transform", transformIdentifier)
	if err != nil && err != jsonparser.KeyPathNotFoundError {
		return nil, fmt.Errorf("failed to extract raw instance transform: %v", err)
	} else if len(rawTransformInstruction) == 0 {
		return nil, nil
	}

	splits := strings.Split(path, ".")
	parentPath := strings.Join(splits[:len(splits)-1], ".")

	var tis transformInstructions
	if err := json.Unmarshal(rawTransformInstruction, &tis); err != nil {
		return nil, fmt.Errorf("failed to unmarshal instance transform: %v", err)
	}
	// replaces the @. format
	tis.replaceJSONPathPrefix("@.", parentPath+".")
	// replaces the @[] format
	tis.replaceJSONPathPrefix("@[", parentPath+"[")

	return &tis, nil
}

// schemaDefault determines the default for an instance based on the JSONSchema.
// If no default is defined nil is returned.
func schemaDefault(schema json.RawMessage) (interface{}, error) {
	ifields := struct {
		Default interface{} `json:"default"`
	}{}
	if err := json.Unmarshal(schema, &ifields); err != nil {
		return nil, fmt.Errorf("failed to extract schema default: %v", err)
	}

	// Try the default
	if ifields.Default != nil {
		return ifields.Default, nil
	}

	return nil, nil
}

// replaceIndex takes a path which may include array index values like `a[0].b.c[23].d` with the index values replaced
// with "*", ie `a[*].b.c[*].d`
func replaceIndex(path string) string {
	return indexRe.ReplaceAllString(path, "[*]")
}
