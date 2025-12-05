package Packet

import (
	"encoding/binary"

	"github.com/udan-jayanith/Quick/varint"
)

// packet number space is the context in which a packet can be processed and acknowledged.
//
// Initial packets can only be sent with Initial packet protection keys and acknowledged in packets that are also Initial packets.
//
// Similarly, Handshake packets are sent at the Handshake encryption level and can only be acknowledged in Handshake packets.
type PacketNumberSpace uint8

const (
	//All Initial packets (Section 17.2.2) are in this space.
	//
	//https://datatracker.ietf.org/doc/html/rfc9000#section-12.3-5.2.1
	InitialSpace PacketNumberSpace = 0 + iota
	//All Handshake packets (Section 17.2.4) are in this space.
	//
	//https://datatracker.ietf.org/doc/html/rfc9000#section-12.3-5.4.1
	HandshakeSpace
	//All 0-RTT (Section 17.2.3) and 1-RTT (Section 17.3.1) packets are in this space.
	//
	//https://datatracker.ietf.org/doc/html/rfc9000#section-12.3-5.6.1
	ApplicationDataSpace
)

type PacketNumber = varint.Int62

const (
	None PacketNumber = varint.MaxInt62 + 1
)

func EncodePacketNumber(packetNumber, largestAcknowledgedPacketNumber PacketNumber) ([]byte, error) {
	bytes := packetNumberLength(packetNumber, largestAcknowledgedPacketNumber)
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(packetNumber))

	// bytes is never 0.
	return b[len(b)-bytes:], nil
}

// From https://go-review.googlesource.com/c/net/+/301451
func packetNumberLength(packetNumber, largestAcknowledgedPacketNumber PacketNumber) int {
	d := packetNumber - largestAcknowledgedPacketNumber
	switch {
	case d < 0x80:
		return 1
	case d < 0x8000:
		return 2
	case d < 0x800000:
		return 3
	default:
		return 4
	}
}
