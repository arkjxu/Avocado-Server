package main

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestUserRepoAdd(t *testing.T) {
	tester := "tester-" + testSessionID
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
	up := NewUserRepository(mysql)
	e = up.Add(&User{
		ID:           tester,
		Email:        tester,
		FirstName:    "tester",
		LastName:     "tester",
		RefreshToken: tester,
	})
	if e != nil {
		t.Fatal(e)
	}
}

func TestUserRepoGet(t *testing.T) {
	tester := "tester-" + testSessionID
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
	up := NewUserRepository(mysql)
	u, e := up.Get(&User{
		ID:           tester,
		Email:        tester,
		FirstName:    tester,
		LastName:     tester,
		RefreshToken: tester,
	})
	if e != nil {
		t.Fatal(e)
	}
	if u.Email != tester {
		t.Fatalf("Expecting tester: %s, but got %s", tester, u.Email)
	}
}
