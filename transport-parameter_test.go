package Quick_test

import (
	"testing"

	quick "github.com/udan-jayanith/Quick"
)

func TestTransportParameterIDEnums(t *testing.T) {
	if quick.Original_destination_connection_id != 0x00 {
		t.Fatal("Expected Original_destination_connection_id to be 0x00 but got", quick.Original_destination_connection_id)
	} else if quick.Retry_source_connection_id != 0x10 {
		t.Fatal("Expected Retry_source_connection_id to be 0x10 but got", quick.Retry_source_connection_id)
	}
}
