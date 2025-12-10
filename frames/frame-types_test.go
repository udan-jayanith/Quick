package Frame_test

import (
	"bufio"
	"bytes"
	"testing"

	QuicErr "github.com/udan-jayanith/Quick/errors"
	Frame "github.com/udan-jayanith/Quick/frames"
	Testing "github.com/udan-jayanith/Quick/internal/testing"
)

var (
	frameValueToTypeTestcases = [...]struct {
		FrameValue byte
		FrameType  Frame.FrameType
		Qerr       QuicErr.Err
	}{
		{
			FrameValue: 0x00,
			FrameType:  Frame.Padding,
			Qerr:       QuicErr.NO_ERROR,
		},

		{
			FrameValue: 0x06,
			FrameType:  Frame.Crypto,
			Qerr:       QuicErr.NO_ERROR,
		},
		{
			FrameValue: 0x11,
			FrameType:  Frame.MaxStreamData,
			Qerr:       QuicErr.NO_ERROR,
		},
		{
			FrameValue: 0x13,
			FrameType:  Frame.MaxStreams,
			Qerr:       QuicErr.NO_ERROR,
		},
		{
			FrameValue: 0x12,
			FrameType:  Frame.MaxStreams,
			Qerr:       QuicErr.NO_ERROR,
		},
		{
			FrameValue: 0x17,
			FrameType:  Frame.StreamsBlocked,
			Qerr:       QuicErr.NO_ERROR,
		},
		{
			FrameValue: 0x1e,
			FrameType:  Frame.HandshakeDone,
			Qerr:       QuicErr.NO_ERROR,
		},
		{
			FrameValue: 0x1e + 1,
			FrameType:  0,
			Qerr:       QuicErr.FRAME_ENCODING_ERROR,
		},
	}
)

func TestFrameValueToType(t *testing.T) {
	for i, testcase := range frameValueToTypeTestcases {
		frameType, err := Frame.FrameValueToType(testcase.FrameValue)
		if err != testcase.Qerr {
			t.Fatalf("Test %v failed.\nExpected error '%s' but got '%s'\n%s", i, testcase.Qerr.Error(), err.Error(), Testing.ToFormattedJson(testcase))
		} else if frameType != testcase.FrameType {
			t.Fatalf("Test %v failed.\nExpected %v but got %v\n%s", i, testcase.FrameType, frameType, Testing.ToFormattedJson(testcase))
		}
	}
}

var (
	readFrameTypeTestcases = [...]struct {
		FrameType  Frame.FrameType
		FrameValue uint8
		Err        QuicErr.Err
		Input      []byte
	}{
		{
			FrameType:  Frame.Padding,
			FrameValue: 0,
			Err:        QuicErr.NO_ERROR,
			Input:      Testing.Int62ToVarint(0),
		},
		{
			FrameType:  Frame.StreamsBlocked,
			FrameValue: 0x16,
			Err:        QuicErr.NO_ERROR,
			Input:      Testing.Int62ToVarint(0x16),
		},
		{
			FrameType:  0,
			FrameValue: 0,
			Err:        QuicErr.FRAME_ENCODING_ERROR,
			Input:      append(Testing.Int62ToVarint(0x1e+1), make([]byte, 10)...),
		},
	}
)

func TestReadFrameType(t *testing.T) {
	for i, testcase := range readFrameTypeTestcases {
		rd := bufio.NewReader(bytes.NewReader(testcase.Input))
		ft, fv, err := Frame.ReadFrameType(rd)
		if err != testcase.Err {
			t.Fatalf("Test %v failed.\nExpected error '%s' but got '%s'\n%s", i, testcase.Err.Error(), err.Error(), Testing.ToFormattedJson(testcase))
		} else if ft != testcase.FrameType {
			t.Fatalf("Test %v failed.\nExpected frame type %v but got %v\n%s", i, testcase.FrameType, ft, Testing.ToFormattedJson(testcase))
		} else if fv != testcase.FrameValue {
			t.Fatalf("Test %v failed.\nExpected frame type %v but got %v\n%s", i, testcase.FrameValue, fv, Testing.ToFormattedJson(testcase))
		}
	}
}
