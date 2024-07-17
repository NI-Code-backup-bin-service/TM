// This test touches the database, so only run it when specifically asked to, via go test -tags=integration
//go:build db
// +build db

package dal

import (
	RpcHelper "bitbucket.org/network-international/nextgen-libs/nextgen-helpers/rpcHelper"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"sort"
	"testing"
)

const (
	dataElementIdAlipayMid     = 77
	dataElementIdAlipayPartner = 75
	dataElementIdLanguage      = 91
	dataElementIdMerchant      = 1
	dataElementIdSiteName      = 3
	dataElementIdRetryTime     = 12

	userName  = "bob"
	chainName = "UNIT_TEST_CHAIN"

	defaultRetryTime    = "3000"
	overriddenRetryTime = "4500"
)

// This function is actually in SiteManagement.go, but it is here as a point of comparison
// for GetChainData
func TestGetSiteData(t *testing.T) {
	setUpDb()

	chainProfileId, siteProfileId, siteId := saveNewSite(t)
	siteGroups, chainGroups, acquirerGroups, globalGroups, _, _, _, _, _, err := GetSiteData(
		int(siteId), int(siteProfileId), 10, 1, "")
	if err != nil {
		t.Error(err)
	}

	if len(siteGroups) == 0 {
		t.Error("Site groups must not be empty")
	}

	if len(chainGroups) != 0 {
		t.Error("Chain groups will be empty because it is a new chain")
	}

	checkGlobalDataSameAsDefault(t, globalGroups)
	checkAcquirerDataIsSameAsDefault(t, acquirerGroups)
	// Chain groups will be empty because it is a newly created chain
	checkChainGroupsAreEmpty(t, chainGroups)

	// Override retry time; we expect the retry entry to "move" from
	// the acquirer groups to the chain groups
	overrideRetryTime(t, chainProfileId)
	siteGroups, chainGroups, acquirerGroups, globalGroups, _, _, _, _, _, err = GetSiteData(
		int(siteId), int(siteProfileId), 10, 1, "")
	checkGlobalDataSameAsDefault(t, globalGroups)
	checkAcquirerDataIsSameAsDefaultButWithoutRetryTime(t, acquirerGroups)
	checkChainGroupsContainOnlyOverriddenRetryTime(t, chainGroups)

	// Insert some unapproved data; we expect no change  in the returned results
	// because it has not yet been approved
	insertUnapprovedData(t, chainProfileId, dataElementIdAlipayPartner, "new partner")
	siteGroups, chainGroups, acquirerGroups, globalGroups, _, _, _, _, _, err = GetSiteData(
		int(siteId), int(siteProfileId), 10, 1, "")
	checkGlobalDataSameAsDefault(t, globalGroups)
	checkAcquirerDataIsSameAsDefaultButWithoutRetryTime(t, acquirerGroups)
	checkChainGroupsContainOnlyOverriddenRetryTime(t, chainGroups)
}

func TestGetChainData_NewlyCreatedChain(t *testing.T) {
	setUpDb()

	chainProfileId := saveNewChain(t)
	chainGroups, acquirerGroups, globalGroups, err := GetChainData(int(chainProfileId))
	if err != nil {
		t.Error(err)
	}

	checkGlobalDataSameAsDefault(t, globalGroups)
	checkAcquirerDataIsSameAsDefault(t, acquirerGroups)
	// Chain groups will be empty because it is a newly created chain
	checkChainGroupsAreEmpty(t, chainGroups)

	// Override retry time; we expect the retry entry to "move" from
	// the acquirer groups to the chain groups
	overrideRetryTime(t, chainProfileId)
	chainGroups, acquirerGroups, globalGroups, err = GetChainData(int(chainProfileId))
	checkGlobalDataSameAsDefault(t, globalGroups)
	checkAcquirerDataIsSameAsDefaultButWithoutRetryTime(t, acquirerGroups)
	checkChainGroupsContainOnlyOverriddenRetryTime(t, chainGroups)

	// Insert some unapproved data; we expect no change  in the returned results
	// because it has not yet been approved
	insertUnapprovedData(t, chainProfileId, dataElementIdAlipayPartner, "new partner")
	chainGroups, acquirerGroups, globalGroups, err = GetChainData(int(chainProfileId))
	checkGlobalDataSameAsDefault(t, globalGroups)
	checkAcquirerDataIsSameAsDefaultButWithoutRetryTime(t, acquirerGroups)
	checkChainGroupsContainOnlyOverriddenRetryTime(t, chainGroups)
}

