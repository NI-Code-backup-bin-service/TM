package authentication

import (
	ConfigurationSupport "bitbucket.org/network-international/nextgen-libs/nextgen-helpers/configHelper"
	"errors"
	"gopkg.in/ldap.v2"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestModifyPasswordIntegrationTest(t *testing.T) {
	/*
	 In order to test this against UAT (in particular so that we can test against Active Directory, as opposed to
	 OpenLDAP which we use on AWS) we can make use of golang’s capability to build an executable binary containing
	 unit tests. Using the script UatLdapTestUpload.sh contained in the same directory as this file, an EXE can be built,
	 uploaded to Artifactory, downloaded onto the UAT ansible box, then executed on that box.
	 This allows testing of the code while not affecting the existing TMS deployment,
	 although the configured test user to be modified needs to be cleared with NI first in case anything goes wrong,
	 e.g. the password gets changed to some unknown/unrecoverable value.
	 */
	t.Skip("Integration test; comment out to use")

	ldap.DefaultTimeout = 3 * time.Second

	tests := []struct {
		name            string
		settingsFn      func()
		userName        string
		currentPassword string
		newPassword     string
		wantErr         bool
	}{
		{
			name:            "TestLocal",
			settingsFn:      setSettingsForLocal,
			userName:        "NISuper",
			currentPassword: "iYnloK1JneDdkTkZyelB",
			newPassword:     "iYnloK1JneDdkTkZyelB",
			wantErr:         false,
		},
		{
			name:            "TestUAT",
			settingsFn:	     setSettingsForUAT,
			// Enter test username & passwords to change; of course we will be changing actual accounts on UAT here,
			// so contact NI in order to obtain some cleared account(s) which you can use for this
			userName:        "INTUTEST1",
			currentPassword: "Welcome1$",
			newPassword:     "Welcome1$",
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.settingsFn()

			if err := ModifyPassword(tt.userName, tt.currentPassword, tt.newPassword); (err != nil) != tt.wantErr {
				t.Errorf("ModifyPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLogon(t *testing.T) {
	/*
	 In order to test this against UAT (in particular so that we can test against Active Directory, as opposed to
	 OpenLDAP which we use on AWS) we can make use of golang’s capability to build an executable binary containing
	 unit tests. Using the script UatLdapTestUpload.sh contained in the same directory as this file, an EXE can be built,
	 uploaded to Artifactory, downloaded onto the UAT ansible box, then executed on that box.
	 This allows testing of the code while not affecting the existing TMS deployment.
	*/
	t.Skip("Integration test; comment out to use")

	ldap.DefaultTimeout = 3 * time.Second

	tests := []struct {
		name            string
		settingsFn      func()
		userName        string
		password string
		wantErr         bool
	}{
		{
			name:       "TestUatLogin",
			settingsFn: setSettingsForUAT,
			userName:   "BINABCU1",
			password:   "P@ssw0rd123",
			wantErr:    false,
		},
		//{
		//	name:       "TestUatLogin2",
		//	settingsFn: setSettingsForUAT,
		//	userName:   "BINABCU2",
		//	password:   "Welcome1$",
		//	wantErr:    false,
		//},
	}
	for _, tt := range tests {
		tt.settingsFn()

		t.Run(tt.name, func(t *testing.T) {
			_, _, err := ValidateUser(tt.userName, tt.password)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func setSettingsForLocal() {
	settings.UseTLS = false
	settings.UseLDAPS = false
	settings.UseMutualAuthLDAPS = false
	settings.IPAddress = "localhost"
	settings.Port = "389"

	settings.Username = "CN=admin,DC=NETWORKINTL,DC=COM"
	settings.Password = "admin"
	settings.BaseDN = "cn=Users,dc=NETWORKINTL,dc=COM"
	settings.ObjectClass = "person"
	settings.IdentityAttribute = "cn"
	settings.UserGroup = "cn=TMSUsers"
	settings.AdminGroup = "cn=TMSAdministrators"
	settings.GlobalAdminGroup = "cn=TMSGlobalAdmins"
}

func setSettingsForUAT() {
	settings.UseTLS = true
	settings.UseLDAPS = false
	settings.UseMutualAuthLDAPS = false
	settings.IPAddress = "10.119.16.51"
	settings.Port = "389"

	settings.Username = "CN=next gen,OU=Service Accounts,DC=UAT,DC=NETWORK,DC=COM"
	settings.Password = strings.TrimSpace(ConfigurationSupport.Decrypt(ConfigurationSupport.GetKey(), "HZMoWVhc-49K8uWSDlYlANmUcXA2nRx1eRI="))

	settings.BaseDN = "dc=UAT,dc=NETWORK,dc=COM"
	settings.ObjectClass = "user"
	settings.IdentityAttribute = "sAMAccountName"
	settings.UserGroup = "CN=TMSUsers"
	settings.AdminGroup = "CN=TMSAdministrators"
	settings.GlobalAdminGroup = "CN=TMSGlobalAdmins"
}

func Test_modifyPasswordForActiveDirectoryServer_successCase(t *testing.T) {
	l := &MockLDAPConnection{}

	if err := modifyPasswordUsingConnection(l, "testUser", "password1", "password2"); err != nil {
		t.Errorf("modifyPasswordUsingConnection() error = %v", err)
	}

	modifyRequests := l.ModifyRequestsReceived
	if len(modifyRequests) != 1 {
		t.Fatalf("Expected a modify request to have been received; actual: %v", len(modifyRequests))
	}

	actualModifyRequest := modifyRequests[0]
	expectedModifyRequest := &ldap.ModifyRequest{
		DN: "cn=testUser,cn=Users,dc=NETWORKINTL,dc=COM",
		// The contents of unicodePwd have to be enclosed in quotes and encoded in utf-16, little-endian
		DeleteAttributes: []ldap.PartialAttribute{
			{"unicodePwd", []string{"\"\000p\000a\000s\000s\000w\000o\000r\000d\0001\000\"\000"}}},
		AddAttributes: []ldap.PartialAttribute{
			{"unicodePwd", []string{"\"\000p\000a\000s\000s\000w\000o\000r\000d\0002\000\"\000"}}},
	}
	if !reflect.DeepEqual(expectedModifyRequest, modifyRequests[0]) {
		t.Errorf("Expected modify request: %v\ngot\n%v", expectedModifyRequest, actualModifyRequest)
	}
}

func Test_modifyPasswordForActiveDirectoryServer_passwordConstraintViolation_returnsCustomError(t *testing.T) {
	l := &MockLDAPConnection{ModifyRequestErrorToReturn:
	ldap.NewError(19, errors.New(""+
		"LDAP Result Code 19 \"Constraint Violation\": 0000052D: AtrErr: DSID-03191083, #1:\n\n"+
		"                0: 0000052D: DSID-03191083, problem 1005 (CONSTRAINT_ATT_TYPE), data 0, Att 9005a (unicodePwd)"))}

	var err error
	if err = modifyPasswordUsingConnection(l, "testUser", "password1", "password2"); err == nil {
		t.Fatal("Expecting an error to have been returned")
	}

	if errType, ok := err.(*PasswordConstraintError); !ok {
		t.Errorf("Expecting an error of type PasswordConstraintError, but was '%s'", errType)
	}
}

func Test_modifyPasswordForActiveDirectoryServer_passwordExpirationError_stillPerformsPasswordModification(t *testing.T) {
	l := &MockLDAPConnection{BindFn: func(userDn, password string) error {
		// Only throw and error for the user whose password is changed, because prior to the Bind call for that user
		// a Bind call will be made for the admin user
		if userDn == "cn=testUser,cn=Users,dc=NETWORKINTL,dc=COM" {
			return createPasswordExpirationError()
		} else {
			return nil
		}
	}}

	if err := modifyPasswordUsingConnection(l, "testUser", "password1", "password2"); err != nil {
		t.Fatalf("modifyPasswordUsingConnection() error = %v", err)
	}

	modifyRequests := l.ModifyRequestsReceived
	if len(modifyRequests) != 1 {
		t.Fatalf("Expected a modify request to have been received; actual: %v", len(modifyRequests))
	}
}

func Test_modifyPasswordForActiveDirectoryServer_miscellaneousLdapError_returnsErrorAsIs(t *testing.T) {
	errorText := "misc error"
	l := &MockLDAPConnection{ModifyRequestErrorToReturn: ldap.NewError(123, errors.New(errorText))}

	var err error
	if err = modifyPasswordUsingConnection(l, "testUser", "password1", "password2"); err == nil {
		t.Fatal("Expecting an error to have been returned")
	}

	if err.Error() != "LDAP Result Code 123 \"\": misc error" {
		t.Errorf("Expecting error text to be '%v', but was '%v'", errorText, err.Error())
	}
}

func Test_modifyPasswordForActiveDirectoryServer_miscellaneousError_returnsErrorAsIs(t *testing.T) {
	errorText := "misc error"
	l := &MockLDAPConnection{ModifyRequestErrorToReturn: errors.New(errorText)}

	var err error
	if err = modifyPasswordUsingConnection(l, "testUser", "password1", "password2"); err == nil {
		t.Fatal("Expecting an error to have been returned")
	}

	if err.Error() != errorText {
		t.Errorf("Expecting error text to be '%v', but was '%v'", errorText, err.Error())
	}
}

func Test_modifyPasswordForOpenLDAPServer_fallsBackSuccessfully(t *testing.T) {
	l := &MockLDAPConnection{ModifyRequestErrorToReturn:
	ldap.NewError(17, errors.New(""+
		"LDAP Result Code 17 \"Undefined Attribute Type\": unicodePwd: attribute type undefined"))}

	if err := modifyPasswordUsingConnection(l, "testUser", "password1", "password2"); err != nil {
		t.Errorf("modifyPasswordUsingConnection() error = %v", err)
	}

	modifyRequests := l.PasswordModifyRequestsReceived
	if len(modifyRequests) != 1 {
		t.Fatalf("Expected a password modify request to have been received; actual: %v", len(modifyRequests))
	}

	actualModifyRequest := modifyRequests[0]
	expectedModifyRequest := &ldap.PasswordModifyRequest{
		OldPassword: "password1",
		NewPassword: "password2",
	}

	if !reflect.DeepEqual(expectedModifyRequest, modifyRequests[0]) {
		t.Errorf("Expected modify request: %v\ngot\n%v", expectedModifyRequest, actualModifyRequest)
	}
}

func Test_errorIsPasswordExpirationError(t *testing.T) {
	err := createPasswordExpirationError()

	if result := isPasswordExpirationError(err); !result {
		t.Error("Expected 'isPasswordExpirationError' to return 'true'; got 'false'")
	}
}

func createPasswordExpirationError() error {
	return ldap.NewError(49, errors.New("80090308: LdapErr: DSID-0C090446, comment: AcceptSecurityContext error, data 532, v2580"))
}

func Test_errorIsPasswordMustChangeError(t *testing.T) {
	err := ldap.NewError(49, errors.New("80090308: LdapErr: DSID-0C090446, comment: AcceptSecurityContext error, data 773, v2580"))

	if result := isPasswordMustChangeError(err); !result {
		t.Error("Expected 'isPasswordMustChangeError' to return 'true'; got 'false'")
	}
}

type MockLDAPConnection struct {
	BindFn                         func(username, password string) error
	ModifyRequestErrorToReturn     error
	ModifyRequestsReceived         []*ldap.ModifyRequest
	PasswordModifyRequestsReceived []*ldap.PasswordModifyRequest
}

func (c *MockLDAPConnection) Bind(username, password string) error {
	if c.BindFn == nil {
		return nil
	} else {
		return c.BindFn(username, password)
	}
}

func (c *MockLDAPConnection) Modify(modifyRequest *ldap.ModifyRequest) error {
	if c.ModifyRequestErrorToReturn == nil {
		c.ModifyRequestsReceived = append(c.ModifyRequestsReceived, modifyRequest)
		return nil
	} else {
		return c.ModifyRequestErrorToReturn
	}
}

func (c *MockLDAPConnection) PasswordModify(passwordModifyRequest *ldap.PasswordModifyRequest) (*ldap.PasswordModifyResult, error) {
	c.PasswordModifyRequestsReceived = append(c.PasswordModifyRequestsReceived, passwordModifyRequest)

	return nil, nil
}

func (c *MockLDAPConnection) Search(_ *ldap.SearchRequest) (*ldap.SearchResult, error) {
	sr := &ldap.SearchResult{
		Entries: []*ldap.Entry{{
			DN: "cn=testUser,cn=Users,dc=NETWORKINTL,dc=COM",
		}},
	}

	return sr, nil
}