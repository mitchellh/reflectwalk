package reflectwalk

import (
	"reflect"
	"testing"
)

type TestStructWalker struct {
	StructFieldFunc func(reflect.StructField, reflect.Value) error
}

func (t TestStructWalker) StructField(sf reflect.StructField, v reflect.Value) error {
	return t.StructFieldFunc(sf, v)
}

func TestTestStructs(t *testing.T) {
	var raw interface{}
	raw = TestStructWalker{}
	if _, ok := raw.(StructWalker); !ok {
		t.Fatal("StructWalker is bad")
	}
}

func TestWalk_Struct(t *testing.T) {
	fields := make([]string, 0)
	w := &TestStructWalker{
		StructFieldFunc: func(sf reflect.StructField, v reflect.Value) error {
			fields = append(fields, sf.Name)
			return nil
		},
	}

	type S struct {
		Foo string
		Bar string
	}

	data := &S{
		Foo: "foo",
		Bar: "bar",
	}

	err := Walk(data, w)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := []string{"Foo", "Bar"}
	if !reflect.DeepEqual(fields, expected) {
		t.Fatalf("bad: %#v", fields)
	}
}
