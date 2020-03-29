package examples

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/richardartoul/molecule"
	"github.com/richardartoul/molecule/src/codec"
	"github.com/stretchr/testify/require"

	"github.com/richardartoul/molecule/src/proto"
)

func TestExampleSelectAField(t *testing.T) {
	m := &simple.Simple{String_: "hello world!"}
	marshaled, err := proto.Marshal(m)
	require.NoError(t, err)

	buffer := codec.NewBuffer(marshaled)
	molecule.MessageEach(buffer, func(fieldNum int32, value molecule.Value) bool {
		if fieldNum == 14 {
			str, err := value.AsStringUnsafe()
			require.NoError(t, err)
			require.Equal(t, "hello world!", str)
			// Found it, stop scanning.
			return false
		}
		// Continue scanning.
		return true
	})
}

func TestExampleSelectAPackedRepeatedField(t *testing.T) {
	int64s := []int64{1, 2, 3, 4, 5, 6, 7}
	m := &simple.Simple{RepeatedInt64Packed: int64s}
	marshaled, err := proto.Marshal(m)
	require.NoError(t, err)

	var (
		buffer          = codec.NewBuffer(marshaled)
		unmarshaledInts = []int64{}
	)
	molecule.MessageEach(buffer, func(fieldNum int32, value molecule.Value) bool {
		if fieldNum == 16 {
			packedArr, err := value.AsBytesUnsafe()
			require.NoError(t, err)

			buffer := codec.NewBuffer(packedArr)
			molecule.PackedArrayEach(buffer, codec.FieldDescriptorProto_TYPE_INT64, func(v molecule.Value) bool {
				vInt64, err := v.AsInt64()
				require.NoError(t, err)
				unmarshaledInts = append(unmarshaledInts, vInt64)
				return true
			})

			// Found it, stop scanning.
			return false
		}
		// Continue scanning.
		return true
	})

	require.Equal(t, int64s, unmarshaledInts)
}
