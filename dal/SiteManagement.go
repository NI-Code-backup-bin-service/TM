package dal

import (
	sliceHelper "bitbucket.org/network-international/nextgen-libs/nextgen-helpers/TypeComparisonHelpers/SliceComparisonHelpers"
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"html"
	"html/template"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	dataGroup "nextgen-tms-website/DataGroup"
	"nextgen-tms-website/config"
	"nextgen-tms-website/crypt"
	"nextgen-tms-website/entities"
	"nextgen-tms-website/fileServer"
	"nextgen-tms-website/models"
	"nextgen-tms-website/resultCodes"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const (
	ApproveNewElement      = 1
	ApproveRemoveOvverride = 2
	ApproveDelete          = 3
	AcquirerProfile        = -1
	ChainProfile           = -2
	ErrorReturn            = -3
	TidProfile             = -4

	ApproveCreate = 5

	//MySQL server error types
	DuplicateEntryError = 1062

	//Friendly server error texts
	DuplicateUserText = "Duplicate user"
)

type SiteManagementDal struct {
}

type siteUserValidationError struct {
	message string
}

var Permissions map[string]*Permission

type PermissionMap struct {
	ModuleName  string     `json:"ModuleName"`
	Permissions Permission `json:"Permissions"`
}

type Permission struct {
	PermissionSale   string `json:"PermissionSale,omitempty"`
	PermissionVoid   string `json:"PermissionVoid,omitempty"`
	PermissionRefund string `json:"PermissionRefund,omitempty"`
}

func (e *siteUserValidationError) Error() string {
	return e.message
}

type EODAutoData struct {
	Datavalue     string
	Name          string
	DataElementId int
	RowNum        int
}

func GetEODAutoData(siteId int) ([]EODAutoData, error) {
	db, err := GetDB()
	if err != nil {
		logging.Error("Unable to get the database instance : " + err.Error())
		return nil, err
	}

	rows, err := db.Query("Call eod_auto_data_fetch(?)", siteId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var eodAutoData []EODAutoData
	for rows.Next() {
		var e EODAutoData
		err = rows.Scan(&e.Datavalue, &e.Name, &e.DataElementId, &e.RowNum)
		if err != nil {
			return nil, err
		}

		eodAutoData = append(eodAutoData, e)
	}

	return eodAutoData, nil
}

func GetIsOverriden(siteId int, elementId int) (bool, error) {
	db, err := GetDB()
	if err != nil {
		return false, err
	}

	rows, err := db.Query("CALL get_site_element_source(?,?)", siteId, elementId)
	if err != nil {
		logging.Error(err.Error())
		return false, err
	}
	defer rows.Close()

	var source sql.NullString

	for rows.Next() {
		rows.Scan(&source)
	}

	if source.Valid {
		return source.String != "site", nil
	}

	return false, nil
}

func GetSiteIdFromProfileId(profileId int) (int, error) {
	db, err := GetDB()
	if err != nil {
		return ErrorReturn, err
	}
	//Check if profileId is chain,acquirer or site
	var acqId int
	err = db.QueryRow("SELECT acquirer_id FROM chain_profiles WHERE acquirer_id = ?", profileId).Scan(&acqId)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return ErrorReturn, err
		}
		if acqId == 0 {
			var profType int
			err = db.QueryRow("SELECT profile_type_id FROM profile WHERE profile_id = ?", profileId).Scan(&profType)
			if err != nil {
				if err.Error() != "sql: no rows in result set" {
					return ErrorReturn, err
				}
			}
			if profType == 2 { //acquirer type
				return AcquirerProfile, nil
			}
		}
	}
	if acqId != 0 {
		return AcquirerProfile, nil
	}
	//Check if it is chain since it ain't an acquirer
	var chainId int
	err = db.QueryRow("SELECT chain_profile_id FROM chain_profiles WHERE chain_profile_id = ?", profileId).Scan(&chainId)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return ErrorReturn, err
		}
	}
	if chainId != 0 {
		return ChainProfile, nil
	}

	var tidId int
	err = db.QueryRow("select tid_id from tid_site where tid_profile_id = ?", profileId).Scan(&tidId)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return ErrorReturn, err
		}
	}
	if tidId != 0 {
		return TidProfile, nil
	}

	row := db.QueryRow("select site_id from site_profiles sp where sp.profile_id = ?", profileId)
	var id int
	err = row.Scan(&id)
	if err != nil {
		_, _ = logging.Error(err)
		return ErrorReturn, err
	}
	return id, nil
}

func GetTIdSiteIdFromProfileId(profileId int) (int, int, error) {
	db, err := GetDB()
	if err != nil {
		return ErrorReturn, ErrorReturn, err
	}

	row := db.QueryRow("select tid_id,site_id from tid_site_profiles tsp where tsp.profile_id = ?", profileId)
	var tId, siteId int
	err = row.Scan(&tId, &siteId)
	if err != nil {
		_, _ = logging.Error(err)
		return ErrorReturn, ErrorReturn, err
	}
	return tId, siteId, nil
}

func GetSubmodules(modules []string) []string {
	var newModules = make([]string, 0)
	for _, module := range modules {
		switch module {
		case "alipay":
			newModules = append(newModules, "alipaySale", "alipayVoid", "alipayRefund")
		case "touchPoints":
			newModules = append(newModules, "touchpointRedeem", "touchpointVoid", "touchpointCheckBalance")
		case "preAuth":
			newModules = append(newModules, "preAuthSale", "preAuthCompletion", "preAuthCancel")
		case "visaQr":
			newModules = append(newModules, "visaQrSale", "visaQrRefund")
		case "mastercardQr":
			newModules = append(newModules, "mastercardQrSale")
		case "meezaQR":
			newModules = append(newModules, "meezaQrSale", "meezaQrRtpSale", "meezaQrRefund")
		case "weChatPay":
			newModules = append(newModules, "weChatPaySale", "weChatPayVoid", "weChatPayRefund")
		case "terraPay":
			newModules = append(newModules, "terraPaySale", "terraPayRefund")
		default:
			//If future module permissions are added via feature script to avoid TMS upgrades we can
			//Loop through the PermissionMap and append directly from here
			//Currently only supports Sales / Refunds / Voids - Single entries would just be the module name itself
			//E.g. NPCI is only under "PermissionSale" as "NPCI" (essentially calling the else statement)
			if Permissions[module] != nil {
				if Permissions[module].PermissionSale != "" {
					newModules = append(newModules, Permissions[module].PermissionSale)
				}

				if Permissions[module].PermissionRefund != "" {
					newModules = append(newModules, Permissions[module].PermissionRefund)
				}

				if Permissions[module].PermissionVoid != "" {
					newModules = append(newModules, Permissions[module].PermissionVoid)
				}
			} else {
				newModules = append(newModules, module)
			}
		}
	}

	return newModules
}

func GetDataElementValue(profileId int, dataElementName, dataGroupName string) (string, error) {
	db, err := GetDB()
	if err != nil {
		return "", err
	}

	var disEncrypted sql.NullBool
	isEncrypted := false
	isPassword := false

	res, err := db.Query("Call get_data_element_value(?,?,?)", profileId, dataElementName, dataGroupName)
	if err != nil {
		return "", err
	}
	defer res.Close()

	var currentValue string
	for res.Next() {
		_ = res.Scan(&currentValue, &disEncrypted, &isPassword)
	}

	if disEncrypted.Valid && disEncrypted.Bool {
		isEncrypted = true
	}

	if isEncrypted {
		currentValue, err = crypt.Decrypt(currentValue)
		if err != nil {
			return "", err
		}
	}

	return currentValue, nil
}

func GetTIDdatavalueFromSite(siteId int, dataElementName, dataGroupName string) ([]string, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	var dataValues []string
	var isEncrypted sql.NullBool

	res, err := db.Query("Call get_tid_datavalue_from_site(?,?,?)", siteId, dataElementName, dataGroupName)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var dataValue string
	for res.Next() {
		err = res.Scan(&dataValue, &isEncrypted)
		if err != nil {
			return nil, err
		}

		if isEncrypted.Valid && isEncrypted.Bool {
			dataValue, err = crypt.Decrypt(dataValue)
			if err != nil {
				return nil, err
			}
		}
		dataValues = append(dataValues, dataValue)
	}

	return dataValues, nil
}

func getElementType(elementID int) string {
	db, err := GetDB()
	if err != nil {
		return ""
	}

	rows, err := db.Query("SELECT datatype FROM data_element WHERE data_element_id = ?", elementID)
	if err != nil {
		return ""
	}
	defer rows.Close()

	dataType := ""
	for rows.Next() {
		err := rows.Scan(&dataType)
		if err != nil {
			return ""
		}
	}

	return dataType
}

func isElementPassword(elementID int, is_password *bool) {
	db, err := GetDB()
	if err != nil {
		return
	}

	*is_password = false

	res, err := db.Query("SELECT is_password FROM data_element WHERE data_element_id = ? LIMIT 1", elementID)
	if err != nil {
		logging.Error(err)
		return
	}
	defer res.Close()

	for res.Next() {
		err = res.Scan(is_password)
	}
}

// SaveElementData Saves a new unapproved update to the element table if the element has changed
func SaveElementData(profileId int, elementID int, value string, user string, approved int, overriden int) (err error) {
	db, err := GetDB()
	if err != nil {
		return err
	}

	isEncrypted := false
	isPassword := false

	res, err := db.Query("Call get_element_value(?,?)", profileId, elementID)
	if err != nil {
		println("TODO: Handle this: " + err.Error())
		return err
	}
	defer res.Close()

	var currentValue string
	for res.Next() {
		err = res.Scan(&currentValue, &isEncrypted, &isPassword)
		if err != nil {
			logging.Error("SaveElementData res.Scan error : " + err.Error())
		}
	}

	//This is to handle the existing values stored as clear value
	clearCurrentValue := currentValue
	if isEncrypted && currentValue != "" {
		currentValue, err = crypt.Decrypt(currentValue)
		if err != nil {
			logging.Error("SaveElementData crypt.Decrypt error : " + err.Error())
			currentValue = clearCurrentValue
		}
	}

	storeEncrypted := false

	if currentValue != value {
		if isEncrypted && crypt.UseEncryption {
			value = crypt.Encrypt(value)
			storeEncrypted = true
		}

		_, err = db.Exec("Call store_profile_data(?,?,?,?,?,?,?)", profileId, elementID, value, user, approved, overriden, storeEncrypted)
		if err != nil {
			return err
		}
	}

	return nil
}

func AddDataGroupToTidProfile(profileID int) error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	// this is just call to add data-group at TID level where TID has already overriden.
	rows, err := db.Query("CALL add_site_data_groups_to_tid_profiles(?)", profileID)
	if err != nil {
		log.Println("unable to query add_data_groups_to_tid_profiles", err)
		logging.Error("An error occurred querying procedure add_data_groups_to_tid_profiles - %s", err.Error())
		return err
	}
	defer rows.Close()
	return nil
}

func SaveUnapprovedElement(profileId int, elementID int, value string, user string, overridden int, approvalType int) (err error) {
	db, err := GetDB()
	if err != nil {
		return err
	}

	var isEncrypted, isPassword bool
	var currentValue string

	res, err := db.Query("Call get_element_value(?,?)", profileId, elementID)
	if err != nil {
		return err
	}
	defer res.Close()
	for res.Next() {
		if err = res.Scan(&currentValue, &isEncrypted, &isPassword); err != nil {
			logging.Error("SaveUnapprovedElement res.Scan error : " + err.Error())
		}
	}

	//This is to handle the existing values stored as clear value
	clearCurrentValue := currentValue
	if isEncrypted && currentValue != "" {
		currentValue, err = crypt.Decrypt(currentValue)
		if err != nil {
			logging.Error("SaveUnapprovedElement : Error while Decrypting : " + err.Error())
			currentValue = clearCurrentValue
		}
	}
	//@NEX-12567 below code addedd for [Flag Specific is un checked within single update in Site Configuration Users and Fraud history]
	var userFraudDataElementId int
	if currentValue == "" {
		userDataElementID, err := GetDataElementByName("core", "users")
		if err != nil {
			logging.Error(err.Error())
		}
		if userDataElementID == elementID {
			currentValue, err = getUsersFraudDetailsForApproval(profileId, "users", false)
			if err != nil {
				return err
			}
			userFraudDataElementId = userDataElementID
		}
		//if users not found at site level then check for tid level users
		if currentValue == "[]" && userDataElementID == elementID {
			currentValue, err = getUsersFraudDetailsForApproval(profileId, "users", true)
			if err != nil {
				return err
			}
			userFraudDataElementId = userDataElementID

		}
		fraudDataElementID, err := GetDataElementByName("core", "fraud")
		if err != nil {
			return err
		}
		if fraudDataElementID == elementID {
			currentValue, err = getUsersFraudDetailsForApproval(profileId, "fraud", false)
			if err != nil {
				return err
			}
			userFraudDataElementId = fraudDataElementID
		}
	}

	if approvalType == ApproveRemoveOvverride {
		if isEncrypted && crypt.UseEncryption {
			value = crypt.Encrypt(value)
			isEncrypted = true
		}
		_, err = db.Exec("Call save_pending_element_change(?,?,?,?,?,?,?,?,?)", profileId, elementID, 4, value, user, isPassword, isEncrypted, currentValue, userFraudDataElementId)
		if err != nil {
			return err
		}
	} else {
		if currentValue != value {
			var changeType int
			if approvalType == ApproveNewElement {
				if overridden == 1 {
					changeType = 2
				} else {
					changeType = 1
				}
			}
			if isEncrypted && crypt.UseEncryption {
				value = crypt.Encrypt(value)
				isEncrypted = true
			}
			// NEX-10379 Skipping the data-element to which have default value empty or false while storing in profile data.
			if currentValue == "" && value == "[]" || currentValue == "" && value == "false" {
				return nil
			}
			_, err = db.Exec("Call save_pending_element_change(?,?,?,?,?,?,?,?,?)", profileId, elementID, changeType, value, user, isPassword, isEncrypted, currentValue, userFraudDataElementId)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func RecordSiteToHistory(profileId int, value string, user string, changeType int, approved int) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	mid, err := GetMerchantIdFromProfileId(profileId)
	if err != nil {
		return err
	}

	_, err = db.Exec("Call record_site_to_history(?,?,?,?,?,?)", profileId, changeType, value, user, approved, mid)
	if err != nil {
		return err
	}

	return nil
}

func GetMerchantIdFromProfileId(profileId int) (string, error) {
	db, err := GetDB()
	if err != nil {
		return "", err
	}

	rows, err := db.Query("select pd.datavalue as 'merchant_id' from profile_data pd where pd.profile_id = (?) AND pd.data_element_id = 1  "+
		"AND pd.version = (SELECT MAX(d.version) FROM profile_data d WHERE d.data_element_id = 1 AND d.profile_id = pd.profile_id AND d.approved = 1)", profileId)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var merchantId string
	for rows.Next() {
		rows.Scan(&merchantId)
	}
	if err != nil {
		return "", err
	}

	return merchantId, nil
}

func SaveTidProfileChange(profileId int, tidID string, value string, user string, changeType int, approved int) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	_, err = db.Exec("Call save_pending_profile_change(?,?,?,?,?,?)", profileId, changeType, value, user, tidID, approved)
	if err != nil {
		return err
	}

	return nil
}

func RemoveOverride(siteId int, elementId int) error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	if _, err := db.Exec("CALL remove_override(?,?)", siteId, elementId); err != nil {
		logging.Error(err.Error())
		return err
	}

	return nil
}

