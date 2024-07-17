package main

import (
	"encoding/csv"
	"errors"
	"net/http"
	"nextgen-tms-website/dal"
	"nextgen-tms-website/entities"
	"strconv"
	"time"
)

func changeApprovalHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	r.ParseForm()
	model := buildChangeApprovalModel(tmsUser)
	renderHeader(w, r, tmsUser)
	renderTemplate(w, r, "changeApprovalViewer", model, tmsUser)
}

func filterChangeApprovalHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	r.ParseForm()
	offset, err := strconv.Atoi(r.Form.Get("Offset"))
	if err != nil {
		handleError(w, errors.New("an error occured while retriving/converting the Offset after parsing the form data:"+err.Error()), tmsUser)
		return
	}
	profileType := r.Form.Get("Type")
	profileTypeID := getProfileTypeFromName(profileType)
	model := buildFilteredChangeApprovalModel(offset, profileTypeID, tmsUser)
	model.IdentifierColumn = profileType
	renderPartialTemplate(w, r, "changeApprovalViewerPartial", model, tmsUser)
}

func filterChangeApprovalHistoryHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	r.ParseForm()
	after := r.Form.Get("After")
	name := r.Form.Get("Name")
	user := r.Form.Get("User")
	before := r.Form.Get("Before")
	field := r.Form.Get("Field")
	offset, err := strconv.Atoi(r.Form.Get("Offset"))
	if err != nil {
		handleError(w, errors.New("an error occured while retriving/converting the Offset after parsing the form data"+err.Error()), tmsUser)
		return
	}
	model, err := buildFilteredChangeApprovalHistory(after, name, user, before, field, true, offset, tmsUser)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, changeApprovalHistoryError, http.StatusInternalServerError)
		return
	}
	renderTemplate(w, r, "changeApprovalViewerHistoryPartial", model, tmsUser)
}

func exportFilteredChangeApprovalHistory(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	r.ParseForm()
	after := r.Form.Get("After")
	name := r.Form.Get("Name")
	user := r.Form.Get("User")
	before := r.Form.Get("Before")
	field := r.Form.Get("Field")
	offset, err := strconv.Atoi(r.Form.Get("Offset"))
	if err != nil {
		handleError(w, errors.New("an error occured while retriving/converting the Offset after parsing the form data"), tmsUser)
		return
	}
	model, err := buildFilteredChangeApprovalHistory(after, name, user, before, field, false, offset, tmsUser)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, exportChangeApprovalHistoryError, http.StatusInternalServerError)
		return
	}

	w.Header().Set("fileName", "ChangeApprovalHistory_"+time.Now().Format("02-01-2006-15-04-05")+".csv")
	wr := csv.NewWriter(w)
	var records [][]string

	var record []string
	record = append(record, "Name")
	record = append(record, "Field")
	record = append(record, "Original Value")
	record = append(record, "Update Value")
	record = append(record, "Updated By")
	record = append(record, "Updated At")
	record = append(record, "Reviewed By")
	record = append(record, "Reviewed At")
	record = append(record, "Approved/Discarded")
	records = append(records, record)

	for _, obj := range model.HistoryTab.History {
		var approved string
		if obj.Approved == 1 {
			approved = "Approved"
		} else {
			approved = "Discarded"
		}

		record = []string{}
		record = append(record, obj.Identifier)
		record = append(record, obj.Field)
		record = append(record, obj.OriginalValue)
		record = append(record, obj.ChangeValue)
		record = append(record, obj.ChangedBy)
		record = append(record, obj.ChangedAt)
		record = append(record, obj.ReviewedBy)
		record = append(record, obj.ReviewedAt)
		record = append(record, approved)

		records = append(records, record)
	}

	err = wr.WriteAll(records)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, exportChangeApprovalHistoryError, http.StatusInternalServerError)
		return
	}
}

func approveAllChangesHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	r.ParseForm()
	offset, err := strconv.Atoi(r.Form.Get("Offset"))
	if err != nil {
		handleError(w, errors.New("an error occured while retriving/converting the Offset:"+err.Error()), tmsUser)
		return
	}
	profileType := r.Form.Get("Type")
	profileTypeID := getProfileTypeFromName(profileType)
	model := buildFilteredChangeApprovalModel(offset, profileTypeID, tmsUser)
	model.IdentifierColumn = profileType
	err = dal.ApproveAllChanges(model.History, tmsUser.Username, profileType)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, approveChangesError, http.StatusInternalServerError)
		return
	}
}

func approveChangeHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	r.ParseForm()
	profileDataID, err := strconv.Atoi(r.Form.Get("profileDataID"))
	if err != nil {
		handleError(w, errors.New("an error occured while retriving/converting the profileDataID"+err.Error()), tmsUser)
		return
	}
	profileType := r.Form.Get("Type")
	err = dal.ApproveChange(profileDataID, tmsUser.Username, profileType)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, approveChangesError, http.StatusInternalServerError)
		return
	}
}

func discardAllChangesHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	r.ParseForm()
	offset, err := strconv.Atoi(r.Form.Get("Offset"))
	if err != nil {
		handleError(w, errors.New("an error occured while retriving/converting the Offset after parsing the form data"+err.Error()), tmsUser)
		return
	}
	profileType := r.Form.Get("Type")
	profileTypeID := getProfileTypeFromName(profileType)
	model := buildFilteredChangeApprovalModel(offset, profileTypeID, tmsUser)
	model.IdentifierColumn = profileType
	err = dal.DiscardAllChanges(model.History, tmsUser.Username)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, discardChangesError, http.StatusInternalServerError)
		return
	}
}

func discardChangeHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	r.ParseForm()
	profileDataID, err := strconv.Atoi(r.Form.Get("profileDataID"))
	if err != nil {
		handleError(w, errors.New("an error occured while retriving/converting the profileDataID"+err.Error()), tmsUser)
		return
	}
	err = dal.DiscardChange(profileDataID, tmsUser.Username)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, discardChangesError, http.StatusInternalServerError)
		return
	}
}

func buildChangeApprovalModel(user *entities.TMSUser) ChangeApprovalModel {
	var model ChangeApprovalModel
	siteHistory, length, err := dal.GetPendingChanges(dal.Site, user)
	if err != nil {
		logging.Error("an error occured during retreiving the pending site changes:" + err.Error())
	}
	model.SiteTab = &ChangeApprovalTabModel{History: siteHistory, IdentifierColumn: "Site", Count: length, TabType: "site"}
	model.SiteTab.CurrentUser = user

	chainHistory, length, err := dal.GetPendingChanges(dal.Chain, user)
	if err != nil {
		logging.Error("an error occured during retreiving the pending Chain changes:" + err.Error())
	}
	model.ChainTab = &ChangeApprovalTabModel{History: chainHistory, IdentifierColumn: "Chain", Count: length, TabType: "chain"}
	model.ChainTab.CurrentUser = user

	acquirerHistory, length, err := dal.GetPendingChanges(dal.Acquirer, user)
	if err != nil {
		logging.Error("an error occured during retreiving the pending Acquirer changes:" + err.Error())
	}
	model.AcquirerTab = &ChangeApprovalTabModel{History: acquirerHistory, IdentifierColumn: "Acquirer", Count: length, TabType: "acquirer"}
	model.AcquirerTab.CurrentUser = user

	tidHistory, length, err := dal.GetPendingChanges(dal.Tid, user)
	if err != nil {
		logging.Error("an error occured during retreiving the pending Tid changes:" + err.Error())
	}
	model.TidTab = &ChangeApprovalTabModel{History: tidHistory, IdentifierColumn: "TID", Count: length, TabType: "tid"}
	model.TidTab.CurrentUser = user

	OthersHistory, length, err := dal.GetPaymentServicePendingChanges(dal.Others, user)
	if err != nil {
		logging.Error("an error occured during retreiving the pending PaymentService changes:" + err.Error())
	}
	model.OthersTab = &ChangeApprovalTabModel{History: OthersHistory, IdentifierColumn: "Others", Count: length, TabType: "others"}
	model.OthersTab.CurrentUser = user

	model.HistoryTab = &ChangeApprovalTabModel{}
	model.HistoryTab.CurrentUser = user
	return model
}

func buildFilteredChangeApprovalModel(offset int, profileTypeId int, user *entities.TMSUser) ChangeApprovalTabModel {
	var model ChangeApprovalTabModel
	if profileTypeId != dal.Others {
		model.History, _ = dal.GetFilteredPendingChanges(offset, profileTypeId, user)
	} else {
		changes, _, _ := dal.GetPaymentServicePendingChanges(profileTypeId, user)
		start := offset
		end := offset + 50
		length := len(changes)

		// range completely within array bounds
		if start <= length && end <= length {
			changes = changes[start:end]
		}

		// range overruns bounds
		if start <= length && end > length {
			changes = changes[start:]
		}
		model.History = changes
	}
	model.CurrentUser = user
	return model
}

func buildFilteredChangeApprovalHistory(after string, name string, user string, before string, field string, limit bool, offset int, tmsUser *entities.TMSUser) (ChangeApprovalModel, error) {
	var model ChangeApprovalModel
	modelHistory, err := dal.BuildFilteredChangeApprovalHistory(after, name, user, before, field, limit, offset, tmsUser)
	if err != nil {
		return ChangeApprovalModel{}, err
	}
	model.HistoryTab = &ChangeApprovalTabModel{History: modelHistory, IdentifierColumn: "Name", TabType: "change"}
	model.HistoryTab.CurrentUser = tmsUser
	return model, nil
}

func getProfileTypeFromName(typeName string) int {
	switch typeName {
	case "Site":
		return dal.Site
	case "Chain":
		return dal.Chain
	case "Acquirer":
		return dal.Acquirer
	case "TID":
		return dal.Tid
	case "Others":
		return dal.Others
	default:
		return -1
	}
}

func (t ChangeApprovalModel) GetTableRowStyle(approved int) string {
	if approved == 1 {
		return "table-approved"
	} else if approved == -1 {
		return "table-disapproved"
	} else {
		return ""
	}
}
