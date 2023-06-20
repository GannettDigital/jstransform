package generate

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/GannettDigital/jstransform/jsonschema"

	"golang.org/x/exp/slices"
)

var (
	typeJSONGraphQL = map[string]string{
		"boolean":   "Boolean",
		"number":    "Float",
		"integer":   "Int",
		"string":    "String",
		"date-time": "DateTime",
		"object":    "(",
	}
)

// Developer Note
// This file started as a copy of `struct.go` and inherits some of its behavior
// that isn't necessary for GraphQL schema, but hasn't been removed or cleaned
// up yet.

// ----- Types -----

// goGQL represents the contents of a single go file to be generated based on the given JSON schema.
type goGQL struct {
	packageName   string
	args          BuildArgs
	rootStruct    *generatedGraphQLObject
	nestedStructs map[string]*generatedGraphQLObject // The key is derived from the path used by the walk function for the given struct
}

// generatedGraphQLObject represents a (not necessarily top-level) GraphQL type.
type generatedGraphQLObject struct {
	gqlExtractedField

	buildType string
}

// gqlExtractedField represents a Golang struct field as extracted from a JSON schema file. It is an intermediate format
// that is populated while parsing the JSON schema file then used when generating the Golang code for the struct.
type gqlExtractedField struct {
	arguments      []string
	array          bool
	arrayNullable  bool
	description    string
	fieldOrder     []string
	fields         gqlExtractedFields
	jsonName       string
	jsonType       string
	name           string
	nullable       bool
	implements     string
	target         string
	requiredFields map[string]bool
	args           BuildArgs
}

// gqlExtractedFields is a map of fields keyed on the field name.
type gqlExtractedFields map[string]*gqlExtractedField

// -----

