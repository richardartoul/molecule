package moleculetest

import (
	"testing"

	"github.com/richardartoul/molecule/src"
	"github.com/richardartoul/molecule/src/proto"

	"github.com/golang/protobuf/proto"
	"github.com/google/gofuzz"
	"github.com/stretchr/testify/require"
)

func TestMoleculeSimple(t *testing.T) {
	var (
		seed   = int64(1)
		fuzzer = fuzz.NewWithSeed(seed)
		m      = &simple.Simple{}
	)
	fuzzer.Fuzz(&m)

	marshaled, err := proto.Marshal(m)
	require.NoError(t, err)
	if len(marshaled) == 0 {
		t.SkipNow()
	}
	newM := &simple.Simple{}
	require.NoError(t, proto.Unmarshal(marshaled, newM))

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
		// case 10:
		// 	v, err := value.AsFixed64()
		// 	require.NoError(t, err)
		// 	require.Equal(t, m.Fixed64, v)
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
