package main

import (
	"archive/zip"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"html"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	_ "mime/multipart"
	"net"
	"net/http"
	"net/url"
	auth "nextgen-tms-website/authentication"
	"nextgen-tms-website/common"
	"nextgen-tms-website/config"
	"nextgen-tms-website/crypt"
	"nextgen-tms-website/dal"
	"nextgen-tms-website/entities"
	"nextgen-tms-website/fileServer"
	"nextgen-tms-website/logger"
	"nextgen-tms-website/models"
	"nextgen-tms-website/routers"
	"nextgen-tms-website/services"
	"nextgen-tms-website/validation"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"bitbucket.org/network-international/nextgen-libs/nextgen-helpers/TypeComparisonHelpers/SliceComparisonHelpers"
	cfg "bitbucket.org/network-international/nextgen-libs/nextgen-helpers/configHelper"
	pbPrintDiags "bitbucket.org/network-international/nextgen-libs/nextgen-helpers/diagnosticsHelper"
	rpcHelp "bitbucket.org/network-international/nextgen-libs/nextgen-helpers/rpcHelper"
	"bitbucket.org/network-international/nextgen-libs/nextgen-helpers/stringHelpers"
	TLSUtils "bitbucket.org/network-international/nextgen-libs/nextgen-helpers/tlsutils"
	pbDiag "bitbucket.org/network-international/nextgen-libs/nextgen-tg-protobuf/Diagnostics"
	txn "bitbucket.org/network-international/nextgen-libs/nextgen-tg-protobuf/Transaction"
	sharedDAL "bitbucket.org/network-international/nextgen-tms/web-shared/dal"
	"github.com/NYTimes/gziphandler"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/jszwec/csvutil"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type UserHandleFunction func(http.ResponseWriter, *http.Request, *entities.TMSUser)

// Please have a look at StoredProcedures.md for cross validating the Stored procedures used in other repositories
// If you are going to update any of those stored procedures, please do the corresponding changes in those repositories as well
// Update this to specify which version of the DB the TMS should be running on
const (
	DbVersion        = 348
	WebsiteVersion   = "1.10.05" // For Display
	DbEncryptVersion = 38
)

type UserPermission int

const (
	NoSignIn            UserPermission = -1
	None                UserPermission = 0
	SiteWrite           UserPermission = 1
	SiteDelete          UserPermission = 2
	ChangeApprovalRead  UserPermission = 3
	ChangeApprovalWrite UserPermission = 4
	AddCreate           UserPermission = 5
	BulkUpdates         UserPermission = 6
	UserManagement      UserPermission = 8
	TransactionViewer   UserPermission = 9
	DirectQuery         UserPermission = 10
	Reporting           UserPermission = 11
	ChangeHistoryView   UserPermission = 12
	EditPasswords       UserPermission = 13
	Fraud               UserPermission = 14
	PermissionGroups    UserPermission = 15
	UserManagementAudit UserPermission = 16
	TidLogs             UserPermission = 17
	BulkImport          UserPermission = 18
	ContactEdit         UserPermission = 19
	OfflinePIN          UserPermission = 20
	DbBackup            UserPermission = 21
	TerminalFlagging    UserPermission = 22
	BulkChangeApproval  UserPermission = 23
	ChainDuplication    UserPermission = 24
	PaymentServices     UserPermission = 25
	EditToken           UserPermission = 26
	LogoManagement      UserPermission = 27
	FileManagement      UserPermission = 30
	SouhoolaLogin       UserPermission = 31
)

type PageData struct {
	CurrentUser *entities.TMSUser
	PageModel   interface{}
	CSRFField   template.HTML
}

// ProfileMaintenanceModel Struct
type ProfileMaintenanceModel struct {
	ProfileTypeName  string
	ProfileType      int
	ProfileId        int
	SiteId           int
	IsSite           bool
	NewProfile       bool
	DGUpdated        bool
	DataGroups       []DataGroupModel
	ProfileGroups    []*dal.DataGroup
	ChainGroups      []*dal.DataGroup
	AcquirerGroups   []*dal.DataGroup
	GlobalGroups     []*dal.DataGroup
	FraudGroups      []*dal.DataGroup
	TIDs             []*dal.TIDData
	Packages         []*dal.PackageData
	History          []*dal.ProfileChangeHistory
	DefaultTidGroups []*dal.DataGroup
	VelocityLimits   []*dal.DataGroup
	AvailableSchemes map[int]string
	TIDPagination    dal.TidPagination
}

type AddTidModel struct {
	TID              int
	ProfileId        int
	SiteId           int
	DefaultTidGroups []*dal.DataGroup
	IsNew            bool
	IsDuplicate      bool
	DuplicatedFrom   int
}

type TidUpdatesModel struct {
	TID                    string
	Serial                 string
	SiteID                 int
	ProfileID              string
	MinimumSoftwareVersion string
	Updates                []*models.TIDUpdateData
	Packages               []*dal.PackageData
	ThirdPartyApks         []*dal.ThirdPartyApk
	ThirdPartyModeActive   bool
	PreFixNames            []string
	BindThirdPartyTarget   []*dal.BindThirdPartyTarget
}

// SearchModel Struct
type SearchModel struct {
	SiteResults     []*dal.SiteList
	ChainResults    []*dal.ChainList
	AcquirerResults []*dal.AcquirerList
	TIDResults      []*dal.TIDData
	Packages        []*dal.PackageData
	PendingExports  map[string]string
}

type AddSiteModel struct {
	Acquirers  []*dal.ProfileData
	Chains     []*dal.ProfileData
	DataGroups []DataGroupModel
}

type DataGroupModel struct {
	Group       dal.DataGroup
	Selected    bool
	PreSelected bool
}

type SignOnModel struct {
	Username        string
	Password        string
	Error           error
	DbVersion       string
	DbTargetVersion string
	WebsiteVersion  string
	CRSFToken       string
}

type GratuityDetails struct {
	TextCharacterLimit int    `json:"textCharacterLimit"`
	English            string `json:"english"`
	French             string `json:"french,omitempty"`
	Arabic             string `json:"arabic,omitempty"`
	PreDefineTip1      int    `json:"preDefineTip1,omitempty"`
	PreDefineTip2      int    `json:"preDefineTip2,omitempty"`
	PreDefineTip3      int    `json:"preDefineTip3,omitempty"`
	PreDefineTip4      int    `json:"preDefineTip4,omitempty"`
	PreDefineTip5      int    `json:"preDefineTip5,omitempty"`
	PreDefineTip6      int    `json:"preDefineTip6,omitempty"`
	PreDefineTip7      int    `json:"preDefineTip7,omitempty"`
	PreDefineTip8      int    `json:"preDefineTip8,omitempty"`
	PreDefineTip9      int    `json:"preDefineTip9,omitempty"`
	PreDefineTip10     int    `json:"preDefineTip10,omitempty"`
}

type FilterComparison string

const (
	Equals               FilterComparison = "=="
	NotEquals            FilterComparison = "!="
	LessThan             FilterComparison = "<"
	LessThanOrEqualTo    FilterComparison = "<="
	GreaterThan          FilterComparison = ">"
	GreaterThanOrEqualTo FilterComparison = ">="
	SERVERLOGGER                          = "TMSLogging"
)

type GroupComparison string
type DiagnosticLayerRPCserver struct{}

const (
	AND GroupComparison = "AND"
	OR  GroupComparison = "OR"
)

type TransactionViewerFieldModel struct {
	Name string
	Type string
	Tag  string
}

type TransactionViewerFilterModel struct {
	FilterID   string
	Field      TransactionViewerFieldModel
	Comparison FilterComparison
	Value      string
}

type TransactionViewerFilterGroupModel struct {
	GroupID    string
	SubGroups  []TransactionViewerFilterGroupModel
	Filters    []TransactionViewerFilterModel
	Comparator GroupComparison
}

type ChangeApprovalModel struct {
	SiteTab     *ChangeApprovalTabModel
	ChainTab    *ChangeApprovalTabModel
	AcquirerTab *ChangeApprovalTabModel
	TidTab      *ChangeApprovalTabModel
	OthersTab   *ChangeApprovalTabModel
	HistoryTab  *ChangeApprovalTabModel
}

type ChangeApprovalTabModel struct {
	History          []*dal.ChangeApprovalHistory
	IdentifierColumn string
	CurrentUser      *entities.TMSUser
	Count            int
	TabType          string
}

type UserAuditHistoryModel struct {
	History          []sharedDAL.UserAuditDisplay
	IdentifierColumn string
	CurrentUser      *entities.TMSUser
	Count            int
}

type UserManagementModel struct {
	Users            []*dal.UserData
	PermissionGroups []*dal.PermissionsGroupsData
	Permissions      []*dal.PermissionsData
	Acquirers        []*dal.PermissionsGroupAcquirerData
	AuditHistory     UserAuditHistoryModel
	EditableGroup    bool
}

type UploadedFile struct {
	FieldId  int
	FileName string
	File     []byte
}

type FileListEntry struct {
	Name string
}

type ChooseFileModel struct {
	ButtonText string
	Files      []FileListEntry
	FileType   string
}

type SitesUsersModel struct {
	Users                 []entities.SiteUser
	Modules               []string
	FriendlyModules       []string
	SuperPins             []string
	HasSavePermission     bool
	HasPasswordPermission bool
}

type TidUsersModel struct {
	Users                 []entities.SiteUser
	Modules               []string
	FriendlyModules       []string
	Overriden             bool
	HasSavePermission     bool
	HasPasswordPermission bool
}

type TidModalModel struct {
	Tid              string
	AvailableSchemes map[int]string
	FraudGroups      []*dal.DataGroup
	CurrentUser      *entities.TMSUser
}

type TransactionUploadModel struct {
	Results [][]string
}

type CashbackEntry struct {
	BIN       string `json:"BIN"`
	MIN       string `json:"MIN Purchase Amount"`
	MAX       string `json:"MAX Cashback Amount"`
	CheckType string `json:"CheckType"`
}

var (
	ApplicationName string
	Version         string
	Build           string
	logging         rpcHelp.LoggingClient
	listeningPort   string
	HTTPSKeyFile    string
	HTTPSCertFile   string
	ReportDir       string
	XMPPHost        string
	XMPPDomain      string
	XMPPPassword    string
	RPIKeyFile      string
	RPICertFile     string
	UserTimeout     int
	PasswordKey     string

	pendingExports        = make(map[string]string, 0)
	cancelExport          = make(map[string]string, 0)
	GRPCclients           map[string]*rpcHelp.GRPCclient
	diagnosticGRPCclients map[string]*rpcHelp.GRPCclient

	userSessionTokenMap      = make(map[string]*entities.TMSUser, 0)
	userSessionTokenMapMutex sync.Mutex

	socketConnections map[string]*websocket.Conn
	upgrader          = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return r.Header.Get("Origin") == cfg.GetString("TMSHost", "https://localhost:")+listeningPort
		},
	}

	// Parse all templates and create function allowing maps to be passed as models
	templates = template.Must(template.New("").Funcs(template.FuncMap{
		"Now": func() string { return fmt.Sprintf("%d", time.Now().Unix()) },
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values) == 0 {
				return nil, errors.New("invalid dict call")
			}

			dict := make(map[string]interface{})

			for i := 0; i < len(values); i++ {
				key, isset := values[i].(string)
				if !isset {
					if reflect.TypeOf(values[i]).Kind() == reflect.Map {
						m := values[i].(map[string]interface{})
						for i, v := range m {
							dict[i] = v
						}
					} else {
						return nil, errors.New("dict values must be maps")
					}
				} else {
					i++
					if i == len(values) {
						return nil, errors.New("specify the key for non array values")
					}
					dict[key] = values[i]
				}

			}
			return dict, nil
		},
		"Capitalise": func(s string) string { return strings.ToUpper(s[:1]) + s[1:] },
	}).ParseGlob("assets/templates/*.html"))

	userManagementFileDirectory string
)

func init() {
	/*
	   Safety net for 'too many open files' issue on legacy code.
	   Set a sane timeout duration for the http.DefaultClient, to ensure idle connections are terminated.
	   Reference: https://stackoverflow.com/questions/37454236/net-http-server-too-many-open-files-error
	*/
	http.DefaultClient.Timeout = time.Minute * 10
}

func renderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, data interface{}, user *entities.TMSUser) {
	addHeaderSecurityItems(w, r)
	p := &PageData{CurrentUser: user}
	p.PageModel = data

	if r != nil {
		p.CSRFField = csrf.TemplateField(r)
	} else {
		p.CSRFField = ""
	}

	err := templates.ExecuteTemplate(w, tmpl, p)
	if err != nil {
		logging.Error(err)
	}
}

func renderPartialTemplate(w http.ResponseWriter, r *http.Request, tmpl string, data interface{}, user *entities.TMSUser) {
	addHeaderSecurityItems(w, r)
	p := &PageData{CurrentUser: user}

	p.PageModel = data

	if r != nil {
		p.CSRFField = csrf.TemplateField(r)
	} else {
		p.CSRFField = ""
	}

	err := templates.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		logging.Error(err)
	}
}

func ajaxResponse(w http.ResponseWriter, data interface{}) {
	addAjaxSecurityItems(w)
	w.Header().Set("Content-type", "application/json")
	err := json.NewEncoder(w).Encode(&data)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, responseError, http.StatusInternalServerError)
	}

	return
}

func renderHeader(w http.ResponseWriter, r *http.Request, user *entities.TMSUser) {
	BuildUserPermissionsModel(w, r, user)
	renderTemplate(w, r, "header", nil, user)
}

func addHeaderSecurityItems(w http.ResponseWriter, r *http.Request) {
	addAjaxSecurityItems(w)
	if r != nil {
		w.Header().Set("X-CSRF-Token", csrf.Token(r))
	}
}

func addAjaxSecurityItems(w http.ResponseWriter) {
	w.Header().Set("X-FRAME-OPTIONS", "SAMEORIGIN")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	v := "max-age=" + cfg.GetString("HSTSMaxAge", "60")
	if cfg.GetBool("HSTSIncludeSubdomains", false) {
		v += "; includeSubDomains"
	}
	if cfg.GetBool("HSTSPreload", false) {
		v += "; preload"
	}
	w.Header().Set("Strict-Transport-Security", v)

	// Cache Control
	v = cfg.GetString("CacheControlDirective", "no-cache")
	v += ", max-age=" + strconv.Itoa(cfg.GetInt("UserTimeout", 30)*60)
	w.Header().Set("Cache-Control", v)
}

func searchHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	var err error

	if r.Method == "GET" {
		renderHeader(w, r, tmsUser)
		renderTemplate(w, r, "search", &SearchModel{}, tmsUser)
		return
	}

	if err := r.ParseForm(); err != nil {
		logging.Warning(err.Error())
		http.Error(w, searchError, http.StatusInternalServerError)
		return
	}

	searchTerm := r.Form.Get("search[value]")
	requestType := r.Form.Get("requestType")
	offset := r.Form.Get("start")
	amount := r.Form.Get("length")
	orderedColumn := r.Form.Get("order[0][column]")
	orderDirection := r.Form.Get("order[0][dir]")

	drawString := r.Form.Get("draw")
	draw := -1
	if drawString != "" {
		draw, err = strconv.Atoi(r.Form.Get("draw"))
		if err != nil {
			logging.Warning(err.Error())
			http.Error(w, searchError, http.StatusInternalServerError)
			return
		}
	}

	switch requestType {
	case "site":
		siteResults, total, filtered, err := dal.GetSitePage(searchTerm, offset, amount, orderedColumn, orderDirection, tmsUser)

		if err != nil {
			logging.Warning(err.Error())
			http.Error(w, searchError, http.StatusInternalServerError)
			return
		}

		var siteBytes []byte
		if siteResults == nil {
			siteBytes = []byte("[]")
		} else {
			siteBytes, err = json.Marshal(siteResults)
			if err != nil {
				logging.Warning(err.Error())
				http.Error(w, searchError, http.StatusInternalServerError)
				return
			}
		}

		bytesToWrite := append([]byte("{\"draw\":"+strconv.Itoa(draw)+","+
			"\"recordsTotal\":"+strconv.Itoa(total)+","+
			"\"recordsFiltered\":"+strconv.Itoa(filtered)+","+
			"\"data\":"), append(siteBytes, []byte("}")...)...)

		if _, err := w.Write(bytesToWrite); err != nil {
			logging.Warning(err.Error())
			http.Error(w, searchError, http.StatusInternalServerError)
			return
		}

	case "tid":
		tidResults, total, filtered, err := dal.GetTIDPage(searchTerm, offset, amount, orderedColumn, orderDirection, tmsUser)

		if err != nil {
			logging.Error(err.Error())
			http.Error(w, searchError, http.StatusInternalServerError)
			return
		}

		var tidBytes []byte
		if tidResults == nil {
			tidBytes = []byte("[]")
		} else {
			tidBytes, err = json.Marshal(tidResults)
			if err != nil {
				logging.Error(err.Error())
				http.Error(w, searchError, http.StatusInternalServerError)
				return
			}
		}

		bytesToWrite := append([]byte("{\"draw\":"+strconv.Itoa(draw)+","+
			"\"recordsTotal\":"+strconv.Itoa(total)+","+
			"\"recordsFiltered\":"+strconv.Itoa(filtered)+","+
			"\"data\":"), append(tidBytes, []byte("}")...)...)

		if _, err := w.Write(bytesToWrite); err != nil {
			logging.Error(err.Error())
			http.Error(w, searchError, http.StatusInternalServerError)
			return
		}

	case "chain":
		chainResults, total, filtered, err := dal.GetChainPage(searchTerm, offset, amount, orderedColumn, orderDirection, tmsUser)

		if err != nil {
			logging.Error(err.Error())
		}

		var chainBytes []byte
		if chainResults == nil {
			chainBytes = []byte("[]")
		} else {
			chainBytes, err = json.Marshal(chainResults)
			if err != nil {
				logging.Error(err.Error())
				http.Error(w, searchError, http.StatusInternalServerError)
				return
			}
		}

		bytesToWrite := append([]byte("{\"draw\":"+strconv.Itoa(draw)+","+
			"\"recordsTotal\":"+strconv.Itoa(total)+","+
			"\"recordsFiltered\":"+strconv.Itoa(filtered)+","+
			"\"data\":"), append(chainBytes, []byte("}")...)...)

		if _, err := w.Write(bytesToWrite); err != nil {
			logging.Error(err.Error())
			http.Error(w, searchError, http.StatusInternalServerError)
			return
		}

	case "acquirer":
		// There is no return in the below error handling because the default value for intOffset and intAmount
		// are 0 which will result in all acquirers the user has permission to being returned.
		intOffset, err := strconv.Atoi(offset)
		if err != nil {
			logging.Error(fmt.Sprintf("An error occurred parsing offset '%s' to an int", offset))
		}
		intAmount, err := strconv.Atoi(amount)
		if err != nil {
			logging.Error(fmt.Sprintf("An error occurred parsing amount '%s' to an int", amount))
		}
		acquirerResults, totalAcquirerCount, err := dal.GetAcquirerList(searchTerm, tmsUser, intOffset, intAmount)

		if err != nil {
			logging.Error(err.Error())
		}

		var acquirerBytes []byte
		if acquirerResults == nil {
			acquirerBytes = []byte("[]")
		} else {
			acquirerBytes, err = json.Marshal(acquirerResults)
			if err != nil {
				logging.Error(err.Error())
				http.Error(w, searchError, http.StatusInternalServerError)
				return
			}
		}

		if err != nil {
			logging.Error(err.Error())
			http.Error(w, searchError, http.StatusInternalServerError)
			return
		}

		bytesToWrite := append([]byte("{\"draw\":"+strconv.Itoa(draw)+","+
			"\"recordsTotal\":"+strconv.Itoa(totalAcquirerCount)+","+
			"\"recordsFiltered\":"+strconv.Itoa(totalAcquirerCount)+","+
			"\"data\":"), append(acquirerBytes, []byte("}")...)...)

		if _, err := w.Write(bytesToWrite); err != nil {
			logging.Error(err.Error())
			http.Error(w, searchError, http.StatusInternalServerError)
			return
		}
	}
}

func getElementDataHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	r.ParseForm()

	elementId, err := strconv.Atoi(r.Form.Get("ElementId"))
	if err != nil {
		handleError(w, errors.New(retrieveElementDataError), tmsUser)
		return
	}
	profileId, err := strconv.Atoi(r.Form.Get("ProfileId"))
	if err != nil {
		logging.Warning(fmt.Sprintf("An error occurred retrieving profile ID from formdata in getElementDataHandler: %s", err.Error()))
		profileId = 0
	}
	elementMetadata, err := dal.GetDataElementMetadata(elementId, profileId)

	if err != nil {
		handleError(w, errors.New(retrieveElementDataError), tmsUser)
		return
	}

	elementMetadata.Overriden = true

	deWrapper := make(map[string]interface{})

	deWrapper["de"] = elementMetadata
	deWrapper["currentUser"] = tmsUser
	deWrapper["group"] = r.Form.Get("Group")

	renderPartialTemplate(w, r, "dataElementRow", deWrapper, tmsUser)
}

func profileMaintenanceHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	query := r.URL.Query()
	id, err := strconv.Atoi(query.Get("profileId"))

	// Used to determine if we are refreshing the profile page as a result of a change to data groups
	DGUpdated, _ := strconv.ParseBool(query.Get("DGUpdated"))

	pageSize := 10         // The number of TIDs to display on a site page at once
	pageNumber := 1        // The current page of TIDs we are on
	tidPageChange := false // Always false unless we are changing page for site TIDs
	tidSearchTerm := ""

	fErr := r.ParseForm()
	if fErr == nil {

		// Check the form for a search term
		tidSearchTerm = r.Form.Get("searchTerm")

		// Check the form to see if we are changing page
		pageChange, pcErr := strconv.ParseBool(r.Form.Get("pageChange"))
		if pcErr == nil {
			tidPageChange = pageChange
		}

		// Check the form to see if an updated pageSize has been passed in
		newPageSize, psErr := strconv.Atoi(r.Form.Get("pageSize"))
		if psErr == nil {
			// If there are no errors, then a new pageSize has been sent
			pageSize = newPageSize
		}

		// Check the form to see if an updated pageNumber has been passed in
		newPageNumber, psErr := strconv.Atoi(r.Form.Get("pageNumber"))
		if psErr == nil {
			// If there are no errors, then a new pageSize has been sent
			pageNumber = newPageNumber
		}

		var pID int
		pID, fErr = strconv.Atoi(r.Form.Get("profileID"))
		if fErr == nil {
			id = pID
		}

		DGForm, DGErr := strconv.ParseBool(r.Form.Get("DGUpdated"))
		if DGErr == nil {
			DGUpdated = DGForm
		}
	}
	profileType := query.Get("type")
	if fErr != nil {

		if err != nil {
			handleError(w, errors.New("no profile id provided"), tmsUser)
			return
		}

		type idData struct {
			ID        int
			DGUpdated bool
			Type      string
		}
		renderHeader(w, r, tmsUser)
		renderTemplate(w, r, "loader", idData{ID: id, DGUpdated: DGUpdated, Type: profileType}, tmsUser)
		return
	}

	if profileType == "" {
		profileType, err = dal.GetTypeForProfile(id)
		if err != nil {
			handleError(w, errors.New("no profile type found with provided id"), tmsUser)
			return
		}
	}

	siteId, err := dal.GetSiteFromProfile(id)
	if err != nil {
		handleError(w, errors.New("no siteId found with provided id"), tmsUser)
		return
	}

	var acquirerName string
	switch profileType {
	case "site":
		acquirerName, err = dal.GetAcquirerName(id)
		if err != nil {
			handleError(w, err, tmsUser)
			return
		}
	case "chain":
		acquirerName, err = dal.GetAcquirerNameFromChainProfileId(id)
		if err != nil {
			handleError(w, err, tmsUser)
			return
		}
	case "acquirer":
		acquirerName, err = dal.GetAcquirerNameFromAcquirerProfileId(id)
		if err != nil {
			handleError(w, errors.New("no acquirer found with provided id"), tmsUser)
			return
		}
	}

	permitted, err := checkUserAcquirerPermissions(tmsUser, acquirerName)
	if err != nil {
		handleError(w, err, tmsUser)
		return
	}

	if !permitted {
		NoPermissionsRedirect(w, r)
		return
	}

	var p ProfileMaintenanceModel
	p = buildProfileMaintenanceModel(w, profileType, id, tmsUser, pageSize, pageNumber, tidSearchTerm, siteId)
	p.ProfileTypeName = profileType
	p.DGUpdated = DGUpdated

	template := "profileMaintenance"
	if tidPageChange {
		template = "profileMaintenanceTIDs"
	}
	renderTemplate(w, r, template, p, tmsUser)
}

func profileMaintenanceChangeHistoryHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	err := r.ParseForm()
	if err != nil {
		handleError(w, errors.New("failed during parsing the form"), tmsUser)
		return
	}

	profileType := r.Form.Get("ProfileType")
	profileID, err := strconv.Atoi(r.Form.Get("profileID"))
	if err != nil {
		handleError(w, errors.New("no profile id provided"), tmsUser)
		return
	}

	siteID, err := strconv.Atoi(r.Form.Get("siteID"))
	if err != nil {
		handleError(w, errors.New("no site found"), tmsUser)
		return
	}

	pageSize := 50
	pageNumber := 1

	newPageSize, psErr := strconv.Atoi(r.Form.Get("pageSize"))
	if psErr == nil {
		pageSize = newPageSize
	}

	newPageNumber, psErr := strconv.Atoi(r.Form.Get("pageNumber"))
	if psErr == nil {
		pageNumber = newPageNumber
	}

	var p ProfileMaintenanceModel
	p.History, p.TIDPagination, err = dal.GetSiteLevelProfileChangeHistory(profileID, siteID, pageSize, pageNumber, (pageNumber-1)*50, pageSize)
	if err != nil {
		handleError(w, errors.New("no site history found with provided id"), tmsUser)
		return
	}

	p.History = SortHistory(p.History)

	p.ProfileTypeName = profileType

	renderTemplate(w, r, "profileChangeHistory", p, tmsUser)
}

func getElementFields(keyPairs url.Values, ignoreEmpty bool) map[int]string {
	dataElements := make(map[int]string)
	doubleElements := make(map[int]string)

	for key := range keyPairs {
		if strings.HasPrefix(key, "data.") {
			groupKey, err := strconv.Atoi(string(key[5:]))
			if err != nil {
				if strings.Contains(key, "***") {
					groupKey, err = strconv.Atoi(string(strings.TrimSuffix(key[5:], "***")))
					if err == nil {
						doubleElements[groupKey] = keyPairs.Get(key)
					} else {
						logging.Error(err.Error())
					}
				}
				continue
			}

			if !ignoreEmpty || keyPairs.Get(key) != "" {
				dataElements[groupKey] = keyPairs.Get(key)
			}
		} else if strings.HasPrefix(key, "multi.") {
			var groupKey, _ = strconv.Atoi(string(key[6:]))
			var multiKeys = keyPairs[key]
			dataElements[groupKey] = ""
			if len(multiKeys) > 0 {
				dataElements[groupKey] = "["
				for i := range multiKeys {
					if i > 0 && i < len(multiKeys) {
						dataElements[groupKey] += ","
					}
					dataElements[groupKey] += "\"" + multiKeys[i] + "\""
				}
				dataElements[groupKey] += "]"
			} else {
				dataElements[groupKey] = "[]"
			}
		} else if strings.HasPrefix(key, "hidden.") {
			var hiddenKey, _ = strconv.Atoi(string(key[7:]))
			if keyPairs.Get("data."+string(key[7:])) == "" &&
				keyPairs.Get("multi."+string(key[7:])) == "" {
				dataElements[hiddenKey] = keyPairs.Get(key)
			}
		}
	}

	for k, v := range doubleElements {
		if dataElements[k] == "" && v == "" {
			dataElements[k] = ""
		} else {
			dataElements[k] = dataElements[k] + " | " + v
		}
	}

	return dataElements
}

func getDataGroupFields(keyPairs url.Values) []string {
	dataGroups := make([]string, 0)
	for key := range keyPairs {
		if strings.HasPrefix(key, "dg.") {
			var groupKey = key[3:]
			dataGroups = append(dataGroups, groupKey)
		}
	}
	return dataGroups
}

func validateDataElements(validationDal dal.ValidationDal, dataElements map[int]string, profileId int, profileType string, firstFailReturn bool, requestFrom string) ([]string, int) {
	var errors []string
	var siteId int
	var element models.DataElementsAndGroup
	var primaryMID, AAIBid int
	var ok bool
	var err error
	if profileId > 0 {
		siteId, err = dal.GetSiteIdFromProfileId(profileId)
		if err != nil {
			if siteId == -3 {
				errors = append(errors, err.Error())
				if firstFailReturn {
					return errors, 0
				}
			}
		}
	}
	dataElementDetails, err := dal.GetDataAllElementID()
	if err != nil {
		errors = append(errors, err.Error())
		if firstFailReturn {
			return errors, 0
		}
	}

	metaDataElements, err := dal.GetAllDataElementsMetadata(profileId)
	if err != nil {
		errors = append(errors, err.Error())
		if firstFailReturn {
			return errors, 0
		}
	}

	if element, ok = dataElementDetails["store"+"-"+"merchantNo"]; ok {
		primaryMID = element.DataElementID
	}

	if element, ok = dataElementDetails["instalments"+"-"+"EPPAAIB"]; ok {
		AAIBid = element.DataElementID
	}

	for key, newValue := range dataElements {
		element := metaDataElements[key]
		if requestFrom == "saveProfile" {
			if element.Name == "merchantNo" || element.Name == "secondaryMid" {
				currentMid, err := dal.GetDataValue(profileId, key)
				if err != nil {
					errors = append(errors, err.Error())
					if firstFailReturn {
						return errors, key
					}
				}
				if currentMid == newValue {
					continue
				}
			}
		}
		if dataElements[key] == "" {
			isRequiredAtAcquirerLevel, err := dal.CheckAcquirerLevelRequiredDataElement(key)
			if err != nil {
				errors = append(errors, err.Error())
				if firstFailReturn {
					return errors, key
				}
			}
			if isRequiredAtAcquirerLevel == 1 {
				errors = append(errors, fmt.Sprintf("%s field cannot be left empty", element.DisplayName))
				continue
			}
		}

		switch element.Name {
		case "secondaryMid":
			// Check that MIDs in the same data elements set don't match
			if primaryMID != 0 && dataElements[element.ElementId] == dataElements[primaryMID] {
				errors = append(errors, "secondaryMid: secondary Mid must not be the same as the primary MID")
				if firstFailReturn {
					return errors, key
				}
			}

		case "thirdPartyPackageList": // this is data element have dynamic value. We have a validation at front end level.
			continue
		case "merchant_endpoint":
			u, err := url.ParseRequestURI(dataElements[element.ElementId])
			if !(err == nil && u.Scheme != "" && u.Host != "") {
				errors = append(errors, fmt.Sprintf("Merchant Endpoint: is not valid %s :", dataElements[element.ElementId]))
				if firstFailReturn {
					return errors, key
				}
			}
		case "eppEnabled":
			if dataElements[element.ElementId] != "" {
				eppEnabled, err := strconv.ParseBool(dataElements[element.ElementId])
				if err != nil {
					errors = append(errors, err.Error())
					if firstFailReturn {
						return errors, key
					}
				} else if eppEnabled && dataElements[AAIBid] == dataElements[element.ElementId] {
					errors = append(errors, "eppEnabled: Cannot have both B24 EPP enabled and C+ instalments EPP enabled")
					if firstFailReturn {
						return errors, key
					}
				}
			}
		case "time":
			if dataElements[key] == "" {
				errors = append(errors, "time : Cannot be empty")
				continue
			}
			timeRange := strings.Split(dataElements[key], " | ")
			if len(timeRange) == 1 {
				_, err = time.Parse("15:04", timeRange[0])
				if err != nil {
					errors = append(errors, "time : Not a valid time")
				}
			} else {
				startTime, err := time.Parse("15:04", timeRange[0])
				if err != nil {
					errors = append(errors, "time : From is not a valid time")
				}
				endTime, err := time.Parse("15:04", timeRange[1])
				if err != nil {
					errors = append(errors, "time : To is not a valid time")
				}
				if endTime.Sub(startTime) > 1*time.Hour {
					errors = append(errors, "time : The difference between From and To should be less than 1 hour")
				}
				if endTime.Sub(startTime) < 30*time.Minute && endTime.Sub(startTime) > -23*time.Hour {
					errors = append(errors, "time : The difference between From and To should be greater than 30 minute")
				}
			}
			continue
		}

		validated, message := validation.New(validationDal).ValidateDataElement(element, dataElements[key], profileId, config.FileserverURL)
		if !validated {
			// Allow chains and acquirers to accept blank fields
			if (profileType == "site" || profileType == "tid") ||
				((profileType == "chain" || profileType == "acquirer") && dataElements[key] != "") {
				var name string
				if element.DisplayName != "" {
					name = element.DisplayName
				} else {
					name = element.Name
				}
				errors = append(errors, name+": "+message)
				if firstFailReturn {
					return errors, key
				}
			}
		}
	}

	// Sort the errors into alphabetical order due to golang's map iteration being random
	sort.Slice(errors, func(i, j int) bool {
		return errors[i] < errors[j]
	})
	return errors, -1
}

func saveDataElements(profileId int, siteId int, dataElements map[int]string, approved int, user *entities.TMSUser, useTemplateSiteOverrides bool, templateSiteId int, profileTypeName ...string) bool {
	return saveDataElementsWithType(profileTypeName, profileId, siteId, dataElements, approved, user, useTemplateSiteOverrides, templateSiteId)
}

func saveDataElementsWithType(profileType []string, profileId int, siteId int, dataElements map[int]string, approved int, user *entities.TMSUser, useTemplateSiteOverrides bool, templateSiteId int) bool {
	var profileTypeName = ""
	//NEX-9987 Unable to remove overrides from Chain
	//make this dynamic for site chain and create site/chain/acquirer
	if len(profileType) > 0 {
		profileTypeName = profileType[0]
	}
	if profileTypeName == "" && siteId > 0 {
		profileTypeName = "site"
	}

	for key := range dataElements {
		var overriden = 0
		var isOverriden bool
		var err error

		switch profileTypeName {
		case "site":
			if useTemplateSiteOverrides {
				isOverriden, err = dal.GetIsOverriden(templateSiteId, key)
			} else {
				isOverriden, err = dal.GetIsOverriden(siteId, key)
			}

		case "chain":
			isOverriden, err = dal.GetIsOverridenForChain(profileId, key)
		default:
			isOverriden, err = false, nil
		}

		if err != nil {
			logging.Error(err)
			break
		}

		if isOverriden {
			overriden = 1
		}

		if approved == 1 {
			err := dal.SaveElementData(profileId, key, dataElements[key], user.Username, approved, overriden)
			if err != nil {
				logging.Error(err.Error())
				return false
			}
		} else {
			err := dal.SaveUnapprovedElement(profileId, key, dataElements[key], user.Username, overriden, dal.ApproveNewElement)
			if err != nil {
				logging.Error(err.Error())
				return false
			}
		}
	}

	return true
}

func IsPermittedToSaveDataElement(id int, user *entities.TMSUser, passwordFieldIds []int) bool {
	// Checks that the current field is not a password field. If it is, check that
	// the user trying to save it has the correct permissions to edit passwords.
	for _, passwordId := range passwordFieldIds {
		if passwordId == id && !user.UserPermissions.EditPasswords {
			log.Print("user does not have permission to edit passwords")
			return false
		}
	}
	//TODO: Future server-side validation on posted HTML data elements to go here
	return true
}

func saveProfileHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	if err := r.ParseMultipartForm(10000); err != nil {
		logging.Error(err.Error())
	}
	var buf bytes.Buffer
	files := make([]UploadedFile, 0)
	for name, header := range r.MultipartForm.File {
		switch name {
		case "flagged-tids-file":
			continue
		default:
			file, err := header[0].Open()
			if err != nil {
				handleError(w, errors.New(saveProfileError), tmsUser)
				return
			} else {
				defer file.Close()
				if _, err := io.Copy(&buf, file); err != nil {
					handleError(w, errors.New(saveProfileError), tmsUser)
					return
				}
				id, err := strconv.Atoi(name)
				if err != nil {
					handleError(w, errors.New(saveProfileError), tmsUser)
					return
				}
				files = append(files, UploadedFile{FieldId: id, FileName: header[0].Filename, File: buf.Bytes()})
			}
		}
	}

	parseHtmlElementId := func(id string) (dataGroupName, dataElementName string, err error) {
		splitterIndex := strings.Index(id, "-")
		if splitterIndex == -1 {
			err = errors.New("splitter `-` not found in ID, could not get element group or name")
			return dataGroupName, dataElementName, err
		}
		dataGroupName = id[:splitterIndex]
		dataElementName = id[splitterIndex+1:]
		return dataGroupName, dataElementName, err
	}
	profileID, err := strconv.Atoi(r.Form.Get("profileID"))
	if err != nil {
		handleError(w, errors.New("no profile id provided/An error occured during parsing it to int"), tmsUser)
		return
	}
	siteId, err := strconv.Atoi(r.Form.Get("siteID"))
	if err != nil {
		handleError(w, errors.New("no siteID id provided/An error occured during parsing it to int"), tmsUser)
		return
	}
	var profileTypeName = r.Form.Get("profileTypeName") //NEX-9987 to get dynamic profile type value
	elements := getElementFields(r.Form, false)
	removeOverrides := r.Form.Get("removeOverrides")
	removeOverrideIds := strings.Split(removeOverrides, ",")
	profileType, err := dal.GetTypeForProfile(profileID)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, saveProfileError, http.StatusInternalServerError)
		return
	}

	validationMessages, _ := validateDataElements(dal.NewValidationDal(), elements, profileID, profileType, false, "saveProfile")
	if len(validationMessages) > 0 {
		validationMessagesJsonBytes := &bytes.Buffer{}
		enc := json.NewEncoder(validationMessagesJsonBytes)
		enc.SetEscapeHTML(false) // Special characters are escaping using json.marshal so using encode
		err := enc.Encode(validationMessages)
		if err != nil {
			_, _ = logging.Error(err)
			http.Error(w, "An unexpected error has occurred", http.StatusInternalServerError)
			return
		}
		_, _ = logging.Error(fmt.Sprintf("Validation Errors: %+v", validationMessages))
		http.Error(w, validationMessagesJsonBytes.String(), http.StatusUnprocessableEntity)
		w.Header().Set("content-type", "application/json")
		return
	}

	dataElementDetails, err := dal.GetDataAllElementID()
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, saveProfileError, http.StatusInternalServerError)
		return
	}

	if err := processFlagging(r, profileID, siteId, profileTypeName, tmsUser); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	for _, rId := range removeOverrideIds {
		switch rId {
		case "":
			continue
		default:
			groupName, elementName, err := parseHtmlElementId(rId)
			if err != nil {
				logging.Error(err.Error())
				http.Error(w, saveProfileError, http.StatusInternalServerError)
				return
			}

			//Get the name from the database
			element, ok := dataElementDetails[groupName+"-"+elementName]
			if !ok {
				logging.Error(err.Error())
				http.Error(w, saveProfileError, http.StatusInternalServerError)
				return
			}
			val, err := dal.GetNewRemovedOverrideValue(siteId, element.DataElementID)
			if err != nil {
				logging.Error(err.Error())
				http.Error(w, saveProfileError, http.StatusInternalServerError)
				return
			}
			err = dal.SaveUnapprovedElement(profileID, element.DataElementID, val, tmsUser.Username, 0, dal.ApproveRemoveOvverride)
			if err != nil {
				logging.Error(err.Error())
				http.Error(w, saveProfileError, http.StatusInternalServerError)
				return
			}
		}
	}

	// Remove elements that have just had their overrides removed from the modified elements
	for _, htmlElementId := range removeOverrideIds {
		switch htmlElementId {
		case "":
			continue
		default:
			groupName, elementName, err := parseHtmlElementId(htmlElementId)
			if err != nil {
				logging.Error(err.Error())
				http.Error(w, saveProfileError, http.StatusInternalServerError)
				return
			}
			dataElement, ok := dataElementDetails[groupName+"-"+elementName]
			if !ok {
				logging.Error(err.Error())
				http.Error(w, saveProfileError, http.StatusInternalServerError)
				return
			}
			if _, ok := elements[dataElement.DataElementID]; ok {
				delete(elements, dataElement.DataElementID)
			}
		}
	}

	// Remove elements if the user is not permitted to edit them (passwords)
	passwordFieldIds := dal.GetPasswordFieldIds()
	for key := range elements {
		if _, ok := elements[key]; ok {
			if !IsPermittedToSaveDataElement(key, tmsUser, passwordFieldIds) {
				log.Print("user is not permitted to save ", elements[key], " to data field ID, ", key)
				delete(elements, key)
			}
		}
	}
	if err := dal.AddDataGroupToTidProfile(profileID); err != nil {
		logging.Error(err)
	}

	go registerWeChatPaySubMerchant(profileID, siteId, elements, tmsUser, dataElementDetails)
	saveDataElements(profileID, siteId, elements, 0, tmsUser, false, 0, profileTypeName)
	http.Redirect(w, r, "/search", http.StatusFound)
}

