package PED

import (
	"fmt"
)

type XmmppMessageType int
const (
	XMPP_MESSAGE_TYPE_UPLOAD_LOGS XmmppMessageType = iota
)
func GenerateXmppQuery(messageType XmmppMessageType, params... interface{}) string {
	var query string
	var requestType string
	switch messageType {
	case XMPP_MESSAGE_TYPE_UPLOAD_LOGS:
		requestType = "sendlogsrequest"
		maxUploadInKb := params[0].(int)
		//multiply by 1000 as the param requires bytes
		query = fmt.Sprintf("<maxLogcatSize>%d</maxLogcatSize>", maxUploadInKb * 1000)
	default:
		return ""
	}

	return fmt.Sprintf("<query xmlns='jabber:iq:%s'>%s</query>", requestType, query)
}
