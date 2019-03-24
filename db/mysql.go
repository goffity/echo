package db

import (
	"database/sql"
	"fmt"
)

func GetDatabaseConnection(username, password, host, port, databaseName string) (*sql.DB, error) {
	var connectionString = ""

	if host == "" && port == "" {
		connectionString = fmt.Sprintf("%s:%s/%s", username, password, databaseName)
	} else {
		connectionString = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, databaseName)
	}

	connection, err := sql.Open("mysql", connectionString)

	return connection, err
}
