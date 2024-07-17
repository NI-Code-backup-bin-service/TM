package dal

import (
	le "bitbucket.org/network-international/nextgen-libs/nextgen-tg-protobuf/Logging"
	txn "bitbucket.org/network-international/nextgen-libs/nextgen-tg-protobuf/Transaction"
	sharedDAL "bitbucket.org/network-international/nextgen-tms/web-shared/dal"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"regexp"
	"strings"
	"time"
)

type AuditFilter struct {
	Name  string
	Value string
}

func BuildUserAuditHistory(ctx context.Context, filters *sharedDAL.UserAuditFilterGroupModel, offsetRows int, userID int) ([]sharedDAL.UserAuditDisplay, error) {
	pageSize := 50

	var auditFilter []AuditFilter
	auditFilter = append(auditFilter, AuditFilter{
		Name:  "extra",
		Value: "AUDIT_UMN",
	})

	for _, filter := range filters.Filters {
		auditFilter = append(auditFilter, AuditFilter{
			Name:  "audit." + filter.Name,
			Value: filter.Value,
		})
	}

	var pipe []bson.M
	find := bson.M{"$match": buildMongoAuditQuery(auditFilter)}
	sort := bson.M{"$sort": bson.M{"date": -1}}
	project := bson.M{"$project": buildMongoAuditProject()}
	limit := bson.M{"$limit": pageSize}
	skip := bson.M{"$skip": offsetRows}

	pipe = append(pipe, find)
	pipe = append(pipe, sort)
	pipe = append(pipe, project)
	pipe = append(pipe, skip)
	pipe = append(pipe, limit)

	userAudits, err := AggregateAuditCollection(ctx, pipe)
	if err != nil {
		return nil, err
	}

	return userAudits, nil
}

// Limits the returned data to only that we need
func buildMongoAuditProject() bson.M {
	match := bson.M{
		"_id":           0,
		"acquirer":      "$audit.acquirer",
		"name":          "$audit.name",
		"module":        "$audit.module",
		"originalvalue": "$audit.originalvalue",
		"updatedvalue":  "$audit.updatedvalue",
		"updatedby":     "$audit.updatedby",
		"updatedat":     "$audit.updatedat",
	}

	return match
}

func AggregateAuditCollection(ctx context.Context, pipe []bson.M) ([]sharedDAL.UserAuditDisplay, error) {
	result := make([]*le.UserManagementAudit, 50)
	friendlyResult := make([]sharedDAL.UserAuditDisplay, 50)
	client, err := GetMongoClient()
	if err != nil {
		return friendlyResult, err
	}

	cursor, err := client.Database(mongoSettings.LoggingDatabase).Collection("audit").Aggregate(ctx, pipe)
	if err != nil {
		return friendlyResult, err
	}

	defer cursor.Close(ctx)

	err = cursor.All(ctx, &result)
	if err != nil {
		return friendlyResult, err
	}

	friendlyResult = ConvertToDisplayFriendly(result)

	return friendlyResult, err
}

func ConvertToDisplayFriendly(audits []*le.UserManagementAudit) []sharedDAL.UserAuditDisplay {
	var displayFriendlyAudits []sharedDAL.UserAuditDisplay

	for _, audit := range audits {
		var displayFriendlyAudit sharedDAL.UserAuditDisplay
		displayFriendlyAudit.Acquirer = audit.Acquirer
		displayFriendlyAudit.Name = audit.Name
		displayFriendlyAudit.Module = audit.Module
		displayFriendlyAudit.OriginalValue = strings.Split(audit.OriginalValue, ".")
		displayFriendlyAudit.UpdatedValue = strings.Split(audit.UpdatedValue, ".")
		displayFriendlyAudit.UpdatedBy = audit.UpdatedBy
		displayFriendlyAudit.UpdatedAt = audit.UpdatedAt.Local()
		displayFriendlyAudits = append(displayFriendlyAudits, displayFriendlyAudit)
	}
	return displayFriendlyAudits
}

func buildMongoAuditQuery(auditFilters []AuditFilter) bson.M {
	match := bson.M{}
	timeInputFormat := "2006/01/02 15:04"

	var or []bson.M

	if len(or) > 0 {
		match["$or"] = or
	}

	var filterGroup []bson.M

	for _, filter := range auditFilters {
		//Trim spaces
		filter.Value = strings.Trim(filter.Value, " ")

		_filter := bson.M{}

		if len(filter.Value) != 0 {
			if filter.Name == "audit.before" {
				timeValue, _ := time.Parse(timeInputFormat, filter.Value)
				_filter["audit.updatedat"] = bson.M{"$lt": timeValue}
			} else if filter.Name == "audit.after" {
				timeValue, _ := time.Parse(timeInputFormat, filter.Value)
				_filter["audit.updatedat"] = bson.M{"$gt": timeValue}
			} else {
				_filter[filter.Name] = bson.M{"$regex": ".*" + regexp.QuoteMeta(filter.Value) + ".*"}
			}
			filterGroup = append(filterGroup, _filter)
		}
	}

	match["$and"] = filterGroup

	return match
}

