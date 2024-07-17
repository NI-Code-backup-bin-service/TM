package services

import (
	"nextgen-tms-website/dal"
	"nextgen-tms-website/entities"
	"nextgen-tms-website/logger"
	"strconv"
)

func SaveTerminalPaymentServiceConfig(tid, json string, tmsUser *entities.TMSUser) bool {
	db, err := dal.GetDB()
	if err != nil {
		return false
	}

	if len(tid) == 0 {
		return false
	}

	tidInt, err := strconv.Atoi(tid)
	if err != nil {
		return false
	}

	siteId, err := dal.GetSiteIDForTid(tidInt)
	if err != nil {
		logger.GetLogger().Information(err.Error())
		return false
	}

	profileExists, _, err := dal.GetTidProfileIdForSiteId(tidInt, siteId)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return false
	}

	if !profileExists {
		profileTypeId, err := dal.GetProfileTypeId("tid", tid, 1, tmsUser.Username)
		if err != nil {
			logger.GetLogger().Error(err.Error())
			return false
		}
		err = dal.CreateTidoverrideAndSaveNewprofileChange(siteId, profileTypeId, tmsUser.Username, dal.ApproveCreate, "Override Created", tid, tidInt, 1)
		if err != nil {
			logger.GetLogger().Error(err.Error())
			return false
		}
	}



	result, err := db.Exec("Call upsert_bulk_payment_services(?,?)", tid, json)
	if err != nil {
		return false
	}
	affected, err := result.RowsAffected()
	if err != nil || affected == 0 {
		return false
	}
	return true
}
