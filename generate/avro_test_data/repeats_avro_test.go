package avro_test_data

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/GannettDigital/jstransform/generate"
	"github.com/GannettDigital/jstransform/generate/avro_test_data/avro/repeats"

	"github.com/actgardner/gogen-avro/v7/container"
)

func TestRepeats_WriteAvroCF(t *testing.T) {
	z := &Repeats{}

	// compile error if for some reason z does not implement generate.AvroCFWriter
	var _ generate.AvroCFWriter = z

	// write to a CF file
	now := time.Now()
	buf := &bytes.Buffer{}
	err := z.WriteAvroCF(buf, now)
	if err != nil {
		t.Fatalf("Unexpected error writing to a CF file: %v", err)
	}

	containerReader, err := container.NewReader(buf)
	if err != nil {
		t.Fatalf("Failed containers from OCF file: %v\n", err)
	}

	read, err := repeats.DeserializeRepeatsFromSchema(containerReader, string(containerReader.AvroContainerSchema()))
	if err != nil {
		t.Fatalf("Failed deserializing OCF file: %v\n", err)
	}

	if got, want := read.AvroWriteTime, generate.AvroTime(now); got != want {
		t.Errorf("Time is wrong, got %d, want %d", got, want)
	}

	if read.AvroDeleted != false {
		t.Error("OCF reports deleted true expected false")
	}
}

func TestRepeats_WriteAvroDeletedCF(t *testing.T) {
	z := &Repeats{} // In real usage at minimum an ID field should be populated

	// compile error if for some reason z does not implement generate.AvroCFDeleter
	var _ generate.AvroCFDeleter = z

	// write to a CF file
	now := time.Now()
	buf := &bytes.Buffer{}
	err := z.WriteAvroDeletedCF(buf, now)
	if err != nil {
		t.Fatalf("Unexpected error writing to a CF file: %v", err)
	}

	containerReader, err := container.NewReader(buf)
	if err != nil {
		t.Fatalf("Failed containers from OCF file: %v\n", err)
	}

	read, err := repeats.DeserializeRepeatsFromSchema(containerReader, string(containerReader.AvroContainerSchema()))
	if err != nil {
		t.Fatalf("Failed deserializing OCF file: %v\n", err)
	}

	if got, want := read.AvroWriteTime, generate.AvroTime(now); got != want {
		t.Errorf("Time is wrong, got %d, want %d", got, want)
	}

	if read.AvroDeleted != true {
		t.Error("OCF reports deleted false expected true")
	}
}

func ExampleRepeatsBulkAvroWriter() {
	input := []*Repeats{{}, {}, {}}
	inputChan := make(chan *Repeats)

	devnull, _ := os.Open("/dev/null")
	defer devnull.Close()

	errChan := RepeatsBulkAvroWriter(devnull, time.Now(), inputChan)

	for _, item := range input {
		select {
		case err := <-errChan:
			fmt.Print(err)
			return
		case inputChan <- item:
		}
	}

	// Check for any final errors, the errorChan should be closed when the BulkWriter is finished processing
	for err := range errChan {
		if err != nil {
			fmt.Print(err)
			return
		}
	}
}
