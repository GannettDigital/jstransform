package jsonschema

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/franela/goreq"
)

// jsonRef represents a JSON Reference source and targetRef
type jsonRef struct {
	Source []string
	Target string
}

// Dereference parse JSON string and replaces all $ref with the referenced data.
func Dereference(schemaPath string, input []byte) ([]byte, error) {
	if !strings.Contains(string(input), "$ref") {
		return input, nil
	}

	var data interface{}
	json.Unmarshal([]byte(input), &data)
	refs, err := walkInterface(data, []string{}, []jsonRef{})
	if err != nil {
		return input, fmt.Errorf("unable to walk interface %s: %v", schemaPath, err)
	}

	for _, ref := range refs {
		top := data
		for i, item := range ref.Source {
			if i < len(ref.Source)-1 {
				// assuming integer item is slice[index] instead of map[string]
				if intKey, err := strconv.Atoi(item); err == nil {
					top = top.([]interface{})[intKey]
				} else {
					top = top.(map[string]interface{})[item]
				}
			} else {
				targetRef, err := buildReference(schemaPath, data, ref.Target)
				if err != nil {
					return input, fmt.Errorf("unable to build reference from %s: %v", ref.Target, err)
				}
				targetKeys := reflect.ValueOf(targetRef).MapKeys()
				if len(targetKeys) > 1 {
					// assuming integer item is slice[index] instead of map[string]
					if intKey, err := strconv.Atoi(item); err == nil {
						top.([]interface{})[intKey] = targetRef
					} else {
						top.(map[string]interface{})[item] = targetRef
					}
				} else {
					// when targetRef = single KV pair, set the value using the key instead of overwriting entire map
					key := targetKeys[0].Interface().(string)
					// assuming integer item is slice[index] instead of map[string]
					if intKey, err := strconv.Atoi(item); err == nil {
						top.([]interface{})[intKey].(map[string]interface{})[key] = targetRef.(map[string]interface{})[key]
						delete(top.([]interface{})[intKey].(map[string]interface{}), "$ref")
					} else {
						top.(map[string]interface{})[item].(map[string]interface{})[key] = targetRef.(map[string]interface{})[key]
						delete(top.(map[string]interface{})[item].(map[string]interface{}), "$ref")
					}
				}
			}
		}
	}

	return json.Marshal(data)
}

// walkInterface traverses the map[string]interface{} to located json references
func walkInterface(node interface{}, source []string, refs []jsonRef) ([]jsonRef, error) {
	var err error
	for key, val := range node.(map[string]interface{}) {
		switch reflect.TypeOf(val).Kind() {
		case reflect.String:
			if key == "$ref" {
				refs = append(refs, jsonRef{
					Source: source,
					Target: val.(string),
				})
			}
		case reflect.Slice:
			for i, item := range val.([]interface{}) {
				if reflect.TypeOf(item).Kind() == reflect.Map {
					refs, err = walkInterface(item, append(source, key, strconv.Itoa(i)), refs)
					if err != nil {
						return nil, fmt.Errorf("unable to walk slice interface: %v", err)
					}
				}
			}
		case reflect.Map:
			refs, err = walkInterface(node.(map[string]interface{})[key], append(source, key), refs)
			if err != nil {
				return nil, fmt.Errorf("unable to walk map interface: %v", err)
			}
		}
	}
	return refs, nil
}

// HttpReferenceClient isolates the HTTP call for testing purposes
func HttpReferenceClient(url string) (*goreq.Response, error) {
	return goreq.Request{Uri: url}.Do()
}

// buildReference constructs the json reference: internal, file or http
func buildReference(schemaPath string, top interface{}, ref string) (interface{}, error) {
	target := strings.Split(ref, "#")
	if len(target) < 2 {
		target = append(target, "/")
	}
	var source interface{}

	switch {
	case len(target[0]) == 0:
		source = top
	case strings.HasPrefix(target[0], "http"):
		res, err := HttpReferenceClient(target[0])
		if err != nil {
			return nil, fmt.Errorf("unable to get reference from %s: %v", target[0], err)
		}
		res.Body.FromJsonTo(&source)
	default:
		refPath, err := filepath.Abs(path.Dir(schemaPath) + "/" + target[0])
		if err != nil {
			return nil, fmt.Errorf("unable to expand reference filepath %s: %v", target[0], err)
		}
		data, err := ioutil.ReadFile(refPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read reference file %q: %v", refPath, err)
		}
		data, err = Dereference(refPath, data)
		if err != nil {
			return nil, fmt.Errorf("failed to dereference refPath %s: %v", refPath, err)
		}
		json.Unmarshal([]byte(data), &source)
	}
	return parseReference(source, strings.Split(target[1], "/")[1:]), nil
}

// parseReference recursively parses the given reference path
func parseReference(source interface{}, refPaths []string) interface{} {
	if len(refPaths) > 1 {
		return parseReference(source.(map[string]interface{})[refPaths[0]], refPaths[1:])
	} else {
		if refPaths[0] != "" {
			return source.(map[string]interface{})[refPaths[0]]
		} else {
			return source
		}
	}
}
