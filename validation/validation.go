package validation

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"nextgen-tms-website/dal"
	"nextgen-tms-website/resultCodes"
	"regexp"
	"strconv"

	rpcHelper "bitbucket.org/network-international/nextgen-libs/nextgen-helpers/rpcHelper"
)

var logger rpcHelper.LoggingClient

type Validator struct {
	dal dal.ValidationDal
}

func New(validationDal dal.ValidationDal) *Validator {
	return &Validator{dal: validationDal}
}

func (v *Validator) ValidateTid(tid string) (bool, error) {
	if tid == "" {
		return false, errors.New(resultCodes.GetErrorMsgByCode(resultCodes.TID_NOT_PRESENT))
	}
	// test for all 0s
	match, err := regexp.MatchString("^[0]+$", tid)
	if err != nil {
		return false, err
	}
	if match {
		return false, nil
	}

	match, err = regexp.MatchString("^\\d{8}$", tid)
	if err != nil {
		return false, err
	}

	if !match {
		return false, errors.New("TID must be numeric and must be of length 8 Digits")
	}

	tidInt, err := strconv.Atoi(tid)
	if err != nil {
		return false, errors.New("invalid TID entry")
	}

	if tidExists, resCode, _ := v.dal.CheckThatTidExists(tidInt); tidExists {
		return false, errors.New(resultCodes.GetErrorMsgByCode(resCode))
	}
	return match, nil
}

func (v *Validator) ValidateTidFormat(tid string) (bool, error) {
	if tid == "" {
		return false, errors.New(resultCodes.GetErrorMsgByCode(resultCodes.TID_NOT_PRESENT))
	}
	// test for all 0s
	match, err := regexp.MatchString("^[0]+$", tid)
	if err != nil {
		return false, err
	}
	if match {
		return false, nil
	}

	match, err = regexp.MatchString("^\\d{8}$", tid)
	if err != nil {
		return false, err
	}

	if !match {
		return false, errors.New("TID must be numeric and must be of length 8 Digits")
	}

	_, err = strconv.Atoi(tid)
	if err != nil {
		return false, errors.New("invalid TID entry")
	}
	return match, nil
}

func (v *Validator) ValidateSerialNumber(serialNumber string) (bool, error) {
	match, err := regexp.MatchString("^.{1,10}$", serialNumber)
	if err != nil {
		return false, err
	}
	return match, nil
}

func (v *Validator) validateUniqueness(element dal.DataElement, value string, profileId int) (bool, error) {

	isUnique, err := dal.GetIsUnique(element.ElementId, value, profileId)
	if err != nil {
		return false, err
	}

	return isUnique, nil

}

