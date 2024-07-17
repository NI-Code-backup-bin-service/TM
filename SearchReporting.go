package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"nextgen-tms-website/TMSExportHandler"
	"nextgen-tms-website/dal"
	"nextgen-tms-website/entities"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	exporter "bitbucket.org/network-international/nextgen-libs/nextgen-helpers/exportHandler"
	"github.com/gorilla/mux"
)

const (
	SiteTab     = "#sites"
	ChainTab    = "#chains"
	AcquirerTab = "#acquirers"
	TIDSTab     = "#tids"
	TimeFormat  = "02-01-2006-15-04-05"
)

type ExportableResultItems struct {
	SiteResults     []*dal.SiteList
	ChainResults    []*dal.ChainList
	AcquirerResults []*dal.AcquirerList
}

// getAccurateResult returns the wanted site position in site list
func getAccurateResult(results []*dal.SiteList, getByName string) (int, error) {
	for c, d := range results {
		if d.SiteName == getByName {
			return c, nil
		}
	}
	return 0, errors.New("Site Name does not exist")
}

// getCsv read csv file from request and return all the file lines
func getCsv(r *http.Request, fileName string) (csvLines [][]string, err error) {
	err = r.ParseMultipartForm(5 << 20) // maxMemory 5MB
	if err != nil {
		logging.Error(err)
		return nil, err
	}

	csvFile, _, err := r.FormFile(fileName)
	if err != nil {
		logging.Error(err)
		return nil, errors.New("error reading file information")
	}
	defer csvFile.Close()

	csvLines, err = csv.NewReader(csvFile).ReadAll()
	if err != nil {
		logging.Error(err)
		return nil, errors.New(fmt.Sprintf("not valid .csv file, reading file resulted in error: %v", err))
	}

	return csvLines, nil
}

func exportSearchHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	start := time.Now()
	query := r.URL.Query()
	searchTerm := query.Get("SearchTerm")
	activeTab := query.Get("ActiveTab")
	filtered := query.Get("Filtered")
	boolValue, boolErr := strconv.ParseBool(filtered)

	var exportableItems exporter.ExportableItems
	var exportableResultItems ExportableResultItems
	var err error

	switch activeTab {
	case SiteTab:
		exportableResultItems.SiteResults, err = dal.GetSiteList(searchTerm, tmsUser)
		if len(exportableResultItems.SiteResults) == 0 || err != nil {
			handleExportError(w, err)
			return
		}
	case ChainTab:
		exportableResultItems.ChainResults, _, _, err = dal.GetChainPage(searchTerm, "0", "-1", "0", "asc", tmsUser)
		if len(exportableResultItems.ChainResults) == 0 || err != nil {
			handleExportError(w, err)
			return
		}
	case AcquirerTab:
		exportableResultItems.AcquirerResults, _, err = dal.GetAcquirerList(searchTerm, tmsUser, 0, 0)
		if len(exportableResultItems.AcquirerResults) == 0 || err != nil {
			handleExportError(w, err)
			return
		}
	case TIDSTab:
		acquirers, err := dal.GetUserAcquirerPermissions(tmsUser)
		if err == nil {
			exportHandler := TMSExportHandler.NewHandler(dal.NewPEDRepository())
			exportableItems, err = exportHandler.ExportPEDs(searchTerm, acquirers)
			if err != nil {
				handleExportError(w, err)
			}
		}
	default:
		http.Error(w, "Invalid active tab", http.StatusInternalServerError)
		return
	}

	pendingExports[tmsUser.Username] = tmsUser.Username
	if err := os.MkdirAll(ReportDir, os.ModePerm); err != nil {
		handleExportError(w, err)
		return
	}
	if activeTab == TIDSTab {
		cancelled := handleTIDSExport(w, tmsUser, exportableItems)
		delete(pendingExports, tmsUser.Username)

		if cancelled {
			logging.Information("Export Cancelled after " + time.Since(start).String())
			return
		}
	}

	reportString, cancelled := handleOtherExports(w, activeTab, exportableResultItems, tmsUser, boolErr, boolValue, searchTerm)
	delete(pendingExports, tmsUser.Username)

	if cancelled {
		logging.Information("Export Cancelled after " + time.Since(start).String())
	} else {
		logging.Information("Total Time to Generate " + reportString + " Report: " + time.Since(start).String())
	}
}

