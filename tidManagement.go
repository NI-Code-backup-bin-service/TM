package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"nextgen-tms-website/common"
	"nextgen-tms-website/dal"
	"nextgen-tms-website/entities"
	"nextgen-tms-website/logger"
	"nextgen-tms-website/models"
	"nextgen-tms-website/services"
	"nextgen-tms-website/validation"
	"strconv"
	"strings"
	"time"
)

func addNewTIDHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	r.ParseForm()

	tidIndex, err := strconv.Atoi(r.Form.Get("TidIndex"))
	if err != nil {
		handleError(w, errors.New("An error occured during parsing TidIndex to int"), tmsUser)
		return
	}
	siteId, err := strconv.Atoi(r.Form.Get("Site"))
	if err != nil {
		handleError(w, errors.New(siteConversionFailed), tmsUser)
		return
	}

	profileId, err := dal.GetProfileFromSite(siteId)
	if err != nil {
		handleError(w, errors.New("An error occured during retriving the profile from site"), tmsUser)
		return
	}

	model := buildAddNewTidModel(w, profileId, tidIndex, siteId)
	model.IsDuplicate = false // ignore model.DuplicatedFrom intentionally since it isn't duplicated

	renderTemplate(w, r, "addTid", model, tmsUser)
}

func buildAddNewTidModel(w http.ResponseWriter, profileID int, tidIndex, siteID int) AddTidModel {
	var model AddTidModel

	DefaultTidGroups, err := dal.GetTIDSiteData(siteID, profileID)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, addTidError, http.StatusInternalServerError)
		return model
	}

	model.TID = tidIndex
	model.SiteId = siteID
	model.ProfileId = profileID
	model.DefaultTidGroups = SortGroups(DefaultTidGroups)
	model.IsNew = true
	return model
}

