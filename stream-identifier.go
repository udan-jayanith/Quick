package Quick

import (
	"encoding/binary"
	"errors"
)

const (
	MaxStreamID int64 = (int64(1) << 62) - 1 // 2^62-1
)

type StreamType = int64

const (
	ClientInitiatedBidi StreamType = iota + 0b_00 // 0b_00
	ServerInitiatedBidi                           // 0b_01
	ClientInitiatedUndi                           // 0b_10
	ServerInitiatedUndi                           // 0b_11
)

type StreamID struct {
	streamID int64
}

func NewStreamID(streamType StreamType) StreamID {
	streamId := StreamID{
		streamID: int64(streamType),
	}

	return streamId
}

var (
	IntegerOverflow error = errors.New("Integer overflow")
)

func (si *StreamID) Increment() error {
	si.streamID += 4
	if si.streamID > MaxStreamID {
		si.streamID -= 4
		return IntegerOverflow
	}
	return nil
}

func (si *StreamID) StreamType() StreamType {
	return si.streamID & 0b_11
}

func (si *StreamID) ToVariableLengthInt() ([]byte, error) {
	return ToVarint(si.streamID)
}

func ToVarint(v int64) ([]byte, error) {
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
	buf := make([]byte, bytes)

	_, err := binary.Encode(buf, binary.BigEndian, v)
	if err != nil {
		return buf, err
	}

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

func VarintToInt64(b []byte) int64 {
	panic("Unimplemented")
}
