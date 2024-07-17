package main

import (
	rpcHelp "bitbucket.org/network-international/nextgen-libs/nextgen-helpers/rpcHelper"
	txn "bitbucket.org/network-international/nextgen-libs/nextgen-tg-protobuf/Transaction"
	"bytes"
	"errors"
	"fmt"
	"html"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"nextgen-tms-website/common"
	"nextgen-tms-website/config"
	"nextgen-tms-website/dal"
	"nextgen-tms-website/entities"
	"nextgen-tms-website/fileServer"
	"path/filepath"
	"strings"
)

const (
	FailedFileRead     = "Failed to read file : "
	FailedRetrieveFile = "Error retrieving file : "
	EncodingError      = "an encoding error occurred"

	// SoftUi File Upload types
	MenuConfiguration   = "menuConfig"
	SoftUIConfiguration = "softUIConfig"
	CustomerReceipt     = "customerReceipt"
	MerchantReceipt     = "merchantReceipt"
	IconImage           = "icon"
	FileManagementType  = "FileManagement"
	FileManagementTAG   = "File Management"
	LogoManagementType  = "LogoManagement"
	LogoManagementTAG   = "Logo Management"

	merchantReceipt         = "merchantReceipt"
	merchantReceiptFilePath = "SoftUI/WebViewPages/application/common/receipts/MerchantReceipts"
	customerReceipt         = "customerReceipt"
	customerReceiptFilePath = "SoftUI/WebViewPages/application/common/receipts/CustomerReceipts"
)

var (
	allowedTextFileTypes = []string{
		"txt",
		"dat",
		"db",
		"rmu",
		"json",
	}

	allowedImageFileTypes = []string{
		"image/jpeg",
		"image/jpg",
		"image/gif",
		"image/png",
	}

	allowedMenuFileTypes = []string{
		"json",
	}
	allowedReceiptConfigFileTypes = []string{
		"json",
	}
	allowedReceiptFileTypes = []string{
		"rmu",
		"txt",
	}
	allowedIconFileTypes = []string{
		"png",
		"jpg",
	}
	allowedMnoFileTypes = []string{
		"json",
		"png",
		"jpg",
		"gif",
	}
)

func getFileListHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	if err := r.ParseForm(); err != nil {
		logging.Error(err.Error())
		http.Error(w, retrieveFileListError, http.StatusInternalServerError)
		return
	}
	fileType := r.Form.Get("type")
	var files []string
	var err error
	var allowedFormats []string
	switch fileType {
	case "customerreceipt":
		files, err = fileServer.NewFsReader(config.FileserverURL).GetAllCustomerReceiptFiles()
		allowedFormats = allowedReceiptFileTypes
	case "merchantreceipt":
		files, err = fileServer.NewFsReader(config.FileserverURL).GetAllMerchantReceiptFiles()
		allowedFormats = allowedReceiptFileTypes
	case "menu":
		files, err = fileServer.NewFsReader(config.FileserverURL).GetAllMenuFiles()
		allowedFormats = allowedMenuFileTypes
	case "softUIConfig":
		files, err = fileServer.NewFsReader(config.FileserverURL).GetAllSoftUIConfigFiles()
		allowedFormats = allowedMenuFileTypes
	case "receiptConfig":
		files, err = fileServer.NewFsReader(config.FileserverURL).GetAllReceiptConfigFiles()
		allowedFormats = allowedReceiptConfigFileTypes
	case "mnoLogo":
		files, err = fileServer.NewFsReader(config.FileserverURL).GetAllMnoLogoFiles()
		allowedFormats = allowedMnoFileTypes
	default:
		switch fileType {
		case "image":
			allowedFormats = allowedImageFileTypes
		case "text":
			allowedFormats = allowedTextFileTypes
		default:
			allowedFormats = nil
		}
		files, err = fileServer.NewFsReader(config.FileserverURL).GetAllFiles()
	}
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, retrieveFileListError, http.StatusInternalServerError)
		return
	}

	var filteredFiles []FileListEntry
	for _, file := range files {
		if isAllowedFormat(file, allowedFormats) {
			filteredFiles = append(filteredFiles, FileListEntry{Name: file})
		}
	}

	model := ChooseFileModel{
		ButtonText: r.Form.Get("ButtonText"),
		Files:      filteredFiles,
		FileType:   fileType}

	renderPartialTemplate(w, r, "chooseFileDialog", model, tmsUser)
}

func fileUploadHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	renderHeader(w, r, tmsUser)
	renderTemplate(w, r, "fileUpload", nil, tmsUser)
}

func logoUploadHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	logging.Debug("logoUploadHandler()")
	renderHeader(w, r, tmsUser)
	renderTemplate(w, r, "logoUpload", nil, tmsUser)
}

func uploadSoftUiFileHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {

	// Parse multi part form
	logging.Debug(TAG, "SoftUI file upload initiated, attempting to parse form")
	err := r.ParseMultipartForm(10000)
	if err != nil {
		logging.Error(TAG, "An error has been thrown attempting to parse SoftUI file upload: "+err.Error())
		http.Error(w, uploadFileError, http.StatusBadRequest)
		return
	}

	// Extract the upload type. This will denote the destination for the file
	logging.Debug(TAG, "SoftUI form parse successful. Attempting to obtain SoftUI file upload type from http.Request")
	uploadType := r.PostFormValue("fileType")
	// r.PostFormValue does not return an error, instead if something goes wrong it returns an empty string
	if uploadType == "" {
		logging.Error(TAG, "An error has been thrown attempting to parse SoftUI file upload type: "+err.Error())
		http.Error(w, softUiFileTypeError, http.StatusBadRequest)
		return
	}
	logging.Debug(TAG, "Successfully obtained SoftUI upload filetype, continuing with upload of filetype: "+uploadType)

	// Validate that a file has been attached
	if len(r.MultipartForm.File) < 1 {
		logging.Warning(TAG, "SoftUI file upload initiated without file being present")
		http.Error(w, softUiFileMissingError, http.StatusBadRequest)
		return
	}

	var fileName string
	var buf bytes.Buffer
	var outboundReq bytes.Buffer
	fileWriter := multipart.NewWriter(&outboundReq)

	for _, header := range r.MultipartForm.File {
		file, err := header[0].Open()
		fileName = header[0].Filename
		if err != nil {
			logging.Error(err.Error())
			http.Error(w, uploadFileError, http.StatusInternalServerError)
			return
		} else {
			// Check the filetype is one of those permitted for the upload type
			if !softUiFileValidator(fileName, uploadType) {
				logging.Error("Invalid filetype for chosen upload")
				http.Error(w, invalidFileError, http.StatusInternalServerError)
				return
			}

			if _, err := io.Copy(&buf, file); err != nil {
				logging.Error(err.Error())
				http.Error(w, uploadFileError, http.StatusInternalServerError)
				return
			}

			buff := make([]byte, 512)
			if _, err = file.Seek(0, 0); err != nil {
				logging.Error(err.Error())
				http.Error(w, html.EscapeString(FailedFileRead+fileName), http.StatusInternalServerError)
				file.Close()
				return
			}

			if _, err = file.Read(buff); err != nil {
				logging.Error(err.Error())
				http.Error(w, html.EscapeString(FailedFileRead+fileName), http.StatusInternalServerError)
				file.Close()
				return
			}

			part, err := fileWriter.CreateFormFile("file."+header[0].Filename, header[0].Filename)
			if err != nil {
				logging.Error(err.Error())
				http.Error(w, uploadFileError, http.StatusInternalServerError)
				return
			}
			_, err = part.Write(buf.Bytes())
			if err != nil {
				logging.Error(err.Error())
				http.Error(w, uploadFileError, http.StatusInternalServerError)
				return
			}
		}
	}

	fileWriter.Close()

	filePath := ""

	switch uploadType {
	case merchantReceipt:
		filePath = merchantReceiptFilePath
	case customerReceipt:
		filePath = customerReceiptFilePath
	}

	if filePath != "" {
		grpcClient, clientFound := GRPCclients["TMFileUpload"]
		if !clientFound {
			logging.Error("TMFileUpload : Error obtaining GRPC Client")
			http.Error(w, uploadFileError, http.StatusInternalServerError)
			return
		} else {
			logging.Debug("TMFileUpload : client found, connection state: " + grpcClient.GetConnection().GetState().String())

			request := &txn.FileUploadRequest{
				FileFragment: buf.Bytes(),
				FileName:     fileName,
				FilePath:     filePath,
			}

			grpcReply := new(txn.FileUploadResponse)

			err = rpcHelp.ExecuteGRPC(grpcClient, request, grpcReply, logging)
			if err != nil {
				logging.Error(err.Error())
				http.Error(w, uploadFileError, http.StatusInternalServerError)
				return
			}
		}
	}

	req, err := http.NewRequest(http.MethodPost, config.FileserverURL+"/uploadSoftUiFile?type="+uploadType, &outboundReq)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, uploadFileError, http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", fileWriter.FormDataContentType())
	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, uploadFileError, http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		respBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logging.Error(err.Error())
			respBytes = append(respBytes, []byte(err.Error())...)
		}
		http.Error(w, uploadFileError, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/search", http.StatusSeeOther)
	return
}

func deleteFileHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	r.ParseForm()
	fileName := r.Form.Get("fileName")

	var fileCount int
	var TAG, Type string
	db, err := dal.GetDB()
	if err != nil {
		logging.Error("An error occured while connecting to database :" + err.Error())
		http.Error(w, "Unable to connect to database", http.StatusInternalServerError)
		return
	}
	directory := r.Form.Get("directory")
	values := make(map[string][]string, 0)
	values["FileName"] = []string{fileName}
	values["Directory"] = []string{directory}
	if directory == "mnoLogo" {
		TAG = LogoManagementTAG
		Type = LogoManagementType
	} else {
		TAG = FileManagementTAG
		Type = FileManagementType
	}

	err = db.QueryRow("select count(*) from bulk_approvals where filename=? AND approved= 0", fileName).Scan(&fileCount)
	if err != nil {
		logging.Error(TAG + err.Error())
		http.Error(w, "An error occured while fetching the file Counts", http.StatusInternalServerError)
		return
	}

	if fileCount > 0 {
		logging.Debug(TAG + "file already exists:" + fileName + "updating its details")
		_, err = db.Exec(" UPDATE bulk_approvals set created_by=created_by,created_at=NOW() where filename = ? AND approved= 0", fileName)
		if err != nil {
			logging.Error(TAG, "Failed to update file details:"+err.Error())
			http.Error(w, "unable to update file details:"+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		err := dal.InsertBulkApproval(fileName, Type, tmsUser.Username, common.DeleteChangetype, 0)
		if err != nil {
			logging.Error(TAG, "Failed to insert into bulkApproval"+err.Error())
			http.Error(w, DatabaseTxnError, http.StatusInternalServerError)
			return
		}
	}

}

func uploadFileHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	err := r.ParseMultipartForm(10000)
	if err != nil {
		logging.Error(err.Error())
	}

	var buf bytes.Buffer
	var outboundReq bytes.Buffer
	mw := multipart.NewWriter(&outboundReq)

	if len(r.MultipartForm.File) < 1 {
		http.Error(w, "Please choose a file to upload", http.StatusBadRequest)
		return
	}

	for _, header := range r.MultipartForm.File {
		file, err := header[0].Open()
		fileName := header[0].Filename
		if err != nil {
			logging.Error(err.Error())
			http.Error(w, uploadFileError, http.StatusInternalServerError)
			return
		} else {
			if _, err := io.Copy(&buf, file); err != nil {
				logging.Error(err.Error())
				http.Error(w, uploadFileError, http.StatusInternalServerError)
				return
			}

			// 512 bytes only because DetectContentType (used in ValidateFileType) only reads up to 512 bytes
			buff := make([]byte, 512)
			if _, err = file.Seek(0, 0); err != nil {
				logging.Error(err.Error())
				http.Error(w, html.EscapeString(FailedFileRead+fileName), http.StatusInternalServerError)
				file.Close()
				return
			}

			if _, err = file.Read(buff); err != nil {
				logging.Error(err.Error())
				http.Error(w, html.EscapeString(FailedFileRead+fileName), http.StatusInternalServerError)
				file.Close()
				return
			}

			// Validates the file type, based on MIME type
			var filetype string
			var isTextFile bool
			if err, filetype, isTextFile = ValidateFileType(buff, fileName); err != nil {
				logging.Error(err.Error())
				http.Error(w, uploadFileError, http.StatusInternalServerError)
				file.Close()
				return
			}

			if !isTextFile {
				// Copy the file bytes into a new buffer and validate these instead of the original to ensure that
				// the validation code cannot modify the file contents we will save
				fileBytes := buf.Bytes()

				// Validate the file contents, based on MIME type
				if err = ValidateFileContent(fileBytes, filetype); err != nil {
					logging.Error(err.Error())
					http.Error(w, uploadFileError, http.StatusInternalServerError)
					file.Close()
					return
				}
			}

			part, err := mw.CreateFormFile("file."+header[0].Filename, header[0].Filename)
			_, err = part.Write(buf.Bytes())
			if err != nil {
				logging.Error(err.Error())
				http.Error(w, uploadFileError, http.StatusInternalServerError)
				return
			}
		}
	}

	mw.Close()
	req, err := http.NewRequest(http.MethodPost, config.FileserverURL+"/uploadFile", &outboundReq)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, uploadFileError, http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", mw.FormDataContentType())
	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, uploadFileError, http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		respBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logging.Error(err.Error())
			respBytes = append(respBytes, []byte(err.Error())...)
		}
		http.Error(w, uploadFileError, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/search", http.StatusSeeOther)
	return
}

func getFileHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	r.ParseForm()
	fileName := r.Form.Get("FileName")
	fileType := r.Form.Get("Directory")
	fileBytes, err := fileServer.NewFsReader(config.FileserverURL).GetFile(fileName, fileType)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, html.EscapeString(FailedRetrieveFile+fileName), http.StatusInternalServerError)
		return
	}
	bodyString := string(fileBytes)

	addAjaxSecurityItems(w)
	io.WriteString(w, bodyString)
}

func isAllowedFormat(filename string, allowedFileTypes []string) bool {
	// nil for allowedFileTypes shows all filetypes
	if allowedFileTypes == nil {
		return true
	}
	splitName := strings.Split(filename, ".")
	fileType := splitName[len(splitName)-1]
	for _, allowedFileType := range allowedFileTypes {
		if strings.HasSuffix(allowedFileType, fileType) {
			return true
		}
	}

	return false
}

// SoftUI Files have their own filetype validation as they can be proprietry non-mime formats
func softUiFileValidator(filename string, filetype string) bool {
	switch filetype {
	case MenuConfiguration:
		fallthrough
	case SoftUIConfiguration:
		for _, allowedType := range allowedMenuFileTypes {
			// Get the file extension, make it lower case and drop the "."
			if strings.Replace(filepath.Ext(strings.ToLower(filename)), ".", "", -1) == allowedType {
				return true
			}
		}
	case MerchantReceipt:
		fallthrough
	case CustomerReceipt:
		for _, allowedType := range allowedReceiptFileTypes {
			// Get the file extension, make it lower case and drop the "."
			if strings.Replace(filepath.Ext(strings.ToLower(filename)), ".", "", -1) == allowedType {
				return true
			}
		}
	case IconImage:
		for _, allowedType := range allowedIconFileTypes {
			// Get the file extension, make it lower case and drop the "."
			if strings.Replace(filepath.Ext(strings.ToLower(filename)), ".", "", -1) == allowedType {
				return true
			}
		}
	}

	return false
}

func ValidateFileType(buff []byte, filename string) (error, string, bool) {
	filetype := http.DetectContentType(buff)

	for _, allowedType := range allowedImageFileTypes {
		if filetype == allowedType {
			return nil, filetype, false
		}
	}

	for _, allowedType := range allowedTextFileTypes {
		if filetype == allowedType {
			return nil, filetype, true
		}
	}

	// Allows .db and .dat files to be validated from the filename
	if isAllowedFormat(filename, allowedTextFileTypes) {
		return nil, filetype, true
	}
	return errors.New("file type not supported"), filetype, false
}

