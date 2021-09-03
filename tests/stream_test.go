package moleculetest

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/richardartoul/molecule"
	testproto "github.com/richardartoul/molecule/src/proto"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
)

// values copied from the .proto file
const fieldDub int = 1
const fieldFlo int = 2
const fieldI32 = 3
const fieldI64 = 4
const fieldU32 = 5
const fieldU64 = 6
const fieldS32 = 7
const fieldS64 = 8
const fieldF32 = 9
const fieldF64 = 10
const fieldSF32 = 11
const fieldSF64 = 12
const fieldBoo = 13
const fieldStr = 14
const fieldByt = 15
const fieldRepI64Packed = 16

// Test that a simple message encoding properly decodes using the generated
// protofbuf code.  This function should test all proto types.
func TestSimpleEncoding(t *testing.T) {
	output := bytes.NewBuffer([]byte{})
	ps := molecule.NewProtoStream()
	ps.Reset(output)

	require.NoError(t, ps.Double(fieldDub, 3.14))
	require.NoError(t, ps.Float(fieldFlo, 3.14))
	require.NoError(t, ps.Int32(fieldI32, int32(-10)))
	require.NoError(t, ps.Int64(fieldI64, int64(-11)))
	require.NoError(t, ps.Uint32(fieldU32, uint32(12)))
	require.NoError(t, ps.Uint64(fieldU64, uint64(13)))
	require.NoError(t, ps.Sint32(fieldS32, int32(-12)))
	require.NoError(t, ps.Sint64(fieldS64, int64(-13)))
	require.NoError(t, ps.Fixed32(fieldF32, uint32(22)))
	require.NoError(t, ps.Fixed64(fieldF64, uint64(23)))
	require.NoError(t, ps.Sfixed32(fieldSF32, int32(-22)))
	require.NoError(t, ps.Sfixed64(fieldSF64, int64(-23)))
	require.NoError(t, ps.Bool(fieldBoo, true))
	require.NoError(t, ps.String(fieldStr, "s"))
	require.NoError(t, ps.Bytes(fieldByt, []byte("b")))
	require.NoError(t, ps.Int64Packed(fieldRepI64Packed, []int64{-1, 2, 3}))

	buf := output.Bytes()
	var res testproto.Simple

	require.NoError(t, proto.Unmarshal(buf, &res))

	assert.Equal(t, float64(3.14), res.Double)
	assert.Equal(t, float32(3.14), res.Float)
	assert.Equal(t, int32(-10), res.Int32)
	assert.Equal(t, int64(-11), res.Int64)
	assert.Equal(t, uint32(12), res.Uint32)
	assert.Equal(t, uint64(13), res.Uint64)
	assert.Equal(t, int32(-12), res.Sint32)
	assert.Equal(t, int64(-13), res.Sint64)
	assert.Equal(t, uint32(22), res.Fixed32)
	assert.Equal(t, uint64(23), res.Fixed64)
	assert.Equal(t, int32(-22), res.Sfixed32)
	assert.Equal(t, int64(-23), res.Sfixed64)
	assert.Equal(t, true, res.Bool)
	assert.Equal(t, "s", res.String_)
	assert.Equal(t, "b", string(res.Bytes))
	assert.Equal(t, int64(-1), res.RepeatedInt64Packed[0])
	assert.Equal(t, int64(2), res.RepeatedInt64Packed[1])
	assert.Equal(t, int64(3), res.RepeatedInt64Packed[2])
}

