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

func (p *Photo) ToDatabase() database.Photo {

	return database.Photo{
		Id:       p.Id,
		UUID:     p.UUID,
		Datetime: p.Datetime,
		Likes:    p.Likes,
		PhotoUrl: p.PhotoUrl,
		UserId:   p.UserId,
	}
}

func (p *Photo) FromDatabase(d database.Photo) {

	p.Id = d.Id
	p.UUID = d.UUID
	p.Datetime = d.Datetime
	p.Likes = d.Likes
	p.PhotoUrl = d.PhotoUrl
	p.UserId = d.UserId

}

func (c *CommentRequest) ToDatabase() database.Comment {

	return database.Comment{
		Comment: c.Text,
	}
}

type Profile struct {
	User      *User   `json:"user"`
	Photos    []Photo `json:"photos"`
	Post      int     `json:"post"`
	Follower  int     `json:"follower"`
	Following int     `json:"following"`
}

type Photo struct {
	Id       uint64    `json:"id"`
	UUID     string    `json:"uuid"`
	Datetime time.Time `json:"datetime"`
	UserId   uint64    `json:"userid"`
	Likes    uint64    `json:"likes"`
	PhotoUrl string    `json:"photourl"`
}

type CommentResponse struct {
	Id       uint64    `json:"id"`
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
