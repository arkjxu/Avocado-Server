package main

import (
	"context"
	"database/sql"
	"time"
)

type Session struct {
	ID      string `json:"id"`
	UserId  string `json:"userId"`
	UUID    string `json:"uuid"`
	Email   string `json:"email"`
	Token   string `json:"token"`
	Created string `json:"created"`
}

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (s *SessionRepository) Add(sess *Session) (e error) {
	ctx, cnl := context.WithTimeout(context.Background(), 3*time.Second)
	defer cnl()
	_, e = s.db.ExecContext(ctx, "CALL Session_Add(?,?)", sess.Email, sess.UUID)
	return e
}

func (s *SessionRepository) AddToken(sess *Session) (id int, e error) {
	ctx, cnl := context.WithTimeout(context.Background(), 3*time.Second)
	defer cnl()
	row := s.db.QueryRowContext(ctx, "CALL Session_Token_Add(?,?)", sess.Email, sess.Token)
	if row.Err() != nil {
		return id, row.Err()
	}
	row.Scan(&id)
	return id, e
}

func (b *SessionRepository) Remove(sess *Session) (e error) {
	ctx, cnl := context.WithTimeout(context.Background(), 3*time.Second)
	defer cnl()
	_, e = b.db.ExecContext(ctx, "CALL Session_Delete(?,?)", sess.Email, sess.UUID)
	return e
}

func (b *SessionRepository) Get(sess *Session) (u *Session, e error) {
	ctx, cnl := context.WithTimeout(context.Background(), 6*time.Second)
	defer cnl()
	row := b.db.QueryRowContext(ctx, "CALL Session_Get(?,?)", sess.Email, sess.UUID)
	if row.Err() != nil {
		return u, row.Err()
	}
	u = &Session{}
	e = row.Scan(&u.ID, &u.UserId, &u.UUID, &u.Email, &u.Token, &u.Created)
	return u, e
}
