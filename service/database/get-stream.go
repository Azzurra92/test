package database

import "time"

func (db *appdbimpl) GetStream(userId uint64) ([]Photo, error) {

	stm, err := db.c.Prepare("SELECT id, uuid, userId, date,likes, photoUrl FROM 'photos' WHERE userid IN (SELECT followedId FROM 'followers' WHERE followerId=?) AND userid NOT IN (SELECT bannedUser FROM 'bans' WHERE userId = ?) ORDER BY date DESC")
	if err != nil {
		return nil, err
	}

	rows, err := stm.Query(userId, userId)
	if err != nil {
		return nil, err
	}

	photos := make([]Photo, 0)

	for rows.Next() {
		var uuid string
		var photoUserId uint64
		var id uint64
		var datetime time.Time
		var likes uint64
		var photoUrl string

		err := rows.Scan(&id, &uuid, &photoUserId, &datetime, &likes, &photoUrl)
		if err != nil {
			return nil, err
		}
		p := Photo{
			Id:       id,
			UUID:     uuid,
			Datetime: datetime,
			UserId:   photoUserId,
			Likes:    likes,
			PhotoUrl: photoUrl,
		}
		photos = append(photos, p)

	}
	return photos, nil
}
