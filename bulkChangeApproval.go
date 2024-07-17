package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"nextgen-tms-website/services"
	"strconv"

	"nextgen-tms-website/common"
	"nextgen-tms-website/config"
	"nextgen-tms-website/dal"
	"nextgen-tms-website/entities"
	"nextgen-tms-website/fileServer"
)

const (
	BCATAG = "Bulk Change Approval : "
)

func bulkChangeApprovalHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	logging.Debug(BCATAG, "bulkChangeApprovalHandler()")
	renderHeader(w, r, tmsUser)
	renderTemplate(w, r, "bulkChangeApproval", nil, tmsUser)
}

func bulkChangeApprovalHistory(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	history, err := dal.GetAllApprovedAndRejectedBulkApprovals()
	if err != nil {
		http.Error(w, DatabaseTxnError, http.StatusBadRequest)
		return
	}

	renderTemplate(w, r, "bulkChangeApprovalHistory", history, tmsUser)
}

func unapprovedbulkChangeApproval(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	unapproved, err := dal.GetAllUnapprovedBulkApprovals()
	if err != nil {
		http.Error(w, DatabaseTxnError, http.StatusBadRequest)
		return
	}

	renderTemplate(w, r, "bulkChangeApprovalApprove", unapproved, tmsUser)
}

func approveBulkChangeApprovalHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	fileName := r.Form.Get("FileName")
	fileType := r.Form.Get("FileType")
	changeType := r.Form.Get("ChangeType")

	if fileName == "" {
		unapproved, err := dal.GetAllUnapprovedBulkApprovals()
		if err != nil {
			http.Error(w, DatabaseTxnError, http.StatusBadRequest)
			return
		}

		for _, record := range unapproved {
			if record.ChangeType != "Delete" {
				err := approveBulkApproval(record.Filename, record.FileType, tmsUser)
				if err != nil {
					http.Error(w, DatabaseTxnError, http.StatusBadRequest)
					return
				}
			} else {
				err := approveBulkDeleteFileType(record.Filename, record.FileType, tmsUser.Username)
				if err != nil {
					http.Error(w, DatabaseTxnError, http.StatusBadRequest)
					return
				}
			}
		}
	} else if fileName != "" && changeType == "Delete" {
		err := approveBulkDeleteFileType(fileName, fileType, tmsUser.Username)
		if err != nil {
			http.Error(w, DatabaseTxnError, http.StatusBadRequest)
			return
		}
	} else {
		err := approveBulkApproval(fileName, fileType, tmsUser)
		if err != nil {
			http.Error(w, DatabaseTxnError, http.StatusBadRequest)
			return
		}
	}

}

func approveBulkApproval(fileName, fileType string, tmsUser *entities.TMSUser) error {

	var configDir string
	switch fileType {
	case TerminalFlaggingType:
		configDir = config.FlaggingFileDirectory
	case BulkSiteUpdateType:
		configDir = config.BulkSiteUpdateDirectory
	case BulkTidUpdateType:
		configDir = config.BulkTidUpdateDirectory
	case BulkTidDeleteType:
		configDir = config.BulkTidDeleteDirectory
	case BulkSiteDeleteType:
		configDir = config.BulkSiteDeleteDirectory
	case BulkPaymentServiceUploadType:
		configDir = config.BulkPaymentUploadDirectory
	case BulkPaymentTidUploadType:
		configDir = config.BulkPaymentUploadDirectory
	}

	fileAsBase64Encoded, err := fileServer.NewFsReader(config.FileserverURL).GetFile(fileName, configDir)
	if err != nil {
		logging.Error("Unable to get file from file server : " + err.Error())
		return err
	}

	fileAsBytes, err := common.ConvertBase64FileToBytes(string(fileAsBase64Encoded))
	if err != nil {
		logging.Error("Unable to convert Base64 File to bytes : " + err.Error())
		return err
	}

	csvReader := csv.NewReader(bytes.NewBuffer(fileAsBytes))
	csvReader.LazyQuotes = true
	records, err := csvReader.ReadAll()
	if err != nil {
		logging.Error("Unable to read the file as csv : " + err.Error())
		return err
	}

	// converting array to slice
	n := records[:]

	switch fileType {
	case TerminalFlaggingType:
		err = dal.ApproveBulkApprovalTerminalFlagging(fileName, tmsUser.Username, n)
		if err != nil {
			return err
		}
	case BulkSiteUpdateType:
		err = dal.ApproveBulkApprovalBulkSiteUpdate(fileName, tmsUser.Username, n)
		if err != nil {
			return err
		}
	case BulkTidUpdateType:
		err = dal.ApproveBulkApprovalBulkTidUpdate(fileName, tmsUser.Username, n)
		if err != nil {
			return err
		}
	case BulkTidDeleteType:
		err = dal.ApproveBulkApprovalBulkDelete(fileName, tmsUser.Username, n, BulkTidDeleteType)
		if err != nil {
			return err
		}
	case BulkSiteDeleteType:
		err = dal.ApproveBulkApprovalBulkDelete(fileName, tmsUser.Username, n, BulkSiteDeleteType)
		if err != nil {
			return err
		}
	case BulkPaymentServiceUploadType:
		err = approveBulkPaymentServiceConfigUpload(fileName, tmsUser.Username, n, BulkPaymentServiceUploadType)
		if err != nil {
			return err
		}
	case BulkPaymentTidUploadType:
		err = approveBulkPaymentTIDConfigUpload(fileName, n, BulkPaymentTidUploadType, tmsUser)
		if err != nil {
			return err
		}
	}

	buff := &bytes.Buffer{}
	w := csv.NewWriter(buff)
	w.WriteAll(n)

	if fileType == TerminalFlaggingType {
		err = sendFileToFileServer(buff.Bytes(), "report_"+fileName, TerminalFlaggingType)
		if err != nil {
			logging.Error(TFTAG, "Terminal Flagging File Upload Failed : ", err.Error())
			return err
		}
	}

	return nil
}

