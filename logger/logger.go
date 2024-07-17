package logger

import (
	rpcHelper "bitbucket.org/network-international/nextgen-libs/nextgen-helpers/rpcHelper"
	"log"
	"nextgen-tms-website/crypt"
)

var logging rpcHelper.LoggingClient

func Init(applicationName string) {
	var err error
	logging, err = rpcHelper.NewLoggingClient(applicationName, "TMSLoggingPort", "61022")
	if err != nil {
		// do not prevent service running because no logger defined.
		log.Print(err)
	}
	crypt.SetLogging(logging)
}

func GetLogger() rpcHelper.LoggingClient {
	return logging
}