func GetSiteFromProfile(profileID int) (int, error) {
	db, err := GetDB()
	if err != nil {
		return -1, err
	}
	rows, err := db.Query("Call get_site_id_from_profile_id(?)", profileID)
	if err != nil {
		return -1, err
	}
	defer rows.Close()
	var siteId int
	for rows.Next() {
		rows.Scan(&siteId)
	}
	if err != nil {
		return -1, err
	}
	return siteId, nil
}

func GetProfileFromSite(siteId int) (int, error) {
	db, err := GetDB()
	if err != nil {
		return -1, err
	}

	rows, err := db.Query("SELECT sp.profile_Id FROM site_profiles sp "+
		"LEFT JOIN profile p ON p.profile_id = sp.profile_id "+
		"LEFT JOIN profile_type pt ON pt.profile_type_id = p.profile_type_id "+
		"WHERE sp.site_id = ? AND pt.priority = 2 limit 1", siteId)
	if err != nil {
		return -1, err
	}
	defer rows.Close()

	var profileID int
	for rows.Next() {
		rows.Scan(&profileID)
	}

	if err != nil {
		return -1, err
	}

	return profileID, nil
}

func GetTypeForProfile(profileID int) (profileType string, err error) {
	db, err := GetDB()
	if err != nil {
		return "", err
	}

	if err := db.QueryRow("SELECT pt.name FROM profile AS p LEFT JOIN profile_type pt ON p.profile_type_id = pt.profile_type_id WHERE p.profile_id = ?", profileID).Scan(&profileType); err != nil {
		return "", err
	}

	return profileType, nil
}

func SaveNewSite(name string, version int, user string, chainId int) (int64, int64, error) {
	profileId, err := SaveNewProfile("site", name, version, user)
	if err != nil {
		return 0, 0, err
	}

	db, err := GetDB()
	if err != nil {
		return -1, -1, err
	}

	var siteId int64
	rows, err := db.Query("Call site_store(?,?,?)", -1, version, user)
	if err != nil {
		return -1, -1, err
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&siteId)
	}

	_, err = db.Exec("Call site_profiles_store(?,?,?,?,?)", -1, siteId, profileId, version, user)
	if err != nil {
		return -1, -1, err
	}

	// Get the acquirer id of the acquirer the selected chain is configured for
	acquirerId, err := GetAcquirerIdFromChainId(chainId)
	if err != nil {
		return -1, -1, err
	}

	_, err = db.Exec("Call site_profiles_store(?,?,?,?,?)", -1, siteId, acquirerId, version, user)
	if err != nil {
		return -1, -1, err
	}

	_, err = db.Exec("Call site_profiles_store(?,?,?,?,?)", -1, siteId, chainId, version, user)

	if err != nil {
		return -1, -1, err
	}

	_, err = db.Exec("Call site_profiles_store(?,?,?,?,?)", -1, siteId, 1, version, user)

	if err != nil {
		return -1, -1, err
	}

	return profileId, siteId, err
}

func DeleteSite(profileID string) error {

	db, err := GetDB()
	if err != nil {
		return err
	}

	_, err = db.Exec("CALL delete_site_update_serial_number(?)", profileID)
	if err != nil {
		return err
	}

	return nil
}

// Removes all Site and TID velocity limits from a Site
func removeSiteVelocityLimits(siteId int) error {

	// Delete the Site's non-scheme velocity limits (level 3)
	err := DeleteSiteVelocityLimits(siteId, 3, -1)
	if err != nil {
		return err
	}

	// Delete the Site's scheme velocity limits (level 1)
	err = DeleteSiteVelocityLimits(siteId, 1, -1)
	if err != nil {
		return err
	}

	return nil
}

// Obtains a list of TID IDs for a given site
func GetSiteTidIDs(profileID int, siteIds []string) ([]string, int, int, error) {
	var siteTidIDs []string
	db, err := GetDB()
	if err != nil {
		return nil, 0, 0, err
	}

	var siteId int
	err = db.QueryRow("SELECT site_id FROM site_profiles WHERE profile_id = ?", profileID).Scan(&siteId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, 0, 0, fmt.Errorf("no site found for profile %d", profileID)
		}
	}

	query := "Call get_site_tid_ids(?)"
	if len(siteIds) > 0 {
		query = "SELECT count(t.tid_id) FROM tid t LEFT JOIN tid_site ts ON ts.tid_id = t.tid_id  WHERE "
		query += "ts.site_id IN (?" + strings.Repeat(",?", len(siteIds)-1) + ")"

		queryArgs := make([]interface{}, len(siteIds))
		for i, tid := range siteIds {
			queryArgs[i] = tid
		}

		var tidID int
		err := db.QueryRow(query, queryArgs...).Scan(&tidID)
		if err != nil {
			return nil, siteId, 0, err
		}

		return siteTidIDs, siteId, tidID, err
	}

	rows, err := db.Query(query, siteId)
	if err != nil {
		return nil, siteId, 0, err
	}
	defer rows.Close()

	var tidID string
	for rows.Next() {
		err = rows.Scan(&tidID)
		if err != nil {
			return nil, siteId, 0, err
		}

		siteTidIDs = append(siteTidIDs, tidID)
	}

	return siteTidIDs, siteId, 0, nil
}

func GetChainSiteIDsAndTids(chainId int) (int, int, error) {
	db, err := GetDB()
	if err != nil {
		return 0, 0, err
	}

	var tidCount, siteCount int
	err = db.QueryRow("SELECT COUNT(tid_id) FROM tid_site ts, site_profiles sp WHERE ts.site_id=sp.site_id AND sp.profile_id IN (SELECT profile_id FROM profile WHERE profile_type_id=(SELECT profile_type_id FROM profile_type WHERE priority = 3 ) AND sp.profile_id = ?) GROUP BY sp.profile_id", chainId).Scan(&tidCount)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, 0, err
	}

	err = db.QueryRow("SELECT COUNT(site_id) FROM site_profiles WHERE profile_id = ?", chainId).Scan(&siteCount)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, 0, err
	}

	return siteCount, tidCount, nil
}

func DeleteChain(chainId int) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	_, err = db.Exec("CALL chain(?)", chainId)

	return err
}

func DeleteAcquirer(acquirerId int) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	_, err = db.Exec("CALL delete_acquirer(?)", acquirerId)

	return err
}

func SaveNewChain(profileType string, name string, version int, user string, acquirerId int) (int64, error) {
	profileId, err := SaveNewProfile(profileType, name, version, user)
	if err != nil {
		return 0, err
	}

	db, err := GetDB()
	if err != nil {
		return -1, err
	}

	_, err = db.Exec("insert into chain_profiles(chain_profile_id, acquirer_id) values (?,?)", profileId, acquirerId)
	if err != nil {
		return -1, err
	}

	return profileId, err
}

// Get profile_id through profile Type
func GetProfileTypeId(profile_type string, name string, version int, user string) (int, error) {
	db, err := GetDB()
	if err != nil {
		return -1, err
	}

	rows, err := db.Query("CALL get_profile_type_by_name(?)", profile_type)
	if err != nil {
		return -1, err
	}
	defer rows.Close()

	var type_id int
	for rows.Next() {
		rows.Scan(&type_id)
	}

	return type_id, nil
}

func SaveNewProfile(profile_type string, name string, version int, user string) (int64, error) {
	db, err := GetDB()
	if err != nil {
		return -1, err
	}

	rows, err := db.Query("CALL get_profile_type_by_name(?)", profile_type)
	if err != nil {
		return -1, err
	}
	defer rows.Close()

	var type_id int
	for rows.Next() {
		rows.Scan(&type_id)
	}

	_, err = db.Exec("Call profile_store(?,?,?,?,?)", -1, type_id, name, version, user)

	if err != nil {
		return -1, err
	}

	var profileID int64
	// Horrifc workaround for the go db driver not supporting out params or stored procs in genreal.
	rows, err = db.Query("SELECT profile_id FROM profile WHERE name = ?", name)
	if err != nil {
		return -1, err
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&profileID)
	}
	return profileID, err
}

func AddDataGroupAndDataElement(duplicateChainProfileId, parentChainProfileId int, createdBy string) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	_, err = db.Exec("Call save_duplicate_chain_profile(?, ?, ?)", duplicateChainProfileId, parentChainProfileId, createdBy)
	if err != nil {
		return err
	}

	return nil
}

func AddDataGroupToProfile(profileId int, datagroupId string, user string) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	_, err = db.Exec("Call profile_data_group_store(?, ?, ?, ?, ?)", -1, profileId, datagroupId, -1, user)

	return err

}

// Retrieves the card/QR schemes available for the passed in TID/Site
func GetAvailableSchemesForSiteId(siteId int) (map[int]string, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	qrSchemes, err := getQrBasedSchemesBySiteId(db, siteId)
	if err != nil {
		return nil, err
	}

	cardSchemes, err := getCardBasedSchemesBySiteId(db, siteId)
	if err != nil {
		return nil, err
	}

	// Join the qr and card schemes together
	allSchemes := make(map[int]string, 0)
	for k, v := range cardSchemes {
		allSchemes[k] = v
	}

	for k, v := range qrSchemes {
		allSchemes[k] = v
	}

	return allSchemes, nil
}

// Retrieves an array of the card based schemes available to the given profile
func getCardBasedSchemesBySiteId(db *sql.DB, siteId int) (map[int]string, error) {
	schemes := make(map[int]string, 0)
	var configuredCardSchemes []string

	rows, err := callGetCardSchemesBySiteId(db, siteId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cardDefinitions string
	for rows.Next() {
		err = rows.Scan(&cardDefinitions)
		if err != nil {
			return nil, err
		}
	}

	if cardDefinitions == "" {
		logging.Error("No card definitions for siteId = " + strconv.Itoa(siteId))
		return nil, nil
	}
	var definitions []map[string]interface{}
	err = json.Unmarshal([]byte(cardDefinitions), &definitions)
	if err != nil {
		return nil, err
	}

	for _, definition := range definitions {
		cardName := definition["cardName"]
		configuredCardSchemes = append(configuredCardSchemes, strings.ToUpper(cardName.(string)))
	}

	allSchemesRows, err := callGetAllSchemes(db)
	if err != nil {
		return nil, err
	}
	defer allSchemesRows.Close()

	for allSchemesRows.Next() {
		var schemeId int
		var schemeName string
		err = allSchemesRows.Scan(&schemeId, &schemeName)
		if err != nil {
			return nil, err
		}

		if sliceHelper.SlicesOfStringContains(configuredCardSchemes, strings.ToUpper(schemeName)) {
			schemes[schemeId] = schemeName
		}
	}

	return schemes, nil
}

// Retrieves an array of the QR based schemes available to the given profile
func getQrBasedSchemesBySiteId(db *sql.DB, siteId int) (map[int]string, error) {
	qrSchemes := make(map[int]string, 0)

	rows, err := callQrSchemesFromSiteId(db, siteId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var schemeId int
		var qrScheme string
		err = rows.Scan(&schemeId, &qrScheme)
		if err != nil {
			return nil, err
		}
		qrSchemes[schemeId] = qrScheme
	}
	return qrSchemes, nil
}

// DB call to retrieve all schemes and their IDs
func callGetAllSchemes(db *sql.DB) (*sql.Rows, error) {
	return db.Query("SELECT scheme_id, scheme_name FROM schemes")
}

// DB call to retrieve all card schemes for a given profile
func callGetCardSchemesBySiteId(db *sql.DB, siteId int) (*sql.Rows, error) {
	sql := `SELECT pde.value
			FROM site_profiles sp
			INNER JOIN profile_data_elements pde on sp.profile_id = pde.profile_id
			WHERE sp.site_id = ?
				AND pde.data_group_name = 'emv'
				AND pde.data_element_name = 'cardDefinitions'
			ORDER BY pde.profile_type_priority asc
			LIMIT 1`
	return db.Query(sql, siteId)
}

// DB call to retrieve all QR schemes for a given profile
func callQrSchemesFromSiteId(db *sql.DB, siteId int) (*sql.Rows, error) {
	return db.Query(`CALL get_qr_schemes_from_site_id(?)`, siteId)
}

func GetSiteData(siteId int, profileId int, pageSize int, pageNumber int, tidSearchTerm string) ([]*DataGroup, []*DataGroup, []*DataGroup, []*DataGroup, []*TIDData, []*PackageData, []*DataGroup, map[string]int, TidPagination, error) {
	db, err := GetDB()
	var pagination TidPagination
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, pagination, err
	}

	rows, err := db.Query("CALL get_site_data(?, ?)", profileId, siteId)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, pagination, err
	}
	defer rows.Close()

	var siteGroups MappedDataGroups
	var chainGroups MappedDataGroups
	var acquirerGroups MappedDataGroups
	var globalGroups MappedDataGroups
	var tidGroups MappedDataGroups
	var defaultTidGroups MappedDataGroups
	var emptyDependencies []SiteData
	tidOverRideDataElement, err := GetTIDOverRideDataElement()
	if err != nil {
		logging.Error(fmt.Sprintf("Error retrieving data element tid over ride enable %s", err.Error()))
		return nil, nil, nil, nil, nil, nil, nil, nil, pagination, err
	}

	preAllocated := make(map[int]bool)
	activeGroups := make(map[string]bool, 0)
	profiles := make(map[string]int, 0)
	for rows.Next() {
		var site SiteData
		var siteOptions string

		err = rows.Scan(&site.TIDID, &site.DataGroupID, &site.DataGroup, &site.DataGroupDisplayName, &site.DataElementID,
			&site.Name, &site.Tooltip, &site.Source, &site.OverridePriority, &site.DataValue, &site.Overriden, &site.IsNotOverridable, &site.DataType, &site.IsAllowEmpty, &site.MaxLength,
			&site.ValidationExpression, &site.ValidationMessage, &site.FrontEndValidate, &siteOptions, &site.DisplayName, &site.SortOrderInGroup, &site.IsPassword, &site.IsEncrypted, &site.FileMaxSize, &site.FileMinRatio, &site.FileMaxRatio, &site.TIDOverridable, &site.IsReadOnlyAtCreation, &site.IsRequiredAtAcquireLevel, &site.IsRequiredAtChainLevel, &site.ProfileID, &site.ProfileName)
		if err != nil {
			return nil, nil, nil, nil, nil, nil, nil, nil, pagination, err
		}

		if site.IsEncrypted.Valid && site.IsEncrypted.Bool && site.DataValue.Valid {
			site.DataValue.String, err = crypt.Decrypt(site.DataValue.String)
			if err != nil {
				return nil, nil, nil, nil, nil, nil, nil, nil, pagination, err
			}
		}

		activeGroups[site.DataGroup] = true
		profiles[site.ProfileName] = site.ProfileID
		site.Options, site.OptionSelectable = BuildOptionsData(siteOptions, site.DataValue.String, site.DataGroup, site.Name, profileId)

		AddDefaultTidGroup(site, &tidGroups)

		if !site.Source.Valid || !site.DataValue.Valid {
			site.DataValue.String = ""
			if !tidOverRideDataElement[strconv.Itoa(site.DataElementID)] {
				emptyDependencies = append(emptyDependencies, site)
			}
		} else {
			removeElement := func(mdg MappedDataGroups, dgName, deName string) {
				if _, groupExists := mdg[dgName]; groupExists {
					delete(mdg[dgName].DataElements, deName)
					if len(mdg[dgName].DataElements) == 0 {
						delete(mdg, dgName)
					}
				}
			}
			elementPresent := func(mdg MappedDataGroups, dgName, deName string) bool {
				if _, groupExists := mdg[dgName]; groupExists {
					_, elementExists := mdg[dgName].DataElements[deName]
					return elementExists
				} else {
					return false
				}
			}
			switch site.Source.String {
			case "site":
				// If it exists at other levels then remove it from their levels
				removeElement(chainGroups, site.DataGroup, site.Name)
				removeElement(acquirerGroups, site.DataGroup, site.Name)
				removeElement(globalGroups, site.DataGroup, site.Name)
				if !tidOverRideDataElement[strconv.Itoa(site.DataElementID)] {
					siteGroups.AddElementFromSiteData(site)
				}
				preAllocated[site.DataElementID] = true
			case "chain":
				if !elementPresent(siteGroups, site.DataGroup, site.Name) {
					removeElement(acquirerGroups, site.DataGroup, site.Name)
					removeElement(globalGroups, site.DataGroup, site.Name)
					chainGroups.AddElementFromSiteData(site)
					preAllocated[site.DataElementID] = true
				}
			case "acquirer":
				if !elementPresent(siteGroups, site.DataGroup, site.Name) && !elementPresent(chainGroups, site.DataGroup, site.Name) {
					removeElement(globalGroups, site.DataGroup, site.Name)
					acquirerGroups.AddElementFromSiteData(site)
					preAllocated[site.DataElementID] = true
				}
			case "global":
				if !elementPresent(siteGroups, site.DataGroup, site.Name) && !elementPresent(chainGroups, site.DataGroup, site.Name) && !elementPresent(acquirerGroups, site.DataGroup, site.Name) {
					globalGroups.AddElementFromSiteData(site)
					preAllocated[site.DataElementID] = true
				}
			}
		}
	}

	defaultTidGroups = tidGroups

	for i := range emptyDependencies {
		if preAllocated[emptyDependencies[i].DataElementID] == false {
			siteGroups.AddElementFromSiteData(emptyDependencies[i])
		}
	}

	packages, err := GetPackages()
	if err != nil {
		logging.Error(err.Error())
		return nil, nil, nil, nil, nil, nil, nil, nil, pagination, err
	}

	tids, pagination, err := GetTids(db, siteId, activeGroups, pageSize, pageNumber, tidSearchTerm)
	if err != nil {
		logging.Error(err.Error())
		return nil, nil, nil, nil, nil, nil, nil, nil, pagination, err
	}

	// If the tid hasn't been overridden it'll have no profile to draw data from
	// so create a default set of fields using site data
	for _, tid := range tids {
		if len(tid.TIDProfileGroups) == 0 {
			tid.TIDProfileGroups = tidGroups.toArrays()
		}
	}

	return siteGroups.toArrays(), chainGroups.toArrays(), acquirerGroups.toArrays(), globalGroups.toArrays(), tids, packages, defaultTidGroups.toArrays(), profiles, pagination, nil
}