func discardBulkChangeApproval(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	fileName := r.Form.Get("FileName")
	fileType := r.Form.Get("FileType")
	err := dal.DiscardBulkApproval(fileName, fileType, tmsUser.Username)
	if err != nil {
		http.Error(w, DatabaseTxnError, http.StatusBadRequest)
		return
	}
}

func sendFileToFileServer(buff []byte, fileName string, fileType string) error {
	var outboundReq bytes.Buffer
	mw := multipart.NewWriter(&outboundReq)

	part, err := mw.CreateFormFile("file."+fileName, fileName)
	if err != nil {
		logging.Error(err.Error())
		return err
	}

	_, err = part.Write(buff)
	if err != nil {
		logging.Error(err.Error())
		return err
	}

	logging.Debug("Form file created with file name : " + fileName)
	err = mw.Close()
	if err != nil {
		logging.Error(err.Error())
		return err
	}

	req, err := http.NewRequest(http.MethodPost, config.FileserverURL+"/uploadFile", &outboundReq)
	if err != nil {
		logging.Error(err.Error())
		return err
	}

	req.Header.Set("Content-Type", mw.FormDataContentType())

	switch fileType {
	case TerminalFlaggingType:
		req.Header.Set("Folder", config.FlaggingFileDirectory)
	case BulkSiteUpdateType:
		req.Header.Set("Folder", config.BulkSiteUpdateDirectory)
	case BulkTidUpdateType:
		req.Header.Set("Folder", config.BulkTidUpdateDirectory)
	case BulkTidDeleteType:
		req.Header.Set("Folder", config.BulkTidDeleteDirectory)
	case BulkSiteDeleteType:
		req.Header.Set("Folder", config.BulkSiteDeleteDirectory)
	case BulkPaymentServiceUploadType:
		req.Header.Set("Folder", config.BulkPaymentUploadDirectory)
	case BulkPaymentTidUploadType:
		req.Header.Set("Folder", config.BulkPaymentUploadDirectory)
	}

	client := http.Client{}
	logging.Debug("Sending uploadFile request to : " + config.FileserverURL)
	resp, err := client.Do(req)
	if err != nil {
		logging.Error(err.Error())
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logging.Error(err.Error())
			return err
		} else {
			logging.Error(string(respBytes))
			return errors.New(string(respBytes))
		}
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logging.Error(err.Error())
		return err
	}
	logging.Debug("Uploaded file to fileserver : " + string(respBytes))
	return nil
}

func approveBulkDeleteFileType(fileName, fileType, username string) error {
	values := make(map[string][]string, 0)
	values["FileName"] = []string{fileName}
	if fileType == LogoManagementType {
		values["Directory"] = []string{"mnoLogo"}
	}
	response, err := http.PostForm(config.FileserverURL+"/deleteFile", values)
	if err != nil {
		logging.Error(err.Error())
		return err
	}
	defer response.Body.Close()
	err = dal.ApproveBulkApprovalBulkDelete(fileName, username, nil, fileType)
	if err != nil {
		return err
	}

	return nil
}

