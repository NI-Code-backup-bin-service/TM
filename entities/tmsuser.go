package entities

import (
	"time"
)

const (
	RoleTMSAdministrator = 1
	RoleTMSUser          = 2
	RoleTMSSuperUser     = 3

	// Password change responses. There is no particular significance to the values chosen here;
	// the fact that certain numbers are missing is most likely simply due to some values being
	// removed during initial development.
	NoChangeRequired            = 9
	PasswordExpired             = 0
	FirstTimeLogon              = 2
	InsufficientPasswordQuality = 5
	PasswordTooShort            = 6
	PasswordTooYoung            = 7
	PasswordInHistory           = 8
)

type TMSUser struct {
	UserID          int
	Username        string
	Token           string
	RoleID          int
	LoggedIn        bool
	PasswordChange  bool
	UserPermissions UserPermissions
	Expires         time.Time
}

type UserPermissions struct {
	SiteWrite           bool
	SiteDelete          bool
	ChangeApprovalRead  bool
	ChangeApprovalWrite bool
	BulkChangeApproval  bool
	AddCreate           bool
	BulkUpdates         bool
	UserManagement      bool
	TransactionViewer   bool
	DirectQuery         bool
	Reporting           bool
	ChangeHistoryView   bool
	EditPasswords       bool
	Fraud               bool
	PermissionGroups    bool
	UserManagementAudit bool
	BulkImport          bool
	ContactEdit         bool
	OfflinePIN          bool
	DbBackup            bool
	TerminalFlagging    bool
	APIAutomation       bool
	PaymentServices     bool
	ChainDuplication    bool
	EditToken           bool
	FileManagement      bool
	LogoManagement      bool
	SouhoolaLogin       bool
}

type SiteUserData struct {
	UserId   int
	Username string
	PIN      string
	Modules  []string
}

type SiteUser struct {
	UserId   int
	Username string
	PIN      string
	Modules  []string
	SiteId   int
	TidId    int
}

func (su *SiteUser) IsEqualTo(user SiteUser) bool {
	//TODO check if the modules match

	if su.PIN != user.PIN {
		return false
	}
	if su.UserId != user.UserId {
		return false
	}
	if su.Username != user.Username {
		return false
	}
	return true
}

type SiteUserExportModel struct {
	Username string
	PIN      string
	Modules  map[string]bool
	Tid      string
}

func (u TMSUser) IsAdmin() bool {
	if u.RoleID == RoleTMSAdministrator {
		return true
	} else {
		return false
	}
}

func (u TMSUser) IsUser() bool {
	if u.RoleID == RoleTMSUser {
		return true
	} else {
		return false
	}
}

func (u TMSUser) IsSuperUser() bool {
	if u.RoleID == RoleTMSSuperUser {
		return true
	} else {
		return false
	}
}
