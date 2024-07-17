package main

import (
	"net/http"
	auth "nextgen-tms-website/authentication"
	"nextgen-tms-website/dal"
	"nextgen-tms-website/entities"

	rl "bitbucket.org/network-international/nextgen-libs/nextgen-helpers/rateLimitHelpers"
	web_shared "bitbucket.org/network-international/nextgen-tms/web-shared"
	sharedDAL "bitbucket.org/network-international/nextgen-tms/web-shared/dal"
	"bitbucket.org/network-international/nextgen-tms/web-shared/userManagement"
	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
)

func createTMSWebsiteRouter() *mux.Router {
	r := mux.NewRouter()

	r.Handle("/profileMaintenance", webHandlers(profileMaintenanceHandler, None))
	r.Handle("/profileMaintenanceChangeHistory", webHandlers(profileMaintenanceChangeHistoryHandler, None))
	r.Handle("/profileMaintenance/", webHandlers(profileMaintenanceHandler, None))
	r.Handle("/saveProfile", webHandlers(saveProfileHandler, SiteWrite))
	r.Handle("/saveTidProfile", webHandlers(saveTidProfileHandler, SiteWrite))
	r.Handle("/deleteTidProfile", webHandlers(deleteTidProfileHandler, SiteWrite))
	r.Handle("/getElementData", webHandlers(getElementDataHandler, None))
	r.Handle("/removeOverride", webHandlers(removeOverrideHandler, None))
	r.Handle("/updateDataGroups", webHandlers(updateDataGroupsHandler, SiteWrite))

	r.Handle("/search", webHandlers(searchHandler, None))
	r.Handle("/search/", webHandlers(searchHandler, None))
	r.Handle("/exportSearch", webHandlers(exportSearchHandler, SiteWrite))
	r.Handle("/cancelExport", webHandlers(cancelExportHandler, SiteWrite))
	r.Handle("/downloadExportedReport/{fileName}", webHandlers(downloadExportedReport, SiteWrite))

	r.Handle("/signon", webHandlers(signonHandler, NoSignIn))
	r.Handle("/logout", webHandlers(logoutHandler, None))
	r.Handle("/", webHandlers(signonHandler, NoSignIn))

	handleFilePath(r, "/assets/", "./assets")
	handleFilePath(r, "/FusionCharts/", "./FusionCharts")
	handleFilePath(r, "/static/", "./smart-energy-monitoring-dashboard/static")

	r.Handle("/addNewTID", webHandlers(addNewTIDHandler, SiteWrite))
	r.Handle("/addNewDuplicatedTidOverride", webHandlers(addNewDuplicatedTidOverride, SiteWrite))
	r.Handle("/addTID", webHandlers(addTIDHandler, SiteWrite))
	r.Handle("/deleteTID", webHandlers(deleteTIDHandler, SiteWrite))
	r.Handle("/updateSerialNumber", webHandlers(updateSerialNumber, SiteWrite))
	r.Handle("/generateOTP", webHandlers(generatePINHandler, SiteWrite))
	r.Handle("/getTidDetails", webHandlers(getTIDDetailsHandler, None))
	r.Handle("/updatesTID", webHandlers(updatesTIDHandler, SiteWrite))
	r.Handle("/updatesThirdParty", webHandlers(updatesThirdPartyHandler, SiteWrite))
	r.Handle("/updatesThirdParty/Select", webHandlers(updatesThirdPartySelectHandler, SiteWrite))
	r.Handle("/updatesThirdPartyApks", webHandlers(handleUpdatesThirdPartyApks, SiteWrite))
	r.Handle("/getThirdPartyTarget", webHandlers(getThirdPartyTarget, SiteWrite))
	r.Handle("/updatesSN", webHandlers(updatesSNHandler, SiteWrite))
	r.Handle("/addTIDUpdate", webHandlers(AddTIDUpdateHandler, SiteWrite))
	r.Handle("/ApplyTIDUpdates", webHandlers(ApplyTIDUpdatesHandler, SiteWrite))
	r.Handle("/DeleteTIDUpdate", webHandlers(DeleteTIDUpdatesHandler, SiteWrite))
	r.Handle("/GetSiteUsers", webHandlers(getUsersForSite, None))
	r.Handle("/SaveSiteUsers", webHandlers(saveSiteUsersHandler, SiteWrite))
	r.Handle("/UploadUserCSV", webHandlers(handleUserCsvUpload, SiteWrite))
	r.Handle("/ExportUploadUserCSVResult", webHandlers(handleUserCsvUploadResultExport, None))
	r.Handle("/showTidUserModal", webHandlers(showTidUserModal, None))
	r.Handle("/GetTidUsers", webHandlers(getUsersForTid, None))
	r.Handle("/SaveTidUsers", webHandlers(saveTidUsersHandler, SiteWrite))
	r.Handle("/ClearTidUsers", webHandlers(clearTidUsers, SiteWrite))
	r.Handle("/ExportUserCsv", webHandlers(handleUserCsvExport, None))

	r.Handle("/GetSiteVelocityLimits", webHandlers(getSiteVelocityLimits, Fraud))
	r.Handle("/SaveSiteVelocityLimits", webHandlers(saveSiteVelocityLimits, Fraud))
	r.Handle("/velocityLimitRow", webHandlers(handleVelocityFilterRow, Fraud))
	r.Handle("/showTidFraudModal", webHandlers(showTidFraudModal, Fraud))
	r.Handle("/TidFraudClose", webHandlers(handleTidFraudClose, None))
	r.Handle("/deleteSiteVelocityLimits", webHandlers(deleteTIDVelocityOverrides, Fraud))

	r.Handle("/addSite", webHandlers(addSiteHandler, AddCreate))
	r.Handle("/getAddSiteFields", webHandlers(getAddProfileFieldsHandler, AddCreate))

	r.Handle("/saveNewProfile", rl.MiddlewareRateLimit(webHandlers(saveNewProfileHandler, AddCreate), 5))
	r.Handle("/addNewDuplicatedChain", webHandlers(addNewDuplicatedChainHandler, ChainDuplication))
	r.Handle("/getDataGroups", webHandlers(getDataGroupsHandler, AddCreate))

	r.Handle("/userManual", webHandlers(userManualHandler, None))
	r.Handle("/submitContactUsForm", webHandlers(submitContactUsFormHandler, ContactEdit))

	r.Handle("/deleteSite", webHandlers(deleteSiteHandler, SiteWrite))
	r.Handle("/deleteChain", webHandlers(deleteChainHandler, SiteWrite))
	r.Handle("/deleteAcquirer", webHandlers(deleteAcquirerHandler, SiteWrite))

	r.Handle("/changeApproval", webHandlers(changeApprovalHandler, ChangeApprovalRead))
	r.Handle("/changeApproval/", webHandlers(changeApprovalHandler, ChangeApprovalRead))
	r.Handle("/approveChange", webHandlers(approveChangeHandler, ChangeApprovalWrite))
	r.Handle("/discardChange", webHandlers(discardChangeHandler, ChangeApprovalWrite))
	r.Handle("/approveAllChanges", webHandlers(approveAllChangesHandler, ChangeApprovalWrite))
	r.Handle("/discardAllChanges", webHandlers(discardAllChangesHandler, ChangeApprovalWrite))
	r.Handle("/filterChangeApproval", webHandlers(filterChangeApprovalHandler, ChangeApprovalRead))
	r.Handle("/filterChangeApprovalHistory", webHandlers(filterChangeApprovalHistoryHandler, ChangeApprovalRead))
	r.Handle("/exportChangeApprovalHistory", webHandlers(exportFilteredChangeApprovalHistory, ChangeApprovalWrite))

	r.Handle("/bulkUpdates", webHandlers(bulkImportHandler, BulkImport))
	r.Handle("/uploadSites", webHandlers(siteUploadHandler, BulkImport))
	r.Handle("/updateSites", webHandlers(siteUpdateHandler, BulkImport))
	r.Handle("/commitBulkSiteUpload", webHandlers(commitBulkSiteUploadHandler, BulkImport))
	r.Handle("/uploadTids", webHandlers(bulkTidImportHandler, BulkImport))
	r.Handle("/updateTids", webHandlers(bulkTidUpdateHandler, BulkImport))
	r.Handle("/commitBulkTidUpload", webHandlers(commitBulkTidUploadHandler, BulkImport))
	r.Handle("/bulkDelete", webHandlers(bulkDelete, SiteWrite))
	r.Handle("/uploadPaymentServicesGroups", webHandlers(bulkPaymentServiceImportHandler, BulkImport))
	r.Handle("/commitPaymentServicesGroups", webHandlers(commitPaymentServiceImportHandler, BulkImport))
	r.Handle("/uploadPaymentServicesTids", webHandlers(bulkPaymentServiceTidImportHandler, BulkImport))
	r.Handle("/commitPaymentServicesTids", webHandlers(commitPaymentServiceTidImportHandler, BulkImport))

	r.Handle("/fileUpload", webHandlers(fileUploadHandler, AddCreate))
	r.Handle("/txnUpload", webHandlers(uploadTxnHandler, None))
	r.Handle("/uploadFile", webHandlers(uploadFileHandler, AddCreate))
	r.Handle("/getFileList", webHandlers(getFileListHandler, AddCreate))
	r.Handle("/deleteFile", webHandlers(deleteFileHandler, AddCreate))
	r.Handle("/getFile", webHandlers(getFileHandler, AddCreate))
	r.Handle("/softUiFileUpload", webHandlers(uploadSoftUiFileHandler, AddCreate))

	r.Handle("/downloadRpiCertificate", webHandlers(generateRpiCertificate, SiteWrite))

	r.Handle("/backupDatabase", webHandlers(backupDatabaseHandler, DbBackup))
	r.Handle("/exportTooltips", webHandlers(exportTooltipsHandler, None))

	r.Handle("/reporting", webHandlers(reportingHandler, Reporting))
	r.Handle("/reporting/totalMIDs", webHandlers(TotalMIDs, Reporting))
	r.Handle("/reporting/totalTIDs", webHandlers(TotalTIDs, Reporting))
	r.Handle("/reporting/TransactingMIDs", webHandlers(TransactingMIDs, Reporting))
	r.Handle("/reporting/TransactingTIDs", webHandlers(TransactingTIDs, Reporting))
	r.Handle("/reporting/TotalTransactions", webHandlers(TotalTransactions, Reporting))
	r.Handle("/reporting/TotalTransactionValue", webHandlers(TotalTransactionValue, Reporting))
	r.Handle("/reporting/ApprovedTransactions", webHandlers(ApprovedTransactions, Reporting))
	r.Handle("/reporting/DeclinedTransactions", webHandlers(DeclinedTransactions, Reporting))
	r.Handle("/reporting/drawTransactionVolume", webHandlers(TransactionVolume, Reporting))
	r.Handle("/reporting/drawApprovedTransactionVolume", webHandlers(ApprovedTransactionVolume, Reporting))
	r.Handle("/reporting/drawDeclinedTransactionVolume", webHandlers(DeclinedTransactionVolume, Reporting))
	r.Handle("/reporting/drawTransactionValue", webHandlers(TransactionValue, Reporting))
	r.Handle("/reporting/drawApprovedTransactionValue", webHandlers(ApprovedTransactionValue, Reporting))
	r.Handle("/reporting/drawDeclinedTransactionValue", webHandlers(DeclinedTransactionValue, Reporting))
	r.Handle("/reporting/drawTop10DeclineReasons", webHandlers(Top10DeclineReasons, Reporting))
	r.Handle("/reporting/drawTop10TransactingMIDs", webHandlers(Top10TransactingMIDs, Reporting))
	r.Handle("/reporting/drawTop10TransactingTIDs", webHandlers(Top10TransactingTIDs, Reporting))
	r.Handle("/reporting/drawTop10CardTypes", webHandlers(Top10CardTypes, Reporting))

	r.Handle("/reporting/drawTop10TransactingMIDValues", webHandlers(Top10TransactingMIDValues, Reporting))
	r.Handle("/reporting/drawTop10TransactingTIDValues", webHandlers(Top10TransactingTIDValues, Reporting))
	r.Handle("/reporting/drawTop10CardTypeValues", webHandlers(Top10CardTypeValues, Reporting))
	r.Handle("/changeUserPassword", webHandlers(ChangeUserPassword, NoSignIn))

	r.Handle("/modalTemplate/cashbackBinDefinitions", webHandlers(GetCashbackEditModal, SiteWrite))
	r.Handle("/modalTemplate/moduleGratuityConfigs", webHandlers(GetGratuityEditModal, SiteWrite))
	r.Handle("/modalTemplate/dpoMomoCountryDetails", webHandlers(GetDpoMomoEditModal, SiteWrite))
	r.Handle("/modalTemplate/softUIMCCDetails", webHandlers(GetSoftUIEditModal, SiteWrite))

	r.Handle("/offlinePIN", webHandlers(offlinePIN, OfflinePIN))
	r.Handle("/generateOfflinePIN", webHandlers(generateOfflinePIN, OfflinePIN))
	r.Handle("/offlinePINImportCSV", webHandlers(offlinePINImportCSV, OfflinePIN))

	r.Handle("/getFlagStatus", webHandlers(getFlagStatusHandler, AddCreate))
	r.Handle("/downloadFile", webHandlers(downloadFile, AddCreate))

	r.Handle("/terminalFlagging", webHandlers(terminalFlaggingHandler, TerminalFlagging))
	r.Handle("/bulkChangeApproval", webHandlers(bulkChangeApprovalHandler, BulkChangeApproval))
	r.Handle("/terminalFlagging/upload", webHandlers(terminalFlaggingUploadHandler, TerminalFlagging))
	r.Handle("/bulkChangeApproval/approve", webHandlers(approveBulkChangeApprovalHandler, AddCreate))
	r.Handle("/bulkChangeApproval/discard", webHandlers(discardBulkChangeApproval, AddCreate))
	r.Handle("/bulkChangeApproval/unapproved", webHandlers(unapprovedbulkChangeApproval, AddCreate))
	r.Handle("/bulkChangeApproval/history", webHandlers(bulkChangeApprovalHistory, AddCreate))

	// Payment Services:
	r.Handle("/paymentServicesEditModal", webHandlers(getPaymentServicesEditModal, AddCreate))
	r.Handle("/paymentServices", webHandlers(paymentServicesHandler, PaymentServices))
	r.Handle("/paymentServicesManagement", webHandlers(paymentServicesManageHandler, PaymentServices))
	r.Handle("/paymentServicesSearch", webHandlers(paymentServicesSearchHandler, PaymentServices))
	r.Handle("/paymentServicesDelete", webHandlers(paymentServicesDelete, PaymentServices))
	r.Handle("/paymentServicesDeleteGroup", webHandlers(paymentServicesDeleteGroup, PaymentServices))
	r.Handle("/paymentServicesAddGroup", webHandlers(paymentServicesAddGroup, PaymentServices))
	r.Handle("/paymentServicesAddService", webHandlers(paymentServicesAddService, PaymentServices))

	web_shared.PrepareLibraryPages(r)

	// NEX-9413 - handler transformation is required due to the way OPS and TMS handle permissions differently
	handlerTransformer := func(handler web_shared.HandlerFunction) UserHandleFunction {
		return func(w http.ResponseWriter, r *http.Request, user *entities.TMSUser) {
			var tmsUser map[string]interface{}
			if err := mapstructure.Decode(user, &tmsUser); err != nil {
				logging.Error(err.Error())
				http.Error(w, "unable to get current user", http.StatusInternalServerError)
				return
			}

			handler(w, r, tmsUser)
		}
	}

	userMgmt := userManagement.Container{
		RenderHeader: func(w http.ResponseWriter, r *http.Request, tmsUser map[string]interface{}) {
			var currentUser entities.TMSUser
			if err := mapstructure.Decode(tmsUser, &currentUser); err != nil {
				logging.Error(err.Error())
				return
			}

			renderHeader(w, r, &currentUser)
		},
		DALManager: sharedDAL.Container{
			SiteMode:                    sharedDAL.TMS,
			Logging:                     logging,
			GetDB:                       dal.GetDB,
			SaveUser:                    dal.SaveUser,
			DeleteUser:                  dal.DeleteUser,
			BuildUserAuditEntry:         dal.BuildUserAuditEntry,
			BuildGroupDeleteAuditEntity: dal.BuildGroupDeleteAuditEntity,
			BuildAddGroupAuditEntry:     dal.BuildAddGroupAuditEntry,
			BuildGroupRenameAuditEntry:  dal.BuildGroupRenameAuditEntry,
			BuildGroupChangeAuditEntry:  dal.BuildGroupChangeAuditEntry,
			BuildUserAuditHistory:       dal.BuildUserAuditHistory,
			UpdateUserGroupRequired:     dal.UpdateUserGroupRequired,
		},
		Logging:           logging,
		LoggingIdentifier: "UMN",
		GetUserAcquirerPermissions: func(tmsUser map[string]interface{}) ([]string, error) {
			var currentUser entities.TMSUser
			if err := mapstructure.Decode(tmsUser, &currentUser); err != nil {
				logging.Error(err.Error())
				return nil, err
			}

			return getUserAcquirerPermissions(&currentUser)
		},
		UpdateUsers: auth.UpdateUsers,
	}

	r.Handle("/userManagement", webHandlers(handlerTransformer(userMgmt.UserManagementHandler), UserManagement))
	r.Handle("/userManagement/", webHandlers(handlerTransformer(userMgmt.UserManagementHandler), UserManagement))
	r.Handle("/userManagement/SaveUserGroup", webHandlers(handlerTransformer(userMgmt.SaveUserGroupHandler), UserManagement))
	r.Handle("/userManagement/AddGroup", webHandlers(handlerTransformer(userMgmt.AddGroupHandler), UserManagement))
	r.Handle("/userManagement/DeleteGroup", webHandlers(handlerTransformer(userMgmt.DeleteGroupHandler), UserManagement))
	r.Handle("/userManagement/RenameGroup", webHandlers(handlerTransformer(userMgmt.RenameGroupHandler), UserManagement))
	r.Handle("/userManagement/SaveGroupPermissions", webHandlers(handlerTransformer(userMgmt.SaveGroupPermissionsHandler), UserManagement))
	r.Handle("/userManagement/Select", webHandlers(handlerTransformer(userMgmt.SelectHandler), UserManagement))
	r.Handle("/userManagement/addUser", webHandlers(handlerTransformer(userMgmt.AddUserHandler), UserManagement))
	r.Handle("/userManagement/deleteUser", webHandlers(handlerTransformer(userMgmt.DeleteUserHandler), UserManagement))
	r.HandleFunc("/userChangeAuditHistory", logHandler(authHandler(handlerTransformer(userMgmt.UserAuditHistoryHandler), UserManagementAudit)))
	r.HandleFunc("/exportUserChangeAuditHistory", logHandler(authHandler(handlerTransformer(userMgmt.ExportFilteredUserManagementAudit), UserManagementAudit)))

	r.Handle("/logoUpload", webHandlers(logoUploadHandler, LogoManagement))
	r.Handle("/getDpoMomoFieldsData", webHandlers(getDpoMomoFieldsDataHandler, None))
	r.Handle("/uploadMnoLogo", webHandlers(uploadMnoLogoHandler, LogoManagement))
	return r
}
