package jsonschema

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/GannettDigital/jsonparser"
)

const (
	fromRef      = "fromRef"
	refKey       = "$ref"
	transformKey = "transform"
)

// dereference parses JSON and replaces all $ref with the referenced data.
// If $ref refers to a file schemaFromFile is called and in this way references in referenced files are handled
// recursively along with other processing done by schemaFromFile.
func dereference(schemaPath string, data json.RawMessage, oneOfType string) (json.RawMessage, error) {
	refs, err := findRefs(data)
	if err != nil {
		return nil, fmt.Errorf("failed when finding refs: %v", err)
	}

	for _, refPath := range refs {
		ref, err := jsonparser.GetString(data, refPath...)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve ref at path %v: %v", refPath, err)
		}

		resolved, err := resolveRef(ref, data, schemaPath, oneOfType)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve ref %q at path %v: %v", ref, refPath, err)
		}

		destPath := refPath[:len(refPath)-1]
		// It is necessary to delete the refKey reference so they are not refound
		data = jsonparser.Delete(data, append(destPath, refKey)...)

		if len(destPath) != 0 {
			// Attempt to read the `transform` object from the source data so that we can apply it after the ref resolving
			// wipes out that object in the data. Ensures that transform is an object and errors if not.
			transformPath := append(destPath, transformKey)
			transform, dataType, _, err := jsonparser.Get(data, transformPath...)
			if err != nil && err != jsonparser.KeyPathNotFoundError {
				return nil, fmt.Errorf("failed to read transform object on source data at path %v: %v", transformPath, err)
			}
			if transform != nil && dataType != jsonparser.Object {
				return nil, fmt.Errorf("transform object is wrong type %q, should be object", dataType)
			}

			// Set the resolved ref contents on the data. This wipes out existing fields in that object
			data, err = jsonparser.Set(data, resolved, destPath...)
			if err != nil {
				return nil, fmt.Errorf("failed to update data with resolved ref %q at path %v: %v", ref, refPath, err)
			}

			// If we found a transform inside that object in the source data, apply that back since setting of the ref
			// would have cleared it
			if transform != nil {
				data, err = jsonparser.Set(data, transform, transformPath...)
				if err != nil {
					return nil, fmt.Errorf("failed to update data transform with resolved ref %q at path %v: %v", ref, refPath, err)
				}
			}
		}
	}

	// TODO sometimes refs remain because of refs with refs in them and the order not being right. Getting ordering
	// right is not easy as the info needed to order is not available until late in the process. Still something
	// better than reprocessing would be nice
	remaining, err := findRefs(data)
	if err != nil {
		return nil, fmt.Errorf("failed checking for remaining refs: %v", err)
	}
	if len(remaining) > 0 {
		return dereference(schemaPath, data, oneOfType)
	}

	return data, nil
}

// findRefs searches through the given JSON finding the location in the structure of all refKeys.
func findRefs(data json.RawMessage) ([][]string, error) {
	refs := make([][]string, 0)

	err := jsonparser.ObjectEach(data, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		sKey := string(key)
		switch dataType {
		case jsonparser.String:
			if sKey == refKey {
				refs = append(refs, []string{sKey})
			}
		case jsonparser.Object:
			childRefs, err := findRefs(value)
			if err != nil {
				return err
			}
			for _, r := range childRefs {
				combined := append([]string{sKey}, r...)
				refs = append(refs, combined)
			}
		case jsonparser.Array:
			var index int
			_, err := jsonparser.ArrayEach(value, func(avalue []byte, adataType jsonparser.ValueType, aoffset int, aerr error) {
				currentIndex := fmt.Sprintf("[%d]", index)
				index++
				if adataType != jsonparser.Object {
					return
				}
				cRefs, err := findRefs(avalue)
				if err != nil {
					return
				}
				for _, r := range cRefs {
					combined := append([]string{sKey, currentIndex}, r...)
					refs = append(refs, combined)
				}
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return refs, nil
}

// resolveRef looks at the reference value passed in as ref and resolves it to a set of JSON.
// The reference may refer to a definition within the given data or a file reference.
// For files schemaPath is used to resolve relative references then SchemaFromFile is used to build the file.
// oneOfType is used by schemaFromFile to select a specific oneOfType.
func resolveRef(ref string, data json.RawMessage, schemaPath string, oneOfType string) (json.RawMessage, error) {
	// TODO there is nothing here to stop circular references other than self references
	var sourcePath, target string
	splits := strings.SplitN(ref, "#", 2)
	if len(splits) != 2 {
		if strings.Contains(ref, "#") {
			target = strings.Trim(ref, "#/")
		} else {
			sourcePath = ref
		}
	} else {
		sourcePath = splits[0]
		target = strings.Trim(splits[1], "/")
	}

	var source json.RawMessage

	switch {
	case sourcePath == "":
		source = data
	case strings.HasPrefix(sourcePath, "http"):
		resp, err := http.Get(sourcePath)
		if err != nil {
			return nil, fmt.Errorf("unable to get reference from %q: %v", sourcePath, err)
		}
		source, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read body from %q: %v", sourcePath, err)
		}
		// TODO since SchemaFromFile does the current allOf/oneOf processing this data does not go through that processing
	default: // Default to assuming it is a file reference
		// TODO this could be rather inefficient if there are multiple references to the same sourcePath but a different
		// target as it will currently reprocess the source file everytime
		refPath, err := filepath.Abs(filepath.Join(filepath.Dir(schemaPath), sourcePath))
		if err != nil {
			return nil, fmt.Errorf("unable to expand reference filepath %q: %v", sourcePath, err)
		}
		if schemaPath == refPath {
			source = data
			break
		}
		schema, err := SchemaFromFile(refPath, oneOfType)
		if err != nil {
			return nil, fmt.Errorf("failed to process reference file %q: %v", refPath, err)
		}
		source, err = json.Marshal(schema)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal schema from file %q: %v", refPath, err)
		}
	}

	var err error
	if target == "" {
		data = source
	} else {
		data, _, _, err = jsonparser.Get(source, strings.Split(target, "/")...)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve ref %q: %v", target, err)
		}
	}

	data = []byte(strings.TrimSpace(string(data)))
	data, err = jsonparser.Set(data, []byte(fmt.Sprintf("%q", ref)), fromRef)
	if err != nil {
		return nil, fmt.Errorf("failed to set fromRef for reference %q: %v", ref, err)
	}

	return data, nil
}
