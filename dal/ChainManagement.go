package dal

import (
	"database/sql"
	"fmt"
	"html"
	"nextgen-tms-website/crypt"
	"nextgen-tms-website/entities"
	"sort"
	"strconv"
	"strings"
)

func GetChainData(profileId int) ([]*DataGroup, []*DataGroup, []*DataGroup, error) {
	db, err := GetDB()
	if err != nil {
		return nil, nil, nil, err
	}

	acquirerID, err := GetAcquirerIdFromChainProfileId(profileId)
	if err != nil {
		return nil, nil, nil, err
	}

	rows, err := db.Query("call get_chain_data(?, ?)", profileId, acquirerID)
	if err == nil {
		defer rows.Close()
	} else {
		return nil, nil, nil, err
	}

	tidOverRideDataElement, err := GetTIDOverRideDataElement()
	if err != nil {
		logging.Error(fmt.Sprintf("Error retrieving data element tid over ride enable %s", err.Error()))
		return nil, nil, nil, err
	}

	var chainGroups []*DataGroup
	var acquirerGroups []*DataGroup
	var globalGroups []*DataGroup
	for rows.Next() {
		var site SiteData
		var siteOptions string
		err = rows.Scan(&site.DataGroupID, &site.DataGroup, &site.DataGroupDisplayName, &site.DataElementID,
			&site.Name, &site.Tooltip, &site.Source, &site.DataValue, &site.Overriden, &site.DataType, &site.IsAllowEmpty,
			&site.MaxLength, &site.ValidationExpression, &site.ValidationMessage, &site.FrontEndValidate, &siteOptions, &site.SortOrderInGroup,
			&site.DisplayName, &site.IsEncrypted, &site.IsPassword, &site.IsNotOverridable, &site.IsRequiredAtAcquireLevel, &site.IsRequiredAtChainLevel)
		if err != nil {
			return nil, nil, nil, err
		}

		if site.IsEncrypted.Valid && site.IsEncrypted.Bool && site.DataValue.Valid {
			site.DataValue.String, err = crypt.Decrypt(site.DataValue.String)
			if err != nil {
				return nil, nil, nil, err
			}
		}

		site.Options, site.OptionSelectable = BuildOptionsData(siteOptions, site.DataValue.String, site.DataGroup, site.Name, profileId)

		switch site.Source.String {
		case "":
			fallthrough
		case "chain":
			if !tidOverRideDataElement[strconv.Itoa(site.DataElementID)] {
				chainGroups = addDataElement(site, chainGroups)
			}
		case "acquirer":
			if !tidOverRideDataElement[strconv.Itoa(site.DataElementID)] {
				acquirerGroups = addDataElement(site, acquirerGroups)
			}
		case "global":
			if !tidOverRideDataElement[strconv.Itoa(site.DataElementID)] {
				globalGroups = addDataElement(site, globalGroups)
			}
		}
	}

	return chainGroups, acquirerGroups, globalGroups, nil
}

