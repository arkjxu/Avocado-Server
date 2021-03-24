package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
)

type UserHandler struct {
	oauth2   *oauth2.Config
	users    *UserRepository
	sessions *SessionRepository
}

type LoginRequest struct {
	Code      string `json:"code"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Avatar        string `json:"picture"`
	Locale        string `json:"locale"`
}

func NewUserHandler(oauthConfig *oauth2.Config, db *sql.DB) *UserHandler {
	return &UserHandler{
		users:    NewUserRepository(db),
		sessions: NewSessionRepository(db),
		oauth2:   oauthConfig}
}

func (u *UserHandler) Join(c *fiber.Ctx) (e error) {
	code := c.Query("code")
	if len(code) == 0 {
		return c.SendStatus(http.StatusUnauthorized)
	}
	t, e := u.oauth2.Exchange(context.Background(), code)
	if e != nil {
		return LogError(e)
	}
	oauthAPI := u.oauth2.Client(c.Context(), &oauth2.Token{AccessToken: t.AccessToken, RefreshToken: t.RefreshToken})
	oauthUserInfo, e := oauthAPI.Get("https://www.googleapis.com/oauth2/v1/userinfo?alt=json")
	if e != nil {
		return LogError(e)
	}
	defer oauthUserInfo.Body.Close()
	if oauthUserInfo.StatusCode != 200 {
		return c.SendStatus(http.StatusUnauthorized)
	}
	var gUser GoogleUserInfo
	e = json.NewDecoder(oauthUserInfo.Body).Decode(&gUser)
	if e != nil {
		return LogError(e)
	}
	existingUser, e := u.users.Get(&User{ID: gUser.Email, Email: gUser.Email})
	if e != nil && e != sql.ErrNoRows {
		return LogError(e)
	}
	if existingUser == nil || len(existingUser.Email) == 0 {
		e = u.users.Add(&User{
			Email:        gUser.Email,
			FirstName:    gUser.GivenName,
			LastName:     gUser.FamilyName,
			RefreshToken: t.RefreshToken,
		})
		if e != nil {
			return LogError(e)
		}
	}
	_, e = u.sessions.AddToken(&Session{Email: gUser.Email, Token: t.AccessToken})
	if e != nil {
		return LogError(e)
	}
	newSessionId, e := GetRandomUUID()
	if e != nil {
		return LogError(e)
	}
	e = u.sessions.Add(&Session{
		UUID:  newSessionId,
		Email: gUser.Email,
		Token: t.AccessToken,
	})
	if e != nil {
		return LogError(e)
	}
	userJson, e := json.Marshal(gUser)
	if e != nil {
		return LogError(e)
	}
	c.Cookie(&fiber.Cookie{
		Name:   "_session_",
		Value:  newSessionId,
		MaxAge: int((7 * 24 * time.Hour) / time.Millisecond),
	})
	c.Send(userJson)
	return e
}
