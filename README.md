# jstransform
[![Go Docs](https://godoc.org/github.com/GannettDigital/jstransform?status.svg)](https://godoc.org/github.com/GannettDigital/jstransform)
[![Build Status](https://travis-ci.org/GannettDigital/jstransform.svg)](https://travis-ci.org/GannettDigital/jstransform)
[![Go Report Card](https://goreportcard.com/badge/github.com/GannettDigital/jstransform)](https://goreportcard.com/report/github.com/GannettDigital/jstransform)
[![Coverage Status](https://coveralls.io/repos/github/GannettDigital/jstransform/badge.svg?branch=master)](https://coveralls.io/github/GannettDigital/jstransform?branch=master)


This repo provides an extension to [JSON Schema](http://json-schema.org/) which defines a `transform` section which can be added for each field.
This transform section is then used to guide a transformation process which converts JSON or XML input into the format defined by the schema.
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

As well as generating Go structs serialization libraries for [Message Pack](https://msgpack.org/) and [Avro](https://avro.apache.org/) are available.
For Avro the schema as well as serialization and helper functions are generated.
The Avro schema also adds two metadata fields, `AvroWriteTime` and `AvroDeleted`.
The generated serialization functions are based these libraries:

* https://github.com/GannettDigital/msgp a fork with minor fixes from https://github.com/tinylib/msgp
* https://github.com/actgardner/gogen-avro

## Building/Testing
This project uses Go modules for dependency management. You need to have a working Go environment with version 1.11 or greater installed. 

Testing is done using standard go tooling, ie `go test ./...`
