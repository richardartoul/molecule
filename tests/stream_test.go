package moleculetest

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	fuzz "github.com/google/gofuzz"
	"github.com/richardartoul/molecule"
	simple "github.com/richardartoul/molecule/src/proto"
	testproto "github.com/richardartoul/molecule/src/proto"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"
)

// values copied from the .proto file
const fieldDouble int = 1
const fieldFloat int = 2
const fieldInt32 = 3
const fieldInt64 = 4
const fieldUint32 = 5
const fieldUint64 = 6
const fieldSint32 = 7
const fieldSint64 = 8
const fieldFixed32 = 9
const fieldFixed64 = 10
const fieldSfixed32 = 11
const fieldSfixed64 = 12
const fieldBool = 13
const fieldString = 14
const fieldBytes = 15
const fieldRepeatedInt64Packed = 16

// Test that a simple message encoding properly decodes using the generated
// protofbuf code.  This function should test all proto types.
func TestSimpleEncoding(t *testing.T) {
	output := bytes.NewBuffer([]byte{})
	ps := molecule.NewProtoStream(output)

	require.NoError(t, ps.Double(fieldDouble, 3.14))
	require.NoError(t, ps.Float(fieldFloat, 3.14))
	require.NoError(t, ps.Int32(fieldInt32, int32(-10)))
	require.NoError(t, ps.Int64(fieldInt64, int64(-11)))
	require.NoError(t, ps.Uint32(fieldUint32, uint32(12)))
	require.NoError(t, ps.Uint64(fieldUint64, uint64(13)))
	require.NoError(t, ps.Sint32(fieldSint32, int32(-12)))
	require.NoError(t, ps.Sint64(fieldSint64, int64(-13)))
	require.NoError(t, ps.Fixed32(fieldFixed32, uint32(22)))
	require.NoError(t, ps.Fixed64(fieldFixed64, uint64(23)))
	require.NoError(t, ps.Sfixed32(fieldSfixed32, int32(-22)))
	require.NoError(t, ps.Sfixed64(fieldSfixed64, int64(-23)))
	require.NoError(t, ps.Bool(fieldBool, true))
	require.NoError(t, ps.String(fieldString, "s"))
	require.NoError(t, ps.Bytes(fieldBytes, []byte("b")))
	require.NoError(t, ps.Int64Packed(fieldRepeatedInt64Packed, []int64{-1, 2, 3}))

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

	// test Reset by sending subsequent output to a new buffer
	output2 := bytes.NewBuffer([]byte{})
	ps.Reset(output2)
	require.NoError(t, ps.String(fieldString, "reset"))
	require.NoError(t, proto.Unmarshal(output2.Bytes(), &res))
	assert.Equal(t, "reset", res.String_)

	// test Write
	output3 := bytes.NewBuffer([]byte{})
	ps.Reset(output3)
	n, err := ps.Write([]byte{1, 2, 3, 4})
	require.NoError(t, err)
	require.Equal(t, 4, n)
}

// Test that the zero values for each field do not result in any encoded data.
func TestSimpleEncodingZero(t *testing.T) {
	output := bytes.NewBuffer([]byte{})
	ps := molecule.NewProtoStream(output)

	require.NoError(t, ps.Double(fieldDouble, 0.0))
	require.NoError(t, ps.Float(fieldFloat, 0.0))
	require.NoError(t, ps.Int32(fieldInt32, int32(0)))
	require.NoError(t, ps.Int64(fieldInt64, int64(0)))
	require.NoError(t, ps.Uint32(fieldUint32, uint32(0)))
	require.NoError(t, ps.Uint64(fieldUint64, uint64(0)))
	require.NoError(t, ps.Sint32(fieldSint32, int32(0)))
	require.NoError(t, ps.Sint64(fieldSint64, int64(0)))
	require.NoError(t, ps.Fixed32(fieldFixed32, uint32(0)))
	require.NoError(t, ps.Fixed64(fieldFixed64, uint64(0)))
	require.NoError(t, ps.Sfixed32(fieldSfixed32, int32(0)))
	require.NoError(t, ps.Sfixed64(fieldSfixed64, int64(0)))
	require.NoError(t, ps.Bool(fieldBool, false))
	require.NoError(t, ps.String(fieldString, ""))
	require.NoError(t, ps.Bytes(fieldBytes, []byte("")))

	buf := output.Bytes()
	// all of those are zero values, so nothing should have been written
	require.Equal(t, 0, len(buf))
}

