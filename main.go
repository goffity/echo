package main

import (
	"net/http"
	"os"

	"github.com/op/go-logging"

	"github.com/labstack/echo"
)

type User struct {
	Username string `json:"username" form:"name"`
	Password string `json:"pwd" form:"pwd"`
}

var log = logging.MustGetLogger("API")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func main() {
	// For demo purposes, create two backend for os.Stderr.
	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)

	// For messages written to backend2 we want to add some additional
	// information to the output, including the used log level and the name of
	// the function.
	backend2Formatter := logging.NewBackendFormatter(backend2, format)

	// Only errors and more severe messages should be sent to backend1
	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(logging.ERROR, "")

	// Set the backends to be used.
	logging.SetBackend(backend1Leveled, backend2Formatter)

	e := echo.New()

	e.GET("/", index)
	e.POST("/login", login)

	log.Fatal(e.Start(":1234"))
}

func index(c echo.Context) error {
	return c.String(http.StatusOK, "Index")
}

func login(c echo.Context) error {
	log.Info("login")
	var ret error

	u := new(User)

	if err := c.Bind(u); err != nil {
		log.Panic(err)
		ret = err
	} else {
		ret = c.JSON(http.StatusOK, u)
	}

	return ret
}
