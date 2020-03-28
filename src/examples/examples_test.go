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
			str, err := value.AsString()
			require.NoError(t, err)
			require.Equal(t, "hello world!", str)
			// Found it, stop scanning.
			return false
		}
		// Continue scanning.
		return true
	})
}
