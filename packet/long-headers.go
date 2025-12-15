package Packet

import (
	Version "github.com/udan-jayanith/Quick/version"
)

type PacketType = uint8

// mask: 00110000
const (
	Initial   PacketType = 0
	RTT0      PacketType = 0b_00010000
	Handshake PacketType = 0b_00100000
	Retry     PacketType = 0b_00110000
)

/*
Long headers are used for packets that are sent prior to the establishment of 1-RTT keys. Once 1-RTT keys are available, a sender switches to sending packets using the short header (Section 17.3). The long form allows for special packets -- such as the Version Negotiation packet -- to be represented in this uniform fixed-length packet format. Packets that use the long header contain the
following fields:

	Long Header Packet {
	  Header Form (1) = 1,
	  Fixed Bit (1) = 1,
	  Long Packet Type (2),
	  Type-Specific Bits (4), //0
	  Version (32),
	  Destination Connection ID Length (8),
	  Destination Connection ID (0..160),
	  Source Connection ID Length (8),
	  Source Connection ID (0..160),
	  Type-Specific Payload (..),
	}

//0x0c mask in byte 0 is protected with header protection and reserved for packet type initial and handshake. (Value of this field plus one) * 8 is the length of the packet number filed.
//0x03 mask in byte 0 represent the length of the packet number as an uint2.
*/
type LongHeader struct {
	//The most significent bit is set to 1.
	//0x80
	HeaderForm bool
	/*
		The next bit (0x40) of byte 0 is set to 1, unless the packet is a Version Negotiation packet. Packets containing a zero value for this bit are not valid packets in this version and MUST be discarded. A value of 1 for this bit allows QUIC to coexist with other protocols
	*/
	FixedBit bool
	/*The next two bits (those with a mask of 0x30) of byte 0 contain a packet type. Packet types are listed in Table 5.
	 */
	LongPacketType PacketType
	/*The semantics of the lower four bits (those with a mask of 0x0f) of byte 0 are determined by the packet type.*/
	TypeSpecificBits uint8
	/*The QUIC Version is a 32-bit field that follows the first byte. This field indicates the version of QUIC that is in use and determines how the rest of the protocol fields are interpreted.*/
	Version Version.QuickVersion
	//DestinationConnectionId is for the packet reciver.
	DestinationConnectionIdLength uint8
	DestinationConnectionId       uint32 //160..0
	//Source connection id is for packet sender
	SourceConnectionIdLength uint8
	SourceConnectionId       uint32 //160..0
	TypeSpecificPayload      []byte
}

/*
	 Initial Packet {
		 Header Form (1) = 1,
		 Fixed Bit (1) = 1,
		 Long Packet Type (2) = 0,
		 --Type specific bits--
		 Reserved Bits (2),
		 Packet Number Length (2),
		 --Type specific bits end--
		 Version (32),
		 Destination Connection ID Length (8),
		 Destination Connection ID (0..160),
		 Source Connection ID Length (8),
		 Source Connection ID (0..160),
		 --Type specific payload--
		 Token Length (i),
		 Token (..),
		 Length (i),
		 Packet Number (8..32),
		 Packet Payload (8..),
		 --Type specific payload end--

}
*/
type InitialPacket struct {
	LongHeader
}
