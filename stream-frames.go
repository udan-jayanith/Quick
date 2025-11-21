package Quick

import (
	"bufio"
	"bytes"
	"io"
)

// 8 bits long
type StreamFrameType uint8

func (sft StreamFrameType) IsValid() bool {
	if sft > 15 || sft < 8 {
		return false
	}
	return true
}

func (sft StreamFrameType) GetOffset() bool {
	return sft&0b_00_000_100 == 0b_00_000_100
}

func (sft StreamFrameType) GetLength() bool {
	return sft&0b_00_00_00_10 == 0b_00_00_00_10
}

// The FIN bit (0x01) indicates that the frame marks the end of the stream.
// The final size of the stream is the sum of the offset and the length of this frame.
func (sft StreamFrameType) GetFin() bool {
	return sft&0b_00_00_00_01 == 0b_00_00_00_01
}

func NewStreamFrameType() StreamFrameType {
	return 0
}

func (sft StreamFrameType) SetOffset(v bool) StreamFrameType {
	//0b_00_000_100
	if v {
		return sft | 0b_00_000_100
	} else {
		return sft & 0b_11_111_011
	}
}

func (sft StreamFrameType) SetLength(v bool) StreamFrameType {
	//0b_00_00_00_10
	if v {
		return sft | 0b_00_00_00_10
	} else {
		return sft & 0b_11_11_11_01
	}
}

func (sft StreamFrameType) SetFin(v bool) StreamFrameType {
	//0b_00_00_00_01
	if v {
		return sft | 0b_00_00_00_01
	} else {
		return sft & 0b_11_11_11_10
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
	StreamID StreamID

	// Offset is starting index which the StreamData should be place in the stream.
	Offset Int62
	// Offset is optional and uses variable length encoding.
	// If Offset of the frame is not specified in the Type of the frame. Offset of the frame is considered 0.

	// When the Length is 0, the Offset in the STREAM frame is the offset of the next byte that would be sent.
	Length Int62 // Length is optional and uses variable length encoding.
	// Offset of the the stream and the Length of the frame cannot overflow int62.

	StreamData *bytes.Reader
}

func ReadStreamFrame(rd *bufio.Reader) (StreamFrame, QuickTransportError) {
	sf := StreamFrame{}

	//Decode the frame type.
	if v, err := ReadVarint62(rd); err != nil {
		return sf, FLOW_CONTROL_ERROR
	} else {
		sf.Type = StreamFrameType(v)
	}

	//Decode stream id
	if v, err := ReadVarint62(rd); err != nil {
		return sf, FLOW_CONTROL_ERROR
	} else {
		sf.StreamID = NewStreamID(v)
	}

	//Decode offset if it's in the frame.
	if sf.Type.GetOffset() {
		if v, err := ReadVarint62(rd); err != nil {
			return sf, FLOW_CONTROL_ERROR
		} else {
			sf.Offset = v
		}
	}

	//Decode length if it's in the frame.
	if sf.Type.GetLength() {
		if v, err := ReadVarint62(rd); err != nil {
			return sf, FLOW_CONTROL_ERROR
		} else {
			sf.Length = v
		}
	}

	/*
		-- the sum of the offset and data length -- cannot exceed 262-1, as it is not possible to provide flow control credit for that data.
		Receipt of a frame that exceeds this limit MUST be treated as a connection error of type FRAME_ENCODING_ERROR or FLOW_CONTROL_ERROR.
	*/
	if (sf.Offset + sf.Length).IsOverflowing() {
		return sf, FLOW_CONTROL_ERROR
	}

	if sf.Length == 0 {
		return sf, NO_ERROR
	}

	// Read the stream data
	buf := make([]byte, sf.Length)
	if _, err := io.ReadFull(rd, buf); err != nil {
		return sf, FLOW_CONTROL_ERROR
	}

	sf.StreamData = bytes.NewReader(buf)
	return sf, NO_ERROR
}
