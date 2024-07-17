package dal

import (
	"database/sql"
	pd "nextgen-tms-website/ProfileData"
	"nextgen-tms-website/entities"
	"sync"
)

type profileDataRepository struct{}

func NewProfileDataRepository() pd.Repository {
	return new(profileDataRepository)
}

var mutex = sync.Mutex{}

func (p *profileDataRepository) SetDataValueByElementIdAndProfileIdWithoutApproval(dataElementId, profileId int, value string, user entities.TMSUser) error {

	mutex.Lock()
	defer mutex.Unlock()
	db, err := GetDB()
	if err != nil {
		logging.Error(err)
		return err
	}

	dataExists, _, err := p.getDataValueByElementIdAndProfileId(dataElementId, profileId)
	if err != nil {
		logging.Error(err)
		return err
	}

	if dataExists {
		_, err = db.Exec("UPDATE profile_data SET datavalue = ?, updated_at = NOW(), updated_by = ? WHERE data_element_id = ? AND profile_id = ?",
			value, user.Username, dataElementId, profileId)
	} else {
		_, err = db.Exec(`INSERT INTO profile_data 
    							(profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden)
    							VALUES
    							(?, ?, ?, 1, NOW(), ?, NOW(), ?, 1, 0)`,
			profileId, dataElementId, value, user.Username, user.Username)
	}
	if err != nil {
		logging.Error(err)
		return err
	}
	return nil
}

func (p *profileDataRepository) getDataValueByElementIdAndProfileId(dataElementId, profileId int) (dataElementExists bool, value string, err error) {
	db, err := GetDB()
	if err != nil {
		logging.Error(err)
		return
	}

	rows, err := db.Query("SELECT datavalue FROM profile_data WHERE data_element_id = ? AND profile_id = ?", dataElementId, profileId)
	if err != nil {
		logging.Error(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var nullableVal sql.NullString
		err = rows.Scan(&nullableVal)
		if err != nil {
			logging.Error(err)
			return
		}

		dataElementExists = true
		if nullableVal.Valid {
			value = nullableVal.String
		}
	}
	return
}
