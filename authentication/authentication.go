package authentication

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"
	"nextgen-tms-website/dal"
	"nextgen-tms-website/entities"
	"strconv"
	"strings"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/text/encoding/unicode"
	"gopkg.in/ldap.v2"
)

type LDAPSettings struct {
	IPAddress                   string
	Port                        string
	Username                    string
	Password                    string
	BaseDN                      string
	ObjectClass                 string
	SearchAttribute             string
	IdentityAttribute           string
	UseTLS                      bool
	AuthenticateUserPermissions bool
	UserGroup                   string
	AdminGroup                  string
	GlobalAdminGroup            string
	UseLDAPS                    bool
	LDAPSClientCertPath         string
	LDAPSClientKeyPath          string
	UseMutualAuthLDAPS          bool
}

var settings LDAPSettings

func SetLDAPSettings(s LDAPSettings) {
	settings = s
}

func getClientCertificates() ([]tls.Certificate, error) {
	cert, err := tls.LoadX509KeyPair(settings.LDAPSClientCertPath, settings.LDAPSClientKeyPath)
	if err != nil {
		return nil, err
	}

	certificates := make([]tls.Certificate, 1)
	certificates[0] = cert

	return certificates, nil
}

func GetLdapConn() (*ldap.Conn, error) {
	var l *ldap.Conn
	var err error
	if settings.UseTLS && settings.UseLDAPS {

		var clientCerts []tls.Certificate
		if settings.UseMutualAuthLDAPS {
			clientCerts, err = getClientCertificates()
			if err != nil {
				return nil, fmt.Errorf("error retreiving TLS certificate information")
			}
		} else {
			clientCerts = nil
		}

		l, err = ldap.DialTLS("tcp", fmt.Sprintf("%s:%s", settings.IPAddress, settings.Port), &tls.Config{
			InsecureSkipVerify: true,
			Certificates:       clientCerts})
		if err != nil {
			if err != nil {
				return nil, fmt.Errorf("error creating LDAPS connection")
			}
		}
	} else {
		l, err = ldap.Dial("tcp", fmt.Sprintf("%s:%s", settings.IPAddress, settings.Port))

		if err != nil {
			return nil, fmt.Errorf("Error connecting to LDAP Server: %s:%s", settings.IPAddress, settings.Port)
		}

		if settings.UseTLS {
			//Reconnect with TLS
			err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
			if err != nil {
				return nil, fmt.Errorf("error reconnecting with TLS")
			}
		}
	}

	return l, nil
}

// Indicates that a password constraint error was returned from Active Directory.
// This can be caused by a variety of reasons beyond the complexity requirements of
// the password itself, e.g. minimum password age or password history constraints.
// See https://activedirectorypro.com/how-to-configure-a-domain-password-policy/
// for details of the settings involved.
type PasswordConstraintError struct {
	Message string
}

func (e *PasswordConstraintError) Error() string {
	return e.Message
}

// Interfaces with the LDAP server to update a user's password
func ModifyPassword(username, oldPassword, newPassword string) error {
	l, err := GetLdapConn()
	if err != nil {
		return err
	}

	defer l.Close()

	err = modifyPasswordUsingConnection(l, username, oldPassword, newPassword)
	return err
}

type ldapConnection interface {
	Bind(username, password string) error
	Modify(modifyRequest *ldap.ModifyRequest) error
	PasswordModify(passwordModifyRequest *ldap.PasswordModifyRequest) (*ldap.PasswordModifyResult, error)
	Search(searchRequest *ldap.SearchRequest) (*ldap.SearchResult, error)
}

