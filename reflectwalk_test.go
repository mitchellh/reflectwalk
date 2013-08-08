package reflectwalk

import (
	"reflect"
	"testing"
)

type TestPrimitiveWalker struct {
	Value reflect.Value
}

func (t *TestPrimitiveWalker) Primitive(v reflect.Value) error {
	t.Value = v
	return nil
}

type TestSliceWalker struct {
	Count int
}

func (t *TestSliceWalker) SliceElem(int, reflect.Value) error {
	t.Count++
	return nil
}

type TestStructWalker struct {
	Fields []string
}

func (t *TestStructWalker) StructField(sf reflect.StructField, v reflect.Value) error {
	if t.Fields == nil {
		t.Fields = make([]string, 0, 1)
	}

	t.Fields = append(t.Fields, sf.Name)
	return nil
}

func TestTestStructs(t *testing.T) {
	var raw interface{}
	raw = new(TestPrimitiveWalker)
	if _, ok := raw.(PrimitiveWalker); !ok {
		t.Fatal("PrimitiveWalker is bad")
	}

	raw = new(TestSliceWalker)
	if _, ok := raw.(SliceWalker); !ok {
		t.Fatal("SliceWalker is bad")
	}

	raw = new(TestStructWalker)
	if _, ok := raw.(StructWalker); !ok {
		t.Fatal("StructWalker is bad")
	}
}

func TestWalk_Basic(t *testing.T) {
	w := new(TestPrimitiveWalker)

	type S struct {
		Foo string
	}

	data := &S{
		Foo: "foo",
	}

	err := Walk(data, w)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if w.Value.Kind() != reflect.String {
		t.Fatalf("bad: %#v", w.Value)
	}
}

func TestWalk_Slice(t *testing.T) {
	w := new(TestSliceWalker)

	type S struct {
		Foo []string
	}

	data := &S{
		Foo: []string{"a", "b", "c"},
	}

	err := Walk(data, w)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if w.Count != 3 {
		t.Fatalf("Bad count: %d", w.Count)
	}
}

func TestWalk_Struct(t *testing.T) {
	w := new(TestStructWalker)

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
	if !reflect.DeepEqual(w.Fields, expected) {
		t.Fatalf("bad: %#v", w.Fields)
	}
}
