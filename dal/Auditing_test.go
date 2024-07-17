package dal

import (
	le "bitbucket.org/network-international/nextgen-libs/nextgen-tg-protobuf/Logging"
	txn "bitbucket.org/network-international/nextgen-libs/nextgen-tg-protobuf/Transaction"
	sharedDAL "bitbucket.org/network-international/nextgen-tms/web-shared/dal"
	"nextgen-tms-website/entities"
	"reflect"
	"testing"
	"time"
)

var (
	oneFoo    = []string{"foo"}
	oneBar    = []string{"bar"}
	twoFooBar = []string{"foo", "bar"}
	twoBarFoo = []string{"bar", "foo"}
	twoBarBar = []string{"bar", "bar"}

	emptyAudit = txn.UserAuditHistory{
		Acquirer:      "",
		Name:          "",
		Module:        "",
		OriginalValue: "",
		UpdatedValue:  "",
		UpdatedBy:     "",
	}

	singleAudit = txn.UserAuditHistory{
		Acquirer:      "foo",
		Name:          "name",
		Module:        "module",
		OriginalValue: "foo",
		UpdatedValue:  "foo",
		UpdatedBy:     "user",
	}

	multiAudit = txn.UserAuditHistory{
		Acquirer:      "foo, bar",
		Name:          "name",
		Module:        "module",
		OriginalValue: "foo. bar",
		UpdatedValue:  "foo. bar",
		UpdatedBy:     "user",
	}

	failedAddAudit = txn.UserAuditHistory{
		Acquirer:      "",
		Name:          "Test Group",
		Module:        "group",
		OriginalValue: "",
		UpdatedValue:  "Error adding user group Test Group",
		UpdatedBy:     "Dummy User Name",
	}

	succeedAddAudit = txn.UserAuditHistory{
		Acquirer:      "",
		Name:          "Test Group",
		Module:        "group",
		OriginalValue: "",
		UpdatedValue:  "User group Test Group added",
		UpdatedBy:     "Dummy User Name",
	}

	dummyUser = entities.TMSUser{
		UserID:          0,
		Username:        "Dummy User Name",
		Token:           "",
		RoleID:          0,
		LoggedIn:        false,
		UserPermissions: entities.UserPermissions{},
		Expires:         time.Time{},
	}

	singleAuditLE = le.UserManagementAudit{
		Acquirer:      "foo",
		Name:          "name",
		Module:        "module",
		OriginalValue: "foo",
		UpdatedValue:  "foo",
		UpdatedBy:     "user",
		UpdatedAt:     time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	singleFriendlyAuditLE = sharedDAL.UserAuditDisplay{
		Acquirer:      "foo",
		Name:          "name",
		Module:        "module",
		OriginalValue: []string{"foo"},
		UpdatedValue:  []string{"foo"},
		UpdatedBy:     "user",
		UpdatedAt:     time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
	}

	multiAuditLE = le.UserManagementAudit{
		Acquirer:      "foo",
		Name:          "name",
		Module:        "module",
		OriginalValue: "foo.bar",
		UpdatedValue:  "foo.bar.test",
		UpdatedBy:     "user",
		UpdatedAt:     time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	multiFriendlyAuditLE = sharedDAL.UserAuditDisplay{
		Acquirer:      "foo",
		Name:          "name",
		Module:        "module",
		OriginalValue: []string{"foo", "bar"},
		UpdatedValue:  []string{"foo", "bar", "test"},
		UpdatedBy:     "user",
		UpdatedAt:     time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
	}
)

func TestConvertToDisplayFriendly(t *testing.T) {
	tests := []struct {
		name     string
		audits   []*le.UserManagementAudit
		expected []sharedDAL.UserAuditDisplay
	}{
		{"Nil input", nil, nil},
		{"Single Input, Single Value", []*le.UserManagementAudit{&singleAuditLE}, []sharedDAL.UserAuditDisplay{singleFriendlyAuditLE}},
		{"Single Input, Multiple Value", []*le.UserManagementAudit{&multiAuditLE}, []sharedDAL.UserAuditDisplay{multiFriendlyAuditLE}},
		{"Multiple Input, Multiple Value", []*le.UserManagementAudit{&singleAuditLE, &multiAuditLE}, []sharedDAL.UserAuditDisplay{singleFriendlyAuditLE, multiFriendlyAuditLE}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertToDisplayFriendly(tt.audits)

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ConvertToDisplayFriendly() = %v, want %v", got, tt.expected)
			}
		})
	}
}

