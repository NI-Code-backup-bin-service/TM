package dal

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DataGroupInfo struct {
	DataGroupID  int
	Name         string
	ElementNames []string
}

type DataGroupInfoTidExport struct {
	Name               string
	ElementName        string
	ExportDisplayIndex string
	DisplayName        string
	DisplayNameEn      string
}

type CountItem struct {
	Label string
	Value int
}

type SumItem struct {
	Label string
	Value string
}

func GetDataGroupInfo() []DataGroupInfo {
	db, err := GetDB()
	if err != nil {
		return nil
	}

	rows, err := db.Query("SELECT data_group_id FROM data_group")
	if err != nil {
		return nil
	}
	defer rows.Close()

	var dataGroupIds []int
	for rows.Next() {
		var dataGroupId int
		rows.Scan(&dataGroupId)
		dataGroupIds = append(dataGroupIds, dataGroupId)
	}

	var dgi = make([]DataGroupInfo, 0)
	for i := 0; i < len(dataGroupIds); i++ {
		var dg DataGroupInfo

		rows, err := db.Query("CALL get_data_group_info(?)", dataGroupIds[i])
		if err != nil {
			return nil
		}

		var elementName sql.NullString
		for rows.Next() {
			rows.Scan(&dg.DataGroupID, &dg.Name, &elementName)

			if elementName.Valid {
				dg.ElementNames = append(dg.ElementNames, elementName.String)
			}
		}
		dgi = append(dgi, dg)
		rows.Close()
	}
	return dgi
}

func GetDataGroupInfoForTIDExport() ([]DataGroupInfoTidExport, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	var dgi = make([]DataGroupInfoTidExport, 0)
	rows, err := db.Query("CALL get_data_group_info_for_TID_export()")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var dg DataGroupInfoTidExport
		err = rows.Scan(&dg.Name, &dg.ElementName, &dg.DisplayNameEn, &dg.ExportDisplayIndex)
		if err != nil {
			return nil, err
		}
		dg.DisplayName = dg.Name + "." + dg.ElementName
		dgi = append(dgi, dg)
	}

	return dgi, nil
}

func GetReportData(siteId int, profileType string, profileId int) ([]*DataGroup, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	if siteId == 0 {
		rows, err := db.Query("Call get_report_data_using_profile(?)", profileId)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var siteGroups []*DataGroup
		for rows.Next() {
			var site SiteData
			var siteOptions string
			err = rows.Scan(&site.DataGroupID, &site.DataGroup, &site.DataElementID,
				&site.Name, &site.DataValue, &site.DataType, &site.IsAllowEmpty,
				&site.MaxLength, &site.ValidationExpression, &site.ValidationMessage, &site.FrontEndValidate, &siteOptions, &site.DisplayName)
			if err != nil {
				return nil, err
			}
			siteGroups = addDataElement(site, siteGroups)
		}

		rows.Close()
		return siteGroups, nil
	} else {
		rows, err := db.Query("Call search_report_data(?, ?, ?)", siteId, profileType, 0)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var siteGroups []*DataGroup
		for rows.Next() {
			var site SiteData
			var siteOptions string
			err = rows.Scan(&site.TIDID, &site.DataGroupID, &site.DataGroup, &site.DataElementID,
				&site.Name, &site.Source, &site.DataValue, &site.Overriden, &site.DataType, &site.IsAllowEmpty,
				&site.MaxLength, &site.ValidationExpression, &site.ValidationMessage, &site.FrontEndValidate, &siteOptions, &site.DisplayName)
			if err != nil {
				return nil, err
			}

			siteGroups = addDataElement(site, siteGroups)
		}

		rows.Close()

		return siteGroups, nil
	}

}
func GetReportBatchData(sites []int, searchTerm string) (map[int][]*DataGroup, error) {
	siteArr := make([]SiteData, 0)
	var err error
	siteDataGroups := make(map[int][]*DataGroup)
	if searchTerm != "" {
		siteArr, err = fetchSiteDataWithFilters(sites)
		if err != nil {
			return nil, err
		}
	} else {
		siteArr, err = fetchSiteDataWithOutFilters()
		if err != nil {
			return nil, err
		}
	}

	for _, site := range siteArr {
		siteGroups := addDataElement(site, siteDataGroups[site.TIDID])
		siteDataGroups[site.TIDID] = siteGroups
	}

	return siteDataGroups, nil
}

func fetchSiteDataWithFilters(sites []int) ([]SiteData, error) {
	siteArr := make([]SiteData, 0)
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	var sitess = strings.Trim(strings.Join(strings.Fields(fmt.Sprint(sites)), ","), "[]")
	rows, err := db.Query("Call search_report_data_batch_with_site_ids(?, ?)", sitess, 0)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var site SiteData
		var siteOptions string
		var profileId int
		err = rows.Scan(&site.TIDID, &profileId, &site.DataGroupID, &site.DataGroup, &site.DataElementID,
			&site.Name, &site.Source, &site.DataValue, &site.Overriden, &site.DataType, &site.IsAllowEmpty,
			&site.MaxLength, &site.ValidationExpression, &site.ValidationMessage, &site.FrontEndValidate, &siteOptions, &site.DisplayName)
		if err != nil {
			return nil, err
		}

		siteArr = append(siteArr, site)
	}

	return siteArr, nil
}

