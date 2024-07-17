package fileServer

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
)

func NewFsReader(fsAddress string) FsReader {
	return &fileServer{address: fsAddress}
}

type FsReader interface {
	GetAllFilesByPattern(pattern string) ([]string, error)
	GetAllFiles() ([]string, error)
	GetAllCustomerReceiptFiles() ([]string, error)
	GetAllMerchantReceiptFiles() ([]string, error)
	GetAllMenuFiles() ([]string, error)
	GetAllSoftUIConfigFiles() ([]string, error)
	GetAllReceiptConfigFiles() ([]string, error)
	GetFile(name string, directory string) ([]byte, error)
	MoveFile(oldPath, newPath string) error
	GetAllMnoLogoFiles() ([]string, error)
}

type fileServer struct {
	address string
}

func (fileServer *fileServer) GetAllFilesByPattern(pattern string) ([]string, error) {
	files, err := fileServer.GetAllFiles()
	if err != nil {
		return nil, err
	}
	filteredFiles := make([]string, 0)
	for _, file := range files {
		match, err := regexp.MatchString(pattern, file)
		if err != nil {
			return nil, err
		}
		if match {
			filteredFiles = append(filteredFiles, file)
		}
	}
	return filteredFiles, nil
}

func (fileServer *fileServer) GetAllFiles() ([]string, error) {
	result, err := http.Get(fileServer.address + "/getFileList")
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	if result.StatusCode != http.StatusOK {
		return nil, err
	}

	type fileListEntry struct {
		Name string
	}
	files := make([]fileListEntry, 0)
	err = json.NewDecoder(result.Body).Decode(&files)
	if err != nil {
		return nil, err
	}

	returnFiles := make([]string, 0)
	for _, file := range files {
		returnFiles = append(returnFiles, file.Name)
	}
	return returnFiles, nil
}

func (fileServer *fileServer) GetAllMnoLogoFiles() ([]string, error) {
	result, err := http.Get(fileServer.address + "/getFileList?directory=mnoLogo")
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	if result.StatusCode != http.StatusOK {
		return nil, err
	}

	type fileListEntry struct {
		Name string
	}
	files := make([]fileListEntry, 0)
	err = json.NewDecoder(result.Body).Decode(&files)
	if err != nil {
		return nil, err
	}

	returnFiles := make([]string, 0)
	for _, file := range files {
		returnFiles = append(returnFiles, file.Name)
	}
	return returnFiles, nil
}

func (fileServer *fileServer) GetAllSoftUIConfigFiles() ([]string, error) {
	result, err := http.Get(fileServer.address + "/getSoftUIConfigFileList")
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	if result.StatusCode != http.StatusOK {
		return nil, err
	}

	type fileListEntry struct {
		Name string
	}
	files := make([]fileListEntry, 0)
	err = json.NewDecoder(result.Body).Decode(&files)
	if err != nil {
		return nil, err
	}

	returnFiles := make([]string, 0)
	for _, file := range files {
		returnFiles = append(returnFiles, file.Name)
	}
	return returnFiles, nil
}

func (fileServer *fileServer) GetAllMenuFiles() ([]string, error) {
	result, err := http.Get(fileServer.address + "/getMenuFileList")
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	if result.StatusCode != http.StatusOK {
		return nil, err
	}

	type fileListEntry struct {
		Name string
	}
	files := make([]fileListEntry, 0)
	err = json.NewDecoder(result.Body).Decode(&files)
	if err != nil {
		return nil, err
	}

	returnFiles := make([]string, 0)
	for _, file := range files {
		returnFiles = append(returnFiles, file.Name)
	}
	return returnFiles, nil
}

func (fileServer fileServer) GetAllReceiptConfigFiles() ([]string, error) {
	result, err := http.Get(fileServer.address + "/getReceiptConfigFileList")
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	if result.StatusCode != http.StatusOK {
		return nil, err
	}

	type fileListEntry struct {
		Name string
	}
	files := make([]fileListEntry, 0)
	err = json.NewDecoder(result.Body).Decode(&files)
	if err != nil {
		return nil, err
	}

	returnFiles := make([]string, 0)
	for _, file := range files {
		returnFiles = append(returnFiles, file.Name)
	}
	return returnFiles, nil
}

func (fileServer *fileServer) GetAllMerchantReceiptFiles() ([]string, error) {
	result, err := http.Get(fileServer.address + "/getMerchantReceiptFileList")
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	if result.StatusCode != http.StatusOK {
		return nil, err
	}

	type fileListEntry struct {
		Name string
	}
	files := make([]fileListEntry, 0)
	err = json.NewDecoder(result.Body).Decode(&files)
	if err != nil {
		return nil, err
	}

	returnFiles := make([]string, 0)
	for _, file := range files {
		returnFiles = append(returnFiles, file.Name)
	}
	return returnFiles, nil
}
func (fileServer *fileServer) GetAllCustomerReceiptFiles() ([]string, error) {
	result, err := http.Get(fileServer.address + "/getCustomerReceiptFileList")
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()
	if result.StatusCode != http.StatusOK {
		return nil, err
	}

	type fileListEntry struct {
		Name string
	}
	files := make([]fileListEntry, 0)
	err = json.NewDecoder(result.Body).Decode(&files)
	if err != nil {
		return nil, err
	}

	returnFiles := make([]string, 0)
	for _, file := range files {
		returnFiles = append(returnFiles, file.Name)
	}
	return returnFiles, nil
}
func (fileServer *fileServer) GetFile(name string, directory string) ([]byte, error) {
	values := make(map[string][]string, 0)
	values["FileName"] = []string{name}
	values["Directory"] = []string{directory}

	result, err := http.PostForm(fileServer.address+"/getFile", values)
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()
	if result.StatusCode != http.StatusOK {
		return nil, errors.New("failed to retrieve file")
	}

	return ioutil.ReadAll(result.Body)
}

func (fileServer *fileServer) MoveFile(oldPath, newPath string) error {
	values := make(map[string][]string, 0)
	values["oldPath"] = []string{oldPath}
	values["newPath"] = []string{newPath}

	result, err := http.PostForm(fileServer.address+"/moveFile", values)
	if err != nil {
		return err
	}
	defer result.Body.Close()
	if result.StatusCode != http.StatusOK {
		return errors.New("failed to move file")
	}

	return nil
}