func TestGetIsOverridenForChain(t *testing.T) {
	setUpDb()
	chainProfileId := saveNewChain(t)

	overridden, err := GetIsOverridenForChain(int(chainProfileId), dataElementIdAlipayMid)
	if err != nil {
		t.Errorf("Failed to determine if chain is overridden: %v", err)
	}
	if overridden {
		t.Error("A data element with no values for any profile should not " +
			"be reported as overridden")
	}

	// Data element 'alipayMid'; this should not be reported as overridden because there are
	// no values defined its acquirer or the global profile
	saveElement(t, chainProfileId, dataElementIdAlipayMid, "foo")
	overridden, err = GetIsOverridenForChain(int(chainProfileId), dataElementIdAlipayMid)
	if err != nil {
		t.Errorf("Failed to determine if chain is overridden: %v", err)
	}
	if overridden {
		t.Error("A data element with a value for the chain but not acquirer or global " +
			"should not be reported as overridden")
	}

	// Data element 'language'; this should be reported as overridden
	// because it is defined in the global profile
	overridden, err = GetIsOverridenForChain(int(chainProfileId), dataElementIdLanguage)
	if err != nil {
		t.Errorf("Failed to determine if chain is overridden: %v", err)
	}
	if !overridden {
		t.Error("A data element with a value in the global profile " +
			"should be reported as overridden")
	}

	// Data element 'retryTime'; this should be reported as overridden
	// because it is defined in the acquirer profile
	overridden, err = GetIsOverridenForChain(int(chainProfileId), dataElementIdRetryTime)
	if err != nil {
		t.Errorf("Failed to determine if chain is overridden: %v", err)
	}
	if !overridden {
		t.Error("A data element with a value in the acquirer profile " +
			"should be reported as overridden")
	}
}

func checkGlobalDataSameAsDefault(t *testing.T, globalGroups []*DataGroup) {
	if len(globalGroups) == 0 {
		t.Fatal("Global data must not be empty")
	}

	if len(globalGroups) <= 2 {
		t.Fatal("Expecting at least 3 global data items")
	}

	globalGroups = sortGroups(globalGroups)

	expectedEod := &DataGroup{
		DataGroupID: 5,
		DataGroup:   "endOfDay",
		DataElements: []DataElement{
			{
				ElementId:            66,
				Name:                 "xReadMaxPrints",
				Type:                 "INTEGER",
				IsAllowedEmpty:       false,
				DataValue:            "10",
				MaxLength:            -1,
				ValidationExpression: "^\\d+$",
				ValidationMessage:    "Must not be blank",
				FrontEndValidate:     false,
				Unique:               false,
				Overriden:            false,
				Options:              make([]OptionData, 0),
				OptionSelectable:     false,
				DisplayName:          "",
				Image:                "",
			},
			{
				ElementId:            67,
				Name:                 "zReadMaxPrints",
				Type:                 "INTEGER",
				IsAllowedEmpty:       false,
				DataValue:            "10",
				MaxLength:            -1,
				ValidationExpression: "^\\d+$",
				ValidationMessage:    "Must not be blank",
				FrontEndValidate:     false,
				Unique:               false,
				Overriden:            false,
				Options:              make([]OptionData, 0),
				OptionSelectable:     false,
				DisplayName:          "",
				Image:                "",
			},
		},
	}
	actualEod := globalGroups[2]
	if !reflect.DeepEqual(actualEod, expectedEod) {
		t.Errorf("End of day data group; got: %v; wanted: %v", actualEod, expectedEod)
	}
}

func checkAcquirerDataIsSameAsDefault(t *testing.T, acquirerGroups []*DataGroup) {
	if len(acquirerGroups) == 0 {
		t.Fatal("Acquirer data must not be empty")
	}

	if len(acquirerGroups) != 2 {
		t.Fatal("Expecting at 2 acquirer data groups")
	}

	acquirerGroups = sortGroups(acquirerGroups)

	expectedReversalData := &DataGroup{
		DataGroupID: 4,
		DataGroup:   "reversal",
		DataElements: []DataElement{
			defaultRetriesDataElement(),
			defaultRetryTimeDataElement(),
		},
	}

	actualReversalData := acquirerGroups[0]
	if !reflect.DeepEqual(actualReversalData, expectedReversalData) {
		t.Errorf("Reversal data group; \n   Got: %v\nWanted: %v", actualReversalData, expectedReversalData)
	}

	expectedStoreData := &DataGroup{
		DataGroupID: 1,
		DataGroup:   "store",
		DataElements: []DataElement{
			defaultAcquiringInstituteIdDataElement(),
		},
	}

	actualStoreData := acquirerGroups[1]
	if !reflect.DeepEqual(actualStoreData, expectedStoreData) {
		t.Errorf("Store data group; \n   Got: %v\nWanted: %v", actualStoreData, expectedStoreData)
	}
}

