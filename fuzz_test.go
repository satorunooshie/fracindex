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
		idx := newIndexerForTest(t)
		runOperationSequence(t, idx, data, 128)
	})
}

func FuzzCustomAlphabetOperationSequence(f *testing.F) {
	f.Add(byte(0), []byte{0, 1, 2, 3, 4, 5})
	f.Add(byte(1), []byte("base16"))
	f.Add(byte(2), []byte("0123456789abcdef"))
	f.Add(byte(3), []byte{255, 0, 128, 64, 32})

	f.Fuzz(func(t *testing.T, alphabet byte, data []byte) {
		idx := customIndexerForTest(t, alphabet)
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
				t.Errorf("NKeysBetween(%q, %q, %d) error = %v", prev, next, n, err)
				return
			}
			checkKeysBetween(t, idx, prev, next, batch)
			if t.Failed() {
				return
			}
			keys = slices.Insert(keys, pos, batch...)
			step += 2
		} else {
			key, err := idx.KeyBetween(prev, next)
			if err != nil {
				t.Errorf("KeyBetween(%q, %q) error = %v", prev, next, err)
				return
			}
			checkKeysBetween(t, idx, prev, next, []string{key})
			if t.Failed() {
				return
			}
			keys = slices.Insert(keys, pos, key)
			step++
		}

		checkKeysBetween(t, idx, "", "", keys)
		if t.Failed() {
			return
		}
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

func customIndexerForTest(tb testing.TB, alphabet byte) *Indexer {
	tb.Helper()

	switch alphabet % 4 {
	case 0:
		return newIndexerForTest(tb)
	case 1:
		return newIndexerForTest(tb, WithAlphabetSet(Base16))
	case 2:
		return newIndexerForTest(tb, WithDigitAlphabetSet(Base94))
	default:
		return newIndexerForTest(tb, WithDigitAlphabetSet(Base95))
	}
}