func processFlagging(r *http.Request, profileID, siteID int, profileTypeName string, tmsUser *entities.TMSUser) error {

	newValue := ""

	switch r.Form.Get("flagStatus") {
	case "":
		return errors.New("You must select a flagging option")
	case "all":
		newValue = "all"
	case "file":
		// Validate that a file has been attached
		if len(r.MultipartForm.File) < 1 {
			logging.Warning("File upload initiated without file being present")
			return errors.New("Please choose a flagging file")
		}

		// Extract the file from the request and obtain the name
		logging.Debug("Attempting to extract file from http.Request")

		file, handler, err := r.FormFile("flagged-tids-file")
		if err != nil {
			logging.Error(err.Error())
			return errors.New(uploadFileError)
		}
		defer file.Close()

		fileName := handler.Filename

		logging.Debug(fmt.Sprintf("File: %v has been uploaded", fileName))

		buff := make([]byte, 512)
		if _, err = file.Seek(0, 0); err != nil {
			logging.Error(err.Error())
			return errors.New(html.EscapeString(FailedFileRead + fileName))
		}

		if _, err = file.Read(buff); err != nil {
			logging.Error(err.Error())
			return errors.New(html.EscapeString(FailedFileRead + fileName))
		}

		// Validates the file type, based on File type
		isCSVFile := strings.HasSuffix(strings.ToLower(fileName), CsvSuffix)
		if !isCSVFile {
			logging.Warning("Incorrect filetype uploaded")
			return errors.New(IncorrectFileTypeCSV)
		}

		logging.Debug(fmt.Sprintf("File %v has passed type validation", fileName))

		logging.Debug("Resetting file read offset")
		// Need to reset the offset after checking for filetype
		if _, err = file.Seek(0, 0); err != nil {
			logging.Error(err.Error())
			return errors.New(html.EscapeString(FailedFileRead + fileName))
		}

		logging.Debug("Parsing CSV data")
		// Parse the entries from the csv along with the column headers
		csvReader := csv.NewReader(file)
		records, err := csvReader.ReadAll()
		if err != nil {
			logging.Error(err.Error())
			return errors.New(html.EscapeString(FailedFileRead + fileName))
		} else if len(records) == 0 || len(records[1:]) < 1 {
			logging.Error("No records found in uploaded CSV file")
			return errors.New(NoColumnsFound)
		}

		//some special characters are getting appended in the first row; inorder to remove that bytes.Trim is used
		header := string(bytes.TrimPrefix([]byte(strings.ToLower(strings.TrimSpace(records[0][0]))), common.ByteOrderMark))

		if header != "tid" {
			logging.Error("Invalid File Header", records[0][0])
			return errors.New("Invalid File Header")
		}

		buf := &bytes.Buffer{}
		writer := csv.NewWriter(buf)
		err = writer.WriteAll(records[:])
		if err != nil {
			logging.Error(TFTAG, "Terminal Flagging File Upload Failed while writing file data to buf : ", err.Error())
			return errors.New(html.EscapeString(writeDataError + err.Error()))
		}

		fileName = time.Now().Format("20060102150405_") + fileName
		newValue = "file : " + fileName
		logging.Debug("Renamed file to : " + fileName)

		if err := sendFileToFileServer(buf.Bytes(), fileName, TerminalFlaggingType); err != nil {
			logging.Error(err.Error())
			return errors.New(uploadFileError)
		}
	case "specific":
		tids := []string{}
		for key, value := range r.Form {
			if strings.HasPrefix(key, "fs.") {
				if value[0] == "true" {
					tids = append(tids, strings.TrimLeft(key, "fs."))
				}
			}
		}

		flaggedTids, err := dal.GetFlaggedTids(strconv.Itoa(siteID))
		if err != nil {
			logging.Error(err.Error())
			return errors.New("Unable to get flagged tids")
		}

		tids = append(tids, flaggedTids...)

		if len(tids) == 0 {
			return nil
		}
		newValue = strings.Join(tids, ", ")
	default:
		logging.Error("invalid flagging option")
		return errors.New("invalid flagging option")
	}
	dataElementID, err := dal.GetDataElementByName("core", "flagStatus")
	if err != nil {
		logging.Error(err.Error())
		return errors.New("Unable to get dataElementID")
	}

	elements := map[int]string{
		dataElementID: newValue,
	}

	saveDataElements(profileID, siteID, elements, 0, tmsUser, false, 0, profileTypeName)
	return nil
}

func buildProfileMaintenanceModel(w http.ResponseWriter, profileType string, profileID int, tmsUser *entities.TMSUser, pageSize int, pageNumber int, tidSearchTerm string, siteId int) ProfileMaintenanceModel {
	switch profileType {
	case "site":
		return buildSiteProfileMaintenanceModel(w, profileID, tmsUser, pageSize, pageNumber, tidSearchTerm, siteId)
	case "chain":
		return buildChainProfileMaintenanceModel(w, profileID, tmsUser)
	default:
		return buildOtherProfileMaintenanceModel(w, profileID, tmsUser)
	}
}

func buildSiteProfileMaintenanceModel(w http.ResponseWriter, profileID int, tmsUser *entities.TMSUser, pageSize int, pageNumber int, tidSearchTerm string, siteID int) ProfileMaintenanceModel {
	var p ProfileMaintenanceModel
	profiles := make(map[string]int, 0)
	var err error
	p.ProfileGroups, p.ChainGroups, p.AcquirerGroups, p.GlobalGroups, p.TIDs, p.Packages, p.DefaultTidGroups, profiles, p.TIDPagination, err = dal.GetSiteData(siteID, profileID, pageSize, pageNumber, tidSearchTerm)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, editError, http.StatusInternalServerError)
		return p
	}
	p.AvailableSchemes, err = dal.GetAvailableSchemesForSiteId(siteID)
	if err != nil {
		// Log the error here but don't return it; an error would occur
		// if the JSON for the cardDefinitions config is invalid, but if we return an
		// error then the page can't load and so the user can no longer update the
		// cardDefinitions in order to resolve this.
		logging.Error(err.Error())
	}

	p.ProfileGroups, p.ChainGroups, p.AcquirerGroups, p.GlobalGroups = SortGroups(p.ProfileGroups), SortGroups(p.ChainGroups), SortGroups(p.AcquirerGroups), SortGroups(p.GlobalGroups)
	populateAllConfigFiles(p)

	dataGroups, err := dal.GetProfileDataForTabByProfileId(profileID, "fraud")
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Error retrieving site profiles", http.StatusInternalServerError)
		return p
	}

	p.FraudGroups = SortGroups(dataGroups)
	for _, i := range p.TIDs {
		i.TIDProfileGroups = SortGroups(i.TIDProfileGroups)
	}
	p.DefaultTidGroups = SortGroups(p.DefaultTidGroups)

	p.ProfileId = profileID
	p.SiteId = siteID
	p.IsSite = true
	var chainId = -1
	var aquirerId = -1
	for p, id := range profiles {
		if p == "chain" {
			chainId = id
		} else if p == "acquirer" {
			aquirerId = id
		}
	}

	groups, err := getMaintenanceModelGroupsWithProfileId(profileID, aquirerId, chainId, false)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Error retrieving data groups", http.StatusInternalServerError)
		return p
	}
	p.DataGroups = groups
	return p
}

func buildChainProfileMaintenanceModel(w http.ResponseWriter, profileId int, tmsUser *entities.TMSUser) ProfileMaintenanceModel {
	var p ProfileMaintenanceModel

	p.ProfileGroups, p.History = dal.GetProfileData(profileId)

	chainGroups, acquirerGroups, globalGroups, err := dal.GetChainData(profileId)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Error retrieving chain data groups", http.StatusInternalServerError)
		return p
	}

	chainGroups = SortGroups(chainGroups)

	p.ProfileGroups, p.ChainGroups, p.AcquirerGroups, p.GlobalGroups = chainGroups, chainGroups, SortGroups(acquirerGroups), SortGroups(globalGroups)

	populateAllConfigFiles(p)

	groups, err := getMaintenanceModelGroupsForChain(profileId)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Error retrieving data groups", http.StatusInternalServerError)
		return p
	}
	p.DataGroups = groups
	p.ProfileId = profileId
	p.IsSite = false
	return p
}

func buildOtherProfileMaintenanceModel(w http.ResponseWriter, profileId int, tmsUser *entities.TMSUser) ProfileMaintenanceModel {
	var p ProfileMaintenanceModel

	p.ProfileGroups, p.History = dal.GetProfileData(profileId)

	groups, err := getMaintenanceModelGroups(profileId, -1, -1, false)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Error retrieving data groups", http.StatusInternalServerError)
		return p
	}

	p.ProfileGroups = SortGroups(p.ProfileGroups)
	p.DataGroups = groups
	populateAllConfigFiles(p)
	p.ProfileId = profileId
	p.IsSite = false
	return p
}

func getMaintenanceModelGroupsForChain(profileId int) ([]DataGroupModel, error) {
	acquirerId, err := dal.GetAcquirerIdFromChainId(profileId)
	if err != nil {
		return nil, err
	}

	return getMaintenanceModelGroups(profileId, profileId, acquirerId, true)
}

func getMaintenanceModelGroups(profileID int, chainId int, acquirerId int, isChain bool) ([]DataGroupModel, error) {
	dataGroups, err := getDataGroups(strconv.Itoa(chainId), strconv.Itoa(acquirerId), isChain)
	if err != nil {
		return nil, err
	}

	profileGroups, err := dal.GetGroupsForProfile(strconv.Itoa(profileID))
	if err != nil {
		return nil, err
	}

	for _, pg := range profileGroups {
		for dg := range dataGroups {
			if dataGroups[dg].Group.DataGroupID == pg.DataGroupID {
				dataGroups[dg].Selected = true
				break
			}
		}
	}

	return dataGroups, nil
}

func getMaintenanceModelGroupsWithProfileId(profileID, acquirerId, chainId int, isChain bool) ([]DataGroupModel, error) {
	dataGroups, err := dal.GetDataGroupsWithProfileId(acquirerId, chainId, profileID, isChain)
	if err != nil {
		return nil, err
	}
	var dataGroupModels = make([]DataGroupModel, 0)
	for group := range dataGroups {
		dataGroupModels = append(dataGroupModels, DataGroupModel{
			Group:       dataGroups[group],
			PreSelected: dataGroups[group].PreSelected,
			Selected:    dataGroups[group].IsSelected},
		)
	}
	sort.Slice(dataGroupModels, func(i, j int) bool {
		return dataGroupModels[i].Group.DataGroup < dataGroupModels[j].Group.DataGroup
	})
	return dataGroupModels, nil
}

func addSiteHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	var model AddSiteModel
	var err error
	model.Acquirers, err = dal.GetProfilesByTypeName("acquirer", tmsUser)
	if err != nil {
		handleError(w, errors.New("an error occured while retriving the profiles for acquirer"), tmsUser)
		return
	}
	model.Chains, err = dal.GetProfilesByTypeName("chain", tmsUser)
	if err != nil {
		handleError(w, errors.New("an error occured while retriving the chain profiles for chain"), tmsUser)
		return
	}

	if len(model.Acquirers) > 0 && len(model.Chains) > 0 {
		model.DataGroups, _ = getDataGroups(strconv.Itoa(model.Chains[0].ID), strconv.Itoa(model.Acquirers[0].ID), false)
	}

	renderHeader(w, r, tmsUser)
	renderTemplate(w, r, "addSite", model, tmsUser)
}

// User Manual handling
func userManualHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	acquirers, err := dal.GetUserAcquirerPermissions(tmsUser)
	if err != nil {
		logging.Warning(TAG, fmt.Sprintf("Failed to fetch contact us fields for the template site. Error: %v", err.Error()))
		http.Error(w, ErrorFetchingSiteData, http.StatusInternalServerError)
		return
	}

	rawAcquirersList := strings.Split(strings.TrimSpace(acquirers), ",")
	var acquirersList []string

	for _, acquirer := range rawAcquirersList {
		if acquirer != "" {
			acquirersList = append(acquirersList, acquirer)
		}
	}

	// We use the NI acquirer as a default
	data := dal.ContactUsData{}
	defaultID, err := dal.GetAcquirerFromAcquirerName("NI")
	acquirerName := "NI"

	// Store the user's Acquirer
	if acquirersList != nil {
		acquirerName = acquirersList[0]

		if err != nil {
			logging.Warning(TAG, fmt.Sprintf("Failed to fetch contact us fields for the template site. Error: %v", err.Error()))
			http.Error(w, ErrorFetchingSiteData, http.StatusInternalServerError)
			return
		}

		acquirerID := -1

		if len(acquirersList) > 1 {
			// If the user has more than one Acquirer then return the NI details
			acquirerID = defaultID
		} else {
			// If the user has one Acquirer then return the specific details
			acquirerID, err = dal.GetAcquirerFromAcquirerName(acquirerName)
		}

		if err != nil {
			logging.Warning(TAG, fmt.Sprintf("Failed to fetch contact us fields for the template site. Error: %v", err.Error()))
			http.Error(w, ErrorFetchingSiteData, http.StatusInternalServerError)
			return
		}

		data, err = dal.GetContactUsFields(acquirerID)
		if err != nil {
			logging.Warning(TAG, fmt.Sprintf("Failed to fetch contact us fields for the template site. Error: %v", err.Error()))
			http.Error(w, ErrorFetchingSiteData, http.StatusInternalServerError)
			return
		}
	}

	// If the data returned is blank then fallback and fetch the NI data
	if !data.Valid {
		data, err = dal.GetContactUsFields(defaultID)
		if err != nil {
			logging.Warning(TAG, fmt.Sprintf("Failed to fetch contact us fields for the template site. Error: %v", err.Error()))
			http.Error(w, ErrorFetchingSiteData, http.StatusInternalServerError)
			return
		}
	}

	data.UserAcquirer = acquirerName

	renderTemplate(w, r, "userManual.html", data, tmsUser)
}

func submitContactUsFormHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {

	acquirers, err := dal.GetUserAcquirerPermissions(tmsUser)
	if err != nil {
		logging.Warning(TAG, fmt.Sprintf("Failed to fetch contact us fields for the template site. Error: %v", err.Error()))
		http.Error(w, ErrorFetchingSiteData, http.StatusInternalServerError)
		return
	}

	rawAcquirersList := strings.Split(strings.TrimSpace(acquirers), ",")
	var acquirersList []string

	for _, acquirer := range rawAcquirersList {
		if acquirer != "" {
			acquirersList = append(acquirersList, acquirer)
		}
	}

	acquirerID := -1
	var error error

	if len(acquirersList) > 1 {
		acquirerID, error = dal.GetAcquirerFromAcquirerName("NI")
	} else {
		acquirerID, error = dal.GetAcquirerFromAcquirerName(acquirersList[0])
	}

	if error != nil {
		logging.Warning(TAG, fmt.Sprintf("Failed to fetch contact us fields for the template site. Error: %v", err.Error()))
		http.Error(w, ErrorFetchingSiteData, http.StatusInternalServerError)
		return
	}
	r.ParseForm()
	inputEmail := r.Form.Get("form[inputEmail]")
	inputPrimaryPhone := r.Form.Get("form[inputPrimaryPhone]")
	inputSecondaryPhone := r.Form.Get("form[inputSecondaryPhone]")
	inputAlOne := r.Form.Get("form[inputALOne]")
	inputAlTwo := r.Form.Get("form[inputALTwo]")
	inputAlThree := r.Form.Get("form[inputALThree]")
	inputFurtherInfo := r.Form.Get("form[inputFurtherInformation]")

	var formData dal.ContactUsData
	formData.AcquirerName = acquirersList[0]
	formData.AcquirerEmail = inputEmail
	formData.AcquirerPrimaryPhone = inputPrimaryPhone
	formData.AcquirerSecondaryPhone = inputSecondaryPhone
	formData.AcquirerAddressLineOne = inputAlOne
	formData.AcquirerAddressLineTwo = inputAlTwo
	formData.AcquirerAddressLineThree = inputAlThree
	formData.FurtherInformation = inputFurtherInfo

	err = dal.SetContactUsFields(acquirerID, formData)
	if err != nil {
		logging.Error(TAG, err.Error())
		http.Error(w, uploadFileError, http.StatusInternalServerError)
		return
	}
}

func getAdditionalGroups(keyPairs url.Values) []string {
	dataGroups := make([]string, 0)
	for key := range keyPairs {
		if strings.HasPrefix(key, "dg") {
			dataGroups = append(dataGroups, key[3:])
		}
	}

	return dataGroups
}

func downloadFile(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	r.ParseForm()
	fileName := r.Form.Get("FileName")
	fileType := r.Form.Get("FileType")
	var configDir string
	if fileType == TerminalFlaggingType {
		configDir = config.FlaggingFileDirectory
	} else if fileType == BulkSiteUpdateType {
		configDir = config.BulkSiteUpdateDirectory
	} else if fileType == BulkTidUpdateType {
		configDir = config.BulkTidUpdateDirectory
	} else if fileType == BulkTidDeleteType {
		configDir = config.BulkTidDeleteDirectory
	} else if fileType == BulkSiteDeleteType {
		configDir = config.BulkSiteDeleteDirectory
	} else if fileType == BulkPaymentServiceUploadType || fileType == BulkPaymentTidUploadType {
		configDir = config.BulkPaymentUploadDirectory
	}
	fileAsBase64Encoded, err := fileServer.NewFsReader(config.FileserverURL).GetFile(fileName, configDir)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, html.EscapeString(FailedRetrieveFile+fileName), http.StatusInternalServerError)
		return
	}

	addAjaxSecurityItems(w)
	// Encode as base64.
	fileAsBytes, err := common.ConvertBase64FileToBytes(string(fileAsBase64Encoded))
	if err != nil {
		logging.Error(err)
		http.Error(w, EncodingError, http.StatusInternalServerError)
		return
	}

	io.WriteString(w, string(fileAsBytes))
}

func getFlagStatusHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	r.ParseForm()
	var siteID = r.PostForm.Get("siteID")

	model, err := dal.GetFlagStatus(siteID)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	renderTemplate(w, r, "flagStatus", model, tmsUser)
}

func getDataGroupsHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	var model AddSiteModel
	var err error
	r.ParseForm()
	var acquirerId = r.PostForm.Get("AcquirerId")
	var chainId = r.PostForm.Get("ChainId")

	if acquirerId == "" {
		acquirerId = "-1"
	}

	if chainId != "-1" {
		acquirerId = dal.GetAcquirerIdForChain(chainId)
	}

	if chainId == "" {
		chainId = "-1"
	}

	model.DataGroups, err = getDataGroups(chainId, acquirerId, false)
	if err != nil {
		handleError(w, errors.New("an error occured while retriving the data groups in getDataGroupsHandler"), tmsUser)
		return
	}

	renderTemplate(w, r, "dataGroups", model, tmsUser)
}

