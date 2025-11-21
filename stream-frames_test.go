package Quick_test

import (
	"testing"

	quick "github.com/udan-jayanith/Quick"
)

var (
	fTypeTestcases = [...]struct {
		offset, length, fin bool
		fType               quick.StreamFrameType
	}{
		{
			offset: true,
			length: false,
			fin:    false,
			fType:  0b_00_000_100,
		},
		{
			offset: false,
			length: true,
			fin:    false,
			fType:  0b_00_000_010,
		},
		{
			offset: false,
			length: false,
			fin:    true,
			fType:  0b_00_000_001,
		},
		{
			offset: false,
			length: false,
			fin:    false,
			fType:  0b_00_000_000,
		},
		{
			offset: true,
			length: true,
			fin:    false,
			fType:  0b_00_000_110,
		},
		{
			offset: true,
			length: false,
			fin:    true,
			fType:  0b_00_000_101,
		},
		{
			offset: true,
			length: true,
			fin:    true,
			fType:  0b_00_000_111,
		},
	}
)

func TestStreamFrameTypeGetters(t *testing.T) {
	for _, fType := range fTypeTestcases {
		if fType.offset != fType.fType.GetOffset() {
			t.Fatal("Expected offset", fType.offset, "but got", fType.fType.GetOffset())
		} else if fType.length != fType.fType.GetLength() {
			t.Fatal("Expected length", fType.length, "but got", fType.fType.GetLength())
		} else if fType.fin != fType.fType.GetFin() {
			t.Fatal("Expected fin", fType.fin, "but got", fType.fType.GetFin())
		}
	}
}

func TestStreamFrameTypeSetters(t *testing.T) {
	for i, fType := range fTypeTestcases {
		frameType := quick.NewStreamFrameType()
		frameType = frameType.SetLength(fType.length)
		frameType = frameType.SetOffset(fType.offset)
		frameType = frameType.SetFin(fType.fin)

		if fType.fType != frameType {
			t.Log("Failed", i+1, "testcase")
			t.Fatalf("Expected %b but got %b", fType.fType, frameType)
		}
	}
}
