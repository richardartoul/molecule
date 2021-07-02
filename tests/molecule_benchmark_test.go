package moleculetest

import (
	"testing"

	"github.com/richardartoul/molecule"
	"github.com/richardartoul/molecule/src/codec"
	simple "github.com/richardartoul/molecule/src/proto"

	"github.com/golang/protobuf/proto"
	fuzz "github.com/google/gofuzz"
)

func BenchmarkMolecule(b *testing.B) {
	var (
		seed   = int64(1623963202)
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

	b.Run("unmarshal all", func(b *testing.B) {
		msgBuffer := codec.NewBuffer(marshaled)
		m := &simple.Simple{}
		for i := 0; i < b.N; i++ {
			msgBuffer.Reset(marshaled)
			err := molecule.MessageEach(msgBuffer, func(fieldNum int32, value molecule.Value) (bool, error) {
				var err error
				switch fieldNum {
				case 1:
					m.Double, err = value.AsDouble()
				case 2:
					m.Float, err = value.AsFloat()
				case 3:
					m.Int32, err = value.AsInt32()
				case 4:
					m.Int64, err = value.AsInt64()
				case 5:
					m.Uint32, err = value.AsUint32()
				case 6:
					m.Uint64, err = value.AsUint64()
				case 7:
					m.Sint32, err = value.AsSint32()
				case 8:
					m.Sint64, err = value.AsSint64()
				case 9:
					m.Fixed32, err = value.AsFixed32()
				case 10:
					m.Fixed64, err = value.AsFixed64()
				case 11:
					m.Sfixed32, err = value.AsSFixed32()
				case 12:
					m.Sfixed64, err = value.AsSFixed64()
				case 13:
					m.Bool, err = value.AsBool()
				case 14:
					m.String_, err = value.AsStringUnsafe()
				case 15:
					m.Bytes, err = value.AsBytesUnsafe()
				case 16:

				}

				return err == nil, err
			})
			noErr(err)
		}
	})

	b.Run("unmarshal loop", func(b *testing.B) {
		msgBuffer := codec.NewBuffer(marshaled)
		m := &simple.Simple{}
		for i := 0; i < b.N; i++ {
			msgBuffer.Reset(marshaled)
			value := molecule.Value{}
			for !msgBuffer.EOF() {
				fieldNum, err := molecule.Next(msgBuffer, &value)
				noErr(err)

				switch fieldNum {
				case 1:
					m.Double, err = value.AsDouble()
				case 2:
					m.Float, err = value.AsFloat()
				case 3:
					m.Int32, err = value.AsInt32()
				case 4:
					m.Int64, err = value.AsInt64()
				case 5:
					m.Uint32, err = value.AsUint32()
				case 6:
					m.Uint64, err = value.AsUint64()
				case 7:
					m.Sint32, err = value.AsSint32()
				case 8:
					m.Sint64, err = value.AsSint64()
				case 9:
					m.Fixed32, err = value.AsFixed32()
				case 10:
					m.Fixed64, err = value.AsFixed64()
				case 11:
					m.Sfixed32, err = value.AsSFixed32()
				case 12:
					m.Sfixed64, err = value.AsSFixed64()
				case 13:
					m.Bool, err = value.AsBool()
				case 14:
					m.String_, err = value.AsStringUnsafe()
				case 15:
					m.Bytes, err = value.AsBytesUnsafe()
				case 16:

				}
				noErr(err)
			}
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