//Builds a user audit object using the given attributes
func TestBuildAuditEntry(t *testing.T) {
	tests := []struct {
		name       string
		acquirers  []string
		moduleName string
		module     string
		original   []string
		updated    []string
		tmsUser    string
		expected   txn.UserAuditHistory
	}{
		{"All fields empty", nil, "", "", nil, nil, "", emptyAudit},
		{"All arrays single element", oneFoo, "name", "module", oneFoo, oneFoo, "user", singleAudit},
		{"All arrays multi element", twoFooBar, "name", "module", twoFooBar, twoFooBar, "user", multiAudit},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildAuditEntry(tt.acquirers, tt.moduleName, tt.module, tt.original, tt.updated, tt.tmsUser)
			if got.Acquirer != tt.expected.Acquirer {
				t.Errorf("BuildAuditEntry() Acquirer = %v, want %v", got.Acquirer, tt.expected.Acquirer)
			}
			if got.Name != tt.expected.Name {
				t.Errorf("BuildAuditEntry() Name = %v, want %v", got.Name, tt.expected.Name)
			}
			if got.Module != tt.expected.Module {
				t.Errorf("BuildAuditEntry() Module = %v, want %v", got.Module, tt.expected.Module)
			}
			if got.OriginalValue != tt.expected.OriginalValue {
				t.Errorf("BuildAuditEntry() Original = %v, want %v", got.OriginalValue, tt.expected.OriginalValue)
			}
			if got.UpdatedValue != tt.expected.UpdatedValue {
				t.Errorf("BuildAuditEntry() Updated = %v, want %v", got.UpdatedValue, tt.expected.UpdatedValue)
			}
			if got.UpdatedBy != tt.expected.UpdatedBy {
				t.Errorf("BuildAuditEntry() UpdatedBy = %v, want %v", got.UpdatedBy, tt.expected.UpdatedBy)
			}
		})
	}
}

//Compares the values of two string arrays and returns an array of elements found in the 1st but not the 2nd
func TestAddedOrRemovedElements(t *testing.T) {
	tests := []struct {
		name     string
		original []string
		new      []string
		expected []string
	}{
		{"Two nil arrays", nil, nil, nil},
		{"Original nil, new populated", nil, oneBar, nil},
		{"Original single, new nil", oneBar, nil, oneBar},
		{"Original multiple, new nil", twoFooBar, nil, twoFooBar},
		{"Original single, new identical", oneBar, oneBar, nil},
		{"Original multiple, new identical", twoFooBar, twoFooBar, nil},
		{"Original multiple, new same contents different order", twoFooBar, twoBarFoo, nil},
		{"Original multiple, new single", twoBarFoo, oneBar, oneFoo},
		{"Original multiple, new multiple different contents", twoFooBar, twoBarBar, oneFoo},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AddedOrRemovedElements(tt.original, tt.new)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("AddedOrRemovedElements() = %v, want %v", got, tt.expected)
			}
		})
	}
}

//Compares two string arrays to see if the content differs. Order is not important. TRUE = different
func TestCompareGroupElements(t *testing.T) {
	tests := []struct {
		name     string
		old      []string
		new      []string
		expected bool
		removed  []string
		added    []string
	}{
		{"Old nil New not nil", nil, oneFoo, true, nil, oneFoo},
		{"Old not nil New nil", oneFoo, nil, true, oneFoo, nil},
		{"Old nil New nil", nil, nil, false, nil, nil},
		{"Old len 2 New len 1", twoFooBar, oneFoo, true, oneBar, nil},
		{"Old len 1 New len 2", oneFoo, twoFooBar, true, nil, oneBar},
		{"Old and New identical", twoFooBar, twoFooBar, false, nil, nil},
		{"Old and New same content different order", twoFooBar, twoBarFoo, false, nil, nil},
		{"Old and New different content", twoFooBar, twoBarBar, true, oneFoo, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotAdded, gotRemoved := CompareGroupElements(tt.old, tt.new)
			if got != tt.expected {
				t.Errorf("CompareGroupElements() = %v, want %v", got, tt.expected)
			}
			if !reflect.DeepEqual(gotRemoved, tt.removed) {
				t.Errorf("CompareGroupElements() = %v, want %v", gotRemoved, tt.removed)
			}
			if !reflect.DeepEqual(gotAdded, tt.added) {
				t.Errorf("CompareGroupElements() = %v, want %v", gotAdded, tt.added)
			}
		})
	}
}
