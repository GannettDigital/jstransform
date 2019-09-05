package generate

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/actgardner/gogen-avro/generator"
	"github.com/actgardner/gogen-avro/schema"
	"golang.org/x/tools/go/ast/astutil"
)

const (
	metadataFields = `{"name":"AvroWriteTime","doc":"The timestamp when this avro data is written. Useful for identifying the newest row of data sharing keys.","type":"long","logicalType":"timestamp-millis"},{"name":"AvroDeleted","doc":"This is set to true when the Avro data is recording a delete in the source data.","default":false,"type":"boolean"},`
)

type avroConfig struct {
	dir           string
	excludeFields map[string]bool
	namespace     []string
	writer        io.Writer
}

// buildAvroSerializationFunctions will run gogen-avro on the given schema to build a serialization and deserialization
// code for the schema. It will create the directory `avro/<schema name>` and add the generated files there. The directory
// is necessary because generated code for different types will clobber each other.
func buildAvroSerializationFunctions(schemaPath string) error {
	schemaName := strings.Split(filepath.Base(schemaPath), ".")[0]
	pkg := generator.NewPackage(schemaName)
	namespace := schema.NewNamespace(false)

	schema, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema file at path %q: %v", schemaPath, err)
	}

	if _, err := namespace.TypeForSchema(schema); err != nil {
		return fmt.Errorf("failed to decode schema file: %v", err)
	}

	err = namespace.AddToPackage(pkg, "// Code generated by github.com/actgardner/gogen-avro. DO NOT EDIT.")
	if err != nil {
		return fmt.Errorf("failed generating serialization code: %v", err)
	}

	destination := filepath.Join(filepath.Dir(schemaPath), "avro", schemaName)
	if err := os.MkdirAll(destination, os.ModeDir|os.ModePerm); err != nil && !os.IsExist(err) {
		return fmt.Errorf("failed to create directory %q: %v", destination, err)
	}
	if err := pkg.WriteFiles(destination); err != nil {
		return fmt.Errorf("failed writing serialization source files to dir %q: %v", destination, err)
	}

	return nil
}

// buildAvroSchemaFile will generate an Avro schema based on the go struct with the given name in the file at the
// given path. I will write it to the same directory as the go file as `<name>.avsc` but with the name lowercased.
// http://avro.apache.org/docs/current/spec.html
//
// By default the file is written with now whitespace to minimize the size, choose pretty for better formatting
//
// Note: Avro can't handle maps with a key other than a string, http://avro.apache.org/docs/current/spec.html#Maps
// Neither can JSON schema, https://json-schema.org/understanding-json-schema/reference/object.html so this only
// becomes relevant if it is used with go structs which weren't just generated from JSON schema
func buildAvroSchemaFile(name, goSourcePath string, pretty bool) (string, error) {
	// Step 1 create the Avro Schema file
	// there are 3 ways to approach this, walk the JSON schema, walk the AST for the go struct or load the Go struct up and do reflection
	// reflection isn't generally clear but is probably the most compact, walking the JSON schema will likely require building a
	// representation of the data in a new set of structs like is done for the go struct building initially.
	// Walking the go AST is able to take advantage of all the standard library AST methods and though it has its share
	// or complication is easier to do in one pass
	spec, err := parseGoStruct(name, goSourcePath)
	if err != nil {
		return "", fmt.Errorf("failed to parse Go struct: %v", err)
	}
	if spec == nil {
		return "", errors.New("type spec is nil")
	}
	dir := filepath.Dir(goSourcePath)
	outPath := filepath.Join(dir, strings.ToLower(name)+".avsc")
	specFile, err := os.Create(outPath)
	if err != nil {
		return outPath, fmt.Errorf("failed to open file %q: %v", outPath, err)
	}

	var writer io.Writer
	var buf *bytes.Buffer
	writer = specFile
	if pretty {
		buf = &bytes.Buffer{}
		writer = buf
	}

	cfg := avroConfig{
		dir:    dir,
		writer: writer,
	}

	fmt.Fprint(cfg.writer, "{")
	astutil.Apply(spec, writeAvroStruct(cfg, name, metadataFields), nil)
	fmt.Fprint(cfg.writer, "]}")

	if pretty {
		rawJSON := json.RawMessage(buf.Bytes())
		enc := json.NewEncoder(specFile)
		enc.SetIndent("", "  ")
		if err := enc.Encode(rawJSON); err != nil {
			return "", fmt.Errorf("failed writing pretty output to file: %v", err)
		}
	}

	if err := specFile.Close(); err != nil {
		return "", fmt.Errorf("failed closing file: %v", err)
	}

	return outPath, nil
}