func downloadExportedReport(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	if err := r.ParseForm(); err != nil {
		logging.Error(err.Error())
		http.Error(w, exportFailedError, http.StatusInternalServerError)
		return
	}

	fileName, ok := mux.Vars(r)["fileName"]
	if !ok || fileName == "" {
		logging.Error(errors.New("no fileName parameter provided when trying to download an exported report"))
		http.Error(w, exportFailedError, http.StatusInternalServerError)
		return
	}

	// The purpose of the filepath.Base call is to remove any ..\ chars which could put the website at risk.
	filePath := filepath.Join(ReportDir, filepath.Base(fileName))
	if _, err := os.Stat(filePath); err != nil {
		logging.Error(err.Error())
		http.Error(w, exportFailedError, http.StatusInternalServerError)
		return
	}

	defer os.Remove(filePath)
	http.ServeFile(w, r, filePath)
}

func cancelExportHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	cancelExport[tmsUser.Username] = tmsUser.Username
}

func exportSites(w http.ResponseWriter, siteResults []*dal.SiteList, user *entities.TMSUser, filtered bool, searchTerm string) bool {
	fileName := fmt.Sprintf("Sites_%s.csv", time.Now().Format(TimeFormat))
	file, err := openFileForWrite(fileName)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "an error has occurred", http.StatusInternalServerError)
		return true
	}
	defer file.Close()

	w.Header().Set("fileName", fileName)
	csvWriter := csv.NewWriter(file)
	var records [][]string

	dataGroupInfo := dal.GetDataGroupInfo()

	if filtered {
		whiteListed := []string{"store.merchantNo", "store.name", "store.addressLine1", "store.addressLine2",
			"store.receiptFooterLine1", "store.receiptFooterLine2", "store.amexMid", "store.timezoneName",
			"store.disableAutomaticUpdates", "endOfDay.time", "endOfDay.auto", "endOfDay.softLimit",
			"endOfDay.hardLimit", "endOfDay.xReadMaxPrints", "endOfDay.zReadMaxPrints",
			"endOfDay.userwiseReceiptEnabled", "modules.active", "modules.dccEnabled", "modules.manualEntryEnabled",
			"modules.supervisorOnly", "modules.gratuityTier", "modules.gratuityMax", "modules.preAuthMax",
			"modules.dccProvider", "modules.mode", "modules.eppEnabled", "modules.PINRestrictedModules",
			"modules.preAuthActions", "modules.taxiEpos", "modules.preAuthWithRRNPrintPAN", "modules.dccECOflow",
			"userMgmt.available", "userMgmt.tmsUsers", "core.RequiredSoftwareVersion", "core.fraudEnabled",
			"core.allowPartialAuth", "alipay.alipayMid", "alipay.alipayMcc", "upiQr.categoryCode",
			"dualCurrency.secondaryCurrency", "dualCurrency.secondaryMid", "dualCurrency.secondaryTid",
			"dualCurrency.enabled", "dualCurrency.ctlsCvmLimit", "dualCurrency.ctlsTxnLimit",
			"dualCurrency.activeSchemes", "dualCurrency.dccEnabled", "dualCurrency.dccCtls", "dualCurrency.dccMinValue",
			"dualCurrency.dccMaxValue", "dualCurrency.dccProvider", "dualCurrency.terminalCountryCode",
			"dualCurrency.activeModules", "nol.nolEnabled", "nol.AID", "nol.nolTaxiBin", "nol.driverId",
			"nol.nolMerBin", "nol.maxCardBalLimitType1", "nol.maxCardBalLimitType2", "nol.maxCardBalLimitType4",
			"nol.locationId", "nol.businessEntityId", "nol.nolDeviceId", "instalments.EPPAAIB",
			"transactionRetrieval.siteID", "qps.enabled", "qps.limit", "vps.enabled", "vps.limit",
		}
		whiteListGroups := make(map[string][]string)
		for _, v := range whiteListed {
			str := strings.Split(v, ".")
			whiteListGroups[str[0]] = append(whiteListGroups[str[0]], str[1])
		}
		var tempDataGroupInfo []dal.DataGroupInfo
		for i, dg := range dataGroupInfo {
			if elements, ok := whiteListGroups[dg.Name]; ok {
				dataGroupInfo[i].ElementNames = elements
				tempDataGroupInfo = append(tempDataGroupInfo, dataGroupInfo[i])
			}
		}
		dataGroupInfo = tempDataGroupInfo
	}

	records = append(records, buildReportHeaders(dataGroupInfo, "Site Name"))
	var list []int
	siteIdsName := make(map[int]string, 0)
	for _, site := range siteResults {
		_, present := cancelExport[user.Username]
		if present {
			delete(cancelExport, user.Username)
			return true
		}
		list = append(list, site.SiteID)
		if _, ok := siteIdsName[site.SiteID]; !ok {
			siteIdsName[site.SiteID] = site.SiteName
		}
	}

	recordsHold := buildBatchReport(w, list, dataGroupInfo, user, searchTerm, siteIdsName)
	for _, record := range recordsHold {
		records = append(records, record)
	}

	err = csvWriter.WriteAll(records)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, exportFailedError, http.StatusInternalServerError)
		return false
	}

	return false
}