func getDataGroups(chainId string, acquirerId string, isChain bool) ([]DataGroupModel, error) {
	dataGroups, err := dal.GetDataGroups()
	if err != nil {
		return nil, err
	}

	groupMap := make(map[int]DataGroupModel, 0)
	for group := range dataGroups {
		groupModel := DataGroupModel{Group: dataGroups[group], PreSelected: false}
		groupMap[dataGroups[group].DataGroupID] = groupModel
	}

	preSelected := make([]dal.DataGroup, 0)

	acquirerGroups, err := dal.GetGroupsForProfile(acquirerId)
	if err != nil {
		return nil, err
	}

	preSelected = append(preSelected, acquirerGroups...)
	if !isChain {
		chainGroups, _ := dal.GetGroupsForProfile(chainId)
		preSelected = append(preSelected, chainGroups...)
	}

	for preSelection := range preSelected {
		id := preSelected[preSelection].DataGroupID
		var dg = groupMap[id]
		dg.PreSelected = true
		groupMap[id] = dg
	}

	var dataGroupModels = make([]DataGroupModel, 0)
	for _, val := range groupMap {
		dataGroupModels = append(dataGroupModels, val)
	}
	sort.Slice(dataGroupModels, func(i, j int) bool {
		return dataGroupModels[i].Group.DataGroup < dataGroupModels[j].Group.DataGroup
	})
	return dataGroupModels, nil
}

func updateDataGroupsHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	r.ParseForm()
	profileId := r.Form.Get("profileID")
	groups := getAdditionalGroups(r.Form)
	pid, err := strconv.Atoi(profileId)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Invalid profile Id", http.StatusInternalServerError)
		return
	}

	profileType, err := dal.GetTypeForProfile(pid)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Invalid profile type", http.StatusInternalServerError)
		return
	}

	if profileType == "site" {
		dataGroupId, err := dal.GetDataGroupByName("store")
		if err != nil {
			logging.Error(err)
			http.Error(w, "An unexpected error has occurred", http.StatusInternalServerError)
			return
		}

		if r.FormValue("dg."+strconv.Itoa(dataGroupId)) != "true" {
			logging.Error(errors.New("Store data group must be enabled"))
			http.Error(w, "Store data group must be enabled", http.StatusInternalServerError)
			return
		}
	}

	err = disableDualCurrencyIfGroupDisabled(profileType, pid, groups, tmsUser)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Get all current data groups
	dataGroups, err := dal.NewDataGroupRepository().FindForSiteByProfileId(pid)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, updateDataGroupError, http.StatusInternalServerError)
		return
	}
	// Compare against new groups
	for _, i := range groups {
		for index, x := range dataGroups {
			d, _ := strconv.Atoi(i)
			if x.ID == d {
				dataGroups = append(dataGroups[:index], dataGroups[index+1:]...)
			}
		}

	}
	// Remove all data within any group not in new groups
	if len(dataGroups) > 0 {
		if err = dal.ClearDataElementsForDisabled(profileId, dataGroups); err != nil {
			logging.Error(err.Error())
			http.Error(w, updateDataGroupError, http.StatusInternalServerError)
			return
		}
	}

	if err := dal.ClearDataGroupsForProfile(profileId); err != nil {
		logging.Error(err.Error())
		http.Error(w, updateDataGroupError, http.StatusInternalServerError)
		return
	}
	services.AddDataGroupsToProfile(pid, groups, tmsUser)
}

func disableDualCurrencyIfGroupDisabled(profileType string, pid int, enabledGroups []string, tmsUser *entities.TMSUser) error {
	if profileType == "site" {
		siteGroups, err := dal.NewDataGroupRepository().FindForSiteByProfileId(pid)
		if err != nil {
			logging.Error(err.Error())
			return errors.New("error finding enabled groups for site")
		}

		// If the dual currency group was selected (not inherited or preselected) for the site but is no longer selected
		// then we need to set dualCurrency/enabled to false.
		for _, currentGroup := range siteGroups {
			if currentGroup.Name == "dualCurrency" && currentGroup.Selected && !currentGroup.Preselected {
				if !SliceComparisonHelpers.SlicesOfStringContains(enabledGroups, strconv.Itoa(currentGroup.ID)) {
					logging.Information(fmt.Sprintf("Setting dualCurrency/enabled to false for profileId '%v'", pid))

					dcEnabledDataGroupId, err := dal.GetDataElementByName("dualCurrency", "enabled")
					if err != nil {
						logging.Error(err.Error())
						return errors.New("error finding dualCurrency/enabled group id")
					}

					err = dal.NewProfileDataRepository().SetDataValueByElementIdAndProfileIdWithoutApproval(dcEnabledDataGroupId, pid, "false", *tmsUser)
					if err != nil {
						logging.Error(err.Error())
						return errors.New("error setting dualCurrency/enabled to false")
					}
				}
			}
		}
	}
	return nil
}

func getAddProfileFieldsHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	r.ParseForm()

	profileType := r.Form.Get("type")
	acquirerId, err := strconv.Atoi(r.Form.Get("acquirerId"))
	if err != nil {
		handleError(w, errors.New("an error occured while retriving/converting the acquirerId"), tmsUser)
		return
	}
	chainId, err := strconv.Atoi(r.Form.Get("chainId"))
	if err != nil {
		handleError(w, errors.New("an error occured while retriving/converting the chainId"), tmsUser)
		return
	}
	mid := r.Form.Get("mid")

	// Site creation on MID validation
	if profileType == "site" {
		acquirerId, err = strconv.Atoi(dal.GetAcquirerIdForChain(strconv.Itoa(chainId)))
		if err != nil {
			handleError(w, errors.New("an error occured during conversion at GetAcquirerIdForChain in getAddProfileFieldsHandler"+err.Error()), tmsUser)
			return
		}

		dataElements := make(map[int]string)
		dataElements[1] = mid
		validationMessages, _ := validateDataElements(dal.NewValidationDal(), dataElements, -1, profileType, false, "")
		if len(validationMessages) > 0 {
			validationMessagesJsonBytes, err := json.Marshal(validationMessages)
			if err != nil {
				logging.Error(err)
				http.Error(w, "An unexpected error has occurred", http.StatusInternalServerError)
				return
			}
			logging.Error(fmt.Sprintf("Validation Errors: %+v", validationMessages))
			http.Error(w, string(validationMessagesJsonBytes), http.StatusUnprocessableEntity)
			w.Header().Set("content-type", "application/json")
			return
		}
	}

	if profileType == "chain" {
		chainId = -1
		if profileType == "acquirer" {
			acquirerId = -1
		}
	}

	dataGroups := getAdditionalGroups(r.Form)

	var p ProfileMaintenanceModel
	if profileType == "site" || profileType == "chain" {
		p.ProfileGroups, p.ChainGroups, p.AcquirerGroups, p.GlobalGroups = dal.GetProfileFields(acquirerId, chainId, dataGroups)
	} else if profileType == "acquirer" {
		parentProfileId := 1
		switch profileType {
		case "site":
			parentProfileId = chainId
		case "chain":
			parentProfileId = acquirerId
		}

		for group := range dataGroups {
			dataGroup, err := dal.GetDataElementsForGroup(dataGroups[group], parentProfileId)
			if err != nil {
				logging.Error(err)
			}

			p.ProfileGroups = append(p.ProfileGroups, dataGroup)
		}
	}

	p.ProfileGroups = SortGroups(p.ProfileGroups)

	// Set the merchant ID to that entered in MID field
	if profileType == "site" {
		for _, val := range p.ProfileGroups {
			if val.DataGroup == "store" {
				for i, ele := range val.DataElements {
					if ele.ElementId == 1 {
						val.DataElements[i].DataValue = mid
					}
				}
			}
		}
	}

	p.ChainGroups = SortGroups(p.ChainGroups)
	p.AcquirerGroups = SortGroups(p.AcquirerGroups)
	p.GlobalGroups = SortGroups(p.GlobalGroups)
	p.IsSite = profileType == "site" || profileType == "chain"
	p.NewProfile = true
	p.ProfileTypeName = profileType
	populateAllConfigFiles(p)

	renderTemplate(w, r, "profileMaintenance", p, tmsUser)
}

func saveNewProfileHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	r.ParseForm()
	profileType := r.Form.Get("type")
	switch profileType {
	case "site":
		saveNewSite(w, r, tmsUser)
		break
	case "chain":
		chainName := r.Form.Get("name")
		if chainExists, err := dal.CheckChainNameExists(chainName); err != nil || !chainExists {
			logging.Error("chain name already exists")
			http.Error(w, "chain name already exists", http.StatusInternalServerError)
			return
		}
		fallthrough
	case "acquirer":
		saveNewProfile(w, r, profileType, tmsUser)
		break
	default:
		http.Error(w, html.EscapeString("Unknown type: "+profileType), http.StatusInternalServerError)
	}
}

func addNewDuplicatedChainHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	r.ParseForm()
	chainName := r.Form.Get("newChainName")
	acquirerName := r.Form.Get("acquirerName")
	chainProfileId, err := strconv.Atoi(r.Form.Get("chainProfileId"))
	if err != nil {
		logging.Error("unable to convert chain profile Id int to string")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	valid := validation.New(dal.NewValidationDal()).ValidateString(chainName)
	if !valid {
		logging.Error(AcquirerValidationError)
		http.Error(w, AcquirerValidationError, http.StatusInternalServerError)
		return
	}

	if chainExists, err := dal.CheckChainNameExists(chainName); err != nil || !chainExists {
		logging.Error("chain name already exists")
		http.Error(w, "chain name already exists", http.StatusInternalServerError)
		return
	}

	acquirerId, err := dal.GetAcquirerIdFromChainProfileId(chainProfileId)
	if err != nil {
		logging.Error("unable to find acquirer id for given chain id")
		handleError(w, errors.New("unable to find acquirer id for given chain id"), tmsUser)
		return
	}

	profileID, err := dal.SaveNewChain("chain", chainName, 1, tmsUser.Username, acquirerId)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, createProfileError, http.StatusInternalServerError)
		return
	}

	err = dal.AddDataGroupAndDataElement(int(profileID), chainProfileId, tmsUser.Username)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "an error occurred during updating data groups and data element", http.StatusInternalServerError)
		return
	}

	approvalId, err := dal.DuplicateChainChangeApproval(tmsUser, chainProfileId, chainName, acquirerName)
	if err != nil {
		logging.Error("unable to create change approval")
		http.Error(w, fmt.Sprintf("unable to create change approval %s", err.Error()), http.StatusInternalServerError)
		return
	}

	ajaxResponse(w, approvalId)
}

func saveNewProfile(w http.ResponseWriter, r *http.Request, profileType string, user *entities.TMSUser) {
	r.ParseForm()
	name := r.Form.Get("name")
	acquirerId, _ := strconv.Atoi(r.Form.Get("acquirer"))
	//NEX-9987 Field is coming so we can use it
	var profileTypeName = r.Form.Get("profileTypeName")
	elements := getElementFields(r.Form, false)
	groups := getDataGroupFields(r.Form)

	valid := validation.New(dal.NewValidationDal()).ValidateString(name)
	if !valid {
		logging.Error(AcquirerValidationError)
		http.Error(w, AcquirerValidationError, http.StatusInternalServerError)
		return
	}
	var validationMessages []string
	if profileType == "chain" {
		acqProfileId := acquirerId * -1
		validationMessages, _ = validateDataElements(dal.NewValidationDal(), elements, acqProfileId, profileType, false, "")
	} else {
		validationMessages, _ = validateDataElements(dal.NewValidationDal(), elements, -1, profileType, false, "")
	}

	if len(validationMessages) > 0 {
		validationMessagesJsonBytes, err := json.Marshal(validationMessages)
		if err != nil {
			logging.Error(err)
			http.Error(w, "An unexpected error has occurred", http.StatusInternalServerError)
			return
		}
		logging.Error(fmt.Sprintf("Validation Errors: %+v", validationMessages))
		http.Error(w, string(validationMessagesJsonBytes), http.StatusUnprocessableEntity)
		w.Header().Set("content-type", "application/json")
		return
	}

	var profileID int64
	var err error
	if profileType == "chain" {
		profileID, err = dal.SaveNewChain(profileType, name, 1, user.Username, acquirerId)
	} else if profileType == "acquirer" {
		if acquirerExists, err := dal.CheckAcquirerNameExists(name); err != nil || acquirerExists {
			logging.Error("acquirer name already exists")
			http.Error(w, "acquirer name already exists", http.StatusInternalServerError)
			return
		}
		profileID, err = dal.SaveNewProfile(profileType, name, 1, user.Username)
	}

	var result bool
	if err == nil {
		result = saveDataElements(int(profileID), -1, elements, 1, user, false, 0, profileTypeName)
	} else {
		logging.Error(err.Error())
		http.Error(w, createProfileError, http.StatusInternalServerError)
	}

	if result {
		services.AddDataGroupsToProfile(int(profileID), groups, user)
	}

	http.Redirect(w, r, "/search", http.StatusFound)
}

func removeOverrideHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	r.ParseForm()

	siteId, err := strconv.Atoi(r.Form.Get("SiteId"))
	if err != nil {
		handleError(w, errors.New(siteConversionFailed), tmsUser)
		return
	}
	elementId, err := strconv.Atoi(r.Form.Get("ElementId"))
	if err != nil {
		handleError(w, errors.New("an error occured while retrving the ElementID or conversion to int"), tmsUser)
		return
	}
	if err := dal.RemoveOverride(siteId, elementId); err != nil {
		handleError(w, errors.New(failedToRemoveOverrideError), tmsUser)
	}
	ajaxResponse(w, true)
}

func saveNewSite(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	r.ParseForm()
	chainId, _ := strconv.Atoi(r.Form.Get("chain"))
	//NEX-9987 Field is coming so we can use it
	var profileTypeName = r.Form.Get("profileTypeName")
	// Check ids aren't 0 or invalid ids have been sent from browser
	if chainId == 0 {
		logging.Error(siteNotSavedError)
		http.Error(w, siteNotSavedError, http.StatusInternalServerError)
		return
	}

	elements := getElementFields(r.Form, false)
	groups := getDataGroupFields(r.Form)

	chainIdForValidation := -1 * chainId
	validationMessages, _ := validateDataElements(dal.NewValidationDal(), elements, chainIdForValidation, "site", false, "")
	if len(validationMessages) > 0 {
		validationMessagesJsonBytes, err := json.Marshal(validationMessages)
		if err != nil {
			logging.Error(err)
			http.Error(w, "An unexpected error has occurred", http.StatusInternalServerError)
			return
		}
		logging.Error(fmt.Sprintf("Validation Errors: %+v", validationMessages))
		http.Error(w, string(validationMessagesJsonBytes), http.StatusUnprocessableEntity)
		w.Header().Set("content-type", "application/json")
		return
	}
	dataGroupId, err := dal.GetDataGroupByName("store")
	if err != nil {
		logging.Error(err)
		http.Error(w, "An unexpected error has occurred", http.StatusInternalServerError)
		return
	}

	if r.FormValue("dg."+strconv.Itoa(dataGroupId)) != "true" {
		validationMessages := []string{"Store data group must be enabled"}
		validationMessagesJsonBytes, err := json.Marshal(validationMessages)
		if err != nil {
			logging.Error(err)
			http.Error(w, "An unexpected error has occurred", http.StatusInternalServerError)
			return
		}
		logging.Error(fmt.Sprintf("Validation Errors: %+v", validationMessages))
		http.Error(w, string(validationMessagesJsonBytes), http.StatusUnprocessableEntity)
		w.Header().Set("content-type", "application/json")
		return
	} else {
		dataElementId, err := dal.GetDataElementByNameAndGroupID("name", dataGroupId)
		if err != nil {
			logging.Error(err)
			http.Error(w, "An unexpected error has occurred", http.StatusInternalServerError)
			return
		}

		if r.FormValue("data."+strconv.Itoa(dataElementId)) == "" {
			validationMessages := []string{"Store data group elements must be added"}
			validationMessagesJsonBytes, err := json.Marshal(validationMessages)
			if err != nil {
				logging.Error(err)
				http.Error(w, "An unexpected error has occurred", http.StatusInternalServerError)
				return
			}
			logging.Error(fmt.Sprintf("Validation Errors: %+v", validationMessages))
			http.Error(w, string(validationMessagesJsonBytes), http.StatusUnprocessableEntity)
			w.Header().Set("content-type", "application/json")
			return
		}
	}

	name := elements[3]
	profileID, siteID, err := dal.SaveNewSite(name, 1, tmsUser.Username, chainId)

	var result bool
	if err == nil {
		result = saveDataElements(int(profileID), int(siteID), elements, 1, tmsUser, false, 0, profileTypeName)
		dal.RecordSiteToHistory(int(profileID), "Site Created", tmsUser.Username, dal.ApproveCreate, 1)
	} else {
		logging.Error(err.Error())
		http.Error(w, saveSiteError, http.StatusInternalServerError)
	}

	if result {
		services.AddDataGroupsToProfile(int(profileID), groups, tmsUser)
	}

	http.Redirect(w, r, "/search", http.StatusFound)
}

func deleteSiteHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	r.ParseForm()
	siteString := r.Form.Get("siteId")

	siteProfileId, err := strconv.Atoi(siteString)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Invalid site Id", http.StatusInternalServerError)
	}

	err = dal.RecordSiteToHistory(siteProfileId, "Site Deleted", tmsUser.Username, dal.ApproveDelete, 0)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, deleteSiteError, http.StatusInternalServerError)
		return
	}
}

func deleteChainHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	r.ParseForm()
	chainString := r.Form.Get("chainId")
	chainId, err := strconv.Atoi(chainString)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Invalid chain Id", http.StatusInternalServerError)
	}
	err = dal.DeleteChain(chainId)

	if err != nil {
		logging.Error(err.Error())
		http.Error(w, deleteChainError, http.StatusInternalServerError)
	}

}

func deleteAcquirerHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	r.ParseForm()
	acquirerString := r.Form.Get("acquirerId")
	acquirerId, err := strconv.Atoi(acquirerString)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Invalid acquirer Id", http.StatusInternalServerError)
	}
	err = dal.DeleteAcquirer(acquirerId)

	if err != nil {
		logging.Error(err.Error())
		http.Error(w, deleteAcquirerError, http.StatusInternalServerError)
	}
}

func checkUserAcquirePermsBySite(user *entities.TMSUser, siteId int) (bool, error) {
	if acqName, err := dal.GetAcquirerNameForSite(siteId); err == nil {
		return checkAcquirerPermissions(user, acqName)
	} else {
		return false, err
	}
}

func checkUserAcquirePermsByTid(user *entities.TMSUser, tid string) (bool, error) {
	if acqName, err := dal.GetTidAcquirer(tid); err == nil {
		return checkAcquirerPermissions(user, acqName)
	} else {
		return false, err
	}
}

func checkUserAcquirerPermissions(user *entities.TMSUser, acquirerName string) (bool, error) {
	acquirers, err := dal.GetUserAcquirers(user)
	if err != nil {
		return false, err
	}

	return SliceComparisonHelpers.SlicesOfStringContains(acquirers, acquirerName), nil
}

func checkAcquirerPermissions(user *entities.TMSUser, acquirerName string) (bool, error) {

	permissions, err := dal.GetUserAcquirerPermissions(user)
	permisionsList := strings.Split(permissions, ",")
	if err != nil {
		return false, err
	}

	for _, permission := range permisionsList {
		if permission == acquirerName && permission != "" {
			return true, nil
		}
	}
	return false, nil
}

func uploadTxnHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	err := r.ParseMultipartForm(10000)
	if err != nil {
		logging.Error(err.Error())
	}

	txns := make(map[string]*bytes.Buffer, 0)

	for _, header := range r.MultipartForm.File {
		for _, formFile := range header {
			if filepath.Ext(formFile.Filename) != "" {
				logging.Error("Invalid File Type: " + formFile.Filename)
				http.Error(w, "Invalid File Type: "+formFile.Filename, http.StatusBadRequest)
				return
			}
			if match, err := regexp.MatchString("^\\d+_[^\\..]+$", formFile.Filename); err != nil {
				logging.Error(err.Error())
				http.Error(w, txnUploadError, http.StatusBadRequest)
				return
			} else if !match {
				logging.Error("Invalid File Name: " + formFile.Filename)
				http.Error(w, "Invalid File Name: "+formFile.Filename, http.StatusBadRequest)
				return
			}
		}
	}

	for _, header := range r.MultipartForm.File {
		for _, formFile := range header {
			var buf bytes.Buffer
			file, err := formFile.Open()
			if err != nil {
				logging.Error(err.Error())
				http.Error(w, txnUploadError, http.StatusInternalServerError)
				return
			} else {
				if _, err := io.Copy(&buf, file); err != nil {
					logging.Error(err.Error())
					http.Error(w, txnUploadError, http.StatusInternalServerError)
					return
				}
				txns[formFile.Filename] = &buf
			}
		}
	}

	uploadResults := make([][]string, 0)

	for k, v := range txns {
		checksum, err := getCheckSum(v.Bytes())
		if err != nil {
			logging.Error(err.Error())
			http.Error(w, txnUploadError, http.StatusInternalServerError)
			return
		}

		unique, conflictFile, err := dal.CheckUploadedChecksum(checksum)
		if err != nil {
			logging.Error(err.Error())
			http.Error(w, txnUploadError, http.StatusInternalServerError)
			return
		}

		if !unique {
			logging.Error(fmt.Sprintf("%s Conflicted with %s", k, conflictFile))
			uploadResults = append(uploadResults, []string{k, "false", fmt.Sprintf("Conflicted with %s", conflictFile)})
			continue
		}

		//Decode hex string
		encryptedTxn, err := hex.DecodeString(string(v.Bytes()))
		if err != nil || len(encryptedTxn) < 18 {
			logging.Error(fmt.Sprintf("Invalid Txn %s", k))
			http.Error(w, "Invalid txn", http.StatusBadRequest)
			return
		}
		key, _ := hex.DecodeString("234D6C800F15EDE56A9D36690B879962BDE9DA38C26EF614C067BACADA6A803F")
		c, _ := aes.NewCipher(key[:])

		//Get IV from front of byte array
		mode := cipher.NewCBCDecrypter(c, encryptedTxn[0:16])
		decryptBytes := make([]byte, len(encryptedTxn)-16)
		mode.CryptBlocks(decryptBytes, encryptedTxn[16:])

		//Remove PKCS7 Padding bytes and marshal
		paddingByte := decryptBytes[len(decryptBytes)-1]
		padSize := int(paddingByte)
		request := txn.TransactionRequest{}
		err = proto.UnmarshalText(string(decryptBytes[0:len(decryptBytes)-padSize]), &request)
		if err != nil {
			saveTxnFile(k, v.Bytes(), "FailedTxns")
			logging.Error(err.Error())
			http.Error(w, "Invalid txn", http.StatusBadRequest)
			return
		}
		client, clientFound := GRPCclients["SaleAdvice"]
		if !clientFound {
			logging.Warning("ProtoTransaction: Error decoding Action Value")
			uploadResults = append(uploadResults, []string{k, "false", "ProtoTransaction: Error decoding Action Value"})
		} else {
			// make call to relevant service and copy result to protoReply if successful
			logging.Debug("ProtoTransaction: client found, connection state: " + client.GetConnection().GetState().String())
			grpcReply := new(txn.TransactionResponse)
			startTime := time.Now()
			err = rpcHelp.ExecuteGRPC(client, &request, grpcReply, logging)
			if err != nil {
				logging.Warning("ProtoTransaction: ExecuteGRPC failed, error:", err.Error())
				uploadResults = append(uploadResults, []string{k, "false", err.Error()})
				saveTxnFile(k, v.Bytes(), "FailedTxns")
			} else {
				txnLog := txn.LogTxnMessage{}
				txnLog.Txn = grpcReply
				endTime := time.Now()

				txnLog.StartTime = startTime.Format(time.RFC3339Nano)
				txnLog.EndTime = endTime.Format(time.RFC3339Nano)
				txnLog.Duration = fmt.Sprintf("%f", float64(endTime.Sub(startTime).Nanoseconds())/float64(time.Millisecond))

				logging.LogTxn(&txnLog)
				var successText = "false"
				if grpcReply.Response.Completed {
					successText = "true"
					err = dal.AddTxnChecksum(k, checksum)
					if err != nil {
						//noinspection GoUnhandledErrorResult
						logging.Error(fmt.Sprintf("Unable to save transaction file %s checksum: %s", k, err.Error()))
					}

					saveTxnFile(k, v.Bytes(), "ProcessedTxns")

				} else {
					saveTxnFile(k, v.Bytes(), "FailedTxns")
				}
				uploadResults = append(uploadResults, []string{k, successText, grpcReply.Response.ResultText})
			}
		}

	}
	ajaxResponse(w, TransactionUploadModel{Results: uploadResults})
	return

}

func saveTxnFile(filename string, data []byte, folder string) {
	err := dal.SaveFileToDir(filename, data, folder)
	if err != nil {
		if os.IsExist(err) {
			logging.Error("file " + filename + "already exists in " + folder)
		} else {
			logging.Error(err.Error())
		}
	}
}

func getCheckSum(fileData []byte) (string, error) {
	check := sha256.New()
	_, err := check.Write(fileData)
	if err != nil {
		return "", err
	}
	checkSum := check.Sum(nil)
	return hex.EncodeToString(checkSum), nil
}

func getUserAcquirerPermissions(tmsUser *entities.TMSUser) ([]string, error) {
	permissions, err := dal.GetUserAcquirerPermissions(tmsUser)
	if err != nil {
		return nil, err
	}

	acquirers := strings.Split(permissions, ",")
	responseAcquirers := make([]string, 0)
	for _, acquirer := range acquirers {
		if acquirer != "" {
			responseAcquirers = append(responseAcquirers, acquirer)
		}
	}
	return responseAcquirers, nil
}

func generatePINHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	r.ParseForm()
	tid := r.Form.Get("tid")
	intentInt, err := strconv.Atoi(r.Form.Get("intent"))
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, generatePINError, http.StatusInternalServerError)
		return
	}
	intent := dal.OTPIntent(intentInt)

	hasPermission, err := checkUserAcquirePermsByTid(tmsUser, tid)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, generatePINError, http.StatusInternalServerError)
		return
	}
	if !hasPermission {
		logging.Error(userUnauthorisedError)
		http.Error(w, "User is not permitted to perform this action", http.StatusForbidden)
		return
	}
	data, err := services.GenerateOTP(tid, intent, tmsUser)

	if err != nil {
		logging.Error("Error occured while executing generateOTP:" + err.Error())
		http.Error(w, generatePINError, http.StatusInternalServerError)
		return
	}
	ajaxResponse(w, data)
}

func backupDatabaseHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	w.Header().Set("fileName", "nextgen_tms "+time.Now().Format("02-01-2006-15-04-05")+".sql")
	dal.DumpAndDownloadDB(w)
}

func exportTooltipsHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	err := dal.ExportTooltips(w)
	if err != nil {
		logging.Error("An error has been thrown while trying to export TMS tooltips: " + err.Error())
		http.Error(w, exportFailedError, http.StatusInternalServerError)
	}
}

func signonHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {

	if r.URL.Path != "/signon" && r.URL.Path != "/signon/" && r.URL.Path != "" && r.URL.Path != "/" { // only 4 legitimate endpoints for signon. Prevents robots.txt / sitemap.xml etc being picked up.
		http.NotFound(w, r)
		return
	}

	if r.Method == "GET" {
		var p SignOnModel
		p.DbTargetVersion = "1." + strconv.Itoa(DbVersion)
		p.DbVersion = "1." + strconv.Itoa(dal.GetDbVersion())
		p.WebsiteVersion = WebsiteVersion

		renderTemplate(w, r, "signon", p, &entities.TMSUser{LoggedIn: false})
	} else {
		r.ParseForm()
		username := r.Form.Get("username")
		password := r.Form.Get("password")
		user, passwordStatus, err := auth.ValidateUser(username, password)
		if err != nil && (passwordStatus == http.StatusTooManyRequests || passwordStatus == http.StatusUnauthorized || passwordStatus == http.StatusNotFound) {
			w.WriteHeader(passwordStatus)
			passwordStatus = entities.NoChangeRequired
		}
		if err != nil {
			logging.SpecificAudit("LOGON", "Failed logon for "+username)
			p := &SignOnModel{Error: err}
			p.DbTargetVersion = "1." + strconv.Itoa(DbVersion)
			p.DbVersion = "1." + strconv.Itoa(dal.GetDbVersion())
			p.WebsiteVersion = WebsiteVersion

			cookie := http.Cookie{Name: "tms", Value: "", HttpOnly: true, SameSite: http.SameSiteStrictMode, Secure: true}
			http.SetCookie(w, &cookie)

			renderTemplate(w, r, "signon", p, &entities.TMSUser{LoggedIn: false})

		} else if passwordStatus != entities.NoChangeRequired {

			p := &ChangePasswordModel{
				Username:     username,
				ChangeReason: fetchPasswordChangeResponse(passwordStatus),
			}
			renderTemplate(w, r, "passwordChange", p, &entities.TMSUser{LoggedIn: false, PasswordChange: true})

		} else {
			logging.SpecificAudit("LOGON", username+" logged on")
			user.Expires = time.Now().Add(time.Minute * time.Duration(UserTimeout))
			//NEX-9715, Remove the Cookie expiry as after closing browser it still present as security concern found.
			//As Session will work as it is.
			cookie := http.Cookie{Name: "tms", Value: user.Token, HttpOnly: true, SameSite: http.SameSiteStrictMode, Secure: true}
			http.SetCookie(w, &cookie)

			userSessionTokenMapMutex.Lock()
			userSessionTokenMap[user.Token] = user
			userSessionTokenMapMutex.Unlock()

			http.Redirect(w, r, "/search", http.StatusSeeOther)
		}
	}
}

// Creates a new websocket connection for each user upon successful login
func createSessionHandlerWebsocket(w http.ResponseWriter, r *http.Request, user *entities.TMSUser) {
	if user != nil {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logging.Error(err)
		}
		socketConnections[user.Token] = conn
	}
}

// Sends a message to the browser to redirect an expired user to logon
func redirectExpiredUser(user *entities.TMSUser) {
	logging.SpecificAudit("LOGON", user.Username+" logged out due to session expiry")
	conn := socketConnections[user.Token]
	if conn != nil {
		conn.WriteMessage(1, nil)
		conn.Close()
		delete(socketConnections, user.Token)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request, user *entities.TMSUser) {
	logging.SpecificAudit("LOGON", user.Username+" logged out")

	cookie := http.Cookie{Name: "tms", Value: "", HttpOnly: true, SameSite: http.SameSiteStrictMode, Secure: true}

	userSessionTokenMapMutex.Lock()
	delete(userSessionTokenMap, user.Token)
	userSessionTokenMapMutex.Unlock()

	user.Token = ""
	user.LoggedIn = false

	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/signon", http.StatusSeeOther)
}

func handleError(w http.ResponseWriter, err error, user *entities.TMSUser) {
	logging.Error(err)
	renderTemplate(w, nil, "error.html", err.Error(), user)
}

func logHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logging.Information(r.RequestURI)
		h(w, r)
		//TODO use defer to log out time taken

	}
}

func authHandler(h UserHandleFunction, requiredPermission UserPermission) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if requiredPermission == NoSignIn {
			h(w, r, nil)
			return
		}

		var token string
		var user *entities.TMSUser
		var ok bool
		for _, cookie := range r.Cookies() {
			if cookie.Name == "tms" {
				token = cookie.Value

				// Check if the db can be found
				_, err := dal.GetDB()
				if err != nil {
					addHeaderSecurityItems(w, r)
					http.Redirect(w, r, "/signon", http.StatusNotFound)
					return
				}

				// Check the database is the correct version for the TMS
				if !dal.CheckDBVersionMatch(DbVersion) {
					addHeaderSecurityItems(w, r)
					http.Redirect(w, r, "/signon", http.StatusUpgradeRequired)
					return
				}

				if user, ok = userSessionTokenMap[token]; ok {
					user.UserPermissions, err = dal.GetUserPermissions(user.Username)
					if err != nil {
						addHeaderSecurityItems(w, r)
						http.Redirect(w, r, "/signon", http.StatusNotFound)
						return
					}

					user.Expires = time.Now().Add(time.Minute * time.Duration(UserTimeout))
					cookie := http.Cookie{Name: "tms", Value: user.Token, HttpOnly: true, SameSite: http.SameSiteStrictMode, Secure: true}
					http.SetCookie(w, &cookie)

					userSessionTokenMapMutex.Lock()
					userSessionTokenMap[user.Token] = user
					userSessionTokenMapMutex.Unlock()
				} else {
					token = ""
				}
			}
		}

		if token == "" {
			addHeaderSecurityItems(w, r)
			http.Redirect(w, r, "/signon", http.StatusSeeOther)
		} else {
			if checkUserPermissions(requiredPermission, user) {
				h(w, r, user)
			} else {
				NoPermissionsRedirect(w, r)
			}
		}
	}
}

func checkUserPermissions(requiredPermission UserPermission, user *entities.TMSUser) bool {
	authorised := false
	switch requiredPermission {
	case None:
		// No permission required
		authorised = true
	case SiteWrite:
		if user.UserPermissions.SiteWrite {
			authorised = true
		}
	case SiteDelete:
		if user.UserPermissions.SiteDelete {
			authorised = true
		}
	case ChangeApprovalRead:
		if user.UserPermissions.ChangeApprovalRead {
			authorised = true
		}
	case ChangeApprovalWrite:
		if user.UserPermissions.ChangeApprovalWrite {
			authorised = true
		}
	case BulkChangeApproval:
		{
			if user.UserPermissions.BulkChangeApproval {
				authorised = true
			}
		}
	case AddCreate:
		if user.UserPermissions.AddCreate {
			authorised = true
		}
	case BulkUpdates:
		if user.UserPermissions.BulkUpdates {
			authorised = true
		}
	case UserManagement:
		if user.UserPermissions.UserManagement {
			authorised = true
		}
	case TransactionViewer:
		if user.UserPermissions.TransactionViewer {
			authorised = true
		}
	case Reporting:
		if user.UserPermissions.Reporting {
			authorised = true
		}
	case ChangeHistoryView:
		if user.UserPermissions.ChangeHistoryView {
			authorised = true
		}
	case EditPasswords:
		if user.UserPermissions.EditPasswords {
			authorised = true
		}
	case Fraud:
		if user.UserPermissions.Fraud {
			authorised = true
		}
	case PermissionGroups:
		if user.UserPermissions.PermissionGroups {
			authorised = true
		}
	case UserManagementAudit:
		if user.UserPermissions.UserManagementAudit {
			authorised = true
		}
	case BulkImport:
		if user.UserPermissions.BulkImport {
			authorised = true
		}
	case ContactEdit:
		if user.UserPermissions.ContactEdit {
			authorised = true
		}
	case OfflinePIN:
		if user.UserPermissions.OfflinePIN {
			authorised = true
		}
	case DbBackup:
		if user.UserPermissions.DbBackup {
			authorised = true
		}
	case TerminalFlagging:
		if user.UserPermissions.TerminalFlagging {
			authorised = true
		}
	case ChainDuplication:
		if user.UserPermissions.ChainDuplication {
			authorised = true
		}
	case EditToken:
		if user.UserPermissions.EditToken {
			authorised = true
		}
	case PaymentServices:
		if user.UserPermissions.PaymentServices {
			authorised = true
		}
	case FileManagement:
		if user.UserPermissions.FileManagement {
			authorised = true
		}
	case LogoManagement:
		if user.UserPermissions.LogoManagement {
			authorised = true
		}
	case SouhoolaLogin:
		if user.UserPermissions.SouhoolaLogin {
			authorised = true
		}

	}
	return authorised
}

func NoPermissionsRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/search", http.StatusUnauthorized)
}

func parseConfig() error {
	ApplicationName = cfg.GetApplicationName("TMSWebSite")
	listeningPort = cfg.GetString("Port", "8081")
	authenticateUserGroups := cfg.GetBool("authenticateUserGroups", true)

	fs := os.Getenv("NextGenFS")
	if fs == "" {
		return errors.New("NextGenFS environment variable not set. Cannot read configuration")
	}
	configState := os.Getenv("NextGenConfigState")
	if configState == "" {
		return errors.New("NextGenConfigState environment variable not set. Cannot read configuration")
	}

	HTTPSKeyFile = filepath.Join(fs, configState, cfg.GetString("HTTPSKeyFile", "TMSWebSite/key.pem"))
	HTTPSCertFile = filepath.Join(fs, configState, cfg.GetString("HTTPSCertFile", "TMSWebSite/cert.pem"))
	ReportDir = filepath.Join(fs, configState, cfg.GetString("ReportDir", "TMSWebSite/Reports"))
	XMPPHost = cfg.GetString("XMPPHost", "192.168.2.167")
	XMPPDomain = cfg.GetString("XMPPDomain", "@ngenius.network.ae")
	XMPPPassword = cfg.GetString("XMPPPassword", "12345")

	shortenedfs := "../" + filepath.Base(fs)
	RPIKeyFile = filepath.Join(shortenedfs, configState, cfg.GetString("RPIKeyFile", "RPI/key.pem"))
	RPICertFile = filepath.Join(shortenedfs, configState, cfg.GetString("RPICertFile", "RPI/cert.pem"))

	useLdaps := cfg.GetBool("UseLDAPS", false)
	mutualAuthLdaps := cfg.GetBool("UseMutualAuthLDAPS", false)

	var ldapsCertPath string
	var ldapsKeyPath string

	if useLdaps && mutualAuthLdaps {
		ldapsCertPath = filepath.Join(fs, configState, cfg.GetString("LDAPSClientCert", ""))
		ldapsKeyPath = filepath.Join(fs, configState, cfg.GetString("LDAPSClientKey", ""))
	}

	ldapSettings := auth.LDAPSettings{
		IPAddress:                   cfg.GetString("LDAPIPAddress", "192.168.2.125"),
		Port:                        cfg.GetString("LDAPPort", "389"),
		Username:                    cfg.GetString("LDAPUsername", "CN=Read Only,CN=Users,DC=NETWORKINTL,DC=COM"),
		Password:                    cfg.GetString("LDAPPassword", "Csc_12345"),
		BaseDN:                      cfg.GetString("LDAPBaseDN", "cn=Users,dc=NETWORKINTL,dc=COM"),
		ObjectClass:                 cfg.GetString("LDAPObjectClass", "user"),
		SearchAttribute:             cfg.GetString("LDAPSearchAttribute", "sAMAccountName"),
		IdentityAttribute:           cfg.GetString("LDAPIdentityAttribute", "sAMAccountName"),
		UseTLS:                      cfg.GetBool("LDAPUseTLS", true),
		AuthenticateUserPermissions: authenticateUserGroups,
		UserGroup:                   cfg.GetString("LDAPUserGroup", "CN=TMSUsers"),
		AdminGroup:                  cfg.GetString("LDAPAdminGroup", "CN=TMSAdministrators"),
		GlobalAdminGroup:            cfg.GetString("LDAPGlobalAdminGroup", "CN=GlobalAdmins"),
		UseLDAPS:                    useLdaps,
		LDAPSClientCertPath:         ldapsCertPath,
		LDAPSClientKeyPath:          ldapsKeyPath,
		UseMutualAuthLDAPS:          mutualAuthLdaps,
	}

	auth.SetLDAPSettings(ldapSettings)

	mongoTimeout, err := strconv.Atoi(cfg.GetString("MongoDBTimeout", "60"))
	if err != nil {
		log.Println("MongoDBTimeout Error, check config")
	}

	mongoSettings := dal.MongoSettings{
		MongoDBAddress:    cfg.GetString("MongoDBAddress", "192.168.2.128:27017"),
		LoggingDatabase:   cfg.GetString("LoggingDatabase", "logStore"),
		TxnCollection:     cfg.GetString("TxnCollection", "Txn2"),
		EcomTxnCollection: cfg.GetString("EcomTxnCollection", "EcomTxn"),
		Username:          cfg.GetString("MongoDBUsername", ""),
		Password:          cfg.GetString("MongoDBPassword", ""),
		Timeout:           mongoTimeout,
	}

	dal.SetMongoSettings(mongoSettings)
	dal.MaxDBConnections = cfg.GetInt("MaxDatabaseConnections", 200)
	UserTimeout = cfg.GetInt("UserTimeout", 30)
	dal.SetConnectionString(cfg.GetString("SQLDatabaseConnectionString", "admin:Csc_12345@tcp(localhost:3306)/NextGen_TMS"))
	PasswordKey = cfg.GetString("PasswordKey", "")

	automationUserString := cfg.GetString("AutomationUsers", "automationGlobal,automationAcquirer,automationUser")
	dal.AutomationUsers = strings.Split(automationUserString, ",")

	userManagementFileDirectory = filepath.Join(fs, configState, cfg.GetString("UserManagementFileDir", "FileServer/UserManagement"))

	permissionMap := cfg.Get("PermissionMap", dal.PermissionMap{})

	//If we have a PermissionMap in the Config, loop through and assign to the Permissions Map
	// that will then be used to add Permissions during SubModule without requiring a TMS change
	if permissionMap != nil {
		dal.Permissions = make(map[string]*dal.Permission)
		value, err := json.Marshal(permissionMap)
		if err != nil {
			return nil
		}

		permissions := make([]dal.PermissionMap, 0)
		_ = json.Unmarshal(value, &permissions)

		for _, entry := range permissions {
			dal.Permissions[entry.ModuleName] = &dal.Permission{
				PermissionSale:   entry.Permissions.PermissionSale,
				PermissionVoid:   entry.Permissions.PermissionVoid,
				PermissionRefund: entry.Permissions.PermissionRefund,
			}
		}
	}
	dal.SetBackupLocation(cfg.GetString("SQLDatabaseBackupLocation", "backups"))

	return nil
}

