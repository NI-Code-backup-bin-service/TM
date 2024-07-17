package main

type DataElement struct {
	ElementId int
	GroupName string
	Name      string
	Data      string
	Encrypted bool
	Unique    bool
	Overriden bool
	Ignore    bool
}

type BulkTidUploadModel struct {
	ValidationFailed bool
	Passes           NewSitesElements
	Failure          FailedProfileValidation
	Columns          []DataColumn
}

func (b *BulkTidUploadModel) getColumns() []DataColumn {
	return b.Columns
}

func (b *BulkTidUploadModel) addColumn(column DataColumn) {
	b.Columns = append(b.getColumns(), column)
}

func (b *BulkTidUploadModel) setFailure(failure FailedProfileValidation) {
	b.Failure = failure
}

func (b *BulkTidUploadModel) setPasses(passes NewSitesElements) {
	b.Passes = passes
}

func (b *BulkTidUploadModel) setValidationResult(failed bool) {
	b.ValidationFailed = failed
}

type BulkSiteUploadModel struct {
	ValidationFailed bool
	ColumnsRemoved   bool
	Passes           NewSitesElements
	Failure          FailedProfileValidation
	Columns          []DataColumn
	UnusedColumns    []DataColumn
}

func (b *BulkSiteUploadModel) setColumnsRemoved(removed bool) {
	b.ColumnsRemoved = removed
}

func (b *BulkSiteUploadModel) getUnusedColumns() []DataColumn {
	return b.UnusedColumns
}

func (b *BulkSiteUploadModel) addColumn(column DataColumn) {
	b.Columns = append(b.getColumns(), column)
}

func (b *BulkSiteUploadModel) getColumns() []DataColumn {
	return b.Columns
}

func (b *BulkSiteUploadModel) setFailed(failed bool) {
	b.ValidationFailed = failed
}

func (b *BulkSiteUploadModel) isFailed() bool {
	return b.ValidationFailed
}

func (b *BulkSiteUploadModel) setPasses(passes NewSitesElements) {
	b.Passes = passes
}

func (b *BulkSiteUploadModel) getPasses() NewSitesElements {
	return b.Passes
}

func (b *BulkSiteUploadModel) setFailure(failure FailedProfileValidation) {
	b.Failure = failure
}

func (b *BulkSiteUploadModel) getFailure() FailedProfileValidation {
	return b.Failure
}

// FailedProfileValidation - Determines which data element from which new profile has failed validation
type FailedProfileValidation struct {
	FailureReason     string
	FailedElementId   int
	FailedElementName string
	Site              NewProfileEntry
}

func (fsv *FailedProfileValidation) setFailedElementName(name string) {
	fsv.FailedElementName = name
}

func (fsv *FailedProfileValidation) getFailureReason() string {
	return fsv.FailureReason
}

func (fsv *FailedProfileValidation) setFailureReason(reason string) {
	fsv.FailureReason = reason
}

func (fsv *FailedProfileValidation) getFailedElementId() int {
	return fsv.FailedElementId
}

func (fsv *FailedProfileValidation) setFailedElementId(elementId int) {
	fsv.FailedElementId = elementId
}

func (fsv *FailedProfileValidation) getSite() NewProfileEntry {
	return fsv.Site
}

func (fsv *FailedProfileValidation) setSite(site NewProfileEntry) {
	fsv.Site = site
}

// END FailedProfileValidation

// NewSitesElements - Stores an array of new sites to be validated and added to DB
type NewSitesElements struct {
	NewSites []NewProfileEntry
}

func (nce *NewSitesElements) getNewSites() []NewProfileEntry {
	return nce.NewSites
}

// END NewSitesElements

// NewProfileEntry - Stores the data elements for a new profile (site/tid)
type NewProfileEntry struct {
	Ref           int
	Mid           string
	SecondaryMid  string
	SiteName      string
	Serial        string
	Tid           string
	SecondaryTid  string
	SiteProfileId int // Used for bulk TID import in order to look up the site ID quickly
	DataElements  []DataElement
}

func (nse *NewProfileEntry) getTid() string {
	return nse.Tid
}

func (nse *NewProfileEntry) setTid(tid string) {
	nse.Tid = tid
}

func (nse *NewProfileEntry) getSecondaryTid() string {
	return nse.SecondaryTid
}

func (nse *NewProfileEntry) setSecondaryTid(secondaryTid string) {
	nse.SecondaryTid = secondaryTid
}

func (nse *NewProfileEntry) getSerial() string {
	return nse.Serial
}

func (nse *NewProfileEntry) setSerial(serial string) {
	nse.Serial = serial
}

func (nse *NewProfileEntry) getMid() string {
	return nse.Mid
}

func (nse *NewProfileEntry) setMid(mid string) {
	nse.Mid = mid
}

func (nse *NewProfileEntry) getSecondaryMid() string {
	return nse.SecondaryMid
}

func (nse *NewProfileEntry) setSecondaryMid(secondaryMid string) {
	nse.SecondaryMid = secondaryMid
}

func (nse *NewProfileEntry) setRef(ref int) {
	nse.Ref = ref
}

func (nse *NewProfileEntry) getRef() int {
	return nse.Ref + 1
}

func (nse *NewProfileEntry) getDataElements() []DataElement {
	return nse.DataElements
}

func (nse *NewProfileEntry) getSiteName() string {
	return nse.SiteName
}

func (nse *NewProfileEntry) setSiteName(siteName string) {
	nse.SiteName = siteName
}

func (nse *NewProfileEntry) getDataElementsAsMap() map[int]string {
	elementMap := make(map[int]string)
	for _, element := range nse.DataElements {

		// Skip over elements with out of range ids or ignore set.
		if element.Ignore || element.ElementId < 0 {
			continue
		}

		elementMap[element.ElementId] = element.Data
	}
	return elementMap
}

// END NewProfileEntry

// DataColumn - Used when converting csv row data into data elements
type DataColumn struct {
	DataGroup   string
	Name        string
	Position    int
	ElementId   int
	IsEncrypted bool
	IsUnique    bool
	IsOverriden bool
	DisplayName string
	Ignore      bool
}

func (dc *DataColumn) setIgnore(ignore bool) {
	dc.Ignore = ignore
}

func (dc *DataColumn) setDisplayName(name string) {
	dc.DisplayName = name
}

func (dc *DataColumn) setDataGroup(name string) {
	dc.DataGroup = name
}

func (dc *DataColumn) getDataGroup() string {
	return dc.DataGroup
}

func (dc *DataColumn) setName(name string) {
	dc.Name = name
}

func (dc *DataColumn) getName() string {
	return dc.Name
}

func (dc *DataColumn) setPosition(entry int) {
	dc.Position = entry
}

func (dc *DataColumn) getPosition() int {
	return dc.Position
}

func (dc *DataColumn) setElementId(entry int) {
	dc.ElementId = entry
}

func (dc *DataColumn) getElementId() int {
	return dc.ElementId
}

func (dc *DataColumn) getIgnore() bool {
	return dc.Ignore
}

func (dc *DataColumn) setOverriden(overriden bool) {
	dc.IsOverriden = overriden
}

func (dc *DataColumn) getOverriden() bool {
	return dc.IsOverriden
}

func (dc *DataColumn) setEncrypted(enc bool) {
	dc.IsEncrypted = enc
}

func (dc *DataColumn) getEncrypted() bool {
	return dc.IsEncrypted
}

func (dc *DataColumn) setUnique(unique bool) {
	dc.IsUnique = unique
}

func (dc *DataColumn) getUnique() bool {
	return dc.IsUnique
}

// END DataColumn
