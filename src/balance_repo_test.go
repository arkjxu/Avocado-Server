package main

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestBalanceRepoAdd(t *testing.T) {
	tester := "tester"
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
	up := NewBalanceRepository(mysql)
	e = up.Add(&Balance{
		ID:      tester,
		Type:    "Asset",
		Email:   tester,
		Name:    tester,
		Balance: 100.0,
	})
	if e != nil {
		t.Fatal(e)
	}
}

func TestBalanceRepoGetAll(t *testing.T) {
	tester := "tester"
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
	up := NewBalanceRepository(mysql)
	u, e := up.GetAll(&User{
		ID:    tester,
		Email: tester,
	})
	if e != nil {
		t.Fatal(e)
	}
	if len(u) == 0 {
		t.Fatalf("No data for balance get all")
	}
}

func TestBalanceRepoRemove(t *testing.T) {
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
	up := NewBalanceRepository(mysql)
	e = up.Remove(&Balance{
		ID:    "100",
		Email: tester,
	})
	if e != nil {
		t.Fatal(e)
	}
}
