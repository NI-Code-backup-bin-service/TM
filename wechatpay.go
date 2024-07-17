package main

import (
	rpcHelp "bitbucket.org/network-international/nextgen-libs/nextgen-helpers/rpcHelper"
	txn "bitbucket.org/network-international/nextgen-libs/nextgen-tg-protobuf/Transaction"
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"nextgen-tms-website/dal"
	"nextgen-tms-website/entities"
	"nextgen-tms-website/models"
)

func registerWeChatPaySubMerchant(profileId int, siteId int, elements map[int]string, user *entities.TMSUser, dataElementDetails map[string]models.DataElementsAndGroup) {
	request := txn.TransactionRequest{}
	request.Header = rpcHelp.SetUpHeader(Version, ApplicationName)
	request.Request = &txn.Request{}
	wcpDetails := txn.WeChatPaySubMerchantOnboardRequest{}

	siteGroups, chainGroups, acquirerGroups, globalGroups, err := dal.GetSiteGroupsData(siteId, profileId)
	if err != nil {
		_, _ = logging.Error(err.Error())
		return
	}

	wcpDetails.SubMerchantID = getDataElementValue(getElementIdFromDataElementMap(dataElementDetails, "weChatPay-wcpSubMerchantId"), elements, siteGroups, chainGroups, acquirerGroups, globalGroups)
	wcpDetails.BusinessCategory = getDataElementValue(dataElementDetails["weChatPay-wcpBcc"].DataElementID, elements, siteGroups, chainGroups, acquirerGroups, globalGroups)
	if len(wcpDetails.BusinessCategory) == 0 {
		return
	}

	wcpDetails.MerchantCategoryCode = getDataElementValue(getElementIdFromDataElementMap(dataElementDetails, "weChatPay-wcpMcc"), elements, siteGroups, chainGroups, acquirerGroups, globalGroups)
	if len(wcpDetails.MerchantCategoryCode) == 0 {
		return
	}

	wcpDetails.OfficePhone = getDataElementValue(getElementIdFromDataElementMap(dataElementDetails, "weChatPay-wcpContactPhone"), elements, siteGroups, chainGroups, acquirerGroups, globalGroups)
	wcpDetails.MobilePhoneNo = getDataElementValue(dataElementDetails["weChatPay-wcpContactPhone"].DataElementID, elements, siteGroups, chainGroups, acquirerGroups, globalGroups)
	if len(wcpDetails.OfficePhone) == 0 {
		return
	}

	wcpDetails.FullName = getDataElementValue(getElementIdFromDataElementMap(dataElementDetails, "weChatPay-wcpContactName"), elements, siteGroups, chainGroups, acquirerGroups, globalGroups)
	if len(wcpDetails.FullName) == 0 {
		return
	}

	wcpDetails.ContactEmail = getDataElementValue(getElementIdFromDataElementMap(dataElementDetails, "weChatPay-wcpContactEmail"), elements, siteGroups, chainGroups, acquirerGroups, globalGroups)
	if len(wcpDetails.ContactEmail) == 0 {
		return
	}

	wcpDetails.MerchantName = getDataElementValue(getElementIdFromDataElementMap(dataElementDetails, "store-name"), elements, siteGroups, chainGroups, acquirerGroups, globalGroups)
	wcpDetails.MerchantShortName = getDataElementValue(getElementIdFromDataElementMap(dataElementDetails, "store-name"), elements, siteGroups, chainGroups, acquirerGroups, globalGroups)
	if len(wcpDetails.MerchantName) == 0 {
		return
	}

	wcpDetails.MerchantRemark = getDataElementValue(getElementIdFromDataElementMap(dataElementDetails, "store-merchantNo"), elements, siteGroups, chainGroups, acquirerGroups, globalGroups)
	if len(wcpDetails.MerchantRemark) == 0 {
		return
	}

	wcpDetails.StoreAddress = getDataElementValue(getElementIdFromDataElementMap(dataElementDetails, "store-addressLine1"), elements, siteGroups, chainGroups, acquirerGroups, globalGroups)
	if len(wcpDetails.StoreAddress) == 0 {
		return
	}

	val := getDataElementValue(getElementIdFromDataElementMap(dataElementDetails, "emv-terminalCountryCode"), elements, siteGroups, chainGroups, acquirerGroups, globalGroups)
	if len(val) == 0 {
		return
	}

	wcpDetails.MerchantCountryCode = val[1:]

	wcpDetailsAny, err := ptypes.MarshalAny(&wcpDetails)
	if err != nil {
		_, _ = logging.Error(err.Error())
		return
	}
	request.Request.Details = append(request.Request.Details, wcpDetailsAny)
	_, _ = logging.Debug(fmt.Sprintf("registerWeChatPaySubMerchant: Registering site %d", siteId))
	client, clientFound := GRPCclients["WeChatPay"]
	if !clientFound {
		_, _ = logging.Warning("registerWeChatPaySubMerchant: Error obtaining GRPC Client")
		err = errors.New("registerWeChatPaySubMerchant: Error obtaining GRPC Client")
		_, _ = logging.Error(err.Error())
		return
	} else {
		_, _ = logging.Debug("Register WCP SubMerchant: client found, connection state: " + client.GetConnection().GetState().String())
		grpcReply := new(txn.TransactionResponse)
		err = rpcHelp.ExecuteGRPC(client, &request, grpcReply, logging)
		if err != nil {
			err = errors.New("registerWeChatPaySubMerchant: ExecuteGRPC failed, error")
			if err != nil {
				_, _ = logging.Error(err.Error())
				return
			}
			return
		} else {
			resp := new(txn.WeChatPaySubMerchantOnboardResponse)
			for _, pItem := range grpcReply.Response.Details {
				test := pItem.GetTypeUrl()
				_, _ = logging.Debug(test)
				if test == "type.googleapis.com/TransactionServer.WeChatPay/WeChatPaySubMerchantOnboardResponse" {
					err = ptypes.UnmarshalAny(pItem, resp)
					if err != nil {
						_, _ = logging.Error(err.Error())
						return
					}
				}
			}
			if len(resp.SubMerchantId) > 0 {
				_, _ = logging.Debug(fmt.Sprintf("registerWeChatPaySubMerchant: Site %d registered successfully: %s", siteId, resp.SubMerchantId))
				err = dal.SaveElementData(profileId, getElementIdFromDataElementMap(dataElementDetails, "weChatPay-wcpSubMerchantId"), resp.SubMerchantId, user.Username, 1, 0)
				if err != nil {
					_, _ = logging.Error(err.Error())
				}
			}
		}
	}
}

func getElementIdFromDataElementMap(dataElementDetails map[string]models.DataElementsAndGroup, key string) int {
	if element, ok := dataElementDetails[key]; ok {
		return element.DataElementID
	}

	return 0
}

func getDataElementValue(dataElementId int, elements map[int]string, siteGroups, chainGroups, acquirerGroups, globalGroups []*dal.DataGroup) string {
	val, ok := elements[dataElementId]
	if ok {
		return val
	}
	for _, g := range siteGroups {
		for _, e := range g.DataElements {
			if e.ElementId == dataElementId {
				return e.DataValue
			}
		}
	}
	for _, g := range chainGroups {
		for _, e := range g.DataElements {
			if e.ElementId == dataElementId {
				return e.DataValue
			}
		}
	}
	for _, g := range acquirerGroups {
		for _, e := range g.DataElements {
			if e.ElementId == dataElementId {
				return e.DataValue
			}
		}
	}
	for _, g := range globalGroups {
		for _, e := range g.DataElements {
			if e.ElementId == dataElementId {
				return e.DataValue
			}
		}
	}
	return ""
}
