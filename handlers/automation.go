package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"nextgen-tms-website/dal"
	"nextgen-tms-website/entities"
	"nextgen-tms-website/logger"
	"nextgen-tms-website/models"
	"nextgen-tms-website/services"
	"strconv"

	"github.com/gorilla/context"
)

func CreateTID(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*entities.TMSUser)

	request := models.CreateTIDRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		respondJSON(models.Error{Code: http.StatusBadRequest, Message: err.Error()}, http.StatusBadRequest, w)
		return
	}
	defer r.Body.Close()
	logger.GetLogger().Information(fmt.Sprintf("Request received for CreateTID with parameter's TID: %s,MID: %s,Serial Number: %s,APK version: %s, Third party APK versions: %v", request.TID, request.SerialNumber, request.MID, request.ApkVersion, request.TPApkVersion))

	err = validate.Struct(request)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		respondJSON(models.Error{Code: http.StatusBadRequest, Message: err.Error()}, http.StatusBadRequest, w)
		return
	}

	siteId, err := services.GetSiteIDFromMerchantID(request.MID)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		respondJSON(models.Error{Code: http.StatusNotFound, Message: err.Error()}, http.StatusNotFound, w)
		return
	}

	profileId, err := services.GetProfileFromSite(siteId)
	if err != nil {
		errMegs := fmt.Sprintf("cannot find the profile id for the site Id (%d)", siteId)
		logger.GetLogger().Error(errMegs)
		respondJSON(models.Error{Code: http.StatusNotFound, Message: errMegs}, http.StatusNotFound, w)
		return
	}

	err = services.ValidateTidCreation(request.TID, request.SerialNumber)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		respondJSON(models.Error{Code: http.StatusBadRequest, Message: err.Error()}, http.StatusBadRequest, w)
		return
	}

	tid, err := strconv.Atoi(request.TID)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		respondJSON(models.Error{Code: http.StatusBadRequest, Message: err.Error()}, http.StatusBadRequest, w)
		return
	}

	err = services.AddAPK(tid, siteId, request.ApkVersion, request.TPApkVersion)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		respondJSON(models.Error{Code: http.StatusInternalServerError, Message: err.Error()}, http.StatusInternalServerError, w)
		return
	}

	err = services.AddTid(request.TID, request.SerialNumber, user.Username, siteId, profileId)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		respondJSON(models.Error{Code: http.StatusInternalServerError, Message: err.Error()}, http.StatusInternalServerError, w)
		return
	}

	result, err := services.GenerateOTP(request.TID, dal.OtpIntentEnrolment, user)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		respondJSON(models.Error{Code: http.StatusInternalServerError, Message: err.Error()}, http.StatusInternalServerError, w)
		return
	}

	response := models.GenerateOTPResponse{
		PIN:        result.PIN,
		ExpiryTime: result.ExpiryTime,
	}

	logger.GetLogger().Information(fmt.Sprintf("Successfully Created TID with number: %s and generated OTP.", request.TID))
	respondJSON(response, http.StatusCreated, w)
}

