package test_data

import (
	"bytes"
	"testing"
	"time"

	"github.com/GannettDigital/jstransform/generate"
	"github.com/GannettDigital/jstransform/generate/test_data/avro/simple"
)

func TestSimple_WriteAvroCF(t *testing.T) {
	z := &Simple{}

	// compile error if for some reason z does not implement generate.AvroCFWriter
	var _ generate.AvroCFWriter = z

	// write to a CF file
	now := time.Now()
	buf := &bytes.Buffer{}
	err := z.WriteAvroCF(buf, now)
	if err != nil {
		t.Fatalf("Unexpected error writing to a CF file: %v", err)
	}

	ocfReader, err := simple.NewSimpleReader(buf)
	if err != nil {
		t.Fatalf("Error creating OCF file reader: %v\n", err)
	}

	read, err := ocfReader.Read()
	if err != nil {
		t.Fatalf("Failed reading from OCF file reader: %v\n", err)
	}

	if got, want := read.AvroWriteTime, generate.AvroTime(now); got != want {
		t.Errorf("Time is wrong, got %d, want %d", got, want)
	}

	if read.AvroDeleted != false {
		t.Error("OCF reports deleted true expected false")
	}
}