// Test that the zero values for packed fields do not result in any encoded data.
func TestPackedEncodingZero(t *testing.T) {
	output := bytes.NewBuffer([]byte{})
	ps := molecule.NewProtoStream(output)

	require.NoError(t, ps.DoublePacked(fieldDouble, []float64{}))
	require.NoError(t, ps.FloatPacked(fieldFloat, []float32{}))
	require.NoError(t, ps.Int32Packed(fieldInt32, []int32{}))
	require.NoError(t, ps.Int64Packed(fieldInt64, []int64{}))
	require.NoError(t, ps.Uint32Packed(fieldUint32, []uint32{}))
	require.NoError(t, ps.Uint64Packed(fieldUint64, []uint64{}))
	require.NoError(t, ps.Sint32Packed(fieldSint32, []int32{}))
	require.NoError(t, ps.Sint64Packed(fieldSint64, []int64{}))
	require.NoError(t, ps.Fixed32Packed(fieldFixed32, []uint32{}))
	require.NoError(t, ps.Fixed64Packed(fieldFixed64, []uint64{}))
	require.NoError(t, ps.Sfixed32Packed(fieldSfixed32, []int32{}))
	require.NoError(t, ps.Sfixed64Packed(fieldSfixed64, []int64{}))

	buf := output.Bytes()
	// all of those are zero values, so nothing should have been written
	require.Equal(t, 0, len(buf))
}