// parseGoStruct parses the go file(s) at path returning the named struct type definition as an *ast.TypeSpec
func parseGoStruct(name, path string) (*ast.TypeSpec, error) {
	fileSet := token.NewFileSet()

	var goFile *ast.File
	if strings.HasSuffix(path, ".go") {
		var err error
		goFile, err = parser.ParseFile(fileSet, path, nil, parser.AllErrors)
		if err != nil {
			return nil, fmt.Errorf("failed to parse go file %q: %v", path, err)
		}
		if !ast.FilterFile(goFile, func(itemName string) bool { return itemName == name }) {
			return nil, fmt.Errorf("a struct named %q was not found in file %q", name, path)
		}
	} else {
		var err error
		goFile, err = filterPackage(name, path)
		if err != nil {
			return nil, err
		}
	}

	d, ok := goFile.Decls[0].(*ast.GenDecl)
	if !ok {
		return nil, errors.New("filtered declaration is of unknown ast type")
	}
	if len(d.Specs) != 1 {
		return nil, errors.New("unexpected number of specs in declaration")
	}
	t, ok := d.Specs[0].(*ast.TypeSpec)
	if !ok {
		return nil, errors.New("filtered type is of unknown ast type")
	}

	if t == nil {
		return nil, errors.New("go struct not found in parsed data")
	}
	return t, nil
}

// filterPackage will filter all the go files found at the path returning a *ast.File containing name.
func filterPackage(name, path string) (*ast.File, error) {
	fileSet := token.NewFileSet()
	pkgmap, err := parser.ParseDir(fileSet, path, nil, parser.AllErrors)
	if err != nil {
		return nil, fmt.Errorf("failed to parse go files at path %q: %v", path, err)
	}
	var pkgs []*ast.Package
	for _, pkg := range pkgmap {
		pkgs = append(pkgs, pkg)
	}
	if length := len(pkgs); length != 1 {
		return nil, fmt.Errorf("expected 1 package for go files at path %q, found %d", path, length)
	}
	pkg := pkgs[0]
	if !ast.FilterPackage(pkg, func(itemName string) bool { return itemName == name }) {
		return nil, fmt.Errorf("a struct named %q was not found in file %q", name, path)
	}
	var goFile *ast.File
	for _, f := range pkg.Files {
		if len(f.Decls) == 1 {
			if goFile != nil {
				return nil, fmt.Errorf("name %q is not unique in go files at path %q", name, path)
			}
			goFile = f
		}
	}
	if goFile == nil {
		return nil, fmt.Errorf("failed to find %q in go files at path %q", name, path)
	}
	return goFile, nil
}

func parseStructTag(literal *ast.BasicLit) (name, description string, omitEmpty bool) {
	if literal == nil {
		return
	}
	tag := reflect.StructTag(strings.Trim(literal.Value, "`"))

	description = tag.Get("description")
	jsonValue := tag.Get("json")
	jsonSplits := strings.Split(jsonValue, ",")
	name = jsonSplits[0]
	if len(jsonSplits) > 1 {
		for _, split := range jsonSplits[1:] {
			if strings.ToLower(split) == "omitempty" {
				omitEmpty = true
				return
			}
		}
	}
	return
}

// writeAvroStruct returns an apply function intended to be called for the start of each node.
// The corresponding Post function will be called after the children of the node have been traversed.
func writeAvroStruct(cfg avroConfig, name, defaultFields string) astutil.ApplyFunc {
	return func(c *astutil.Cursor) bool {
		if c.Name() == "Node" {
			var jsonNamespace string
			if len(cfg.namespace) != 0 {
				jsonNamespace = fmt.Sprintf(`"namespace":%q,`, strings.Join(cfg.namespace, "."))
			}
			if _, err := fmt.Fprintf(cfg.writer, `"name":%q,%s"type":"record","fields":[%s`, name, jsonNamespace, defaultFields); err != nil {
				return false
			}
			return true
		}
		if _, ok := c.Parent().(*ast.TypeSpec); ok {
			// The header includes all needed info for the initial type spec so skip other nodes until we get to the contents
			return true
		}
		n := c.Node()
		if n == nil {
			return false
		}

		if list, ok := n.(*ast.FieldList); ok {
			writeAvroFields(cfg, list)
			return false
		}

		return true
	}
}

