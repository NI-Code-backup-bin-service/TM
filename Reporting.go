package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"nextgen-tms-website/dal"
	"nextgen-tms-website/entities"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type ReportingPageModel struct {
	Acquirers []string
}

func reportingHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	var pageModel ReportingPageModel

	acqs, err := dal.GetUserAcquirerPermissions(tmsUser)
	if err != nil {
		handleError(w, errors.New("an error occured while retriving the user permission:"+err.Error()), tmsUser)
		return
	}
	acqs = strings.Trim(acqs, "'")
	acqs = strings.Trim(acqs, ",")
	acquirers := strings.Split(acqs, ",")
	pageModel.Acquirers = acquirers

	renderHeader(w, r, tmsUser)
	renderTemplate(w, r, "Reporting", pageModel, tmsUser)
}

func getTimeRange(r *http.Request) bson.M {
	r.ParseForm()
	before := r.Form.Get("Before")
	after := r.Form.Get("After")
	tpe := r.Form.Get("Type")
	_, _, _ = before, after, tpe

	timeFormat := "2006/01/02 15:04"

	parameters := bson.M{}

	b := time.Now().Format(timeFormat)
	var a string

	switch tpe {
	case "month":
		a = time.Now().AddDate(0, -1, 0).Format(timeFormat)
		beforeDate, _ := time.Parse(timeFormat, b)
		afterDate, _ := time.Parse(timeFormat, a)
		parameters["startTime"] = bson.M{"startTime": bson.M{"$gte": afterDate, "$lte": beforeDate}}

	case "week":
		a = time.Now().AddDate(0, 0, -7).Format(timeFormat)
		beforeDate, _ := time.Parse(timeFormat, b)
		afterDate, _ := time.Parse(timeFormat, a)
		parameters["startTime"] = bson.M{"startTime": bson.M{"$gte": afterDate, "$lte": beforeDate}}

	case "day":
		year, month, day := time.Now().Date()

		a = time.Date(year, month, day, 0, 0, 0, 0, time.Now().Location()).Format(timeFormat)
		beforeDate, _ := time.Parse(timeFormat, b)
		afterDate, _ := time.Parse(timeFormat, a)
		parameters["startTime"] = bson.M{"startTime": bson.M{"$gte": afterDate, "$lte": beforeDate}}

	case "custom":
		if before != "" && after != "" {
			beforeDate, _ := time.Parse(timeFormat, before)
			afterDate, _ := time.Parse(timeFormat, after)
			parameters["startTime"] = bson.M{"startTime": bson.M{"$gte": afterDate, "$lte": beforeDate}}
		} else {
			if before != "" {
				beforeDate, _ := time.Parse(timeFormat, before)
				parameters["startTime"] = bson.M{"startTime": bson.M{"$lte": beforeDate}}
			}
			if after != "" {
				afterDate, _ := time.Parse(timeFormat, after)
				parameters["startTime"] = bson.M{"startTime": bson.M{"$gte": afterDate}}
			}
			if before == "" && after == "" {
				parameters["startTime"] = bson.M{}
			}
		}
	}
	return parameters
}

func getAcquirerPermissions(r *http.Request) []string {
	r.ParseForm()
	keyPairs := r.Form
	var acquirers []string
	for key := range keyPairs {
		if key == "Acquirers[]" {
			var multiKeys = keyPairs[key]
			if len(multiKeys) > 0 {
				for i := range multiKeys {
					acquirers = append(acquirers, multiKeys[i])
				}
			}
		}
	}
	return acquirers
}

func TotalMIDs(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	totalMIDs, err := dal.GetTotalMIDs()
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, reportingError, http.StatusBadRequest)
	}

	results, err := json.Marshal(totalMIDs)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, reportingError, http.StatusBadRequest)
	}

	w.Write(results)
}

func TotalTIDs(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	totalTIDs, err := dal.GetTotalTIDs()
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, reportingError, http.StatusBadRequest)
	}

	results, err := json.Marshal(totalTIDs)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, reportingError, http.StatusBadRequest)
	}
	w.Write(results)
}

