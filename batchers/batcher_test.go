package batchers

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExecutor(t *testing.T) {
	t.Run("should execute incomplete batch", func(t *testing.T) {
		var (
			called int
			count  uint64
		)

		e := NewBatchExecutor[int](5, func(in []int) {
			if len(in) == 0 {
				return
			}
			count += uint64(len(in))
			called++
		})

		for _, item := range []int{1, 2, 3} {
			e.Push(item)
		}
		e.Flush()

		require.Equal(t, 1, called)
		require.Equal(t, uint64(3), count)
		require.Equal(t, count, e.Executed())
	})

	t.Run("should just work", func(t *testing.T) {
		var (
			called int
			count  uint64
		)

		e := NewBatchExecutor[*int](5, func(in []*int) {
			if len(in) == 0 {
				return
			}
			count += uint64(len(in))
			called++
		})

		for _, item := range []int{1, 2, 3} {
			e.Push(&item)
		}
		e.Flush()

		require.Equal(t, 1, called)
		require.Equal(t, uint64(3), count)
		require.Equal(t, count, e.Executed())
	})

	t.Run("should execute exact matching batch", func(t *testing.T) {
		var (
			called int
			count  uint64
		)

		e := NewBatchExecutor[int](2, func(in []int) {
			if len(in) == 0 {
				return
			}
			count += uint64(len(in))
			called++
		})

		for _, item := range []int{1, 2} {
			e.Push(item)
		}
		e.Flush()

		require.Equal(t, 1, called)
		require.Equal(t, uint64(2), count)
		require.Equal(t, count, e.Executed())
	})

	t.Run("pointer executor should execute exact matching batch", func(t *testing.T) {
		var (
			called int
			count  uint64
		)

		e := NewBatchPointerExecutor[int](2, func(in []*int) {
			if len(in) == 0 {
				return
			}
			count += uint64(len(in))
			called++
		})

		for _, toPin := range []int{1, 2} {
			item := toPin

			e.Push(&item)
		}
		e.Flush()

		require.Equal(t, 1, called)
		require.Equal(t, uint64(2), count)
		require.Equal(t, count, e.Executed())
	})

	t.Run("pointer executor should skip nil values", func(t *testing.T) {
		var (
			called int
			count  uint64
		)

		e := NewBatchPointerExecutor[int](2, func(in []*int) {
			if len(in) == 0 {
				return
			}
			count += uint64(len(in))
			called++
		})

		for _, toPin := range []int{1, 0, 2, 0} {
			item := toPin

			if item != 0 {
				e.Push(&item)
			} else {
				e.Push(nil)
			}
		}
		e.Flush()

		require.Equal(t, 1, called)
		require.Equal(t, uint64(2), count)
		require.Equal(t, count, e.Executed())
	})

	t.Run("should carry out options (dummy)", func(t *testing.T) {
		var called int

		dummyOption := func() Option {
			return func(o *options) {
				called++
			}
		}

		_ = NewBatchExecutor[int](2, func(in []int) {}, dummyOption())

		require.Equal(t, 1, called)
	})

	t.Run("executor should be goroutine-safe", func(t *testing.T) {
		var (
			called int
			count  uint64
			sum    int
		)

		e := NewBatchExecutor[int](3, func(in []int) {
			if len(in) == 0 {
				return
			}
			count += uint64(len(in))
			called++

			for _, element := range in {
				sum += element
			}
		})

		var wg sync.WaitGroup
		const n = 7
		for i := 0; i < n; i++ {
			v := i
			wg.Add(1)
			go func() {
				defer wg.Done()

				e.Push(v)
			}()
		}
		wg.Wait()
		e.Flush()

		require.Equal(t, 3, called)
		require.Equal(t, uint64(n), count)
		require.Equal(t, count, e.Executed())
		require.Equal(t, n*(n-1)/2, sum)
	})
}
