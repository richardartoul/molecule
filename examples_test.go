package molecule

import (
	"fmt"

	"github.com/richardartoul/molecule/src/codec"
	simple "github.com/richardartoul/molecule/src/proto"

	"github.com/golang/protobuf/proto"
)

// Example demonstrates how the molecule library can be used to parse a protobuf message.
func Example() {
	// Proto definitions:
	//
	//   message Test {
	//     string string_field = 1;
	//     int64 int64_field = 2;
	//     repeated int64 repeated_int64_field = 3;
	//   }

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

// Example_nested demonstrates how to use the MessageEach function to
// decode a nested message.
func Example_nested() {
	// Proto definitions:
	//
	//   message Test {
	//       string string_field = 1;
	//       int64 int64_field = 2;
	//       repeated int64 repeated_int64_field = 3;
	//   }
	//
	//   message Nested {
	//       Test nested_message = 1;
	//   }

	var (
		test   = &simple.Test{StringField: "Hello world!"}
		nested = &simple.Nested{NestedMessage: test}
	)
	marshaled, err := proto.Marshal(nested)
	if err != nil {
		panic(err)
	}

	var (
		buffer = codec.NewBuffer(marshaled)
		strVal Value
	)
	err = MessageEach(buffer, func(fieldNum int32, value Value) bool {
		if fieldNum == 1 {
			packedArr, err := value.AsBytesUnsafe()
			if err != nil {
				panic(err)
			}

			buffer := codec.NewBuffer(packedArr)
			err = MessageEach(buffer, func(fieldNum int32, value Value) bool {
				if fieldNum == 1 {
					strVal = value
				}
				// Found it, stop scanning.
				return false
			})
			if err != nil {
				panic(err)
			}

			// Found it, stop scanning.
			return false
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

	fmt.Println("NestedMessage.StringField:", str)

	// Output:
	// NestedMessage.StringField: Hello world!
}

// Example_repeated demonstrates how to use the PackedRepeatedEach function to
// decode a repeated field encoded in the packed (proto 3) format.
func Example_repeated() {
	// Proto definitions:
	//
	//   message Test {
	//     string string_field = 1;
	//     int64 int64_field = 2;
	//     repeated int64 repeated_int64_field = 3;
	//   }

	int64s := []int64{1, 2, 3, 4, 5, 6, 7}
	m := &simple.Test{RepeatedInt64Field: int64s}
	marshaled, err := proto.Marshal(m)
	if err != nil {
		panic(err)
	}

	var (
		buffer          = codec.NewBuffer(marshaled)
		unmarshaledInts = []int64{}
	)
	err = MessageEach(buffer, func(fieldNum int32, value Value) bool {
		if fieldNum == 3 {
			packedArr, err := value.AsBytesUnsafe()
			if err != nil {
				panic(err)
			}

			buffer := codec.NewBuffer(packedArr)
			PackedRepeatedEach(buffer, codec.FieldType_INT64, func(v Value) bool {
				vInt64, err := v.AsInt64()
				if err != nil {
					panic(err)
				}
				unmarshaledInts = append(unmarshaledInts, vInt64)
				return true
			})

			// Found it, stop scanning.
			return false
		}
		// Continue scanning.
		return true
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("Int64s:", unmarshaledInts)

	// Output:
	// Int64s: [1 2 3 4 5 6 7]
}

// ExampleMessageEach_SelectAField desmonates how the MessageEach function can
// be used to select an individual field.
func ExampleMessageEach_selectAField() {
	// Proto definitions:
	//
	//   message Test {
	//     string string_field = 1;
	//     int64 int64_field = 2;
	//     repeated int64 repeated_int64_field = 3;
	//   }

	m := &simple.Test{StringField: "hello world!"}
	marshaled, err := proto.Marshal(m)
	if err != nil {
		panic(err)
	}

	var (
		buffer = codec.NewBuffer(marshaled)
		strVal Value
	)
	err = MessageEach(buffer, func(fieldNum int32, value Value) bool {
		if fieldNum == 1 {
			strVal = value
			// Found it, stop scanning.
			return false
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

	fmt.Println("StringField:", str)

	// Output:
	// StringField: hello world!
}

// ExamplePackedRepeatedEach demonstrates how to use the PackedRepeatedEach function to
// decode a repeated field encoded in the packed (proto 3) format.
func ExamplePackedRepeatedEach() {
	// Proto definitions:
	//
	//   message Test {
	//     string string_field = 1;
	//     int64 int64_field = 2;
	//     repeated int64 repeated_int64_field = 3;
	//   }

	int64s := []int64{1, 2, 3, 4, 5, 6, 7}
	m := &simple.Test{RepeatedInt64Field: int64s}
	marshaled, err := proto.Marshal(m)
	if err != nil {
		panic(err)
	}

	var (
		buffer          = codec.NewBuffer(marshaled)
		unmarshaledInts = []int64{}
	)
	err = MessageEach(buffer, func(fieldNum int32, value Value) bool {
		if fieldNum == 3 {
			packedArr, err := value.AsBytesUnsafe()
			if err != nil {
				panic(err)
			}

			buffer := codec.NewBuffer(packedArr)
			err = PackedRepeatedEach(buffer, codec.FieldType_INT64, func(v Value) bool {
				vInt64, err := v.AsInt64()
				if err != nil {
					panic(err)
				}
				unmarshaledInts = append(unmarshaledInts, vInt64)
				return true
			})
			if err != nil {
				panic(err)
			}

			// Found it, stop scanning.
			return false
		}
		// Continue scanning.
		return true
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("Int64s:", unmarshaledInts)

	// Output:
	// Int64s: [1 2 3 4 5 6 7]
}