func TransactingMIDs(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	parameters := getTimeRange(r)
	acquirers := getAcquirerPermissions(r)
	find := []bson.M{}
	if acquirers != nil {
		find = []bson.M{
			{"$match": bson.M{"$and": []interface{}{
				bson.M{"request.amount": bson.M{"$gte": 1}},
				parameters["startTime"],
				bson.M{"siteid": bson.M{"$ne": nil}},
				bson.M{"siteid": bson.M{"$ne": ""}},
				bson.M{"acquirer": bson.M{"$in": acquirers}},
			}},
			},
			{"$group": bson.M{"_id": bson.M{"id": "$siteid"}, "count": bson.M{"$sum": 1}}},
			{"$count": "count"},
		}
	} else {
		find = []bson.M{
			{"$match": bson.M{"$and": []interface{}{
				bson.M{"request.amount": bson.M{"$gte": 1}},
				parameters["startTime"],
				bson.M{"siteid": bson.M{"$ne": nil}},
				bson.M{"siteid": bson.M{"$ne": ""}},
			}},
			},
			{"$group": bson.M{"_id": bson.M{"id": "$siteid"}, "count": bson.M{"$sum": 1}}},
			{"$count": "count"},
		}
	}
	TransactingMids, err := dal.GetTotals(r.Context(), find)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, reportingError, http.StatusBadRequest)
	}

	results, err := json.Marshal(TransactingMids)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, reportingError, http.StatusBadRequest)
	}
	w.Write(results)
}

func TransactingTIDs(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	parameters := getTimeRange(r)
	acquirers := getAcquirerPermissions(r)
	find := []bson.M{}
	if acquirers != nil {
		find = []bson.M{
			{"$match": bson.M{"$and": []interface{}{
				bson.M{"request.amount": bson.M{"$gte": 1}},
				parameters["startTime"],
				bson.M{"acquirer": bson.M{"$in": acquirers}},
			}}},
			{"$group": bson.M{"_id": bson.M{"tid": "$tid"}, "count": bson.M{"$sum": 1}}},
			{"$count": "count"},
		}
	} else {
		find = []bson.M{
			{"$match": bson.M{"$and": []interface{}{
				bson.M{"request.amount": bson.M{"$gte": 1}},
				parameters["startTime"],
			}}},
			{"$group": bson.M{"_id": bson.M{"tid": "$tid"}, "count": bson.M{"$sum": 1}}},
			{"$count": "count"},
		}
	}

	TransactingTids, err := dal.GetTotals(r.Context(), find)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, reportingError, http.StatusBadRequest)
	}

	results, err := json.Marshal(TransactingTids)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, reportingError, http.StatusBadRequest)
	}
	w.Write(results)
}

func TotalTransactions(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	parameters := getTimeRange(r)
	acquirers := getAcquirerPermissions(r)
	find := []bson.M{}
	if acquirers != nil {
		find = []bson.M{
			{"$match": bson.M{"$and": []interface{}{
				bson.M{"request.amount": bson.M{"$gte": 1}},
				parameters["startTime"],
				bson.M{"siteid": bson.M{"$ne": nil}},
				bson.M{"siteid": bson.M{"$ne": ""}},
				bson.M{"acquirer": bson.M{"$in": acquirers}},
			}}},
			{"$group": bson.M{"_id": "", "count": bson.M{"$sum": 1}}},
		}
	} else {
		find = []bson.M{
			{"$match": bson.M{"$and": []interface{}{
				bson.M{"request.amount": bson.M{"$gte": 1}},
				parameters["startTime"],
				bson.M{"siteid": bson.M{"$ne": nil}},
				bson.M{"siteid": bson.M{"$ne": ""}},
			}}},
			{"$group": bson.M{"_id": "", "count": bson.M{"$sum": 1}}},
		}
	}

	TotalTransactions, err := dal.GetTotals(r.Context(), find)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, reportingError, http.StatusBadRequest)
	}

	results, err := json.Marshal(TotalTransactions)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, reportingError, http.StatusBadRequest)
	}
	w.Write(results)
}

func TotalTransactionValue(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	parameters := getTimeRange(r)
	acquirers := getAcquirerPermissions(r)
	find := []bson.M{}
	if acquirers != nil {
		find = []bson.M{
			{"$match": bson.M{"$and": []interface{}{
				bson.M{"request.amount": bson.M{"$gte": 1}},
				parameters["startTime"],
				bson.M{"siteid": bson.M{"$ne": nil}},
				bson.M{"siteid": bson.M{"$ne": ""}},
				bson.M{"acquirer": bson.M{"$in": acquirers}},
			}}},
			{"$group": bson.M{"_id": "", "count": bson.M{"$sum": "$request.amount"}}},
		}
	} else {
		find = []bson.M{
			{"$match": bson.M{"$and": []interface{}{
				bson.M{"request.amount": bson.M{"$gte": 1}},
				parameters["startTime"],
				bson.M{"siteid": bson.M{"$ne": nil}},
				bson.M{"siteid": bson.M{"$ne": ""}},
			}}},
			{"$group": bson.M{"_id": "", "count": bson.M{"$sum": "$request.amount"}}},
		}
	}

	TotalTransactionValue, err := dal.GetMonetaryTotals(r.Context(), find)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, reportingError, http.StatusBadRequest)
	}

	results, err := json.Marshal(TotalTransactionValue)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, reportingError, http.StatusBadRequest)
	}
	w.Write(results)
}

