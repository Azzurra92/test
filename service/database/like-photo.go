package database

func (db *appdbimpl) LikePhoto(userId uint64, photoId uint64) error {

	var resL int
	errL := db.c.QueryRow(`SELECT COUNT(*) FROM likes WHERE userId=? AND photoId=?`, userId, photoId).Scan(&resL)
	if errL != nil {
		return errL
	}
	if resL > 0 {
		return ErrLikesExists
	}

	_, err := db.c.Exec(`INSERT INTO likes (userId,photoId) VALUES (?, ?)`,
		userId, photoId)
	if err != nil {
		return err
	}

	var cnt int
	errP := db.c.QueryRow(`SELECT COUNT(*) FROM likes WHERE photoId=?`, photoId).Scan(&cnt)
	if errP != nil {
		return errP
	}

	_, errU := db.c.Exec(`UPDATE photos SET likes=? WHERE id=?`, &cnt, photoId)
	if errU != nil {
		return errU
	}
	return nil
}

func (db *appdbimpl) DeleteLike(userId uint64, photoId uint64) error {
	_, err := db.c.Exec(`DELETE FROM likes WHERE userId=? AND photoId=?`, userId, photoId)
	if err != nil {
		return err
	}

	var cnt int
	err = db.c.QueryRow(`SELECT COUNT(*) FROM likes WHERE photoId=?`, photoId).Scan(&cnt)
	if err != nil {
		return err
	}

	res, err := db.c.Exec(`UPDATE photos SET likes=? WHERE id=?`, &cnt, photoId)
	if err != nil {
		return err
	}

	_, errR := res.RowsAffected()
	if errR != nil {
		return errR
	}

	return nil
}
