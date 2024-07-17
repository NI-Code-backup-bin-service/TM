package main

import (
	"errors"
	"nextgen-tms-website/dal"
	"nextgen-tms-website/entities"
	"reflect"
	"strings"
	"testing"

	"bitbucket.org/network-international/nextgen-libs/nextgen-helpers/TypeComparisonHelpers/SliceComparisonHelpers"
)

var (
	noPermissions = entities.TMSUser{
		UserPermissions: entities.UserPermissions{
			SiteWrite:           false,
			SiteDelete:          false,
			ChangeApprovalRead:  false,
			ChangeApprovalWrite: false,
			AddCreate:           false,
			BulkUpdates:         false,
			UserManagement:      false,
			TransactionViewer:   false,
			DirectQuery:         false,
			Reporting:           false,
			ChangeHistoryView:   false,
			EditPasswords:       false,
			Fraud:               false,
			PermissionGroups:    false,
			UserManagementAudit: false,
			BulkImport:          false,
			ContactEdit:         false,
			PaymentServices:     false,
			SouhoolaLogin:       false,
		},
	}

	allPermissions = entities.TMSUser{
		UserPermissions: entities.UserPermissions{
			SiteWrite:           true,
			SiteDelete:          true,
			ChangeApprovalRead:  true,
			ChangeApprovalWrite: true,
			AddCreate:           true,
			BulkUpdates:         true,
			UserManagement:      true,
			TransactionViewer:   true,
			DirectQuery:         true,
			Reporting:           true,
			ChangeHistoryView:   true,
			EditPasswords:       true,
			Fraud:               true,
			PermissionGroups:    true,
			UserManagementAudit: true,
			BulkImport:          true,
			ContactEdit:         true,
			PaymentServices:     true,
			SouhoolaLogin:       true,
		},
	}

	userOne = entities.SiteUser{Username: "User1", PIN: "1"}
	userTwo = entities.SiteUser{Username: "User2", PIN: "2"}
)

