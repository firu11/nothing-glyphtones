package database

import (
	"database/sql"
	"log"
	"math"
	"strings"

	"github.com/blockloop/scan/v2"
	"github.com/lib/pq"
)

var resultsPerPage int = 20

func GetRingtones(search string, phones []int, effects []int, page int) ([]RingtoneModel, int, error) {
	var ringtones []RingtoneModel
	var rows *sql.Rows
	var err error

	rows, err = DB.Query(`WITH ringtones_matched AS ( SELECT id, name, phone, effect FROM ringtone WHERE LOWER(name) LIKE '%' || $1 || '%' AND phone = ANY ($2) AND effect = ANY ($3) ) SELECT rm.id, rm.name, p.name as phone_name, e.name as effect_name, COUNT(*) OVER () as results FROM ringtones_matched rm INNER JOIN phone p ON rm.phone = p.id INNER JOIN effect e ON rm.effect = e.id LIMIT $4 OFFSET $5;`, search, pq.Array(phones), pq.Array(effects), resultsPerPage, (page-1)*resultsPerPage)
	if err != nil {
		return ringtones, 0, err
	}

	err = scan.Rows(&ringtones, rows)
	if err != nil {
		return ringtones, 0, err
	}

	var numberOfPages int = 0
	if len(ringtones) != 0 {
		numberOfPages = int(math.Ceil(float64(ringtones[0].NumberOfResults) / float64(resultsPerPage)))
	}

	return ringtones, numberOfPages, nil
}

func GetPhones() ([]PhoneModel, error) {
	var phones []PhoneModel
	var rows *sql.Rows
	var err error

	rows, err = DB.Query(`SELECT * FROM phone ORDER BY name;`)
	if err != nil {
		log.Println(err.Error())
		return phones, err
	}

	err = scan.Rows(&phones, rows)
	if err != nil {
		log.Println(2)
		return phones, err
	}

	return phones, nil
}

func GetEffects() ([]EffectModel, error) {
	var effects []EffectModel
	var rows *sql.Rows
	var err error

	rows, err = DB.Query(`SELECT * FROM effect ORDER BY id;`)
	if err != nil {
		log.Println(err.Error())
		return effects, err
	}

	err = scan.Rows(&effects, rows)
	if err != nil {
		log.Println(2)
		return effects, err
	}

	return effects, nil
}

func GetUser(id int) (UserModel, error) {
	var user UserModel
	row, err := DB.Query(`SELECT * FROM user WHERE id = $1 AND NOT smazany;`, id)
	if err != nil {
		return user, err
	}
	err = scan.Row(&user, row)
	return user, err
}

func CreateUser(name string, email string) (int, error) {
	email = strings.ToLower(email)

	var userID int
	var deleted bool
	err := DB.QueryRow(`SELECT id, deleted FROM "user" WHERE email = $1;`, email).Scan(&userID, &deleted)

	if err == nil && !deleted { // already in the db
		return userID, nil

	} else if err == nil && deleted { // in the db but deleted
		_, err := DB.Exec(`UPDATE "user" SET deleted = false, name = $1 RETURNING id;`, name)
		return userID, err

	} else if err == sql.ErrNoRows { // not in the db
		err = DB.QueryRow(`INSERT INTO "user" (name, email) VALUES ($1, $2) RETURNING id;`, name, email).Scan(&userID)
		if err != nil {
			return 0, err
		}
		return userID, nil

	} else { // other error than NoRows
		return 0, err
	}
}
