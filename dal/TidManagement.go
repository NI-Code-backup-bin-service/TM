package dal

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"log"
	"math"
	"nextgen-tms-website/common"
	"nextgen-tms-website/crypt"
	"nextgen-tms-website/entities"
	"nextgen-tms-website/logger"
	"nextgen-tms-website/models"
	"nextgen-tms-website/resultCodes"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

type OTPIntent int

const (
	OtpIntentEnrolment OTPIntent = iota
	OtpIntentReset
)

type TidPagination struct {
	FirstRecord  int
	LastRecord   int
	TotalRecords int
	CurrentPage  int
	PageSize     int
	PageCount    int
	More         bool
	Less         bool
	SearchTerm   string
	Pages        []PaginationPage
}

type PaginationPage struct {
	PageNumber string
	Selected   bool
	Active     bool
}

func StoreOneTimePIN(tid string, PIN string, intent OTPIntent, user *entities.TMSUser) (string, error) {
	db, err := GetDB()
	if err != nil {
		return "", err
	}

	var time string
	switch intent {
	case OtpIntentEnrolment:
		_, err = db.Exec("CALL store_enrolment_PIN(?,?,?)", tid, PIN, 15)
		if err != nil {
			return "", err
		}

		rows, err := db.Query("SELECT ExpiryDate FROM tid WHERE tid_id = ?", tid)
		if err != nil {
			return "", err
		}
		defer rows.Close()

		for rows.Next() {
			rows.Scan(&time)
		}
	case OtpIntentReset:
		_, err = db.Exec("CALL store_reset_PIN(?,?,?)", tid, PIN, 5)
		if err != nil {
			return "", err
		}

		rows, err := db.Query("SELECT reset_pin_expiry_date FROM tid WHERE tid_id = ?", tid)
		if err != nil {
			return "", err
		}
		defer rows.Close()

		for rows.Next() {
			rows.Scan(&time)
		}

		tidInt, err := strconv.Atoi(tid)
		if err != nil {
			return "", err
		}
		err = LogTidChange(tidInt, 1, ApproveCreate, "", "Reset PIN Generated", user.Username, true)
		if err != nil {
			return "", err
		}
	}

	return time, err
}

// GetTids needs to be passed the sitewide data groups, obtained using `FindForSiteByProfileId`
// It is passed in as a parameter rather than queried within this function as the groups may
// need to referenced outside of this function (such as with GetSiteData), and it is an expensive query
// so duplication of the query is especially egregious
// without tidchangehistory data
func GetTids(db *sql.DB, siteId int, activeGroups map[string]bool, pageSize int, pageNumber int, tidSearchTerm string) ([]*TIDData, TidPagination, error) {
	var allTids = []*TIDData{}
	var pagination TidPagination

	rows, err := db.Query("Call get_site_tids(?)", siteId)
	if err != nil {
		logging.Error(err)
		return nil, pagination, err
	}
	defer rows.Close()

	for rows.Next() {
		var tidData TIDData
		var enrolmentPIN sql.NullString
		var resetPIN sql.NullString
		var expiryTime sql.NullString
		var activationTime sql.NullString
		var merchantID sql.NullString
		var userOverrides int
		var fraudOverride int
		var resetPinExpiryTime sql.NullString
		var tidProfileID int

		err = rows.Scan(&tidData.TID, &tidData.Serial, &enrolmentPIN, &expiryTime, &resetPIN, &resetPinExpiryTime, &activationTime, &merchantID, &userOverrides, &fraudOverride, &tidProfileID)
		if err != nil {
			logging.Error(err)
			return nil, pagination, err
		}

		//My kingdom for a ternary operator
		if enrolmentPIN.Valid {
			tidData.EnrolmentPIN = enrolmentPIN.String
		}
		if expiryTime.Valid {
			tidData.ExpiryTime = expiryTime.String
		}
		if activationTime.Valid {
			tidData.ActivationTime = activationTime.String
		}
		if resetPIN.Valid {
			tidData.ResetPIN = resetPIN.String
		}
		if resetPinExpiryTime.Valid {
			tidData.ResetPinExpiryTime = resetPinExpiryTime.String
		}
		if merchantID.Valid {
			tidData.MerchantID = merchantID.String
		}

		tidData.UserOverrides = userOverrides > 0

		tidData.FraudOverride = fraudOverride > 0

		tidData.TIDProfileID = tidProfileID

		tidData.SiteId = siteId

		allTids = append(allTids, &tidData)
	}

	pagination.SearchTerm = tidSearchTerm
	// Filter the TIDs to display based on the search term entered
	if tidSearchTerm != "" || len(tidSearchTerm) > 0 {
		allTids = filterTidsBySearchTerm(allTids, tidSearchTerm)
	}

	var selectedTids []*TIDData

	// Set up the pagination navigation
	// Calculate the total number of pages
	pageCount := int(math.Ceil(float64(len(allTids)) / float64(pageSize)))

	if pageCount == 1 || pageCount == 0 { // If there is only one page then there is no next or previous page to view
		pagination.Less = false
		pagination.More = false

		var page PaginationPage
		page.PageNumber = "1"
		page.Selected = true
		page.Active = true
		pagination.Pages = append(pagination.Pages, page)
	} else {
		if pageNumber == 1 { // If this is the first page, then only Next should display
			pagination.Less = false
			pagination.More = true
		} else if pageNumber == pageCount { // If this is the last page, then only Previous should display
			pagination.Less = true
			pagination.More = false
		} else { // If this is neither the first or last page, show Previous and Next
			pagination.Less = true
			pagination.More = true
		}
		pagination.Pages = buildPageChangeOptions(pageCount, pageNumber)
	}

	// If the page size is larger than the total TIDs on the site or is set to -1 then we need to read all the TIDs
	if pageSize > len(allTids) || pageSize == -1 {
		selectedTids = allTids
		if len(allTids) == 0 {
			pagination.FirstRecord = 0 // If there are no PEDs set the first record to 0
		} else {
			pagination.FirstRecord = 1 // This will always be 1 when viewing all PEDs where there is at least one PED
		}
		pagination.LastRecord = len(allTids)
		pagination.TotalRecords = len(allTids)
	} else {
		// To cut down on the number of DB calls we need to limit the number of TIDs we populate
		// Sort all of the TIDs by TID number
		sort.Slice(allTids, func(i, j int) bool { return allTids[i].TID < allTids[j].TID })
		// Now that the TIDs are in a sensible order, we need to calculate the offset we use to find the first TID we want
		// i.e for page 3 and a page size of 10 we need records starting at position 20 - 10 * (3 - 1) = 20
		offset := pageSize * (pageNumber - 1)
		// Calculate the last record we want, based on the current page size
		lastRecord := offset + (pageSize)
		if lastRecord > len(allTids) {
			lastRecord = len(allTids)
		}
		// Obtain the desired records as a slice of the total
		selectedTids = allTids[offset:lastRecord]

		pagination.FirstRecord = offset + 1
		pagination.LastRecord = lastRecord
		pagination.TotalRecords = len(allTids)
	}

	pagination.CurrentPage = pageNumber
	pagination.PageSize = pageSize
	pagination.PageCount = pageCount

	worker := func(ctx context.Context, src <-chan *TIDData, out chan<- error) {
		for {
			select {
			case <-ctx.Done():
				return
			case tid, ok := <-src:
				if !ok {
					return
				}
				out <- AddTidProfileData(tid, db, activeGroups)
			}
		}
	}

	numCPU := runtime.NumCPU() / 2
	if numCPU < 1 {
		numCPU = 1
	}
	sources := make(chan *TIDData, numCPU)
	outputs := make(chan error, numCPU)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for i := 0; i < numCPU; i++ {
		go worker(ctx, sources, outputs)
	}
	publisher := func(src chan<- *TIDData, tids []*TIDData) {
		defer close(src)
		for _, tid := range tids {
			src <- tid
		}
	}

	go publisher(sources, selectedTids)
	for i := 0; i < len(selectedTids); i++ {
		if e := <-outputs; e != nil {
			cancel()
			return nil, pagination, e
		}
	}

	return selectedTids, pagination, nil
}