// Compares the requested acquirer with those the user is permitted to view.
func GetAcquirerForSearch(requested string, permitted string) (acquirers []string) {

	permittedList := strings.Split(permitted, ",")

	if requested == "" {
		acquirers = permittedList
	} else {
		for _, acquirer := range permittedList {
			if strings.ToLower(acquirer) == strings.ToLower(requested) {
				acquirers = append(acquirers, acquirer)
			}
		}
	}
	return
}

// Builds a user management change audit object to be logged to mongo
func BuildUserAuditEntry(user string, updatedValue []string, module string, username string) (txn.UserAuditHistory, error) {
	userId, err := GetUserID(user)
	if err != nil {
		return txn.UserAuditHistory{}, err
	}

	acquirers, err := GetUserAcquirer(userId)
	if err != nil {
		return txn.UserAuditHistory{}, err
	}
	groups, err := GetUserGroupNames(userId)
	if err != nil {
		return txn.UserAuditHistory{}, err
	}

	auditHistory := buildAuditEntry(acquirers, user, module, groups, updatedValue, username)

	return auditHistory, nil
}

// Checks to see if an update is required to the user group.
func UpdateUserGroupRequired(user string, updatedGroups []string) (bool, error) {
	userId, err := GetUserID(user)
	if err != nil {
		return false, err
	}

	originalGroups, err := GetUserGroupNames(userId)
	if err != nil {
		return false, err
	}

	updateRequired, _, _ := CompareGroupElements(originalGroups, updatedGroups)
	return updateRequired, nil
}

func buildAuditEntry(acquirers []string, name string, module string, originalValue []string, updatedValue []string, tmsUser string) txn.UserAuditHistory {

	var auditHistory txn.UserAuditHistory

	auditHistory.Acquirer = strings.Join(acquirers, ", ")
	auditHistory.Name = name
	auditHistory.Module = module
	auditHistory.OriginalValue = strings.Join(originalValue, ". ")
	auditHistory.UpdatedValue = strings.Join(updatedValue, ". ")
	auditHistory.UpdatedBy = tmsUser
	auditHistory.UpdatedAt = time.Now().Format(time.RFC3339Nano)

	return auditHistory
}

