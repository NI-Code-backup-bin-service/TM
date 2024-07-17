package models

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type GenerateOTPResponse struct {
	PIN        string `json:"pin"`
	ExpiryTime string `json:"expiryTime"`
}
