package main

import (
	"net/http"
	"os"
	"webframework_echo/db"

	"github.com/labstack/echo"
	"github.com/op/go-logging"

	_ "github.com/go-sql-driver/mysql"
)

// login request
type loginRequest struct {
	Username string `json:"username" form:"name"`
	Password string `json:"password" form:"password"`
}

// Person is person object
type Person struct {
	ID        int    `json:"id" form:"id"`
	Firstname string `json:"firstname" form:"firstname"`
	Lastname  string `json:"lastname" form:"lastname"`
	Username  string `json:"username" form:"username"`
	Password  string `json:"password" form:"password"`
}

const (
	dbName = "gotest"
	dbPass = "root"
	dbHost = "localhost"
	dbPort = "3307"
)

var log = logging.MustGetLogger("API")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func initLog() {
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
}

func main() {

	initLog()

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

	u := new(loginRequest)
	person := new(Person)

	if err := c.Bind(u); err != nil {
		log.Error(err)
		ret = err
	} else {

		connection, errConnection := db.GetDatabaseConnection("root", "root", "127.0.0.1", "3307", "gotest")

		if errConnection != nil {
			log.Fatal(errConnection)
		}
		defer connection.Close()

		errConnection = connection.Ping()
		if errConnection != nil {
			log.Fatal(errConnection)
		}

		log.Debug("Username ", u.Username)
		log.Debug("Password ", u.Password)
		row := connection.QueryRow("SELECT * FROM person WHERE first_name=?", u.Username)

		selERR := row.Scan(&person.ID, &person.Firstname, &person.Lastname, &person.Username, &person.Password)

		if selERR != nil {
			log.Error("Query Error ", selERR)
		}
		log.Debug("person ", person)

		ret = c.JSON(http.StatusOK, person)
	}

	return ret
}
