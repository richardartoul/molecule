package molecule

import (
	"fmt"
	"io"
	"math"

	"github.com/richardartoul/molecule/src/codec"
)

type MessageEachFn func(fieldNum int32, value Value) error

func MessageEach(b []byte, fn MessageEachFn) error {
	buffer := codec.NewBuffer(b)
	for {
		fieldNum, wireType, err := buffer.DecodeTagAndWireType()
		if err == io.EOF || fieldNum == 0 {
			return nil
		}

		value := Value{
			WireType: wireType,
		}

		switch wireType {
		case codec.WireVarint:
			varint, err := buffer.DecodeVarint()
			if err != nil {
				return fmt.Errorf(
					"MessageEach: error decoding varint for fieldNum: %d, err: %v", fieldNum, err)
			}
			value.Number = varint
		case codec.WireFixed32:
			fixed32, err := buffer.DecodeFixed32()
			if err != nil {
				return fmt.Errorf(
					"MessageEach: error decoding fixed32 for fieldNum: %d, err: %v", fieldNum, err)
			}
			value.Number = fixed32
		case codec.WireFixed64:
			fixed64, err := buffer.DecodeFixed64()
			if err != nil {
				return fmt.Errorf(
					"MessageEach: error decoding fixed64 for fieldNum: %d, err: %v", fieldNum, err)
			}
			value.Number = fixed64
		case codec.WireBytes:
			b, err := buffer.DecodeRawBytes(false)
			if err != nil {
				return fmt.Errorf(
					"MessageEach: error decoding raw bytes for fieldNum: %d, err: %v", fieldNum, err)
			}
			value.Bytes = b
		case codec.WireStartGroup, codec.WireEndGroup:
			return fmt.Errorf(
				"MessageEach: encountered group wire type: %d for fieldNum: %d. Groups not supported",
				wireType, fieldNum)
		default:
			return fmt.Errorf(
				"MessageEach: unknown wireType: %d for fieldNum: %d",
				wireType, fieldNum)
		}

		if err := fn(fieldNum, value); err != nil {
			return err
		}
	}
}

// fieldNum := fd.GetNumber()
// 	val := unmarshalValue{fieldNumber: fieldNum, v: v}
// 	switch fd.GetType() {
// 	case dpb.FieldDescriptorProto_TYPE_BOOL,
// 		dpb.FieldDescriptorProto_TYPE_UINT64,
// 		dpb.FieldDescriptorProto_TYPE_FIXED64,
// 		dpb.FieldDescriptorProto_TYPE_INT64,
// 		dpb.FieldDescriptorProto_TYPE_SFIXED64,
// 		dpb.FieldDescriptorProto_TYPE_DOUBLE:
// 		return val, nil

// 	case dpb.FieldDescriptorProto_TYPE_UINT32,
// 		dpb.FieldDescriptorProto_TYPE_FIXED32:
// 		if v > math.MaxUint32 {
// 			return zeroValue, fmt.Errorf("%d (field num %d) overflows uint32", v, fieldNum)
// 		}
// 		return val, nil

// 	case dpb.FieldDescriptorProto_TYPE_INT32,
// 		dpb.FieldDescriptorProto_TYPE_ENUM:
// 		s := int64(v)
// 		if s > math.MaxInt32 {
// 			return zeroValue, fmt.Errorf("%d (field num %d) overflows int32", v, fieldNum)
// 		}
// 		if s < math.MinInt32 {
// 			return zeroValue, fmt.Errorf("%d (field num %d) underflows int32", v, fieldNum)
// 		}
// 		return val, nil

// 	case dpb.FieldDescriptorProto_TYPE_SFIXED32:
// 		if v > math.MaxUint32 {
// 			return zeroValue, fmt.Errorf("%d (field num %d) overflows int32", v, fieldNum)
// 		}
// 		return val, nil

// 	case dpb.FieldDescriptorProto_TYPE_SINT32:
// 		if v > math.MaxUint32 {
// 			return zeroValue, fmt.Errorf("%d (field num %d) overflows int32", v, fieldNum)
// 		}
// 		val.v = uint64(decodeZigZag32(v))
// 		return val, nil

// 	case dpb.FieldDescriptorProto_TYPE_SINT64:
// 		val.v = uint64(decodeZigZag64(v))
// 		return val, nil

// 	case dpb.FieldDescriptorProto_TYPE_FLOAT:
// 		if v > math.MaxUint32 {
// 			return zeroValue, fmt.Errorf("%d (field num %d) overflows uint32", v, fieldNum)
// 		}
// 		float32Val := math.Float32frombits(uint32(v))
// 		float64Bits := math.Float64bits(float64(float32Val))
// 		val.v = float64Bits
// 		return val, nil

// 	default:
// 		// bytes, string, message, and group cannot be represented as a simple numeric value.
// 		return zeroValue, fmt.Errorf("bad input; field %s requires length-delimited wire type", fd.GetFullyQualifiedName())
// 	}

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

func (v *Value) AsDouble() (float64, error) {
	return math.Float64frombits(v.Number), nil
}

func (v *Value) AsFloat() (float32, error) {
	if v.Number > math.MaxUint32 {
		return 0, fmt.Errorf("AsFloat: %d overflows float32", v.Number)
	}
	return math.Float32frombits(uint32(v.Number)), nil
}

func (v *Value) AsInt32() (int32, error) {
	if v.Number > math.MaxUint32 {
		return 0, fmt.Errorf("AsInt32: %d overflows int32", v.Number)
	}
	return int32(v.Number), nil
}

func (v *Value) AsInt64() (int64, error) {
	return int64(v.Number), nil
}

func (v *Value) AsUint32() (uint32, error) {
	if v.Number > math.MaxUint32 {
		return 0, fmt.Errorf("AsUInt32: %d overflows uint32", v.Number)
	}
	return uint32(v.Number), nil
}

func (v *Value) AsUint64() (uint64, error) {
	return v.Number, nil
}

func (v *Value) AsSint32() (int32, error) {
	if v.Number > math.MaxUint32 {
		return 0, fmt.Errorf("AsSint32: %d overflows int32", v.Number)
	}
	return codec.DecodeZigZag32(v.Number), nil
}

func (v *Value) AsSint64() (int64, error) {
	return codec.DecodeZigZag64(v.Number), nil
}

func (v *Value) AsFixed32() (uint32, error) {
	if v.Number > math.MaxUint32 {
		return 0, fmt.Errorf("AsFixed32: %d overflows int32", v.Number)
	}
	return uint32(v.Number), nil
}

func (v *Value) AsFixed64() (uint64, error) {
	return uint64(v.Number), nil
}

func (v *Value) AsSFixed32() (int32, error) {
	if v.Number > math.MaxUint32 {
		return 0, fmt.Errorf("AsSFixed32: %d overflows int32", v.Number)
	}
	return int32(v.Number), nil
}

func (v *Value) AsSFixed64() (int64, error) {
	return int64(v.Number), nil
}

func (v *Value) AsBool() (bool, error) {
	return v.Number == 1, nil
}

func (v *Value) AsString() (string, error) {
	// TODO: Do unsafe conversion here.
	return string(v.Bytes), nil
}

func (v *Value) AsBytes() ([]byte, error) {
	return v.Bytes, nil
}
