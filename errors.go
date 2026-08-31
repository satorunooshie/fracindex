package fracindex

import "errors"

var (
	// ErrInvalidOption reports an invalid option passed to New.
	ErrInvalidOption = errors.New("invalid option")

	// ErrInvalidAlphabet reports an invalid digit or head alphabet.
	ErrInvalidAlphabet = errors.New("invalid alphabet")

	// ErrInvalidKey reports a malformed key for an Indexer's key space.
	ErrInvalidKey = errors.New("invalid key")

	// ErrInvalidRange reports bounds where prev does not sort before next.
	ErrInvalidRange = errors.New("invalid range")

	// ErrInvalidCount reports a negative count passed to NKeysBetween.
	ErrInvalidCount = errors.New("invalid count")

	// ErrRangeUnderflow reports that no smaller integer key can be generated.
	ErrRangeUnderflow = errors.New("range underflow")

	// ErrRangeOverflow reports that no larger integer key can be generated.
	ErrRangeOverflow = errors.New("range overflow")

	// ErrKeyspaceExhausted reports that a key hit the configured key space
	// constraints, usually WithMaxLength.
	ErrKeyspaceExhausted = errors.New("key space exhausted")
)
