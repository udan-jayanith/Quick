package Quick_test

import (
	"bufio"
	"bytes"
	"io"
	"math/rand"
	"slices"
	"testing"

	quick "github.com/udan-jayanith/Quick"
	"github.com/udan-jayanith/Quick/varint"
	"github.com/xyproto/randomstring"
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
			fType:  0b_00_001_100,
		},
		{
			offset: false,
			length: true,
			fin:    false,
			fType:  0b_00_001_010,
		},
		{
			offset: false,
			length: false,
			fin:    true,
			fType:  0b_00_001_001,
		},
		{
			offset: false,
			length: false,
			fin:    false,
			fType:  0b_00_001_000,
		},
		{
			offset: true,
			length: true,
			fin:    false,
			fType:  0b_00_001_110,
		},
		{
			offset: true,
			length: false,
			fin:    true,
			fType:  0b_00_001_101,
		},
		{
			offset: true,
			length: true,
			fin:    true,
			fType:  0b_00_001_111,
		},
	}
)

func TestStreamFrameTypeGetters(t *testing.T) {
	for _, fType := range fTypeTestcases {
		if !fType.fType.IsValid() {
			t.Fatal("Invalid stream frame type", fType.fType)
		}

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

func TestReadStreamFrame(t *testing.T) {
	streamFrame := quick.StreamFrame{
		Type:     quick.NewStreamFrameType().SetFin(true),
		StreamID: quick.NewStreamID(quick.ClientInitiatedUndi),
	}
	streamFrame.StreamID.Increment()

	buf := make([]byte, 0)
	if b, err := varint.Int62ToVarint(varint.Int62(streamFrame.Type)); err != nil {
		t.Fatal(err.Error())
	} else {
		buf = append(buf, b...)
	}

	if b, err := streamFrame.StreamID.ToVariableLength(); err != nil {
		t.Fatal(err.Error())
	} else {
		buf = append(buf, b...)
	}

	// Add few extra bytes to the end the reader to check if ReadStreamFrame reads bytes more then it should.
	rd := bufio.NewReader(bytes.NewReader(append(buf, make([]byte, 10)...)))
	newFrame, err := quick.ReadStreamFrame(rd)
	if err != quick.NO_ERROR {
		t.Fatal("Quick frame parsing error", err)
	}

	if n, err := io.ReadFull(rd, make([]byte, 11)); err == nil {
		t.Fatal("Expected a error but got no error.")
	} else if n != 10 {
		t.Fatal("Expected n == 10 but n is", n)
	}

	if streamFrame != newFrame {
		t.Log("Expected")
		t.Log(ToFormattedJson(streamFrame))
		t.Log("but got")
		t.Log(ToFormattedJson(newFrame))
		t.FailNow()
	}
}

func generateStreamFrames() []quick.StreamFrame {
	res := make([]quick.StreamFrame, 0, len(fTypeTestcases))
	for _, t := range fTypeTestcases {
		sf := quick.StreamFrame{
			Type: t.fType,
		}
		sf.StreamID = quick.NewStreamID(varint.Int62(rand.Int() % int(quick.MaxStreamID)))

		if t.offset {
			sf.Offset = varint.Int62(rand.Int63() % int64(varint.MaxInt62))
		}
		if t.length {
			str := randomstring.HumanFriendlyEnglishString(int(rand.Int63() % 100))
			sf.Length = varint.Int62(len(str))
			sf.StreamData = bytes.NewReader([]byte(str))
		}
		res = append(res, sf)
	}
	return res
}

func TestStreamFrameEncodeAndDecode(t *testing.T) {
	for _, sf := range generateStreamFrames() {
		b, rd, err := sf.Encode()
		if err != nil {
			t.Fatal(err.Error())
		}

		//rd might be nil check for that
		if rd != nil {
			b1, err := io.ReadAll(rd)
			if err != nil {
				t.Fatal(err.Error())
			}
			b = append(b, b1...)
		}

		newSf, qErr := quick.ReadStreamFrame(bufio.NewReader(bytes.NewReader(b)))
		if qErr != quick.NO_ERROR {
			t.Fatal("Unexpected error", qErr)
		}

		if newSf.StreamData == nil && sf.StreamData != nil {
			t.Fatal("Missing stream data")
		} else if newSf.StreamData != nil {
			b2, err := io.ReadAll(newSf.StreamData)
			if err != nil {
				t.Fatal(err.Error())
			}

			b1 := b[len(b)-int(sf.Length):]
			if slices.Compare(b1, b2) != 0 {
				t.Log(string(b1))
				t.Log(string(b2))
				t.Fatal("Stream data of the both stream frames are different")
			}
		}

		sf.StreamData = nil
		newSf.StreamData = nil
		if sf != newSf {
			t.Log("Expected")
			t.Log(ToFormattedJson(&sf))
			t.Log("but got")
			t.Log(ToFormattedJson(&newSf))
		}
	}
}