func updateSerialNumber(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	if err := r.ParseForm(); err != nil {
		logging.Error(err.Error())
		http.Error(w, "Form parse error", http.StatusInternalServerError)
		return
	}

	tid := r.Form.Get("TID")
	siteProfileId := r.Form.Get("SITE_PROFILE_ID")
	oldSn := r.Form.Get("OLD_SN")
	newSn := r.Form.Get("NEW_SN")

	hasPermission, err := checkUserAcquirePermsByTid(tmsUser, tid)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, applyTidUpdatesError, http.StatusInternalServerError)
		return
	}
	if !hasPermission {
		logging.Error(userUnauthorisedError)
		http.Error(w, userUnauthorisedError, http.StatusForbidden)
		return
	}
	updates := getUpdateFields(r.Form)

	valid, msg := handleUpdateSerialNumberValidation(updates, newSn)
	if !valid {
		logger.GetLogger().Error(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	err = services.UpdateTIDHandler(tmsUser, oldSn, newSn, siteProfileId, updates, tid)
	if err != nil {
		logging.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func addTIDHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	if err := r.ParseMultipartForm(10000); err != nil {
		logging.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var tidIndex = r.Form.Get("tidID")
	siteId, err := strconv.Atoi(r.Form.Get("siteID"))
	if err != nil {
		handleError(w, errors.New(siteConversionFailed), tmsUser)
		return
	}
	profileId, err := strconv.Atoi(r.Form.Get("profileID"))
	if err != nil {
		handleError(w, errors.New(profileConversionFailed), tmsUser)
		return
	}
	var tid = r.Form.Get("newTidInput" + tidIndex)
	var serial = r.Form.Get("newSerialInput" + tidIndex)

	hasPermission, err := checkUserAcquirePermsBySite(tmsUser, siteId)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "unable to check user permissions", http.StatusInternalServerError)
		return
	}
	if !hasPermission {
		logging.Error(userUnauthorisedError)
		http.Error(w, userUnauthorisedError, http.StatusForbidden)
		return
	}

	if err := services.ValidateTidCreation(tid, serial); err != nil {
		logger.GetLogger().Information(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := services.AddTid(tid, serial, tmsUser.Username, siteId, profileId); err != nil {
		logging.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Form.Get("IsDuplicate") == "true" {
		tidId, _ := strconv.Atoi(tid)
		newProfileId, err := dal.SaveNewProfile("tid", tid, 1, tmsUser.Username)
		if err != nil {
			logging.Error("error saving duplicate tid profile", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = dal.AddTidProfileLink(tidId, siteId, int(newProfileId))
		if err != nil {
			logging.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = services.AddTidOverrideDataGroups(int(newProfileId), siteId, tmsUser)
		if err != nil {
			logging.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//Create a silent approved override for history auditing
		duplicatedFrom := r.Form.Get("DuplicatedFrom")
		note := "Tid duplicated from " + duplicatedFrom
		dataElements := getElementFields(r.Form, false)
		saveDataElements(int(newProfileId), siteId, dataElements, 1, tmsUser, false, 0)
		err = dal.SaveTidProfileChange(int(newProfileId), tid, note, tmsUser.Username, dal.ApproveNewElement, 1)
		if err != nil {
			logging.Error("failed to save profile duplicate terminal ", tid, err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	p := buildSiteProfileMaintenanceModel(w, profileId, tmsUser, 10, 1, "", siteId)
	renderTemplate(w, r, "profileMaintenanceTIDs", p, tmsUser)
}

func deleteTIDHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	if err := r.ParseForm(); err != nil {
		logging.Error(err.Error())
		http.Error(w, "an error occurred deleting tid", http.StatusInternalServerError)
		return
	}

	tid := r.Form.Get("TID")
	siteId, err := strconv.Atoi(r.Form.Get("Site"))
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "an error occurred deleting tid", http.StatusInternalServerError)
		return
	}

	hasPermission, err := checkUserAcquirePermsBySite(tmsUser, siteId)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "an error occurred deleting tid", http.StatusInternalServerError)
		return
	}
	if !hasPermission {
		logging.Error(err.Error())
		http.Error(w, userUnauthorisedError, http.StatusForbidden)
		return
	}

	id, err := dal.GetProfileFromSite(siteId)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "an error occurred deleting tid", http.StatusInternalServerError)
	}

	t, err := strconv.Atoi(tid)
	if err != nil {
		handleError(w, errors.New(tidConversionFailed), tmsUser)
		return
	}

	tidString := dal.GetPaddedTidId(t)
	err = dal.SaveTidProfileChange(id, tidString, "TID Deleted", tmsUser.Username, dal.ApproveDelete, 0)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "an error occurred saving profile change approval", http.StatusInternalServerError)
	}
	p := buildSiteProfileMaintenanceModel(w, id, tmsUser, 10, 1, "", siteId)
	renderTemplate(w, r, "profileMaintenanceTIDs", p, tmsUser)
}

func BuildUpdatesTIDModel(w http.ResponseWriter, tidId string, siteID int, partialPackageName string) TidUpdatesModel {
	var model TidUpdatesModel
	var updates []*models.TIDUpdateData
	var packages []*dal.PackageData

	updates, err := dal.GetTIDUpdates(tidId)
	if err != nil {
		logging.Error(err.Error())
	}

	packages, err = dal.GetPackages()
	if err != nil {
		logging.Error(err.Error())
	}

	var thirdPartyApks []*dal.ThirdPartyApk
	var preFixNames []string

	thirdPartyApks, err = dal.GetThirdPartyApks(partialPackageName)
	if err != nil {
		logging.Error(err.Error())
	}
	model.ThirdPartyModeActive = true

	preFixName, err := dal.GetProfileIdForTIDAndSiteID(tidId, siteID)
	if err != nil {
		logging.Error(err.Error())
	}
	if err == nil {
		err = json.Unmarshal([]byte(preFixName), &preFixNames)
		if err != nil {
			logging.Error(err.Error())
		}
	}

	tidIdInt, err := strconv.Atoi(tidId)
	if err != nil {
		logging.Error(err.Error())
	}

	profileExists, tidProfileId, err := dal.GetTidProfileIdForSiteId(tidIdInt, siteID)
	if err != nil {
		logging.Error(err.Error())
	}
	if profileExists {
		model.ProfileID = strconv.FormatInt(tidProfileId, 10)
	}

	model.TID = tidId
	model.SiteID = siteID
	model.Updates = updates
	model.Packages = packages
	model.ThirdPartyApks = thirdPartyApks
	model.PreFixNames = preFixNames

	return model
}

func updatesTIDHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	r.ParseForm()
	tid := r.Form.Get("TID")
	siteId, err := strconv.Atoi(r.Form.Get("Site"))
	if err != nil {
		handleError(w, errors.New(siteConversionFailed), tmsUser)
		return
	}

	hasPermission, err := checkUserAcquirePermsByTid(tmsUser, tid)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, updateTidError, http.StatusInternalServerError)
		return
	}
	if !hasPermission {
		logging.Error(userUnauthorisedError)
		http.Error(w, userUnauthorisedError, http.StatusForbidden)
		return
	}

	model := BuildUpdatesTIDModel(w, tid, siteId, "")

	minimumSoftwareVersion, err := dal.GetMinimumRequiredSoftwareVersionForSite(siteId)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, updateTidError, http.StatusInternalServerError)
	}
	model.MinimumSoftwareVersion = minimumSoftwareVersion

	renderTemplate(w, r, "updateDetailsTID", model, tmsUser)
}

func updatesThirdPartyHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	r.ParseForm()
	tid := r.Form.Get("TID")
	siteId, err := strconv.Atoi(r.Form.Get("Site"))
	if err != nil {
		handleError(w, errors.New(siteConversionFailed), tmsUser)
		return
	}
	partialPackageName := r.Form.Get("partialPackageName")

	hasPermission, err := checkUserAcquirePermsByTid(tmsUser, r.Form.Get("TID"))
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, updateTidError, http.StatusInternalServerError)
		return
	}
	if !hasPermission {
		logging.Error(userUnauthorisedError)
		http.Error(w, userUnauthorisedError, http.StatusForbidden)
		return
	}
	model := BuildUpdatesTIDModel(w, tid, siteId, partialPackageName)

	minimumSoftwareVersion, err := dal.GetMinimumRequiredSoftwareVersionForSite(siteId)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, updateTidError, http.StatusInternalServerError)
	}
	model.MinimumSoftwareVersion = minimumSoftwareVersion

	if model.Updates == nil {
		t := time.Now()
		newUpdate := models.TIDUpdateData{
			UpdateID:   0,
			PackageID:  model.Packages[0].PackageID,
			UpdateDate: common.FormatTime(t),
			IsTPA:      true,
		}
		model.Updates = append(model.Updates, &newUpdate)
	}
	dict := map[string]interface{}{
		"update": model.Updates,
		"model":  model,
		"show":   true,
	}
	renderTemplate(w, r, "updatesThirdPartyTID", dict, tmsUser)
}

func updatesThirdPartySelectHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	r.ParseForm()
	partialPackageName := r.Form.Get("partialPackageName")

	hasPermission, err := checkUserAcquirePermsByTid(tmsUser, r.Form.Get("TID"))
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, updateTidError, http.StatusInternalServerError)
		return
	}
	if !hasPermission {
		logging.Error(userUnauthorisedError)
		http.Error(w, userUnauthorisedError, http.StatusForbidden)
		return
	}

	thirdPartyApks, err := dal.GetThirdPartyApks(partialPackageName)
	if err != nil {
		logging.Error(err.Error())
	}

	ajaxResponse(w, thirdPartyApks)
}

func handleUpdatesThirdPartyApks(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	err := r.ParseForm()
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, filterVelocityLimitsError, http.StatusInternalServerError)
		return
	}

	tid, err := strconv.Atoi(r.Form.Get("TID"))
	if err != nil {
		handleError(w, errors.New(tidConversionFailed), tmsUser)
		return
	}
	siteId, err := strconv.Atoi(r.Form.Get("Site"))
	if err != nil {
		handleError(w, errors.New(siteConversionFailed), tmsUser)
		return
	}
	tidUpdateID, err := strconv.Atoi(r.Form.Get("tidUpdateID"))
	if err != nil {
		handleError(w, errors.New("an error occured while retrivin/converting the tidUpdateID"), tmsUser)
		return
	}

	packageIds := r.Form.Get("packageIds")

	currentUpdate, err := dal.GetTIDUpdate(tid, tidUpdateID)
	oldValue := "[" + strings.Join(getThirdPartyPackageNameByID("", currentUpdate.Options), " ") + "]"

	isValid, err := dal.UpdateTIDThirdPartyAPks(tid, tidUpdateID, packageIds)
	if err != nil {
		logging.Error(err)
	}

	if isValid {
		var profileExists bool
		var tidProfileId int64
		if profileExists, tidProfileId, err = dal.GetTidProfileIdForSiteId(tid, siteId); err != nil {
			logging.Error(err.Error())
		}
		newValue := "[" + strings.Join(getThirdPartyPackageNameByID(packageIds, nil), " ") + "]"
		if !profileExists && newValue != "" && newValue != "[]" {
			tidStr := dal.GetPaddedTidId(tid)
			profileTypeId, err := dal.GetProfileTypeId("tid", tidStr, 1, tmsUser.Username)
			if err != nil {
				logging.Error(err.Error())
			}
			err = dal.CreateTidoverrideAndSaveNewprofileChange(siteId, profileTypeId, tmsUser.Username, dal.ApproveCreate, "Override Created", tidStr, tid, 1)
			if err != nil {
				logging.Error(err.Error())
			}

			tidProfileId, err = dal.GetProfileIdFromTID(strings.TrimSpace(tidStr))
			if err != nil {
				logging.Error("An error occured during executing GetProfileIdFromTID", err.Error())
			}
			profileExists = true
		}
		if profileExists && oldValue != packageIds {
			SaveThirdPartyAudit(tidProfileId, oldValue, newValue, tid, tmsUser)
		}

	}

	ajaxResponse(w, isValid)
}