func AddTidProfileData(tid *TIDData, db *sql.DB, activeGroups map[string]bool) error {
	if tid.TIDProfileID != 0 {
		pd, err := GetTidProfileData(db, tid.TIDProfileID, activeGroups)
		if err != nil {
			logging.Error(err)
			return err
		}
		tid.TIDProfileGroups = pd
		tid.Overridden = true
	} else {
		tid.Overridden = false
	}

	return nil
}

func GetSiteTids(db *sql.DB, siteId int) ([]int, error) {
	rows, err := db.Query("select tid_id from tid_site where site_id= ?", siteId)
	if err != nil {
		logging.Error(err)
		return nil, err
	}
	defer rows.Close()

	tids := make([]int, 0)
	for rows.Next() {
		var tid int
		err = rows.Scan(&tid)
		if err != nil {
			logging.Error(err)
			return nil, err
		}

		tids = append(tids, tid)
	}

	return tids, nil
}

func GetTIDOverRideDataElement() (map[string]bool, error) {
	dataElementId := make(map[string]bool, 0)
	db, err := GetDB()
	if err != nil {
		logging.Error("An error occurred querying procedure get_tid_override_data_element_ids - %s", err.Error())
		return dataElementId, err
	}

	rows, err := db.Query("CALL get_tid_override_data_element_ids()")
	if err != nil {
		log.Println("unable to query get_tid_override_data_element_ids", err)
		logging.Error("An error occurred querying procedure get_tid_override_data_element_ids - %s", err.Error())
		return dataElementId, err
	}
	defer rows.Close()

	for rows.Next() {
		var deId string
		err = rows.Scan(&deId)
		if err != nil {
			logging.Error("An error occurred querying procedure get_tid_override_data_element_ids - %s", err.Error())
			return dataElementId, err
		}

		if _, ok := dataElementId[deId]; !ok {
			dataElementId[deId] = true
		}
	}
	return dataElementId, nil
}

func filterTidsBySearchTerm(allTids []*TIDData, searchTerm string) []*TIDData {
	var filteredTids []*TIDData

	for _, tid := range allTids {
		tidString := GetPaddedTidId(tid.TID)
		serial := tid.Serial
		// If either the TID or serial contain the searchTerm then include it in the results
		if strings.Contains(tidString, searchTerm) || strings.Contains(serial, searchTerm) {
			filteredTids = append(filteredTids, tid)
		}
	}
	return filteredTids
}

/*
If there are 5 or less pages of TIDs then we need to display all the pages.
If there are more than 5 pages of TIDs then we need to display the current page,
two pages either side, as well as a first/last page
*/
func buildPageChangeOptions(pageCount int, pageNumber int) []PaginationPage {
	var pages []PaginationPage
	showAll := pageCount <= 5

	var startPage int
	var endPage int

	if showAll {
		endPage = pageCount
		startPage = 1
	} else {
		switch pageNumber {
		case 1, 2, 3:
			startPage = 1
			endPage = startPage + 4
		case pageCount, pageCount - 1, pageCount - 2:
			endPage = pageCount
			startPage = pageCount - 4
		default:
			startPage = pageNumber - 2
			endPage = pageNumber + 2
		}
	}

	// If we have more than 5 pages, and the current page number is not 1, 2 or 3
	if !showAll && pageNumber-2 > 1 {
		var page PaginationPage
		page.PageNumber = "1"
		page.Selected = false
		page.Active = true
		pages = append(pages, page)

		// This will append an ellipsis so the pagination navigation appears as it does on search page
		pages = append(pages, makeEllipsisPage())
	}

	for i := startPage; i <= endPage; i++ {
		var page PaginationPage
		page.PageNumber = strconv.Itoa(i)
		if i == pageNumber {
			page.Selected = true
		} else {
			page.Selected = false
		}
		page.Active = true
		pages = append(pages, page)
	}

	// If we have more than 5 pages, and the current page number is not one of the last 3 pages
	if !showAll && pageNumber < pageCount-2 {

		// This will append an ellipsis so the pagination navigation appears as it does on search page
		pages = append(pages, makeEllipsisPage())

		var page PaginationPage
		page.PageNumber = strconv.Itoa(pageCount)
		page.Selected = false
		page.Active = true
		pages = append(pages, page)
	}

	return pages
}

/*
*
Constructs a non-interactable "..." page for the Site TID pagination
*/
func makeEllipsisPage() PaginationPage {
	var page PaginationPage
	page.PageNumber = "..."
	page.Selected = false
	page.Active = false

	return page
}

