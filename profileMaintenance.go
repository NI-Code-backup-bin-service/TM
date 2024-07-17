package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"nextgen-tms-website/dal"
	"nextgen-tms-website/entities"
	"nextgen-tms-website/services"
	"strconv"
	"strings"
)

// Fetches the velocity limits for a given site
func getSiteVelocityLimits(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	err := r.ParseForm()
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, velocityLimitsError, http.StatusInternalServerError)
		return
	}

	//Retrieve the siteId from the http request
	siteId := r.Form.Get("siteId")
	site, err := strconv.Atoi(siteId)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Invalid site Id", http.StatusBadRequest)
		return
	}

	//Retrieve the tidId from the http request
	tidId := r.Form.Get("tidId")
	tid, err := strconv.Atoi(tidId)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Invalid TID Id", http.StatusBadRequest)
		return
	}

	siteLimits, velocityLimits, err := dal.GetSiteVelocityLimits(site, tid)

	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Unable to retrieve Velocity Limits", http.StatusInternalServerError)
		return
	}

	availableTxns, err := dal.GetAvailableTransactions()
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Unable to retrieve transaction types", http.StatusInternalServerError)
		return
	}

	availableLimits, err := dal.GetAvailableLimits()
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Unable to retrieve limit types", http.StatusInternalServerError)
		return
	}

	var pageModel entities.SiteFraudLimitModel
	pageModel.SiteLimits = siteLimits
	pageModel.Limits = velocityLimits
	pageModel.AvailableTransactions = availableTxns
	pageModel.AvailableLimits = availableLimits
	pageModel.HasSavePermission = checkUserPermissions(None, tmsUser)

	ajaxResponse(w, pageModel)

}

func saveSiteVelocityLimits(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	if err := r.ParseForm(); err != nil {
		logging.Error(err.Error())
		http.Error(w, saveVelocityLimitsError, http.StatusInternalServerError)
		return
	}

	//Retrieve the siteId from the http request
	siteId := r.Form.Get("siteId")
	site, err := strconv.Atoi(siteId)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Invalid site Id", http.StatusBadRequest)
		return
	}

	//Retrieve the tidId from the http request (not required for site level saving)
	tidId := r.Form.Get("tidId")
	tid, err := strconv.Atoi(tidId)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Invalid TID Id", http.StatusBadRequest)
		return
	}

	//Retrieve the velocity limits from the request
	//@NEX-12567
	limits := r.Form.Get("limits")
	limits = strings.Replace(limits, "\"{", "{", -1)
	limits = strings.Replace(limits, "}\"", "}", -1)
	limits = "[" + limits + "]"
	//TODO: Clean this up to be more generic and not be a specific field, just key: value pairs.
	dailyTxnCleanseTime := r.Form.Get("dailyTxnCleanseTime")
	profileId := 0
	profileType := ""
	if tid > 0 {
		profileType = "tid"
		profileExists, tidProfileId, err := dal.GetTidProfileIdForSiteId(tid, site)
		if err != nil {
			logging.Error(err.Error())
			http.Error(w, saveVelocityLimitsError, http.StatusInternalServerError)
			return
		}

		if !profileExists {
			tidProfileId, err = dal.SaveNewProfile("tid", strconv.Itoa(tid), 1, tmsUser.Username)
			if err != nil {
				logging.Error("an error occured during SaveNewProfile" + err.Error())
				http.Error(w, saveVelocityLimitsError, http.StatusInternalServerError)
				return
			}

			err = dal.AddTidProfileLink(tid, site, int(tidProfileId))
			if err != nil {
				logging.Error("an error occured during tid profile link" + err.Error())
				http.Error(w, saveVelocityLimitsError, http.StatusInternalServerError)
				return
			}

			groups, err := dal.GetDataGroups()
			if err != nil {
				logging.Error(err.Error())
				http.Error(w, saveVelocityLimitsError, http.StatusInternalServerError)
				return
			}

			opi := ""

			for _, g := range groups {
				if g.DataGroup == "opi" {
					opi = fmt.Sprintf("%d", g.DataGroupID)
					break
				}
			}
			if opi == "" {
				logging.Error("Could not locate opi group")
				http.Error(w, saveVelocityLimitsError, http.StatusInternalServerError)
				return
			}

			dataGroups := []string{"1", "7", "8", "9", opi}
			services.AddDataGroupsToProfile(int(tidProfileId), dataGroups, tmsUser)

			// Create a silent approved change for history auditing
			err = dal.SaveTidProfileChange(int(tidProfileId), strconv.Itoa(tid), "Override Created", tmsUser.Username, dal.ApproveCreate, 1)
			if err != nil {
				logging.Error(err.Error())
				http.Error(w, saveVelocityLimitsError, http.StatusInternalServerError)
				return
			}
		}
		profileId = int(tidProfileId)
	} else {
		profileType = "site"
		profileId, err = dal.GetProfileIdForSite(site)
	}

	if err != nil {
		logging.Error(err.Error())
		http.Error(w, saveVelocityLimitsError, http.StatusInternalServerError)
		return
	}

	elements := make(map[int]string, 0)
	elementId, err := dal.GetDataElementByName("core", "dailyTxnCleanseTime")
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, saveVelocityLimitsError, http.StatusInternalServerError)
		return
	}
	elements[elementId] = dailyTxnCleanseTime
	validationMessages, _ := validateDataElements(dal.NewValidationDal(), elements, profileId, profileType, false, "")
	if len(validationMessages) > 0 {
		validationMessagesJsonBytes, err := json.Marshal(validationMessages)
		if err != nil {
			logging.Error(err)
			http.Error(w, "An unexpected error has occurred", http.StatusInternalServerError)
			return
		}
		logging.Error(fmt.Sprintf("Validation Errors: %+v", validationMessages))
		http.Error(w, string(validationMessagesJsonBytes), http.StatusUnprocessableEntity)
		w.Header().Set("content-type", "application/json")
		return
	}

	var profileTypeName = r.Form.Get("profileTypeName")
	if tidId == "-1" {
		if err := processFlagging(r, profileId, site, profileTypeName, tmsUser); err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
	}

	//Save velocity limits for this site approval
	limitString := r.Form.Get("limitLevel")
	if err := processFraudFlagging(profileId, site, tid, profileTypeName, limitString, dailyTxnCleanseTime, limits, tmsUser); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	if err := dal.UpdateTIDFlag(tid); err != nil {
		logging.Error(err)
	}
}

