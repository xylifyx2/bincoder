package bincoder

import (
	"bufio"
	"reflect"
	"testing"
)

// A is a type
type foo struct {
	a uint16
	b uint32
}

type bar struct {
	x       uint32
	foo     foo
	y       uint16
	foos    []foo
	name    string
	z       uint64
	integer int
	data    [30]byte
}

func (f *bar) encode(w testCoder) {
	w.UI16(&f.y)
	w.foo(&f.foo)
	w.UI32(&f.x)
	w.fooSlice(&f.foos)
	w.Bytes(
		func() int { return 30 },
		func() []byte { return f.data[:] },
		func(data []byte) {
			buf := [30]byte{}
			copy(buf[:], data)
			f.data = buf
		})
	w.String(&f.name)
	w.UI64(&f.z)
	w.Int(&f.integer)
}
func (wire *BinReader) bar(f *bar) {
	f.encode(wire)
}

func (wire *BinWriter) bar(f *bar) {
	f.encode(wire)
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

func (f *foo) encode(w Bincoder) {
	w.UI16(&f.a)
	w.UI32(&f.b)
}

func TestFoo_marshall(t *testing.T) {
	o := foo{
		a: 10,
		b: 20,
	}
	got := Marshall(func(writer *bufio.Writer) {
		w := BinWriter{Target: writer}
		w.foo(&o)
	})
	want := []byte{10, 0, 20, 0, 0, 0}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("%q. got %v, want %v", "version", got, want)
	}
}

func TestFoo_unmarshall(t *testing.T) {

	var got foo

	marshalled := []byte{87, 0, 42, 0, 0, 0}
	Unmarshall(func(reader *bufio.Reader) {
		r := BinReader{Source: reader}
		r.foo(&got)
	}, marshalled)

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
		x: 42, y: 87, z: 1024 * 1024 * 1024 * 1024, integer: -87,

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
		data: [30]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 128, 129, 130},
		name: "Wire Marshall",
	}

	marshalled := Marshall(func(writer *bufio.Writer) {
		w := BinWriter{Target: writer}
		w.bar(&want)
	})

	var got bar

	Unmarshall(func(reader *bufio.Reader) {
		r := BinReader{Source: reader}
		r.bar(&got)
	}, marshalled)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("%q. got %v, want %v", "version", got, want)
	}
}