func exportChains(w http.ResponseWriter, chainResults []*dal.ChainList, user *entities.TMSUser, profileType string) bool {
	fileName := fmt.Sprintf("Chains_%s.csv", time.Now().Format(TimeFormat))
	file, err := openFileForWrite(fileName)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "an error has occurred", http.StatusInternalServerError)
		return true
	}
	defer file.Close()
	w.Header().Set("fileName", fileName)
	wr := csv.NewWriter(file)
	var records [][]string

	dataGroupInfo := dal.GetDataGroupInfo()
	records = append(records, buildReportHeaders(dataGroupInfo, "Chain Name"))

	for _, chain := range chainResults {
		_, present := cancelExport[user.Username]
		if present {
			delete(cancelExport, user.Username)
			return true
		}

		records = append(records, buildReport(w, chain.ChainProfileID, dataGroupInfo, user))
	}

	err = wr.WriteAll(records)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, exportFailedError, http.StatusInternalServerError)
		return false
	}

	return false
}

func exportAcquirers(w http.ResponseWriter, acquirerResults []*dal.AcquirerList, user *entities.TMSUser, profileType string) bool {
	fileName := fmt.Sprintf("Acquirers_%s.csv", time.Now().Format("02-01-2006-15-04-05"))
	file, err := openFileForWrite(fileName)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "an error has occurred", http.StatusInternalServerError)
		return true
	}
	defer file.Close()

	w.Header().Set("fileName", fileName)
	wr := csv.NewWriter(file)
	var records [][]string

	dataGroupInfo := dal.GetDataGroupInfo()
	records = append(records, buildReportHeaders(dataGroupInfo, "Acquirer Name"))

	for _, acquirer := range acquirerResults {
		_, present := cancelExport[user.Username]
		if present {
			delete(cancelExport, user.Username)
			return true
		}

		records = append(records, buildReport(w, acquirer.AcquirerProfileID, dataGroupInfo, user))
	}

	err = wr.WriteAll(records)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, exportFailedError, http.StatusInternalServerError)
		return false
	}

	return false
}

func buildReportHeaders(dataGroupInfo []dal.DataGroupInfo, nameField string) []string {
	var recordHeaders []string
	recordHeaders = append(recordHeaders, nameField)

	if nameField == "Chain Name" {
		recordHeaders = append(recordHeaders, "totalSites", "totaltids")
	}

	for _, dg := range dataGroupInfo {
		for _, ele := range dg.ElementNames {
			recordHeaders = append(recordHeaders, fmt.Sprintf("%s.%s", dg.Name, ele))
		}
	}
	return recordHeaders
}

func buildBatchReport(w http.ResponseWriter, sites []int, dataGroupInfo []dal.DataGroupInfo, user *entities.TMSUser, searchTerm string, siteNamesMap map[int]string) [][]string {

	t, err := dal.GetReportBatchData(sites, searchTerm)
	if err != nil {
		logging.Error(err.Error())
		return nil
	}

	for profileId := range t {
		// Sort the profile groups into id order
		sort.Slice(t[profileId], func(i, j int) bool {
			return t[profileId][i].DataGroupID < t[profileId][j].DataGroupID
		})
		// Sort the data elements in id order
		for _, g := range t[profileId] {
			sort.Slice(g.DataElements, func(i, j int) bool {
				return g.DataElements[i].ElementId < g.DataElements[j].ElementId
			})
		}
	}

	var records [][]string
	for _, id := range sites {
		var record []string
		var profileName string
		if name, ok := siteNamesMap[id]; ok {
			profileName = name
		} else {
			handleError(w, errors.New(exportFailedError), user)
			return nil
		}

		var siteConfig = make([][]string, len(dataGroupInfo))
		for i, dg := range dataGroupInfo {
			siteConfig[i] = make([]string, len(dg.ElementNames))
			for _, obj := range t[id] {
				if dg.DataGroupID == obj.DataGroupID {
					for j, ele := range dg.ElementNames {
						for _, pele := range obj.DataElements {
							if ele == pele.Name {
								siteConfig[i][j] = pele.DataValue
								continue
							}
						}
					}
					continue
				}
			}

			if i == 0 {
				siteConfig[i] = append([]string{profileName}, siteConfig[i]...)
			}
		}
		for _, obj := range siteConfig {
			record = append(record, obj...)
		}
		records = append(records, record)
	}
	return records
}