// buildGraphQLFile generates the specified struct file.
// TODO: Most of this logic should go into 'newGeneratedGraphQLFile()' so that this function only opens the file and writes the parsed GraphQL object to it.
func buildGraphQLFile(schemaPath, name, packageName string, args BuildArgs) error {
	outPath := filepath.Join(args.OutputDirGraphQL, name+".graphqls")
	gfile, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("failed to open file %q: %w", outPath, err)
	}
	if _, err := fmt.Fprintf(gfile, "# %s\n", disclaimer); err != nil {
		return fmt.Errorf("failed writing GraphQL: %w", err)
	}

	if !filepath.IsAbs(schemaPath) {
		schemaPath = filepath.Join(filepath.Dir(args.SchemaPath), schemaPath)
	}
	schema, err := jsonschema.SchemaFromFile(schemaPath, name)
	if err != nil {
		return err
	}

	// If the schema has one or more AllOf definitions then this has an Interface type.
	var common []*goGQL
	var commonFields int
	if len(schema.AllOf) != 0 {
		common = make([]*goGQL, 0, len(schema.AllOf))
		for _, allOfSchema := range schema.AllOf {
			baseName := strings.Split(filepath.Base(schema.AllOf[0].FromRef), ".")[0]
			generated, err := newGeneratedGraphQLFile(allOfSchema, baseName, packageName, false, args)
			if err != nil {
				return fmt.Errorf("failed to build generated struct: %w", err)
			}
			common = append(common, generated)
		}
		// TODO: Doesn't handle multiple AllOf properly with regard to mapping to Go type name.  Ex: `BaseLinks`.
		// TODO: Having to merge all the sub-parsed schema files into a single one synthesized here is regretable.
		obj := goGQL{
			packageName: packageName,
			args:        args,
			rootStruct: &generatedGraphQLObject{
				buildType: "type",
				gqlExtractedField: gqlExtractedField{
					args:           args,
					array:          false,
					description:    "",
					fieldOrder:     []string{},
					fields:         map[string]*gqlExtractedField{},
					jsonName:       packageName,
					jsonType:       "object",
					name:           name,
					requiredFields: map[string]bool{},
				},
			},
			nestedStructs: map[string]*generatedGraphQLObject{},
		}
		if len(common) != 0 {
			obj.rootStruct.gqlExtractedField.jsonName = exportedName(packageName)
			name := strings.Split(filepath.Base(schema.AllOf[0].FromRef), ".")[0]
			if newName, ok := args.StructNameMap[name]; ok {
				name = newName
			} else {
				name = exportedName(name)
			}
			obj.rootStruct.gqlExtractedField.name = name
			for _, com := range common {
				commonFields += len(com.rootStruct.fields)
				for fk, fv := range com.rootStruct.fields {
					obj.rootStruct.fields[fk] = fv
					obj.rootStruct.fieldOrder = append(obj.rootStruct.fieldOrder, fk)
				}
				for rk, rv := range com.rootStruct.requiredFields {
					obj.rootStruct.requiredFields[rk] = rv
				}
				for nk, nv := range com.nestedStructs {
					obj.nestedStructs[nk] = nv
				}
			}
			if len(schema.OneOf) != 0 {
				if newName, ok := args.GraphQLTypeNameMap[packageName]; ok {
					obj.rootStruct.gqlExtractedField.jsonName = newName
				} else {
					obj.rootStruct.gqlExtractedField.jsonName = exportedName(packageName)
				}
				obj.rootStruct.gqlExtractedField.name = obj.rootStruct.gqlExtractedField.jsonName
				obj.rootStruct.buildType = "interface"
			}
		}
		if err := obj.write(gfile); err != nil {
			return fmt.Errorf("failed to write allOf struct: %w", err)
		}
	}

	// If the schema has one or more OneOf definitions then these are the Implementations of the Interface type.
	for _, oneOfSchema := range schema.OneOf {
		oneOfName := strings.Split(filepath.Base(oneOfSchema.FromRef), ".")[0]
		generated, err := newGeneratedGraphQLFile(oneOfSchema, oneOfName, packageName, len(schema.AllOf) != 0, args)
		if err != nil {
			return fmt.Errorf("failed to build generated struct: %w", err)
		}
		if len(schema.AllOf) != 0 {
			// Only an implementation of an interface if there was also an AllOf schema.
			var implements string
			if newName, ok := args.GraphQLTypeNameMap[name]; ok {
				implements = newName
			} else {
				implements = exportedName(name)
			}
			generated.rootStruct.implements = implements
		}

		// Must merge into this implementation the fields defined in the AllOf blocks
		// because in GraphQL each implementation repeats the fields from the interface
		// definition.
		lfd := generated.rootStruct.fields
		lfo := generated.rootStruct.fieldOrder
		generated.rootStruct.fields = make(map[string]*gqlExtractedField, len(generated.rootStruct.fields)+commonFields)
		generated.rootStruct.fieldOrder = make([]string, 0, len(generated.rootStruct.fields)+commonFields)
		for _, com := range common {
			for fk, fv := range com.rootStruct.fields {
				generated.rootStruct.fields[fk] = fv
				generated.rootStruct.fieldOrder = append(generated.rootStruct.fieldOrder, fk)
			}
			for rk, rv := range com.rootStruct.requiredFields {
				generated.rootStruct.requiredFields[rk] = rv
			}
		}
		for _, lf := range lfo {
			generated.rootStruct.fields[lf] = lfd[lf]
			generated.rootStruct.fieldOrder = append(generated.rootStruct.fieldOrder, lf)
		}

		gfile := gfile
		outPath := outPath
		if args.InterfaceFiles {
			outName := oneOfName + ".graphqls"
			if !strings.HasPrefix(oneOfName, name) {
				outName = name + "_" + outName
			}
			outPath = filepath.Join(args.OutputDirGraphQL, outName)
			gfile, err = os.Create(outPath)
			if err != nil {
				return fmt.Errorf("failed to open file %q: %w", outPath, err)
			}
			if _, err := fmt.Fprintf(gfile, "# %s\n", disclaimer); err != nil {
				return fmt.Errorf("failed writing GraphQL: %w", err)
			}
		}
		if err := generated.write(gfile); err != nil {
			return fmt.Errorf("failed to write oneOf struct: %w", err)
		}
		if args.InterfaceFiles {
			if err := gfile.Close(); err != nil {
				return fmt.Errorf("failed to close file %q: %w", outPath, err)
			}
		}
	}

	// The properties in this schema.  As in, not AllOf or OneOf subschema definitions.
	if len(schema.AllOf) == 0 && len(schema.OneOf) == 0 {
		schemaName := strings.Split(filepath.Base(schemaPath), ".")[0]
		generated, err := newGeneratedGraphQLFile(schema.Instance, schemaName, packageName, false, args)
		if err != nil {
			return fmt.Errorf("failed to build generated struct: %w", err)
		}
		if err := generated.write(gfile); err != nil {
			return fmt.Errorf("failed to write properties struct: %w", err)
		}
	}

	if err := gfile.Close(); err != nil {
		return fmt.Errorf("failed to close file %q: %w", outPath, err)
	}
	return nil
}

