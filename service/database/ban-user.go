package database

func (db *appdbimpl) BanUser(userId uint64, bannedUser uint64) error {

	_, err := db.c.Exec(`INSERT INTO bans (userId,bannedUser) VALUES (?, ?)`,
		userId, bannedUser)
	if err != nil {
		return err
	}

	return nil
}

func (db *appdbimpl) DeleteBan(userId uint64, bannedUser uint64) error {
	_, err := db.c.Exec(`DELETE FROM bans WHERE userId=? AND bannedUser=?`, userId, bannedUser)
	if err != nil {
		return err
	}

	return nil
}
