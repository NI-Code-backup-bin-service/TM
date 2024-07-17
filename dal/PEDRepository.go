package dal

import (
	"database/sql"
	"fmt"
	"nextgen-tms-website/PED"
	"nextgen-tms-website/common"
	"strconv"
	"strings"

	exporter "bitbucket.org/network-international/nextgen-libs/nextgen-helpers/exportHandler"
)

type pedRepository struct{}

var (
	pedCommonQuery = `SELECT  
    	ts.tid_profile_id,
    	t.tid_id AS 'tid',
        t.serial AS 'serial',
		t.PIN AS 'pin',
        t.ExpiryDate AS 'expiryTime',
        t.ActivationDate AS 'activationTime',
        s.site_id AS 'siteId',
        pdn.datavalue AS 'siteName',
        p3.profile_id AS 'chainId',
		p3.name AS 'chainName',
        pd.datavalue AS 'merchantId',
        t.software_version AS 'appVer',
        t.firmware_version AS 'firmwareVer',
        date_format(t.last_transaction_time, "%Y-%m-%dT%H:%i:%sZ") AS 'lastTransaction',
        date_format(from_unixtime(t.last_checked_time / 1000), "%Y-%m-%dT%H:%i:%sZ") AS 'lastCheckedTime',
        date_format(from_unixtime(t.confirmed_time / 1000), "%Y-%m-%dT%H:%i:%sZ") AS 'confirmedTime',
        date_format(t.last_apk_download, "%Y-%m-%dT%H:%i:%sZ") AS 'lastAPKDownload',
		t.ip_address AS 'ipaddress',
		t.ip_addresses AS 'ipAddresses',
        t.sim_card_serial_number AS 'simCardSerialNumber',
        t.flag_status AS 'flagStatus',
        t.flagged_date AS 'flaggedDate',
        t.eod_auto AS 'eodAuto',
        t.auto_time AS 'autoTime',
		t.coordinates AS 'pedCoordinates',
        t.accuracy AS 'pedAccuracy',
        t.last_coordinate_time AS 'pedCoordinatesLastUpdated',
        t.free_internal_storage AS 'freeInternalStorage',
        t.total_internal_storage AS 'totalInternalStorage',
        t.softui_last_downloaded_file_name AS 'softuiLastDownloadedFileName',
        t.softui_last_downloaded_file_hash AS 'softuiLastDownloadedFileHash',
        t.softui_last_downloaded_file_list AS 'softuiLastDownloadedFileList',
        date_format(t.softui_last_downloaded_file_date_time, "%Y-%m-%dT%H:%i:%sZ") AS 'softuiLastDownloadedFileName'

    FROM tid t
      LEFT JOIN tid_site ts ON ts.tid_id = t.tid_id
      LEFT JOIN site s ON s.site_id = ts.site_id
      LEFT JOIN (
        site_profiles tp2
          join profile p2 on p2.profile_id = tp2.profile_id
          join profile_type pt2 on pt2.profile_type_id = p2.profile_type_id and pt2.priority = 2
      ) on tp2.site_id = s.site_id
      LEFT JOIN profile_data pd ON pd.profile_id = p2.profile_id
        AND pd.data_element_id = 1
        AND pd.version = (
          SELECT MAX(d.version)
          FROM profile_data d
          WHERE d.data_element_id = 1
            AND d.profile_id = p2.profile_id
            AND d.approved = 1
        )
      LEFT JOIN profile_data pdn ON pdn.profile_id = p2.profile_id
        AND pdn.data_element_id = 3
        AND pdn.version = (
          SELECT MAX(d.version)
          FROM profile_data d
          WHERE d.data_element_id = 3
            AND d.profile_id = p2.profile_id
            AND d.approved = 1
        )
		  LEFT JOIN (
		    site_profiles tp3
					join profile p3 on p3.profile_id = tp3.profile_id
					join profile_type pt3 on pt3.profile_type_id = p3.profile_type_id and pt3.priority = 3
		  ) on tp3.site_id = s.site_id
			LEFT JOIN (
			  site_profiles tp4
          join profile p4 on p4.profile_id = tp4.profile_id
          join profile_type pt4 on pt4.profile_type_id = p4.profile_type_id and pt4.priority = 4
		  ) on tp4.site_id = s.site_id
    where ts.tid_id = t.tid_id
      and FIND_IN_SET(p4.name, ?)
      and (upper(t.tid_id) like ?
      or upper(t.serial) like ?
      or upper(pdn.datavalue) like ?)
    group by
      t.tid_id,
      t.serial,
      t.PIN,
      t.ExpiryDate,
      t.ActivationDate,
      s.site_id,
      pdn.datavalue,
      pd.datavalue,
      p3.profile_id,
      p3.name;`

	pedOverrideQuery = `select 
    de.name,
	dg.name,
		pd.datavalue
		from profile_data_group pdg
        	join data_group dg on dg.data_group_id = pdg.data_group_id
         	join data_element de on de.data_group_id = dg.data_group_id
         	left join profile_data pd on pd.data_element_id = de.data_element_id and pd.profile_id = pdg.profile_id
         	left join profile p ON p.profile_id = pd.profile_id
		where pdg.profile_id = ?;`
)