// Test that *Packed functions work.
func TestPacking(t *testing.T) {
	output := bytes.NewBuffer([]byte{})
	ps := molecule.NewProtoStream(output)

	assertBytes := func(t *testing.T, exp []byte, got []byte) {
		assert.Equal(t, fmt.Sprintf("%#v", exp), fmt.Sprintf("%#v", got))
	}

	key := func(fieldNumber int) uint8 {
		const wtLengthDelimited = 2
		return uint8((fieldNumber << 3) + wtLengthDelimited)
	}

	t.Run("DoublePacked", func(t *testing.T) {
		output.Reset()
		require.NoError(t, ps.DoublePacked(fieldDouble, []float64{3.14, 1.414}))
		assertBytes(t, []byte{
			key(fieldDouble),
			0x10,                                          // length
			0x1f, 0x85, 0xeb, 0x51, 0xb8, 0x1e, 0x9, 0x40, // 3.14 as fixed64
			0x39, 0xb4, 0xc8, 0x76, 0xbe, 0x9f, 0xf6, 0x3f, // 1.414 as fixed64
		}, output.Bytes())
	})

	t.Run("FloatPacked", func(t *testing.T) {
		output.Reset()
		require.NoError(t, ps.FloatPacked(fieldFloat, []float32{3.14, 1.414}))
		assertBytes(t, []byte{
			key(fieldFloat),
			0x8,                    // length
			0xc3, 0xf5, 0x48, 0x40, // 3.14 as fixed32
			0xf4, 0xfd, 0xb4, 0x3f, // 1.414 as fixed32
		}, output.Bytes())
	})

	t.Run("Int32Packed", func(t *testing.T) {
		output.Reset()
		require.NoError(t, ps.Int32Packed(fieldInt32, []int32{int32(-12), int32(12), int32(-13)}))
		assertBytes(t, []byte{
			key(fieldInt32),
			0x15,                                                       // length
			0xf4, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01, // -12 (little-endian varint)
			0x0c,                                                       // 12
			0xf3, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01, // -13
		}, output.Bytes())
	})

	t.Run("Int64Packed", func(t *testing.T) {
		output.Reset()
		require.NoError(t, ps.Int64Packed(fieldInt64, []int64{int64(-12), int64(12), int64(-13)}))
		assertBytes(t, []byte{
			key(fieldInt64),
			0x15,                                                      // length
			0xf4, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x1, // -12
			0x0c,                                                      // 12
			0xf3, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x1, // -13
		}, output.Bytes())
	})

	t.Run("Uint32Packed", func(t *testing.T) {
		output.Reset()
		require.NoError(t, ps.Uint32Packed(fieldUint32, []uint32{uint32(1), uint32(2), uint32(3)}))
		assertBytes(t, []byte{
			key(fieldUint32),
			0x3, // length
			0x1,
			0x2,
			0x3,
		}, output.Bytes())
	})

	t.Run("Uint64Packed", func(t *testing.T) {
		output.Reset()
		require.NoError(t, ps.Uint64Packed(fieldUint64, []uint64{uint64(1), uint64(2), uint64(3)}))
		assertBytes(t, []byte{
			key(fieldUint64),
			0x3, // length
			0x1,
			0x2,
			0x3,
		}, output.Bytes())
	})

	t.Run("Sint32Packed", func(t *testing.T) {
		output.Reset()
		require.NoError(t, ps.Sint32Packed(fieldSint32, []int32{int32(-12), int32(12), int32(-13)}))
		assertBytes(t, []byte{
			key(fieldSint32),
			0x03, // length
			0x17, // zigzag encoding of -12
			0x18, // 12
			0x19, // 13
		}, output.Bytes())
	})

	t.Run("Sint64Packed", func(t *testing.T) {
		output.Reset()
		require.NoError(t, ps.Sint64Packed(fieldSint64, []int64{int64(-12), int64(12), int64(-13)}))
		assertBytes(t, []byte{
			key(fieldSint64),
			0x03, // length
			0x17,
			0x18,
			0x19,
		}, output.Bytes())
	})

	t.Run("Fixed32Packed", func(t *testing.T) {
		output.Reset()
		require.NoError(t, ps.Fixed32Packed(fieldFixed32, []uint32{uint32(12), uint32(13), uint32(14)}))
		assertBytes(t, []byte{
			key(fieldFixed32),
			0xc,                // length
			0xc, 0x0, 0x0, 0x0, // 12
			0xd, 0x0, 0x0, 0x0, // 13
			0xe, 0x0, 0x0, 0x0, // 14
		}, output.Bytes())
	})

	t.Run("Fixed64Packed", func(t *testing.T) {
		output.Reset()
		require.NoError(t, ps.Fixed64Packed(fieldFixed64, []uint64{uint64(12), uint64(13), uint64(14)}))
		assertBytes(t, []byte{
			key(fieldFixed64),
			0x18,                                   // length
			0xc, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, // 12
			0xd, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, // 13
			0xe, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, // 14
		}, output.Bytes())
	})

	t.Run("Sfixed32Packed", func(t *testing.T) {
		output.Reset()
		require.NoError(t, ps.Sfixed32Packed(fieldSfixed32, []int32{int32(12), int32(-13), int32(14)}))
		assertBytes(t, []byte{
			key(fieldSfixed32),
			0xc,                // length
			0xc, 0x0, 0x0, 0x0, // 12
			0xf3, 0xff, 0xff, 0xff, // -13
			0xe, 0x0, 0x0, 0x0, // 14
		}, output.Bytes())
	})

	t.Run("Sfixed64Packed", func(t *testing.T) {
		output.Reset()
		require.NoError(t, ps.Sfixed64Packed(fieldSfixed64, []int64{int64(12), int64(-13), int64(14)}))
		assertBytes(t, []byte{
			key(fieldSfixed64),
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
	ps := molecule.NewProtoStream(output)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		output.Reset()
		ps.Double(fieldDouble, 3.14)
		ps.Float(fieldFloat, 3.14)
		ps.Int32(fieldInt32, int32(-10))
		ps.Int64(fieldInt64, int64(-11))
		ps.Uint32(fieldUint32, uint32(12))
		ps.Uint64(fieldUint64, uint64(13))
		ps.Sint32(fieldSint32, int32(-12))
		ps.Sint64(fieldSint64, int64(-13))
		ps.Fixed32(fieldFixed32, uint32(22))
		ps.Fixed64(fieldFixed64, uint64(23))
		ps.Sfixed32(fieldSfixed32, int32(-22))
		ps.Sfixed64(fieldSfixed64, int64(-23))
		ps.Bool(fieldBool, true)
		ps.String(fieldString, "s")
		ps.Bytes(fieldBytes, []byte("b"))
		ps.Int64Packed(fieldRepeatedInt64Packed, []int64{99, -99})
	}
}

// Microbenchmark packing performance
func BenchmarkPacking(b *testing.B) {
	output := bytes.NewBuffer([]byte{})
	ps := molecule.NewProtoStream(output)

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
		err := ps.DoublePacked(fieldDouble, floats)
		if err != nil {
			panic(err)
		}
	}
}

// Test ps.Embedded embedding a repeated message
func TestEmbedding(t *testing.T) {
	output := bytes.NewBuffer([]byte{})
	ps := molecule.NewProtoStream(output)

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

func TestProtoStreamFuzzing(t *testing.T) {
	var (
		seed      = time.Now().UnixNano()
		fuzzer    = fuzz.NewWithSeed(seed)
		numFuzzes = 100000
	)
	defer func() {
		// Log the seed to make debugging failures easier.
		t.Logf("Running test with seed: %d", seed)
	}()
	// Limit slice size to prevent tests from taking too long.
	fuzzer.NumElements(0, 100)

	for i := 0; i < numFuzzes; i++ {
		m := &simple.Simple{}
		fuzzer.Fuzz(&m)
		if m == nil {
			continue
		}

		marshaled := bytes.NewBuffer([]byte{})
		ps := molecule.NewProtoStream(marshaled)

		require.NoError(t, ps.Double(fieldDouble, m.Double))
		require.NoError(t, ps.Float(fieldFloat, m.Float))
		require.NoError(t, ps.Int32(fieldInt32, m.Int32))
		require.NoError(t, ps.Int64(fieldInt64, m.Int64))
		require.NoError(t, ps.Uint32(fieldUint32, m.Uint32))
		require.NoError(t, ps.Uint64(fieldUint64, m.Uint64))
		require.NoError(t, ps.Sint32(fieldSint32, m.Sint32))
		require.NoError(t, ps.Sint64(fieldSint64, m.Sint64))
		require.NoError(t, ps.Fixed32(fieldFixed32, m.Fixed32))
		require.NoError(t, ps.Fixed64(fieldFixed64, m.Fixed64))
		require.NoError(t, ps.Sfixed32(fieldSfixed32, m.Sfixed32))
		require.NoError(t, ps.Sfixed64(fieldSfixed64, m.Sfixed64))
		require.NoError(t, ps.Bool(fieldBool, m.Bool))
		require.NoError(t, ps.String(fieldString, m.String_))
		require.NoError(t, ps.Bytes(fieldBytes, m.Bytes))
		require.NoError(t, ps.Int64Packed(fieldRepeatedInt64Packed, m.RepeatedInt64Packed))

		m2 := &simple.Simple{}
		require.NoError(t, proto.Unmarshal(marshaled.Bytes(), m2))

		require.True(t, proto.Equal(m, m2))
	}
}
