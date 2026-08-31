package fracindex

import (
	"slices"
	"testing"
)

func newIndexerForTest(tb testing.TB, opts ...Option) *Indexer {
	tb.Helper()

	idx, err := New(opts...)
	if err != nil {
		tb.Fatalf("New() error = %v", err)
	}
	return idx
}

func checkKeysBetween(t *testing.T, idx *Indexer, prev, next string, keys []string) bool {
	t.Helper()

	ok := true
	if !slices.IsSorted(keys) {
		t.Errorf("keys are not sorted: %#v", keys)
		ok = false
	}
	for i, key := range keys {
		if err := idx.Validate(key); err != nil {
			t.Errorf("key %q did not validate: %v", key, err)
			ok = false
		}
		if prev != "" && key <= prev {
			t.Errorf("key[%d] = %q is not after prev %q", i, key, prev)
			ok = false
		}
		if next != "" && key >= next {
			t.Errorf("key[%d] = %q is not before next %q", i, key, next)
			ok = false
		}
		if i > 0 && key <= keys[i-1] {
			t.Errorf("key[%d] = %q is not after key[%d] = %q", i, key, i-1, keys[i-1])
			ok = false
		}
	}
	return ok
}
