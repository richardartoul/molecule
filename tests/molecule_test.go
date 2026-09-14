package moleculetest

import (
	"google.golang.org/protobuf/encoding/protowire"
	"testing"
	"time"

	"github.com/richardartoul/molecule"
	"github.com/richardartoul/molecule/src/codec"
	simple "github.com/richardartoul/molecule/src/proto"

	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

// TODO: Support and test enums.
// TODO: Add test-cases for repeated fields of all types.
func TestMoleculeSimple(t *testing.T) {
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

		marshaled, err := proto.Marshal(m)
		require.NoError(t, err)

		buffer := codec.NewBuffer(marshaled)
		err = molecule.MessageEach(buffer, func(fieldNum int32, value molecule.Value) (bool, error) {
			switch fieldNum {
			case 1:
				v, err := value.AsDouble()
				require.NoError(t, err)
				require.Equal(t, m.Double, v)
			case 2:
				v, err := value.AsFloat()
				require.NoError(t, err)
				require.Equal(t, m.Float, v)
			case 3:
				v, err := value.AsInt32()
				require.NoError(t, err)
				require.Equal(t, m.Int32, v)
			case 4:
				v, err := value.AsInt64()
				require.NoError(t, err)
				require.Equal(t, m.Int64, v)
			case 5:
				v, err := value.AsUint32()
				require.NoError(t, err)
				require.Equal(t, m.Uint32, v)
			case 6:
				v, err := value.AsUint64()
				require.NoError(t, err)
				require.Equal(t, m.Uint64, v)
			case 7:
				v, err := value.AsSint32()
				require.NoError(t, err)
				require.Equal(t, m.Sint32, v)
			case 8:
				v, err := value.AsSint64()
				require.NoError(t, err)
				require.Equal(t, m.Sint64, v)
			case 9:
				v, err := value.AsFixed32()
				require.NoError(t, err)
				require.Equal(t, m.Fixed32, v)
			case 10:
				v, err := value.AsFixed64()
				require.NoError(t, err)
				require.Equal(t, m.Fixed64, v)
			case 11:
				v, err := value.AsSFixed32()
				require.NoError(t, err)
				require.Equal(t, m.Sfixed32, v)
			case 12:
				v, err := value.AsSFixed64()
				require.NoError(t, err)
				require.Equal(t, m.Sfixed64, v)
			case 13:
				v, err := value.AsBool()
				require.NoError(t, err)
				require.Equal(t, m.Bool, v)
			case 14:
				v, err := value.AsStringUnsafe()
				require.NoError(t, err)
				require.Equal(t, m.String_, v)
			case 15:
				v, err := value.AsBytesUnsafe()
				require.NoError(t, err)
				require.Equal(t, m.Bytes, v)
				require.Equal(t, len(v), cap(v))
			case 16:
				packedArr, err := value.AsBytesUnsafe()
				require.NoError(t, err)

				var (
					int64s = []int64{}
					buffer = codec.NewBuffer(packedArr)
				)
				err = molecule.PackedRepeatedEach(buffer, codec.FieldType_INT64, func(value molecule.Value) (bool, error) {
					v, err := value.AsInt64()
					require.NoError(t, err)
					int64s = append(int64s, v)
					return true, nil
				})
				require.NoError(t, err)

				require.Equal(t, m.RepeatedInt64Packed, int64s)
			default:
				t.Errorf("unknown field number: %d", fieldNum)
			}
			return true, nil
		})
		require.NoError(t, err)
	}
}

func TestMoleculeProto2(t *testing.T) {
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
	fuzzer.NilChance(0)

	for i := 0; i < numFuzzes; i++ {
		m := &simple.MessageWithGroup{}
		fuzzer.Fuzz(&m)
		if m == nil {
			continue
		}

		marshaled, err := proto.Marshal(m)
		require.NoError(t, err)

		group := &simple.MessageWithGroup_Group{
			Double:   new(float64),
			Float:    new(float32),
			Int32:    new(int32),
			Int64:    new(int64),
			Uint32:   new(uint32),
			Uint64:   new(uint64),
			Sint32:   new(int32),
			Sint64:   new(int64),
			Fixed32:  new(uint32),
			Fixed64:  new(uint64),
			Sfixed32: new(int32),
			Sfixed64: new(int64),
			Bool:     new(bool),
			String_:  new(string),
		}

		var (
			buffer      = codec.NewBuffer(marshaled)
			groupBuffer = codec.NewBuffer(nil)
		)
		decodeGroupField := func(fieldNum int32, value molecule.Value) (bool, error) {
			var err error
			switch fieldNum {
			case 1:
				*group.Double, err = value.AsDouble()
			case 2:
				*group.Float, err = value.AsFloat()
			case 3:
				*group.Int32, err = value.AsInt32()
			case 4:
				*group.Int64, err = value.AsInt64()
			case 5:
				*group.Uint32, err = value.AsUint32()
			case 6:
				*group.Uint64, err = value.AsUint64()
			case 7:
				*group.Sint32, err = value.AsSint32()
			case 8:
				*group.Sint64, err = value.AsSint64()
			case 9:
				*group.Fixed32, err = value.AsFixed32()
			case 10:
				*group.Fixed64, err = value.AsFixed64()
			case 11:
				*group.Sfixed32, err = value.AsSFixed32()
			case 12:
				*group.Sfixed64, err = value.AsSFixed64()
			case 13:
				*group.Bool, err = value.AsBool()
			case 14:
				*group.String_, err = value.AsStringUnsafe()
			case 15:
				group.Bytes, err = value.AsBytesUnsafe()
			case 16:
				packedArr, err := value.AsBytesUnsafe()
				require.NoError(t, err)

				var (
					buffer = codec.NewBuffer(packedArr)
				)
				err = molecule.PackedRepeatedEach(buffer, codec.FieldType_INT64, func(value molecule.Value) (bool, error) {
					v, err := value.AsInt64()
					require.NoError(t, err)
					group.RepeatedInt64Packed = append(group.RepeatedInt64Packed, v)
					return true, nil
				})
				require.NoError(t, err)
				
			default:
				t.Errorf("unknown field number: %d", fieldNum)

			}

			return err == nil, err
		}
		decodeWholeMessage := func(fieldNum int32, value molecule.Value) (bool, error) {
			switch fieldNum {
			case 1:
				groupBytes, err := value.AsBytesUnsafe()
				noErr(err)

				groupBuffer.Reset(groupBytes)
				err = molecule.MessageEach(groupBuffer, decodeGroupField)
				noErr(err)
			}
			return true, nil
		}

		err = molecule.MessageEach(buffer, decodeWholeMessage)
		require.NoError(t, err)

		require.True(t, proto.Equal(m.Group, group))
	}
}

