package fracindex

import (
	"errors"
	"slices"
	"testing"
)

// Compatibility vectors follow rocicorp/fracdex, which is published under CC0-1.0.
func TestRocicorpKeyBetweenCompatibility(t *testing.T) {
	idx := newIndexerForTest(t)

	tests := []struct {
		name string
		prev string
		next string
		want string
	}{
		{name: "first", prev: "", next: "", want: "a0"},
		{name: "before first", prev: "", next: "a0", want: "Zz"},
		{name: "before negative integer", prev: "", next: "Zz", want: "Zy"},
		{name: "after first", prev: "a0", next: "", want: "a1"},
		{name: "after second", prev: "a1", next: "", want: "a2"},
		{name: "between consecutive integers", prev: "a0", next: "a1", want: "a0V"},
		{name: "between next consecutive integers", prev: "a1", next: "a2", want: "a1V"},
		{name: "between fractional and integer", prev: "a0V", next: "a1", want: "a0l"},
		{name: "between negative and zero", prev: "Zz", next: "a0", want: "ZzV"},
		{name: "between negative and positive", prev: "Zz", next: "a1", want: "a0"},
		{name: "roll over positive integer", prev: "az", next: "", want: "b00"},
		{name: "roll under negative integer", prev: "", next: "Z0", want: "Yzz"},
		{name: "roll under long negative integer", prev: "", next: "Y00", want: "Xzzz"},
		{name: "roll over long positive integer", prev: "bzz", next: "", want: "c000"},
		{name: "between integer and fraction", prev: "a0", next: "a0V", want: "a0G"},
		{name: "between integer and lower fraction", prev: "a0", next: "a0G", want: "a08"},
		{name: "between non-consecutive integers", prev: "b125", next: "b129", want: "b127"},
		{name: "between previous and fractional next integer", prev: "a0", next: "a1V", want: "a1"},
		{name: "between negative and positive fractional", prev: "Zz", next: "a01", want: "a0"},
		{name: "before positive fractional", prev: "", next: "a0V", want: "a0"},
		{name: "before long positive integer", prev: "", next: "b999", want: "b99"},
		{name: "between dense fractional", prev: "aV", next: "aV0V", want: "aV0G"},
		{
			name: "above smallest representable integer",
			prev: "",
			next: "A000000000000000000000000001",
			want: "A000000000000000000000000000V",
		},
		{
			name: "append at maximum integer predecessor",
			prev: "zzzzzzzzzzzzzzzzzzzzzzzzzzy",
			next: "",
			want: "zzzzzzzzzzzzzzzzzzzzzzzzzzz",
		},
		{
			name: "append after maximum integer",
			prev: "zzzzzzzzzzzzzzzzzzzzzzzzzzz",
			next: "",
			want: "zzzzzzzzzzzzzzzzzzzzzzzzzzzV",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := idx.KeyBetween(tt.prev, tt.next)
			if err != nil {
				t.Errorf("KeyBetween(%q, %q) error = %v", tt.prev, tt.next, err)
				return
			}
			if got != tt.want {
				t.Errorf("KeyBetween(%q, %q) = %q, want %q", tt.prev, tt.next, got, tt.want)
			}
		})
	}
}

func TestRocicorpKeyBetweenInvalidInputCompatibility(t *testing.T) {
	idx := newIndexerForTest(t)

	tests := []struct {
		name    string
		prev    string
		next    string
		wantErr error
	}{
		{
			name:    "smallest integer",
			prev:    "",
			next:    "A00000000000000000000000000",
			wantErr: ErrInvalidKey,
		},
		{name: "invalid previous key", prev: "a00", next: "", wantErr: ErrInvalidKey},
		{name: "invalid previous key with upper bound", prev: "a00", next: "a1", wantErr: ErrInvalidKey},
		{name: "invalid head", prev: "0", next: "1", wantErr: ErrInvalidKey},
		{name: "invalid range", prev: "a1", next: "a0", wantErr: ErrInvalidRange},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := idx.KeyBetween(tt.prev, tt.next)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("KeyBetween(%q, %q) error = %v, want %v", tt.prev, tt.next, err, tt.wantErr)
			}
		})
	}
}

func TestRocicorpNKeysBetweenCompatibility(t *testing.T) {
	idx := newIndexerForTest(t)

	tests := []struct {
		name string
		prev string
		next string
		n    int
		want []string
	}{
		{name: "initial", prev: "", next: "", n: 5, want: []string{"a0", "a1", "a2", "a3", "a4"}},
		{name: "append", prev: "a4", next: "", n: 10, want: []string{"a5", "a6", "a7", "a8", "a9", "aA", "aB", "aC", "aD", "aE"}},
		{name: "prepend", prev: "", next: "a0", n: 5, want: []string{"Zv", "Zw", "Zx", "Zy", "Zz"}},
		{
			name: "between sparse integers",
			prev: "a0",
			next: "a2",
			n:    20,
			want: []string{
				"a04", "a08", "a0G", "a0K", "a0O",
				"a0V", "a0Z", "a0d", "a0l", "a0t",
				"a1", "a14", "a18", "a1G", "a1O",
				"a1V", "a1Z", "a1d", "a1l", "a1t",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := idx.NKeysBetween(tt.prev, tt.next, tt.n)
			if err != nil {
				t.Errorf("NKeysBetween(%q, %q, %d) error = %v", tt.prev, tt.next, tt.n, err)
				return
			}
			if !slices.Equal(got, tt.want) {
				t.Errorf("NKeysBetween(%q, %q, %d) = %#v, want %#v", tt.prev, tt.next, tt.n, got, tt.want)
			}
		})
	}
}
