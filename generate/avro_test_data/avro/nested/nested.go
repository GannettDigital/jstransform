// Code generated by github.com/actgardner/gogen-avro/v7. DO NOT EDIT.
/*
 * SOURCE:
 *     nested.avsc.out
 */
package nested

import (
	"github.com/actgardner/gogen-avro/v7/compiler"
	"github.com/actgardner/gogen-avro/v7/vm"
	"github.com/actgardner/gogen-avro/v7/vm/types"
	"io"
)

type Nested struct {
	// The timestamp when this avro data is written. Useful for identifying the newest row of data sharing keys.
	AvroWriteTime int64 `json:"AvroWriteTime"`
	// This is set to true when the Avro data is recording a delete in the source data.
	AvroDeleted bool `json:"AvroDeleted"`

	FactCheckClaims []*FactCheckClaims_record `json:"factCheckClaims"`
}

const NestedAvroCRC64Fingerprint = "\x17\xf9\x83\xcc\x03\xfa\xfb)"

func NewNested() *Nested {
	return &Nested{}
}

func DeserializeNested(r io.Reader) (*Nested, error) {
	t := NewNested()
	deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
	if err != nil {
		return nil, err
	}

	err = vm.Eval(r, deser, t)
	if err != nil {
		return nil, err
	}
	return t, err
}

func DeserializeNestedFromSchema(r io.Reader, schema string) (*Nested, error) {
	t := NewNested()

	deser, err := compiler.CompileSchemaBytes([]byte(schema), []byte(t.Schema()))
	if err != nil {
		return nil, err
	}

	err = vm.Eval(r, deser, t)
	if err != nil {
		return nil, err
	}
	return t, err
}

func writeNested(r *Nested, w io.Writer) error {
	var err error
	err = vm.WriteLong(r.AvroWriteTime, w)
	if err != nil {
		return err
	}
	err = vm.WriteBool(r.AvroDeleted, w)
	if err != nil {
		return err
	}
	err = writeArrayFactCheckClaims_record(r.FactCheckClaims, w)
	if err != nil {
		return err
	}
	return err
}

func (r *Nested) Serialize(w io.Writer) error {
	return writeNested(r, w)
}

func (r *Nested) Schema() string {
	return "{\"fields\":[{\"doc\":\"The timestamp when this avro data is written. Useful for identifying the newest row of data sharing keys.\",\"logicalType\":\"timestamp-millis\",\"name\":\"AvroWriteTime\",\"type\":\"long\"},{\"default\":false,\"doc\":\"This is set to true when the Avro data is recording a delete in the source data.\",\"name\":\"AvroDeleted\",\"type\":\"boolean\"},{\"name\":\"factCheckClaims\",\"type\":{\"items\":{\"fields\":[{\"name\":\"appearanceURLs\",\"namespace\":\"factCheckClaims\",\"type\":{\"items\":{\"fields\":[{\"default\":false,\"name\":\"original\",\"namespace\":\"factCheckClaims.appearanceURLs\",\"type\":\"boolean\"},{\"name\":\"url\",\"namespace\":\"factCheckClaims.appearanceURLs\",\"type\":\"string\"}],\"name\":\"appearanceURLs_record\",\"namespace\":\"factCheckClaims.appearanceURLs\",\"type\":\"record\"},\"type\":\"array\"}},{\"name\":\"author\",\"namespace\":\"factCheckClaims\",\"type\":[\"null\",\"string\"]},{\"name\":\"claim\",\"namespace\":\"factCheckClaims\",\"type\":[\"null\",\"string\"]},{\"name\":\"date\",\"namespace\":\"factCheckClaims\",\"type\":[\"null\",\"string\"]},{\"name\":\"rating\",\"namespace\":\"factCheckClaims\",\"type\":[\"null\",\"string\"]}],\"name\":\"factCheckClaims_record\",\"namespace\":\"factCheckClaims\",\"type\":\"record\"},\"type\":\"array\"}}],\"name\":\"Nested\",\"type\":\"record\"}"
}

func (r *Nested) SchemaName() string {
	return "Nested"
}

func (_ *Nested) SetBoolean(v bool)    { panic("Unsupported operation") }
func (_ *Nested) SetInt(v int32)       { panic("Unsupported operation") }
func (_ *Nested) SetLong(v int64)      { panic("Unsupported operation") }
func (_ *Nested) SetFloat(v float32)   { panic("Unsupported operation") }
func (_ *Nested) SetDouble(v float64)  { panic("Unsupported operation") }
func (_ *Nested) SetBytes(v []byte)    { panic("Unsupported operation") }
func (_ *Nested) SetString(v string)   { panic("Unsupported operation") }
func (_ *Nested) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r *Nested) Get(i int) types.Field {
	switch i {
	case 0:
		return &types.Long{Target: &r.AvroWriteTime}
	case 1:
		return &types.Boolean{Target: &r.AvroDeleted}
	case 2:
		r.FactCheckClaims = make([]*FactCheckClaims_record, 0)

		return &ArrayFactCheckClaims_recordWrapper{Target: &r.FactCheckClaims}
	}
	panic("Unknown field index")
}

func (r *Nested) SetDefault(i int) {
	switch i {
	case 1:
		r.AvroDeleted = false
		return
	}
	panic("Unknown field index")
}

func (r *Nested) NullField(i int) {
	switch i {
	}
	panic("Not a nullable field index")
}

func (_ *Nested) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ *Nested) AppendArray() types.Field         { panic("Unsupported operation") }
func (_ *Nested) Finalize()                        {}

func (_ *Nested) AvroCRC64Fingerprint() []byte {
	return []byte(NestedAvroCRC64Fingerprint)
}