func getThirdPartyTarget(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	r.ParseForm()

	bindThirdPartyTarget := make([]dal.BindThirdPartyTarget, 0)
	tid, err := strconv.Atoi(r.Form.Get("TID"))
	if err != nil {
		handleError(w, errors.New(tidConversionFailed), tmsUser)
		return
	}

	tidUpdateId, err := strconv.Atoi(r.Form.Get("tidUpdateId"))
	if err != nil {
		handleError(w, errors.New("an error occured while retriving the tidUpdateID"), tmsUser)
		return
	}

	packageIDs := r.Form.Get("packageIDs")
	apkPackageIds := make([]int, 0)
	if packageIDs != "[]" && packageIDs != "" {
		var selectedPackageIDs []int
		err := json.Unmarshal([]byte(packageIDs), &selectedPackageIDs)
		if err != nil {
			logging.Error(err.Error())
		}
		apkPackageIds = append(apkPackageIds, selectedPackageIDs...)
	}

	hasPermission, err := checkUserAcquirePermsByTid(tmsUser, r.Form.Get("TID"))
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, generatePINError, http.StatusInternalServerError)
		return
	}
	if !hasPermission {
		logging.Error(userUnauthorisedError)
		http.Error(w, "User is not permitted to perform this action", http.StatusForbidden)
		return
	}

	if len(apkPackageIds) == 0 {
		update, err := dal.GetTIDUpdate(tid, tidUpdateId)
		if err != nil {
			logging.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if update.Options != nil {
			apkPackageIds = append(apkPackageIds, update.Options...)
		}
	}

	if len(apkPackageIds) > 0 {
		thirdPartyApks, err := dal.GetThirdPartyApksById(apkPackageIds)
		if err != nil {
			logging.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, tpApk := range thirdPartyApks {
			var tpTarget = dal.BindThirdPartyTarget{}
			tpTarget.ApkID = tpApk.ApkID
			tpTarget.ThirdPartyApkID = tpApk.Apk
			bindThirdPartyTarget = append(bindThirdPartyTarget, tpTarget)
		}
	}

	ajaxResponse(w, bindThirdPartyTarget)
}

func updatesSNHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	r.ParseForm()
	tid := r.Form.Get("TID")
	serial := r.Form.Get("SerialNo")
	siteId, err := strconv.Atoi(r.Form.Get("Site"))
	if err != nil {
		handleError(w, errors.New(siteConversionFailed), tmsUser)
		return
	}

	hasPermission, err := checkUserAcquirePermsByTid(tmsUser, r.Form.Get("TID"))
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, updateTidError, http.StatusInternalServerError)
		return
	}
	if !hasPermission {
		logging.Error(userUnauthorisedError)
		http.Error(w, userUnauthorisedError, http.StatusForbidden)
		return
	}

	model := BuildUpdatesTIDModel(w, tid, siteId, "")
	model.Serial = serial

	minimumSoftwareVersion, err := dal.GetMinimumRequiredSoftwareVersionForSite(siteId)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, updateTidError, http.StatusInternalServerError)
	}
	model.MinimumSoftwareVersion = minimumSoftwareVersion

	renderTemplate(w, r, "updateDetailsSN", model, tmsUser)
}

func AddTIDUpdateHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	r.ParseForm()
	tidID := r.Form.Get("TID")
	updateID, err := strconv.Atoi(r.Form.Get("UpdateID"))
	if err != nil {
		handleError(w, errors.New("an error occured while retriving the UpdateID/conversion to int"), tmsUser)
		return
	}
	serialNo := r.Form.Get("SerialNo")
	serial := r.Form.Get("Serial")
	siteId, err := strconv.Atoi(r.Form.Get("siteId"))
	if err != nil {
		handleError(w, errors.New(siteConversionFailed), tmsUser)
		return
	}

	model := BuildUpdatesTIDModel(w, tidID, siteId, "")
	model.Serial = serialNo
	var newUpdate models.TIDUpdateData
	newUpdate.UpdateID = updateID
	var t = time.Now()
	newUpdate.PackageID = model.Packages[0].PackageID
	newUpdate.UpdateDate = common.FormatTime(t)
	newUpdate.IsTPA = true
	dict := make(map[string]interface{})
	dict["update"] = newUpdate
	dict["model"] = model
	dict["show"] = true
	if serial == "SN" {
		renderPartialTemplate(w, r, "updateDetailsSNRow", dict, tmsUser)
	} else {
		renderPartialTemplate(w, r, "updateDetailsTIDRow", dict, tmsUser)
	}
}

type TPA struct {
	IsTPA bool
}

func ApplyTIDUpdatesHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	if err := r.ParseMultipartForm(10000); err != nil {
		logging.Error(err.Error())
		http.Error(w, applyTidUpdatesError, http.StatusInternalServerError)
		return
	}

	hasPermission, err := checkUserAcquirePermsByTid(tmsUser, r.Form.Get("TID"))
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, applyTidUpdatesError, http.StatusInternalServerError)
		return
	}
	if !hasPermission {
		logging.Error(userUnauthorisedError)
		http.Error(w, userUnauthorisedError, http.StatusForbidden)
		return
	}

	tidID, err := strconv.Atoi(r.Form.Get("TID"))
	if err != nil {
		handleError(w, errors.New(tidConversionFailed), tmsUser)
		return
	}
	siteId, err := strconv.Atoi(r.Form.Get("siteId"))
	if err != nil {
		handleError(w, errors.New(siteConversionFailed), tmsUser)
		return
	}
	tidUpdateID, err := strconv.Atoi(r.Form.Get("tidUpdateID"))
	if err != nil {
		handleError(w, errors.New("an error occured while retriving/converting the tidUpdateID"), tmsUser)
		return
	}

	profileExists, tidProfileId, err := dal.GetTidProfileIdForSiteId(tidID, siteId)
	if err != nil {
		logging.Error(err.Error())
	}
	currentUpdate, err := dal.GetTIDUpdate(tidID, tidUpdateID)
	if err != nil {
		logging.Error(err.Error())
	}
	oldV := getThirdPartyPackageNameByID("", currentUpdate.Options)
	oldValue := strings.Join(oldV, " ")
	oldValue = "[" + oldValue + "]"
	for _, update := range getUpdateFields(r.Form) {
		if update.PackageID == 0 {
			msg := "Please select the target software"
			logging.Error(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		newV, err := addAndUpdateTIDFlagAndThirdPartyPackage(update.UpdateID, tidID, update.PackageID, update.UpdateDate, update.ThirdPartyApkID, nil)
		if err != nil {
			logging.Error(err)
		}
		newValue := strings.Join(newV, " ")
		newValue = "[" + newValue + "]"
		if !profileExists && newValue != "" && newValue != "[]" {
			tidStr := dal.GetPaddedTidId(tidID)
			profileTypeId, err := dal.GetProfileTypeId("tid", tidStr, 1, tmsUser.Username)
			if err != nil {
				logging.Error(err.Error())
			}
			err = dal.CreateTidoverrideAndSaveNewprofileChange(siteId, profileTypeId, tmsUser.Username, dal.ApproveCreate, "Override Created", tidStr, tidID, 1)
			if err != nil {
				logging.Error(err.Error())
			}

			tidProfileId, err = dal.GetProfileIdFromTID(strings.TrimSpace(tidStr))
			if err != nil {
				logging.Error("An error occured during executing GetProfileIdFromTID", err.Error())
			}
			profileExists = true
		}
		if profileExists && newValue != "" {
			SaveThirdPartyAudit(tidProfileId, oldValue, newValue, tidID, tmsUser)
		}
	}
}

func addAndUpdateTIDFlagAndThirdPartyPackage(tidUpdateId int, tidId int, packageId int, updateDate string, thirdPartyApk string, thirdPartyApkArray []int) ([]string, error) {
	var selectedPackageIDs []int
	var thirdPartyPackageList []string
	if thirdPartyApk != "[]" && thirdPartyApk != "" {
		err := json.Unmarshal([]byte(thirdPartyApk), &selectedPackageIDs)
		if err != nil {
			logging.Error(err.Error())
		}
	}
	if selectedPackageIDs == nil {
		selectedPackageIDs = thirdPartyApkArray
	}
	thirdPartyApks, err := dal.AddUpdateTidAndGetThirdPartyApks(tidUpdateId, tidId, packageId, updateDate, thirdPartyApk, "")
	if err != nil {
		logging.Error(err.Error())
	}
	if selectedPackageIDs != nil {
		for _, id := range selectedPackageIDs {
			for _, tpApk := range thirdPartyApks {
				if id == tpApk.ApkID {
					thirdPartyPackageList = append(thirdPartyPackageList, tpApk.Apk)
				}
			}
		}
	}
	return thirdPartyPackageList, nil
}

func SaveThirdPartyAudit(tidProfileId int64, oldValue string, newValue string, tid int, user *entities.TMSUser) {
	err := dal.AddThirdPartyAuditHistory(tidProfileId, oldValue, newValue, tid, user)
	if err != nil {
		return
	}
	elementId, err := dal.GetElementIdFromGroupNameElementName("thirdParty", "thirdPartyPackageList")
	if err != nil {
		logging.Error(err.Error())
	}
	newValue = strings.Replace(newValue, " ", ",", -1)
	err = dal.NewProfileDataRepository().SetDataValueByElementIdAndProfileIdWithoutApproval(elementId, int(tidProfileId), newValue, *user)
	if err != nil {
		logging.Error(err.Error())
	}
}

func getThirdPartyPackageNameByID(thirdPartyApkString string, thirdPartyApkArray []int) []string {
	var selectedPackageIDs []int
	var thirdPartyPackageList []string

	if thirdPartyApkString != "[]" && thirdPartyApkString != "" {
		err := json.Unmarshal([]byte(thirdPartyApkString), &selectedPackageIDs)
		if err != nil {
			logging.Error(err.Error())
		}
	}
	if selectedPackageIDs == nil {
		selectedPackageIDs = thirdPartyApkArray
	}
	if len(selectedPackageIDs) > 0 {
		thirdPartyApks, err := dal.GetThirdPartyApksById(selectedPackageIDs)
		if err != nil {
			logging.Error(err.Error())
		}
		for _, tpApk := range thirdPartyApks {
			thirdPartyPackageList = append(thirdPartyPackageList, tpApk.Apk)
		}
	}
	return thirdPartyPackageList
}

func DeleteTIDUpdatesHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	r.ParseForm()

	hasPermission, err := checkUserAcquirePermsByTid(tmsUser, r.Form.Get("TID"))
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, deleteTidUpdatesError, http.StatusInternalServerError)
		return
	}
	if !hasPermission {
		logging.Error(userUnauthorisedError)
		http.Error(w, userUnauthorisedError, http.StatusForbidden)
		return
	}

	updateID, err := strconv.Atoi(r.Form.Get("UpdateID"))
	if err != nil {
		handleError(w, errors.New("an error occured while retreiving/converting the tidUpdateID"), tmsUser)
		return
	}
	tidID, err := strconv.Atoi(r.Form.Get("TID"))
	if err != nil {
		handleError(w, errors.New(tidConversionFailed), tmsUser)
		return
	}
	currentUpdate, err := dal.GetTIDUpdate(tidID, updateID)
	if err != nil {
		logging.Error(err.Error())
	}

	oldValue := "[" + strings.Join(getThirdPartyPackageNameByID("", currentUpdate.Options), " ") + "]"

	err = dal.DeleteUpdateFromTID(updateID, tidID)
	if err != nil {
		handleError(w, errors.New("an error occured while executing DeleteUpdateFromTID "+err.Error()), tmsUser)
		return
	}

	newUpdate, err := dal.GetTIDUpdate(tidID, updateID)
	if err != nil {
		logging.Error(err.Error())
	}
	newValue := "[" + strings.Join(getThirdPartyPackageNameByID("", newUpdate.Options), " ") + "]"

	tidProfileIdStr := r.Form.Get("ProfileId")
	if tidProfileIdStr != "" {
		tidProfileID, err := strconv.Atoi(tidProfileIdStr)
		if err != nil {
			handleError(w, errors.New("an error occured while retreiving/converting the tidProfileId"), tmsUser)
			return
		}
		err = dal.UpdateOrSetThirdPartyPackageList(tidProfileID, oldValue, newValue, tmsUser, tidID)
		if err != nil {
			handleError(w, errors.New("an error occured while executing UpdateOrSetThirdPartyPackageList"+err.Error()), tmsUser)
			return
		}
	}

}

func saveTidProfileHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	err := r.ParseMultipartForm(10000)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	siteId, err := strconv.Atoi(r.Form.Get("siteID"))
	if err != nil {
		handleError(w, errors.New(siteConversionFailed), tmsUser)
		return
	}
	tidId, err := strconv.Atoi(r.Form.Get("tidID"))
	if err != nil {
		handleError(w, errors.New(tidConversionFailed), tmsUser)
		return
	}
	tidString := dal.GetPaddedTidId(tidId)

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
	hasPermission, err := checkUserAcquirePermsByTid(tmsUser, tidString)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "unable to check user permission", http.StatusInternalServerError)
		return
	}
	if !hasPermission {
		logging.Error(userUnauthorisedError)
		http.Error(w, userUnauthorisedError, http.StatusForbidden)
		return
	}

	profileExists, profileID, err := dal.GetTidProfileIdForSiteId(tidId, siteId)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "unable to find TID", http.StatusInternalServerError)
		return
	}

	elements := getElementFields(r.Form, false)
	validationMessages, _ := validateDataElements(dal.NewValidationDal(), elements, int(profileID), "tid", false, "")
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

	if !profileExists {
		//Workaround for query which does not support inside stored procedure
		profileTypeId, err := dal.GetProfileTypeId("tid", tidString, 1, tmsUser.Username)
		if err != nil {
			handleError(w, errors.New("unable to save TID profile while executing SaveNewProfile"), tmsUser)
			return
		}
		err = dal.CreateTidoverrideAndSaveNewprofileChange(siteId, profileTypeId, tmsUser.Username, dal.ApproveCreate, "Override Created", tidString, tidId, 1)
		if err != nil {
			logging.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}

	saveDataElements(int(profileID), siteId, elements, 0, tmsUser, false, 0)
	id, err := dal.GetProfileFromSite(siteId)
	if err != nil {
		handleError(w, errors.New("unable to get profile id while executing GetProfileFromSite"), tmsUser)
		return
	}
	p := buildSiteProfileMaintenanceModel(w, id, tmsUser, pageSize, pageNumber, tidSearchTerm, siteId)
	renderTemplate(w, r, "profileMaintenanceTIDs", p, tmsUser)
}

func deleteTidProfileHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	err := r.ParseMultipartForm(10000)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, deletingTidProfileError, http.StatusInternalServerError)
		return
	}

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

	hasPermission, err := checkUserAcquirePermsByTid(tmsUser, r.Form.Get("tidID"))
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, deletingTidProfileError, http.StatusInternalServerError)
		return
	}
	if !hasPermission {
		logging.Error(userUnauthorisedError)
		http.Error(w, userUnauthorisedError, http.StatusForbidden)
		return
	}

	tidId, err := strconv.Atoi(r.Form.Get("tidID"))
	if err != nil {
		handleError(w, errors.New(tidConversionFailed), tmsUser)
		return
	}
	siteId, err := strconv.Atoi(r.Form.Get("siteID"))
	if err != nil {
		handleError(w, errors.New(siteConversionFailed), tmsUser)
		return
	}
	profileExists, profileID, err := dal.GetTidProfileIdForSiteId(tidId, siteId)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, deletingTidProfileError, http.StatusInternalServerError)
		return
	}

	if profileExists {
		tidString := dal.GetPaddedTidId(tidId)
		err = dal.SaveTidProfileChange(int(profileID), tidString, "Override Deleted", tmsUser.Username, dal.ApproveDelete, 0)
		if err != nil {
			handleError(w, errors.New("an error occured while executing SaveTidProfileChange:"+err.Error()), tmsUser)
			return
		}
	}

	id, err := dal.GetProfileFromSite(siteId)
	if err != nil {
		handleError(w, errors.New("an error occured while executing GetProfileFromSite:"+err.Error()), tmsUser)
		return
	}
	p := buildSiteProfileMaintenanceModel(w, id, tmsUser, pageSize, pageNumber, tidSearchTerm, siteId)
	renderTemplate(w, r, "profileMaintenanceTIDs", p, tmsUser)
}

func getTIDDetailsHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	r.ParseForm()

	tid := r.Form.Get("TID")

	hasPermission, err := checkUserAcquirePermsByTid(tmsUser, tid)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, retrievingTidDetailsError, http.StatusInternalServerError)
		return
	}
	if !hasPermission {
		logging.Error(userUnauthorisedError)
		http.Error(w, userUnauthorisedError, http.StatusForbidden)
		return
	}

	details, err := dal.GetTidDetails(tid)
	if err != nil {
		handleError(w, errors.New(retrievingTidDetailsError), tmsUser)
	}

	renderPartialTemplate(w, r, "TidDetailsModal", details, tmsUser)
}

