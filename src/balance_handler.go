package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/gofiber/fiber/v2"
)

type BalanceHandler struct {
	balance *BalanceRepository
	session *SessionRepository
}

func NewBalanceHandler(db *sql.DB) *BalanceHandler {
	return &BalanceHandler{
		session: NewSessionRepository(db),
		balance: NewBalanceRepository(db)}
}

func (b *BalanceHandler) GetBalances(c *fiber.Ctx) (e error) {
	session := c.Cookies("_session_", "")
	userId, e := url.QueryUnescape(c.Cookies("_userId_", ""))
	if e != nil {
		return c.SendStatus(http.StatusBadRequest)
	}
	if len(session) == 0 || len(userId) == 0 {
		return c.SendStatus(http.StatusUnauthorized)
	}
	sessionDetails, e := b.session.Get(&Session{UUID: session, Email: userId})
	if e != nil {
		return LogError(nil)
	}
	if len(sessionDetails.Email) == 0 || sessionDetails.Email != userId {
		return c.SendStatus(http.StatusUnauthorized)
	}
	balances, e := b.balance.GetAll(&User{ID: sessionDetails.UserId, Email: sessionDetails.Email})
	if e != nil {
		return LogError(e)
	}
	balancesJSON, e := json.Marshal(balances)
	if e != nil {
		return LogError(e)
	}
	return c.Send(balancesJSON)
}

func (b *BalanceHandler) AddBalance(c *fiber.Ctx) (e error) {
	session := c.Cookies("_session_", "")
	userId, e := url.QueryUnescape(c.Cookies("_userId_", ""))
	if e != nil {
		return c.SendStatus(http.StatusBadRequest)
	}
	if len(session) == 0 || len(userId) == 0 {
		return c.SendStatus(http.StatusUnauthorized)
	}
	sessionDetails, e := b.session.Get(&Session{UUID: session, Email: userId})
	if e != nil {
		return LogError(e)
	}
	if len(sessionDetails.Email) == 0 || sessionDetails.Email != userId {
		return c.SendStatus(http.StatusUnauthorized)
	}
	var newBal Balance
	e = json.Unmarshal(c.Body(), &newBal)
	if e != nil {
		return LogError(e)
	}
	if len(newBal.Type) == 0 || len(newBal.Name) == 0 {
		return c.SendStatus(http.StatusBadRequest)
	}
	if newBal.Balance < 0 {
		return c.SendStatus(http.StatusBadRequest)
	}
	newBalanceRow := &Balance{
		Email:   sessionDetails.Email,
		Type:    newBal.Type,
		Name:    newBal.Name,
		Balance: newBal.Balance,
		Created: newBal.Created,
	}
	e = b.balance.Add(newBalanceRow)
	if e != nil {
		return LogError(e)
	}
	newBalJSON, e := json.Marshal(newBalanceRow)
	if e != nil {
		return LogError(e)
	}
	return c.Send(newBalJSON)
}

func (b *BalanceHandler) DeleteBalance(c *fiber.Ctx) (e error) {
	session := c.Cookies("_session_", "")
	userId, e := url.QueryUnescape(c.Cookies("_userId_", ""))
	if e != nil {
		return c.SendStatus(http.StatusBadRequest)
	}
	if len(session) == 0 || len(userId) == 0 {
		return c.SendStatus(http.StatusUnauthorized)
	}
	sessionDetails, e := b.session.Get(&Session{UUID: session, Email: userId})
	if e != nil {
		return LogError(e)
	}
	if len(sessionDetails.Email) == 0 || sessionDetails.Email != userId {
		return c.SendStatus(http.StatusUnauthorized)
	}
	var newBal Balance
	e = c.BodyParser(&newBal)
	if e != nil {
		return LogError(e)
	}
	if len(newBal.ID) == 0 {
		return c.SendStatus(http.StatusBadRequest)
	}
	e = b.balance.Remove(&Balance{
		ID:    newBal.ID,
		Email: sessionDetails.Email})
	if e != nil {
		return LogError(e)
	}
	return e
}
