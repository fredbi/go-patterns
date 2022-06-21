//nolint: forbidigo
package batchers_test

import (
	"fmt"
	"testing"

	"github.com/fredbi/go-patterns/batchers"
)

type testItem struct {
	A int
}

func makeTestItems(n int) []testItem {
	fixtures := make([]testItem, n)

	for i := 0; i < n; i++ {
		fixtures[i] = testItem{
			A: i,
		}
	}

	return fixtures
}

func TestExecutor(_ *testing.T) {
	const n = 42

	batchExecutor := batchers.NewExecutor[testItem](10, func(in batchers.Batch[testItem]) {
		if len(in) == 0 {
			return
		}

		fmt.Printf("processing batch [%d items]: [%d-%d]\n", len(in), in[0].A, in[len(in)-1].A)
	})

	batchExecutorWithPointers := batchers.NewPointerExecutor[testItem](10, func(in batchers.Batch[*testItem]) {
		if len(in) == 0 {
			return
		}

		fmt.Printf("processing batch [%d pointer items]: [%d-%d]\n", len(in), in[0].A, in[len(in)-1].A)
	})

	for _, item := range makeTestItems(n) {
		batchExecutor.Push(item)
		batchExecutorWithPointers.Push(&item) //nolint:gosec  // we actually clone the value pointed to immediatly, so we may safely pass the iterated (constant) pointer here
	}

	batchExecutor.Flush()
	batchExecutorWithPointers.Flush()

	// Output:
	// processing batch [10 items]: [0-9]
	// processing batch [10 pointer items]: [0-9]
	// processing batch [10 items]: [10-19]
	// processing batch [10 pointer items]: [10-19]
	// processing batch [10 items]: [20-29]
	// processing batch [10 pointer items]: [20-29]
	// processing batch [10 items]: [30-39]
	// processing batch [10 pointer items]: [30-39]
	// processing batch [2 items]: [40-41]
	// processing batch [2 pointer items]: [40-41]
}
