// Code generated by github.com/actgardner/gogen-avro. DO NOT EDIT.
package complex

import (
	"github.com/actgardner/gogen-avro/v7/compiler"
	"github.com/actgardner/gogen-avro/v7/vm"
	"github.com/actgardner/gogen-avro/v7/vm/types"
	"io"
)

type Contributors_record struct {
	ContributorId *UnionNullString `json:"contributorId"`

	Id string `json:"id"`

	Name string `json:"name"`
}

const Contributors_recordAvroCRC64Fingerprint = "/-\xd0\f3\x80C\x94"

func NewContributors_record() *Contributors_record {
	return &Contributors_record{}
}

func DeserializeContributors_record(r io.Reader) (*Contributors_record, error) {
	t := NewContributors_record()
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

func DeserializeContributors_recordFromSchema(r io.Reader, schema string) (*Contributors_record, error) {
	t := NewContributors_record()

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

func writeContributors_record(r *Contributors_record, w io.Writer) error {
	var err error
	err = writeUnionNullString(r.ContributorId, w)
	if err != nil {
		return err
	}
	err = vm.WriteString(r.Id, w)
	if err != nil {
		return err
	}
	err = vm.WriteString(r.Name, w)
	if err != nil {
		return err
	}
	return err
}

func (r *Contributors_record) Serialize(w io.Writer) error {
	return writeContributors_record(r, w)
}

func (r *Contributors_record) Schema() string {
	return "{\"fields\":[{\"name\":\"contributorId\",\"namespace\":\"Simple.contributors\",\"type\":[\"null\",\"string\"]},{\"name\":\"id\",\"namespace\":\"Simple.contributors\",\"type\":\"string\"},{\"name\":\"name\",\"namespace\":\"Simple.contributors\",\"type\":\"string\"}],\"name\":\"Simple.contributors.contributors_record\",\"type\":\"record\"}"
}

func (r *Contributors_record) SchemaName() string {
	return "Simple.contributors.contributors_record"
}

func (_ *Contributors_record) SetBoolean(v bool)    { panic("Unsupported operation") }
func (_ *Contributors_record) SetInt(v int32)       { panic("Unsupported operation") }
func (_ *Contributors_record) SetLong(v int64)      { panic("Unsupported operation") }
func (_ *Contributors_record) SetFloat(v float32)   { panic("Unsupported operation") }
func (_ *Contributors_record) SetDouble(v float64)  { panic("Unsupported operation") }
func (_ *Contributors_record) SetBytes(v []byte)    { panic("Unsupported operation") }
func (_ *Contributors_record) SetString(v string)   { panic("Unsupported operation") }
func (_ *Contributors_record) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r *Contributors_record) Get(i int) types.Field {
	switch i {
	case 0:
		r.ContributorId = NewUnionNullString()

		return r.ContributorId
	case 1:
		return &types.String{Target: &r.Id}
	case 2:
		return &types.String{Target: &r.Name}
	}
	panic("Unknown field index")
}

func (r *Contributors_record) SetDefault(i int) {
	switch i {
	}
	panic("Unknown field index")
}

func (r *Contributors_record) NullField(i int) {
	switch i {
	case 0:
		r.ContributorId = nil
		return
	}
	panic("Not a nullable field index")
}

func (_ *Contributors_record) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ *Contributors_record) AppendArray() types.Field         { panic("Unsupported operation") }
func (_ *Contributors_record) Finalize()                        {}

func (_ *Contributors_record) AvroCRC64Fingerprint() []byte {
	return []byte(Contributors_recordAvroCRC64Fingerprint)
}