// Fetches a db connection and call getTidFraudOverride
func TidFraudOverrideStatus(tidID int) (bool, error) {
	db, err := GetDB()
	if err != nil {
		return false, err
	}

	return getTidFraudOverride(db, tidID)
}

func getTidFraudOverride(db *sql.DB, tidID int) (bool, error) {
	var rowcount int

	//Check to see if TID overrides have been set
	err := db.QueryRow("SELECT COUNT(*) FROM velocity_limits WHERE tid_id = ?", tidID).Scan(&rowcount)

	if err != nil {
		return false, err
	}

	if rowcount < 1 {
		return false, nil
	} else {
		return true, nil
	}
}

func GetTidProfileData(db *sql.DB, profileId int, activeGroups map[string]bool) ([]*DataGroup, error) {
	rows, err := db.Query("Call profile_data_fetch(?)", profileId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profileGroups []*DataGroup

	for rows.Next() {
		var site SiteData
		var siteOptions string
		err = rows.Scan(&site.TIDID, &site.DataGroupID, &site.DataGroupDisplayName, &site.DataElementID, &site.Name,
			&site.DataGroup, &site.TIDOverridable, &site.DataType, &site.Tooltip, &site.DataValue, &site.IsAllowEmpty,
			&site.MaxLength, &site.ValidationExpression, &site.ValidationMessage, &site.FrontEndValidate, &siteOptions, &site.DisplayName, &site.IsPassword, &site.IsEncrypted, &site.SortOrderInGroup, &site.IsReadOnlyAtCreation, &site.IsRequiredAtAcquireLevel, &site.IsRequiredAtChainLevel)
		if err != nil {
			return nil, err
		}

		//I do not THINK that the DV needs decrypting here.
		if site.DataValue.Valid && site.IsEncrypted.Valid && site.IsEncrypted.Bool {
			site.DataValue.String, err = crypt.Decrypt(site.DataValue.String)
			if err != nil {
				return nil, err
			}
		}

		// Build options model for data element and compare with data value to default selections
		site.Options, site.OptionSelectable = BuildOptionsData(siteOptions, site.DataValue.String, site.DataGroup, site.Name, profileId)

		if site.TIDOverridable && (site.DataValue.Valid || activeGroups[site.DataGroup]) {
			profileGroups = addDataElement(site, profileGroups)
		}
	}

	return profileGroups, nil
}

func FetchActiveGroupNames(profileId int) (map[string]bool, error) {
	db, err := GetDB()
	if err != nil {
		logging.Error(err)
		return nil, err
	}

	rows, err := db.Query("Call fetch_active_group_by_profileId(?)", profileId)
	if err != nil {
		logging.Error(err)
		return nil, err
	}
	defer rows.Close()

	activeGroups := map[string]bool{}

	var name string
	for rows.Next() {
		err = rows.Scan(&name)
		if err != nil {
			logging.Error(err)
			return nil, err
		}
		activeGroups[name] = true
	}

	return activeGroups, nil
}

func getFilteredTIDCount(searchTerm string, acquirers string) (int, error) {
	db, err := GetDB()
	if err != nil {
		logging.Error(err)
		return 0, err
	}

	logging.Information("Calling get_tid_count_filtered")

	rows, err := db.Query("CALL get_tid_count_filtered(?,?)", searchTerm, acquirers)
	if err != nil {
		logging.Error(err)
		return 0, err
	}
	defer rows.Close()

	var tidCount int
	for rows.Next() {
		rows.Scan(&tidCount)
	}

	return tidCount, nil
}

func GetTIDPage(searchTerm string, offset string, amount string, orderedColumn string, orderDirection string, user *entities.TMSUser) (page []*TIDData, total int, filtered int, err error) {
	// Find user acquirers to limit search results
	acquirers, err := GetUserAcquirerPermissions(user)
	if err != nil {
		return nil, 0, 0, err
	}

	total, err = getFilteredTIDCount(searchTerm, acquirers)
	if err != nil {
		logging.Error(err)
		return nil, 0, 0, err
	}

	// If the amount is -1 then we want all the TIDs
	if amount == "-1" {
		amount = strconv.Itoa(total)
	}

	db, err := GetDB()
	if err != nil {
		return nil, 0, 0, err
	}

	logging.Information("Calling get_tid_page")

	rows, err := db.Query("CALL get_tid_page(?, ?, ?, ?)", searchTerm, offset, amount, acquirers)
	if err != nil {
		logging.Error(err)
		return nil, 0, 0, err
	}
	defer rows.Close()

	var tids []*TIDData

	logging.Information("Looping through rows")

	for rows.Next() {
		var tidData TIDData
		var enrolmentPIN sql.NullString
		var resetPIN sql.NullString
		var activationTime sql.NullString
		var expiryTime sql.NullString
		var resetPinExpiryTime sql.NullString
		var acquirerName string
		var profileID sql.NullInt64
		err = rows.Scan(&tidData.TID, &tidData.Serial, &enrolmentPIN, &expiryTime, &resetPIN, &resetPinExpiryTime, &activationTime, &tidData.SiteId, &tidData.SiteName, &tidData.MerchantID, &acquirerName, &profileID)
		if err != nil {
			logging.Warning("pageQuery error: " + err.Error())
		}

		if enrolmentPIN.Valid {
			tidData.EnrolmentPIN = enrolmentPIN.String
		}
		if expiryTime.Valid {
			tidData.ExpiryTime = expiryTime.String
		}
		if resetPIN.Valid {
			tidData.ResetPIN = resetPIN.String
		}
		if resetPinExpiryTime.Valid {
			tidData.ResetPinExpiryTime = resetPinExpiryTime.String
		}
		if activationTime.Valid {
			tidData.ActivationTime = activationTime.String
		}
		if profileID.Valid {
			tidData.TIDProfileID = int(profileID.Int64)
		}
		tidData.SiteName = html.EscapeString(tidData.SiteName)

		tids = append(tids, &tidData)
	}

	logging.Information("Finished looping through %d rows", total)

	//sort the slice
	sort.Slice(tids, func(i, j int) bool {
		switch orderedColumn {
		case "1":
			if strings.ToUpper(orderDirection) == "ASC" {
				return tids[i].TID < tids[j].TID
			} else {
				return tids[i].TID > tids[j].TID
			}
		case "2":
			if strings.ToUpper(orderDirection) == "ASC" {
				return strings.ToLower(tids[i].Serial) < strings.ToLower(tids[j].Serial)
			} else {
				return strings.ToLower(tids[i].Serial) > strings.ToLower(tids[j].Serial)
			}
		case "3":
			if strings.ToUpper(orderDirection) == "ASC" {
				return tids[i].EnrolmentPIN < tids[j].EnrolmentPIN
			} else {
				return tids[i].EnrolmentPIN > tids[j].EnrolmentPIN
			}
		case "4":
			if strings.ToUpper(orderDirection) == "ASC" {
				return tids[i].ResetPIN < tids[j].ResetPIN
			} else {
				return tids[i].ResetPIN > tids[j].ResetPIN
			}
		case "5":
			if strings.ToUpper(orderDirection) == "ASC" {
				return tids[i].ActivationTime < tids[j].ActivationTime
			} else {
				return tids[i].ActivationTime > tids[j].ActivationTime
			}
		case "6":
			if strings.ToUpper(orderDirection) == "ASC" {
				return tids[i].MerchantID < tids[j].MerchantID
			} else {
				return tids[i].MerchantID > tids[j].MerchantID
			}
		case "7":
			if strings.ToUpper(orderDirection) == "ASC" {
				return strings.ToLower(tids[i].SiteName) < strings.ToLower(tids[j].SiteName)
			} else {
				return strings.ToLower(tids[i].SiteName) > strings.ToLower(tids[j].SiteName)
			}
		}
		return false
	})

	filtered = total
	logging.Information("Finished building TID Page")
	return tids, total, filtered, nil
}

func generateTIDSearchWhere(searchTerm string) string {
	if searchTerm == "" {
		return ""
	}

	searchTerm = "%" + strings.ToUpper(searchTerm) + "%"

	return fmt.Sprintf(`WHERE (upper(l.tid_id) LIKE "%[1]v"
	OR upper(l.serial) LIKE "%[1]v"
	OR upper(l.site_name) LIKE "%[1]v")`, searchTerm)
}

func CheckSerialInUse(serial string) error {
	currentAssignment, err := GetTidBySerial(serial)
	if err != nil {
		return err
	}

	if currentAssignment != "" {
		return errors.New("serial number already assigned")
	}
	return nil
}

func AddTidToSiteAndSaveTidProfileChange(tidID, serial string, site int, auto bool, autoTime string, profileId int, userName string, changeType int, approved int) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	tidInt, err := strconv.Atoi(tidID)
	if err != nil {
		return err
	}

	if _, err := db.Exec("Call save_tid_to_site_and_pending_profile_change(?,?,?,?,?,?,?,?,?,?,?)", profileId, changeType, "TID Created", userName, tidID, approved, serial, site, auto, autoTime, tidInt); err != nil {
		return err
	}
	return nil
}

