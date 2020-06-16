package test_data

import (
	"errors"
	"io"
	"time"

	"github.com/GannettDigital/jstransform/generate"
	"github.com/GannettDigital/jstransform/generate/test_data/avro/simple"

	"github.com/actgardner/gogen-avro/v7/container"
)

// WriteAvroCF writes an Avro Containter File to the given io.Writer using snappy compression for the data.
// The time is used as the AvroWriteTime, if the time is the Zero value then the current time is used.
// NOTE: If the type has a field in an embedded struct with the same name as a field not in the embedded struct the
// value will be pulled from the field not in the embedded struct.
func (z *Simple) WriteAvroCF(writer io.Writer, writeTime time.Time) error {
	if z == nil {
		return errors.New("unable to write a nil pointer")
	}
	if writeTime.IsZero() {
		writeTime = time.Now()
	}
	avroWriter, err := container.NewWriter(writer, container.Snappy, 1, simple.NewSimple().Schema())
	if err != nil {
		return err
	}

	if err := avroWriter.WriteRecord(z.convertToAvro(writeTime)); err != nil {
		return err
	}

	return avroWriter.Flush()
}

// WriteAvroDeletedCF works nearly identically to WriteAvroCF but sets the AvroDeleted metadata field to true.
func (z *Simple) WriteAvroDeletedCF(writer io.Writer, writeTime time.Time) error {
	if z == nil {
		return errors.New("unable to write a nil pointer")
	}
	if writeTime.IsZero() {
		writeTime = time.Now()
	}
	avroWriter, err := container.NewWriter(writer, container.Snappy, 1, simple.NewSimple().Schema())
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

func (z *Simple) convertToAvro(writeTime time.Time) *simple.Simple {
	aTime := generate.AvroTime(writeTime)
	if z == nil {
		return &simple.Simple{AvroWriteTime: aTime, AvroDeleted: true}
	}

	Contributors_recordSlice := func(in []struct {
		ContributorId string `json:"contributorId,omitempty"`
		Id            string `json:"id"`
		Name          string `json:"name"`
	}) []*simple.Contributors_record {
		converted := make([]*simple.Contributors_record, len(in))
		for i, z := range in {
			converted[i] = &simple.Contributors_record{
				ContributorId: &simple.UnionNullString{String: z.ContributorId, UnionType: simple.UnionNullStringTypeEnumString},
				Id:            z.Id,
				Name:          z.Name,
			}
		}
		return converted
	}

	return &simple.Simple{
		AvroWriteTime: aTime,
		Contributors:  Contributors_recordSlice(z.Contributors),
		Height:        &simple.UnionNullLong{Long: z.Height, UnionType: simple.UnionNullLongTypeEnumLong},
		SomeDateObj:   &simple.UnionNullSomeDateObj_record{SomeDateObj_record: &simple.SomeDateObj_record{Dates: generate.AvroTimeSlice(z.SomeDateObj.Dates)}, UnionType: simple.UnionNullSomeDateObj_recordTypeEnumSomeDateObj_record},
		Type:          z.Type,
		Visible:       z.Visible,
		Width:         &simple.UnionNullDouble{Double: z.Width, UnionType: simple.UnionNullDoubleTypeEnumDouble},
	}
}

// SimpleBulkAvroWriter will begin a go routine writing an Avro Container File to the writer and add each item from the
// request channel. If an error is encountered it will be sent on the returned error channel.
// The given writeTime will be used for all data items written by this function.
// When the returned request channel is closed this function will finalize the Container File and exit.
// The returned error channel will be closed just before the go routine exits.
// Note: That though a nil item will be written as delete it will also be written without an ID or other identifying
// field and so this is of limited value. In general deletes should be done using WriteAvroDeletedCF.
func SimpleBulkAvroWriter(writer io.Writer, writeTime time.Time, request <-chan *Simple) <-chan error {
	if writeTime.IsZero() {
		writeTime = time.Now()
	}
	errors := make(chan error, 1)

	go func() {
		defer close(errors)

		avroWriter, err := container.NewWriter(writer, container.Snappy, 1, simple.NewSimple().Schema())
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
