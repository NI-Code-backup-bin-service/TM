package entities

type PaymentServiceGroup struct {
	Id           uint
	Name         string
	ServiceCount string
}

type PaymentService struct {
	Id      uint
	GroupId uint
	Name    string
}

type PaymentServiceGroupImportModel struct {
	GroupsCreated   int
	ServicesCreated int
	FailedRows      int
	Groups          map[string]PaymentServiceImportGroup
}

type PaymentServiceImportGroup struct {
	Id       int
	Services map[string]bool
}
