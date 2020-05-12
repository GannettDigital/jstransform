// Code generated by github.com/actgardner/gogen-avro. DO NOT EDIT.
package arrays

import (
	"github.com/actgardner/gogen-avro/v7/compiler"
	"github.com/actgardner/gogen-avro/v7/vm"
	"github.com/actgardner/gogen-avro/v7/vm/types"
	"io"
)

type Info_record struct {
	Name string `json:"name"`

	Age int64 `json:"age"`
}

const Info_recordAvroCRC64Fingerprint = "\x1f\xf5\xe3Eʞ\xf31"

func NewInfo_record() *Info_record {
	return &Info_record{}
}

func DeserializeInfo_record(r io.Reader) (*Info_record, error) {
	t := NewInfo_record()
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

func DeserializeInfo_recordFromSchema(r io.Reader, schema string) (*Info_record, error) {
	t := NewInfo_record()

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

func writeInfo_record(r *Info_record, w io.Writer) error {
	var err error
	err = vm.WriteString(r.Name, w)
	if err != nil {
		return err
	}
	err = vm.WriteLong(r.Age, w)
	if err != nil {
		return err
	}
	return err
}

func (r *Info_record) Serialize(w io.Writer) error {
	return writeInfo_record(r, w)
}

func (r *Info_record) Schema() string {
	return "{\"fields\":[{\"name\":\"name\",\"namespace\":\"parents.info\",\"type\":\"string\"},{\"name\":\"age\",\"namespace\":\"parents.info\",\"type\":\"long\"}],\"name\":\"parents.info.info_record\",\"type\":\"record\"}"
}

func (r *Info_record) SchemaName() string {
	return "parents.info.info_record"
}

func (_ *Info_record) SetBoolean(v bool)    { panic("Unsupported operation") }
func (_ *Info_record) SetInt(v int32)       { panic("Unsupported operation") }
func (_ *Info_record) SetLong(v int64)      { panic("Unsupported operation") }
func (_ *Info_record) SetFloat(v float32)   { panic("Unsupported operation") }
func (_ *Info_record) SetDouble(v float64)  { panic("Unsupported operation") }
func (_ *Info_record) SetBytes(v []byte)    { panic("Unsupported operation") }
func (_ *Info_record) SetString(v string)   { panic("Unsupported operation") }
func (_ *Info_record) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r *Info_record) Get(i int) types.Field {
	switch i {
	case 0:
		return &types.String{Target: &r.Name}
	case 1:
		return &types.Long{Target: &r.Age}
	}
	panic("Unknown field index")
}

func (r *Info_record) SetDefault(i int) {
	switch i {
	}
	panic("Unknown field index")
}

func (r *Info_record) NullField(i int) {
	switch i {
	}
	panic("Not a nullable field index")
}

func (_ *Info_record) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ *Info_record) AppendArray() types.Field         { panic("Unsupported operation") }
func (_ *Info_record) Finalize()                        {}

func (_ *Info_record) AvroCRC64Fingerprint() []byte {
	return []byte(Info_recordAvroCRC64Fingerprint)
}
