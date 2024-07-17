package handlers

import (
	cfg "bitbucket.org/network-international/nextgen-libs/nextgen-helpers/configHelper"
	"encoding/json"
	"github.com/gorilla/context"
	"net/http"
	"nextgen-tms-website/authentication"
	"nextgen-tms-website/dal"
	"nextgen-tms-website/entities"
	"nextgen-tms-website/logger"
	"nextgen-tms-website/models"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := BasicAuth(r)
		if ok {
			user, passwordStatus, err := authentication.ValidateUser(username, password)
			if err != nil {
				if passwordStatus == entities.NoChangeRequired {
					passwordStatus = http.StatusUnauthorized
				}
				logger.GetLogger().Information("Error from ldap", err.Error())
				respondJSON(models.Error{Code: passwordStatus, Message: err.Error()}, passwordStatus, w)
				return
			} else if passwordStatus == entities.PasswordExpired {
				respondJSON(models.Error{Code: http.StatusNotAcceptable, Message: "Please change the password"}, http.StatusNotAcceptable, w)
				return
			}

			user.UserPermissions, err = dal.GetUserPermissions(user.Username)
			if err != nil {
				logger.GetLogger().Information("Error while checking permission", err.Error())
				respondJSON(models.Error{Code: http.StatusInternalServerError, Message: "Error while checking permission"}, http.StatusInternalServerError, w)
				return
			}

			if user.UserPermissions.APIAutomation {
				logger.GetLogger().Information(username + " logged in...")
				context.Set(r, "user", user)
				next.ServeHTTP(w, r)
				return
			}
			logger.GetLogger().Information(username + " does not have permission.")
		} else {
			logger.GetLogger().Information("Basic auth failed")
		}

		respondJSON(models.Error{Code: http.StatusUnauthorized, Message: "Unauthorized"}, http.StatusUnauthorized, w)
	})
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedOrigins := cfg.Get("AllowedOrigins", []string{}).([]string)
		if len(allowedOrigins) != 0 {
			w.Header().Set("Access-Control-Allow-Origin", strings.Join(allowedOrigins, ", "))
		}
		logger.GetLogger().Information(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func respondJSON(response interface{}, statusCode int, w http.ResponseWriter) {
	jsonResponse, jsonError := json.Marshal(response)
	if jsonError != nil {
		logger.GetLogger().Error("Unable to Marshal JSON : ", jsonError.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonResponse)
}
