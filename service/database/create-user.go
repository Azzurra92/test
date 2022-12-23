package database

func (db *appdbimpl) CreateUser(u User) (User, error) {

	var cnt int
	err := db.c.QueryRow(`SELECT COUNT(*) FROM users WHERE username=?`, u.Username).Scan(&cnt)
	if err != nil {
		return u, err
	}

	if cnt > 0 {
		return u, ErrUserExists
	}

	res, err := db.c.Exec(`INSERT INTO users (id,username) VALUES (NULL, ?)`,
		u.Username)
	if err != nil {
		return u, err
	}

	lastInsertID, err := res.LastInsertId()
	if err != nil {
		return u, err
	}

	u.ID = uint64(lastInsertID)
	return u, nil
}
