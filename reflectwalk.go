// reflectwalk is a package that allows you to "walk" complex structures
// similar to how you may "walk" a filesystem: visiting every element one
// by one and calling callback functions allowing you to handle and manipulate
// those elements.
package reflectwalk

import (
	"reflect"
)

// PrimitiveWalker implementations are able to handle primitive values
// within complex structures. Primitive values are numbers, strings,
// booleans, funcs, chans.
//
// These primitive values are often members of more complex
// structures (slices, maps, etc.) that are walkable by other interfaces.
type PrimitiveWalker interface {
	Primitive(reflect.Value) error
}

// MapWalker implementations are able to handle individual elements
// found within a map structure.
type MapWalker interface {
	MapElem(k, v reflect.Value) error
}

// SliceWalker implementations are able to handle slice elements found
// within complex structures.
type SliceWalker interface {
	SliceElem(int, reflect.Value) error
}

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
	k := v.Kind()
	if k >= reflect.Int && k <= reflect.Complex128 {
		k = reflect.Int
	}

	switch k {
	// Primitives
	case reflect.Bool:
		fallthrough
	case reflect.Chan:
		fallthrough
	case reflect.Func:
		fallthrough
	case reflect.Int:
		fallthrough
	case reflect.String:
		return walkPrimitive(v, w)
	case reflect.Map:
		return walkMap(v, w)
	case reflect.Slice:
		return walkSlice(v, w)
	case reflect.Struct:
		return walkStruct(v, w)
	default:
		panic("unsupported type: " + k.String())
	}
}

func walkMap(v reflect.Value, w interface{}) error {
	for _, k := range v.MapKeys() {
		kv := v.MapIndex(k)

		if mw, ok := w.(MapWalker); ok {
			if err := mw.MapElem(k, kv); err != nil {
				return err
			}
		}

		if err := walk(k, w); err != nil {
			return err
		}

		if err := walk(kv, w); err != nil {
			return err
		}
	}

	return nil
}

func walkPrimitive(v reflect.Value, w interface{}) error {
	if pw, ok := w.(PrimitiveWalker); ok {
		return pw.Primitive(v)
	}

	return nil
}

func walkSlice(v reflect.Value, w interface{}) (err error) {
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)

		if sw, ok := w.(SliceWalker); ok {
			if err := sw.SliceElem(i, elem); err != nil {
				return err
			}
		}

		if err := walk(elem, w); err != nil {
			return err
		}
	}

	return nil
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
