package main

import (
	"database/sql"
	"fmt"
	"nextgen-tms-website/crypt"
	"nextgen-tms-website/dal"
)

func EncryptEverything() {
	// Encrypt EVERYTHING!

	db, err := dal.GetDB()

	if err != nil {
		logging.Error(err)
		return
	}

	err = EncryptTIDUserOverride(db)
	if err != nil {
		return
	}

	err = EncryptSiteUserOverride(db)
	if err != nil {
		return
	}

	err = EncryptProfileData(db)
	if err != nil {
		return
	}

	err = EncryptApprovals(db)
	if err != nil {
		return
	}
	}

func EncryptTIDUserOverride(db *sql.DB) error {
	logging.Information("EncryptTIDUserOverride Started")

	rows, err := db.Query("SELECT tid_user_id, PIN FROM tid_user_override WHERE is_encrypted=0")
	if err != nil {
		logging.Error(err)
		return err
	}
	defer rows.Close()

	var PIN sql.NullString
	var tid_id sql.NullString

	for rows.Next() {
		err = rows.Scan(&tid_id, &PIN)

		if err != nil {
			logging.Error(err)
			return err
		}

		if tid_id.Valid && PIN.Valid {
			_, err = db.Exec("UPDATE tid_user_override SET is_encrypted=1, PIN=? WHERE tid_user_id=? LIMIT 1;\n",crypt.Encrypt(PIN.String),tid_id.String )

			if err != nil {
				logging.Error(err)
				return err
			}
		}
	}

	logging.Information("EncryptTIDUserOverride Completed")

	return nil
}

///

func EncryptSiteUserOverride(db *sql.DB) error {
	logging.Information("EncryptSiteUserOverride Started")

	rows, err := db.Query("SELECT user_id, PIN FROM site_level_users WHERE is_encrypted=0")
	if err != nil {
		logging.Error(err)
		return err
	}
	defer rows.Close()

	var PIN sql.NullString
	var user_id sql.NullString

	for rows.Next() {
		err = rows.Scan(&user_id, &PIN)

		if err != nil {
			logging.Error(err)
			return err
		}

		if user_id.Valid && PIN.Valid {
			_, err = db.Exec("UPDATE site_level_users SET is_encrypted=1, PIN=? WHERE user_id=? LIMIT 1;\n",crypt.Encrypt(PIN.String),user_id.String )

			if err != nil {
				logging.Error(err)
				return err
			}
		}
	}

	logging.Information("EncryptSiteUserOverride Completed")

	return nil
}

///

func EncryptProfileData(db *sql.DB) error {
	logging.Information("EncryptProfileData Started")

	rowsGRP, err := db.Query("SELECT GROUP_CONCAT(`data_element_id`) FROM data_element WHERE `is_password` = 1")
	if err != nil {
		logging.Error(err)
		return err
	}
	defer rowsGRP.Close()

	var group_data_element_ids sql.NullString
	for rowsGRP.Next() {
		err = rowsGRP.Scan(&group_data_element_ids)
		if err != nil {
			logging.Error(err)
			return err
		}
	}

	if !group_data_element_ids.Valid {
		return nil // no password fields
	}

	rowsGRP.Close()
	dataElementIds := fmt.Sprintf("%v", group_data_element_ids.String)
	rows, err := db.Query("SELECT profile_data_id, datavalue FROM profile_data WHERE (is_encrypted=0 OR ISNULL(is_encrypted)) && data_element_id IN (?)", dataElementIds)
	if err != nil {
		logging.Error(err)
		return err
	}
	defer rows.Close()

	var datavalue sql.NullString
	var profile_data_id sql.NullString

	for rows.Next() {
		err = rows.Scan(&profile_data_id, &datavalue)

		if err != nil {
			logging.Error(err)
			return err
		}

		if profile_data_id.Valid && datavalue.Valid {
			_, err = db.Exec("UPDATE profile_data SET is_encrypted=1, datavalue=? WHERE profile_data_id=? LIMIT 1;\n",crypt.Encrypt(datavalue.String),profile_data_id.String )

			if err != nil {
				logging.Error(err)
				return err
			}
		}
	}

	logging.Information("EncryptProfileData Completed")
	return nil
}

///// Approvals
func EncryptApprovals(db *sql.DB) error {
	logging.Information("EncryptApprovals Started")

	rows, err := db.Query("SELECT approval_id, current_value, new_value FROM approvals WHERE is_encrypted=0 AND is_password=1")
	if err != nil {
		logging.Error(err)
		return err
	}
	defer rows.Close()

	var approval_id sql.NullString
	var current_value sql.NullString
	var new_value sql.NullString

	for rows.Next() {
		err = rows.Scan(&approval_id, &current_value, &new_value)

		if err != nil {
			logging.Error(err)
			return err
		}

		if approval_id.Valid {
			_, err = db.Exec("UPDATE approvals SET is_encrypted=1, current_value=?, new_value=? WHERE approval_id=? LIMIT 1;",crypt.Encrypt(current_value.String), crypt.Encrypt(new_value.String), approval_id.String)

			if err != nil {
				logging.Error(err)
				return err
			}
		}
	}

	logging.Information("EncryptApprovals Completed")
	return nil
}
