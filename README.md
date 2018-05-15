# jstransform

This repo provides an extension to [JSON Schema](http://json-schema.org/) which defines a `transform` section which can be added for each field.
This transform section is then used to guide a transformation process which converts JSON input into the format defined by the schema.
The result is that you can write one JSON schema that defines both the desired result and how to transform a known type of data into the defined result.

The code also provides some utilities for walking a JSON schema file section by section and generating Golang structs from a JSON schema file.

## JSON Schema Transform extension
Details on this are found in this [doc](./transform.adoc) and this [schema file](./transformSchema.json)

## Usage
For details on using the project as a library for transformations or JSON schema walking refer to the godocs.

The Golang struct generation portion of this code based is intended to be used with [go generate](https://blog.golang.org/generate).

### Go Generate Examples

To use the struct generation with go generate include a generate line in a go source file for example:

    //go:generate go run ../../vendor/github.com/GannettDigital/jstransform/main.go myschema.json $PWD

or if you have compiled the tool and have it in your path rather than vendoring the source:

    //go:generate jstransform myschema.json $PWD

then simply run `go generate`.

## Building/Testing
This project uses the Go package management tool [Dep](https://github.com/golang/dep) for dependencies.
To leverage this tool to install dependencies, run the following command from the project root:

    dep ensure

Testing is done using standard go tooling, ie `go test ./...`
