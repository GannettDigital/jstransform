package generate

var (
	avroTemplate = `
package {{ .pkgName }}

import (
	"io"
	"time"

	"github.com/GannettDigital/jstransform/generate"
    "{{ .avroImport }}"

	"github.com/actgardner/gogen-avro/container"
)

// WriteAvroCF writes an Avro Containter File to the given io.Writer using snappy compression for the data.
// The time is used as the AvroWriteTime, if the time is the Zero value then the current time is used.
// If z is nil then the data will be a delete as indicated by the AvroDeleted field.
// NOTE: If the type has a field in an embedded struct with the same name as a field not in the embedded struct the
// value will be pulled from the field not in the embedded struct.
func (z *{{ .name }}) WriteAvroCF(writer io.Writer, writeTime time.Time) error {
	if writeTime.IsZero() {
		writeTime = time.Now()
	}
	avroWriter, err := {{ .avroPackage }}.New{{ .name }}Writer(writer, container.Snappy, 1)
	if err != nil {
		return err
	}

	if err := avroWriter.WriteRecord(z.convertToAvro(writeTime)); err != nil {
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
`
	avroStructSliceTemplate = `
{{ .funcName }} := func(in []struct{
{{ .structDef }}
}) []*{{ .typeID }} {
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
)
