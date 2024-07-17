package dal

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"nextgen-tms-website/common"
	"nextgen-tms-website/config"
	"nextgen-tms-website/crypt"
	"nextgen-tms-website/entities"
	"nextgen-tms-website/fileServer"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func GetPendingChanges(profileType int, user *entities.TMSUser) ([]*ChangeApprovalHistory, int, error) {
	var changes = make([]*ChangeApprovalHistory, 0)
	db, err := GetDB()
	if err != nil {
		_, _ = logging.Error(fmt.Sprintf("Failed to get approvals for profile type '%v'; Failed to get DB: %v", profileType, err))
		return nil, 0, err
	}

	// Find user acquirers to limit search results
	acquirers, err := GetUserAcquirerPermissions(user)
	if err != nil {
		_, _ = logging.Error(fmt.Sprintf("Failed to get approvals for profile type '%v'; Failed to get aquirers: %v", profileType, err))
		return nil, 0, err
	}

	var rows *sql.Rows
	switch profileType {
	case Site:
		rows, err = db.Query("Call get_pending_site_change_approvals(?)", acquirers)
	case Chain:
		rows, err = db.Query("Call get_pending_chain_change_approvals(?)", acquirers)
	case Acquirer:
		rows, err = db.Query("Call get_pending_acquirer_change_approvals(?)", acquirers)
	case Tid:
		rows, err = db.Query("Call get_pending_tid_change_approvals(?)", acquirers)
	}

	if err != nil {
		_, _ = logging.Error(fmt.Sprintf("Failed to get approvals for profile type '%v' and acquirers '%v'; %v", profileType, acquirers, err))
		return nil, 0, err
	}

	var entityName sql.NullString
	var originalValue sql.NullString
	var updatedValue sql.NullString
	var reviewedBy sql.NullString
	var reviewedAt sql.NullString
	var tidID sql.NullString
	var merchantID sql.NullString
	var isEncrypted sql.NullBool
	var isPassword sql.NullBool

	defer rows.Close()
	limit := 0
	for rows.Next() {
		if limit < 50 {
			var changeHistory = &ChangeApprovalHistory{}
			err = rows.Scan(
				&changeHistory.ProfileDataID,
				&entityName,
				&changeHistory.Field,
				&changeHistory.ChangeType,
				&originalValue,
				&updatedValue,
				&changeHistory.ChangedBy,
				&changeHistory.ChangedAt,
				&changeHistory.Approved,
				&reviewedBy,
				&reviewedAt,
				&tidID,
				&merchantID,
				&isEncrypted,
				&isPassword)

			if err != nil {
				_, _ = logging.Error(fmt.Sprintf("Failed to get approvals for profile type '%v' and acquirers '%v'; %v", profileType, acquirers, err))
				return nil, limit, err
			}

			if entityName.Valid {
				changeHistory.Identifier = entityName.String
			} else {
				// If the site name is null, then this means that the entity (e.g. site, TID) has been deleted, so skip the row
				continue
			}

			if originalValue.Valid {
				if isEncrypted.Valid && isEncrypted.Bool && originalValue.String != "" {
					changeHistory.OriginalValue, err = crypt.Decrypt(originalValue.String)
					if err != nil {
						logging.Error("SaveUnapprovedElement : Error while Decrypting : " + err.Error())
						changeHistory.OriginalValue = originalValue.String
					}
				} else {
					changeHistory.OriginalValue = strings.ReplaceAll(strings.ReplaceAll(originalValue.String, `"TidId":0`, ``), `"tidId":-1`, ``)
				}
			}

			if updatedValue.Valid {
				if isEncrypted.Valid && isEncrypted.Bool {
					changeHistory.ChangeValue, err = crypt.Decrypt(updatedValue.String)
					if err != nil {
						return nil, limit, err
					}
				} else {

					re := regexp.MustCompile(`"PIN":"\d+"`)
					changeHistory.ChangeValue = re.ReplaceAllString(updatedValue.String, `"PIN":"*****"`)
					changeHistory.ChangeValue = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(changeHistory.ChangeValue, `"SiteId":0`, ``), `"tidId":-1`, ``), `"TidId":0`, ``)

				}
			}

			if isEncrypted.Valid {
				changeHistory.IsEncrypted = isEncrypted.Bool
			}

			if isPassword.Valid {
				changeHistory.IsPassword = isPassword.Bool
			}

			if reviewedBy.Valid {
				changeHistory.ReviewedBy = reviewedBy.String
			}

			if reviewedAt.Valid {
				changeHistory.ReviewedAt = reviewedAt.String
			}

			if tidID.Valid {
				changeHistory.TidId = tidID.String
			}

			if merchantID.Valid {
				changeHistory.MID = merchantID.String
			}

			changes = append(changes, changeHistory)
		}
		limit++
	}

	return changes, limit, nil
}