func AddTidToSite(tid string, serial string, site int, auto bool, autoTime string) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	_, err = db.Exec("Call save_tid(?, ?, ?, ?, ?)", tid, serial, site, auto, autoTime)
	if err != nil {
		return err
	}

	return nil
}

func GetTidBySerial(serialNumber string) (string, error) {
	db, err := GetDB()
	if err != nil {
		return "", err
	}

	rows, err := db.Query("SELECT tid_id FROM tid where serial = ? ", serialNumber)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var tid string
	for rows.Next() {
		err = rows.Scan(&tid)
		if err != nil {
			return "", err
		}
	}
	return tid, nil
}

func GetSerialByTid(tid string) (string, error) {
	db, err := GetDB()

	if err != nil {
		return "", err
	}
	rows, err := db.Query("SELECT serial FROM tid where tid_id = ? ", tid)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var serial string
	for rows.Next() {
		err = rows.Scan(&serial)
		if err != nil {
			return "", err
		}
	}

	if serial == "" {
		return "", errors.New("Serial Number does not Exist for given TID")
	}
	return serial, nil

}

func RemoveTidFromSite(tid string, site string) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	_, err = db.Exec("Call delete_tid(?, ?)", tid, site)
	if err != nil {
		return err
	}

	// With the TID deleted the overridden TID velocity limits now need to be deleted
	err = removeTidVelocityLimits(tid, site)
	if err != nil {
		return err
	}

	return nil
}

// Ensures that Tid velocity limits are deleted when the deletion of a TID is approved in Change Approval
func removeTidVelocityLimits(tid string, site string) error {

	// Convert the incoming strings to int
	tidAsInt, err := strconv.Atoi(tid)
	if err != nil {
		return err
	}

	siteAsInt, err := strconv.Atoi(site)
	if err != nil {
		return err
	}

	// Delete the non-scheme velocity limits (level 4)
	err = DeleteSiteVelocityLimits(siteAsInt, 4, tidAsInt)
	if err != nil {
		return err
	}

	// Delete the scheme velocity limits (level 2)
	err = DeleteSiteVelocityLimits(siteAsInt, 2, tidAsInt)
	if err != nil {
		return err
	}

	return nil
}

func GetTidAcquirer(tid string) (string, error) {
	db, err := GetDB()
	if err != nil {
		return "", err
	}

	rows, err := db.Query("CALL get_tid_acquirer(?)", tid)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var acqName sql.NullString
	for rows.Next() {
		rows.Scan(&acqName)
	}

	return acqName.String, nil
}

