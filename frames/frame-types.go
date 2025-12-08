package Frame

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
