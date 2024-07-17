package dal

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"nextgen-tms-website/config"
	"nextgen-tms-website/resultCodes"
	"time"

	rpcHelp "bitbucket.org/network-international/nextgen-libs/nextgen-helpers/rpcHelper"
)

const RollbackFailedMessage string = "Rollback failed (Please contact IT)."

type DataGroupType struct {
	Id           int
	DataElements map[string]int
}

var dataGroupsMap map[string]DataGroupType

func InitConstants() (err error) {
	db, err := GetDB()
	if err != nil {
		logging.Error(err)
		return
	}
	if dataGroupsMap, err = getDataGroupMap(db); err != nil {
		logging.Error(err)
		return
	} else if len(dataGroupsMap) <= 0 {
		err = errors.New("no Data Groups Found in Map")
		logging.Error(err)
		return
	}
	return
}

func getDataGroupMap(db *sql.DB) (map[string]DataGroupType, error) {
	rows, err := db.Query(`SELECT dg.data_group_id, dg.name, de.data_element_id, de.name 
                           FROM data_group dg
                           INNER JOIN data_element de on dg.data_group_id = de.data_group_id`)
	if err != nil {
		logging.Error(err)
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logging.Error(err)
		}
	}()

	dataGroupsObj := make(map[string]DataGroupType)
	for rows.Next() {
		var dgId int
		var dgName string
		var deId int
		var deName string
		err = rows.Scan(&dgId, &dgName, &deId, &deName)
		if err != nil {
			logging.Error(err)
			return nil, err
		}
		if _, present := dataGroupsObj[dgName]; !present {
			dataGroupsObj[dgName] = DataGroupType{
				Id:           dgId,
				DataElements: make(map[string]int),
			}
		}
		dataGroupsObj[dgName].DataElements[deName] = deId
	}
	return dataGroupsObj, nil
}

// Searches for a data element within a list of data groups
// return arg 1 bool: true if the data element was found within a data group
// return arg 2 int: the index of the data group that the data element belongs to if the element is present, or -1 if not
// return arg 3 int: the index of the data element within the data group that the data element belongs to if the element is present, or -1 if not
func DataGroupsContainsDataElement(groups []*DataGroup, elementId int) (dataElementFound bool, dataGroupIndex int, dataElementIndex int) {
	for i, group := range groups {
		for j, element := range group.DataElements {
			if element.ElementId == elementId {
				return true, i, j
			}
		}
	}

	return false, -1, -1
}

func RemoveDataElementIfPresentInGroups(dataGroups []*DataGroup, elementId int) []*DataGroup {
	present, dataGroupIndex, dataElementIndex := DataGroupsContainsDataElement(dataGroups, elementId)
	if !present {
		return dataGroups
	}
	removeDataElement := func(s []DataElement, i int) []DataElement {
		s[len(s)-1], s[i] = s[i], s[len(s)-1]
		return s[:len(s)-1]
	}
	removeGroup := func(s []*DataGroup, i int) []*DataGroup {
		s[len(s)-1], s[i] = s[i], s[len(s)-1]
		return s[:len(s)-1]
	}

	dataGroups[dataGroupIndex].DataElements = removeDataElement(dataGroups[dataGroupIndex].DataElements, dataElementIndex)
	if len(dataGroups[dataGroupIndex].DataElements) == 0 {
		dataGroups = removeGroup(dataGroups, dataGroupIndex)
	}
	return dataGroups
}

type MappedDataGroups map[string]MappedDataGroup
type MappedDataGroup struct {
	DataGroupID  int
	DataGroup    string
	DisplayName  string
	DataElements map[string]DataElement
}

