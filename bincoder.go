package bincoder

import "encoding/binary"

// Bincoder that support basic builtin types
type Bincoder interface {
	UI16(f *uint16)
	UI32(f *uint32)
	UI64(f *uint64)
	Int(f *int)
	String(f *string)
	Slice(
		length func() int,
		constructor func(int),
		iterate func(int),
	)
	Bytes(func() int, func() []byte, func([]byte))
}

// Reader interface
type Reader interface {
	Read(p []byte) (n int, err error)
}

// Writer interface, required functions
type Writer interface {
	Write(p []byte) (n int, err error)
	Flush() error
}

// BinReader holds a bufio.Reader that is the Source of unmarshalling
type BinReader struct {
	err    error
	Source Reader
}

func (coder *BinReader) Read(p []byte) (n int, err error) {
	if coder.err != nil {
		return 0, coder.err
	}
	n, err = coder.Source.Read(p)
	if err != nil {
		coder.err = err
	}
	return n, err
}

func (coder *BinWriter) Write(p []byte) (n int, err error) {
	if coder.err != nil {
		return 0, coder.err
	}
	n, err = coder.Target.Write(p)
	if err != nil {
		coder.err = err
	}
	return n, err
}

// BinWriter holds a bufio.Writer that is the Target of marshalling
type BinWriter struct {
	err    error
	Target Writer
}

// Flush output to io.Writer
func (coder *BinWriter) Flush() {
	coder.Target.Flush()
}

// UI16 uint16 reader
func (coder *BinReader) UI16(f *uint16) {
	buf := [2]byte{}
	coder.Read(buf[0:2])
	*f = binary.LittleEndian.Uint16(buf[0:2])
}

// UI16 uint16 writer
func (coder *BinWriter) UI16(f *uint16) {
	buf := [2]byte{}
	binary.LittleEndian.PutUint16(buf[:], *f)
	coder.Write(buf[:])
}

// UI32 uint32 reader
func (coder *BinReader) UI32(f *uint32) {
	buf := [4]byte{}
	coder.Read(buf[:])
	*f = binary.LittleEndian.Uint32(buf[:])
}

// UI32 uint32 writer
func (coder *BinWriter) UI32(f *uint32) {
	buf := [4]byte{}
	binary.LittleEndian.PutUint32(buf[:], *f)
	coder.Write(buf[:])
}

// UI64 uint64 reader
func (coder *BinReader) UI64(f *uint64) {
	buf := [8]byte{}
	coder.Read(buf[:])
	*f = binary.LittleEndian.Uint64(buf[:])
}

// UI64 uint64 writer
func (coder *BinWriter) UI64(f *uint64) {
	buf := [8]byte{}
	binary.LittleEndian.PutUint64(buf[:], *f)
	coder.Write(buf[:])
}

// Int int reader
func (coder *BinReader) Int(f *int) {
	buf := [8]byte{}
	coder.Read(buf[:])
	*f = int(binary.LittleEndian.Uint64(buf[:]))
}

// Int int writer
func (coder *BinWriter) Int(f *int) {
	buf := [8]byte{}
	binary.LittleEndian.PutUint64(buf[:], uint64(*f))
	coder.Write(buf[:])
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

// reads size [4]byte, content [size]byte from Source
func (coder *BinReader) String(f *string) {
	var size int
	coder.Int(&size)
	c := make([]byte, size)
	coder.Read(c)
	*f = string(c)
}

// writes size [4]byte, content [size]byte to Target
func (coder *BinWriter) String(f *string) {
	c := []byte(*f)
	size := len(c)
	coder.Int(&size)
	coder.Write(c)
}

// Byte Slice
func byteSliceCoder(f *[]byte, coder Bincoder, length int) {
	coder.Bytes(func() int {
		return length
	}, func() []byte { return *f }, func(value []byte) {
		*f = value
	})
}

func (coder *BinReader) ByteSlice(f *[]byte, length int) {
	byteSliceCoder(f, coder, length)
}

func (coder *BinWriter) ByteSlice(f *[]byte, length int) {
	byteSliceCoder(f, coder, length)
}

// Bytes reader
func (coder *BinReader) Bytes(
	length func() int,
	getter func() []byte,
	setter func([]byte)) {
	buf := make([]byte, length())
	coder.Read(buf)
	setter(buf)
}

// Bytes writer
func (coder *BinWriter) Bytes(length func() int,
	getter func() []byte, setter func([]byte)) {
	coder.Write(getter())
}