func processFraudFlagging(profileId, siteID, tidId int, profileTypeName, limitString, dailyTxnCleanseTime, velocityLimits string, tmsUser *entities.TMSUser) error {
	dataElementID, err := dal.GetDataElementByName("core", "fraud")
	if err != nil {
		logging.Error(err.Error())
		return errors.New("Unable to get dataElementID")
	}
	elements := map[int]string{
		dataElementID: getFormattedString(siteID, tidId, limitString, dailyTxnCleanseTime, velocityLimits),
	}
	saveDataElements(profileId, siteID, elements, 0, tmsUser, false, 0, profileTypeName)
	return nil
}

func handleVelocityFilterRow(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	txnLimits := &entities.TransactionLimitGroupModel{}

	err := r.ParseForm()
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, filterVelocityLimitsError, http.StatusInternalServerError)
		return
	}

	limitData := r.Form.Get("TxnLimitData")
	err = json.Unmarshal([]byte(limitData), txnLimits)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, filterVelocityLimitsError, http.StatusBadRequest)
		return
	}

	override := r.Form.Get("override")
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Invalid override status", http.StatusBadRequest)
		return
	}

	availableTxns, err := dal.GetAvailableTransactions()
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Unable to retrieve transaction types", http.StatusInternalServerError)
		return
	}

	availableLimits, err := dal.GetAvailableLimits()
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Unable to retrieve limit types", http.StatusInternalServerError)
		return
	}

	pageModel := buildTxnLimitModel(txnLimits)
	pageModel.AvailableTransactions = availableTxns
	pageModel.AvailableLimits = availableLimits
	pageModel.Override = override

	renderTemplate(w, r, "txnLimitRow", pageModel, tmsUser)
}

func buildTxnLimitModel(txnLimits *entities.TransactionLimitGroupModel) entities.TransactionLimitMaintenanceModel {

	var transactionTypes []entities.TransactionLimitGroup

	//List the transaction types that currently have velocity limits
	for _, txnLimit := range txnLimits.TxnLimits {
		transactionTypes = appendIfUnique(transactionTypes, txnLimit.TxnType, txnLimit.TxnTypeReadable)
	}

	transactionTypes = buildLimits(transactionTypes, txnLimits)
	var txnLimitModel entities.TransactionLimitMaintenanceModel
	txnLimitModel.TransactionTypes = transactionTypes

	return txnLimitModel
}

func appendIfUnique(slice []entities.TransactionLimitGroup, txnName string, txnNameReadable string) []entities.TransactionLimitGroup {
	for _, ele := range slice {
		if ele.TxnLimitGroup == txnNameReadable {
			return slice
		}
	}
	newTxn := entities.TransactionLimitGroup{}
	newTxn.TxnLimitGroup = txnNameReadable
	return append(slice, newTxn)
}

func buildLimits(transactionTypes []entities.TransactionLimitGroup, txnLimits *entities.TransactionLimitGroupModel) []entities.TransactionLimitGroup {

	for i, txnType := range transactionTypes {
		var limits []entities.TxnLimit
		for _, txnLimit := range txnLimits.TxnLimits {
			if txnType.TxnLimitGroup == txnLimit.TxnTypeReadable {
				limits = append(limits, txnLimit)
			}
		}
		transactionTypes[i].TxnLimits = limits
	}
	return transactionTypes
}

