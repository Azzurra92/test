/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

var ErrUserExists = errors.New("user exists")
var ErrUserNotExists = errors.New("user not exists")

type User struct {
	ID       uint64
	Username string
}

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	// CreateUser creates a new user if he/she doesn't exist
	CreateUser(User) (User, error)
	// UpdateUser updates the user, replacing every value with those provided in the argument
	UpdateUser(User) (User, error)

	BanUser(uint64, uint64) error
	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}
	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE users (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}
	err1 := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='bans';`).Scan(&tableName)
	if errors.Is(err1, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE bans (
			userId INTEGER NOT NULL,
			bannedUser INTEGER NOT NULL,
			PRIMARY KEY(userId, bannedUser),
			FOREIGN KEY(userId) REFERENCES users(id),
			FOREIGN KEY(bannedUser) REFERENCES users(id));`
		_, err1 = db.Exec(sqlStmt)
		if err1 != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err1)
		}
	}
	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
