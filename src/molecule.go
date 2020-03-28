package molecule

import (
	"fmt"
	"io"
	"math"

	"github.com/richardartoul/molecule/src/codec"
)

// MessageEachFn is a function that will be called for each top-level field in a
// message passed to MessageEach.
type MessageEachFn func(fieldNum int32, value Value) error

// MessageEach iterates over each top-level field in message b and calls fn on
// each one.
func MessageEach(buffer *codec.Buffer, fn MessageEachFn) error {
	for !buffer.EOF() {
		fieldNum, wireType, err := buffer.DecodeTagAndWireType()
		if err == io.EOF {
			return nil
		}

		value, err := readValueFromBuffer(wireType, buffer)
		if err != nil {
			return fmt.Errorf("MessageEach: error reading value from buffer: %v", err)
		}

		if err := fn(fieldNum, value); err != nil {
			return err
		}
	}
	return nil
}

// PackedRepeatedEachFn is a function that is called for each value in a repeated field.
type PackedRepeatedEachFn func(value Value) error

// PackedArrayEach iterates over each value in packed repeated field b and calls fn on
// each one.
//
// PackedArrayEach only supports repeated fields encoded using packed encoding.
func PackedArrayEach(buffer *codec.Buffer, wireType int8, fn PackedRepeatedEachFn) error {
	for !buffer.EOF() {
		value, err := readValueFromBuffer(wireType, buffer)
		if err != nil {
			return fmt.Errorf("ArrayEach: error reading value from buffer: %v", err)
		}
		if err := fn(value); err != nil {
			return nil
		}
	}

	return nil
}

func readValueFromBuffer(wireType int8, buffer *codec.Buffer) (Value, error) {
	value := Value{
		WireType: wireType,
	}

	switch wireType {
	case codec.WireVarint:
		varint, err := buffer.DecodeVarint()
		if err != nil {
			return Value{}, fmt.Errorf(
				"MessageEach: error decoding varint: %v", err)
		}
		value.Number = varint
	case codec.WireFixed32:
		fixed32, err := buffer.DecodeFixed32()
		if err != nil {
			return Value{}, fmt.Errorf(
				"MessageEach: error decoding fixed32: %v", err)
		}
		value.Number = fixed32
	case codec.WireFixed64:
		fixed64, err := buffer.DecodeFixed64()
		if err != nil {
			return Value{}, fmt.Errorf(
				"MessageEach: error decoding fixed64: %v", err)
		}
		value.Number = fixed64
	case codec.WireBytes:
		b, err := buffer.DecodeRawBytes(false)
		if err != nil {
			return Value{}, fmt.Errorf(
				"MessageEach: error decoding raw bytes: %v", err)
		}
		value.Bytes = b
	case codec.WireStartGroup, codec.WireEndGroup:
		return Value{}, fmt.Errorf(
			"MessageEach: encountered group wire type: %d. Groups not supported",
			wireType)
	default:
		return Value{}, fmt.Errorf(
			"MessageEach: unknown wireType: %d", wireType)
	}

	return value, nil
}

// Value represents a protobuf value. It contains the original wiretype that the value
// was encoded with as well as a variety of helper methods for interpreting the raw
// value based on the field's actual type.
type Value struct {
	// WireType is the protobuf wire type that was used to encode the field.
	WireType int8
	// Number will contain the value for any fields encoded with the
	// following wire types:
	//
	// 1. varint
	// 2. Fixed32
	// 3. Fixed64
	Number uint64
	// BytesValue will contain the value for any fields encoded with the
	// following wire types:
	//
	// 1. bytes
	Bytes []byte
}

// AsDouble interprets the value as a double.
func (v *Value) AsDouble() (float64, error) {
	return math.Float64frombits(v.Number), nil
}

// AsFloat interprets the value as a float.
func (v *Value) AsFloat() (float32, error) {
	if v.Number > math.MaxUint32 {
		return 0, fmt.Errorf("AsFloat: %d overflows float32", v.Number)
	}
	return math.Float32frombits(uint32(v.Number)), nil
}

// AsInt32 interprets the value as an int32.
func (v *Value) AsInt32() (int32, error) {
	s := int64(v.Number)
	if s > math.MaxInt32 {
		return 0, fmt.Errorf("AsInt32: %d overflows int32", s)
	}
	if s < math.MinInt32 {
		return 0, fmt.Errorf("AsInt32: %d underflows int32", s)
	}
	return int32(v.Number), nil
}

// AsInt64 interprets the value as an int64.
func (v *Value) AsInt64() (int64, error) {
	return int64(v.Number), nil
}

// AsUint32 interprets the value as a uint32.
func (v *Value) AsUint32() (uint32, error) {
	if v.Number > math.MaxUint32 {
		return 0, fmt.Errorf("AsUInt32: %d overflows uint32", v.Number)
	}
	return uint32(v.Number), nil
}

// AsUint64 interprets the value as a uint64.
func (v *Value) AsUint64() (uint64, error) {
	return v.Number, nil
}

// AsSint32 interprets the value as a sint32.
func (v *Value) AsSint32() (int32, error) {
	if v.Number > math.MaxUint32 {
		return 0, fmt.Errorf("AsSint32: %d overflows int32", v.Number)
	}
	return codec.DecodeZigZag32(v.Number), nil
}

// AsSint64 interprets the value as a sint64.
func (v *Value) AsSint64() (int64, error) {
	return codec.DecodeZigZag64(v.Number), nil
}

// AsFixed32 interprets the value as a fixed32.
func (v *Value) AsFixed32() (uint32, error) {
	if v.Number > math.MaxUint32 {
		return 0, fmt.Errorf("AsFixed32: %d overflows int32", v.Number)
	}
	return uint32(v.Number), nil
}

// AsFixed64 interprets the value as a fixed64.
func (v *Value) AsFixed64() (uint64, error) {
	return uint64(v.Number), nil
}

// AsSFixed32 interprets the value as a SFixed32.
func (v *Value) AsSFixed32() (int32, error) {
	if v.Number > math.MaxUint32 {
		return 0, fmt.Errorf("AsSFixed32: %d overflows int32", v.Number)
	}
	return int32(v.Number), nil
}

// AsSFixed64 interprets the value as a SFixed64.
func (v *Value) AsSFixed64() (int64, error) {
	return int64(v.Number), nil
}

// AsBool interprets the value as a bool.
func (v *Value) AsBool() (bool, error) {
	return v.Number == 1, nil
}

// AsString interprets the value as a string.
func (v *Value) AsString() (string, error) {
	// TODO: Do unsafe conversion here.
	return string(v.Bytes), nil
}

// AsBytes interprets the value as bytes.
func (v *Value) AsBytes() ([]byte, error) {
	return v.Bytes, nil
}