// GetChainList Method for searching chains
func GetChainList(searchTerm string, user *entities.TMSUser) ([]*ChainList, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	// Find user acquirers to limit search results
	acquirers, err := GetUserAcquirerPermissions(user)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("Call chain_list_fetch(?, ?)", searchTerm, acquirers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chains []*ChainList

	for rows.Next() {
		var chain ChainList
		err = rows.Scan(
			&chain.ChainProfileID,
			&chain.ChainName,
			&chain.AcquirerName)

		if err != nil {
			return nil, err
		}

		chains = append(chains, &chain)
	}

	return chains, nil
}

func GetIsOverridenForChain(profileId int, elementId int) (bool, error) {
	db, err := GetDB()
	if err != nil {
		return false, err
	}

	query := `
select source
from chain_data
where profile_id = ? 
and data_element_id = ?
`

	rows, err := db.Query(query, profileId, elementId)
	if err == nil {
		defer rows.Close()
	} else {
		return false, err
	}

	if rows.Next() {
		var source sql.NullString
		err = rows.Scan(&source)
		if err != nil {
			return false, err
		}

		return source.Valid && source.String != "chain", nil
	} else {
		return false, nil
	}
}

func GetChainPage(searchTerm string, offset string, amount string, orderedColumn string, orderDirection string, user *entities.TMSUser) (page []*ChainList, total int, filtered int, err error) {
	db, err := GetDB()
	if err != nil {
		return nil, 0, 0, err
	}

	// Find user acquirers to limit search results
	acquirers, err := GetUserAcquirerPermissions(user)
	if err != nil {
		return nil, 0, 0, err
	}

	rows, err := db.Query("CALL get_chain_page(?, ?)", searchTerm, acquirers)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()

	var chains []*ChainList
	for rows.Next() {
		var chain ChainList
		err = rows.Scan(&chain.ChainProfileID, &chain.ChainName, &chain.AcquirerName, &chain.ChainTIDCount)
		if err != nil {
			return nil, 0, 0, err
		}

		chain.AcquirerName = html.EscapeString(chain.AcquirerName)
		chain.ChainName = html.EscapeString(chain.ChainName)

		chains = append(chains, &chain)
	}
	total = len(chains)

	//sort the slice
	sort.Slice(chains, func(i, j int) bool {
		switch orderedColumn {
		case "0":
			if strings.ToUpper(orderDirection) == "ASC" {
				return chains[i].ChainProfileID < chains[j].ChainProfileID
			} else {
				return chains[i].ChainProfileID > chains[j].ChainProfileID
			}
		case "1":
			if strings.ToUpper(orderDirection) == "ASC" {
				return strings.ToLower(chains[i].ChainName) < strings.ToLower(chains[j].ChainName)
			} else {
				return strings.ToLower(chains[i].ChainName) > strings.ToLower(chains[j].ChainName)
			}
		case "2":
			if strings.ToUpper(orderDirection) == "ASC" {
				return strings.ToLower(chains[i].AcquirerName) < strings.ToLower(chains[j].AcquirerName)
			} else {
				return strings.ToLower(chains[i].AcquirerName) > strings.ToLower(chains[j].AcquirerName)
			}
		}
		return false
	})

	//Apply the limiting
	if amount != "-1" {
		offsetInt, err := strconv.Atoi(offset)
		if err != nil {
			return nil, 0, 0, err
		}
		amountInt, err := strconv.Atoi(amount)
		if err != nil {
			return nil, 0, 0, err
		}

		results := offsetInt + amountInt
		if amountInt > len(chains) || results > len(chains) {
			results = len(chains)
		}
		chains = chains[offsetInt:results]
	}

	filtered = total

	return chains, total, filtered, nil
}

func generateChainSearchWhere(searchTerm string) string {
	if searchTerm == "" {
		return ""
	}

	searchTerm = "%" + strings.ToUpper(searchTerm) + "%"

	return "AND (upper(p.name) like \"" + searchTerm + "\" or p.profile_id like \"" + searchTerm + "\" or p2.name like \"" + searchTerm + "\")"

}

func GetAcquirerIdFromChainProfileId(profileId int) (int, error) {
	db, err := GetDB()
	if err != nil {
		return -1, err
	}

	row := db.QueryRow("select acquirer_id from chain_profiles where chain_profile_id = ?", profileId)
	var id int
	err = row.Scan(&id)
	if err != nil {
		logging.Error(err)
		return -1, err
	}
	return id, nil
}

func CheckChainNameExists(chainName string) (bool, error) {
	db, err := GetDB()
	if err != nil {
		return false, err
	}

	var chainProfileId int
	rows, err := db.Query("SELECT profile_id FROM profile WHERE name = ?", chainName)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&chainProfileId)
	}

	if chainProfileId == 0 {
		return true, nil
	}

	acquirerId, err := GetAcquirerIdFromChainProfileId(chainProfileId)
	if err != nil {
		return false, err
	}

	if acquirerId != 0 {
		return false, err
	}

	return true, err
}

func CheckAcquirerNameExists(name string) (bool, error) {
	db, err := GetDB()
	if err != nil {
		return false, err
	}

	var acquirerProfileID int
	row := db.QueryRow("SELECT profile_id FROM profile WHERE profile_type_id = (select profile_type_id from profile_type where name = 'acquirer') AND name = ? LIMIT 1", name)

	err = row.Scan(&acquirerProfileID)
	if err != nil && err.Error() != "sql: no rows in result set" {
		logging.Error(err)
		return false, err
	}
	if acquirerProfileID == 0 {
		return false, nil
	}

	return true, err
}
