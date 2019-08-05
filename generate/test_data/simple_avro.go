package test_data

import (
	"io"
	"time"

	"github.com/GannettDigital/jstransform/generate"
	"github.com/GannettDigital/jstransform/generate/test_data/avro/simple"

	"github.com/actgardner/gogen-avro/container"
)

// WriteAvroCF writes an Avro Containter File to the given io.Writer using snappy compression for the data.
// The time is used as the AvroWriteTime, if the time is the Zero value then the current time is used.
// If z is nil then the data will be a delete as indicated by the AvroDeleted field.
// NOTE: If the type has a field in an embedded struct with the same name as a field not in the embedded struct the
// value will be pulled from the field not in the embedded struct.
func (z *Simple) WriteAvroCF(writer io.Writer, writeTime time.Time) error {
	if writeTime.IsZero() {
		writeTime = time.Now()
	}
	avroWriter, err := simple.NewSimpleWriter(writer, container.Snappy, 1)
	if err != nil {
		return err
	}

	if err := avroWriter.WriteRecord(z.convertToAvro(writeTime)); err != nil {
		return err
	}

	return avroWriter.Flush()
}

func (z *Simple) convertToAvro(writeTime time.Time) *simple.Simple {
	aTime := generate.AvroTime(writeTime)
	if z == nil {
		return &simple.Simple{AvroWriteTime: aTime, AvroDeleted: true}
	}

	return &simple.Simple{
		AvroWriteTime: aTime,
		Height:        &simple.UnionNullLong{Long: z.Height, UnionType: simple.UnionNullLongTypeEnumLong},
		SomeDateObj:   &simple.SomeDateObj_record{Dates: generate.AvroTimeSlice(z.SomeDateObj.Dates)},
		Type:          z.Type,
		Visible:       z.Visible,
		Width:         &simple.UnionNullDouble{Double: z.Width, UnionType: simple.UnionNullDoubleTypeEnumDouble},
	}
}