// write outputs the Golang representation of this field to the writer with prefix before each line.
// It handles inline structs by calling this method recursively adding a new \t to the prefix for each layer.
// If required is set to false 'omitempty' is added in the JSON struct tag for the field.
func (ef *gqlExtractedField) write(w io.Writer, prefix string, required, descriptionAsStructTag, pointers bool) error {
	var attribute string
	description := ef.description
	if strings.HasPrefix(description, "DEPRECATED:") {
		// First line is the deprecation reason and the rest is the normal description.
		lines := strings.Split(description, "\n")
		attribute = fmt.Sprintf(" @deprecated(reason: %q)", strings.TrimSpace(strings.TrimPrefix(lines[0], "DEPRECATED:")))
		description = strings.TrimSpace(strings.Join(lines[1:], "\n"))
	}
	if description != "" {
		// Multi-line descriptions get the """ comment syntax.
		if strings.IndexRune(description, '\n') < 0 {
			description = fmt.Sprintf("%s\"%s\"\n", prefix, description)
		} else {
			newDescription := prefix + `"""` + "\n"
			for _, line := range strings.Split(description, "\n") {
				if line == "" {
					newDescription += "\n"
				} else {
					newDescription += prefix + line + "\n"
				}
			}
			newDescription += prefix + `"""` + "\n"
			description = newDescription
		}
		if _, err := fmt.Fprint(w, description); err != nil {
			return fmt.Errorf("error writing field %q description: %w", ef.name, err)
		}
	}

	// Simple field type.
	if ef.jsonType != "object" {
		gqlArgs, gqlType := ef.graphqlType(required, pointers)
		if _, err := fmt.Fprintf(w, "%s%s%s: %s%s\n", prefix, ef.jsonName, gqlArgs, gqlType, attribute); err != nil {
			return fmt.Errorf("error writing field %q definition: %w", ef.name, err)
		}
		return nil
	}

	// Object field type.  (Not generally used for GraphQL as the types aren't nested.)
	gqlArgs, gqlType := ef.graphqlType(required, pointers)
	if _, err := fmt.Fprintf(w, "%s%s%s\t%s {\n", prefix, ef.jsonName, gqlArgs, gqlType); err != nil {
		return err
	}

	sort.Strings(ef.fieldOrder)
	for _, fieldName := range ef.fieldOrder {
		field := ef.fields[fieldName]
		if err := field.write(w, prefix+"\t", ef.requiredFields[field.jsonName], descriptionAsStructTag, pointers); err != nil {
			return fmt.Errorf("failed writing field %q: %w", field.name, err)
		}
	}

	if _, err := fmt.Fprintf(w, "%s\t}%s", prefix, attribute); err != nil {
		return err
	}
	return nil
}

// Sorted will return the fields in a sorted list. The sort is a string sort on the keys.
func (efs gqlExtractedFields) Sorted() []*gqlExtractedField {
	var sortedKeys sort.StringSlice
	fieldsByName := make(map[string]*gqlExtractedField, len(efs))
	for _, f := range efs {
		sortedKeys = append(sortedKeys, f.name)
		fieldsByName[f.name] = f
	}

	sortedKeys.Sort()

	sorted := make([]*gqlExtractedField, 0, len(sortedKeys))
	for _, key := range sortedKeys {
		sorted = append(sorted, fieldsByName[key])
	}

	return sorted
}

