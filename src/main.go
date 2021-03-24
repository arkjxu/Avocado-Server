package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	oauthConfigs := &oauth2.Config{
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	mysql, e := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_KEY"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_DB")))
	if e != nil {
		log.Fatal(e)
	}
	defer mysql.Close()

	userHandler := NewUserHandler(oauthConfigs, mysql)
	sessionHandler := NewSessionHandler(oauthConfigs, mysql)
	balanceHandler := NewBalanceHandler(mysql)

	app := fiber.New(fiber.Config{
		ReadTimeout: 12 * time.Second,
	})
	defer app.Shutdown()

	app.Use(cors.New())
	app.Use(sessionHandler.ValidateSession)

	app.Post("/user/authorize", userHandler.Join)
	app.Get("/session/activate", sessionHandler.Activate)
	app.Post("/session/logout", sessionHandler.Logout)

	app.Get("/balances", balanceHandler.GetBalances)
	app.Post("/balances", balanceHandler.AddBalance)
	app.Delete("/balances", balanceHandler.DeleteBalance)

	e = app.Listen(fmt.Sprintf(":%s", os.Getenv("PORT")))
	if e != nil {
		log.Fatal(e)
	}
}
