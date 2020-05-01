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
	var (
		seed   = time.Now().UnixNano()
		fuzzer = fuzz.NewWithSeed(seed)
	)
	// Limit slice size to prevent tests from taking too long.
	fuzzer.NumElements(0, 100)
	fuzzer.NilChance(0)

	m := &simple.Simple{}
	fuzzer.Fuzz(&m)
	marshaled, err := proto.Marshal(m)
	noErr(err)

	b.Run("standard unmarshal", func(b *testing.B) {
		into := &simple.Simple{}
		for i := 0; i < b.N; i++ {
			err := proto.Unmarshal(marshaled, into)
			noErr(err)
		}
	})

	b.Run("unmarshal single with molecule", func(b *testing.B) {
		msgBuffer := codec.NewBuffer(marshaled)
		for i := 0; i < b.N; i++ {
			msgBuffer.Reset(marshaled)
			err := molecule.MessageEach(msgBuffer, func(fieldNum int32, value molecule.Value) (bool, error) {
				switch fieldNum {
				case 14:
					_, err := value.AsStringUnsafe()
					noErr(err)
					return false, nil
				}

				return true, nil
			})
			noErr(err)
		}
	})

	b.Run("unmarshal multiple with molecule", func(b *testing.B) {
		var (
			msgBuffer   = codec.NewBuffer(marshaled)
			arrayBuffer = codec.NewBuffer(nil)
			int64s      []int64
		)
		for i := 0; i < b.N; i++ {
			msgBuffer.Reset(marshaled)
			int64s = int64s[:0]

			err := molecule.MessageEach(msgBuffer, func(fieldNum int32, value molecule.Value) (bool, error) {
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
					err = molecule.PackedRepeatedEach(arrayBuffer, codec.FieldType_INT64, func(value molecule.Value) (bool, error) {
						v, err := value.AsInt64()
						noErr(err)
						int64s = append(int64s, v)
						return true, nil
					})
					noErr(err)
				}

				return true, nil
			})
			noErr(err)
		}
	})
}

func noErr(err error) {
	if err != nil {
		panic(err)
	}
}