func (mdg *MappedDataGroups) toArrays() []*DataGroup {
	groups := make([]*DataGroup, len(*mdg))
	i := 0
	for _, dg := range *mdg {
		groups[i] = &DataGroup{
			DataGroupID:  dg.DataGroupID,
			DataGroup:    dg.DataGroup,
			DisplayName:  dg.DisplayName,
			DataElements: make([]DataElement, len(dg.DataElements)),
		}
		j := 0
		for _, de := range dg.DataElements {
			groups[i].DataElements[j] = DataElement{
				ElementId:                de.ElementId,
				Name:                     de.Name,
				Type:                     de.Type,
				IsAllowedEmpty:           de.IsAllowedEmpty,
				DataValue:                de.DataValue,
				MaxLength:                de.MaxLength,
				ValidationExpression:     de.ValidationExpression,
				ValidationMessage:        de.ValidationMessage,
				FrontEndValidate:         de.FrontEndValidate,
				Unique:                   de.Unique,
				Overriden:                de.Overriden,
				Options:                  de.Options,
				OptionSelectable:         de.OptionSelectable,
				DisplayName:              de.DisplayName,
				Image:                    de.Image,
				IsPassword:               de.IsPassword,
				IsEncrypted:              de.IsEncrypted,
				SortOrderInGroup:         de.SortOrderInGroup,
				IsNotOverrideable:        de.IsNotOverrideable,
				OverridePriority:         de.OverridePriority,
				Tooltip:                  de.Tooltip,
				FileMaxSize:              de.FileMaxSize,
				FileMinRatio:             de.FileMinRatio,
				FileMaxRatio:             de.FileMaxRatio,
				IsReadOnlyAtCreation:     de.IsReadOnlyAtCreation,
				IsRequiredAtAcquireLevel: de.IsRequiredAtAcquireLevel,
				IsRequiredAtChainLevel:   de.IsRequiredAtChainLevel,
			}
			j++
		}
		i++
	}
	return groups
}

func (mdg *MappedDataGroups) AddElementFromSiteData(sd SiteData) {
	if *mdg == nil {
		*mdg = make(map[string]MappedDataGroup, 0)
	}
	if _, exists := (*mdg)[sd.DataGroup]; !exists {
		(*mdg)[sd.DataGroup] = MappedDataGroup{
			DataGroupID:  sd.DataGroupID,
			DataGroup:    sd.DataGroup,
			DisplayName:  sd.DataGroupDisplayName,
			DataElements: make(map[string]DataElement),
		}
	}
	if elem, exists := (*mdg)[sd.DataGroup].DataElements[sd.Name]; exists && sd.OverridePriority.Valid && int(sd.OverridePriority.Int32) == elem.OverridePriority {
		return
	}

	var maxLength int = -1
	var validationExpression = ""
	var validationMessage = ""
	var overridePriority = -1
	var tooltip = ""

	if sd.MaxLength.Valid {
		maxLength = int(sd.MaxLength.Int64)
	}

	if sd.ValidationExpression.Valid {
		validationExpression = sd.ValidationExpression.String
	}

	if sd.ValidationMessage.Valid {
		validationMessage = sd.ValidationMessage.String
	}

	if sd.OverridePriority.Valid {
		overridePriority = int(sd.OverridePriority.Int32)
	}

	if sd.Tooltip.Valid {
		tooltip = sd.Tooltip.String
	}

	var image template.URL
	if sd.DataType == "FILE" {
		image = GetFile(sd.DataValue.String, config.FileserverURL)
	}

	(*mdg)[sd.DataGroup].DataElements[sd.Name] = DataElement{
		ElementId:                sd.DataElementID,
		Name:                     sd.Name,
		Type:                     sd.DataType,
		DataValue:                sd.DataValue.String,
		MaxLength:                maxLength,
		ValidationExpression:     validationExpression,
		ValidationMessage:        validationMessage,
		FrontEndValidate:         sd.FrontEndValidate != 0,
		Overriden:                sd.Overriden.Valid && sd.Overriden.Int64 != 0,
		Options:                  sd.Options,
		OptionSelectable:         sd.OptionSelectable,
		DisplayName:              sd.DisplayName.String,
		Image:                    image,
		IsPassword:               sd.IsPassword,
		IsEncrypted:              sd.IsEncrypted.Valid && sd.IsEncrypted.Bool,
		SortOrderInGroup:         sd.SortOrderInGroup,
		IsNotOverrideable:        sd.IsNotOverridable.Valid && sd.IsNotOverridable.Bool,
		OverridePriority:         overridePriority,
		Tooltip:                  tooltip,
		IsAllowedEmpty:           sd.IsAllowEmpty,
		FileMaxSize:              int(sd.FileMaxSize.Int64),
		FileMinRatio:             sd.FileMinRatio.Float64,
		FileMaxRatio:             sd.FileMaxRatio.Float64,
		IsReadOnlyAtCreation:     sd.IsReadOnlyAtCreation,
		IsRequiredAtAcquireLevel: sd.IsRequiredAtAcquireLevel,
		IsRequiredAtChainLevel:   sd.IsRequiredAtChainLevel,
	}
}

// DataGroup struct
type DataGroup struct {
	DataGroupID  int
	DataGroup    string
	DisplayName  string
	PreSelected  bool
	IsSelected   bool
	DataElements []DataElement
}

