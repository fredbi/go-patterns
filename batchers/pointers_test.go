package batchers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPointers(t *testing.T) {
	t.Run("should be pointers", func(t *testing.T) {
		require.True(t, isPointer[*int]())
		require.True(t, isPointer[*struct{ a int }]())
		require.True(t, isPointer[*[]int]())
	})

	t.Run("should not be pointers", func(t *testing.T) {
		require.False(t, isPointer[int]())
		require.False(t, isPointer[struct{ a int }]())
		require.False(t, isPointer[[]int]())
	})

	t.Run("should shallow clone value of pointer type", func(t *testing.T) {
		cloner := clone[*int]()

		a := 1
		skip, b := cloner(&a)

		t.Log("skipped", skip)
		t.Logf("a=%d [%p]", a, &a)
		t.Logf("b=%d [%p]", *b, b)

	})
}