func checkAcquirerDataIsSameAsDefaultButWithoutRetryTime(t *testing.T, acquirerGroups []*DataGroup) {
	if len(acquirerGroups) == 0 {
		t.Fatal("Acquirer data must not be empty")
	}

	if len(acquirerGroups) != 2 {
		t.Fatal("Expecting 2 acquirer data groups")
	}

	acquirerGroups = sortGroups(acquirerGroups)

	expectedReversalData := &DataGroup{
		DataGroupID: 4,
		DataGroup:   "reversal",
		DataElements: []DataElement{
			defaultRetriesDataElement(),
		},
	}

	actualReversalData := acquirerGroups[0]
	if !reflect.DeepEqual(actualReversalData, expectedReversalData) {
		t.Errorf("Reversal data group; \n   Got: %v\nWanted: %v", actualReversalData, expectedReversalData)
	}

	expectedStoreData := &DataGroup{
		DataGroupID: 1,
		DataGroup:   "store",
		DataElements: []DataElement{
			defaultAcquiringInstituteIdDataElement(),
		},
	}

	actualStoreData := acquirerGroups[1]
	if !reflect.DeepEqual(actualStoreData, expectedStoreData) {
		t.Errorf("Store data group; \n   Got: %v\nWanted: %v", actualStoreData, expectedStoreData)
	}
}

func checkChainGroupsAreEmpty(t *testing.T, chainGroups []*DataGroup) {
	if len(chainGroups) != 0 {
		t.Fatal("Chain groups must be empty")
	}
}

func checkChainGroupsContainOnlyOverriddenRetryTime(t *testing.T, chainGroups []*DataGroup) {
	if len(chainGroups) != 1 {
		t.Fatal("Expecting a single chain group")
	}

	chainGroups = sortGroups(chainGroups)

	expectedChainGroup := &DataGroup{
		DataGroupID: 4,
		DataGroup:   "reversal",
		DataElements: []DataElement{
			overriddenRetryTimeDataElement(),
		},
	}

	actualChainGroup := chainGroups[0]
	if !reflect.DeepEqual(actualChainGroup, expectedChainGroup) {
		t.Errorf("Store data group; \n   Got: %v\nWanted: %v", actualChainGroup, expectedChainGroup)
	}
}

func saveNewChain(t *testing.T) int64 {
	chainProfileId, err := SaveNewChain("chain", chainName, 1, userName, 2)
	if err != nil {
		t.Error(err)
	}

	return chainProfileId
}

func saveNewSite(t *testing.T) (int64, int64, int64) {
	chainProfileId := saveNewChain(t)

	siteName := "UNIT_TEST_SITE"
	siteProfileId, siteId, err := SaveNewSite(siteName, 1, userName, int(chainProfileId))
	if err != nil {
		t.Error(err)
	}

	merchantNumber := "001180220031"
	err = SaveElementData(int(siteProfileId), dataElementIdMerchant, merchantNumber, userName, 1, 0)
	if err != nil {
		t.Error(err)
	}

	err = SaveElementData(int(siteProfileId), dataElementIdSiteName, siteName, userName, 1, 0)
	if err != nil {
		t.Error(err)
	}

	return chainProfileId, siteProfileId, siteId
}

func overrideRetryTime(t *testing.T, chainProfileId int64) {
	saveElement(t, chainProfileId, dataElementIdRetryTime, overriddenRetryTime)
}

func saveElement(t *testing.T, profileId int64, dataElementId int, dataValue string) {
	err := SaveUnapprovedElement(
		int(profileId), dataElementId, dataValue, userName, 1, ApproveNewElement)
	if err != nil {
		t.Fatal(err)
	}

	approvalId, err := getLatestApprovalId()
	if err != nil {
		t.Fatal(err)
	}

	err = ApproveChange(approvalId, userName)
	if err != nil {
		t.Fatal(err)
	}
}

