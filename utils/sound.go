package utils

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"glyphtones/database"
	"io"
	"log"
	"os"
	"os/exec"
)

var RingtonesDir string = "./sounds"
var TemporaryDir string = "./tmp"

func CheckFile(file *os.File, phones []database.PhoneModel) ([]int, string, bool) {
	var phonesResult []int

	cmd := exec.Command("ffprobe", "-i", file.Name(), "-show_streams", "-select_streams", "a", "-v", "quiet", "-of", "json")

	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		fmt.Println("Error running FFprobe:", err)
		return phonesResult, "", false
	}

	var result map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		fmt.Println("Error running Json Unmarshal:", err)
		return phonesResult, "", false
	}

	streams, ok := result["streams"].([]interface{})
	if !ok || len(streams) == 0 {
		return phonesResult, "", false
	}
	firstStream, ok := streams[0].(map[string]interface{})
	if !ok {
		return phonesResult, "", false
	}
	if firstStream["codec_name"] != "opus" {
		return phonesResult, "", false
	}
	tags, ok := firstStream["tags"].(map[string]interface{})
	if !ok {
		return phonesResult, "", false
	}

	author, ok := tags["AUTHOR"].(string)
	if !ok {
		return phonesResult, "", false
	}

	// im not sure what the difference between StdEncoding and RawStdEncoding is, but sometimes the file doesn't decode with the normal one so I have to try raw too...
	decoded, err := base64.StdEncoding.DecodeString(author) // decode from base64 to bytes
	if err != nil {
		decoded, err = base64.RawStdEncoding.DecodeString(author) // decode from base64 to bytes
		if err != nil {
			return phonesResult, "", false
		}
	}
	reader, err := zlib.NewReader(bytes.NewReader(decoded)) // decode zlib compression
	if err != nil {
		return phonesResult, "", false
	}
	defer reader.Close()

	var decompressed bytes.Buffer
	_, err = io.Copy(&decompressed, reader) // copy the result to buffer
	if err != nil {
		return phonesResult, "", false
	}

	csvReader := csv.NewReader(&decompressed)
	csvReader.TrimLeadingSpace = true
	record, err := csvReader.Read()
	if err != nil {
		return phonesResult, "", false
	}
	columns := len(record)
	for i := len(record) - 1; i >= 0; i-- {
		if record[i] == "" {
			columns -= 1
		} else {
			break
		}
	}

	for _, v := range phones {
		if v.NumberOfColumns == columns || v.NumberOfColumns2 == columns {
			phonesResult = append(phonesResult, v.ID)
		}
	}

	file.Seek(0, 0)

	// finally if everything is ok, return true
	return phonesResult, author, true
}

func CreateRingtoneFile(src *os.File, ringtoneID int) error {
	dst, err := os.Create(fmt.Sprintf("%s/%d.ogg", RingtonesDir, ringtoneID))
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

func CreateTemporaryFile(src io.Reader) (*os.File, error) {
	dst, err := os.CreateTemp(TemporaryDir, "upload")
	if err != nil {
		log.Println(3)
		return dst, err
	}

	_, err = io.Copy(dst, src)
	if err != nil {
		log.Println(4)
		return dst, err
	}

	dst.Seek(0, 0)

	return dst, nil
}

func DeleteFile(name string) {
	if err := os.Remove(name); err != nil {
		log.Println(err)
	}
}