func GetPaymentServicePendingChanges(profileType int, user *entities.TMSUser) ([]*ChangeApprovalHistory, int, error) {
	var changes = make([]*ChangeApprovalHistory, 0)
	db, err := GetDB()
	if err != nil {
		_, _ = logging.Error(fmt.Sprintf("Failed to get approvals for profile type '%v'; Failed to get DB: %v", profileType, err))
		return nil, 0, err
	}

	// Find user acquirers to limit search results
	acquirers, err := GetUserAcquirerPermissions(user)
	if err != nil {
		_, _ = logging.Error(fmt.Sprintf("Failed to get approvals for profile type '%v'; Failed to get aquirers: %v", profileType, err))
		return nil, 0, err
	}

	var rows *sql.Rows
	rows, err = db.Query("Call get_pending_others_change_approvals(?)", acquirers)

	if err != nil {
		_, _ = logging.Error(fmt.Sprintf("Failed to get approvals for profile type '%v' and acquirers '%v'; %v", profileType, acquirers, err))
		return nil, 0, err
	}

	var entityName sql.NullString
	var originalValue sql.NullString
	var updatedValue sql.NullString
	var reviewedBy sql.NullString
	var reviewedAt sql.NullString
	var tidID sql.NullString
	var merchantID sql.NullString
	var isEncrypted sql.NullBool
	var isPassword sql.NullBool
	var PaymentServiceGroupName sql.NullString
	var PaymentServiceParentGroupName sql.NullString
	var PaymentServiceName sql.NullString

	defer rows.Close()
	limit := 0
	for rows.Next() {
		if limit < 50 {
			var changeHistory = &ChangeApprovalHistory{}
			err = rows.Scan(
				&changeHistory.ProfileDataID,
				&entityName,
				&changeHistory.Field,
				&changeHistory.ChangeType,
				&originalValue,
				&updatedValue,
				&changeHistory.ChangedBy,
				&changeHistory.ChangedAt,
				&changeHistory.Approved,
				&reviewedBy,
				&reviewedAt,
				&tidID,
				&merchantID,
				&isEncrypted,
				&isPassword,
				&PaymentServiceGroupName,
				&PaymentServiceName,
				&PaymentServiceParentGroupName)

			if err != nil {
				_, _ = logging.Error(fmt.Sprintf("Failed to get approvals for profile type '%v' and acquirers '%v'; %v", profileType, acquirers, err))
				return nil, limit, err
			}

			if entityName.Valid {
				changeHistory.Identifier = entityName.String
			} else {
				// If the site name is null, then this means that the entity (e.g. site, TID) has been deleted, so skip the row
				continue
			}

			if originalValue.Valid {
				if isEncrypted.Valid && isEncrypted.Bool && originalValue.String != "" {
					changeHistory.OriginalValue, err = crypt.Decrypt(originalValue.String)
					if err != nil {
						//This is to handle the existing values stored as clear value
						logging.Error("GetPaymentServicePendingChanges : Error while Decrypting : " + err.Error())
						changeHistory.OriginalValue = originalValue.String
					}
				} else {
					changeHistory.OriginalValue = originalValue.String
				}
			}

			if updatedValue.Valid {
				if isEncrypted.Valid && isEncrypted.Bool {
					changeHistory.ChangeValue, err = crypt.Decrypt(updatedValue.String)
					if err != nil {
						return nil, limit, err
					}
				} else {
					changeHistory.ChangeValue = updatedValue.String
				}
			}

			if isEncrypted.Valid {
				changeHistory.IsEncrypted = isEncrypted.Bool
			}

			if isPassword.Valid {
				changeHistory.IsPassword = isPassword.Bool
			}

			if reviewedBy.Valid {
				changeHistory.ReviewedBy = reviewedBy.String
			}

			if reviewedAt.Valid {
				changeHistory.ReviewedAt = reviewedAt.String
			}

			if tidID.Valid {
				changeHistory.TidId = tidID.String
			}

			if merchantID.Valid {
				changeHistory.MID = merchantID.String
			}

			if PaymentServiceGroupName.Valid {
				changeHistory.PaymentServiceGroupName = PaymentServiceGroupName.String
			}

			if PaymentServiceName.Valid {
				changeHistory.PaymentServiceName = PaymentServiceName.String
			}

			if PaymentServiceParentGroupName.Valid {
				changeHistory.PaymentServiceGroupName = PaymentServiceParentGroupName.String
			}

			changes = append(changes, changeHistory)
		}
		limit++
	}

	return changes, limit, nil
}