func GetTIDSiteData(siteId, profileId int) ([]*DataGroup, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("CALL get_tid_default_data_group(?, ?)", profileId, siteId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var defaultTidGroups MappedDataGroups
	for rows.Next() {
		var site SiteData
		var siteOptions string

		err = rows.Scan(&site.TIDID, &site.DataGroupID, &site.DataGroup, &site.DataGroupDisplayName, &site.DataElementID,
			&site.Name, &site.Tooltip, &site.Source, &site.OverridePriority, &site.DataValue, &site.Overriden, &site.IsNotOverridable, &site.DataType, &site.IsAllowEmpty, &site.MaxLength,
			&site.ValidationExpression, &site.ValidationMessage, &site.FrontEndValidate, &siteOptions, &site.DisplayName, &site.SortOrderInGroup, &site.IsPassword, &site.IsEncrypted, &site.FileMaxSize, &site.FileMinRatio, &site.FileMaxRatio, &site.TIDOverridable, &site.IsReadOnlyAtCreation, &site.IsRequiredAtAcquireLevel, &site.IsRequiredAtChainLevel, &site.ProfileID, &site.ProfileName)
		if err != nil {
			return nil, err
		}

		if site.IsEncrypted.Valid && site.IsEncrypted.Bool && site.DataValue.Valid {
			site.DataValue.String, err = crypt.Decrypt(site.DataValue.String)
			if err != nil {
				return nil, err
			}
		}

		site.Options, site.OptionSelectable = BuildOptionsData(siteOptions, site.DataValue.String, site.DataGroup, site.Name, profileId)

		AddDefaultTidGroup(site, &defaultTidGroups)
	}

	return defaultTidGroups.toArrays(), nil
}

func GetSiteGroupsData(siteId int, profileId int) ([]*DataGroup, []*DataGroup, []*DataGroup, []*DataGroup, error) {
	db, err := GetDB()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	rows, err := db.Query("CALL get_site_data(?, ?)", profileId, siteId)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	defer rows.Close()

	var siteGroups MappedDataGroups
	var chainGroups MappedDataGroups
	var acquirerGroups MappedDataGroups
	var globalGroups MappedDataGroups
	var emptyDependencies []SiteData
	tidOverRideDataElement, err := GetTIDOverRideDataElement()
	if err != nil {
		logging.Error(fmt.Sprintf("Error retrieving data element tid over ride enable %s", err.Error()))
		return nil, nil, nil, nil, err
	}

	preAllocated := make(map[int]bool)
	for rows.Next() {
		var site SiteData
		var siteOptions string

		err = rows.Scan(&site.TIDID, &site.DataGroupID, &site.DataGroup, &site.DataGroupDisplayName, &site.DataElementID,
			&site.Name, &site.Tooltip, &site.Source, &site.OverridePriority, &site.DataValue, &site.Overriden, &site.IsNotOverridable, &site.DataType, &site.IsAllowEmpty, &site.MaxLength,
			&site.ValidationExpression, &site.ValidationMessage, &site.FrontEndValidate, &siteOptions, &site.DisplayName, &site.SortOrderInGroup, &site.IsPassword, &site.IsEncrypted, &site.FileMaxSize, &site.FileMinRatio, &site.FileMaxRatio, &site.TIDOverridable, &site.IsReadOnlyAtCreation, &site.IsRequiredAtAcquireLevel, &site.IsRequiredAtChainLevel, &site.ProfileID, &site.ProfileName)
		if err != nil {
			return nil, nil, nil, nil, err
		}

		if site.IsEncrypted.Valid && site.IsEncrypted.Bool && site.DataValue.Valid {
			site.DataValue.String, err = crypt.Decrypt(site.DataValue.String)
			if err != nil {
				return nil, nil, nil, nil, err
			}
		}
		site.Options, site.OptionSelectable = BuildOptionsData(siteOptions, site.DataValue.String, site.DataGroup, site.Name, profileId)
		if !site.Source.Valid || !site.DataValue.Valid {
			site.DataValue.String = ""
			if !tidOverRideDataElement[strconv.Itoa(site.DataElementID)] {
				emptyDependencies = append(emptyDependencies, site)
			}
		} else {
			removeElement := func(mdg MappedDataGroups, dgName, deName string) {
				if _, groupExists := mdg[dgName]; groupExists {
					delete(mdg[dgName].DataElements, deName)
					if len(mdg[dgName].DataElements) == 0 {
						delete(mdg, dgName)
					}
				}
			}
			elementPresent := func(mdg MappedDataGroups, dgName, deName string) bool {
				if _, groupExists := mdg[dgName]; groupExists {
					_, elementExists := mdg[dgName].DataElements[deName]
					return elementExists
				} else {
					return false
				}
			}
			switch site.Source.String {
			case "site":
				// If it exists at other levels then remove it from their levels
				removeElement(chainGroups, site.DataGroup, site.Name)
				removeElement(acquirerGroups, site.DataGroup, site.Name)
				removeElement(globalGroups, site.DataGroup, site.Name)
				if !tidOverRideDataElement[strconv.Itoa(site.DataElementID)] {
					siteGroups.AddElementFromSiteData(site)
				}
				preAllocated[site.DataElementID] = true
			case "chain":
				if !elementPresent(siteGroups, site.DataGroup, site.Name) {
					removeElement(acquirerGroups, site.DataGroup, site.Name)
					removeElement(globalGroups, site.DataGroup, site.Name)
					chainGroups.AddElementFromSiteData(site)
					preAllocated[site.DataElementID] = true
				}
			case "acquirer":
				if !elementPresent(siteGroups, site.DataGroup, site.Name) && !elementPresent(chainGroups, site.DataGroup, site.Name) {
					removeElement(globalGroups, site.DataGroup, site.Name)
					acquirerGroups.AddElementFromSiteData(site)
					preAllocated[site.DataElementID] = true
				}
			case "global":
				if !elementPresent(siteGroups, site.DataGroup, site.Name) && !elementPresent(chainGroups, site.DataGroup, site.Name) && !elementPresent(acquirerGroups, site.DataGroup, site.Name) {
					globalGroups.AddElementFromSiteData(site)
					preAllocated[site.DataElementID] = true
				}
			}
		}
	}

	for i := range emptyDependencies {
		if preAllocated[emptyDependencies[i].DataElementID] == false {
			siteGroups.AddElementFromSiteData(emptyDependencies[i])
		}
	}
	return siteGroups.toArrays(), chainGroups.toArrays(), acquirerGroups.toArrays(), globalGroups.toArrays(), nil
}

func computeSelectionOptions(profileId int, dataGroupName, dataElementName string) (optionsString string, computedField bool) {
	activeSchemesOptions := func(profileId int) string {
		elementId, err := GetDataElementByName("dualCurrency", "cardDefinitions")
		if err != nil {
			logging.Error(fmt.Sprintf("An error occurred fetching data element by name: %v", err.Error()))
			return ""
		}

		cardDefinitions, err := getPriorityDataByElementName(profileId, elementId)
		if err != nil {
			logging.Error(fmt.Sprintf("An error occurred fetching card definitions: %v", err.Error()))
			return ""
		}

		if cardDefinitions == "" {
			return ""
		}

		type CardScheme struct {
			CardName string `json:"cardName"`
		}
		schemes := new([]CardScheme)
		err = json.Unmarshal([]byte(cardDefinitions), schemes)
		if err != nil {
			logging.Error(fmt.Sprintf("An error ocucurred deserialzing card definitions: %v", err.Error()))
			return ""
		}

		var options string
		for _, scheme := range *schemes {
			if options == "" {
				options = scheme.CardName
			} else {
				options = fmt.Sprintf("%s|%s", options, scheme.CardName)
			}
		}
		return options
	}

	// Returns an option string of files on the file server that match the provided pattern
	fsFilesOptions := func(pattern string) string {
		files, err := fileServer.NewFsReader(config.FileserverURL).GetAllFilesByPattern(pattern)
		if err != nil {
			logging.Error(fmt.Sprintf("An error occurred retrieving files from the file server: %s", err.Error()))
			return ""
		}
		return strings.Join(files, "|")
	}

	switch dataGroupName {
	case "dualCurrency":
		switch dataElementName {
		case "activeSchemes":
			optionsString = activeSchemesOptions(profileId)
			computedField = true
		}
	case "transactionRetrieval":
		switch dataElementName {
		case "merchantReceiptTemplate", "customerReceiptTemplate":
			optionsString = "|" + fsFilesOptions(".*\\.rmu")
			computedField = true
		case "parameterDefinitions":
			optionsString = "|" + fsFilesOptions("tr_parameter_definitions_.*\\.json")
			computedField = true
		}
	}
	return optionsString, computedField
}

func getPriorityDataByElementName(profileId, dataElementId int) (string, error) {
	db, err := GetDB()
	if err != nil {
		logging.Error("An error occurred executing GetDB() - %s", err.Error())
		return "", err
	}
	rows, err := db.Query("SELECT get_parent_datavalue(?, ?)", dataElementId, profileId)
	if err != nil {
		logging.Error("An error occurred executing procedure GET_TID_BY_TID - %s", err.Error())
		return "", err
	}
	defer rows.Close()

	var foundValue string
	for rows.Next() {
		rows.Scan(&foundValue)
	}

	return foundValue, nil
}

// When adding TID override this switch needs updating
func AddDefaultTidGroup(site SiteData, tidGroups *MappedDataGroups) {
	if tidGroups != nil {
		if _, groupPresent := (*tidGroups)[site.DataGroup]; groupPresent {
			if elem, elementPresent := (*tidGroups)[site.DataGroup].DataElements[site.Name]; elementPresent && site.OverridePriority.Valid {
				if elem.OverridePriority < int(site.OverridePriority.Int32) {
					return
				}
			}
		}
	}
	if site.TIDOverridable {
		tidGroups.AddElementFromSiteData(site)
	}
}

func BuildOptionsData(optionString, selectOptionsString, dataGroupName, dataElementName string, profileId int) ([]OptionData, bool) {
	computedOptionsString, isComputedField := computeSelectionOptions(profileId, dataGroupName, dataElementName)
	// If no data was available to populate the options string but it is a calculated field then set the value to __NO_DATA__
	// the UI handles __NO_DATA__ by rendering a multi-select field with a single disabled option which reads No Data Available
	if isComputedField {
		if computedOptionsString == "" {
			optionString = "__NO_DATA__"
		} else {
			optionString = computedOptionsString
		}
	}

	options := make([]OptionData, 0)
	// Build options model for data element and compare with data value to default selections
	var optionsSplit = strings.Split(optionString, "|")

	if optionString != "" {
		hasOptionsSelected := false
		for i := range optionsSplit {
			var od OptionData
			od.Option = strings.TrimSpace(optionsSplit[i])
			od.Selected = false

			od.Selected = strings.Contains(selectOptionsString, optionsSplit[i])
			if od.Selected && optionsSplit[i] != "" {
				hasOptionsSelected = true
			}
			options = append(options, od)
		}

		// In the case that an option has been removed but is still selected in the config, we need to make the
		// config value blank
		if selectOptionsString != "" && !hasOptionsSelected {
			selectOptionsString = ""
		}
	}

	return options, strings.Contains(optionString, selectOptionsString)
}

func GetPasswordFieldIds() []int {
	passwordIds := make([]int, 0)

	db, err := GetDB()
	if err != nil {
		logging.Error("error connecting to db, ", err)
		return passwordIds
	}

	rows, err := db.Query("SELECT `data_element_id` FROM data_element WHERE `is_password` = 1;")
	if err != nil {
		logging.Error("error fetching password field IDs, ", err)
		return passwordIds
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			logging.Error("error scanning row for password id, ", err)
			return passwordIds
		}

		passwordIds = append(passwordIds, id)
	}

	return passwordIds
}

