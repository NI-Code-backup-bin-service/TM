package models

type TIDUpdateData struct {
	UpdateID        int
	PackageID       int
	ThirdPartyApkID string
	UpdateDate      string
	IsTPA           bool
	Options         []int
}

type GenerateOTPResult struct {
	result     bool
	PIN        string
	ExpiryTime string
}

type DataElementsAndGroup struct {
	DataElementID   int
	DataElementName string
	DataGroupName   string
	Options         string
}
