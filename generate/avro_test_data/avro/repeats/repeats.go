// Code generated by github.com/actgardner/gogen-avro. DO NOT EDIT.
package repeats

import (
	"github.com/actgardner/gogen-avro/v7/compiler"
	"github.com/actgardner/gogen-avro/v7/vm"
	"github.com/actgardner/gogen-avro/v7/vm/types"
	"io"
)

type Repeats struct {
	// The timestamp when this avro data is written. Useful for identifying the newest row of data sharing keys.
	AvroWriteTime int64 `json:"AvroWriteTime"`
	// This is set to true when the Avro data is recording a delete in the source data.
	AvroDeleted bool `json:"AvroDeleted"`

	Height *UnionNullLong `json:"height"`

	SomeDateObj *UnionNullSomeDateObj_record `json:"someDateObj"`

	Type string `json:"type"`

	Visible bool `json:"visible"`

	Width *UnionNullDouble `json:"width"`
}

const RepeatsAvroCRC64Fingerprint = "\xb4\xd3\x17r\x13\xfe\x91\x8c"

func NewRepeats() *Repeats {
	return &Repeats{}
}

func DeserializeRepeats(r io.Reader) (*Repeats, error) {
	t := NewRepeats()
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

func DeserializeRepeatsFromSchema(r io.Reader, schema string) (*Repeats, error) {
	t := NewRepeats()

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

func writeRepeats(r *Repeats, w io.Writer) error {
	var err error
	err = vm.WriteLong(r.AvroWriteTime, w)
	if err != nil {
		return err
	}
	err = vm.WriteBool(r.AvroDeleted, w)
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

func (r *Repeats) Serialize(w io.Writer) error {
	return writeRepeats(r, w)
}

func (r *Repeats) Schema() string {
	return "{\"fields\":[{\"doc\":\"The timestamp when this avro data is written. Useful for identifying the newest row of data sharing keys.\",\"logicalType\":\"timestamp-millis\",\"name\":\"AvroWriteTime\",\"type\":\"long\"},{\"default\":false,\"doc\":\"This is set to true when the Avro data is recording a delete in the source data.\",\"name\":\"AvroDeleted\",\"type\":\"boolean\"},{\"name\":\"height\",\"type\":[\"null\",\"long\"]},{\"name\":\"someDateObj\",\"type\":[\"null\",{\"fields\":[{\"name\":\"type\",\"namespace\":\"someDateObj\",\"type\":\"string\"},{\"default\":false,\"name\":\"visible\",\"namespace\":\"someDateObj\",\"type\":\"boolean\"}],\"name\":\"someDateObj_record\",\"namespace\":\"someDateObj\",\"type\":\"record\"}]},{\"name\":\"type\",\"type\":\"string\"},{\"default\":false,\"name\":\"visible\",\"type\":\"boolean\"},{\"name\":\"width\",\"type\":[\"null\",\"double\"]}],\"name\":\"Repeats\",\"type\":\"record\"}"
}

func (r *Repeats) SchemaName() string {
	return "Repeats"
}

func (_ *Repeats) SetBoolean(v bool)    { panic("Unsupported operation") }
func (_ *Repeats) SetInt(v int32)       { panic("Unsupported operation") }
func (_ *Repeats) SetLong(v int64)      { panic("Unsupported operation") }
func (_ *Repeats) SetFloat(v float32)   { panic("Unsupported operation") }
func (_ *Repeats) SetDouble(v float64)  { panic("Unsupported operation") }
func (_ *Repeats) SetBytes(v []byte)    { panic("Unsupported operation") }
func (_ *Repeats) SetString(v string)   { panic("Unsupported operation") }
func (_ *Repeats) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r *Repeats) Get(i int) types.Field {
	switch i {
	case 0:
		return &types.Long{Target: &r.AvroWriteTime}
	case 1:
		return &types.Boolean{Target: &r.AvroDeleted}
	case 2:
		r.Height = NewUnionNullLong()

		return r.Height
	case 3:
		r.SomeDateObj = NewUnionNullSomeDateObj_record()

		return r.SomeDateObj
	case 4:
		return &types.String{Target: &r.Type}
	case 5:
		return &types.Boolean{Target: &r.Visible}
	case 6:
		r.Width = NewUnionNullDouble()

		return r.Width
	}
	panic("Unknown field index")
}

func (r *Repeats) SetDefault(i int) {
	switch i {
	case 1:
		r.AvroDeleted = false
		return
	case 5:
		r.Visible = false
		return
	}
	panic("Unknown field index")
}

func (r *Repeats) NullField(i int) {
	switch i {
	case 2:
		r.Height = nil
		return
	case 3:
		r.SomeDateObj = nil
		return
	case 6:
		r.Width = nil
		return
	}
	panic("Not a nullable field index")
}

func (_ *Repeats) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ *Repeats) AppendArray() types.Field         { panic("Unsupported operation") }
func (_ *Repeats) Finalize()                        {}

func (_ *Repeats) AvroCRC64Fingerprint() []byte {
	return []byte(RepeatsAvroCRC64Fingerprint)
}