func GetFilteredPendingChanges(offset int, profileType int, user *entities.TMSUser) ([]*ChangeApprovalHistory, error) {
	var changes = make([]*ChangeApprovalHistory, 0)
	db, err := GetDB()
	if err != nil {
		_, _ = logging.Error(fmt.Sprintf("Failed to get approvals for profile type '%v'; Failed to get DB: %v", profileType, err))
		return nil, err
	}

	// Find user acquirers to limit search results
	acquirers, err := GetUserAcquirerPermissions(user)
	if err != nil {
		_, _ = logging.Error(fmt.Sprintf("Failed to get approvals for profile type '%v'; Failed to get aquirers: %v", profileType, err))
		return nil, err
	}

	var rows *sql.Rows
	switch profileType {
	case Site:
		rows, err = db.Query("Call get_pending_site_change_approvals(?)", acquirers)
	case Chain:
		rows, err = db.Query("Call get_pending_chain_change_approvals(?)", acquirers)
	case Acquirer:
		rows, err = db.Query("Call get_pending_acquirer_change_approvals(?)", acquirers)
	case Tid:
		rows, err = db.Query("Call get_pending_tid_change_approvals(?)", acquirers)
	}

	if err != nil {
		_, _ = logging.Error(fmt.Sprintf("Failed to get approvals for profile type '%v' and acquirers '%v'; %v", profileType, acquirers, err))
		return nil, err
	}

	defer rows.Close()

	var entityName sql.NullString
	var originalValue sql.NullString
	var updatedValue sql.NullString
	var reviewedBy sql.NullString
	var reviewedAt sql.NullString
	var tidID sql.NullString
	var merchantID sql.NullString
	var isEncrypted sql.NullBool
	var isPassword sql.NullBool

	for rows.Next() {
		var changeHistory = &ChangeApprovalHistory{}
		err = rows.Scan(
			&changeHistory.ProfileDataID,
			&entityName,
			&changeHistory.Field,
			&changeHistory.ChangeType,
			&originalValue,
			&updatedValue,
			&changeHistory.ChangedBy,
			&changeHistory.ChangedAt,
			&changeHistory.Approved,
			&reviewedBy,
			&reviewedAt,
			&tidID,
			&merchantID,
			&isEncrypted,
			&isPassword)

		if err != nil {
			_, _ = logging.Error(fmt.Sprintf("Failed to get approvals for profile type '%v' and acquirers '%v'; %v", profileType, acquirers, err))
			return nil, err
		}

		if entityName.Valid {
			changeHistory.Identifier = entityName.String
		} else {
			// If the site name is null, then this means that the entity (e.g. site, TID) has been deleted, so skip the row
			continue
		}

		if originalValue.Valid {
			if isEncrypted.Valid && isEncrypted.Bool {
				changeHistory.OriginalValue, err = crypt.Decrypt(originalValue.String)
				if err != nil {
					//This is to handle the existing values stored as clear value
					logging.Error("GetFilteredPendingChanges : Error while Decrypting : " + err.Error())
					changeHistory.OriginalValue = originalValue.String
				}
			} else {
				changeHistory.OriginalValue = originalValue.String
			}
		}

		if updatedValue.Valid {
			if isEncrypted.Valid && isEncrypted.Bool {
				changeHistory.ChangeValue, err = crypt.Decrypt(updatedValue.String)
				if err != nil {
					return nil, err
				}
			} else {
				changeHistory.ChangeValue = updatedValue.String
			}
		}

		if isPassword.Valid {
			changeHistory.IsPassword = isPassword.Bool
		}

		if isEncrypted.Valid {
			changeHistory.IsEncrypted = isEncrypted.Bool
		}

		if reviewedBy.Valid {
			changeHistory.ReviewedBy = reviewedBy.String
		}

		if reviewedAt.Valid {
			changeHistory.ReviewedAt = reviewedAt.String
		}

		if tidID.Valid {
			changeHistory.TidId = tidID.String
		}

		if merchantID.Valid {
			changeHistory.MID = merchantID.String
		}

		changes = append(changes, changeHistory)
	}

	start := offset
	end := offset + 50
	length := len(changes)

	// range completely within array bounds
	if start <= length && end <= length {
		return changes[start:end], nil
	}

	// range overruns bounds
	if start <= length && end > length {
		return changes[start:], nil
	}

	// start position outside array bounds
	return nil, nil
}