func approveBulkPaymentServiceConfigUpload(fileName, currentUser string, records [][]string, fileType string) error {
	_, groupIdx, serviceIdx, _, _, _, _ := getPaymentServiceColumnIndexes(records[0])
	model := &entities.PaymentServiceGroupImportModel{
		GroupsCreated:   0,
		ServicesCreated: 0,
		FailedRows:      0,
		Groups:          map[string]entities.PaymentServiceImportGroup{},
	}
	for _, record := range records[1:] {
		groupName := record[groupIdx]
		serviceName := record[serviceIdx]
		if _, ok := model.Groups[groupName]; !ok {
			groupId := dal.GroupIdFromName(groupName)
			if groupId == -1 {
				model.GroupsCreated++
			}
			model.Groups[groupName] = entities.PaymentServiceImportGroup{
				Id:       groupId,
				Services: map[string]bool{},
			}
		}
		group := model.Groups[groupName]
		if group.Id != -1 && dal.ServiceAlreadyExists(serviceName, strconv.Itoa(group.Id), nil) {
			model.FailedRows++
			continue
		}
		group.Services[serviceName] = true
		model.ServicesCreated++
	}
	PsUpload = model
	if PsUpload == nil || len(PsUpload.Groups) == 0 {
		return nil
	}
	var err error
	successfulRows := 0
	for key, val := range PsUpload.Groups {
		if val.Id == -1 {
			err = dal.AddPaymentServiceGroup(key)
			if err != nil {
				continue
			}
			val.Id = dal.GroupIdFromName(key)
			if val.Id == -1 {
				continue
			}
		}
		for serviceName, _ := range val.Services {
			err = dal.AddPaymentService(serviceName, strconv.Itoa(val.Id))
			if err != nil {
				continue
			}
			successfulRows++
		}
	}
	err = dal.ApproveBulkPaymentUpload(fileName, currentUser, nil, fileType)
	if err != nil {
		return err
	}
	return nil
}

func approveBulkPaymentTIDConfigUpload(fileName string, records [][]string, fileType string, tmsUser *entities.TMSUser) error {

	siteConf := map[string]siteInfo{}
	PsTermUpload = map[tidIdentifier]map[int]*dal.PaymentService{}
	successfulRows := 0

	_, groupIdx, serviceIdx, mMidIdx, mTidIdx, sMidIdx, sTidIdx := getPaymentServiceColumnIndexes(records[0])

	for _, record := range records[1:] {
		if anyColumnsEmpty(record, groupIdx, serviceIdx, mMidIdx, mTidIdx, sMidIdx, sTidIdx) {
			continue
		}

		mid := record[mMidIdx]
		tid := record[mTidIdx]
		sTid := record[sTidIdx]

		conf, ok := siteConf[mid]
		// we've already validated this mid - so we just need to validate the group and service names
		if ok {
			if len(conf.groupName) > 0 && conf.siteId > 0 && conf.groupName == record[groupIdx] && dal.IsServiceInGroup(record[groupIdx], record[serviceIdx]) {
				if appendToPsUploadModel(record[groupIdx], record[serviceIdx], record[sMidIdx], conf.siteId, tid, sTid) {
					successfulRows++
				}
			}
			continue
		}

		// we haven't encountered this mid yet, so we need to validate it
		profileIdStr, err := dal.GetProfileIdFromMID(mid)
		if err != nil || len(profileIdStr) == 0 {
			siteConf[mid] = siteInfo{-1, ""}
			continue
		}
		profileId, err := strconv.Atoi(profileIdStr)
		if err != nil || profileId < 1 {
			siteConf[mid] = siteInfo{-1, ""}
			continue
		}

		groupName := dal.GetProfileAssignedServiceGroupName(profileIdStr)
		siteConf[mid] = siteInfo{profileId, groupName}
		if len(groupName) > 0 && groupName == record[groupIdx] && dal.IsServiceInGroup(record[groupIdx], record[serviceIdx]) {
			if appendToPsUploadModel(record[groupIdx], record[serviceIdx], record[sMidIdx], profileId, tid, sTid) {
				successfulRows++
			}
		}
	}

	if PsTermUpload == nil || len(PsTermUpload) == 0 {
		return nil
	}

	rowsUpdated := 0
	for id, upload := range PsTermUpload {
		serviceJson := MergePaymentsWithConfigured(upload, id.tid, id.siteId)
		bytes, err := json.Marshal(serviceJson)
		if err != nil {
			continue
		}

		if services.SaveTerminalPaymentServiceConfig(id.tid, string(bytes), tmsUser) {
			rowsUpdated++
		}
	}

	err := dal.ApproveBulkPaymentUpload(fileName, tmsUser.Username, nil, fileType)
	if err != nil {
		return err
	}

	return nil

}
