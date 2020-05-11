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
func (ef *extractedField) write(w io.Writer, prefix string, required, descriptionAsStructTag bool) error {
	var omitempty string
	if !required {
		omitempty = ",omitempty"
	}
	jsonTag := fmt.Sprintf(`json:"%s%s"`, ef.jsonName, omitempty)
	var description string
	if descriptionAsStructTag && ef.description != "" {
		description = fmt.Sprintf(`description:"%s"`, strings.Split(ef.description, "\n")[0])
	}
	structTag := fmt.Sprintf("`%s`\n", strings.Trim(strings.Join([]string{jsonTag, description}, " "), " "))

	if !descriptionAsStructTag && ef.description != "" {
		for _, line := range strings.Split(ef.description, "\n") {
			if _, err := w.Write([]byte(fmt.Sprintf("// %s\n", line))); err != nil {
				return err
			}
		}
	}

	if ef.jsonType != "object" {
		_, err := w.Write([]byte(fmt.Sprintf("%s%s\t%s\t%s", prefix, ef.name, goType(ef.jsonType, ef.array), structTag)))
		return err
	}

	if _, err := w.Write([]byte(fmt.Sprintf("%s%s\t%s {\n", prefix, ef.name, goType(ef.jsonType, ef.array)))); err != nil {
		return err
	}

	for _, field := range ef.fields.Sorted() {
		fieldRequired := ef.requiredFields[field.jsonName]
		if err := field.write(w, prefix+"\t", fieldRequired, descriptionAsStructTag); err != nil {
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

// goFile represents the contents of a single go file to be generated based on the given JSON schema.
type goFile struct {
	packageName   string
	args          BuildArgs
	rootStruct    *generatedStruct
	nestedStructs map[string]*generatedStruct // The key is derived from the path used by the walk function for the given struct
}

// newGeneratedGoFile creates a go file based on the given JSON schema.
// The write function can be used to write out the value of the file, which will end up with either a single struct
// or multiple depending on the presence of nested structs and the value of the NoNestedStructs build argument.
func newGeneratedGoFile(schema *jsonschema.Schema, name, packageName string, embeds []string, args BuildArgs) (*goFile, error) {
	required := map[string]bool{}
	for _, fname := range schema.Required {
		required[fname] = true
	}

	gof := &goFile{
		packageName:   packageName,
		args:          args,
		nestedStructs: make(map[string]*generatedStruct),
	}

	gof.rootStruct = gof.newGeneratedStruct(name, required)
	gof.rootStruct.embededStructs = embeds

	if err := jsonschema.Walk(schema, gof.walkFunc); err != nil {
		return nil, fmt.Errorf("failed to walk schema for %q: %v", name, err)
	}

	return gof, nil
}

func (gof *goFile) newGeneratedStruct(name string, requiredFields map[string]bool) *generatedStruct {
	return &generatedStruct{
		extractedField: extractedField{
			name:           name,
			fields:         make(map[string]*extractedField),
			requiredFields: requiredFields,
		},
		args: gof.args,
	}
}

func (gof *goFile) structs() []*generatedStruct {
	if len(gof.nestedStructs) < 1 {
		return []*generatedStruct{gof.rootStruct}
	}
	nested := make([]*generatedStruct, len(gof.nestedStructs))
	var i int
	for _, s := range gof.nestedStructs {
		nested[i] = s
		i++
	}

	// order with root first and nested in a consistent following order
	sort.Slice(nested, func(i, j int) bool {
		return nested[i].name < nested[j].name
	})
	structs := []*generatedStruct{gof.rootStruct}
	structs = append(structs, nested...)

	return structs
}

// walkFunc is a jsonschema.WalkFunc which builds the fields for generatedStructFile within the gofile as the
// JSON schema file is walked.
func (gof *goFile) walkFunc(path string, i jsonschema.Instance) error {
	gen := gof.rootStruct

	if !gof.args.NoNestedStructs {
		return addField(gen.fields, splitJSONPath(path), i, gen.args.FieldNameMap)
	}

	parts := []string{exportedName(gof.rootStruct.name)}
	for _, part := range splitJSONPath(path) {
		parts = append(parts, exportedName(part))
	}

	// Find parent struct or if none use root struct
	parentKey := strings.Join(parts[:len(parts)-1], "")
	name := splitJSONPath(path)[len(parts)-2]
	gen = gof.nestedStructs[parentKey]
	if gen == nil {
		gen = gof.rootStruct
	}

	// If the types is an object create a new generated struct for it
	if i.Type == "object" {
		key := strings.Join(parts, "")

		requiredFields := make(map[string]bool)
		for _, name := range i.Required {
			requiredFields[name] = true
		}

		gof.nestedStructs[key] = gof.newGeneratedStruct(key, requiredFields)
		return addField(gen.fields, []string{name}, jsonschema.Instance{Description: i.Description, Type: key}, gen.args.FieldNameMap)
	}

	return addField(gen.fields, []string{name}, i, gen.args.FieldNameMap)
}

// write will write the generated file to the given io.Writer.
func (gof *goFile) write(w io.Writer) error {
	buf := &bytes.Buffer{} // the formatter uses the entire output, so buffer for that

	if _, err := buf.Write([]byte(fmt.Sprintf("package %s\n\n%s\n\n", gof.packageName, disclaimer))); err != nil {
		return fmt.Errorf("failed writing struct: %v", err)
	}

	var includeTime bool
	for _, s := range gof.structs() {
		if s.fields.IncludeTime() {
			includeTime = true
			break
		}
	}

	if includeTime {
		if _, err := buf.Write([]byte("import \"time\"\n")); err != nil {
			return fmt.Errorf("failed writing imports: %v", err)
		}
	}

	for _, s := range gof.structs() {
		if _, err := buf.Write([]byte("\n\n")); err != nil {
			return fmt.Errorf("failed writing struct %q: %v", s.name, err)
		}
		if err := s.write(buf); err != nil {
			return fmt.Errorf("failed writing struct %q: %v", s.name, err)
		}
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

type generatedStruct struct {
	extractedField

	args           BuildArgs
	embededStructs []string
}

// write will write the generated file to the given io.Writer.
func (gen *generatedStruct) write(w io.Writer) error {
	embeds := strings.Join(gen.embededStructs, "\n")
	if embeds != "" {
		embeds += "\n\n"
	}
	if _, err := w.Write([]byte(fmt.Sprintf("type %s struct {\n%s", exportedName(gen.name), embeds))); err != nil {
		return fmt.Errorf("failed writing struct: %v", err)
	}

	for _, field := range gen.fields.Sorted() {
		req := gen.requiredFields[field.jsonName]
		if err := field.write(w, "\t", req, gen.args.DescriptionAsStructTag); err != nil {
			return fmt.Errorf("failed writing field %q: %v", field.name, err)
		}
	}

	if _, err := w.Write([]byte("}")); err != nil {
		return fmt.Errorf("failed writing struct: %v", err)
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
