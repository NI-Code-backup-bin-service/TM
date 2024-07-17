package dal

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

const (
	GET_SITE_ID_FROM_MID             = "Select profile_id FROM profile_data WHERE datavalue = ? AND data_element_id= (select data_element_id from data_element where name = 'merchantNo')"
	GET_PROFILE_ID_FROM_TID          = "SELECT tid_profile_id FROM tid_site WHERE tid_id = ?;"
	GET_DISPLAY_NAME_FROM_ELEMENT_ID = "Select displayname_en FROM data_element WHERE data_element_id = ?"
	GET_CHAIN_ID_FROM_SITE_ID        = "SELECT sp.profile_id FROM site_profiles sp JOIN profile p ON p.profile_id = sp.profile_id WHERE sp.site_id = (SELECT site_id FROM site_profiles WHERE profile_id = ?) AND p.profile_type_id = 3;"
	// GET_DATA_GROUPS_FROM_MID - To get the correct/single profileID against the dataValue/MID & data_element_id = 1 i.e., merchantNo
	GET_DATA_GROUPS_FROM_MID = "SELECT dg.name FROM data_group dg JOIN profile_data_group pdg ON dg.data_group_id = pdg.data_group_id WHERE pdg.profile_id = (SELECT profile_id FROM profile_data WHERE datavalue = ? AND data_element_id = (select data_element_id from data_element where name = 'merchantNo'));"
	PUSHPAYMENTS             = "pushpayments"
)

func FetchTidDataGroups(tid int) ([]string, error) {
	var dataGroups []string
	db, err := GetDB()
	if err != nil {
		return dataGroups, err
	}

	rows, err := db.Query("CALL get_data_groups_from_tid(?)", tid)
	if err != nil {
		return dataGroups, err
	}
	defer rows.Close()

	for rows.Next() {
		var dataGroup string
		err := rows.Scan(&dataGroup)
		if err != nil {
			return nil, err
		}

		dataGroups = append(dataGroups, dataGroup)
	}

	return dataGroups, nil
}

// retrieve active data group id's with tid overrideable enabled
func GetTidOveridableActiveDataGroupIds(siteId int) ([]string, error) {
	var dataGroupIds []string
	db, err := GetDB()
	if err != nil {
		return dataGroupIds, err
	}

	rows, err := db.Query("CALL get_tid_overridable_active_data_group_ids(?)", siteId)
	if err != nil {
		return dataGroupIds, err
	}
	defer rows.Close()

	for rows.Next() {
		var dataGroupId string
		err := rows.Scan(&dataGroupId)
		if err != nil {
			return nil, err
		}
		dataGroupIds = append(dataGroupIds, dataGroupId)
	}

	return dataGroupIds, nil
}
func FetchDataGroupsWithTIDOverRideEnabledAtDataElement(activeGroupIds []string) ([]string, error) {
	var dataGroups []string
	db, err := GetDB()
	if err != nil {
		return dataGroups, err
	}

	queryArgs := make([]interface{}, len(activeGroupIds))
	for i, tid := range activeGroupIds {
		queryArgs[i] = tid
	}

	var sb strings.Builder
	sb.WriteString("SELECT DISTINCT data_group_id FROM data_element WHERE tid_overridable=1 AND data_group_id IN (?")
	sb.WriteString(strings.Repeat(",?", len(activeGroupIds)-1))
	sb.WriteString(")")
	queryStr := sb.String()
	rows, err := db.Query(queryStr, queryArgs...)
	if err != nil {
		return dataGroups, err
	}

	defer rows.Close()

	for rows.Next() {
		var dataGroup string
		err := rows.Scan(&dataGroup)
		if err != nil {
			return nil, err
		}
		dataGroups = append(dataGroups, dataGroup)
	}

	return dataGroups, nil
}

