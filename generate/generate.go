// Package generate implements a tooling to generate Golang structs from a JSON schema file.
// It is intended to be used with the go generate, https://blog.golang.org/generate
package generate



 
import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/GannettDigital/jstransform/jsonschema"
	"github.com/GannettDigital/msgp/gen"
	"github.com/GannettDigital/msgp/parse"
	"github.com/GannettDigital/msgp/printer"
)

const disclaimer = "// Code generated by github.com/GannettDigital/jstransform; DO NOT EDIT."
const msgpSuffix = "_msgp"
const msgpMode = gen.Encode | gen.Decode | gen.Marshal | gen.Unmarshal | gen.Size | gen.Test

// BuildArgs contains information used to build the structs for a JSONschema.
//  SchemaPath is the path tot he jsonSchema file to use generate the Go struct representations
//
//  OutputDir is the destination for the generated files
//
//  NoNestedStructs will create structs that have no unnamed nested structs in them but rather defined types for each
//  nested struct
//
//  Pointers will create non-required objects and date/time fields with pointers thus allowing the JSON to support null for those fields.
//
//  GenerateAvro is a flag that defines if Avro serializing code should be built.
//  The Avro generated code will only use a single field in the case where a field name is defined in a oneOf and
//  elsewhere in the JSON schema. When converting the most specific version of such a field will be used. In general
//  conflicting names like this should be avoided in the JSON schema.
//
//  GenerateMessagePack is a flag that defines if message pack serializing code should be built.
//
//  StructNameMap allows specifying the type name of the struct for each JSON file.
//
//  FieldNameMap is used to provide alternate names for fields in the resulting structs.
//  The property names in the JSON tags for these structs remains the same as supplied.
//  This can be used to accommodate names that are valid JSON but not valid Go identifiers
type BuildArgs struct {
	SchemaPath             string
	OutputDir              string
	DescriptionAsStructTag bool
	NoNestedStructs        bool
	Pointers               bool
	GenerateAvro           bool
	GenerateMessagePack    bool
	ImportPath             string
	StructNameMap          map[string]string
	FieldNameMap           map[string]string
}

// BuildStructs is a backward-compatibility wrapper for BuildStructsWithArgs.
func BuildStructs(schemaPath string, outputDir string, useMessagePack bool) error {
	return BuildStructsWithArgs(BuildArgs{
		SchemaPath:             schemaPath,
		OutputDir:              outputDir,
		GenerateMessagePack:    useMessagePack,
		DescriptionAsStructTag: true,
	})
}

// BuildStructsRename is a backward-compatibility wrapper for BuildStructsWithArgs.
func BuildStructsRename(schemaPath string, outputDir string, useMessagePack bool, nameMap map[string]string) error {
	return BuildStructsWithArgs(BuildArgs{
		SchemaPath:             schemaPath,
		OutputDir:              outputDir,
		GenerateMessagePack:    useMessagePack,
		StructNameMap:          nameMap,
		DescriptionAsStructTag: true,
	})
}

// BuildStructsWithArgs takes a JSON Schema and generates Golang structs that match the schema.
// The generated structs include struct tags for marshaling/unmarshaling to/from JSON.
// One file will be created for each included allOf/oneOf file in the root schema with any allOf files resulting in
// structs which are embedded in the oneOf files.
//
// The JSON schema can specify more information than the structs enforce (like field size) and so validation of
// any JSON generated from the structs is still necessary.
//
// The args parameter is a BuildArgs struct that defines the settings for this function
//
// If undefined args.OutputDir defaults to the current working directory.
//
// The package name is set to the args.OutputDir directory name.
//
// NOTE: If oneOf/allOf entries exist than any JSON schema instances in the root schema file will be skipped.
func BuildStructsWithArgs(args BuildArgs) error {
	if args.OutputDir == "" {
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to determine working directory: %v", err)
		}
		args.OutputDir = wd
	}

	packageName := filepath.Base(args.OutputDir)

	allOfTypes, oneOfTypes, properties, err := jsonschema.SchemaTypes(args.SchemaPath)
	if err != nil {
		return fmt.Errorf("failed to discover oneOfTypes: %v", err)
	}

	if len(allOfTypes) == 0 && len(oneOfTypes) == 0 || len(properties) > 0 {
		path, err := filepath.Abs(args.SchemaPath)
		if err != nil {
			return fmt.Errorf("failed to determine absolute path of %q: %v", args.SchemaPath, err)
		}
		allOfTypes = append(allOfTypes, path)
	}

	var embeds []string
	for _, allOfPath := range allOfTypes {
		name := strings.Split(filepath.Base(allOfPath), ".")[0]
		if newName, ok := args.StructNameMap[name]; ok {
			name = newName
		} else {
			name = exportedName(name)
		}
		embeds = append(embeds, name)

		path, err := buildStructFile(allOfPath, name, packageName, nil, args)
		if err != nil {
			return fmt.Errorf("failed to build struct file for %q: %v", name, err)
		}
		if args.GenerateAvro {
			if err := buildAvro(name, path, args.ImportPath); err != nil {
				return fmt.Errorf("failed to build Avro files for %q: %v", packageName, err)
			}
		}
	}

	for _, oneOfPath := range oneOfTypes {
		name := strings.Split(filepath.Base(oneOfPath), ".")[0]
		if newName, ok := args.StructNameMap[name]; ok {
			name = newName
		}

		path, err := buildStructFile(oneOfPath, name, packageName, embeds, args)
		if err != nil {
			return fmt.Errorf("failed to build struct file for %q: %v", name, err)
		}
		if args.GenerateAvro {
			if err := buildAvro(name, path, args.ImportPath); err != nil {
				return fmt.Errorf("failed to build Avro files for %q: %v", packageName, err)
			}
		}
	}

	if args.GenerateMessagePack {
		if err := buildMessagePackFile(args.OutputDir, msgpMode); err != nil {
			return fmt.Errorf("failed to build MessagePack file for %q: %v", packageName, err)
		}
	}

	return nil
}

