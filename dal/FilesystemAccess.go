package dal

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func SaveFileToDir(filename string, contents []byte, directory string) error{
	_, err := os.Stat(directory)
	if os.IsNotExist(err) {
		os.Mkdir(directory, os.ModePerm)
	} else if err != nil {
		return err
	}

	newFilePath := filepath.Join(directory, filename)

	_, err = os.Stat(newFilePath)
	if err == nil {
		return os.ErrExist
	}

	return ioutil.WriteFile(newFilePath,contents,os.ModePerm)
}
