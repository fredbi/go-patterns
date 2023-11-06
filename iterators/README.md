# iterators

An iterator is essentially something that knows what comes `Next() bool`, then collects this `Item() T`.

The iterator pattern is useful to manipulate streams of data, and traverse several layers without undue buffering.

For example, we may fetch from a database an arbitrary number of rows.

The `iterators` package exposes 3 generic variants:
  1. A `SimpleIterator` built on top of a slice `[]T` (e.g. to build mocks, etc)
  2. A SQL rows iterator using github.com/jmoiron/sqlx.Rows and the `StructScan(interface{}) error` method.
     (this is used to iterate and unmarshal structs scanned from a SQL cursor).
  3. A `ChanIterator` that joins a collection of input iterators in parallel (the result is unordered).
  4. A `TransformIterator` that applies a data transform on the items iterated over anoother base iterator.


Sample code:  

See also the full [testable example](iterators/example_rows_iterator_test.go)

```go
    // create DB with some data
    ...

	// Open a cursor selecting over some test data.
    // Rows come from a SQL query. Here we use `sqlx.Rows` to directly
    // unmarshal the row using `StructScan` (this uses struct tags to decode SQL columns).
	rows, err := testdb.OpenDBCursor(db)
	if err != nil {
		log.Fatalf("could not query DB: %v", err)
	}

    // testdb.DummyRow is the go type receiving unmarshalled data
	iterator := iterators.NewRowsIterator[*sqlx.Rows, testdb.DummyRow](rows)
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
```

