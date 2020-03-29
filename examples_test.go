package molecule

import (
	"fmt"

	"github.com/richardartoul/molecule/src/codec"
	simple "github.com/richardartoul/molecule/src/proto"

	"github.com/golang/protobuf/proto"
)

// Example demonstrates how the molecule library can be used to parse a protobuf message.
func Example() {
	m := &simple.Test{
		StringField: "hello world!",
		Int64Field:  10,
	}
	marshaled, err := proto.Marshal(m)
	if err != nil {
		panic(err)
	}

	var (
		buffer   = codec.NewBuffer(marshaled)
		strVal   Value
		int64Val Value
	)
	err = MessageEach(buffer, func(fieldNum int32, value Value) bool {
		if fieldNum == 1 {
			strVal = value
		}
		if fieldNum == 2 {
			int64Val = value
		}

		// Continue scanning.
		return true
	})
	if err != nil {
		panic(err)
	}

	str, err := strVal.AsStringUnsafe()
	if err != nil {
		panic(err)
	}
	int64V, err := int64Val.AsInt64()
	if err != nil {
		panic(err)
	}

	fmt.Println("StringField:", str)
	fmt.Println("Int64Field:", int64V)

	// Output:
	// StringField: hello world!
	// Int64Field: 10
}
