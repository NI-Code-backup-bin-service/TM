package dal

import (
	"database/sql"
	"errors"
	"fmt"
	"nextgen-tms-website/entities"
	"strconv"
	"strings"
)

const (
	// SQL TABLE/COLUMN NAMES:
	groupTbl        = "payment_service_group"
	serviceTbl      = "payment_service"
	groupIdColumn   = "group_id"
	serviceIdColumn = "service_id"
	nameColumn      = "name"
)

// SearchServiceGroups - get service groups that match the search string, paginated with an offset/pagesize
func SearchServiceGroups(searchTerm, orderDirection, offset, pagesize string) (count string, groups []*entities.PaymentServiceGroup) {
	db, err := GetDB()
	if err != nil {
		return
	}

	if strings.ToUpper(orderDirection) != "DESC" {
		orderDirection = ""
	}

	count = CountGroups(searchTerm, db)
	if count == "0" {
		return
	}

	rows, err := db.Query(
		fmt.Sprintf("select * from %s where %s like concat('%%', ?, '%%') order by %s %s limit ?,?", groupTbl, nameColumn, nameColumn, orderDirection),
		searchTerm,
		offset,
		pagesize)

	if err != nil {
		return "0", groups
	}
	defer closeRows(rows)

	for rows.Next() {
		var group entities.PaymentServiceGroup
		err = rows.Scan(&group.Id, &group.Name)
		if err != nil {
			continue
		}
		group.ServiceCount = CountServicesInGroup(group.Id, "", db)
		groups = append(groups, &group)
	}
	return
}

// GetPaymentServiceGroup - get the single service group identified by groupId
func GetPaymentServiceGroup(groupId string) (*entities.PaymentServiceGroup, error) {
	var group entities.PaymentServiceGroup
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	row := db.QueryRow(fmt.Sprintf("select * from %s where %s = ? limit 1", groupTbl, groupIdColumn), groupId)
	err = row.Scan(&group.Id, &group.Name)
	return &group, err
}

// SearchServicesInGroup - get services that match the search string and belong to the group with groupId, paginated with an offset/pagesize
func SearchServicesInGroup(groupId, searchTerm, orderDirection, offset, pagesize string) (count string, services []*entities.PaymentService) {
	db, err := GetDB()
	if err != nil {
		return
	}

	if strings.ToUpper(orderDirection) != "DESC" {
		orderDirection = ""
	}

	id, err := strconv.Atoi(groupId)
	if err != nil {
		return
	}

	count = CountServicesInGroup(uint(id), searchTerm, db)
	if count == "0" {
		return
	}

	rows, err := db.Query(
		fmt.Sprintf("select * from %s where %s = ? and %s like concat('%%', ?, '%%') order by %s %s limit ?,?", serviceTbl, groupIdColumn, nameColumn, nameColumn, orderDirection),
		groupId,
		searchTerm,
		offset,
		pagesize)

	if err != nil {
		return "0", services
	}
	defer closeRows(rows)

	for rows.Next() {
		var service entities.PaymentService
		err = rows.Scan(&service.Id, &service.GroupId, &service.Name)
		if err != nil {
			continue
		}
		services = append(services, &service)
	}
	return
}

// CountServicesInGroup - get the number of services that belong to the group identified by groupId
// (either the total number or the number matching a search string)
func CountServicesInGroup(groupId uint, searchKey string, db *sql.DB) string {
	var err error
	count := 0
	if db == nil {
		db, err = GetDB()
		if err != nil {
			return "0"
		}
	}

	var row *sql.Row
	if len(searchKey) == 0 {
		row = db.QueryRow(fmt.Sprintf("select count(*) from %s where %s = ?", serviceTbl, groupIdColumn), groupId)
	} else {
		row = db.QueryRow(fmt.Sprintf("select count(*) from %s where %s = ? and %s like concat('%%', ?, '%%')", serviceTbl, groupIdColumn, nameColumn), groupId, searchKey)
	}

	if row == nil {
		return "0"
	}
	err = row.Scan(&count)
	if err != nil {
		return "0"
	}
	return strconv.Itoa(count)
}

// CountGroups - get the number of service groups (either the total number or the number matching a search string)
func CountGroups(searchKey string, db *sql.DB) string {
	var err error
	count := 0
	if db == nil {
		db, err = GetDB()
		if err != nil {
			return "0"
		}
	}

	var row *sql.Row
	if len(searchKey) == 0 {
		row = db.QueryRow(fmt.Sprintf("select count(*) from %s", groupTbl))
	} else {
		row = db.QueryRow(fmt.Sprintf("select count(*) from %s where %s like concat('%%', ?, '%%')", groupTbl, nameColumn), searchKey)
	}

	if row == nil {
		return "0"
	}
	err = row.Scan(&count)
	if err != nil {
		return "0"
	}
	return strconv.Itoa(count)
}

// AddPaymentServiceGroup - add a new payment service group with the provided name
func AddPaymentServiceGroup(name string) error {
	db, err := GetDB()
	if err != nil {
		return errors.New("an internal error occurred")
	}

	count := 0
	row := db.QueryRow(fmt.Sprintf("select count(*) from %s where %s = ?", groupTbl, nameColumn), name)
	err = row.Scan(&count)

	if err != nil {
		return errors.New("an internal error occurred")
	}
	if count > 0 {
		return fmt.Errorf("a group with the name '%s' already exists", name)
	}

	result, err := db.Exec(fmt.Sprintf("insert into %s (`%s`) values (?)", groupTbl, nameColumn), name)
	if err != nil {
		return errors.New("an internal error occurred")
	}

	affected, err := result.RowsAffected()
	if err != nil || affected == 0 {
		return errors.New("an internal error occurred")
	}

	err = updatePaymentServiceOptions(db)
	if err != nil {
		return fmt.Errorf("an internal error occurred: %w", err)
	}
	return nil
}

