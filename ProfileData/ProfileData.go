package ProfileData

import "nextgen-tms-website/entities"

type Repository interface {
	SetDataValueByElementIdAndProfileIdWithoutApproval(dataElementId, profileId int, value string, user entities.TMSUser) error
}