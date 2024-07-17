package TMSExportHandler

import (
	exporter "bitbucket.org/network-international/nextgen-libs/nextgen-helpers/exportHandler"
	"golang.org/x/exp/maps"
	"nextgen-tms-website/PED"
)

func NewHandler(repository PED.Repository) Handler {
	return &handler{
		PEDs: repository,
	}
}

type Handler interface {
	ExportPEDs(searchTerm string, availableAcquirers string) (exporter.ExportableItems, error)
}

type handler struct {
	PEDs PED.Repository
}

func (h *handler) ExportPEDs(searchTerm string, availableAcquirers string) (exporter.ExportableItems, error) {
	resultExportableItems := exporter.ExportableItems{}
	foundPEDs, err := h.PEDs.FindBySearchTermAndAcquirer(searchTerm, availableAcquirers)
	if err != nil {
		return resultExportableItems, err
	}

	resultExportableItems.Items = make([]exporter.ExportableItem, 0)

	for _, PED := range foundPEDs {
		pedDataMap, err := exporter.StructToExportableItem(*PED)
		if err != nil {
			return resultExportableItems, err
		}

		maps.Copy(PED.PEDInfo, pedDataMap)
		resultExportableItems.Items = append(resultExportableItems.Items, PED.PEDInfo)
	}

	return resultExportableItems, nil
}
