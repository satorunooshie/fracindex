package fracindex

import (
	"errors"
	"slices"
	"testing"
)

func TestAlphabetSets(t *testing.T) {
	tests := []struct {
		name string
		set  AlphabetSet
		want string
	}{
		{name: "base16", set: Base16, want: "a0"},
		{name: "base32", set: Base32, want: "a0"},
		{name: "base36", set: Base36, want: "a0"},
		{name: "base62", set: Base62, want: "a0"},
		{name: "base94", set: Base94, want: "a!"},
		{name: "base95", set: Base95, want: "a "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := newIndexerForTest(t, WithDigitAlphabetSet(tt.set))

			key, err := idx.KeyBetween("", "")
			if err != nil {
				t.Errorf("KeyBetween(%q, %q) error = %v", "", "", err)
				return
			}
			if key != tt.want {
				t.Errorf("first key = %q, want %q", key, tt.want)
			}
		})
	}
}

func TestSharedAlphabetSets(t *testing.T) {
	tests := []struct {
		name string
		set  AlphabetSet
	}{
		{name: "base16", set: Base16},
		{name: "base32", set: Base32},
		{name: "base36", set: Base36},
		{name: "base62", set: Base62},
		{name: "base94", set: Base94},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := newIndexerForTest(t, WithAlphabetSet(tt.set))

			keys, err := idx.NKeysBetween("", "", 4)
			if err != nil {
				t.Errorf("NKeysBetween(%q, %q, %d) error = %v", "", "", 4, err)
				return
			}
			checkKeysBetween(t, idx, "", "", keys)
		})
	}
}

func TestAlphabetSetBase95CannotBeShared(t *testing.T) {
	_, err := New(WithAlphabetSet(Base95))
	if !errors.Is(err, ErrInvalidAlphabet) {
		t.Errorf("New() error = %v, want ErrInvalidAlphabet", err)
	}
}

func TestInvalidAlphabetSet(t *testing.T) {
	_, err := New(WithDigitAlphabetSet(AlphabetSet(255)))
	if !errors.Is(err, ErrInvalidOption) {
		t.Errorf("New() error = %v, want ErrInvalidOption", err)
	}
}

func TestCustomAlphabetIsSelfHeaded(t *testing.T) {
	idx := newIndexerForTest(t, WithAlphabet("0123456789"))

	keys, err := idx.NKeysBetween("", "", 4)
	if err != nil {
		t.Errorf("NKeysBetween(%q, %q, %d) error = %v", "", "", 4, err)
		return
	}
	want := []string{"50", "51", "52", "53"}
	if !slices.Equal(keys, want) {
		t.Errorf("NKeysBetween with decimal custom alphabet = %#v, want %#v", keys, want)
	}
}

func TestCustomDigitAlphabetKeepsDefaultHeadAlphabet(t *testing.T) {
	idx := newIndexerForTest(t, WithDigitAlphabet("0123456789"))

	keys, err := idx.NKeysBetween("", "", 4)
	if err != nil {
		t.Errorf("NKeysBetween(%q, %q, %d) error = %v", "", "", 4, err)
		return
	}
	want := []string{"a0", "a1", "a2", "a3"}
	if !slices.Equal(keys, want) {
		t.Errorf("NKeysBetween with decimal custom digits and classic heads = %#v, want %#v", keys, want)
	}
}

func TestBase95DigitAlphabet(t *testing.T) {
	idx := newIndexerForTest(t, WithDigitAlphabetSet(Base95))

	keys, err := idx.NKeysBetween("", "", 4)
	if err != nil {
		t.Errorf("NKeysBetween(%q, %q, %d) error = %v", "", "", 4, err)
		return
	}
	want := []string{"a ", "a!", "a\"", "a#"}
	if !slices.Equal(keys, want) {
		t.Errorf("NKeysBetween with base95 digits = %#v, want %#v", keys, want)
	}
}

func TestGeneratedKeysAreValidSortedAndBounded(t *testing.T) {
	idx := newIndexerForTest(t)

	tests := []struct {
		prev string
		next string
		n    int
	}{
		{prev: "", next: "", n: 10},
		{prev: "", next: "a0", n: 10},
		{prev: "a0", next: "", n: 10},
		{prev: "a0", next: "a1", n: 10},
		{prev: "a0V", next: "a1", n: 10},
	}

	for _, tt := range tests {
		keys, err := idx.NKeysBetween(tt.prev, tt.next, tt.n)
		if err != nil {
			t.Errorf("NKeysBetween(%q, %q, %d) error = %v", tt.prev, tt.next, tt.n, err)
			continue
		}
		if len(keys) != tt.n {
			t.Errorf("got %d keys, want %d", len(keys), tt.n)
		}
		checkKeysBetween(t, idx, tt.prev, tt.next, keys)
	}
}

func TestValidateErrors(t *testing.T) {
	idx := newIndexerForTest(t)

	tests := []string{
		"",
		"0",
		"a",
		"a00",
		"a0V0",
		"A00000000000000000000000000",
	}

	for _, key := range tests {
		if err := idx.Validate(key); !errors.Is(err, ErrInvalidKey) {
			t.Errorf("Validate(%q) error = %v, want ErrInvalidKey", key, err)
		}
	}
}

func TestInvalidRange(t *testing.T) {
	idx := newIndexerForTest(t)

	_, err := idx.KeyBetween("a1", "a0")
	if !errors.Is(err, ErrInvalidRange) {
		t.Errorf("KeyBetween invalid range error = %v, want ErrInvalidRange", err)
	}
}

func TestInvalidCount(t *testing.T) {
	idx := newIndexerForTest(t)

	_, err := idx.NKeysBetween("", "", -1)
	if !errors.Is(err, ErrInvalidCount) {
		t.Errorf("NKeysBetween invalid count error = %v, want ErrInvalidCount", err)
	}
}

func TestMaxLength(t *testing.T) {
	idx := newIndexerForTest(t, WithMaxLength(2))

	key, err := idx.KeyBetween("", "")
	if err != nil {
		t.Errorf("KeyBetween(%q, %q) error = %v", "", "", err)
		return
	}
	if key != "a0" {
		t.Errorf("first key = %q, want a0", key)
	}

	_, err = idx.KeyBetween("a0", "a1")
	if !errors.Is(err, ErrKeyspaceExhausted) {
		t.Errorf("KeyBetween length error = %v, want ErrKeyspaceExhausted", err)
	}
}

func TestInvalidAlphabet(t *testing.T) {
	tests := []struct {
		name string
		opts []Option
	}{
		{name: "too short", opts: []Option{WithDigitAlphabet("0")}},
		{name: "unsorted", opts: []Option{WithDigitAlphabet("10")}},
		{name: "duplicate", opts: []Option{WithDigitAlphabet("001")}},
		{name: "odd shared alphabet", opts: []Option{WithAlphabet("012")}},
		{name: "odd head alphabet", opts: []Option{WithHeadAlphabet("ABC")}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.opts...)
			if !errors.Is(err, ErrInvalidAlphabet) {
				t.Errorf("New() error = %v, want ErrInvalidAlphabet", err)
			}
		})
	}
}
