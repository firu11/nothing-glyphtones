package database

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/blockloop/scan/v2"
	"github.com/lib/pq"
)

var resultsPerPage int = 10

func GetRingtones(search string, phones []int, effects []int, page int) ([]RingtoneModel, int, error) {
	var ringtones []RingtoneModel
	var rows *sql.Rows
	var err error

	log.Println(len(phones) == 0 && len(effects) == 0)
	if len(phones) == 0 && len(effects) == 0 {
		rows, err = DB.Query(`WITH ringtones_matched AS ( SELECT id, name, phone, effect, downloads, ( similarity (name, $1) * 0.7 + ( downloads::FLOAT / MAX(downloads) OVER () ) * 0.1 - not_working::FLOAT / MAX(not_working) OVER () * 0.2 ) AS score FROM ringtone WHERE similarity (name, $1) > 0.05 ) SELECT rm.id, rm.name, rm.score, rm.downloads, p.name AS phone_name, e.name AS effect_name, COUNT(*) OVER () AS results FROM ringtones_matched rm INNER JOIN phone p ON rm.phone = p.id INNER JOIN effect e ON rm.effect = e.id ORDER BY score DESC LIMIT $2 OFFSET $3;`, search, resultsPerPage, (page-1)*resultsPerPage)
	} else if len(phones) != 0 && len(effects) != 0 {
		rows, err = DB.Query(`WITH ringtones_matched AS ( SELECT id, name, phone, effect FROM ringtone WHERE LOWER(name) LIKE '%' || LOWER($1) || '%' AND phone = ANY ($2) AND effect = ANY ($3) ) SELECT rm.id, rm.name, p.name as phone_name, e.name as effect_name, COUNT(*) OVER () as results FROM ringtones_matched rm INNER JOIN phone p ON rm.phone = p.id INNER JOIN effect e ON rm.effect = e.id LIMIT $4 OFFSET $5;`, search, pq.Array(phones), pq.Array(effects), resultsPerPage, (page-1)*resultsPerPage)
	} else if len(effects) != 0 {
		rows, err = DB.Query(`WITH ringtones_matched AS ( SELECT id, name, phone, effect FROM ringtone WHERE LOWER(name) LIKE '%' || LOWER($1) || '%' effect = ANY ($2) ) SELECT rm.id, rm.name, p.name as phone_name, e.name as effect_name, COUNT(*) OVER () as results FROM ringtones_matched rm INNER JOIN phone p ON rm.phone = p.id INNER JOIN effect e ON rm.effect = e.id LIMIT $3 OFFSET $4;`, search, pq.Array(effects), resultsPerPage, (page-1)*resultsPerPage)
	} else /* if len(phones) != 0 */ {
		rows, err = DB.Query(`WITH ringtones_matched AS ( SELECT id, name, phone, effect FROM ringtone WHERE LOWER(name) LIKE '%' || LOWER($1) || '%' phone = ANY ($2) ) SELECT rm.id, rm.name, p.name as phone_name, e.name as effect_name, COUNT(*) OVER () as results FROM ringtones_matched rm INNER JOIN phone p ON rm.phone = p.id INNER JOIN effect e ON rm.effect = e.id LIMIT $3 OFFSET $4;`, search, pq.Array(phones), resultsPerPage, (page-1)*resultsPerPage)
	}
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

func GetPopularRingtones(page int) ([]RingtoneModel, int, error) {
	var ringtones []RingtoneModel
	var rows *sql.Rows
	var err error

	rows, err = DB.Query(`WITH ringtones_matched AS ( SELECT id, name, phone, effect, author_id, downloads, (downloads::FLOAT - 2 * not_working::FLOAT) / MAX(downloads) OVER () AS score FROM ringtone ) SELECT rm.id, rm.name, rm.score, u.id AS author_id, u.name AS author_name, rm.downloads, p.name AS phone_name, e.name AS effect_name, COUNT(*) OVER () AS results FROM ringtones_matched rm INNER JOIN phone p ON rm.phone = p.id INNER JOIN effect e ON rm.effect = e.id INNER JOIN "user" u ON rm.author_id = u.id ORDER BY score DESC LIMIT $1 OFFSET $2;`, resultsPerPage, (page-1)*resultsPerPage)
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

func GetRingtonesByUser(userID int, page int) ([]RingtoneModel, int, error) {
	var ringtones []RingtoneModel
	var rows *sql.Rows
	var err error

	rows, err = DB.Query(`WITH ringtones_matched AS ( SELECT id, name, phone, effect, downloads, downloads::FLOAT / MAX(downloads) OVER () - 5 * not_working::FLOAT / MAX(not_working) OVER () AS score FROM ringtone WHERE author_id = $1 ) SELECT rm.id, rm.name, rm.score, rm.downloads, p.name AS phone_name, e.name AS effect_name, COUNT(*) OVER () AS results FROM ringtones_matched rm INNER JOIN phone p ON rm.phone = p.id INNER JOIN effect e ON rm.effect = e.id ORDER BY score DESC LIMIT $2 OFFSET $3;`, userID, resultsPerPage, (page-1)*resultsPerPage)
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

func CreateRingtone(name string, phone int, effect int, authorID int) (int, error) {
	var ringtoneID int
	err := DB.QueryRow(`INSERT INTO ringtone (name, phone, effect, author_id) VALUES ($1, $2, $3, $4) RETURNING id;`, name, phone, effect, authorID).Scan(&ringtoneID)
	if err != nil {
		return 0, err
	}
	return ringtoneID, nil
}

func RingtoneIncreaseDownload(id int) (string, error) {
	var name string
	var phone string
	err := DB.QueryRow(`UPDATE ringtone SET downloads = downloads + 1 WHERE id = $1 RETURNING name, ( SELECT p.name FROM phone p INNER JOIN ringtone r ON p.id = r.phone WHERE r.id = $1 );`, id).Scan(&name, &phone)
	return fmt.Sprintf("%s - %s.ogg", name, phone), err
}

func RingtoneIncreaseNotWorking(id int) error {
	_, err := DB.Exec(`UPDATE ringtone SET not_working = not_working + 1 WHERE id = $1;`, id)
	return err
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
	row, err := DB.Query(`SELECT * FROM "user" WHERE id = $1 AND NOT deleted;`, id)
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
		err = DB.QueryRow(`UPDATE "user" SET deleted = false, name = $1 RETURNING id;`, name).Scan(&userID)
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

func RenameUser(id int, newName string) (string, error) {
	var email string
	err := DB.QueryRow(`UPDATE "user" SET name = $1 WHERE id = $2 RETURNING email;`, newName, id).Scan(&email)
	return email, err
}
