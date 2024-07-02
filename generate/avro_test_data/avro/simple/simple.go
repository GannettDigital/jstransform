// Code generated by github.com/actgardner/gogen-avro. DO NOT EDIT.
package simple

import (
	"github.com/actgardner/gogen-avro/v7/compiler"
	"github.com/actgardner/gogen-avro/v7/vm"
	"github.com/actgardner/gogen-avro/v7/vm/types"
	"io"
)

type Simple struct {
	// The timestamp when this avro data is written. Useful for identifying the newest row of data sharing keys.
	AvroWriteTime int64 `json:"AvroWriteTime"`
	// This is set to true when the Avro data is recording a delete in the source data.
	AvroDeleted bool `json:"AvroDeleted"`

	AMap map[string]string `json:"aMap"`

	Contributors []*Contributors_record `json:"contributors"`

	Height *UnionNullLong `json:"height"`

	SomeDateObj *UnionNullSomeDateObj_record `json:"someDateObj"`

	Type string `json:"type"`

	Visible bool `json:"visible"`

	Width *UnionNullDouble `json:"width"`
}

const SimpleAvroCRC64Fingerprint = "\xdc\xfc\xfd\t[\xb9\x98\x01"

func NewSimple() *Simple {
	return &Simple{}
}

func DeserializeSimple(r io.Reader) (*Simple, error) {
	t := NewSimple()
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

func DeserializeSimpleFromSchema(r io.Reader, schema string) (*Simple, error) {
	t := NewSimple()

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

func writeSimple(r *Simple, w io.Writer) error {
	var err error
	err = vm.WriteLong(r.AvroWriteTime, w)
	if err != nil {
		return err
	}
	err = vm.WriteBool(r.AvroDeleted, w)
	if err != nil {
		return err
	}
	err = writeMapString(r.AMap, w)
	if err != nil {
		return err
	}
	err = writeArrayContributors_record(r.Contributors, w)
	if err != nil {
		return err
	}
	err = writeUnionNullLong(r.Height, w)
	if err != nil {
		return err
	}
	err = writeUnionNullSomeDateObj_record(r.SomeDateObj, w)
	if err != nil {
		return err
	}
	err = vm.WriteString(r.Type, w)
	if err != nil {
		return err
	}
	err = vm.WriteBool(r.Visible, w)
	if err != nil {
		return err
	}
	err = writeUnionNullDouble(r.Width, w)
	if err != nil {
		return err
	}
	return err
}

func (r *Simple) Serialize(w io.Writer) error {
	return writeSimple(r, w)
}

func (r *Simple) Schema() string {
	return "{\"fields\":[{\"doc\":\"The timestamp when this avro data is written. Useful for identifying the newest row of data sharing keys.\",\"logicalType\":\"timestamp-millis\",\"name\":\"AvroWriteTime\",\"type\":\"long\"},{\"default\":false,\"doc\":\"This is set to true when the Avro data is recording a delete in the source data.\",\"name\":\"AvroDeleted\",\"type\":\"boolean\"},{\"name\":\"aMap\",\"type\":{\"type\":\"map\",\"values\":\"string\"}},{\"name\":\"contributors\",\"type\":{\"items\":{\"fields\":[{\"name\":\"contributorId\",\"namespace\":\"contributors\",\"type\":[\"null\",\"string\"]},{\"name\":\"id\",\"namespace\":\"contributors\",\"type\":\"string\"},{\"name\":\"name\",\"namespace\":\"contributors\",\"type\":\"string\"}],\"name\":\"contributors_record\",\"namespace\":\"contributors\",\"type\":\"record\"},\"type\":\"array\"}},{\"name\":\"height\",\"type\":[\"null\",\"long\"]},{\"name\":\"someDateObj\",\"type\":[\"null\",{\"fields\":[{\"name\":\"dates\",\"namespace\":\"someDateObj\",\"type\":{\"items\":{\"logicalType\":\"timestamp-millis\",\"type\":\"long\"},\"type\":\"array\"}}],\"name\":\"someDateObj_record\",\"namespace\":\"someDateObj\",\"type\":\"record\"}]},{\"name\":\"type\",\"type\":\"string\"},{\"default\":false,\"name\":\"visible\",\"type\":\"boolean\"},{\"name\":\"width\",\"type\":[\"null\",\"double\"]}],\"name\":\"Simple\",\"type\":\"record\"}"
}

func (r *Simple) SchemaName() string {
	return "Simple"
}

func (_ *Simple) SetBoolean(v bool)    { panic("Unsupported operation") }
func (_ *Simple) SetInt(v int32)       { panic("Unsupported operation") }
func (_ *Simple) SetLong(v int64)      { panic("Unsupported operation") }
func (_ *Simple) SetFloat(v float32)   { panic("Unsupported operation") }
func (_ *Simple) SetDouble(v float64)  { panic("Unsupported operation") }
func (_ *Simple) SetBytes(v []byte)    { panic("Unsupported operation") }
func (_ *Simple) SetString(v string)   { panic("Unsupported operation") }
func (_ *Simple) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r *Simple) Get(i int) types.Field {
	switch i {
	case 0:
		return &types.Long{Target: &r.AvroWriteTime}
	case 1:
		return &types.Boolean{Target: &r.AvroDeleted}
	case 2:
		r.AMap = make(map[string]string)

		return &MapStringWrapper{Target: &r.AMap}
	case 3:
		r.Contributors = make([]*Contributors_record, 0)

		return &ArrayContributors_recordWrapper{Target: &r.Contributors}
	case 4:
		r.Height = NewUnionNullLong()

		return r.Height
	case 5:
		r.SomeDateObj = NewUnionNullSomeDateObj_record()

		return r.SomeDateObj
	case 6:
		return &types.String{Target: &r.Type}
	case 7:
		return &types.Boolean{Target: &r.Visible}
	case 8:
		r.Width = NewUnionNullDouble()

		return r.Width
	}
	panic("Unknown field index")
}

func (r *Simple) SetDefault(i int) {
	switch i {
	case 1:
		r.AvroDeleted = false
		return
	case 7:
		r.Visible = false
		return
	}
	panic("Unknown field index")
}

func (r *Simple) NullField(i int) {
	switch i {
	case 4:
		r.Height = nil
		return
	case 5:
		r.SomeDateObj = nil
		return
	case 8:
		r.Width = nil
		return
	}
	panic("Not a nullable field index")
}

func (_ *Simple) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ *Simple) AppendArray() types.Field         { panic("Unsupported operation") }
func (_ *Simple) Finalize()                        {}

func (_ *Simple) AvroCRC64Fingerprint() []byte {
	return []byte(SimpleAvroCRC64Fingerprint)
}