// Test that the zero values for each field do not result in any encoded data.
func TestSimpleEncodingZero(t *testing.T) {
	output := bytes.NewBuffer([]byte{})
	ps := molecule.NewProtoStream()
	ps.Reset(output)

	require.NoError(t, ps.Double(fieldDub, 0.0))
	require.NoError(t, ps.Float(fieldFlo, 0.0))
	require.NoError(t, ps.Int32(fieldI32, int32(0)))
	require.NoError(t, ps.Int64(fieldI64, int64(0)))
	require.NoError(t, ps.Uint32(fieldU32, uint32(0)))
	require.NoError(t, ps.Uint64(fieldU64, uint64(0)))
	require.NoError(t, ps.Sint32(fieldS32, int32(0)))
	require.NoError(t, ps.Sint64(fieldS64, int64(0)))
	require.NoError(t, ps.Fixed32(fieldF32, uint32(0)))
	require.NoError(t, ps.Fixed64(fieldF64, uint64(0)))
	require.NoError(t, ps.Sfixed32(fieldSF32, int32(0)))
	require.NoError(t, ps.Sfixed64(fieldSF64, int64(0)))
	require.NoError(t, ps.Bool(fieldBoo, false))
	require.NoError(t, ps.String(fieldStr, ""))
	require.NoError(t, ps.Bytes(fieldByt, []byte("")))

	buf := output.Bytes()
	// all of those are zero values, so nothing should have been written
	require.Equal(t, 0, len(buf))
}

// Test that the zero values for packed fields do not result in any encoded data.
func TestPackedEncodingZero(t *testing.T) {
	output := bytes.NewBuffer([]byte{})
	ps := molecule.NewProtoStream()
	ps.Reset(output)

	require.NoError(t, ps.DoublePacked(fieldDub, []float64{}))
	require.NoError(t, ps.FloatPacked(fieldFlo, []float32{}))
	require.NoError(t, ps.Int32Packed(fieldI32, []int32{}))
	require.NoError(t, ps.Int64Packed(fieldI64, []int64{}))
	require.NoError(t, ps.Uint32Packed(fieldU32, []uint32{}))
	require.NoError(t, ps.Uint64Packed(fieldU64, []uint64{}))
	require.NoError(t, ps.Sint32Packed(fieldS32, []int32{}))
	require.NoError(t, ps.Sint64Packed(fieldS64, []int64{}))
	require.NoError(t, ps.Fixed32Packed(fieldF32, []uint32{}))
	require.NoError(t, ps.Fixed64Packed(fieldF64, []uint64{}))
	require.NoError(t, ps.Sfixed32Packed(fieldSF32, []int32{}))
	require.NoError(t, ps.Sfixed64Packed(fieldSF64, []int64{}))

	buf := output.Bytes()
	// all of those are zero values, so nothing should have been written
	require.Equal(t, 0, len(buf))
}

