// Code generated by github.com/actgardner/gogen-avro. DO NOT EDIT.
package nested

import (
	"github.com/actgardner/gogen-avro/v7/compiler"
	"github.com/actgardner/gogen-avro/v7/vm"
	"github.com/actgardner/gogen-avro/v7/vm/types"
	"io"
)

type FactCheckClaims_record struct {
	AppearanceURLs []*AppearanceURLs_record `json:"appearanceURLs"`

	Author *UnionNullString `json:"author"`

	Claim *UnionNullString `json:"claim"`

	Date *UnionNullString `json:"date"`

	Rating *UnionNullString `json:"rating"`
}

const FactCheckClaims_recordAvroCRC64Fingerprint = "\xa2YuPD\x84\xc4,"

func NewFactCheckClaims_record() *FactCheckClaims_record {
	return &FactCheckClaims_record{}
}

func DeserializeFactCheckClaims_record(r io.Reader) (*FactCheckClaims_record, error) {
	t := NewFactCheckClaims_record()
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

func DeserializeFactCheckClaims_recordFromSchema(r io.Reader, schema string) (*FactCheckClaims_record, error) {
	t := NewFactCheckClaims_record()

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

func writeFactCheckClaims_record(r *FactCheckClaims_record, w io.Writer) error {
	var err error
	err = writeArrayAppearanceURLs_record(r.AppearanceURLs, w)
	if err != nil {
		return err
	}
	err = writeUnionNullString(r.Author, w)
	if err != nil {
		return err
	}
	err = writeUnionNullString(r.Claim, w)
	if err != nil {
		return err
	}
	err = writeUnionNullString(r.Date, w)
	if err != nil {
		return err
	}
	err = writeUnionNullString(r.Rating, w)
	if err != nil {
		return err
	}
	return err
}

func (r *FactCheckClaims_record) Serialize(w io.Writer) error {
	return writeFactCheckClaims_record(r, w)
}

func (r *FactCheckClaims_record) Schema() string {
	return "{\"fields\":[{\"name\":\"appearanceURLs\",\"namespace\":\"factCheckClaims\",\"type\":{\"items\":{\"fields\":[{\"default\":false,\"name\":\"original\",\"namespace\":\"factCheckClaims.appearanceURLs\",\"type\":\"boolean\"},{\"name\":\"url\",\"namespace\":\"factCheckClaims.appearanceURLs\",\"type\":\"string\"}],\"name\":\"appearanceURLs_record\",\"namespace\":\"factCheckClaims.appearanceURLs\",\"type\":\"record\"},\"type\":\"array\"}},{\"name\":\"author\",\"namespace\":\"factCheckClaims\",\"type\":[\"null\",\"string\"]},{\"name\":\"claim\",\"namespace\":\"factCheckClaims\",\"type\":[\"null\",\"string\"]},{\"name\":\"date\",\"namespace\":\"factCheckClaims\",\"type\":[\"null\",\"string\"]},{\"name\":\"rating\",\"namespace\":\"factCheckClaims\",\"type\":[\"null\",\"string\"]}],\"name\":\"factCheckClaims.factCheckClaims_record\",\"type\":\"record\"}"
}

func (r *FactCheckClaims_record) SchemaName() string {
	return "factCheckClaims.factCheckClaims_record"
}

func (_ *FactCheckClaims_record) SetBoolean(v bool)    { panic("Unsupported operation") }
func (_ *FactCheckClaims_record) SetInt(v int32)       { panic("Unsupported operation") }
func (_ *FactCheckClaims_record) SetLong(v int64)      { panic("Unsupported operation") }
func (_ *FactCheckClaims_record) SetFloat(v float32)   { panic("Unsupported operation") }
func (_ *FactCheckClaims_record) SetDouble(v float64)  { panic("Unsupported operation") }
func (_ *FactCheckClaims_record) SetBytes(v []byte)    { panic("Unsupported operation") }
func (_ *FactCheckClaims_record) SetString(v string)   { panic("Unsupported operation") }
func (_ *FactCheckClaims_record) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r *FactCheckClaims_record) Get(i int) types.Field {
	switch i {
	case 0:
		r.AppearanceURLs = make([]*AppearanceURLs_record, 0)

		return &ArrayAppearanceURLs_recordWrapper{Target: &r.AppearanceURLs}
	case 1:
		r.Author = NewUnionNullString()

		return r.Author
	case 2:
		r.Claim = NewUnionNullString()

		return r.Claim
	case 3:
		r.Date = NewUnionNullString()

		return r.Date
	case 4:
		r.Rating = NewUnionNullString()

		return r.Rating
	}
	panic("Unknown field index")
}

func (r *FactCheckClaims_record) SetDefault(i int) {
	switch i {
	}
	panic("Unknown field index")
}

func (r *FactCheckClaims_record) NullField(i int) {
	switch i {
	case 1:
		r.Author = nil
		return
	case 2:
		r.Claim = nil
		return
	case 3:
		r.Date = nil
		return
	case 4:
		r.Rating = nil
		return
	}
	panic("Not a nullable field index")
}

func (_ *FactCheckClaims_record) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ *FactCheckClaims_record) AppendArray() types.Field         { panic("Unsupported operation") }
func (_ *FactCheckClaims_record) Finalize()                        {}

func (_ *FactCheckClaims_record) AvroCRC64Fingerprint() []byte {
	return []byte(FactCheckClaims_recordAvroCRC64Fingerprint)
}