// Returns an array of acquirer names for a given user ID
func GetUserAcquirer(userId int) ([]string, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("Call get_user_acquirer_permissions(?)", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var acquirers []string
	for rows.Next() {
		var acquirer string
		err = rows.Scan(&acquirer)
		if err != nil {
			return nil, err
		}
		acquirers = append(acquirers, acquirer)
	}
	return acquirers, nil
}

// Retrieves the groups to which a user has permission
func GetUserGroupNames(userId int) ([]string, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("Call get_user_group_names(?)", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []string
	for rows.Next() {
		var group string
		err = rows.Scan(&group)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}
	return groups, nil
}

// Builds an add group audit entry
func BuildAddGroupAuditEntry(groupName string, success bool, username string) txn.UserAuditHistory {

	var updatedValue []string

	if success {
		updatedValue = append(updatedValue, "User group "+groupName+" added")
	} else {
		updatedValue = append(updatedValue, "Error adding user group "+groupName)
	}

	acquirers, err := GetOriginalGroupAcquirers(groupName)

	if err != nil {
		logging.Error(err)
	}

	auditHistory := buildAuditEntry(acquirers, groupName, "group", nil, updatedValue, username)

	return auditHistory
}

// Compares the saved group permissions/acquirers and builds an audit object detailing the changes
func BuildGroupChangeAuditEntry(groupName string, newPermissions []string, newAcquirers []string, username string) (txn.UserAuditHistory, bool, error) {

	var originalValue []string
	var updatedValue []string
	var groupChanged bool

	//Fetch the original permissions
	originalPermissions, err := GetOriginalGroupPermissions(groupName)
	if err != nil {
		return txn.UserAuditHistory{}, false, err
	}

	//Compare the permissions received against those stored in the database
	permissionsChanged, added, removed := CompareGroupElements(originalPermissions, newPermissions)
	if permissionsChanged {
		if len(originalPermissions) != 0 {
			originalValue = append(originalValue, "ORIGINAL_GROUP_PERMISSIONS: "+strings.Join(originalPermissions, ", "))
		} else {
			originalValue = append(originalValue, "No previous Permissions")
		}
		if len(newPermissions) != 0 {
			if len(added) != 0 {
				updatedValue = append(updatedValue, "PERMISSIONS_ADDED: "+strings.Join(added, ", "))
			}
			if len(removed) != 0 {
				updatedValue = append(updatedValue, "PERMISSIONS_REMOVED: "+strings.Join(removed, ", "))
			}
		} else {
			updatedValue = append(updatedValue, "All Permissions removed")
		}

	}

	//Fetch the original acquirers for the group
	originalAcquirers, err := GetOriginalGroupAcquirers(groupName)
	if err != nil {
		return txn.UserAuditHistory{}, false, err
	}

	acquirersChanged, added, removed := CompareGroupElements(originalAcquirers, newAcquirers)
	if acquirersChanged {
		if len(originalAcquirers) != 0 {
			originalValue = append(originalValue, "ORIGINAL_GROUP_ACQUIRERS: "+strings.Join(originalAcquirers, ", "))
		} else {
			originalValue = append(originalValue, "No previous Acquirers")
		}
		if len(newAcquirers) != 0 {
			if len(added) != 0 {
				updatedValue = append(updatedValue, "ACQUIRERS_ADDED: "+strings.Join(added, ", "))
			}
			if len(removed) != 0 {
				updatedValue = append(updatedValue, "ACQUIRERS_REMOVED: "+strings.Join(removed, ", "))
			}
		} else {
			updatedValue = append(updatedValue, "All Acquirers removed")
		}
	}

	if acquirersChanged || permissionsChanged {
		groupChanged = true
	}

	auditHistory := buildAuditEntry(newAcquirers, groupName, "group", originalValue, updatedValue, username)

	return auditHistory, groupChanged, nil
}

func GetOriginalGroupAcquirers(groupName string) ([]string, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("CALL get_group_acquirers_by_name(?)", groupName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var acquirers []string

	for rows.Next() {
		var acquirer string

		err := rows.Scan(&acquirer)
		if err != nil {
			return nil, err
		}

		acquirers = append(acquirers, acquirer)
	}
	return acquirers, nil
}

func GetOriginalGroupPermissions(groupName string) ([]string, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("CALL get_group_permissions_by_name(?)", groupName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []string

	for rows.Next() {
		var permission string

		err := rows.Scan(&permission)
		if err != nil {
			return nil, err
		}

		permissions = append(permissions, permission)
	}
	return permissions, nil
}

func CompareGroupElements(old []string, new []string) (bool, []string, []string) {

	//Find the elements that have been added
	added := AddedOrRemovedElements(new, old)
	//Find the elements that have been removed
	removed := AddedOrRemovedElements(old, new)

	//If one is nil, then the other must be nil
	if (old == nil) != (new == nil) {
		return true, added, removed
	}

	//If they are not of the same length they are not equal
	if len(old) != len(new) {
		return true, added, removed
	}

	difference := make(map[string]int, len(old))
	for _, oldPermission := range old {
		difference[oldPermission]++
	}
	for _, newPermission := range new {
		if _, ok := difference[newPermission]; !ok {
			return true, added, removed
		}
		difference[newPermission] -= 1
		if difference[newPermission] == 0 {
			delete(difference, newPermission)
		}
	}

	if len(difference) == 0 {
		return false, nil, nil
	}

	return true, added, removed
}

// Compares the values of two string arrays and returns an array of elements found in the 1st but not the 2nd
func AddedOrRemovedElements(outer []string, inner []string) []string {
	var differentElements []string

	for _, outerElem := range outer {
		var different bool
		for _, innerElem := range inner {
			if outerElem == innerElem {
				different = true
			}
		}
		if !different {
			differentElements = append(differentElements, outerElem)
		}
	}

	return differentElements
}

func BuildGroupDeleteAuditEntity(groupName string, username string) (txn.UserAuditHistory, error) {
	acquirers, err := GetOriginalGroupAcquirers(groupName)
	if err != nil {
		return txn.UserAuditHistory{}, err
	}

	updateMessage := []string{groupName + " deleted from groups"}
	auditEntry := buildAuditEntry(acquirers, groupName, "group", nil, updateMessage, username)

	return auditEntry, nil
}

func BuildGroupRenameAuditEntry(groupName string, newName string, username string) (txn.UserAuditHistory, error) {
	acquirers, err := GetOriginalGroupAcquirers(groupName)
	if err != nil {
		return txn.UserAuditHistory{}, err
	}

	var updatedMessage []string
	updatedMessage = append(updatedMessage, "Name changed to ")
	updatedMessage = append(updatedMessage, newName)

	auditEntry := buildAuditEntry(acquirers, newName, "group", []string{groupName}, updatedMessage, username)
	return auditEntry, nil
}
