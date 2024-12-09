package database

import "github.com/rickb777/date/v2"

type RingtoneModel struct {
	Id   int
	Name string
}

type PhoneModel struct {
	Id   int
	Name string
}

type EffectModel struct {
	Id   int
	Name string
}

type UserModel struct {
	Id         int
	Name       string
	Email      string
	DateJoined date.Date `db:"date_joined"`
}
