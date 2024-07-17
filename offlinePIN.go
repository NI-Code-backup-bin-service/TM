package main

import (
	pb "bitbucket.org/network-international/nextgen-libs/nextgen-tg-protobuf/Transaction"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"nextgen-tms-website/entities"
	"strconv"
	"strings"
	"time"
)

const offlinePINSecret = "0aeadf2eb87fbeccdb95bc0bbc78d15d7d42011a9b3b175aca1a397b2075403e"
const resetPINSecret = "b87d74f1032bf5168505376acf0129da40e39cbbff1c476b4aef67b8c010954c"

type OfflinePINData struct {
	SerialNumber      string `json:"serial"`
	IMEI              string `json:"imei"`
	OfflineExpiryDays string `json:"expiry1"`
	ResetExpiryDays   string `json:"expiry2"`
	OfflinePINSN      string `json:"offline-sn"`
	OfflinePINIMEI    string `json:"offline-imei"`
	ResetPINSN        string `json:"reset-sn"`
	ResetPINIMEI      string `json:"reset-imei"`
	OfflineExpiryDate string `json:"expiry-date-engineering"`
	ResetExpiryDate   string `json:"expiry-date-reset"`
}

func offlinePIN(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	renderHeader(w, r, tmsUser)
	renderTemplate(w, r, "offlinePIN", nil, tmsUser)
}

func generateOfflinePIN(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	logging.Information("User " + tmsUser.Username + " attempting to generate Offline PIN(s)")

	if err := r.ParseForm(); err != nil {
		logging.Error(err.Error())
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	dataJSON := r.Form.Get("pedData")
	var data []OfflinePINData
	if err := json.Unmarshal([]byte(dataJSON), &data); err != nil {
		logging.Error(err.Error())
		http.Error(w, "Data not in correct format", http.StatusBadRequest)
		return
	}

	generatedType := "Unknown"
	generateMode := r.Form.Get("mode")
	switch generateMode {
	case "both":
		generatedType = "Engineering and Factory Reset"
	case "reset":
		generatedType = "Factory Reset"
	case "offline":
		generatedType = "Engineering"
	}
	logging.Information("User " + tmsUser.Username + " generating " + generatedType + " PIN(s)")

	for i := range data {
		const hashFormat = "20060102"
		const outputFormat = "02/01/2006"

		days, err := strconv.Atoi(data[i].OfflineExpiryDays)
		if err != nil {
			logging.Error("Engineering ExpiryDays for column " + strconv.Itoa(i) + " invalid:" + err.Error())
		} else if generateMode == "offline" || generateMode == "both" {
			// 1 day is only valid for current day, so we will need to reduce days to increment by 1
			days = days - 1

			expiryDate := time.Now()
			expiryDate = expiryDate.AddDate(0, 0, days)

			if data[i].SerialNumber != "" {
				data[i].OfflinePINSN = generatePIN(offlinePINSecret, data[i].SerialNumber, expiryDate.Format(hashFormat))
			}

			if data[i].IMEI != "" {
				data[i].OfflinePINIMEI = generatePIN(offlinePINSecret, data[i].IMEI, expiryDate.Format(hashFormat))
			}

			data[i].OfflineExpiryDate = expiryDate.Format(outputFormat)
		}

		days, err = strconv.Atoi(data[i].ResetExpiryDays)
		if err != nil {
			logging.Error("Reset ExpiryDays for column " + strconv.Itoa(i) + " invalid:" + err.Error())
		} else if generateMode == "reset" || generateMode == "both" {
			// 1 day is only valid for current day, so we will need to reduce days to increment by 1
			days = days - 1

			expiryDate := time.Now()
			expiryDate = expiryDate.AddDate(0, 0, days)

			if data[i].SerialNumber != "" {
				data[i].ResetPINSN = generatePIN(resetPINSecret, data[i].SerialNumber, expiryDate.Format(hashFormat))
			}

			if data[i].IMEI != "" {
				data[i].ResetPINIMEI = generatePIN(resetPINSecret, data[i].IMEI, expiryDate.Format(hashFormat))
			}

			data[i].ResetExpiryDate = expiryDate.Format(outputFormat)
		}
	}

	logging.Information("User " + tmsUser.Username + " generated PIN(s) for " + strconv.Itoa(len(data)) + " TIDs")

	// Send a logging message for Terminal Monitoring of Offline PIN generation
	logData := new(pb.LoggingRequest)
	logMessage := pb.LogMessage{
		Level:   4,
		Source:  "TMSWebSite",
		LogText: tmsUser.Username + ": Generated " + generatedType + " PIN(s) for " + strconv.Itoa(len(data)) + " TIDs",
		Extra:   "MONITOR",
	}

	logData.Messages = append(logData.Messages, &logMessage)
	if err := logging.PassThroughLog(logData); err != nil {
		logging.Debug(err)
	}

	returnJSON, err := json.Marshal(data)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Failed to generate JSON return data", http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(returnJSON); err != nil {
		logging.Error(err.Error())
		http.Error(w, "Failed to write return data", http.StatusInternalServerError)
		return
	}
}