// get Data
func GetProfileData(profileId int) ([]*DataGroup, []*ProfileChangeHistory) {
	db, err := GetDB()
	if err != nil {
		return nil, nil
	}

	rows, err := db.Query("Call profile_data_fetch(?)", profileId)
	if err != nil {
		return nil, nil
	}
	defer rows.Close()

	var profileGroups []*DataGroup

	var emptyDependencies []SiteData

	preAllocated := make(map[int]bool)
	tidOverRideDataElement, err := GetTIDOverRideDataElement()
	if err != nil {
		logging.Error(fmt.Sprintf("Error retrieving data element tid over ride enable %s", err.Error()))
		return nil, nil
	}

	for rows.Next() {

		var site SiteData
		var siteOptions string
		err = rows.Scan(&site.TIDID, &site.DataGroupID, &site.DataGroupDisplayName, &site.DataElementID, &site.Name,
			&site.DataGroup, &site.TIDOverridable, &site.DataType, &site.Tooltip, &site.DataValue, &site.IsAllowEmpty,
			&site.MaxLength, &site.ValidationExpression, &site.ValidationMessage, &site.FrontEndValidate, &siteOptions, &site.DisplayName, &site.IsPassword, &site.IsEncrypted, &site.SortOrderInGroup, &site.IsReadOnlyAtCreation, &site.IsRequiredAtAcquireLevel, &site.IsRequiredAtChainLevel)
		if err != nil {
			return nil, nil
		}

		if site.DataValue.Valid && site.IsEncrypted.Valid && site.IsEncrypted.Bool {
			site.DataValue.String, err = crypt.Decrypt(site.DataValue.String)
			if err != nil {
				return nil, nil
			}
		}

		// Build options model for data element and compare with data value to default selections
		site.Options, site.OptionSelectable = BuildOptionsData(siteOptions, site.DataValue.String, site.DataGroup, site.Name, profileId)

		if !site.DataValue.Valid {
			site.DataValue.String = ""
			if !tidOverRideDataElement[strconv.Itoa(site.DataElementID)] {
				emptyDependencies = append(emptyDependencies, site)
			}
		} else {
			profileGroups = addDataElement(site, profileGroups)
			preAllocated[site.DataElementID] = true
		}
	}

	for i := range emptyDependencies {
		if preAllocated[emptyDependencies[i].DataElementID] == false {
			profileGroups = addDataElement(emptyDependencies[i], profileGroups)
		}
	}

	changeHistory, err := GetProfileChangeHistory(profileId)

	if err != nil {
		return nil, nil
	}

	return profileGroups, changeHistory
}

func GetProfileChangeHistory(profileID int) ([]*ProfileChangeHistory, error) {
	const getProfileChangeHistoryCall = "get_profile_change_history"
	return GetChangeHistory(profileID, getProfileChangeHistoryCall)
}

func GetTIDChangeHistory(tidID int) ([]*ProfileChangeHistory, error) {
	const getTIDChangeHistoryCall = "get_tid_change_history"
	return GetChangeHistory(tidID, getTIDChangeHistoryCall)
}

func GetChangeHistory(id int, call string) ([]*ProfileChangeHistory, error) {
	var changes = make([]*ProfileChangeHistory, 0)

	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("CALL "+call+"(?)", id)
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

// Returns the element ID for the given group ID and element name
func GetElementId(groupId int, elementName string) (int, error) {
	db, err := GetDB()
	if err != nil {
		return -1, err
	}

	var elementId int

	rows, err := db.Query("Call get_element_id(?,?)", groupId, elementName)
	if err != nil {
		return -1, err
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&elementId)
	}

	return elementId, nil
}

// Returns the element ID for the given group name and element name
func GetElementIdFromGroupNameElementName(groupName string, elementName string) (int, error) {
	db, err := GetDB()
	if err != nil {
		return -1, err
	}

	var elementId int

	rows, err := db.Query("Call get_element_id_by_group_element_name(?,?)", groupName, elementName)
	if err != nil {
		return -1, err
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&elementId)
	}

	return elementId, nil
}

func GetThirdPartyEnabled(tidId, siteID int) (bool, error) {

	db, err := GetDB()
	if err != nil {
		return false, err
	}

	//Check to see if the TID is overridden
	exists, tidProfileId, err := GetTidProfileIdForSiteId(tidId, siteID)
	profileId := int(tidProfileId)

	var elementValue string
	//If the TID is overridden then use the tid profile
	if exists {

		//Fetch the value of mode for this profile
		elementValue, err = GetDataElementValue(profileId, "mode", "modules")
		if err != nil {
			return false, err
		}
	} else {
		//If the TID is not overridden then use the site profile
		if !exists {
			err = db.QueryRow("Call profile_data_element_fetch(?,?,?)", siteID, "mode", "modules").Scan(&elementValue)
			if err != nil {
				return false, err
			}

		}
	}

	if elementValue == "Third Party Application" {
		return true, err
	}

	return false, nil
}

func GetSiteIDForTid(tidId int) (int, error) {
	db, err := GetDB()
	if err != nil {
		return -1, err
	}

	return getSiteIDForTid(db, tidId)
}

func CheckTidExistsForMID(mid, tidId string) (int, error) {
	db, err := GetDB()
	if err != nil {
		return -1, err
	}

	count := 0
	siteIDrows := db.QueryRow("Call check_tid_exists_for_mid(?,?)", mid, tidId)
	if err != nil {
		return -1, err
	}

	err = siteIDrows.Scan(&count)
	if err != nil {
		return -1, err
	}

	return count, nil
}

func getSiteIDForTid(db *sql.DB, tidId int) (int, error) {
	siteIDrows, err := db.Query(fmt.Sprintf("SELECT site_id from tid_site WHERE tid_id = (%d);", tidId))
	if err != nil {
		return -1, err
	}
	defer siteIDrows.Close()

	var siteID int
	for siteIDrows.Next() {
		err = siteIDrows.Scan(&siteID)
		if err != nil {
			return -1, err
		}
	}

	return siteID, nil
}

func GetThirdPartyApksById(apkIds []int) ([]*ThirdPartyApk, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	var thirdPartyApks []*ThirdPartyApk
	queryArgs := make([]interface{}, len(apkIds))
	for i, tid := range apkIds {
		queryArgs[i] = tid
	}
	rows, err := db.Query("SELECT apk_id, `name` FROM third_party_apks WHERE apk_id IN (?"+strings.Repeat(",?", len(apkIds)-1)+")", queryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var apk ThirdPartyApk
		if err := rows.Scan(&apk.ApkID, &apk.Apk); err != nil {
			return nil, err
		}
		thirdPartyApks = append(thirdPartyApks, &apk)
	}
	return thirdPartyApks, nil
}

func GetThirdPartyApks(dataValue string) ([]*ThirdPartyApk, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	var thirdPartyApks []*ThirdPartyApk

	rows, err := db.Query("Call get_third_party_apks(?)", dataValue)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var apk ThirdPartyApk

		err = rows.Scan(&apk.ApkID, &apk.Apk)
		if err != nil {
			return nil, err
		}

		thirdPartyApks = append(thirdPartyApks, &apk)
	}
	return thirdPartyApks, nil
}

func AddUpdateTidAndGetThirdPartyApks(tidUpdateId int, tidId int, packageId int, updateDate string, thirdPartyApk string, dataValue string) ([]*ThirdPartyApk, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	var thirdPartyApks []*ThirdPartyApk
	rows, err := db.Query("Call add_and_update_tid_updates_and_flag(?,?,?,?,?,?)", tidUpdateId, tidId, packageId, updateDate, thirdPartyApk, dataValue)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var apk ThirdPartyApk
		err = rows.Scan(&apk.ApkID, &apk.Apk)
		if err != nil {
			return nil, err
		}
		thirdPartyApks = append(thirdPartyApks, &apk)
	}
	return thirdPartyApks, nil
}

// DB holds 3 tables packages, apks & a join table. This method is used to build a data
// structure from the 3 tables
func GetPackages() ([]*PackageData, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	// get packages
	var packages = []*PackageData{}
	packageRows, err := db.Query("Call get_all_packages()")
	if err != nil {
		return nil, err
	}
	defer packageRows.Close()

	for packageRows.Next() {
		var packageData PackageData
		var apks sql.NullString
		err = packageRows.Scan(&packageData.PackageID, &packageData.Version, &apks)
		if err != nil {
			return nil, err
		}

		packageData.Apks = strings.Split(apks.String, ",")

		packages = append(packages, &packageData)
	}

	return packages, nil
}

// I dont think this needs to have any special case code added to this.
func addDataElement(site SiteData, dataGroups []*DataGroup) []*DataGroup {
	//do we have data group or do we need to create one
	var currentDataGroup *DataGroup
	for i := range dataGroups {
		if dataGroups[i].DataGroupID == site.DataGroupID {
			currentDataGroup = dataGroups[i]
			break
		}
	}

	if currentDataGroup == nil {
		currentDataGroup = &DataGroup{DataGroupID: site.DataGroupID, DataGroup: site.DataGroup, DisplayName: site.DataGroupDisplayName}
		dataGroups = append(dataGroups, currentDataGroup)
	} else {
		// check for duplicate
		for i := range currentDataGroup.DataElements {
			if currentDataGroup.DataElements[i].ElementId == site.DataElementID {
				return dataGroups
			}
		}
	}
	var maxLength int = -1
	var validationExpression = ""
	var validationMessage = ""
	var tooltip = ""

	if site.MaxLength.Valid {
		maxLength = int(site.MaxLength.Int64)
	}

	if site.ValidationExpression.Valid {
		validationExpression = site.ValidationExpression.String
	}

	if site.ValidationMessage.Valid {
		validationMessage = site.ValidationMessage.String
	}

	var image template.URL
	if site.DataType == "FILE" && site.DataValue.String != "" {
		image = GetFile(site.DataValue.String, config.FileserverURL)
	}

	if site.Tooltip.Valid {
		tooltip = site.Tooltip.String
	}

	dataElement := DataElement{
		ElementId:                site.DataElementID,
		Name:                     site.Name,
		Type:                     site.DataType,
		DataValue:                site.DataValue.String,
		MaxLength:                maxLength,
		ValidationExpression:     validationExpression,
		ValidationMessage:        validationMessage,
		FrontEndValidate:         site.FrontEndValidate != 0,
		Overriden:                site.Overriden.Valid && site.Overriden.Int64 != 0,
		Options:                  site.Options,
		OptionSelectable:         site.OptionSelectable,
		DisplayName:              site.DisplayName.String,
		Image:                    image,
		IsPassword:               site.IsPassword,
		IsEncrypted:              site.IsEncrypted.Valid && site.IsEncrypted.Bool,
		SortOrderInGroup:         site.SortOrderInGroup,
		IsNotOverrideable:        site.IsNotOverridable.Valid && site.IsNotOverridable.Bool,
		Tooltip:                  tooltip,
		IsAllowedEmpty:           site.IsAllowEmpty,
		FileMaxSize:              int(site.FileMaxSize.Int64),
		FileMinRatio:             site.FileMinRatio.Float64,
		FileMaxRatio:             site.FileMaxRatio.Float64,
		IsReadOnlyAtCreation:     site.IsReadOnlyAtCreation,
		IsRequiredAtAcquireLevel: site.IsRequiredAtAcquireLevel,
		IsRequiredAtChainLevel:   site.IsRequiredAtChainLevel,
	}
	currentDataGroup.DataElements = append(currentDataGroup.DataElements, dataElement)

	return dataGroups
}