// newGeneratedGraphQLFile creates a GraphQL schema file based on the given JSON schema.
// The write function can be used to write out the value of the file, which will end up with either a single struct
// or multiple depending on the presence of nested structs and the value of the NoNestedStructs build argument.
func newGeneratedGraphQLFile(schema jsonschema.Instance, name, packageName string, implements bool, args BuildArgs) (*goGQL, error) {
	required := make(map[string]bool, len(schema.Required))
	for _, fname := range schema.Required {
		required[fname] = true
	}

	gof := &goGQL{
		packageName:   packageName,
		args:          args,
		nestedStructs: make(map[string]*generatedGraphQLObject),
	}

	// Apply renaming for Go structure and GraphQL type.
	gqlName := name
	if newName, ok := args.StructNameMap[name]; ok {
		name = newName
	} else {
		name = exportedName(name)
	}
	if newName, ok := args.GraphQLTypeNameMap[gqlName]; ok {
		gqlName = newName
	} else if implements {
		gqlName = strings.ToLower(gqlName)
	} else {
		gqlName = exportedName(name)
	}

	gof.rootStruct = gof.newGeneratedGraphQLObject(name, gqlName, required)
	if len(schema.OneOf) != 0 {
		gof.rootStruct.buildType = "interface"
	} else {
		gof.rootStruct.buildType = "type"
	}

	if err := jsonschema.Walk(&jsonschema.Schema{Instance: schema}, gof.walkFunc); err != nil {
		return nil, fmt.Errorf("failed to walk schema for %q: %w", name, err)
	}

	return gof, nil
}

func (gof *goGQL) newGeneratedGraphQLObject(name, gqlName string, requiredFields map[string]bool) *generatedGraphQLObject {
	return &generatedGraphQLObject{
		buildType: "type",
		gqlExtractedField: gqlExtractedField{
			jsonName:       gqlName,
			name:           name,
			fields:         make(map[string]*gqlExtractedField),
			requiredFields: requiredFields,
			args:           gof.args,
		},
	}
}

func (gof *goGQL) structs() []*generatedGraphQLObject {
	if len(gof.nestedStructs) == 0 {
		return []*generatedGraphQLObject{gof.rootStruct}
	}
	nested := make([]*generatedGraphQLObject, len(gof.nestedStructs))
	var i int
	for _, s := range gof.nestedStructs {
		nested[i] = s
		i++
	}

	// order with root first and nested in a consistent following order
	sort.Slice(nested, func(i, j int) bool {
		return nested[i].name < nested[j].name
	})
	structs := []*generatedGraphQLObject{gof.rootStruct}
	structs = append(structs, nested...)

	return structs
}

// walkFunc is a jsonschema.WalkFunc which builds the fields within the gofile as the
// JSON schema file is walked.
func (gof *goGQL) walkFunc(path string, i jsonschema.Instance) error {
	gen := gof.rootStruct

	if !gof.args.NoNestedStructs {
		return gen.addField(splitJSONPath(path), nil, i)
	}

	parts := []string{exportedName(gof.rootStruct.name)}
	for _, part := range splitJSONPath(path) {
		parts = append(parts, exportedName(part))
	}

	// Find parent struct or if none use root struct
	parentKey := strings.Join(parts[:len(parts)-1], "")
	name := splitJSONPath(path)[len(parts)-2]

	// There's a convention that "implements" types are all lowercase,
	// but sub-objects use the normal Pascal-case format like other
	// GraphQL objects do.  This uses the exported name as the base
	// level GraphQL name prefix if this is the first level off the
	// root schema tree.
	var gqlTypeName string
	gen = gof.nestedStructs[parentKey]
	if gen == nil {
		gen = gof.rootStruct
		gqlTypeName = exportedName(gen.jsonName)
	} else {
		gqlTypeName = gen.jsonName
	}
	gqlTypeName += exportedName(name)

	// If the types is an object create a new generated struct for it
	if slices.Contains(i.Type, "object") {
		key := strings.Join(parts, "")
		structType := key

		/*
		   // nullable nested structs could be pointers
		   if !gen.requiredFields[name] && gof.args.Pointers || slices.Contains(i.Type, "null") {
		           structType = "*" + key
		   }
		*/

		requiredFields := make(map[string]bool, len(i.Required))
		for _, name := range i.Required {
			requiredFields[name] = true
		}

		obj := gof.newGeneratedGraphQLObject(key, gqlTypeName, requiredFields)
		obj.arguments = i.GraphQLArguments
		obj.target = i.Target
		gof.nestedStructs[key] = obj

		// Don't add GraphQL schema for sub-objects of a targeted object.
		if gen.target != "" || gen.buildType == "ignored" {
			obj.buildType = "ignored"
		}

		// Generate proper GraphQL type name target, taking into account renaming.
		gqlTypeList := []string{gqlTypeName}
		if slices.Contains(i.Type, "null") {
			gqlTypeList = append(gqlTypeList, "null")
		}

		return gen.addField([]string{name}, gqlTypeList, jsonschema.Instance{Description: i.Description, Type: []string{structType}, Target: i.Target, GraphQLArguments: i.GraphQLArguments})
	}

	return gen.addField([]string{name}, nil, i)
}

