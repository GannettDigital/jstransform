package test_data

import (
	"io"
	"time"

	"github.com/GannettDigital/jstransform/generate"
	"github.com/GannettDigital/jstransform/generate/test_data/avro/repeats"

	"github.com/actgardner/gogen-avro/container"
)

// WriteAvroCF writes an Avro Containter File to the given io.Writer using snappy compression for the data.
// The time is used as the AvroWriteTime, if the time is the Zero value then the current time is used.
// If z is nil then the data will be a delete as indicated by the AvroDeleted field.
// NOTE: If the type has a field in an embedded struct with the same name as a field not in the embedded struct the
// value will be pulled from the field not in the embedded struct.
func (z *Repeats) WriteAvroCF(writer io.Writer, writeTime time.Time) error {
	if writeTime.IsZero() {
		writeTime = time.Now()
	}
	avroWriter, err := repeats.NewRepeatsWriter(writer, container.Snappy, 1)
	if err != nil {
		return err
	}

	if err := avroWriter.WriteRecord(z.convertToAvro(writeTime)); err != nil {
		return err
	}

	return avroWriter.Flush()
}

func (z *Repeats) convertToAvro(writeTime time.Time) *repeats.Repeats {
	aTime := generate.AvroTime(writeTime)
	if z == nil {
		return &repeats.Repeats{AvroWriteTime: aTime, AvroDeleted: true}
	}

	return &repeats.Repeats{
		AvroWriteTime: aTime,
		Height:        &repeats.UnionNullLong{Long: z.Height, UnionType: repeats.UnionNullLongTypeEnumLong},
		SomeDateObj: &repeats.SomeDateObj_record{Type: z.SomeDateObj.Type,
			Visible: z.SomeDateObj.Visible},
		Type:    z.Type,
		Visible: z.Visible,
		Width:   &repeats.UnionNullDouble{Double: z.Width, UnionType: repeats.UnionNullDoubleTypeEnumDouble},
	}
}
