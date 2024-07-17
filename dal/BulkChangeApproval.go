package dal

import (
	"bytes"
	"errors"
	"fmt"
	"nextgen-tms-website/common"
	"nextgen-tms-website/crypt"
	"nextgen-tms-website/models"
	"strconv"
	"strings"
	"time"
)

const ASCIINULL = "\x00"

func InsertBulkApproval(fileName, fileType, currentUser string, changeType int, approved int) error {
	db, err := GetDB()
	if err != nil {
		logging.Error("Unable to get the database instance : " + err.Error())
		return err
	}

	_, err = db.Exec("Call add_to_bulk_approvals(?, ?, ?, ?, ?)", fileName, fileType, currentUser, changeType, approved)
	if err != nil {
		logging.Error("Error while inserting bulk_approvals : " + err.Error())
		return err
	}

	return nil
}

func GetAllApprovedAndRejectedBulkApprovals() ([]models.BulkApproval, error) {
	db, err := GetDB()
	if err != nil {
		logging.Error("Unable to get the database instance : " + err.Error())
		return nil, err
	}

	bulkApprovals := []models.BulkApproval{}
	rows, err := db.Query("SELECT filename, filetype, change_type, created_by, approved_by, created_at, approved_at, approved FROM bulk_approvals WHERE approved!=0 ORDER BY approved_at DESC")
	if err != nil {
		logging.Error("Error while getting bulk_tid_flagging : " + err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ba models.BulkApproval
		err = rows.Scan(&ba.Filename, &ba.FileType, &ba.ChangeType, &ba.CreatedBy, &ba.ApprovedBy, &ba.CreatedAt, &ba.ApprovedAt, &ba.Approved)
		if err != nil {
			logging.Error("Error while scanning bulk_tid_flagging : " + err.Error())
			return nil, err
		}
		//To convert String ChangeType value into int to get ChangeType string value.
		intChangeType, _ := strconv.Atoi(ba.ChangeType)
		ba.ChangeType = common.GetChangeType()(intChangeType)

		bulkApprovals = append(bulkApprovals, ba)
	}

	return bulkApprovals, nil
}

func GetAllUnapprovedBulkApprovals() ([]models.BulkApproval, error) {
	db, err := GetDB()
	if err != nil {
		logging.Error("Unable to get the database instance : " + err.Error())
		return nil, err
	}

	unapprovedBulkApprovals := []models.BulkApproval{}
	rows, err := db.Query("SELECT filename, filetype, change_type, created_by, created_at FROM bulk_approvals WHERE approved=0 ORDER BY created_at DESC")
	if err != nil {
		logging.Error("Error while getting bulk_tid_flagging : " + err.Error())
		return nil, err
	}
	defer rows.Close()

	var fileName, createdBy, createdAt, fileType string
	var ChangeType int
	for rows.Next() {
		err = rows.Scan(&fileName, &fileType, &ChangeType, &createdBy, &createdAt)
		if err != nil {
			logging.Error("Error while scanning bulk_tid_flagging : " + err.Error())
			return nil, err
		}

		BulkApproval := models.BulkApproval{
			Filename:   fileName,
			ChangeType: common.GetChangeType()(ChangeType),
			FileType:   fileType,
			CreatedBy:  createdBy,
			CreatedAt:  createdAt,
		}

		unapprovedBulkApprovals = append(unapprovedBulkApprovals, BulkApproval)
	}

	return unapprovedBulkApprovals, nil
}

func ApproveBulkApprovalTerminalFlagging(fileName, currentUser string, records [][]string) error {
	db, err := GetDB()
	if err != nil {
		logging.Error("Unable to get the database instance : " + err.Error())
		return err
	}

	i := 0
	for j, r := range records {
		if j == 0 {
			// Appending header
			records[j] = append(records[j], "Error")
			continue
		}

		if strings.TrimSpace(r[0]) == "" || strings.TrimSpace(r[1]) == "" || strings.TrimSpace(r[2]) == "" {
			logging.Error("TID / Serial Number / APK Version is missing")
			records[j] = append(records[j], "TID / Serial Number / APK Version is missing")
			continue
		}

		tid, err := strconv.Atoi(strings.TrimSpace(r[0]))
		if err != nil {
			logging.Error("Invalid tid : " + r[0] + " ; " + err.Error())
			records[j] = append(records[j], "Invalid tid : "+err.Error())
			continue
		}

		var packageID int
		err = db.QueryRow("SELECT package_id FROM package WHERE version=?", strings.TrimSpace(r[2])).Scan(&packageID)
		if err != nil {
			logging.Error("Unable to find the Target APK version : " + r[2] + " ; " + err.Error())
			records[j] = append(records[j], "Unable to find the Target APK version : "+err.Error())
			continue
		}

		if packageID == 0 {
			logging.Error("Unable to find the Target APK version : " + r[2])
			records[j] = append(records[j], "Unable to find the Target APK version")
			continue
		}

		var apkID string
		//some null values are getting appended last row; inorder to remove that bytes.Trim is used
		tpApk := string(bytes.Trim([]byte(strings.TrimSpace(r[3])), ASCIINULL))
		if tpApk != "" {
			tpApks := strings.Split(tpApk, "|")
			apkID, err = GetThirdPartyApk(tpApks, tid, 0)
			if err != nil {
				logging.Error("Unable to find the TP APK : " + r[3] + " ; " + err.Error())
				records[j] = append(records[j], "Unable to find the TP APK : "+err.Error())
				continue
			}
		}

		result, err := db.Exec("UPDATE tid SET flag_status=true, flagged_date=CURRENT_TIMESTAMP WHERE tid_id=? AND serial=?", tid, strings.TrimSpace(r[1]))
		if err != nil {
			logging.Error("Unable to update flag status : " + err.Error())
			records[j] = append(records[j], "Unable to update flag status : "+err.Error())
			continue
		}

		n, _ := result.RowsAffected()
		if n < 1 {
			logging.Error("Unable to find the record matching tid : "+r[0]+"and serial number : "+r[1]+" ; RowsAffected : ", n)
			records[j] = append(records[j], "Unable to find the record matching tid and serial number")
			continue
		}

		_, err = db.Exec("CALL insert_tid_update_ignore_TPAPK(?, ?, ?, ?)", tid, packageID, time.Now(), apkID)
		if err != nil {
			logging.Error("Unable to insert tid update : " + err.Error())
			records[j] = append(records[j], "Unable to insert tid update : "+err.Error())
			continue
		}
		i++
	}
	logging.Debug(fmt.Sprintf("Flag status updated for %d tids", i))

	_, err = db.Exec("UPDATE bulk_approvals SET approved_by=?, approved_at=NOW(), approved=1 WHERE filename=?", currentUser, fileName)
	if err != nil {
		logging.Error("Error while updating bulk_tid_flagging : " + err.Error())
		return err
	}

	return nil
}

func ApproveBulkApprovalBulkSiteUpdate(fileName, currentUser string, records [][]string) error {

	db, err := GetDB()
	if err != nil {
		logging.Error("Unable to get the database instance : " + err.Error())
		return err
	}

	header := records[0]
	for _, rows := range records[1:] {
		profileIDstr, err := GetProfileIdFromMID(rows[0]) //profileID is string here
		if err != nil {
			logging.Error("Unable to get the profile id from TID : " + err.Error())
			return err
		}

		profileID, err := strconv.Atoi(profileIDstr)
		if err != nil {
			logging.Error("Error converting profileID to integer" + err.Error())
			return err
		}

		for column, value := range rows {
			if column == 0 || value == "" {
				continue
			}

			value = strings.TrimSpace(value)
			dataGroupElement := strings.Split(strings.TrimSpace(header[column]), ".")
			dataElementId, err := GetDataElementByName(dataGroupElement[0], dataGroupElement[1])
			if err != nil {
				logging.Error("Unable to get the data element id : " + err.Error())
				return err
			}

			isEncrypted := false
			isPassword := false

			res, err := db.Query("Call get_element_value(?,?)", profileID, dataElementId)
			if err != nil {
				return err
			}

			var currentValue string
			for res.Next() {
				err = res.Scan(&currentValue, &isEncrypted, &isPassword)
				if err != nil {
					res.Close()
					logging.Error(err.Error())
					return err
				}
			}
			res.Close()

			//to handle the existing stored clear values
			clearCurrentValue := currentValue
			if isEncrypted && currentValue != "" {
				currentValue, err = crypt.Decrypt(currentValue)
				if err != nil {
					logging.Error("ApproveBulkApprovalBulkSiteUpdate crypt.Decrypt error : " + err.Error())
					currentValue = clearCurrentValue
				}
			}

			if column == len(rows)-1 {
				value = string(bytes.Trim([]byte(value), ASCIINULL))
				if strings.HasSuffix(value, "\"") {
					value = value[0 : len(value)-1]
				}
			}

			if strings.ToLower(value) == "true" || strings.ToLower(value) == "false" {
				value = strings.ToLower(value)
			}

			if currentValue != value {
				if value == "" {
					continue
				}

				if isEncrypted && crypt.UseEncryption {
					value = crypt.Encrypt(value)
					isEncrypted = true
				}

				_, err = db.Exec("Call bulk_site_profile_data_update(?,?,?,?,?)", profileID, dataElementId, value, currentUser, isEncrypted)
				if err != nil {
					logging.Error("Error performing bulk_site_profile_data_update : " + err.Error())
					return err
				}
			}
		}
	}

	_, err = db.Exec("UPDATE bulk_approvals SET approved_by=?, approved_at=NOW(), approved=1 WHERE filename=?", currentUser, fileName)
	if err != nil {
		logging.Error("Error while updating bulk_tid_flagging : " + err.Error())
		return err
	}

	return nil
}

func ApproveBulkApprovalBulkTidUpdate(fileName, currentUser string, records [][]string) error {
	db, err := GetDB()
	if err != nil {
		logging.Error("Unable to get the database instance : " + err.Error())
		return err
	}

	header := records[0]

	for _, rows := range records[1:] {
		tidString := rows[0]
		tidId, err := strconv.Atoi(tidString)
		if err != nil {
			logging.Error("An error occured during tid conversion", err.Error())
			return err
		}

		tidProfileExists, profileID, siteId, err := CheckTidProfileExistsAndGetSiteId(tidId)
		if err != nil {
			logging.Error("An error occured during executing GetTidProfileIdForSiteId ", err.Error())
			return err
		}

		if !tidProfileExists {
			logging.Information("Creating tid override:", tidString)
			profileTypeId, err := GetProfileTypeId("tid", tidString, 1, currentUser)
			if err != nil {
				logging.Error("An error occured during executing GetProfileTypeId", err.Error())
				return err
			}

			err = CreateTidoverrideAndSaveNewprofileChange(siteId, profileTypeId, currentUser, ApproveCreate, "Override Created", tidString, tidId, 1)
			if err != nil {
				logging.Error("An error occured during executing CreateTidoverrideAndSaveNewprofileChange ", err.Error())
				return err
			}

			profileID, err = GetProfileIdFromTID(strings.TrimSpace(tidString))
			if err != nil {
				logging.Error("An error occured during executing GetProfileIdFromTID", err.Error())
				return err
			}
		}

		for column, value := range rows {
			if column == 0 || value == "" {
				continue
			}
			value = strings.TrimSpace(value)
			dataGroupElement := strings.Split(strings.TrimSpace(header[column]), ".")
			dataElementId, err := GetDataElementByName(dataGroupElement[0], dataGroupElement[1])
			if err != nil {
				logging.Error("Unable to get the data element id : " + err.Error())
				return err
			}

			isEncrypted := false
			isPassword := false

			res, err := db.Query("Call get_element_value(?,?)", int(profileID), dataElementId)
			if err != nil {
				logging.Error("Error while fetching data element value : " + err.Error())
				return err
			}
			res.Close()

			var currentValue string
			for res.Next() {
				err = res.Scan(&currentValue, &isEncrypted, &isPassword)
				if err != nil {
					res.Close()
					logging.Error(err.Error())
				}
			}
			res.Close()

			//This is to handle the existing values stored as clear value
			clearCurrentValue := currentValue
			if isEncrypted && currentValue != "" {
				currentValue, err = crypt.Decrypt(currentValue)
				if err != nil {
					logging.Error("Error while Decrypting : " + err.Error())
					currentValue = clearCurrentValue
				}
			}

			if column == len(rows)-1 {
				value = string(bytes.Trim([]byte(value), ASCIINULL))
				if strings.HasSuffix(value, "\"") {
					value = value[0 : len(value)-1]
				}
			}

			if strings.ToLower(value) == "true" || strings.ToLower(value) == "false" {
				value = strings.ToLower(value)
			}

			if currentValue != value {
				if value == "" {
					continue
				}

				if isEncrypted && crypt.UseEncryption {
					value = crypt.Encrypt(value)
					isEncrypted = true
				}

				_, err = db.Exec("Call bulk_tid_profile_data_update(?,?,?,?,?)", int(profileID), dataElementId, value, currentUser, isEncrypted)
				if err != nil {
					logging.Error("Error performing bulk_tid_profile_data_update : " + err.Error())
					return err
				}
			}
		}
	}

	_, err = db.Exec("UPDATE bulk_approvals SET approved_by=?, approved_at=NOW(), approved=1 WHERE filename=?", currentUser, fileName)
	if err != nil {
		logging.Error("Error while updating bulk_approvals : " + err.Error())
		return err
	}

	return nil
}

func ApproveBulkApprovalBulkDelete(fileName, currentUser string, records [][]string, deleteType string) error {
	db, err := GetDB()
	if err != nil {
		_, _ = logging.Error("Unable to get the database instance : " + err.Error())
		return err
	}

	if deleteType == "BulkTidDelete" {
		for col, val := range records {
			value := strings.TrimSpace(val[0])
			if col != 0 {
				value = string(bytes.TrimPrefix([]byte(strings.ToLower(value)), common.ByteOrderMark))
				tidID, err := strconv.Atoi(value)
				if err != nil {
					_, _ = logging.Error("Error while Tid conversion ")
					return err
				}
				tidString := GetPaddedTidId(tidID)
				if err != nil {
					_, _ = logging.Error("Error fetching Site ID for TID")
					return err
				}
				success, err := NewPEDRepository().DeleteByTid(tidString)
				if err != nil {
					return err
				} else if !success {
					return errors.New(fmt.Sprintf("Failure deleting tid %s", tidString))
				}

			}
		}
	} else if deleteType != "FileManagement" {
		for col, val := range records {
			value := strings.TrimSpace(val[0])
			if col != 0 {
				value = string(bytes.TrimPrefix([]byte(strings.ToLower(value)), common.ByteOrderMark))
				var profileID string
				profileID, err := GetProfileIDFromSiteName(value, 4)
				if err != nil {
					_, _ = logging.Error("Error while fetching profileID for site" + err.Error())
					return err
				}
				err = DeleteSite(profileID)
				if err != nil {
					_, _ = logging.Error("Error while attempting to delete Site : " + err.Error())
					return err
				}
			}
		}
	}

	_, err = db.Exec("UPDATE bulk_approvals SET approved_by=?, approved_at=NOW(), approved=1 WHERE filename=? AND approved=0", currentUser, fileName)
	if err != nil {
		_, _ = logging.Error("Error while updating bulk_approvals : " + err.Error())
		return err
	}

	return nil
}

func DiscardBulkApproval(fileName, fileType, currentUser string) error {
	db, err := GetDB()
	if err != nil {
		logging.Error("Unable to get the database instance : " + err.Error())
		return err
	}

	if fileName == "" {
		_, err = db.Exec("UPDATE bulk_approvals SET approved_by=?, approved_at=NOW(), approved=-1 WHERE approved=0", currentUser)
	} else {
		_, err = db.Exec("UPDATE bulk_approvals SET approved_by=?, approved_at=NOW(), approved=-1 WHERE filename=? AND filetype = ? AND approved=0", currentUser, fileName, fileType)
	}
	if err != nil {
		logging.Error("Error while updating bulk_approvals : " + err.Error())
		return err
	}

	return nil
}

func ApproveBulkPaymentUpload(fileName, currentUser string, records [][]string, deleteType string) error {
	db, err := GetDB()
	if err != nil {
		logging.Error("Unable to get the database instance : " + err.Error())
		return err
	}

	_, err = db.Exec("UPDATE bulk_approvals SET approved_by=?, approved_at=NOW(), approved=1 WHERE filename=? AND approved=0", currentUser, fileName)
	if err != nil {
		logging.Error("Error while updating bulk_approvals : " + err.Error())
		return err
	}
	return nil
}