func modifyPasswordUsingConnection(l ldapConnection, username, oldPassword, newPassword string) error {
	sr, err := lookupUserByUsername(l, username)
	if err != nil {
		return err
	}

	userDn := sr.Entries[0].DN

	// Bind as the user whose password we want to change (since they should have permission to change their own password)
	// If this bind operation returns a 'password expired' error, then we should carry on since the user still needs
	// to change their password, and that error can only be received if the correct password was entered.
	err = l.Bind(userDn, oldPassword)
	if err != nil {
		if isErrorRequiringPasswordChange(err) {
			// If the password expired or is forced to change, we need to bind again as the admin user;
			// this is because if we don't then Active Directory would respond with the following error:
			// LDAP Result Code 1 "Operations Error": 000004DC: LdapErr: DSID-0C090D02, comment: "In order to perform
			//     this operation a successful bind must be completed on the connection., data 0, v2580"
			err := l.Bind(settings.Username, settings.Password)
			if err != nil {
				return errors.New("error binding default user")
			}
		} else {
			return err
		}
	}

	err = modifyPasswordForActiveDirectory(l, oldPassword, newPassword, userDn)
	if err != nil {
		if ldapErr, ok := err.(*ldap.Error); ok &&
			ldapErr.ResultCode == ldap.LDAPResultUndefinedAttributeType {

			err = modifyPasswordForOpenLDAP(l, oldPassword, newPassword)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

/**
 * Attempt a password modification operation suitable for Windows Activity Directory (which is what NI use in
 * UAT/preprod/prod.
 */
func modifyPasswordForActiveDirectory(l ldapConnection, oldPassword string, newPassword string, userDn string) error {
	/*
	 * Sadly AD does not support the PasswordModify function provided by the LDAP library, so we have to
	 * do it a bit more manually.
	 *
	 * The implementation here has been derived from https://github.com/go-ldap/ldap/issues/106, in combination with
	 * the docs at https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-adts/6e803168-f140-4d23-b2d3-c3a8ab5917d2
	 */

	oldPwdEncoded, err := encodePasswordForActiveDirectory(oldPassword)
	if err != nil {
		return err
	}

	newPwdEncoded, err := encodePasswordForActiveDirectory(newPassword)
	if err != nil {
		return err
	}

	passReq := &ldap.ModifyRequest{
		DN: userDn,
		DeleteAttributes: []ldap.PartialAttribute{
			{"unicodePwd", []string{oldPwdEncoded}},
		},
		AddAttributes: []ldap.PartialAttribute{
			{"unicodePwd", []string{newPwdEncoded}},
		},
	}

	err = l.Modify(passReq)
	if err != nil {
		if ldapErr, ok := err.(*ldap.Error); ok {
			if ldapErr.ResultCode == ldap.LDAPResultConstraintViolation &&
				strings.Contains(ldapErr.Error(), "unicodePwd") {
				return &PasswordConstraintError{ldapErr.Error()}
			} else {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

func modifyPasswordForOpenLDAP(l ldapConnection, oldPassword string, newPassword string) error {
	passwordModifyRequest := ldap.NewPasswordModifyRequest("", oldPassword, newPassword)
	_, err := l.PasswordModify(passwordModifyRequest)
	if err != nil {
		return err
	}

	return nil
}

func encodePasswordForActiveDirectory(password string) (string, error) {
	utf16 := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	// NB the password needs to be enclosed in quotes
	pwdEncoded, err := utf16.NewEncoder().String("\"" + password + "\"")
	if err != nil {
		return "", err
	}

	return pwdEncoded, nil
}

func ValidateUser(username, password string) (*entities.TMSUser, int, error) {
	// This will only change if a user's password has expired or if they are logging in for the first time
	passwordStatus := entities.NoChangeRequired

	/* Debug Code for logging into TMS as Global Admin with no LDAP access
	u1 := uuid.Must(uuid.NewV4())
	user := entities.TMSUser{Username: username, Token: u1.String(), LoggedIn: true}
	user.RoleID = entities.RoleTMSSuperUser
	err := dal.SaveUser(user)
	if err != nil {
		return nil, entities.NoChangeRequired, err
	}

	// Assign some default permissions to the user based on user role
	// if none already exist i.e. first time login
	user.UserID = dal.GetUserId(user.Username)
	dal.SetDefaultUserPermissions(user)

	return &user, entities.NoChangeRequired, nil
	*/

	// Check if the user is in the DB, if not, return login failure see: NEX-9413
	userID, err := dal.GetUserID(username)
	if err != nil {
		return nil, entities.NoChangeRequired, err
	}
	// if -1 or 0 is returned, this means the user does not exist
	if userID <= 0 {
		log.Print("Login failed: user does not exist. Status code:", http.StatusNotFound)
		return nil, http.StatusNotFound, fmt.Errorf("login failed: invalid credential")
	}

	l, err := GetLdapConn()
	if err != nil {
		return nil, entities.NoChangeRequired, err
	}
	defer l.Close()

	sr, err := lookupUserByUsername(l, username)
	if err != nil {
		return nil, entities.NoChangeRequired, err
	}

	// Check the search request for a flag indicating password change is required
	if sr.Controls != nil {
		controlIssue, err := checkLdapControls(sr.Controls)
		if err != nil {
			return nil, entities.NoChangeRequired, err
		}
		if controlIssue != entities.NoChangeRequired {
			return nil, controlIssue, nil
		}
	}

	entry := sr.Entries[0]
	userdn := entry.DN

	// Bind as the user to verify their password
	if err := l.Bind(userdn, password); err != nil {
		if isErrorRequiringPasswordChange(err) {
			passwordStatus = entities.PasswordExpired
		} else {
			return LDAPErrorHandler(err)
		}
	}

	//Generate user
	u1 := uuid.Must(uuid.NewV4())
	user := entities.TMSUser{
		Username: username,
		UserID:   userID,
		Token:    u1.String(),
		LoggedIn: true,
	}

	for _, memberOf := range entry.GetAttributeValues("memberOf") {
		groups := strings.Split(memberOf, ",")
		if !settings.AuthenticateUserPermissions {
			user.RoleID = entities.RoleTMSUser
		} else {
			for _, group := range groups {
				if group == settings.UserGroup {
					user.RoleID = entities.RoleTMSUser
				} else if group == settings.AdminGroup {
					user.RoleID = entities.RoleTMSAdministrator
				} else if group == settings.GlobalAdminGroup {
					user.RoleID = entities.RoleTMSSuperUser
				}
			}
		}
	}
	if user.RoleID == 0 {
		return nil, entities.NoChangeRequired, fmt.Errorf("%s doesn't have access to this system", username)
	}

	if err := dal.SaveUser(user.Username, user.RoleID, 1); err != nil {
		return nil, entities.NoChangeRequired, err
	}

	// Assign some default permissions to the user based on user role
	// if none already exist i.e. first time login
	user.LoggedIn = true
	dal.SetDefaultUserPermissions(user)

	if pwChangeRequired, err := userRequiresPasswordChange(user); err != nil {
		return nil, entities.NoChangeRequired, err
	} else if pwChangeRequired {
		passwordStatus = entities.FirstTimeLogon
	}

	return &user, passwordStatus, nil
}

func userRequiresPasswordChange(user entities.TMSUser) (bool, error) {
	// Set the first logon to false as the user has now signed in. NI requested that upon initial sign-in users no
	// longer are required to change their passwords.
	err := dal.ToggleFirstTimeLogon(user.Username, false)
	return false, err
}

// Indicates whether the given error is an LDAP result which then requires that the user
// will immediately need to change their password, for example if their password has expired
// or if a flag has been set forcing the password to be changed.
func isErrorRequiringPasswordChange(err error) bool {
	return isPasswordMustChangeError(err) || isPasswordExpirationError(err)
}

func isPasswordExpirationError(err error) bool {
	if ldapErr, ok := err.(*ldap.Error); ok {
		if ldapErr.ResultCode == ldap.LDAPResultInvalidCredentials &&
			// Data 532 (1330 in decimal) indicates that this is a password expiry error;
			// see https://ldapwiki.com/wiki/Common%20Active%20Directory%20Bind%20Errors
			strings.Contains(ldapErr.Error(), "AcceptSecurityContext error, data 532") {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func isPasswordMustChangeError(err error) bool {
	if ldapErr, ok := err.(*ldap.Error); ok {
		if ldapErr.ResultCode == ldap.LDAPResultInvalidCredentials &&
			// Data 773 (1907 in decimal) indicates that this is a 'password must change' error;
			// see https://ldapwiki.com/wiki/Common%20Active%20Directory%20Bind%20Errors
			strings.Contains(ldapErr.Error(), "AcceptSecurityContext error, data 773") {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

// Searches for a user with a given name. Assumes that a unique result will be found;
// an error will be returned if no unique result is found. Any LDAP query errors
// will NOT be returned directly in order to conform with Veracode requirements.
func lookupUserByUsername(l ldapConnection, username string) (*ldap.SearchResult, error) {
	// First bind with a read only user
	err := l.Bind(settings.Username, settings.Password)
	if err != nil {
		return nil, fmt.Errorf("error binding default user")
	}

	searchRequest := ldap.NewSearchRequest(
		settings.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectclass=%s)(%s=%s))", settings.ObjectClass, settings.IdentityAttribute, username),
		[]string{"memberof"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil || len(sr.Entries) == 0 || len(sr.Entries) > 1 {
		return nil, fmt.Errorf("Unable to find user: %s", username)
	}

	return sr, err
}

func checkLdapControls(controls []ldap.Control) (int, error) {
	for _, control := range controls {
		controlType := control.GetControlType()
		switch controlType {
		case ldap.ControlTypeBeheraPasswordPolicy:
			controlString := control.String()
			// controlString is in the format "Control Type: %s (%q)  Criticality: %t  Expire: %d  Grace: %d  Error: %d, ErrorString: %s"
			// We only need the error code
			codeString := strings.Split(strings.Split(controlString, "Error: ")[1], ",")[0]
			code, err := strconv.Atoi(codeString)
			if err != nil {
				return entities.NoChangeRequired, err
			}
			return code, nil
			// Currently the only controlType we are interested in is the above, this can be extended if there is future need
		}
	}
	return entities.NoChangeRequired, nil
}

func UpdateUsers() error {
	l, err := GetLdapConn()
	if err != nil {
		return err
	}

	defer l.Close()

	// First bind with a read only user
	err = l.Bind(settings.Username, settings.Password)

	if err != nil {
		return fmt.Errorf("error binding default user")
	}

	searchRequest := ldap.NewSearchRequest(
		settings.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectclass=%s)(%s=%s))", settings.ObjectClass, "memberOf", settings.GlobalAdminGroup+","+settings.AdminGroup+","+settings.BaseDN),
		[]string{"*"},
		nil,
	)

	//Lookup user
	sr, err := l.Search(searchRequest)
	if err != nil || len(sr.Entries) < 1 {
		return err
	}

	for _, u := range sr.Entries {
		//Generate user
		u1 := uuid.Must(uuid.NewV4())
		user := entities.TMSUser{Username: u.GetAttributeValue(settings.IdentityAttribute), Token: u1.String(), LoggedIn: true}
		err = dal.SaveUser(user.Username, user.RoleID, 1)
		if err != nil {
			//TODO handle errorcase
		}
	}

	return nil
}
