package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/GannettDigital/jstransform/generate"
)

// mapFlags allows for "-opt key=value" flags.
type mapFlags struct {
	kv map[string]string
}

// String converts the provided options into the input format.
func (mf *mapFlags) String() string {
	kvs := make([]string, 0, len(mf.kv))
	for k, v := range mf.kv {
		kvs = append(kvs, fmt.Sprintf("-%s=%s", k, v))
	}
	return strings.Join(kvs, " ")
}

// Set stores the provided options into the data structure. Multiple key-values pairs can be provided if separated by a comma.
func (mf *mapFlags) Set(value string) error {
	for _, pair := range strings.Split(value, ",") {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			return fmt.Errorf("value in 'key=value' format: %v", value)
		}
		mf.kv[kv[0]] = kv[1]
	}
	return nil
}

func main() {
	renameStructs := mapFlags{kv: make(map[string]string)}
	renameFields := mapFlags{kv: make(map[string]string)}
	renameGQLType := mapFlags{kv: make(map[string]string)}

	nestedStructs := flag.Bool("nestedStructs", true, "Build struct with unnamed nested structs, if false each nested struct is made its own type.")
	pointers := flag.Bool("pointers", false, "Build non-required JSON objects and date time fields as struct pointers")
	descriptionAsStructTag := flag.Bool("descriptionAsStructTag", true, "Include the description as a struct tag, rather than a comment")
	flag.Var(&renameStructs, "rename", "Override generated name of structure; use '-rename old=new'.")
	flag.Var(&renameFields, "renameFields", "Override generated name of structure; use '-renameFields old=new'.")
	flag.Var(&renameGQLType, "renameGraphQLType", "Override generated name of GraphQL type; use '-renameGraphQLType old1=new1,old2=new2' where old is the schema file name.")
	genAvro := flag.Bool("avro", false, "generate Avro schema and serialization methods")
	genMessagePack := flag.Bool("msgp", false, "generate MessagePack serialization methods")
	genGraphQL := flag.Bool("graphql", false, "generate GraphQL schema")
	outputPathGraphQL := flag.String("outputPathGraphQL", "", "The output path where to write GraphQL schema.")
	importPath := flag.String("importPath", "", "The import path used as the base for generated code, required for Avro")

	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		fmt.Printf("Usage: %s [-avro] [-importPath a/b/c] [-msgp] [-rename k=v] [-renameFields k=v] [-renameGraphQLType k=v] [-graphql] [-outputPathGraphQL a/b/c] <JSON Schema Path> [output directory]\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *genAvro && *importPath == "" {
		fmt.Printf("Avro requires specifying an import path.\n")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *genGraphQL && *importPath == "" {
		fmt.Printf("GraphQL requires specifying an import path.\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	inputPath, err := filepath.Abs(args[0])
	if err != nil {
		fmt.Printf("Input directory \"%s\" error: %v", args[0], err)
		os.Exit(2)
	}

	var outputPath string
	if len(args) > 1 {
		outputPath, err = filepath.Abs(args[1])
		if err != nil {
			fmt.Printf("Output directory \"%s\" error: %v", args[1], err)
			os.Exit(3)
		}
	}

	if err = generate.BuildStructsWithArgs(generate.BuildArgs{
		SchemaPath:             inputPath,
		OutputDir:              outputPath,
		OutputDirGraphQL:       *outputPathGraphQL,
		GenerateAvro:           *genAvro,
		GenerateMessagePack:    *genMessagePack,
		GenerateGraphQL:        *genGraphQL,
		ImportPath:             *importPath,
		DescriptionAsStructTag: *descriptionAsStructTag,
		NoNestedStructs:        !*nestedStructs,
		Pointers:               *pointers,
		StructNameMap:          renameStructs.kv,
		FieldNameMap:           renameFields.kv,
		GraphQLTypeNameMap:     renameGQLType.kv,
	}); err != nil {
		fmt.Printf("Golang Struct generation failed: %v\n", err)
		os.Exit(4)
	}
}
