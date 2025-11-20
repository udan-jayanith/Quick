package Quick_test

import (
	"testing"

	quick "github.com/udan-jayanith/Quick"
)

func TestVarInt(t *testing.T) {
	for _, num := range [...]int64{0, 60, 63, 16383, 1073741823, 4611686018427387903, 4611686018427387903 / 2} {
		b, err := quick.Int62ToVarint(num)
		if err != nil {
			t.Fatal(err)
		}

		v, err := quick.VarintToInt62(b)
		if err != nil {
			t.Fatal(err)
		} else if v != num {
			t.Fatal("Expected", num, "but got", v)
		}
	}

	if _, err := quick.Int62ToVarint(4611686018427387903 + 1); err == nil {
		t.Fatal("Expected a error but got nil")
	}
}

func TestReadVarint62(t *testing.T){
}