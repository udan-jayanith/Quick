package Quick

import (
	"github.com/udan-jayanith/Quick/varint"
)

// MSB bits is needed for variable sized int62 encoding and decoding.
// StreamType bits is the LSB.

type StreamType = varint.Int62

const (
	MaxStreamID varint.Int62 = varint.MaxInt62
)

const (
	ClientInitiatedBidi StreamType = iota + 0b_00 // 0b_00
	ServerInitiatedBidi                           // 0b_01
	ClientInitiatedUni                           // 0b_10
	ServerInitiatedUni                           // 0b_11
)

// StreamID can convert into variable length int62 by calling ToVariableLength method
type StreamID struct {
	streamID varint.Int62
}

func NewStreamID(streamType StreamType) StreamID {
	streamId := StreamID{
		streamID: streamType,
	}

	return streamId
}

func (si *StreamID) Increment() error {
	si.streamID += 4
	if si.streamID.IsOverflowing() {
		si.streamID -= 4
		return varint.IntegerOverflow
	}
	return nil
}

func (si *StreamID) StreamType() StreamType {
	return StreamType(si.streamID & 0b_11)
}

func (si *StreamID) ToVariableLength() ([]byte, error) {
	return varint.Int62ToVarint(si.streamID)
}
