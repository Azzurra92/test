package database

func (db *appdbimpl) BanUser(userId uint64, bannedUser uint64) error {

	_, err := db.c.Exec(`INSERT INTO bans (userId,bannedUser) VALUES (?, ?)`,
		userId, bannedUser)
	if err != nil {
		return err
	}

	return nil
}
