[![GoDoc](https://godoc.org/github.com/richardartoul/molecule?status.png)](https://godoc.org/github.com/richardartoul/molecule)
[![C.I](https://github.com/richardartoul/molecule/workflows/Go/badge.svg)](https://github.com/richardartoul/molecule/actions)

# Molecule

Molecule is a Go library for parsing protobufs in an efficient and zero-allocation manner. The API is loosely based on [this excellent](https://github.com/buger/jsonparser) Go JSON parsing library.

This library is in alpha and the API could change.

## Rationale

The standard `Unmarshal` protobuf interface in Go makes it difficult to manually control allocations when parsing protobufs. In addition, its common to only require access to a subset of an individual protobuf's fields. These issues make it hard to use protobuf in performance critical paths.

This library attempts to solve those problems by introducing a streaming, zero-allocation interface that allows users to have complete control over which fields are parsed, and how/when objects are allocated.

The downside, of course, is that `molecule` is more difficult to use (and easier to misuse) than the standard protobuf libraries so its recommended that it only be used in situations where performance is important. It is not a general purpose replacement for `proto.Unmarshal()`. It is recommended that users familiarize themselves with the [proto3 encoding](https://developers.google.com/protocol-buffers/docs/encoding) before attempting to use this library.

## Features

1. Unmarshal all protobuf primitive types with a streaming, zero-allocation API.
2. Support for iterating through protobuf messages in a streaming fashion.
3. Support for iterating through packed protobuf repeated fields (arrays) in a streaming fashion.

## Not Supported

1. Proto2 syntax (some things will probably work, but nothing is tested).
2. Repeated fields encoded not using the "packed" encoding (although in theory they can be parsed using this library, there just aren't any special helpers).
3. Map fields (and probably lots of other things).

## Examples

The `/src/examples/examples_test.go` file has a few examples (including one that demonstrates how to work with `repeated` fields), but the example below demonstrates a brief example of how the API can be used:

```proto3
message Test {
    string string_field = 1;
    int64 int64_field = 2;
}
```

```golang
    m := &Test{StringField: "hello world!"}
    marshaled, err := proto.Marshal(m)
    if err != nil {
        panic(err)
    }

    var (
        buffer = codec.NewBuffer(marshaled)
        strVal molecule.Value
        int64Val molcule.Value
    )
    err := molecule.MessageEach(buffer, func(fieldNum int32, value molecule.Value) bool {
        if fieldNum == 1 {
            strVal = value
        }
        if fieldNum == 2 {
            int64Val = value
        }

        // Continue scanning.
        return true
    })
    if err != nil {
        panic(err)
    }

    str, err := strVal.AsStringUnsafe()
    if err != nil {
        panic(err)
    }
    int64V, err := int64Val.AsInt64()
    if err != nil {
        panic(err)
    }

    fmt.Println("StringField: ", str)
    fmt.Println("Int64Field: ", int64V)
```

Note that in the example above the `str` variable in an "unsafe" view over the `marshaled` bytes. If those bytes were to be modified, pool, or reused in any way the value of the `str` variable would be undefined. If a safe value is required use the `AsStringSafe()` API instead, however, be aware that this will allocate a new string.

## Attributions

This library is mostly a thin wrapper around other people's work:

1. The interface was inspired by this [jsonparser](https://github.com/buger/jsonparser) library.
2. The codec for interacting with protobuf streams was lifted from this [protobuf reflection library](https://github.com/jhump/protoreflect). The code was manually vendored instead of imported to reduce dependencies.

## Dependencies
The core `molecule` library has zero external dependencies. The `go.sum` file does contain some dependencies introduced from the tests package, however,
those *should* not be included transitively when using this library.
