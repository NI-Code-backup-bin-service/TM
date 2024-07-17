package dal

import (
	dg "nextgen-tms-website/DataGroup"
)

type dataGroupRepository struct{}

// Returns all data groups and if they're selected or preselected for a given site profile ID
func (r *dataGroupRepository) FindForSiteByProfileId(siteProfileId int) ([]dg.DataGroup, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	sql :=
		// Left join all data groups with the data groups we can inner join with
		// the site in question (which is where the main bulk of the query logic is).
		// This avoids the need to perform left outer joins within the main query logic
		// (i.e. within 'dgs_for_site') which is very expensive (see NEX-6806).
		`SELECT all_dgs.data_group_id, all_dgs.name, all_dgs.displayname_en,
			IFNULL(dgs_for_site.selected, false) as selected,
			IFNULL(dgs_for_site.preselected, false) as preselected
   		FROM data_group all_dgs
		LEFT JOIN (
			SELECT dg.data_group_id,
				dg.name,
				dg.displayname_en,
				IF(COUNT(sp.profile_id) > 0, true, false) selected,
				IF(COUNT(sp.profile_id) > 1, true, false) preselected
			FROM data_group dg
			INNER JOIN profile_data_group pdg ON dg.data_group_id = pdg.data_group_id
			INNER JOIN (
				SELECT sp2.profile_id
				FROM site_profiles sp
				INNER JOIN site_profiles sp2 ON sp.site_id = sp2.site_id
				WHERE sp.profile_id = ?
				  # We don't care what data groups global has selected
				  AND sp2.profile_id != 1
			) sp ON sp.profile_id = pdg.profile_id
			GROUP BY dg.data_group_id
		) dgs_for_site on all_dgs.data_group_id = dgs_for_site.data_group_id;`

	rows, err := db.Query(sql, siteProfileId)
	if err != nil {
		logging.Error(err)
		return nil, err
	}

	defer rows.Close()

	dataGroups := make([]dg.DataGroup, 0)
	for rows.Next() {
		var dataGroup dg.DataGroup
		err := rows.Scan(&dataGroup.ID, &dataGroup.Name, &dataGroup.DisplayName, &dataGroup.Selected, &dataGroup.Preselected)
		if err != nil {
			logging.Error(err)
			return nil, err
		}
		dataGroups = append(dataGroups, dataGroup)
	}
	return dataGroups, nil
}

// Returns all data groups and if they're selected or preselected for a given TID profile ID
func (r *dataGroupRepository) FindForTidByProfileId(tidProfileId int) ([]dg.DataGroup, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	sql :=
		// Left join all data groups with the data groups we can inner join with
		// the site in question (which is where the main bulk of the query logic is).
		// This avoids the need to perform left outer joins within the main query logic
		// (i.e. within 'dgs_for_tid') which is very expensive (see NEX-6806).
		`SELECT all_dgs.data_group_id, all_dgs.name, all_dgs.displayname_en,
			IFNULL(dgs_for_tid.selected, false) as selected,
			IFNULL(dgs_for_tid.preselected, false) as preselected
   		FROM data_group all_dgs
		LEFT JOIN (
			SELECT dg.data_group_id,
				dg.name,
				dg.displayname_en,
				IF(COUNT(sp.profile_id) > 0, true, false) selected,
				IF(COUNT(sp.profile_id) > 1, true, false) preselected
			FROM data_group dg
			INNER JOIN profile_data_group pdg ON dg.data_group_id = pdg.data_group_id
			INNER JOIN (
				SELECT
					sp2.profile_id
				FROM
					site_profiles sp
					INNER JOIN site_profiles sp2 ON sp.site_id = sp2.site_id
					INNER JOIN tid_site ts ON sp.site_id = ts.site_id
				WHERE
					ts.tid_profile_id = ?
					-- We don't care what data groups global has selected
					AND sp2.profile_id != 1
			) sp ON sp.profile_id = pdg.profile_id
			GROUP BY dg.data_group_id
		) dgs_for_tid on all_dgs.data_group_id = dgs_for_tid.data_group_id;`

	rows, err := db.Query(sql, tidProfileId)
	if err != nil {
		logging.Error(err)
		return nil, err
	}

	defer rows.Close()

	dataGroups := make([]dg.DataGroup, 0)
	for rows.Next() {
		var dataGroup dg.DataGroup
		err := rows.Scan(&dataGroup.ID, &dataGroup.Name, &dataGroup.DisplayName, &dataGroup.Selected, &dataGroup.Preselected)
		if err != nil {
			logging.Error(err)
			return nil, err
		}
		dataGroups = append(dataGroups, dataGroup)
	}
	return dataGroups, nil
}

// Returns the DataGroup an element belongs to by the element ID
func (r *dataGroupRepository) FindByDataElementId(elementId int) (dg.DataGroup, error) {
	var dataGroup dg.DataGroup
	db, err := GetDB()
	if err != nil {
		return dataGroup, err
	}

	rows, err := db.Query(`SELECT dg.data_group_id, dg.name, dg.displayname_en
								FROM data_element de
								INNER JOIN data_group dg ON de.data_group_id = dg.data_group_id
								WHERE de.data_element_id = ?`, elementId)
	if err != nil {
		logging.Error(err)
		return dataGroup, err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&dataGroup.ID, &dataGroup.Name, &dataGroup.DisplayName)
		if err != nil {
			logging.Error(err)
			return dataGroup, err
		}
	}
	return dataGroup, nil
}

func NewDataGroupRepository() dg.Repository {
	return new(dataGroupRepository)
}
