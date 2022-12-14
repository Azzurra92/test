package database

func (db *appdbimpl) UpdateUser(u User) (User, error) {
	res, err := db.c.Exec(`UPDATE users SET username=? WHERE id=?`,
		u.Username, u.ID)
	if err != nil {
		return u, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return u, err
	} else if affected == 0 {

		return u, ErrUserNotExists
	}
	return u, nil
}
