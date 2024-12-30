package utils

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"glyphtones/database"
	"io"
	"log"
	"os"

	"github.com/dhowden/tag"
)

var ringtonesDir string = "./sounds"

func CheckFile(file *os.File, phones []database.PhoneModel) ([]int, bool) {
	var phonesResult []int

	m, err := tag.ReadFrom(file)
	if err != nil {
		return phonesResult, false
	}

	if m.Format() != tag.VORBIS {
		return phonesResult, false
	}
	if m.FileType() != tag.OGG {
		return phonesResult, false
	}

	if m.Raw()["author"] == nil {
		return phonesResult, false
	}
	author := m.Raw()["author"].(string)

	// im not sure what the difference between StdEncoding and RawStdEncoding is, but sometimes the file doesn't decode with the normal one so I have to try raw too...
	decoded, err := base64.StdEncoding.DecodeString(author) // decode from base64 to bytes
	if err != nil {
		decoded, err = base64.RawStdEncoding.DecodeString(author) // decode from base64 to bytes
		if err != nil {
			return phonesResult, false
		}
	}
	reader, err := zlib.NewReader(bytes.NewReader(decoded)) // decode zlib compression
	if err != nil {
		return phonesResult, false
	}
	defer reader.Close()

	var decompressed bytes.Buffer
	_, err = io.Copy(&decompressed, reader) // copy the result to buffer
	if err != nil {
		return phonesResult, false
	}

	csvReader := csv.NewReader(&decompressed)
	csvReader.TrimLeadingSpace = true
	record, err := csvReader.Read()
	if err != nil {
		return phonesResult, false
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
	return phonesResult, true
}

func CreateRingtoneFile(src *os.File, ringtoneID int) error {
	dst, err := os.Create(fmt.Sprintf("%s/%d.ogg", ringtonesDir, ringtoneID))
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

func CreateTemporaryFile(src io.Reader) (*os.File, error) {
	dst, err := os.CreateTemp("./tmp", "upload")
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

func DeleteTemporaryFile(name string) {
	log.Println(os.Remove(name))
}
