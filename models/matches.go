package models

import (
	"ridenow/matcher"
	"time"
)

type Float32Range struct {
	Lower float32
	Upper float32
}

type Int32Range struct {
	Lower int32
	Upper int32
}

type User struct {
	Id               int64
	Username         string
	Name             string
	Surname          string
	Email            string
	WaveHeightRange  Float32Range
	AllowedTimeRange Int32Range
	Created          time.Time
}

type Location struct {
	Id   int64
	Name string
}

type Match struct {
	User        *User
	Location    *Location
	WaveHeightM float64
	Time        time.Time
	Created     time.Time
}

func (db *DB) MatchUsers(fc *matcher.Forecast) ([]*Match, error) {
	query := `SELECT up.id, up.username, up.name, up.surname, up.email, user_location.location_id
	          FROM user_profile AS up
              JOIN user_location ON up.id = user_location.user_profile_id 
              WHERE user_location.location_id = $1 AND $2::numeric <@ up.wave_height_range;`
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(fc.GetLocationId(), fc.GetWaveHeightM())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	matches := make([]*Match, 0)
	for rows.Next() {
		loc := new(Location)
		user := new(User)
		err := rows.Scan(&user.Id, &user.Username, &user.Name, &user.Surname, &user.Email, &loc.Id)
		if err != nil {
			return nil, err
		}
		match := &Match{
			User:        user,
			Location:    loc,
			WaveHeightM: fc.GetWaveHeightM(),
			Time:        time.Unix(0, fc.GetTime()), // unix nsec -> time
			Created:     time.Now()}
		matches = append(matches, match)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return matches, nil
}
