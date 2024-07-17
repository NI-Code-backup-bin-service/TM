package PED

import exporter "bitbucket.org/network-international/nextgen-libs/nextgen-helpers/exportHandler"

// IMPORTANT - If you add/remove any fields from these structs ensure that all fields have a unique exportDisplayIndex

type TIDDataElement struct {
	Name         string
	DisplayName  string
	DisplayIndex string
	Value        string
}

type PED struct {
	TID                              int    `json:"tid" exportable:"true" displayName:"TID" exportDisplayIndex:"0"`
	Serial                           string `json:"serial" exportable:"true" displayName:"Serial" exportDisplayIndex:"1"`
	PIN                              string `json:"pin" exportable:"true" displayName:"PIN" exportDisplayIndex:"2"`
	ExpiryTime                       string `json:"expiryTime" exportable:"true" displayName:"ExpiryTime" exportDisplayIndex:"3"`
	ActivationTime                   string `json:"activationTime" exportable:"true" displayName:"ActivationTime" exportDisplayIndex:"4"`
	SiteId                           int    `json:"siteId" exportable:"true" displayName:"SiteId" exportDisplayIndex:"5"`
	SiteName                         string `json:"siteName" exportable:"true" displayName:"SiteName" exportDisplayIndex:"6"`
	ChainId                          int    `json:"chainId" exportable:"true" displayName:"ChainId" exportDisplayIndex:"7"`
	ChainName                        string `json:"chainName" exportable:"true" displayName:"ChainName" exportDisplayIndex:"8"`
	MerchantID                       string `json:"merchantId" exportable:"true" displayName:"MerchantID" exportDisplayIndex:"9"`
	AppVer                           int    `json:"appVer" exportable:"true" displayName:"AppVer" exportDisplayIndex:"10"`
	FirmwareVer                      string `json:"firmwareVer" exportable:"true" displayName:"FirmwareVer" exportDisplayIndex:"11"`
	LastTransaction                  string `json:"lastTransaction" exportable:"true" displayName:"LastTransaction" exportDisplayIndex:"12"`
	LastCheckedTime                  string `json:"lastCheckedTime" exportable:"true" displayName:"LastCheckedTime" exportDisplayIndex:"13"`
	ConfirmedTime                    string `json:"confirmedTime" exportable:"true" displayName:"ConfirmedTime" exportDisplayIndex:"14"`
	LastAPKDownloadTime              string `json:"lastAPKDownload" exportable:"true" displayName:"LastAPKDownloadTime" exportDisplayIndex:"15"`
	IPAddress                        string `json:"ipaddress" exportable:"true" displayName:"IP Address" exportDisplayIndex:"16"`
	IPAddresses                      string `json:"ipAddresses" exportable:"true" displayName:"IP Addresses" exportDisplayIndex:"17"`
	SIMCardSerialNumber              string `json:"simCardSerialNumber" exportable:"true" displayName:"SIMCardSerialNumber" exportDisplayIndex:"18"`
	FlagStatus                       bool   `json:"flagStatus" exportable:"true" displayName:"FlagStatus" exportDisplayIndex:"19"`
	FlaggedDate                      string `json:"flaggedDate" exportable:"true" displayName:"FlaggedDate" exportDisplayIndex:"20"`
	EODAuto                          bool   `json:"eodAuto" exportable:"true" displayName:"EODAuto" exportDisplayIndex:"21"`
	AutoTime                         string `json:"autoTime" exportable:"true" displayName:"AutoTime" exportDisplayIndex:"22"`
	PEDCoordinates                   string `json:"pedCoordinates" exportable:"true" displayName:"PEDCoordinates" exportDisplayIndex:"45"`
	PEDCoordinatesAccuracy           string `json:"pedCoordinatesAccuracy" exportable:"true" displayName:"PEDCoordinatesAccuracy" exportDisplayIndex:"46"`
	PEDCoordinatesLastUpdated        string `json:"pedCoordinatesLastUpdated" exportable:"true" displayName:"PEDCoordinatesLastUpdated" exportDisplayIndex:"47"`
	FreeInternalStorage              string `json:"freeInternalStorage" exportable:"true" displayName:"FreeInternalStorage" exportDisplayIndex:"48"`
	TotalInternalStorage             string `json:"totalInternalStorage" exportable:"true" displayName:"TotalInternalStorage" exportDisplayIndex:"49"`
	SoftuiLastDownloadedFileName     string `json:"softuiLastDownloadedFileName" exportable:"true" displayName:"SoftuiLastDownloadedFileName" exportDisplayIndex:"50"`
	SoftuiLastDownloadedFileHash     string `json:"softuiLastDownloadedFileHash" exportable:"true" displayName:"SoftuiLastDownloadedFileHash" exportDisplayIndex:"51"`
	SoftuiLastDownloadedFileList     string `json:"softuiLastDownloadedFileList" exportable:"true" displayName:"SoftuiLastDownloadedFileList" exportDisplayIndex:"52"`
	SoftuiLastDownloadedFileDateTime string `json:"softuiLastDownloadedFileDateTime" exportable:"true" displayName:"SoftuiLastDownloadedFileDateTime" exportDisplayIndex:"53"`
}

// IMPORTANT - If you add/remove any fields from these structs ensure that all fields have a unique exportDisplayIndex

type PEDDetailed struct {
	PED `exportable:"true"`

	PEDInfo exporter.ExportableItem
	// These are all references to the TMS settings core/language etc.
}

type Repository interface {
	// FindBySearchTermAndAcquirer Finds and returns a slice of PEDDetailed for a given acquirer by given search terms
	FindBySearchTermAndAcquirer(searchTerm string, acquirers string) ([]*PEDDetailed, error)

	// DeleteByTid Deletes a TID and any corresponding overrides, if no tid was found then tidDeleted
	// will return false but err wil be nil.
	DeleteByTid(tid string) (tidDeleted bool, err error)

	// DeleteOverrideByTid Deletes the overrides of a given TID, if no override was found then overrideDeleted
	// will return false but err wil be nil.
	DeleteOverrideByTid(tid string) (overrideDeleted bool, err error)

	// DeleteFraudOverrideByTid Deletes the fraud overrides of a given TID, if no override was found then fraudOverrideDeleted
	// will return false but err wil be nil.
	DeleteFraudOverrideByTid(tid string) (fraudOverrideDeleted bool, err error)

	// DeleteUserOverrideByTid Deletes the user overrides of a given TID, if no override was found then userOverrideDeleted
	// will return false but err wil be nil.
	DeleteUserOverrideByTid(tid string) (userOverrideDeleted bool, err error)
}
