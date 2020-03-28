# Molecule

Molecule is a Go library for parsing protobufs in an efficient and zero-allocation manner. The API is loosely based on [this excellent](https://github.com/buger/jsonparser) Go JSON parsing library.

This library is in alpha and the API could change.

## Rationale

The standard `Unmarshal` protobuf interface in Go makes it extremely difficult to manually control allocations when parsing protobufs. In addition, its common to only require access to a subset of an individual protobuf's fields. These issues make using protobufs in performance critical paths difficult.

This library attempts to solve those problems by introducing a streaming, zero-allocation interface that allows users to have complete control over which fields are parsed, and how objects are allocated.

The downside, of course, is that `molecule` is more difficult to use (and easier to misuse) than the standard protobuf libraries so its recommended that it only be used in situations where performance is important. It is not a general purpose replacement for `proto.Unmarshal()`.

## Features

1. Unmarshal all protobuf primitive types with a streaming, zero-allocation API.
2. Support for iterating through protobuf messages in a streaming fashion.
3. Support for iterating through packed protobuf repeated fields (arrays) in a streaming fashion.

## Not Supported

1. Proto2 syntax (some thing will probably work, but nothing is tested).
2. Repeated fields encoded not using the "packed" encoding (although in theory they can be parsed using this library, there just aren't any special helpers).
3. Map fields (and probably lots of other things).

## Examples

```golang

```

## Attributions

This library is mostly a thin wrapper around other people's work:

1. The interface was inspired by this [jsonparser](https://github.com/buger/jsonparser) library.
2. The codec for interacting protobuf streams was lifted from this [protobuf reflection library](https://github.com/jhump/protoreflect). The code was manually vendored instead of imported to reduce dependencies.

## Dependencies
The core `molecule` library has zero external dependencies, although the go.sum file does contain some dependencies introduced from the tests package, however,
those should not be included transitively when using this library.