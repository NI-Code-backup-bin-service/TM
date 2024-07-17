package config

import (
	cfg "bitbucket.org/network-international/nextgen-libs/nextgen-helpers/configHelper"
)

var (
	FileserverURL              string
	FlaggingFileDirectory      string
	BulkSiteUpdateDirectory    string
	BulkTidUpdateDirectory     string
	BulkTidDeleteDirectory     string
	BulkSiteDeleteDirectory    string
	BulkPaymentUploadDirectory string
)

func ParseConfig() {
	FileserverURL = cfg.GetString("FileServerURL", "localhost:3654")
	FlaggingFileDirectory = cfg.GetString("FlaggingFileDirectory", "TerminalFlagging")
	BulkSiteUpdateDirectory = cfg.GetString("BulkSiteUpdateDirectory", "BulkSiteUpdate")
	BulkTidUpdateDirectory = cfg.GetString("BulkTidUpdateDirectory", "BulkTidUpdate")
	BulkSiteDeleteDirectory = cfg.GetString("BulkSiteDeleteDirectory", "BulkSiteDelete")
	BulkTidDeleteDirectory = cfg.GetString("BulkTidDeleteDirectory", "BulkTidDelete")
	BulkPaymentUploadDirectory = cfg.GetString("BulkPaymentUploadDirectory", "BulkPaymentUpload")
}
