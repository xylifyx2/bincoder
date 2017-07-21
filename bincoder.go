package bincoder

import "bufio"
import "encoding/binary"

// Bincoder that support basic builtin types
type Bincoder interface {
	ui16(f *uint16)
	ui32(f *uint32)
	i(f *int)
	slice(
		length func() int,
		constructor func(int),
		iterate func(int),
	)
	string(f *string)
}

// BinReader holds a bufio.Reader that is the source of unmarshalling
type BinReader struct {
	source *bufio.Reader
}

// BinWriter holds a bufio.Writer that is the target of marshalling
type BinWriter struct {
	target *bufio.Writer
}

func (coder *BinReader) ui16(f *uint16) {
	buf := [2]byte{}
	coder.source.Read(buf[0:2])
	*f = binary.LittleEndian.Uint16(buf[0:2])
}

func (coder *BinWriter) ui16(f *uint16) {
	buf := [2]byte{}
	binary.LittleEndian.PutUint16(buf[:], *f)
	coder.target.Write(buf[:])
}

func (coder *BinReader) ui32(f *uint32) {
	buf := [4]byte{}
	coder.source.Read(buf[0:4])
	*f = binary.LittleEndian.Uint32(buf[0:4])
}

func (coder *BinWriter) ui32(f *uint32) {
	buf := [4]byte{}
	binary.LittleEndian.PutUint32(buf[:], *f)
	coder.target.Write(buf[:])
}

func (coder *BinReader) i(f *int) {
	buf := [4]byte{}
	coder.source.Read(buf[0:4])
	*f = int(binary.LittleEndian.Uint32(buf[0:4]))
}

func (coder *BinWriter) i(f *int) {
	buf := [4]byte{}
	binary.LittleEndian.PutUint32(buf[:], uint32(*f))
	coder.target.Write(buf[:])
}

// reads length [4]byte, entries [length*sizeof(E)]byte
func (coder *BinReader) slice(
	length func() int, constructor func(int), iterator func(int)) {
	var size int
	coder.i(&size)
	constructor(size)
	for i := 0; i < size; i++ {
		iterator(i)
	}
}

// writes length [4]byte, entries [length*sizeof(E)]byte
func (coder *BinWriter) slice(length func() int, constructor func(int),
	iterator func(int)) {
	size := length()
	coder.i(&size) // writes the size of the slice
	for i := 0; i < size; i++ {
		iterator(i)
	}
}

// reads size [4]byte, content [size]byte from source
func (coder *BinReader) string(f *string) {
	var size int
	coder.i(&size)
	c := make([]byte, size)
	coder.source.Read(c)
	*f = string(c)
}

// writes size [4]byte, content [size]byte to target
func (coder *BinWriter) string(f *string) {
	c := []byte(*f)
	size := len(c)
	coder.i(&size)
	coder.target.Write(c)
}
