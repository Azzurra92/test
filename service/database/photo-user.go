package database

import "time"

func (db *appdbimpl) CreatePhoto(p Photo) (Photo, error) {

	res, err := db.c.Exec(`INSERT INTO photos (id,date,userid,uuid,likes, photourl) VALUES (NULL, ?,?,?,?,?)`,
		p.Datetime, p.UserId, p.UUID, p.Likes, p.PhotoUrl)
	if err != nil {
		return p, err
	}

	lastInsertID, err := res.LastInsertId()
	if err != nil {
		return p, err
	}

	p.Id = uint64(lastInsertID)

	return p, nil
}

func (db *appdbimpl) DeletePhoto(userId uint64, photoId uint64) error {
	_, err := db.c.Exec(`DELETE FROM photos WHERE userid = ? AND id = ?`, userId, photoId)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) GetPhoto(userid uint64, id uint64) (*Photo, error) {

	stm, err := db.c.Prepare("SELECT uuid,date,photoUrl,likes FROM photos WHERE userid=? AND id = ?")
	if err != nil {
		return nil, err
	}

	var uuid string
	var likes uint64
	date := time.Now()
	photoUrl := ""
	if err := stm.QueryRow(userid, id).Scan(&uuid, &date, &photoUrl, &likes); err != nil {
		return nil, err
	}
	return &Photo{
		Id:       id,
		UUID:     uuid,
		Datetime: date,
		UserId:   userid,
		Likes:    likes,
		PhotoUrl: photoUrl,
	}, nil

}
