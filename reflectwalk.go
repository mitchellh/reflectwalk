package reflectwalk

import (
	"reflect"
)

// StructWalker is an interface that has methods that are called for
// structs when a Walk is done.
type StructWalker interface {
	StructField(reflect.StructField, reflect.Value) error
}

func Walk(data, walker interface{}) error {
	v := reflect.Indirect(reflect.ValueOf(data))
	return walk(v, walker)
}

func walk(v reflect.Value, w interface{}) error {
	switch v.Kind() {
	case reflect.Struct:
		return walkStruct(v, w)
	default:
		return nil
	}
}

func walkStruct(v reflect.Value, w interface{}) (err error) {
	vt := v.Type()
	for i := 0; i < vt.NumField(); i++ {
		sf := vt.Field(i)
		f := v.FieldByIndex([]int{i})

		if sw, ok := w.(StructWalker); ok {
			err = sw.StructField(sf, f)
			if err != nil {
				return
			}
		}

		err = walk(f, w)
		if err != nil {
			return
		}
	}

	return nil
}
