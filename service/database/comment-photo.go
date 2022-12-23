package database

import "time"

func (db *appdbimpl) CommentPhoto(userId uint64, photoId uint64, c Comment) (*Comment, error) {

	var username string
	var date = time.Now()
	errU := db.c.QueryRow(`SELECT username FROM users WHERE id=?`, userId).Scan(&username)
	if errU != nil {
		return nil, ErrUserNotExists
	}

	res, err := db.c.Exec(`INSERT INTO comments (id,userId,photoId,date,comment) VALUES (NULL, ?,?,?,?)`,
		userId, photoId, date, c.Comment)
	if err != nil {
		return &c, err
	}

	lastInsertID, err := res.LastInsertId()
	if err != nil {
		return &c, err
	}

	c.Id = uint64(lastInsertID)

	u := User{
		Username: username,
	}
	return &Comment{
		Id:       c.Id,
		User:     &u,
		Datetime: date,
		Comment:  c.Comment,
	}, nil
}

func (db *appdbimpl) DeleteComment(commentId uint64, userId uint64, photoId uint64) error {
	_, err := db.c.Exec(`DELETE FROM comments WHERE id=? AND userId=? AND photoId=?`, commentId, userId, photoId)
	if err != nil {
		return err
	}

	return nil
}