func ApprovedTransactions(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	parameters := getTimeRange(r)
	acquirers := getAcquirerPermissions(r)
	find := []bson.M{}
	if acquirers != nil {
		find = []bson.M{
			{"$match": bson.M{"$and": []interface{}{
				bson.M{"request.amount": bson.M{"$gte": 1}},
				parameters["startTime"],
				bson.M{"emvResponse.resultcode": bson.M{"$regex": "00"}},
				bson.M{"acquirer": bson.M{"$in": acquirers}},
			}}},
			{"$group": bson.M{"_id": "", "count": bson.M{"$sum": 1}}},
		}
	} else {
		find = []bson.M{
			{"$match": bson.M{"$and": []interface{}{
				bson.M{"request.amount": bson.M{"$gte": 1}},
				parameters["startTime"],
				bson.M{"emvResponse.resultcode": bson.M{"$regex": "00"}},
			}}},
			{"$group": bson.M{"_id": "", "count": bson.M{"$sum": 1}}},
		}
	}
	ApprovedTransactions, err := dal.GetTotals(r.Context(), find)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, reportingError, http.StatusBadRequest)
	}

	results, err := json.Marshal(ApprovedTransactions)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, reportingError, http.StatusBadRequest)
	}
	w.Write(results)
}

func DeclinedTransactions(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	parameters := getTimeRange(r)
	acquirers := getAcquirerPermissions(r)
	find := []bson.M{}
	if acquirers != nil {
		find = []bson.M{
			{"$match": bson.M{"$and": []interface{}{
				bson.M{"request.amount": bson.M{"$gte": 1}},
				parameters["startTime"],
				bson.M{"emvResponse": bson.M{"$ne": nil}},
				bson.M{"emvResponse.resultcode": bson.M{"$regex": "^((?!00).)*$"}},
				bson.M{"acquirer": bson.M{"$in": acquirers}},
			}},
			},
			{"$group": bson.M{"_id": "", "count": bson.M{"$sum": 1}}}}
	} else {
		find = []bson.M{
			{"$match": bson.M{"$and": []interface{}{
				bson.M{"request.amount": bson.M{"$gte": 1}},
				parameters["startTime"],
				bson.M{"emvResponse": bson.M{"$ne": nil}},
				bson.M{"emvResponse.resultcode": bson.M{"$regex": "^((?!00).)*$"}},
			}},
			},
			{"$group": bson.M{"_id": "", "count": bson.M{"$sum": 1}}}}
	}
	DeclinedTransactions, err := dal.GetTotals(r.Context(), find)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, reportingError, http.StatusBadRequest)
	}

	results, err := json.Marshal(DeclinedTransactions)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, reportingError, http.StatusBadRequest)
	}
	w.Write(results)
}

