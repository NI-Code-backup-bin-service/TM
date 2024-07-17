package main

import (
	"net/http"
	"nextgen-tms-website/authentication"
	"nextgen-tms-website/dal"
	"nextgen-tms-website/entities"
)

type ChangePasswordModel struct {
	Username     string
	ChangeReason string
}

// Handler func for changing a TMS user password
func ChangeUserPassword(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	r.ParseForm()

	username := r.Form.Get("username")
	newPassword := r.Form.Get("newPassword")
	oldPassword := r.Form.Get("oldPassword")

	userID, err := dal.GetUserID(username)
	if err != nil {
		http.Error(w, InvalidUsernameError, http.StatusInternalServerError)
		return
	}

	switch userID {
	case -1:
		http.Error(w, PasswordChangeError, http.StatusInternalServerError)
		return
	case 0:
		http.Error(w, InvalidUsernameError, http.StatusInternalServerError)
		return
	}

	if err := processPasswordChange(username, oldPassword, newPassword); err != nil {
		logging.Error(err.Error())
		if _, ok := err.(*authentication.PasswordConstraintError); ok {
			http.Error(w, PasswordConstraintError, http.StatusBadRequest)
		} else {
			http.Error(w, PasswordChangeError, http.StatusInternalServerError)
		}
		return
	}

	renderTemplate(w, r, "passwordChange", nil, &entities.TMSUser{LoggedIn: false})
}

// Interfaces with LDAP server to process the password change
func processPasswordChange(username string, oldPassword, newPassword string) error {

	err := authentication.ModifyPassword(username, oldPassword, newPassword)
	if err != nil {
		return err
	}

	// Update the db flag to show the user has updated their password on first logon
	return dal.ToggleFirstTimeLogon(username, false)
}

// Compares the reason a password change is required and returns an explanatory string for the site user
func fetchPasswordChangeResponse(changeCode int) string {
	switch changeCode {
	case entities.PasswordExpired:
		return PasswordExpiryMessage
	case entities.FirstTimeLogon:
		return PasswordFirstTimeLogon
	case entities.InsufficientPasswordQuality:
		return PasswordInsufficientQuality
	case entities.PasswordTooShort:
		return PasswordTooShort
	case entities.PasswordTooYoung:
		return PasswordTooYoung
	case entities.PasswordInHistory:
		return PasswordInHistory
	}
	return ""
}
