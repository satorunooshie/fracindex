# fracindex

Fractional index keys for Go.

`fracindex` generates lexicographically sortable string keys for ordered
sequences. The default key space is compatible with Rocicorp-style base62
fractional indexing, while custom digit and head alphabets are available for
callers that need a different byte alphabet.

## Requirements

- Go 1.27 or later
- No third-party dependencies

## Install

```sh
go get github.com/satorunooshie/fracindex
```

## Features

Feature | Description
--- | ---
Rocicorp-compatible defaults | `New()` starts with the classic base62 key space.
Built-in alphabet sets | Choose base16, base32, base36, base62, base94, or base95.
Custom alphabets | Drop down to raw byte alphabets when presets are not enough.
Closed key spaces | An `Indexer` carries the key-space contract explicitly.
Public validation | Validate external or persisted keys with the same rules used by generation.
Typed errors | Sentinel errors support `errors.Is`.
Fuzzed invariants | Operation-sequence fuzz tests exercise ordering and validation properties.
Benchmarks | Public API benchmarks cover append, prepend, dense intervals, validation, and batching.

## Synopsis

```go
package main

import (
	"fmt"

	"github.com/satorunooshie/fracindex"
)

func main() {
	idx, err := fracindex.New()
	if err != nil {
		panic(err)
	}

	first, _ := idx.KeyBetween("", "")
	second, _ := idx.KeyBetween(first, "")
	middle, _ := idx.KeyBetween(first, second)

	fmt.Println(first)
	fmt.Println(middle)
	fmt.Println(second)
}
```

Output:

```text
a0
a0V
a1
```

Use `NKeysBetween` when creating more than one key in the same interval:

```go
keys, _ := idx.NKeysBetween("", "", 5)
// []string{"a0", "a1", "a2", "a3", "a4"}
```

Empty bounds represent the beginning or end of the key space.

## Alphabet Sets

`New()` creates this default key space:

```text
digit alphabet: 0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz
head alphabet:  ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz
```

Use `WithDigitAlphabetSet` to change only the digit alphabet. Base95 uses
printable ASCII bytes from space to tilde, so generated keys may contain spaces
and punctuation.

```go
idx, _ := fracindex.New(fracindex.WithDigitAlphabetSet(fracindex.Base95))

keys, _ := idx.NKeysBetween("", "", 4)
// []string{"a ", "a!", "a\"", "a#"}
```

Built-in sets:

Set | Alphabet
--- | ---
`Base16` | `0123456789ABCDEF`
`Base32` | `0123456789ABCDEFGHJKMNPQRSTVWXYZ`
`Base36` | `0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ`
`Base62` | `0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz`
`Base94` | Printable ASCII without space
`Base95` | Printable ASCII with space

## Custom Alphabets

Use raw alphabets when a preset does not fit.

```go
idx, _ := fracindex.New(fracindex.WithAlphabet("0123456789"))

keys, _ := idx.NKeysBetween("", "", 4)
// []string{"50", "51", "52", "53"}
```

Alphabets must be sorted by byte value, contain no duplicate bytes, contain at
least two bytes, and contain no NUL bytes. Head alphabets must have an even
length.

## API

```go
func New(opts ...Option) (*Indexer, error)

func WithAlphabetSet(set AlphabetSet) Option
func WithDigitAlphabetSet(set AlphabetSet) Option
func WithHeadAlphabetSet(set AlphabetSet) Option

func WithAlphabet(chars string) Option
func WithDigitAlphabet(chars string) Option
func WithHeadAlphabet(chars string) Option
func WithMaxLength(n int) Option

func (x *Indexer) KeyBetween(prev, next string) (string, error)
func (x *Indexer) NKeysBetween(prev, next string, n int) ([]string, error)
func (x *Indexer) Validate(key string) error
```

## Compatibility

The zero-option `New()` mode is intended to be byte-compatible with
Rocicorp-style base62 keys for the covered compatibility cases.

Compatibility tests live in `compat_test.go`. Core ordering invariants are
covered by unit tests and fuzz tests.

## Quality

Run the full local check:

```sh
make check
```

Run longer fuzzing:

```sh
make fuzz FUZZTIME=30s
```

Run benchmarks:

```sh
make bench
```

## Documentation

- [Specification](./SPEC.md)
- [Runnable examples](./example_test.go)
- [API reference](https://pkg.go.dev/github.com/satorunooshie/fracindex)

## Non-goals

- No package-level default generator
- No float conversion
- No database integration
- No automatic rebalance
- No CRDT conflict-resolution policy

## Prior Art

This package follows the variable-length integer approach described by David
Greenspan and used by Rocicorp's fractional-indexing libraries.