func TransactionVolume(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	parameters := getTimeRange(r)
	acquirers := getAcquirerPermissions(r)
	find := []bson.M{}
	if acquirers != nil {
		find = []bson.M{
			{"$match": bson.M{"$and": []interface{}{
				bson.M{"request.amount": bson.M{"$gte": 1}},
				parameters["startTime"],
				bson.M{"acquirer": bson.M{"$in": acquirers}},
			}},
			},
			{"$group": bson.M{"_id": bson.M{"date": bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d", "date": "$startTime"}}}, "count": bson.M{"$sum": 1}}},
			{"$sort": bson.M{"_id.date": 1}}}
	} else {
		find = []bson.M{
			{"$match": bson.M{"$and": []interface{}{
				bson.M{"request.amount": bson.M{"$gte": 1}},
				parameters["startTime"],
			}},
			},
			{"$group": bson.M{"_id": bson.M{"date": bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d", "date": "$startTime"}}}, "count": bson.M{"$sum": 1}}},
			{"$sort": bson.M{"_id.date": 1}}}
	}
	totalVolume, err := dal.GetTransactionVolume(r.Context(), find)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, reportingError, http.StatusBadRequest)
	}

	results, err := json.Marshal(totalVolume)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, reportingError, http.StatusBadRequest)
	}
	w.Write(results)
}

func ApprovedTransactionVolume(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	parameters := getTimeRange(r)
	acquirers := getAcquirerPermissions(r)
	find := []bson.M{}
	if acquirers != nil {
		find = []bson.M{
			{"$match": bson.M{"$and": []interface{}{
				bson.M{"request.amount": bson.M{"$gte": 1}},
				parameters["startTime"],
				bson.M{"emvResponse.resultcode": bson.M{"$regex": "00"}},
				bson.M{"acquirer": bson.M{"$in": acquirers}},
			}},
			},
			{"$group": bson.M{"_id": bson.M{"date": bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d", "date": "$startTime"}}}, "count": bson.M{"$sum": 1}}},
			{"$sort": bson.M{"_id.date": 1}}}
	} else {
		find = []bson.M{
			{"$match": bson.M{"$and": []interface{}{
				bson.M{"request.amount": bson.M{"$gte": 1}},
				parameters["startTime"],
				bson.M{"emvResponse.resultcode": bson.M{"$regex": "00"}},
			}},
			},
			{"$group": bson.M{"_id": bson.M{"date": bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d", "date": "$startTime"}}}, "count": bson.M{"$sum": 1}}},
			{"$sort": bson.M{"_id.date": 1}}}
	}
	totalVolume, err := dal.GetTransactionVolume(r.Context(), find)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, reportingError, http.StatusBadRequest)
	}

	results, err := json.Marshal(totalVolume)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, reportingError, http.StatusBadRequest)
	}
	w.Write(results)
}

func DeclinedTransactionVolume(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	parameters := getTimeRange(r)
	acquirers := getAcquirerPermissions(r)
	find := []bson.M{}
	if acquirers != nil {
		find = []bson.M{
			{"$match": bson.M{"$and": []interface{}{
				bson.M{"request.amount": bson.M{"$gte": 1}},
				parameters["startTime"],
				bson.M{"emvResponse.resultcode": bson.M{"$regex": "^((?!00).)*$"}},
				bson.M{"acquirer": bson.M{"$in": acquirers}},
			}},
			},
			{"$group": bson.M{"_id": bson.M{"date": bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d", "date": "$startTime"}}}, "count": bson.M{"$sum": 1}}},
			{"$sort": bson.M{"_id.date": 1}}}
	} else {
		find = []bson.M{
			{"$match": bson.M{"$and": []interface{}{
				bson.M{"request.amount": bson.M{"$gte": 1}},
				parameters["startTime"],
				bson.M{"emvResponse.resultcode": bson.M{"$regex": "^((?!00).)*$"}},
			}},
			},
			{"$group": bson.M{"_id": bson.M{"date": bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d", "date": "$startTime"}}}, "count": bson.M{"$sum": 1}}},
			{"$sort": bson.M{"_id.date": 1}}}
	}
	totalVolume, err := dal.GetTransactionVolume(r.Context(), find)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, reportingError, http.StatusBadRequest)
	}

	results, err := json.Marshal(totalVolume)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, reportingError, http.StatusBadRequest)
	}
	w.Write(results)
}

func TransactionValue(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	parameters := getTimeRange(r)
	acquirers := getAcquirerPermissions(r)
	find := []bson.M{}
	if acquirers != nil {
		find = []bson.M{
			{"$match": bson.M{"$and": []interface{}{
				bson.M{"request.amount": bson.M{"$gte": 1}},
				parameters["startTime"],
				bson.M{"acquirer": bson.M{"$in": acquirers}},
			}},
			},
			{"$group": bson.M{"_id": bson.M{"date": bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d", "date": "$startTime"}}}, "count": bson.M{"$sum": "$request.amount"}}},
			{"$sort": bson.M{"_id.date": 1}}}
	} else {
		find = []bson.M{
			{"$match": bson.M{"$and": []interface{}{
				bson.M{"request.amount": bson.M{"$gte": 1}},
				parameters["startTime"],
			}},
			},
			{"$group": bson.M{"_id": bson.M{"date": bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d", "date": "$startTime"}}}, "count": bson.M{"$sum": "$request.amount"}}},
			{"$sort": bson.M{"_id.date": 1}}}
	}
	totalVolume, err := dal.GetMonetaryValues(r.Context(), find)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, reportingError, http.StatusBadRequest)
	}

	results, err := json.Marshal(totalVolume)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, reportingError, http.StatusBadRequest)
	}
	w.Write(results)
}

func ApprovedTransactionValue(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	parameters := getTimeRange(r)
	acquirers := getAcquirerPermissions(r)
	find := []bson.M{}
	if acquirers != nil {
		find = []bson.M{
			{"$match": bson.M{"$and": []interface{}{
				bson.M{"request.amount": bson.M{"$gte": 1}},
				parameters["startTime"],
				bson.M{"emvResponse.resultcode": bson.M{"$regex": "00"}},
				bson.M{"acquirer": bson.M{"$in": acquirers}},
			}},
			},
			{"$group": bson.M{"_id": bson.M{"date": bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d", "date": "$startTime"}}}, "count": bson.M{"$sum": "$request.amount"}}},
			{"$sort": bson.M{"_id.date": 1}}}
	} else {
		find = []bson.M{
			{"$match": bson.M{"$and": []interface{}{
				bson.M{"request.amount": bson.M{"$gte": 1}},
				parameters["startTime"],
				bson.M{"emvResponse.resultcode": bson.M{"$regex": "00"}},
			}},
			},
			{"$group": bson.M{"_id": bson.M{"date": bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d", "date": "$startTime"}}}, "count": bson.M{"$sum": "$request.amount"}}},
			{"$sort": bson.M{"_id.date": 1}}}
	}
	totalVolume, err := dal.GetMonetaryValues(r.Context(), find)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, reportingError, http.StatusBadRequest)
	}

	results, err := json.Marshal(totalVolume)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, reportingError, http.StatusBadRequest)
	}
	w.Write(results)
}

func DeclinedTransactionValue(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	parameters := getTimeRange(r)
	acquirers := getAcquirerPermissions(r)
	find := []bson.M{}
	if acquirers != nil {
		find = []bson.M{
			{"$match": bson.M{"$and": []interface{}{
				bson.M{"request.amount": bson.M{"$gte": 1}},
				parameters["startTime"],
				bson.M{"emvResponse.resultcode": bson.M{"$regex": "^((?!00).)*$"}},
				bson.M{"acquirer": bson.M{"$in": acquirers}},
			}},
			},
			{"$group": bson.M{"_id": bson.M{"date": bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d", "date": "$startTime"}}}, "count": bson.M{"$sum": "$request.amount"}}},
			{"$sort": bson.M{"_id.date": 1}}}
	} else {
		find = []bson.M{
			{"$match": bson.M{"$and": []interface{}{
				bson.M{"request.amount": bson.M{"$gte": 1}},
				parameters["startTime"],
				bson.M{"emvResponse.resultcode": bson.M{"$regex": "^((?!00).)*$"}},
			}},
			},
			{"$group": bson.M{"_id": bson.M{"date": bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d", "date": "$startTime"}}}, "count": bson.M{"$sum": "$request.amount"}}},
			{"$sort": bson.M{"_id.date": 1}}}
	}
	totalVolume, err := dal.GetMonetaryValues(r.Context(), find)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, reportingError, http.StatusBadRequest)
	}

	results, err := json.Marshal(totalVolume)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, reportingError, http.StatusBadRequest)
	}
	w.Write(results)
}

func Top10DeclineReasons(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	parameters := getTimeRange(r)
	acquirers := getAcquirerPermissions(r)
	// Top 10 Decline reasons
	find := []bson.M{}
	if acquirers != nil {
		find = []bson.M{
			{"$match": bson.M{"$and": []interface{}{
				bson.M{"request.amount": bson.M{"$gte": 1}},
				parameters["startTime"],
				bson.M{"emvResponse": bson.M{"$ne": nil}},
				bson.M{"emvResponse.resultcode": bson.M{"$regex": "^((?!00).)*$"}},
				bson.M{"acquirer": bson.M{"$in": acquirers}},
			}},
			},
			{"$group": bson.M{"_id": bson.M{"id": "$emvResponse.resultcode"}, "count": bson.M{"$sum": 1}}},
			{"$sort": bson.M{"count": -1}},
			{"$limit": 10},
		}
	} else {
		find = []bson.M{
			{"$match": bson.M{"$and": []interface{}{
				bson.M{"request.amount": bson.M{"$gte": 1}},
				parameters["startTime"],
				bson.M{"emvResponse": bson.M{"$ne": nil}},
				bson.M{"emvResponse.resultcode": bson.M{"$regex": "^((?!00).)*$"}},
			}},
			},
			{"$group": bson.M{"_id": bson.M{"id": "$emvResponse.resultcode"}, "count": bson.M{"$sum": 1}}},
			{"$sort": bson.M{"count": -1}},
			{"$limit": 10},
		}
	}

	TopTransacting, err := dal.GetTopTransacting(r.Context(), find)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, reportingError, http.StatusBadRequest)
	}

	// Convert Result codes to string descriptions
	for i, obj := range TopTransacting {
		TopTransacting[i].Label = GetResultCodeDescription(obj.Label)
	}

	results, err := json.Marshal(TopTransacting)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, reportingError, http.StatusBadRequest)
	}
	w.Write(results)
}

func Top10TransactingMIDs(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	parameters := getTimeRange(r)
	acquirers := getAcquirerPermissions(r)
	find := []bson.M{}
	if acquirers != nil {
		find = []bson.M{
			{"$match": bson.M{"$and": []interface{}{
				bson.M{"request.amount": bson.M{"$gte": 1}},
				parameters["startTime"],
				bson.M{"siteid": bson.M{"$ne": nil}},
				bson.M{"siteid": bson.M{"$ne": ""}},
				bson.M{"acquirer": bson.M{"$in": acquirers}},
			}},
			},
			{"$group": bson.M{"_id": bson.M{"id": "$siteid"}, "count": bson.M{"$sum": 1}}},
			{"$sort": bson.M{"count": -1}},
			{"$limit": 10},
		}
	} else {
		find = []bson.M{
			{"$match": bson.M{"$and": []interface{}{
				bson.M{"request.amount": bson.M{"$gte": 1}},
				parameters["startTime"],
				bson.M{"siteid": bson.M{"$ne": nil}},
				bson.M{"siteid": bson.M{"$ne": ""}},
			}},
			},
			{"$group": bson.M{"_id": bson.M{"id": "$siteid"}, "count": bson.M{"$sum": 1}}},
			{"$sort": bson.M{"count": -1}},
			{"$limit": 10},
		}
	}

	TopTransacting, err := dal.GetTopTransacting(r.Context(), find)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, reportingError, http.StatusBadRequest)
	}

	// Convert Result codes to string descriptions
	for i, obj := range TopTransacting {
		var siteName = GetMerchantName(obj.Label)
		if siteName != "" {
			TopTransacting[i].Label = siteName
		}
	}

	results, err := json.Marshal(TopTransacting)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, reportingError, http.StatusBadRequest)
	}
	w.Write(results)
}

func Top10TransactingTIDs(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	parameters := getTimeRange(r)
	acquirers := getAcquirerPermissions(r)
	find := []bson.M{}
	if acquirers != nil {
		find = []bson.M{
			{"$match": bson.M{"$and": []interface{}{
				bson.M{"request.amount": bson.M{"$gte": 1}},
				parameters["startTime"],
				bson.M{"acquirer": bson.M{"$in": acquirers}},
			}},
			},
			{"$group": bson.M{"_id": bson.M{"id": "$tid"}, "count": bson.M{"$sum": 1}}},
			{"$sort": bson.M{"count": -1}},
			{"$limit": 10},
		}
	} else {
		find = []bson.M{
			{"$match": bson.M{"$and": []interface{}{
				bson.M{"request.amount": bson.M{"$gte": 1}},
				parameters["startTime"],
			}},
			},
			{"$group": bson.M{"_id": bson.M{"id": "$tid"}, "count": bson.M{"$sum": 1}}},
			{"$sort": bson.M{"count": -1}},
			{"$limit": 10},
		}
	}

	TopTransacting, err := dal.GetTopTransacting(r.Context(), find)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, reportingError, http.StatusBadRequest)
	}

	results, err := json.Marshal(TopTransacting)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, reportingError, http.StatusBadRequest)
	}
	w.Write(results)
}

