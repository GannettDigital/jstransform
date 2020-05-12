// Code generated by github.com/actgardner/gogen-avro. DO NOT EDIT.
package arrays

import (
	"github.com/actgardner/gogen-avro/v7/compiler"
	"github.com/actgardner/gogen-avro/v7/vm"
	"github.com/actgardner/gogen-avro/v7/vm/types"
	"io"
)

type Arrays struct {
	// The timestamp when this avro data is written. Useful for identifying the newest row of data sharing keys.
	AvroWriteTime int64 `json:"AvroWriteTime"`
	// This is set to true when the Avro data is recording a delete in the source data.
	AvroDeleted bool `json:"AvroDeleted"`

	Heights []int64 `json:"heights"`

	Parents []*Parents_record `json:"parents"`
}

const ArraysAvroCRC64Fingerprint = "\x97pH\xcf\x06ˬ)"

func NewArrays() *Arrays {
	return &Arrays{}
}

func DeserializeArrays(r io.Reader) (*Arrays, error) {
	t := NewArrays()
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

func DeserializeArraysFromSchema(r io.Reader, schema string) (*Arrays, error) {
	t := NewArrays()

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

func writeArrays(r *Arrays, w io.Writer) error {
	var err error
	err = vm.WriteLong(r.AvroWriteTime, w)
	if err != nil {
		return err
	}
	err = vm.WriteBool(r.AvroDeleted, w)
	if err != nil {
		return err
	}
	err = writeArrayLong(r.Heights, w)
	if err != nil {
		return err
	}
	err = writeArrayParents_record(r.Parents, w)
	if err != nil {
		return err
	}
	return err
}

func (r *Arrays) Serialize(w io.Writer) error {
	return writeArrays(r, w)
}

func (r *Arrays) Schema() string {
	return "{\"fields\":[{\"doc\":\"The timestamp when this avro data is written. Useful for identifying the newest row of data sharing keys.\",\"logicalType\":\"timestamp-millis\",\"name\":\"AvroWriteTime\",\"type\":\"long\"},{\"default\":false,\"doc\":\"This is set to true when the Avro data is recording a delete in the source data.\",\"name\":\"AvroDeleted\",\"type\":\"boolean\"},{\"name\":\"heights\",\"type\":{\"items\":{\"type\":\"long\"},\"type\":\"array\"}},{\"name\":\"parents\",\"type\":{\"items\":{\"fields\":[{\"name\":\"count\",\"namespace\":\"parents\",\"type\":\"long\"},{\"name\":\"children\",\"namespace\":\"parents\",\"type\":{\"items\":{\"type\":\"string\"},\"type\":\"array\"}},{\"name\":\"date\",\"namespace\":\"parents\",\"type\":{\"logicalType\":\"timestamp-millis\",\"type\":\"long\"}},{\"name\":\"info\",\"namespace\":\"parents\",\"type\":{\"fields\":[{\"name\":\"name\",\"namespace\":\"parents.info\",\"type\":\"string\"},{\"name\":\"age\",\"namespace\":\"parents.info\",\"type\":\"long\"}],\"name\":\"info_record\",\"namespace\":\"parents.info\",\"type\":\"record\"}}],\"name\":\"parents_record\",\"namespace\":\"parents\",\"type\":\"record\"},\"type\":\"array\"}}],\"name\":\"Arrays\",\"type\":\"record\"}"
}

func (r *Arrays) SchemaName() string {
	return "Arrays"
}

func (_ *Arrays) SetBoolean(v bool)    { panic("Unsupported operation") }
func (_ *Arrays) SetInt(v int32)       { panic("Unsupported operation") }
func (_ *Arrays) SetLong(v int64)      { panic("Unsupported operation") }
func (_ *Arrays) SetFloat(v float32)   { panic("Unsupported operation") }
func (_ *Arrays) SetDouble(v float64)  { panic("Unsupported operation") }
func (_ *Arrays) SetBytes(v []byte)    { panic("Unsupported operation") }
func (_ *Arrays) SetString(v string)   { panic("Unsupported operation") }
func (_ *Arrays) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r *Arrays) Get(i int) types.Field {
	switch i {
	case 0:
		return &types.Long{Target: &r.AvroWriteTime}
	case 1:
		return &types.Boolean{Target: &r.AvroDeleted}
	case 2:
		r.Heights = make([]int64, 0)

		return &ArrayLongWrapper{Target: &r.Heights}
	case 3:
		r.Parents = make([]*Parents_record, 0)

		return &ArrayParents_recordWrapper{Target: &r.Parents}
	}
	panic("Unknown field index")
}

func (r *Arrays) SetDefault(i int) {
	switch i {
	case 1:
		r.AvroDeleted = false
		return
	}
	panic("Unknown field index")
}

func (r *Arrays) NullField(i int) {
	switch i {
	}
	panic("Not a nullable field index")
}

func (_ *Arrays) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ *Arrays) AppendArray() types.Field         { panic("Unsupported operation") }
func (_ *Arrays) Finalize()                        {}

func (_ *Arrays) AvroCRC64Fingerprint() []byte {
	return []byte(ArraysAvroCRC64Fingerprint)
}
