package utils

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
)

var ErrInvalidRingtoneFile = errors.New("invalid ringtone file")

type PhoneSpec struct {
	ID              int
	NumberOfColumns int
	AlternateCols   int
}

type FileCheckResult struct {
	PhoneIDs  []int
	GlyphData []byte
}

func CheckFile(file *os.File, phones []PhoneSpec) (FileCheckResult, error) {
	var result FileCheckResult

	cmd := exec.Command("ffprobe", "-i", file.Name(), "-show_streams", "-select_streams", "a", "-v", "quiet", "-of", "json")

	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return result, fmt.Errorf("run ffprobe: %w", err)
	}

	var probe struct {
		Streams []struct {
			CodecName string         `json:"codec_name"`
			Tags      map[string]any `json:"tags"`
		} `json:"streams"`
	}
	if err := json.Unmarshal(out.Bytes(), &probe); err != nil {
		return result, fmt.Errorf("decode ffprobe output: %w", err)
	}

	if len(probe.Streams) == 0 {
		return result, ErrInvalidRingtoneFile
	}
	stream := probe.Streams[0]
	if stream.CodecName != "opus" {
		return result, ErrInvalidRingtoneFile
	}

	author, ok := stream.Tags["AUTHOR"].(string)
	if !ok || author == "" {
		return result, ErrInvalidRingtoneFile
	}

	decoded, err := base64.StdEncoding.DecodeString(author)
	if err != nil {
		decoded, err = base64.RawStdEncoding.DecodeString(author)
		if err != nil {
			return result, ErrInvalidRingtoneFile
		}
	}

	reader, err := zlib.NewReader(bytes.NewReader(decoded))
	if err != nil {
		return result, ErrInvalidRingtoneFile
	}
	defer reader.Close()

	var decompressed bytes.Buffer
	if _, err := io.Copy(&decompressed, reader); err != nil {
		return result, ErrInvalidRingtoneFile
	}

	csvReader := csv.NewReader(&decompressed)
	csvReader.TrimLeadingSpace = true
	record, err := csvReader.Read()
	if err != nil {
		return result, ErrInvalidRingtoneFile
	}

	columns := len(record)
	for i := len(record) - 1; i >= 0; i-- {
		if record[i] == "" {
			columns--
			continue
		}
		break
	}

	result.GlyphData = decoded
	for _, phone := range phones {
		if phone.NumberOfColumns == columns || phone.AlternateCols == columns {
			result.PhoneIDs = append(result.PhoneIDs, phone.ID)
		}
	}
	if len(result.PhoneIDs) == 0 {
		return FileCheckResult{}, ErrInvalidRingtoneFile
	}

	if _, err := file.Seek(0, 0); err != nil {
		return FileCheckResult{}, fmt.Errorf("rewind ringtone file: %w", err)
	}

	return result, nil
}

func CreateRingtoneFile(src *os.File, dstDir string, name string) error {
	dst, err := os.Create(fmt.Sprintf("%s/%s.ogg", dstDir, name))
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

func CreateTemporaryFile(tmpDir string, src io.Reader) (*os.File, error) {
	dst, err := os.CreateTemp(tmpDir, "upload")
	if err != nil {
		return nil, fmt.Errorf("create temp file: %w", err)
	}

	if _, err := io.Copy(dst, src); err != nil {
		dst.Close()
		return nil, fmt.Errorf("copy upload to temp file: %w", err)
	}

	if _, err := dst.Seek(0, 0); err != nil {
		dst.Close()
		return nil, fmt.Errorf("rewind temp file: %w", err)
	}

	return dst, nil
}

func DeleteFile(name string) error {
	if err := os.Remove(name); err != nil {
		return err
	}
	return nil
}

func LogDeleteFileError(path string, err error) {
	if err != nil {
		slog.Error("failed to delete file", "path", path, "error", err)
	}
}
