package bincoder

import (
	"bufio"
	"bytes"
	"reflect"
	"testing"
)

// A is a type
type foo struct {
	a uint16
	b uint32
}

type bar struct {
	x    uint32
	foo  foo
	y    uint16
	foos []foo
	name string
}

type testCoder interface {
	Bincoder
	foo(f *foo)
	bar(f *bar)
	fooSlice(f *[]foo)
}

func (wire *BinReader) foo(f *foo) {
	f.encode(wire)
}

func (wire *BinWriter) foo(f *foo) {
	f.encode(wire)
}

func encodeFooSlice(f *[]foo, w testCoder) {
	w.Slice(
		func() int { return len(*f) },
		func(size int) { *f = make([]foo, size) },
		func(i int) { w.foo(&(*f)[i]) },
	)
}

func (wire *BinReader) fooSlice(f *[]foo) {
	encodeFooSlice(f, wire)
}

func (wire *BinWriter) fooSlice(f *[]foo) {
	encodeFooSlice(f, wire)
}

func (wire *BinReader) bar(f *bar) {
	f.encode(wire)
}

func (wire *BinWriter) bar(f *bar) {
	f.encode(wire)
}

func (f *foo) encode(w Bincoder) {
	w.UI16(&f.a)
	w.UI32(&f.b)
}

func (f *bar) encode(w testCoder) {
	w.UI16(&f.y)
	w.foo(&f.foo)
	w.UI32(&f.x)
	w.fooSlice(&f.foos)
	w.String(&f.name)
}

func TestFoo_marshall(t *testing.T) {
	o := foo{
		a: 10,
		b: 20,
	}

	var b bytes.Buffer
	w := BinWriter{Target: bufio.NewWriter(&b)}
	w.foo(&o)
	want := []byte{10, 0, 20, 0, 0, 0}
	w.Target.Flush()

	got := b.Bytes()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%q. got %v, want %v", "version", got, want)
	}
}

func TestFoo_unmarshall(t *testing.T) {
	var b bytes.Buffer
	b.Write([]byte{87, 0, 42, 0, 0, 0})

	w := BinReader{Source: bufio.NewReader(&b)}
	got := foo{}
	w.foo(&got)

	want := foo{
		a: 87,
		b: 42,
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("%q. got %v, want %v", "version", got, want)
	}
}

func TestBar_marshall(t *testing.T) {
	want := bar{
		x: 42, y: 87,
		foo: foo{
			a: 10,
			b: 20,
		},
		foos: []foo{
			{
				a: 11,
				b: 12,
			}, {
				a: 13,
				b: 14,
			},
		},
		name: "Wire Marshall",
	}
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	w := BinWriter{Target: writer}
	w.bar(&want)
	writer.Flush()

	reader := bufio.NewReader(&b)
	r := BinReader{Source: reader}
	got := bar{}
	r.bar(&got)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("%q. got %v, want %v", "version", got, want)
	}
}
