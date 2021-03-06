package generate

var (
	avroTemplate = `
package {{ .pkgName }}

import (
    "errors"
	"io"
	"time"

	"github.com/GannettDigital/jstransform/generate"
    "{{ .avroImport }}"

	"github.com/actgardner/gogen-avro/v7/container"
)

// WriteAvroCF writes an Avro Containter File to the given io.Writer using snappy compression for the data.
// The time is used as the AvroWriteTime, if the time is the Zero value then the current time is used.
// NOTE: If the type has a field in an embedded struct with the same name as a field not in the embedded struct the
// value will be pulled from the field not in the embedded struct.
func (z *{{ .name }}) WriteAvroCF(writer io.Writer, writeTime time.Time) error {
	if z == nil {
		return errors.New("unable to write a nil pointer")
	}
	if writeTime.IsZero() {
		writeTime = time.Now()
	}
	avroWriter, err := container.NewWriter(writer, container.Snappy, 1, {{ .avroPackage }}.New{{ .name }}().Schema())
	if err != nil {
		return err
	}

	if err := avroWriter.WriteRecord(z.convertToAvro(writeTime)); err != nil {
		return err
	}

	return avroWriter.Flush()
}

// WriteAvroDeletedCF works nearly identically to WriteAvroCF but sets the AvroDeleted metadata field to true.
func (z *{{ .name }}) WriteAvroDeletedCF(writer io.Writer, writeTime time.Time) error {
	if z == nil {
		return errors.New("unable to write a nil pointer")
	}
	if writeTime.IsZero() {
		writeTime = time.Now()
	}
	avroWriter, err := container.NewWriter(writer, container.Snappy, 1, {{ .avroPackage }}.New{{ .name }}().Schema())
	if err != nil {
		return err
	}

	converted := z.convertToAvro(writeTime)
	converted.AvroDeleted = true
	if err := avroWriter.WriteRecord(converted); err != nil {
		return err
	}

	return avroWriter.Flush()
}

func (z *{{ .name }}) convertToAvro(writeTime time.Time) *{{ .avroPackage}}.{{ .name }} {
	aTime := generate.AvroTime(writeTime)
	if z == nil {
		return &{{ .avroPackage}}.{{ .name }}{AvroWriteTime: aTime, AvroDeleted: true}
	}

  {{ .preProcessing }}

	return &{{ .avroPackage}}.{{ .name }}{
	  AvroWriteTime: aTime,
	  {{ .fieldMapping }}
	}
}

// {{ .name }}BulkAvroWriter will begin a go routine writing an Avro Container File to the writer and add each item from the
// request channel. If an error is encountered it will be sent on the returned error channel.
// The given writeTime will be used for all data items written by this function.
// When the returned request channel is closed this function will finalize the Container File and exit.
// The returned error channel will be closed just before the go routine exits.
// Note: That though a nil item will be written as delete it will also be written without an ID or other identifying
// field and so this is of limited value. In general deletes should be done using WriteAvroDeletedCF.
func {{ .name }}BulkAvroWriter(writer io.Writer, writeTime time.Time, request <-chan *{{ .name }}) <-chan error {
	if writeTime.IsZero() {
		writeTime = time.Now()
	}
	errors := make(chan error, 1)

	go func() {
		defer close(errors)

		avroWriter, err := container.NewWriter(writer, container.Snappy, 1, {{ .avroPackage }}.New{{ .name }}().Schema())
		if err != nil {
			errors <- err
			return
		}

		for item := range request {
			if err := avroWriter.WriteRecord(item.convertToAvro(writeTime)); err != nil {
				errors <- err
				return
			}
		}

		if err := avroWriter.Flush(); err != nil {
			errors <- err
			return
		}
	}()
	return errors
}
`
	avroTestTemplate = `
package {{ .pkgName }}

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/GannettDigital/jstransform/generate"
    "{{ .avroImport }}"

	"github.com/actgardner/gogen-avro/v7/container"
)

func Test{{ .name }}_WriteAvroCF(t *testing.T) {
	z := &{{ .name }}{}

	// compile error if for some reason z does not implement generate.AvroCFWriter
	var _ generate.AvroCFWriter = z

	// write to a CF file
	now := time.Now()
	buf := &bytes.Buffer{}
	err := z.WriteAvroCF(buf, now)
	if err != nil {
		t.Fatalf("Unexpected error writing to a CF file: %v", err)
	}

	containerReader, err := container.NewReader(buf)
	if err != nil {
		t.Fatalf("Failed containers from OCF file: %v\n", err)
	}

	read, err := {{ .avroPackage}}.Deserialize{{ .name }}FromSchema(containerReader, string(containerReader.AvroContainerSchema()))
	if err != nil {
		t.Fatalf("Failed deserializing OCF file: %v\n", err)
	}

	if got, want := read.AvroWriteTime, generate.AvroTime(now); got != want {
		t.Errorf("Time is wrong, got %d, want %d", got, want)
	}

	if read.AvroDeleted != false {
		t.Error("OCF reports deleted true expected false")
	}
}

func Test{{ .name }}_WriteAvroDeletedCF(t *testing.T) {
	z := &{{ .name }}{} // In real usage at minimum an ID field should be populated

	// compile error if for some reason z does not implement generate.AvroCFDeleter
	var _ generate.AvroCFDeleter = z

	// write to a CF file
	now := time.Now()
	buf := &bytes.Buffer{}
	err := z.WriteAvroDeletedCF(buf, now)
	if err != nil {
		t.Fatalf("Unexpected error writing to a CF file: %v", err)
	}

	containerReader, err := container.NewReader(buf)
	if err != nil {
		t.Fatalf("Failed containers from OCF file: %v\n", err)
	}

	read, err := {{ .avroPackage}}.Deserialize{{ .name }}FromSchema(containerReader, string(containerReader.AvroContainerSchema()))
	if err != nil {
		t.Fatalf("Failed deserializing OCF file: %v\n", err)
	}

	if got, want := read.AvroWriteTime, generate.AvroTime(now); got != want {
		t.Errorf("Time is wrong, got %d, want %d", got, want)
	}

	if read.AvroDeleted != true {
		t.Error("OCF reports deleted false expected true")
	}
}

func Example{{ .name }}BulkAvroWriter() {
	input := []*{{ .name }}{ {}, {}, {} }
	inputChan := make(chan *{{ .name }})

	devnull, _ := os.Open("/dev/null")
	defer devnull.Close()

	errChan := {{ .name }}BulkAvroWriter(devnull, time.Now(), inputChan)

	for _, item := range input {
		select {
		case err := <-errChan:
			fmt.Print(err)
			return
		case inputChan <- item:
		}
	}

	// Check for any final errors, the errorChan should be closed when the BulkWriter is finished processing
	for err := range errChan {
		if err != nil {
			fmt.Print(err)
			return
		}
	}
}
`
	avroStructSliceTemplate = `
{{ .funcName }} := func(in []{{ .structDef }}) []*{{ .typeID }} {
	converted := make([]*{{ .typeID }}, len(in))
	for i, z := range in {
		converted[i] = &{{ .typeID }} {
			{{ .structFields }}
		}
	}
	return converted
}
`
	unionNullDouble = `{{ .packageName }}.UnionNullDouble{Double: {{ .value }}, UnionType: {{ .packageName }}.UnionNullDoubleTypeEnumDouble}`
	unionNullLong   = `{{ .packageName }}.UnionNullLong{Long: {{ .value }}, UnionType: {{ .packageName }}.UnionNullLongTypeEnumLong}`
	unionNullString = `{{ .packageName }}.UnionNullString{String: {{ .value }}, UnionType: {{ .packageName }}.UnionNullStringTypeEnumString}`
	unionNullStruct = `{{ .packageName }}.UnionNull{{ .typeName }} { {{ .typeName }}: &{{ .value }}, UnionType: {{ .packageName }}.UnionNull{{ .typeName }}TypeEnum{{ .typeName }}}`
)