func GetTidDetails(tid string) (*TidDetails, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	tidDetails := &TidDetails{}
	var lastTransaction sql.NullString
	var lastAPKDownloadTime sql.NullString
	var lastChecked sql.NullInt64
	var confirmedTime sql.NullInt64
	var coordinates sql.NullString
	var accuracy sql.NullString
	var lastCoordinateTime sql.NullString
	var freeInternalStorage sql.NullString
	var totalInternalStorage sql.NullString
	var softuiLastDownloadedFilename sql.NullString
	var softuiLastDownloadedFileHash sql.NullString
	var softuiLastDownloadedFileList sql.NullString
	var softuiLastDownloadedFileDateTime sql.NullString

	row := db.QueryRow("CALL fetch_tid_details(?)", tid)

	err = row.Scan(&tidDetails.AppVer, &tidDetails.FirmwareVer, &lastTransaction, &lastChecked, &lastAPKDownloadTime, &confirmedTime, &tidDetails.EODAuto, &tidDetails.AutoTime, &coordinates, &accuracy, &lastCoordinateTime, &freeInternalStorage, &totalInternalStorage, &softuiLastDownloadedFilename, &softuiLastDownloadedFileHash, &softuiLastDownloadedFileList, &softuiLastDownloadedFileDateTime)
	if err != nil {
		logging.Error(err)
		return nil, err
	}

	if lastTransaction.Valid {
		tidDetails.LastTransaction = lastTransaction.String
	} else {
		tidDetails.LastTransaction = "None"
	}

	if lastAPKDownloadTime.Valid {
		tidDetails.LastAPKDownloadTime = lastAPKDownloadTime.String
	} else {
		tidDetails.LastAPKDownloadTime = "None"
	}

	if lastChecked.Valid {
		tidDetails.LastCheckedTime = time.Unix(lastChecked.Int64/1000, 0).String()
	} else {
		tidDetails.LastCheckedTime = "None"
	}

	if confirmedTime.Valid {
		tidDetails.ConfirmedTime = time.Unix(confirmedTime.Int64/1000, 0).String()
	} else {
		tidDetails.ConfirmedTime = "None"
	}

	if coordinates.Valid {
		tidDetails.Coordinates = coordinates.String
	} else {
		tidDetails.Coordinates = "None"
	}

	if accuracy.Valid {
		tidDetails.Accuracy = accuracy.String
	} else {
		tidDetails.Accuracy = "None"
	}

	if lastCoordinateTime.Valid {
		tidDetails.LastCoordinateTime = lastCoordinateTime.String
	} else {
		tidDetails.LastCoordinateTime = "None"
	}

	if freeInternalStorage.Valid {
		tidDetails.FreeInternalStorage = freeInternalStorage.String
	} else {
		tidDetails.FreeInternalStorage = "None"
	}

	if totalInternalStorage.Valid {
		tidDetails.TotalInternalStorage = totalInternalStorage.String
	} else {
		tidDetails.TotalInternalStorage = "None"
	}

	if softuiLastDownloadedFilename.Valid {
		tidDetails.SoftuiLastDownloadedFileName = softuiLastDownloadedFilename.String
	} else {
		tidDetails.SoftuiLastDownloadedFileName = "None"
	}

	if softuiLastDownloadedFileHash.Valid {
		tidDetails.SoftuiLastDownloadedFileHash = softuiLastDownloadedFileHash.String
	} else {
		tidDetails.SoftuiLastDownloadedFileHash = "None"
	}

	if softuiLastDownloadedFileDateTime.Valid {
		tidDetails.SoftuiLastDownloadedFileDateTime = softuiLastDownloadedFileDateTime.String
	} else {
		tidDetails.SoftuiLastDownloadedFileDateTime = "None"
	}

	if softuiLastDownloadedFileList.Valid && strings.TrimSpace(softuiLastDownloadedFileList.String) != "" {
		softuiLastDownloadedFileListMap := map[string][]map[string]string{}
		err = json.Unmarshal([]byte(softuiLastDownloadedFileList.String), &softuiLastDownloadedFileListMap)
		if err != nil {
			logging.Error(err)
			return nil, err
		}

		var ok bool
		tidDetails.SoftuiLastDownloadedFileList, ok = softuiLastDownloadedFileListMap["SoftUI"]
		if !ok {
			tidDetails.SoftuiLastDownloadedFileList = []map[string]string{}
		}
	} else {
		tidDetails.SoftuiLastDownloadedFileList = []map[string]string{}
	}

	return tidDetails, nil
}

func GetDataValueForTidProfileId(tidId int, siteId int) (string, error) {
	db, err := GetDB()
	if err != nil {
		return "", err
	}

	rows, err := db.Query("select pd.datavalue as dataValue from tid_site ts left join profile_data pd ON ts.tid_profile_id = pd.profile_id left join data_element de ON pd.data_element_id = de.data_element_id WHERE ts.tid_id = (?) and ts.site_id = (?) and de.name = (?)", tidId, siteId, "partialPackageName")
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var dataValue string
	for rows.Next() {
		rows.Scan(&dataValue)
	}

	return dataValue, nil
}

func GetDataValueForSiteTidAndElementName(tidId int, siteId int, dataElementName string) (string, error) {
	db, err := GetDB()
	if err != nil {
		return "", err
	}

	rows, err := db.Query("select pd.datavalue as dataValue from tid_site ts left join profile_data pd ON ts.tid_profile_id = pd.profile_id left join data_element de ON pd.data_element_id = de.data_element_id WHERE ts.tid_id = (?) and ts.site_id = (?) and de.name = (?)", tidId, siteId, dataElementName)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var dataValue string
	for rows.Next() {
		rows.Scan(&dataValue)
	}

	return dataValue, nil
}

func GetTIDUpdates(tidId string) ([]*models.TIDUpdateData, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	var updateData []*models.TIDUpdateData

	rows, err := db.Query("SELECT tid_update_id, target_package_id, DATE_FORMAT(update_date, '%Y-%m-%d %H:%i'), third_party_apk FROM tid_updates where tid_id = ? ", tidId)
	if err != nil {
		return updateData, err
	}
	defer rows.Close()

	for rows.Next() {
		var update models.TIDUpdateData
		var thirdPartyApkID sql.NullString
		var finalOptionTpApkID []int
		err = rows.Scan(&update.UpdateID, &update.PackageID, &update.UpdateDate, &thirdPartyApkID)
		if err != nil {
			return updateData, err
		}
		if thirdPartyApkID.Valid {
			update.ThirdPartyApkID = thirdPartyApkID.String
			err := json.Unmarshal([]byte(update.ThirdPartyApkID), &finalOptionTpApkID)
			if err != nil {
				logging.Error(fmt.Sprintf("Unmarshalling failed with data: %s and the error: %s", update.ThirdPartyApkID, err.Error()))
				return updateData, err
			}
		}
		update.Options = append(update.Options, finalOptionTpApkID...)
		updateData = append(updateData, &update)
	}
	return updateData, nil
}

