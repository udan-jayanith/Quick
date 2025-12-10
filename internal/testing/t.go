package Testing

import (
	"encoding/json"
	"fmt"

	"github.com/udan-jayanith/Quick/varint"
)

func ToFormattedJson(v any) string {
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		panic("ToFormattedJson function in main_test.go failed.")
	}
	return string(b)
}

func Int62ToVarint(v varint.Int62) []byte {
	b, err := varint.Int62ToVarint(v)
	if err != nil {
		panic(fmt.Sprintf("Int62ToVarint retuned a error\n%s\nFrom internal/Testing pkg.", err.Error()))
	}
	return b
}