func Top10CardTypes(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	parameters := getTimeRange(r)
	acquirers := getAcquirerPermissions(r)
	find := []bson.M{}
	if acquirers != nil {
		find = []bson.M{
			{"$match": bson.M{"$and": []interface{}{
				bson.M{"request.amount": bson.M{"$gte": 1}},
				parameters["startTime"],
				bson.M{"emvRequest": bson.M{"$ne": nil}},
				bson.M{"emvRequest.cardScheme": bson.M{"$ne": nil}},
				bson.M{"emvRequest.cardScheme": bson.M{"$ne": ""}},
				bson.M{"acquirer": bson.M{"$in": acquirers}},
			}},
			},
			{"$group": bson.M{"_id": bson.M{"id": "$emvRequest.cardScheme"}, "count": bson.M{"$sum": 1}}},
			{"$sort": bson.M{"count": -1}},
			{"$limit": 10},
		}
	} else {
		find = []bson.M{
			{"$match": bson.M{"$and": []interface{}{
				bson.M{"request.amount": bson.M{"$gte": 1}},
				parameters["startTime"],
				bson.M{"emvRequest": bson.M{"$ne": nil}},
				bson.M{"emvRequest.cardScheme": bson.M{"$ne": nil}},
				bson.M{"emvRequest.cardScheme": bson.M{"$ne": ""}},
			}},
			},
			{"$group": bson.M{"_id": bson.M{"id": "$emvRequest.cardScheme"}, "count": bson.M{"$sum": 1}}},
			{"$sort": bson.M{"count": -1}},
			{"$limit": 10},
		}
	}

	TopTransacting, err := dal.GetTopTransacting(r.Context(), find)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, reportingError, http.StatusBadRequest)
	}

	results, err := json.Marshal(TopTransacting)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, reportingError, http.StatusBadRequest)
	}
	w.Write(results)
}

