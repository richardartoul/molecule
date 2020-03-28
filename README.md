# Molecule

Molecule is a Go library for parsing protobufs in an efficient and zero-allocation manner. The API is loosely based on [this excellent](https://github.com/buger/jsonparser) Go JSON parsing library.

This library is alpha quality and the API could change.

## Features

1. Unmarshal all protobuf primitive types with a streaming, zero-allocation API.
2. Support for iterating through protobuf messages in a streaming fashion.
3. Support for iterating through protobuf repeated fields (arrays) in a streaming fashion.

## Not Supported

1. Proto2 syntax (some thing will probably work, but nothing is tested).
2. Repeated fields encoded not using the "packed" encoding.