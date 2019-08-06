package test_data

import (
	"io"
	"time"

	"github.com/GannettDigital/jstransform/generate"
	"github.com/GannettDigital/jstransform/generate/test_data/avro/arrays"

	"github.com/actgardner/gogen-avro/container"
)

// WriteAvroCF writes an Avro Containter File to the given io.Writer using snappy compression for the data.
// The time is used as the AvroWriteTime, if the time is the Zero value then the current time is used.
// If z is nil then the data will be a delete as indicated by the AvroDeleted field.
// NOTE: If the type has a field in an embedded struct with the same name as a field not in the embedded struct the
// value will be pulled from the field not in the embedded struct.
func (z *Arrays) WriteAvroCF(writer io.Writer, writeTime time.Time) error {
	if writeTime.IsZero() {
		writeTime = time.Now()
	}
	avroWriter, err := arrays.NewArraysWriter(writer, container.Snappy, 1)
	if err != nil {
		return err
	}

	if err := avroWriter.WriteRecord(z.convertToAvro(writeTime)); err != nil {
		return err
	}

	return avroWriter.Flush()
}

func (z *Arrays) convertToAvro(writeTime time.Time) *arrays.Arrays {
	aTime := generate.AvroTime(writeTime)
	if z == nil {
		return &arrays.Arrays{AvroWriteTime: aTime, AvroDeleted: true}
	}

	Parents_recordSlice := func(in []struct {
		Count    int64    `json:"count"`
		Children []string `json:"children"`
	}) []*arrays.Parents_record {
		converted := make([]*arrays.Parents_record, len(in))
		for i, z := range in {
			converted[i] = &arrays.Parents_record{
				Count:    z.Count,
				Children: z.Children,
			}
		}
		return converted
	}

	return &arrays.Arrays{
		AvroWriteTime: aTime,
		Heights:       z.Heights,
		Parents:       Parents_recordSlice(z.Parents),
	}
}
