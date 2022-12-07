package database

func (db *appdbimpl) CreateUser(u User) (User, error) {
	tx, err := db.c.Begin()
	if err != nil {
		return u, err
	}

	var cnt int
	err = tx.QueryRow(`SELECT COUNT(*) FROM users WHERE username=?`, u.Username).Scan(&cnt)
	if err != nil {
		_ = tx.Rollback()
		return u, err
	} else if cnt > 0 {
		_ = tx.Rollback()
		return u, ErrUserExists
	}

	_, err = tx.Exec(`INSERT INTO users (id,username) VALUES (?, ?)`,
		u.ID, u.Username)
	if err != nil {
		_ = tx.Rollback()
		return u, err
	}
	return u, tx.Commit()
}