type FuzzSplitBytes struct {
	bytes    []byte
	selected int
}

func (f *FuzzSplitBytes) Fuzz(c fuzz.Continue) {
	f.selected = c.Intn(len(f.bytes))
}

func fuzzSplitBytes(fuzzer *fuzz.Fuzzer, bytes []byte) []byte {
	splitter := &FuzzSplitBytes{bytes, 0}
	fuzzer.Fuzz(splitter)
	return bytes[:splitter.selected]
}

func TestMoleculeTruncatedShouldNotPanic(t *testing.T) {
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

		marshaled, err := proto.Marshal(m)
		require.NoError(t, err)
		require.NotEmpty(t, marshaled)

		// split the buffer: iterating over it should not panic
		// there was a bug in DecodeVarint that caused panics
		splitMarshaled := fuzzSplitBytes(fuzzer, marshaled)
		buffer := codec.NewBuffer(splitMarshaled)
		err = molecule.MessageEach(buffer, func(fieldNum int32, value molecule.Value) (bool, error) {
			return true, nil
		})
		// sometimes this split will actually generate a "correct" truncation: err == nil is okay
		if err != nil {
			require.Error(t, err, "unexpected EOF", "wtf %#v %#v", string(splitMarshaled), m)
		}
	}
}

func TestMoleculeTruncated(t *testing.T) {
	payload := &simple.Simple{Int64: 42}
	serialized, err := proto.Marshal(payload)
	require.NoError(t, err)

	// truncate the payload at the 1 byte mark: parsing should return an error, not panic
	// this is a simplified case of a bug found by the TestMoleculeTruncatedShouldNotPanic fuzzer
	buffer := codec.NewBuffer(serialized[:1])
	err = molecule.MessageEach(buffer, func(fieldNum int32, value molecule.Value) (bool, error) {
		return true, nil
	})
	require.Error(t, err, "unexpected EOF")
}

func TestMoleculeGroups(t *testing.T) {
	encodeGroup := func(startFieldNum, endFieldNum protowire.Number, payload []byte) []byte {
		data := protowire.AppendTag(nil, startFieldNum, protowire.StartGroupType)
		data = append(data, payload...)
		data = protowire.AppendTag(data, endFieldNum, protowire.EndGroupType)
		return data
	}
	encodeVarintField := func(fieldNum protowire.Number, value int) []byte {
		data := protowire.AppendTag(nil, fieldNum, protowire.VarintType)
		data = protowire.AppendVarint(data, uint64(value))
		return data
	}

	tests := []struct {
		name        string
		data        []byte
		expectError bool
	}{
		{
			name: "valid group",
			data: encodeGroup(1, 1, encodeVarintField(1, 123)),
		},
		{
			name: "valid nested group",
			data: encodeGroup(1, 1, encodeGroup(3, 3, encodeVarintField(1, 123))),
		},
		{
			name:        "invalid group",
			data:        encodeGroup(1, 2, encodeVarintField(1, 123)),
			expectError: true,
		},
		{
			name:        "invalid nested group",
			data:        encodeGroup(1, 1, encodeGroup(3, 4, encodeVarintField(1, 123))),
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := molecule.MessageEach(
				codec.NewBuffer(test.data),
				func(fieldNum int32, value molecule.Value) (bool, error) {
					return true, nil
				})
			if test.expectError {
				require.Error(t, err)
				require.ErrorContains(t, err, "mismatched end group")
			} else {
				require.NoError(t, err)
			}
		})
	}
}