func UpdateTID(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*entities.TMSUser)
	var siteId int
	var tid int

	request := models.UpdateTIDRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		respondJSON(models.Error{Code: http.StatusInternalServerError, Message: err.Error()}, http.StatusInternalServerError, w)
		return
	}
	defer r.Body.Close()
	logger.GetLogger().Information(fmt.Sprintf("Request received for UpdateTID with parameter's TID: %s, MID: %s, Serial Number: %s,APK version: %s, Third party APK versions: %v", request.TID, request.MID, request.SerialNumber, request.ApkVersion, request.TPApkVersion))

	err = validate.Struct(request)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		respondJSON(models.Error{Code: http.StatusBadRequest, Message: err.Error()}, http.StatusBadRequest, w)
		return
	}

	// NEX-11786: Validate MID and TID //
	if siteId, err = services.GetSiteIDFromMerchantID(request.MID); err != nil {
		logger.GetLogger().Error(err.Error())
		respondJSON(models.Error{Code: http.StatusNotFound, Message: err.Error()}, http.StatusNotFound, w)
		return
	}
	if tid, err = strconv.Atoi(request.TID); err != nil {
		logger.GetLogger().Error(err.Error())
		respondJSON(models.Error{Code: http.StatusBadRequest, Message: err.Error()}, http.StatusBadRequest, w)
		return
	}
	err = services.ValidateProfileIdForTidSiteId(tid, siteId)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		respondJSON(models.Error{Code: http.StatusBadRequest, Message: err.Error()}, http.StatusBadRequest, w)
		return
	}
	// End  //

	OldSerialNumber, siteProfileID, update, err := services.GetTIDUpdateFields(request)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		respondJSON(models.Error{Code: http.StatusInternalServerError, Message: err.Error()}, http.StatusInternalServerError, w)
		return
	}

	updates := []models.TIDUpdateData{*update}
	err = services.UpdateTIDHandler(user, OldSerialNumber, request.SerialNumber, siteProfileID, updates, request.TID)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		respondJSON(models.Error{Code: http.StatusInternalServerError, Message: err.Error()}, http.StatusInternalServerError, w)
		return
	}

	result, err := services.GenerateOTP(request.TID, dal.OtpIntentEnrolment, user)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		respondJSON(models.Error{Code: http.StatusInternalServerError, Message: err.Error()}, http.StatusInternalServerError, w)
		return
	}

	response := models.GenerateOTPResponse{
		PIN:        result.PIN,
		ExpiryTime: result.ExpiryTime,
	}
	logger.GetLogger().Information(fmt.Sprintf("Successfully Updated TID: %s", request.TID))
	respondJSON(response, http.StatusOK, w)
}

func GenerateOTP(w http.ResponseWriter, r *http.Request) {
	tmsUser := context.Get(r, "user").(*entities.TMSUser)
	var siteId int
	var request models.GenerateOTPRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		respondJSON(models.Error{Code: http.StatusBadRequest, Message: err.Error()}, http.StatusBadRequest, w)
		return
	}
	logger.GetLogger().Information(fmt.Sprintf("Request received for GenerateOTP with TID: %s", request.TID))
	defer r.Body.Close()
	err = validate.Struct(request)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		respondJSON(models.Error{Code: http.StatusBadRequest, Message: err.Error()}, http.StatusBadRequest, w)
		return
	}

	tidInt, err := strconv.Atoi(request.TID)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		respondJSON(models.Error{Code: http.StatusInternalServerError, Message: err.Error()}, http.StatusInternalServerError, w)
		return
	}

	if tidExists, _, _ := dal.CheckThatTidExists(tidInt); !tidExists {
		errStr := "TID does not exist"
		logger.GetLogger().Error(errStr)
		respondJSON(models.Error{Code: http.StatusNotFound, Message: errStr}, http.StatusNotFound, w)
		return
	}

	// NEX-11786: Validate MID and TID //
	if siteId, err = services.GetSiteIDFromMerchantID(request.MID); err != nil {
		logger.GetLogger().Error(err.Error())
		respondJSON(models.Error{Code: http.StatusNotFound, Message: err.Error()}, http.StatusNotFound, w)
		return
	}
	err = services.ValidateProfileIdForTidSiteId(tidInt, siteId)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		respondJSON(models.Error{Code: http.StatusBadRequest, Message: err.Error()}, http.StatusBadRequest, w)
		return
	}
	// End  //

	result, err := services.GenerateOTP(request.TID, dal.OtpIntentEnrolment, tmsUser)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		respondJSON(models.Error{Code: http.StatusInternalServerError, Message: err.Error()}, http.StatusInternalServerError, w)
		return
	}

	profileId, err := services.GetSiteProfileFromTid(request.TID)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		respondJSON(models.Error{Code: http.StatusInternalServerError, Message: err.Error()}, http.StatusInternalServerError, w)
		return
	}

	err = services.ChangeHistory(profileId, request.TID, "Enrolment PIN Generated", tmsUser.Username)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		respondJSON(models.Error{Code: http.StatusInternalServerError, Message: err.Error()}, http.StatusInternalServerError, w)
		return
	}

	response := models.GenerateOTPResponse{
		PIN:        result.PIN,
		ExpiryTime: result.ExpiryTime,
	}
	logger.GetLogger().Information(fmt.Sprintf("Successfully Generated OTP for TID: %s", request.TID))
	respondJSON(response, http.StatusCreated, w)
}
