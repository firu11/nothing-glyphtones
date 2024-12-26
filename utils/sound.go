package utils

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/dhowden/tag"
)

var ringtonesDir string = "./sounds"

func CheckFile(file *os.File) (bool, error) {
	m, err := tag.ReadFrom(file)
	if err != nil {
		return false, err
	}

	if m.Format() != tag.VORBIS {
		return false, nil
	}
	if m.FileType() != tag.OGG {
		return false, nil
	}

	if m.Raw()["author"] == nil {
		return false, nil
	}
	author := m.Raw()["author"].(string)

	// im not sure what the difference between StdEncoding and RawStdEncoding is, but sometimes the file doesn't decode with the normal one so I have to try raw too...
	decoded, err := base64.StdEncoding.DecodeString(author) // decode from base64 to bytes
	if err != nil {
		decoded, err = base64.RawStdEncoding.DecodeString(author) // decode from base64 to bytes
		if err != nil {
			return false, nil
		}
	}
	reader, err := zlib.NewReader(bytes.NewReader(decoded)) // decode zlib compression
	if err != nil {
		return false, nil
	}
	defer reader.Close()

	var decompressed bytes.Buffer
	_, err = io.Copy(&decompressed, reader) // copy the result to buffer
	if err != nil {
		return false, nil
	}

	// i can read the csv but no need now
	/* csvReader := csv.NewReader(&decompressed)
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(record)
	} */

	file.Seek(0, 0)

	// finally if everything is ok, return true
	return true, nil
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
