package iterators_test

import (
	"fmt"
	"strings"
)

type (
	SampleStruct struct {
		A int
		B string
	}

	SortableStructs []SampleStruct
)

func testSlice() []SampleStruct {
	return []SampleStruct{
		{
			A: 1,
			B: "x",
		},
		{
			A: 2,
			B: "y",
		},
	}
}

func (s SortableStructs) String() string {
	var w strings.Builder
	fmt.Fprint(&w, "[")

	if len(s) > 1 {
		for _, elem := range s[:len(s)-1] {
			fmt.Fprintf(&w, "%v, ", elem)
		}
	}

	if len(s) > 0 {
		fmt.Fprintf(&w, "%v", s[len(s)-1])
	}

	fmt.Fprint(&w, "]")

	return w.String()
}

func (s SortableStructs) Len() int {
	return len(s)
}

func (s SortableStructs) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SortableStructs) Less(i, j int) bool {
	if s[i].A == s[j].A {
		return s[i].B < s[j].B
	}

	return s[i].A < s[j].A
}
