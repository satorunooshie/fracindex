package fracindex

import (
	"fmt"
	"slices"
	"strings"
)

const (
	// Base16DigitAlphabet contains uppercase hexadecimal digits.
	Base16DigitAlphabet = "0123456789ABCDEF"

	// Base32DigitAlphabet contains a Crockford-style uppercase base32 alphabet.
	Base32DigitAlphabet = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

	// Base36DigitAlphabet contains decimal digits and uppercase Latin letters.
	Base36DigitAlphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// Base62DigitAlphabet is the default digit alphabet used by New.
	Base62DigitAlphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	// Base94DigitAlphabet contains printable ASCII bytes from exclamation mark
	// to tilde.
	Base94DigitAlphabet = "!\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~"

	// Base95DigitAlphabet contains printable ASCII bytes from space to tilde.
	Base95DigitAlphabet = " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~"

	// Base52HeadAlphabet is the default head alphabet used by New.
	Base52HeadAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

// AlphabetSet identifies a built-in byte alphabet.
type AlphabetSet uint8

const (
	// Base16 identifies Base16DigitAlphabet.
	Base16 AlphabetSet = iota

	// Base32 identifies Base32DigitAlphabet.
	Base32

	// Base36 identifies Base36DigitAlphabet.
	Base36

	// Base62 identifies Base62DigitAlphabet.
	Base62

	// Base94 identifies Base94DigitAlphabet.
	Base94

	// Base95 identifies Base95DigitAlphabet.
	Base95
)

// Indexer generates and validates keys for one fractional-index key space.
type Indexer struct {
	digits      string
	heads       string
	maxLength   int
	split       int
	zero        byte
	maxDigit    byte
	zeroKey     string
	smallestInt string
	digitIndex  [256]int
	headIndex   [256]int
}

// Option configures an Indexer created by New.
type Option interface {
	apply(*config) error
}

type optionFunc func(*config) error

func (f optionFunc) apply(c *config) error {
	return f(c)
}

type config struct {
	digitAlphabet string
	headAlphabet  string
	maxLength     int
}

// New returns an Indexer. Without options it creates a Rocicorp-compatible
// base62 key space.
func New(opts ...Option) (*Indexer, error) {
	c := config{
		digitAlphabet: Base62DigitAlphabet,
		headAlphabet:  Base52HeadAlphabet,
	}

	for _, opt := range opts {
		if opt == nil {
			return nil, fmt.Errorf("%w: nil option", ErrInvalidOption)
		}
		if err := opt.apply(&c); err != nil {
			return nil, err
		}
	}

	digitIndex, err := validateAlphabet(c.digitAlphabet, "digit alphabet")
	if err != nil {
		return nil, err
	}
	headIndex, err := validateAlphabet(c.headAlphabet, "head alphabet")
	if err != nil {
		return nil, err
	}
	if len(c.headAlphabet)%2 != 0 {
		return nil, fmt.Errorf("%w: head alphabet length must be even", ErrInvalidAlphabet)
	}

	split := len(c.headAlphabet) / 2
	zero := c.digitAlphabet[0]

	return &Indexer{
		digits:      c.digitAlphabet,
		heads:       c.headAlphabet,
		maxLength:   c.maxLength,
		split:       split,
		zero:        zero,
		maxDigit:    c.digitAlphabet[len(c.digitAlphabet)-1],
		zeroKey:     string(c.headAlphabet[split]) + string(zero),
		smallestInt: string(c.headAlphabet[0]) + strings.Repeat(string(zero), split),
		digitIndex:  digitIndex,
		headIndex:   headIndex,
	}, nil
}

// WithAlphabetSet sets the same built-in alphabet for digit and head positions.
// The selected alphabet must satisfy both digit and head alphabet constraints.
func WithAlphabetSet(set AlphabetSet) Option {
	return optionFunc(func(c *config) error {
		chars, err := alphabetForSet(set)
		if err != nil {
			return err
		}
		c.digitAlphabet = chars
		c.headAlphabet = chars
		return nil
	})
}

// WithAlphabet sets the same byte alphabet for digit and head positions. The
// alphabet must satisfy both digit and head alphabet constraints.
func WithAlphabet(chars string) Option {
	return optionFunc(func(c *config) error {
		c.digitAlphabet = chars
		c.headAlphabet = chars
		return nil
	})
}

// WithDigitAlphabetSet sets the built-in alphabet used for integer and
// fractional digits.
func WithDigitAlphabetSet(set AlphabetSet) Option {
	return optionFunc(func(c *config) error {
		chars, err := alphabetForSet(set)
		if err != nil {
			return err
		}
		c.digitAlphabet = chars
		return nil
	})
}

// WithDigitAlphabet sets the byte alphabet used for integer and fractional
// digits.
func WithDigitAlphabet(chars string) Option {
	return optionFunc(func(c *config) error {
		c.digitAlphabet = chars
		return nil
	})
}

// WithHeadAlphabetSet sets the built-in alphabet used for integer-part heads.
// The selected alphabet must have an even length.
func WithHeadAlphabetSet(set AlphabetSet) Option {
	return optionFunc(func(c *config) error {
		chars, err := alphabetForSet(set)
		if err != nil {
			return err
		}
		c.headAlphabet = chars
		return nil
	})
}

// WithHeadAlphabet sets the byte alphabet used for integer-part heads. The head
// alphabet is split in half: the first half represents negative integer lengths,
// and the second half represents positive integer lengths.
func WithHeadAlphabet(chars string) Option {
	return optionFunc(func(c *config) error {
		c.headAlphabet = chars
		return nil
	})
}

// WithMaxLength rejects generated and validated keys longer than n bytes.
func WithMaxLength(n int) Option {
	return optionFunc(func(c *config) error {
		if n <= 0 {
			return fmt.Errorf("%w: max length must be positive", ErrInvalidOption)
		}
		c.maxLength = n
		return nil
	})
}

func alphabetForSet(set AlphabetSet) (string, error) {
	switch set {
	case Base16:
		return Base16DigitAlphabet, nil
	case Base32:
		return Base32DigitAlphabet, nil
	case Base36:
		return Base36DigitAlphabet, nil
	case Base62:
		return Base62DigitAlphabet, nil
	case Base94:
		return Base94DigitAlphabet, nil
	case Base95:
		return Base95DigitAlphabet, nil
	default:
		return "", fmt.Errorf("%w: unknown alphabet set %d", ErrInvalidOption, set)
	}
}

// KeyBetween returns a key that sorts lexicographically between prev and next.
// Empty prev means the beginning of the key space. Empty next means the end.
func (x *Indexer) KeyBetween(prev, next string) (string, error) {
	if err := x.validateBound(prev); err != nil {
		return "", err
	}
	if err := x.validateBound(next); err != nil {
		return "", err
	}
	if prev != "" && next != "" && prev >= next {
		return "", fmt.Errorf("%w: prev must sort before next", ErrInvalidRange)
	}

	var key string
	var err error

	switch {
	case prev == "" && next == "":
		key = x.zeroKey
	case prev == "":
		key, err = x.keyBefore(next)
	case next == "":
		key, err = x.keyAfter(prev)
	default:
		key, err = x.keyBetween(prev, next)
	}
	if err != nil {
		return "", err
	}
	return x.enforceMaxLength(key)
}

// NKeysBetween returns n keys that sort lexicographically between prev and next.
// The returned keys are sorted in ascending order.
func (x *Indexer) NKeysBetween(prev, next string, n int) ([]string, error) {
	if n < 0 {
		return nil, fmt.Errorf("%w: n must be non-negative", ErrInvalidCount)
	}
	if n == 0 {
		return []string{}, nil
	}
	if n == 1 {
		key, err := x.KeyBetween(prev, next)
		if err != nil {
			return nil, err
		}
		return []string{key}, nil
	}

	if next == "" {
		key, err := x.KeyBetween(prev, next)
		if err != nil {
			return nil, err
		}
		keys := make([]string, 0, n)
		keys = append(keys, key)
		for i := 1; i < n; i++ {
			key, err = x.KeyBetween(key, next)
			if err != nil {
				return nil, err
			}
			keys = append(keys, key)
		}
		return keys, nil
	}

	if prev == "" {
		key, err := x.KeyBetween(prev, next)
		if err != nil {
			return nil, err
		}
		keys := make([]string, 0, n)
		keys = append(keys, key)
		for i := 1; i < n; i++ {
			key, err = x.KeyBetween(prev, key)
			if err != nil {
				return nil, err
			}
			keys = append(keys, key)
		}
		slices.Reverse(keys)
		return keys, nil
	}

	mid := n / 2
	key, err := x.KeyBetween(prev, next)
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0, n)
	left, err := x.NKeysBetween(prev, key, mid)
	if err != nil {
		return nil, err
	}
	keys = append(keys, left...)
	keys = append(keys, key)

	right, err := x.NKeysBetween(key, next, n-mid-1)
	if err != nil {
		return nil, err
	}
	keys = append(keys, right...)
	return keys, nil
}

