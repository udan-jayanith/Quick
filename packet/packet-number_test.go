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
			PacketNumber:           0xac5c02,
			LargestAckPacketNumber: 0xabe8b3,
			ExpectedPacketNumber:   0xac5c02,
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
