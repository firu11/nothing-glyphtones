package database

import (
	"database/sql"
	"log"
	"math"
	"strings"

	"github.com/blockloop/scan/v2"
	"github.com/lib/pq"
)

var resultsPerPage int = 10

func GetRingtones(search string, category int, sortBy string, phones []int, effects []int, page int) ([]RingtoneModel, int, error) {
	var ringtones []RingtoneModel
	var rows *sql.Rows
	var err error

	query := `WITH all_ringtones AS ( SELECT id, name, ARRAY( SELECT p.name FROM phone_and_ringtone par INNER JOIN phone p ON p.id = par.phone_id WHERE par.ringtone_id = ringtone.id ) as phone_names, ARRAY( SELECT par.phone_id FROM phone_and_ringtone par WHERE par.ringtone_id = ringtone.id ) as phone_ids, effect_id, author_id, downloads, time_added, category, ( CASE WHEN $1 = '' THEN ( downloads::FLOAT - 2 * not_working::FLOAT ) / GREATEST(MAX(downloads) OVER (), 1) ELSE similarity (name, $1) * 0.7 + ( downloads::FLOAT / GREATEST(MAX(downloads) OVER (), 1) ) * 0.1 - not_working::FLOAT / GREATEST( MAX(not_working) OVER () * 0.2, 1 ) END ) AS score FROM ringtone ), ringtones_matched AS ( SELECT * FROM all_ringtones WHERE ( similarity (name, $1) > 0.05 OR $1 = '' ) AND ( phone_ids && $2 OR COALESCE(array_length($2, 1), 0) = 0 ) AND ( effect_id = ANY ($3) OR COALESCE(array_length($3, 1), 0) = 0 ) AND ( $4 = 0 OR category = $4 ) ) SELECT rm.id, rm.name, rm.score, u.id AS author_id, u.name AS author_name, rm.downloads, rm.phone_names, rm.category, e.name AS effect_name, COUNT(*) OVER () AS results FROM ringtones_matched rm INNER JOIN effect e ON rm.effect_id = e.id INNER JOIN author u ON rm.author_id = u.id`

	if search != "" {
		switch sortBy {
		case "popular":
			query += ` ORDER BY rm.score DESC, rm.name`
		case "latest":
			query += ` ORDER BY rm.score DESC, rm.time_added DESC, rm.name`
		case "name (a-z)":
			query += ` ORDER BY rm.score DESC, rm.name`
		default:
			query += ` ORDER BY rm.score DESC, rm.name`
		}
	} else {
		switch sortBy {
		case "popular":
			query += ` ORDER BY rm.score DESC, rm.name`
		case "latest":
			query += ` ORDER BY rm.time_added DESC, rm.name`
		case "name (a-z)":
			query += ` ORDER BY rm.name, rm.score DESC`
		default:
			query += ` ORDER BY rm.score DESC, rm.name`
		}
	}
	query += ` LIMIT $5 OFFSET $6;`

	rows, err = DB.Query(query, search, pq.Array(phones), pq.Array(effects), category, resultsPerPage, (page-1)*resultsPerPage)
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

func GetRingtonesByAuthor(authorName string, page int) ([]RingtoneModel, int, error) {
	var ringtones []RingtoneModel
	var rows *sql.Rows
	var err error

	rows, err = DB.Query(`WITH ringtones_matched AS ( SELECT id, name, ARRAY( SELECT p.name FROM phone_and_ringtone par INNER JOIN phone p ON p.id = par.phone_id WHERE par.ringtone_id = ringtone.id ) as phone_names, effect_id, author_id, downloads, ( downloads::FLOAT - 2 * not_working::FLOAT ) / GREATEST(MAX(downloads) OVER (), 1) AS score FROM ringtone WHERE author_id = ( SELECT id FROM author WHERE name = $1 ) ) SELECT rm.id, rm.name, rm.score, u.id AS author_id, u.name AS author_name, rm.downloads, rm.phone_names, e.name AS effect_name, COUNT(*) OVER () AS results FROM ringtones_matched rm INNER JOIN effect e ON rm.effect_id = e.id INNER JOIN author u ON rm.author_id = u.id ORDER BY score DESC LIMIT $2 OFFSET $3;`, authorName, resultsPerPage, (page-1)*resultsPerPage)
	if err != nil {
		log.Println(1)
		return ringtones, 0, err
	}

	err = scan.Rows(&ringtones, rows)
	if err != nil {
		log.Println(2)
		return ringtones, 0, err
	}

	var numberOfPages int = 0
	if len(ringtones) != 0 {
		numberOfPages = int(math.Ceil(float64(ringtones[0].NumberOfResults) / float64(resultsPerPage)))
	}

	return ringtones, numberOfPages, nil
}

func CreateRingtone(name string, category int, phones []int, effect int, authorID int) (int, error) {
	var ringtoneID int
	err := DB.QueryRow(`INSERT INTO ringtone (name, category, effect_id, author_id) VALUES ($1, $2, $3, $4) RETURNING id;`, name, category, effect, authorID).Scan(&ringtoneID)
	if err != nil {
		return 0, err
	}
	_, err = DB.Exec(`INSERT INTO phone_and_ringtone (ringtone_id, phone_id) SELECT $1, UNNEST($2::int[])`, ringtoneID, pq.Array(phones))
	if err != nil {
		return 0, err
	}
	return ringtoneID, nil
}

func DeleteRingtone(ringtoneID int, authorID int) error {
	_, err := DB.Exec(`DELETE FROM ringtone WHERE id = $1 AND author_id = $2;`, ringtoneID, authorID)
	return err
}

func GetRingtone(ringtoneID int) (RingtoneModel, error) {
	var ringtone RingtoneModel
	rows, err := DB.Query(`SELECT r.id, r.name, ARRAY ( SELECT p.name FROM phone_and_ringtone par INNER JOIN phone p ON p.id = par.phone_id WHERE par.ringtone_id = r.id ) as phone_names, u.id AS author_id, u.name AS author_name, e.name AS effect_name, r.downloads FROM ringtone r INNER JOIN effect e ON r.effect_id = e.id INNER JOIN author u ON r.author_id = u.id WHERE r.id = $1;`, ringtoneID)
	if err != nil {
		log.Println(err.Error())
		return ringtone, err
	}
	err = scan.Row(&ringtone, rows)
	return ringtone, err
}

func RenameRingtone(ringtoneID int, name string, authorID int) error {
	_, err := DB.Exec(`UPDATE ringtone SET name = $1 WHERE id = $2 AND (author_id = $3 OR $3 = 1);`, name, ringtoneID, authorID)
	return err
}

func RingtoneIncreaseDownload(id int) error {
	_, err := DB.Exec(`UPDATE ringtone SET downloads = downloads + 1 WHERE id = $1;`, id)
	return err
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

func GetAuthor(id int) (AuthorModel, error) {
	var author AuthorModel
	rows, err := DB.Query(`SELECT * FROM author WHERE id = $1;`, id)
	if err != nil {
		return author, err
	}
	err = scan.Row(&author, rows)
	return author, err
}

func CreateAuthor(name string, email string) (int, error) {
	var authorID int

	email = strings.ToLower(email)

	err := DB.QueryRow(`WITH res AS (INSERT INTO author (name, email) VALUES ($1, $2) ON CONFLICT(email) DO NOTHING RETURNING id) SELECT id FROM res UNION ALL SELECT id FROM author WHERE email = $2 LIMIT 1;`, name, email).Scan(&authorID)
	if err != nil {
		return 0, err
	}
	return authorID, nil

}

func RenameAuthor(id int, newName string) (string, error) {
	var email string
	err := DB.QueryRow(`UPDATE author SET name = $1 WHERE id = $2 RETURNING email;`, newName, id).Scan(&email)
	return email, err
}