// Test that *Packed functions work.
func TestPacking(t *testing.T) {
	output := bytes.NewBuffer([]byte{})
	ps := molecule.NewProtoStream()
	ps.Reset(output)

	assertBytes := func(t *testing.T, exp []byte, got []byte) {
		assert.Equal(t, fmt.Sprintf("%#v", exp), fmt.Sprintf("%#v", got))
	}

	key := func(fieldNumber int) uint8 {
		const wtLengthDelimited = 2
		return uint8((fieldNumber << 3) + wtLengthDelimited)
	}

	t.Run("DoublePacked", func(t *testing.T) {
		output.Reset()
		require.NoError(t, ps.DoublePacked(fieldDub, []float64{3.14, 1.414}))
		assertBytes(t, []byte{
			key(fieldDub),
			0x10,                                          // length
			0x1f, 0x85, 0xeb, 0x51, 0xb8, 0x1e, 0x9, 0x40, // 3.14 as fixed64
			0x39, 0xb4, 0xc8, 0x76, 0xbe, 0x9f, 0xf6, 0x3f, // 1.414 as fixed64
		}, output.Bytes())
	})

	t.Run("FloatPacked", func(t *testing.T) {
		output.Reset()
		require.NoError(t, ps.FloatPacked(fieldFlo, []float32{3.14, 1.414}))
		assertBytes(t, []byte{
			key(fieldFlo),
			0x8,                    // length
			0xc3, 0xf5, 0x48, 0x40, // 3.14 as fixed32
			0xf4, 0xfd, 0xb4, 0x3f, // 1.414 as fixed32
		}, output.Bytes())
	})

	t.Run("Int32Packed", func(t *testing.T) {
		output.Reset()
		require.NoError(t, ps.Int32Packed(fieldI32, []int32{int32(-12), int32(12), int32(-13)}))
		assertBytes(t, []byte{
			key(fieldI32),
			0x15,                                                       // length
			0xf4, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01, // -12 (little-endian varint)
			0x0c,                                                       // 12
			0xf3, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01, // -13
		}, output.Bytes())
	})

	t.Run("Int64Packed", func(t *testing.T) {
		output.Reset()
		require.NoError(t, ps.Int64Packed(fieldI64, []int64{int64(-12), int64(12), int64(-13)}))
		assertBytes(t, []byte{
			key(fieldI64),
			0x15,                                                      // length
			0xf4, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x1, // -12
			0x0c,                                                      // 12
			0xf3, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x1, // -13
		}, output.Bytes())
	})

	t.Run("Uint32Packed", func(t *testing.T) {
		output.Reset()
		require.NoError(t, ps.Uint32Packed(fieldU32, []uint32{uint32(1), uint32(2), uint32(3)}))
		assertBytes(t, []byte{
			key(fieldU32),
			0x3, // length
			0x1,
			0x2,
			0x3,
		}, output.Bytes())
	})

	t.Run("Uint64Packed", func(t *testing.T) {
		output.Reset()
		require.NoError(t, ps.Uint64Packed(fieldU64, []uint64{uint64(1), uint64(2), uint64(3)}))
		assertBytes(t, []byte{
			key(fieldU64),
			0x3, // length
			0x1,
			0x2,
			0x3,
		}, output.Bytes())
	})

	t.Run("Sint32Packed", func(t *testing.T) {
		output.Reset()
		require.NoError(t, ps.Sint32Packed(fieldS32, []int32{int32(-12), int32(12), int32(-13)}))
		assertBytes(t, []byte{
			key(fieldS32),
			0x03, // length
			0x17, // zigzag encoding of -12
			0x18, // 12
			0x19, // 13
		}, output.Bytes())
	})

	t.Run("Sint64Packed", func(t *testing.T) {
		output.Reset()
		require.NoError(t, ps.Sint64Packed(fieldS64, []int64{int64(-12), int64(12), int64(-13)}))
		assertBytes(t, []byte{
			key(fieldS64),
			0x03, // length
			0x17,
			0x18,
			0x19,
		}, output.Bytes())
	})

	t.Run("Fixed32Packed", func(t *testing.T) {
		output.Reset()
		require.NoError(t, ps.Fixed32Packed(fieldF32, []uint32{uint32(12), uint32(13), uint32(14)}))
		assertBytes(t, []byte{
			key(fieldF32),
			0xc,                // length
			0xc, 0x0, 0x0, 0x0, // 12
			0xd, 0x0, 0x0, 0x0, // 13
			0xe, 0x0, 0x0, 0x0, // 14
		}, output.Bytes())
	})

	t.Run("Fixed64Packed", func(t *testing.T) {
		output.Reset()
		require.NoError(t, ps.Fixed64Packed(fieldF64, []uint64{uint64(12), uint64(13), uint64(14)}))
		assertBytes(t, []byte{
			key(fieldF64),
			0x18,                                   // length
			0xc, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, // 12
			0xd, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, // 13
			0xe, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, // 14
		}, output.Bytes())
	})

	t.Run("Sfixed32Packed", func(t *testing.T) {
		output.Reset()
		require.NoError(t, ps.Sfixed32Packed(fieldSF32, []int32{int32(12), int32(-13), int32(14)}))
		assertBytes(t, []byte{
			key(fieldSF32),
			0xc,                // length
			0xc, 0x0, 0x0, 0x0, // 12
			0xf3, 0xff, 0xff, 0xff, // -13
			0xe, 0x0, 0x0, 0x0, // 14
		}, output.Bytes())
	})

	t.Run("Sfixed64Packed", func(t *testing.T) {
		output.Reset()
		require.NoError(t, ps.Sfixed64Packed(fieldSF64, []int64{int64(12), int64(-13), int64(14)}))
		assertBytes(t, []byte{
			key(fieldSF64),
			0x18,                                   // length
			0xc, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, // 12
			0xf3, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, // -13
			0xe, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, // 14
		}, output.Bytes())
	})
}

