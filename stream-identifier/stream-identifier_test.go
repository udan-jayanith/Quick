package StreamIdentifier_test

import (
	"testing"

	StreamIdentifier "github.com/udan-jayanith/Quick/stream-identifier"
)

func TestStreamID_Increment(t *testing.T) {
	{
		streamID := StreamIdentifier.NewStreamID(StreamIdentifier.ClientInitiatedBidi)
		if err := streamID.Increment(); err != nil {
			t.Fatal("Unexpected error", err.Error())
		}

		if err := streamID.Increment(); err != nil {
			t.Fatal("Unexpected error", err.Error())
		}
	}

	{
		streamID := StreamIdentifier.NewStreamID(StreamIdentifier.MaxStreamID)
		if err := streamID.Increment(); err == nil {
			t.Fatal("Expected a error but got no error")
		}
	}

	{
		streamID := StreamIdentifier.NewStreamID(StreamIdentifier.MaxStreamID + 1)
		if err := streamID.Increment(); err == nil {
			t.Fatal("Expected a error but got no error")
		}
	}
}

func TestStreamID_StreamType(t *testing.T) {
	for _, testcase := range []struct {
		streamID StreamIdentifier.StreamID
		output   StreamIdentifier.StreamType
	}{
		{
			streamID: StreamIdentifier.NewStreamID(StreamIdentifier.ClientInitiatedBidi),
			output:   StreamIdentifier.ClientInitiatedBidi,
		},
		{
			streamID: StreamIdentifier.NewStreamID(StreamIdentifier.ClientInitiatedUni),
			output:   StreamIdentifier.ClientInitiatedUni,
		},
		{
			streamID: StreamIdentifier.NewStreamID(StreamIdentifier.ServerInitiatedBidi),
			output:   StreamIdentifier.ServerInitiatedBidi,
		},
		{
			streamID: StreamIdentifier.NewStreamID(StreamIdentifier.ServerInitiatedUni),
			output:   StreamIdentifier.ServerInitiatedUni,
		},
	} {
		st := testcase.streamID.StreamType()
		if st != testcase.output {
			t.Fatal("Expected", testcase.output, "but got", st)
		}
	}

	{
		streamID := StreamIdentifier.NewStreamID(StreamIdentifier.ServerInitiatedBidi)
		for range 22 {
			if err := streamID.Increment(); err != nil {
				t.Fatal("Unexpected error", err.Error())
			}
		}

		st := streamID.StreamType()
		if st != StreamIdentifier.ServerInitiatedBidi {
			t.Fatal("Expected", StreamIdentifier.ServerInitiatedBidi, "but got", st)
		}
	}

	{
		streamID := StreamIdentifier.NewStreamID(StreamIdentifier.MaxStreamID)
		st := streamID.StreamType()
		if st != StreamIdentifier.ServerInitiatedUni {
			t.Fatal("Expected", StreamIdentifier.ServerInitiatedUni, "but got", st)
		}
	}

	{
		streamID := StreamIdentifier.NewStreamID(StreamIdentifier.MaxStreamID + 1)
		st := streamID.StreamType()
		if st != StreamIdentifier.ClientInitiatedBidi {
			t.Fatal("Expected", StreamIdentifier.ClientInitiatedBidi, "but got", st)
		}
	}
}
