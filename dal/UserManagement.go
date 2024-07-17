package dal

import (
	"database/sql"
	"nextgen-tms-website/entities"

	"bitbucket.org/network-international/nextgen-libs/nextgen-helpers/TypeComparisonHelpers/SliceComparisonHelpers"
)

type TransactionPurpose int

var AutomationUsers []string

const (
	AddOrUpdateSiteUser TransactionPurpose = iota
	DeleteSiteUser
	AddOrUpdateTidUserOverride
	DeleteTidUserOverride
)

type UserUpdateResult struct {
	Result DbTransactionResult
	User   entities.SiteUser
	Action string
}

func (ur *UserUpdateResult) Equals(result UserUpdateResult) bool {
	if ur.Action != result.Action {
		return false
	}

	if ur.Result.Success != result.Result.Success {
		return false
	}
	if ur.Result.ErrorMessage != result.Result.ErrorMessage {
		return false
	}

	if ur.User.Username != result.User.Username {
		return false
	}
	if ur.User.PIN != result.User.PIN {
		return false
	}
	if ur.User.UserId != result.User.UserId {
		return false
	}
	if ur.User.SiteId != result.User.SiteId {
		return false
	}
	if ur.User.TidId != result.User.TidId {
		return false
	}
	if !SliceComparisonHelpers.SlicesOfStringAreEqual(ur.User.Modules, result.User.Modules, true) {
		return false
	}
	return true
}

func SaveUserGroup(user string, groups []string) {
	db, err := GetDB()
	if err != nil {
		return
	}

	var groupId int

	userId, err := GetUserID(user)
	if err != nil {
		logging.Error(err)
		return
	}

	_, err = db.Exec("delete from user_permissiongroup where user_id = ?", userId)
	if err != nil {
		logging.Error(err)
		return
	}

	for i := range groups {
		var groupName = groups[i]
		rows, err := db.Query("select group_id from permissiongroup pg where pg.name = (?)", groupName)
		if err != nil {
			logging.Error(err)
			return
		}

		for rows.Next() {
			err = rows.Scan(&groupId)
			if err != nil {
				logging.Error(err)
				rows.Close()
				return
			}
		}
		rows.Close()

		_, err = db.Exec("insert into user_permissiongroup (user_id, permission_group_id) values (?, ?)", userId, groupId)
		if err != nil {
			logging.Error(err)
			return
		}
	}
}

func RemoveGlobalAdminGroup(user string) {
	db, err := GetDB()
	if err != nil {
		return
	}

	userId, err := GetUserID(user)
	if err != nil {
		logging.Error(err)
		return
	}

	_, err = db.Exec("delete user_permissiongroup "+
		"from user_permissiongroup "+
		"left join permissiongroup pg on pg.group_id = user_permissiongroup.permission_group_id "+
		"where pg.name = 'GlobalAdmin' and user_id = (?)", userId)
	if err != nil {
		logging.Error(err)
		return
	}
}

func GetUserPermissions(userName string) (entities.UserPermissions, error) {
	var userPermissions entities.UserPermissions
	db, err := GetDB()
	if err != nil {
		return userPermissions, err
	}

	// Get user ID
	userId, err := GetUserID(userName)
	if err != nil {
		logging.Error(err)
		return userPermissions, err
	}

	// Get User Groups
	var groupId int
	rows, err := db.Query("Call get_user_groups(?)", userId)
	if err != nil {
		logging.Error(err)
		return userPermissions, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&groupId)
		if err != nil {
			logging.Error(err)
			return userPermissions, err
		}
		// Get group permissions
		var permissionId int
		userGroupPermissionsRows, err := db.Query("Call get_group_permissions(?)", groupId)
		if err != nil {
			logging.Error(err)
			return userPermissions, err
		}
		for userGroupPermissionsRows.Next() {
			err = userGroupPermissionsRows.Scan(&permissionId)
			if err != nil {
				logging.Error(err)
				userGroupPermissionsRows.Close()
				return userPermissions, err
			}
			switch permissionId {
			case 1:
				userPermissions.SiteWrite = true
			case 2:
				userPermissions.SiteDelete = true
			case 3:
				userPermissions.ChangeApprovalRead = true
			case 4:
				userPermissions.ChangeApprovalWrite = true
			case 5:
				userPermissions.AddCreate = true
			case 6:
				userPermissions.BulkUpdates = true
			case 7:
			case 8:
				userPermissions.UserManagement = true
			case 9:
				userPermissions.TransactionViewer = true
			case 10:
				userPermissions.DirectQuery = true
			case 11:
				userPermissions.Reporting = true
			case 12:
				userPermissions.ChangeHistoryView = true
			case 13:
				userPermissions.EditPasswords = true
			case 14:
				userPermissions.Fraud = true
			case 15:
				userPermissions.PermissionGroups = true
			case 16:
				userPermissions.UserManagementAudit = true
			case 19:
				userPermissions.BulkImport = true
			case 20:
				userPermissions.ContactEdit = true
			case 21:
				userPermissions.OfflinePIN = true
			case 22:
				userPermissions.DbBackup = true
			case 23:
				userPermissions.TerminalFlagging = true
			case 24:
				userPermissions.APIAutomation = true
			case 25:
				userPermissions.BulkChangeApproval = true
			case 26:
				userPermissions.ChainDuplication = true
			case 27:
				userPermissions.EditToken = true
			case 28:
				userPermissions.PaymentServices = true
			case 29:
				userPermissions.LogoManagement = true
			case 30:
				userPermissions.FileManagement = true
			case 31:
				userPermissions.SouhoolaLogin = true
			}
		}
		userGroupPermissionsRows.Close()
	}
	return userPermissions, nil
}

