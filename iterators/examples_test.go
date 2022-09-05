package iterators_test

type SampleStruct struct {
	A int
	B string
}

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