// ProfileData struct
type ProfileData struct {
	ID     int
	Name   string
	TypeId int
	Type   string
}

// DataElement struct
type DataElement struct {
	ElementId                int
	Name                     string
	Type                     string
	IsAllowedEmpty           bool
	DataValue                string
	MaxLength                int
	ValidationExpression     string
	ValidationMessage        string
	FrontEndValidate         bool
	Unique                   bool
	Overriden                bool
	Options                  []OptionData
	OptionSelectable         bool
	DisplayName              string
	Image                    template.URL
	IsPassword               bool
	IsEncrypted              bool
	SortOrderInGroup         int
	IsNotOverrideable        bool
	OverridePriority         int // The priority of the profile type which is overriding this element. This refers to profile_type.priority
	Tooltip                  string
	FileMaxSize              int
	FileMinRatio             float64
	FileMaxRatio             float64
	IsReadOnlyAtCreation     bool
	IsRequiredAtAcquireLevel bool
	IsRequiredAtChainLevel   bool
}

type SiteList struct {
	SiteID            int
	SiteProfileID     int
	SiteName          string
	ChainProfileID    int
	ChainName         string
	AcquirerProfileID int
	AcquirerName      string
	GlobalProfileID   int
	GlobalName        string
	MerchantId        string
}

type ChainList struct {
	ChainProfileID int
	ChainName      string
	ChainTIDCount  int
	AcquirerName   string
}

type AcquirerList struct {
	AcquirerProfileID int
	AcquirerName      string
}

type SiteData struct {
	TIDID                    int // todo: get rid of this
	DataGroupID              int
	DataGroup                string
	DataGroupDisplayName     string
	DataElementID            int
	Tooltip                  sql.NullString
	Name                     string
	Source                   sql.NullString
	OverridePriority         sql.NullInt32
	DataValue                sql.NullString
	DataType                 string
	IsAllowEmpty             bool
	MaxLength                sql.NullInt64
	ValidationExpression     sql.NullString
	ValidationMessage        sql.NullString
	FrontEndValidate         int
	Overriden                sql.NullInt64
	Options                  []OptionData
	OptionSelectable         bool
	DisplayName              sql.NullString
	IsPassword               bool
	IsEncrypted              sql.NullBool
	SortOrderInGroup         int
	RequiredAtSiteLevel      bool
	IsNotOverridable         sql.NullBool
	FileMaxSize              sql.NullInt64
	FileMinRatio             sql.NullFloat64
	FileMaxRatio             sql.NullFloat64
	TIDOverridable           bool
	IsReadOnlyAtCreation     bool
	IsRequiredAtAcquireLevel bool
	IsRequiredAtChainLevel   bool
	ProfileID                int
	ProfileName              string
}

type OptionData struct {
	Option   string
	Selected bool
}

// TIDData struct
type TIDData struct {
	TID                int
	TIDProfileID       int
	Serial             string
	EnrolmentPIN       string
	ExpiryTime         string
	ResetPIN           string
	ResetPinExpiryTime string
	ActivationTime     string
	SiteId             int
	SiteName           string
	TIDProfileGroups   []*DataGroup
	Overridden         bool
	UserOverrides      bool
	MerchantID         string
	FraudOverride      bool
}

// PackageData struct
type PackageData struct {
	PackageID int
	Version   string
	Apks      []string
}

type ThirdPartyApk struct {
	ApkID   int
	Version string
	Apk     string
}

type BindThirdPartyTarget struct {
	ApkID           int
	ThirdPartyApkID string
}

type ProfileChangeHistory struct {
	RowNo         int
	Field         string
	ChangeType    int
	OriginalValue string
	ChangeValue   string
	ChangedBy     string
	ChangedAt     string
	Approved      int
	TidId         string
	IsEncrypted   bool
	IsPassword    bool
}

type ChangeApprovalHistory struct {
	ProfileDataID           string
	Identifier              string
	Field                   string
	ChangeType              int
	OriginalValue           string
	ChangeValue             string
	ChangedBy               string
	ChangedAt               string
	Approved                int
	ReviewedBy              string
	ReviewedAt              string
	TidId                   string
	MID                     string
	IsEncrypted             bool
	IsPassword              bool
	PaymentServiceGroupName string
	PaymentServiceName      string
}