// ValidateFileContent /* Validates the contents of the file based on its MIME type */
func ValidateFileContent(fileBytes []byte, filetype string) error {
	var err error
	switch filetype {
	case "image/png":
		_, err = png.Decode(bytes.NewReader(fileBytes))
	case "image/gif":
		_, err = gif.Decode(bytes.NewReader(fileBytes))
	default:
		_, err = jpeg.Decode(bytes.NewReader(fileBytes))
	}
	return err
}

func uploadMnoLogoHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	err := r.ParseMultipartForm(10000)
	if err != nil {
		logging.Error(err.Error())
	}

	var buf bytes.Buffer
	var outboundReq bytes.Buffer
	mw := multipart.NewWriter(&outboundReq)
	mnoName := strings.TrimSpace(r.FormValue("mnoName"))
	if mnoName == "select" || mnoName == "" {
		http.Error(w, "Please select a mno name to upload", http.StatusBadRequest)
		return
	}

	if len(r.MultipartForm.File) < 1 {
		http.Error(w, "Please choose a file to upload", http.StatusBadRequest)
		return
	}

	for _, header := range r.MultipartForm.File {
		file, err := header[0].Open()
		fileName := header[0].Filename
		if err != nil {
			logging.Error(err.Error())
			http.Error(w, uploadFileError, http.StatusInternalServerError)
			return
		} else {
			if _, err := io.Copy(&buf, file); err != nil {
				logging.Error(err.Error())
				http.Error(w, uploadFileError, http.StatusInternalServerError)
				return
			}

			// 512 bytes only because DetectContentType (used in ValidateFileType) only reads up to 512 bytes
			buff := make([]byte, 512)
			if _, err = file.Seek(0, 0); err != nil {
				logging.Error(err.Error())
				http.Error(w, html.EscapeString(FailedFileRead+fileName), http.StatusInternalServerError)
				file.Close()
				return
			}
			if _, err = file.Read(buff); err != nil {
				logging.Error(err.Error())
				http.Error(w, html.EscapeString(FailedFileRead+fileName), http.StatusInternalServerError)
				file.Close()
				return
			}
			// Validates the file type, based on MIME type
			var filetype string
			var isTextFile bool
			if err, filetype, isTextFile = ValidateFileType(buff, fileName); err != nil {
				logging.Error(err.Error())
				http.Error(w, uploadFileError, http.StatusInternalServerError)
				file.Close()
				return
			}
			if !isTextFile {
				// Copy the file bytes into a new buffer and validate these instead of the original to ensure that
				// the validation code cannot modify the file contents we will save
				fileBytes := buf.Bytes()
				// Validate the file contents, based on MIME type
				if err = ValidateFileContent(fileBytes, filetype); err != nil {
					logging.Error(err.Error())
					http.Error(w, uploadFileError, http.StatusInternalServerError)
					file.Close()
					return
				}
			} else {
				http.Error(w, "Please choose only png and gif image", http.StatusBadRequest)
				return
			}
			part, err := mw.CreateFormFile("file."+fmt.Sprintf("%s.%s", mnoName, "png"), fmt.Sprintf("%s.%s", mnoName, "png"))
			if err != nil {
				logging.Error(err.Error())
				http.Error(w, uploadFileError, http.StatusInternalServerError)
				return
			}
			_, err = part.Write(buf.Bytes())
			if err != nil {
				logging.Error(err.Error())
				http.Error(w, uploadFileError, http.StatusInternalServerError)
				return
			}
		}
	}
	mw.Close()
	req, err := http.NewRequest(http.MethodPost, config.FileserverURL+"/uploadFile", &outboundReq)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, uploadFileError, http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())
	req.Header.Set("Folder", "MnoLogo")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, uploadFileError, http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		respBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logging.Error(err.Error())
			respBytes = append(respBytes, []byte(err.Error())...)
		}
		http.Error(w, uploadFileError, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/search", http.StatusSeeOther)
	return
}