func GetProfilesByTypeName(name string, user *entities.TMSUser) ([]*ProfileData, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	// Find user acquirers to limit search results
	acquirers, err := GetUserAcquirerPermissions(user)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("Call profile_list_fetch_by_type_name(?, ?)", name, acquirers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dataList = []*ProfileData{}
	for rows.Next() {
		var data ProfileData
		err = rows.Scan(&data.ID, &data.Name)
		if err != nil {
			return nil, err
		}
		dataList = append(dataList, &data)
	}
	return dataList, nil
}

func GetProfileFields(acquirerId int, chainId int, selectedDataGroups []string) ([]*DataGroup, []*DataGroup, []*DataGroup, []*DataGroup) {
	db, err := GetDB()
	if err != nil {
		return nil, nil, nil, nil
	}

	rows, err := db.Query("Call get_data_elements_for_new_site(?, ?)", acquirerId, chainId)
	if err != nil {
		logging.Error(err)
		return nil, nil, nil, nil
	}
	defer rows.Close()

	var emptyDependencies []SiteData
	var siteGroups []*DataGroup
	var chainGroups []*DataGroup
	var acquirerGroups []*DataGroup
	var globalGroups []*DataGroup

	preAllocated := make(map[int]bool)

	for rows.Next() {

		var site SiteData
		var options string
		err = rows.Scan(&site.DataGroupID, &site.DataGroup, &site.DataGroupDisplayName, &site.DataElementID, &site.Tooltip, &site.Name, &site.Source, &site.DataValue, &site.DataType, &site.IsAllowEmpty, &options, &site.DisplayName, &site.IsPassword, &site.IsEncrypted, &site.SortOrderInGroup, &site.RequiredAtSiteLevel, &site.IsNotOverridable, &site.IsReadOnlyAtCreation, &site.IsRequiredAtAcquireLevel, &site.IsRequiredAtChainLevel)
		if err != nil {
			logging.Error(err)
			return nil, nil, nil, nil
		}

		// If the found element is not one of the selected data groups then continue, it should not be displayed so do not
		// include the group or it's elements in the returned data.
		if !sliceHelper.SlicesOfStringContains(selectedDataGroups, strconv.Itoa(site.DataGroupID)) {
			continue
		}

		if site.IsEncrypted.Valid && site.IsEncrypted.Bool && site.DataValue.Valid {
			site.DataValue.String, err = crypt.Decrypt(site.DataValue.String)
			if err != nil {
				return nil, nil, nil, nil
			}
		}

		site.Options, site.OptionSelectable = BuildOptionsData(options, "", site.DataGroup, site.Name, chainId)

		if !site.Source.Valid || !site.DataValue.Valid {
			site.DataValue.String = ""
			emptyDependencies = append(emptyDependencies, site)
		} else {
			if (site.IsAllowEmpty == false && site.DataValue.String == "") || site.RequiredAtSiteLevel {
				siteGroups = addDataElement(site, siteGroups)
				preAllocated[site.DataElementID] = true
			} else {
				switch site.Source.String {
				case "chain":
					// If the element is being added at chain level, then remove it from the acquirer and global elements
					// as chain is higher priority than acquirer or global
					acquirerGroups = RemoveDataElementIfPresentInGroups(acquirerGroups, site.DataElementID)
					globalGroups = RemoveDataElementIfPresentInGroups(globalGroups, site.DataElementID)

					chainGroups = addDataElement(site, chainGroups)
					preAllocated[site.DataElementID] = true
				case "acquirer":
					// If the element is being added at acquirer level, then remove it from the global elements as
					// acquirer is higher priority than global
					globalGroups = RemoveDataElementIfPresentInGroups(globalGroups, site.DataElementID)

					// If the element is already present at a chain level then it should not be added at acquirer level
					// as chain is higher priority.
					if present, _, _ := DataGroupsContainsDataElement(chainGroups, site.DataElementID); !present {
						acquirerGroups = addDataElement(site, acquirerGroups)
						preAllocated[site.DataElementID] = true
					}
				case "global":
					// If the element is already present at a chain level or acquirer level then it should not be added
					//at global level as chain and acquirer they are both higher priority.
					presentInChainGroups, _, _ := DataGroupsContainsDataElement(chainGroups, site.DataElementID)
					presentInAcquirerGroups, _, _ := DataGroupsContainsDataElement(acquirerGroups, site.DataElementID)
					if !presentInChainGroups && !presentInAcquirerGroups {
						globalGroups = addDataElement(site, globalGroups)
						preAllocated[site.DataElementID] = true
					}
				case "tid":
					siteGroups = addDataElement(site, siteGroups)
					preAllocated[site.DataElementID] = true
				case "site":
					siteGroups = addDataElement(site, siteGroups)
					preAllocated[site.DataElementID] = true

				default:
					//TODO handle error
				}
			}
		}
	}

	for i := range emptyDependencies {
		if preAllocated[emptyDependencies[i].DataElementID] == false {
			siteGroups = addDataElement(emptyDependencies[i], siteGroups)
		}
	}

	return siteGroups, chainGroups, acquirerGroups, globalGroups
}

func GetProfileIdForSite(siteId int) (int, error) {
	db, err := GetDB()
	if err != nil {
		return -1, err
	}

	rows, err := db.Query("Select site_profile(?)", siteId)
	if err != nil {
		return -1, err
	}
	defer rows.Close()

	var profileId int
	for rows.Next() {
		err = rows.Scan(&profileId)
	}

	if err != nil {
		return -1, err
	}

	return profileId, nil
}

func GetProfileIdForTIDAndSiteID(tidId string, siteId int) (string, error) {
	db, err := GetDB()
	if err != nil {
		return "", err
	}

	var dataValue sql.NullString
	rows, err := db.Query("Call get_thirdparty_partialpackagename_datavalue(?,?)", siteId, tidId)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	if rows.Next() {
		if err = rows.Scan(&dataValue); err != nil {
			return "", err
		}
	}
	if dataValue.Valid {
		return dataValue.String, nil
	}

	return dataValue.String, nil
}

func GetContactUsFields(acquirerId int) (ContactUsData, error) {
	db, err := GetDB()
	if err != nil {
		return ContactUsData{}, err
	}

	rows, err := db.Query("SELECT AcquirerName, AcquirerPrimaryPhone, AcquirerSecondaryPhone, AcquirerEmail, AcquirerAddressLineOne, AcquirerAddressLineTwo, AcquirerAddressLineThree, FurtherInformation FROM contact_us WHERE AcquirerId = ?", acquirerId)
	if err != nil {
		logging.Error(err)
		return ContactUsData{}, err
	}
	defer rows.Close()

	var data ContactUsData

	for rows.Next() {
		err = rows.Scan(&data.AcquirerName,
			&data.AcquirerPrimaryPhone,
			&data.AcquirerSecondaryPhone,
			&data.AcquirerEmail,
			&data.AcquirerAddressLineOne,
			&data.AcquirerAddressLineTwo,
			&data.AcquirerAddressLineThree,
			&data.FurtherInformation)

		data.Valid = true
	}
	if err != nil {
		return ContactUsData{}, err
	}

	return data, nil
}

func SetContactUsFields(acquirerId int, data ContactUsData) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	_, error := db.Exec("CALL set_contact_us_data(?, ?, ?, ?, ?, ?, ?, ?, ?)",
		acquirerId,
		data.AcquirerName,
		data.AcquirerPrimaryPhone,
		data.AcquirerSecondaryPhone,
		data.AcquirerEmail,
		data.AcquirerAddressLineOne,
		data.AcquirerAddressLineTwo,
		data.AcquirerAddressLineThree,
		data.FurtherInformation)
	if error != nil {
		logging.Error(error)
		return error
	}
	return nil
}

func GetAllDataElementsMetadata(profileId int) (map[int]DataElement, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("CALL get_all_metadata_elements()")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	returnMap := make(map[int]DataElement, 0)
	for rows.Next() {
		var element DataElement
		var maxLength sql.NullInt64
		var validationExpression sql.NullString
		var validationMessage sql.NullString
		var frontEndValidate int
		var unique int
		var options string
		var displayName sql.NullString
		var isPassword sql.NullBool
		var isEncrypted sql.NullBool
		var tooltip sql.NullString
		var fileMaxSize sql.NullInt64
		var fileMinRatio sql.NullFloat64
		var fileMaxRatio sql.NullFloat64
		err = rows.Scan(&element.ElementId, &element.Name, &element.Type, &element.IsAllowedEmpty, &maxLength,
			&validationExpression, &validationMessage, &frontEndValidate, &unique, &options, &displayName, &isPassword, &isEncrypted,
			&tooltip, &fileMaxSize, &fileMinRatio, &fileMaxRatio, &element.IsReadOnlyAtCreation, &element.IsRequiredAtAcquireLevel, &element.IsRequiredAtChainLevel)

		group, err := NewDataGroupRepository().FindByDataElementId(element.ElementId)
		if err != nil {
			return nil, err
		}
		element.Options, element.OptionSelectable = BuildOptionsData(options, "", group.Name, element.Name, profileId)

		if maxLength.Valid {
			element.MaxLength = int(maxLength.Int64)
		}

		if validationExpression.Valid {
			element.ValidationExpression = validationExpression.String
		}

		if validationMessage.Valid {
			element.ValidationMessage = validationMessage.String
		}

		if displayName.Valid {
			element.DisplayName = displayName.String
		}

		if isPassword.Valid {
			element.IsPassword = isPassword.Bool
		}

		if isEncrypted.Valid && isEncrypted.Bool {
			element.IsEncrypted = isEncrypted.Bool
			var decryptError error
			element.DataValue, decryptError = crypt.Decrypt(element.DataValue)
			if decryptError != nil {
				return nil, decryptError
			}
		}

		element.FrontEndValidate = frontEndValidate > 0
		element.Unique = unique > 0
		if tooltip.Valid {
			element.Tooltip = tooltip.String
		}
		if fileMaxSize.Valid {
			element.FileMaxSize = int(fileMaxSize.Int64)
		}
		if fileMinRatio.Valid {
			element.FileMinRatio = fileMinRatio.Float64
		}
		if fileMaxRatio.Valid {
			element.FileMaxRatio = fileMaxRatio.Float64
		}
		returnMap[element.ElementId] = element
	}
	return returnMap, nil
}

// profileId parameter is optional but is used for computed data elements, pass 0 as an id to skip this.
func GetDataElementMetadata(dataElementId int, profileId int) (DataElement, error) {
	var element DataElement

	db, err := GetDB()
	if err != nil {
		return element, err
	}

	rows, err := db.Query("CALL fetch_data_element_metadata(?)", dataElementId)
	if err != nil {
		return element, err
	}
	defer rows.Close()

	var maxLength sql.NullInt64
	var validationExpression sql.NullString
	var validationMessage sql.NullString
	var frontEndValidate int
	var unique int
	var options string
	var displayName sql.NullString
	var isPassword sql.NullBool
	var isEncrypted sql.NullBool
	var tooltip sql.NullString
	var fileMaxSize sql.NullInt64
	var fileMinRatio sql.NullFloat64
	var fileMaxRatio sql.NullFloat64

	for rows.Next() {
		err = rows.Scan(&element.ElementId, &element.Name, &element.Type, &element.IsAllowedEmpty, &maxLength,
			&validationExpression, &validationMessage, &frontEndValidate, &unique, &options, &displayName, &isPassword, &isEncrypted, &tooltip, &fileMaxSize, &fileMinRatio, &fileMaxRatio, &element.IsReadOnlyAtCreation, &element.IsRequiredAtAcquireLevel, &element.IsRequiredAtChainLevel)
		if err != nil {
			return element, err
		}
	}

	group, err := NewDataGroupRepository().FindByDataElementId(element.ElementId)
	if err != nil {
		return element, err
	}

	element.Options, element.OptionSelectable = BuildOptionsData(options, "", group.Name, element.Name, profileId)

	if maxLength.Valid {
		element.MaxLength = int(maxLength.Int64)
	}

	if validationExpression.Valid {
		element.ValidationExpression = validationExpression.String
	}

	if validationMessage.Valid {
		element.ValidationMessage = validationMessage.String
	}

	if displayName.Valid {
		element.DisplayName = displayName.String
	}

	if isPassword.Valid {
		element.IsPassword = isPassword.Bool
	}

	if isEncrypted.Valid && isEncrypted.Bool {
		var decryptError error
		element.DataValue, decryptError = crypt.Decrypt(element.DataValue)
		if decryptError != nil {
			return element, err
		}
	}
	if tooltip.Valid {
		element.Tooltip = tooltip.String
	}

	element.FrontEndValidate = frontEndValidate > 0
	element.Unique = unique > 0

	if fileMaxSize.Valid {
		element.FileMaxSize = int(fileMaxSize.Int64)
	}
	if fileMinRatio.Valid {
		element.FileMinRatio = fileMinRatio.Float64
	}
	if fileMaxRatio.Valid {
		element.FileMaxRatio = fileMaxRatio.Float64
	}

	if err != nil {
		return element, err
	}

	return element, nil
}

func GetIsUnique(elementId int, elementValue string, profile int) (bool, error) {
	db, err := GetDB()
	if err != nil {
		return false, err
	}

	rows, err := db.Query("CALL count_matched_values(?,?, ?)", elementId, elementValue, profile)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	matchedRows := 0

	for rows.Next() {
		err = rows.Scan(&matchedRows)
	}

	if err != nil {
		return false, err
	}

	return matchedRows == 0, nil

}

