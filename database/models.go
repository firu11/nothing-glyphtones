package database

import "github.com/rickb777/date/v2"

type RingtoneModel struct {
	Id              int    `db:"id"`
	Name            string `db:"name"`
	PhoneName       string `db:"phone_name"`
	EffectName      string `db:"effect_name"`
	Downloaded      int    `db:"downloaded"`
	NumberOfResults int    `db:"results"`
}

type PhoneModel struct {
	Id       int    `db:"id"`
	Name     string `db:"name"`
	Selected bool   `db:""`
}

type EffectModel struct {
	Id       int    `db:"id"`
	Name     string `db:"name"`
	Selected bool   `db:""`
}

type UserModel struct {
	Id         int       `db:"id"`
	Name       string    `db:"name"`
	Email      string    `db:"email"`
	DateJoined date.Date `db:"date_joined"`
}
