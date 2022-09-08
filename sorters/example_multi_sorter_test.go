package sorters_test

import (
	"encoding/json"
	"fmt"

	"github.com/fredbi/go-patterns/sorters"
	"github.com/go-openapi/swag"
)

type SampleUser struct {
	ID        string   `json:"id"`
	FirstName *string  `json:"first_name,omitempty"`
	LastName  string   `json:"last_name"`
	Age       *float64 `json:"age,omitempty"`
}

func ExampleMulti_Sort() {
	// builds a sortable collection with multiple sorting criteria.

	s := sorters.NewMulti[SampleUser](
		sampleUsers(),
		// order by:
		// * LastName ASC
		// * FirstName ASC
		// * Age ASC
		// * ID DESC
		func(a, b SampleUser) int {
			// Alternatively, we may use WithLocaleTag(language.French) to specify the collating sequence.
			return sorters.StringsComparator(sorters.WithLocale("fr"))(a.LastName, b.LastName)
		},
		func(a, b SampleUser) int {
			// Pointer logic is that nil < !nil , nil == nil. Usual comparison occurs whenever
			// both arguments are not nil
			return sorters.StringsPtrComparator(sorters.WithLocale("fr"))(a.FirstName, b.FirstName)
		},
		func(a, b SampleUser) int {
			// Ordered are go numerical types
			return sorters.OrderedPtrComparator[float64]()(a.Age, b.Age)
		},
		func(a, b SampleUser) int {
			return sorters.Reverse(
				// The default string comparison operates here (no language collation)
				sorters.StringsComparator(),
			)(a.ID, b.ID)
		},
	)

	// alternatively, we can call the sort package: sort.Sort(s)
	s.Sort()

	jazon, err := json.MarshalIndent(s.Collection(), "", " ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	fmt.Println(string(jazon))

	// Output:
	// [
	//  {
	//   "id": "10",
	//   "first_name": "Fred",
	//   "last_name": "B.",
	//   "age": 25
	//  },
	//  {
	//   "id": "4",
	//   "first_name": "Fred",
	//   "last_name": "B.",
	//   "age": 49
	//  },
	//  {
	//   "id": "2",
	//   "first_name": "Fred",
	//   "last_name": "B.",
	//   "age": 49
	//  },
	//  {
	//   "id": "1",
	//   "first_name": "Fred",
	//   "last_name": "B.",
	//   "age": 49
	//  },
	//  {
	//   "id": "2",
	//   "first_name": "Mathieu",
	//   "last_name": "B.",
	//   "age": 15
	//  },
	//  {
	//   "id": "5",
	//   "first_name": "Thomas",
	//   "last_name": "B.",
	//   "age": 13
	//  },
	//  {
	//   "id": "5",
	//   "last_name": "L.",
	//   "age": 45
	//  },
	//  {
	//   "id": "0",
	//   "first_name": "Enzo",
	//   "last_name": "L."
	//  }
	// ]
}

func sampleUsers() []SampleUser {
	return []SampleUser{
		{
			ID:        "1",
			FirstName: swag.String("Fred"),
			LastName:  "B.",
			Age:       swag.Float64(49),
		},
		{
			ID:        "4",
			FirstName: swag.String("Fred"),
			LastName:  "B.",
			Age:       swag.Float64(49),
		},
		{
			ID:        "10",
			FirstName: swag.String("Fred"),
			LastName:  "B.",
			Age:       swag.Float64(25),
		},
		{
			ID:       "5",
			LastName: "L.",
			Age:      swag.Float64(45),
		},
		{
			ID:        "0",
			FirstName: swag.String("Enzo"),
			LastName:  "L.",
		},
		{
			ID:        "2",
			FirstName: swag.String("Fred"),
			LastName:  "B.",
			Age:       swag.Float64(49),
		},
		{
			ID:        "5",
			FirstName: swag.String("Thomas"),
			LastName:  "B.",
			Age:       swag.Float64(13),
		},
		{
			ID:        "2",
			FirstName: swag.String("Mathieu"),
			LastName:  "B.",
			Age:       swag.Float64(15),
		},
	}
}