func GetSiteList(searchTerm string, user *entities.TMSUser) ([]*SiteList, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	// Find user acquirers to limit search results
	acquirers, err := GetUserAcquirerPermissions(user)
	if err != nil {
		return nil, err
	}

	// Now get search results using assigned user acquirer permissions
	rows, err := db.Query("Call site_list_fetch(?,?)", searchTerm, acquirers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sites []*SiteList
	var errorSites map[string]string
	errorSites = make(map[string]string)

	for rows.Next() {
		var site SiteList

		err = rows.Scan(
			&site.SiteID,
			&site.SiteProfileID,
			&site.SiteName,
			&site.ChainProfileID,
			&site.ChainName,
			&site.AcquirerProfileID,
			&site.AcquirerName,
			&site.GlobalProfileID,
			&site.GlobalName,
			&site.MerchantId)

		if err != nil {
			errorSites[site.SiteName] = ""
		}

		sites = append(sites, &site)
	}

	if len(errorSites) > 0 {
		err = InvalidSiteDataError(errorSites)
	}

	return sites, err
}

func getFilteredMIDCount(searchTerm string, acquirers string) (int, error) {
	db, err := GetDB()
	if err != nil {
		logging.Error(err)
		return 0, err
	}

	logging.Information("Calling get_mid_count_filtered")

	row := db.QueryRow("CALL get_mid_count_filtered(?, ?)", searchTerm, acquirers)

	var midCount int

	row.Scan(&midCount)

	return midCount, nil
}

func GetSitePage(searchTerm string, offset string, amount string, orderedColumn string, orderDirection string, user *entities.TMSUser) (page []*SiteList, total int, filtered int, err error) {

	db, err := GetDB()
	if err != nil {
		return nil, 0, 0, err
	}

	// Find user acquirers to limit search results
	acquirers, err := GetUserAcquirerPermissions(user)
	if err != nil {
		return nil, 0, 0, err
	}

	total, err = getFilteredMIDCount(searchTerm, acquirers)
	if err != nil {
		logging.Error(err)
		return nil, 0, 0, err
	}

	// If the amount is -1 then we want all the sites
	if amount == "-1" {
		amount = strconv.Itoa(total)
	}

	rows, err := db.Query("CALL get_site_page(?,?, ?, ?)", searchTerm, acquirers, offset, amount)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()

	var sites []*SiteList
	errorSites := make(map[string]string)

	for rows.Next() {
		var site SiteList

		err = rows.Scan(
			&site.SiteID,
			&site.SiteProfileID,
			&site.SiteName,
			&site.ChainProfileID,
			&site.ChainName,
			&site.AcquirerProfileID,
			&site.AcquirerName,
			&site.GlobalProfileID,
			&site.GlobalName,
			&site.MerchantId)

		if err != nil {
			errorSites[site.SiteName] = ""
		}

		site.AcquirerName = html.EscapeString(site.AcquirerName)
		site.ChainName = html.EscapeString(site.ChainName)
		site.SiteName = html.EscapeString(site.SiteName)

		sites = append(sites, &site)
	}

	//sort the slice
	sort.Slice(sites, func(i, j int) bool {
		switch orderedColumn {
		case "0":
			if strings.ToUpper(orderDirection) == "ASC" {
				return sites[i].MerchantId < sites[j].MerchantId
			} else {
				return sites[i].MerchantId > sites[j].MerchantId
			}
		case "1":
			if strings.ToUpper(orderDirection) == "ASC" {
				return strings.ToLower(sites[i].SiteName) < strings.ToLower(sites[j].SiteName)
			} else {
				return strings.ToLower(sites[i].SiteName) > strings.ToLower(sites[j].SiteName)
			}
		case "2":
			if strings.ToUpper(orderDirection) == "ASC" {
				return strings.ToLower(sites[i].ChainName) < strings.ToLower(sites[j].ChainName)
			} else {
				return strings.ToLower(sites[i].ChainName) > strings.ToLower(sites[j].ChainName)
			}
		}
		return false
	})

	filtered = total

	return sites, total, filtered, nil
}

func versionFor(name string) (sqlSelect string) {
	return "SELECT MAX(d.version) FROM profile_data AS d WHERE d.data_element_id = (" + dataElementIDFor(name) + ") AND d.profile_id = p2.profile_id AND d.approved = 1"
}

func dataElementIDFor(name string) (sqlSelect string) {
	return "SELECT de.data_element_id FROM data_element AS de WHERE de.name = \"" + name + "\""
}

func siteProfileLeftJoin(index string, siteName string) (leftJoin string) {
	return fmt.Sprintf(`LEFT JOIN (
		site_profiles AS tp%[1]v
		
		JOIN profile AS p%[1]v
		ON p%[1]v.profile_id = tp%[1]v.profile_id
		
		JOIN profile_type AS pt%[1]v
		ON pt%[1]v.profile_type_id = p%[1]v.profile_type_id
		AND pt%[1]v.priority = %[1]v
	
	) ON tp%[1]v.site_id = %[2]v.site_id`, index, siteName)

}

func profileDataLeftJoin(profileDataName string, dataElementName string, profileIndex string) (leftJoin string) {
	return fmt.Sprintf(`LEFT JOIN profile_data AS %[1]v
	ON %[1]v.profile_id = p%[4]v.profile_id
	AND %[1]v.data_element_id = (%[2]v)
	AND %[1]v.version = (%[3]v)`, profileDataName, dataElementIDFor(dataElementName), versionFor(dataElementName), profileIndex)
}

func InvalidSiteDataError(errorSites map[string]string) error {
	var siteList string
	for k := range errorSites {
		if k == "" {
			return errors.New("Sites exist with no name assigned")
		}
		siteList += k + ", "
	}

	return errors.New("The following sites contain invalid data: " + siteList)
}

func GetDataGroups() ([]DataGroup, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("Call fetch_data_groups()")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []DataGroup

	groups = make([]DataGroup, 0)

	for rows.Next() {
		var dg = DataGroup{}
		rows.Scan(&dg.DataGroupID, &dg.DataGroup, &dg.DisplayName)
		groups = append(groups, dg)
	}

	return groups, nil
}

func GetDataGroupsWithProfileId(acquirerId, chainId, profileID int, isChain bool) ([]DataGroup, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("Call get_datagroup_with_profile(?,?,?,?)", acquirerId, chainId, profileID, isChain)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var groups []DataGroup
	groups = make([]DataGroup, 0)
	for rows.Next() {
		var dg = DataGroup{}
		rows.Scan(&dg.DataGroupID, &dg.DataGroup, &dg.DisplayName, &dg.PreSelected, &dg.IsSelected)
		groups = append(groups, dg)
	}
	return groups, nil
}

// profileId parameter is optional but is used for computed data elements, pass 0 as an id to skip this.
func GetDataElementsForGroup(groupId string, profileId int) (*DataGroup, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(`
			SELECT
				dg.data_group_id,
				dg.name,
				dg.displayname_en, 
				e.data_element_id,
				e.name,
				e.displayname_en,
				e.datatype,
				e.max_length,
				e.validation_expression,
				e.validation_message,
				e.front_end_validate,
				e.options,
				e.is_password,
				e.is_encrypted,
				e.sort_order_in_group,
				e.tooltip,
				e.is_allow_empty,
			    IFNULL(e.file_max_size, 0),
			    IFNULL(e.file_min_ratio, 0),
			    IFNULL(e.file_max_ratio, 0),
			    is_read_only_at_creation,
				required_at_acquirer_level,
				required_at_chain_level
			FROM data_group dg
			LEFT JOIN data_element e ON e.data_group_id = dg.data_group_id
			WHERE dg.data_group_id = ?`, groupId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	group := DataGroup{}
	group.DataGroupID, _ = strconv.Atoi(groupId)
	group.DataElements = make([]DataElement, 0)

	for rows.Next() {
		var (
			dgDisplayName, deDisplayName, validationExpression, validationMessage sql.NullString
			maxLength                                                             sql.NullInt64
			optionString                                                          string
			isPassword, isEncrypted                                               sql.NullBool
			de                                                                    DataElement
		)

		err = rows.Scan(&group.DataGroupID,
			&group.DataGroup,
			&dgDisplayName,
			&de.ElementId,
			&de.Name,
			&deDisplayName,
			&de.Type,
			&maxLength,
			&validationExpression,
			&validationMessage,
			&de.FrontEndValidate,
			&optionString,
			&isPassword,
			&isEncrypted,
			&de.SortOrderInGroup,
			&de.Tooltip,
			&de.IsAllowedEmpty,
			&de.FileMaxSize,
			&de.FileMinRatio,
			&de.FileMaxRatio,
			&de.IsReadOnlyAtCreation,
			&de.IsRequiredAtAcquireLevel,
			&de.IsRequiredAtChainLevel)

		if err != nil {
			return nil, err
		}

		if optionString != "" {
			de.Options, de.OptionSelectable = BuildOptionsData(optionString, "", group.DataGroup, de.Name, profileId)
		}

		if dgDisplayName.Valid {
			group.DisplayName = dgDisplayName.String
		} else {
			group.DisplayName = group.DataGroup
		}

		if deDisplayName.Valid {
			de.DisplayName = deDisplayName.String
		} else {
			de.DisplayName = de.Name
		}

		if maxLength.Valid {
			de.MaxLength = int(maxLength.Int64)
		}
		if validationExpression.Valid {
			de.ValidationExpression = validationExpression.String
		}
		if validationMessage.Valid {
			de.ValidationMessage = validationMessage.String
		}

		if isPassword.Valid {
			de.IsPassword = isPassword.Bool
		}

		group.DataElements = append(group.DataElements, de)
	}

	return &group, nil
}

func GetDataElementByNameAndGroupID(elementName string, groupID int) (int, error) {
	db, err := GetDB()
	if err != nil {
		return 0, err
	}

	var dataElementId int
	err = db.QueryRow(`SELECT data_element_id FROM data_element 
                       WHERE name=? AND data_group_id=?`, elementName, groupID).Scan(&dataElementId)
	if err != nil {
		return 0, err
	}

	return dataElementId, nil
}

func GetDataGroupByName(groupName string) (int, error) {
	db, err := GetDB()
	if err != nil {
		return 0, err
	}

	var dataGroupId int
	err = db.QueryRow(`SELECT dg.data_group_id 
								FROM data_group dg 
								WHERE dg.name = ?`, groupName).Scan(&dataGroupId)
	if err != nil {
		return 0, err
	}

	return dataGroupId, nil
}

func GetDpoMomoFieldsData() (string, error) {
	db, err := GetDB()
	if err != nil {
		return "NoData", err
	}
	rows, err := db.Query("SELECT options from data_element where data_group_id=(select data_group_id from data_group where name='dpoMomo') and name ='mno'")
	if err != nil {
		return "NoData", err
	}
	defer rows.Close()
	var options string
	for rows.Next() {
		rows.Scan(&options)
	}
	if strings.TrimSpace(options) == "" {
		return "NoData", nil
	}
	return options, nil
}

func GetSiteProfiles(siteId int) ([]ProfileData, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("CALL get_site_profiles(?)", siteId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []ProfileData

	profiles = make([]ProfileData, 0)

	for rows.Next() {
		var p = ProfileData{}
		rows.Scan(&p.ID, &p.TypeId, &p.Type)
		profiles = append(profiles, p)
	}

	return profiles, nil
}

func GetAcquirerNameForSite(siteId int) (string, error) {
	db, err := GetDB()
	if err != nil {
		return "", err
	}

	var acqName sql.NullString
	if rows, err := db.Query("CALL fetch_site_acquirer(?)", siteId); err == nil {
		defer rows.Close()
		for rows.Next() {
			rows.Scan(&acqName)
		}
	} else {
		return "", err
	}

	return acqName.String, nil

}

func GetAcquirerName(profileID int) (string, error) {
	db, err := GetDB()
	if err != nil {
		return "", err
	}

	var acqName sql.NullString
	err = db.QueryRow("CALL get_acquirer_name(?)", profileID).Scan(&acqName)
	if err != nil {
		return "", err
	}

	return acqName.String, nil
}

func GetAcquirerIdForChain(chainId string) string {
	var acquirerId string

	db, err := GetDB()
	if err != nil {
		return "-1"
	}

	rows, err := db.Query("SELECT cp.acquirer_id FROM chain_profiles cp WHERE cp.chain_profile_id = ?", chainId)
	if err != nil {
		return "-1"
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&acquirerId)
	}

	return acquirerId
}

func GetGroupsForProfile(profileId string) ([]DataGroup, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("CALL get_data_groups_by_profile_id(?)", profileId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []DataGroup

	groups = make([]DataGroup, 0)

	for rows.Next() {
		var dg = DataGroup{}
		rows.Scan(&dg.DataGroupID, &dg.DataGroup)
		groups = append(groups, dg)
	}

	return groups, nil
}

func ClearDataGroupsForProfile(profileId string) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	_, err = db.Exec("CALL clear_profile_datagroups(?)", profileId)

	return err
}

func ClearDataElementsForDisabled(profileID string, dataGroups []dataGroup.DataGroup) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	var queryArgs []interface{}
	queryArgs = append(queryArgs, profileID)
	for _, dataGroup := range dataGroups {
		queryArgs = append(queryArgs, dataGroup.ID)
	}

	_, err = db.Exec(`DELETE FROM 
		profile_data 
		WHERE profile_id = ? 
		AND data_element_id 
		IN (SELECT data_element_id 
		FROM  data_element 
		WHERE data_group_id IN (?`+strings.Repeat(",?", len(queryArgs)-2)+"))", queryArgs...)

	return err
}

func GetNewRemovedOverrideValue(siteId int, elementId int) (string, error) {
	db, err := GetDB()
	if err != nil {
		return "", err
	}

	rows, err := db.Query("CALL get_removed_override_value(?,?)", siteId, elementId)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var count = 0
	var dataValue sql.NullString
	for rows.Next() {
		count++
		rows.Scan(&dataValue)
	}

	if count < 1 || !dataValue.Valid {
		return "", nil // No overridable value underlying override may have been removed.
	}

	return dataValue.String, err

}

func (d *SiteManagementDal) GetUsersForSite(siteID int) ([]entities.SiteUser, error) {
	return GetUsersForSite(siteID)
}

func GetUsersForSite(siteID int) ([]entities.SiteUser, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("CALL get_site_user_data(?)", siteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []entities.SiteUser{}
	for rows.Next() {
		user := entities.SiteUser{}
		var siteid int
		var mods string
		var isEncrypted bool
		err = rows.Scan(&user.UserId, &siteid, &user.Username, &user.PIN, &mods, &isEncrypted)
		if err != nil {
			return nil, err
		}

		if isEncrypted {
			user.PIN, err = crypt.Decrypt(user.PIN)
			if err != nil {
				return nil, err
			}
		}

		user.Modules = strings.Split(mods, ",")
		users = append(users, user)
	}

	return users, nil
}

func (d *SiteManagementDal) GetUsersForTid(tidId int) ([]entities.SiteUser, error) {
	return GetUsersForTid(tidId)
}

func GetUsersForTidFromSiteId(siteId int) ([]entities.SiteUser, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("CALL get_tid_user_data_of_site(?)", siteId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []entities.SiteUser{}
	for rows.Next() {
		user := entities.SiteUser{}
		var mods string
		var isEncrypted bool

		err = rows.Scan(&user.UserId, &user.TidId, &user.Username, &user.PIN, &mods, &isEncrypted)
		if err != nil {
			return nil, err
		}

		if isEncrypted {
			user.PIN, err = crypt.Decrypt(user.PIN)
			if err != nil {
				return nil, err
			}
		}

		user.Modules = strings.Split(mods, ",")
		users = append(users, user)
	}

	return users, nil
}

func GetUsersForTid(tid int) ([]entities.SiteUser, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("CALL get_tid_user_data(?)", tid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []entities.SiteUser{}
	for rows.Next() {
		user := entities.SiteUser{}
		var site int
		var mods string
		var isEncrypted bool

		err = rows.Scan(&user.UserId, &site, &user.Username, &user.PIN, &mods, &isEncrypted)
		if err != nil {
			return nil, err
		}

		if isEncrypted {
			user.PIN, err = crypt.Decrypt(user.PIN)
			if err != nil {
				return nil, err
			}
		}

		user.Modules = strings.Split(mods, ",")
		users = append(users, user)
	}

	return users, nil
}

func (d *SiteManagementDal) GetTidUsersForSite(siteID int) ([]entities.SiteUser, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	tidUsers := make([]entities.SiteUser, 0)
	tids, err := GetSiteTids(db, siteID)
	if err != nil {
		return nil, err
	}

	for _, tid := range tids {
		users, err := d.GetUsersForTid(tid)
		if err != nil {
			return nil, err
		}

		tidUsers = append(tidUsers, users...)
	}

	return tidUsers, nil
}

func GetAvailableModules() ([]string, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	var modules string

	err = db.QueryRow("CALL get_available_modules()").Scan(&modules)
	if err != nil {
		return nil, err
	}

	return strings.Split(modules, "|"), nil
}

// Create multiple users against a site
func (d *SiteManagementDal) AddOrUpdateSiteUsers(siteId int, users []*entities.SiteUser) ([]*UserUpdateResult, error) {
	return updateUsers(siteId, 0, users, AddOrUpdateSiteUser)
}

// Create multiple users against a tid
func (d *SiteManagementDal) AddOrUpdateTidUserOverride(tid int, users []*entities.SiteUser) ([]*UserUpdateResult, error) {
	return updateUsers(0, tid, users, AddOrUpdateTidUserOverride)
}

// Create multiple users against a site
func (d *SiteManagementDal) DeleteSiteUsers(siteId int, userIds []int) ([]*UserUpdateResult, error) {
	users := make([]*entities.SiteUser, 0)
	for _, userId := range userIds {
		users = append(users, &entities.SiteUser{UserId: userId})
	}
	return updateUsers(siteId, 0, users, DeleteSiteUser)
}

// Create multiple users against a tid
func (d *SiteManagementDal) DeleteTidUsers(userIds []int) ([]*UserUpdateResult, error) {
	users := make([]*entities.SiteUser, 0)
	for _, userId := range userIds {
		users = append(users, &entities.SiteUser{UserId: userId})
	}
	return updateUsers(0, 0, users, DeleteTidUserOverride)
}

func (d *SiteManagementDal) GetUserForId(userId int) (*entities.SiteUser, error) {
	return GetUserForId(userId)
}

func GetUserForId(userId int) (*entities.SiteUser, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("CALL get_site_user_by_id(?)", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user := new(entities.SiteUser)
	for rows.Next() {
		var site int
		var mods string
		var isEncrypted bool

		err = rows.Scan(&user.UserId, &site, &user.Username, &user.PIN, &mods, &isEncrypted)
		if err != nil {
			return nil, err
		}

		if isEncrypted {
			user.PIN, err = crypt.Decrypt(user.PIN)
			if err != nil {
				return nil, err
			}
		}

		user.Modules = strings.Split(mods, ",")
	}

	return user, nil
}

// Adds all the given users to a site or a tid
// if a user's ID <= 0 then the user is created else it is updated
func updateUsers(siteId int, tid int, users []*entities.SiteUser, action TransactionPurpose) ([]*UserUpdateResult, error) {
	updateResults := make([]*UserUpdateResult, 0)

	//get the db connection, continue to the next user if error
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	//This is done as one transaction per user
	for _, user := range users {
		//Init the current result object
		updateResult := &UserUpdateResult{Result: DbTransactionResult{Success: true}, User: *user}
		//begin transaction, continue to the next user if error
		tx, err := db.Begin()
		if err != nil {
			updateResult.Result.SetError(errors.New("Unable to start transaction"))
			updateResults = append(updateResults, updateResult)
			continue
		}

		//encrypt the pin if its required
		PIN := user.PIN
		isEncrypted := false
		if crypt.UseEncryption {
			PIN = crypt.Encrypt(PIN)
			isEncrypted = true
		}

		switch action {
		case AddOrUpdateTidUserOverride:
			//Add tid users
			var returnRows *sql.Rows
			if tidExists, _, _ := CheckThatTidExists(tid); tidExists {
				returnRows, err = addTidUser(tx, user.UserId, tid, user.Username, PIN, userModulesToString(*user), isEncrypted)
			} else {
				updateResult.Result.SetErrorByString("TID does not exist")
			}
			if updateResult.User.UserId <= 0 {
				updateResult.Action = "Add TID User Override"
			} else {
				updateResult.Action = "Update TID User Override"
			}
			if err == nil {
				for returnRows.Next() {
					_ = returnRows.Scan(&updateResult.User.UserId)
				}
				returnRows.Close()
			}
			updateResult.User.TidId = tid
		case AddOrUpdateSiteUser:
			//Add site users
			var returnRows *sql.Rows
			returnRows, err = addSiteUsers(tx, user.UserId, siteId, user.Username, PIN, userModulesToString(*user), isEncrypted)

			if updateResult.User.UserId <= 0 {
				updateResult.Action = "Add Site User"
			} else {
				updateResult.Action = "Update Site User"
			}
			if err == nil {
				for returnRows.Next() {
					_ = returnRows.Scan(&updateResult.User.UserId)
				}
			}
			updateResult.User.SiteId = siteId

			if returnRows != nil {
				returnRows.Close()
			}

		case DeleteSiteUser:
			foundUser, _ := GetUserForId(user.UserId)
			updateResult.User = *foundUser
			updateResult.Action = "Delete Site User"
			//Delete users from site
			updateResult.User.SiteId = siteId
			_, err = deleteSiteUser(tx, user.UserId)
		case DeleteTidUserOverride:
			foundUser, _ := GetUserForId(user.UserId)
			updateResult.User = *foundUser
			updateResult.Action = "Delete TID User Override"
			//Delete users from tid
			updateResult.User.TidId = tid
			_, err = deleteTidUser(tx, user.UserId)
		}

		//check the procedure response, continue to the next user if error
		if err != nil {
			// log failure reason
			updateResult.Result = DbTransactionResult{
				Success:      false,
				ErrorMessage: err.Error(),
			}
			updateResults = append(updateResults, updateResult)

			// Log actual sql error
			logging.Error(err.Error())
			updateResult.Result.SetError(errors.New(makeUserUpdateErrorsFriendly(err)))
			err = tx.Rollback()
			if err != nil {
				updateResult.Result.SetError(errors.New("Unable to start transaction " + RollbackFailedMessage))
				updateResults = append(updateResults, updateResult)
				continue
			}
			updateResults = append(updateResults, updateResult)
			continue
		}

		//commit the transaction, continue to the next user if error
		err = tx.Commit()
		if err != nil {
			err = tx.Rollback()
			if err != nil {
				updateResult.Result.SetError(errors.New("Failed to commit txn " + RollbackFailedMessage))
				updateResults = append(updateResults, updateResult)
				continue
			}
			updateResult.Result.SetError(errors.New("Failed to commit txn"))
			updateResults = append(updateResults, updateResult)
			continue
		}

		//create the update result for this insert
		updateResults = append(updateResults, updateResult)
	}
	return updateResults, nil
}

/*
Assigns a user friendly error message for display on TMS, additional error messages can be added as needed.
MySQL error codes can be found here: https://dev.mysql.com/doc/refman/8.0/en/server-error-reference.html
*/
func makeUserUpdateErrorsFriendly(err error) string {
	switch err.(*mysql.MySQLError).Number {
	case DuplicateEntryError:
		return DuplicateUserText
	default:
		return "Unable to start transaction"
	}
}

// Run the user to site procedure
func addSiteUsers(transaction *sql.Tx, userId int, siteId int, username string, pin string, modulesString string, isEncrypted bool) (*sql.Rows, error) {
	return transaction.Query("CALL add_site_user(?,?,?,?,?,?)", userId, siteId, username, pin, modulesString, isEncrypted)
}

// Run the user to site tid
func addTidUser(transaction *sql.Tx, tidUserId int, tid int, username string, pin string, modulesString string, isEncrypted bool) (*sql.Rows, error) {
	return transaction.Query("CALL add_tid_user(?,?,?,?,?,?)", tidUserId, tid, username, pin, modulesString, isEncrypted)
}

// Run the user to site tid
func deleteTidUser(transaction *sql.Tx, tidUserId int) (sql.Result, error) {
	return transaction.Exec("CALL delete_tid_user(?)", tidUserId)
}

// Run the user to site tid
func deleteSiteUser(transaction *sql.Tx, userId int) (sql.Result, error) {
	return transaction.Exec("CALL delete_site_user(?)", userId)
}

func ClearTidUsers(tid int) error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	_, err = db.Exec("CALL clear_tid_users(?)", tid)
	if err != nil {
		return err
	}

	return nil

}

func GetMinimumRequiredSoftwareVersionForSite(siteId int) (string, error) {
	db, err := GetDB()
	if err != nil {
		return "", err
	}

	rows, err := db.Query("CALL required_software_version_fetch(?,?)", siteId, 0)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var version sql.NullString
	if rows.Next() {
		if err := rows.Scan(&version); err != nil {
			return "", err
		}
	}

	return version.String, nil
}

func ExportSiteUsers(siteId int) (*[]entities.SiteUserExportModel, error) {
	exportModels := make([]entities.SiteUserExportModel, 0)

	siteUsers, err := GetUsersForSite(siteId)
	if err != nil {
		logging.Error(err)
		return nil, err
	}
	for _, siteUser := range siteUsers {
		moduleMap := make(map[string]bool, 0)
		for _, module := range siteUser.Modules {
			moduleMap[module] = true
		}
		exportModels = append(exportModels, entities.SiteUserExportModel{
			Username: siteUser.Username,
			PIN:      siteUser.PIN,
			Modules:  moduleMap,
			Tid:      ""})
	}

	tidUsers, err := GetUsersForTidFromSiteId(siteId)
	if err != nil {
		logging.Error(err)
		return nil, err
	}
	for _, tidUser := range tidUsers {
		moduleMap := make(map[string]bool, 0)
		for _, module := range tidUser.Modules {
			moduleMap[module] = true
		}

		exportModels = append(exportModels, entities.SiteUserExportModel{
			Username: tidUser.Username,
			PIN:      tidUser.PIN,
			Modules:  moduleMap,
			Tid:      GetPaddedTidId(tidUser.TidId),
		})
	}

	return &exportModels, nil
}

// TODO: This method isnt used - perhaps delete it?
func GetDataValue(profileId int, elementId int) (string, error) {
	db, err := GetDB()
	if err != nil {
		return "", err
	}

	isEncrypted := false
	isPassword := false

	rows, err := db.Query("CALL get_element_value(?,?)", profileId, elementId)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	var value string
	for rows.Next() {
		err = rows.Scan(&value, &isEncrypted, &isPassword)
		if err != nil {
			return "", err
		}

		//this is for handling existing clear values
		clearValue := value
		if isEncrypted && value != "" {
			value, err = crypt.Decrypt(value)
			if err != nil {
				logging.Error("GetDataValue crypt.Decrypt error : " + err.Error())
				value = clearValue
			}
		}
	}

	return value, nil
}

func GetSiteNameFromMerchantID(merchantId string) (string, error) {
	db, err := GetDB()
	if err != nil {
		return "", err
	}

	rows, err := db.Query("CALL Get_site_name_from_merchantID(?)", merchantId)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	var siteName string
	for rows.Next() {
		err = rows.Scan(&siteName)
		if err != nil {
			return "", err
		}
	}
	return siteName, nil
}

func GetProfileIDFromMerchantID(merchantId string) (string, error) {
	db, err := GetDB()
	if err != nil {
		return "", err
	}

	rows, err := db.Query("call get_profile_id_from_mid(?)", merchantId)

	if err != nil {
		return "", err
	}
	defer rows.Close()

	var profileId string
	for rows.Next() {
		err = rows.Scan(&profileId)
		if err != nil {
			return "", err
		}
	}
	return profileId, nil
}

func GetSiteIDFromMerchantID(merchantId string) (string, error) {
	db, err := GetDB()
	if err != nil {
		return "", err
	}

	rows, err := db.Query("call get_site_id_from_mid(?)", merchantId)

	if err != nil {
		return "", err
	}
	defer rows.Close()

	var profileId string
	for rows.Next() {
		err = rows.Scan(&profileId)
		if err != nil {
			return "", err
		}
	}
	return profileId, nil
}

func GetFile(fileName string, fileserverURL string) template.URL {
	if fileName == "" {
		return ""
	}

	values := make(map[string][]string, 0)
	values["FileName"] = []string{fileName}

	result, err := http.PostForm(fileserverURL+"/getFile", values)
	if err != nil {
		return ""
	}
	defer result.Body.Close()
	if result.StatusCode != http.StatusOK {
		return ""
	}

	bodyBytes, err := ioutil.ReadAll(result.Body)
	if err != nil {
		logging.Error(err)
	}
	bodyString := string(bodyBytes)
	return template.URL(bodyString)
}

func GetDataElement(groupName string, elementName string) (sql.NullString, sql.NullString, bool, string, string, int, error) {
	db, err := GetDB()

	var valExp sql.NullString
	var valMsg sql.NullString
	var isAllowEmpty bool
	var dataType string
	var options string
	var dataElementID int
	if err != nil {
		return valExp, valMsg, isAllowEmpty, dataType, options, dataElementID, err
	}

	rows, err := db.Query(`SELECT de.validation_expression,de.validation_message,de.is_allow_empty,de.datatype,de.options,de.data_element_id
								FROM data_element de
								INNER JOIN data_group dg ON de.data_group_id = dg.data_group_id
								WHERE de.name = ? AND dg.name = ?;`, elementName, groupName)
	if err != nil {
		return valExp, valMsg, isAllowEmpty, dataType, options, dataElementID, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&valExp, &valMsg, &isAllowEmpty, &dataType, &options, &dataElementID)
		if err != nil {
			return valExp, valMsg, isAllowEmpty, dataType, options, dataElementID, err
		}
	}

	return valExp, valMsg, isAllowEmpty, dataType, options, dataElementID, nil
}

func GetDataElementByName(groupName string, elementName string) (int, error) {
	db, err := GetDB()
	if err != nil {
		return 0, err
	}

	rows, err := db.Query(`SELECT de.data_element_id
								FROM data_element de
								INNER JOIN data_group dg ON de.data_group_id = dg.data_group_id
								WHERE de.name = ? AND dg.name = ?`, elementName, groupName)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	var dataElementID int
	for rows.Next() {
		err = rows.Scan(&dataElementID)
		if err != nil {
			return 0, err
		}
	}
	return dataElementID, nil
}

func GetDataAllElementID() (map[string]models.DataElementsAndGroup, error) {
	dataElementMap := make(map[string]models.DataElementsAndGroup, 0)
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(`CALL get_all_data_elements_and_group_name()`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var e models.DataElementsAndGroup
	for rows.Next() {
		err = rows.Scan(&e.DataElementID, &e.DataElementName, &e.Options, &e.DataGroupName)
		if err != nil {
			return nil, err
		}

		dataElementMap[e.DataGroupName+"-"+e.DataElementName] = e
	}
	return dataElementMap, nil
}

func userModulesToString(u entities.SiteUser) string {
	return strings.Join(u.Modules, ",")
}

func CleanseJSON(jsonString string) string {
	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, []byte(jsonString)); err != nil {
		return jsonString
	}

	return string(buffer.Bytes())
}

func CheckThatMidExists(MID string) (bool, resultCodes.ResultCode) {
	if dataGroupsMap["store"].DataElements["merchantNo"] == 0 || dataGroupsMap["dualCurrency"].DataElements["secondaryMid"] == 0 {
		err := InitConstants()
		if err != nil {
			_, _ = logging.Error("An error occurred while calling InitConstants() - %s", err.Error())
			return false, resultCodes.DATA_GROUPS_MAPPINGS
		}
	}

	db, err := GetDB()
	if err != nil {
		_, _ = logging.Error("An error occurred executing GetDB() - %s", err.Error())
		return false, resultCodes.DATABASE_CONNECTION_ERROR
	}
	log.Println("MID", MID, "store", dataGroupsMap["store"].DataElements["merchantNo"], "dualCurrency", dataGroupsMap["dualCurrency"].DataElements["secondaryMid"])
	rows, err := db.Query("CALL GET_MID_BY_MID(?, ?, ?)", MID, dataGroupsMap["store"].DataElements["merchantNo"], dataGroupsMap["dualCurrency"].DataElements["secondaryMid"])
	if err != nil {
		_, _ = logging.Error("An error occurred executing procedure GET_MID_BY_MID - %s", err.Error())
		return false, resultCodes.DATABASE_QUERY_ERROR
	}
	defer rows.Close()

	var foundMid int
	var midType string
	var primaryMid int
	for rows.Next() {
		err := rows.Scan(&foundMid, &midType, &primaryMid)
		if err != nil {
			_, _ = logging.Error("An error occurred while scan rows", err.Error())
			return false, 0
		}
	}
	log.Println("foundMid", foundMid, "midType", midType, "primaryMid", primaryMid)
	if foundMid == 0 {
		return false, resultCodes.MID_DOES_NOT_EXIST
	} else {
		if midType == "primaryMid" {

			return true, resultCodes.MID_NOT_UNIQUE_PRIMARY_MID_DUPLICATE
		} else {
			return true, resultCodes.MID_NOT_UNIQUE_SECONDARY_MID_DUPLICATE
		}
	}
}

type FlagStatus struct {
	TID        int
	FlagStatus bool
}

func GetFlagStatus(siteID string) ([]FlagStatus, error) {
	db, err := GetDB()
	if err != nil {
		logging.Error("An error occurred executing GetDB() - %s", err.Error())
		return nil, err
	}

	rows, err := db.Query(`SELECT t.tid_id, t.flag_status FROM tid t 
		LEFT JOIN tid_site ts ON t.tid_id = ts.tid_id 
        WHERE ts.site_id =  ?`, siteID)
	if err != nil {
		logging.Error("Error thrown executing DB query", err.Error())
		return nil, err
	}
	defer rows.Close()

	var flagStatus []FlagStatus

	for rows.Next() {
		var fs FlagStatus
		err = rows.Scan(&fs.TID, &fs.FlagStatus)
		if err != nil {
			logging.Error("Error thrown attempting to scan current row.", err.Error())
			return nil, err
		}
		flagStatus = append(flagStatus, fs)
	}

	return flagStatus, nil
}

func GetFlaggedTids(siteID string) ([]string, error) {
	db, err := GetDB()
	if err != nil {
		logging.Error("An error occurred executing GetDB() - %s", err.Error())
		return nil, err
	}

	rows, err := db.Query(`SELECT t.tid_id FROM tid t 
		LEFT JOIN tid_site ts ON t.tid_id = ts.tid_id 
        WHERE ts.site_id =  ? and t.flag_status = ?`, siteID, true)
	if err != nil {
		logging.Error("Error thrown executing DB query", err.Error())
		return nil, err
	}
	defer rows.Close()

	var tids []string

	for rows.Next() {
		var tid string
		err = rows.Scan(&tid)
		if err != nil {
			logging.Error("Error thrown attempting to scan current row.", err.Error())
			return nil, err
		}
		tids = append(tids, tid)
	}

	return tids, nil
}

func GetFlagStatusUsingProfileID(profileID string) ([]FlagStatus, error) {
	db, err := GetDB()
	if err != nil {
		logging.Error("An error occurred executing GetDB() - %s", err.Error())
		return nil, err
	}

	rows, err := db.Query(`SELECT t.tid_id, t.flag_status FROM tid t 
		LEFT JOIN tid_site ts ON t.tid_id = ts.tid_id 
        LEFT JOIN site_profiles sp ON ts.site_id = sp.site_id
        WHERE sp.profile_id =  ?`, profileID)

	if err != nil {
		logging.Error("Error thrown executing DB query", err.Error())
		return nil, err
	}
	defer rows.Close()

	var flagStatus []FlagStatus

	for rows.Next() {
		var fs FlagStatus
		err = rows.Scan(&fs.TID, &fs.FlagStatus)
		if err != nil {
			logging.Error("Error thrown attempting to scan current row.", err.Error())
			return nil, err
		}
		flagStatus = append(flagStatus, fs)
	}

	return flagStatus, nil
}

func GetSitePaymentServices(siteId int) (map[int]*PaymentService, error) {
	db, err := GetDB()
	if err != nil {
		logging.Error("An error occurred executing GetDB() - %s", err.Error())
		return nil, err
	}

	servicesMap := make(map[int]*PaymentService)
	rowServices, err := db.Query("call get_site_profile_payment_services(?)", siteId)
	if err != nil {
		logging.Error("Error thrown executing DB query", err.Error())
		return nil, err
	}
	defer rowServices.Close()

	for rowServices.Next() {
		var paymentService PaymentService
		err = rowServices.Scan(&paymentService.ServiceId, &paymentService.Name)
		if err != nil {
			logging.Error("Error thrown executing DB query", err.Error())
			return nil, err
		}
		paymentService.TID = ""
		paymentService.MID = ""
		servicesMap[paymentService.ServiceId] = &paymentService
	}

	return servicesMap, nil
}

func CheckMandatoryDataElement(dataElementId int) (int, error) {
	db, err := GetDB()
	if err != nil {
		logging.Error("An error occurred executing GetDB() - %s", err.Error())
		return -1, err
	}
	var isAllowEmpty int
	if err := db.QueryRow("SELECT is_allow_empty FROM data_element WHERE data_element_id = ?", dataElementId).Scan(&isAllowEmpty); err != nil {
		return -1, err
	}

	return isAllowEmpty, nil
}

func CheckAcquirerLevelRequiredDataElement(dataElementId int) (int, error) {
	db, err := GetDB()
	if err != nil {
		logging.Error("An error occurred executing GetDB() - %s", err.Error())
		return -1, err
	}
	var isRequiredAtAcquireLevel int
	if err := db.QueryRow("SELECT required_at_acquirer_level FROM data_element WHERE data_element_id = ?", dataElementId).Scan(&isRequiredAtAcquireLevel); err != nil {
		return -1, err
	}

	return isRequiredAtAcquireLevel, nil
}

func GetSiteHistoryData(profileId int) ([]*ProfileChangeHistory, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	tidChangeHistory, err := GetSiteTidsChangeHistory(db, profileId)
	if err != nil {
		return nil, err
	}

	changeHistory, err := GetProfileChangeHistory(profileId)
	if err != nil {
		return nil, err
	}

	for _, history := range tidChangeHistory {
		duplicate := false
		for _, ch := range changeHistory {
			if ch.ChangeType == history.ChangeType &&
				ch.OriginalValue == history.OriginalValue &&
				ch.ChangeValue == history.ChangeValue &&
				ch.ChangedBy == history.ChangedBy &&
				ch.ChangedAt == history.ChangedAt &&
				ch.Approved == history.Approved &&
				ch.TidId == history.TidId &&
				ch.IsEncrypted == history.IsEncrypted &&
				ch.IsPassword == history.IsPassword {
				duplicate = true
				break
			}
		}
		if !duplicate {
			changeHistory = append(changeHistory, history)
		}
	}

	// remove duplicate history items
	keys := map[string]bool{}
	var uniqueHistory []*ProfileChangeHistory
	for i, entry := range changeHistory {
		key := fmt.Sprintf("%s,%s,%s,%s,%s,%d,%d", entry.Field, entry.TidId, entry.ChangeValue, entry.ChangedBy, entry.ChangedAt, entry.ChangeType, entry.Approved)
		if _, value := keys[key]; !value {
			keys[key] = true
			uniqueHistory = append(uniqueHistory, changeHistory[i])
		}
	}
	return uniqueHistory, nil
}

func GetSiteLevelProfileChangeHistory(profileID, siteID, pageSize, pageNumber, offset, limit int) ([]*ProfileChangeHistory, TidPagination, error) {
	const getSiteLevelChangeHistoryCall = "get_site_level_change_history"
	return GetSiteLevelChangeHistory(profileID, siteID, pageSize, pageNumber, offset, limit, getSiteLevelChangeHistoryCall)
}

func GetSiteLevelChangeHistory(id, siteID, pageSize, pageNumber, offset, limit int, call string) ([]*ProfileChangeHistory, TidPagination, error) {
	var changes = make([]*ProfileChangeHistory, 0)
	var pagination TidPagination

	db, err := GetDB()
	if err != nil {
		return nil, pagination, err
	}

	var totalCount int
	err = db.QueryRow("CALL get_site_level_change_history_count(?, ?)", id, siteID).Scan(&totalCount)
	if err != nil {
		return nil, pagination, err
	}

	rows, err := db.Query("CALL "+call+"(?,?,?,?)", id, siteID, offset, limit)
	if err != nil {
		return nil, pagination, err
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
			re := regexp.MustCompile(`"PIN":"\d+"`)
			changeHistory.OriginalValue = re.ReplaceAllString(originalVal.String, `"PIN":"*****"`)
			changeHistory.OriginalValue = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(changeHistory.OriginalValue, `"SiteId":0`, ``), `"tidId":-1`, ``), `"TidId":0`, ``)

		}

		if updatedValue.Valid {
			re := regexp.MustCompile(`"PIN":"\d+"`)
			changeHistory.ChangeValue = re.ReplaceAllString(updatedValue.String, `"PIN":"*****"`)
			changeHistory.ChangeValue = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(changeHistory.ChangeValue, `"SiteId":0`, ``), `"tidId":-1`, ``), `"TidId":0`, ``)

		}

		if tidId.Valid {
			changeHistory.TidId = tidId.String
		}

		if isEncrypted.Valid {
			if isEncrypted.Bool {
				if originalVal.Valid {
					changeHistory.OriginalValue, err = crypt.Decrypt(changeHistory.OriginalValue)
					if err != nil {
						return nil, pagination, err
					}
				}
				if updatedValue.Valid {
					changeHistory.ChangeValue, err = crypt.Decrypt(changeHistory.ChangeValue)
					if err != nil {
						return nil, pagination, err
					}
				}
			}

			changeHistory.IsEncrypted = isEncrypted.Bool
		}

		if isPassword.Valid {
			changeHistory.IsPassword = isPassword.Bool
		}

		if err != nil {
			return nil, pagination, err
		}

		changeHistory.RowNo = rowCount
		rowCount++

		changes = append(changes, changeHistory)
	}

	pagination = buildPagination(pageSize, totalCount, pageNumber)

	return changes, pagination, nil
}

