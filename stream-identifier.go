package Quick

import (
	"errors"
)

// MSB bits is needed for variable sized int62 encoding and decoding.
// StreamType bits is the LSB. 
type Int62 uint64

const (
	MaxStreamID Int62 = Int62((uint64(1) << 62) - 1) // 2^62-1
)

func (int62 Int62) IsOverflowing() bool {
	return int62 > MaxStreamID
}

type StreamType = Int62

const (
	ClientInitiatedBidi StreamType = iota + 0b_00 // 0b_00
	ServerInitiatedBidi                           // 0b_01
	ClientInitiatedUndi                           // 0b_10
	ServerInitiatedUndi                           // 0b_11
)

type StreamID struct {
	streamID Int62
}

func NewStreamID(streamType StreamType) StreamID {
	streamId := StreamID{
		streamID: streamType,
	}

	return streamId
}

var (
	IntegerOverflow error = errors.New("Integer overflow")
)

func (si *StreamID) Increment() error {
	si.streamID += 4
	if si.streamID.IsOverflowing() {
		si.streamID -= 4
		return IntegerOverflow
	}
	return nil
}

func (si *StreamID) StreamType() StreamType {
	return StreamType(si.streamID & 0b_11)
}

func (si *StreamID) ToVariableLength() ([]byte, error) {
	return Int62ToVarint(si.streamID)
}
