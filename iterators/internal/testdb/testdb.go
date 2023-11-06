// Package testdb expose internal utilities for testing
// against a postgres database.
package testdb

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"sync"

	sq "github.com/Masterminds/squirrel"
	// initialize DB driver
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

const (
	base   = "postgresql://postgres@localhost:5432?sslmode=disable"
	dbName = "test_iterator_%d"
)

type DummyRow struct {
	A int    `db:"a"`
	B string `db:"b"`
}

var (
	testIndex int
	mx        sync.Mutex
)

func UniqueDBName() string {
	mx.Lock()
	defer mx.Unlock()
	testIndex++

	return fmt.Sprintf(dbName, testIndex)
}

func OpenDB(dbName string) (*sqlx.DB, error) {
	u, err := url.Parse(base)
	if err != nil {
		return nil, err
	}
	u.Path = dbName

	info := u.User
	if info == nil {
		info = url.User("postgres")
	}

	// take overrides from env, e.g. for CI
	if user := os.Getenv("PGUSER"); user != "" {
		pwd, _ := info.Password()
		info = url.UserPassword(user, pwd)
		u.User = info
	}
	if pwd := os.Getenv("PGPASSWORD"); pwd != "" {
		user := info.Username()
		info = url.UserPassword(user, pwd)
		u.User = info
	}
	host := os.Getenv("PGHOST")
	port := os.Getenv("PGPORT")
	switch {
	case host != "" && port == "":
		u.Host = fmt.Sprintf("%s:%s", host, "5432")
	case host != "" && port != "":
		u.Host = fmt.Sprintf("%s:%s", host, port)
	case host == "" && port != "":
		u.Host = fmt.Sprintf("%s:%s", "localhost", port)
	}

	log.Printf("DEBUG: postgres URL: %s", u.String())

	db, err := sqlx.Open("pgx", u.String())
	if err != nil {
		return nil, err
	}

	return db, nil
}

func CloseDB(db *sqlx.DB) error {
	return db.Close()
}

func CreateDBAndData(dbName string) (*sqlx.DB, error) {
	db, err := createDB(dbName)
	if err != nil {
		return nil, err
	}

	if err = createData(db); err != nil {
		return nil, err
	}

	return db, err
}

func CreateDBWithWrongData(dbName string) (*sqlx.DB, error) {
	db, err := createDB(dbName)
	if err != nil {
		return nil, err
	}

	if err = createWrongData(db); err != nil {
		return nil, err
	}

	return db, err
}

func createDB(dbName string) (*sqlx.DB, error) {
	emptyDB, err := OpenDB("")
	if err != nil {
		return nil, err
	}

	_, err = emptyDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName))
	if err != nil {
		return nil, err
	}

	_, err = emptyDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
	if err != nil {
		return nil, err
	}

	if err = emptyDB.Close(); err != nil {
		return nil, err
	}

	return OpenDB(dbName)
}

func createData(db *sqlx.DB) error {
	_, err := db.Exec("CREATE TABLE dummy(a integer, b character varying)")
	if err != nil {
		return err
	}

	insert, args := sq.
		Insert(
			"dummy",
		).
		Columns(
			"a", "b",
		).
		Values(1, "x").
		Values(2, "y").
		PlaceholderFormat(sq.Dollar).
		MustSql()

	_, err = db.Exec(insert, args...)

	return err
}

func createWrongData(db *sqlx.DB) error {
	_, err := db.Exec("CREATE TABLE dummy(a integer, b character varying)")
	if err != nil {
		return err
	}

	insert, args := sq.
		Insert(
			"dummy",
		).
		Columns(
			"a", "b",
		).
		Values(1, "x").
		Values(2, nil).
		PlaceholderFormat(sq.Dollar).
		MustSql()

	_, err = db.Exec(insert, args...)

	return err
}

func OpenDBCursor(db *sqlx.DB) (*sqlx.Rows, error) {
	query, args := sq.Select(
		"a", "b",
	).From(
		"dummy",
	).OrderBy(
		"a",
	).MustSql()

	return db.QueryxContext(context.Background(), query, args...)
}

func EmptyDBCursor(db *sqlx.DB) (*sqlx.Rows, error) {
	query, args := sq.Select(
		"a", "b",
	).From(
		"dummy",
	).OrderBy(
		"a",
	).Where(
		"1 = 0",
	).MustSql()

	return db.QueryxContext(context.Background(), query, args...)
}