func (v *Validator) ValidateDataElement(element dal.DataElement, value string, profileId int, fileServerURL string) (bool, string) {
	// Store the element's original validation message, since some parts of this function modify the element struct.
	originalValidationMessage := element.ValidationMessage

	match, maxLength, unique := true, true, true

	if value == "" && element.IsAllowedEmpty {
		return true, ""
	}

	if element.Type == "STRING" && !element.IsAllowedEmpty && value == "" {
		if element.ValidationMessage != "" {
			return false, element.ValidationMessage
		} else {
			return false, "element must not be empty"
		}
	}

	if element.Type == "JSON" {
		if !element.IsAllowedEmpty && value == "" {
			return true, element.ValidationMessage
		} else if !v.isValidJson(value) {
			return false, "Invalid JSON"
		}
	}

	if element.ValidationExpression != "" {
		var err error
		match, err = regexp.MatchString(element.ValidationExpression, value)
		if err != nil {
			return false, err.Error()
		}
	}

	if element.MaxLength > 0 {
		if element.Type == "INTEGER" || element.Type == "LONG" {
			intVal, err := strconv.Atoi(value)
			if err != nil {
				return false, err.Error()
			}
			maxLength = intVal <= element.MaxLength
		} else {
			maxLength = len(value) <= element.MaxLength
		}
	}

	if element.Unique && match && maxLength {
		var err error
		unique, err = v.validateUniqueness(element, value, profileId)
		if err != nil {
			return false, err.Error()
		}
		if !unique {
			element.ValidationMessage = "Element must be unique"
		}
	}

	if match {
		switch element.Name {
		case `secondaryTid`:
			secondaryTid, err := strconv.Atoi(value)
			if err != nil {
				return false, err.Error()
			}

			if secondaryTid == 0 {
				element.ValidationMessage = originalValidationMessage
				match = false
			} else if tidExists, resCode, foundTidProfileId := v.dal.CheckThatTidExists(secondaryTid); tidExists {
				// We need to also check if the profileId of the found TID matches that of the TID we're dealing with now
				// otherwise we will be checking for uniqueness against the secondaryTID already belonging to this TID.
				// We should also check that the found tid is not the primary TID as the primary and secondary tid
				// cannot be the same.
				if foundTidProfileId == profileId && resCode == resultCodes.TID_NOT_UNIQUE_PRIMARY_TID_DUPLICATE {
					element.ValidationMessage = resultCodes.GetErrorMsgByCode(resCode)
					match = false
				}
			}
		case `secondaryMid`, `merchantNo`:
			if midExists, resCode, _ := v.dal.CheckThatMidExists(value); midExists && resCode != resultCodes.MID_DOES_NOT_EXIST {
				element.ValidationMessage = resultCodes.GetErrorMsgByCode(resCode)
				match = false
			}
		}

		// Validate File Size:
		if element.Type == "FILE" && element.FileMaxSize > 0 {
			limit := element.FileMaxSize
			match = v.ValidateFileSize(value, fileServerURL, limit)
			if !match {
				fileLimit := int(limit / 1000)
				element.ValidationMessage = "Image File too Large, must be " + strconv.Itoa(fileLimit) + "KB or less"
			}
		}

		// Validate Aspect Ratio:
		if element.Type == "FILE" && len(element.Options) == 0 && element.FileMinRatio > 0 && element.FileMaxRatio > 0 {
			if match {
				match = v.ValidateAspectRatio(value, fileServerURL, element.FileMinRatio, element.FileMaxRatio)
				if !match {
					element.ValidationMessage = fmt.Sprintf("Invalid Aspect Ratio, width / height should be between %.2f and %.2f", element.FileMinRatio, element.FileMaxRatio)

				}
			}
		}
	}

	return match && maxLength && unique, element.ValidationMessage
}

func (v *Validator) isValidJson(jsonString string) bool {
	return json.Valid([]byte(jsonString))
}

func (v *Validator) ValidateDirectQueryLimit(limit int) bool {

	if limit < 1 || limit > 10000 {
		return false
	}

	return true
}

func (v *Validator) ValidateString(name string) bool {
	if len(name) < 1 || len(name) > 30 {
		return false
	}

	match, err := regexp.MatchString("[!@#$%^&*(),.?\":{}|<>£_\\-'\\=\\+\\[\\]\\~\\;\\/]", name)
	if err != nil {
		return false
	}

	return !match
}

func (v *Validator) ValidateFileSize(name string, fileserverURL string, sizeLimit int) bool {
	fileName := name

	values := make(map[string][]string, 0)
	values["FileName"] = []string{fileName}

	response, err := http.PostForm(fileserverURL+"/getFileSize", values)
	if err != nil {
		return false
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return false
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	val, _ := strconv.Atoi(string(bodyBytes))
	if val > sizeLimit {
		return false
	}

	return true
}

func (v *Validator) ValidateAspectRatio(name string, fileserverURL string, minRatio, maxRatio float64) bool {
	fileName := name

	values := make(map[string][]string, 0)
	values["FileName"] = []string{fileName}

	response, err := http.PostForm(fileserverURL+"/getFileDimensions", values)
	if err != nil {
		return false
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return false
	}

	var data map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		logger.Error(err)
		return false
	}

	width := data["width"]
	height := data["height"]
	aspectRatio := width.(float64) / height.(float64)
	valid := aspectRatio >= minRatio && aspectRatio <= maxRatio
	return valid
}

func (v *Validator) ValidateMID(mid string) (bool, error) {
	match, err := regexp.MatchString("^[a-zA-Z0-9]{6,15}$", mid)
	if !match || err != nil {
		return false, err
	}

	match, err = regexp.MatchString("0{6,15}$", mid)
	if match || err != nil {
		return false, err
	}

	if midExists, resCode, _ := v.dal.CheckThatMidExists(mid); midExists {
		return false, errors.New(resultCodes.GetErrorMsgByCode(resCode))
	}
	return true, nil
}
