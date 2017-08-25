package bincoder

import (
	"encoding/binary"
	"fmt"
	"log"
)

// Bincoder that support basic builtin types
type Bincoder interface {
	SetError(err error)
	Error() error
	UI16(f *uint16)
	UI32(f *uint32)
	I32(f *int32)
	UI64(f *uint64)
	I64(f *int64)
	Int(f *int)
	VarInt(f *uint64)
	// String coded as length followed by raw byte data
	String(f *string)
	// ByteSlice codes a []byte field
	ByteSlice(f *[]byte, length int)
	// Slice coder that codes any slice type
	Slice(
		length int,
		constructor func(int),
		iterate func(int),
	)
	// Codes raw bytes
	Bytes(int, func() []byte, func([]byte))
}

// SetReader updates the Source of the BinReader
func (coder *BinReader) SetReader(reader Reader) {
	coder.source = reader
}

// SetWriter updates the target of the BinWriter
func (coder *BinWriter) SetWriter(writer Writer) {
	coder.target = writer
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

// CoderBase is a base type of a Bincoder
type CoderBase struct {
	err error
}

// SetError on the reader or writer
func (coder *CoderBase) SetError(err error) {
	if coder.err == nil {
		coder.err = err
		log.Print(fmt.Sprintf("SetError: %v", err))
	}
}

// Error get the error of a reader or writer
func (coder *CoderBase) Error() error {
	return coder.err
}

// BinReader holds a bufio.Reader that is the Source of unmarshalling
type BinReader struct {
	CoderBase
	source Reader
}

// BinWriter holds a bufio.Writer that is the target of marshalling
type BinWriter struct {
	CoderBase
	target Writer
}

func (coder *BinReader) Read(p []byte) (int, error) {
	if coder.err != nil {
		return 0, coder.err
	}
	n := 0
	for n < len(p) {
		c, err := coder.source.Read(p[n:])
		if err != nil {
			coder.SetError(err)
			return n + c, err
		}
		n += c
	}

	return n, nil
}

func (coder *BinWriter) Write(p []byte) (n int, err error) {
	if coder.err != nil {
		return 0, coder.err
	}
	n, err = coder.target.Write(p)
	if n != len(p) {
		coder.SetError(fmt.Errorf("Wrote %d bytes expected %d", n, len(p)))
	} else if err != nil {
		coder.SetError(err)
	}
	return n, err
}

// Flush output to io.Writer
func (coder *BinWriter) Flush() {
	coder.target.Flush()
}

// UI16 uint16 reader
func (coder *BinReader) UI16(f *uint16) {
	buf := [2]byte{}
	_, err := coder.Read(buf[:])
	if err != nil {
		return
	}
	*f = binary.LittleEndian.Uint16(buf[0:2])
}

// UI16 uint16 writer
func (coder *BinWriter) UI16(f *uint16) {
	buf := [2]byte{}
	binary.LittleEndian.PutUint16(buf[:], *f)
	_, err := coder.Write(buf[:])
	if err != nil {
		return
	}
}

// UI32 uint32 reader
func (coder *BinReader) UI32(f *uint32) {
	buf := [4]byte{}
	_, err := coder.Read(buf[:])
	if err != nil {
		return
	}
	*f = binary.LittleEndian.Uint32(buf[:])
}

// UI32 uint32 writer
func (coder *BinWriter) UI32(f *uint32) {
	if coder.err != nil {
		return
	}
	buf := [4]byte{}
	binary.LittleEndian.PutUint32(buf[:], *f)
	coder.Write(buf[:])
}

// UI32 uint32 reader
func (coder *BinReader) I32(f *int32) {
	var buf [4]byte
	_, err := coder.Read(buf[:])
	if err != nil {
		return
	}
	*f = int32(binary.LittleEndian.Uint32(buf[:]))
}

// UI32 uint32 writer
func (coder *BinWriter) I32(f *int32) {
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], uint32(*f))
	coder.Write(buf[:])
}

// UI64 uint64 reader
func (coder *BinReader) UI64(f *uint64) {
	var buf [8]byte
	_, err := coder.Read(buf[:])
	if err != nil {
		return
	}
	*f = binary.LittleEndian.Uint64(buf[:])
}

// UI64 uint64 writer
func (coder *BinWriter) UI64(f *uint64) {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], *f)
	coder.Write(buf[:])
}

// Int int reader
func (coder *BinReader) Int(f *int) {
	var buf [8]byte
	_, err := coder.Read(buf[:])
	if err != nil {
		return
	}
	*f = int(binary.LittleEndian.Uint64(buf[:]))
}

