package dal

import (
	"database/sql"
	"nextgen-tms-website/entities"
	"regexp"
	"strings"

	"github.com/rs/xid"
)

// Deletes the site velocity limits
func DeleteSiteVelocityLimits(siteId int, level int, tidId int) error {

	db, err := GetDB()
	if err != nil {
		return err
	}

	//Need to delete the nested txn velocity limits first due to use of foreign key
	_, err = db.Exec("CALL delete_txn_velocity_limit(?, ?, ?)", siteId, level, tidId)
	if err != nil {
		return err
	}

	//Now delete the site's velocity limits
	_, err = db.Exec("CALL delete_velocity_limit(?, ?, ?)", siteId, level, tidId)
	if err != nil {
		return err
	}

	return nil
}

func DeleteSiteVelocityLimitsAndTxn(siteId, schemeLimitLevel, limitLevel, tidId int) error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	_, err = db.Exec("CALL delete_velocity_and_txn_velocity_limit(?, ?, ?, ?)", siteId, schemeLimitLevel, limitLevel, tidId)
	if err != nil {
		return err
	}
	return nil
}

func SaveProfileDataElementByName(profileId int, elementName, elementValue string) error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	_, err = db.Exec("CALL save_profile_data_element_by_name(?, ?, ?)", profileId, elementName, elementValue)
	if err != nil {
		return err
	}
	return nil
}

// Saves the site velocity limits
func SetSiteVelocityLimits(siteID int, tidID int, limitList []entities.VelocityLimit) error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	for _, velocity := range limitList {
		id := velocity.ID
		limitLevel := velocity.Level
		scheme := velocity.Scheme
		dailyTL := velocity.DailyCount
		batchTL := velocity.BatchCount
		singleTL := velocity.SingleTransLimit
		dailyLimit := velocity.DailyLimit
		batchLimit := velocity.BatchLimit
		index := velocity.Index
		if _, err = db.Exec("CALL set_velocity_limit(?,?,?,?,?,?,?,?,?,?,?)", id, siteID, tidID, limitLevel, scheme, dailyTL, batchTL, singleTL, dailyLimit, batchLimit, index); err != nil {
			return err
		}
		for _, txnVelocity := range velocity.TxnLimits {
			txnId := txnVelocity.TxnLimitID
			limitType := txnVelocity.LimitType
			txnType := txnVelocity.TxnType
			value := txnVelocity.Value
			if _, err = db.Exec("CALL set_txn_velocity_limit(?,?,?,?,?)", txnId, id, limitType, txnType, value); err != nil {
				return err
			}
		}
	}
	return nil
}

