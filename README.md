# fracindex

Fractional index keys for Go.

`fracindex` generates bytewise lexicographic order keys for mutable sequences.
`New()` keeps the classic Rocicorp-compatible base62 key space by default;
options create explicit key spaces with built-in or custom byte alphabets.
Errors wrap stable sentinel values, so callers can branch with `errors.Is`.

The algorithm follows the variable-length string key approach described in
[Implementing Fractional Indexing](https://observablehq.com/@dgreensp/implementing-fractional-indexing)
by [David Greenspan](https://github.com/dgreensp), a technique used for
[realtime editing of ordered sequences](https://www.figma.com/blog/realtime-editing-of-ordered-sequences/).

## Install

```sh
go get github.com/satorunooshie/fracindex
```

## Example

```go
idx, err := fracindex.New()
if err != nil {
	return err
}

first, _ := idx.KeyBetween("", "")
second, _ := idx.KeyBetween(first, "")
middle, _ := idx.KeyBetween(first, second)

fmt.Println(first, middle, second)
// a0 a0V a1
```

Use `NKeysBetween` when creating more than one key in the same interval:

```go
keys, _ := idx.NKeysBetween("", "", 5)
// []string{"a0", "a1", "a2", "a3", "a4"}
```

Empty bounds mean the beginning or end of the key space.

## Alphabets

Built-in alphabet sets are available for common byte alphabets:

```go
idx, _ := fracindex.New(fracindex.WithDigitAlphabetSet(fracindex.Base95))
```

Supported sets are `Base16`, `Base32`, `Base36`, `Base62`, `Base94`, and
`Base95`.

Raw alphabets are still supported:

```go
idx, _ := fracindex.New(fracindex.WithAlphabet("0123456789"))
```

See [SPEC.md](./SPEC.md) for alphabet constraints and exact behavior.

## Links

- [Specification](./SPEC.md)
- [Package documentation](https://pkg.go.dev/github.com/satorunooshie/fracindex)
