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

// Parameter values must derived from the same packet number space.
func EncodePacketNumber(packetNumber, largestAcknowledgedPacketNumber PacketNumber) ([]byte, error) {
	if packetNumber.IsOverflowing() {
		return []byte{}, varint.IntegerOverflow
	}

	/*
		var numUnacked PacketNumber
		if largestAcknowledgedPacketNumber == None {
			numUnacked = packetNumber + 1
		} else {
			numUnacked = packetNumber - largestAcknowledgedPacketNumber
		}
	*/

	bytes := PacketNumberLength(packetNumber, largestAcknowledgedPacketNumber)
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(packetNumber))

	// bytes is never 0.
	return b[len(b)-bytes:], nil
}

// From https://go-review.googlesource.com/c/net/+/301451
// This works only for where packetNumber is larger then arg2
func PacketNumberLength(packetNumber, largestAcknowledgedPacketNumber PacketNumber) int {
	a := make([]byte, 0, 8)
	a = binary.BigEndian.AppendUint64(a, uint64(packetNumber))

	b := make([]byte, 0, 8)
	b = binary.BigEndian.AppendUint64(b, uint64(largestAcknowledgedPacketNumber))

	var cutBytes int
	for i := range a {
		if a[i] != b[i] {
			cutBytes = len(a) - (i)
			break
		}
	}
	if cutBytes == 0 {
		return 1
	}
	return cutBytes
}

// The DecodePacketNumber function takes three arguments:
// largestPacketNumber is the largest packet number that has been successfully processed in the current packet number space.
// packetNumber is the value of the Packet Number field.
// This is not the right solution
// This only works if packetNumber doesn't exceed 2^4 and packetNumber greater then largestPacketNumber.
func DecodePacketNumber(packetNumber []byte, largestPacketNumber PacketNumber) (PacketNumber, error) {
	/*
		var returnVal varint.Int62
		for range 2 {
			b := make([]byte, 0, 8)
			b = binary.BigEndian.AppendUint64(b, uint64(largestPacketNumber))
			copy(b[len(b)-len(packetNumber):], packetNumber)
			returnVal = PacketNumber(binary.BigEndian.Uint64(b))
			if returnVal < largestPacketNumber {
				largestPacketNumber = largestPacketNumber << 1
				continue
			}
			break
		}
		return returnVal, nil
	*/
	length := len(packetNumber)
	expected := largestPacketNumber + 1
	win := PacketNumber(1 << (length * 8))
	hwin := win / 2
	mask := win - 1
	candidate := (expected & ^mask) | PacketNumber(binary.BigEndian.Uint64(fillUpTo8Bytes(packetNumber)))
	if candidate <= expected-hwin && candidate < 1<<62-win {
		return candidate + win, nil
	}
	if candidate > expected+hwin && candidate >= win {
		return candidate - win, nil
	}
	return candidate, nil
}

// fillUpTo8Bytes prepend 0 byte values to b to make up for length of b to be 8. len(b) must be less then 9.

func fillUpTo8Bytes(b []byte) []byte {
	return append(make([]byte, 8-len(b), 8), b...)
}
