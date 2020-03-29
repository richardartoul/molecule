package moleculetest

import (
	"testing"
	"time"

	"github.com/richardartoul/molecule"
	"github.com/richardartoul/molecule/src/codec"
	"github.com/richardartoul/molecule/src/proto"

	"github.com/golang/protobuf/proto"
	"github.com/google/gofuzz"
)

func BenchmarkMolecule(b *testing.B) {
	b.StopTimer()
	var (
		seed   = time.Now().UnixNano()
		fuzzer = fuzz.NewWithSeed(seed)
	)
	// Limit slice size to prevent tests from taking too long.
	fuzzer.NumElements(0, 10)
	fuzzer.NilChance(0)

	m := &simple.Simple{}
	fuzzer.Fuzz(&m)
	marshaled, err := proto.Marshal(m)
	noErr(err)

	var (
		msgBuffer   = codec.NewBuffer(marshaled)
		arrayBuffer = codec.NewBuffer(nil)
		int64s      []int64
	)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		msgBuffer.Reset(marshaled)
		int64s = int64s[:0]

		err := molecule.MessageEach(msgBuffer, func(fieldNum int32, value molecule.Value) bool {
			switch fieldNum {
			case 14:
				_, err := value.AsStringUnsafe()
				noErr(err)
			case 15:
				_, err := value.AsBytesUnsafe()
				noErr(err)
			case 16:
				packedArr, err := value.AsBytesUnsafe()
				noErr(err)

				arrayBuffer.Reset(packedArr)
				err = molecule.PackedRepeatedEach(arrayBuffer, codec.FieldType_INT64, func(value molecule.Value) bool {
					v, err := value.AsInt64()
					noErr(err)
					int64s = append(int64s, v)
					return true
				})
				noErr(err)
			}

			return true
		})
		noErr(err)
	}
}

func noErr(err error) {
	if err != nil {
		panic(err)
	}
}
