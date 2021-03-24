package main

import (
	"database/sql"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
)

type SessionHandler struct {
	oauth2   *oauth2.Config
	sessions *SessionRepository
	users    *UserRepository
}

func NewSessionHandler(oauthConfig *oauth2.Config, db *sql.DB) *SessionHandler {
	return &SessionHandler{
		users:    NewUserRepository(db),
		sessions: NewSessionRepository(db),
		oauth2:   oauthConfig}
}

func (s *SessionHandler) ValidateSession(c *fiber.Ctx) error {
	paths := strings.Split(c.Path(), "/")
	if len(paths) > 1 && (paths[1] == "session") {
		userId, e := url.QueryUnescape(c.Cookies("_userId_"))
		if e != nil {
			return c.SendStatus(http.StatusBadRequest)
		}
		sessionId := c.Cookies("_session_")
		if len(sessionId) == 0 {
			return c.SendStatus(http.StatusUnauthorized)
		}
		if len(userId) == 0 {
			return c.SendStatus(http.StatusUnauthorized)
		}
		sess, e := s.sessions.Get(&Session{UUID: sessionId, Email: userId})
		if e != nil {
			if e == sql.ErrNoRows {
				return c.SendStatus(http.StatusUnauthorized)
			}
			return LogError(e)
		}
		userInfo, e := s.users.Get(&User{ID: userId, Email: userId})
		if e != nil {
			return LogError(e)
		}
		c.Locals("accessToken", sess.Token)
		c.Locals("user", userInfo)
	}
	return c.Next()
}

func (s *SessionHandler) Activate(c *fiber.Ctx) (e error) {
	accessToken := c.Locals("accessToken").(string)
	user := c.Locals("user").(*User)
	oauthClient := s.oauth2.Client(c.Context(), &oauth2.Token{
		AccessToken:  accessToken,
		RefreshToken: user.RefreshToken,
	})
	userInfoRes, e := oauthClient.Get("https://www.googleapis.com/oauth2/v1/userinfo?alt=json")
	if e != nil {
		return LogError(e)
	}
	defer userInfoRes.Body.Close()
	userInfo, e := ioutil.ReadAll(userInfoRes.Body)
	if e != nil {
		return LogError(e)
	}
	c.Response().Header.Set("Content-Type", "application/json")
	return c.Send(userInfo)
}

func (s *SessionHandler) Logout(c *fiber.Ctx) (e error) {
	sessionId := c.Cookies("_session_", "")
	userId, e := url.QueryUnescape(c.Cookies("_userId_", ""))
	if e != nil {
		return e
	}
	e = s.sessions.Remove(&Session{
		UUID:  sessionId,
		Email: userId})
	if e != nil {
		_ = LogError(e)
	}
	return e
}