func showTidFraudModal(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	r.ParseForm()
	tid := r.Form.Get("TID")

	hasPermission, err := checkUserAcquirePermsByTid(tmsUser, tid)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, tidFraudError, http.StatusInternalServerError)
		return
	}

	if !hasPermission {
		logging.Error(userUnauthorisedError)
		http.Error(w, "User is not permitted to perform this action", http.StatusForbidden)
		return
	}

	// Collate the information passed from the request
	siteProfileId, err := strconv.Atoi(r.Form.Get("SiteProfileId"))
	if err != nil {
		handleError(w, errors.New(tidFraudError), tmsUser)
		return
	}
	siteId, err := strconv.Atoi(r.Form.Get("SiteId"))
	if err != nil {
		handleError(w, errors.New(tidFraudError), tmsUser)
		return
	}
	tidId, err := strconv.Atoi(tid)
	if err != nil {
		handleError(w, errors.New(tidFraudError), tmsUser)
		return
	}
	modalModel := TidModalModel{Tid: tid}

	// Check to see if the TID has its own profile ID, this will only be the case if the TIDs config is overridden
	validTidProfile, tidProfileId, err := dal.GetTidProfileIdForSiteId(tidId, siteId)
	if err != nil {
		handleError(w, errors.New(tidFraudError), tmsUser)
		return
	}

	// If the TID does not have its own profile, then use the site profile
	if !validTidProfile {
		tidProfileId = int64(siteProfileId)
	}

	// Find the schemes applicable to this TID
	availableSchemes, err := dal.GetAvailableSchemesForSiteId(siteId)
	if err != nil {
		handleError(w, errors.New(tidFraudError), tmsUser)
		return
	}

	modalModel.AvailableSchemes = availableSchemes

	// If the TID does not have fraud overrides set, then we need to use the Site's cleanse time
	tidFraudOverridden, err := dal.TidFraudOverrideStatus(tidId)
	if err != nil {
		handleError(w, errors.New(tidFraudError), tmsUser)
		return
	}
	if !tidFraudOverridden {
		tidProfileId = int64(siteProfileId)
	}

	dataGroups, err := dal.GetProfileDataForTabByProfileId(int(tidProfileId), "fraud")
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Error retrieving fraud groups", http.StatusInternalServerError)
		return
	}
	modalModel.FraudGroups = SortGroups(dataGroups)
	modalModel.CurrentUser = tmsUser
	renderPartialTemplate(w, r, "tidVelocityLimits", modalModel, tmsUser)
}

func handleTidFraudClose(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	r.ParseForm()
	pageSize := 10  // The number of TIDs to display on a site page at once
	pageNumber := 1 // The current page of TIDs we are on
	tidSearchTerm := r.Form.Get("searchTerm")

	// Check the form to see if an updated pageSize has been passed in
	newPageSize, psErr := strconv.Atoi(r.Form.Get("pageSize"))
	if psErr == nil {
		// If there are no errors, then a new pageSize has been sent
		pageSize = newPageSize
	}

	// Check the form to see if an updated pageNumber has been passed in
	newPageNumber, psErr := strconv.Atoi(r.Form.Get("pageNumber"))
	if psErr == nil {
		// If there are no errors, then a new pageSize has been sent
		pageNumber = newPageNumber
	}

	idString := r.Form.Get("profileId")
	id, err := strconv.Atoi(idString)
	if err != nil {
		handleError(w, errors.New(profileConversionFailed), tmsUser)
		return
	}

	profileType, err := dal.GetTypeForProfile(id)
	if err != nil {
		logging.Error(err.Error())
		handleError(w, errors.New(tidFraudError), tmsUser)
		return
	}

	siteId, err := dal.GetSiteFromProfile(id)
	if err != nil {
		handleError(w, errors.New("no siteId found with provided id"), tmsUser)
		return
	}

	p := buildSiteProfileMaintenanceModel(w, id, tmsUser, pageSize, pageNumber, tidSearchTerm, siteId)
	p.ProfileTypeName = profileType

	renderTemplate(w, r, "profileMaintenanceTIDs", p, tmsUser)
}

func deleteTIDVelocityOverrides(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	err := r.ParseForm()
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, deleteTIDVelocityOverrideError, http.StatusInternalServerError)
		return
	}

	//Retrieve the tidId from the http request
	tidId := r.Form.Get("tidId")
	logging.Information(fmt.Sprintf("Deleting TID velocity override for TID '%v'", tidId))
	_, err = dal.NewPEDRepository().DeleteFraudOverrideByTid(tidId)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, deleteTIDVelocityOverrideError, http.StatusInternalServerError)
		return
	}
	ajaxResponse(w, nil)
}

func getFormattedString(siteID, tidId int, limitString, dailyTxnCleanseTime, velocityLimits string) string {
	return fmt.Sprintf("{\"siteID\":%d, \"tidId\":%d, \"limitString\":%s, \"dailyTxnCleanseTime\":\"%s\",\"velocityLimits\":%s}", siteID, tidId, limitString, dailyTxnCleanseTime, velocityLimits)
}