func Top10TransactingMIDValues(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	parameters := getTimeRange(r)
	acquirers := getAcquirerPermissions(r)
	find := []bson.M{}
	if acquirers != nil {
		find = []bson.M{
			{"$match": bson.M{"$and": []interface{}{
				bson.M{"request.amount": bson.M{"$gte": 1}},
				parameters["startTime"],
				bson.M{"siteid": bson.M{"$ne": nil}},
				bson.M{"siteid": bson.M{"$ne": ""}},
				bson.M{"acquirer": bson.M{"$in": acquirers}},
			}},
			},
			{"$group": bson.M{"_id": bson.M{"id": "$siteid"}, "count": bson.M{"$sum": "$request.amount"}}},
			{"$sort": bson.M{"count": -1}},
			{"$limit": 10},
		}
	} else {
		find = []bson.M{
			{"$match": bson.M{"$and": []interface{}{
				bson.M{"request.amount": bson.M{"$gte": 1}},
				parameters["startTime"],
				bson.M{"siteid": bson.M{"$ne": nil}},
				bson.M{"siteid": bson.M{"$ne": ""}},
			}},
			},
			{"$group": bson.M{"_id": bson.M{"id": "$siteid"}, "count": bson.M{"$sum": "$request.amount"}}},
			{"$sort": bson.M{"count": -1}},
			{"$limit": 10},
		}
	}

	TopTransacting, err := dal.GetMonetaryTopTransacting(r.Context(), find)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, reportingError, http.StatusBadRequest)
	}

	// Convert Result codes to string descriptions
	for i, obj := range TopTransacting {
		var siteName = GetMerchantName(obj.Label)
		if siteName != "" {
			TopTransacting[i].Label = siteName
		}
	}

	results, err := json.Marshal(TopTransacting)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, reportingError, http.StatusBadRequest)
	}
	w.Write(results)
}

