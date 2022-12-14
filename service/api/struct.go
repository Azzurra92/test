package api

import (
	"regexp"
	"sapienza/azzurra/wasaphoto/service/database"
	"time"
)

var usernameRx = regexp.MustCompile(`^[A-Za-z0-9_-]*$`)

func (u *User) IsValid() bool {
	return len(u.Username) >= 3 && len(u.Username) <= 16 && usernameRx.MatchString(u.Username)
}

type User struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
}

func (u *User) ToDatabase() database.User {
	return database.User{
		ID:       u.ID,
		Username: u.Username,
	}
}

func (u *User) FromDatabase(user database.User) {
	u.ID = user.ID
}

type Profile struct {
	User      *User   `json:"user"`
	Photos    []Photo `json:"photos"`
	Post      int     `json:"post"`
	Follower  int     `json:"follower"`
	Following int     `json:"following"`
}

type Photo struct {
	Id       int               `json:"user"`
	Datetime time.Time         `json:"datetime"`
	Likes    int               `json:"like"`
	Comments []CommentResponse `json:"comments"`
	PhotoUrl string            `json:"photourl"`
}

type CommentResponse struct {
	Id       int       `json:"id"`
	From     *User     `json:"from"`
	Comment  string    `json:"comment"`
	Datetime time.Time `json:"datetime"`
}

type CommentRequest struct {
	Text string `json:"text"`
}

type Stream struct {
	Photos []Photo `json:"photos"`
}

type ApiResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func parseDate(t time.Time) time.Time {
	s, _ := time.Parse(time.RFC3339, time.RFC3339)
	return s
}
