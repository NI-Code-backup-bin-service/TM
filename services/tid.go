package services

import (
	"errors"
	"fmt"
	"nextgen-tms-website/logger"
	"nextgen-tms-website/validation"
	"strings"

	"math/rand"
	"nextgen-tms-website/dal"
	"nextgen-tms-website/entities"
	"nextgen-tms-website/models"
	"strconv"
	"time"
)

func GetEODAutoTime(siteId, tid int) (bool, string, error) {
	auto := false
	autoTime := ""
	eodAutoDatas, err := dal.GetEODAutoData(siteId)
	if err != nil {
		return false, "", err
	}

	for _, eodAutoData := range eodAutoDatas {
		switch eodAutoData.Name {
		case "auto":
			auto, err = strconv.ParseBool(eodAutoData.Datavalue)
			if err != nil {
				return false, "", err
			}
		case "time":
			timeRange := strings.Split(eodAutoData.Datavalue, " | ")
			if len(timeRange) == 1 {
				autoTime = eodAutoData.Datavalue
			} else {
				startTime, err := time.Parse("15:04", timeRange[0])
				if err != nil {
					return false, "", err
				}

				endTime, err := time.Parse("15:04", timeRange[1])
				if err != nil {
					return false, "", err
				}

				timeDiff := endTime.Sub(startTime).Minutes()
				if timeDiff < 0 {
					endTime = endTime.AddDate(0, 0, 1)
				}
				timeDiff = endTime.Sub(startTime).Minutes()
				autoTime = startTime.Add(time.Duration((tid%1000)%int(timeDiff)) * time.Minute).Format("15:04")
			}
		}
	}

	return auto, autoTime, nil
}

func ValidateTidCreation(tidString, serialNumber string) error {
	var valid = true

	valid, err := validation.New(dal.NewValidationDal()).ValidateTid(tidString)
	if err != nil {
		logger.GetLogger().Information(err.Error())
		return err
	}

	if !valid {
		errMegs := "TID must be 8 digits long and cannot be all zeros"
		logger.GetLogger().Information(errMegs)
		return errors.New(errMegs)
	}

	valid, err = validation.New(dal.NewValidationDal()).ValidateSerialNumber(serialNumber)
	if err != nil {
		logger.GetLogger().Information(err.Error())
		return err
	}
	if !valid {
		errMegs := "Serial number must consist of up to 10 characters"
		logger.GetLogger().Information(errMegs)
		return errors.New(errMegs)
	}

	err = dal.CheckSerialInUse(serialNumber)
	if err != nil {
		logger.GetLogger().Information(err.Error())
		return err
	}

	return nil
}

func GetSiteIDFromMerchantID(MerchantId string) (int, error) {
	var errMegs string
	siteIdStr, err := dal.GetSiteIDFromMerchantID(MerchantId)
	if err != nil || siteIdStr == "" {
		errMegs = fmt.Sprintf("cannot find the site id for the mid %s:", MerchantId)
		logger.GetLogger().Information(errMegs)
		return 0, errors.New(errMegs)
	}
	siteId, err := strconv.Atoi(siteIdStr)
	if err != nil {
		errMegs = fmt.Sprintf("cannot convert the profile id (%s)", siteIdStr)
		logger.GetLogger().Information(errMegs)
		return 0, errors.New(errMegs)
	}

	return siteId, nil
}

func GetSiteProfileFromTid(tidString string) (int, error) {
	tid, err := strconv.Atoi(tidString)
	if err != nil {
		logger.GetLogger().Information(err.Error())
		return 0, err
	}

	siteId, err := dal.GetSiteIDForTid(tid)
	if err != nil {
		logger.GetLogger().Information(err.Error())
		return 0, err
	}

	profileId, err := GetProfileFromSite(siteId)
	if err != nil {
		logger.GetLogger().Information(err.Error())
		return 0, err
	}

	return profileId, nil
}

func GetProfileFromSite(siteId int) (int, error) {
	profileId, err := dal.GetProfileFromSite(siteId)
	if err != nil {
		errMegs := fmt.Sprintf("cannot find the profile id for the site Id (%d)", siteId)
		logger.GetLogger().Information(errMegs)
		return 0, errors.New(errMegs)
	}
	return profileId, nil
}

func AddTidOverrideDataGroups(profileId, siteId int, tmsUser *entities.TMSUser) error {
	activeGroups, err := dal.GetTidOveridableActiveDataGroupIds(siteId)
	if err != nil {
		return err
	}

	if !AddDataGroupsToProfile(int(profileId), activeGroups, tmsUser) {
		return errors.New("error adding profile groups")
	}
	return nil
}