// Validate reports whether key is well-formed for this Indexer's key space.
func (x *Indexer) Validate(key string) error {
	if key == "" {
		return fmt.Errorf("%w: empty key", ErrInvalidKey)
	}
	if err := x.enforceMaxLengthOnly(key); err != nil {
		return err
	}
	if key == x.smallestInt {
		return fmt.Errorf("%w: %q is reserved", ErrInvalidKey, key)
	}

	intPart, err := x.getIntPart(key)
	if err != nil {
		return err
	}
	if err := x.validateInt(intPart); err != nil {
		return err
	}

	frac := key[len(intPart):]
	for i := 0; i < len(frac); i++ {
		if x.digitIndex[frac[i]] < 0 {
			return fmt.Errorf("%w: invalid digit %q in %q", ErrInvalidKey, frac[i], key)
		}
	}
	if len(frac) > 0 && frac[len(frac)-1] == x.zero {
		return fmt.Errorf("%w: fractional part must not end with zero digit", ErrInvalidKey)
	}

	return nil
}

func (x *Indexer) validateBound(key string) error {
	if key == "" {
		return nil
	}
	return x.Validate(key)
}

func (x *Indexer) keyBefore(next string) (string, error) {
	intPart, err := x.getIntPart(next)
	if err != nil {
		return "", err
	}
	frac := next[len(intPart):]

	if intPart == x.smallestInt {
		mid, err := x.midpoint("", frac)
		if err != nil {
			return "", err
		}
		return intPart + mid, nil
	}
	if intPart < next {
		return intPart, nil
	}

	key, err := x.decrementInt(intPart)
	if err != nil {
		return "", err
	}
	if key == "" {
		return "", ErrRangeUnderflow
	}
	return key, nil
}