func Top10TransactingTIDValues(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	parameters := getTimeRange(r)
	acquirers := getAcquirerPermissions(r)
	find := []bson.M{}
	if acquirers != nil {
		find = []bson.M{
			{"$match": bson.M{"$and": []interface{}{
				bson.M{"request.amount": bson.M{"$gte": 1}},
				parameters["startTime"],
				bson.M{"acquirer": bson.M{"$in": acquirers}},
			}},
			},
			{"$group": bson.M{"_id": bson.M{"id": "$tid"}, "count": bson.M{"$sum": "$request.amount"}}},
			{"$sort": bson.M{"count": -1}},
			{"$limit": 10},
		}
	} else {
		find = []bson.M{
			{"$match": bson.M{"$and": []interface{}{
				bson.M{"request.amount": bson.M{"$gte": 1}},
				parameters["startTime"],
			}},
			},
			{"$group": bson.M{"_id": bson.M{"id": "$tid"}, "count": bson.M{"$sum": "$request.amount"}}},
			{"$sort": bson.M{"count": -1}},
			{"$limit": 10},
		}
	}

	TopTransacting, err := dal.GetMonetaryTopTransacting(r.Context(), find)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, reportingError, http.StatusBadRequest)
	}

	results, err := json.Marshal(TopTransacting)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, reportingError, http.StatusBadRequest)
	}
	w.Write(results)
}