// FindBySearchTermAndAcquirer Finds and returns a slice of PEDDetailed for a given acquirer by given search terms
func (r *pedRepository) FindBySearchTermAndAcquirer(searchTerm string, acquirers string) ([]*PED.PEDDetailed, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	dataGroupInfo, err := GetDataGroupInfoForTIDExport()
	if err != nil {
		logging.Error(fmt.Sprintf("unable to fecth data group info :%s", err.Error()))
		return nil, err
	}

	searchTerm = "%" + strings.ToUpper(searchTerm) + "%"

	rows, err := db.Query(pedCommonQuery, acquirers, searchTerm, searchTerm, searchTerm)
	if err != nil {
		logging.Error(err)
		return nil, err
	}

	defer rows.Close()

	resultPedData := make([]*PED.PEDDetailed, 0)

	var profileId sql.NullInt64
	var tidId sql.NullInt64
	var serial sql.NullString
	var pin sql.NullString
	var expiryDate sql.NullString
	var activationDate sql.NullString
	var siteId sql.NullInt64
	var siteName sql.NullString
	var chainId sql.NullInt64
	var chainName sql.NullString
	var merchantId sql.NullString
	var appVer sql.NullInt64
	var firmwareVer sql.NullString
	var lastTransaction sql.NullString
	var lastCheckedTime sql.NullString
	var confirmedTime sql.NullString
	var lastAPKDownload sql.NullString
	var ipaddress sql.NullString
	var ipAddresses sql.NullString
	var simCardSerialNumber sql.NullString
	var flagStatus sql.NullBool
	var flaggedDate sql.NullString
	var eodAuto sql.NullBool
	var autoTime sql.NullString
	var pedCoordinates sql.NullString
	var pedCoordinatesAccuracy sql.NullString
	var pedCoordinatesLastUpdated sql.NullString
	var freeInternalStorage sql.NullString
	var totalInternalStorage sql.NullString
	var softuiLastDownloadedFileName sql.NullString
	var softuiLastDownloadedFileHash sql.NullString
	var softuiLastDownloadedFileList sql.NullString
	var softuiLastDownloadedFileDateTime sql.NullString

	for rows.Next() {
		rowPedData := new(PED.PEDDetailed)
		tidDataElementDetails := make(exporter.ExportableItem, 0)
		err := rows.Scan(&profileId, &tidId, &serial, &pin, &expiryDate, &activationDate,
			&siteId, &siteName, &chainId, &chainName, &merchantId, &appVer, &firmwareVer,
			&lastTransaction, &lastCheckedTime, &confirmedTime, &lastAPKDownload, &ipaddress, &ipAddresses,
			&simCardSerialNumber, &flagStatus, &flaggedDate, &eodAuto, &autoTime, &pedCoordinates, &pedCoordinatesAccuracy,
			&pedCoordinatesLastUpdated, &freeInternalStorage, &totalInternalStorage, &softuiLastDownloadedFileName, &softuiLastDownloadedFileHash,
			&softuiLastDownloadedFileList, &softuiLastDownloadedFileDateTime)
		if err != nil {
			logging.Error(err)
			return nil, err
		}

		// First we need to check the profile ID to see if we need to fetch Override data
		// If it is NULL then the PED does not have an overridden configuration
		if profileId.Valid {
			// Placing this nested query within an anonymous function to ensure we avoid resource leaks
			func() {
				overrideRows, err := db.Query(pedOverrideQuery, profileId.Int64)
				if err != nil {
					logging.Error(err)
					return
				}
				defer overrideRows.Close()

				var overrideMap = map[string]string{}

				for overrideRows.Next() {

					var elementName sql.NullString
					var elementGroup sql.NullString
					var elementData sql.NullString

					err = overrideRows.Scan(&elementName, &elementGroup, &elementData)
					if err != nil {
						logging.Error(err)
						return
					}

					// If the element data is null then we can move on
					if !elementData.Valid {
						continue
					}

					// Now we know the element data is valid we construct the JSON key This needs to be in the format element.Group.elementName
					// e.g. userMgmt.wifiPIN to match the JSON tags in PED.go
					elementKey := elementGroup.String + "." + elementName.String

					// Add our dataValue to the map using the combined group.element name as the key
					overrideMap[elementKey] = elementData.String
				}

				for _, dataElement := range dataGroupInfo {
					if _, ok := overrideMap[dataElement.DisplayName]; ok {
						exportDisplayIndex, err := strconv.Atoi(dataElement.ExportDisplayIndex)
						if err != nil {
							logging.Error(fmt.Sprintf("An error occurred coverting index value string to int, index '%s' : error :%s", dataElement.ExportDisplayIndex, err.Error()))
							return
						}

						tidDataElementDetails[exporter.ExportableItemHeader{
							exportDisplayIndex,
							dataElement.DisplayNameEn,
						}] = overrideMap[dataElement.DisplayName]
					} else {
						exportDisplayIndex, err := strconv.Atoi(dataElement.ExportDisplayIndex)
						if err != nil {
							logging.Error(fmt.Sprintf("An error occurred coverting index value string to int index '%s' : error :%s", dataElement.ExportDisplayIndex, err.Error()))
							return
						}

						tidDataElementDetails[exporter.ExportableItemHeader{
							exportDisplayIndex,
							dataElement.DisplayNameEn,
						}] = overrideMap[dataElement.DisplayName]
					}
				}
			}()

		}

		rowPedData.PEDInfo = tidDataElementDetails
		// Then we can add the data that is present for all PEDs
		rowPedData.TID = int(tidId.Int64)
		rowPedData.Serial = serial.String
		rowPedData.PIN = pin.String
		rowPedData.ExpiryTime = expiryDate.String
		rowPedData.ActivationTime = activationDate.String
		rowPedData.SiteId = int(siteId.Int64)
		rowPedData.SiteName = siteName.String
		rowPedData.ChainId = int(chainId.Int64)
		rowPedData.ChainName = chainName.String
		rowPedData.MerchantID = merchantId.String
		rowPedData.AppVer = int(appVer.Int64)
		rowPedData.FirmwareVer = firmwareVer.String
		rowPedData.LastTransaction = lastTransaction.String
		rowPedData.LastCheckedTime = lastCheckedTime.String
		rowPedData.ConfirmedTime = confirmedTime.String
		rowPedData.LastAPKDownloadTime = lastAPKDownload.String
		rowPedData.IPAddress = ipaddress.String
		rowPedData.IPAddresses = ipAddresses.String
		rowPedData.SIMCardSerialNumber = simCardSerialNumber.String
		rowPedData.PEDCoordinates = pedCoordinates.String
		rowPedData.PEDCoordinatesAccuracy = pedCoordinatesAccuracy.String
		rowPedData.PEDCoordinatesLastUpdated = pedCoordinatesLastUpdated.String
		rowPedData.FreeInternalStorage = freeInternalStorage.String
		rowPedData.TotalInternalStorage = totalInternalStorage.String
		rowPedData.SoftuiLastDownloadedFileDateTime = common.CheckStringIsValid(softuiLastDownloadedFileDateTime)
		rowPedData.SoftuiLastDownloadedFileList = common.CheckStringIsValid(softuiLastDownloadedFileList)
		rowPedData.SoftuiLastDownloadedFileHash = common.CheckStringIsValid(softuiLastDownloadedFileHash)
		rowPedData.SoftuiLastDownloadedFileName = common.CheckStringIsValid(softuiLastDownloadedFileName)
		rowPedData.FlagStatus = common.CheckBoolIsValid(flagStatus)
		rowPedData.FlaggedDate = common.CheckStringIsValid(flaggedDate)
		rowPedData.EODAuto = common.CheckBoolIsValid(eodAuto)
		rowPedData.AutoTime = common.CheckStringIsValid(autoTime)

		resultPedData = append(resultPedData, rowPedData)
	}

	return resultPedData, nil
}

