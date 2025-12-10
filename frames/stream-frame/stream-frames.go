package StreamFrame

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"

	QuicErr "github.com/udan-jayanith/Quick/errors"
	StreamIdentifier "github.com/udan-jayanith/Quick/stream-identifier"
	"github.com/udan-jayanith/Quick/varint"
)

// 8 bits long, StreamFrameType can convert to Int62 by doing
//
//	Int62(StreamFrameType)
//	---
type StreamFrameType uint8

var (
	InvalidStreamFrameType error = errors.New("StreamFrameType is invalid")
)

func (sft StreamFrameType) IsValid() bool {
	if sft > 15 || sft < 8 {
		return false
	}
	return true
}

func (sft StreamFrameType) GetOffset() bool {
	var mask StreamFrameType = 0b_00_001_100
	return sft&mask == mask
}

func (sft StreamFrameType) GetLength() bool {
	var mask StreamFrameType = 0b_00_001_010
	return sft&mask == mask
}

// The FIN bit (0x01) indicates that the frame marks the end of the stream.
// The final size of the stream is the sum of the offset and the length of this frame.
func (sft StreamFrameType) GetFin() bool {
	var mask StreamFrameType = 0b_00_00_10_01
	return sft&mask == mask
}

func NewStreamFrameType() StreamFrameType {
	return 0b_1000
}

func (sft StreamFrameType) SetOffset(v bool) StreamFrameType {
	//0b_00_000_100
	var mask StreamFrameType = 0b_00_001_100
	if v {
		return sft | mask
	} else {
		return sft & 0b_00_001_011
	}
}

func (sft StreamFrameType) SetLength(v bool) StreamFrameType {
	//0b_00_00_00_10
	var mask StreamFrameType = 0b_00_001_010
	if v {
		return sft | mask
	} else {
		return sft & 0b_00_001_101
	}
}

func (sft StreamFrameType) SetFin(v bool) StreamFrameType {
	//0b_00_00_00_01
	var mask StreamFrameType = 0b_00_001_001
	if v {
		return sft | mask
	} else {
		return sft & 0b_00_001_110
	}
}

/*
STREAM Frame {
  Type (i) = 0x08..0x0f,
  Stream ID (i),
  [Offset (i)],
  [Length (i)],
  Stream Data (..),
}
*/

// STREAM frames implicitly create a stream and carry stream data.
type StreamFrame struct {
	Type     StreamFrameType //Half the byte is empty. Only LS 4 bytes is in use.
	StreamID StreamIdentifier.StreamID

	// Offset is starting index which the StreamData should be place in the stream.
	Offset varint.Int62
	// Offset is optional and uses variable length encoding.
	// If Offset of the frame is not specified in the Type of the frame. Offset of the frame is considered 0.

	// When the Length is 0, the Offset in the STREAM frame is the offset of the next byte that would be sent.
	Length varint.Int62 // Length is optional and uses variable length encoding.
	// Offset of the the stream and the Length of the frame cannot overflow int62.

	StreamData *bytes.Reader
}

/*
	STREAM Frame {
	  Type (i) = 0x08..0x0f,
	  Stream ID (i),
	  [Offset (i)],
	  [Length (i)],
	  Stream Data (..),
	}
*/

// Encode returns
/*
	Type (i),
	Stream ID (i),
	[Offset (i)],
	[Length (i)],
*/
// as []byte and then returns Stream Data (..) as a reader and an error if any.
// bytes reader may me empty depending on the StreamFrame.
// All the return values are invalid if Encode returns an error there for user first must check for errors.
func (sf *StreamFrame) Encode() ([]byte, *bytes.Reader, error) {
	buf := make([]byte, 0, 40)

	// StreamFrame type
	{
		if !sf.Type.IsValid() {
			fmt.Printf("%b\n", sf.Type)
			fmt.Printf("%v\n", sf.Type)
			return []byte{}, nil, InvalidStreamFrameType
		}

		b, err := varint.Int62ToVarint(varint.Int62(sf.Type))
		if err != nil {
			return []byte{}, nil, err
		}
		buf = append(buf, b...)
	}

	//Stream ID
	if b, err := sf.StreamID.ToVariableLength(); err != nil {
		return []byte{}, nil, err
	} else {
		buf = append(buf, b...)
	}

	//Offset
	if sf.Type.GetOffset() {
		b, err := varint.Int62ToVarint(sf.Offset)
		if err != nil {
			return []byte{}, nil, err
		}
		buf = append(buf, b...)
	}

	//Length
	if sf.Type.GetLength() {
		b, err := varint.Int62ToVarint(sf.Length)
		if err != nil {
			return []byte{}, nil, err
		}
		buf = append(buf, b...)
	}

	return buf, sf.StreamData, nil
}

func ReadStreamFrame(rd *bufio.Reader) (StreamFrame, QuicErr.Err) {
	sf := StreamFrame{}

	//Decode the frame type.
	if v, err := varint.ReadVarint62(rd); err != nil {
		return sf, QuicErr.FRAME_ENCODING_ERROR
	} else {
		//Check for eligibility of the frame type
		sf.Type = StreamFrameType(v)
	}

	//Decode stream id
	if v, err := varint.ReadVarint62(rd); err != nil {
		return sf, QuicErr.FRAME_ENCODING_ERROR
	} else {
		sf.StreamID = StreamIdentifier.NewStreamID(v)
	}

	//Decode offset if it's in the frame.
	if sf.Type.GetOffset() {
		if v, err := varint.ReadVarint62(rd); err != nil {
			return sf, QuicErr.FRAME_ENCODING_ERROR
		} else {
			sf.Offset = v
		}
	}

	//Decode length if it's in the frame.
	if sf.Type.GetLength() {
		if v, err := varint.ReadVarint62(rd); err != nil {
			return sf, QuicErr.FRAME_ENCODING_ERROR
		} else {
			sf.Length = v
		}
	}

	/*
		-- the sum of the offset and data length -- cannot exceed 262-1, as it is not possible to provide flow control credit for that data.
		Receipt of a frame that exceeds this limit MUST be treated as a connection error of type FRAME_ENCODING_ERROR or FLOW_CONTROL_ERROR.
	*/
	if (sf.Offset + sf.Length).IsOverflowing() {
		return sf, QuicErr.FRAME_ENCODING_ERROR
	}

	if sf.Length == 0 {
		return sf, QuicErr.NO_ERROR
	}

	// Read the stream data
	buf := make([]byte, sf.Length)
	if _, err := io.ReadFull(rd, buf); err != nil {
		return sf, QuicErr.FRAME_ENCODING_ERROR
	}

	sf.StreamData = bytes.NewReader(buf)
	return sf, QuicErr.NO_ERROR
}
