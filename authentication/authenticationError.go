package authentication

import (
	"fmt"
	"log"
	"net/http"
	"nextgen-tms-website/entities"
	"strings"
)

func LDAPErrorHandler(err error) (*entities.TMSUser, int, error) {
	if strings.Contains(err.Error(), "data 52e") {
		log.Print("Login failed: invalid credential. Status code:", http.StatusUnauthorized)
		return nil, http.StatusUnauthorized, fmt.Errorf("login failed: invalid credential")
	}
	if strings.Contains(err.Error(), "data 775") {
		log.Print("Login failed: account locked out. Status code:", http.StatusTooManyRequests)
		return nil, http.StatusTooManyRequests, fmt.Errorf("login failed: invalid credential")
	}
	return nil, entities.NoChangeRequired, fmt.Errorf("login failed: invalid credential")
}
