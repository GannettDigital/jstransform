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

// RepeatsBulkAvroWriter will begin a go routine writing an Avro Container File to the writer and add each item from the
// request channel. If an error is encountered it will be sent on the returned error channel.
// The given writeTime will be used for all data items written by this function.
// When the returned request channel is closed this function will finalize the Container File and exit.
// The returned error channel will be closed just before the go routine exits.
func RepeatsBulkAvroWriter(writer io.Writer, writeTime time.Time, request <-chan *Repeats) <-chan error {
	if writeTime.IsZero() {
		writeTime = time.Now()
	}
	errors := make(chan error, 1)

	go func() {
		defer close(errors)

		avroWriter, err := repeats.NewRepeatsWriter(writer, container.Snappy, 1)
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