type UserAuditHistory struct {
	ID            string
	Acquirer      string
	Name          string
	Module        string
	OriginalValue string
	UpdatedValue  string
	UpdatedBy     string
	UpdatedAt     string
}

type UserData struct {
	UserId   int
	Username string
	Selected bool
}

type ContactUsData struct {
	AcquirerName             string
	AcquirerPrimaryPhone     string
	AcquirerSecondaryPhone   string
	AcquirerEmail            string
	AcquirerAddressLineOne   string
	AcquirerAddressLineTwo   string
	AcquirerAddressLineThree string
	FurtherInformation       string
	Valid                    bool
	UserAcquirer             string
}

type PermissionsGroupsData struct {
	GroupId                  int
	Name                     string
	DefaultGroup             int
	Editable                 bool
	UserGroupSelected        bool
	GroupPermissionsSelected bool
}

type PermissionsData struct {
	PermissionId int
	Name         string
	Selected     bool
}

type PermissionsGroupAcquirerData struct {
	ProfileId int
	Name      string
	Selected  bool
}

type TidDetails struct {
	AppVer                           int
	FirmwareVer                      string
	LastTransaction                  string
	LastCheckedTime                  string
	ConfirmedTime                    string
	LastAPKDownloadTime              string
	EODAuto                          bool
	AutoTime                         string
	Coordinates                      string
	Accuracy                         string
	LastCoordinateTime               string
	FreeInternalStorage              string
	TotalInternalStorage             string
	SoftuiLastDownloadedFileName     string
	SoftuiLastDownloadedFileHash     string
	SoftuiLastDownloadedFileList     []map[string]string
	SoftuiLastDownloadedFileDateTime string
}

type ToolTipExport struct {
	DisplayName string
	Name        string
	DataGroup   string
	ToolTip     string
	UpdateTime  string
	UpdateBy    string
}

type DbTransactionResult struct {
	Success      bool
	ErrorMessage string
}

type PaymentServiceGroup struct {
	GroupID  int
	Name     string
	Enabled  bool
	AuthType string
	Services []PaymentService
}

type PaymentService struct {
	ServiceId int
	Name      string
	MID       string
	TID       string
}

func (res *DbTransactionResult) SetError(err error) {
	res.Success = false
	res.ErrorMessage = err.Error()
}
func (res *DbTransactionResult) SetErrorByString(err string) {
	res.Success = false
	res.ErrorMessage = err
}

const (
	Others   = 6
	Tid      = 5
	Site     = 4
	Chain    = 3
	Acquirer = 2
	History  = -1
)

var (
	dataBase         *sql.DB
	logging          rpcHelp.LoggingClient
	MaxDBConnections int
)

var connectionString string
var backupLocation string

func SetConnectionString(s string) {
	connectionString = s
}

func SetBackupLocation(s string) {
	backupLocation = s
}

func Connect(l rpcHelp.LoggingClient) {
	logging = l
	logging.Information("Initialising Connection to Database")

	// sql.Open only validates the driverName, not the connection string ∴ sql.Ping needs to be called to check connection
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		criticalError := rpcHelp.BuildSQLCritError(err.Error())
		logging.Error(criticalError)
		return
	}

	dataBase = db
	if MaxDBConnections > 0 {
		dataBase.SetMaxOpenConns(MaxDBConnections)
	}

	err = db.Ping()
	if err != nil {
		criticalError := rpcHelp.BuildSQLCritError(err.Error())
		logging.Error(criticalError)
		return
	}

	logging.Information("Connection Successful")
}

func ConnectAncillaryScripts(l rpcHelp.LoggingClient) (*sql.DB, error) {
	logging = l
	logging.Information("Initialising Ancillary Connection to Database")
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		criticalError := rpcHelp.BuildSQLCritError(err.Error())
		logging.Error(criticalError)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		criticalError := rpcHelp.BuildSQLCritError(err.Error())
		logging.Error(criticalError)
		return nil, err
	}
	logging.Information("Connection Successful")

	return db, nil
}

