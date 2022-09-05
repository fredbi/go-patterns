package iterators_test

import (
	"fmt"

	"github.com/fredbi/go-patterns/iterators"
)

func ExampleSliceIterator_Next() {
	iterator := iterators.NewSliceIterator[SampleStruct](testSlice())
	defer func() {
		_ = iterator.Close()
	}()
	count := 0

	for iterator.Next() {
		item, err := iterator.Item()
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}
		count++
		fmt.Printf("item: %#v\n", item)
	}

	fmt.Printf("count: %d\n", count)

	// Output:
	// item: iterators_test.SampleStruct{A:1, B:"x"}
	// item: iterators_test.SampleStruct{A:2, B:"y"}
	// count: 2
}

func ExampleSliceIterator_Collect() {
	iterator := iterators.NewSliceIterator[SampleStruct](testSlice())
	items, err := iterator.Collect()
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	fmt.Printf("items: %#v\n", items)
	fmt.Printf("count: %d\n", len(items))

	// Output:
	// items: []iterators_test.SampleStruct{iterators_test.SampleStruct{A:1, B:"x"}, iterators_test.SampleStruct{A:2, B:"y"}}
	// count: 2
}

func ExampleSliceIterator_CollectPtr() {
	iterator := iterators.NewSliceIterator[SampleStruct](testSlice())
	items, err := iterator.CollectPtr()
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	fmt.Printf("count: %d\n", len(items))

	// Output:
	// count: 2
}
