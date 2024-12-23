package database

import "github.com/rickb777/date/v2"

type RingtoneModel struct {
	ID              int    `db:"id"`
	Name            string `db:"name"`
	PhoneName       string `db:"phone_name"`
	EffectName      string `db:"effect_name"`
	Downloads       int    `db:"downloads"`
	AuthorName      string `db:"author_name"`
	AuthorID        int    `db:"author_id"`
	NumberOfResults int    `db:"results"`
}

type PhoneModel struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	Selected bool   `db:""`
}

type EffectModel struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	Selected bool   `db:""`
}

type UserModel struct {
	ID         int       `db:"id"`
	Name       string    `db:"name"`
	Email      string    `db:"email"`
	DateJoined date.Date `db:"date_joined"`
}
