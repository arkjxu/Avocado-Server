package main

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestSessionRepoAdd(t *testing.T) {
	tester := testSessionID
	mysql, e := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_KEY"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_DB")))
	if e != nil {
		t.Fatal(e)
	}
	defer mysql.Close()
	up := NewSessionRepository(mysql)
	e = up.Add(&Session{
		ID:     tester,
		UserId: tester,
		Email:  tester,
		UUID:   tester,
		Token:  tester,
	})
	if e != nil {
		t.Fatal(e)
	}
}

func TestSessionRepoGet(t *testing.T) {
	tester := testSessionID
	mysql, e := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_KEY"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_DB")))
	if e != nil {
		t.Fatal(e)
	}
	defer mysql.Close()
	up := NewSessionRepository(mysql)
	u, e := up.Get(&Session{
		ID:     tester,
		UserId: tester,
		Email:  tester,
		UUID:   tester,
		Token:  tester,
	})
	if e != nil {
		t.Fatal(e)
	}
	if u.UUID != tester {
		t.Fatalf("Expecting session: %s, but got %s", tester, u.UUID)
	}
}

func TestSessionRepoAddToken(t *testing.T) {
	tester := testSessionID
	mysql, e := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_KEY"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_DB")))
	if e != nil {
		t.Fatal(e)
	}
	defer mysql.Close()
	up := NewSessionRepository(mysql)
	_, e = up.AddToken(&Session{
		ID:     tester,
		UserId: tester,
		Email:  tester,
		UUID:   tester,
		Token:  tester,
	})
	if e != nil {
		t.Fatal(e)
	}
}

func TestSessionRepoRemove(t *testing.T) {
	tester := testSessionID
	mysql, e := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_KEY"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_DB")))
	if e != nil {
		t.Fatal(e)
	}
	defer mysql.Close()
	up := NewSessionRepository(mysql)
	e = up.Remove(&Session{
		ID:     tester,
		UserId: tester,
		Email:  tester,
		UUID:   tester,
		Token:  tester,
	})
	if e != nil {
		t.Fatal(e)
	}
}
