package moleculetest

import (
	"testing"

	"github.com/richardartoul/molecule/src"
	"github.com/richardartoul/molecule/src/proto/gen/pb-go"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/require"
)

func TestMoleculeSimple(t *testing.T) {
	m := &simple.Simple{
		Double:   1.0,
		Float:    2.0,
		Int32:    int32(3),
		Int64:    int64(4),
		Uint32:   uint32(5),
		Uint64:   uint64(6),
		Sint32:   int32(7),
		Sint64:   int64(8),
		Fixed32:  uint32(9),
		Fixed64:  uint64(10),
		Sfixed32: int32(11),
		Sfixed64: int64(12),
		Bool:     true,
		String_:  "fourteen",
		Bytes:    []byte("fifteen"),
	}
	marshaled, err := proto.Marshal(m)
	require.NoError(t, err)

	err = molecule.MessageEach(marshaled, func(fieldNum int32, value molecule.Value) error {
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
			v, err := value.AsString()
			require.NoError(t, err)
			require.Equal(t, m.String_, v)
		case 15:
			v, err := value.AsBytes()
			require.NoError(t, err)
			require.Equal(t, m.Bytes, v)
		}
		return nil
	})
	require.NoError(t, err)
}
