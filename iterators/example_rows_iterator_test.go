package iterators_test

import (
	"fmt"
	"log"

	"github.com/fredbi/go-patterns/iterators"
	"github.com/fredbi/go-patterns/iterators/internal/testdb"
	"github.com/jmoiron/sqlx"
)

func ExampleRowsIterator_Next() {
	dbName := testdb.UniqueDBName()

	// create a DB and fill a table with some data
	db, err := testdb.CreateDBAndData(dbName)
	if err != nil {
		log.Fatalf("could not create test DB: %v", err)
	}

	// open a cursor selecting over the test data
	rows, err := testdb.OpenDBCursor(db)
	if err != nil {
		log.Fatalf("could not query DB: %v", err)
	}

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

	fmt.Printf("count: %d\n", count)

	// Output:
	// item: testdb.DummyRow{A:1, B:"x"}
	// item: testdb.DummyRow{A:2, B:"y"}
	// count: 2
}

func ExampleRowsIterator_Collect() {
	dbName := testdb.UniqueDBName()

	db, err := testdb.CreateDBAndData(dbName)
	if err != nil {
		log.Fatalf("could not create test DB: %v", err)
	}

	rows, err := testdb.OpenDBCursor(db)
	if err != nil {
		log.Fatalf("could not create test DB: %v", err)
	}

	iterator := iterators.NewSqlxIterator[testdb.DummyRow](rows)
	items, err := iterator.Collect()
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	fmt.Printf("items: %#v\n", items)
	fmt.Printf("count: %d\n", len(items))

	// Output:
	// items: []testdb.DummyRow{testdb.DummyRow{A:1, B:"x"}, testdb.DummyRow{A:2, B:"y"}}
	// count: 2
}

func ExampleRowsIterator_CollectPtr() {
	dbName := testdb.UniqueDBName()

	db, err := testdb.CreateDBAndData(dbName)
	if err != nil {
		log.Fatalf("could not create test DB: %v", err)
	}

	rows, err := testdb.OpenDBCursor(db)
	if err != nil {
		log.Fatalf("could not create test DB: %v", err)
	}

	iterator := iterators.NewSqlxIterator[testdb.DummyRow](rows)
	items, err := iterator.CollectPtr()
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	fmt.Printf("count: %d\n", len(items))

	// Output:
	// count: 2
}
