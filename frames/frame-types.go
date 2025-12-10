package Frame

import (
	"bufio"

	QuicErr "github.com/udan-jayanith/Quick/errors"
	"github.com/udan-jayanith/Quick/varint"
)

type FrameType uint8

const (
	//0x00
	Padding FrameType = 0+iota
	//0x01
	Ping
	//0x02-0x03
	Ack
	//0x04
	ResetStream
	//0x05
	StopSending
	//0x06
	Crypto
	//0x07
	NewToken
	//0x08-0x0f
	Stream
	//0x10
	MaxData
	//0x11
	MaxStreamData
	//0x12-0x13
	MaxStreams
	//0x14
	DataBlocked
	//0x15
	StreamDataBlocked
	//0x16-0x17
	StreamsBlocked
	//0x18
	NewConnectionId
	//0x19
	RetierConnectionId
	//0x1a
	PathChallenge
	//0x1b
	PathResponse
	//0x1c-0x1d
	ConnectionClose
	//0x1e
	HandshakeDone
)

func FrameValueToType(frameValue byte) (FrameType, QuicErr.Err) {
	if frameValue == 0{
		return Padding, QuicErr.NO_ERROR
	}else if frameValue == 1{
		return Ping, QuicErr.NO_ERROR
	}else if frameValue == 2 || frameValue == 3{
		return Ack, QuicErr.NO_ERROR
	}else if frameValue == 4{
		return ResetStream, QuicErr.NO_ERROR
	}else if frameValue ==5 {
		return StopSending, QuicErr.NO_ERROR
	}else if frameValue ==6 {
		return Crypto, QuicErr.NO_ERROR
	}else if frameValue ==7 {
		return NewToken, QuicErr.NO_ERROR
	}else if frameValue >= 8 && frameValue <= 0x0f {
		return Stream, QuicErr.NO_ERROR
	}else if frameValue ==0x10{
		return MaxData, QuicErr.NO_ERROR
	}else if frameValue ==0x11 {
		return MaxStreamData, QuicErr.NO_ERROR
	}else if frameValue >= 0x12 && frameValue <= 0x13 {
		return MaxStreams, QuicErr.NO_ERROR
	}else if frameValue ==0x14 {
		return DataBlocked, QuicErr.NO_ERROR
	}else if frameValue ==0x15 {
		return StreamDataBlocked, QuicErr.NO_ERROR
	}else if frameValue >= 0x16 && frameValue <= 0x17 {
		return StreamsBlocked, QuicErr.NO_ERROR
	}else if frameValue ==0x18 {
		return NewConnectionId, QuicErr.NO_ERROR
	}else if frameValue ==0x19 {
		return RetierConnectionId, QuicErr.NO_ERROR
	}else if frameValue ==0x1a {
		return PathChallenge, QuicErr.NO_ERROR
	}else if frameValue ==0x1b {
		return PathResponse, QuicErr.NO_ERROR
	}else if frameValue >= 0x1c && frameValue <= 0x1d {
		return ConnectionClose, QuicErr.NO_ERROR
	}else if frameValue == 0x1e{
		return HandshakeDone, QuicErr.NO_ERROR
	}
	return 0, QuicErr.FRAME_ENCODING_ERROR
}

// ReadFrameType returns the type of the frame, frame type value and an quick error.
// If any error occurred ReadFrameType returns 0, 0 and a error.
func ReadFrameType(rd *bufio.Reader) (FrameType, uint8, QuicErr.Err){
	frameValue, err := varint.ReadVarint62(rd)
	if err != nil {
		return 0, 0, QuicErr.FRAME_ENCODING_ERROR
	}

	frameType, qerr := FrameValueToType(byte(frameValue))
	if qerr != QuicErr.NO_ERROR {
		return 0, 0, qerr
	}

	return frameType, uint8(frameValue), QuicErr.NO_ERROR
}
