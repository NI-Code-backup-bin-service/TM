package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"html"
	"net/http"
	"nextgen-tms-website/common"
	"nextgen-tms-website/dal"
	"nextgen-tms-website/entities"
	"strings"
	"time"
)

const (
	TFTAG                = "Terminal Flagging : "
	TerminalFlaggingType = "TerminalFlagging"
)

func terminalFlaggingHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	logging.Debug(TFTAG, "terminalFlaggingHandler()")
	renderHeader(w, r, tmsUser)
	renderTemplate(w, r, "terminalFlagging", nil, tmsUser)
}

func terminalFlaggingUploadHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	logging.Debug(TFTAG, "Terminal Flagging upload initiated")

	err := r.ParseMultipartForm(10000)
	if err != nil {
		logging.Error(TFTAG, err.Error())
		http.Error(w, uploadFileError, http.StatusBadRequest)
		return
	}

	// Validate that a file has been attached
	if len(r.MultipartForm.File) < 1 {
		logging.Warning(TFTAG, "Site upload initiated without file being present")
		http.Error(w, "Please choose a file", http.StatusBadRequest)
		return
	}

	// Extract the file from the request and obtain the name
	logging.Debug(TFTAG, "Attempting to extract file from http.Request")

	file, handler, err := r.FormFile("file")
	if err != nil {
		logging.Error(TFTAG, err.Error())
		http.Error(w, uploadFileError, http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileName := handler.Filename

	logging.Debug(TFTAG, fmt.Sprintf("File: %v has been uploaded", fileName))

	// Validate that the filetype is csv
	// 512 bytes only because DetectContentType (used in ValidateFileType) only reads up to 512 bytes
	buff := make([]byte, 512)
	if _, err = file.Seek(0, 0); err != nil {
		logging.Error(TFTAG, err.Error())
		http.Error(w, html.EscapeString(FailedFileRead+fileName), http.StatusInternalServerError)
		return
	}

	if _, err = file.Read(buff); err != nil {
		logging.Error(TFTAG, err.Error())
		http.Error(w, html.EscapeString(FailedFileRead+fileName), http.StatusInternalServerError)
		return
	}

	// Validates the file type, based on File type
	isCSVFile := strings.HasSuffix(strings.ToLower(fileName), ".csv")
	if !isCSVFile {
		logging.Warning(TFTAG, "Incorrect filetype uploaded")
		http.Error(w, IncorrectFileTypeCSV, http.StatusInternalServerError)
		return
	}

	logging.Debug(TFTAG, fmt.Sprintf("File %v has passed type validation", fileName))

	// we know that the file is of csv format

	logging.Debug(TFTAG, "Resetting file read offset")
	// Need to reset the offset after checking for filetype
	if _, err = file.Seek(0, 0); err != nil {
		logging.Error(TFTAG, err.Error())
		http.Error(w, html.EscapeString(FailedFileRead+fileName), http.StatusInternalServerError)
		return
	}

	logging.Debug(TFTAG, "Parsing CSV data")
	// Parse the entries from the csv along with the column headers
	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		logging.Error(TFTAG, err.Error())
		http.Error(w, html.EscapeString(FailedFileRead+fileName), http.StatusInternalServerError)
		return
	} else if len(records) == 0 || len(records[1:]) < 1 {
		logging.Error(TFTAG, "No records found in uploaded CSV file")
		http.Error(w, NoColumnsFound, http.StatusInternalServerError)
		return
	}

	if len(records[0]) < 3 {
		logging.Error(TFTAG, "Invalid File Header", records[0])
		http.Error(w, "Invalid File Header", http.StatusInternalServerError)
		return
	}

	//some special characters are getting appended in the first row; inorder to remove that bytes.Trim is used
	header := string(bytes.TrimPrefix([]byte(strings.ToLower(strings.TrimSpace(records[0][0]))), common.ByteOrderMark))

	if header != "tid" || strings.ToLower(strings.TrimSpace(records[0][1])) != "serial number" || strings.ToLower(strings.TrimSpace(records[0][2])) != "apk version" {
		logging.Error(TFTAG, "Invalid File Header", records[0])
		http.Error(w, "Invalid File Header", http.StatusInternalServerError)
		return
	}

	tpApkExists := false
	for _, r2 := range records[1:] {
		if len(r2) > 3 && strings.TrimSpace(r2[3]) != "" {
			if strings.ToLower(strings.TrimSpace(records[0][3])) != "tp apk" {
				logging.Error(TFTAG, "Invalid File Header", records[0][3])
				http.Error(w, "Invalid File Header", http.StatusInternalServerError)
				return
			}
			tpApkExists = true
		}
	}

	if !tpApkExists {
		if strings.ToLower(strings.TrimSpace(records[0][3])) != "tp apk" && strings.TrimSpace(records[0][3]) != "" {
			logging.Error(TFTAG, "Invalid File Header", records[0][3])
			http.Error(w, "Invalid File Header", http.StatusInternalServerError)
			return
		}
	}

	buf := &bytes.Buffer{}
	writer := csv.NewWriter(buf)
	err = writer.WriteAll(records[:])
	if err != nil {
		logging.Error(TFTAG, "Terminal Flagging File Upload Failed while writing data to buf : ", err.Error())
		http.Error(w, html.EscapeString(writeDataError+err.Error()), http.StatusInternalServerError)
		return
	}

	fileName = time.Now().Format("20060102150405_") + fileName
	logging.Debug(TFTAG, "Renamed file to : "+fileName)

	if err := sendFileToFileServer(buf.Bytes(), fileName, TerminalFlaggingType); err != nil {
		logging.Error(TFTAG, "Terminal Flagging File Upload Failed : ", err.Error())
		http.Error(w, uploadFileError, http.StatusInternalServerError)
		return
	}

	err = dal.InsertBulkApproval(fileName, TerminalFlaggingType, tmsUser.Username, common.UploadChangeType, 0)
	if err != nil {
		http.Error(w, DatabaseTxnError, http.StatusInternalServerError)
		return
	}

	type TerminalFlagUpload struct {
		UploadStatus bool
	}
	var FlagUpload TerminalFlagUpload
	FlagUpload.UploadStatus = true
	logging.Debug(TFTAG, "Terminal Flagging Upload Successful")
	renderPartialTemplate(w, r, "terminalFlaggingUploadResults", FlagUpload, tmsUser)
}

func validateCsvFileType(buff []byte, desiredFiletype string) bool {
	filetype := http.DetectContentType(buff)
	return strings.Contains(filetype, desiredFiletype)
}
