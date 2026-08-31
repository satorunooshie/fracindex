package fracindex

import (
	"slices"
	"testing"
)

func TestRocicorpKeyBetweenCompatibility(t *testing.T) {
	idx, err := New()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		prev string
		next string
		want string
	}{
		{name: "first", prev: "", next: "", want: "a0"},
		{name: "before first", prev: "", next: "a0", want: "Zz"},
		{name: "after first", prev: "a0", next: "", want: "a1"},
		{name: "after second", prev: "a1", next: "", want: "a2"},
		{name: "between consecutive integers", prev: "a0", next: "a1", want: "a0V"},
		{name: "between next consecutive integers", prev: "a1", next: "a2", want: "a1V"},
		{name: "between negative and zero", prev: "Zz", next: "a0", want: "ZzV"},
		{name: "roll over positive integer", prev: "az", next: "", want: "b00"},
		{name: "roll under negative integer", prev: "", next: "Z0", want: "Yzz"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := idx.KeyBetween(tt.prev, tt.next)
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Fatalf("KeyBetween(%q, %q) = %q, want %q", tt.prev, tt.next, got, tt.want)
			}
		})
	}
}

func TestRocicorpNKeysBetweenCompatibility(t *testing.T) {
	idx, err := New()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		prev string
		next string
		n    int
		want []string
	}{
		{name: "initial", prev: "", next: "", n: 5, want: []string{"a0", "a1", "a2", "a3", "a4"}},
		{name: "append", prev: "a1", next: "", n: 2, want: []string{"a2", "a3"}},
		{name: "prepend", prev: "", next: "a0", n: 5, want: []string{"Zv", "Zw", "Zx", "Zy", "Zz"}},
		{name: "between", prev: "a0", next: "a1", n: 2, want: []string{"a0G", "a0V"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := idx.NKeysBetween(tt.prev, tt.next, tt.n)
			if err != nil {
				t.Fatal(err)
			}
			if !slices.Equal(got, tt.want) {
				t.Fatalf("NKeysBetween(%q, %q, %d) = %#v, want %#v", tt.prev, tt.next, tt.n, got, tt.want)
			}
		})
	}
}