func BuildFilteredChangeApprovalHistory(after string, name string, user string,
	before string, field string, limit bool, offset int, tmsUser *entities.TMSUser) ([]*ChangeApprovalHistory, error) {
	var changes = make([]*ChangeApprovalHistory, 0)
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	var afterDate sql.NullString
	var beforeDate sql.NullString

	if after != "" {
		afterDate.String = after
		afterDate.Valid = true
	}

	if before != "" {
		beforeDate.String = before
		beforeDate.Valid = true
	}

	// Find user acquirers to limit search results
	acquirers, err := GetUserAcquirerPermissions(tmsUser)
	if err != nil {
		return nil, err
	}
	acquirers = strings.Trim(acquirers, "'")
	acquirers = strings.Trim(acquirers, ",")

	var rows *sql.Rows
	rows, err = db.Query("Call get_change_approval_history(?, ?, ?, ?, ?, ?)", afterDate, name, user, beforeDate, field, acquirers)

	if err != nil {
		return nil, err
	}

	var originalValue sql.NullString
	var updatedValue sql.NullString
	var reviewedBy sql.NullString
	var reviewedAt sql.NullString
	var changedBy sql.NullString
	var tidId sql.NullString
	var merchantId sql.NullString
	var isPassword sql.NullBool
	var isEncrypted sql.NullBool

	defer rows.Close()
	for rows.Next() {
		var changeHistory = &ChangeApprovalHistory{}
		err = rows.Scan(
			&changeHistory.ProfileDataID,
			&changeHistory.Identifier,
			&changeHistory.Field,
			&originalValue,
			&updatedValue,
			&changedBy,
			&changeHistory.ChangedAt,
			&changeHistory.Approved,
			&reviewedBy,
			&reviewedAt,
			&tidId,
			&merchantId,
			&isPassword,
			&isEncrypted,
			&changeHistory.ChangeType)

		if err != nil {
			return nil, err
		}

		if originalValue.Valid {
			if isEncrypted.Valid && isEncrypted.Bool && originalValue.String != "" {
				changeHistory.OriginalValue, err = crypt.Decrypt(originalValue.String)
				if err != nil {
					//This is to handle the existing values stored as clear value
					logging.Error("GetPaymentServicePendingChanges : Error while Decrypting : " + err.Error())
					changeHistory.OriginalValue = originalValue.String
				}
			} else {
				re := regexp.MustCompile(`"PIN":"\d+"`)
				changeHistory.OriginalValue = re.ReplaceAllString(originalValue.String, `"PIN":"*****"`)
				changeHistory.OriginalValue = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(changeHistory.OriginalValue, `"SiteId":0`, ``), `"tidId":-1`, ``), `"TidId":0`, ``)
			}
		}

		if updatedValue.Valid {
			if isEncrypted.Valid && isEncrypted.Bool {
				changeHistory.ChangeValue, err = crypt.Decrypt(updatedValue.String)
				if err != nil {
					return nil, err
				}
			} else {
				re := regexp.MustCompile(`"PIN":"\d+"`)
				changeHistory.ChangeValue = re.ReplaceAllString(updatedValue.String, `"PIN":"*****"`)
				changeHistory.ChangeValue = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(changeHistory.ChangeValue, `"SiteId":0`, ``), `"tidId":-1`, ``), `"TidId":0`, ``)
			}
		}

		if isEncrypted.Valid && isEncrypted.Bool {
			changeHistory.IsEncrypted = isEncrypted.Bool
		}

		if isPassword.Valid && isPassword.Bool {
			changeHistory.IsPassword = isPassword.Bool
		}

		if reviewedBy.Valid {
			changeHistory.ReviewedBy = reviewedBy.String
		}

		if reviewedAt.Valid {
			changeHistory.ReviewedAt = reviewedAt.String
		}

		if changedBy.Valid {
			changeHistory.ChangedBy = changedBy.String
		}

		if tidId.Valid {
			changeHistory.TidId = tidId.String
		}

		if merchantId.Valid {
			changeHistory.MID = merchantId.String
		}

		changes = append(changes, changeHistory)
	}

	if limit {
		start := offset
		end := offset + 50
		length := len(changes)

		// range completely within array bounds
		if start <= length && end <= length {
			return changes[start:end], nil
		}

		// range overruns bounds
		if start <= length && end > length {
			return changes[start:], nil
		}
	} else {
		return changes, nil
	}

	// start position outside array bounds
	return nil, nil
}

