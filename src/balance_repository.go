package main

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go/aws"
)

type Balance struct {
	ID      string  `json:"id"`
	Email   string  `json:"email,omitempty"`
	Type    string  `json:"type"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
	Created string  `json:"created"`
}

type BalanceRepository struct {
	db    *sql.DB
	table *string
}

func NewBalanceRepository(DB *sql.DB) *BalanceRepository {
	return &BalanceRepository{db: DB, table: aws.String("avocado-balances")}
}

func (b *BalanceRepository) Get(bal *Balance) (e error) {
	return errors.New("Not implemented")
}

func (b *BalanceRepository) Add(bal *Balance) (e error) {
	ctx, cnl := context.WithTimeout(context.Background(), 3*time.Second)
	defer cnl()
	row := b.db.QueryRowContext(ctx, "CALL Balances_Add(?,?,?,?)", bal.Email, bal.Name, bal.Type, bal.Balance)
	if row.Err() != nil {
		return e
	}
	e = row.Scan(&bal.ID, &bal.Created)
	return e
}

func (b *BalanceRepository) Remove(bal *Balance) (e error) {
	ctx, cnl := context.WithTimeout(context.Background(), 3*time.Second)
	defer cnl()
	_, e = b.db.ExecContext(ctx, "CALL Balances_Delete(?)", bal.ID)
	return e
}

func (b *BalanceRepository) GetAll(user *User) (ub []Balance, e error) {
	ctx, cnl := context.WithTimeout(context.Background(), 6*time.Second)
	defer cnl()
	rows, e := b.db.QueryContext(ctx, "CALL Balances_GetAll(?)", user.Email)
	if e != nil && e != sql.ErrNoRows {
		return ub, e
	}
	ub = []Balance{}
	for rows.Next() {
		var b Balance
		e = rows.Scan(&b.ID, &b.Name, &b.Type, &b.Balance, &b.Created)
		if e != nil {
			return ub, e
		}
		ub = append(ub, b)
	}
	return ub, e
}