// Returns a pointer to a sql.DB object. NB as per its documentation, it
// represents a *pool* of connections, and is safe to be used by multiple
// goroutines.
//
// In general the returned DB should NOT be closed; only close the individual connections
// created, e.g. returned "rows" instances returned from queries.
func GetDB() (*sql.DB, error) {
	err := dataBase.Ping()
	if err != nil {
		logging.Information("Connection to database lost, attempting to re-connect...")
		dataBase.Close()
		dataBase, err = sql.Open("mysql", connectionString)
		if err != nil {
			criticalError := rpcHelp.BuildSQLCritError(err.Error())
			logging.Error(criticalError)
			return nil, err
		}

		if MaxDBConnections > 0 {
			dataBase.SetMaxOpenConns(MaxDBConnections)
		}

		err = dataBase.Ping()
		if err != nil {
			criticalError := rpcHelp.BuildSQLCritError(err.Error())
			logging.Error(criticalError)
			return nil, errors.New("error: unable to connect to MySQL database")
		}

		logging.Information("Connection Successful")
	}

	dbStats := dataBase.Stats()
	if dbStats.OpenConnections >= dbStats.MaxOpenConnections {
		return nil, errors.New("all database connections in use, please wait and try again")
	}

	return dataBase, nil
}

func CloseDB() {
	if dataBase.Ping() == nil {
		dataBase.Close()
	}
}

func SaveUser(username string, roleID int, firstLogon int) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	// TODO: Remove unused user_store procedure
	statement := `
	insert into user(
        username, roleId, updated_at, updated_by, created_at, created_by, first_logon
    ) values (
    	?, ?, current_timestamp, ?, current_timestamp, ?, ?)
	ON DUPLICATE KEY UPDATE
		username = ?, roleId = ?, updated_at = current_timestamp, updated_by = ?;`

	result, err := db.Exec(
		statement,
		username, roleID,
		username, // updated by
		username, // created by
		firstLogon,

		// On duplicate key update...
		username,
		roleID,
		// updated by
		username)

	if err != nil {
		return err
	}
	ra, err := result.RowsAffected()

	if ra == 0 {
		return fmt.Errorf("Unable to save user")
	}

	return nil
}

func DeleteUser(username string) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	statement := "DELETE FROM user WHERE username=?;"
	result, err := db.Exec(statement, username)
	if err != nil {
		return err
	}

	if ra, _ := result.RowsAffected(); ra == 0 {
		// error ignored explicitly as result should be able to return RA, if not we don't know if the call was a success anyway
		return errors.New("unable to delete user")
	}

	return nil
}

func CheckDBVersionMatch(dbVersion int) bool {
	db, err := GetDB()
	if err != nil {
		return false
	}

	versionMatch := true
	currentVersion := 0
	rows, err := db.Query("SELECT version FROM db_version")
	if err != nil {
		logging.Information(fmt.Sprintf("Database Version Mismatch, Expected: %d, Current: %d", dbVersion, currentVersion))
		versionMatch = false
	} else {
		defer rows.Close()
		for rows.Next() {
			rows.Scan(&currentVersion)
			if currentVersion != dbVersion {
				logging.Information(fmt.Sprintf("Database Version Mismatch, Expected: %d, Current: %d", dbVersion, currentVersion))
				versionMatch = false
			}
		}
	}
	return versionMatch
}

func AddTxnChecksum(filename string, checksum string) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	_, err = db.Exec("CALL add_uploaded_txn(?,?)", filename, checksum)
	if err != nil {
		return err
	}
	return nil
}

func CheckUploadedChecksum(checksum string) (bool, string, error) {
	db, err := GetDB()
	if err != nil {
		return false, "", err
	}

	rows, err := db.Query("CALL checkUploadedTxn(?)", checksum)
	if err != nil {
		return false, "", err
	}
	defer rows.Close()
	var conflictName sql.NullString
	if rows.Next() {
		rows.Scan(&conflictName)
		return false, conflictName.String, nil
	}

	return true, "", nil

}