func DuplicateChainChangeApproval(user *entities.TMSUser, chainProfileId int, newChainName, acquirerName string) (int, error) {
	db, err := GetDB()
	if err != nil {
		return 0, err
	}

	res, err := db.Exec(`INSERT INTO approvals (profile_id,data_element_id, change_type, current_value, new_value, created_at, approved, created_by, acquirer,approved_by,approved_at  )
				   VALUE
				   (?, 1, 7,?, ?, NOW(), 1, ?, ?,?,NOW())`,
		chainProfileId, newChainName, "duplicate chain creation", user.Username, acquirerName, user.Username)
	if err != nil {
		logging.Error(err.Error())
		return 0, err
	}

	id, err := res.LastInsertId()
	return int(id), err
}

func PaymentServiceCreationChangeApproval(userName string, chainProfileId, changeType int, newPaymentServiceGroupName, newPaymentServiceName, acquirerName string) (int, error) {
	db, err := GetDB()
	if err != nil {
		return 0, err
	}

	res, err := db.Exec(`INSERT INTO approvals (profile_id,data_element_id, change_type, current_value, new_value, created_at, approved, created_by, acquirer,approved_by,approved_at  )
				   VALUE
				   (?, 1, ?,?, ?, NOW(), 1, ?, ?,?,NOW())`,
		chainProfileId, changeType, newPaymentServiceGroupName, newPaymentServiceName, userName, acquirerName, userName)
	if err != nil {
		logging.Error(err.Error())
		return 0, err
	}

	id, err := res.LastInsertId()
	return int(id), err
}
func PaymentServicesDeletionChangeApproval(userName string, chainProfileId, changeType int, serviceId, newPaymentServiceName, acquirerName string) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	_, err = db.Exec("Call payment_service_deletion(?,?,?,?,?,?)", chainProfileId, changeType, serviceId, newPaymentServiceName, userName, acquirerName)
	if err != nil {
		return err
	}

	return nil
}

func PaymentServiceGroupDeletionChangeApproval(userName string, chainProfileId, changeType int, groupId, newPaymentServiceName, acquirerName string) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	_, err = db.Exec("Call payment_service_group_deletion(?,?,?,?,?,?)", chainProfileId, changeType, groupId, newPaymentServiceName, userName, acquirerName)
	if err != nil {
		return err
	}

	return nil
}

func ApproveAllChanges(model []*ChangeApprovalHistory, currentUser, profileType string) error {
	for _, change := range model {
		id, _ := strconv.Atoi(change.ProfileDataID)
		err := ApproveChange(id, currentUser, profileType)
		if err != nil {
			return err
		}
	}
	return nil
}

func ApproveChange(profileDataID int, currentUser, profileType string) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	rows, err := db.Query("SELECT profile_id, change_type, tid_id, merchant_id, data_element_id, new_value, current_value FROM approvals WHERE approval_id = ?", profileDataID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var sqlProfileID, changeType, dataElementID sql.NullInt64
	var tidID, merchantID, newValue, currentValue sql.NullString

	for rows.Next() {
		if err = rows.Scan(&sqlProfileID, &changeType, &tidID, &merchantID, &dataElementID, &newValue, &currentValue); err != nil {
			return err
		}
		var dataElement sql.NullString
		if sqlProfileID.Valid && changeType.Valid {
			if err = db.QueryRow("select name from data_element where data_element_id=?", dataElementID.Int64).Scan(&dataElement); err != nil {
				return err
			}

			switch dataElement.String {
			case "fraud":
				err = handleFraud(db, sqlProfileID, profileDataID, newValue, currentUser)
			case "users":
				err = handleUser(db, sqlProfileID, profileDataID, newValue, currentUser)
			case "flagStatus":
				err = handleFlagStatus(db, sqlProfileID, profileDataID, newValue, currentUser)
			case "time":
				err = handleTime(db, sqlProfileID, profileDataID, newValue, currentUser)
			case "auto":
				err = handleAuto(db, sqlProfileID, profileDataID, newValue, currentUser)
			default:
				if changeType.Int64 == 10 || changeType.Int64 == 11 {
					err = handleServiceDeletion(db, profileDataID, changeType, currentValue, currentUser)
				} else {
					_, err = db.Exec("CALL approve_change(?, ?)", profileDataID, currentUser)
				}
			}
			if err != nil {
				return err
			}
		}

		if profileType == "TID" {
			err = UpdateTIDFlagWithProfileID(sqlProfileID.Int64)
			if err != nil {
				_, _ = logging.Error(err)
			}
		} else if profileType == "Site" && dataElement.String == "name" {
			err = UpdateProfileNameWithProfileID(newValue.String, sqlProfileID.Int64)
			if err != nil {
				_, _ = logging.Error(err)
			}
		}

		if !(sqlProfileID.Valid && changeType.Valid && changeType.Int64 == ApproveDelete) {
			break
		}

		if tidID.Valid {
			err = handleTIDDeletion(tidID, profileType)
		}

		if merchantID.Valid {
			err = handleSiteDeletion(sqlProfileID)
		}
	}

	return err
}