func (x *Indexer) keyAfter(prev string) (string, error) {
	intPart, err := x.getIntPart(prev)
	if err != nil {
		return "", err
	}
	frac := prev[len(intPart):]

	key, err := x.incrementInt(intPart)
	if err != nil {
		return "", err
	}
	if key == "" {
		mid, err := x.midpoint(frac, "")
		if err != nil {
			return "", err
		}
		return intPart + mid, nil
	}
	return key, nil
}

func (x *Indexer) keyBetween(prev, next string) (string, error) {
	prevInt, err := x.getIntPart(prev)
	if err != nil {
		return "", err
	}
	prevFrac := prev[len(prevInt):]

	nextInt, err := x.getIntPart(next)
	if err != nil {
		return "", err
	}
	nextFrac := next[len(nextInt):]

	if prevInt == nextInt {
		mid, err := x.midpoint(prevFrac, nextFrac)
		if err != nil {
			return "", err
		}
		return prevInt + mid, nil
	}

	key, err := x.incrementInt(prevInt)
	if err != nil {
		return "", err
	}
	if key == "" {
		return "", ErrRangeOverflow
	}
	if key < next {
		return key, nil
	}

	mid, err := x.midpoint(prevFrac, "")
	if err != nil {
		return "", err
	}
	return prevInt + mid, nil
}

func (x *Indexer) midpoint(prev, next string) (string, error) {
	if next != "" {
		i := 0
		for ; i < len(next); i++ {
			c := x.zero
			if len(prev) > i {
				c = prev[i]
			}
			if c != next[i] {
				break
			}
		}
		if i > 0 {
			midPrev := ""
			if i <= len(prev) {
				midPrev = prev[i:]
			}
			mid, err := x.midpoint(midPrev, next[i:])
			if err != nil {
				return "", err
			}
			return next[:i] + mid, nil
		}
	}

	prevDigit := 0
	if prev != "" {
		prevDigit = x.digitIndex[prev[0]]
		if prevDigit < 0 {
			return "", fmt.Errorf("%w: invalid digit %q", ErrInvalidKey, prev[0])
		}
	}

	nextDigit := len(x.digits)
	if next != "" {
		nextDigit = x.digitIndex[next[0]]
		if nextDigit < 0 {
			return "", fmt.Errorf("%w: invalid digit %q", ErrInvalidKey, next[0])
		}
	}

	if nextDigit-prevDigit > 1 {
		return string(x.digits[(prevDigit+nextDigit+1)/2]), nil
	}
	if len(next) > 1 {
		return next[:1], nil
	}

	rest := ""
	if len(prev) > 0 {
		rest = prev[1:]
	}
	mid, err := x.midpoint(rest, "")
	if err != nil {
		return "", err
	}
	return string(x.digits[prevDigit]) + mid, nil
}

func (x *Indexer) getIntPart(key string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("%w: empty key", ErrInvalidKey)
	}
	intLen, err := x.getIntLen(key[0])
	if err != nil {
		return "", err
	}
	if intLen > len(key) {
		return "", fmt.Errorf("%w: integer part is shorter than its head declares", ErrInvalidKey)
	}
	return key[:intLen], nil
}