func setUpDb() {
	// Assuming here that DB is already upgraded to latest
	SetConnectionString("admin:Csc_12345@tcp(localhost:3306)/NextGen_TMS")
	Connect(RpcHelper.LoggingClient{}, 10)
}

func defaultRetriesDataElement() DataElement {
	return DataElement{
		ElementId:            13,
		Name:                 "retries",
		Type:                 "INTEGER",
		IsAllowedEmpty:       false,
		DataValue:            "4",
		MaxLength:            -1,
		ValidationExpression: "",
		ValidationMessage:    "",
		FrontEndValidate:     false,
		Unique:               false,
		Overriden:            false,
		Options:              make([]OptionData, 0),
		OptionSelectable:     false,
		DisplayName:          "",
		Image:                "",
	}
}

func defaultRetryTimeDataElement() DataElement {
	return DataElement{
		ElementId:            12,
		Name:                 "retryTime",
		Type:                 "LONG",
		IsAllowedEmpty:       false,
		DataValue:            defaultRetryTime,
		MaxLength:            -1,
		ValidationExpression: "",
		ValidationMessage:    "",
		FrontEndValidate:     false,
		Unique:               false,
		Overriden:            false,
		Options:              make([]OptionData, 0),
		OptionSelectable:     false,
		DisplayName:          "",
		Image:                "",
	}
}

func overriddenRetryTimeDataElement() DataElement {
	return DataElement{
		ElementId:            12,
		Name:                 "retryTime",
		Type:                 "LONG",
		IsAllowedEmpty:       false,
		DataValue:            overriddenRetryTime,
		MaxLength:            -1,
		ValidationExpression: "",
		ValidationMessage:    "",
		FrontEndValidate:     false,
		Unique:               false,
		Overriden:            true,
		Options:              make([]OptionData, 0),
		OptionSelectable:     false,
		DisplayName:          "",
		Image:                "",
	}
}

func defaultAcquiringInstituteIdDataElement() DataElement {
	return DataElement{
		ElementId:            63,
		Name:                 "acquiringInstituteId",
		Type:                 "STRING",
		IsAllowedEmpty:       false,
		DataValue:            "11111111",
		MaxLength:            -1,
		ValidationExpression: "",
		ValidationMessage:    "",
		FrontEndValidate:     false,
		Unique:               false,
		Overriden:            false,
		Options:              make([]OptionData, 0),
		OptionSelectable:     false,
		DisplayName:          "",
		Image:                "",
	}
}

func getLatestApprovalId() (int, error) {
	db, err := GetDB()
	if err != nil {
		return -1, err
	}

	sql := `
SELECT approval_id
FROM NextGen_TMS.approvals
order by approval_id desc
limit 1;
`

	rows, err := db.Query(sql)
	if err != nil {
		return -1, err
	}
	defer rows.Close()

	var approvalId int
	if rows.Next() {
		rows.Scan(&approvalId)
		return approvalId, nil
	} else {
		return 0, errors.New("Could not find latest approval ID")
	}
}

func insertUnapprovedData(t *testing.T, profileId int64, dataElementId int, value string) {
	err := insertProfileData(profileId, dataElementId, value, 0, 1)
	if err != nil {
		t.Fatal(err)
	}
}

func insertProfileData(
	profileId int64, dataElementId int, value string, approved int, version int) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	sql := `
INSERT INTO profile_data (
	profile_id, data_element_id, datavalue, version,
	updated_at, updated_by, created_at, created_by,
	approved, overriden)
VALUES (
	?, ?, ?, ?,
	'2018-01-30 12:00:00', 'bob', '2018-01-30 12:00:00', 'bob',
	?, 1)
`
	_, err = db.Exec(sql,
		profileId, dataElementId, value, version,
		approved)
	return err
}

// TODO: Rather than this copy-paste, is there a way to reference this from Sort.go?
// The fact that the file is in the 'main' package may be an issue
func sortGroups(groups []*DataGroup) []*DataGroup {
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].DataGroup < groups[j].DataGroup
	})

	for _, g := range groups {
		g.DataElements = sortElements(g.DataElements)
	}
	return groups
}

func sortElements(elements []DataElement) []DataElement {
	sort.Slice(elements, func(i, j int) bool {
		return elements[i].Name < elements[j].Name
	})

	return elements
}