func UpdateOrSetThirdPartyPackageList(ProfileID int, oldValue string, newValue string, user *entities.TMSUser, tid int) error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	// approval_type_id 12 Has been dedicated for the TID update deletion. See NEX-13071
	changeType := 12
	if oldValue != newValue {
		_, err = db.Exec("CALL update_or_set_thirdPartyPackageList_and_insert_into_approval(?,?,?,?,?,?)", ProfileID, oldValue, newValue, user.Username, tid, changeType)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetTIDUpdate(tidId int, tidUpdateId int) (models.TIDUpdateData, error) {
	var updateData = models.TIDUpdateData{}
	db, err := GetDB()
	if err != nil {
		return updateData, err
	}

	var updateID sql.NullInt64
	var packageID sql.NullInt64
	var updateDate sql.NullString
	var thirdPartyApk sql.NullString

	err = db.QueryRow("call get_tid_updates(?,?)", tidUpdateId, tidId).Scan(&updateID, &packageID, &updateDate, &thirdPartyApk)
	if err != nil && err != sql.ErrNoRows {
		return updateData, err
	}

	if updateID.Valid {
		updateData.UpdateID = int(updateID.Int64)
	}

	if packageID.Valid {
		updateData.PackageID = int(packageID.Int64)
	}

	if updateDate.Valid {
		updateData.UpdateDate = updateDate.String
	}

	var finalOptionTpApkID []int

	if thirdPartyApk.Valid {
		updateData.ThirdPartyApkID = thirdPartyApk.String
		err = json.Unmarshal([]byte(updateData.ThirdPartyApkID), &finalOptionTpApkID)
		if err != nil {
			return updateData, err
		}
	}
	updateData.Options = append(updateData.Options, finalOptionTpApkID...)

	return updateData, nil
}

func AddUpdateToTID(tidUpdateId int, tidId int, packageId int, updateDate string, thirdPartyApk string) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	_, err = db.Exec("CALL insert_tid_update(?, ?, ?, ?, ?)", tidUpdateId, tidId, packageId, updateDate, thirdPartyApk)
	if err != nil {
		return err
	}
	return nil
}

func GetTIDUpdateFields(tid, siteId int, apkVersion string, tpApkVersion []string) (*models.TIDUpdateData, error) {

	db, err := GetDB()
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return nil, err
	}

	var updateID int
	res, err := db.Query("SELECT COUNT(*) FROM tid_updates WHERE tid_id = ?", tid)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return nil, err
	}

	defer res.Close()

	for res.Next() {
		err = res.Scan(&updateID)
		if err != nil {
			logger.GetLogger().Error(err.Error())
			return nil, err
		}
	}
	updateID += 1

	var packageID int
	res, err = db.Query("SELECT package_id FROM package WHERE version = ?", apkVersion)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return nil, err
	}
	defer res.Close()

	for res.Next() {
		err = res.Scan(&packageID)
		if err != nil {
			logger.GetLogger().Error(err.Error())
			return nil, err
		}
	}

	if packageID == 0 && apkVersion != "" {
		logger.GetLogger().Error("Incorrect APK version ")
		return nil, errors.New("Incorrect APK version ")
	}

	currentTime := time.Now().Local()
	updateDate := common.FormatTime(currentTime)

	var update models.TIDUpdateData
	update.UpdateID = updateID
	update.UpdateDate = updateDate
	update.PackageID = packageID

	if len(tpApkVersion) != 0 {
		update.ThirdPartyApkID, err = GetThirdPartyApk(tpApkVersion, tid, siteId)
		if err != nil {
			return nil, err
		}
	} else {
		update.ThirdPartyApkID = "[]"
	}

	return &update, nil
}

func GetThirdPartyApk(tpApkVersion []string, tid, siteId int) (string, error) {

	db, err := GetDB()
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return "", err
	}

	if siteId == 0 {
		siteId, err = GetSiteIDForTid(tid)
		if err != nil {
			logger.GetLogger().Information(err.Error())
			return "", errors.New("error fetching SiteId for given TID")
		}
	}

	// Check to see if third party enabled for TID provided
	checkThirdParty, err := GetThirdPartyEnabled(tid, siteId)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return "", err
	}

	if !checkThirdParty {
		logger.GetLogger().Error("Third Party not Enabled for TID")
		return "", errors.New("Third Party not enabled for TID")
	}

	dataValue, err := GetDataValueForTidProfileId(tid, siteId)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return "", errors.New("error fetching Partial Package Name")
	}
	replacer := strings.NewReplacer("[", "", "]", "", `"`, "")
	dataValue = replacer.Replace(dataValue)
	dataValueArray := strings.Split(dataValue, ",")

	finalTPAPKMap := make(map[string]bool, 0)
	for _, dv := range dataValueArray {
		version := 0
		for _, tpAPK := range tpApkVersion {
			if strings.Contains(tpAPK, dv) {
				version++
				if version >= 2 {
					return "", errors.New("multiple apks are selected for single merchant name")
				}
				finalTPAPKMap[tpAPK] = true
			}
		}
	}

	var finalTpApkVersion []string
	for tpApk := range finalTPAPKMap {
		finalTpApkVersion = append(finalTpApkVersion, tpApk)
	}
	var thirdPartyID string
	var thirdPartyIDs []string
	ids := strings.Join(finalTpApkVersion, "','")
	query := fmt.Sprintf(`SELECT apk_id FROM third_party_apks WHERE name IN ('%s')`, ids)
	res, err := db.Query(query)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return "", err
	}
	defer res.Close()

	for res.Next() {
		err = res.Scan(&thirdPartyID)
		if err != nil {
			logger.GetLogger().Error(err.Error())
			return "", err
		}
		thirdPartyIDs = append(thirdPartyIDs, thirdPartyID)
	}

	if len(thirdPartyID) == 0 {
		logger.GetLogger().Error("Third Party APK provided does not exist")
		return "", errors.New("Third Party APK provided does not exist")
	}
	thirdPartyIDStr := strings.Join(thirdPartyIDs, ",")
	thirdPartyIDStr = "[" + thirdPartyIDStr + "]"
	return thirdPartyIDStr, nil
}