func (x *Indexer) getIntLen(head byte) (int, error) {
	headPos := x.headIndex[head]
	if headPos < 0 {
		return 0, fmt.Errorf("%w: invalid head %q", ErrInvalidKey, head)
	}
	if headPos < x.split {
		return x.split - headPos + 1, nil
	}
	return headPos - x.split + 2, nil
}

func (x *Indexer) validateInt(intPart string) error {
	if intPart == "" {
		return fmt.Errorf("%w: empty integer part", ErrInvalidKey)
	}
	expected, err := x.getIntLen(intPart[0])
	if err != nil {
		return err
	}
	if len(intPart) != expected {
		return fmt.Errorf("%w: malformed integer part %q", ErrInvalidKey, intPart)
	}
	for i := 1; i < len(intPart); i++ {
		if x.digitIndex[intPart[i]] < 0 {
			return fmt.Errorf("%w: invalid digit %q in integer part", ErrInvalidKey, intPart[i])
		}
	}
	return nil
}

func (x *Indexer) incrementInt(intPart string) (string, error) {
	if err := x.validateInt(intPart); err != nil {
		return "", err
	}

	headPos := x.headIndex[intPart[0]]
	digits := []byte(intPart[1:])
	carry := true
	for i := len(digits) - 1; carry && i >= 0; i-- {
		digit := x.digitIndex[digits[i]]
		if digit < 0 {
			return "", fmt.Errorf("%w: invalid digit %q", ErrInvalidKey, digits[i])
		}
		digit++
		if digit == len(x.digits) {
			digits[i] = x.zero
		} else {
			digits[i] = x.digits[digit]
			carry = false
		}
	}

	if !carry {
		return string(intPart[0]) + string(digits), nil
	}
	if headPos == x.split-1 {
		return x.zeroKey, nil
	}
	if headPos == len(x.heads)-1 {
		return "", nil
	}

	nextHeadPos := headPos + 1
	if nextHeadPos > x.split {
		digits = append(digits, x.zero)
	} else {
		digits = digits[1:]
	}
	return string(x.heads[nextHeadPos]) + string(digits), nil
}

func (x *Indexer) decrementInt(intPart string) (string, error) {
	if err := x.validateInt(intPart); err != nil {
		return "", err
	}

	headPos := x.headIndex[intPart[0]]
	digits := []byte(intPart[1:])
	borrow := true
	for i := len(digits) - 1; borrow && i >= 0; i-- {
		digit := x.digitIndex[digits[i]]
		if digit < 0 {
			return "", fmt.Errorf("%w: invalid digit %q", ErrInvalidKey, digits[i])
		}
		digit--
		if digit == -1 {
			digits[i] = x.maxDigit
		} else {
			digits[i] = x.digits[digit]
			borrow = false
		}
	}

	if !borrow {
		return string(intPart[0]) + string(digits), nil
	}
	if headPos == x.split {
		return string(x.heads[x.split-1]) + string(x.maxDigit), nil
	}
	if headPos == 0 {
		return "", nil
	}

	prevHeadPos := headPos - 1
	if prevHeadPos < x.split-1 {
		digits = append(digits, x.maxDigit)
	} else {
		digits = digits[1:]
	}
	return string(x.heads[prevHeadPos]) + string(digits), nil
}

func (x *Indexer) enforceMaxLength(key string) (string, error) {
	if err := x.enforceMaxLengthOnly(key); err != nil {
		return "", err
	}
	return key, nil
}

func (x *Indexer) enforceMaxLengthOnly(key string) error {
	if x.maxLength > 0 && len(key) > x.maxLength {
		return fmt.Errorf("%w: length %d exceeds max length %d", ErrKeyspaceExhausted, len(key), x.maxLength)
	}
	return nil
}

func validateAlphabet(chars, name string) ([256]int, error) {
	var index [256]int
	for i := range index {
		index[i] = -1
	}
	if len(chars) < 2 {
		return index, fmt.Errorf("%w: %s must contain at least two bytes", ErrInvalidAlphabet, name)
	}
	for i := 0; i < len(chars); i++ {
		if chars[i] == 0 {
			return index, fmt.Errorf("%w: %s must not contain NUL", ErrInvalidAlphabet, name)
		}
		if i > 0 && chars[i-1] >= chars[i] {
			return index, fmt.Errorf("%w: %s must be strictly sorted by byte value with no duplicates", ErrInvalidAlphabet, name)
		}
		index[chars[i]] = i
	}
	return index, nil
}
