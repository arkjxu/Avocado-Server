package main

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type User struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	RefreshToken string `json:"token"`
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) Add(user *User) (e error) {
	ctx, cnl := context.WithTimeout(context.Background(), 3*time.Second)
	defer cnl()
	_, e = u.db.ExecContext(ctx, "CALL Users_Add(?,?,?,?)", user.Email, user.FirstName, user.LastName, user.RefreshToken)
	return e
}

func (u *UserRepository) Remove(user *User) (e error) {
	return errors.New("Not Implemented")
}

func (u *UserRepository) Get(user *User) (usr *User, e error) {
	ctx, cnl := context.WithTimeout(context.Background(), 3*time.Second)
	defer cnl()
	row := u.db.QueryRowContext(ctx, "CALL Users_Get(?)", user.Email)
	if row.Err() != nil && row.Err() != sql.ErrNoRows {
		return usr, row.Err()
	}
	if row.Err() == sql.ErrNoRows {
		return usr, e
	}
	usr = &User{}
	e = row.Scan(&usr.ID, &usr.Email, &usr.FirstName, &usr.LastName, &usr.RefreshToken)
	return usr, e
}
