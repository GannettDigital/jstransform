package test_data

import (
	"io"
	"time"

	"github.com/GannettDigital/jstransform/generate"
	"github.com/GannettDigital/jstransform/generate/test_data/avro/complex"

	"github.com/actgardner/gogen-avro/container"
)

// WriteAvroCF writes an Avro Containter File to the given io.Writer using snappy compression for the data.
// The time is used as the AvroWriteTime, if the time is the Zero value then the current time is used.
// If z is nil then the data will be a delete as indicated by the AvroDeleted field.
// NOTE: If the type has a field in an embedded struct with the same name as a field not in the embedded struct the
// value will be pulled from the field not in the embedded struct.
func (z *Complex) WriteAvroCF(writer io.Writer, writeTime time.Time) error {
	if writeTime.IsZero() {
		writeTime = time.Now()
	}
	avroWriter, err := complex.NewComplexWriter(writer, container.Snappy, 1)
	if err != nil {
		return err
	}

	if err := avroWriter.WriteRecord(z.convertToAvro(writeTime)); err != nil {
		return err
	}

	return avroWriter.Flush()
}

func (z *Complex) convertToAvro(writeTime time.Time) *complex.Complex {
	aTime := generate.AvroTime(writeTime)
	if z == nil {
		return &complex.Complex{AvroWriteTime: aTime, AvroDeleted: true}
	}

	Crops_recordSlice := func(in []struct {
		Height       float64 `json:"height"`
		Name         string  `json:"name"`
		Path         string  `json:"path" description:"full path to the cropped image file"`
		RelativePath string  `json:"relativePath" description:"a long"`
		Width        float64 `json:"width"`
	}) []*complex.Crops_record {
		converted := make([]*complex.Crops_record, len(in))
		for i, z := range in {
			converted[i] = &complex.Crops_record{
				Height:       z.Height,
				Name:         z.Name,
				Path:         z.Path,
				RelativePath: z.RelativePath,
				Width:        z.Width,
			}
		}
		return converted
	}

	return &complex.Complex{
		AvroWriteTime:  aTime,
		Height:         &complex.UnionNullLong{Long: z.Height, UnionType: complex.UnionNullLongTypeEnumLong},
		SomeDateObj:    &complex.SomeDateObj_record{Dates: generate.AvroTimeSlice(z.SomeDateObj.Dates)},
		Visible:        z.Visible,
		Width:          &complex.UnionNullDouble{Double: z.Width, UnionType: complex.UnionNullDoubleTypeEnumDouble},
		Caption:        z.Caption,
		Credit:         z.Credit,
		Crops:          Crops_recordSlice(z.Crops),
		Cutline:        &complex.UnionNullString{String: z.Cutline, UnionType: complex.UnionNullStringTypeEnumString},
		DatePhotoTaken: z.DatePhotoTaken,
		Orientation:    z.Orientation,
		OriginalSize: &complex.OriginalSize_record{Height: z.OriginalSize.Height,
			Width: z.OriginalSize.Width},
		Type: z.Type,
		URL: &complex.URL_record{Absolute: z.URL.Absolute,
			Publish: z.URL.Publish},
	}
}
