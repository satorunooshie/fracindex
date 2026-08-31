package fracindex

import "testing"

func BenchmarkKeyBetweenAppend(b *testing.B) {
	idx := mustNewForBenchmark(b)
	prev := ""

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		key, err := idx.KeyBetween(prev, "")
		if err != nil {
			b.Fatal(err)
		}
		prev = key
	}
}

func BenchmarkKeyBetweenPrepend(b *testing.B) {
	idx := mustNewForBenchmark(b)
	next := ""

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		key, err := idx.KeyBetween("", next)
		if err != nil {
			b.Fatal(err)
		}
		next = key
	}
}

func BenchmarkKeyBetweenDenseMiddle(b *testing.B) {
	idx := mustNewForBenchmark(b)
	prev, next := denseBoundsForBenchmark(b, idx, 64)

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		if _, err := idx.KeyBetween(prev, next); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkNKeysBetween32(b *testing.B) {
	idx := mustNewForBenchmark(b)

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		if _, err := idx.NKeysBetween("a0", "a1", 32); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkValidate(b *testing.B) {
	idx := mustNewForBenchmark(b)
	key := "a0V"

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		if err := idx.Validate(key); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCustomAlphabetKeyBetween(b *testing.B) {
	idx, err := New(WithAlphabet("0123456789"))
	if err != nil {
		b.Fatal(err)
	}
	prev, next := "50", "51"

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		if _, err := idx.KeyBetween(prev, next); err != nil {
			b.Fatal(err)
		}
	}
}

func denseBoundsForBenchmark(tb testing.TB, idx *Indexer, depth int) (string, string) {
	tb.Helper()

	prev, next := "a0", "a1"
	for range depth {
		key, err := idx.KeyBetween(prev, next)
		if err != nil {
			tb.Fatal(err)
		}
		prev = key
	}
	return prev, next
}

func mustNewForBenchmark(tb testing.TB) *Indexer {
	tb.Helper()

	idx, err := New()
	if err != nil {
		tb.Fatal(err)
	}
	return idx
}