// Finds the profileId for a given TID if an override exists, if no profile exists for the given tid then profileId
// will be 0.
func (r *pedRepository) findOverrideProfileIdByTid(tid string) (profileExists bool, profileId int, err error) {
	db, err := GetDB()
	if err != nil {
		return profileExists, profileId, err
	}

	rows, err := db.Query("SELECT tid_profile_id FROM tid_site WHERE tid_id = ?", tid)
	if err != nil {
		return profileExists, profileId, err
	}

	defer rows.Close()
	for rows.Next() {
		var nullableProfileId sql.NullInt32
		err = rows.Scan(&nullableProfileId)
		if err != nil {
			return profileExists, profileId, err
		}

		if nullableProfileId.Valid {
			profileId = int(nullableProfileId.Int32)
			profileExists = true
		}
	}
	return profileExists, profileId, err
}

// DeleteByTid Deletes a TID and any corresponding overrides if they exist
func (r *pedRepository) DeleteByTid(tid string) (tidDeleted bool, err error) {
	db, err := GetDB()
	if err != nil {
		return false, err
	}

	_, err = r.DeleteFraudOverrideByTid(tid)
	if err != nil {
		return false, err
	}

	_, err = r.DeleteOverrideByTid(tid)
	if err != nil {
		return false, err
	}

	_, err = r.DeleteUserOverrideByTid(tid)
	if err != nil {
		return false, err
	}

	logging.Information(fmt.Sprintf("Deleting TID '%v'", tid))

	// We handle the deletion of the tid_site and tid record within the same db txn.
	// Because the tid_site table record deletion occurs first, if the deletion of the record from the tid table fails
	// then we need to roll back the deletion.
	// If we don't roll it back then we will get orphaned TIDs in the tid table not linked to any site because the
	// tid_site record has been deleted.
	// If this happens then TID and SN validation will occur against TIDs which do not exist in the UI.
	txn, err := db.Begin()
	if err != nil {
		return false, err
	}

	result, err := txn.Exec("DELETE FROM tid_site WHERE tid_id = ?", tid)
	if err != nil {
		logging.Error(txn.Rollback())
		return false, err
	}

	if rowsAffected, err := result.RowsAffected(); err != nil {
		logging.Error(txn.Rollback())
		return false, err
	} else if rowsAffected == 0 {
		logging.Error(txn.Rollback())
		return false, err
	} else {
		tidDeleted = true
	}

	result, err = txn.Exec("DELETE FROM tid WHERE tid_id = ?", tid)
	if err != nil {
		logging.Error(txn.Rollback())
		return false, err
	}

	if rowsAffected, err := result.RowsAffected(); err != nil {
		logging.Error(txn.Rollback())
		return false, err
	} else if rowsAffected == 0 {
		logging.Error(txn.Rollback())
		return false, err
	} else {
		tidDeleted = true
	}

	//NEX-10276- Delete from tid updates which keeps the package details
	result, err = txn.Exec("DELETE FROM tid_updates WHERE tid_id = ?", tid)
	if err != nil {
		logging.Error(txn.Rollback())
		return false, err
	}

	if rowsAffected, err := result.RowsAffected(); err != nil {
		logging.Error(txn.Rollback())
		return false, err
	} else {
		logging.Information(strconv.Itoa(int(rowsAffected)) + " Rows removed from tid_updates")
		tidDeleted = true
	}

	err = txn.Commit()
	if err != nil {
		logging.Error(txn.Rollback())
		return false, err
	}

	return tidDeleted, err
}

