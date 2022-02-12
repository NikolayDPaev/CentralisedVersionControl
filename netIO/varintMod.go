package netio

// Modified code from the binary package
// the function modified is ReadUvarint, it was made to work with io.Reader instead of io.ByteReader
// because bytereader caches some bytes and that causes problems with consecutive reads.
import (
	"errors"
	"io"
)

const MaxVarintLen64 = 10

var errOverflow = errors.New("binary: varint overflows a 64-bit integer")

// ReadUvarint reads an encoded unsigned integer from r and returns it as a uint64.
func ReadUvarint(r io.Reader) (uint64, error) {
	var x uint64
	var s uint
	buf := make([]byte, 1)
	for i := 0; i < MaxVarintLen64; i++ {
		r.Read(buf)
		b := buf[0]

		if b < 0x80 {
			if i == MaxVarintLen64-1 && b > 1 {
				return x, errOverflow
			}
			return x | uint64(b)<<s, nil
		}
		x |= uint64(b&0x7f) << s
		s += 7
	}
	return x, errOverflow
}

// ReadVarint reads an encoded signed integer from r and returns it as an int64.
func ReadVarint(r io.Reader) (int64, error) {
	ux, err := ReadUvarint(r) // ok to continue in presence of error
	x := int64(ux >> 1)
	if ux&1 != 0 {
		x = ^x
	}
	return x, err
}
