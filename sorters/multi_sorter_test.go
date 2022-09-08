package sorters

import (
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type sampleStruct struct {
	A string
	B int
}

func TestMulti(t *testing.T) {
	s := NewMulti(testCollection(),
		func(a, b sampleStruct) int {
			return strings.Compare(a.A, b.A)
		},
		func(a, b sampleStruct) int {
			if a.B == b.B {
				return 0
			}

			if a.B < b.B {
				return -1
			}

			return 1
		},
	)

	s.Sort()

	result := s.Collection()

	t.Logf("result: %v", result)
	require.EqualValues(t, expected(), result)
}

func TestNoCriteria(t *testing.T) {
	s := NewMulti(testCollection())

	s.Sort()

	require.EqualValues(t, testCollection(), s.Collection())
}

func TestSortPackage(t *testing.T) {
	s := NewMulti(testCollection(),
		func(a, b sampleStruct) int {
			return strings.Compare(a.A, b.A)
		},
		func(a, b sampleStruct) int {
			if a.B == b.B {
				return 0
			}

			if a.B < b.B {
				return -1
			}

			return 1
		},
	)

	sort.Stable(s)
	require.EqualValues(t, expected(), s.Collection())
}

func TestCompoundCriteria(t *testing.T) {
	comparator := CompoundCriteria[sampleStruct](
		func(a, b sampleStruct) int {
			return strings.Compare(a.A, b.A)
		},
		func(a, b sampleStruct) int {
			if a.B == b.B {
				return 0
			}

			if a.B < b.B {
				return -1
			}

			return 1
		},
	)

	a := sampleStruct{
		A: "a",
		B: 1,
	}
	b := sampleStruct{
		A: "b",
		B: 1,
	}
	c := sampleStruct{
		A: "a",
		B: 2,
	}

	require.Equal(t, 0, comparator(a, a))
	require.Equal(t, -1, comparator(a, b))
	require.Equal(t, 1, comparator(b, a))
	require.Equal(t, -1, comparator(a, c))
}

func expected() []sampleStruct {
	expected := testCollection()
	sort.Slice(expected,
		func(i, j int) bool {
			r := strings.Compare(expected[i].A, expected[j].A)
			if r == 0 {
				return expected[i].B < expected[j].B
			}

			return r < 0
		},
	)

	return expected
}

func testCollection() []sampleStruct {
	return []sampleStruct{
		{
			A: "c",
			B: 1,
		},
		{
			A: "a",
			B: 1,
		},
		{
			A: "a",
			B: 4,
		},
		{
			A: "b",
			B: 2,
		},
	}
}