func DiscardAllChanges(model []*ChangeApprovalHistory, currentUser string) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	for _, change := range model {
		_, err := db.Exec("Call discard_change(?,?)", change.ProfileDataID, currentUser)
		if err != nil {
			return err
		}
	}
	return nil
}

func DiscardChange(profileDataID int, currentUser string) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	_, err = db.Exec("Call discard_change(?,?)", profileDataID, currentUser)
	return err
}

func LogTidChange(tid, dataElementId, changeType int, currentValue, updatedValue, currentUser string, approved bool) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	var approvedInt int
	if approved {
		approvedInt = 1
	} else {
		approvedInt = 0
	}
	_, err = db.Exec("CALL log_tid_change(?, ?, ?, ?, ?, ?, ?)", tid, dataElementId, changeType, currentValue, updatedValue, currentUser, approvedInt)
	if err != nil {
		return err
	}
	return nil
}

func handleFlagStatus(db *sql.DB, sqlProfileID sql.NullInt64, profileDataID int, newValue sql.NullString, currentUser string) error {
	query := `UPDATE tid t
						JOIN tid_site ts
							ON ts.tid_id = t.tid_id
						JOIN site_profiles sp 
							ON sp.site_id = ts.site_id
						SET t.flag_status=true, t.flagged_date=CURRENT_TIMESTAMP
						WHERE sp.profile_id = ? `
	var tids []int
	skip := false
	if strings.HasPrefix(newValue.String, "file : ") {
		fileName := newValue.String[7:]
		fileAsBase64Encoded, err := fileServer.NewFsReader(config.FileserverURL).GetFile(fileName, config.FlaggingFileDirectory)
		if err != nil {
			return err
		}

		fileAsBytes, err := common.ConvertBase64FileToBytes(string(fileAsBase64Encoded))
		if err != nil {
			return err
		}

		csvReader := csv.NewReader(bytes.NewBuffer(fileAsBytes))
		records, err := csvReader.ReadAll()
		if err != nil {
			return err
		}

		for _, r := range records[1:] {
			// some null values are getting appended last row; inorder to remove that bytes.Trim is used
			r[0] = string(bytes.Trim([]byte(strings.TrimSpace(r[0])), ASCIINULL))
			if strings.TrimSpace(r[0]) == "" {
				continue
			}
			tid, err := strconv.Atoi(r[0])
			if err != nil {
				logging.Error("Invalid tid : " + r[0])
				continue
			}
			tids = append(tids, tid)
		}

		if len(tids) == 0 {
			skip = true
		} else {
			query += " AND t.tid_id IN (?" + strings.Repeat(",?", len(tids)-1) + ")"
		}
	} else if newValue.String != "all" {
		tidsArray := strings.Split(newValue.String, ", ")
		for _, r := range tidsArray {
			if strings.TrimSpace(r) == "" {
				continue
			}

			tid, err := strconv.Atoi(strings.TrimSpace(r))
			if err != nil {
				logging.Error("Invalid tid : " + r)
				continue
			}
			tids = append(tids, tid)
		}

		if len(tids) == 0 {
			skip = true
		} else {
			query += " AND t.tid_id IN (?" + strings.Repeat(",?", len(tids)-1) + ")"
		}
	}

	if !skip {
		queryArgs := make([]interface{}, len(tids)+1)
		queryArgs[0] = sqlProfileID.Int64
		for i, tid := range tids {
			queryArgs[i+1] = tid
		}

		r, err := db.Exec(query, queryArgs...)
		if err != nil {
			return err
		}

		n, err := r.RowsAffected()
		if err != nil {
			return err
		}
		logging.Information(fmt.Sprintf("Change approved, updated flag status for %d tids in profile %d", n, sqlProfileID.Int64))
	}

	_, err := db.Exec("UPDATE approvals a SET approved = 1, approved_by = ?, approved_at = NOW() WHERE a.approval_id = ?", currentUser, profileDataID)
	if err != nil {
		return err
	}
	return nil
}

