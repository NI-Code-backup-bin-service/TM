package main

const (
	// Error Strings used to replace internal errors with generic ones in http responses
	genericServerError = "a server error occurred"
	searchError        = "an error occurred performing search"

	addTidError                     = "an error occurred adding tid"
	updateTidError                  = "an error occurred updating tid"
	applyTidUpdatesError            = "an error occurred adding tid updates"
	deleteTidUpdatesError           = "an error occurred deleting tid updates"
	savingTidProfileError           = "an error occurred saving tid override"
	deletingTidProfileError         = "an error occurred deleted tid override"
	retrievingTidDetailsError       = "an error occurred retrieving tid details"
	duplicateTidOverrideError       = "an error occurred duplicating tid override"
	addDataGroupsToTidProfilesError = "an error occurred during adding data groups to tid profiles"

	uploadFileError        = "an error occurred uploading file"
	writeDataError         = "an error occured while writing file data to buffer"
	invalidFileError       = "invalid file type"
	retrieveFileListError  = "an error retrieving file list"
	deleteFileError        = "an error occurred deleting file"
	exportFailedError      = "an error occurred during export"
	softUiFileTypeError    = "no SoftUI filetype selected"
	softUiFileMissingError = "no SoftUI file selected"

	responseError = "an error occurred generating response"

	createProfileError = "an error occurred creating profile"
	saveProfileError   = "an error occurred saving profile"
	editError          = "an error occurred during Edit"

	updateDataGroupError = "an error occurred during updating data groups"

	velocityLimitsError            = "an error occurred retrieving velocity Limits"
	saveVelocityLimitsError        = "an error occurred saving velocity Limits"
	filterVelocityLimitsError      = "an error occurred filtering velocity Limits"
	tidFraudError                  = "an error occurred with tid fraud"
	deleteTIDVelocityOverrideError = "an error occurred deleting tid velocity override"

	saveSiteError       = "an error occurred saving site"
	deleteSiteError     = "an error occurred deleting site"
	deleteChainError    = "an error occurred deleting chain"
	deleteAcquirerError = "an error occurred deleting acquirer"

	saveUserGroupError        = "an error occurred saving user group"
	deleteUserGroupError      = "an error occurred deleting user group"
	addUserGroupError         = "an error occurred adding user group"
	renameUserGroupError      = "an error occurred renaming user group"
	saveGroupPermissionsError = "an error occurred saving group permissions"
	userAuditHistoryEditError = "an error occurred with user audit history"
	exportUserAuditError      = "an error occurred exporting user audit history"

	txnUploadError   = "an error occurred uploading txns"
	generatePINError = "an error occurred generating pin"
	reportingError   = "an error occurred retrieving statistics for reporting"

	userCsvUploadError = "an error occurred uploading user csv file"

	changeApprovalHistoryError       = "an error occurred retrieving change approval history"
	exportChangeApprovalHistoryError = "an error occurred exporting change approval history"
	approveChangesError              = "an error occurred approving change(s)"
	discardChangesError              = "an error occurred discarding change(s)"

	saveSiteUsersError = "an error occurred saving site users"
	saveTidUsersError  = "an error occurred saving tid users"

	retrieveElementDataError    = "an error occurred retrieving element data"
	failedToRemoveOverrideError = "an error occurred removing override"
	tidConversionFailed         = "no tidID id provided/An error occured during parsing it to int"
	siteConversionFailed        = "no siteID id provided/An error occured during parsing to int"
	profileConversionFailed     = "no profileID id provided/An error occured during parsing to int"
	userUnauthorisedError       = "user is not authorised to perform this action"

	AcquirerValidationError = "Name does not meet rules. Must not contain special characters. Must not be empty. Must not be greater than 30 characters"
	siteNotSavedError       = "Site not saved, ChainId is invalid, please clear browser cache and try again"

	InvalidUsernameError    = "the username entered is not recognised"
	PasswordChangeError     = "an error occurred changing the user password"
	PasswordConstraintError = "The new password was rejected by the server. This may be because " +
		"the new password has been used previously too recently, or that the " +
		"password has already been changed very recently."
	PasswordExpiryMessage       = "Your password has expired."
	PasswordFirstTimeLogon      = "You must change your password before logging on the first time."
	PasswordInsufficientQuality = "Your chosen password does not meet complexity requirements"
	PasswordTooShort            = "Your chosen password must be at least six characters in length"
	PasswordTooYoung            = "You are trying to change your password too often"
	PasswordInHistory           = "Please choose a password that you have not used before"
	BulksiteUpdateMissingFile   = "Please choose a file for Site update"
	SiteBulkUpdateMissingFields = "Please choose a file and enter a template MID"
	TidBulkUpdateMissingFile    = "Please choose a file for TID import"
	bulkTidUpdateMissingFile    = "Please choose a file for TID update"
	BulkTidDeleteMissingFile    = "Please choose a file for TID delete"
	BulkSiteDeleteMissingFile   = "Please choose a file for Site delete"
	MidNotValid                 = "Please enter a valid MID"
	IncorrectFileTypeCSV        = "Please upload a valid CSV file"
	FetchProfileIdError         = "An error has occurred when attempting to retrieve profile ID for the supplied MID"
	NoProfileIdforMid           = "No record found for the supplied MID, please enter a valid MID and retry"
	TidInvalidFormat            = "TID must be unique, 8 digits long and cannot be all zeros"
	SerialInvalidFormat         = "Serial number must be unique and consist of up to 10 characters"
	ErrorFetchingSiteData       = "An error has occurred when retrieving template site data"
	InsufficientUserPermissions = "The current user does not have permission to perform this action"
	NoColumnsFound              = "No data columns found, please ensure the uploaded file is in valid CSV format"
	MID_EMPTY_IN_CSV            = "No MID found for record(s) %v. Please ensure that all records have a valid MID."
	UNIQUE_FIELDS_EXPECTED      = "The fields %v are expected to be unique. Currently record(s) %v are not unique"
	DatabaseAccessError         = "An error has occurred accessing the database"
	DatabaseTxnError            = "An error has occurred during a database transaction"
)
