package main

import (
	"archive/zip"
	TLSUtils "bitbucket.org/network-international/nextgen-libs/nextgen-helpers/tlsutils"
	"io/ioutil"
	"os"
	"testing"
)

func TestCertificateGenerationAsZip(t *testing.T) {
	RPICertFile := "test/RPI/key.pem"
	RPIKeyFile := "test/RPI/cert.pem"
	CN := "123456"
	password := "ben"

	_ = os.RemoveAll(CN + "/") // to match the environment
	clientCert, err := TLSUtils.GenerateEncryptedClientCertificate(RPICertFile, RPIKeyFile, CN, password)

	if err!=nil {
		t.Errorf("Cert failed %v",err)
	}

	caCert, err := ioutil.ReadFile(RPICertFile)

	newZipFile, err := os.Create("test/test_example.zip")
	if err != nil {
		t.Error(err)
		return
	}
	defer newZipFile.Close()

	wZip := zip.NewWriter(newZipFile)
	defer wZip.Close()

	clientCertFile, err := wZip.Create("Server0Q.pfx")
	if err != nil {
		t.Errorf("Error generating certificate (Zipping Certificate)")
		return
	}
	_, err = clientCertFile.Write(clientCert)
	if err != nil {
		t.Errorf( "Error generating certificate (Zipping Certificate)")
		return
	}

	caCertFile, err := wZip.Create("Server0QRoot.cer")
	if err != nil {
		t.Errorf( "Error generating certificate (Zipping Certificate)")
		return
	}
	_, err = caCertFile.Write(caCert)
	if err != nil {
		t.Errorf( "Error generating certificate (Zipping Certificate)")
		return
	}
}

