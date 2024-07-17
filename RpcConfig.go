package main

import (
	cfg "bitbucket.org/network-international/nextgen-libs/nextgen-helpers/configHelper"
	"time"
)

func getRpcConfig() []cfg.RpcClientConfig {
	return []cfg.RpcClientConfig{
		{
			ProcessName:          "SaleAdvice",
			ConnectionTimeout:    time.Duration(cfg.GetInt("ShortConnectionTimeout", 2)) * time.Second,
			IPAddress:            cfg.GetString("TMServerIP", "localhost") + ":" + cfg.GetString("B24Port", "61010"),
			RpcFunction:          "/TransactionServer.Base24/Advice",
			RpcProcessingTimeout: time.Duration(cfg.GetInt("ShortProcessingTimeout", 36)) * time.Second,
		},
		{
			ProcessName:          "WeChatPay",
			ConnectionTimeout:    time.Duration(cfg.GetInt("ShortConnectionTimeout", 2)) * time.Second,
			IPAddress:            cfg.GetString("TMServerIP", "localhost") + ":" + cfg.GetString("WeChatPayPort", "61018"),
			RpcFunction:          "/TransactionServer.WeChatPay/WeChatPaySubMerchantOnboard",
			RpcProcessingTimeout: time.Duration(cfg.GetInt("ShortProcessingTimeout", 36)) * time.Second,
		},
		{
			ProcessName:          "TMFileUpload",
			ConnectionTimeout:    time.Duration(cfg.GetInt("ShortConnectionTimeout", 2)) * time.Second,
			IPAddress:            cfg.GetString("TMServerIP", "localhost") + ":" + cfg.GetString("TMFileServerPort", "61150"),
			RpcFunction:          "/TransactionServer.TMFileServer/FileUpload",
			RpcProcessingTimeout: time.Duration(cfg.GetInt("ShortProcessingTimeout", 36)) * time.Second,
		},
	}
}

func getDiagnosticRpcConfig() []cfg.RpcClientConfig {
	return []cfg.RpcClientConfig{
		{
			ProcessName:          "Base 24",
			ConnectionTimeout:    time.Duration(cfg.GetInt("ShortConnectionTimeout", 2)) * time.Second,
			IPAddress:            cfg.GetString("TMServerIP", "localhost") + ":" + cfg.GetString("B24Port", "61010"),
			RpcFunction:          "/TransactionServer.Base24/Advice",
			RpcProcessingTimeout: time.Duration(cfg.GetInt("ShortProcessingTimeout", 36)) * time.Second,
		}}
}
