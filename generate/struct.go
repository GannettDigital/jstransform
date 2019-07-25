package generate

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"sort"
	"strings"

	"github.com/GannettDigital/jstransform/jsonschema"
)

// extractedField represents a Golang struct field as extracted from a JSON schema file. It is an intermediate format
// that is populated while parsing the JSON schema file then used when generating the Golang code for the struct.
type extractedField struct {
	array          bool
	description    string
	fields         extractedFields
	jsonName       string
	jsonType       string
	name           string
	requiredFields map[string]bool
}

// write outputs the Golang representation of this field to the writer with prefix before each line.
// It handles inline structs by calling this method recursively adding a new \t to the prefix for each layer.
// If required is set to false 'omitempty' is added in the JSON struct tag for the field
func (ef *extractedField) write(w io.Writer, prefix string, required bool) error {
	var omitempty string
	if !required {
		omitempty = ",omitempty"
	}
	jsonTag := fmt.Sprintf(`json:"%s%s"`, ef.jsonName, omitempty)
	var description string
	if ef.description != "" {
		description = fmt.Sprintf(`description:"%s"`, strings.Split(ef.description, "\n")[0])
	}
	structTag := fmt.Sprintf("`%s`\n", strings.Trim(strings.Join([]string{jsonTag, description}, " "), " "))

	if ef.jsonType != "object" {
		_, err := w.Write([]byte(fmt.Sprintf("%s%s\t%s\t%s", prefix, ef.name, goType(ef.jsonType, ef.array), structTag)))
		return err
	}

	if _, err := w.Write([]byte(fmt.Sprintf("%s%s\t%s {\n", prefix, ef.name, goType(ef.jsonType, ef.array)))); err != nil {
		return err
	}

	for _, field := range ef.fields.Sorted() {
		fieldRequired := ef.requiredFields[field.jsonName]
		if err := field.write(w, prefix+"\t", fieldRequired); err != nil {
			return fmt.Errorf("failed writing field %q: %v", field.name, err)
		}
	}

	if _, err := w.Write([]byte(fmt.Sprintf("%s\t}\t%s", prefix, structTag))); err != nil {
		return err
	}
	return nil
}

// extractedFields is a map of fields keyed on the field name.
type extractedFields map[string]*extractedField

// IncludeTime does a depth-first recursive search to see if any field or child field is of type "date-time"
func (efs extractedFields) IncludeTime() bool {
	for _, field := range efs {
		if field.fields != nil {
			if field.fields.IncludeTime() {
				return true
			}
		}
		if field.jsonType == "date-time" {
			return true
		}
	}
	return false
}

// Sorted will return the fields in a sorted list. The sort is a string sort on the keys
func (efs extractedFields) Sorted() []*extractedField {
	var sorted []*extractedField
	var sortedKeys sort.StringSlice
	fieldsByName := make(map[string]*extractedField)
	for _, f := range efs {
		sortedKeys = append(sortedKeys, f.name)
		fieldsByName[f.name] = f
	}

	sortedKeys.Sort()

	for _, key := range sortedKeys {
		sorted = append(sorted, fieldsByName[key])
	}

	return sorted
}

type generatedStruct struct {
	extractedField

	packageName    string
	embededStructs []string
	fieldNameMap   map[string]string
}

func newGeneratedStruct(schema *jsonschema.Schema, name, packageName string, embeds []string, fieldNameMap map[string]string) (*generatedStruct, error) {
	required := map[string]bool{}
	for _, fname := range schema.Required {
		required[fname] = true
	}
	generated := &generatedStruct{
		extractedField: extractedField{
			name:           name,
			fields:         make(map[string]*extractedField),
			requiredFields: required,
		},
		fieldNameMap:   fieldNameMap,
		embededStructs: embeds,
		packageName:    packageName,
	}
	if err := jsonschema.Walk(schema, generated.walkFunc); err != nil {
		return nil, fmt.Errorf("failed to walk schema for %q: %v", name, err)
	}

	return generated, nil
}