func buildReport(w http.ResponseWriter, id int, dataGroupInfo []dal.DataGroupInfo, user *entities.TMSUser) []string {

	profileName, profileType, siteID, err := dal.GetDetailsByProfileID(id)
	if err != nil {
		handleError(w, errors.New(exportFailedError), user)
		return nil
	}

	var p ProfileMaintenanceModel
	t, err := dal.GetReportData(int(siteID.Int64), profileType, id)

	if err != nil {
		logging.Error(err.Error())
		return nil
	}
	p.ProfileGroups = t

	// Sort the profile groups into id order
	sort.Slice(p.ProfileGroups, func(i, j int) bool {
		return p.ProfileGroups[i].DataGroupID < p.ProfileGroups[j].DataGroupID
	})

	// Sort the data elements in id order
	for _, g := range p.ProfileGroups {
		sort.Slice(g.DataElements, func(i, j int) bool {
			return g.DataElements[i].ElementId < g.DataElements[j].ElementId
		})
	}

	var record []string
	var siteConfig = make([][]string, len(dataGroupInfo))
	for i, dg := range dataGroupInfo {
		siteConfig[i] = make([]string, len(dg.ElementNames))
		for _, obj := range p.ProfileGroups {

			if dg.DataGroupID == obj.DataGroupID {

				for j, ele := range dg.ElementNames {
					for _, pele := range obj.DataElements {

						if ele == pele.Name {
							siteConfig[i][j] = pele.DataValue
						}
					}
				}
			}
		}

	}

	record = append(record, profileName)

	if profileType == "chain" {
		chainSiteCount, chainTidCount, err := dal.GetChainSiteIDsAndTids(id)
		if err != nil {
			logging.Error(err.Error())
			http.Error(w, "unable to fetch chain site list", http.StatusInternalServerError)
			return nil
		}
		record = append(record, strconv.Itoa(chainSiteCount), strconv.Itoa(chainTidCount))
	}

	for _, obj := range siteConfig {
		record = append(record, obj...)
	}
	return record
}

func handleTIDSExport(w http.ResponseWriter, tmsUser *entities.TMSUser, exportableItems exporter.ExportableItems) bool {
	// TODO: This should really be replaced by some context cancellation logic
	// Has the user cancelled the request?
	if _, present := cancelExport[tmsUser.Username]; present {
		delete(cancelExport, tmsUser.Username)
		return true
	}

	fileName := fmt.Sprintf("TIDS_%s.csv", time.Now().Format(TimeFormat))
	file, err := openFileForWrite(fileName)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "an error has occurred", http.StatusInternalServerError)
		return true
	}
	defer file.Close()

	// Set the file name
	w.Header().Set("fileName", fileName)
	fieldFormattingRules := make(exporter.FieldNameFormattingRules, 0)
	typeFormattingRules := make(exporter.TypeFormattingRules, 0)
	typeFormattingRules[reflect.TypeOf(time.Time{})] = func(inObj interface{}) string {
		inTime := inObj.(time.Time)
		if inTime.IsZero() {
			return ""
		}
		return inTime.String()
	}
	err = exportableItems.WriteAsCsv(file, typeFormattingRules, fieldFormattingRules)
	if err != nil {
		http.Error(w, "an error has occurred", http.StatusInternalServerError)
		return true
	}

	return false
}

func handleOtherExports(w http.ResponseWriter, activeTab string, exportableResultItems ExportableResultItems, tmsUser *entities.TMSUser, boolErr error, boolValue bool, searchTerm string) (string, bool) {
	var reportString string
	var cancelled bool

	profileType := strings.Replace(activeTab, "#", "", -1)

	switch activeTab {
	case SiteTab:
		reportString = "Sites"
		cancelled = exportSitesData(w, exportableResultItems.SiteResults, tmsUser, boolErr, boolValue, searchTerm)
	case ChainTab:
		reportString = "Chains"
		cancelled = exportChains(w, exportableResultItems.ChainResults, tmsUser, profileType)
	case AcquirerTab:
		reportString = "Acquirers"
		cancelled = exportAcquirers(w, exportableResultItems.AcquirerResults, tmsUser, profileType)
	}

	return reportString, cancelled
}

func exportSitesData(w http.ResponseWriter, siteResults []*dal.SiteList, tmsUser *entities.TMSUser, boolErr error, boolValue bool, searchTerm string) bool {
	if boolErr != nil {
		return exportSites(w, siteResults, tmsUser, false, searchTerm)
	}
	return exportSites(w, siteResults, tmsUser, boolValue, searchTerm)
}

func handleExportError(w http.ResponseWriter, err error) {
	logging.Error(err)
	http.Error(w, exportFailedError, http.StatusInternalServerError)
}

func openFileForWrite(fileName string) (*os.File, error) {
	file, err := os.OpenFile(filepath.Join(ReportDir, fileName), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0700)
	if err != nil {
		logging.Error(err.Error())
		return nil, err
	}
	return file, nil
}