// AddPaymentService - add a new payment service with the provided name to the group identified by groupId
func AddPaymentService(name, groupId string) error {
	db, err := GetDB()
	if err != nil {
		return errors.New("an internal error occurred")
	}

	if ServiceAlreadyExists(name, groupId, db) {
		return fmt.Errorf("a service with the name '%s' already exists", name)
	}

	result, err := db.Exec(fmt.Sprintf("insert into %s (`%s`, `%s`) values (?, ?)", serviceTbl, groupIdColumn, nameColumn), groupId, name)
	if err != nil {
		return errors.New("an internal error occurred")
	}

	affected, err := result.RowsAffected()
	if err != nil || affected == 0 {
		return errors.New("an internal error occurred")
	}
	return nil
}

func ServiceAlreadyExists(name, groupId string, db *sql.DB) bool {
	if db == nil {
		var err error
		db, err = GetDB()
		if err != nil {
			return true
		}
	}

	count := 0
	row := db.QueryRow(fmt.Sprintf("select count(*) from %s where %s = ? and %s = ?", serviceTbl, groupIdColumn, nameColumn), groupId, name)
	err := row.Scan(&count)
	if err != nil {
		return true
	}
	return count > 0
}

func GroupIdFromName(name string) int {
	groupId := -1
	db, err := GetDB()
	if err != nil {
		return groupId
	}

	row := db.QueryRow(fmt.Sprintf("select %s from %s where %s = ?", groupIdColumn, groupTbl, nameColumn), name)
	err = row.Scan(&groupId)
	if err != nil {
		return -1
	}
	return groupId
}

// DeleteServiceGroup - delete the service group identified by groupId
func DeleteServiceGroup(groupId string) bool {
	return paymentServicesDelete(groupId, groupTbl, groupIdColumn)
}

// DeleteService - delete the service identified by serviceId
func DeleteService(serviceId string) bool {
	return paymentServicesDelete(serviceId, serviceTbl, serviceIdColumn)
}

func IsServiceInGroup(groupName, serviceName string) bool {
	groupId := GroupIdFromName(groupName)
	if groupId < 1 {
		return false
	}

	db, err := GetDB()
	if err != nil {
		return false
	}

	count := 0
	row := db.QueryRow(fmt.Sprintf("select count(*) from %s where %s = ? and %s = ?", serviceTbl, groupIdColumn, nameColumn), groupId, serviceName)
	err = row.Scan(&count)
	if err != nil {
		return false
	}

	return count > 0
}

func GetServiceIdFromNames(groupName, serviceName string) int {
	id := -1
	groupId := GroupIdFromName(groupName)
	if groupId == -1 {
		return id
	}

	db, err := GetDB()
	if err != nil {
		return id
	}
	row := db.QueryRow(fmt.Sprintf("select %s from %s where %s = ? and %s = ?", serviceIdColumn, serviceTbl, groupIdColumn, nameColumn), groupId, serviceName)
	err = row.Scan(&id)
	if err != nil {
		return -1
	}
	return id
}

func GetProfileAssignedServiceGroupName(profileId string) string {
	db, err := GetDB()
	if err != nil {
		return ""
	}

	result := ""
	row := db.QueryRow(`select datavalue from profile_data as pd
			inner join data_element as de on pd.data_element_id = de.data_element_id
			where pd.profile_id = ?
			and de.name = 'paymentServiceGroup';`, profileId)

	err = row.Scan(&result)
	if err != nil {
		return ""
	}
	return result
}

func paymentServicesDelete(id, tbl, column string) bool {
	db, err := GetDB()
	if err != nil {
		return false
	}

	result, err := db.Exec(fmt.Sprintf("delete from %s where %s = ?", tbl, column), id)
	if err != nil {
		return false
	}

	count, err := result.RowsAffected()
	if err != nil {
		return false
	}

	if tbl == groupTbl {
		err = updatePaymentServiceOptions(db)
	}
	return count > 0 && err == nil
}

func closeRows(rows *sql.Rows) {
	if rows == nil {
		return
	}
	_ = rows.Close()
}

func updatePaymentServiceOptions(db *sql.DB) error {
	rows, err := db.Query(fmt.Sprintf("select %s from %s", nameColumn, groupTbl))
	if err != nil {
		return errors.New("an internal error occurred")
	}
	defer rows.Close()

	var optionsList []string
	for rows.Next() {
		var service string
		err = rows.Scan(&service)
		if err != nil {
			continue
		}
		optionsList = append(optionsList, service)
	}

	options := strings.Join(optionsList, "|")
	result, err := db.Exec("UPDATE data_element de INNER JOIN data_group dg ON dg.data_group_id = de.data_group_id  SET options = ? where dg.name = ? and de.name = ?", options, "paymentServices", "paymentServiceGroup")
	if err != nil {
		return errors.New("an internal error occurred")
	}

	affected, err := result.RowsAffected()
	if err != nil || affected == 0 {
		return errors.New("data element options not updated")
	}
	return nil
}