func handleTime(db *sql.DB, sqlProfileID sql.NullInt64, profileDataID int, newValue sql.NullString, currentUser string) error {

	var err error
	timeRange := strings.Split(newValue.String, " | ")

	if len(timeRange) == 1 {
		_, err = db.Exec("Call update_eod_auto_time(?,?)", sqlProfileID.Int64, timeRange[0])
		if err != nil {
			return err
		}
	} else {
		var startTime, endTime time.Time
		startTime, err = time.Parse("15:04", timeRange[0])
		if err != nil {
			return err
		}

		endTime, err = time.Parse("15:04", timeRange[1])
		if err != nil {
			return err
		}

		timeDiff := endTime.Sub(startTime).Minutes()
		if timeDiff < 0 {
			endTime = endTime.AddDate(0, 0, 1)
		}
		timeDiff = endTime.Sub(startTime).Minutes()

		_, err = db.Exec("Call update_eod_auto_time_range(?, ?, ?, ?)", sqlProfileID.Int64, startTime.Hour(), startTime.Minute(), int(timeDiff))
		if err != nil {
			return err
		}
	}
	_, err = db.Exec("Call approve_change(?, ?)", profileDataID, currentUser)
	if err != nil {
		return err
	}

	return nil
}

func handleAuto(db *sql.DB, sqlProfileID sql.NullInt64, profileDataID int, newValue sql.NullString, currentUser string) error {
	var err error
	auto := false
	// for handling removeOverride, while empty
	if newValue.String != "" {
		auto, err = strconv.ParseBool(newValue.String)
		if err != nil {
			return err
		}
	}

	_, err = db.Exec("Call update_eod_auto(?,?)", sqlProfileID.Int64, auto)
	if err != nil {
		return err
	}

	_, err = db.Exec("Call approve_change(?, ?)", profileDataID, currentUser)
	if err != nil {
		return err
	}
	return nil
}

func handleServiceDeletion(db *sql.DB, profileDataID int, changeType sql.NullInt64, currentValue sql.NullString, currentUser string) error {
	var err error

	_, err = db.Exec("Call approve_change(?, ?)", profileDataID, currentUser)
	if err != nil {
		return err
	}
	if changeType.Int64 == 10 {
		if !DeleteServiceGroup(currentValue.String) {
			return errors.New("unable to delete payment service group")
		}
	}

	if changeType.Int64 == 11 {
		if !DeleteService(currentValue.String) {
			return errors.New("unable to delete payment service")
		}
	}
	return nil
}

func handleTIDDeletion(tidID sql.NullString, profileType string) error {
	if profileType == "TID" {
		logging.Information(fmt.Sprintf("Change approved, deleting TID override for TID '%v'", tidID.String))
		success, err := NewPEDRepository().DeleteOverrideByTid(tidID.String)
		if err != nil {
			return err
		} else if !success {
			return errors.New(fmt.Sprintf("Failure deleting override for tid %s", tidID.String))
		}
	} else {
		logging.Information(fmt.Sprintf("Change approved, deleting TID for TID '%v'", tidID.String))
		success, err := NewPEDRepository().DeleteByTid(tidID.String)
		if err != nil {
			return err
		} else if !success {
			return errors.New(fmt.Sprintf("Failure deleting tid %s", tidID.String))
		}
	}
	return nil
}

func handleSiteDeletion(sqlProfileID sql.NullInt64) error {
	err := DeleteSite(strconv.FormatInt(sqlProfileID.Int64, 10))
	if err != nil {
		return err
	}
	return nil
}

func handleFraud(db *sql.DB, sqlProfileID sql.NullInt64, profileDataID int, newValue sql.NullString, currentUser string) error {
	fraudData := make(map[string]interface{})
	if err := json.Unmarshal([]byte(newValue.String), &fraudData); err != nil {
		return err
	}

	limits, err := json.Marshal(fraudData["velocityLimits"])
	if err != nil {
		return err
	}

	var limitList []entities.VelocityLimit
	if err = json.Unmarshal(limits, &limitList); err != nil {
		return err
	}

	dailyTxnCleanseTime := fraudData["dailyTxnCleanseTime"].(string)
	if err := SaveProfileDataElementByName(int(sqlProfileID.Int64), "dailyTxnCleanseTime", dailyTxnCleanseTime); err != nil {
		return err
	}

	limitLevel := int(fraudData["limitString"].(float64))
	var schemeLimitLevel int
	//Determine from the limit level (3 for site, 4 for ped) what level the scheme limits are
	if limitLevel == 3 {
		schemeLimitLevel = 1
	} else {
		schemeLimitLevel = 2
	}
	site := int(fraudData["siteID"].(float64))
	tid := int(fraudData["tidId"].(float64))
	//Clear the previously stored scheme velocity limits
	if err = DeleteSiteVelocityLimitsAndTxn(site, schemeLimitLevel, limitLevel, tid); err != nil {
		return err
	}

	if err = SetSiteVelocityLimits(site, tid, limitList); err != nil {
		return err
	}

	if _, err = db.Exec("UPDATE approvals a SET tid_id= ?, approved = 1, approved_by = ?, approved_at = NOW() WHERE a.approval_id = ?", tid, currentUser, profileDataID); err != nil {
		return err
	}
	return nil
}

