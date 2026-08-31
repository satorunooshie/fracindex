# Specification

## Key Space

An `Indexer` represents one key space. Keys should only be compared with keys
from an indexer created with the same options.

`New()` uses:

```text
digit alphabet: 0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz
head alphabet:  ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz
```

The default output is intended to match Rocicorp-style base62 fractional index
keys for the compatibility cases in `compat_test.go`.

## Ordering

Keys are ordered by bytewise lexicographic comparison. For any generated key
`k`:

```text
prev == "" || prev < k
next == "" || k < next
```

`NKeysBetween(prev, next, n)` returns `n` keys in ascending order. Empty bounds
mean the beginning or end of the key space. If both bounds are non-empty,
`prev` must sort before `next`.

## Format

A key is an integer part followed by an optional fractional part:

```text
a0V
|-| integer part
  ` fractional part
```

The first byte is a head byte. It encodes the integer-part length and side of
the key space. Remaining integer bytes and all fractional bytes come from the
digit alphabet.

Fractional parts must not end with the zero digit. This keeps keys canonical.

## Alphabets

Built-in sets:

Set | Alphabet
--- | ---
`Base16` | `0123456789ABCDEF`
`Base32` | `0123456789ABCDEFGHJKMNPQRSTVWXYZ`
`Base36` | `0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ`
`Base62` | `0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz`
`Base94` | Printable ASCII bytes from `!` (`0x21`) to `~` (`0x7e`)
`Base95` | Printable ASCII bytes from space (`0x20`) to `~` (`0x7e`)

Options:

- `WithAlphabetSet` sets digit and head alphabets from the same built-in set.
- `WithDigitAlphabetSet` sets only the digit alphabet from a built-in set.
- `WithHeadAlphabetSet` sets only the head alphabet from a built-in set.
- `WithAlphabet` sets digit and head alphabets from the same raw byte string.
- `WithDigitAlphabet` sets only the raw digit alphabet.
- `WithHeadAlphabet` sets only the raw head alphabet.

Alphabet constraints:

- At least two bytes
- Strictly sorted by byte value
- No duplicate bytes
- No NUL bytes
- Head alphabet length must be even

Alphabets are byte alphabets, not Unicode character alphabets.

## Max Length

`WithMaxLength(n)` rejects generated and validated keys longer than `n` bytes
with `ErrKeyspaceExhausted`.

## Errors

Errors wrap sentinel values and can be checked with `errors.Is`.

- `ErrInvalidOption`
- `ErrInvalidAlphabet`
- `ErrInvalidKey`
- `ErrInvalidRange`
- `ErrInvalidCount`
- `ErrRangeUnderflow`
- `ErrRangeOverflow`
- `ErrKeyspaceExhausted`