// Int int writer
func (coder *BinWriter) Int(f *int) {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], uint64(*f))
	coder.Write(buf[:])
}

// Int int reader
func (coder *BinReader) I64(f *int64) {
	var buf [8]byte
	_, err := coder.Read(buf[:])
	if err != nil {
		return
	}
	*f = int64(binary.LittleEndian.Uint64(buf[:]))
}

// Int int writer
func (coder *BinWriter) I64(f *int64) {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], uint64(*f))
	coder.Write(buf[:])
}

// Slice reads length [4]byte, entries [length*sizeof(E)]byte
func (coder *BinReader) Slice(
	length int, constructor func(int), iterator func(int)) {
	if coder.err != nil {
		return
	}
	constructor(length)
	for i := 0; i < length; i++ {
		iterator(i)
	}
}

// Slice writes length [4]byte, entries [length*sizeof(E)]byte
func (coder *BinWriter) Slice(length int, constructor func(int),
	iterator func(int)) {
	if coder.err != nil {
		return
	}
	for i := 0; i < length; i++ {
		iterator(i)
	}
}

// reads size [4]byte, content [size]byte from Source
func (coder *BinReader) String(f *string) {
	if coder.err != nil {
		return
	}
	var size uint64
	coder.VarInt(&size)
	c := make([]byte, size)
	coder.Read(c)
	*f = string(c)
}

// writes size [4]byte, content [size]byte to target
func (coder *BinWriter) String(f *string) {
	if coder.err != nil {
		return
	}
	c := []byte(*f)
	size := uint64(len(c))
	coder.VarInt(&size)
	coder.Write(c)
}

// Byte Slice
func byteSliceCoder(f *[]byte, coder Bincoder, length int) {
	coder.Bytes(length,
		func() []byte { return *f }, func(value []byte) {
			*f = value
		})
}

// ByteSlice field reader
func (coder *BinReader) ByteSlice(f *[]byte, length int) {
	byteSliceCoder(f, coder, length)
}

// ByteSlice field writer
func (coder *BinWriter) ByteSlice(f *[]byte, length int) {
	byteSliceCoder(f, coder, length)
}

// Bytes reader
func (coder *BinReader) Bytes(
	length int,
	getter func() []byte,
	setter func([]byte)) {
	if coder.err != nil {
		return
	}
	buf := make([]byte, length)
	_, err := coder.Read(buf)
	if err == nil {
		setter(buf)
	}
}

// Bytes writer
func (coder *BinWriter) Bytes(length int,
	getter func() []byte, setter func([]byte)) {
	if coder.err != nil {
		return
	}
	b := getter()
	if length > len(b) {
		larger := make([]byte, length)
		copy(larger, b)
		b = larger
	} else if length < len(b) {
		b = b[0:length]
	}
	coder.Write(b)
}

// VarInt reader
func (coder *BinReader) VarInt(f *uint64) {
	if coder.err != nil {
		return
	}
	buf := [9]byte{}
	_, err := coder.Read(buf[0:1])
	if err != nil {
		return
	}
	d := buf[0]
	if d < 0xFD {
		*f = uint64(d)
	} else if d == 0xFD {
		_, err = coder.Read(buf[1:3])
		if err != nil {
			return
		}
		*f = uint64(binary.LittleEndian.Uint16(buf[1:3]))
	} else if d == 0xFE {
		_, err = coder.Read(buf[1:5])
		if err != nil {
			return
		}
		*f = uint64(binary.LittleEndian.Uint32(buf[1:5]))
	} else {
		_, err = coder.Read(buf[1:9])
		if err != nil {
			return
		}
		*f = binary.LittleEndian.Uint64(buf[1:9])
	}
}

// VarInt writer
func (coder *BinWriter) VarInt(f *uint64) {
	n := *f
	buf := []byte{}
	if n < 0xFD {
		buf = []byte{byte(n)}
	} else if n < 0xFFFF {
		// <= 0xFFFF	3	0xFD followed by the length as uint16_t
		buf = make([]byte, 3)
		buf[0] = 0xFD
		binary.LittleEndian.PutUint16(buf[1:], uint16(n))
	} else if n < 0xFFFFFFFF {
		// <= 0xFFFF FFFF	5	0xFE followed by the length as uint32_t
		buf = make([]byte, 5)
		buf[0] = 0xFE
		binary.LittleEndian.PutUint32(buf[1:], uint32(n))

	} else {
		// -	9	0xFF followed by the length as uint64_t
		buf = make([]byte, 9)
		buf[0] = 0xFF
		binary.LittleEndian.PutUint64(buf[1:], uint64(n))
	}

	coder.Write(buf)
}
