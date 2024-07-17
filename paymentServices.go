package main

import (
	"encoding/json"
	"net/http"
	"nextgen-tms-website/dal"
	"nextgen-tms-website/entities"
	"strconv"
)

type PaymentServicesPageModel struct {
	Group *entities.PaymentServiceGroup
	Error bool
}

func paymentServicesHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	renderHeader(w, r, tmsUser)
	renderTemplate(w, r, "paymentServices", nil, tmsUser)
}

func paymentServicesManageHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	var pageModel PaymentServicesPageModel
	var err error
	groupId := r.URL.Query().Get("groupId")
	pageModel.Group, err = dal.GetPaymentServiceGroup(groupId)
	pageModel.Error = err != nil
	renderHeader(w, r, tmsUser)
	renderTemplate(w, r, "paymentServicesManagement", pageModel, tmsUser)
}

func paymentServicesSearchHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	searchTerm := r.Form.Get("search[value]")
	offset := r.Form.Get("start")
	pagesize := r.Form.Get("length")
	orderedColumn := r.Form.Get("order[0][column]")
	orderDirection := r.Form.Get("order[0][dir]")
	requestType := r.Form.Get("requestType")

	if orderedColumn != "0" {
		orderDirection = ""
	}

	draw, err := strconv.Atoi(r.Form.Get("draw"))
	if err != nil {
		logging.Warning(err.Error())
		http.Error(w, searchError, http.StatusInternalServerError)
		return
	}

	searchSize := "0"
	searchBytes := []byte("[]")
	if requestType == "group" {
		var groups []*entities.PaymentServiceGroup
		searchSize, groups = dal.SearchServiceGroups(searchTerm, orderDirection, offset, pagesize)

		if groups != nil && len(groups) > 0 {
			searchBytes, err = json.Marshal(&groups)
		}
	} else {
		groupId := r.Form.Get("groupId")
		if len(groupId) == 0 {
			_, _ = logging.Warning(err.Error())
			http.Error(w, searchError, http.StatusInternalServerError)
			return
		}
		var services []*entities.PaymentService
		searchSize, services = dal.SearchServicesInGroup(groupId, searchTerm, orderDirection, offset, pagesize)

		if services != nil && len(services) > 0 {
			searchBytes, err = json.Marshal(&services)
		}
	}

	if err != nil {
		_, _ = logging.Warning(err.Error())
		http.Error(w, searchError, http.StatusInternalServerError)
		return
	}

	bytesToWrite := append([]byte("{\"draw\":"+strconv.Itoa(draw)+","+
		"\"recordsTotal\":"+searchSize+","+
		"\"recordsFiltered\":"+searchSize+","+
		"\"data\":"), append(searchBytes, []byte("}")...)...)

	if _, err = w.Write(bytesToWrite); err != nil {
		_, _ = logging.Error(err.Error())
		http.Error(w, searchError, http.StatusInternalServerError)
		return
	}
}

func paymentServicesDeleteGroup(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	groupId := r.Form.Get("groupId")

	groupDetais, err := dal.GetPaymentServiceGroup(groupId)
	if err != nil {
		_, _ = logging.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = dal.PaymentServiceGroupDeletionChangeApproval(tmsUser.Username, 1, 10, groupId, groupDetais.Name, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("success")); err != nil {
		_, _ = logging.Error(err.Error())
		http.Error(w, searchError, http.StatusInternalServerError)
		return
	}
}

func paymentServicesDelete(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	serviceId := r.Form.Get("serviceId")

	err := dal.PaymentServicesDeletionChangeApproval(tmsUser.Username, 1, 11, serviceId, "", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("success")); err != nil {
		_, _ = logging.Error(err.Error())
		http.Error(w, searchError, http.StatusInternalServerError)
		return
	}
}

func paymentServicesAddGroup(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	name := r.Form.Get("name")
	err := dal.AddPaymentServiceGroup(name)
	if err != nil {
		_, _ = logging.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = dal.PaymentServiceCreationChangeApproval(tmsUser.Username, 1, 8, name, "", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write([]byte("success")); err != nil {
		_, _ = logging.Error(err.Error())
		http.Error(w, "an internal error occurred", http.StatusInternalServerError)
		return
	}
}

func paymentServicesAddService(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	groupId := r.Form.Get("groupId")
	name := r.Form.Get("name")

	err := dal.AddPaymentService(name, groupId)
	if err != nil {
		_, _ = logging.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	groupDetais, err := dal.GetPaymentServiceGroup(groupId)
	if err != nil {
		_, _ = logging.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = dal.PaymentServiceCreationChangeApproval(tmsUser.Username, 1, 9, name, groupDetais.Name, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write([]byte("success")); err != nil {
		_, _ = logging.Error(err.Error())
		http.Error(w, "an internal error occurred", http.StatusInternalServerError)
		return
	}
}
