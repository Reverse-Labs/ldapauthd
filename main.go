package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	cli "github.com/jawher/mow.cli"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	app  = cli.App("dlb", "Dynamic Load Balancer")
	host = env("HOST", "0.0.0.0")
	port = env("PORT", "8080")

	baseDN       = env("BASE_DN", "cn=accounts,cn=users,dc=example,dc=com")
	bindDN       = env("BIND_DN", "cn=read-only-admin,dc=example,dc=com")
	ldapPort     = env("LDAP_PORT", "389")
	ldapHost     = env("LDAP_HOST", "ldap.example.com")
	bindPassword = env("BIND_PASSWORD", "password")
	ldapFilter   = env("LDAP_FILTER", "(uid=%s)")

	e = echo.New()
)

func getAuthentication(ctx echo.Context) error {
	return ctx.NoContent(http.StatusOK)
}

func main() {

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.BasicAuth(func(username, password string, ctx echo.Context) (bool, error) {
		err := authenticateLdap(username, password)
		return err == nil, err
	}))

	e.Any("*", getAuthentication)

	app.Command("serve", "serve proxy", cmdServe)

	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%s", host, port)))
	if err := app.Run(os.Args); err != nil {
		log.Println(err.Error())
	}
}

func cmdServe(cmd *cli.Cmd) {

	cmd.Action = func() {

		e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%s", host, port)))
	}
}

func env(key, fallback string) (value string) {
	value = os.Getenv(key)

	if len(value) < 1 {
		return fallback
	}

	return
}
