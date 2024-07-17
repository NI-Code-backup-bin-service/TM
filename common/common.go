package common

import (
	"database/sql"
	"encoding/base64"
	"strings"
	"time"
)

var ByteOrderMark = []byte{239, 187, 191}

const (
	UploadChangeType = 2
	DeleteChangetype = 3
)

func GetChangeType() func(int) string {

	changeTypeMap := map[int]string{
		1: "Update",
		2: "Upload",
		3: "Delete",
		5: "Create",
	}

	return func(key int) string {
		return changeTypeMap[key]
	}
}

func ConvertBase64FileToBytes(base64EncodedFile string) ([]byte, error) {
	splitString := strings.Split(base64EncodedFile, ";base64,")
	base64EncodedFile = splitString[len(splitString)-1]
	rawBytes := make([]byte, base64.StdEncoding.DecodedLen(len(base64EncodedFile)))
	length, err := base64.StdEncoding.Decode(rawBytes, []byte(base64EncodedFile))
	if err != nil {
		return []byte{}, err
	}
	//now set the raw bytes to the actual decoded size
	rawBytes = rawBytes[:length]
	return rawBytes, nil
}

func CheckStringIsValid(value sql.NullString) string {
	if value.Valid {
		return value.String
	}

	return ""
}

func CheckBoolIsValid(value sql.NullBool) bool {
	if value.Valid {
		return value.Bool
	}

	return false
}

func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04")
}