func GetUserID(username string) (int, error) {
	db, err := GetDB()
	if err != nil {
		return -1, err
	}

	// Get user ID
	var userId int
	rows, err := db.Query("select user_id from user u where u.username = (?)", username)
	if err != nil {
		logging.Error(err)
		return -1, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&userId)
		if err != nil {
			logging.Error(err)
			return -1, err
		}
	}
	return userId, nil
}

// TODO: We need to look at the lack of error handling in this function
func SetDefaultUserPermissions(user entities.TMSUser) {
	db, err := GetDB()
	if err != nil {
		return
	}

	// Get User Groups
	rows, err := db.Query("Call get_user_groups(?)", user.UserID)
	if err != nil {
		return
	}
	defer rows.Close()

	// If there are existing groups do nothing and return otherwise set default groups based on user role
	// Otherwise if GlobalAdmin assign all permissions to ensure permissions for existing users are
	// corrected if they are made a global admin
	if rows.Next() {
		if user.RoleID == 3 {
			var groups []string
			var globalAdminString = "GlobalAdmin"
			groups = append(groups, globalAdminString)
			SaveUserGroup(user.Username, groups)
		} else {
			RemoveGlobalAdminGroup(user.Username)
		}
	} else {
		var groups []string
		var adminString = "NI Admin"
		if user.RoleID == 3 {
			groups = append(groups, adminString)
		}
		SaveUserGroup(user.Username, groups)
	}
}

// Changes a users first time logon status.
// NB unlike the CheckFirstTimeLogon, this only operates on the
// user table corresponding to this particular site.
func ToggleFirstTimeLogon(username string, firstLogon bool) error {

	userId, err := GetUserID(username)
	if err != nil {
		return err
	}

	db, err := GetDB()
	if err != nil {
		return err
	}

	_, err = toggleFirstTimeLogon(db, userId, firstLogon)
	if err != nil {
		return err
	}
	return nil
}

// Updates first_logon value in user table
func toggleFirstTimeLogon(db *sql.DB, userId int, firstLogon bool) (sql.Result, error) {
	return db.Exec("UPDATE user SET first_logon = (?) WHERE user_id = (?)", firstLogon, userId)
}

func GetUserAcquirerPermissions(user *entities.TMSUser) (string, error) {
	db, err := GetDB()
	if err != nil {
		return "", err
	}

	// Find user acquirers to limit search results
	UserAqRows, err := db.Query("Call get_user_acquirer_permissions(?)", user.UserID)
	if err != nil {
		return "", err
	}
	defer UserAqRows.Close()

	acquirers := ""
	var aq string
	for UserAqRows.Next() {
		err = UserAqRows.Scan(&aq)
		if err != nil {
			return "", err
		}

		acquirers += aq + ","
	}
	UserAqRows.Close()

	return acquirers, nil
}

func GetUserAcquirers(user *entities.TMSUser) ([]string, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	// Find user acquirers to limit search results
	UserAqRows, err := db.Query("Call get_user_acquirer_permissions(?)", user.UserID)
	if err != nil {
		return nil, err
	}
	defer UserAqRows.Close()

	var acquirers []string
	var acquirer string
	for UserAqRows.Next() {
		err = UserAqRows.Scan(&acquirer)
		if err != nil {
			return nil, err
		}

		acquirers = append(acquirers, acquirer)
	}

	return acquirers, nil
}
