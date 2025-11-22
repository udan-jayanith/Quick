package varint

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
)

type Int62 uint64

const (
	MaxInt62 Int62 = Int62((uint64(1) << 62) - 1) // 2^62-1
)

var (
	IntegerOverflow error = errors.New("Integer overflow")
)

func (int62 Int62) IsOverflowing() bool {
	return int62 > MaxInt62
}

func Int62ToVarint(v Int62) ([]byte, error) {
	/*
		2MSB	Length	Usable Bits	 Range
		00		1		6			 0-63
		01		2		14			 0-16383
		10		4		30			 0-1073741823
		11		8		62			 0-4611686018427387903
	*/
	var bytes int
	if v <= 63 {
		bytes = 1
	} else if v <= 16383 {
		bytes = 2
	} else if v <= 1073741823 {
		bytes = 4
	} else if v <= 4611686018427387903 {
		bytes = 8
	} else {
		return []byte{}, IntegerOverflow
	}
	buf := make([]byte, 8)

	_, err := binary.Encode(buf, binary.BigEndian, v)
	if err != nil {
		return buf, err
	}
	buf = buf[len(buf)-bytes:]

	switch bytes {
	case 2:
		buf[0] = buf[0] | 0b_01_00_00_00
	case 4:
		buf[0] = buf[0] | 0b_10_00_00_00
	case 8:
		buf[0] = buf[0] | 0b_11_00_00_00
	}

	return buf, nil
}

// This is actually int62
func VarintToInt62(b []byte) (Int62, error) {
	if len(b) == 0 || len(b) > 8 {
		return 0, IntegerOverflow
	}

	b[0] = b[0] & 0b_00_11_11_11

	buf := make([]byte, 8-len(b))
	buf = append(buf, b...)

	var v Int62
	_, err := binary.Decode(buf, binary.BigEndian, &v)
	return v, err
}

func ReadVarint62(rd *bufio.Reader) (Int62, error) {
	/*
		2MSB	Length	Usable Bits	 Range
		00		1		6			 0-63
		01		2		14			 0-16383
		10		4		30			 0-1073741823
		11		8		62			 0-4611686018427387903
	*/

	b, err := rd.Peek(1)
	if err != nil {
		rd.ReadByte()
		return 0, err
	}

	var length int
	switch b[0] & 0b_11_00_00_00 {
	case 0b_0:
		length = 1
	case 0b_01_00_00_00:
		length = 2
	case 0b_10_00_00_00:
		length = 4
	case 0b_11_00_00_00:
		length = 8
	}

	if length == 0 {
		return 0, IntegerOverflow
	}

	p := make([]byte, length)
	if _, err := io.ReadFull(rd, p); err != nil {
		return 0, err
	}
	return VarintToInt62(p)
}
