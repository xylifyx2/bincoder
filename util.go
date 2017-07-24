package bincoder

import (
	"bufio"
	"bytes"
)

// Marshall to bytes
func Marshall(m func(w *bufio.Writer)) []byte {
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	m(writer)
	writer.Flush()
	return b.Bytes()
}

// Unmarshall from bytes
func Unmarshall(m func(w *bufio.Reader), bin []byte) {
	var b bytes.Buffer
	b.Write(bin)
	reader := bufio.NewReader(&b)
	m(reader)
}
