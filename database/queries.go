package database

import (
	"database/sql"
	"log"
	"sort"

	"github.com/blockloop/scan/v2"
)

func GetRingtones(search string, phone int, effect int) ([]RingtoneModel, error) {
	var ringtones []RingtoneModel
	var rows *sql.Rows
	var err error

	log.Println(phone, effect)
	rows, err = DB.Query(`SELECT * FROM ringtone WHERE name LIKE '%' || $1 || '%';`, search) //, phone, effect AND ($2 == 0 OR phone == $2) AND ($3 == 0 OR effect == $3);
	if err != nil {
		log.Println(err.Error())
		return ringtones, err
	}

	err = scan.Rows(&ringtones, rows)
	if err != nil {
		log.Println(2)
		return ringtones, err
	}

	return ringtones, nil
}

func GetPhones() ([]PhoneModel, error) {
	var phones []PhoneModel
	var rows *sql.Rows
	var err error

	rows, err = DB.Query(`SELECT * FROM phone;`)
	if err != nil {
		log.Println(err.Error())
		return phones, err
	}

	err = scan.Rows(&phones, rows)
	if err != nil {
		log.Println(2)
		return phones, err
	}

	sort.Slice(phones, func(i, j int) bool {
		return phones[i].Name < phones[j].Name
	})

	return phones, nil
}

func GetEffects() ([]EffectModel, error) {
	var phones []EffectModel
	var rows *sql.Rows
	var err error

	rows, err = DB.Query(`SELECT * FROM effect;`)
	if err != nil {
		log.Println(err.Error())
		return phones, err
	}

	err = scan.Rows(&phones, rows)
	if err != nil {
		log.Println(2)
		return phones, err
	}

	sort.Slice(phones, func(i, j int) bool {
		return phones[i].Name < phones[j].Name
	})

	return phones, nil
}
