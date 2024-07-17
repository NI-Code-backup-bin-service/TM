package models

type CreateTIDRequest struct {
	MID          string   `json:"mId" validate:"required""`
	TID          string   `json:"tId" validate:"required"`
	SerialNumber string   `json:"serialNumber" validate:"required"`
	ApkVersion   string   `json:"apkVersion" validate:"required"`
	TPApkVersion []string `json:"tpApkVersion,omitempty"`
}

type UpdateTIDRequest struct {
	MID          string   `json:"mId" validate:"required""`
	TID          string   `json:"tId" validate:"required"`
	SerialNumber string   `json:"serialNumber" validate:"required"`
	ApkVersion   string   `json:"apkVersion" validate:"required"`
	TPApkVersion []string `json:"tpApkVersion,omitempty"`
}

type GenerateOTPRequest struct {
	MID string `json:"mId" validate:"required""`
	TID string `json:"tId" validate:"required"`
}
