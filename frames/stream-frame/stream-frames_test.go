package StreamFrame_test

import (
	"bufio"
	"bytes"
	"io"
	"math/rand"
	"slices"
	"testing"

	quick "github.com/udan-jayanith/Quick"
	StreamFrame "github.com/udan-jayanith/Quick/frames/stream-frame"
	StreamIdentifier "github.com/udan-jayanith/Quick/stream-identifier"
	"github.com/udan-jayanith/Quick/varint"

	//testing libraries
	Quick_test "github.com/udan-jayanith/Quick/internal/testing"
	"github.com/xyproto/randomstring"
)

var (
	fTypeTestcases = [...]struct {
		offset, length, fin bool
		fType               StreamFrame.StreamFrameType
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
		frameType := StreamFrame.NewStreamFrameType()
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
	streamFrame := StreamFrame.StreamFrame{
		Type:     StreamFrame.NewStreamFrameType().SetFin(true),
		StreamID: StreamIdentifier.NewStreamID(StreamIdentifier.ClientInitiatedUni),
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
	newFrame, err := StreamFrame.ReadStreamFrame(rd)
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
		t.Log(Quick_test.ToFormattedJson(streamFrame))
		t.Log("but got")
		t.Log(Quick_test.ToFormattedJson(newFrame))
		t.FailNow()
	}
}

func generateStreamFrames() []StreamFrame.StreamFrame {
	res := make([]StreamFrame.StreamFrame, 0, len(fTypeTestcases))
	for _, t := range fTypeTestcases {
		sf := StreamFrame.StreamFrame{
			Type: t.fType,
		}
		sf.StreamID = StreamIdentifier.NewStreamID(varint.Int62(rand.Int() % int(StreamIdentifier.MaxStreamID)))

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

		newSf, qErr := StreamFrame.ReadStreamFrame(bufio.NewReader(bytes.NewReader(b)))
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
			t.Log(Quick_test.ToFormattedJson(&sf))
			t.Log("but got")
			t.Log(Quick_test.ToFormattedJson(&newSf))
		}
	}
}
