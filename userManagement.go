package main

import (
	"net/http"
	"nextgen-tms-website/dal"
	"nextgen-tms-website/entities"
)

func BuildUserPermissionsModel(w http.ResponseWriter, r *http.Request, user *entities.TMSUser) {
	user.UserPermissions, _ = dal.GetUserPermissions(user.Username)
}
