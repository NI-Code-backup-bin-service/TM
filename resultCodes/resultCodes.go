package resultCodes

type ResultCode int

const (
	SUCCESS ResultCode = iota
	TID_NOT_PRESENT
	TID_DOES_NOT_EXIST
	MID_DOES_NOT_EXIST
	TID_NOT_UNIQUE_PRIMARY_TID_DUPLICATE
	TID_NOT_UNIQUE_SECONDARY_TID_DUPLICATE
	MID_NOT_UNIQUE_PRIMARY_MID_DUPLICATE
	MID_NOT_UNIQUE_SECONDARY_MID_DUPLICATE
	DATABASE_CONNECTION_ERROR
	DATABASE_QUERY_ERROR
	DATA_GROUPS_MAPPINGS
)

var resultMessagesTable = map[ResultCode]string{
	SUCCESS:                                "",
	TID_NOT_PRESENT:                        "TID not present",
	TID_DOES_NOT_EXIST:                     "TID does not exist",
	MID_DOES_NOT_EXIST:                     "MID does not exist",
	TID_NOT_UNIQUE_PRIMARY_TID_DUPLICATE:   "TID already exists as a primary TID and must be unique",
	TID_NOT_UNIQUE_SECONDARY_TID_DUPLICATE: "TID already exists as a secondary TID and must be unique",
	MID_NOT_UNIQUE_PRIMARY_MID_DUPLICATE:   "MID already exists as a primary MID and must be unique",
	MID_NOT_UNIQUE_SECONDARY_MID_DUPLICATE: "MID already exists as a secondary MID and must be unique",
	DATABASE_CONNECTION_ERROR:              "An unexpected error occurred connecting to the database",
	DATABASE_QUERY_ERROR:                   "An unexpected error occurred executing the database query",
	DATA_GROUPS_MAPPINGS:                   "Data group/element mapping failed",
}

func GetErrorMsgByCode(code ResultCode) string {
	return resultMessagesTable[code]
}