// walkFunc is a jsonschema.WalkFunc which builds the fields in the generatedStructFile as the JSON schema file is
// walked.
func (gen *generatedStruct) walkFunc(path string, i jsonschema.Instance) error {
	if err := addField(gen.fields, splitJSONPath(path), i, gen.fieldNameMap); err != nil {
		return err
	}
	return nil
}

// write will write the generated file to the given io.Writer.
func (gen *generatedStruct) write(w io.Writer) error {
	buf := &bytes.Buffer{} // the formatter uses the entire output, so buffer for that

	if _, err := buf.Write([]byte(fmt.Sprintf("package %s\n\n%s\n\n", gen.packageName, disclaimer))); err != nil {
		return fmt.Errorf("failed writing struct: %v", err)
	}

	if gen.fields.IncludeTime() {
		if _, err := buf.Write([]byte("import \"time\"\n")); err != nil {
			return fmt.Errorf("failed writing struct: %v", err)
		}
	}

	embeds := strings.Join(gen.embededStructs, "\n")
	if embeds != "" {
		embeds += "\n"
	}
	if _, err := buf.Write([]byte(fmt.Sprintf("type %s struct {\n%s\n", exportedName(gen.name), embeds))); err != nil {
		return fmt.Errorf("failed writing struct: %v", err)
	}

	for _, field := range gen.fields.Sorted() {
		req := gen.requiredFields[field.jsonName]
		if err := field.write(buf, "\t", req); err != nil {
			return fmt.Errorf("failed writing field %q: %v", field.name, err)
		}
	}

	if _, err := buf.Write([]byte("}")); err != nil {
		return fmt.Errorf("failed writing struct: %v", err)
	}

	final, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("failed to format source: %v", err)
	}

	if _, err := w.Write(final); err != nil {
		return fmt.Errorf("error writing to io.Writer: %v", err)
	}
	return nil
}

// addField will create a new field or add to an existing field in the extractedFields.
// Nested fields are handled by recursively calling this function until the leaf field is reached.
// For all fields the name and jsonType are set, for arrays the array bool is set for true and for JSON objects,
// the fields map is created and if it exists the requiredFields section populated.
// fields will be renamed if a matching entry is supplied in the fieldRenameMap
func addField(fields extractedFields, tree []string, inst jsonschema.Instance, fieldRenameMap map[string]string) error {
	if len(tree) > 1 {
		if f, ok := fields[tree[0]]; ok {
			return addField(f.fields, tree[1:], inst, fieldRenameMap)
		}
		f := &extractedField{jsonName: tree[0], jsonType: "object", name: exportedName(tree[0]), fields: make(map[string]*extractedField)}
		fields[tree[0]] = f
		if err := addField(f.fields, tree[1:], inst, fieldRenameMap); err != nil {
			return fmt.Errorf("failed field %q: %v", tree[0], err)
		}
		return nil
	}

	if len(tree) > 0 {
		fieldName, ok := fieldRenameMap[tree[0]]
		if !ok {
			fieldName = tree[0]
		}
		f := &extractedField{
			description: inst.Description,
			name:        exportedName(fieldName),
			jsonName:    tree[0],
			jsonType:    inst.Type,
		}
		// Second processing of an array type
		if exists, ok := fields[f.jsonName]; ok {
			f = exists
			if f.array && f.jsonType == "" {
				f.jsonType = inst.Type
			} else {
				return fmt.Errorf("field %q already exists but is not an array field", f.name)
			}
		}
		if inst.Type == "string" && inst.Format == "date-time" {
			f.jsonType = "date-time"
		}

		switch f.jsonType {
		case "array":
			f.jsonType = ""
			f.array = true
		case "object":
			f.requiredFields = make(map[string]bool)
			for _, name := range inst.Required {
				f.requiredFields[name] = true
			}
			f.fields = make(map[string]*extractedField)
		}

		fields[tree[0]] = f
	}

	return nil
}