func fetchSiteDataWithOutFilters() ([]SiteData, error) {
	siteArr := make([]SiteData, 0)
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("Call search_report_data_batch()")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var site SiteData
		var siteOptions string
		var profileId int
		err = rows.Scan(&site.TIDID, &profileId, &site.DataGroupID, &site.DataGroup, &site.DataElementID,
			&site.Name, &site.DataValue, &site.Overriden, &site.DataType, &site.IsAllowEmpty,
			&site.MaxLength, &site.ValidationExpression, &site.ValidationMessage, &site.FrontEndValidate, &siteOptions, &site.DisplayName)
		if err != nil {
			return nil, err
		}

		siteArr = append(siteArr, site)
	}

	return siteArr, nil
}

func GetTotalMIDs() (int, error) {
	db, err := GetDB()
	if err != nil {
		return 0, err
	}

	rows, err := db.Query("SELECT COUNT(*) FROM site")
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var midCount int
	for rows.Next() {
		rows.Scan(&midCount)
	}

	return midCount, nil
}

func GetTotalTIDs() (int, error) {
	db, err := GetDB()
	if err != nil {
		return 0, err
	}

	rows, err := db.Query("SELECT COUNT(*) FROM tid")
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var tidCount int
	for rows.Next() {
		rows.Scan(&tidCount)
	}

	return tidCount, nil
}

func GetTotals(ctx context.Context, filters []bson.M) (int32, error) {
	results, err := AggregateMongoQuery(ctx, filters)
	if err != nil {
		return 0, err
	}

	for _, value := range results {
		values := value.(primitive.D)
		for _, t := range values {
			if t.Key == "count" {
				return t.Value.(int32), nil
			}
		}
	}

	return 0, nil
}

func GetMonetaryTotals(ctx context.Context, filters []bson.M) (float64, error) {
	results, err := AggregateMongoQuery(ctx, filters)
	if err != nil {
		return 0, err
	}

	for _, value := range results {
		values := value.(primitive.D)
		for _, t := range values {
			if t.Key == "count" {
				s := t.Value
				val := fmt.Sprintf("%.2f", s)
				s, err = strconv.ParseFloat(val, 64)
				if err != nil {
					s = "0.0"
				}
				return s.(float64), nil
			}
		}
	}

	return 0, nil
}

func GetTopTransacting(ctx context.Context, filters []bson.M) ([]CountItem, error) {
	var items = make([]CountItem, 0)

	results, err := AggregateMongoQuery(ctx, filters)
	if err != nil {
		return items, err
	}

	for _, value := range results {
		values := value.(primitive.D)
		var id string
		var count int

		for _, t := range values {
			if t.Key == "_id" {
				ids := t.Value.(primitive.D)
				for _, _id := range ids {
					if _id.Key == "id" {
						id = _id.Value.(string)
					}
				}
			}

			if t.Key == "count" {
				count = int(t.Value.(int32))
			}
		}

		item := CountItem{Label: id, Value: count}
		items = append(items, item)
	}

	return items, nil
}

func GetMonetaryTopTransacting(ctx context.Context, filters []bson.M) ([]SumItem, error) {
	var items = make([]SumItem, 0)

	results, err := AggregateMongoQuery(ctx, filters)
	if err != nil {
		return items, err
	}

	for _, value := range results {
		values := value.(primitive.D)
		var id, count string

		for _, t := range values {
			if t.Key == "_id" {
				ids := t.Value.(primitive.D)
				for _, _id := range ids {
					if _id.Key == "id" {
						id = _id.Value.(string)
					}
				}
			}

			if t.Key == "count" {
				count = strconv.FormatFloat(t.Value.(float64), 'f', -1, 64)
			}
		}

		item := SumItem{Label: id, Value: count}
		items = append(items, item)
	}

	return items, nil
}

func GetTransactionVolume(ctx context.Context, filters []bson.M) ([]SumItem, error) {
	var items = make([]SumItem, 0)

	results, err := AggregateMongoQuery(ctx, filters)
	if err != nil {
		return items, err
	}

	for _, value := range results {
		values := value.(primitive.D)
		var date, count string

		for _, t := range values {
			if t.Key == "_id" {
				ids := t.Value.(primitive.D)
				for _, id := range ids {
					if id.Key == "date" {
						date = id.Value.(string)
					}
				}
			}

			if t.Key == "count" {
				count = strconv.Itoa(int(t.Value.(int32)))
			}
		}

		item := SumItem{Label: date, Value: count}
		items = append(items, item)
	}
	return items, err
}

func GetMonetaryValues(ctx context.Context, filters []bson.M) ([]SumItem, error) {
	var items = make([]SumItem, 0)

	results, err := AggregateMongoQuery(ctx, filters)
	if err != nil {
		return items, err
	}

	for _, value := range results {
		values := value.(primitive.D)
		var date, count string

		for _, t := range values {
			if t.Key == "_id" {
				ids := t.Value.(primitive.D)
				for _, id := range ids {
					if id.Key == "date" {
						date = id.Value.(string)
					}
				}
			}

			if t.Key == "count" {
				count = strconv.FormatFloat(t.Value.(float64), 'f', -1, 64)
			}
		}

		item := SumItem{Label: date, Value: count}
		items = append(items, item)
	}

	return items, nil
}
