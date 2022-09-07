package iterators_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/fredbi/go-patterns/iterators"
)

type TransformedStruct struct {
	X int
	Y string
}

func ExampleTransformIterator_Next() {
	baseIterator := iterators.NewSliceIterator[SampleStruct](testSlice())

	transformer := func(ctx context.Context, in SampleStruct) (TransformedStruct, error) {
		ictx := iterators.GetIteratorContext(ctx)
		var index int

		// this retrieves the current iterated count from the context
		if ictx != nil {
			index = ictx.Iterated
		}

		fmt.Printf("transforming iteration %d\n", index)

		return TransformedStruct{
			X: in.A + index,
			Y: strings.Repeat(in.B, 2),
		}, nil
	}

	iterator := iterators.NewTransformIterator[SampleStruct, TransformedStruct](context.Background(), baseIterator, transformer)
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
	// transforming iteration 1
	// item: iterators_test.TransformedStruct{X:2, Y:"xx"}
	// transforming iteration 2
	// item: iterators_test.TransformedStruct{X:4, Y:"yy"}
	// count: 2
}

func ExampleTransformIterator() {
	// this example interrupts the iterations after 1 iteration, using
	// a transformer and the iterator's context
	baseIterator := iterators.NewSliceIterator[SampleStruct](testSlice())

	transformer := func(ctx context.Context, in SampleStruct) (SampleStruct, error) {
		ictx := iterators.GetIteratorContext(ctx)
		var index int

		// this retrieves the current iterated count from the context
		if ictx != nil {
			index = ictx.Iterated
		}
		if index > 1 {
			return SampleStruct{}, io.EOF
		}

		return in, nil
	}

	iterator := iterators.NewTransformIterator[SampleStruct, SampleStruct](context.Background(), baseIterator, transformer)
	defer func() {
		_ = iterator.Close()
	}()
	count := 0

	for iterator.Next() {
		item, err := iterator.Item()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			fmt.Printf("err: %v\n", err)
			break
		}

		count++
		fmt.Printf("item: %#v\n", item)
	}

	fmt.Printf("count: %d\n", count)

	// Output:
	// item: iterators_test.SampleStruct{A:1, B:"x"}
	// count: 1
}
