package dal

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"nextgen-tms-website/crypt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const dbVersionAddWifiPin = 71

var (
	dbEncryptVersion int
)

func SetDbEncryptVersion(s int) {
	dbEncryptVersion = s
}

// CheckDatabaseVersion Check the current database version and see if we need to upgrade
func CheckDatabaseVersion(databaseVersion int, encryptCallback func()) {
	db, err := GetDB()
	if err != nil {
		logging.Information("Upgrade Failed")
		return
	}

	currentVersion := 0
	rows, err := db.Query("SELECT version FROM db_version")
	if err != nil {
		// Unversioned database
		onUpgrade(db, currentVersion, databaseVersion, encryptCallback)
	} else {
		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&currentVersion)
		}
		onUpgrade(db, currentVersion, databaseVersion, encryptCallback)
	}
}

func GetDbVersion() int {
	db, err := GetDB()
	if err != nil {
		return 0
	}

	currentVersion := 0
	rows, err := db.Query("SELECT version FROM db_version")
	if err != nil {
		return currentVersion
	} else {
		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&currentVersion)
		}
	}
	return currentVersion
}

func onUpgrade(db *sql.DB, oldVersion int, newVersion int, encryptCallback func()) {
	if oldVersion < newVersion {
		err := DumpDB()
		if err != nil {
			logging.Error(err)
		}

		for oldVersion < newVersion {
			logging.Information("onUpgrade oldVersion: " + strconv.Itoa(oldVersion) + " newVersion: " + strconv.Itoa(newVersion))
			oldVersion++

			// Update the db tables
			safePath := filepath.Join("TMSUpdate", "UpgradeV")
			err = updateDb(safePath + strconv.Itoa(oldVersion))
			if err != nil {
				logging.Error("Upgrade DB ", err)
				logging.Information("Upgrade Failed")
				return
			}

			switch oldVersion {
			case dbEncryptVersion:
				// WARNING SPECIAL CASE FOR HANDLING THE FULL ENCRYPT STEP
				encryptCallback()
			case dbVersionAddWifiPin:
				if err = populateWifiPINs(db); err != nil {
					logging.Error("Upgrade DB ", err)
					logging.Information("Upgrade Failed")
					return
				}
			}

			_, err = db.Exec("Update db_version SET version = ? Where id = 0", oldVersion)

			if err != nil {
				logging.Error("Upgrade DB ", err)
				logging.Information("Upgrade Failed")
				return
			}
		}

		logging.Information("Upgrade Completed Successfully")
	}
}