func (r *pedRepository) DeleteFraudOverrideByTid(tid string) (fraudOverrideDeleted bool, err error) {
	overrideExists, profileId, err := r.findOverrideProfileIdByTid(tid)
	if err != nil {
		return fraudOverrideDeleted, err
	}

	if overrideExists {
		logging.Information(fmt.Sprintf("Deleting tid fraud override for tid '%v' with profileId '%v'", tid, profileId))

		tidInt, err := strconv.Atoi(tid)
		if err != nil {
			return false, err
		}

		siteId, err := GetSiteIDForTid(tidInt)
		if err != nil {
			return false, err
		}

		err = DeleteSiteVelocityLimits(siteId, 4, tidInt)
		if err != nil {
			return false, err
		}

		//Delete the scheme velocity limits
		err = DeleteSiteVelocityLimits(siteId, 2, tidInt)
		if err != nil {
			return false, err
		}

		db, err := GetDB()
		if err != nil {
			return false, err
		}

		// Delete any tid config overrides relating to fraud.
		_, err = db.Exec("CALL remove_tid_fraud_config_override(?)", profileId)
		if err != nil {
			return false, err
		}
		return true, err
	} else {
		return false, nil
	}
}

// DeleteOverrideByTid Deletes the overrides of a given TID
func (r *pedRepository) DeleteOverrideByTid(tid string) (overrideDeleted bool, err error) {
	overrideExists, profileId, err := r.findOverrideProfileIdByTid(tid)
	if err != nil {
		return overrideDeleted, err
	}

	if overrideExists {
		db, err := GetDB()
		if err != nil {
			return false, err
		}

		logging.Information(fmt.Sprintf("Deleting tid override for TID '%v' with profileId '%v'", tid, profileId))

		result, err := db.Exec("CALL remove_tid_override(?)", profileId)
		if err != nil {
			return false, err
		}

		if rowsAffected, err := result.RowsAffected(); err != nil {
			return false, err
		} else {
			overrideDeleted = rowsAffected != 0
		}

		return overrideDeleted, err
	} else {
		return false, nil
	}
}

