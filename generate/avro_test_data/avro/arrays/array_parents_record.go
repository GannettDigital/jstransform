// Code generated by github.com/actgardner/gogen-avro. DO NOT EDIT.
package arrays

import (
	"io"

	"github.com/actgardner/gogen-avro/v7/vm"
	"github.com/actgardner/gogen-avro/v7/vm/types"
)

func writeArrayParents_record(r []*Parents_record, w io.Writer) error {
	err := vm.WriteLong(int64(len(r)), w)
	if err != nil || len(r) == 0 {
		return err
	}
	for _, e := range r {
		err = writeParents_record(e, w)
		if err != nil {
			return err
		}
	}
	return vm.WriteLong(0, w)
}

type ArrayParents_recordWrapper struct {
	Target *[]*Parents_record
}

func (_ *ArrayParents_recordWrapper) SetBoolean(v bool)     { panic("Unsupported operation") }
func (_ *ArrayParents_recordWrapper) SetInt(v int32)        { panic("Unsupported operation") }
func (_ *ArrayParents_recordWrapper) SetLong(v int64)       { panic("Unsupported operation") }
func (_ *ArrayParents_recordWrapper) SetFloat(v float32)    { panic("Unsupported operation") }
func (_ *ArrayParents_recordWrapper) SetDouble(v float64)   { panic("Unsupported operation") }
func (_ *ArrayParents_recordWrapper) SetBytes(v []byte)     { panic("Unsupported operation") }
func (_ *ArrayParents_recordWrapper) SetString(v string)    { panic("Unsupported operation") }
func (_ *ArrayParents_recordWrapper) SetUnionElem(v int64)  { panic("Unsupported operation") }
func (_ *ArrayParents_recordWrapper) Get(i int) types.Field { panic("Unsupported operation") }
func (_ *ArrayParents_recordWrapper) AppendMap(key string) types.Field {
	panic("Unsupported operation")
}
func (_ *ArrayParents_recordWrapper) Finalize()        {}
func (_ *ArrayParents_recordWrapper) SetDefault(i int) { panic("Unsupported operation") }
func (r *ArrayParents_recordWrapper) NullField(i int) {
	panic("Unsupported operation")
}

func (r *ArrayParents_recordWrapper) AppendArray() types.Field {
	var v *Parents_record
	v = NewParents_record()

	*r.Target = append(*r.Target, v)

	return (*r.Target)[len(*r.Target)-1]
}
