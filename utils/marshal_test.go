package utils_test

import (
	"fmt"
	"testing"

	"github.com/datewu/gtea/utils"
)

type anCache []string

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (l anCache) MarshalBinary() ([]byte, error) {
	type a anCache
	a1 := a(l)
	return utils.RedisMarshal(a1)
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (l *anCache) UnmarshalBinary(data []byte) error {
	type a anCache
	a1 := (*a)(l)
	return utils.RedisUnMarshal(data, a1)
}

func ExampleRedisMarshal() {
	data := anCache{"hello", "world"}
	bs, err := data.MarshalBinary()
	// before go1.21.0
	// [15 255 129 2 1 1 1 97 1 255 130 0 1 12 0 0 16 255 130 0 2 5 104 101 108 108 111 5 119 111 114 108 100]
	if err != nil {
		panic(err)
	}
	fmt.Println(bs)

	// Output:
	// [14 127 2 1 1 1 97 1 255 128 0 1 12 0 0 16 255 128 0 2 5 104 101 108 108 111 5 119 111 114 108 100]

}

func TestUnMarshal(t *testing.T) {
	data := anCache{"hello", "world"}
	bs, err := data.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	a := &anCache{}
	err = a.UnmarshalBinary(bs)
	if err != nil {
		t.Fatal(err)
	}
	for i, v := range data {
		if v != (*a)[i] {
			t.Fatal("not match", v, (*a)[i])
		}
	}
}
