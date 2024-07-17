package models

import "database/sql"

type ValidationSet struct {
	EntryNo       string
	DataElement   string
	FailureReason string
}

type UpdateDataElement struct {
	DataGroupName string
	ValExp        sql.NullString
	ValMsg        sql.NullString
	IsAllowEmpty  bool
	DataType      string
	Options       string
	DataElementId int
}

type BulkUpdateVal struct {
	UpdateStatus bool
	Validations  []ValidationSet
}

type BulkApproval struct {
	ID         int
	Filename   string
	FileType   string
	ChangeType string
	CreatedBy  string
	CreatedAt  string
	ApprovedBy string
	ApprovedAt string
	Approved   int
}