// writeAvroField will traverse the field extracting an JSON struct tags and using them to determine the name and
// nullability of the given Avro field. With this information it will then write out the field name.
// If the field was skipped false is returned.
func writeAvroField(cfg avroConfig, f *ast.Field, fieldMap map[string]bool) bool {
	if len(f.Names) == 0 { // An embedded struct
		t, ok := f.Type.(*ast.Ident)
		if !ok {
			return false
		}
		newcfg := cfg
		newcfg.namespace = append(cfg.namespace, t.Name)
		newcfg.excludeFields = fieldMap
		writeEmbeddedStructFields(newcfg)
		return true
	}
	name := f.Names[0].Name

	if _, ok := cfg.excludeFields[name]; ok {
		return false
	}

	tagName, tagDescription, omitEmpty := parseStructTag(f.Tag)

	if tagName != "" {
		name = tagName
	}

	avroType := fmt.Sprintf(`"type":%s`, convertToAvroType(cfg, f.Type, name, omitEmpty))

	avroName := fmt.Sprintf(`"name":%q`, name)
	if len(cfg.namespace) != 0 {
		avroName += fmt.Sprintf(`,"namespace":%q`, strings.Join(cfg.namespace, "."))
	}
	if tagDescription != "" {
		avroName += fmt.Sprintf(`,"doc":%q`, tagDescription)
	}

	fmt.Fprintf(cfg.writer, "{%s,%s}", avroName, avroType)
	return true
}

// writeAvroFields loops through a list of fields in a struct writing out the Avro values for the fields and adding
// trailing commas as appropriate for JSON.
func writeAvroFields(cfg avroConfig, list *ast.FieldList) {
	length := list.NumFields()

	// Build a list of all fields so that embedded fields with the same name at this level can be excluded
	fieldMap := make(map[string]bool, length)
	for _, f := range list.List {
		if len(f.Names) > 0 {
			fieldMap[f.Names[0].Name] = true
		}
	}
	for i, f := range list.List {
		processed := writeAvroField(cfg, f, fieldMap)
		if processed && i+1 != length {
			fmt.Fprint(cfg.writer, ",")
		}
	}
}

func writeEmbeddedStructFields(cfg avroConfig) {
	structName := cfg.namespace[len(cfg.namespace)-1]
	spec, err := parseGoStruct(structName, filepath.Join(cfg.dir, strings.ToLower(structName)+".go"))
	if err != nil {
		fmt.Fprint(cfg.writer, `{"type":"embedded struct not found"}`)
		return
	}

	astutil.Apply(spec, func(c *astutil.Cursor) bool {
		n := c.Node()
		if list, ok := n.(*ast.FieldList); ok {
			writeAvroFields(cfg, list)
			return false
		}
		return true
	}, nil)
}

// convertToAvroType returns the avro type definition for a go type
func convertToAvroType(cfg avroConfig, expr ast.Expr, name string, nullable bool) string {
	// Note: the go code generated from JSON schema does not include maps and they are not handled here
	switch t := expr.(type) {
	case *ast.Ident:
		var typeName string
		switch t.Name {
		case "bool":
			typeName = "boolean"
		case "int", "uint", "int64", "uint64":
			typeName = "long"
		case "int8", "int16", "int32", "uint8", "uint16", "uint32":
			typeName = "int"
		case "float32":
			typeName = "float"
		case "float64":
			typeName = "double"
		case "byte":
			typeName = "bytes"
		case "string":
			typeName = "string"
		default:
			typeName = "unknown"
		}
		if nullable {
			if typeName == "boolean" {
				return fmt.Sprintf(`%q,"default":false`, typeName)
			}
			return fmt.Sprintf(`["null",%q]`, typeName)
		} else {
			return fmt.Sprintf(`%q`, typeName)
		}
	case *ast.ArrayType:
		itemType := convertToAvroType(cfg, t.Elt, name, false)
		if strings.HasPrefix(itemType, "{") {
			return fmt.Sprintf(`{"type":"array","items":%s}`, itemType)
		} else {
			return fmt.Sprintf(`{"type":"array","items":{"type": %s}}`, itemType)
		}
	case *ast.StructType:
		// recursively handle this struct
		buf := &bytes.Buffer{}
		newcfg := cfg
		newcfg.namespace = append(cfg.namespace, name)
		newcfg.writer = buf
		newcfg.excludeFields = nil
		// nested structs get _struct appended on their name
		astutil.Apply(t, writeAvroStruct(newcfg, name+"_record", ""), nil)
		if nullable {
			return fmt.Sprintf(`["null",{%s]}]`, buf.String())
		}
		return fmt.Sprintf("{%s]}", buf.String())
	case *ast.SelectorExpr:
		if t.Sel.Name == "Time" {
			return `{"type":"long","logicalType":"timestamp-millis"}`
		}
		return fmt.Sprintf(`unsupported type %q`, t.Sel.Name)
	default:
		// This should break the Avro schema but with some indication as to why it failed
		return fmt.Sprintf(`unsupported type %q`, astutil.NodeDescription(t))
	}

}
