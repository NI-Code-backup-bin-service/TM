package entities

type SiteFraudLimitModel struct {
	SiteLimits            VelocityLimit
	Limits                []VelocityLimit
	HasSavePermission     bool
	AvailableTransactions []VelocityTransactions
	AvailableLimits       []VelocityLimitTypes
}

type TransactionLimitGroupModel struct {
	SchemeID  string
	TxnLimits []TxnLimit
}

type TransactionLimitMaintenanceModel struct {
	TransactionTypes      []TransactionLimitGroup
	AvailableTransactions []VelocityTransactions
	AvailableLimits       []VelocityLimitTypes
	Override              string
}

type VelocityTransactions struct {
	TxnType         string
	TxnTypeReadable string
}

type VelocityLimitTypes struct {
	Limittype		string
	Identifier		string
}

type TransactionLimitGroup struct {
	TxnLimitID    int
	TxnLimitGroup string
	TxnLimits     []TxnLimit
}

type VelocityLimit struct {
	ID               string
	Scheme           string
	DailyCount       int
	DailyLimit       int
	BatchCount       int
	BatchLimit       int
	SingleTransLimit int
	TxnLimits        []TxnLimit
	Level            int
	Index			 int
}

type TxnLimit struct {
	TxnLimitID      string
	TxnType         string
	TxnTypeReadable string
	LimitType       string
	Value           int
}