func populateWifiPINs(db *sql.DB) error {
	rows, err := db.Query(`SELECT JSON_ARRAYAGG(profile_id) FROM profile
		WHERE profile_type_id=(SELECT profile_type_id FROM profile_type WHERE name='site')`)
	if err != nil {
		return err
	}
	defer rows.Close()

	var jsonResult sql.NullString
	if rows.Next() {
		err := rows.Scan(&jsonResult)
		if err != nil {
			return err
		}
	}

	if !jsonResult.Valid {
		// no sites to process
		return nil
	}

	var siteProfileIDs []int
	if err := json.Unmarshal([]byte(jsonResult.String), &siteProfileIDs); err != nil {
		return err
	}

	for _, siteID := range siteProfileIDs {
		tidProfiles, err := getPEDProfilesForSite(db, siteID)
		if err != nil {
			return err
		}

		rawPin, err := generate5DigitPin()
		if err != nil {
			return err
		}
		pin := crypt.Encrypt(rawPin)

		profileIDs := append(tidProfiles, siteID)
		for _, profileID := range profileIDs {
			_, err = db.Exec(
				`INSERT IGNORE into profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden, is_encrypted) 
			values (?, (SELECT data_element_id FROM data_element WHERE name = 'wifiPIN'), ?, 1, NOW(), 'system', NOW(), 'system', 1, 0, 1);`,
				profileID, pin)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func getPEDProfilesForSite(db *sql.DB, siteID int) ([]int, error) {
	rows, err := db.Query(`
		SELECT JSON_ARRAYAGG(tp.profile_id)
		FROM profile AS tp
		INNER JOIN profile_type tpt ON tpt.profile_type_id = tp.profile_type_id
		INNER JOIN tid_site AS ts ON tp.name=ts.tid_id
		WHERE tpt.name = 'tid'
		AND ts.site_id=(SELECT site_id FROM site_profiles WHERE profile_id=?)`, siteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jsonResult sql.NullString
	if rows.Next() {
		err := rows.Scan(&jsonResult)
		if err != nil {
			return nil, err
		}
	}
	if !jsonResult.Valid {
		// no tid profiles for site
		return nil, nil
	}

	var tidProfiles []int
	if err := json.Unmarshal([]byte(jsonResult.String), &tidProfiles); err != nil {
		return nil, err
	}
	return tidProfiles, nil
}

func generate5DigitPin() (string, error) {
	PIN, err := rand.Int(rand.Reader, big.NewInt(99999))
	if err != nil {
		return "", err
	}
	PINstr := strconv.FormatInt(PIN.Int64(), 10)
	leads := 5 - len(PINstr)
	for leads > 0 {
		PINstr = "0" + PINstr
		leads--
	}
	return PINstr, nil
}

func DumpDB() error {
	// Split the connection string into components for a db command
	fallbackBackupLocation := "backups"
	t := strings.Split(connectionString, "@")
	userPassword := strings.Split(t[0], ":")
	user := userPassword[0]
	password := userPassword[1]

	ipPort := strings.Split(t[1], "/")
	ipPort = strings.Split(ipPort[0], ":")
	ipPort[0] = strings.Replace(ipPort[0], "tcp(", "", -1)
	ipPort[1] = strings.Replace(ipPort[1], ")", "", -1)
	ip := ipPort[0]
	port := ipPort[1]

	backupPath, dirError := ValidateDirectory(backupLocation, true, fallbackBackupLocation)
	if dirError != nil {
		message := fmt.Sprintf("Unable to create Backup directory: %s", dirError.Error())
		return errors.New(message)
	}

	cmd := exec.Command("mysqldump", "-P"+port, "-h"+ip, "-u"+user, "-p"+password, "--routines", "NextGen_TMS")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logging.Error(err)
	}

	outfile, err := os.Create(backupPath + "/nextgen_tms " + time.Now().Format("02-01-2006-15-04-05") + ".sql")
	if err != nil {
		logging.Error(err)
		// File create fail re-attempt using fallback directory
		backupPath, dirError = ValidateDirectory(fallbackBackupLocation, true, "")
		if dirError != nil {
			message := fmt.Sprintf("Unable to create Backup directory: %s", dirError.Error())
			return errors.New(message)
		}

		outfile, err = os.Create(backupPath + "/nextgen_tms " + time.Now().Format("02-01-2006-15-04-05") + ".sql")
		if err != nil {
			return err
		}
	}
	defer outfile.Close()

	// start the command after having set up the pipe
	if err := cmd.Start(); err != nil {
		logging.Error(err)
	}

	// read command's stdout line by line
	in := bufio.NewWriter(outfile)
	defer in.Flush()

	io.Copy(outfile, stdout)

	logging.Information(fmt.Sprintf("Backup Created: '%s'", outfile.Name()))
	return nil
}

func DumpAndDownloadDB(w http.ResponseWriter) {
	// Split the connection string into components for a db command
	t := strings.Split(connectionString, "@")
	userPassword := strings.Split(t[0], ":")
	user := userPassword[0]
	password := userPassword[1]

	ipPort := strings.Split(t[1], "/")
	ipPort = strings.Split(ipPort[0], ":")
	ipPort[0] = strings.Replace(ipPort[0], "tcp(", "", -1)
	ipPort[1] = strings.Replace(ipPort[1], ")", "", -1)
	ip := ipPort[0]
	port := ipPort[1]

	cmd := exec.Command("mysqldump", "-P"+port, "-h"+ip, "-u"+user, "-p"+password, "--routines", "NextGen_TMS")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logging.Error(err)
	}

	// start the command after having set up the pipe
	if err := cmd.Start(); err != nil {
		logging.Error(err)
	}

	io.Copy(w, stdout)
}

func updateDb(directory string) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(directory)
	if err != nil {
		logging.Error(err)
		return err
	}

	for _, f := range files {
		logging.Information("Applying update: " + directory + "/" + f.Name())

		file, err := os.Open(directory + "/" + f.Name())
		if err != nil {
			logging.Error(err)
			return err
		}

		fileinfo, err := file.Stat()
		if err != nil {
			file.Close()
			logging.Error(err)
			return err
		}

		filesize := fileinfo.Size()
		buffer := make([]byte, filesize)

		_, err = file.Read(buffer)
		if err != nil {
			file.Close()
			logging.Error(err)
			return err
		}

		t := strings.Split(string(buffer), "\n")
		// Check for procedure otherwise execute commands on a per line basis
		if strings.Contains(t[0], "--multiline") {
			err := parseMultilineCommand(t, db)
			if err != nil {
				file.Close()
				logging.Error(err)
				return err
			}
		} else {
			for _, f := range t {
				if f != "" {
					_, err := db.Exec(f)
					if err != nil {
						file.Close()
						logging.Error(err)
						return err
					}
				}
			}
		}
		file.Close()
	}

	return nil
}

// For a Procedure or multi line command parse the entire file as a single command
// NB. procedure upgrade files should contain a single procedure
func parseMultilineCommand(commands []string, db *sql.DB) error {
	var buffer bytes.Buffer
	for _, f := range commands {
		if f != "" && !strings.Contains(f, "--multiline") {
			_, err := buffer.WriteString(f + "\n")
			if err != nil {
				return err
			}
		}
	}
	_, err := db.Exec(buffer.String())
	if err != nil {
		return err
	}
	return nil
}

// CheckAncillaryScripts searches in ancillary folder and applies any additional scripts
func CheckAncillaryScripts() {
	safeSrcPath := filepath.Join("TMSUpdate", "PendingAncillaryScripts")
	safeDestPath := filepath.Join("ProcessedAncillaryScripts")
	safeFailedPath := filepath.Join("FailedAncillaryScripts")

	// Create directories if they don't already exist
	if _, err := os.Stat(safeSrcPath); os.IsNotExist(err) {
		os.Mkdir(safeSrcPath, 0777)
	}
	if _, err := os.Stat(safeDestPath); os.IsNotExist(err) {
		os.Mkdir(safeDestPath, 0777)
	}
	if _, err := os.Stat(safeFailedPath); os.IsNotExist(err) {
		os.Mkdir(safeFailedPath, 0777)
	}

	files, err := ioutil.ReadDir(safeSrcPath)
	if err != nil {
		logging.Error(err)
		return
	}

	if len(files) != 0 {
		logging.Information("Ancillary scripts found, processing updates...")
		applyAncillaryScripts(safeSrcPath, safeDestPath, safeFailedPath, files)
		logging.Information("Finished processing ancillary scripts")
	}
}

func applyAncillaryScripts(src string, dest string, fail string, files []os.FileInfo) {
	db, err := ConnectAncillaryScripts(logging)
	if err != nil {
		return
	}

	// backup before applying any ancillary scripts incase something goes wrong
	err = DumpDB()
	if err != nil {
		logging.Error(err)
		return
	}

	// loop the files and apply all the scripts
	for _, f := range files {
		logging.Information("Applying update: " + f.Name())

		file, err := os.Open(src + "/" + f.Name())
		if err != nil {
			logging.Error(err)
			return
		}

		fileinfo, err := file.Stat()
		if err != nil {
			file.Close()
			logging.Error(err)
			return
		}

		filesize := fileinfo.Size()
		buffer := make([]byte, filesize)

		_, err = file.Read(buffer)
		if err != nil {
			file.Close()
			logging.Error(err)
			return
		}

		lines := strings.Split(string(buffer), "\n")
		// Check for procedure otherwise execute commands on a per line basis
		var failed = false
		var sqlError error
		if strings.Contains(lines[0], "--multiline") {
			sqlError = parseMultilineCommand(lines, db)
			if sqlError != nil {
				logging.Error(sqlError)
				copyAndDeleteFile(src, fail, f.Name(), buffer, file)
				logging.Information("Failed: " + f.Name() + " has been moved to " + fail)
				failed = true
				ancillaryAudit(f.Name(), failed, sqlError, db)
				continue
			}
		} else {
			for _, line := range lines {
				if line != "" {
					_, sqlError = db.Exec(line)
					if sqlError != nil {
						logging.Error(sqlError)
						copyAndDeleteFile(src, fail, f.Name(), buffer, file)
						logging.Information("Failed: " + f.Name() + " has been moved to " + fail)
						failed = true
						ancillaryAudit(f.Name(), failed, sqlError, db)
						break
					}
				}
			}
		}

		if !failed {
			copyAndDeleteFile(src, dest, f.Name(), buffer, file)
			logging.Information("Success: The script has been moved to " + dest)
			ancillaryAudit(f.Name(), failed, sqlError, db)
			file.Close()
		}
	}

	// Close the connection pool after running the feature scripts
	err = db.Close()
	if err != nil {
		logging.Error(err)
	}
}

//Records to DB when an ancillary script is ran, both failed and successful
func ancillaryAudit(name string, status bool, reason error, db *sql.DB) {

	dt := time.Now()
	var executionStatus int
	if status {
		executionStatus = 0
	} else {
		executionStatus = 1
	}
	var failReason string
	if reason != nil {
		failReason = reason.Error()
	}
	rows, err := db.Query("CALL add_feature_script_audit(?,?,?,?,?)", name, dt, dt, executionStatus, failReason)
	if err != nil {
		logging.Error(err)
	}
	defer rows.Close()
}

func copyAndDeleteFile(src string, dest string, fileName string, buffer []byte, file *os.File) {
	outfile, err := os.Create(dest + "/" + fileName)
	if err != nil {
		logging.Error(err)
		return
	}

	defer outfile.Close()
	_, err = outfile.Write(buffer)
	if err != nil {
		logging.Error(err)
		return
	}

	// Close and delete the src file
	file.Close()
	err = os.Remove(src + "/" + fileName)
	if err != nil {
		logging.Error(err)
		return
	}
}

func ValidateDirectory(path string, createIfNotExist bool, fallback string) (string, error) {

	useFallback := false

	fullPath, absError := filepath.Abs(path)
	if absError != nil {
		logging.Error("Unable to get directory path")
		useFallback = true
	} else {
		info, err := os.Stat(fullPath)
		if os.IsNotExist(err) {
			if createIfNotExist {
				dirCreateError := os.MkdirAll(fullPath, 0777)
				if dirCreateError != nil {
					logging.Error(fmt.Sprintf("Unable to create directory: '%s', Error: '%s'", backupLocation, dirCreateError.Error()))
					useFallback = true
				}
			} else {
				useFallback = true
			}

		} else if !info.IsDir() {
			logging.Error("Path exists but is not a directory")
			useFallback = true
		}
	}

	if useFallback {

		if fallback != "" {
			logging.Error(fmt.Sprintf("Falling back to default: '%s'", fallback))
			return ValidateDirectory(fallback, true, "")
		}

		// no fallback specified, return error
		return "", errors.New("directory invalid")
	}

	absolutePath, absError := filepath.Abs(path)
	if absError != nil {
		logging.Error("Unable to get directory path")
	}

	return absolutePath, nil
}
