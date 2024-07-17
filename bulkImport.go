package main

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"log"
	"mime/multipart"
	"net/http"
	"nextgen-tms-website/common"
	"nextgen-tms-website/crypt"
	"nextgen-tms-website/dal"
	"nextgen-tms-website/entities"
	"nextgen-tms-website/models"
	"nextgen-tms-website/resultCodes"
	"nextgen-tms-website/services"
	"nextgen-tms-website/validation"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context"

	sliceHelper "bitbucket.org/network-international/nextgen-libs/nextgen-helpers/TypeComparisonHelpers/SliceComparisonHelpers"
)

type bulkImportError struct {
	statusCode int
	respErrMsg string
	logErrMsg  string
}

type tidIdentifier struct {
	tid    string
	siteId int
}

type siteInfo struct {
	siteId    int
	groupName string
}

const (
	TAG                          = "BulkImport: "
	BulkSiteUpdateType           = "BulkSiteUpdate"
	BulkTidUpdateType            = "BulkTidUpdate"
	BulkTidDeleteType            = "BulkTidDelete"
	BulkSiteDeleteType           = "BulkSiteDelete"
	BulkPaymentTidUploadType     = "BulkPaymentTidUpload"
	BulkPaymentServiceUploadType = "BulkPaymentServiceUpload"
	// Permitted bulk import file types
	FiletypeCsv = "text/plain"
	CsvSuffix   = ".csv"

	// Profile types
	Site = "site"
	Tid  = "tid"

	// These are all mapped to lower case
	SiteName            = "sitename"
	TID                 = "tid"
	SecondaryTID        = "dualcurrency.secondarytid"
	MerchantID          = "merchantid"
	SecondaryMerchantID = "dualcurrency.secondarymid"
	Serial              = "serial"

	// Using negative numbers to identify these so there are no clashes with "real" data elements
	TID_ELEMENT_ID = -iota
	SERIAL_NUMBER_ID
	TYPE_MID      = "mid"
	TYPE_TID      = "tid"
	TID_MAXLENGTH = 8
)

var (
	SiteUpload    BulkSiteUploadModel
	TidUpload     BulkTidUploadModel
	PsUpload      *entities.PaymentServiceGroupImportModel
	PsTermUpload  map[tidIdentifier]map[int]*dal.PaymentService
	DataGroups    map[int]string
	ProfileId     int
	ByteOrderMark = []byte{239, 187, 191}
)

func validateNewTids(newTids NewSitesElements, columns []DataColumn) (BulkTidUploadModel, error) {
	var validationResults BulkTidUploadModel
	var validationPasses NewSitesElements
	failed := false

	validationDal := new(siteValidatorDal)
	// The below TID map and sn map are used exclusively to keep track of the TIDs and SNs we're creating so that we can
	// identify rows with duplicate TIDs or SNs
	tidMap := make(map[string]string, len(newTids.getNewSites()))
	snMap := make(map[string]interface{}, len(newTids.getNewSites()))
	uniqueElementsMap := make(map[int]map[string]interface{}, len(newTids.getNewSites()))
	// Iterate through each new TID and validate its details and data elements
	for _, tid := range newTids.getNewSites() {
		if tidType, tidPresent := tidMap[tid.getTid()]; tidPresent {
			failed = true
			var validationFailure FailedProfileValidation
			validationFailure.setFailureReason("TID is not unique within the CSV")
			validationFailure.setFailedElementId(TID_ELEMENT_ID)
			validationFailure.setFailedElementName(tidType)
			validationFailure.setSite(tid)
			validationResults.setFailure(validationFailure)
			break
		}
		if _, snPresent := snMap[tid.getSerial()]; snPresent {
			failed = true
			var validationFailure FailedProfileValidation
			validationFailure.setFailureReason("Serial number is not unique within the CSV")
			validationFailure.setFailedElementId(SERIAL_NUMBER_ID)
			validationFailure.setFailedElementName("serial")
			validationFailure.setSite(tid)
			validationResults.setFailure(validationFailure)
			break
		}
		if tid.getTid() == tid.getSecondaryTid() {
			failed = true
			var validationFailure FailedProfileValidation
			validationFailure.setFailureReason("Primary TID is equal to the secondary TID")
			validationFailure.setFailedElementId(TID_ELEMENT_ID)
			validationFailure.setFailedElementName("TID and SecondaryTID")
			validationFailure.setSite(tid)
			validationResults.setFailure(validationFailure)
			break
		}

		tidMap[tid.getTid()] = "TID"
		if tid.getSecondaryTid() != "" {
			tidMap[tid.getSecondaryTid()] = "secondaryTID"
		}
		snMap[tid.getSerial()] = ""
		for _, element := range tid.getDataElements() {
			if element.Unique {
				if elem, present := uniqueElementsMap[element.ElementId]; present && !element.Ignore {
					if _, valuePresent := elem[element.Data]; valuePresent {
						failed = true
						var validationFailure FailedProfileValidation
						validationFailure.setFailureReason("Unique element duplicated within CSV")
						validationFailure.setFailedElementId(element.ElementId)
						validationFailure.setFailedElementName(element.Name)
						validationFailure.setSite(tid)
						validationResults.setFailure(validationFailure)
						break
					}
					// The data element is now assigned and so should not be present again in the CSV
					uniqueElementsMap[element.ElementId][element.Data] = ""
				} else {
					uniqueElementsMap[element.ElementId] = make(map[string]interface{}, 0)
					uniqueElementsMap[element.ElementId][element.Data] = ""
				}
			}
		}

		// First validate the MID the tid is present and numerical
		valid, siteProfileId, err := validateTemplateMid(validationDal, tid.getMid())
		if !valid {
			logging.Error(TAG, "MID validation failed as supplied MID is invalid")
			failed = true
			var validationFailure FailedProfileValidation
			validationFailure.setFailureReason(err.Error())
			validationFailure.setFailedElementName("MID")
			validationFailure.setSite(tid)
			validationResults.setFailure(validationFailure)
			break
		}
		// This will be used later on to lookup the site
		tid.SiteProfileId = siteProfileId

		// Ensure that the TID is valid and unique
		if valid, err := validation.New(validationDal).ValidateTid(tid.getTid()); !valid {
			logging.Error(TAG, err.Error())
			failed = true
			var validationFailure FailedProfileValidation
			validationFailure.setFailureReason(err.Error())
			validationFailure.setFailedElementId(TID_ELEMENT_ID)
			validationFailure.setFailedElementName("TID")
			validationFailure.setSite(tid)
			validationResults.setFailure(validationFailure)
			break
		}

		// Ensure that the serial number is valid and unique
		valid, err = validation.New(validationDal).ValidateSerialNumber(tid.getSerial())
		if err != nil {
			logging.Error(TAG, "Error thrown when attempting to validate Serial Number")
			return validationResults, errors.New(DatabaseAccessError)
		}
		if !valid {
			logging.Error(TAG, SerialInvalidFormat)
			failed = true
			var validationFailure FailedProfileValidation
			validationFailure.setFailureReason(SerialInvalidFormat)
			validationFailure.setFailedElementId(SERIAL_NUMBER_ID)
			validationFailure.setFailedElementName("TID")
			validationFailure.setSite(tid)
			validationResults.setFailure(validationFailure)
			break
		}

		exists, err := validationDal.CheckThatSerialNumberExists(tid.getSerial())
		if err != nil {
			logging.Error(TAG, err.Error())
			valid = false
			failed = true
			break
		}
		if exists {
			logging.Error(TAG, "SN already exists")
			failed = true
			var validationFailure FailedProfileValidation
			validationFailure.setFailureReason("SerialNumber is already in use")
			validationFailure.setFailedElementId(SERIAL_NUMBER_ID)
			validationFailure.setFailedElementName("SerialNumber")
			validationFailure.setSite(tid)
			validationResults.setFailure(validationFailure)
			break
		}

		// And finally validate the actual elements
		logging.Debug(TAG, fmt.Sprintf("Validating data elements for entry number %d", tid.getRef()))
		validationErrors, failedIndex := validateDataElements(validationDal, tid.getDataElementsAsMap(), -1, Tid, true, "")

		if len(validationErrors) > 0 {
			logging.Error(TAG, fmt.Sprintf("Validation of data element %v for new TID entry number %v has failed", failedIndex, tid.getRef()))
			var validationFailure FailedProfileValidation

			validationFailure.setFailedElementId(failedIndex)
			validationFailure.setFailureReason(validationErrors[0]) // Will always be in position 0 as the first validation failure is returned
			// Find the data element name for the failed validation
			failedElementName := ""
			for _, column := range columns {
				if column.getElementId() == failedIndex {
					failedElementName = column.getName()
				}
			}
			validationFailure.setFailedElementName(failedElementName)
			validationFailure.setSite(tid)

			failed = true
			validationResults.setFailure(validationFailure)
			// As soon as we hit a validation failure we want to cease further validations
			break
		} else {
			logging.Debug(TAG, fmt.Sprintf("All data elements and attributes successfully validated for entry number %v", tid.getRef()))
			validationPasses.NewSites = append(validationPasses.getNewSites(), tid)
		}
	}
	validationResults.setValidationResult(failed)
	validationResults.setPasses(validationPasses)

	columns = trimColumnsForTidImport(columns)

	// Add the columns to the pageModel
	for _, column := range columns {
		validationResults.addColumn(column)
	}

	// Prevents any sillyness with Go not iterating maps in order
	validationResults = sortTidResults(validationResults)

	return validationResults, nil
}

func trimColumnsForTidImport(columns []DataColumn) []DataColumn {
	// The objective here is to remove the 3 special case column
	// serial	merchantId	tid
	var trimmedColumns []DataColumn

	logging.Debug(TAG, "trimColumnsForTidImport")
	for _, column := range columns {
		if column.getName() == "serial" || column.getName() == "merchantId" || column.getName() == "tid" {
			continue
		}

		trimmedColumns = append(trimmedColumns, column)
	}

	return trimmedColumns
}

func sortTidResults(results BulkTidUploadModel) BulkTidUploadModel {
	// sort the data elements of the new TIDs by data element id
	for _, site := range results.Passes.getNewSites() {
		sort.Slice(site.DataElements, func(i, j int) bool { return site.DataElements[i].ElementId < site.DataElements[j].ElementId })
	}

	// Sort the columns by data element id
	sort.Slice(results.Columns, func(i, j int) bool { return results.Columns[i].ElementId < results.Columns[j].ElementId })

	return results
}

func buildNewTids(newTidElements NewSitesElements) (NewSitesElements, error) {
	var newTids NewSitesElements

	// Each of the rows of the passed in data represents a new TID to be added
	for _, entry := range newTidElements.NewSites {
		var newTid NewProfileEntry
		// Copy the reference over to the new TID
		newTid.setRef(entry.getRef())

		// We need to extract the identifiers (MID, TID and serial) from the "real" data elements
		for _, element := range entry.DataElements {
			switch element.ElementId {
			case TID_ELEMENT_ID:
				newTid.setTid(element.Data)
			case SERIAL_NUMBER_ID:
				newTid.setSerial(element.Data)
			default: //If the element ID does not relate to one of our special cases then it is a "real" data element
				switch element.GroupName + "." + element.Name {
				case "dualCurrency.secondaryTid":
					newTid.setSecondaryTid(element.Data)
				case "dualCurrency.secondaryMid":
					newTid.setSecondaryMid(element.Data)
					element.Ignore = true
				case "store.merchantNo":
					newTid.setMid(element.Data)
					element.Ignore = true
				case "store.name":
					newTid.setSiteName(element.Data)
				}
				newTid.DataElements = append(newTid.DataElements, element)
			}
		}

		newTids.NewSites = append(newTids.NewSites, newTid)
	}

	// As the bulk upload facility can be used to upload TIDs to multiple sites at once with potentially different
	// enabled data groups we need to ensure that any non-enabled data group elements are not stored
	var err error
	newTids.NewSites, err = ignoreDisabledDataGroupElements(newTids.NewSites)
	if err != nil {
		return NewSitesElements{}, err
	}

	return newTids, nil
}

// Obtains a list of unique MIDs from the passed in list of TIDs to be bulk uploaded
func buildUniqueMIDList(elements []NewProfileEntry) []string {
	var uniqueMIDs []string

	for _, element := range elements {
		present := false
		for _, mid := range uniqueMIDs {
			if element.Mid == mid {
				present = true
			}
		}
		if !present {
			uniqueMIDs = append(uniqueMIDs, element.Mid)
		}
	}

	return uniqueMIDs
}

// Obtains enabled data groups for each MID in the slice
func buildMIDDataGroups(mids []string) (map[string][]string, error) {
	var midDataGroups = map[string][]string{}

	for _, mid := range mids {
		// For each unique MID, obtain its active data groups
		siteDataGroups, err := dal.FetchSiteDataGroups(mid)
		if err != nil {
			logging.Error(TAG, err.Error())
			return midDataGroups, err
		}
		midDataGroups[mid] = siteDataGroups
	}

	return midDataGroups, nil
}

// Ensures that any data elements belonging to a non-enable data group are ignored
func ignoreDisabledDataGroupElements(newTids []NewProfileEntry) ([]NewProfileEntry, error) {
	var processedTIDs []NewProfileEntry

	uniqueMids := buildUniqueMIDList(newTids)

	activeDataGroups, err := buildMIDDataGroups(uniqueMids)
	if err != nil {
		logging.Error(TAG, err.Error())
		return processedTIDs, nil
	}

	for _, tid := range newTids {
		var tidElements []DataElement
		for _, element := range tid.DataElements {
			// If the element is already ignored then we do not need to do this
			if !element.Ignore {
				// Assume the data group is not enabled
				activeDataGroup := false
				for _, dataGroup := range activeDataGroups[tid.getMid()] {
					if dataGroup == element.GroupName {
						// If we find it is enabled, then flag as such
						activeDataGroup = true
					}
				}
				// If the element is from a non-active data group then ignore it
				element.Ignore = !activeDataGroup
			}
			tidElements = append(tidElements, element)
		}
		// Overwrite the TIDs elements with the processed elements
		tid.DataElements = tidElements
		// Add the TID to the return array
		processedTIDs = append(processedTIDs, tid)
	}

	return processedTIDs, nil
}

func bulkTidConstructor(records [][]string) (BulkTidUploadModel, error) {
	var results BulkTidUploadModel

	// the first entry in records contains column headers. These are written in the following format:
	// dataGroup.dataElement for example store.merchantNo, the only exceptions to this are the MID, TID and serial number, which do not
	// have a data groups
	// NOTE: these patterns exactly match that of the site import functionality on TMS to make it easier for NI to
	// create bulk upload templates

	// Extract the column headers and convert them into a usable struct
	columns := convertRecordsToColumns(records[0], true)
	// Obtain the data element id for all columns
	columns, err := fetchColumnDataElementIds(columns)
	if err != nil {
		return results, err
	}

	// Convert the data rows into map[int]string to match the format of the template
	newTidElements, err := convertDataRowsToElements(records[1:], columns, false)
	if err != nil {
		return results, err
	}

	// Convert the raw tidElements into representations of a single TID
	newTids, err := buildNewTids(newTidElements)
	if err != nil {
		return results, err
	}

	// Validate all of the data elements as well as the MID, TID and serial of each new TID
	return validateNewTids(newTids, columns)
}

func validateExp(str string, exp string) bool {
	Regex := regexp.MustCompile(exp)
	return Regex.MatchString(str)
}

func validateOptions(options []string, dataElementValue []interface{}) bool {
	for _, val := range dataElementValue {
		if !sliceHelper.SlicesOfStringContains(options, val.(string)) {
			return false
		}
	}
	return true
}

// Validation function for Bulk Tid/Site Update results
func bulkUpdateValidationFunc(entry string, dataelement string, reason string) models.ValidationSet {
	var FailureVal models.ValidationSet
	FailureVal.EntryNo = entry
	FailureVal.DataElement = dataelement
	FailureVal.FailureReason = reason
	return FailureVal
}

func bulkTidUpdateHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	logging.Debug(TAG, "TID bulk update initiated")

	var FailureValSet []models.ValidationSet
	var TidUpdate models.BulkUpdateVal

	err := r.ParseMultipartForm(10000)
	if err != nil {
		logging.Error(TAG, err.Error())
		http.Error(w, uploadFileError, http.StatusBadRequest)
		return
	}

	// Validate that a file has been attached
	if len(r.MultipartForm.File) < 1 {
		logging.Warning(TAG, "TID upload initiated without file being present")
		http.Error(w, bulkTidUpdateMissingFile, http.StatusBadRequest)
		return
	}

	// Extract the file from the request and obtain the name
	logging.Debug(TAG, "Attempting to extract file from http.Request")

	file, handler, err := r.FormFile("file")
	if err != nil {
		logging.Error(TAG, err.Error())
		http.Error(w, uploadFileError, http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileName := handler.Filename

	logging.Debug(TAG, fmt.Sprintf("File: %v has been uploaded", fileName))

	// Validate that the filetype is csv
	// 512 bytes only because DetectContentType (used in ValidateFileType) only reads up to 512 bytes
	buff := make([]byte, 512)
	if _, err = file.Seek(0, 0); err != nil {
		logging.Error(TAG, err.Error())
		http.Error(w, html.EscapeString(FailedFileRead+fileName), http.StatusInternalServerError)
		return
	}

	if _, err = file.Read(buff); err != nil {
		logging.Error(TAG, err.Error())
		http.Error(w, html.EscapeString(FailedFileRead+fileName), http.StatusInternalServerError)
		return
	}

	// Validates the file type, based on File type
	isCSVFile := strings.HasSuffix(strings.ToLower(fileName), CsvSuffix)
	if !isCSVFile {
		logging.Warning(TAG, "Incorrect filetype uploaded")
		http.Error(w, IncorrectFileTypeCSV, http.StatusInternalServerError)
		return
	}

	logging.Debug(TAG, fmt.Sprintf("File %v has passed type validation", fileName))

	logging.Debug(TAG, "Resetting file read offset")
	// Need to reset the offset after checking for filetype
	if _, err = file.Seek(0, 0); err != nil {
		logging.Error(TAG, err.Error())
		http.Error(w, html.EscapeString(FailedFileRead+fileName), http.StatusInternalServerError)
		return
	}

	logging.Debug(TAG, "Parsing CSV data")
	// Parse the entries from the csv along with the column
	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		logging.Error(TAG, err.Error())
		http.Error(w, html.EscapeString(FailedFileRead+fileName), http.StatusInternalServerError)
		return
	} else if len(records) == 0 || len(records[1:]) < 1 {
		logging.Error(TAG, "No records found in uploaded CSV file")
		http.Error(w, NoColumnsFound, http.StatusInternalServerError)
		return
	}

	n := records[:]

	buf := &bytes.Buffer{}
	writer := csv.NewWriter(buf)
	writer.WriteAll(n)

	dataElementMap := make(map[int]models.UpdateDataElement, 0)
	for index, header := range records[0] {
		header = strings.TrimSpace(header)
		var updateDataElement models.UpdateDataElement
		if index == 0 {
			header = string(bytes.TrimPrefix([]byte(strings.ToLower(header)), common.ByteOrderMark))
			if header != "tid" {
				///validation error
				logging.Error(TAG, "Tid header validation failed in uploaded CSV file")
				http.Error(w, "Please provide the correct header name as tid", http.StatusInternalServerError)
				return
			}
		} else { //header columns as keys
			if header == "" || !strings.Contains(header, ".") {
				logging.Error(TAG, "Data group header validation failed in uploaded CSV file")
				http.Error(w, "Data Group Name cannot be empty in file header or must be of length 2; Separated by .", http.StatusInternalServerError)
				return
			}
			elementName := strings.Split(header, ".")
			valExp, valMsg, isAllowEmpty, dataType, options, elementId, err := dal.GetDataElement(elementName[0], elementName[1])

			//Also check the fun does not return the empty value's in case of incorrect dataGroup name or data element
			if err != nil || (valExp.String == "" && valMsg.String == "" && dataType == "") {
				///validation error
				logging.Error(TAG, "failed to locate data elements from DB")
				http.Error(w, "Please provide the correct data group or data element name", http.StatusInternalServerError)
				return
			}
			updateDataElement.DataGroupName = elementName[0]
			updateDataElement.DataType = dataType
			updateDataElement.IsAllowEmpty = isAllowEmpty
			updateDataElement.ValExp = valExp
			updateDataElement.ValMsg = valMsg
			updateDataElement.DataElementId = elementId
			updateDataElement.Options = options
			dataElementMap[index] = updateDataElement
		}
	}
	//key's are TID and values are data group's associated with that TID
	dataGroupMap := make(map[int][]string)

	header := records[0]

	for _, rows := range records[1:] {
		var tidInt int
		for column, value := range rows {
			if column == 0 {
				//validate tid
				tidInt, err = strconv.Atoi(value)
				if err != nil {
					logging.Error(TAG, "Error while Tid conversion ")
					FailureValSet = append(FailureValSet, bulkUpdateValidationFunc(value, header[column], "Error while Tid conversion"))
					continue
				}
				tidExits, _, _ := dal.CheckThatTidExists(tidInt)
				if !tidExits {
					logging.Error(TAG, "Tid does not exists")
					FailureValSet = append(FailureValSet, bulkUpdateValidationFunc(rows[0], header[column], "Tid does not exist"))
					break
				}
				continue
			} else {
				elementId := dataElementMap[column].DataElementId
				isElement, _, err := dal.CheckIfDataElementExistsinTidData(tidInt, elementId)
				if err != nil {
					logging.Error(TAG, "Error retrieving data element id")
					continue
				}
				if !isElement {
					logging.Error(TAG, "Data Element does not exist for this TID")
					FailureValSet = append(FailureValSet, bulkUpdateValidationFunc(rows[0], header[column], "Data Element does not exist for this TID"))
					continue
				}
			}

			if value == "" {
				continue
			} else if value == "NULL" {
				//check the data element mandatory
				isAllowEmpty := dataElementMap[column].IsAllowEmpty
				if !isAllowEmpty {
					//log the error and continue
					logging.Error(TAG, "mandatory data element field is empty")
					FailureValSet = append(FailureValSet, bulkUpdateValidationFunc(rows[0], header[column], "Mandatory data element field is empty"))
					continue

				}
			} else {
				dataElement := dataElementMap[column] //getting header from csv though map
				dataGroupsExists, ok := dataGroupMap[tidInt]
				if !ok {
					dataGroupsExists, err = dal.FetchTidDataGroups(tidInt)
					if err != nil {
						logging.Error(TAG, "error while fetching the data group's from DB")
						FailureValSet = append(FailureValSet, bulkUpdateValidationFunc(rows[0], header[column], "Error while fetching the data group's from DB"))

					}
					dataGroupMap[tidInt] = dataGroupsExists
				}

				//check the data group is enable
				if sliceHelper.SlicesOfStringContains(dataGroupsExists, dataElement.DataGroupName) {
					switch dataElement.DataType {
					case "STRING":
						if dataElement.ValExp.Valid {
							if validator := validateExp(value, dataElement.ValExp.String); !validator {
								logging.Error(TAG, "validation expression does not match")
								FailureValSet = append(FailureValSet, bulkUpdateValidationFunc(rows[0], header[column], dataElement.ValMsg.String))

							}
						}
						if dataElement.Options != "" {
							optionsStr := strings.Split(dataElement.Options, "|")
							if !sliceHelper.SlicesOfStringContains(optionsStr, value) {
								FailureValSet = append(FailureValSet, bulkUpdateValidationFunc(rows[0], header[column], "Please provide correct value's or options"))
							}
						}
					case "BOOLEAN":
						if strings.ToLower(value) == "true" || strings.ToLower(value) == "false" {
							logging.Debug("successfully validated the data element type ")
						} else {
							logging.Error(TAG, "could not convert data element value type is mismatched")
							FailureValSet = append(FailureValSet, bulkUpdateValidationFunc(rows[0], header[column], "Could not convert data element as value type expected is true or false"))
						}
					case "INTEGER":
						if _, err := strconv.Atoi(value); err == nil {
							if dataElement.ValExp.Valid {
								if validator := validateExp(value, dataElement.ValExp.String); !validator {
									logging.Error(TAG, "validation expression does not match")
									FailureValSet = append(FailureValSet, bulkUpdateValidationFunc(rows[0], header[column], dataElement.ValMsg.String))
								}
							}
						} else {
							logging.Error(TAG, "could not convert data element value of type Integer")
							FailureValSet = append(FailureValSet, bulkUpdateValidationFunc(rows[0], header[column], "Could not convert data element value as expected type is Integer"))
						}

					case "LONG":
						if _, err := strconv.ParseInt(value, 10, 64); err == nil {
							if dataElement.ValExp.Valid {
								if validator := validateExp(value, dataElement.ValExp.String); !validator {
									logging.Error(TAG, "validation expression does not match")
									FailureValSet = append(FailureValSet, bulkUpdateValidationFunc(rows[0], header[column], dataElement.ValMsg.String))
								}
							}
						} else {
							logging.Error(TAG, "could not convert data element value of type LONG ")
							FailureValSet = append(FailureValSet, bulkUpdateValidationFunc(rows[0], header[column], "Could not convert data element as expected type is LONG"))
						}

					case "JSON":
						var jsonData []interface{}
						if err := json.Unmarshal([]byte(value), &jsonData); err != nil {
							logging.Error(TAG, "could not convert data element value type of JSON")
							FailureValSet = append(FailureValSet, bulkUpdateValidationFunc(rows[0], header[column], "Could not convert data element as expected type is JSON"))
						} else {
							logging.Debug(TAG, "successfully decoded the json data element")
							if dataElement.Options != "" {
								optionsStr := strings.Split(dataElement.Options, "|")
								if !validateOptions(optionsStr, jsonData) {
									FailureValSet = append(FailureValSet, bulkUpdateValidationFunc(rows[0], header[column], "Please provide correct value's inside list data"))
								}
							}
						}

					default:
						logging.Error(TAG, fmt.Sprintf("The other %v type of data element %v is not allowed", tidInt, dataElement.DataType))
						FailureValSet = append(FailureValSet, bulkUpdateValidationFunc(rows[0], header[column], fmt.Sprintf("The other %v type of data element %v is not allowed", tidInt, dataElement.DataType)))

					}
				} else {
					logging.Error(TAG, fmt.Sprintf("data group is not enable for tid: %v", tidInt))
					FailureValSet = append(FailureValSet, bulkUpdateValidationFunc(rows[0], dataElement.DataGroupName, fmt.Sprintf("Data group is not enabled for tid: %v", tidInt)))

				}
			}
		}
	}

	fileName = time.Now().Format("20060102150405_") + fileName
	logging.Debug(TAG, "Renamed file to : "+fileName)

	if err := sendFileToFileServer(buf.Bytes(), fileName, BulkTidUpdateType); err != nil {
		logging.Error(TAG, "Bulk Tid Update File Upload Failed : ", err.Error())
		http.Error(w, uploadFileError, http.StatusInternalServerError)
		return
	}

	if len(FailureValSet) > 0 {
		TidUpdate.UpdateStatus = false
		TidUpdate.Validations = FailureValSet
		err = dal.InsertBulkApproval(fileName, BulkTidUpdateType, tmsUser.Username, common.UploadChangeType, -1)
		if err != nil {
			logging.Error(TAG, "Unable to add to bulk approval: ", err.Error())
			http.Error(w, DatabaseTxnError, http.StatusInternalServerError)
			return
		}
	} else {
		TidUpdate.UpdateStatus = true
		err = dal.InsertBulkApproval(fileName, BulkTidUpdateType, tmsUser.Username, common.UploadChangeType, 0)
		if err != nil {
			logging.Error(TAG, "Unable to add to bulk approval: ", err.Error())
			http.Error(w, DatabaseTxnError, http.StatusInternalServerError)
			return
		}
	}

	renderPartialTemplate(w, r, "bulkTidUpdateResults", TidUpdate, tmsUser)
}

func bulkTidImportHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	logging.Debug(TAG, "TID bulk upload initiated")

	err := r.ParseMultipartForm(10000)
	if err != nil {
		logging.Error(TAG, err.Error())
		http.Error(w, uploadFileError, http.StatusBadRequest)
		return
	}

	// Validate that a file has been attached
	if len(r.MultipartForm.File) < 1 {
		logging.Warning(TAG, "TID upload initiated without file being present")
		http.Error(w, TidBulkUpdateMissingFile, http.StatusBadRequest)
		return
	}

	// Extract the file from the request and obtain the name
	logging.Debug(TAG, "Attempting to extract file from http.Request")

	file, handler, err := r.FormFile("file")
	if err != nil {
		logging.Error(TAG, err.Error())
		http.Error(w, uploadFileError, http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileName := handler.Filename

	logging.Debug(TAG, fmt.Sprintf("File: %v has been uploaded", fileName))

	// Validate that the filetype is csv
	nameSeparated := strings.Split(fileName, ".")
	if nameSeparated[len(nameSeparated)-1] != "csv" {
		logging.Warning(TAG, "Incorrect filetype uploaded")
		http.Error(w, IncorrectFileTypeCSV, http.StatusInternalServerError)
		return
	}
	logging.Debug(TAG, fmt.Sprintf("File %v has passed type validation", fileName))

	logging.Debug(TAG, "Resetting file read offset")
	// Need to reset the offset after checking for filetype
	if _, err = file.Seek(0, 0); err != nil {
		logging.Error(TAG, err.Error())
		http.Error(w, html.EscapeString(FailedFileRead+fileName), http.StatusInternalServerError)
		return
	}

	logging.Debug(TAG, "Parsing CSV data")
	// Parse the entries from the csv along with the column headers
	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		logging.Error(TAG, err.Error())
		http.Error(w, FailedFileRead, http.StatusInternalServerError)
		return
	} else if len(records) == 0 || len(records[1:]) < 1 {
		logging.Error(TAG, "No records found in uploaded CSV file")
		http.Error(w, NoColumnsFound, http.StatusInternalServerError)
		return
	}

	TidUpload, err = bulkTidConstructor(records)
	if err != nil {
		logging.Error(TAG, err.Error())
		http.Error(w, NoColumnsFound, http.StatusInternalServerError)
		return
	}

	renderPartialTemplate(w, r, "bulkTidUploadResultsPartial", TidUpload, tmsUser)
}

/*
*
This method closely follows the flow of the saveNewSite() method within tms.go only adapted for bulk import of sites
*/
func saveSite(name string, version int, tmsUser *entities.TMSUser, chainId int, elements map[int]string, dataGroups []string) error {

	logging.Debug(TAG, fmt.Sprintf("Attempting to save new site named %v", name))

	// Save the new site
	newProfileId, siteId, err := dal.SaveNewSite(name, version, tmsUser.Username, chainId)

	var result bool
	if err == nil {
		//Save the data elements to the new site
		result = saveDataElements(int(newProfileId), int(siteId), elements, 1, tmsUser, true, ProfileId)
		// And then record the new site creation in change history
		dal.RecordSiteToHistory(int(newProfileId), "Site Created", tmsUser.Username, dal.ApproveCreate, 1)
	} else {
		logging.Error(TAG, err.Error())
		return errors.New(saveSiteError)
	}

	// If we have successfully saved the data elements then add the data groups
	if result {
		services.AddDataGroupsToProfile(int(newProfileId), dataGroups, tmsUser)
	}

	logging.Debug(TAG, fmt.Sprintf("Successfully saved new site named %v", name))
	return nil
}

