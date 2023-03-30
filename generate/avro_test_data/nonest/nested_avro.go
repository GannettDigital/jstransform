package nonest

import (
	"errors"
	"io"
	"time"

	"github.com/GannettDigital/jstransform/generate"
	"github.com/GannettDigital/jstransform/generate/avro_test_data/nonest/avro/nested"

	"github.com/actgardner/gogen-avro/v7/container"
)

// WriteAvroCF writes an Avro Containter File to the given io.Writer using snappy compression for the data.
// The time is used as the AvroWriteTime, if the time is the Zero value then the current time is used.
// NOTE: If the type has a field in an embedded struct with the same name as a field not in the embedded struct the
// value will be pulled from the field not in the embedded struct.
func (z *Nested) WriteAvroCF(writer io.Writer, writeTime time.Time) error {
	if z == nil {
		return errors.New("unable to write a nil pointer")
	}
	if writeTime.IsZero() {
		writeTime = time.Now()
	}
	avroWriter, err := container.NewWriter(writer, container.Snappy, 1, nested.NewNested().Schema())
	if err != nil {
		return err
	}

	if err := avroWriter.WriteRecord(z.convertToAvro(writeTime)); err != nil {
		return err
	}

	return avroWriter.Flush()
}

// WriteAvroDeletedCF works nearly identically to WriteAvroCF but sets the AvroDeleted metadata field to true.
func (z *Nested) WriteAvroDeletedCF(writer io.Writer, writeTime time.Time) error {
	if z == nil {
		return errors.New("unable to write a nil pointer")
	}
	if writeTime.IsZero() {
		writeTime = time.Now()
	}
	avroWriter, err := container.NewWriter(writer, container.Snappy, 1, nested.NewNested().Schema())
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

func (z *Nested) convertToAvro(writeTime time.Time) *nested.Nested {
	aTime := generate.AvroTime(writeTime)
	if z == nil {
		return &nested.Nested{AvroWriteTime: aTime, AvroDeleted: true}
	}

	NestedFactCheckClaimsAppearanceURLs_AppearanceURLs_recordSlice := func(in []NestedFactCheckClaimsAppearanceURLs) []*nested.AppearanceURLs_record {
		converted := make([]*nested.AppearanceURLs_record, len(in))
		for i, z := range in {
			converted[i] = &nested.AppearanceURLs_record{
				Original: z.Original,
				Url:      z.Url,
			}
		}
		return converted
	}

	NestedFactCheckClaims_FactCheckClaims_recordSlice := func(in []NestedFactCheckClaims) []*nested.FactCheckClaims_record {
		converted := make([]*nested.FactCheckClaims_record, len(in))
		for i, z := range in {
			converted[i] = &nested.FactCheckClaims_record{
				AppearanceURLs: NestedFactCheckClaimsAppearanceURLs_AppearanceURLs_recordSlice(z.AppearanceURLs),
				Author:         &nested.UnionNullString{String: z.Author, UnionType: nested.UnionNullStringTypeEnumString},
				Claim:          &nested.UnionNullString{String: z.Claim, UnionType: nested.UnionNullStringTypeEnumString},
				Date:           &nested.UnionNullString{String: z.Date, UnionType: nested.UnionNullStringTypeEnumString},
				Rating:         &nested.UnionNullString{String: z.Rating, UnionType: nested.UnionNullStringTypeEnumString},
			}
		}
		return converted
	}

	return &nested.Nested{
		AvroWriteTime:   aTime,
		FactCheckClaims: NestedFactCheckClaims_FactCheckClaims_recordSlice(z.FactCheckClaims),
	}
}

// NestedBulkAvroWriter will begin a go routine writing an Avro Container File to the writer and add each item from the
// request channel. If an error is encountered it will be sent on the returned error channel.
// The given writeTime will be used for all data items written by this function.
// When the returned request channel is closed this function will finalize the Container File and exit.
// The returned error channel will be closed just before the go routine exits.
// Note: That though a nil item will be written as delete it will also be written without an ID or other identifying
// field and so this is of limited value. In general deletes should be done using WriteAvroDeletedCF.
func NestedBulkAvroWriter(writer io.Writer, writeTime time.Time, request <-chan *Nested) <-chan error {
	if writeTime.IsZero() {
		writeTime = time.Now()
	}
	errors := make(chan error, 1)

	go func() {
		defer close(errors)

		avroWriter, err := container.NewWriter(writer, container.Snappy, 1, nested.NewNested().Schema())
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