// Microbenchmark simple encoding performance
func BenchmarkSimple(b *testing.B) {
	output := bytes.NewBuffer([]byte{})
	ps := molecule.NewProtoStream()
	ps.Reset(output)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		output.Reset()
		ps.Double(fieldDub, 3.14)
		ps.Float(fieldFlo, 3.14)
		ps.Int32(fieldI32, int32(-10))
		ps.Int64(fieldI64, int64(-11))
		ps.Uint32(fieldU32, uint32(12))
		ps.Uint64(fieldU64, uint64(13))
		ps.Sint32(fieldS32, int32(-12))
		ps.Sint64(fieldS64, int64(-13))
		ps.Fixed32(fieldF32, uint32(22))
		ps.Fixed64(fieldF64, uint64(23))
		ps.Sfixed32(fieldSF32, int32(-22))
		ps.Sfixed64(fieldSF64, int64(-23))
		ps.Bool(fieldBoo, true)
		ps.String(fieldStr, "s")
		ps.Bytes(fieldByt, []byte("b"))
		ps.Int64Packed(fieldRepI64Packed, []int64{99, -99})
	}
}

// Microbenchmark packing performance
func BenchmarkPacking(b *testing.B) {
	output := bytes.NewBuffer([]byte{})
	ps := molecule.NewProtoStream()
	ps.Reset(output)

	const packSize = 102400
	floats := make([]float64, 0, packSize)
	for i := 0; i < packSize; i++ {
		// note that nothing in the implementation optimizes for
		// repeated values, so there's no need to do anything more
		// interesting than this.
		floats = append(floats, 2.7182818284)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		output.Reset()
		err := ps.DoublePacked(fieldDub, floats)
		if err != nil {
			panic(err)
		}
	}
}

// Test ps.Embedded embedding a repeated message
func TestEmbedding(t *testing.T) {
	output := bytes.NewBuffer([]byte{})
	ps := molecule.NewProtoStream()
	ps.Reset(output)

	// values copied from the .proto file
	const fieldNestedMessage = 1

	makeMsg := func(str string) func(*molecule.ProtoStream) error {
		return func(ps *molecule.ProtoStream) error {
			// values copied from the .proto file
			const fieldStringField = 1
			const fieldInt64Field = 2
			const fieldRepeatedInt64Field = 3
			var err error

			err = ps.String(fieldStringField, str)
			if err != nil {
				return err
			}
			err = ps.Int64(fieldInt64Field, int64(len(str)))
			if err != nil {
				return err
			}
			chars := make([]int64, len(str))
			for i := range chars {
				chars[i] = int64(str[i])
			}
			err = ps.Int64Packed(fieldRepeatedInt64Field, chars)
			if err != nil {
				return err
			}
			return nil
		}
	}

	require.NoError(t, ps.Embedded(fieldNestedMessage, makeMsg("hello")))

	buf := output.Bytes()
	var res testproto.Nested

	require.NoError(t, proto.Unmarshal(buf, &res))

	require.Equal(t, "hello", res.NestedMessage.StringField)
	require.Equal(t, int64(5), res.NestedMessage.Int64Field)
	require.Equal(t, []int64{104, 101, 108, 108, 111}, res.NestedMessage.RepeatedInt64Field)
}
