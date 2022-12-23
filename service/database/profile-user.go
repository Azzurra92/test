package database

func (db *appdbimpl) GetUserProfile(userId uint64) (*Profile, error) {
	var username string
	errU := db.c.QueryRow(`SELECT username FROM users WHERE id=?`, userId).Scan(&username)
	if errU != nil {
		return nil, ErrUserNotExists
	}

	var cntP int
	errP := db.c.QueryRow(`SELECT COUNT(*) FROM photos WHERE userId=?`, userId).Scan(&cntP)
	if errP != nil {
		return nil, errP
	}

	var cntF int
	errF := db.c.QueryRow(`SELECT COUNT(*) FROM followers WHERE followerId=?`, userId).Scan(&cntF)
	if errF != nil {
		return nil, errF
	}

	var cntD int
	errD := db.c.QueryRow(`SELECT COUNT(*) FROM followers WHERE followedId=?`, userId).Scan(&cntD)
	if errD != nil {
		return nil, errD
	}

	stm, err := db.c.Prepare(`SELECT photoUrl,likes FROM photos WHERE userId=?`)
	if err != nil {
		return nil, err
	}
	rows, err := stm.Query(userId)
	if err != nil {
		return nil, err
	}

	photos := make([]Photo, 0)

	for rows.Next() {
		var pathPhoto string
		var likes uint64
		err := rows.Scan(&pathPhoto, &likes)
		if err != nil {
			return nil, err
		}

		p := Photo{
			Likes:    likes,
			PhotoUrl: pathPhoto,
		}
		photos = append(photos, p)
	}

	u := User{
		Username: username,
	}

	return &Profile{
		User:      &u,
		Post:      cntP,
		Follower:  cntD,
		Following: cntF,
		Photos:    photos,
	}, nil
}
