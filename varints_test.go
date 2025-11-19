package Quick_test

import (
	"testing"

	quick "github.com/udan-jayanith/Quick"
)

func TestVarInt(t *testing.T) {
	b, err := quick.ToVarint(60)
	if err != nil {
		t.Fatal(err)
	}

	v, err := quick.VarintToInt64(b)
	if err != nil {
		t.Fatal(err)
	} else if v != 60 {
		t.Fatal("Expected 60 but got", v)
	}
}
