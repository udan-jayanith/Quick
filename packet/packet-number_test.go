package Packet_test

import (
	"testing"

	Testing "github.com/udan-jayanith/Quick/internal/testing"

	Packet "github.com/udan-jayanith/Quick/packet"
)

var (
	packetNumberTestcases = [...]struct {
		PacketNumber, LargestAckPacketNumber, ExpectedPacketNumber Packet.PacketNumber
	}{
		{
			PacketNumber:           255,
			LargestAckPacketNumber: 254,
		},
		{
			PacketNumber:           0,
			LargestAckPacketNumber: 0,
		},
		{
			PacketNumber:           515,
			LargestAckPacketNumber: 500,
		},
		{
			PacketNumber:           0,
			LargestAckPacketNumber: 5252,
		},
		{
			PacketNumber:           5252,
			LargestAckPacketNumber: 582590259,
		}, {
			PacketNumber:           140581,
			LargestAckPacketNumber: 414115,
		},
		{
			PacketNumber:           5,
			LargestAckPacketNumber: 414115,
		}, {
			PacketNumber:           4294967295 + 2,
			LargestAckPacketNumber: 4294967280,
		},
	}
)

func TestPacketNumberEncodingAndDecoding(t *testing.T) {
	for i, testcase := range packetNumberTestcases {
		b, err := Packet.EncodePacketNumber(testcase.PacketNumber, testcase.LargestAckPacketNumber)
		if err != nil {
			t.Fatalf("Error occurred at %v\n%s\n%s", i, Testing.ToFormattedJson(testcase), err.Error())
		}

		packetNumber, err := Packet.DecodePacketNumber(b, testcase.LargestAckPacketNumber)
		if err != nil {
			t.Fatalf("Error occurred at %v\n%s\n%s", i, Testing.ToFormattedJson(testcase), err.Error())
		}

		if packetNumber != testcase.PacketNumber {
			t.Fatalf("Expected %v but got %v at %v\n%s", testcase.PacketNumber, packetNumber, i, Testing.ToFormattedJson(testcase))
		}
	}
}