func AddThirdPartyAuditHistory(tidProfileId int64, oldValue string, newValue string, tid int, user *entities.TMSUser) error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO approvals (profile_id,data_element_id, change_type, current_value, new_value, created_at,approved_at, approved, created_by,approved_by, tid_id,acquirer)
				   VALUE
				   (?, (SELECT data_element_id from data_element where name ='thirdPartyPackageList'), 5, ?, ?, NOW(),NOW(), 1, ?,?, ?,'NI')`,
		tidProfileId, oldValue, newValue, user.Username, user.Username, tid)
	if err != nil {
		return err
	}
	return nil
}

func UpdateTIDThirdPartyAPks(tidId int, tidUpdateID int, packageIds string) (bool, error) {
	db, err := GetDB()
	if err != nil {
		return false, err
	}

	result, err := db.Exec("update tid_updates set update_date = now(), third_party_apk = (?) where tid_update_id = (?) and tid_id = (?)", packageIds, tidUpdateID, tidId)
	if err != nil {
		return false, err
	}
	RowsAffected, _ := result.RowsAffected()
	if RowsAffected == 0 {
		return false, nil
	}

	_, err = db.Exec(`update tid set flag_status=true, flagged_date=CURRENT_TIMESTAMP where tid_id = (?)`, tidId)
	if err != nil {
		return false, nil
	}

	return true, nil
}

func DeleteUpdateFromTID(tidUpdateId int, tidId int) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM tid_updates where tid_update_id = ? and tid_id = ?", tidUpdateId, tidId)
	if err != nil {
		return err
	}
	return nil
}

func AddTidProfileLink(tidId int, siteId int, tidProfileId int) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	_, err = db.Exec("update tid_site set tid_profile_id = (?) where tid_id = (?) and site_id = (?)", tidProfileId, tidId, siteId)
	if err != nil {
		return err
	}

	return nil
}

func UpdateTIDFlag(tidId int) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	_, err = db.Exec(`update tid set flag_status=true, flagged_date=CURRENT_TIMESTAMP where tid_id = (?)`, tidId)
	if err != nil {
		return err
	}

	return nil
}

func UpdateTIDFlagWithProfileID(profileID int64) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	_, err = db.Exec(`UPDATE tid SET flag_status=true, flagged_date=CURRENT_TIMESTAMP
           					WHERE tid_id = (SELECT tid_id FROM tid_site where tid_profile_id= ?)`, profileID)
	if err != nil {
		return err
	}

	return nil
}

func UpdateProfileNameWithProfileID(value string, profileID int64) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	_, err = db.Exec(`UPDATE profile SET name = ? WHERE profile_id = ?`, value, profileID)
	if err != nil {
		return err
	}

	return nil
}

func CheckTidExistsSiteId(tidId int, siteId int) (bool, error) {
	var tidCount int
	db, err := GetDB()
	if err != nil {
		return false, err
	}

	err = db.QueryRow("select count(*) from tid_site where tid_id = (?) and site_id = (?)", tidId, siteId).Scan(&tidCount)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}

	if tidCount > 0 {
		return true, nil
	}
	return false, nil
}

func CheckTidProfileExistsAndGetSiteId(tidId int) (bool, int64, int, error) {
	var profileId sql.NullInt64
	var siteId sql.NullInt64
	db, err := GetDB()
	if err != nil {
		return false, 0, 0, err
	}

	rows, err := db.Query("select tid_profile_id,site_id from tid_site where tid_id = (?) limit 1", tidId)
	if err != nil {
		return false, 0, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&profileId, &siteId)
	}

	if profileId.Valid && siteId.Valid {
		return true, profileId.Int64, int(siteId.Int64), nil
	}

	return false, profileId.Int64, int(siteId.Int64), nil
}

func GetTidProfileIdForSiteId(tidId int, siteId int) (bool, int64, error) {
	db, err := GetDB()
	if err != nil {
		return false, -1, err
	}

	rows, err := db.Query("select tid_profile_id from tid_site where tid_id = (?) and site_id = (?) limit 1", tidId, siteId)
	if err != nil {
		return false, -1, err
	}
	defer rows.Close()

	var profileId sql.NullInt64
	for rows.Next() {
		rows.Scan(&profileId)
	}

	if profileId.Valid {
		return true, profileId.Int64, nil
	}
	return false, -1, nil
}

func DeleteTidOverride(profileId int) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	logging.Information(fmt.Sprintf("Deleting TID override for profileId '%v'", profileId))
	_, err = tx.Exec("CALL remove_tid_override(?)", profileId)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func GetElementOptionsForId(dgName, elementName string, dataElementDetails map[string]models.DataElementsAndGroup) string {
	if element, ok := dataElementDetails[dgName+"-"+elementName]; ok {
		return element.Options
	}
	return ""
}

func (t TIDData) GetPaddedTidId(tidId int) string {
	return fmt.Sprintf("%08d", tidId)
}

func GetPaddedTidId(tidId int) string {
	return fmt.Sprintf("%08d", tidId)
}

func AcquirerNotPermissable(acquirerName string, user *entities.TMSUser) bool {
	// Find user acquirers to limit search results
	acquirers, err := GetUserAcquirerPermissions(user)
	if err != nil {
		logging.Warning("Error getting user acquirer permissions")
		return false
	}

	acquirerPermissions := strings.Split(acquirers, ",")

	for _, permissibleAcquirer := range acquirerPermissions {
		if permissibleAcquirer == acquirerName {
			return false
		}
	}

	return true
}

// Checks if a given TID exists
// Returns:
//
//	tidExists: true if the tid exists and false if not
//	resultCode: the result code of the function which indicates more information about the result such as if the tid
//			exists elsewhere as a primary tid or if it exists elsewhere as a secondary tid.
//	overrideProfileId: the profileId of the TID's override, this will be -1 if no override profile exists (i.e. when
//			tidExists is false) or if resultCode contains an error result (e.g. DATABASE_CONNECTION_ERROR or
//			DATABASE_QUERY_ERROR).
func CheckThatTidExists(TID int) (tidExists bool, resultCode resultCodes.ResultCode, overrideProfileId int) {
	db, err := GetDB()
	if err != nil {
		logging.Error("An error occurred executing GetDB() - %s", err.Error())
		return false, resultCodes.DATABASE_CONNECTION_ERROR, -1
	}
	rows, err := db.Query("CALL GET_TID_BY_TID(?)", TID)
	if err != nil {
		logging.Error("An error occurred executing procedure GET_TID_BY_TID - %s", err.Error())
		return false, resultCodes.DATABASE_QUERY_ERROR, -1
	}
	defer rows.Close()

	var foundTid int
	var tidType string
	var primaryTid int
	for rows.Next() {
		rows.Scan(&foundTid, &tidType, &primaryTid, &overrideProfileId)
	}

	//The tid does not exist
	if foundTid == 0 {
		tidExists = false
		resultCode = resultCodes.TID_DOES_NOT_EXIST
		return
	} else {
		tidExists = true
		if tidType == "primaryTid" {
			resultCode = resultCodes.TID_NOT_UNIQUE_PRIMARY_TID_DUPLICATE
		} else {
			resultCode = resultCodes.TID_NOT_UNIQUE_SECONDARY_TID_DUPLICATE
		}
		return
	}
}

func UpdateSerialNumberPendingApproval(user *entities.TMSUser, siteProfileId, tid, oldSerialNumber, newSerialNumber string) (int, error) {
	db, err := GetDB()
	if err != nil {
		return 0, err
	}
	acquirer, err := GetTidAcquirer(tid)
	if err != nil {
		return 0, err
	}
	res, err := db.Exec(`INSERT INTO approvals (profile_id, data_element_id, change_type, current_value, new_value, created_at, approved, created_by, tid_id, acquirer)
				   VALUE
				   (?, 1, 6, ?, ?, NOW(), 0, ?, ?, ?)`,
		siteProfileId, oldSerialNumber, newSerialNumber, user.Username, tid, acquirer)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func AddDataGroupsToTidProfiles() error {
	db, err := GetDB()
	if err != nil {
		log.Println("unable to get DB instance", err)
		logging.Error("An error occurred executing GetDB() - %s", err.Error())
		return err
	}
	_, err = db.Exec("CALL add_data_groups_to_tid_profiles")
	if err != nil {
		log.Println("unable to execute add_data_groups_to_tid_profiles", err)
		logging.Error("An error occurred executing procedure add_data_groups_to_tid_profiles - %s", err.Error())
		return err
	}
	return nil
}

func CreateTidoverrideAndSaveNewprofileChange(siteId int, profileTypeId int, username string, changeType int, dataValue string, tidStr string, tidInt, approved int) error {
	db, err := GetDB()
	if err != nil {
		log.Println("unable to get DB instance", err)
		logging.Error("An error occurred executing GetDB() - %s", err.Error())
		return err
	}

	_, err = db.Exec("CALL create_tid_override_and_save_profile_change(?,?,?,?,?,?,?,?)", profileTypeId, username, changeType, dataValue, tidStr, tidInt, approved, siteId)
	if err != nil {
		log.Println("unable to execute CreateTidoverrideAndSaveNewprofileChange", err)
		logging.Error("An error occurred executing procedure create_tid_override_and_save_profile_change - %s", err.Error())
		return err
	}

	return nil
}

func CreateTidOverride(siteId int, profileID int, username string) error {
	db, err := GetDB()
	if err != nil {
		log.Println("unable to get DB instance", err)
		logging.Error("An error occurred executing GetDB() - %s", err.Error())
		return err
	}

	_, err = db.Exec("CALL create_tid_override(?,?,?)", siteId, profileID, username)
	if err != nil {
		log.Println("unable to execute create_tid_override", err)
		logging.Error("An error occurred executing procedure create_tid_override - %s", err.Error())
		return err
	}

	return nil
}

func CheckIfDataElementExistsinTidData(tid, dataElementId int) (bool, int64, error) {
	db, err := GetDB()
	if err != nil {
		return false, -1, err
	}

	rows, err := db.Query("CALL check_data_element_exists(?,?)", tid, dataElementId)
	if err != nil {
		return false, -1, err
	}
	defer rows.Close()

	var DataElememtId sql.NullInt64
	for rows.Next() {
		rows.Scan(&DataElememtId)
	}

	if DataElememtId.Valid {
		return true, DataElememtId.Int64, nil
	}
	return false, -1, nil

}

func GetSiteTidsChangeHistory(db *sql.DB, profileId int) ([]*ProfileChangeHistory, error) {
	var changes = make([]*ProfileChangeHistory, 0)

	rows, err := db.Query("CALL get_site_tids_change_history(?)", profileId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var originalVal sql.NullString
	var updatedValue sql.NullString
	var tidId sql.NullString
	var isPassword sql.NullBool
	var isEncrypted sql.NullBool

	rowCount := 0

	for rows.Next() {
		var changeHistory = &ProfileChangeHistory{}
		err = rows.Scan(
			&changeHistory.Field,
			&changeHistory.ChangeType,
			&originalVal,
			&updatedValue,
			&changeHistory.ChangedBy,
			&changeHistory.ChangedAt,
			&changeHistory.Approved,
			&tidId,
			&isPassword,
			&isEncrypted)

		if originalVal.Valid {
			changeHistory.OriginalValue = originalVal.String
		}

		if updatedValue.Valid {
			changeHistory.ChangeValue = updatedValue.String
		}

		if tidId.Valid {
			changeHistory.TidId = tidId.String
		}

		if isEncrypted.Valid {
			if isEncrypted.Bool {
				if originalVal.Valid {
					changeHistory.OriginalValue, err = crypt.Decrypt(changeHistory.OriginalValue)
					if err != nil {
						return nil, err
					}
				}
				if updatedValue.Valid {
					changeHistory.ChangeValue, err = crypt.Decrypt(changeHistory.ChangeValue)
					if err != nil {
						return nil, err
					}
				}
			}

			changeHistory.IsEncrypted = isEncrypted.Bool
		}

		if isPassword.Valid {
			changeHistory.IsPassword = isPassword.Bool
		}

		if err != nil {
			return nil, err
		}

		changeHistory.RowNo = rowCount
		rowCount++

		changes = append(changes, changeHistory)
	}

	return changes, nil
}

func AddBulkUpdateToTID(updates []models.TIDUpdateData, tidID int) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return err
	}
	defer tx.Rollback()

	// stored procedure call
	stmt, err := tx.Prepare("CALL insert_tid_update_bulk(?, ?, ?, ?, ?)")
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return err
	}
	defer stmt.Close()

	var tidUpdateIds, tidIds, targetPackageIds, updateDates, thirdPartyApks string
	for i, update := range updates {
		if i > 0 {
			tidUpdateIds += ","
			tidIds += ","
			targetPackageIds += ","
			updateDates += ","
			thirdPartyApks += ","
		}

		tidUpdateIds += strconv.Itoa(update.UpdateID)
		tidIds += strconv.Itoa(tidID)
		targetPackageIds += strconv.Itoa(update.PackageID)
		updateDates += update.UpdateDate
		thirdPartyApks += update.ThirdPartyApkID
	}

	// Executed the stored procedure
	_, err = stmt.Exec(tidUpdateIds, tidIds, targetPackageIds, updateDates, thirdPartyApks)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return err
	}

	err = tx.Commit()
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return err
	}

	return nil
}
