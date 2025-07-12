package main

import (
	"fmt"
	"glyphtones/database"
	"os"
	"time"

	"github.com/blockloop/scan/v2"
	"github.com/joho/godotenv"
	"github.com/teris-io/shortid"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	database.Init()

	var numOfRingtones int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM ringtone;").Scan(&numOfRingtones)
	if err != nil {
		panic(err)
	}

	var newIDs []string
	for range numOfRingtones {
		id, err := shortid.Generate()
		if err != nil {
			panic(err)
		}
		newIDs = append(newIDs, id)

		time.Sleep(10 * time.Millisecond)
	}

	rows, err := database.DB.Query("SELECT * FROM ringtone;")
	if err != nil {
		panic(err)
	}

	var ringtones []database.RingtoneModel
	err = scan.Rows(&ringtones, rows)
	if err != nil {
		panic(err)
	}

	for i := range numOfRingtones {
		_, err := database.DB.Exec("UPDATE ringtone SET display_id = $1 WHERE id = $2", newIDs[i], ringtones[i].ID)
		if err != nil {
			panic(err)
		}

		stat, err := os.Stat(fmt.Sprintf("sounds/%s.ogg", newIDs[i]))
		if err == nil {
			panic(stat.Name())
		}

		err = os.Rename(fmt.Sprintf("sounds/%d.ogg", ringtones[i].ID), fmt.Sprintf("sounds/%s.ogg", newIDs[i]))
		if err != nil {
			panic(err)
		}
	}

	fmt.Println(newIDs)
}