// write will write the generated file to the given io.Writer.
func (gof *goGQL) write(w io.Writer) error {
	buf := &bytes.Buffer{} // the formatter uses the entire output, so buffer for that

	for _, s := range gof.structs() {
		if s.target == "" && s.buildType != "ignored" {
			if _, err := buf.Write([]byte("\n")); err != nil {
				return fmt.Errorf("failed writing GraphQL %q: %w", s.name, err)
			}
			if err := s.write(buf); err != nil {
				return fmt.Errorf("failed writing GraphQL %q: %w", s.name, err)
			}
		}
	}

	if _, err := w.Write(buf.Bytes()); err != nil {
		return fmt.Errorf("error writing to io.Writer: %w", err)
	}
	return nil
}

// write will write the generated file to the given io.Writer.
func (gen *generatedGraphQLObject) write(w io.Writer) error {
	// If there is an alternate type for this object then don't generate it.
	if gen.target != "" || gen.buildType == "ignored" {
		return nil
	}
	var implements string
	if gen.implements != "" {
		implements = " implements " + gen.implements
	}
	// Some structures are "virtual" in that there is a GraphQL schema with only hydration and no underlying storage.
	var hasModel string
	for _, field := range gen.fields {
		if field.jsonType != "graphql-hydration" {
			hasModel = fmt.Sprintf("@goModel(model: %q) ", strings.Join([]string{gen.args.ImportPath, exportedName(gen.name)}, "."))
			break
		}
	}
	if _, err := fmt.Fprintf(w, "%s %s%s %s{\n", gen.buildType, gen.jsonName, implements, hasModel); err != nil {
		return fmt.Errorf("failed writing GraphQL: %w", err)
	}

	sortedFields := gen.fields.Sorted()
	for idx, field := range sortedFields {
		if err := field.write(w, "  ", gen.requiredFields[field.jsonName], gen.args.DescriptionAsStructTag, gen.args.Pointers); err != nil {
			return fmt.Errorf("failed writing field %q: %w", field.name, err)
		}
		if idx+1 != len(sortedFields) {
			if _, err := io.WriteString(w, "\n"); err != nil {
				return fmt.Errorf("failed writing field %q extra line: %w", field.name, err)
			}
		}
	}

	if _, err := w.Write([]byte("}\n")); err != nil {
		return fmt.Errorf("failed writing GraphQL: %w", err)
	}

	return nil
}

