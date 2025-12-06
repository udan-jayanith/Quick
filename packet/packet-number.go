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
	None PacketNumber = 0
)

// Parameter values must derived from the same packet number space.
func EncodePacketNumber(packetNumber, largestAcknowledgedPacketNumber PacketNumber) ([]byte, error) {
	if packetNumber.IsOverflowing() || largestAcknowledgedPacketNumber.IsOverflowing() {
		return []byte{}, varint.IntegerOverflow
	}

	bytes := packetNumberLength(packetNumber, largestAcknowledgedPacketNumber)
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(packetNumber))

	// bytes is never 0.
	return b[len(b)-bytes:], nil
}

// From https://go-review.googlesource.com/c/net/+/301451
func packetNumberLength(packetNumber, largestAcknowledgedPacketNumber PacketNumber) int {
	diff := max(packetNumber-largestAcknowledgedPacketNumber) - min(packetNumber-largestAcknowledgedPacketNumber)

	if diff < 256 {
		//2^8
		return 1
	} else if diff < 65536 {
		//2^(8*2)
		return 2
	} else if diff < 16777216 {
		//2^(8*3)
		return 3
	}
	return 4
	/*
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
	*/
}

// The DecodePacketNumber function takes three arguments:
// largestPacketNumber is the largest packet number that has been successfully processed in the current packet number space.
// packetNumber is the value of the Packet Number field.
func DecodePacketNumber(packetNumber []byte, largestPacketNumber PacketNumber) (PacketNumber, error) {
	if len(packetNumber) > 4 {
		return 0, varint.IntegerOverflow
	}

	// 0 value bytes to packetNumber to decode the packetNumber using bit endian.
	//packetNumber = append(make([]byte, 8-len(packetNumber)), packetNumber...)

	var mask PacketNumber
	{
		b := make([]byte, len(packetNumber))
		for i, v := range b {
			b[i] = ^v
		}
		b = fillUpTo8Bytes(b)
		mask = ^PacketNumber(binary.BigEndian.Uint64(b))
	}
	packetNumber = fillUpTo8Bytes(packetNumber)
	pn := PacketNumber(binary.BigEndian.Uint64(packetNumber))
	expectedPacketNumber := (largestPacketNumber & mask) | pn
	return expectedPacketNumber, nil
	//estimatedPacketNumber := expectedPacketNumber

	/*
		noOfBits := len(packetNumber) * 8

		pnWin := PacketNumber(1 << noOfBits)
		pnHwin := PacketNumber(pnWin / 2)
		pnMask := PacketNumber(pnWin - 1)

		// The incoming packet number should be greater than
		// expected_pn - pn_hwin and less than or equal to
		// expected_pn + pn_hwin
		//
		// This means we cannot just strip the trailing bits from
		// expected_pn and add the truncated_pn because that might
		// yield a value outside the window.

		truncatedPacketNumber := PacketNumber(binary.BigEndian.Uint64(packetNumber))
		candidatePacketNumber := (expectedPacketNumber & ^pnMask) | truncatedPacketNumber

		if candidatePacketNumber <= expectedPacketNumber - pnWin{

		}
	*/
}

// fillUpTo8Bytes prepend 0 byte values to b to make up for length of b to be 8. len(b) must be less then 9.
func fillUpTo8Bytes(b []byte) []byte {
	return append(make([]byte, 8-len(b), 8), b...)
}
