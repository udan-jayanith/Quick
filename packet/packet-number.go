package Packet

import (
	"encoding/binary"
	"math"

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

const (
	None varint.Int62 = varint.MaxInt62 + 1
)

// The EncodePacketNumber function takes two arguments:
//
// fullPn is the full packet number of the packet being sent.
// largestAcked is the largest packet number that has been acknowledged by the peer in the current packet number space, if not value to largestAcked must me None.
func EncodePacketNumber(fullPn, largestAcked varint.Int62) []byte {
	// The number of bits must be at least one more
	// than the base-2 logarithm of the number of contiguous
	// unacknowledged packet numbers, including the new packet.
	var numUnacked varint.Int62
	if largestAcked == None {
		numUnacked = fullPn + 1
	} else {
		numUnacked = fullPn - largestAcked
	}

	// There is a better way to do this.
	minBits := math.Log2(float64(numUnacked))
	numBytes := uint8(math.Ceil(minBits / 8))

	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], uint64(fullPn))

	return buf[8-numBytes:]
}