func getUsersForTid(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	r.ParseForm()
	tid := r.Form.Get("TID")
	tidID, err := strconv.Atoi(tid)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Error parsing TID", http.StatusBadRequest)
		return
	}

	overriden := true
	users, err := dal.GetUsersForTid(tidID)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Error retrieving users", http.StatusInternalServerError)
		return
	}

	if len(users) < 1 {
		overriden = false
		site := r.Form.Get("siteId")
		siteId, err := strconv.Atoi(site)
		if err != nil {
			logging.Error(err.Error())
			http.Error(w, "Error parsing Site", http.StatusBadRequest)
			return
		}
		users, err = dal.GetUsersForSite(siteId)
		if err != nil {
			logging.Error(err.Error())
			http.Error(w, "Error retrieving users", http.StatusInternalServerError)
			return
		}
	}

	for i, user := range users {
		for j, module := range user.Modules {
			if module == "upiSale" {
				user.Modules[j] = "upiSale,upiVoid,upiRefund"

				//Remove the elements at i+1 and i+2 because these were modules upiVoid and upiRefund
				user.Modules = append(user.Modules[:j+1], user.Modules[j+3:]...)

				users[i] = user
				break
			}
		}
	}

	modules, err := dal.GetAvailableModules()
	modules = dal.GetSubmodules(modules)
	friendlyModules := makeModulesFriendly(modules)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Unable to retrieve modules", http.StatusInternalServerError)
		return
	}

	profileId, err := services.GetSiteProfileFromTid(tid)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Unable to retrieve modules", http.StatusInternalServerError)
		return
	}

	modules = append(modules, "X-Read", "Z-Read", "pinSet", "pinChange", "transactionHistory")
	friendlyModules = append(friendlyModules, "X-Read", "Z-Read", "PIN Set", "PIN Change", "Transaction History")

	isExits, err := dal.CheckDataGroupExistsProfile(profileId, dal.PUSHPAYMENTS)
	if isExits {
		modules = append(modules, "queuedTransaction")
		friendlyModules = append(friendlyModules, "Queued Transaction")
	}

	for i, module := range modules {
		if module == "upi" {
			modules[i] = "upiSale,upiVoid,upiRefund"
		}
	}

	ajaxResponse(w, TidUsersModel{Users: users, Modules: modules, FriendlyModules: friendlyModules, Overriden: overriden, HasSavePermission: checkUserPermissions(SiteWrite, tmsUser), HasPasswordPermission: checkUserPermissions(EditPasswords, tmsUser)})
	return

}

func getUsersForSite(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	r.ParseForm()
	siteId := r.Form.Get("siteId")
	site, err := strconv.Atoi(siteId)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Invalid site Id", http.StatusBadRequest)
		return
	}

	profileId := r.Form.Get("profileId")
	profile, err := strconv.Atoi(profileId)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Invalid profile Id", http.StatusBadRequest)
		return
	}

	users, err := dal.GetUsersForSite(site)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Unable to retrieve users", http.StatusInternalServerError)
		return
	}

	for i, user := range users {
		for j, module := range user.Modules {
			if module == "upiSale" {
				user.Modules[j] = "upiSale,upiVoid,upiRefund"

				//Remove the elements at i+1 and i+2 because these were modules upiVoid and upiRefund
				user.Modules = append(user.Modules[:j+1], user.Modules[j+3:]...)

				users[i] = user
				break
			}
		}
	}

	modules, err := dal.GetAvailableModules()
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Unable to retrieve modules", http.StatusInternalServerError)
		return
	}

	modules = dal.GetSubmodules(modules)
	friendlyModules := makeModulesFriendly(modules)

	modules = append(modules, "X-Read", "Z-Read", "pinSet", "pinChange", "transactionHistory")
	friendlyModules = append(friendlyModules, "X-Read", "Z-Read", "PIN Set", "PIN Change", "Transaction History")

	isExits, err := dal.CheckDataGroupExistsProfile(profile, dal.PUSHPAYMENTS)
	if isExits {
		modules = append(modules, "queuedTransaction")
		friendlyModules = append(friendlyModules, "Queued Transaction")
	}

	for i, module := range modules {
		if module == "upi" {
			modules[i] = "upiSale,upiVoid,upiRefund"
		}
	}

	if !tmsUser.UserPermissions.EditPasswords { // If you cannot edit the passwords, you cannot view them.
		for index, user := range users {
			user.PIN = "*****"
			users[index] = user
		}
	}

	superPins := make([]string, 1)
	siteSuperPin, err := dal.GetDataElementValue(profile, "superPIN", "userMgmt")
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Error getting superPins", http.StatusInternalServerError)
		return
	}

	superPins = append(superPins, siteSuperPin)

	tidSuperPins, err := dal.GetTIDdatavalueFromSite(site, "superPIN", "userMgmt")
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Error getting tid superPins", http.StatusInternalServerError)
		return
	}

	superPins = append(superPins, tidSuperPins...)

	ajaxResponse(w, SitesUsersModel{Users: users, Modules: modules, FriendlyModules: friendlyModules, SuperPins: superPins, HasSavePermission: checkUserPermissions(SiteWrite, tmsUser), HasPasswordPermission: checkUserPermissions(EditPasswords, tmsUser)})
	return
}

func handleUserCsvUploadResultExport(w http.ResponseWriter, r *http.Request, _ *entities.TMSUser) {
	err := r.ParseForm()
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}
	profile := r.Form.Get("ProfileId")
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Bad site id", http.StatusBadRequest)
		return
	}
	profileId, err := strconv.Atoi(profile)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Bad profile id", http.StatusBadRequest)
		return
	}
	merchantId, err := dal.GetDataElementValue(profileId, "merchantNo", "store")

	var uploadResults = make([]dal.UserUpdateResult, 0)
	result := r.Form.Get("UploadResult")
	err = json.Unmarshal([]byte(result), &uploadResults)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Error parsing upload results", http.StatusBadRequest)
		return
	}

	modules, err := dal.GetAvailableModules()
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Unable to get modules", http.StatusInternalServerError)
		return
	}
	modules = dal.GetSubmodules(modules)
	modules = append(modules, "X-Read", "Z-Read", "pinSet", "pinChange", "transactionHistory")

	header := append([]string{"Username", "PIN", "TID"}, modules...)
	header = append(header, "Action", "Successful", "Error")
	title := []string{"Export for Merchant ID:", merchantId}
	csvWriter := csv.NewWriter(w)
	csvWriter.Write(title)
	csvWriter.Write(header)

	for _, result := range uploadResults {
		record := []string{result.User.Username, "=\"" + result.User.PIN + "\""}
		if result.User.TidId != 0 {
			record = append(record, "=\""+strconv.Itoa(result.User.TidId)+"=\"")
		} else {
			record = append(record, "")
		}
		for _, module := range modules {
			if SliceComparisonHelpers.SlicesOfStringContains(result.User.Modules, module) {
				record = append(record, "Y")
			} else {
				record = append(record, "N")
			}
		}
		record = append(record, result.Action)
		if result.Result.Success {
			record = append(record, "Y")
		} else {
			record = append(record, "N")
		}
		if result.Result.ErrorMessage != "" {
			record = append(record, result.Result.ErrorMessage)
		} else {
			record = append(record, "")
		}

		csvWriter.Write(record)
	}
	csvWriter.Flush()
}

func handleUserCsvUpload(w http.ResponseWriter, r *http.Request, _ *entities.TMSUser) {
	err := r.ParseMultipartForm(10000)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Unable to parse CSV", http.StatusBadRequest)
		return
	}
	site := r.Form.Get("SiteId")
	siteId, err := strconv.Atoi(site)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Unable to parse siteId", http.StatusBadRequest)
		return
	}

	var buf bytes.Buffer
	var outboundReq bytes.Buffer
	mw := multipart.NewWriter(&outboundReq)

	for _, header := range r.MultipartForm.File {
		file, err := header[0].Open()
		if err != nil {
			logging.Error(err.Error())
			http.Error(w, userCsvUploadError, http.StatusInternalServerError)
			return
		} else {
			if _, err := io.Copy(&buf, file); err != nil {
				logging.Error(err.Error())
				http.Error(w, userCsvUploadError, http.StatusInternalServerError)
				return
			}

			part, err := mw.CreateFormFile("file."+header[0].Filename, header[0].Filename)
			_, err = part.Write(buf.Bytes())
			if err != nil {
				logging.Error(err.Error())
				http.Error(w, userCsvUploadError, http.StatusInternalServerError)
				return
			}
		}
	}
	bom := make([]byte, 0)
	bom = append(bom, 0xEF, 0xBB, 0xBF)
	var fileBytes []byte
	if len(buf.Bytes()) > 3 && bytes.Compare(buf.Bytes()[0:3], bom) == 0 {
		fileBytes = buf.Bytes()[3:]
	} else {
		fileBytes = buf.Bytes()
	}
	reader := bytes.NewReader(fileBytes)
	usersToAddToSite, usersToAddToTids, userIdsToDeleteFromSite, tidUserIdsToDelete, err := unmarshalUsersUploadCsv(&dal.SiteManagementDal{}, reader, siteId)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, userCsvUploadError, http.StatusBadRequest)
		return
	}

	result, err := saveUserChanges(&dal.SiteManagementDal{}, siteId, usersToAddToSite, userIdsToDeleteFromSite, usersToAddToTids, tidUserIdsToDelete, true)
	json, err := json.Marshal(result)

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)

	if err != nil {
		logging.Error(err.Error())
		http.Error(w, userCsvUploadError, http.StatusBadRequest)
		return
	}
}

// Takes a csv containing user info and applies it to the site with the mid matching that in the file.
func applyCsvToSite(path, name string) bool {
	reader, err := os.Open(path)
	defer reader.Close()
	if err != nil {
		logging.Error(fmt.Sprintf("could not open the file %s", name))
		return false
	}

	csvReader := csv.NewReader(reader)
	record, err := csvReader.Read()
	var mid string
	if err != nil {
		logging.Error(fmt.Sprintf("could not read the file %s", name))
		return false
	}
	if len(record) > 1 {
		mid = record[1]
	} else {
		logging.Error(fmt.Sprintf("cannot read the mid from the csv file %s", name))
		return false
	}

	profileId, err := dal.GetSiteIDFromMerchantID(mid)
	if err != nil || profileId == "" {
		logging.Error(fmt.Sprintf("cannot not resolve the profile id for the mid %s used in file: %s", mid, name))
		return false
	}
	siteId, err := strconv.Atoi(profileId)
	if err != nil {
		logging.Error(fmt.Sprintf("cannot convert the profile id (%s) to numeric for file: %s", profileId, name))
		return false
	}

	newReader, err := os.Open(path)
	defer newReader.Close()
	if err != nil {
		logging.Error(fmt.Sprintf("could not open the file %s", name))
		return false
	}
	usersToAddToSite, usersToAddToTids, userIdsToDeleteFromSite, tidUserIdsToDelete, err := unmarshalUsersUploadCsv(&dal.SiteManagementDal{}, newReader, siteId)
	if err != nil {
		logging.Error(fmt.Sprintf("error importing the file %s - invalid formatting or users", name))
		return false
	}

	failedValidationUsers, err := validatePedUsers(&dal.SiteManagementDal{}, siteId, usersToAddToSite, usersToAddToTids, false)
	if err != nil {
		return false
	}
	if failedValidationUsers != nil {
		if len(failedValidationUsers) > 0 {
			logging.Error(fmt.Sprintf("error importing the file %s - invalid formatting or users", name))
			return false
		}
	}

	_, err = saveUserChanges(&dal.SiteManagementDal{}, siteId, usersToAddToSite, userIdsToDeleteFromSite, usersToAddToTids, tidUserIdsToDelete, false)
	if err != nil {
		logging.Error(fmt.Sprintf("error applying file %s to the site %s", name, mid))
		return false
	}

	return true
}

// Checks every 30 seconds if any new user management files have been uploaded to the server
func listenForUserManagementFiles() {
	for {
		fileInfo, err := listUserManagementFiles()
		fileCount := len(fileInfo)

		// Process any files that have been added since the last check
		if err == nil && fileCount > 0 {
			log.Println(fmt.Sprintf("Found %d user management files to process", fileCount))
			for _, f := range fileInfo {
				path := filepath.Join(userManagementFileDirectory, f.Name())
				newPath := ""
				file, err := os.ReadFile(path)

				if err == nil {
					mimeType := http.DetectContentType(file)
					mimeType = strings.Split(mimeType, ";")[0]
					if mimeType == "text/plain" || mimeType == "text/csv" {
						if err == nil {
							success := applyCsvToSite(path, f.Name())
							if success {
								log.Println(fmt.Sprintf("successfully processed file - %s", f.Name()))
								newPath = filepath.Join(userManagementFileDirectory, "success", f.Name())
							} else {
								newPath = filepath.Join(userManagementFileDirectory, "failed", f.Name())
							}
						} else {
							log.Println(fmt.Sprintf("cannot read file %s", f.Name()))
							newPath = filepath.Join(userManagementFileDirectory, "failed", f.Name())
						}
					} else {
						newPath = filepath.Join(userManagementFileDirectory, "failed", f.Name())
						log.Println(fmt.Sprintf("expected %s to be a csv file, instead received: %s", f.Name(), mimeType))
					}
				} else {
					log.Println(fmt.Sprintf("cannot read file %s", f.Name()))
					newPath = filepath.Join(userManagementFileDirectory, "failed", f.Name())
				}

				err = fileServer.NewFsReader(config.FileserverURL).MoveFile(path, newPath)

				if err != nil {
					log.Println(fmt.Sprintf("Error moving file from %s to %s", path, newPath))
				}
			}

		}
		// Check every 30 seconds to see if new files have been uploaded
		time.Sleep(30 * time.Second)
	}
}

// Lists any user management files in the directory with the name stored in the "userManagementFileDirectory" config item
func listUserManagementFiles() ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir(userManagementFileDirectory)
	if err != nil {
		return nil, err
	}

	fileInfo := make([]os.FileInfo, 0)
	for _, f := range files {
		if !f.IsDir() {
			fileInfo = append(fileInfo, f)
		}
	}
	return fileInfo, nil
}

func handleUserCsvExport(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	err := r.ParseForm()
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}
	site := r.Form.Get("SiteId")
	profile := r.Form.Get("ProfileId")
	siteId, err := strconv.Atoi(site)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Bad site id", http.StatusBadRequest)
		return
	}
	profileId, err := strconv.Atoi(profile)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Bad profile id", http.StatusBadRequest)
		return
	}
	merchantId, err := dal.GetDataElementValue(profileId, "merchantNo", "store")
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "An error occured while retrieving information from database", http.StatusBadRequest)
		return
	}

	userModels, err := dal.ExportSiteUsers(siteId)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Unable to get users", http.StatusInternalServerError)
		return
	}
	csvWriter := csv.NewWriter(w)
	modules, err := dal.GetAvailableModules()
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Unable to get modules", http.StatusInternalServerError)
		return
	}
	modules = dal.GetSubmodules(modules)
	modules = append(modules, "X-Read", "Z-Read", "pinSet", "pinChange", "transactionHistory")
	header := append([]string{"Username", "PIN", "TID"}, modules...)
	title := []string{"Export for Merchant ID:", merchantId}
	csvWriter.Write(title)
	csvWriter.Write(header)
	for _, model := range *userModels {
		record := []string{model.Username, "=\"" + model.PIN + "\"", "=\"" + model.Tid + "\""}
		for _, module := range modules {
			if model.Modules[module] {
				record = append(record, "Y")
			} else if model.Modules["upiSale"] && model.Modules["upiVoid"] && model.Modules["upiRefund"] && module == "upi" {
				record = append(record, "Y")
			} else {
				record = append(record, "N")
			}
		}
		csvWriter.Write(record)
	}
	csvWriter.Flush()
}

func makeModulesFriendly(modules []string) []string {
	var friendlyNames = make([]string, 0)
	for _, mod := range modules {
		if isAllUpperCase(mod) {
			friendlyNames = append(friendlyNames, mod)
			continue
		}
		i := 0
		var words = make([]string, 0)
		cap := strings.ToUpper(mod[0:1])
		mod = cap + mod[1:]

		for name := mod; name != ""; name = name[i:] {
			i = strings.IndexFunc(name[1:], unicode.IsUpper) + 1
			if i <= 0 {
				i = len(name)
			}
			friendlyMod := name[:i]
			friendlyMod = strings.Replace(friendlyMod, "Gratuity", "Tip", -1)
			words = append(words, friendlyMod)
		}
		friendlyNames = append(friendlyNames, strings.Join(words, " "))

	}
	return friendlyNames
}