func handleUser(db *sql.DB, sqlProfileID sql.NullInt64, profileDataID int, newValue sql.NullString, currentUser string) error {
	usersData := make(map[string]interface{})
	var siteId int
	var tidId float64
	var ok bool
	if err := json.Unmarshal([]byte(newValue.String), &usersData); err != nil {
		return err
	}

	users, err := json.Marshal(usersData["updatedUsers"])
	if err != nil {
		return err
	}

	var usersToAddToSite = make([]*entities.SiteUser, 0)
	if err = json.Unmarshal(users, &usersToAddToSite); err != nil {
		return err
	}

	users, err = json.Marshal(usersData["newUsers"])
	if err != nil {
		return err
	}

	var newUsersToAddToSite = make([]*entities.SiteUser, 0)
	if err = json.Unmarshal(users, &newUsersToAddToSite); err != nil {
		return err
	}

	usersToAddToSite = append(usersToAddToSite, newUsersToAddToSite...)

	if val, found := usersData["tidID"]; found {
		if tidId, ok = val.(float64); ok {

			tidUsers, err := GetUsersForTid(int(tidId))
			if err != nil {
				return err
			}
			if len(tidUsers) <= 0 {
				siteId, err := GetSiteIDForTid(int(tidId))
				if err != nil {
					return err
				}
				SiteUsers, err := GetUsersForSite(siteId)
				if err != nil {
					return err
				}

				for _, siteUser := range SiteUsers {
					IsSiteUserIDPresent := false
					for _, checkUsersToAddToSite := range usersToAddToSite {
						if checkUsersToAddToSite.UserId == siteUser.UserId {
							IsSiteUserIDPresent = true
						}
					}
					if !IsSiteUserIDPresent {
						usersToAddToSite = append(usersToAddToSite, &entities.SiteUser{
							UserId:   siteUser.UserId,
							Username: siteUser.Username,
							PIN:      siteUser.PIN,
							Modules:  siteUser.Modules,
							SiteId:   siteUser.SiteId,
							TidId:    siteUser.TidId,
						})
					}
				}

			}

			if len(usersToAddToSite) > 0 {
				for _, user := range usersToAddToSite {
					PIN := user.PIN
					isEncrypted := false
					if crypt.UseEncryption {
						PIN = crypt.Encrypt(PIN)
						isEncrypted = true
					}
					modules := userModulesToString(*user)
					if _, err = db.Exec("CALL add_tid_user(?,?,?,?,?,?)", user.UserId, int(tidId), user.Username, PIN, modules, isEncrypted); err != nil {
						return err
					}
				}
			}
		}
		deletedUsers, err := json.Marshal(usersData["deletedUsers"])
		if err != nil {
			return err
		}

		var usersToDelete = make([]*entities.SiteUser, 0)
		if err = json.Unmarshal(deletedUsers, &usersToDelete); err != nil {
			return err
		}

		for _, deleteUser := range usersToDelete {
			if _, err = db.Exec("CALL delete_tid_user(?)", deleteUser.UserId); err != nil {
				return err
			}
		}
	} else {
		err = db.QueryRow("SELECT site_id FROM site_profiles WHERE profile_id = ?", int(sqlProfileID.Int64)).Scan(&siteId)
		if err != nil {
			return err
		}

		if len(usersToAddToSite) > 0 {
			for _, user := range usersToAddToSite {
				PIN := user.PIN
				isEncrypted := false
				if crypt.UseEncryption {
					PIN = crypt.Encrypt(PIN)
					isEncrypted = true
				}
				modules := userModulesToString(*user)
				if _, err = db.Exec("CALL add_site_user(?,?,?,?,?,?)", user.UserId, siteId, user.Username, PIN, modules, isEncrypted); err != nil {
					return err
				}
			}
		}

		deletedUsers, err := json.Marshal(usersData["deletedUsers"])
		if err != nil {
			return err
		}

		var usersToDelete = make([]*entities.SiteUser, 0)
		if err = json.Unmarshal(deletedUsers, &usersToDelete); err != nil {
			return err
		}

		for _, deleteUser := range usersToDelete {
			if _, err = db.Exec("CALL delete_site_user(?)", deleteUser.UserId); err != nil {
				return err
			}
		}
	}

	if _, err = db.Exec("UPDATE approvals a SET tid_id= ?, approved = 1, approved_by = ?, approved_at = NOW() WHERE a.approval_id = ?", tidId, currentUser, profileDataID); err != nil {
		return err
	}
	return nil
}
