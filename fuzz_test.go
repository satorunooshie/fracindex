package fracindex

import (
	"slices"
	"testing"
)

func FuzzIndexerOperationSequence(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte{0})
	f.Add([]byte{0, 1, 2, 3, 4, 5})
	f.Add([]byte{255, 254, 253, 252, 251})
	f.Add([]byte("append-prepend-middle"))

	f.Fuzz(func(t *testing.T, data []byte) {
		idx := mustNewForTest(t)
		runOperationSequence(t, idx, data, 128)
	})
}

func FuzzCustomAlphabetOperationSequence(f *testing.F) {
	f.Add(byte(0), []byte{0, 1, 2, 3, 4, 5})
	f.Add(byte(1), []byte("base16"))
	f.Add(byte(2), []byte("0123456789abcdef"))
	f.Add(byte(3), []byte{255, 0, 128, 64, 32})

	f.Fuzz(func(t *testing.T, alphabet byte, data []byte) {
		idx := mustNewCustomForTest(t, alphabet)
		runOperationSequence(t, idx, data, 64)
	})
}

func runOperationSequence(t *testing.T, idx *Indexer, data []byte, maxOps int) {
	t.Helper()

	keys := make([]string, 0, min(len(data), maxOps))

	for step := 0; step < len(data) && step < maxOps && len(keys) < 256; {
		pos := int(data[step]) % (len(keys) + 1)
		prev, next := boundsAt(keys, pos)

		if data[step]&0x07 == 0 && step+1 < len(data) {
			n := int(data[step+1]%8) + 1
			batch, err := idx.NKeysBetween(prev, next, n)
			if err != nil {
				t.Fatalf("NKeysBetween(%q, %q, %d): %v", prev, next, n, err)
			}
			assertKeysBetween(t, idx, prev, next, batch)
			keys = slices.Insert(keys, pos, batch...)
			step += 2
		} else {
			key, err := idx.KeyBetween(prev, next)
			if err != nil {
				t.Fatalf("KeyBetween(%q, %q): %v", prev, next, err)
			}
			assertKeysBetween(t, idx, prev, next, []string{key})
			keys = slices.Insert(keys, pos, key)
			step++
		}

		assertKeysBetween(t, idx, "", "", keys)
	}
}

func boundsAt(keys []string, pos int) (prev, next string) {
	if pos > 0 {
		prev = keys[pos-1]
	}
	if pos < len(keys) {
		next = keys[pos]
	}
	return prev, next
}

func assertKeysBetween(t *testing.T, idx *Indexer, prev, next string, keys []string) {
	t.Helper()

	if !slices.IsSorted(keys) {
		t.Fatalf("keys are not sorted: %#v", keys)
	}
	for i, key := range keys {
		if err := idx.Validate(key); err != nil {
			t.Fatalf("key %q did not validate: %v", key, err)
		}
		if prev != "" && key <= prev {
			t.Fatalf("key[%d] = %q is not after prev %q", i, key, prev)
		}
		if next != "" && key >= next {
			t.Fatalf("key[%d] = %q is not before next %q", i, key, next)
		}
		if i > 0 && key <= keys[i-1] {
			t.Fatalf("key[%d] = %q is not after key[%d] = %q", i, key, i-1, keys[i-1])
		}
	}
}

func mustNewForTest(tb testing.TB) *Indexer {
	tb.Helper()

	idx, err := New()
	if err != nil {
		tb.Fatal(err)
	}
	return idx
}

func mustNewCustomForTest(tb testing.TB, alphabet byte) *Indexer {
	tb.Helper()

	var (
		idx *Indexer
		err error
	)
	switch alphabet % 4 {
	case 0:
		idx, err = New()
	case 1:
		idx, err = New(WithAlphabetSet(Base16))
	case 2:
		idx, err = New(WithDigitAlphabetSet(Base94))
	default:
		idx, err = New(WithDigitAlphabetSet(Base95))
	}
	if err != nil {
		tb.Fatal(err)
	}
	return idx
}
