package DataGroup

type DataGroup struct {
	ID          int
	Name        string
	DisplayName string
	Selected    bool
	Preselected bool
	ToolTip     string
}

type Repository interface {
	// Returns all data groups and if they're selected or preselected for a given site profile ID
	FindForSiteByProfileId(siteProfileId int) ([]DataGroup, error)
	// Returns all data groups and if they're selected or preselected for a given TID profile ID
	FindForTidByProfileId(tidProfileId int) ([]DataGroup, error)
	// Returns the DataGroup an element belongs to by the element ID
	FindByDataElementId(elementId int) (DataGroup, error)
}
