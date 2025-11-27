package Testing

import (
	"encoding/json"
)

func ToFormattedJson(v any) string {
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		panic("ToFormattedJson function in main_test.go failed.")
	}
	return string(b)
}
