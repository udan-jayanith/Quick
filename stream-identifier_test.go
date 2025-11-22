package Quick_test

import (
	"testing"

	quick "github.com/udan-jayanith/Quick"
)

func TestStreamID_Increment(t *testing.T) {
	{
		streamID := quick.NewStreamID(quick.ClientInitiatedBidi)
		if err := streamID.Increment(); err != nil {
			t.Fatal("Unexpected error", err.Error())
		}

		if err := streamID.Increment(); err != nil {
			t.Fatal("Unexpected error", err.Error())
		}
	}

	{
		streamID := quick.NewStreamID(quick.MaxStreamID)
		if err := streamID.Increment(); err == nil {
			t.Fatal("Expected a error but got no error")
		}
	}

	{
		streamID := quick.NewStreamID(quick.MaxStreamID + 1)
		if err := streamID.Increment(); err == nil {
			t.Fatal("Expected a error but got no error")
		}
	}
}

func TestStreamID_StreamType(t *testing.T) {
	for _, testcase := range []struct {
		streamID quick.StreamID
		output   quick.StreamType
	}{
		{
			streamID: quick.NewStreamID(quick.ClientInitiatedBidi),
			output:   quick.ClientInitiatedBidi,
		},
		{
			streamID: quick.NewStreamID(quick.ClientInitiatedUni),
			output:   quick.ClientInitiatedUni,
		},
		{
			streamID: quick.NewStreamID(quick.ServerInitiatedBidi),
			output:   quick.ServerInitiatedBidi,
		},
		{
			streamID: quick.NewStreamID(quick.ServerInitiatedUni),
			output:   quick.ServerInitiatedUni,
		},
	} {
		st := testcase.streamID.StreamType()
		if st != testcase.output {
			t.Fatal("Expected", testcase.output, "but got", st)
		}
	}

	{
		streamID := quick.NewStreamID(quick.ServerInitiatedBidi)
		for range 22 {
			if err := streamID.Increment(); err != nil {
				t.Fatal("Unexpected error", err.Error())
			}
		}

		st := streamID.StreamType()
		if st != quick.ServerInitiatedBidi {
			t.Fatal("Expected", quick.ServerInitiatedBidi, "but got", st)
		}
	}

	{
		streamID := quick.NewStreamID(quick.MaxStreamID)
		st := streamID.StreamType()
		if st != quick.ServerInitiatedUni {
			t.Fatal("Expected", quick.ServerInitiatedUni, "but got", st)
		}
	}

	{
		streamID := quick.NewStreamID(quick.MaxStreamID + 1)
		st := streamID.StreamType()
		if st != quick.ClientInitiatedBidi {
			t.Fatal("Expected", quick.ClientInitiatedBidi, "but got", st)
		}
	}
}