func buildPagination(pageSize, total, pageNumber int) TidPagination {

	var pagination TidPagination

	// Calculate the total number of pages
	pageCount := int(math.Ceil(float64(total) / float64(pageSize)))

	if pageCount == 1 || pageCount == 0 { // If there is only one page then there is no next or previous page to view
		pagination.Less = false
		pagination.More = false

		var page PaginationPage
		page.PageNumber = strconv.Itoa(pageNumber)
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

	if pageSize > total || pageSize == -1 {
		if total == 0 {
			pagination.FirstRecord = 0
		} else {
			pagination.FirstRecord = 1
		}
		pagination.LastRecord = total
		pagination.TotalRecords = total
	}

	pagination.CurrentPage = pageNumber
	pagination.PageSize = pageSize
	pagination.PageCount = pageCount

	return pagination
}

// GetDetailsByProfileID used to get the profileName, ProfileType, SiteId
func GetDetailsByProfileID(profileID int) (string, string, sql.NullInt64, error) {
	var profileName, profileType string
	var siteID sql.NullInt64
	db, err := GetDB()
	if err != nil {
		return "", "", sql.NullInt64{}, err
	}

	rows, err := db.Query(`
    SELECT
        p.name AS profile_name,
        pt.name AS profile_type,
        sp.site_id
    FROM
        profile AS p
    LEFT JOIN profile_type AS pt ON p.profile_type_id = pt.profile_type_id
    LEFT JOIN site_profiles AS sp ON p.profile_id = sp.profile_id
    WHERE
        p.profile_id = ?;
`, profileID)
	if err != nil {
		return "", "", sql.NullInt64{}, err
	}

	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&profileName, &profileType, &siteID); err != nil {
			return "", "", sql.NullInt64{}, err
		}
	}

	return profileName, profileType, siteID, nil
}