func TestHandleSubmodules(t *testing.T) {
	tests := []struct {
		name     string
		modules  []string
		expected []string
	}{
		{"Non specified single module", []string{"One"}, []string{"One"}},
		{"Non specified multiple module", []string{"One", "Two"}, []string{"One", "Two"}},
		{"alipay single module", []string{"alipay"}, []string{"alipaySale", "alipayVoid", "alipayRefund"}},
		{"preAuth single module", []string{"preAuth"}, []string{"preAuthSale", "preAuthCompletion", "preAuthCancel"}},
		{"Non specified & alipay multiple module", []string{"One", "alipay"}, []string{"One", "alipaySale", "alipayVoid", "alipayRefund"}},
		{"alipay & Non specified multiple module", []string{"alipay", "One"}, []string{"alipaySale", "alipayVoid", "alipayRefund", "One"}},
		{"Non specified & preAuth multiple module", []string{"One", "preAuth"}, []string{"One", "preAuthSale", "preAuthCompletion", "preAuthCancel"}},
		{"preAuth & Non specified multiple module", []string{"preAuth", "One"}, []string{"preAuthSale", "preAuthCompletion", "preAuthCancel", "One"}},
		{"alipay & preAuth multiple module", []string{"alipay", "preAuth"}, []string{"alipaySale", "alipayVoid", "alipayRefund", "preAuthSale", "preAuthCompletion", "preAuthCancel"}},
		{"preAuth & alipay multiple module", []string{"preAuth", "alipay"}, []string{"preAuthSale", "preAuthCompletion", "preAuthCancel", "alipaySale", "alipayVoid", "alipayRefund"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := dal.GetSubmodules(tt.modules)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("handleSubmodules() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestMakeModulesFriendly(t *testing.T) {
	tests := []struct {
		name     string
		modules  []string
		expected []string
	}{
		{"Capitalise single strings", []string{"one"}, []string{"One"}},
		{"Capitalise multiple strings", []string{"one", "two"}, []string{"One", "Two"}},
		{"Replace Gratuity with Tip (single)", []string{"gratuity"}, []string{"Tip"}},
		{"Replace Gratuity with Tip (multiple)", []string{"one", "gratuity", "two"}, []string{"One", "Tip", "Two"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := makeModulesFriendly(tt.modules)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("makeModulesFriendly() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func Test_ValidateProfileUser_NoUsers(t *testing.T) {
	if validateProfileUser(nil) != nil {
		t.Errorf("validateProfileUser() = %v, want nil", validateProfileUser(nil))
	}
}

func Test_ValidateProfileUser_DuplicatePins(t *testing.T) {
	users := []*entities.SiteUser{&userOne, &userOne}
	if validateProfileUser(users)[0].Result.Success {
		t.Errorf("validateProfileUser().Result.Success = true, want false")
	}
}

func Test_ValidateProfileUser_DifferentPins(t *testing.T) {
	users := []*entities.SiteUser{&userOne, &userTwo}
	if !reflect.DeepEqual(validateProfileUser(users), []*dal.UserUpdateResult{}) {
		t.Errorf("validateProfileUser() = %v, want []*dal.UserUpdateResult{}", validateProfileUser(nil))
	}
}

func TestCheckUserPermissions(t *testing.T) {
	tests := []struct {
		name     string
		required UserPermission
		user     *entities.TMSUser
		expected bool
	}{
		{"None | no permissions", 0, &noPermissions, true},
		{"None | all permissions", 0, &allPermissions, true},
		{"SiteWrite | no permissions", 1, &noPermissions, false},
		{"SiteWrite | all permissions", 1, &allPermissions, true},
		{"SiteDelete | no permissions", 2, &noPermissions, false},
		{"SiteDelete | all permissions", 2, &allPermissions, true},
		{"ChangeApprovalRead | no permissions", 3, &noPermissions, false},
		{"ChangeApprovalRead | all permissions", 3, &allPermissions, true},
		{"ChangeApprovalWrite | no permissions", 4, &noPermissions, false},
		{"ChangeApprovalWrite | all permissions", 4, &allPermissions, true},
		{"AddCreate | no permissions", 5, &noPermissions, false},
		{"AddCreate | all permissions", 5, &allPermissions, true},
		{"BulkUpdates | no permissions", 6, &noPermissions, false},
		{"BulkUpdates | all permissions", 6, &allPermissions, true},
		{"UserManagement | no permissions", 8, &noPermissions, false},
		{"UserManagement | all permissions", 8, &allPermissions, true},
		{"TransactionViewer | no permissions", 9, &noPermissions, false},
		{"TransactionViewer | all permissions", 9, &allPermissions, true},
		{"Reporting | no permissions", 11, &noPermissions, false},
		{"Reporting | all permissions", 11, &allPermissions, true},
		{"ChangeHistoryView | no permissions", 12, &noPermissions, false},
		{"ChangeHistoryView | all permissions", 12, &allPermissions, true},
		{"EditPasswords | no permissions", 13, &noPermissions, false},
		{"EditPasswords | all permissions", 13, &allPermissions, true},
		{"Fraud | no permissions", 14, &noPermissions, false},
		{"Fraud | all permissions", 14, &allPermissions, true},
		{"PermissionGroups | no permissions", 15, &noPermissions, false},
		{"PermissionGroups | all permissions", 15, &allPermissions, true},
		{"UserManagementAudit | no permissions", 16, &noPermissions, false},
		{"UserManagementAudit | all permissions", 16, &allPermissions, true},
		{"SouhoolaLogin | all permissions", 31, &allPermissions, true},
		{"SouhoolaLogin | no permissions", 31, &noPermissions, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checkUserPermissions(tt.required, tt.user)
			if got != tt.expected {
				t.Errorf("checkUserPermissions() = %v, want %v", got, tt.expected)
			}
		})
	}
}

type tidUser struct {
	UserId   int
	Tid      int
	Username string
	PIN      string
	Modules  []string
}

func (t tidUser) toSiteUser() entities.SiteUser {
	return entities.SiteUser{
		UserId:   t.UserId,
		Username: t.Username,
		PIN:      t.PIN,
		Modules:  t.Modules,
	}
}

type mockSaveSiteUsersDal struct {
	SiteUsers []entities.SiteUser
	TidUsers  []entities.SiteUser
}

// For simplicity, in the following implementations assume that we are using the correct site ID

func (d *mockSaveSiteUsersDal) GetUsersForSite(siteId int) ([]entities.SiteUser, error) {
	return d.SiteUsers, nil
}
func (d *mockSaveSiteUsersDal) GetUserForId(userId int) (*entities.SiteUser, error) {
	for _, siteUser := range d.SiteUsers {
		if siteUser.UserId == userId {
			return &siteUser, nil
		}
	}
	return nil, nil
}

func (d *mockSaveSiteUsersDal) GetUsersForTid(tid int) ([]entities.SiteUser, error) {
	users := make([]entities.SiteUser, 0)
	for _, tidUser := range d.TidUsers {
		if tidUser.TidId == tid {
			users = append(users, tidUser)
		}
	}
	return users, nil
}

func (d *mockSaveSiteUsersDal) GetTidUsersForSite(siteId_ int) ([]entities.SiteUser, error) {
	siteUsers := make([]entities.SiteUser, 0)
	for _, tidUser := range d.TidUsers {
		siteUsers = append(siteUsers, tidUser)
	}
	return siteUsers, nil
}

func (d *mockSaveSiteUsersDal) AddOrUpdateSiteUsers(siteId_ int, users []*entities.SiteUser) ([]*dal.UserUpdateResult, error) {
	updateResults := make([]*dal.UserUpdateResult, 0)
	// Update existing users
	for userIndex, existingUser := range d.SiteUsers {
		for _, userPtr := range users {
			updateResult := dal.UserUpdateResult{User: *userPtr, Action: "Update Site User", Result: dal.DbTransactionResult{Success: true}}
			user := *userPtr
			if existingUser.UserId == user.UserId {
				d.SiteUsers[userIndex] = user

				// Update the PINs of corresponding TID users
				for i, tidUser := range d.TidUsers {
					if tidUser.Username == user.Username {
						tidUser.PIN = user.PIN
						d.TidUsers[i] = tidUser
					}
				}
			}
			updateResults = append(updateResults, &updateResult)
		}
	}

	// Add new users
	for _, user := range users {
		if user.UserId <= 0 {
			updateResult := dal.UserUpdateResult{User: *user, Action: "Add Site User", Result: dal.DbTransactionResult{Success: true}}
			d.SiteUsers = append(d.SiteUsers, *user)
			updateResults = append(updateResults, &updateResult)
		}
	}

	return updateResults, nil
}

func (d *mockSaveSiteUsersDal) AddOrUpdateTidUserOverride(tid int, users []*entities.SiteUser) ([]*dal.UserUpdateResult, error) {
	updateResults := make([]*dal.UserUpdateResult, 0)
	for _, userPtr := range users {
		updateResult := dal.UserUpdateResult{User: *userPtr, Action: "Add TID User Override", Result: dal.DbTransactionResult{Success: true}}
		user := *userPtr
		d.TidUsers = append(d.TidUsers, user)
		updateResults = append(updateResults, &updateResult)
	}

	return updateResults, nil
}

func (d *mockSaveSiteUsersDal) DeleteSiteUsers(siteId int, userIdsToDelete []int) ([]*dal.UserUpdateResult, error) {
	updateResults := make([]*dal.UserUpdateResult, 0)
	newUsers := make([]entities.SiteUser, 0)
	for _, user := range d.SiteUsers {
		updateResult := dal.UserUpdateResult{User: user, Action: "Delete Site User", Result: dal.DbTransactionResult{Success: true}}
		if contains(userIdsToDelete, user.UserId) {
			// We also need to delete matching TID overrides with the same name
			newTidUsers := make([]entities.SiteUser, 0)
			for _, tidUser := range d.TidUsers {
				if tidUser.Username != user.Username {
					newTidUsers = append(newTidUsers, tidUser)
				}
			}
			d.TidUsers = newTidUsers
		} else {
			newUsers = append(newUsers, user)
		}
		updateResults = append(updateResults, &updateResult)
	}

	d.SiteUsers = newUsers
	return updateResults, nil
}

func (d *mockSaveSiteUsersDal) DeleteTidUsers(userIdsToDelete []int) ([]*dal.UserUpdateResult, error) {
	updateResults := make([]*dal.UserUpdateResult, 0)
	newUsers := make([]entities.SiteUser, 0)
	for _, user := range d.TidUsers {
		if !contains(userIdsToDelete, user.UserId) {
			updateResult := dal.UserUpdateResult{User: user, Action: "Delete TID User Override", Result: dal.DbTransactionResult{Success: true}}
			newUsers = append(newUsers, user)
			updateResults = append(updateResults, &updateResult)
		}
	}

	d.TidUsers = newUsers
	return updateResults, nil
}

var (
	siteId         = 7894
	defaultModules = []string{"foo", "bar"}

	alice = entities.SiteUser{
		UserId:   1,
		Username: "Alice",
		PIN:      "12345",
		Modules:  []string{"foo", "bar"},
		SiteId:   siteId,
	}

	bob = entities.SiteUser{
		UserId:   2,
		Username: "Bob",
		PIN:      "23456",
		Modules:  []string{"bar", "baz"},
		SiteId:   siteId,
	}

	siteUserWithOverride = entities.SiteUser{
		UserId:   5,
		Username: "Ellis",
		PIN:      "97346",
		Modules:  defaultModules,
		SiteId:   siteId,
	}

	tid1 = 88880001
	tid2 = 88880002

	tidUser1 = entities.SiteUser{
		UserId:   3,
		TidId:    tid1,
		Username: "Carol",
		PIN:      "34567",
		Modules:  []string{"bar", "wibble"},
		SiteId:   siteId,
	}

	tidUser2 = entities.SiteUser{
		UserId:   4,
		TidId:    tid1,
		Username: "Dave",
		PIN:      "45678",
		Modules:  []string{"baz"},
		SiteId:   siteId,
	}

	overriddenSiteUser1 = entities.SiteUser{
		UserId:   6,
		TidId:    tid2,
		Username: "Ellis",
		PIN:      "97346",
		Modules:  []string{"foo"},
		SiteId:   siteId,
	}

	newUser1 = entities.SiteUser{
		UserId:   -1,
		Username: "NewUser",
		PIN:      "97865",
		Modules:  defaultModules,
		SiteId:   siteId,
	}

	newTidUser1 = entities.SiteUser{
		UserId:   -2,
		Username: "NewTidUser1",
		TidId:    tid1,
		PIN:      "76436",
		Modules:  defaultModules,
		SiteId:   siteId,
	}

	initialSiteUsers = []entities.SiteUser{alice, bob, siteUserWithOverride}
	initialTidUsers  = []entities.SiteUser{tidUser1, tidUser2, overriddenSiteUser1}
	finalTidUsers    = []entities.SiteUser{tidUser1, tidUser2, overriddenSiteUser1, newTidUser1}
)

func initialSiteUsersCopy() []entities.SiteUser {
	sliceCopy := make([]entities.SiteUser, len(initialSiteUsers))
	copy(sliceCopy, initialSiteUsers)
	return sliceCopy
}

func initialTidUsersCopy() []entities.SiteUser {
	sliceCopy := make([]entities.SiteUser, len(initialTidUsers))
	copy(sliceCopy, initialTidUsers)
	return sliceCopy
}

func Test_savePedUsers_siteUsers(t *testing.T) {

	tests := []struct {
		name                 string
		newUsers             []*entities.SiteUser
		userIdsToDelete      []int
		expectedSiteUsers    []entities.SiteUser
		expectedTidOverrides []entities.SiteUser
		wantErr              bool
	}{
		{
			name:                 "Passing empty lists has no effect on the stored users",
			newUsers:             nil,
			userIdsToDelete:      nil,
			expectedSiteUsers:    initialSiteUsers,
			expectedTidOverrides: initialTidUsers,
			wantErr:              false,
		},
		{
			name:                 "Site users are deleted if they are contained in the given list of user IDs to deleted; removing Bob",
			newUsers:             nil,
			userIdsToDelete:      []int{bob.UserId},
			expectedSiteUsers:    []entities.SiteUser{alice, siteUserWithOverride},
			expectedTidOverrides: initialTidUsers,
			wantErr:              false,
		},
		{
			name:                 "Site users are deleted if they are contained in the given list of user IDs to deleted; removing Alice",
			newUsers:             nil,
			userIdsToDelete:      []int{alice.UserId},
			expectedSiteUsers:    []entities.SiteUser{bob, siteUserWithOverride},
			expectedTidOverrides: initialTidUsers,
			wantErr:              false,
		},
		{
			name:                 "New site users are added if they are contained in the given list of new users",
			newUsers:             []*entities.SiteUser{&newUser1},
			userIdsToDelete:      nil,
			expectedSiteUsers:    []entities.SiteUser{alice, bob, siteUserWithOverride, newUser1},
			expectedTidOverrides: initialTidUsers,
			wantErr:              false,
		},
		{
			name: "Existing site users are updated if they are contained in the given list of new users",
			newUsers: []*entities.SiteUser{{
				UserId:   bob.UserId,
				Username: bob.Username,
				PIN:      bob.PIN,
				Modules:  []string{"newModuleList"},
				SiteId:   bob.SiteId,
			}},
			userIdsToDelete: nil,
			expectedSiteUsers: []entities.SiteUser{alice, {
				UserId:   bob.UserId,
				Username: bob.Username,
				PIN:      bob.PIN,
				Modules:  []string{"newModuleList"},
				SiteId:   bob.SiteId,
			}, siteUserWithOverride},
			expectedTidOverrides: initialTidUsers,
			wantErr:              false,
		},
		{
			name:                 "TID user overrides are deleted if their site users are deleted",
			newUsers:             nil,
			userIdsToDelete:      []int{siteUserWithOverride.UserId},
			expectedSiteUsers:    []entities.SiteUser{alice, bob},
			expectedTidOverrides: []entities.SiteUser{tidUser1, tidUser2},
			wantErr:              false,
		},
		{
			name: "If a user's PIN is updated for a site then that user's PIN is updated for all corresponding overrides",
			newUsers: []*entities.SiteUser{{
				UserId:   siteUserWithOverride.UserId,
				Username: siteUserWithOverride.Username,
				PIN:      "00033",
				Modules:  siteUserWithOverride.Modules,
				SiteId:   7894,
			}},
			expectedSiteUsers: []entities.SiteUser{alice, bob, {
				UserId:   siteUserWithOverride.UserId,
				Username: siteUserWithOverride.Username,
				PIN:      "00033",
				Modules:  siteUserWithOverride.Modules,
				SiteId:   7894,
			}},
			expectedTidOverrides: []entities.SiteUser{tidUser1, tidUser2, {
				UserId:   overriddenSiteUser1.UserId,
				TidId:    overriddenSiteUser1.TidId,
				Username: overriddenSiteUser1.Username,
				PIN:      "00033",
				Modules:  overriddenSiteUser1.Modules,
				SiteId:   7894,
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		dal := &mockSaveSiteUsersDal{
			SiteUsers: initialSiteUsersCopy(),
			TidUsers:  initialTidUsersCopy(),
		}

		t.Run(tt.name, func(t *testing.T) {
			if _, err := saveUserChanges(dal, siteId, tt.newUsers, tt.userIdsToDelete, nil, nil, false); (err != nil) != tt.wantErr {
				t.Errorf("savePedUserChanges() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr == true {
				return
			}

			if !reflect.DeepEqual(dal.SiteUsers, tt.expectedSiteUsers) {
				t.Errorf("Site users not equal to expected;\n got %v,\nwant %v", dal.SiteUsers, tt.expectedSiteUsers)
			}

			if !reflect.DeepEqual(dal.TidUsers, tt.expectedTidOverrides) {
				t.Errorf("TID overrides not equal to expected;\n got %v,\nwant %v", dal.TidUsers, tt.expectedTidOverrides)
			}
		})
	}
}

func Test_savePedUsers_tidUsers(t *testing.T) {
	tests := []struct {
		name               string
		newTidsAndUsers    map[int][]*entities.SiteUser
		tidUserIdsToDelete []int
		expectedTidUsers   []entities.SiteUser
		wantErr            bool
	}{
		{
			name:               "Passing empty lists has not effect on the stored TID users",
			newTidsAndUsers:    nil,
			tidUserIdsToDelete: nil,
			expectedTidUsers:   initialTidUsers,
			wantErr:            false,
		},
		{
			name: "New TID users can be added",
			newTidsAndUsers: map[int][]*entities.SiteUser{
				tid1: {&newTidUser1},
			},
			tidUserIdsToDelete: nil,
			expectedTidUsers:   finalTidUsers,
			wantErr:            false,
		},
		{
			name:               "TID users are deleted if their IDs are in the given list of IDs",
			newTidsAndUsers:    nil,
			tidUserIdsToDelete: []int{tidUser2.UserId},
			expectedTidUsers:   []entities.SiteUser{tidUser1, overriddenSiteUser1},
			wantErr:            false,
		},
	}
	for _, tt := range tests {
		dal := &mockSaveSiteUsersDal{
			SiteUsers: initialSiteUsersCopy(),
			TidUsers:  initialTidUsersCopy(),
		}

		t.Run(tt.name, func(t *testing.T) {
			if _, err := saveUserChanges(dal, siteId, nil, nil, tt.newTidsAndUsers, tt.tidUserIdsToDelete, false); (err != nil) != tt.wantErr {
				t.Errorf("savePedUserChanges() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr == true {
				return
			}

			if !reflect.DeepEqual(dal.TidUsers, tt.expectedTidUsers) {
				t.Errorf("TID overrides not equal to expected; got %v, want %v", dal.TidUsers, tt.expectedTidUsers)
			}
		})
	}
}

func Test_savePedUsers_validation(t *testing.T) {
	tests := []struct {
		name                string
		dal                 *mockSaveSiteUsersDal
		siteUsers           []*entities.SiteUser
		tidsAndUsers        map[int][]*entities.SiteUser
		siteUsersToDelete   []int
		tidUsersToDelete    []int
		validateAgainstTids bool
		responseIsCorrect   func(results []*dal.UserUpdateResult) bool
	}{
		{
			name:                "No validation error is returned when adding a single site user to an empty site",
			validateAgainstTids: false,
			siteUsers:           []*entities.SiteUser{&newUser1},
			responseIsCorrect: func(results []*dal.UserUpdateResult) bool {
				for _, result := range results {
					if !result.Result.Success {
						return false
					}
					if result.Result.ErrorMessage != "" {
						return false
					}
					if result.Action != "Add Site User" {
						return false
					}
				}
				return true
			},
		},
		{
			name:                "A validation error is returned when trying to add two site users with the same name",
			validateAgainstTids: false,
			siteUsers: []*entities.SiteUser{
				{UserId: 1, Username: "User1", PIN: "12345"},
				{UserId: 2, Username: "User1", PIN: "23456"}},
			responseIsCorrect: func(results []*dal.UserUpdateResult) bool {
				for _, result := range results {
					if result.Result.Success {
						return false
					}
					if result.Result.ErrorMessage != "Duplicate user" {
						return false
					}
				}
				return true
			},
		},
		{
			name:                "A validation error is returned when adding two site users with the same PIN at once",
			validateAgainstTids: false,
			siteUsers: []*entities.SiteUser{
				{UserId: 1, Username: "User1", PIN: "12345"},
				{UserId: 2, Username: "User2", PIN: "12345"}},
			responseIsCorrect: func(results []*dal.UserUpdateResult) bool {
				for _, result := range results {
					if result.Result.Success {
						return false
					}
					if result.Result.ErrorMessage != "Duplicate Pins" {
						return false
					}
				}
				return true
			},
		},
		{
			name:                "A validation error is returned when adding a site user and a TID user with different names that both have the same PIN",
			validateAgainstTids: false,
			siteUsers: []*entities.SiteUser{
				{UserId: 0, Username: "User1", PIN: "12345"}},
			tidsAndUsers: map[int][]*entities.SiteUser{
				88880004: {{UserId: 0, Username: "User2", PIN: "12345"}}},
			responseIsCorrect: func(results []*dal.UserUpdateResult) bool {
				for _, result := range results {
					if result.User.TidId == 88880004 && result.Result.Success {
						return false
					}
					if result.User.TidId == 88880004 && result.Result.ErrorMessage != "Cannot save TID user 'User2'; site user 'User1' has the same PIN" {
						return false
					}
				}
				return true
			},
		},
		{
			name:                "A validation error is returned when adding a site user when a site user with the same PIN already exists",
			validateAgainstTids: false,
			siteUsers:           []*entities.SiteUser{{UserId: 0, Username: "User1", PIN: "12345"}},
			dal:                 &mockSaveSiteUsersDal{SiteUsers: []entities.SiteUser{{UserId: 1, Username: "User2", PIN: "12345"}}},
			responseIsCorrect: func(results []*dal.UserUpdateResult) bool {
				for _, result := range results {
					if result.Result.Success {
						return false
					}
					if result.Result.ErrorMessage != "Cannot save site user 'User1'; site user 'User2' has the same PIN" {
						return false
					}

				}
				return true
			},
		},
		{
			name:                "A validation error is returned when adding a site that has the same PIN as an override for a different user",
			validateAgainstTids: false,
			siteUsers:           []*entities.SiteUser{{UserId: 0, Username: "User1", PIN: "12345"}},
			dal: &mockSaveSiteUsersDal{TidUsers: []entities.SiteUser{
				{UserId: 1, TidId: 88880004, Username: "User2", PIN: "12345"}},
			},
			responseIsCorrect: func(results []*dal.UserUpdateResult) bool {
				for _, result := range results {
					if result.Result.Success {
						return false
					}
					if result.Result.ErrorMessage != "Cannot save site user 'User1'; TID user 'User2' has the same PIN" {
						return false
					}
				}
				return true
			},
		},
		{
			name:                "A validation error is returned when adding two TID users for the same TID with the same PIN at once",
			validateAgainstTids: false,
			tidsAndUsers: map[int][]*entities.SiteUser{
				88880004: {
					{UserId: 0, Username: "User1", PIN: "12345"},
					{UserId: 0, Username: "User2", PIN: "12345"}}},
			responseIsCorrect: func(results []*dal.UserUpdateResult) bool {
				for _, result := range results {
					if result.Result.Success {
						return false
					}
					possibleErrorMessages := []string{
						"Cannot save TID user 'User1'; TID user 'User2' has the same PIN",
						"Cannot save TID user 'User2'; TID user 'User1' has the same PIN",
					}
					if !SliceComparisonHelpers.SlicesOfStringContains(possibleErrorMessages, result.Result.ErrorMessage) {
						return false
					}

				}
				return true
			},
		},
		{
			name:                "A validation error is returned when adding two TID users with different TIDs and the same PIN at once",
			validateAgainstTids: false,
			tidsAndUsers: map[int][]*entities.SiteUser{
				88880004: {{UserId: 0, Username: "User1", PIN: "12345"}},
				88880005: {{UserId: 0, Username: "User2", PIN: "12345"}}},
			responseIsCorrect: func(results []*dal.UserUpdateResult) bool {
				for _, result := range results {
					if result.Result.Success {
						return false
					}
					possibleErrorMessages := []string{
						"Cannot save TID user 'User1'; TID user 'User2' has the same PIN",
						"Cannot save TID user 'User2'; TID user 'User1' has the same PIN",
					}
					if !SliceComparisonHelpers.SlicesOfStringContains(possibleErrorMessages, result.Result.ErrorMessage) {
						return false
					}
				}
				return true
			},
		},
		{
			name:                "A validation error is returned when adding a TID user that has the same PIN as an override for a different user",
			validateAgainstTids: true,
			tidsAndUsers: map[int][]*entities.SiteUser{
				88880003: {{UserId: 0, Username: "User1", PIN: "12345"}}},
			dal: &mockSaveSiteUsersDal{TidUsers: []entities.SiteUser{
				{UserId: 1, TidId: 88880004, Username: "User2", PIN: "12345"}}},
			responseIsCorrect: func(results []*dal.UserUpdateResult) bool {
				for _, result := range results {
					if result.Result.Success {
						return false
					}
					if result.Result.ErrorMessage != "Cannot save TID user 'User1'; TID user 'User2' has the same PIN" {
						return false
					}
				}
				return true
			},
		},
		{
			name:                "A validation error is returned when adding a TID user that has the same PIN as a different site user",
			validateAgainstTids: false,
			tidsAndUsers: map[int][]*entities.SiteUser{
				88880003: {{UserId: 0, Username: "User1", PIN: "12345"}}},
			dal: &mockSaveSiteUsersDal{SiteUsers: []entities.SiteUser{
				{UserId: 1, Username: "User2", PIN: "12345"}}},
			responseIsCorrect: func(results []*dal.UserUpdateResult) bool {
				for _, result := range results {
					if result.Result.Success {
						return false
					}
					if result.Result.ErrorMessage != "Cannot save TID user 'User1'; site user 'User2' has the same PIN" {
						return false
					}
				}
				return true
			},
		},
		{
			name:                "A validation error is returned when trying to change a TID override so that it has a different PIN to its parent site user",
			validateAgainstTids: false,
			tidsAndUsers: map[int][]*entities.SiteUser{
				88880003: {{UserId: 0, Username: "User1", PIN: "22222"}}},
			dal: &mockSaveSiteUsersDal{
				SiteUsers: []entities.SiteUser{{UserId: 1, Username: "User1", PIN: "11111"}},
				TidUsers:  []entities.SiteUser{{UserId: 2, TidId: 88880003, Username: "User1", PIN: "11111"}}},
			responseIsCorrect: func(results []*dal.UserUpdateResult) bool {
				for _, result := range results {
					if result.Result.Success {
						return false
					}
					if result.Result.ErrorMessage != "Cannot save TID user 'User1'; PIN is different to the user it is overriding" {
						return false
					}
				}
				return true
			},
		},
		{
			name:                "It's valid to change a TID override's PIN if it has no parent site user",
			validateAgainstTids: false,
			tidsAndUsers: map[int][]*entities.SiteUser{
				88880003: {{UserId: 0, Username: "User1", PIN: "22222"}}},
			dal: &mockSaveSiteUsersDal{
				TidUsers: []entities.SiteUser{{UserId: 2, TidId: 88880003, Username: "User1", PIN: "11111"}}},
			responseIsCorrect: func(results []*dal.UserUpdateResult) bool {
				for _, result := range results {
					if !result.Result.Success {
						return false
					}
					if result.Action != "Add TID User Override" {
						return false
					}
				}
				return true
			},
		},
		{
			name:                "It's valid to add a new site together with a new TID override with the same PIN",
			validateAgainstTids: false,
			siteUsers: []*entities.SiteUser{{
				UserId: -1, Username: "User1", PIN: "11111", Modules: []string{"foo"}}},
			tidsAndUsers: map[int][]*entities.SiteUser{
				88880003: {{UserId: -2, Username: "User1", PIN: "11111", Modules: []string{"foo", "bar"}}}},
			responseIsCorrect: func(results []*dal.UserUpdateResult) bool {
				for _, result := range results {
					if !result.Result.Success {
						return false
					}
					if result.Action != "Add Site User" && result.Action != "Add TID User Override" {
						return false
					}
				}
				return true
			},
		},
		{
			name:                "It's valid to add a new site together with 2 new TID overrides (for different TIDs) with the same PIN",
			validateAgainstTids: false,
			siteUsers: []*entities.SiteUser{{
				UserId: -1, Username: "User1", PIN: "11111", Modules: []string{"foo"}}},
			tidsAndUsers: map[int][]*entities.SiteUser{
				88880003: {{UserId: -2, Username: "User1", PIN: "11111", Modules: []string{"foo", "bar"}}},
				88880004: {{UserId: -3, Username: "User1", PIN: "11111", Modules: []string{"foo", "bar"}}},
			},
			responseIsCorrect: func(results []*dal.UserUpdateResult) bool {
				for _, result := range results {
					if !result.Result.Success {
						return false
					}
					if (result.User.TidId == 88880003 || result.User.TidId == 88880004) && result.User.Username == "User1" && result.Action != "Add TID User Override" {
						return false
					}
					if result.User.SiteId == siteId && result.Action != "Add Site User" {
						return false
					}
				}
				return true
			},
		},
		{
			name:                "A validation error is returned when adding multiple TID overrides with the same TID and username",
			validateAgainstTids: false,
			siteUsers: []*entities.SiteUser{{
				UserId: -1, Username: "User1", PIN: "11111", Modules: []string{"foo"}}},
			tidsAndUsers: map[int][]*entities.SiteUser{
				88880003: {
					{UserId: -2, Username: "User1", PIN: "11111", Modules: []string{"foo"}},
					{UserId: -3, Username: "User1", PIN: "11111", Modules: []string{"foo", "bar"}},
				},
			},
			responseIsCorrect: func(results []*dal.UserUpdateResult) bool {
				for _, result := range results {
					if result.User.TidId == 88880003 && result.Result.Success {
						return false
					}
					if result.User.TidId == 88880003 && result.Result.Success && result.Result.ErrorMessage != "Cannot save TID user 'User1'; TID user 'User1' has the same username on the same TID" {
						return false
					}
					if result.User.SiteId == siteId && result.Action != "Add Site User" {
						return false
					}
				}
				return true
			},
		},
	}

	for _, tt := range tests {
		dal := tt.dal
		if dal == nil {
			dal = &mockSaveSiteUsersDal{}
		}

		t.Run(tt.name, func(t *testing.T) {
			results, err := saveUserChanges(dal, siteId, tt.siteUsers, tt.siteUsersToDelete, tt.tidsAndUsers, tt.tidUsersToDelete, tt.validateAgainstTids)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err.Error())
			} else if !tt.responseIsCorrect(results) {
				t.Fail()
			}
		})
	}
}

func Test_unmarshalUsersUploadCsv(t *testing.T) {

	defaultModules := []string{"sale", "gratuitySale", "X-Read", "Z-Read"}

	tests := []struct {
		name              string
		csvContents       string
		dal               *mockSaveSiteUsersDal
		siteUsers         []*entities.SiteUser
		tidsAndUsers      map[int][]*entities.SiteUser
		siteUsersToDelete []int
		tidUsersToDelete  []int
		expectedError     error
	}{
		{
			name: "Adding site users",
			csvContents: `Export for Merchant ID:,000122430044
Username,PIN,TID,sale,refund,void,preAuthSale,preAuthCompletion,preAuthCancel,gratuitySale,gratuityCompletion,alipaySale,alipayVoid,alipayRefund,upi,X-Read,Z-Read,delete
NewUser1,"=""13746""","=""""",Y,N,N,N,N,N,Y,N,N,N,N,N,Y,Y,N
NewUser2,"=""13747""","=""""",Y,N,N,N,N,N,Y,N,N,N,N,N,Y,Y,N`,
			siteUsers: []*entities.SiteUser{
				{UserId: 0, Username: "NewUser1", PIN: "13746", Modules: defaultModules},
				{UserId: 0, Username: "NewUser2", PIN: "13747", Modules: defaultModules}},
		},
		{
			name: "Updating site users",
			csvContents: `Export for Merchant ID:,000122430044
Username,PIN,TID,sale,refund,void,preAuthSale,preAuthCompletion,preAuthCancel,gratuitySale,gratuityCompletion,alipaySale,alipayVoid,alipayRefund,upi,X-Read,Z-Read,delete
Existing1,"=""12345""","=""""",Y,N,N,N,N,N,Y,N,N,N,N,N,Y,Y,N`,
			dal: &mockSaveSiteUsersDal{
				SiteUsers: []entities.SiteUser{
					{UserId: 2, Username: "Existing1", PIN: "12345", Modules: defaultModules},
					{UserId: 3, Username: "OtherUser", PIN: "23456", Modules: defaultModules}},
			},
			siteUsers: []*entities.SiteUser{
				{UserId: 2, Username: "Existing1", PIN: "12345", Modules: defaultModules}},
		},
		{
			name: "Adding TID users",
			csvContents: `Export for Merchant ID:,000122430044
Username,PIN,TID,sale,refund,void,preAuthSale,preAuthCompletion,preAuthCancel,gratuitySale,gratuityCompletion,alipaySale,alipayVoid,alipayRefund,upi,X-Read,Z-Read
NewUser1,"=""12345""","=""88880004""",Y,N,N,N,N,N,Y,N,N,N,N,N,Y,Y
NewUser2,"=""23456""","=""88880004""",Y,N,N,N,N,N,Y,N,N,N,N,N,Y,Y
NewUser3,"=""34567""","=""88880005""",Y,N,N,N,N,N,Y,N,N,N,N,N,Y,Y
`,
			tidsAndUsers: map[int][]*entities.SiteUser{
				88880004: {
					{UserId: 0, Username: "NewUser1", PIN: "12345", Modules: defaultModules},
					{UserId: 0, Username: "NewUser2", PIN: "23456", Modules: defaultModules}},
				88880005: {
					{UserId: 0, Username: "NewUser3", PIN: "34567", Modules: defaultModules}},
			},
		},
		{
			name: "Updating TID users",
			csvContents: `Export for Merchant ID:,000122430044
Username,PIN,TID,sale,refund,void,preAuthSale,preAuthCompletion,preAuthCancel,gratuitySale,gratuityCompletion,alipaySale,alipayVoid,alipayRefund,upi,X-Read,Z-Read,delete
TidUser2,45678,88880004,Y,N,N,N,N,N,Y,N,N,N,N,N,Y,Y,N`,
			dal: &mockSaveSiteUsersDal{
				SiteUsers: []entities.SiteUser{
					{UserId: 2, Username: "Bob", PIN: "12345", Modules: defaultModules, SiteId: siteId},
					{UserId: 3, Username: "Alice", PIN: "23456", Modules: defaultModules, SiteId: siteId}},
				TidUsers: []entities.SiteUser{
					{UserId: 4, TidId: 88881111, Username: "TidUser1", PIN: "34567", Modules: defaultModules},
					{UserId: 5, TidId: 88880004, Username: "TidUser2", PIN: "45678", Modules: defaultModules},
					{UserId: 6, TidId: 88889999, Username: "TidUser3", PIN: "56789", Modules: defaultModules}},
			},
			tidsAndUsers: map[int][]*entities.SiteUser{
				88880004: {{UserId: 5, Username: "TidUser2", PIN: "45678", Modules: defaultModules}},
			},
		},
		{
			name: "Deleting site users",
			csvContents: `Export for Merchant ID:,000122430044
Username,PIN,TID,sale,refund,void,preAuthSale,preAuthCompletion,preAuthCancel,gratuitySale,gratuityCompletion,alipaySale,alipayVoid,alipayRefund,upi,X-Read,Z-Read,delete
Bob1,"=""12345""","=""""",Y,N,N,N,N,N,Y,N,N,N,N,N,Y,Y,Y
Bob2,"=""23456""","=""""",Y,N,N,N,N,N,Y,N,N,N,N,N,Y,Y,Y
UserX,"=""99999""","=""""",Y,N,N,N,N,N,Y,N,N,N,N,N,Y,Y,Y`,
			dal: &mockSaveSiteUsersDal{
				SiteUsers: []entities.SiteUser{
					{UserId: 12, Username: "Bob1", PIN: "12345", Modules: defaultModules},
					{UserId: 13, Username: "Bob2", PIN: "23456", Modules: defaultModules},
				},
				TidUsers: nil,
			},
			// Won't have an ID for 'UserX' since it's not in the database
			siteUsersToDelete: []int{12, 13},
		},
		{
			name: "Deleting TID users",
			csvContents: `Export for Merchant ID:,000122430044
Username,PIN,TID,sale,refund,void,preAuthSale,preAuthCompletion,preAuthCancel,gratuitySale,gratuityCompletion,alipaySale,alipayVoid,alipayRefund,upi,X-Read,Z-Read,delete
Bob1,"=""13746""",88880004,Y,N,N,N,N,N,Y,N,N,N,N,N,Y,Y,Y
Bob2,"=""13747""",88880004,Y,N,N,N,N,N,Y,N,N,N,N,N,Y,Y,Y
UserX,"=""99999""",88880005,Y,N,N,N,N,N,Y,N,N,N,N,N,Y,Y,Y
`,
			dal: &mockSaveSiteUsersDal{
				SiteUsers: nil,
				TidUsers: []entities.SiteUser{
					{UserId: 12, TidId: 88880004, Username: "Bob1", PIN: "13746", Modules: defaultModules},
					{UserId: 13, TidId: 88880004, Username: "Bob2", PIN: "13747", Modules: defaultModules}},
			},
			// Won't have an ID for 'UserX' since it's not in the database
			tidUsersToDelete: []int{12, 13},
		},
		{
			name: "Omission of username field results in an error",
			csvContents: `Export for Merchant ID:,000122430044
PIN,TID,sale,refund,void,preAuthSale,preAuthCompletion,preAuthCancel,gratuitySale,gratuityCompletion,alipaySale,alipayVoid,alipayRefund,upi,X-Read,Z-Read,delete
13746""",88880004,Y,N,N,N,N,N,Y,N,N,N,N,N,Y,Y,Y
13747""",88880004,Y,N,N,N,N,N,Y,N,N,N,N,N,Y,Y,Y`,
			expectedError: errors.New("invalid csv format"),
		},
		{
			name: "Omission of PIN field results in an error",
			csvContents: `Export for Merchant ID:,000122430044
Username,TID,sale,refund,void,preAuthSale,preAuthCompletion,preAuthCancel,gratuitySale,gratuityCompletion,alipaySale,alipayVoid,alipayRefund,upi,X-Read,Z-Read,delete
Bob1,88880004,Y,N,N,N,N,N,Y,N,N,N,N,N,Y,Y,Y
Bob2,88880004,Y,N,N,N,N,N,Y,N,N,N,N,N,Y,Y,Y`,
			expectedError: errors.New("invalid csv format"),
		},
		{
			name: "Omission of TID field results in an error",
			csvContents: `Export for Merchant ID:,000122430044
Username,PIN,sale,refund,void,preAuthSale,preAuthCompletion,preAuthCancel,gratuitySale,gratuityCompletion,alipaySale,alipayVoid,alipayRefund,upi,X-Read,Z-Read,delete
Bob1,"=""13746""",Y,N,N,N,N,N,Y,N,N,N,N,N,Y,Y,Y
Bob2,"=""13747""",Y,N,N,N,N,N,Y,N,N,N,N,N,Y,Y,Y`,
			expectedError: errors.New("invalid csv format"),
		},
		{
			name: "Max username length is 10; username of length 10",
			csvContents: `Export for Merchant ID:,000122430044
Username,PIN,TID,sale,refund,void,preAuthSale,preAuthCompletion,preAuthCancel,gratuitySale,gratuityCompletion,alipaySale,alipayVoid,alipayRefund,upi,X-Read,Z-Read,delete
1234567890,"=""13746""",88880004,Y,N,N,N,N,N,Y,N,N,N,N,N,Y,Y,Y`,
			expectedError: nil,
		},
		{
			name: "Max username length is 10; username of length 11",
			csvContents: `Export for Merchant ID:,000122430044
Username,PIN,TID,sale,refund,void,preAuthSale,preAuthCompletion,preAuthCancel,gratuitySale,gratuityCompletion,alipaySale,alipayVoid,alipayRefund,upi,X-Read,Z-Read,delete
12345678901,"=""13746""",88880004,Y,N,N,N,N,N,Y,N,N,N,N,N,Y,Y,Y`,
			expectedError: errors.New("Username must not be more than 10 characters"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.siteUsers == nil {
				tt.siteUsers = make([]*entities.SiteUser, 0)
			}

			if tt.tidsAndUsers == nil {
				tt.tidsAndUsers = make(map[int][]*entities.SiteUser, 0)
			}

			if tt.siteUsersToDelete == nil {
				tt.siteUsersToDelete = make([]int, 0)
			}

			if tt.tidUsersToDelete == nil {
				tt.tidUsersToDelete = make([]int, 0)
			}

			reader := strings.NewReader(tt.csvContents)
			dal := tt.dal
			if dal == nil {
				dal = &mockSaveSiteUsersDal{}
			}

			siteUsers, tidsAndUsers, siteUsersToDelete, tidUsersToDelete, err := unmarshalUsersUploadCsv(dal, reader, siteId)
			if tt.expectedError == nil && err != nil {
				t.Fatalf("Failed to unmarshal CSV: %v", err.Error())
			}

			if tt.expectedError != nil {
				if !reflect.DeepEqual(err, tt.expectedError) {
					t.Errorf("Error not equal to expected; got %v, want %v", err, tt.expectedError)
				}
				return
			}

			if !reflect.DeepEqual(siteUsers, tt.siteUsers) {
				t.Errorf("Site users not equal to expected; got %v, want %v", siteUsers, tt.siteUsers)
			}

			if !reflect.DeepEqual(tidsAndUsers, tt.tidsAndUsers) {
				t.Errorf("TID users not equal to expected; got %v, want %v", tidsAndUsers, tt.tidsAndUsers)
			}

			if !reflect.DeepEqual(siteUsersToDelete, tt.siteUsersToDelete) {
				t.Errorf("Site users to delete not equal to expected; got %v, want %v", siteUsersToDelete, tt.siteUsersToDelete)
			}

			if !reflect.DeepEqual(tidUsersToDelete, tt.tidUsersToDelete) {
				t.Errorf("TID users to delete not equal to expected; got %v, want %v", tidUsersToDelete, tt.tidUsersToDelete)
			}
		})
	}
}

func TestIsAllUpperCase(t *testing.T) {
	tests := []struct {
		name        string
		inputString string
		expected    bool
	}{
		{"All caps", "ALLCAPS", true},
		{"All lower case", "alllower", false},
		{"Mixed case", "MixedCase", false},
		{"Mixed caps and numbers", "1ALLC4PS", true},
		{"Mixed lower and numbers", "1alll0wer", false},
		{"Mixed case and numbers", "1MixedC4se", false},
		{"All numbers", "123456", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isAllUpperCase(tt.inputString)
			if result != tt.expected {
				t.Errorf("isAllUpperCase() = %v, want %v", result, tt.expected)
			}
		})
	}
}