func GetProfileDataForTabByProfileId(profileId int, tabName string) ([]*DataGroup, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	if tabName == "" {
		tabName = "site_configuration"
	}
	rows, err := db.Query("CALL get_profile_data_for_tab_by_profile_id(?, ?)", profileId, tabName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var jsonResponse string
	for rows.Next() {
		err = rows.Scan(&jsonResponse)
		if err != nil {
			return nil, err
		}
	}

	//unmarshall the json into our object
	dataGroups := make([]*DataGroup, 0)
	err = json.Unmarshal([]byte(jsonResponse), &dataGroups)
	if err != nil {
		return nil, err
	}

	return dataGroups, nil
}

type ValidationDal interface {
	GetDataElementMetadata(dataElementId int, profileId int) (DataElement, error)
	GetDataElementByName(groupName string, elementName string) (int, error)
	CheckThatTidExists(TID int) (tidExists bool, resultCode resultCodes.ResultCode, overrideProfileId int)
	CheckThatSerialNumberExists(SN string) (snExists bool, err error)
	CheckThatMidExists(MID string) (midExists bool, resultCode resultCodes.ResultCode, profileId int)
	GetIsUnique(elementId int, elementValue string, profile int) (bool, error)
}

func NewValidationDal() ValidationDal {
	return new(validationDal)
}

type validationDal struct{}

func (v validationDal) CheckThatSerialNumberExists(SN string) (snExists bool, err error) {
	err = CheckSerialInUse(SN)
	return err != nil, err
}

func (v validationDal) GetDataElementMetadata(dataElementId int, profileId int) (DataElement, error) {
	return GetDataElementMetadata(dataElementId, profileId)
}

func (v validationDal) GetDataElementByName(groupName string, elementName string) (int, error) {
	return GetDataElementByName(groupName, elementName)
}

func (v validationDal) CheckThatTidExists(TID int) (tidExists bool, resultCode resultCodes.ResultCode, overrideProfileId int) {
	return CheckThatTidExists(TID)
}

func (v validationDal) CheckThatMidExists(MID string) (midExists bool, resultCode resultCodes.ResultCode, profileId int) {
	midExists, resultCode = CheckThatMidExists(MID)
	profileId = -1
	return
}

func (v validationDal) GetIsUnique(elementId int, elementValue string, profile int) (bool, error) {
	return GetIsUnique(elementId, elementValue, profile)
}

func ExportTooltips(w http.ResponseWriter) error {
	logging.Debug("Export Tooltip request received.")
	// Grab the DB
	db, err := GetDB()
	if err != nil {
		logging.Error("Error thrown obtaining database for Tooltip export.")
		return err
	}

	// Execute query to obtain all tooltips. Simple SELECT so using inline SQL.
	rows, err := db.Query("SELECT de.displayname_en, de.name, dg.name, de.tooltip, de.updated_at, de.updated_by FROM data_element de JOIN data_group dg ON de.data_group_id = dg.data_group_id")
	if err != nil {
		logging.Error("Error thrown executing DB query to obtain Tooltips")
		return err
	}
	defer rows.Close()

	var tooltips []ToolTipExport

	// Loop through the returned rows, build a ToolTipExport object for each one and append to the array
	for rows.Next() {
		var tooltip ToolTipExport
		var updatedAt sql.NullString
		var updatedBy sql.NullString
		err = rows.Scan(&tooltip.DisplayName, &tooltip.Name, &tooltip.DataGroup, &tooltip.ToolTip, &updatedAt, &updatedBy)
		if err != nil {
			logging.Error("Error thrown attempting to scan current row.")
			return err
		}

		// Possibility of the updated at value being null so accounting for that here
		if updatedAt.Valid {
			tooltip.UpdateTime = updatedAt.String
		} else {
			tooltip.UpdateTime = ""
		}

		// Possibility of the updated by value being null so accounting for that here
		if updatedBy.Valid {
			tooltip.UpdateBy = updatedBy.String
		} else {
			tooltip.UpdateBy = ""
		}

		tooltips = append(tooltips, tooltip)
	}

	// Set the file name
	w.Header().Set("fileName", "nextgen_tms_tooltips "+time.Now().Format("02-01-2006-15-04-05")+".csv")
	recordWriter := csv.NewWriter(w)
	var records [][]string

	// Declare the headers, these won't change unless the data we request from the DB also changes
	header := []string{"Display Name", "Element Name", "Data Group", "Tooltip", "Updated At", "Updated By"}

	// Append header to records first so it becomes the headers in the actual csv
	records = append(records, header)

	// Now loop through all the entries and append them as a line each in the csv
	for _, tooltip := range tooltips {
		var record []string
		// Add the display name
		record = append(record, tooltip.DisplayName)
		// Add the element name
		record = append(record, tooltip.Name)
		// Add the data group name
		record = append(record, tooltip.DataGroup)
		// Add the tooltip
		record = append(record, tooltip.ToolTip)
		// Add who last updated the element
		record = append(record, tooltip.UpdateTime)
		// Add when it was last updated
		record = append(record, tooltip.UpdateBy)

		// Finally add the row to the csv builder
		records = append(records, record)
	}

	// Convert the array of records into csv
	err = recordWriter.WriteAll(records)
	if err != nil {
		logging.Error("Error thrown attempting to write Tooltip export records to file.")
		return err
	}

	// Now flush the writer ready for the next export
	recordWriter.Flush()

	return nil
}