func FetchSiteDataGroups(mid string) ([]string, error) {
	var dataGroups []string

	siteprofileId, _ := GetProfileIdFromMID(mid)
	if siteprofileId == "" {
		return dataGroups, nil
	}

	siteId, err := strconv.Atoi(siteprofileId)
	if err != nil {
		logging.Error(fmt.Sprintf("Error retrieving data groups for site id %s - %s", siteprofileId, err.Error()))
		return dataGroups, err
	}
	groups, err := NewDataGroupRepository().FindForSiteByProfileId(siteId)
	if err != nil {
		logging.Error(fmt.Sprintf("Error retrieving data groups for site id %s - %s", siteprofileId, err.Error()))
		return dataGroups, err
	}

	for _, group := range groups {
		if group.Selected {
			dataGroups = append(dataGroups, group.Name)
		}
	}

	return dataGroups, nil
}

/*
*
Retrieves the profile ID of the parent chain for the provided site profile ID
@param id - The profile ID of the site
*/
func GetChainIdFromSiteId(id int) (int, error) {
	db, err := GetDB()
	if err != nil {
		return 0, err
	}

	// Find the chain ID for the supplied site
	rows, err := db.Query(GET_CHAIN_ID_FROM_SITE_ID, id)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var chainId sql.NullInt64
	for rows.Next() {
		err = rows.Scan(&chainId)
	}

	if err != nil {
		return 0, err
	}

	if chainId.Valid {
		return int(chainId.Int64), nil
	}

	return 0, nil
}

func CheckDataGroupExistsProfile(profileId int, dataGroupName string) (bool, error) {
	db, err := GetDB()
	if err != nil {
		return false, err
	}

	var count sql.NullInt64
	err = db.QueryRow("CALL check_data_group_exist_in_profile(?, ?)", profileId, dataGroupName).Scan(&count)
	if err != nil {
		return false, err
	}

	if count.Valid {
		return count.Int64 > 0, nil
	}

	return false, nil
}

/*
*
Retrieves the display name for the provided data element
@param id - The element ID of the data element to be named
*/
func GetElementDisplayNameFromId(id int) (string, error) {
	db, err := GetDB()
	if err != nil {
		return "", err
	}

	// Finds the display name for the supplied data element id
	rows, err := db.Query(GET_DISPLAY_NAME_FROM_ELEMENT_ID, id)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var profileId sql.NullString
	for rows.Next() {
		err = rows.Scan(&profileId)
	}

	if err != nil {
		return "", err
	}

	if profileId.Valid {
		return profileId.String, nil
	}

	return "", nil
}

/*
*
Retrieves the Profile ID for the supplied MID
@Param mid - the MID value for the site in question
*/
func GetProfileIdFromMID(mid string) (string, error) {
	db, err := GetDB()
	if err != nil {
		return "", err
	}

	// Finds the supplied MID and retrieves the site profile ID for it
	rows, err := db.Query(GET_SITE_ID_FROM_MID, mid)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var profileId string
	for rows.Next() {
		err = rows.Scan(&profileId)
	}

	if err != nil {
		return "", err
	}

	return profileId, nil
}

func GetUniqueFields() ([]string, []bool, error) {
	db, err := GetDB()
	if err != nil {
		return nil, nil, err
	}

	// Select the name and whether that field is allowed to be empty
	rows, err := db.Query("SELECT name, is_allow_empty FROM data_element WHERE data_element.unique = 1;")
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var elements []string
	var is_allowed_empty []bool
	var tempName string
	var tempBool bool
	for rows.Next() {
		err = rows.Scan(&tempName, &tempBool)
		if err != nil {
			return nil, nil, err
		}
		// Adding the data to arrays to be returned
		elements = append(elements, tempName)
		is_allowed_empty = append(is_allowed_empty, tempBool)
	}
	return elements, is_allowed_empty, nil
}

func GetProfileIdFromTID(tid string) (int64, error) {
	db, err := GetDB()
	if err != nil {
		return 0, err
	}

	rows, err := db.Query(GET_PROFILE_ID_FROM_TID, tid)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var profileId sql.NullInt64
	for rows.Next() {
		err = rows.Scan(&profileId)
	}

	if err != nil {
		return 0, err
	}

	return profileId.Int64, nil
}