// Retrieves an array of velocity limits set for a given site
func GetSiteVelocityLimits(siteID int, tidID int) (entities.VelocityLimit, []entities.VelocityLimit, error) {

	//Creating a dummy struct here as when an error occurs 'something' has to be sent back, can't be nil
	siteLevelLimits := entities.VelocityLimit{
		ID:               "",
		Scheme:           "",
		DailyCount:       -1,
		DailyLimit:       -1,
		BatchCount:       -1,
		BatchLimit:       -1,
		SingleTransLimit: -1,
		TxnLimits:        nil,
		Index:            0,
	}

	db, err := GetDB()
	if err != nil {
		return siteLevelLimits, nil, err
	}

	var rowcount int
	//Check to see if TID overrides have been set
	err = db.QueryRow("SELECT COUNT(*) FROM velocity_limits WHERE tid_id = ? AND limit_level = 4", tidID).Scan(&rowcount)

	//Determines whether the TID is using site level limits or not
	inheritedLimits := false

	var siteLimits []entities.VelocityLimit
	var siteRows *sql.Rows

	//If the TID has no overrides set, make it use the Site velocity limits
	if rowcount < 1 {
		inheritedLimits = true
		siteRows, err = db.Query("CALL get_ordered_site_velocity_limits(?, ?, ?)", siteID, -1, 3)
	} else {
		siteRows, err = db.Query("CALL get_ordered_site_velocity_limits(?, ?, ?)", siteID, tidID, 4)
	}

	if err != nil {
		return siteLevelLimits, nil, err
	}
	defer siteRows.Close()

	for siteRows.Next() {

		var dailyLimit sql.NullInt32
		var batchLimit sql.NullInt32

		err = siteRows.Scan(&siteLevelLimits.ID, &siteLevelLimits.Scheme, &siteLevelLimits.DailyCount, &siteLevelLimits.BatchCount, &siteLevelLimits.SingleTransLimit, &dailyLimit, &batchLimit, &siteLevelLimits.Index)
		if err != nil {
			return siteLevelLimits, nil, err
		}

		if dailyLimit.Valid {
			siteLevelLimits.DailyLimit = int(dailyLimit.Int32)
		} else {
			siteLevelLimits.DailyLimit = -1
		}

		if batchLimit.Valid {
			siteLevelLimits.BatchLimit = int(batchLimit.Int32)
		} else {
			siteLevelLimits.BatchLimit = -1
		}

		txnLimits, err := getTxnVelocityLimits(siteLevelLimits.ID)
		if err != nil {
			return siteLevelLimits, nil, err
		}

		if inheritedLimits {
			siteLevelLimits.ID = newID()
			for _, txnLimit := range txnLimits {
				txnLimit.TxnLimitID = newID()
			}
		}

		siteLevelLimits.TxnLimits = txnLimits
	}

	var rows *sql.Rows

	//If the TID has no overrides set, make it use the Site velocity limits
	if rowcount < 1 {
		inheritedLimits = true
		rows, err = db.Query("CALL get_ordered_site_velocity_limits(?, ?, ?)", siteID, -1, 1)
	} else {
		rows, err = db.Query("CALL get_ordered_site_velocity_limits(?, ?, ?)", siteID, tidID, 2)
	}

	if err != nil {
		return siteLevelLimits, nil, err
	}
	defer rows.Close()

	for rows.Next() {

		var limits entities.VelocityLimit

		err = rows.Scan(&limits.ID, &limits.Scheme, &limits.DailyCount, &limits.BatchCount, &limits.SingleTransLimit, &limits.DailyLimit, &limits.BatchLimit, &limits.Index)
		if err != nil {
			return siteLevelLimits, nil, err
		}

		txnLimits, err := getTxnVelocityLimits(limits.ID)
		if err != nil {
			return siteLevelLimits, nil, err
		}

		if inheritedLimits {
			limits.ID = newID()
			for _, txnLimit := range txnLimits {
				txnLimit.TxnLimitID = newID()
			}
		}

		limits.TxnLimits = txnLimits
		siteLimits = append(siteLimits, limits)
	}

	return siteLevelLimits, siteLimits, nil
}

func getTxnVelocityLimits(limitID string) ([]entities.TxnLimit, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("CALL get_txn_velocity_limits(?)", limitID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txnLimits []entities.TxnLimit
	for rows.Next() {
		var limit entities.TxnLimit
		err = rows.Scan(&limit.TxnLimitID, &limit.LimitType, &limit.TxnType, &limit.TxnTypeReadable, &limit.Value)
		if err != nil {
			return nil, err
		}
		txnLimits = append(txnLimits, limit)
	}

	return txnLimits, nil
}

// Retrieves the available transaction types from DB
func GetAvailableTransactions() ([]entities.VelocityTransactions, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("CALL get_velocity_transactions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var transactions []entities.VelocityTransactions
	for rows.Next() {
		var transaction entities.VelocityTransactions
		err = rows.Scan(&transaction.TxnType, &transaction.TxnTypeReadable)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

// Retrieves the available limit types from DB
func GetAvailableLimits() ([]entities.VelocityLimitTypes, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("CALL get_velocity_limit_types")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var limits []entities.VelocityLimitTypes
	for rows.Next() {
		var limit entities.VelocityLimitTypes
		err = rows.Scan(&limit.Limittype)
		if err != nil {
			return nil, err
		}
		limit.Identifier, err = formatLimitIdentifier(limit.Limittype)
		if err != nil {
			return nil, err
		}
		limits = append(limits, limit)
	}
	return limits, nil
}

// Sets string to lower case, replaces all non alpha-numeric/non bracket characters with hyphens
func formatLimitIdentifier(limitType string) (string, error) {
	//Make all lower case
	limitIdentifier := strings.ToLower(limitType)

	//Strip the whitespace and replace with dashes
	reg, err := regexp.Compile("[^a-zA-Z0-9()]+")
	if err != nil {
		return "", err
	}
	return reg.ReplaceAllString(limitIdentifier, "-"), nil
}

func newID() string {
	id := xid.New()
	return id.String()
}