// Checks a string to determine if it is all caps or not
func isAllUpperCase(s string) bool {
	for _, r := range s {
		// Returns false if any letter in the string is lower case
		if !unicode.IsUpper(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func unmarshalUsersUploadCsv(d siteUsersDal, reader io.Reader, siteId int) ([]*entities.SiteUser, map[int][]*entities.SiteUser, []int, []int,
	error) {

	read := csv.NewReader(reader)
	read.FieldsPerRecord = -1
	read.TrimLeadingSpace = true

	// Ditch the first line
	_, err := read.Read()
	if err != nil {
		return nil, nil, nil, nil, err
	}
	csvReader := csvutil.Reader(read)

	decoder, err := csvutil.NewDecoder(csvReader)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	hasUsername := false
	hasPIN := false
	hasTid := false
	for _, column := range decoder.Header() {
		if strings.ToLower(column) == "username" {
			hasUsername = true
		} else if strings.ToLower(column) == "pin" {
			hasPIN = true
		} else if strings.ToLower(column) == "tid" {
			hasTid = true
		}
	}

	if !hasUsername || !hasPIN || !hasTid {
		return nil, nil, nil, nil, errors.New("invalid csv format")
	}

	siteUsers := make([]*entities.SiteUser, 0)
	tidUsers := make(map[int][]*entities.SiteUser)
	existingSiteUsers, err := d.GetUsersForSite(siteId)
	if err != nil {
		return nil, nil, nil, nil, errors.New("Failed to get existing site users: " + err.Error())
	}

	existingSiteUsersByUsername := make(map[string]entities.SiteUser)
	for _, user := range existingSiteUsers {
		existingSiteUsersByUsername[user.Username] = user
	}

	existingTidUsers := make(map[int][]entities.SiteUser)
	deletedSiteUserIds := make([]int, 0)
	deletedTidUserIds := make([]int, 0)
	for {
		user := entities.SiteUser{}
		user.Modules = make([]string, 0)
		if err = decoder.Decode(&user); err == io.EOF {
			break
		} else if err != nil {
			return nil, nil, nil, nil, err
		}
		user.PIN = strings.Replace(user.PIN, "=\"", "", -1)
		user.PIN = strings.Replace(user.PIN, "\"", "", -1)

		// Validate PIN length
		if len(user.PIN) < 4 || len(user.PIN) > 5 {
			return nil, nil, nil, nil, errors.New(fmt.Sprintf("PIN %s must be either 4 or 5 characters",
				user.PIN))
		}

		// Validate that the PIN is numeric
		if _, err := strconv.Atoi(user.PIN); err != nil {
			return nil, nil, nil, nil, errors.New(fmt.Sprintf("PIN %s must only contain digits (0-9)",
				user.PIN))
		}

		if len(strings.Trim(user.Username, " ")) < 1 {
			return nil, nil, nil, nil, errors.New("Username must not be blank")
		}

		if len(strings.Trim(user.Username, " ")) > 10 {
			return nil, nil, nil, nil, errors.New("Username must not be more than 10 characters")
		}

		if strings.ContainsAny(user.Username, "!#$%&'()*+,.:;<=>?@[]{}~|`\"/\\") {
			return nil, nil, nil, nil, errors.New("Username must not contain special characters")
		}

		userTidId := -1
		deleteUser := false
		for _, i := range decoder.Unused() {
			if strings.ToUpper(strings.Trim(decoder.Header()[i], " ")) == "TID" {
				userTidId, err = unmarshalCsvTid(decoder.Record()[i])
				if err != nil {
					return nil, nil, nil, nil, err
				}

				if existingTidUsers[userTidId] == nil {
					tidUsers, err := d.GetUsersForTid(userTidId)
					if err != nil {
						return nil, nil, nil, nil, err
					}

					for _, user := range tidUsers {
						existingTidUsers[userTidId] = append(existingTidUsers[userTidId], user)
					}
				}
			} else {
				if strings.ToUpper(decoder.Record()[i]) == "Y" || strings.ToUpper(decoder.Record()[i]) == "TRUE" {
					if decoder.Header()[i] == "delete" {
						deleteUser = true
					} else {
						user.Modules = append(user.Modules, decoder.Header()[i])

						if decoder.Header()[i] == "upi" {
							user.Modules = append(user.Modules, "upiSale")
							user.Modules = append(user.Modules, "upiRefund")
							user.Modules = append(user.Modules, "upiVoid")
						}
					}
				}
			}
		}

		if deleteUser {
			if userTidId == -1 {
				for _, existingUser := range existingSiteUsers {
					if user.Username == existingUser.Username {
						deletedSiteUserIds = append(deletedSiteUserIds, existingUser.UserId)
						break
					}
				}
			} else {
				for tid, existingTidUsers := range existingTidUsers {
					if tid == userTidId {
						for _, existingUser := range existingTidUsers {
							if existingUser.Username == user.Username {
								deletedTidUserIds = append(deletedTidUserIds, existingUser.UserId)
							}
						}
					}
				}
			}
		} else if userTidId == -1 {
			if existingUser, found := existingSiteUsersByUsername[user.Username]; found {
				user.UserId = existingUser.UserId
			}
			siteUsers = append(siteUsers, &user)
		} else {
			if existingUsers, found := existingTidUsers[userTidId]; found {
				for _, existingUser := range existingUsers {
					if existingUser.Username == user.Username {
						user.UserId = existingUser.UserId
						break
					}
				}
			}

			tidUsers[userTidId] = append(tidUsers[userTidId], &user)
		}
	}

	return siteUsers, tidUsers, deletedSiteUserIds, deletedTidUserIds, nil
}

func unmarshalCsvTid(value string) (int, error) {
	tid := strings.Trim(value, " ")
	tid = strings.Replace(tid, "=\"", "", -1)
	tid = strings.Replace(tid, "\"", "", -1)

	if tid == "" {
		return -1, nil
	} else {
		tidId, err := strconv.Atoi(tid)
		if err != nil {
			return -1, errors.New("bad tid: " + tid)
		} else {
			return tidId, nil
		}
	}
}

type siteUserValidationError struct {
	message string
}

func (e *siteUserValidationError) Error() string {
	return e.message
}

func validateProfileUser(users []*entities.SiteUser) []*dal.UserUpdateResult {
	if users == nil || len(users) == 0 {
		return nil
	}
	resultArr := make([]*dal.UserUpdateResult, 0)
	pins := make(map[string]bool, 0)
	names := make(map[string]bool, 0)
	for _, u := range users {
		userResult := dal.UserUpdateResult{User: *u}
		if pins[u.PIN] {
			userResult.Result.SetError(errors.New("Duplicate Pins"))
			resultArr = append(resultArr, &userResult)
			for _, otherUser := range users {
				if otherUser.PIN == u.PIN && otherUser.Username != u.Username {
					otherUserResult := dal.UserUpdateResult{User: *otherUser}
					otherUserResult.Result.SetError(errors.New("Duplicate Pins"))
					resultArr = append(resultArr, &otherUserResult)
				}
			}
			continue
		} else {
			pins[u.PIN] = true
		}

		if names[u.Username] {
			userResult.Result.SetError(fmt.Errorf("Duplicate user"))
			resultArr = append(resultArr, &userResult)
			for _, otherUser := range users {
				if otherUser.Username == u.Username && otherUser.PIN != u.PIN {
					otherUserResult := dal.UserUpdateResult{User: *otherUser}
					otherUserResult.Result.SetError(fmt.Errorf("Duplicate user"))
					resultArr = append(resultArr, &otherUserResult)
				}
			}
			continue
		} else {
			names[u.Username] = true
		}
	}

	return resultArr
}

func saveSiteUsersHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	r.ParseForm()
	siteId, err := strconv.Atoi(r.Form.Get("SiteId"))
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Error with request: "+err.Error(), http.StatusBadRequest)
		return
	}

	users := r.Form.Get("Users")
	var userList = make([]*entities.SiteUser, 0)
	err = json.Unmarshal([]byte(users), &userList)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Error with request: "+err.Error(), http.StatusBadRequest)
		return
	}

	newUsers := r.Form.Get("NewUsers")
	var newUserList = make([]*entities.SiteUser, 0)

	err = json.Unmarshal([]byte(newUsers), &newUserList)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Error with request: "+err.Error(), http.StatusBadRequest)
		return
	}

	userList = append(userList, newUserList...)

	// checking special characters
	re, _ := regexp.Compile(`[^\w]`)
	for i, _ := range userList {
		userList[i].Username = strings.TrimSpace(userList[i].Username)
		found := re.MatchString(userList[i].Username)
		if found {
			http.Error(w, "Username must not contain special characters", http.StatusInternalServerError)
			return
		}
	}

	//@NEX-12567 need this flag status for the same
	var profileID, _ = strconv.Atoi(r.Form.Get("profileID"))
	var profileTypeName = r.Form.Get("profileTypeName")
	if err := processFlagging(r, profileID, siteId, profileTypeName, tmsUser); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	deletedUsers := r.Form.Get("DeletedUsers")
	if err := processUsersFlagging(siteId, profileTypeName, users, newUsers, deletedUsers, tmsUser); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
}

func processUsersOverride(profileTypeName, users string, deletedUsers string, newUsers string, tmsUser *entities.TMSUser, tidUsersOverride bool, tidStr string, tidInt, siteId int) error {
	profileID, err := dal.GetProfileIdFromTID(tidStr)
	if err != nil {
		return err
	}

	if profileID < 1 {
		profileTypeId, err := dal.GetProfileTypeId("tid", tidStr, 1, tmsUser.Username)
		if err != nil {
			logger.GetLogger().Error(err.Error())
			return err
		}
		err = dal.CreateTidoverrideAndSaveNewprofileChange(siteId, profileTypeId, tmsUser.Username, dal.ApproveCreate, "Override Created", tidStr, tidInt, 1)
		if err != nil {
			logger.GetLogger().Error(err.Error())
			return err
		}

		profileID, err = dal.GetProfileIdFromTID(strings.TrimSpace(tidStr))
		if err != nil {
			logging.Error("An error occured during executing GetProfileIdFromTID", err.Error())
			return err
		}
	}

	dataElementID, err := dal.GetDataElementByName("core", "users")
	if err != nil {
		logging.Error(err.Error())
		return errors.New("unable to get dataElementID for core users")
	}
	if users != "[]" || newUsers != "[]" || deletedUsers != "[]" {
		elements := map[int]string{
			dataElementID: fmt.Sprintf("{\"tidID\":%d,\"updatedUsers\":%s,\"newUsers\":%s, \"deletedUsers\":%s}", tidInt, users, newUsers, deletedUsers),
		}
		saveDataElements(int(profileID), 0, elements, 0, tmsUser, false, 0, profileTypeName)
	}
	return nil
}

func processUsersFlagging(siteID int, profileTypeName, users string, newUsers string, deletedUsers string, tmsUser *entities.TMSUser) error {

	profileID, err := dal.GetProfileIdForSite(siteID)
	if err != nil {
		return err
	}

	dataElementID, err := dal.GetDataElementByName("core", "users")
	if err != nil {
		logging.Error(err.Error())
		return errors.New("Unable to get dataElementID")
	}

	if users != "[]" || newUsers != "[]" || deletedUsers != "[]" {
		elements := map[int]string{
			dataElementID: fmt.Sprintf("{\"updatedUsers\":%s, \"newUsers\":%s, \"deletedUsers\":%s}", users, newUsers, deletedUsers),
		}
		saveDataElements(profileID, siteID, elements, 0, tmsUser, false, 0, profileTypeName)
	}
	return nil
}

type siteUsersDal interface {
	GetTidUsersForSite(siteID int) ([]entities.SiteUser, error)
	AddOrUpdateSiteUsers(siteId int, users []*entities.SiteUser) ([]*dal.UserUpdateResult, error)
	AddOrUpdateTidUserOverride(tid int, users []*entities.SiteUser) ([]*dal.UserUpdateResult, error)
	DeleteSiteUsers(siteId int, userIds []int) ([]*dal.UserUpdateResult, error)
	DeleteTidUsers(userIds []int) ([]*dal.UserUpdateResult, error)
	GetUsersForSite(siteId int) ([]entities.SiteUser, error)
	GetUsersForTid(tidId int) ([]entities.SiteUser, error)
	GetUserForId(userId int) (*entities.SiteUser, error)
}

func saveSiteUsers(d siteUsersDal, siteId int, siteUsers []*entities.SiteUser, deletedSiteUserIds []int) ([]*dal.UserUpdateResult, error) {
	return saveUserChanges(d, siteId, siteUsers, deletedSiteUserIds, nil, nil, true)
}

func saveUserChanges(
	d siteUsersDal,
	siteId int,
	usersToAddToSite []*entities.SiteUser,
	userIdsToDeleteFromSite []int,
	usersToAddToTids map[int][]*entities.SiteUser,
	tidUserIdsToDelete []int,
	validateAgainstTids bool) ([]*dal.UserUpdateResult, error) {
	saveChangeResults := make([]*dal.UserUpdateResult, 0)

	//Delete users from the site
	if len(userIdsToDeleteFromSite) > 0 {
		result, err := d.DeleteSiteUsers(siteId, userIdsToDeleteFromSite)
		if err != nil {
			return nil, err
		}
		for i := range result {
			saveChangeResults = append(saveChangeResults, result[i])
		}
	}

	if len(tidUserIdsToDelete) > 0 {
		result, err := d.DeleteTidUsers(tidUserIdsToDelete)
		if err != nil {
			return nil, err
		}
		for i := range result {
			saveChangeResults = append(saveChangeResults, result[i])
		}
	}

	failedValidationUsers, err := validatePedUsers(d, siteId, usersToAddToSite, usersToAddToTids, validateAgainstTids)
	for i := range failedValidationUsers {
		saveChangeResults = append(saveChangeResults, failedValidationUsers[i])
	}

	if err != nil {
		return nil, err
	}

	//checks if the user is valid
	isValid := func(validationResults []*dal.UserUpdateResult, user entities.SiteUser) bool {
		for _, validationResult := range validationResults {
			if validationResult.User.IsEqualTo(user) {
				return false
			}
		}
		return true
	}

	//Add users to the site
	if len(usersToAddToSite) > 0 {
		validUsers := make([]*entities.SiteUser, 0)
		for _, userToAddToSite := range usersToAddToSite {
			if isValid(failedValidationUsers, *userToAddToSite) {
				userToAddToSite.SiteId = siteId
				validUsers = append(validUsers, userToAddToSite)
			}
		}

		if len(validUsers) > 0 {
			result, err := d.AddOrUpdateSiteUsers(siteId, validUsers)
			if err != nil {
				return nil, err
			}
			for i := range result {
				saveChangeResults = append(saveChangeResults, result[i])
			}
		}
	}

	//For-each tid, add the users to it
	if len(usersToAddToTids) > 0 {
		for tid, tidUsers := range usersToAddToTids {
			validUsers := make([]*entities.SiteUser, 0)
			for _, userToAddToTid := range tidUsers {
				if isValid(failedValidationUsers, *userToAddToTid) {
					userToAddToTid.TidId = tid
					validUsers = append(validUsers, userToAddToTid)
				}
			}

			if len(validUsers) > 0 {
				results, err := d.AddOrUpdateTidUserOverride(tid, validUsers)
				if err != nil {
					return nil, err
				}
				for _, result := range results {
					saveChangeResults = append(saveChangeResults, result)
				}
			}
		}
	}

	return saveChangeResults, nil
}

func addValidationError(currentErrors []*dal.UserUpdateResult, newError *dal.UserUpdateResult) []*dal.UserUpdateResult {
	for _, u := range currentErrors {
		if u.User.IsEqualTo(newError.User) {
			return currentErrors
		}
	}
	return append(currentErrors, newError)

}

func validatePedUsers(d siteUsersDal, siteId int,
	siteUsers []*entities.SiteUser, tidsAndUsers map[int][]*entities.SiteUser, validateAgainstTids bool) ([]*dal.UserUpdateResult, error) {

	validationResults := make([]*dal.UserUpdateResult, 0)
	for _, invalidUser := range validateProfileUser(siteUsers) {
		validationResults = addValidationError(validationResults, invalidUser)
	}

	allTidUsersToAdd := make([]*entities.SiteUser, 0)
	for _, usersForTid := range tidsAndUsers {
		for _, tidUser := range usersForTid {
			allTidUsersToAdd = append(allTidUsersToAdd, tidUser)
		}
	}

	// Ensure that the users we're adding don't clash with existing site users
	usersForSite, err := d.GetUsersForSite(siteId)
	if err != nil {
		return nil, errors.New("could not validate site users: " + err.Error())
	}
	for _, userToAdd := range siteUsers {
		userValidationResult := dal.UserUpdateResult{User: *userToAdd}
		for _, existingUser := range usersForSite {
			if userToAdd.PIN == existingUser.PIN && userToAdd.Username != existingUser.Username && userToAdd.UserId != existingUser.UserId {
				userValidationResult.Result.SetError(fmt.Errorf("Cannot save site user '%v'; site user '%v' has the same PIN", userToAdd.Username, existingUser.Username))
				userValidationResult.User.SiteId = siteId
				validationResults = addValidationError(validationResults, &userValidationResult)
				continue
			}
		}

		for _, existingUser := range usersForSite {
			if userToAdd.PIN != existingUser.PIN && userToAdd.Username == existingUser.Username && userToAdd.UserId != existingUser.UserId {
				userValidationResult.Result.SetError(fmt.Errorf("Duplicate user"))
				userValidationResult.User.SiteId = siteId
				validationResults = addValidationError(validationResults, &userValidationResult)
				continue
			}
		}
	}

	// Ensure that users we're adding don't clash with existing tid users
	existingTidUsers, err := d.GetTidUsersForSite(siteId)
	if err != nil {
		return nil, errors.New("could not validate site users: " + err.Error())
	}
	for _, userToAdd := range siteUsers {
		userValidationResult := dal.UserUpdateResult{User: *userToAdd}
		for _, existingUser := range existingTidUsers {
			if userToAdd.PIN == existingUser.PIN && userToAdd.Username != existingUser.Username && userToAdd.UserId != existingUser.UserId {
				userValidationResult.Result.SetError(fmt.Errorf("Cannot save site user '%v'; TID user '%v' has the same PIN", userToAdd.Username, existingUser.Username))
				userValidationResult.User.SiteId = siteId
				validationResults = addValidationError(validationResults, &userValidationResult)
				continue
			}
		}
	}

	// Ensure that TID users we're adding don't clash with other users
	for TID, usersForTid := range tidsAndUsers {
		for i, tidUserToAdd := range usersForTid {
			userValidationResult := dal.UserUpdateResult{User: *tidUserToAdd}

			for j, otherTidUserToAdd := range usersForTid {
				if tidUserToAdd.Username == otherTidUserToAdd.Username && i != j {
					userValidationResult.Result.SetError(fmt.Errorf("Cannot save TID user '%v'; TID user '%v' has the same username on the same TID", tidUserToAdd.Username, otherTidUserToAdd.Username))
					userValidationResult.User.TidId = TID
					validationResults = addValidationError(validationResults, &userValidationResult)
					continue
				}
			}

			for _, existingUser := range usersForSite {
				if tidUserToAdd.PIN == existingUser.PIN && tidUserToAdd.Username != existingUser.Username && tidUserToAdd.UserId != existingUser.UserId {
					userValidationResult.Result.SetError(fmt.Errorf("Cannot save TID user '%v'; site user '%v' has the same PIN", tidUserToAdd.Username, existingUser.Username))
					userValidationResult.User.TidId = TID
					validationResults = addValidationError(validationResults, &userValidationResult)
					continue
				}

				if tidUserToAdd.Username == existingUser.Username && tidUserToAdd.PIN != existingUser.PIN {
					userValidationResult.Result.SetError(fmt.Errorf("Cannot save TID user '%v'; PIN is different to the user it is overriding", tidUserToAdd.Username))
					userValidationResult.User.TidId = TID
					validationResults = addValidationError(validationResults, &userValidationResult)
					continue
				}
			}

			for _, existingUser := range siteUsers {
				if tidUserToAdd.PIN == existingUser.PIN && tidUserToAdd.Username != existingUser.Username {
					userValidationResult.Result.SetError(fmt.Errorf("Cannot save TID user '%v'; site user '%v' has the same PIN", tidUserToAdd.Username, existingUser.Username))
					userValidationResult.User.TidId = TID
					validationResults = append(validationResults, &userValidationResult)
					continue
				}
			}
			if validateAgainstTids {
				for _, existingUser := range existingTidUsers {
					if tidUserToAdd.PIN == existingUser.PIN && tidUserToAdd.Username != existingUser.Username && tidUserToAdd.UserId != existingUser.UserId {
						userValidationResult.Result.SetError(fmt.Errorf("Cannot save TID user '%v'; TID user '%v' has the same PIN", tidUserToAdd.Username, existingUser.Username))
						userValidationResult.User.TidId = TID
						validationResults = addValidationError(validationResults, &userValidationResult)
						continue
					}
				}
			}

			for _, otherTidUserToAdd := range allTidUsersToAdd {
				if tidUserToAdd.PIN == otherTidUserToAdd.PIN && tidUserToAdd.Username != otherTidUserToAdd.Username {
					userValidationResult.Result.SetError(fmt.Errorf("Cannot save TID user '%v'; TID user '%v' has the same PIN", tidUserToAdd.Username, otherTidUserToAdd.Username))
					userValidationResult.User.TidId = TID
					validationResults = addValidationError(validationResults, &userValidationResult)
					continue
				}
			}
		}
	}

	return validationResults, nil
}

func contains(ints []int, x int) bool {
	for _, i := range ints {
		if i == x {
			return true
		}
	}

	return false
}

func clearTidUsers(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	r.ParseForm()
	tid := r.Form.Get("TID")
	tidId, err := strconv.Atoi(tid)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Error with request", http.StatusBadRequest)
		return
	}
	err = dal.ClearTidUsers(tidId)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Error saving users", http.StatusInternalServerError)
		return
	}
}

func saveTidUsersHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	r.ParseForm()
	tid := r.Form.Get("TID")
	siteId, err := strconv.Atoi(r.Form.Get("SiteId"))
	if err != nil {
		logging.Error("Failed to get site ID from formdata inside saveTidUsersHandler", err.Error())
		http.Error(w, saveTidUsersError, http.StatusBadRequest)
		return
	}
	tidId, err := strconv.Atoi(tid)
	if err != nil {
		logging.Error("Failed to get TID from formdata inside saveTidUsersHandler", err.Error())
		http.Error(w, saveTidUsersError, http.StatusBadRequest)
		return
	}

	users := r.Form.Get("Users")
	var userList = make([]*entities.SiteUser, 0)
	err = json.Unmarshal([]byte(users), &userList)
	if err != nil {
		logging.Error("Failed to get users from formdata inside saveTidUsersHandler", err.Error())
		http.Error(w, saveTidUsersError, http.StatusBadRequest)
		return
	}

	newUsers := r.Form.Get("NewUsers")
	var newUserList = make([]*entities.SiteUser, 0)

	err = json.Unmarshal([]byte(newUsers), &newUserList)
	if err != nil {
		logging.Error(err.Error())
		http.Error(w, "Error with request: "+err.Error(), http.StatusBadRequest)
		return
	}

	userList = append(userList, newUserList...)

	// checking special characters
	re, _ := regexp.Compile(`[^\w]`)
	for i, _ := range userList {
		userList[i].Username = strings.TrimSpace(userList[i].Username)
		found := re.MatchString(userList[i].Username)
		if found {
			http.Error(w, "Username must not contain special characters", http.StatusInternalServerError)
			return
		}
	}

	var profileTypeName = r.Form.Get("profileTypeName")

	deletedUsers := r.Form.Get("DeletedUsers")

	err = validateTidUsers(tidId, userList, siteId)
	if err != nil {
		switch err.(type) {
		case *siteUserValidationError:
			logging.Error(err.Error())
			http.Error(w, saveTidUsersError+": "+err.Error(), http.StatusBadRequest)
			return
		default:
			logging.Error(err.Error())
			http.Error(w, saveTidUsersError+": "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := processUsersOverride(profileTypeName, users, deletedUsers, newUsers, tmsUser, true, tid, tidId, siteId); err != nil {
		logging.Error("An error occurred while executing processUsersOverride", err.Error())
		http.Error(w, saveTidUsersError, http.StatusUnprocessableEntity)
		return
	}
}

func validateTidUsers(tidId int, tidUsers []*entities.SiteUser, siteId int) error {
	// TODO: Have the UI send through the site ID
	saveChangeResults := make([]*dal.UserUpdateResult, 0)

	existingUsers, err := dal.GetUsersForTid(tidId)
	if err != nil {
		return err
	}
	existingUserIds := make(map[int]bool)
	for _, existingUser := range existingUsers {
		existingUserIds[existingUser.UserId] = true
	}

	// TODO: Have the UI work out what TID users need to be deleted
	tidUsersToDelete := make(map[int]bool)
	for tidId, _ := range existingUserIds {
		tidUsersToDelete[tidId] = true
	}
	for _, tidUser := range tidUsers {
		if existingUserIds[tidUser.UserId] {
			tidUsersToDelete[tidUser.UserId] = false
		}
	}
	tidUserIdsToDelete := make([]int, 0)
	for tidId, deleteUser := range tidUsersToDelete {
		if deleteUser {
			tidUserIdsToDelete = append(tidUserIdsToDelete, tidId)
		}
	}

	tidsAndUsers := make(map[int][]*entities.SiteUser)
	tidsAndUsers[tidId] = tidUsers

	failedValidationUsers, err := validatePedUsers(&dal.SiteManagementDal{}, siteId, nil, tidsAndUsers, true)
	for i := range failedValidationUsers {
		saveChangeResults = append(saveChangeResults, failedValidationUsers[i])
	}

	if err != nil {
		return err
	}
	var errorString = ""

	for _, u := range failedValidationUsers {
		if !u.Result.Success {
			if errorString != "" {
				errorString += "<br>"
			}
			errorString += u.Result.ErrorMessage
		}
	}

	if errorString != "" {
		return errors.New(errorString)
	}

	return nil
}

func showTidUserModal(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	r.ParseForm()
	tid := r.Form.Get("TID")
	renderPartialTemplate(w, r, "TidUserModal", TidModalModel{Tid: tid}, tmsUser)
}

func generateRpiCertificate(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	mid := r.URL.Query().Get("MID")
	if strings.TrimSpace(mid) == "" {
		http.Error(w, "MID is not populated", http.StatusBadRequest)
		return
	}
	passwordBytes, err := base64.StdEncoding.DecodeString(r.URL.Query().Get("PASSWORD"))
	if err != nil {
		handleError(w, errors.New("An error occured during base64 decoding:"+err.Error()), tmsUser)
		return
	}
	password := string(passwordBytes)

	if RPIKeyFile == "" || RPICertFile == "" {
		http.Error(w, "Error generating certificate (No CA Certificate Found)", http.StatusInternalServerError)
		return
	}

	_ = os.RemoveAll(mid + "/")
	clientCert, err := TLSUtils.GenerateEncryptedClientCertificate(RPIKeyFile, RPICertFile, mid, password)
	if err != nil {
		http.Error(w, fmt.Sprint("Error generating certificate (Certificate Generation)", err), http.StatusInternalServerError)
		return
	}

	caCert, err := ioutil.ReadFile(RPICertFile)
	if err != nil {
		logging.Error("An error occured during reading the RPICertFile:" + err.Error())
		http.Error(w, "An error occured while reading the file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("fileName", stringHelpers.RemoveCRLFChars(mid+"_certs.zip"))
	wZip := zip.NewWriter(w)
	defer wZip.Close()

	clientCertFile, err := wZip.Create("Server0Q.pfx")
	if err != nil {
		http.Error(w, "Error generating certificate (Zipping Certificate)", http.StatusInternalServerError)
		return
	}
	_, err = clientCertFile.Write(clientCert)
	if err != nil {
		http.Error(w, "Error generating certificate (Zipping Certificate)", http.StatusInternalServerError)
		return
	}

	caCertFile, err := wZip.Create("Server0QRoot.cer")
	if err != nil {
		http.Error(w, "Error generating certificate (Zipping Certificate)", http.StatusInternalServerError)
		return
	}
	_, err = caCertFile.Write(caCert)
	if err != nil {
		http.Error(w, "Error generating certificate (Zipping Certificate)", http.StatusInternalServerError)
		return
	}
}

func GetCashbackEditModal(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	if err := r.ParseForm(); err != nil {
		logging.Warning(err.Error())
		http.Error(w, genericServerError, http.StatusInternalServerError)
		return
	}

	renderPartialTemplate(w, r, "cashback", nil, tmsUser)
}

// Generate payment service modal with configured services profile_data
func getPaymentServicesEditModal(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	if err := r.ParseForm(); err != nil {
		logging.Warning(err.Error())
		http.Error(w, genericServerError, http.StatusInternalServerError)
		return
	}
	validator := validation.New(dal.NewValidationDal())
	siteId, err := strconv.Atoi(r.Form.Get("siteId"))
	if err != nil {
		handleError(w, errors.New(siteConversionFailed+err.Error()), tmsUser)
		return
	}
	tidInt, err := strconv.Atoi(r.Form.Get("tid"))
	if err != nil {
		handleError(w, errors.New(tidConversionFailed+err.Error()), tmsUser)
		return
	}
	tid := dal.GetPaddedTidId(tidInt)

	_, err = validator.ValidateTidFormat(tid)
	if err != nil {
		logging.Warning(err.Error())
		http.Error(w, genericServerError, http.StatusInternalServerError)
		return
	}

	//Fetch payment group with it services
	paymentServicesMap, err := dal.GetSitePaymentServices(siteId)
	if err != nil {
		logging.Warning(err.Error())
		http.Error(w, genericServerError, http.StatusInternalServerError)
		return
	}

	//Merge payment group services with already configured services in profile data
	paymentServices := MergePaymentsWithConfigured(paymentServicesMap, tid, siteId)
	serviceJson, err := json.Marshal(paymentServices)

	dict := make(map[string]interface{})
	if err != nil {
		dict["servicesJson"] = "[]"
	} else {
		dict["servicesJson"] = string(serviceJson)
	}
	renderPartialTemplate(w, r, "tidServicesMidTidModal", dict, tmsUser)
}

// Merge empty set of payment services of payment group with already configured profile data of services data_element
func MergePaymentsWithConfigured(paymentServicesMap map[int]*dal.PaymentService, tid string, siteId int) []*dal.PaymentService {
	dataValue := ""
	var err error
	intTid, err := strconv.Atoi(tid)
	if intTid > 0 && err == nil {
		dataValue, err = dal.GetDataValueForSiteTidAndElementName(intTid, siteId, "paymentServicesConfigs")
		if err != nil {
			logging.Warning(err.Error())
		}
	}

	if dataValue != "" {
		var configuredServices []dal.PaymentService
		err := json.Unmarshal([]byte(dataValue), &configuredServices)
		if err == nil {
			//remove item from current services if the item deleted
			for _, configuredService := range configuredServices {
				service, exists := paymentServicesMap[configuredService.ServiceId]
				if exists {
					service.MID = configuredService.MID
					service.TID = configuredService.TID
				}
			}
		}
	}

	var servicesMerged []*dal.PaymentService
	for _, service := range paymentServicesMap {
		servicesMerged = append(servicesMerged, service)
	}
	return servicesMerged
}

func GetGratuityEditModal(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	if err := r.ParseForm(); err != nil {
		logging.Warning(err.Error())
		http.Error(w, genericServerError, http.StatusInternalServerError)
		return
	}

	gd := r.Form.Get("data")
	var gratuity GratuityDetails
	json.Unmarshal([]byte(gd), &gratuity)

	dict := make(map[string]interface{})
	dict["model"] = gratuity
	renderPartialTemplate(w, r, "gratuity", dict, tmsUser)
}

func GetDpoMomoEditModal(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	if err := r.ParseForm(); err != nil {
		logging.Warning(err.Error())
		http.Error(w, genericServerError, http.StatusInternalServerError)
		return
	}

	renderPartialTemplate(w, r, "dpo", nil, tmsUser)
}

func GetSoftUIEditModal(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	if err := r.ParseForm(); err != nil {
		logging.Warning(err.Error())
		http.Error(w, genericServerError, http.StatusInternalServerError)
		return
	}

	renderPartialTemplate(w, r, "softUI", nil, tmsUser)
}

func neuterDirectoryBrowsing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Path) == 0 || strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}

		addHeaderSecurityItems(w, r)
		next.ServeHTTP(w, r)
	})
}

// DiagnosticStatus : hop down through service structure.
func (t *DiagnosticLayerRPCserver) DiagnosticStatus(cxt context.Context, request *txn.DiagnosticsRequest) (*txn.DiagnosticsResponse, error) {
	defer recoveryProcess()
	var err error

	// INTENTIONALLY DOES NOT CHECK CERTIFICATE HERE
	var maxHops int32 = 1

	var clients []rpcHelp.GRPCclient
	for _, client := range diagnosticGRPCclients {
		clients = append(clients, *client)
	}

	pbDiagResp := new(pbDiag.DiagnosticsResponse)
	err = pbPrintDiags.GetDiags(maxHops, pbDiagResp, clients, nil, logging, ApplicationName, listeningPort, SERVERLOGGER)

	if err != nil {
		fmt.Println("Error during GetDiags ", err)
	}

	returnval := pbPrintDiags.TransformHop(pbDiagResp)
	returnval.Hops = append(returnval.Hops, getMySqlResponse())
	returnval.Hops = append(returnval.Hops, getMongoResponse())

	return returnval, nil
}

// Returns a Diagnostic Response to check connectivity to MySQL
func getMySqlResponse() *txn.DiagnosticsResponse {
	mysqlResponse := txn.DiagnosticsResponse{}
	mysqlResponse.Name = "MySQL"

	_, err := dal.GetDB()
	if err != nil {
		mysqlResponse.Errors = append(mysqlResponse.Errors, err.Error())
	}

	return &mysqlResponse
}

// Returns a Diagnostic Response to check connectivity to Mongo
func getMongoResponse() *txn.DiagnosticsResponse {
	mongoResponse := txn.DiagnosticsResponse{}
	mongoResponse.Name = "MongoDB"

	_, err := dal.GetMongoClient()
	if err != nil {
		mongoResponse.Errors = append(mongoResponse.Errors, err.Error())
	}

	return &mongoResponse
}

func recoveryProcess() {
	// Go 1.18 does not allow nil comparison to type "any" and flags this as a compiler error.
	// As we're still using Go 1.14, this won't cause builds to fail, but will show as an error in the IDE
	// This type assertion of nil fixes this
	if r := recover(); r != interface{}(nil) {
		logging.Panic(r)
	}
}

func webHandlers(h UserHandleFunction, perm UserPermission) http.Handler {
	return gziphandler.GzipHandler(logHandler(authHandler(h, perm)))
}

func handleFilePath(r *mux.Router, prefix string, directory string) {
	r.PathPrefix(prefix).Handler(gziphandler.GzipHandler(http.StripPrefix(prefix, neuterDirectoryBrowsing(http.FileServer(http.Dir(directory))))))
}

func SetupDiagnosticLayerRPCServer() {
	// Setup a port for INTERNAL USE ONLY.

	var internalPort = ":" + cfg.GetString("DiagnosticsPort", "8911")

	logger.GetLogger().Information("listening on Diagnostics port ", internalPort)

	internalOnlyl, e := net.Listen("tcp", internalPort)
	if e != nil {
		logger.GetLogger().Error("listen error:", e)
	}
	internalOnlys := grpc.NewServer()
	txn.RegisterDiagnosticLayerRPCServer(internalOnlys, &DiagnosticLayerRPCserver{})
	// Register reflection service on gRPC server.
	reflection.Register(internalOnlys)
	if err := internalOnlys.Serve(internalOnlyl); err != nil {
		logger.GetLogger().Error("Failed to serve: ", err)
	}
}

func main() {
	// Look for a config override flag
	override := flag.String("config", "TMS", "Config Override file")
	flag.Parse()

	err := cfg.InitialRead(os.Args[0], *override)

	crypt.SetupCrypt(cfg.GetKey())

	if err != nil {
		go log.Fatal("Configuration error:", err)
	}

	err = parseConfig()
	if err != nil {
		go log.Fatal("Error parsing configuration: ", err)
	}
	config.ParseConfig()

	pbPrintDiags.Start(Version, ApplicationName, Build)

	logger.Init(ApplicationName)
	logging = logger.GetLogger()

	if err != nil {
		// do not prevent service running because no logger defined.
		log.Print(err)
	} else {
		logging.Information(fmt.Sprintf("Startup. Version : %s Build Date : %s", Version, Build))
		logging.Information(cfg.PrintEnvironment())
		logging.Information(fmt.Sprintf("WebSite. Version : %s DB Version : %s", WebsiteVersion, strconv.Itoa(DbVersion)))
	}

	if err != nil {
		// do not prevent service running because no longer defined.
		log.Print(err)
	}

	dal.Connect(logging)
	defer dal.CloseDB()

	dal.ConnectToMongo()
	defer dal.CloseMongo()
	dal.InitConstants()

	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_AES_128_GCM_SHA256,
		},
	}

	if *override == "API" {

		go SetupDiagnosticLayerRPCServer()

		server := &http.Server{
			Addr:         ":" + listeningPort,
			Handler:      routers.CreateAutomationAPIRouter(),
			ReadTimeout:  10 * time.Minute,
			WriteTimeout: 10 * time.Minute,
			TLSConfig:    cfg,
			TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
		} // setting read/write timeouts to avoid reaching file limit, see: https://github.com/golang/go/issues/28272#issuecomment-478728262

		logging.Information("Listening on Automation API port ", listeningPort)

		if err := server.ListenAndServeTLS(HTTPSCertFile, HTTPSKeyFile); err != nil {
			logging.Error(err.Error())
		}

		return
	}

	socketConnections = make(map[string]*websocket.Conn, 0)

	GRPCclients = rpcHelp.SetupGRPCClientsFromArray(getRpcConfig(), logging)
	// As of now automation api is not having any downstream Services, so keeping this here only
	diagnosticGRPCclients = rpcHelp.SetupGRPCClientsFromArray(getDiagnosticRpcConfig(), logging)

	token := GenerateCSRFToken()
	CSRF := csrf.Protect(token, csrf.Secure(true), csrf.FieldName("csrfmiddlewaretoken"), csrf.HttpOnly(true))

	r := createTMSWebsiteRouter()

	dal.SetDbEncryptVersion(DbEncryptVersion)
	// Do a version check and perform any neccessary upgrades
	dal.CheckDatabaseVersion(DbVersion, EncryptEverything)

	// Monitor for and apply any ancillary scripts every 30 seconds
	ticker := time.NewTicker(30 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				dal.CheckAncillaryScripts()
			}
		}
	}()

	// Clear tokens after the token expires.
	clearTokensTicker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case <-clearTokensTicker.C:
				currentTime := time.Now()

				for token, user := range userSessionTokenMap {
					if currentTime.After(user.Expires) {
						redirectExpiredUser(user)
						userSessionTokenMapMutex.Lock()
						delete(userSessionTokenMap, token)
						userSessionTokenMapMutex.Unlock()
					}
				}
			}
		}
	}()

	// Listen for user management files being uploaded - these should be processed automatically
	go listenForUserManagementFiles()

	go SetupDiagnosticLayerRPCServer()

	http.Handle("/echo", webHandlers(createSessionHandlerWebsocket, None))

	go func() {
		server := &http.Server{
			Addr:         ":5006",
			ReadTimeout:  10 * time.Minute,
			WriteTimeout: 10 * time.Minute,
			TLSConfig:    cfg,
			TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
		}
		if err := server.ListenAndServeTLS(HTTPSCertFile, HTTPSKeyFile); err != nil {
			logging.Error(err.Error())
		}
	}()

	logging.Information("Listening on Web TMS port ", listeningPort)
	server := &http.Server{
		Addr:         ":" + listeningPort,
		Handler:      CSRF(r),
		ReadTimeout:  10 * time.Minute,
		WriteTimeout: 10 * time.Minute,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	} // setting read/write timeouts to avoid reaching file limit, see: https://github.com/golang/go/issues/28272#issuecomment-478728262
	if err := server.ListenAndServeTLS(HTTPSCertFile, HTTPSKeyFile); err != nil {
		logging.Error(err.Error())
	}
}

func getDpoMomoFieldsDataHandler(w http.ResponseWriter, r *http.Request, tmsUser *entities.TMSUser) {
	logging.Debug("@getDpoMomoFieldsDataHandler, attempting to fetch DPO Momo Data")
	res, err := dal.GetDpoMomoFieldsData()
	if err != nil {
		http.Error(w, responseError, http.StatusBadRequest)
		return
	}
	w.Write([]byte(res))
}

func populateAllConfigFiles(p ProfileMaintenanceModel) {
	populateConfigFiles(p, "softUI", "softUIConfiguration", fileServer.NewFsReader(config.FileserverURL).GetAllSoftUIConfigFiles)
	populateConfigFiles(p, "mainMenu", "mainMenuConfiguration", fileServer.NewFsReader(config.FileserverURL).GetAllMenuFiles)
	populateConfigFiles(p, "receipt", "receiptConfiguration", fileServer.NewFsReader(config.FileserverURL).GetAllReceiptConfigFiles)
}

func populateConfigFiles(p ProfileMaintenanceModel, configType, dataElementName string, getConfigFilesFunc func() ([]string, error)) {
	// Create a single collection of all data groups and append individual groups
	allGroups := [][]*dal.DataGroup{
		p.ProfileGroups,
		p.GlobalGroups,
		p.AcquirerGroups,
		p.GlobalGroups,
	}

	// Get all JSON files based on the specified config type
	fileNames, err := getConfigFilesFunc()
	if err != nil {
		logging.Error(err)
		return
	}

	// Loop through all the groups i - each group
	for i := range allGroups {
		// Loop through the different data groups within each i group j - each data group
		for j := range allGroups[i] {
			// Loop through the different data elements within the j data group x - each data element
			for x := range allGroups[i][j].DataElements {
				if allGroups[i][j].DataElements[x].Name == dataElementName {
					options := make([]dal.OptionData, 0)

					// Check file names returns values otherwise populate dropdown with no data
					if fileNames == nil {
						options = append(options, dal.OptionData{Option: "__NO_DATA__"})
						allGroups[i][j].DataElements[x].Options = options
					}

					// Loop through each of the file names returned and if not empty append to the data element x options
					for _, name := range fileNames {
						if name != "" {
							if strings.Contains(allGroups[i][j].DataElements[x].DataValue, name) {
								options = append(options, dal.OptionData{Option: name, Selected: true})
							} else {
								options = append(options, dal.OptionData{Option: name})
							}
							allGroups[i][j].DataElements[x].Options = options
						}
					}
				}
			}
		}
	}
}
