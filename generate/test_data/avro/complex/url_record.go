// Code generated by github.com/actgardner/gogen-avro. DO NOT EDIT.

package complex

import (
	"github.com/actgardner/gogen-avro/compiler"
	"github.com/actgardner/gogen-avro/container"
	"github.com/actgardner/gogen-avro/vm"
	"github.com/actgardner/gogen-avro/vm/types"
	"io"
)

type URL_record struct {
	Absolute string
	Meta     *Meta_record
	Publish  string
}

func NewURL_recordWriter(writer io.Writer, codec container.Codec, recordsPerBlock int64) (*container.Writer, error) {
	str := &URL_record{}
	return container.NewWriter(writer, codec, recordsPerBlock, str.Schema())
}

func DeserializeURL_record(r io.Reader) (*URL_record, error) {
	t := NewURL_record()

	deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
	if err != nil {
		return nil, err
	}

	err = vm.Eval(r, deser, t)
	return t, err
}

func NewURL_record() *URL_record {
	return &URL_record{}
}

func (r *URL_record) Schema() string {
	return "{\"fields\":[{\"name\":\"absolute\",\"namespace\":\"URL\",\"type\":\"string\"},{\"default\":{},\"name\":\"meta\",\"namespace\":\"URL\",\"type\":{\"fields\":[{\"name\":\"description\",\"namespace\":\"URL.meta\",\"type\":\"string\"},{\"name\":\"siteName\",\"namespace\":\"URL.meta\",\"type\":\"string\"}],\"name\":\"meta_record\",\"namespace\":\"URL.meta\",\"type\":\"record\"}},{\"name\":\"publish\",\"namespace\":\"URL\",\"type\":\"string\"}],\"name\":\"URL_record\",\"namespace\":\"URL\",\"type\":\"record\"}"
}

func (r *URL_record) SchemaName() string {
	return "URL.URL_record"
}

func (r *URL_record) Serialize(w io.Writer) error {
	return writeURL_record(r, w)
}

func (_ *URL_record) SetBoolean(v bool)    { panic("Unsupported operation") }
func (_ *URL_record) SetInt(v int32)       { panic("Unsupported operation") }
func (_ *URL_record) SetLong(v int64)      { panic("Unsupported operation") }
func (_ *URL_record) SetFloat(v float32)   { panic("Unsupported operation") }
func (_ *URL_record) SetDouble(v float64)  { panic("Unsupported operation") }
func (_ *URL_record) SetBytes(v []byte)    { panic("Unsupported operation") }
func (_ *URL_record) SetString(v string)   { panic("Unsupported operation") }
func (_ *URL_record) SetUnionElem(v int64) { panic("Unsupported operation") }
func (r *URL_record) Get(i int) types.Field {
	switch i {
	case 0:
		return (*types.String)(&r.Absolute)
	case 1:
		r.Meta = NewMeta_record()
		return r.Meta
	case 2:
		return (*types.String)(&r.Publish)

	}
	panic("Unknown field index")
}
func (r *URL_record) SetDefault(i int) {
	switch i {
	case 1:

		return

	}
	panic("Unknown field index")
}
func (_ *URL_record) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ *URL_record) AppendArray() types.Field         { panic("Unsupported operation") }
func (_ *URL_record) Finalize()                        {}

type URL_recordReader struct {
	r io.Reader
	p *vm.Program
}

func NewURL_recordReader(r io.Reader) (*URL_recordReader, error) {
	containerReader, err := container.NewReader(r)
	if err != nil {
		return nil, err
	}

	t := NewURL_record()
	deser, err := compiler.CompileSchemaBytes([]byte(containerReader.AvroContainerSchema()), []byte(t.Schema()))
	if err != nil {
		return nil, err
	}

	return &URL_recordReader{
		r: containerReader,
		p: deser,
	}, nil
}

func (r *URL_recordReader) Read() (*URL_record, error) {
	t := NewURL_record()
	err := vm.Eval(r.r, r.p, t)
	return t, err
}