func addNewDuplicatedTidOverride(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	r.ParseForm()
	tidIndex, err := strconv.Atoi(r.Form.Get("TidIndex"))
	if err != nil {
		handleError(w, errors.New("an error occured during parsing TidIndex to int"), tmsUser)
		return
	}

	siteId, err := strconv.Atoi(r.Form.Get("Site"))
	if err != nil {
		handleError(w, errors.New(siteConversionFailed), tmsUser)
		return
	}

	profileId, err := strconv.Atoi(r.Form.Get("SiteProfileId"))
	if err != nil {
		handleError(w, errors.New(profileConversionFailed), tmsUser)
		return
	}

	hasPermission, err := checkUserAcquirePermsBySite(tmsUser, siteId)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, duplicateTidOverrideError, http.StatusInternalServerError)
		return
	}
	if !hasPermission {
		logging.Error(userUnauthorisedError)
		http.Error(w, userUnauthorisedError, http.StatusForbidden)
		return
	}

	elements := getElementFields(r.Form, true)
	dataElementDetails, err := dal.GetDataAllElementID()
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, saveProfileError, http.StatusInternalServerError)
		return
	}

	// We used to validate the data elements at this point, but that is unnecessary; validation will
	// occur when the duplicated TID details are attempted to be saved.
	// Moreover, attempting to validate here would cause unique constraint checks to fail because
	// we would have duplicated the values from the TID that was just saved.

	// Use the regular build New tid model to create an add tid with defaults set to the values for the current site
	model := buildAddNewTidModel(w, profileId, tidIndex, siteId)
	model.IsDuplicate = true
	model.DuplicatedFrom, _ = strconv.Atoi(r.Form.Get("tidID"))

	// Replace the default site values with those passed from the previous tid to be duplicated from
	for _, group := range model.DefaultTidGroups {
		for j := range group.DataElements {
			ele := &group.DataElements[j]
			for e, val := range elements {
				if ele.ElementId == e {
					ele.DataValue = val
					ele.Options, ele.OptionSelectable = dal.BuildOptionsData(
						dal.GetElementOptionsForId(group.DataGroup, ele.Name, dataElementDetails), val, group.DataGroup, ele.Name, profileId)
				}
			}
		}
	}

	renderTemplate(w, r, "addTid", model, tmsUser)
}

func getUpdateFields(keyPairs url.Values) []models.TIDUpdateData {
	updateElements := make([]models.TIDUpdateData, 0)

	dates := make(map[int]string)
	versions := make(map[int]string)
	thirdParty := make(map[int]string)

	for key, val := range keyPairs {
		if strings.HasPrefix(key, "date_") {
			var groupKey, _ = strconv.Atoi(string(key[5:]))
			dates[groupKey] = val[0]
		}

		if strings.HasPrefix(key, "version_") {
			var groupKey, _ = strconv.Atoi(string(key[8:]))
			versions[groupKey] = val[0]
		}

		if strings.HasPrefix(key, "thirdParty_") {
			var groupKey, _ = strconv.Atoi(string(key[11:]))
			thirdParty[groupKey] = val[0]
		}
	}

	if len(dates) == len(versions) {
		for key := range dates {
			t := models.TIDUpdateData{}
			t.UpdateID = key
			if thirdParty[key] == "" {
				thirdParty[key] = "[]"
			}
			t.ThirdPartyApkID, _ = thirdParty[key]
			t.UpdateDate = dates[key]
			t.PackageID, _ = strconv.Atoi(versions[key])
			updateElements = append(updateElements, t)
		}
	}
	return updateElements
}

func handleUpdateSerialNumberValidation(updates []models.TIDUpdateData, serialNumber string) (bool, string) {

	valid, err := validation.New(dal.NewValidationDal()).ValidateSerialNumber(serialNumber)
	if err != nil {
		msg := "Error thrown when attempting to validate Serial Number"
		return false, msg
	}

	if !valid {
		msg := "Serial number must consist of up to 10 characters"
		return false, msg
	}

	if err := dal.CheckSerialInUse(serialNumber); err != nil {
		return false, err.Error()
	}

	for _, updateModel := range updates {
		if updateModel.PackageID == 0 {
			msg := "Please select the target software"
			return false, msg
		}
	}
	return true, serialNumber
}