// addField will create a new field or add to an existing field in the gqlExtractedFields.
// Nested fields are handled by recursively calling this function until the leaf field is reached.
// For all fields the name and jsonType are set, for arrays the array bool is set for true and for JSON objects,
// the fields map is created and if it exists the requiredFields section populated.
// fields will be renamed if a matching entry is supplied in the fieldRenameMap.
func (gen *gqlExtractedField) addField(tree []string, gqlTypeName []string, inst jsonschema.Instance) error {
	if len(tree) > 1 {
		if f, ok := gen.fields[tree[0]]; ok {
			return f.addField(tree[1:], nil, inst)
		}
		f := &gqlExtractedField{
			jsonName:  tree[0],
			jsonType:  "object",
			name:      exportedName(tree[0]),
			fields:    make(map[string]*gqlExtractedField),
			args:      gen.args,
			target:    inst.Target,
			arguments: inst.GraphQLArguments,
		}
		gen.fields[tree[0]] = f
		gen.fieldOrder = append(gen.fieldOrder, tree[0])
		if err := f.addField(tree[1:], nil, inst); err != nil {
			return fmt.Errorf("failed field %q: %w", tree[0], err)
		}
		return nil
	}

	if len(tree) > 0 {
		fieldName, ok := gen.args.FieldNameMap[tree[0]]
		if !ok {
			fieldName = tree[0]
		}
		typeSource := inst.Type
		if gqlTypeName != nil {
			typeSource = gqlTypeName
		}
		var jsonType string
		for _, iType := range typeSource {
			if iType == "null" {
				continue
			}
			if jsonType != "" {
				return fmt.Errorf("cannot generate GraphQL schema for union types of %q and %q", jsonType, iType)
			}
			jsonType = iType
		}
		f := &gqlExtractedField{
			description: inst.Description,
			name:        exportedName(fieldName),
			jsonName:    tree[0],
			jsonType:    jsonType,
			args:        gen.args,
			target:      inst.Target,
			arguments:   inst.GraphQLArguments,
			nullable:    slices.Contains(inst.Type, "null"),
		}
		// Second processing of an array type
		if exists, ok := gen.fields[f.jsonName]; ok {
			f = exists
			if f.array && f.jsonType == "" {
				f.jsonType = jsonType
			} else {
				return fmt.Errorf("field %q already exists but is not an array field", f.name)
			}
		}
		if slices.Contains(inst.Type, "string") && inst.Format == "date-time" {
			f.jsonType = "date-time"
			f.nullable = true
		}

		switch f.jsonType {
		case "array":
			f.jsonType = ""
			f.nullable = false
			f.array = true
			f.arrayNullable = f.nullable

			totalDescription := fmt.Sprintf("The total length of the %s list at this same level in the data, this number is unaffected by filtering.", fieldName)
			// If the field this is a total of is deprecated then also use its deprecation message (first line of description).
			if strings.HasPrefix(f.description, "DEPRECATED:") {
				lines := strings.Split(f.description, "\n")
				totalDescription = lines[0] + "\n" + totalDescription
			}
			totalName := "total" + f.name
			gen.fields[totalName] = &gqlExtractedField{
				description: totalDescription,
				name:        exportedName(totalName),
				jsonName:    totalName,
				jsonType:    "integer",
				args:        gen.args,
				target:      "",
				arguments:   nil,
			}
			gen.fieldOrder = append(gen.fieldOrder, totalName)
			gen.requiredFields[totalName] = true
		case "object":
			f.requiredFields = make(map[string]bool, len(inst.Required))
			for _, name := range inst.Required {
				f.requiredFields[name] = true
			}
			f.fields = make(map[string]*gqlExtractedField)
		}

		gen.fields[tree[0]] = f
		gen.fieldOrder = append(gen.fieldOrder, tree[0])
	}

	return nil
}

// graphqlType maps a jsonType to a string representation of the GraphQL type.
// If Array is true it makes the type into an array.
// If the JSON Schema had a type of "string" and a format of "date-time" it is expected the input jsonType will be
// "date-time".
// Non-required times are added as pointers to allow for their values to missing go marshalled JSON.
func (ef *gqlExtractedField) graphqlType(required, pointers bool) (string, string) {
	// This is a "virtual" field in that no storage exists for it, but the GraphQL schema
	// still needs for it to be defined so a resolver can be attached to do the lookup.
	if ef.jsonType == "graphql-hydration" {
		if len(ef.arguments) == 0 {
			return "", ef.target
		}
		arguments := make([]string, 0, len(ef.arguments)+2)
		arguments = append(arguments, "(")
		for _, argStr := range ef.arguments {
			if argStr != "" {
				argStr = "    " + argStr
			}
			arguments = append(arguments, argStr)
		}
		arguments = append(arguments, "  )")
		return strings.Join(arguments, "\n"), ef.target
	}

	var regularType bool
	var graphqlType string
	if ef.target != "" {
		graphqlType = ef.target
	} else if gqlType, isaJSONType := typeJSONGraphQL[ef.jsonType]; isaJSONType {
		graphqlType = gqlType
		regularType = true
	} else {
		graphqlType = ef.jsonType
	}
	if graphqlType != "DateTime" && (regularType && !ef.nullable || required || !pointers) {
		graphqlType += "!"
	}

	var graphqlArguments string
	if ef.array {
		graphqlArguments = `(
    "A List Filter expression such as '{Field: \"position\", Operation: \"<=\", Argument: {Value: 10}}'"
    filter: ListFilter

    "Sort the list, ie '{Field: \"position\", Order: \"ASC\"}'"
    sort: ListSortParams
  )`
		if ef.arrayNullable || pointers && !required {
			graphqlType = "[" + graphqlType + "]"
		} else {
			graphqlType = "[" + graphqlType + "]!"
		}
		graphqlType += " @goField(forceResolver: true)"
	}
	return graphqlArguments, graphqlType
}