/*
*
Commits the now validated new sites to the DB one at a time
*/
func commitBulkSiteUploadHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	logging.Debug(TAG, "Site bulk upload initiated")

	logging.Debug(TAG, fmt.Sprintf("Attempting to fetch chain id for the profile %d", ProfileId))
	// Fetch the id of the chain to which the site belongs
	chainId, err := dal.GetChainIdFromSiteId(ProfileId)
	if err != nil {
		logging.Warning(TAG, fmt.Sprintf("Failed to fetch chain id for the template site. Error: %v", err.Error()))
		http.Error(w, ErrorFetchingSiteData, http.StatusInternalServerError)
		return
	}
	acquirerId, err := dal.GetAcquirerIdFromChainId(chainId)
	if err != nil {
		logging.Warning(TAG, fmt.Sprintf("Failed to fetch acquirer id for the template site. Error: %v", err.Error()))
		http.Error(w, ErrorFetchingSiteData, http.StatusInternalServerError)
		return
	}

	err = insertAllImportedSites(SiteUpload.Passes, chainId, acquirerId, extractKeysToArray(DataGroups), tmsUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func locateColumn(columns []DataColumn, field string) int {
	for index, column := range columns {
		if column.getName() == field {
			return index
		}
	}

	return -1
}

func createTIDs(newTIDs []NewProfileEntry, user *entities.TMSUser) error {
	var preparedStatements []*sql.Stmt
	db, err := dal.GetDB()
	if err != nil {
		return err
	}

	var tidInsertLines []string
	var tidSiteInsertLines []string
	var historyInsertLines []string
	tidInsertArgs := make([]interface{}, 0)
	tidSiteInsertArgs := make([]interface{}, 0)
	historySiteInsertArgs := make([]interface{}, 0)
	for _, tid := range newTIDs {
		tidInsertLines = append(tidInsertLines, "(?, ?, ?, ?)")
		tidSiteInsertLines = append(tidSiteInsertLines, "(?, (SELECT site_id FROM site_profiles WHERE profile_id = ?), NULL, NOW())")
		historyInsertLines = append(historyInsertLines, "(?, 1, 5, '', 'TID Created', NOW(), NOW(), ?, ?, 1, ?, ?, get_site_acquirer_name_by_site_profile_id(?))")

		tidInsertArgs = append(tidInsertArgs, tid.getTid())
		tidInsertArgs = append(tidInsertArgs, tid.getSerial())

		profileId, err := dal.GetSiteIDFromMerchantID(tid.getMid())
		if err != nil || profileId == "" {
			logging.Error(fmt.Sprintf("cannot not resolve the site id for the mid %s", tid.getMid()))
			return err
		}

		siteId, err := strconv.Atoi(profileId)
		if err != nil {
			logging.Error(fmt.Sprintf("cannot convert the site id (%s) to numeric", profileId))
			return err
		}

		tidId, err := strconv.Atoi(tid.getTid())
		if err != nil {
			logging.Error(fmt.Sprintf("cannot convert the tid id (%d) to numeric", tidId))
			return err
		}

		auto, autoTime, err := services.GetEODAutoTime(siteId, tidId)
		if err != nil {
			logging.Error(fmt.Sprintf("Error while getting EODAutoTime for tid id (%d) , site id (%d)", tidId, siteId), err)
			return err
		}

		tidInsertArgs = append(tidInsertArgs, auto)
		tidInsertArgs = append(tidInsertArgs, autoTime)

		tidSiteInsertArgs = append(tidSiteInsertArgs, tid.getTid())
		tidSiteInsertArgs = append(tidSiteInsertArgs, tid.SiteProfileId)

		historySiteInsertArgs = append(historySiteInsertArgs, tid.SiteProfileId)
		historySiteInsertArgs = append(historySiteInsertArgs, user.Username)
		historySiteInsertArgs = append(historySiteInsertArgs, user.Username)
		historySiteInsertArgs = append(historySiteInsertArgs, tid.getTid())
		historySiteInsertArgs = append(historySiteInsertArgs, tid.getMid())
		historySiteInsertArgs = append(historySiteInsertArgs, tid.SiteProfileId)
	}

	txn, err := db.Begin()
	if err != nil {
		return err
	}

	gracefullyExit := func(exitMessage string, exitError error) error {
		if exitMessage != "" {
			logging.Debug(fmt.Sprintf("Finished createTIDs: '%s'", exitMessage))
		}
		if exitError != nil {
			logging.Error(fmt.Sprintf("Rolling back transaction due to error - %s", exitError.Error()))
		}
		if exitError != nil {
			err = txn.Rollback()
			if err != nil {
				logging.Error("Error rolling back transaction")
			}
		} else {
			err = txn.Commit()
			if err != nil {
				logging.Error("Error committing transaction, rolling back transaction")
				err = txn.Rollback()
				if err != nil {
					logging.Error("Error rolling back transaction")
				}
			}
		}
		for _, stmt := range preparedStatements {
			stmt.Close()
		}
		if exitError != nil {
			return errors.New(DatabaseTxnError)
		} else {
			return nil
		}
	}

	logging.Debug(fmt.Sprintf("Inserting '%d' rows into the tid table", len(newTIDs)))
	_, err = txn.Exec(fmt.Sprintf(`
		INSERT INTO tid (tid_id, serial, eod_auto, auto_time)
		VALUES
    	%s;`,
		strings.Join(tidInsertLines, ", ")), tidInsertArgs...)
	if err != nil {
		return gracefullyExit("An error occured inserting the bulk import TIDs into the tid table", err)
	}

	logging.Debug(fmt.Sprintf("Inserting '%d' rows into the tid_site table", len(newTIDs)))
	_, err = txn.Exec(fmt.Sprintf(`
		INSERT INTO tid_site (tid_id, site_id, tid_profile_id, updated_at)
		VALUES
    	%s;`,
		strings.Join(tidSiteInsertLines, ", ")), tidSiteInsertArgs...)
	if err != nil {
		return gracefullyExit("An error occurred inserting the bulk import TIDs into the tid_site table", err)
	}

	logging.Debug(fmt.Sprintf("Inserting '%d' rows into the approvals table", len(newTIDs)))
	_, err = txn.Exec(fmt.Sprintf(`
		INSERT INTO approvals(profile_id, data_element_id, change_type, current_value, new_value, created_at, approved_at, created_by, approved_by, approved, tid_id, merchant_id, acquirer)
    	VALUES
    	%s;`,
		strings.Join(historyInsertLines, ", ")), historySiteInsertArgs...)
	if err != nil {
		return gracefullyExit("An error occurred inserting the bulk import TIDs into the approvals table", err)
	}

	type elementValueStructure struct {
		priority         int
		elementId        int
		dataType         string
		overriddenBySite bool
		value            string
	}
	siteDataElementCache := make(map[string]map[int]elementValueStructure, 0)
	logging.Debug("Preparing siteData statement")
	siteDataStatement, err := txn.Prepare(`
		SELECT
		    sd.priority,
		    sd.data_element_id,
		    sd.datavalue,
		    de.datatype
		FROM site_data sd
		INNER JOIN data_element de ON sd.data_element_id = de.data_element_id
		WHERE sd.site_id = (SELECT sp.site_id FROM site_profiles sp WHERE sp.profile_id = ?)`)
	if err != nil {
		return gracefullyExit("An error occurred preparing the siteData statement", err)
	}
	preparedStatements = append(preparedStatements, siteDataStatement)

	createTidProfileStatement, err := txn.Prepare(`
		INSERT INTO profile (profile_type_id, name, version, updated_at, updated_by, created_at, created_by)
		VALUE 
		(5, ?, 1, NOW(), ?, NOW(), ?)`)
	if err != nil {
		return gracefullyExit("An error occurred preparing the createTidProfileStatement statement", err)
	}
	preparedStatements = append(preparedStatements, createTidProfileStatement)

	assignTidProfileStatement, err := txn.Prepare("UPDATE tid_site SET tid_profile_id = ? WHERE tid_id = ?")
	if err != nil {
		return gracefullyExit("An error occurred preparing the assignTidProfileStatement statement", err)
	}
	preparedStatements = append(preparedStatements, assignTidProfileStatement)

	createOverrideCreatedHistoryStatement, err := txn.Prepare(`
		INSERT INTO approvals(profile_id, data_element_id, change_type, current_value, new_value, created_at, approved_at, created_by, approved_by, approved, tid_id, merchant_id, acquirer)
    	VALUE
		(?, 1, 5, null, 'Override Created', NOW(), NOW(), ?, ?, 1, ?, null, get_site_acquirer_name_by_tid(?))`)
	if err != nil {
		return gracefullyExit("An error occurred preparing the createOverrideCreatedHistoryStatement statement", err)
	}
	preparedStatements = append(preparedStatements, createOverrideCreatedHistoryStatement)

	createTidOverrideDataGroups, err := txn.Prepare(`
		INSERT INTO profile_data_group (profile_id, data_group_id, version, updated_at, updated_by, created_at, created_by)
		SELECT ?, data_group_id, 1, NOW(), ?, NOW(), ?
		FROM data_group
		WHERE name IN ('store', 'modules', 'userMgmt', 'core', 'opi', 'dualCurrency', 'nol')`)
	if err != nil {
		return gracefullyExit("An error occurred preparing the createTidOverrideDataGroups statement", err)
	}
	preparedStatements = append(preparedStatements, createTidOverrideDataGroups)

	for _, tid := range newTIDs {
		if _, present := siteDataElementCache[tid.getMid()]; !present {
			siteDataElementCache[tid.getMid()] = make(map[int]elementValueStructure, 0)
			rows, err := siteDataStatement.Query(tid.SiteProfileId)
			if err != nil {
				return gracefullyExit(fmt.Sprintf("An error occured executing the siteData statement for siteProfileId '%v'", tid.SiteProfileId), err)
			}

			for rows.Next() {
				var val elementValueStructure
				err = rows.Scan(&val.priority, &val.elementId, &val.value, &val.dataType)
				if err != nil {
					return gracefullyExit(fmt.Sprintf("An error occured scaning the result for siteData query for siteProfileId '%v'", tid.SiteProfileId), err)
				}
				// priority 2 is the site level
				val.overriddenBySite = val.priority == 2
				if currentValue, present := siteDataElementCache[tid.getMid()][val.elementId]; !present || (present && currentValue.priority > val.priority) {
					if val.dataType == "JSON" {
						val.value = dal.CleanseJSON(val.value)
					}
					siteDataElementCache[tid.getMid()][val.elementId] = val
				}
			}
		}

		// Check if we need to create an override
		tidOverride := false
		for i, de := range tid.getDataElements() {
			if de.Ignore || de.ElementId < 0 {
				continue
			}
			cachedValue, present := siteDataElementCache[tid.getMid()][de.ElementId]
			if !present && de.Data != "" {
				tidOverride = true
				break
			}
			if cachedValue.dataType == "JSON" {
				tid.DataElements[i].Data = dal.CleanseJSON(de.Data)
				cachedValue.value = dal.CleanseJSON(cachedValue.value)
			}
			if de.Data != cachedValue.value {
				tidOverride = true
				break
			}
		}

		if tidOverride {
			var overrideProfileId int
			if res, err := createTidProfileStatement.Exec(tid.getTid(), user.Username, user.Username); err != nil {
				return gracefullyExit(fmt.Sprintf("An error occured executing the createTidProfileStatement statement for siteProfileId '%v'", tid.SiteProfileId), err)
			} else {
				id, err := res.LastInsertId()
				if err != nil {
					return gracefullyExit("A problem occurred retrieving the profileId of the inserted TID override", err)
				}
				overrideProfileId = int(id)
			}
			if _, err = assignTidProfileStatement.Exec(overrideProfileId, tid.Tid); err != nil {
				return gracefullyExit("A problem occurred assigning the override profileID to the tid_site record", err)
			}
			if _, err = createOverrideCreatedHistoryStatement.Exec(overrideProfileId, user.Username, user.Username, tid.Tid, tid.Tid); err != nil {
				return gracefullyExit("A problem occurred creating the history entry for creating the TID override", err)
			}

			if _, err = createTidOverrideDataGroups.Exec(overrideProfileId, user.Username, user.Username); err != nil {
				return gracefullyExit("An error occurred executing statement createTidOverrideDataGroups", err)
			}

			var overrideProfileDataInsertLines []string
			overrideProfileDataParams := make([]interface{}, 0)
			for _, overrideDe := range tid.getDataElements() {
				if overrideDe.Ignore || overrideDe.ElementId < 0 {
					continue
				}
				if overrideDe.Encrypted {
					overrideDe.Data = crypt.Encrypt(overrideDe.Data)
				}
				insertLine := "(?, ?, ?, 1, NOW(), ?, NOW(), ?, 1, ?, ?)"

				overrideProfileDataInsertLines = append(overrideProfileDataInsertLines, insertLine)
				overrideProfileDataParams = append(overrideProfileDataParams, overrideProfileId)
				overrideProfileDataParams = append(overrideProfileDataParams, overrideDe.ElementId)
				if isDataBoolean(overrideDe.Data) {
					overrideProfileDataParams = append(overrideProfileDataParams, strings.ToLower(overrideDe.Data))
				} else {
					overrideProfileDataParams = append(overrideProfileDataParams, overrideDe.Data)
				}
				overrideProfileDataParams = append(overrideProfileDataParams, user.Username)
				overrideProfileDataParams = append(overrideProfileDataParams, user.Username)
				if overrideDe.Overriden {
					overrideProfileDataParams = append(overrideProfileDataParams, "1")
				} else {
					overrideProfileDataParams = append(overrideProfileDataParams, "0")
				}
				if overrideDe.Encrypted {
					overrideProfileDataParams = append(overrideProfileDataParams, "1")
				} else {
					overrideProfileDataParams = append(overrideProfileDataParams, "0")
				}
			}

			_, err = txn.Exec(fmt.Sprintf(`
				INSERT INTO profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by,created_at, created_by, approved, overriden, is_encrypted)
				VALUES 
				%s;
			`, strings.Join(overrideProfileDataInsertLines, ", ")), overrideProfileDataParams...)
			if err != nil {
				return gracefullyExit("An error occurred inserting the profile_data", err)
			}
		}
	}

	return gracefullyExit("Successfully created all TIDs", nil)
}

// Determines if the data in the element is a boolean. We need this because true != TRUE in SQL but we can't simply
// make all data elements lower case.
func isDataBoolean(data string) bool {
	switch strings.ToLower(data) {
	case "true", "false":
		return true
	}
	return false
}

func commitBulkTidUploadHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	logging.Debug(TAG, "commitBulkTidUploadHandler")

	if err := createTIDs(TidUpload.Passes.NewSites, tmsUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func extractKeysToArray(keymap map[int]string) []string {
	var keyArray []string
	for k, _ := range keymap {
		keyArray = append(keyArray, strconv.Itoa(k))
	}
	return keyArray
}

/*
*
Retrieves the the data element for a given data element
*/
func getDataFromElementName(columns []DataColumn, elements map[int]string, name string) string {
	for _, column := range columns {
		if column.getName() == name {
			return elements[column.getElementId()]
		}
	}
	return ""
}

/*
*
Renders the Bulk Import framework template
*/
func bulkImportHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	logging.Debug(TAG, "bulkImportHandler()")
	renderHeader(w, r, tmsUser)
	renderTemplate(w, r, "bulkUpload", nil, tmsUser)
}

// site update handler
func siteUpdateHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	logging.Debug(TAG, "Site bulk update initiated")

	var FailureValSet []models.ValidationSet
	var SiteUpdate models.BulkUpdateVal

	err := r.ParseMultipartForm(10000)
	if err != nil {
		logging.Error(TAG, err.Error())
		http.Error(w, uploadFileError, http.StatusBadRequest)
		return
	}
	// Validate that a file has been attached
	if len(r.MultipartForm.File) < 1 {
		logging.Warning(TAG, "Site bulk update initiated without file being present")
		http.Error(w, BulksiteUpdateMissingFile, http.StatusBadRequest)
		return
	}

	// Extract the file from the request and obtain the name
	logging.Debug(TAG, "Attempting to extract file from http.Request")

	file, handler, err := r.FormFile("file")

	if err != nil {
		logging.Error(TAG, err.Error())
		http.Error(w, uploadFileError, http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileName := handler.Filename

	logging.Debug(TAG, fmt.Sprintf("File: %v has been uploaded", fileName))

	// Validate that the filetype is csv
	// 512 bytes only because DetectContentType (used in ValidateFileType) only reads up to 512 bytes
	buff := make([]byte, 512)
	if _, err = file.Seek(0, 0); err != nil {
		logging.Error(TAG, err.Error())
		http.Error(w, html.EscapeString(FailedFileRead+fileName), http.StatusInternalServerError)
		return
	}

	if _, err = file.Read(buff); err != nil {
		logging.Error(TAG, err.Error())
		http.Error(w, html.EscapeString(FailedFileRead+fileName), http.StatusInternalServerError)
		return
	}

	// Validates the file type, based on File type
	isCSVFile := strings.HasSuffix(strings.ToLower(fileName), CsvSuffix)
	if !isCSVFile {
		logging.Warning(TAG, "Incorrect filetype uploaded")
		http.Error(w, IncorrectFileTypeCSV, http.StatusInternalServerError)
		return
	}

	logging.Debug(TAG, fmt.Sprintf("File %v has passed type validation", fileName))

	logging.Debug(TAG, "Resetting file read offset")
	// Need to reset the offset after checking for filetype
	if _, err = file.Seek(0, 0); err != nil {
		logging.Error(TAG, err.Error())
		http.Error(w, html.EscapeString(FailedFileRead+fileName), http.StatusInternalServerError)
		return
	}

	logging.Debug(TAG, "Parsing CSV data")
	// Parse the entries from the csv along with the column headers
	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		logging.Error(TAG, err.Error())
		http.Error(w, html.EscapeString(FailedFileRead+fileName), http.StatusInternalServerError)
		return
	} else if len(records) == 0 || len(records[1:]) < 1 {
		logging.Error(TAG, "No records found in uploaded CSV file")
		http.Error(w, NoColumnsFound, http.StatusInternalServerError)
		return
	}

	n := records[:]

	buf := &bytes.Buffer{}
	writer := csv.NewWriter(buf)
	writer.WriteAll(n)

	var badRecordsPositions []string
	for position, record := range records {
		// Do not run the check for the header column
		// Check the first position of the Array which will be the MID
		if position != 0 && record[0] == "" {
			badRecordsPositions = append(badRecordsPositions, strconv.Itoa(position))
		}
	}

	if len(badRecordsPositions) != 0 {
		logging.Error(TAG, "MID is empty in CSV file")
		http.Error(w, fmt.Sprintf(MID_EMPTY_IN_CSV, strings.Join(badRecordsPositions, ", ")), http.StatusInternalServerError)
		return
	}
	dataElementMap := make(map[int]models.UpdateDataElement, 0)
	//header validation
	for position, header := range records[0] {
		header = strings.TrimSpace(header)
		var updateDataElement models.UpdateDataElement
		if position == 0 {
			//some special characters are getting appended in the first row; inorder to remove that bytes.Trim is used
			header = string(bytes.TrimPrefix([]byte(strings.ToLower(header)), common.ByteOrderMark))

			if header != "mid" {
				///validation error
				logging.Error(TAG, "Site validation failed in uploaded CSV file")
				http.Error(w, "Please provide the correct header name as mid", http.StatusInternalServerError)
				return
			}
		} else {
			if header == "" || !strings.Contains(header, ".") {
				logging.Error(TAG, "Data group header validation failed in uploaded CSV file")
				http.Error(w, "Data Group cannot be empty in file header or must be of length 2; Separated by .", http.StatusInternalServerError)
				return
			}
			elementName := strings.Split(header, ".")
			valExp, valMsg, isAllowEmpty, dataType, options, _, err := dal.GetDataElement(elementName[0], elementName[1])
			// handle the case
			if err != nil || (valExp.String == "" && valMsg.String == "" && dataType == "") {
				///validation error
				logging.Error(TAG, "No data elements matched with existing data elements")
				http.Error(w, "Please provide the correct data group or data element name", http.StatusInternalServerError)
				return
			}

			updateDataElement.DataGroupName = elementName[0]
			updateDataElement.DataType = dataType
			updateDataElement.IsAllowEmpty = isAllowEmpty
			updateDataElement.ValExp = valExp
			updateDataElement.ValMsg = valMsg
			updateDataElement.Options = options
			dataElementMap[position] = updateDataElement

		}
	}
	//mapping mID with their data groups
	dataGroupMap := make(map[string][]string)
	header := records[0]
	for _, column := range records[1:] {
		// check merchant exist before updating
		log.Println("MID", column[0])
		mid := column[0]
		midExits, dbResultCode := dal.CheckThatMidExists(mid)
		if !midExits {
			errMsg := fmt.Sprintf("MID doesn't exist, db responded with code %d", dbResultCode)
			logging.Error(TAG, errMsg)
			FailureValSet = append(FailureValSet, bulkUpdateValidationFunc(column[0], header[0], errMsg))
			break
		}

		//go through the columns in each row
		for index, value := range column {
			if value == "" || index == 0 {
				continue
			} else if value == "NULL" {
				//check this element mandatory
				isAllowEmpty := dataElementMap[index].IsAllowEmpty
				if !isAllowEmpty {
					//log the error and continue
					logging.Error(TAG, "mandatory field is NULL")
					FailureValSet = append(FailureValSet, bulkUpdateValidationFunc(column[0], header[index], "Mandatory Field is NULL"))
					continue
				}
			} else { //validate each data elements with their validation expression
				dataElement := dataElementMap[index] //getting header from csv though map

				dataGroupsExists, ok := dataGroupMap[mid]
				if !ok {
					dataGroupsExists, err = dal.FetchSiteDataGroups(mid)
					if err != nil {
						logging.Error(TAG, "error while fetching the data group's from DB")
					}
					dataGroupMap[mid] = dataGroupsExists
				}

				//check the data group is enable and validate values from DB options
				if sliceHelper.SlicesOfStringContains(dataGroupsExists, dataElement.DataGroupName) {
					dataElement := dataElementMap[index] //getting header from csv though map
					switch dataElement.DataType {
					case "STRING":
						if dataElement.ValExp.Valid {
							if validator := validateExp(value, dataElement.ValExp.String); !validator {
								logging.Error(TAG, "validation expression does not match")
								FailureValSet = append(FailureValSet, bulkUpdateValidationFunc(column[0], header[index], dataElement.ValMsg.String))
							}
						}
						if dataElement.Options != "" {
							optionsStr := strings.Split(dataElement.Options, "|")
							if !sliceHelper.SlicesOfStringContains(optionsStr, value) {
								FailureValSet = append(FailureValSet, bulkUpdateValidationFunc(column[0], header[index], "Please provide correct value's or options"))
							}
						}

					case "BOOLEAN":
						if strings.ToLower(value) == "true" || strings.ToLower(value) == "false" {
							logging.Debug("successfully validated the data element type ")
						} else {
							logging.Error(TAG, "could not convert data element value type is mismatched")
							FailureValSet = append(FailureValSet, bulkUpdateValidationFunc(column[0], header[index], "Could not convert data element value as expected type is true or false"))
						}
					case "INTEGER":
						if _, err := strconv.Atoi(value); err == nil {
							if dataElement.ValExp.Valid {
								if validator := validateExp(value, dataElement.ValExp.String); !validator {
									logging.Error(TAG, "validation expression does not match")
									FailureValSet = append(FailureValSet, bulkUpdateValidationFunc(column[0], header[index], dataElement.ValMsg.String))
								}
							}
						} else {
							logging.Error(TAG, "could not convert  data element value type is not Integer")
							FailureValSet = append(FailureValSet, bulkUpdateValidationFunc(column[0], header[index], "Could not convert data element value as expected type is integer"))
						}
					case "LONG":
						if _, err := strconv.ParseInt(value, 10, 64); err == nil {
							if dataElement.ValExp.Valid {
								if validator := validateExp(value, dataElement.ValExp.String); !validator {
									logging.Error(TAG, "validation expression does not match")
									FailureValSet = append(FailureValSet, bulkUpdateValidationFunc(column[0], header[index], dataElement.ValMsg.String))
								}
							}

						} else {
							logging.Error(TAG, "could not convert data element value of type LONG ")
							FailureValSet = append(FailureValSet, bulkUpdateValidationFunc(column[0], header[index], "Could not convert data element value as expected type is LONG"))

						}

					case "JSON":
						var jsonData []interface{}
						if err := json.Unmarshal([]byte(value), &jsonData); err != nil {
							logging.Error(TAG, "could not convert data element value type of JSON")
							FailureValSet = append(FailureValSet, bulkUpdateValidationFunc(column[0], header[index], "Could not convert data element value as expected type is JSON"))
						} else {
							logging.Debug(TAG, "successfully decoded the json data element")
							if dataElement.Options != "" {
								optionsStr := strings.Split(dataElement.Options, "|")
								if !validateOptions(optionsStr, jsonData) {
									FailureValSet = append(FailureValSet, bulkUpdateValidationFunc(column[0], header[index], "Please provide correct value's inside list data"))
								}
							}
						}

					default:
						logging.Error(TAG, fmt.Sprintf("The other %v type of data element %v is not allowed", dataElement.DataType, mid))
						FailureValSet = append(FailureValSet, bulkUpdateValidationFunc(column[0], header[index], fmt.Sprintf("%v type of data element %v is not allowed", dataElement.DataType, mid)))
					}
				} else {
					logging.Error(TAG, fmt.Sprintf("data group is not enable for mid: %v", mid))
					FailureValSet = append(FailureValSet, bulkUpdateValidationFunc(column[0], header[index], fmt.Sprintf("Data group is not enable for mid: %v", mid)))
				}
			}
		}
	}

	fileName = time.Now().Format("20060102150405_") + fileName
	logging.Debug(TAG, "Renamed file to : "+fileName)
	if err := sendFileToFileServer(buf.Bytes(), fileName, BulkSiteUpdateType); err != nil {
		logging.Error(TAG, "Bulk Site Update File Upload Failed : ", err.Error())
		http.Error(w, uploadFileError, http.StatusInternalServerError)
		return
	}

	if len(FailureValSet) > 0 {
		SiteUpdate.UpdateStatus = false
		SiteUpdate.Validations = FailureValSet
		err = dal.InsertBulkApproval(fileName, BulkSiteUpdateType, tmsUser.Username, common.UploadChangeType, -1)
		if err != nil {
			logging.Error(TAG, "Unable to add to bulk approval: ", err.Error())
			http.Error(w, DatabaseTxnError, http.StatusInternalServerError)
			return
		}
	} else {
		SiteUpdate.UpdateStatus = true
		err = dal.InsertBulkApproval(fileName, BulkSiteUpdateType, tmsUser.Username, common.UploadChangeType, 0)
		if err != nil {
			logging.Error(TAG, "Unable to add to bulk approval: ", err.Error())
			http.Error(w, DatabaseTxnError, http.StatusInternalServerError)
			return
		}
	}

	renderPartialTemplate(w, r, "bulkSiteUpdateResults", SiteUpdate, tmsUser)
}

/*
*
Renders the site upload table of the bulk import page
*/
func siteUploadHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	logging.Debug(TAG, "Site bulk upload initiated")

	err := r.ParseMultipartForm(10000)
	if err != nil {
		logging.Error(TAG, err.Error())
		http.Error(w, uploadFileError, http.StatusBadRequest)
		return
	}

	validationDal := new(siteValidatorDal)

	// Extract and validate that the MID is present and in the correct format
	logging.Debug(TAG, "Attempting to extract template MID from http.Request")
	templateMid := r.PostFormValue("mid")
	midValid, siteProfileId, err := validateTemplateMid(validationDal, templateMid)

	if !midValid {
		logging.Error(TAG, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logging.Debug(TAG, fmt.Sprintf("MID %v is of the correct format", templateMid))

	// Validate that a file has been attached
	if len(r.MultipartForm.File) < 1 {
		logging.Warning(TAG, "Site upload initiated without file being present")
		http.Error(w, SiteBulkUpdateMissingFields, http.StatusBadRequest)
		return
	}

	// Extract the file from the request and obtain the name
	logging.Debug(TAG, "Attempting to extract file from http.Request")

	file, handler, err := r.FormFile("file")

	if err != nil {
		logging.Error(TAG, err.Error())
		http.Error(w, uploadFileError, http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileName := handler.Filename

	logging.Debug(TAG, fmt.Sprintf("File: %v has been uploaded", fileName))

	// Validate that the filetype is csv
	// 512 bytes only because DetectContentType (used in ValidateFileType) only reads up to 512 bytes
	buff := make([]byte, 512)
	if _, err = file.Seek(0, 0); err != nil {
		logging.Error(TAG, err.Error())
		http.Error(w, html.EscapeString(FailedFileRead+fileName), http.StatusInternalServerError)
		return
	}

	if _, err = file.Read(buff); err != nil {
		logging.Error(TAG, err.Error())
		http.Error(w, html.EscapeString(FailedFileRead+fileName), http.StatusInternalServerError)
		return
	}

	// Validates the file type, based on File type
	isCSVFile := strings.HasSuffix(strings.ToLower(fileName), CsvSuffix)
	if !isCSVFile {
		logging.Warning(TAG, "Incorrect filetype uploaded")
		http.Error(w, IncorrectFileTypeCSV, http.StatusInternalServerError)
		return
	}

	logging.Debug(TAG, fmt.Sprintf("File %v has passed type validation", fileName))

	logging.Debug(TAG, fmt.Sprintf("MID: %v is valid", templateMid))

	// set the global
	ProfileId = siteProfileId

	// By now we know that the MID is valid and represents a real site, we also know that the file is of csv format

	// Now we need to fetch the template site's data
	templateSite, permitted, err := buildTemplateSiteModel(siteProfileId, tmsUser, w)
	if err != nil {
		logging.Error(TAG, err.Error())
		http.Error(w, ErrorFetchingSiteData, http.StatusInternalServerError)
		return
	} else if !permitted {
		logging.Error(TAG, InsufficientUserPermissions)
		http.Error(w, InsufficientUserPermissions, http.StatusInternalServerError)
		return
	}

	// Storing the template site's data groups here. To be used when saving the new sites and to warn the user
	// if they have included any data elements that will not be saved due said element's group not being active
	DataGroups = parseTemplateSiteDataGroups(templateSite.DataGroups)

	logging.Debug(TAG, "Resetting file read offset")
	// Need to reset the offset after checking for filetype
	if _, err = file.Seek(0, 0); err != nil {
		logging.Error(TAG, err.Error())
		http.Error(w, html.EscapeString(FailedFileRead+fileName), http.StatusInternalServerError)
		return
	}

	logging.Debug(TAG, "Parsing CSV data")
	// Parse the entries from the csv along with the column headers
	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		logging.Error(TAG, err.Error())
		http.Error(w, html.EscapeString(FailedFileRead+fileName), http.StatusInternalServerError)
		return
	} else if len(records) == 0 || len(records[1:]) < 1 {
		logging.Error(TAG, "No records found in uploaded CSV file")
		http.Error(w, NoColumnsFound, http.StatusInternalServerError)
		return
	}

	// Calling the function [ checkUnique ]
	// Only run the code inside the if check if there are non-unique elements
	if failed, sameRecords, uniqueFields := checkUnique(w, records); failed {
		logging.Error(TAG, "Unique fields are not unique in csv")
		http.Error(w, fmt.Sprintf(UNIQUE_FIELDS_EXPECTED, strings.Join(uniqueFields, ", "), strings.Join(sameRecords, ", ")), http.StatusInternalServerError)
		return
	}

	// Format the csv [ TRUE / FALSE ] values to be lowercase
	// Having the values as upper case does not allow them to save
	// Looping through each different record
	for positionY, record := range records {
		// Do not alter the header column
		if positionY != 0 {
			// Looping through each element in the record
			for positionX, recordData := range record {
				// If the element is a boolean type then set it to be lowercase
				lowerCaseData := strings.ToLower(recordData)
				if lowerCaseData == "true" || lowerCaseData == "false" {
					records[positionY][positionX] = lowerCaseData
				}
			}
		}
	}

	var badRecordsPositions []string
	for position, record := range records {
		// Do not run the check for the header column
		// Check the first position of the Array which will be the MID
		if position != 0 && record[0] == "" {
			badRecordsPositions = append(badRecordsPositions, strconv.Itoa(position))
		}
	}

	if len(badRecordsPositions) != 0 {
		logging.Error(TAG, "MID is empty in CSV file")
		http.Error(w, fmt.Sprintf(MID_EMPTY_IN_CSV, strings.Join(badRecordsPositions, ", ")), http.StatusInternalServerError)
		return
	}

	SiteUpload, err = bulkSiteConstructor(validationDal, records, templateSite, siteProfileId)
	if err != nil {
		logging.Error(TAG, err.Error())
		http.Error(w, NoColumnsFound, http.StatusInternalServerError)
		return
	}

	renderPartialTemplate(w, r, "bulkSiteUploadResultsPartial", SiteUpload, tmsUser)
}

func checkUnique(w http.ResponseWriter, records [][]string) (bool, []string, []string) {
	// Keeping the code dynamic by retrieving the unique fields and whether they are allowed to be empty
	uniqueFields, is_allowed_empty, err := dal.GetUniqueFields()
	if err != nil {
		logging.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return true, nil, uniqueFields
	}

	var badRecords []string
	// Looping through the headers of the csv file to find the unique fields
	for headerPositionX, columnName := range records[0] {
		// Formatting the headers to match the database values
		if strings.Contains(columnName, ".") {
			parts := strings.Split(columnName, ".")
			columnName = parts[1]
		}
		// Looping through the unique fields to see if the current selected header is a match
		for uniquePosition, uniqueField := range uniqueFields {
			if columnName == uniqueField {
				// Check the entire columns values to see if they are unique
				if failed, sameRecords := checkColumn(records, headerPositionX, is_allowed_empty[uniquePosition]); failed {
					// Place the returned records from checkColumn on the end of badRecords
					for _, same := range sameRecords {
						badRecords = append(badRecords, same)
					}
				}
			}
		}
	}

	// Removing duplicate bad records from the array for user experience
	badRecords = unique(badRecords)

	if len(badRecords) > 0 {
		return true, badRecords, uniqueFields
	}
	return false, nil, uniqueFields
}

func unique(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func checkColumn(records [][]string, column int, is_allowed_empty bool) (bool, []string) {
	var badRecords []string
	// Looping down the csv column
	for positionY, record := range records {
		// Do not compare the header column
		if positionY == 0 {
			continue
		}

		// The variable column passed into this function is the "x" position of the column
		// Looping through each element in the record
		for columnPositionY, data := range records {
			// Preventing from checking its own column
			// Not checking blank columns
			if positionY == columnPositionY || (is_allowed_empty && data[column] == "") {
				continue
			}
			// If the current element of the csv file matches the selected element then add it to bad records to be returned
			if data[column] == record[column] {
				badRecords = append(badRecords, strconv.Itoa(positionY))
				break
			}
		}
	}

	if len(badRecords) > 0 {
		return true, badRecords
	}
	return false, nil
}

/*
*
Obtains the data groups enables on the template site
*/
func parseTemplateSiteDataGroups(dataGroups []DataGroupModel) map[int]string {
	groupList := make(map[int]string)
	// Iterate through the template site's data groups and store them for the commit process
	for _, dataGroup := range dataGroups {
		if dataGroup.Selected {
			groupList[dataGroup.Group.DataGroupID] = dataGroup.Group.DataGroup
		}
	}
	return groupList
}

type PerformanceTimer struct {
	start time.Time
	label string
}

func (pt *PerformanceTimer) Start(label string) {
	pt.label = label
	pt.start = time.Now()
}

func (pt *PerformanceTimer) Stop() time.Duration {
	duration := time.Since(pt.start)
	fmt.Printf("%s took %vs to execute\n", pt.label, duration.Seconds())
	return duration
}

/*
*
Constructs the new sites from the csv data and template site
*/
func bulkSiteConstructor(validationDal dal.ValidationDal, records [][]string, template ProfileMaintenanceModel, siteProfileId int) (BulkSiteUploadModel, error) {
	var results BulkSiteUploadModel

	// the first entry in records contains column headers. These are written in the following format:
	// dataGroup.dataElement for example store.merchantNo, the only exception to this is the site name, which does not
	// have a data group and so is simply "Site Name"
	// NOTE: these patterns exactly match that of the site export functionality on TMS to make it easier for NI to
	// create bulk upload templates

	// Extract the column headers and convert them into a usable struct
	columns := convertRecordsToColumns(records[0], false)

	// Obtain the data element id for all columns
	columns, err := fetchColumnDataElementIds(columns)
	if err != nil {
		return results, err
	}

	// Convert template into map[int]string where k = element id and v = value
	// This then matches the type used for adding a new site
	templateMap := convertTemplateSiteToMap(template)

	// Convert the data rows into map[int]string to match the format of the template
	newSitesMap, err := convertDataRowsToElements(records[1:], columns, true)
	if err != nil {
		return results, err
	}

	// Use the template and newSitesMap to generate new sites
	newSites := buildNewSites(newSitesMap, templateMap)

	// Validate all of the data elements for the uploaded sites
	results = validateNewSites(validationDal, newSites, columns, siteProfileId)

	return results, nil
}

type siteValidatorDal struct {
	elementsAndMetaData map[int]dal.DataElement
	primaryTidMap       map[int]int
	secondaryTidMap     map[int]int
	snMap               map[string]interface{}
	primaryMidMap       map[string]int
	secondaryMidMap     map[string]int
}

func (s *siteValidatorDal) CheckThatSerialNumberExists(SN string) (snExists bool, err error) {
	if s.snMap == nil {
		db, err := dal.GetDB()
		if err != nil {
			return false, err
		}
		rows, err := db.Query("SELECT t.serial FROM tid t")
		if err != nil {
			return false, err
		}
		defer rows.Close()
		s.snMap = make(map[string]interface{}, 0)
		for rows.Next() {
			var sn string
			if err = rows.Scan(&sn); err != nil {
				return false, err
			}
			s.snMap[sn] = ""
		}
	}
	_, snExists = s.snMap[SN]
	return
}

func (s *siteValidatorDal) GetDataElementMetadata(dataElementId int, profileId int) (dal.DataElement, error) {
	if s.elementsAndMetaData == nil {
		elemes, err := dal.GetAllDataElementsMetadata(profileId)
		if err != nil {
			return dal.DataElement{}, err
		}
		s.elementsAndMetaData = elemes
	}
	return s.elementsAndMetaData[dataElementId], nil
}

func (s *siteValidatorDal) GetDataElementByName(groupName string, elementName string) (int, error) {
	return dal.GetDataElementByName(groupName, elementName)
}

func (s *siteValidatorDal) CheckThatTidExists(TID int) (tidExists bool, resultCode resultCodes.ResultCode, overrideProfileId int) {
	if s.primaryTidMap == nil {
		db, err := dal.GetDB()
		if err != nil {
			return false, resultCodes.DATABASE_CONNECTION_ERROR, -1
		}
		primaryTidRows, err := db.Query(`
			SELECT 
				t.tid_id,
			    ts.tid_profile_id
			FROM tid t
			LEFT JOIN tid_site ts ON
			    t.tid_id = ts.tid_id`)
		if err != nil {
			return false, resultCodes.DATABASE_QUERY_ERROR, -1
		}
		defer primaryTidRows.Close()

		s.primaryTidMap = make(map[int]int, 0)
		for primaryTidRows.Next() {
			var tid int
			var profileId sql.NullInt32
			if err = primaryTidRows.Scan(&tid, &profileId); err != nil {
				return false, resultCodes.DATABASE_QUERY_ERROR, -1
			}
			if profileId.Valid {
				s.primaryTidMap[tid] = int(profileId.Int32)
			} else {
				s.primaryTidMap[tid] = -1
			}
		}

		secondaryTidRows, err := db.Query(`
			SELECT 
				pd.datavalue,
			    pd.profile_id
			FROM data_element de
			INNER JOIN data_group dg ON 
			    de.data_group_id = dg.data_group_id
			INNER JOIN profile_data pd ON
			    pd.data_element_id = de.data_element_id
			WHERE
			    dg.name = 'dualCurrency' AND de.name = 'secondaryTid'`)
		if err != nil {
			return false, resultCodes.DATABASE_QUERY_ERROR, -1
		}
		defer secondaryTidRows.Close()

		s.secondaryTidMap = make(map[int]int, 0)
		for secondaryTidRows.Next() {
			var tid string
			var profileId sql.NullInt32
			if err = secondaryTidRows.Scan(&tid, &profileId); err != nil {
				return false, resultCodes.DATABASE_QUERY_ERROR, -1
			}
			if tid != "" {
				tidInt, err := strconv.Atoi(tid)
				if err != nil {
					return false, resultCodes.DATABASE_QUERY_ERROR, -1
				}

				if profileId.Valid {
					s.secondaryTidMap[tidInt] = int(profileId.Int32)
				} else {
					s.secondaryTidMap[tidInt] = -1
				}
			}
		}
	}
	if profileId, present := s.primaryTidMap[TID]; present {
		return true, resultCodes.TID_NOT_UNIQUE_PRIMARY_TID_DUPLICATE, profileId
	} else if profileId, present := s.secondaryTidMap[TID]; present {
		return true, resultCodes.TID_NOT_UNIQUE_SECONDARY_TID_DUPLICATE, profileId
	} else {
		return false, resultCodes.TID_DOES_NOT_EXIST, -1
	}
}

func (s *siteValidatorDal) CheckThatMidExists(MID string) (midExists bool, resultCode resultCodes.ResultCode, profileId int) {
	if s.primaryMidMap == nil {
		db, err := dal.GetDB()
		if err != nil {
			return false, resultCodes.DATABASE_CONNECTION_ERROR, -1
		}
		primaryMidRows, err := db.Query(`
			SELECT 
				pd.datavalue, pd.profile_id
			FROM data_element de
			INNER JOIN data_group dg ON 
			    de.data_group_id = dg.data_group_id
			INNER JOIN profile_data pd ON
			    pd.data_element_id = de.data_element_id
			WHERE
				dg.name = 'store' AND de.name = 'merchantNo'`)
		if err != nil {
			return false, resultCodes.DATABASE_QUERY_ERROR, -1
		}
		defer primaryMidRows.Close()
		s.primaryMidMap = make(map[string]int, 0)
		for primaryMidRows.Next() {
			var mid string
			var profileId int
			if err = primaryMidRows.Scan(&mid, &profileId); err != nil {
				return false, resultCodes.DATABASE_QUERY_ERROR, -1
			}
			s.primaryMidMap[mid] = profileId
		}

		secondaryMidRows, err := db.Query(`
			SELECT 
				pd.datavalue, pd.profile_id
			FROM data_element de
			INNER JOIN data_group dg ON 
			    de.data_group_id = dg.data_group_id
			INNER JOIN profile_data pd ON
			    pd.data_element_id = de.data_element_id
			WHERE
			    dg.name = 'dualCurrency' AND de.name = 'secondaryMid'`)
		if err != nil {
			return false, resultCodes.DATABASE_QUERY_ERROR, -1
		}
		defer secondaryMidRows.Close()
		s.secondaryMidMap = make(map[string]int, 0)
		for secondaryMidRows.Next() {
			var mid string
			var profileId int
			if err = secondaryMidRows.Scan(&mid, &profileId); err != nil {
				return false, resultCodes.DATABASE_QUERY_ERROR, -1
			}
			s.secondaryMidMap[mid] = profileId
		}
	}

	if profileId, present := s.primaryMidMap[MID]; present {
		return true, resultCodes.MID_NOT_UNIQUE_PRIMARY_MID_DUPLICATE, profileId
	} else if profileId, present := s.secondaryMidMap[MID]; present {
		return true, resultCodes.MID_NOT_UNIQUE_SECONDARY_MID_DUPLICATE, profileId
	} else {
		return false, resultCodes.MID_DOES_NOT_EXIST, -1
	}
}

func (s *siteValidatorDal) GetIsUnique(elementId int, elementValue string, profile int) (bool, error) {
	return dal.GetIsUnique(elementId, elementValue, profile)
}

/*
*
Validates the data elements for each new site
*/
func validateNewSites(validationDal dal.ValidationDal, newSites NewSitesElements, columns []DataColumn, templateSiteProfileId int) BulkSiteUploadModel {
	var validationResults BulkSiteUploadModel
	var validationPasses NewSitesElements
	failed := false

	// Iterate through each new site and validate its data elements
	for _, site := range newSites.getNewSites() {
		validationErrors, failedIndex := validateDataElements(validationDal, site.getDataElementsAsMap(), templateSiteProfileId, Site, true, "")

		if len(validationErrors) > 0 {
			logging.Error(TAG, fmt.Sprintf("Validation of data element %v for new site entry number %v has failed", failedIndex, site.getRef()))
			var validationFailure FailedProfileValidation

			validationFailure.setFailedElementId(failedIndex)
			validationFailure.setFailureReason(validationErrors[0]) // Will always be in position 0 as the first validation failure is returned
			// Find the data element name for the failed validation
			failedElementName := ""
			for _, column := range columns {
				if column.getElementId() == failedIndex {
					failedElementName = column.getName()
				}
			}
			validationFailure.setFailedElementName(failedElementName)
			validationFailure.setSite(site)

			failed = true
			validationResults.setFailure(validationFailure)
			// As soon as we hit a validation failure we want to cease further validations
			break
		}
		validationPasses.NewSites = append(validationPasses.getNewSites(), site)
	}
	logging.Debug(TAG, "All sites successfully validated")

	validationResults.setFailed(failed)
	validationResults.setPasses(validationPasses)

	// Add the columns to the pageModel
	for _, column := range columns {
		validationResults.addColumn(column)
	}

	// Prevents any sillyness with Go not iterating maps in order
	validationResults = sortResults(validationResults)

	return validationResults
}

/*
*
Sorts the data elements of the validated sites and the columns by data element id.
This enables the validation table to have the correct data in the correct columns on TMS
*/
func sortResults(results BulkSiteUploadModel) BulkSiteUploadModel {

	// trim the columns that relate to unused data groups
	results = trimUnusedDataGroupColumns(results)

	// sort the data elements of the new sites by data element id
	for _, site := range results.Passes.getNewSites() {
		sort.Slice(site.DataElements, func(i, j int) bool { return site.DataElements[i].ElementId < site.DataElements[j].ElementId })
	}

	// Sort the columns by data element id
	sort.Slice(results.Columns, func(i, j int) bool { return results.Columns[i].ElementId < results.Columns[j].ElementId })

	return results
}

/*
*
Strips out any data columns not used by the template site
*/
func trimUnusedDataGroupColumns(results BulkSiteUploadModel) BulkSiteUploadModel {
	inUseElements := make(map[int]bool)

	// First lets find all of the elements in use on all sites
	passes := results.getPasses()
	for _, pass := range passes.NewSites {
		for _, element := range pass.DataElements {
			inUseElements[element.ElementId] = true
		}
	}

	var inUseColumns []DataColumn
	// Now strip the columns out that are not in use
	for _, column := range results.Columns {
		if _, ok := inUseElements[column.getElementId()]; ok {
			// Columns that are used
			inUseColumns = append(inUseColumns, column)
		}
	}

	results.UnusedColumns = listIgnoredColumns(results.Columns)
	if len(results.UnusedColumns) > 0 {
		results.setColumnsRemoved(true)
	}

	results.Columns = inUseColumns

	return results
}

func listIgnoredColumns(columns []DataColumn) []DataColumn {
	var unusedColumns []DataColumn

	// iterate through the original columns
	for i, column := range columns {
		enabled := false
		fmt.Sprintf("Item number: %d", i)
		// for each column we need to check if it's parent data group is enabled
		for _, group := range DataGroups {
			if column.getDataGroup() == group {
				enabled = true
			}
		}
		if !enabled {
			unusedColumns = append(unusedColumns, column)
		}
	}

	return unusedColumns
}

/*
*
Overlays the new site data ontop of the template site
*/
func buildNewSites(newSitesMap NewSitesElements, templateMap map[int]DataElement) NewSitesElements {
	var newSitesElements NewSitesElements
	// For each of the new sites to be built, copy the template and then overwrite where data has been entered
	for _, row := range newSitesMap.getNewSites() {
		var newSite NewProfileEntry
		// make a copy of the template
		var templateCopy []DataElement
		for k := range templateMap {
			var element DataElement
			element.ElementId = k
			element.Data = templateMap[k].Data
			element.Encrypted = templateMap[k].Encrypted
			element.Unique = templateMap[k].Unique
			element.Overriden = templateMap[k].Overriden
			templateCopy = append(templateCopy, element)
		}

		// Iterate through the elements in the current data row
		for _, tempElement := range row.getDataElements() {
			if tempElement.Data != "" && len(tempElement.Data) != 0 {
				// find the element in the template and overwrite it
				for i, temp := range templateCopy {
					if temp.ElementId == tempElement.ElementId {
						temp.Data = tempElement.Data
						templateCopy[i] = temp
					}
				}
			}
		}
		newSite.setRef(row.getRef())
		newSite.DataElements = templateCopy
		newSite.setMid(row.getMid())
		newSite.setSiteName(row.getSiteName())
		newSitesElements.NewSites = append(newSitesElements.getNewSites(), newSite)
	}

	return newSitesElements
}

/*
*
Takes the records supplied in csv format and converts them into structures
*/
func convertDataRowsToElements(records [][]string, columns []DataColumn, complainAboutMissingElements bool) (NewSitesElements, error) {
	var newSitesMap NewSitesElements

	logging.Debug(TAG, "Attempting to convert data rows to element map")
	// Iterate through each of the rows of data
	for i, row := range records {
		// For each row of data create a new map[int]string
		var siteMap []DataElement
		// Iterate through all of the data in a row and obtain the position and value
		for j, entry := range row {
			// Use the position of the data to obtain the dataElementId from DataColumn
			col, err := getDataColumnFromPosition(columns, j)
			if err != nil {
				logging.Warning(err.Error())
			}

			if col == nil && complainAboutMissingElements {
				return newSitesMap, errors.New("invalid element position")
			}
			var element DataElement
			element.ElementId = col.ElementId
			element.Unique = col.IsUnique
			element.Encrypted = col.IsEncrypted
			element.GroupName = col.DataGroup
			element.Name = col.Name
			element.Data = entry
			siteMap = append(siteMap, element)
		}
		// Create a NewProfileEntry
		var newSite NewProfileEntry
		newSite.setRef(i)
		newSite.DataElements = siteMap
		newSite.setSiteName(getDataFromElementName(columns, newSite.getDataElementsAsMap(), "name"))
		newSite.setMid(getDataFromElementName(columns, newSite.getDataElementsAsMap(), "merchantNo"))
		newSitesMap.NewSites = append(newSitesMap.getNewSites(), newSite)
	}
	logging.Debug(TAG, "Successfully converted data rows site to element map")
	return newSitesMap, nil
}

/*
*
Fetch the data element ID by position in the data column
*/
func getElementIdFromPosition(columns []DataColumn, position int) int {
	if col, err := getDataColumnFromPosition(columns, position); err == nil {
		return col.ElementId
	} else {
		logging.Warning(err)
		return -1
	}
}

func getDataColumnFromPosition(columns []DataColumn, position int) (*DataColumn, error) {
	for _, column := range columns {
		if column.getPosition() == position {
			return &column, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("No column header found at position %d", position))
}

/*
*
Converts the template site profile maintenance model into a struct common to the new data
*/
func convertTemplateSiteToMap(templateSite ProfileMaintenanceModel) map[int]DataElement {
	templateMap := make(map[int]DataElement)

	logging.Debug(TAG, "Attempting to convert template site to map of data elements")

	// Add site data elements to map
	for _, dataGroup := range templateSite.ProfileGroups {
		for _, templateElement := range dataGroup.DataElements {
			var dataElement DataElement
			dataElement.ElementId = templateElement.ElementId
			dataElement.Data = templateElement.DataValue
			dataElement.Encrypted = templateElement.IsEncrypted
			dataElement.Unique = templateElement.Unique
			dataElement.Overriden = templateElement.Overriden
			templateMap[dataElement.ElementId] = dataElement
		}
	}
	logging.Debug(TAG, "Successfully converted template site to element map")

	return templateMap
}

/*
*
Retrieves the data element ids for each column in the new site csv. If
*/
func fetchColumnDataElementIds(columns []DataColumn) ([]DataColumn, error) {
	var columnsWithId []DataColumn

	metadata, err := dal.GetAllDataElementsMetadata(-1)
	if err != nil {
		return nil, err
	}
	logging.Debug(TAG, "Retrieving data element ids for all uploaded columns")
	// Iterate through all of the column headers and use the data group and element name
	// to obtain the data element id
	for _, column := range columns {
		elementId, err := dal.GetDataElementByName(column.getDataGroup(), column.getName())
		if err != nil {
			logging.Warning(TAG, fmt.Sprintf("Error seen when attempting to obtain element id for %v of the %v data group", column.getName(), column.getDataGroup()))
			return columnsWithId, err
		}
		displayName := metadata[elementId].DisplayName
		if metadata[elementId].DisplayName == "" || len(metadata[elementId].DisplayName) == 0 {
			displayName = column.getName()
		}

		if column.getElementId() >= 0 {
			column.setElementId(elementId)
		}

		column.setEncrypted(metadata[elementId].IsEncrypted)
		column.setUnique(metadata[elementId].Unique)
		column.setDisplayName(displayName)
		columnsWithId = append(columnsWithId, column)
	}

	logging.Debug(TAG, "Successfully retrieved data element ids for all uploaded columns")
	return columnsWithId, nil
}

/*
*
Converts the array of column headers into a usable struct
*/
func convertRecordsToColumns(records []string, addUnknownColumnsAndIgnore bool) []DataColumn {
	var columns []DataColumn

	// Iterate through the records and convert them one by one into DataColumn
	for i, record := range records {
		var column DataColumn

		switch strings.ToLower(record) {
		case TID:
			column.setName(record)
			column.setPosition(i)
			column.setElementId(TID_ELEMENT_ID)
			columns = append(columns, column)
		case Serial:
			column.setName(record)
			column.setPosition(i)
			column.setElementId(SERIAL_NUMBER_ID)
			columns = append(columns, column)
		default:
			if strings.Contains(record, ".") {
				// For all other data elements they are stored in the format dataGroup.dataElement so we need to split them
				splitRecord := strings.Split(record, ".")

				column.setDataGroup(splitRecord[0])
				column.setName(splitRecord[1])
				column.setPosition(i)

				columns = append(columns, column)
			} else if addUnknownColumnsAndIgnore {
				// Putting the data into a pseudo group so that unknown headers are ignored rather than throwing an error.
				// Also all the headers have to be present.
				logging.Debug(TAG, "putting into unknown_do_not_commit."+record)

				column.setName(record)
				column.setPosition(i)
				column.setElementId(-1)
				column.setIgnore(true)

				columns = append(columns, column)
			}
		}
	}
	return columns
}

/*
*
Obtains the template site model
*/
func buildTemplateSiteModel(profileId int, tmsUser *entities.TMSUser, w http.ResponseWriter) (ProfileMaintenanceModel, bool, error) {
	var templateSite ProfileMaintenanceModel

	// Use int profileId to retrieve the siteId
	logging.Debug(TAG, "Attempting to fetch site ID from profile")
	siteId, err := dal.GetSiteFromProfile(profileId)
	if err != nil {
		logging.Warning(TAG, "Failed to fetch siteId from profile")
		return templateSite, false, err
	}

	// Ensure the logged in user has permission to access site data
	permitted, err := checkUserAcquirePermsBySite(tmsUser, siteId)
	if err != nil {
		logging.Warning(TAG, "Unable to ascertain if user has permission to access site data")
		return templateSite, false, err
	} else if !permitted {
		logging.Warning(TAG, "User does not have permission to access site data")
		return templateSite, false, nil
	}

	// Return the site model
	return buildProfileMaintenanceModel(w, Site, profileId, tmsUser, 0, 1, "", siteId), permitted, nil
}

/*
Compares the filetype of the uploaded file with that of the passed in filetype
*/
func validateFileType(buff []byte, desiredFiletype string) bool {
	filetype := http.DetectContentType(buff)
	if strings.Contains(filetype, desiredFiletype) {
		return true
	}
	return false
}

/*
Ensures that the supplied string is not empty and is all numeric characters
*/
func validateTemplateMid(validationDal dal.ValidationDal, mid string) (bool, int, error) {
	logging.Debug(TAG, "Beginning MID validation")

	// Ensure that the MID has been entered
	if len(mid) == 0 || mid == "" {
		err := "No MID entered for bulk upload"
		logging.Warning(TAG, err)
		return false, -1, errors.New(err)
	}

	// Ensure that the MID is numerical
	if !isMidNumerical(mid) {
		err := "The entered MID contains non-numeric characters"
		logging.Warning(TAG, err)
		return false, -1, errors.New(err)
	}

	midExists, resultCode, profileId := validationDal.CheckThatMidExists(mid)
	if !midExists || resultCode != resultCodes.MID_NOT_UNIQUE_PRIMARY_MID_DUPLICATE {
		var err string
		switch resultCode {
		case resultCodes.MID_NOT_UNIQUE_SECONDARY_MID_DUPLICATE:
			err = "The entered MID is a secondary MID, a primary MID must be supplied."
		default:
			err = resultCodes.GetErrorMsgByCode(resultCode)
		}
		logging.Debug(TAG, err)
		return false, -1, errors.New(err)
	}

	logging.Debug(TAG, "MID validation completed")
	return true, profileId, nil
}

/*
*
Checks that the supplied string is comprised entirely of numeric characters
*/
func isMidNumerical(mid string) bool {
	if _, err := strconv.Atoi(mid); err == nil {
		return true
	} else {
		return false
	}
}

func insertAllImportedSites(newSites NewSitesElements, chainProfileId, acquirerProfileId int, dataGroups []string, user *entities.TMSUser) error {
	db, err := dal.GetDB()
	if err != nil {
		return errors.New(DatabaseAccessError)
	}

	logging.Debug(fmt.Sprintf("Starting database transaction to import '%d' new sites", len(newSites.getNewSites())))
	txn, err := db.Begin()
	if err != nil {
		logging.Error(fmt.Sprintf("An error occured starting the bulk import database transaction; '%s'", err.Error()))
		return errors.New(DatabaseAccessError)
	}

	profileInsertArgs := make([]interface{}, len(newSites.getNewSites()))
	for i, name := range newSites.getNewSites() {
		profileInsertArgs[i] = name.SiteName
	}

	var profileInsertLines []string
	var siteCreateLines []string
	for range newSites.getNewSites() {
		profileInsertLines = append(profileInsertLines, "(4, ?, 1, current_timestamp, 'system', current_timestamp, 'system')")
		siteCreateLines = append(siteCreateLines, "(1, current_timestamp, 'system', current_timestamp, 'system')")
	}
	logging.Debug(fmt.Sprintf("Inserting '%d' rows into the profile table", len(profileInsertLines)))
	r, err := txn.Exec(fmt.Sprintf(`
		INSERT INTO profile (profile_type_id, name, version, updated_at, updated_by, created_at, created_by)
		VALUES
		%s;`, strings.Join(profileInsertLines, ", ")), profileInsertArgs...)
	if err != nil {
		logging.Error(fmt.Sprintf("An error occured inserting the bulk import sites into the profile table, rolling back txn; '%s'", err.Error()))
		txn.Rollback()
		return errors.New(DatabaseTxnError)
	}

	// This simply gets all the profileIds that have been generated from the insert
	var newProfiles []int64
	if firstNewProfileId, err := r.LastInsertId(); err == nil {
		if newRows, err := r.RowsAffected(); err == nil {
			newProfiles = make([]int64, newRows)
			for i := 0; i < len(newProfiles); i++ {
				newProfiles[i] = firstNewProfileId + int64(i)
			}
		} else {
			logging.Error(fmt.Sprintf("An error occured getting the rows affected by the insert into the profile table; '%s'", err.Error()))
			txn.Rollback()
			return errors.New(genericServerError)
		}
	} else {
		logging.Error(fmt.Sprintf("An error occured getting the last ID inserted into the profile table; '%s'", err.Error()))
		txn.Rollback()
		return errors.New(genericServerError)
	}

	logging.Debug(fmt.Sprintf("Inserting '%d' rows into the site table", len(siteCreateLines)))
	r, err = txn.Exec(fmt.Sprintf(`
		INSERT INTO site (version, updated_at, updated_by, created_at, created_by)
		VALUES
		%s;`, strings.Join(siteCreateLines, ",")))
	if err != nil {
		logging.Error(fmt.Sprintf("An error occured inserting the bulk import sites into the site table, rolling back txn; '%s'", err.Error()))
		txn.Rollback()
		return errors.New(DatabaseTxnError)
	}

	// This simply gets all the site profileIds that have been generated from the insert
	var newSiteProfileIds []int64
	if firstNewProfileId, err := r.LastInsertId(); err == nil {
		if newRows, err := r.RowsAffected(); err == nil {
			newSiteProfileIds = make([]int64, newRows)
			for i := 0; i < len(newSiteProfileIds); i++ {
				newSiteProfileIds[i] = firstNewProfileId + int64(i)
			}
		} else {
			logging.Error(fmt.Sprintf("An error occured getting the rows affected by the insert into the site table; '%s'", err.Error()))
			txn.Rollback()
			return errors.New(genericServerError)
		}
	} else {
		logging.Error(fmt.Sprintf("An error occured getting the last ID inserted into the site table; '%s'", err.Error()))
		txn.Rollback()
		return errors.New(genericServerError)
	}

	logging.Debug(fmt.Sprintf("Inserting '%d' rows into the site_profiles table", len(newProfiles)*4))
	for i, newProfileId := range newProfiles {
		_, err := txn.Exec(`
			INSERT INTO site_profiles (site_id, profile_id, version, updated_at, updated_by, created_at, created_by) 
			VALUES
			(?, ?, 1, current_timestamp, 'system', current_timestamp, 'system'),
			(?, ?, 1, current_timestamp, 'system', current_timestamp, 'system'),
			(?, ?, 1, current_timestamp, 'system', current_timestamp, 'system'),
			(?, 1, 1, current_timestamp, 'system', current_timestamp, 'system');`,
			newSiteProfileIds[i], newProfileId, newSiteProfileIds[i], chainProfileId, newSiteProfileIds[i], acquirerProfileId, newSiteProfileIds[i])
		if err != nil {
			logging.Error(fmt.Sprintf("An error occured inserting the bulk import sites into the site_profiles table, rolling back txn; '%s'", err.Error()))
			txn.Rollback()
			return errors.New(DatabaseTxnError)
		}
	}

	logging.Debug("Preparing the data_element insert statement")
	insertDataElementsStatement, err := txn.Prepare(`
		INSERT INTO profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by,created_at, created_by, approved, overriden, is_encrypted)
		VALUES
		(?, ?, ?, 1, current_timestamp, 'system', current_timestamp, 'system', 1, ?, ?);`)
	if err != nil {
		logging.Error(fmt.Sprintf("An error occurred preparing the data_element insert statement; '%s'", err.Error()))
		txn.Rollback()
		return errors.New(DatabaseTxnError)
	}

	logging.Debug("Fetching the acquirer name")
	var acquirerName string
	if err := txn.QueryRow("SELECT name FROM profile WHERE profile_id = ?", acquirerProfileId).Scan(&acquirerName); err != nil {
		logging.Error(fmt.Sprintf("An error occurred fetching the acquirer name; '%s'", err.Error()))
		txn.Rollback()
		insertDataElementsStatement.Close()
		return errors.New(DatabaseTxnError)
	}

	logging.Debug("Preparing the site created history insert statement")
	writeSiteCreatedHistoryStatement, err := txn.Prepare(`
		INSERT INTO approvals(profile_id, data_element_id, change_type, current_value, new_value, created_at, approved_at, created_by, approved_by, approved, tid_id, merchant_id, acquirer)
    	VALUES
    	(?, 1, 5, '', 'Site Created', NOW(), NOW(), ?, ?, 1, null, ?, ?);`)
	if err != nil {
		logging.Error(fmt.Sprintf("An error occurred preparing the site created history insert statement; '%s'", err.Error()))
		txn.Rollback()
		insertDataElementsStatement.Close()
		return errors.New(DatabaseTxnError)
	}

	logging.Debug("Preparing the profile data group insert statement")
	writeDataGroupStatement, err := txn.Prepare(`
		insert into profile_data_group(profile_id, data_group_id, version, updated_at, updated_by, created_at, created_by) 
		VALUES 
		(?, ?, 1, current_timestamp, 'system', current_timestamp, 'system');`)
	if err != nil {
		logging.Error(fmt.Sprintf("An error occurred preparing the profile data group insert statement insert statement; '%s'", err.Error()))
		txn.Rollback()
		insertDataElementsStatement.Close()
		writeSiteCreatedHistoryStatement.Close()
		return errors.New(DatabaseTxnError)
	}

	// The sites are inserted in the same order as they are in within newSite so it is safe to simply get the site by
	// index.
	logging.Debug(fmt.Sprintf("Inserting data elements, site created history, and data groups for '%d' sites", len(newProfiles)))
	for i, profileId := range newProfiles {
		site := newSites.getNewSites()[i]
		for _, elem := range site.DataElements {
			data := elem.Data
			if elem.Encrypted {
				data = crypt.Encrypt(data)
			}

			_, err := insertDataElementsStatement.Exec(profileId, elem.ElementId, data, elem.Overriden, elem.Encrypted)
			if err != nil {
				logging.Error(fmt.Sprintf("An error occurred inserting the data elements, rolling back database txn; '%s'", err.Error()))
				txn.Rollback()
				insertDataElementsStatement.Close()
				writeDataGroupStatement.Close()
				writeSiteCreatedHistoryStatement.Close()
				return errors.New(DatabaseTxnError)
			}
		}
		_, err := writeSiteCreatedHistoryStatement.Exec(profileId, user.Username, user.Username, site.getMid(), acquirerName)
		if err != nil {
			logging.Error(fmt.Sprintf("An error occurred inserting the site created history data, rolling back database txn; '%s'", err.Error()))
			txn.Rollback()
			insertDataElementsStatement.Close()
			writeDataGroupStatement.Close()
			writeSiteCreatedHistoryStatement.Close()
			return errors.New(DatabaseTxnError)
		}

		for _, dg := range dataGroups {
			_, err = writeDataGroupStatement.Exec(profileId, dg)
			if err != nil {
				logging.Error(fmt.Sprintf("An error occurred inserting site data group data, rolling back database txn; '%s'", err.Error()))
				txn.Rollback()
				writeDataGroupStatement.Close()
				writeSiteCreatedHistoryStatement.Close()
				insertDataElementsStatement.Close()
				return errors.New(DatabaseTxnError)
			}
		}
	}

	logging.Debug("Bulk import complete, committing transaction to database")
	if err := txn.Commit(); err != nil {
		logging.Error(fmt.Sprintf("An error committing the transaction to the database, rolling back changes; '%s'", err.Error()))
		txn.Rollback()
		insertDataElementsStatement.Close()
		writeSiteCreatedHistoryStatement.Close()
		writeDataGroupStatement.Close()
		return errors.New(DatabaseTxnError)
	}
	logging.Debug("Bulk import completed")
	insertDataElementsStatement.Close()
	writeSiteCreatedHistoryStatement.Close()
	writeDataGroupStatement.Close()
	return nil
}

func bulkDelete(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {

	logging.Debug(TAG, "Bulk delete initiated")

	var FailureValSet []models.ValidationSet
	var BulkDelete models.BulkUpdateVal

	err := r.ParseMultipartForm(10000)
	if err != nil {
		logging.Error(TAG, err.Error())
		http.Error(w, uploadFileError, http.StatusBadRequest)
		return
	}
	deleteType := r.FormValue("DeleteType")

	// Validate that a file has been attached
	if len(r.MultipartForm.File) < 1 {
		switch deleteType {
		case "SiteDelete":
			logging.Warning(TAG, BulkSiteDeleteMissingFile)
			http.Error(w, BulkSiteDeleteMissingFile, http.StatusBadRequest)
		case "TidDelete":
			logging.Warning(TAG, BulkTidDeleteMissingFile)
			http.Error(w, BulkTidDeleteMissingFile, http.StatusBadRequest)
		}
		return
	}

	// Extract the file from the request and obtain the name
	logging.Debug(TAG, "Attempting to extract file from http.Request")

	file, handler, err := r.FormFile("file")
	if err != nil {
		logging.Error(TAG, err.Error())
		http.Error(w, uploadFileError, http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileName := handler.Filename

	logging.Debug(TAG, fmt.Sprintf("File: %v has been uploaded", fileName))

	// Validate that the filetype is csv
	// 512 bytes only because DetectContentType (used in ValidateFileType) only reads up to 512 bytes
	buff := make([]byte, 512)
	if _, err = file.Seek(0, 0); err != nil {
		logging.Error(TAG, err.Error())
		http.Error(w, html.EscapeString(FailedFileRead+fileName), http.StatusInternalServerError)
		return
	}

	if _, err = file.Read(buff); err != nil {
		logging.Error(TAG, err.Error())
		http.Error(w, html.EscapeString(FailedFileRead+fileName), http.StatusInternalServerError)
		return
	}

	isCSVFile := strings.HasSuffix(strings.ToLower(fileName), CsvSuffix)
	// Validates the file type, based on File type
	if !isCSVFile {
		logging.Warning(TAG, "Incorrect filetype uploaded")
		http.Error(w, IncorrectFileTypeCSV, http.StatusInternalServerError)
		return
	}

	logging.Debug(TAG, fmt.Sprintf("File %v has passed type validation", fileName))

	logging.Debug(TAG, "Resetting file read offset")
	// Need to reset the offset after checking for filetype
	if _, err = file.Seek(0, 0); err != nil {
		logging.Error(TAG, err.Error())
		http.Error(w, html.EscapeString(FailedFileRead+fileName), http.StatusInternalServerError)
		return
	}

	logging.Debug(TAG, "Parsing CSV data")
	// Parse the entries from the csv along with the column
	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		logging.Error(TAG, err.Error())
		http.Error(w, html.EscapeString(FailedFileRead+fileName), http.StatusInternalServerError)
		return
	} else if len(records) == 0 || len(records[1:]) < 1 {
		logging.Error(TAG, "No records found in uploaded CSV file")
		http.Error(w, NoColumnsFound, http.StatusInternalServerError)
		return
	}

	n := records[:]

	buf := &bytes.Buffer{}
	writer := csv.NewWriter(buf)
	writer.WriteAll(n)

	for col, val := range records {
		value := strings.TrimSpace(val[0])
		if col == 0 {
			//CHECK HEADER
			value = string(bytes.TrimPrefix([]byte(strings.ToLower(value)), common.ByteOrderMark))
			if value != "tid" && value != "site name" {
				///validation error
				logging.Error(TAG, "Header validation failed in uploaded CSV file")
				http.Error(w, "Please provide the correct header (TID/Site Name)", http.StatusInternalServerError)
				return
			}
		} else {
			if deleteType == "TidDelete" {
				//Validate TID
				tidInt, err := strconv.Atoi(value)
				if err != nil {
					logging.Error(TAG, "Error while Tid conversion ")
					FailureValSet = append(FailureValSet, bulkUpdateValidationFunc(strconv.Itoa(col), value, "Error while Tid conversion"))
					continue
				}
				tidExits, _, _ := dal.CheckThatTidExists(tidInt)
				if !tidExits {
					logging.Error(TAG, "Tid does not exists")
					FailureValSet = append(FailureValSet, bulkUpdateValidationFunc(strconv.Itoa(col), value, "Tid does not exist"))
					continue
				}
			} else {
				//Validate Site
				siteResults, err := dal.GetSiteList(value, tmsUser)
				if err != nil || siteResults == nil {
					if err == nil {
						err = errors.New("Site Name does not exist")
					}
					logging.Error(TAG, err.Error())
					FailureValSet = append(FailureValSet, bulkUpdateValidationFunc(strconv.Itoa(col), value, err.Error()))
					continue
				}

				_, err = getAccurateResult(siteResults, value)
				if err != nil {
					logging.Error(TAG, err.Error())
					FailureValSet = append(FailureValSet, bulkUpdateValidationFunc(strconv.Itoa(col), value, err.Error()))
					continue
				}
			}

		}
	}

	if len(FailureValSet) > 0 {
		BulkDelete.UpdateStatus = false
		BulkDelete.Validations = FailureValSet
	} else {
		BulkDelete.UpdateStatus = true
	}

	if BulkDelete.UpdateStatus {

		fileName = time.Now().Format("20060102150405_") + fileName
		logging.Debug(TAG, "Renamed file to : "+fileName)
		if deleteType == "TidDelete" {
			if err := sendFileToFileServer(buf.Bytes(), fileName, BulkTidDeleteType); err != nil {
				logging.Error(TAG, "Bulk Tid Update File Upload Failed : ", err.Error())
				http.Error(w, uploadFileError, http.StatusInternalServerError)
				return
			}

			err = dal.InsertBulkApproval(fileName, BulkTidDeleteType, tmsUser.Username, common.UploadChangeType, 0)
			if err != nil {
				http.Error(w, DatabaseTxnError, http.StatusInternalServerError)
				return
			}
		} else {
			if err := sendFileToFileServer(buf.Bytes(), fileName, BulkSiteDeleteType); err != nil {
				logging.Error(TAG, "Bulk Tid Update File Upload Failed : ", err.Error())
				http.Error(w, uploadFileError, http.StatusInternalServerError)
				return
			}

			err = dal.InsertBulkApproval(fileName, BulkSiteDeleteType, tmsUser.Username, common.UploadChangeType, 0)
			if err != nil {
				http.Error(w, DatabaseTxnError, http.StatusInternalServerError)
				return
			}
		}

	}
	if deleteType == "TidDelete" {
		renderPartialTemplate(w, r, "bulkTidDeleteResults", BulkDelete, tmsUser)
	} else {
		renderPartialTemplate(w, r, "bulkSiteDeleteResults", BulkDelete, tmsUser)
	}

}

func bulkPaymentServiceImportHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	var BulkUploadConfig models.BulkUpdateVal
	var FailureValSet []models.ValidationSet
	file, header, inErr := getUploadedFile(r)
	if inErr != nil {
		respondWithError(w, inErr)
		return
	}
	defer file.Close()
	logging.Information(TAG, fmt.Sprintf("payment services import file %s uploaded", header.Filename))

	fileName := header.Filename

	inErr = verifyFileIsCsv(file, header.Filename)
	if inErr != nil {
		respondWithError(w, inErr)
		return
	}

	records, inErr := readAllCsvRecords(file, header.Filename)
	if inErr != nil {
		respondWithError(w, inErr)
		return
	}

	n := records[:]
	buf := &bytes.Buffer{}
	writer := csv.NewWriter(buf)
	writer.WriteAll(n)

	_, groupIdx, serviceIdx, _, _, _, _ := getPaymentServiceColumnIndexes(records[0])
	if groupIdx == -1 || serviceIdx == -1 {
		respondWithError(w, &bulkImportError{
			statusCode: http.StatusBadRequest,
			respErrMsg: "the uploaded file does not contain the correct columns",
			logErrMsg:  "the uploaded file does not contain the correct columns",
		})
		return
	}

	model := &entities.PaymentServiceGroupImportModel{
		GroupsCreated:   0,
		ServicesCreated: 0,
		FailedRows:      0,
		Groups:          map[string]entities.PaymentServiceImportGroup{},
	}

	for index, record := range records[1:] {
		groupName := record[groupIdx]
		serviceName := record[serviceIdx]

		if groupName == "" {
			FailureValSet = append(FailureValSet, bulkUpdateValidationFunc(strconv.Itoa(index), "Service_Group", "Service Group cannot be empty"))
			model.FailedRows++
			continue
		}

		if _, ok := model.Groups[groupName]; !ok {
			groupId := dal.GroupIdFromName(groupName)
			if groupId == -1 {
				model.GroupsCreated++
			}

			model.Groups[groupName] = entities.PaymentServiceImportGroup{
				Id:       groupId,
				Services: map[string]bool{},
			}
		}
		group := model.Groups[groupName]

		if _, ok := group.Services[serviceName]; ok {
			FailureValSet = append(FailureValSet, bulkUpdateValidationFunc(strconv.Itoa(index), "Service_Group", "Service does not exists"))
			model.FailedRows++
			continue
		}

		if group.Id != -1 && dal.ServiceAlreadyExists(serviceName, strconv.Itoa(group.Id), nil) {
			FailureValSet = append(FailureValSet, bulkUpdateValidationFunc(strconv.Itoa(index), "Service_Group", "Service Already exists"))
			model.FailedRows++
			continue
		}

		group.Services[serviceName] = true
		model.ServicesCreated++
	}
	PsUpload = model

	if len(FailureValSet) > 0 {
		BulkUploadConfig.UpdateStatus = false
		BulkUploadConfig.Validations = FailureValSet

	} else {
		BulkUploadConfig.UpdateStatus = true
	}
	if BulkUploadConfig.UpdateStatus {

		fileName = time.Now().Format("20060102150405_") + fileName
		logging.Debug(TAG, "Renamed file to : "+fileName)

		if err := sendFileToFileServer(buf.Bytes(), fileName, BulkPaymentServiceUploadType); err != nil {
			logging.Error(TAG, "Bulk Payment Service Upload File Failed : ", err.Error())
			http.Error(w, uploadFileError, http.StatusInternalServerError)
			return
		}

		err := dal.InsertBulkApproval(fileName, BulkPaymentServiceUploadType, tmsUser.Username, common.UploadChangeType, 0)
		if err != nil {
			http.Error(w, DatabaseTxnError, http.StatusInternalServerError)
			return
		}
	}
	renderPartialTemplate(w, r, "bulkPaymentServiceUploadResults", BulkUploadConfig, tmsUser)

}

func commitPaymentServiceImportHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	if PsUpload == nil {
		respondWithError(w, &bulkImportError{
			statusCode: http.StatusBadRequest,
			respErrMsg: "error finding data to import",
			logErrMsg:  "error finding data to import",
		})
		return
	}

	if len(PsUpload.Groups) == 0 {
		respondWithError(w, &bulkImportError{
			statusCode: http.StatusBadRequest,
			respErrMsg: "error finding data to import",
			logErrMsg:  "error finding data to import",
		})
		return
	}

	var err error
	successfulRows := 0
	for key, val := range PsUpload.Groups {
		if val.Id == -1 {
			err = dal.AddPaymentServiceGroup(key)
			if err != nil {
				continue
			}
			val.Id = dal.GroupIdFromName(key)
			if val.Id == -1 {
				continue
			}
		}

		for serviceName, _ := range val.Services {
			err = dal.AddPaymentService(serviceName, strconv.Itoa(val.Id))
			if err != nil {
				continue
			}
			successfulRows++
		}
	}

	bytesToWrite := []byte("{\"importedRows\":" + strconv.Itoa(successfulRows) + "}")
	if _, err = w.Write(bytesToWrite); err != nil {
		_, _ = logging.Error(err.Error())
		http.Error(w, searchError, http.StatusInternalServerError)
		return
	}
}

func bulkPaymentServiceTidImportHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	var FailureValSet []models.ValidationSet
	var BulkUploadConfig models.BulkUpdateVal
	ctx := r.Context()
	file, header, inErr := getUploadedFile(r)
	if inErr != nil {
		respondWithError(w, inErr)
		return
	}
	defer file.Close()
	logging.Information(TAG, fmt.Sprintf("payment services import file %s uploaded", header.Filename))

	fileName := header.Filename

	inErr = verifyFileIsCsv(file, header.Filename)
	if inErr != nil {
		respondWithError(w, inErr)
		return
	}

	records, inErr := readAllCsvRecords(file, fileName)
	if inErr != nil {
		respondWithError(w, inErr)
		return
	}
	n := records[:]
	buf := &bytes.Buffer{}
	writer := csv.NewWriter(buf)
	writer.WriteAll(n)

	idxTID, idxMID := 0, 0
	for idx, column := range records[0] {
		if column == "Service_TID" {
			idxTID = idx
		}
		if column == "Service_Merchant" {
			idxMID = idx
		}
	}

	validator := validation.New(dal.NewValidationDal())
	tidErrors, midErrors := 0, 0
	maxWorkers := runtime.NumCPU()
	logging.Information(TAG, fmt.Sprintf("%d logical CPUs detected", maxWorkers))
	type ValidateJob struct {
		Type  string
		Value string
	}

	type ValidationResponse struct {
		Job     ValidateJob
		Err     error
		IsValid bool
	}

	// will hold the validation jobs in a queue and to be executed in parallel manner
	jobPool := make(chan ValidateJob)

	// definition of what worker will do
	worker := func(ctx context.Context, workerId int, in <-chan ValidateJob, out chan<- ValidationResponse) {
		for {
			select {
			case payload, ok := <-in: // iterate through in and returns when the channel is closed
				if !ok {
					logging.Information(TAG, fmt.Sprintf("worker %d is done", workerId))
					return
				}
				logging.Information(TAG, fmt.Sprintf("worker %d is validating %s %s", workerId, payload.Type, payload.Value))
				switch payload.Type {
				case TYPE_MID:
					isValid, err := validator.ValidateMID(payload.Value)
					out <- ValidationResponse{Job: payload, Err: err, IsValid: isValid}
				case TYPE_TID:
					isValid, err := validator.ValidateTid(payload.Value)
					out <- ValidationResponse{Job: payload, Err: err, IsValid: isValid}
				}
			case <-ctx.Done(): // returns when the user terminates the request
				return
			}
		}
	}

	publisher := func(ctx context.Context, jobs []ValidateJob, in chan<- ValidateJob) {
		defer close(in)
		for _, job := range jobs {
			select {
			case in <- job:
			// do nothing
			case <-ctx.Done():
				return
			}
		}
	}

	var jobs []ValidateJob
	for _, val := range records[1:] {
		jobs = append(jobs, ValidateJob{Type: TYPE_MID, Value: val[idxMID]}, ValidateJob{Type: TYPE_TID, Value: val[idxTID]})
	}

	// validationResponses will update the mid and tid validation errors
	validationResponses := make(chan ValidationResponse, len(jobs))
	defer close(validationResponses)

	for i := 0; i < maxWorkers; i++ {
		go worker(ctx, i, jobPool, validationResponses)
		logging.Information(TAG, fmt.Sprintf("%d of %d worker started", i+1, maxWorkers))
	}

	start := time.Now()
	go publisher(ctx, jobs, jobPool)
	var invalidError []string
	// this will update midErrors and tidErrors upon reporting
	for range jobs {
		select {
		case resp := <-validationResponses:
			switch resp.Job.Type {
			case TYPE_MID:
				if resp.Err != nil || !resp.IsValid {
					midErrors++
					invalidError = append(invalidError, fmt.Sprintf("Invalid %s: %s The MID (Service_Merchant) should consist of alphabetic characters and have a character length between 6 to 15 characters - %v", resp.Job.Type, resp.Job.Value, resp.Err))
					logging.Error(TAG, fmt.Sprintf("Invalid %s: %s - %v", resp.Job.Type, resp.Job.Value, resp.Err))
				} else {
					logging.Information(TAG, fmt.Sprintf("Valid %s: %s", resp.Job.Type, resp.Job.Value))
				}
			case TYPE_TID:
				if resp.Err != nil || !resp.IsValid {
					tidErrors++
					if !regexp.MustCompile(`^[0-9]*$`).MatchString(resp.Job.Value) && len(resp.Job.Value) < TID_MAXLENGTH {
						invalidError = append(invalidError, fmt.Sprintf("Invalid %s: %s - %s", resp.Job.Type, resp.Job.Value, "TID (Service_TID) should be numeric and have a fixed length of 8 characters."))
					} else {
						invalidError = append(invalidError, fmt.Sprintf("Invalid %s: %s - %v", resp.Job.Type, resp.Job.Value, resp.Err))
					}
					logging.Error(TAG, fmt.Sprintf("Invalid %s: %s - %v", resp.Job.Type, resp.Job.Value, resp.Err))
				} else {
					logging.Information(TAG, fmt.Sprintf("Valid %s: %s", resp.Job.Type, resp.Job.Value))
				}
			}
		}
	}
	logging.Information(TAG, "done receiving validation responses")

	duration := time.Since(start)
	logging.Information(TAG, fmt.Sprintf("%d validations are completed in %v", len(records)-1, duration))

	if tidErrors > 0 || midErrors > 0 {
		err := fmt.Errorf(strings.Join(invalidError, ", ")+`. There were %d invalid mid(s) and %d invalid tid(s)`, midErrors, tidErrors)
		respondWithError(w, &bulkImportError{
			statusCode: http.StatusBadRequest,
			respErrMsg: err.Error(),
			logErrMsg:  fmt.Sprintf("%d mid and %d tid validation errors found", midErrors, tidErrors),
		})
		return
	}

	if errCount, err := checkDuplicateRecords(records); err != nil {
		respondWithError(w, &bulkImportError{
			statusCode: http.StatusBadRequest,
			respErrMsg: err.Error(),
			logErrMsg:  fmt.Sprintf("%d Number of duplicate records found !", errCount),
		})
		return
	}

	_, groupIdx, serviceIdx, mMidIdx, mTidIdx, sMidIdx, sTidIdx := getPaymentServiceColumnIndexes(records[0])
	if groupIdx == -1 || serviceIdx == -1 || mMidIdx == -1 || mTidIdx == -1 || sMidIdx == -1 || sTidIdx == -1 {
		respondWithError(w, &bulkImportError{
			statusCode: http.StatusBadRequest,
			respErrMsg: "the uploaded file does not contain the correct columns",
			logErrMsg:  "the uploaded file does not contain the correct columns",
		})
		return
	}

	siteConf := map[string]siteInfo{}
	PsTermUpload = map[tidIdentifier]map[int]*dal.PaymentService{}
	successfulRows := 0

	// validate rows (check that the group matches the one assigned to the site and the service is valid for that group)
	// and build import model consisting of just valid rows
	for index, record := range records[1:] {
		if anyColumnsEmpty(record, groupIdx, serviceIdx, mMidIdx, mTidIdx, sMidIdx, sTidIdx) {
			continue
		}

		mid := record[mMidIdx]
		tid := record[mTidIdx]
		sTid := record[sTidIdx]

		conf, ok := siteConf[mid+"-"+tid]
		// we've already validated this mid - so we just need to validate the group and service names
		if ok {
			if len(conf.groupName) > 0 && conf.siteId > 0 && conf.groupName == record[groupIdx] && dal.IsServiceInGroup(record[groupIdx], record[serviceIdx]) {
				if appendToPsUploadModel(record[groupIdx], record[serviceIdx], record[sMidIdx], conf.siteId, tid, sTid) {
					successfulRows++
				}
			} else {
				FailureValSet = append(FailureValSet, bulkUpdateValidationFunc(strconv.Itoa(index), "Service_Group", "Service does not exists"))
			}
			continue
		}

		// we haven't encountered this mid yet, so we need to validate it
		profileIdStr, err := dal.GetProfileIdFromMID(mid)
		if err != nil || len(profileIdStr) == 0 {
			siteConf[mid+"-"+tid] = siteInfo{-1, ""}
			FailureValSet = append(FailureValSet, bulkUpdateValidationFunc(strconv.Itoa(index), "Master_Merchant", "SiteID does not exists"))
			continue
		}
		profileId, err := strconv.Atoi(profileIdStr)
		if err != nil || profileId < 1 {
			siteConf[mid+"-"+tid] = siteInfo{-1, ""}
			FailureValSet = append(FailureValSet, bulkUpdateValidationFunc(strconv.Itoa(index), "Master_Merchant", "SiteID does not exists"))
			continue
		}

		tidCount, err := dal.CheckTidExistsForMID(mid, tid)
		if err != nil || tidCount <= 0 {
			siteConf[mid+"-"+tid] = siteInfo{-1, ""}
			FailureValSet = append(FailureValSet, bulkUpdateValidationFunc(strconv.Itoa(index), "Master_TID", "TID does not exists to this MID"))
			continue
		}
		groupName := dal.GetProfileAssignedServiceGroupName(profileIdStr)
		siteConf[mid+"-"+tid] = siteInfo{profileId, groupName}
		if len(groupName) > 0 && groupName == record[groupIdx] && dal.IsServiceInGroup(record[groupIdx], record[serviceIdx]) {
			if appendToPsUploadModel(record[groupIdx], record[serviceIdx], record[sMidIdx], profileId, tid, sTid) {
				successfulRows++
			}
		} else {
			FailureValSet = append(FailureValSet, bulkUpdateValidationFunc(strconv.Itoa(index), "Service_Group", "Service does not exists"))
		}
	}
	if len(FailureValSet) > 0 {
		BulkUploadConfig.UpdateStatus = false
		BulkUploadConfig.Validations = FailureValSet

	} else {
		BulkUploadConfig.UpdateStatus = true
	}
	if BulkUploadConfig.UpdateStatus {

		fileName = time.Now().Format("20060102150405_") + fileName
		logging.Debug(TAG, "Renamed file to : "+fileName)

		if err := sendFileToFileServer(buf.Bytes(), fileName, BulkPaymentTidUploadType); err != nil {
			logging.Error(TAG, "Bulk Payment Service Upload File Failed : ", err.Error())
			http.Error(w, uploadFileError, http.StatusInternalServerError)
			return
		}

		err := dal.InsertBulkApproval(fileName, BulkPaymentTidUploadType, tmsUser.Username, common.UploadChangeType, 0)
		if err != nil {
			http.Error(w, DatabaseTxnError, http.StatusInternalServerError)
			return
		}
	}
	renderPartialTemplate(w, r, "bulkPaymentServiceUploadResults", BulkUploadConfig, tmsUser)
}

func checkDuplicateRecords(records [][]string) (int, error) {
	_, _, _, mMidIdx, mTidIdx, sMidIdx, sTidIdx := getPaymentServiceColumnIndexes(records[0])

	var duplicates, tidDuplicate, midDuplicate string
	var errCount int

	values := make(map[string]string)
	midValues := make(map[string]string)
	tidValues := make(map[string]string)

	for _, record := range records[1:] {
		rows := fmt.Sprintf("Master_Merchant: %s, Master_TID: %s, Service_Merchant: %s, Service_TID: %s", record[mMidIdx], record[mTidIdx], record[sMidIdx], record[sTidIdx])
		if _, ok := values[rows]; ok {
			duplicates += rows + "<br/>"
			errCount++
			continue
		}
		midRows := fmt.Sprintf("Service_Merchant: %s", record[sMidIdx])
		if _, ok := midValues[midRows]; ok {
			midDuplicate += midRows + "<br/>"
			errCount++
			continue
		}
		tidRows := fmt.Sprintf("Service_TID: %s", record[sTidIdx])
		if _, ok := tidValues[tidRows]; ok {
			tidDuplicate += tidRows + "<br/>"
			errCount++
			continue
		}
		values[rows] = ""
		midValues[midRows] = ""
		tidValues[tidRows] = ""
	}
	if errCount > 0 {
		return errCount, fmt.Errorf(fmt.Sprintf("Identify (%d) Duplicate Records: <br/> %s %s %s", errCount, duplicates, midDuplicate, tidDuplicate))
	}
	return errCount, nil
}

func appendToPsUploadModel(groupName, serviceName, sMid string, siteId int, tid string, sTid string) bool {
	serviceId := dal.GetServiceIdFromNames(groupName, serviceName)
	if serviceId == -1 {
		return false
	}

	id := tidIdentifier{
		tid:    tid,
		siteId: siteId,
	}

	tidUploads, ok := PsTermUpload[id]
	if !ok {
		tidUploads = map[int]*dal.PaymentService{}
		PsTermUpload[id] = tidUploads
	}

	service, ok := tidUploads[serviceId]
	if !ok {
		service = &dal.PaymentService{
			ServiceId: serviceId,
			Name:      serviceName,
			MID:       sMid,
			TID:       sTid,
		}
		tidUploads[serviceId] = service
		return true
	}
	service.MID = sMid
	service.TID = sTid
	return true
}

func commitPaymentServiceTidImportHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	if PsTermUpload == nil || len(PsTermUpload) == 0 {
		http.Error(w, searchError, http.StatusInternalServerError)
		return
	}

	rowsUpdated := 0
	for id, upload := range PsTermUpload {
		serviceJson := MergePaymentsWithConfigured(upload, id.tid, id.siteId)
		bytes, err := json.Marshal(serviceJson)
		if err != nil {
			continue
		}

		if services.SaveTerminalPaymentServiceConfig(id.tid, string(bytes), tmsUser) {
			rowsUpdated++
		}
	}

	bytesToWrite := []byte("{\"importedRows\":" + strconv.Itoa(rowsUpdated) + "}")
	if _, err := w.Write(bytesToWrite); err != nil {
		_, _ = logging.Error(err.Error())
		http.Error(w, searchError, http.StatusInternalServerError)
		return
	}
}

func getUploadedFile(r *http.Request) (file multipart.File, header *multipart.FileHeader, error *bulkImportError) {
	err := r.ParseMultipartForm(10000)
	if err != nil {
		error = &bulkImportError{
			statusCode: http.StatusBadRequest,
			respErrMsg: uploadFileError,
			logErrMsg:  err.Error(),
		}
		return
	}

	if len(r.MultipartForm.File) < 1 {
		error = &bulkImportError{
			statusCode: http.StatusBadRequest,
			respErrMsg: "please upload a file",
			logErrMsg:  "bulk import started with no file attached",
		}
		return
	}

	file, header, err = r.FormFile("file")
	if err != nil {
		error = &bulkImportError{
			statusCode: http.StatusBadRequest,
			respErrMsg: uploadFileError,
			logErrMsg:  err.Error(),
		}
	}
	return
}

func verifyFileIsCsv(file multipart.File, filename string) *bulkImportError {
	buf := make([]byte, 512)
	if _, err := file.Seek(0, 0); err != nil {
		return &bulkImportError{
			statusCode: http.StatusInternalServerError,
			respErrMsg: html.EscapeString(FailedFileRead + filename),
			logErrMsg:  err.Error(),
		}
	}
	if _, err := file.Read(buf); err != nil {
		return &bulkImportError{
			statusCode: http.StatusInternalServerError,
			respErrMsg: html.EscapeString(FailedFileRead + filename),
			logErrMsg:  err.Error(),
		}
	}
	if !validateFileType(buf, FiletypeCsv) {
		return &bulkImportError{
			statusCode: http.StatusInternalServerError,
			respErrMsg: IncorrectFileTypeCSV,
			logErrMsg:  "Incorrect filetype provided",
		}
	}

	// reset file reader to beginning after validating
	if _, err := file.Seek(0, 0); err != nil {
		return &bulkImportError{
			statusCode: http.StatusInternalServerError,
			respErrMsg: html.EscapeString(FailedFileRead + filename),
			logErrMsg:  err.Error(),
		}
	}
	return nil
}

func readAllCsvRecords(file multipart.File, filename string) ([][]string, *bulkImportError) {
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, &bulkImportError{
			statusCode: http.StatusInternalServerError,
			respErrMsg: html.EscapeString(FailedFileRead + filename),
			logErrMsg:  err.Error(),
		}
	}
	if len(records) == 0 || len(records[1:]) < 1 {
		return nil, &bulkImportError{
			statusCode: http.StatusInternalServerError,
			respErrMsg: NoColumnsFound,
			logErrMsg:  "No records found in uploaded CSV file",
		}
	}
	return records, nil
}

func respondWithError(w http.ResponseWriter, err *bulkImportError) {
	logging.Error(TAG, err.logErrMsg)
	http.Error(w, err.respErrMsg, err.statusCode)
}

func getPaymentServiceColumnIndexes(header []string) (numIdx, groupIdx, serviceIdx, mMidIdx, mTidIdx, sMidIdx, sTidIdx int) {
	numIdx, groupIdx, serviceIdx, mMidIdx, mTidIdx, sMidIdx, sTidIdx = -1, -1, -1, -1, -1, -1, -1
	for i, column := range header {
		switch strings.ToLower(column) {
		case "sl_num":
			numIdx = i
		case "service_group":
			groupIdx = i
		case "service_label":
			serviceIdx = i
		case "master_merchant":
			mMidIdx = i
		case "master_tid":
			mTidIdx = i
		case "service_merchant":
			sMidIdx = i
		case "service_tid":
			sTidIdx = i
		}
	}
	return
}

func anyColumnsEmpty(record []string, idxs ...int) bool {
	for _, i := range idxs {
		if i < 0 || i > len(record) {
			return true
		}
		if len(record[i]) == 0 {
			return true
		}
	}
	return false
}
