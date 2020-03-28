package moleculetest

import (
	"testing"
	"time"

	"github.com/richardartoul/molecule"
	"github.com/richardartoul/molecule/src/codec"
	"github.com/richardartoul/molecule/src/proto"

	"github.com/golang/protobuf/proto"
	"github.com/google/gofuzz"
	"github.com/stretchr/testify/require"
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
		err = molecule.MessageEach(buffer, func(fieldNum int32, value molecule.Value) bool {
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
			case 16:
				packedArr, err := value.AsBytesUnsafe()
				require.NoError(t, err)

				var (
					int64s = []int64{}
					buffer = codec.NewBuffer(packedArr)
				)
				err = molecule.PackedArrayEach(buffer, codec.WireVarint, func(value molecule.Value) bool {
					v, err := value.AsInt64()
					require.NoError(t, err)
					int64s = append(int64s, v)
					return true
				})
				require.NoError(t, err)

				require.Equal(t, m.RepeatedInt64NotPacked, int64s)
			default:
				t.Errorf("unknown field number: %d", fieldNum)
			}
			return true
		})
		require.NoError(t, err)
	}
}
