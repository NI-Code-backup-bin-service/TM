package dal

import (
	"database/sql"
	"html"
	"nextgen-tms-website/entities"
	"strings"
)

// GetAcquirerList Method for searching acquirers
func GetAcquirerList(searchTerm string, user *entities.TMSUser, offset, amount int) ([]*AcquirerList, int, error) {
	db, err := GetDB()
	if err != nil {
		return nil, -1, err
	}

	// Find user acquirers to limit search results
	acquirers, err := GetUserAcquirerPermissions(user)
	if err != nil {
		return nil, -1, err
	}
	if amount == 0 {
		amount = len(strings.Split(acquirers, ","))
	}
	filterArg := "%" + searchTerm + "%"
	rows, err := db.Query("CALL get_aquirer_page(?,?,?,?)", filterArg, acquirers, amount, offset)

	if err != nil {
		return nil, -1, err
	}
	defer rows.Close()

	var acquirerList []*AcquirerList

	for rows.Next() {
		var acquirer AcquirerList
		err = rows.Scan(
			&acquirer.AcquirerProfileID,
			&acquirer.AcquirerName)

		if err != nil {
			return nil, -1, err
		}

		acquirer.AcquirerName = html.EscapeString(acquirer.AcquirerName)
		acquirerList = append(acquirerList, &acquirer)
	}

	var totalAcquirers int
	if db.QueryRow(`CALL get_aquirer_count(?,?)`, filterArg, acquirers).Scan(&totalAcquirers); err != nil {
		return nil, -1, err
	}
	return acquirerList, totalAcquirers, nil
}

func GetAcquirerIdFromChainId(chainId int) (int, error) {
	db, err := GetDB()
	if err != nil {
		return -1, err
	}

	row := db.QueryRow("select acquirer_id from chain_profiles where chain_profile_id = ?", chainId)
	var id int
	err = row.Scan(&id)
	if err != nil {
		logging.Error(err)
		return -1, err
	}
	return id, nil
}

func GetAcquirerNameFromChainProfileId(profileId int) (string, error) {
	db, err := GetDB()
	if err != nil {
		logging.Error(err)
		return "", err
	}

	row := db.QueryRow("SELECT p.name FROM profile p JOIN chain_profiles cp ON p.profile_id = cp.acquirer_id WHERE cp.chain_profile_id = ?", profileId)
	var acquirer sql.NullString
	err = row.Scan(&acquirer)
	if err != nil {
		logging.Error(err)
		return "", err
	}

	return acquirer.String, err
}

func GetAcquirerNameFromAcquirerProfileId(profileId int) (string, error) {
	db, err := GetDB()
	if err != nil {
		logging.Error(err)
		return "", err
	}

	row := db.QueryRow("select name from profile where profile_id = ?", profileId)
	var acquirer sql.NullString
	err = row.Scan(&acquirer)
	if err != nil {
		logging.Error(err)
		return "", err
	}
	return acquirer.String, err
}

func GetAcquirerFromAcquirerName(name string) (int, error) {
	db, err := GetDB()
	if err != nil {
		logging.Error(err)
		return -1, err
	}

	row := db.QueryRow("select profile_id from profile where name = ? AND profile_type_id = 2", name)
	var acquirer int
	err = row.Scan(&acquirer)
	if err != nil {
		logging.Error(err)
		return -1, err
	}
	return acquirer, err
}
