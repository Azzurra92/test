package database

func (db *appdbimpl) FollowerUser(followerId uint64, followedId uint64) error {

	_, err := db.c.Exec(`INSERT INTO followers (followerId,followedId) VALUES (?, ?)`,
		followerId, followedId)
	if err != nil {
		return err
	}

	return nil
}

func (db *appdbimpl) DeleteFollowerUser(followerId uint64, followedId uint64) error {

	_, err := db.c.Exec(`DELETE FROM followers WHERE followerId=? AND followedId=?`, followerId, followedId)
	if err != nil {
		return err
	}

	return nil
}