func offlinePINImportCSV(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	if err := r.ParseMultipartForm(10000); err != nil {
		logging.Error(err.Error())
		http.Error(w, uploadFileError, http.StatusBadRequest)
		return
	}

	// Validate that a file has been attached
	if len(r.MultipartForm.File) < 1 {
		err := "No file provided"
		logging.Error(err)
		http.Error(w, err, http.StatusBadRequest)
		return
	}

	// Extract the file from the request and obtain the name
	logging.Debug("Attempting to extract file from http.Request")
	file, handler, err := r.FormFile("file")
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, uploadFileError, http.StatusBadRequest)
		return
	}
	defer file.Close()

	filename := handler.Filename
	nameParts := strings.Split(filename, ".")
	if nameParts[len(nameParts)-1] != "csv" {
		err := "Uploaded file is not a CSV"
		logging.Error(err)
		http.Error(w, err, http.StatusBadRequest)
		return
	}

	logging.Debug(fmt.Sprintf("File: %v has been uploaded", filename))

	if _, err = file.Seek(0, 0); err != nil {
		logging.Error(err.Error())
		http.Error(w, html.EscapeString(FailedFileRead+filename), http.StatusInternalServerError)
		return
	}

	logging.Debug("Parsing CSV data")
	// Parse the entries from the csv along with the column headers
	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, html.EscapeString(FailedFileRead+filename), http.StatusInternalServerError)
		return
	} else if len(records) == 0 || len(records[1:]) == 0 {
		logging.Error("No records found in uploaded CSV file")
		http.Error(w, NoColumnsFound, http.StatusInternalServerError)
		return
	}

	var data []OfflinePINData
	// Ranging over all records except columns definitions row
	for i, record := range records[1:] {
		if len(record) < 4 {
			logging.Error("Row " + strconv.Itoa(i) + " does not contain enough columns")
			continue
		}

		row := OfflinePINData{
			SerialNumber:      record[0],
			IMEI:              record[1],
			OfflineExpiryDays: record[2],
			ResetExpiryDays:   record[3],
		}
		data = append(data, row)
	}

	returnJSON, err := json.Marshal(data)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Failed to generate JSON return data", http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(returnJSON); err != nil {
		logging.Error(err.Error())
		http.Error(w, "Failed to write return data", http.StatusInternalServerError)
		return
	}
}

func generatePIN(sharedStatic string, serialNumber string, validityDate string) string {
	const OfflinePinLength = 1000000

	toHashStr := sharedStatic + serialNumber + validityDate
	toHashBytes := make([]byte, hex.EncodedLen(len(toHashStr)))
	hex.Encode(toHashBytes, []byte(toHashStr))

	hasher := sha256.New()
	hasher.Write(toHashBytes)
	hash := hasher.Sum(nil)

	offset := hash[len(hash)-1] & 0xf
	part := (uint32(hash[offset])&0x7f)<<24 | (uint32(hash[offset+1])&0xff)<<16 | (uint32(hash[offset+2])&0xff)<<8 | (uint32(hash[offset+3]) & 0xff)
	final := int(part) % OfflinePinLength
	return fmt.Sprintf("%06d", final)
}