func getUsersFraudDetailsForApproval(profileId int, profileType string, isTidUsersOverride bool) (string, error) {
	siteId, err := GetSiteIdFromProfileId(profileId)
	if err != nil {
		return "", err
	}

	tID := -1
	if siteId == -4 {
		tID, siteId, err = GetTIdSiteIdFromProfileId(profileId)
		if err != nil {
			return "", err
		}
	}

	db, err := GetDB()
	if err != nil {
		return "", err
	}

	if profileType == "users" {
		if isTidUsersOverride {
			rows, err := db.Query("SELECT tid_user_id, Username, PIN , Modules from tid_user_override WHERE tid_id=(?)", uint(tID))
			if err != nil {
				_, _ = logging.Error("Error thrown attempting to @get_users_value_by_site_id.", err.Error())
				return "", err
			}
			defer rows.Close()

			var UserId int
			var Username string
			var pin string
			var Modules string
			var userList = make([]*entities.SiteUserData, 0)
			for rows.Next() {
				if err = rows.Scan(&UserId, &Username, &pin, &Modules); err != nil {
					logging.Error("Error thrown attempting to get users value.", err.Error())
					return "", err
				}
				userList = append(userList, &entities.SiteUserData{UserId: UserId, Username: Username, PIN: "*****", Modules: strings.Split(Modules, ",")})
			}

			usersData, err := json.Marshal(userList)
			if err != nil {
				return "", err
			}

			return string(usersData), nil
		} else {
			rows, err := db.Query("CALL get_users_value_by_site_id(?,?,?,?)", siteId, profileId, "", tID)
			if err != nil {
				_, _ = logging.Error("Error thrown attempting to @get_users_value_by_site_id.", err.Error())
				return "", err
			}
			defer rows.Close()

			var UserId int
			var Username string
			var pin string
			var Modules string
			var userList = make([]*entities.SiteUserData, 0)
			for rows.Next() {
				if err = rows.Scan(&UserId, &Username, &pin, &Modules); err != nil {
					logging.Error("Error thrown attempting to get users value.", err.Error())
					return "", err
				}
				userList = append(userList, &entities.SiteUserData{UserId: UserId, Username: Username, PIN: "*****", Modules: strings.Split(Modules, ",")})
			}

			usersData, err := json.Marshal(userList)
			if err != nil {
				return "", err
			}

			return string(usersData), nil
		}
	} else {
		rows, err := db.Query("CALL get_users_value_by_site_id(?,?,?,?)", siteId, profileId, "dailyTxnCleanseTime", tID)
		if err != nil {
			logging.Error("Error thrown attempting to @get_users_value_by_site_id.", err.Error())
			return "", err
		}
		defer rows.Close()

		var dailyTxnCleanseTime string
		for rows.Next() {
			if err = rows.Scan(&dailyTxnCleanseTime); err != nil {
				logging.Error("Error thrown attempting to get users value.", err.Error())
				return "", err
			}
		}

		siteRows, err := db.Query("CALL get_users_value_by_site_id(?,?,?,?)", siteId, profileId, "siteVelocity", tID)
		if err != nil {
			logging.Error("Error thrown attempting to @get_users_value_by_site_id.", err.Error())
			return "", err
		}
		defer siteRows.Close()

		var velocityLimits []entities.VelocityLimit
		for siteRows.Next() {
			var limits entities.VelocityLimit
			if err = siteRows.Scan(&limits.ID, &limits.Scheme, &limits.DailyCount, &limits.BatchCount, &limits.SingleTransLimit, &limits.DailyLimit, &limits.BatchLimit, &limits.Index); err != nil {
				logging.Error("Error thrown attempting to get users value.", err.Error())
				return "", err
			}

			txnLimits, err := getTxnVelocityLimits(limits.ID)
			if err != nil {
				logging.Error("Error thrown attempting to @getTxnVelocityLimits.", err.Error())
				return "", err
			}

			limits.TxnLimits = txnLimits
			velocityLimits = append(velocityLimits, limits)
		}

		velocity, err := json.Marshal(velocityLimits)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("{\"siteID\":%d, \"tidId\":%d, \"limitString\":%s, \"dailyTxnCleanseTime\":\"%s\",\"velocityLimits\":%s}", siteId, -1, "3", dailyTxnCleanseTime, string(velocity)), nil
	}
}

func GetProfileIDFromSiteName(name string, profileTypeID int) (string, error) {
	db, err := GetDB()
	if err != nil {
		_, _ = logging.Error(err)
		return "", err
	}

	row := db.QueryRow("select GROUP_CONCAT(pd.profile_id) as profile_ids_list from profile_data pd LEFT JOIN `profile` p ON p.profile_id = pd.profile_id where pd.datavalue = ? AND p.profile_type_id = ? ", name, profileTypeID)
	var profileIDsList string
	err = row.Scan(&profileIDsList)
	if err != nil {
		_, _ = logging.Error(err)
		return "", err
	}
	return profileIDsList, err
}
