package bincoder

import "bufio"
import "encoding/binary"

// Bincoder that support basic builtin types
type Bincoder interface {
	UI16(f *uint16)
	UI32(f *uint32)
	Int(f *int)
	Slice(
		length func() int,
		constructor func(int),
		iterate func(int),
	)
	String(f *string)
}

// BinReader holds a bufio.Reader that is the source of unmarshalling
type BinReader struct {
	source *bufio.Reader
}

// BinWriter holds a bufio.Writer that is the target of marshalling
type BinWriter struct {
	target *bufio.Writer
}

// UI16 uint16 reader
func (coder *BinReader) UI16(f *uint16) {
	buf := [2]byte{}
	coder.source.Read(buf[0:2])
	*f = binary.LittleEndian.Uint16(buf[0:2])
}

// UI16 uint16 writer
func (coder *BinWriter) UI16(f *uint16) {
	buf := [2]byte{}
	binary.LittleEndian.PutUint16(buf[:], *f)
	coder.target.Write(buf[:])
}

// UI32 uint32 reader
func (coder *BinReader) UI32(f *uint32) {
	buf := [4]byte{}
	coder.source.Read(buf[0:4])
	*f = binary.LittleEndian.Uint32(buf[0:4])
}

// UI32 uint32 writer
func (coder *BinWriter) UI32(f *uint32) {
	buf := [4]byte{}
	binary.LittleEndian.PutUint32(buf[:], *f)
	coder.target.Write(buf[:])
}

// Int int reader
func (coder *BinReader) Int(f *int) {
	buf := [4]byte{}
	coder.source.Read(buf[0:4])
	*f = int(binary.LittleEndian.Uint32(buf[0:4]))
}

// Int int writer
func (coder *BinWriter) Int(f *int) {
	buf := [4]byte{}
	binary.LittleEndian.PutUint32(buf[:], uint32(*f))
	coder.target.Write(buf[:])
}

// Slice reads length [4]byte, entries [length*sizeof(E)]byte
func (coder *BinReader) Slice(
	length func() int, constructor func(int), iterator func(int)) {
	var size int
	coder.Int(&size)
	constructor(size)
	for i := 0; i < size; i++ {
		iterator(i)
	}
}

// Slice writes length [4]byte, entries [length*sizeof(E)]byte
func (coder *BinWriter) Slice(length func() int, constructor func(int),
	iterator func(int)) {
	size := length()
	coder.Int(&size) // writes the size of the slice
	for i := 0; i < size; i++ {
		iterator(i)
	}
}

// reads size [4]byte, content [size]byte from source
func (coder *BinReader) String(f *string) {
	var size int
	coder.Int(&size)
	c := make([]byte, size)
	coder.source.Read(c)
	*f = string(c)
}

// writes size [4]byte, content [size]byte to target
func (coder *BinWriter) String(f *string) {
	c := []byte(*f)
	size := len(c)
	coder.Int(&size)
	coder.target.Write(c)
}
