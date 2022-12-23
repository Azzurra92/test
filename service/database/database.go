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
	"time"
)

var ErrUserExists = errors.New("user exists")
var ErrUserNotExists = errors.New("user not exists")
var ErrLikesExists = errors.New("The user has already liked")

type User struct {
	ID       uint64
	Username string
}

type Photo struct {
	Id       uint64
	Datetime time.Time
	UUID     string
	PhotoUrl string
	Likes    uint64
	UserId   uint64
}

type Profile struct {
	User      *User
	Photos    []Photo
	Post      int
	Follower  int
	Following int
}

type Comment struct {
	Id       uint64
	User     *User
	Datetime time.Time
	Comment  string
}

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	// CreateUser creates a new user if he/she doesn't exist
	CreateUser(User) (User, error)
	// UpdateUser updates the user, replacing every value with those provided in the argument
	UpdateUser(User) (User, error)
	// Insert and Delete ban user with the given ID
	BanUser(uint64, uint64) error
	DeleteBan(uint64, uint64) error
	// Insert and Delete follower user with the given ID
	FollowerUser(uint64, uint64) error
	DeleteFollowerUser(uint64, uint64) error
	// Gets Stream
	GetStream(uint64) ([]Photo, error)
	GetUserProfile(uint64) (*Profile, error)
	LikePhoto(uint64, uint64) error
	DeleteLike(uint64, uint64) error
	// Get Photo
	GetPhoto(uint64, uint64) (*Photo, error)
	// Insert Photo
	CreatePhoto(Photo) (Photo, error)
	// Delete Photo
	DeletePhoto(uint64, uint64) error
	CommentPhoto(uint64, uint64, Comment) (*Comment, error)
	DeleteComment(uint64, uint64, uint64) error
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
	errU := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users';`).Scan(&tableName)
	if errors.Is(errU, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE users (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL);`
		_, errU = db.Exec(sqlStmt)
		if errU != nil {
			return nil, fmt.Errorf("error creating database structure: %w", errU)
		}
	}
	errB := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='bans';`).Scan(&tableName)
	if errors.Is(errB, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE bans (
			userId INTEGER NOT NULL,
			bannedUser INTEGER NOT NULL,
			PRIMARY KEY(userId, bannedUser),
			FOREIGN KEY(userId) REFERENCES users(id),
			FOREIGN KEY(bannedUser) REFERENCES users(id));`
		_, errB = db.Exec(sqlStmt)
		if errB != nil {
			return nil, fmt.Errorf("error creating database structure: %w", errB)
		}
	}
	errF := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='followers';`).Scan(&tableName)
	if errors.Is(errF, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE followers (
			followerId INTEGER NOT NULL,
			followedId INTEGER NOT NULL,
			PRIMARY KEY(followerId, followedId),
			FOREIGN KEY(followerId) REFERENCES users(id),
			FOREIGN KEY(followedId) REFERENCES users(id));`
		_, errF = db.Exec(sqlStmt)
		if errF != nil {
			return nil, fmt.Errorf("error creating database structure: %w", errF)
		}
	}
	errP := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='photos';`).Scan(&tableName)
	if errors.Is(errP, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE photos (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			uuid TEXT NOT NULL,
			date TIMESTAMP NOT NULL,
			userId INTEGER NOT NULL,
			likes INTEGER NOT NULL,
			photoUrl TEXT NOT NULL,
			FOREIGN KEY(userId) REFERENCES users(id));`
		_, errP = db.Exec(sqlStmt)
		if errP != nil {
			return nil, fmt.Errorf("error creating database structure: %w", errP)
		}
	}
	errL := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='likes';`).Scan(&tableName)
	if errors.Is(errL, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE likes (
			userId INTEGER NOT NULL,
			photoId INTEGER NOT NULL,
			FOREIGN KEY(userId) REFERENCES users(id),
			FOREIGN KEY(photoId) REFERENCES photos(id));`
		_, errL = db.Exec(sqlStmt)
		if errL != nil {
			return nil, fmt.Errorf("error creating database structure: %w", errL)
		}
	}
	errC := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='comments';`).Scan(&tableName)
	if errors.Is(errC, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE comments (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			userId INTEGER NOT NULL,
			photoId INTEGER NOT NULL,
			date TIMESTAMP NOT NULL,
			comment TEXT NOT NULL,
			FOREIGN KEY(userId) REFERENCES users(id),
			FOREIGN KEY(photoId) REFERENCES photos(id));`
		_, errC = db.Exec(sqlStmt)
		if errC != nil {
			return nil, fmt.Errorf("error creating database structure: %w", errC)
		}
	}
	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