func Top10CardTypeValues(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	parameters := getTimeRange(r)
	acquirers := getAcquirerPermissions(r)
	find := []bson.M{}
	if acquirers != nil {
		find = []bson.M{
			{"$match": bson.M{"$and": []interface{}{
				bson.M{"request.amount": bson.M{"$gte": 1}},
				parameters["startTime"],
				bson.M{"emvRequest": bson.M{"$ne": nil}},
				bson.M{"emvRequest.cardScheme": bson.M{"$ne": nil}},
				bson.M{"emvRequest.cardScheme": bson.M{"$ne": ""}},
				bson.M{"acquirer": bson.M{"$in": acquirers}},
			}},
			},
			{"$group": bson.M{"_id": bson.M{"id": "$emvRequest.cardScheme"}, "count": bson.M{"$sum": "$request.amount"}}},
			{"$sort": bson.M{"count": -1}},
			{"$limit": 10},
		}
	} else {
		find = []bson.M{
			{"$match": bson.M{"$and": []interface{}{
				bson.M{"request.amount": bson.M{"$gte": 1}},
				parameters["startTime"],
				bson.M{"emvRequest": bson.M{"$ne": nil}},
				bson.M{"emvRequest.cardScheme": bson.M{"$ne": nil}},
				bson.M{"emvRequest.cardScheme": bson.M{"$ne": ""}},
			}},
			},
			{"$group": bson.M{"_id": bson.M{"id": "$emvRequest.cardScheme"}, "count": bson.M{"$sum": "$request.amount"}}},
			{"$sort": bson.M{"count": -1}},
			{"$limit": 10},
		}
	}

	TopTransacting, err := dal.GetMonetaryTopTransacting(r.Context(), find)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, reportingError, http.StatusBadRequest)
	}

	results, err := json.Marshal(TopTransacting)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, reportingError, http.StatusBadRequest)
	}
	w.Write(results)
}

func GetResultCodeDescription(resultCode string) string {
	switch resultCode {
	case "00": // Approved
		return "00 - Approved"
	case "01": // Please Call issuer
		return "01 - Call issuer"
	case "02": // Over floor limit
		return "02 - Call issuer"
	case "03": // Merchant not on file / Invalid Merchant
		return "03 - Invalid Merchant"
	case "05": // Pin Tries Exceeded
		return "05 - Do Not Honor"
	case "12": // Invalid Transaction
		return "12 - Invalid Txn"
	case "13": // Invalid Amount
		return "13 - Invalid Amount"
	case "14": // Invalid Account
		return "14 - Invalid Card"
	case "15":
		return ""
	case "19": // Retry the transaction
		return "19 - Retry the Txn"
	case "25": // Invalid Account (AMEX) / Declined
		return "25 - Declined"
	case "30": // Format Error
		return "30 - Format Error"
	case "31": // Transaction Not supported
		return "31 - Unsupported Txn"
	case "41": // Please Call - Lost Card
		return "41 - Please Call - LC"
	case "43": // Pick-Up Card (AMEX) / Please Call - Captured Card
		return "43 - Please Call - CC"
	case "51": // Declined
		return "51 - Declined"
	case "54": // Expired Card
		return "54 - Expired Card"
	case "55": // Incorrect Pin
		return "55 - Incorrect Pin"
	case "58": // Txn not allowed
		return "58 - Txn not allowed"
	case "65": // Declined - try contact
		return "65 - Perform Contact Txn"
	case "74": // Original Txn not found
		return "74 - Txn Unavailable"
	case "78": // Invalid amount
		return "78 - Invalid amount"
	case "89": // Invalid terminal
		return "89 - Invalid terminal"
	case "91": // Auth timed out
		return "91 - Auth timed out"
	case "94": // Txn already voided
		return "94 - Duplicate TXN"
	case "95": // Non-B24 Transaction Cancelled
		return "95 - Txn Cancelled"
	case "96": // Non-B24 Offline Decline
		return "96 - Declined"
	case "97": // Non-B24 Signature Mismatch
		return "97 - Signature Mismatch"
	case "98": // Non-B24 Card Removed
		return "98 - Card Removed"
	case "99": // Non-B24 Error
		return "99 - Comms Error"
	case "":
		return "Unknown Reason"
	default:
		return resultCode + " - Declined"
	}
}

func GetMerchantName(merchantId string) string {
	siteName, err := dal.GetSiteNameFromMerchantID(merchantId)
	if err != nil {
		return err.Error()
	}
	return siteName
}