func AddDataGroupsToProfile(profileId int, dataGroups []string, user *entities.TMSUser) bool {
	for _, dataGroup := range dataGroups {
		dal.AddDataGroupToProfile(profileId, dataGroup, user.Username)
	}
	return true
}

func ValidateProfileIdForTidSiteId(tid int, siteId int) error {
	exists, err := dal.CheckTidExistsSiteId(tid, siteId)
	if err != nil {
		errMegs := fmt.Sprintf("cannot find the tid for the site Id (%d) and tid(%d)", siteId, tid)
		logger.GetLogger().Information(errMegs)
		return err
	}
	if exists {
		return nil
	}
	errMsg := fmt.Sprintf("cannot find the tid with parameter's TID: %d, SiteID: %d", tid, siteId)
	return errors.New(errMsg)
}

func AddTid(tidString, serialId, userName string, siteId, profileId int) error {
	tid, err := strconv.Atoi(tidString)
	if err != nil {
		logger.GetLogger().Information(errors.New("invalid TID entry"))
		return err
	}
	auto, autoTime, err := GetEODAutoTime(siteId, tid)
	if err != nil {
		logger.GetLogger().Information(err.Error())
		return err
	}
	if err := dal.AddTidToSiteAndSaveTidProfileChange(tidString, serialId, siteId, auto, autoTime, profileId, userName, dal.ApproveCreate, 1); err != nil {
		logger.GetLogger().Information(err.Error())
		return err
	}
	return nil
}

func GetTIDUpdateFields(hardware models.UpdateTIDRequest) (string, string, *models.TIDUpdateData, error) {

	//Check if TID exists
	tid, err := strconv.Atoi(hardware.TID)
	if err != nil {
		logger.GetLogger().Error("Error during string conversion", err.Error())
		return "", "", nil, err
	}

	tidExists, _, _ := dal.CheckThatTidExists(tid)
	if !tidExists {
		logger.GetLogger().Error("TID does not exist")
		return "", "", nil, errors.New("TID does not exist")
	}

	//oldSerialnumber
	oldSerialNumber, err := dal.GetSerialByTid(hardware.TID)
	if err != nil || oldSerialNumber == "" {
		logger.GetLogger().Error("Serial number does not exist for TID", err.Error())
		return "", "", nil, err
	}

	update, err := dal.GetTIDUpdateFields(tid, 0, hardware.ApkVersion, hardware.TPApkVersion)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return "", "", nil, err
	}

	profileId, err := GetSiteProfileFromTid(hardware.TID)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return "", "", nil, err
	}

	return oldSerialNumber, strconv.Itoa(profileId), update, nil
}

func AddAPK(tid, siteId int, apkVersion string, tpApkVersion []string) error {
	update, err := dal.GetTIDUpdateFields(tid, siteId, apkVersion, tpApkVersion)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return err
	}

	err = dal.AddUpdateToTID(update.UpdateID, tid, update.PackageID, update.UpdateDate, update.ThirdPartyApkID)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return err
	}

	return nil
}

func UpdateTIDHandler(tmsUser *entities.TMSUser, oldSn string, newSn string, siteProfileId string, updates []models.TIDUpdateData, tid string) error {
	approvalID, err := dal.UpdateSerialNumberPendingApproval(tmsUser, siteProfileId, tid, oldSn, newSn)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return err
	}

	err = dal.ApproveChange(approvalID, tmsUser.Username, "TID")
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return err
	}

	var tidID, _ = strconv.Atoi(tid)
	err = dal.AddBulkUpdateToTID(updates, tidID)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return err
	}

	err = dal.UpdateTIDFlag(tidID)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return err
	}

	return nil
}

func GenerateOTP(tid string, intent dal.OTPIntent, tmsUser *entities.TMSUser) (*models.GenerateOTPResult, error) {
	pin := createOTP()

	time, err := dal.StoreOneTimePIN(tid, pin, intent, tmsUser)
	if err != nil {
		return nil, err
	}

	var data models.GenerateOTPResult
	data.PIN = pin
	data.ExpiryTime = time

	return &data, nil
}

func ChangeHistory(profileId int, tidString, changeMessage, userName string) error {
	err := dal.SaveTidProfileChange(profileId, tidString, changeMessage, userName, dal.ApproveCreate, 1)
	if err != nil {
		logger.GetLogger().Information(err.Error())
		return errors.New(err.Error())
	}

	return nil
}

func createOTP() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	PIN := r.Intn(99999)
	PINstr := strconv.Itoa(PIN)
	leads := 5 - len(PINstr)
	for leads > 0 {
		PINstr = "0" + PINstr
		leads--
	}

	return PINstr
}
