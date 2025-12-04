package Packet

import (
	"encoding/binary"
	"math"

	"github.com/udan-jayanith/Quick/varint"
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
