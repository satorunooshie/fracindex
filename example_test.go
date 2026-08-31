package fracindex_test

import (
	"fmt"

	"github.com/satorunooshie/fracindex"
)

func ExampleIndexer_KeyBetween() {
	idx, _ := fracindex.New()

	first, _ := idx.KeyBetween("", "")
	second, _ := idx.KeyBetween(first, "")
	middle, _ := idx.KeyBetween(first, second)

	fmt.Println(first)
	fmt.Println(middle)
	fmt.Println(second)
	// Output:
	// a0
	// a0V
	// a1
}

func ExampleIndexer_NKeysBetween() {
	idx, _ := fracindex.New()

	keys, _ := idx.NKeysBetween("", "", 3)

	fmt.Println(keys)
	// Output:
	// [a0 a1 a2]
}

func ExampleWithAlphabet() {
	idx, _ := fracindex.New(fracindex.WithAlphabet("0123456789"))

	keys, _ := idx.NKeysBetween("", "", 3)

	fmt.Println(keys)
	// Output:
	// [50 51 52]
}