// buildAvro creates an Avro schema file, Avro serialization functions and some helper functions which link the structs
// used by the generated Avro serialization with those created by the BuildStructs functions.
// The serialization methods are created with https://github.com/actgardner/gogen-avro
func buildAvro(name, path, importPath string) error {
	name = exportedName(name)

	avroSchemaPath, err := buildAvroSchemaFile(name, path, false)
	if err != nil {
		return fmt.Errorf("failed to build Avro Schema file: %v", err)
	}

	// step 2 build serialization functions from the new avro schema
	// this step parses the just created Avro schema and by doing so acts as validation step as well
	if err := buildAvroSerializationFunctions(avroSchemaPath); err != nil {
		return fmt.Errorf("failed to build Avro serialization code: %v", err)
	}

	// step 3 generate helper functions
	if err := buildAvroHelperFunctions(name, path, importPath); err != nil {
		return fmt.Errorf("failed to build Avro serialization code: %v", err)
	}
	return nil
}

// buildMessagePackFile generates MessagePack serialization methods for the entire package.
func buildMessagePackFile(outputDir string, mode gen.Method) error {
	fs, err := parse.File(outputDir, false)
	if err != nil {
		return err
	}

	if len(fs.Identities) == 0 {
		return nil
	}

	return printer.PrintFile(filepath.Join(outputDir, fs.Package+msgpSuffix+".go"), fs, mode)
}

// buildStructFile generates the specified struct file.
func buildStructFile(childPath, name, packageName string, embeds []string, args BuildArgs) (string, error) {
	if !filepath.IsAbs(childPath) {
		childPath = filepath.Join(filepath.Dir(args.SchemaPath), childPath)
	}
	schema, err := jsonschema.SchemaFromFile(childPath, name)
	if err != nil {
		return "", err
	}

	generated, err := newGeneratedGoFile(schema, name, packageName, embeds, args)
	if err != nil {
		return "", fmt.Errorf("failed to build generated struct: %v", err)
	}

	outPath := filepath.Join(args.OutputDir, strings.Split(filepath.Base(childPath), ".")[0]+".go")
	gfile, err := os.Create(outPath)
	if err != nil {
		return "", fmt.Errorf("failed to open file %q: %v", outPath, err)
	}
	defer gfile.Close()

	err = generated.write(gfile)
	return outPath, err
}

// exportedName returns a name that is usable as an exported field in Go.
// Only minimal checking on naming is done rather it is assumed the name from the JSON schema is reasonable any
// unacceptable names will likely fail during formatting.
func exportedName(name string) string {
	return strings.Title(name)
}

// goType maps a jsonType to a string representation of the go type.
// If Array is true it makes the type into an array.
// If the JSON Schema had a type of "string" and a format of "date-time" it is expected the input jsonType will be
// "date-time".
// Non-required times are added as pointers to allow for their values to missing go marshalled JSON.
func goType(jsonType string, array, required, pointers bool) string {
	var goType string
	switch jsonType {
	case "boolean":
		goType = "bool"
	case "number":
		goType = "float64"
	case "integer":
		goType = "int64"
	case "string":
		goType = "string"
	case "date-time":
		goType = "time.Time"
		if pointers && !required {
			goType = "*time.Time"
		}
	case "object":
		goType = "struct"
	default:
		goType = jsonType
	}

	if array {
		return "[]" + goType
	}

	return goType
}

// splitJSONPath takes a JSON path and returns an array of path items each of which represents a JSON object with the
// name normalized in a way suitable for using it as a Go struct filed name.
func splitJSONPath(path string) []string {
	var tree []string
	for _, split := range strings.Split(path, ".") {
		if split == "$" {
			continue
		}
		if strings.HasSuffix(split, "[*]") {
			split = split[:len(split)-3]
		}

		tree = append(tree, split)
	}

	return tree
}
