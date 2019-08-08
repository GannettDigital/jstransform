package test_data

import (
	"errors"
	"io"
	"time"

	"github.com/GannettDigital/jstransform/generate"
	"github.com/GannettDigital/jstransform/generate/test_data/avro/arrays"

	"github.com/actgardner/gogen-avro/container"
)

// WriteAvroCF writes an Avro Containter File to the given io.Writer using snappy compression for the data.
// The time is used as the AvroWriteTime, if the time is the Zero value then the current time is used.
// NOTE: If the type has a field in an embedded struct with the same name as a field not in the embedded struct the
// value will be pulled from the field not in the embedded struct.
func (z *Arrays) WriteAvroCF(writer io.Writer, writeTime time.Time) error {
	if z == nil {
		return errors.New("unable to write a nil pointer")
	}
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

// WriteAvroDeletedCF works nearly identically to WriteAvroCF but sets the AvroDeleted metadata field to true.
func (z *Arrays) WriteAvroDeletedCF(writer io.Writer, writeTime time.Time) error {
	if z == nil {
		return errors.New("unable to write a nil pointer")
	}
	if writeTime.IsZero() {
		writeTime = time.Now()
	}
	avroWriter, err := arrays.NewArraysWriter(writer, container.Snappy, 1)
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

// ArraysBulkAvroWriter will begin a go routine writing an Avro Container File to the writer and add each item from the
// request channel. If an error is encountered it will be sent on the returned error channel.
// The given writeTime will be used for all data items written by this function.
// When the returned request channel is closed this function will finalize the Container File and exit.
// The returned error channel will be closed just before the go routine exits.
// Note: That though a nil item will be written as delete it will also be written without an ID or other identifying
// field and so this is of limited value. In general deletes should be done using WriteAvroDeletedCF.
func ArraysBulkAvroWriter(writer io.Writer, writeTime time.Time, request <-chan *Arrays) <-chan error {
	if writeTime.IsZero() {
		writeTime = time.Now()
	}
	errors := make(chan error, 1)

	go func() {
		defer close(errors)

		avroWriter, err := arrays.NewArraysWriter(writer, container.Snappy, 1)
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
