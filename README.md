# jstransform

This repo provides an extension to [JSON Schema](http://json-schema.org/) which defines a `transform` section which can be added for each field.
This transform section is then used to guide a transformation process which converts JSON input into the format defined by the schema.
The result is that you can write one JSON schema that defines both the desired result and how to transform a known type of data into the defined result.

The code also provides some utilities for walking a JSON schema file section by section and generating Golang structs from a JSON schema file.

## JSON Schema Transform extension
Details on this are found in this [doc](./transform.adoc) and this [schema file](./jsonschema.json)

## Usage
For details on using the project as a library for transformations refer to the godocs.

This project uses the Go package management tool [Dep](https://github.com/golang/dep) for package versioning.
To leverage this tool to install dependencies, run the following command from the project root:

    dep ensure

Testing is done using standard go tooling, ie `go test ./...`

## Examples

The structs in `github.com/GannettDigital/content-api/model/asset` are built from the JSON schema using this code. To regenerate these structs simple run [go generate](https://blog.golang.org/generate) within that directory. This tool was not intended for command-line use, but regeneration from this repo would be:

    go run main.go ../content-api/schema/v1/asset.json ../content-api/model/asset

or any one of them can be transformed individually:

    go run main.go ../content-api/schema/v1/assets/base.json /tmp

