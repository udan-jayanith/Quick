package varint_test

import (
	"bufio"
	"bytes"
	"io"
	"testing"

	varint "github.com/udan-jayanith/Quick/varint"
)

func TestVarInt(t *testing.T) {
	for _, num := range [...]varint.Int62{0, 60, 63, 16383, 1073741823, 4611686018427387903, 4611686018427387903 / 2} {
		b, err := varint.Int62ToVarint(num)
		if err != nil {
			t.Fatal(err)
		}

		v, err := varint.VarintToInt62(b)
		if err != nil {
			t.Fatal(err)
		} else if v != num {
			t.Fatal("Expected", num, "but got", v)
		}
	}

	if _, err := varint.Int62ToVarint(4611686018427387903 + 1); err == nil {
		t.Fatal("Expected a error but got nil")
	}
}

func TestReadVarint62(t *testing.T) {
	{
		b, err := varint.Int62ToVarint(73)
		if err != nil {
			t.Fatal(err.Error())
		}

		b = append(b, make([]byte, 10)...)
		rd := bufio.NewReader(bytes.NewReader(b))
		v, err := varint.ReadVarint62(rd)
		if err != nil {
			t.Fatal(err.Error())
		} else if v != 73 {
			t.Fatal("Expected 73 but got", v)
		}

		//Test wether more bytes were read then needed.
		buf := make([]byte, 11)
		n, err := io.ReadFull(rd, buf)
		if n != 10 {
			t.Fatal(err.Error())
		}
	}

	{
		var input varint.Int62 = 1073741823
		b, err := varint.Int62ToVarint(input)
		if err != nil {
			t.Fatal(err.Error())
		}

		//Input is a minimum of 4 bytes
		//Take first two bytes from it.
		b = b[:3]
		rd := bufio.NewReader(bytes.NewReader(b))
		if _, err := varint.ReadVarint62(rd); err == nil {
			t.Fatal("Expected a error but got no error")
		}
	}
}