func (r *pedRepository) DeleteUserOverrideByTid(tid string) (userOverrideDeleted bool, err error) {
	db, err := GetDB()
	if err != nil {
		return
	}

	result, err := db.Exec("DELETE FROM tid_user_override WHERE tid_id = ?", tid)
	if err != nil {
		return
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return
	}

	userOverrideDeleted = rowsAffected != 0
	return
}

func (r *pedRepository) DeleteFraudOverrideByTids(tids []string, siteId int) (fraudOverrideDeleted bool, err error) {

	tidIDsString := strings.Join(tids, ",")

	overrideExists, profileId, err := r.findOverrideProfileIdByTids(tids)
	if err != nil {
		return fraudOverrideDeleted, err
	}

	if overrideExists {
		db, err := GetDB()
		if err != nil {
			return false, err
		}

		logging.Information(fmt.Sprintf("Deleting site fraud override for tids '%v' with profileId '%v'", tids, profileId))

		result, err := db.Exec("CALL remove_tid_override(?)", profileId)
		if err != nil {
			return false, err
		}

		if rowsAffected, err := result.RowsAffected(); err != nil {
			return false, err
		} else {
			fraudOverrideDeleted = rowsAffected != 0
		}

		// Delete any tid config overrides relating to fraud.
		_, err = db.Exec("CALL delete_velocity_limit_overide(?)", tidIDsString, siteId)
		if err != nil {
			return false, err
		}

		// Delete any tid config overrides relating to fraud.
		_, err = db.Exec("CALL remove_tid_fraud_config_override(?)", profileId)
		if err != nil {
			return false, err
		}
		return fraudOverrideDeleted, err
	} else {
		return false, nil
	}
}

func (r *pedRepository) findOverrideProfileIdByTids(tids []string) (profileExists bool, profileId int, err error) {
	db, err := GetDB()
	if err != nil {
		return profileExists, profileId, err
	}

	rows, err := db.Query("SELECT tid_profile_id FROM tid_site WHERE tid_id IN (" + strings.Repeat("?,", len(tids)-1) + "?)")
	if err != nil {
		return profileExists, profileId, err
	}

	defer rows.Close()
	for rows.Next() {
		var nullableProfileId sql.NullInt32
		err = rows.Scan(&nullableProfileId)
		if err != nil {
			return profileExists, profileId, err
		}

		if nullableProfileId.Valid {
			profileId = int(nullableProfileId.Int32)
			profileExists = true
		}
	}
	return profileExists, profileId, err
}

func NewPEDRepository() PED.Repository {
	return new(pedRepository)
}
