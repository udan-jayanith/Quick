package QuicErr_test

import (
	"testing"

	quick "github.com/udan-jayanith/Quick/errors"
)

func TestQuickTransportError(t *testing.T) {
	if quick.NO_VIABLE_PATH != 0x10 {
		t.Fatal("Unexpected error code for quick.NO_VIABLE_PATH. Expected error code 0x10")
	} else if quick.NO_ERROR != 0x00 {
		t.Fatal("Unexpected error code for quick.NO_ERROR. Expected error code 0x00")
	}
}
