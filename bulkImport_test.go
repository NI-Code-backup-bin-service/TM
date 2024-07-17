package main

import (
	rpcHelp "bitbucket.org/network-international/nextgen-libs/nextgen-helpers/rpcHelper"
	"errors"
	"reflect"
	"sort"
	"testing"
)

var (
	Invalid_Element_Positions = errors.New("invalid element position")
)

var (
	recordsToColumnsTestOne = []DataColumn{
		{Position: 0, Name: "merchantNo", DataGroup: "store"},
		{Position: 1, Name: "name", DataGroup: "store"},
		{Position: 2, Name: "addressLine1", DataGroup: "store"},
	}
)

func setupTemplateMap() map[int]DataElement {
	returnValue := make(map[int]DataElement)

	returnValue[888] = DataElement{ElementId: 888, Data: "eggs"}
	returnValue[1] = DataElement{ElementId: 1, Data: "bread"}
	returnValue[137] = DataElement{ElementId: 137}
	returnValue[55] = DataElement{ElementId: 55, Data: "elementTest4"}

	return returnValue
}

func setupNewSitesElementsOverwritten() NewSitesElements {
	var returnValue NewSitesElements
	var entryOne NewProfileEntry
	var entryTwo NewProfileEntry
	var entryOneElements []DataElement
	var entryTwoElements []DataElement

	entryOneElements = append(entryOneElements, DataElement{ElementId: 888, Data: "elementTest1"})
	entryOneElements = append(entryOneElements, DataElement{ElementId: 1, Data: "elementTest2"})
	entryOneElements = append(entryOneElements, DataElement{ElementId: 137, Data: "elementTest3"})
	entryOneElements = append(entryOneElements, DataElement{ElementId: 55, Data: "elementTest4"})

	entryOne.setRef(1)
	entryOne.DataElements = entryOneElements

	entryTwoElements = append(entryTwoElements, DataElement{ElementId: 888, Data: "elementTest1"})
	entryTwoElements = append(entryTwoElements, DataElement{ElementId: 1, Data: "elementTest2"})
	entryTwoElements = append(entryTwoElements, DataElement{ElementId: 137})
	entryTwoElements = append(entryTwoElements, DataElement{ElementId: 55, Data: "elementTest4"})

	entryTwo.setRef(2)
	entryTwo.DataElements = entryTwoElements

	returnValue.NewSites = append(returnValue.getNewSites(), entryOne)
	returnValue.NewSites = append(returnValue.getNewSites(), entryTwo)

	return returnValue
}

func setupNewSitesElements() NewSitesElements {
	var returnValue NewSitesElements
	var entryOne NewProfileEntry
	var entryTwo NewProfileEntry
	var entryOneElements []DataElement
	var entryTwoElements []DataElement

	entryOneElements = append(entryOneElements, DataElement{ElementId: 888, Name: "elementTest1", Data: "elementTest1"})
	entryOneElements = append(entryOneElements, DataElement{ElementId: 1, Name: "elementTest2", GroupName: "dataGroup1", Data: "elementTest2"})
	entryOneElements = append(entryOneElements, DataElement{ElementId: 137, Name: "elementTest3", GroupName: "dataGroup1", Data: "elementTest3"})
	entryOneElements = append(entryOneElements, DataElement{ElementId: 55, Name: "elementTest4", GroupName: "dataGroup1", Data: "elementTest4"})

	entryOne.setRef(0)
	entryOne.DataElements = entryOneElements

	entryTwoElements = append(entryTwoElements, DataElement{ElementId: 888, Name: "elementTest1", Data: "elementTest1"})
	entryTwoElements = append(entryTwoElements, DataElement{ElementId: 1, Name: "elementTest2", GroupName: "dataGroup1", Data: "elementTest2"})

	entryTwo.setRef(1)
	entryTwo.DataElements = entryTwoElements

	returnValue.NewSites = append(returnValue.getNewSites(), entryOne)
	returnValue.NewSites = append(returnValue.getNewSites(), entryTwo)

	return returnValue
}

func setupErrorDataRecords() [][]string {
	arrayOne := []string{"elementTest1", "elementTest2", "elementTest3", "elementTest4", "elementTest5"}
	return [][]string{arrayOne}
}

func setupDataRecords() [][]string {
	arrayOne := []string{"elementTest1", "elementTest2", "elementTest3", "elementTest4"}
	arrayTwo := []string{"elementTest1", "elementTest2"}

	return [][]string{arrayOne, arrayTwo}
}

func setupColumnsStructs() []DataColumn {
	return []DataColumn{
		{Position: 0, Name: "elementTest1", DataGroup: "", ElementId: 888},
		{Position: 1, Name: "elementTest2", DataGroup: "dataGroup1", ElementId: 1},
		{Position: 2, Name: "elementTest3", DataGroup: "dataGroup1", ElementId: 137},
		{Position: 3, Name: "elementTest4", DataGroup: "dataGroup1", ElementId: 55},
	}
}

func setupColumnErrorRecords() [][]string {
	arrayOne := []string{"store.merchantNo", "store.name", "store.addressLine1", "store.addressLine2"}
	arrayTwo := []string{"8765433000", "RTA1", "WonderLane", "Bournemouth"}
	arrayThree := []string{"8765433001", "RTA1", "WonderLane", "Bournemouth"}
	arrayFour := []string{"8765433000", "RTA1", "WonderLane", "Bournemouth"}
	arrayFive := []string{"8765433000", "RTA1", "WonderLane", "Bournemouth"}
	arraySix := []string{"8765433002", "RTA1", "WonderLane", "Bournemouth"}

	return [][]string{arrayOne, arrayTwo, arrayThree, arrayFour, arrayFive, arraySix}
}

func setupColumnErrorRecordsEmpty() [][]string {
	arrayOne := []string{"store.merchantNo", "store.name", "store.addressLine1", "store.addressLine2"}
	arrayTwo := []string{"8765433000", "RTA1", "WonderLane", "Bournemouth"}
	arrayThree := []string{"", "RTA1", "WonderLane", "Bournemouth"}
	arrayFour := []string{"", "RTA1", "WonderLane", "Bournemouth"}
	arrayFive := []string{"8765433001", "RTA1", "WonderLane", "Bournemouth"}
	arraySix := []string{"", "RTA1", "WonderLane", "Bournemouth"}

	return [][]string{arrayOne, arrayTwo, arrayThree, arrayFour, arrayFive, arraySix}
}

func setupColumnExpectedOne() []string {
	errorArray := []string{"1", "3", "4"}
	return errorArray
}

func setupColumnDataRecords() [][]string {
	arrayOne := []string{"store.merchantNo", "store.name", "store.addressLine1", "store.addressLine2"}
	arrayTwo := []string{"8765433000", "RTA1", "WonderLane", "Bournemouth"}
	arrayThree := []string{"8765433001", "RTA2", "WonderLane", "Bournemouth"}
	arrayFour := []string{"8765433002", "RTA3", "WonderLane", "Bournemouth"}

	return [][]string{arrayOne, arrayTwo, arrayThree, arrayFour}
}

func TestBuildNewSites(t *testing.T) {
	logging, _ = rpcHelp.NewLoggingClient("test", "", "")

	tests := []struct {
		name        string
		sitesMap    NewSitesElements
		templateMap map[int]DataElement
		expected    NewSitesElements
	}{
		{"", setupNewSitesElements(), setupTemplateMap(), setupNewSitesElementsOverwritten()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildNewSites(tt.sitesMap, tt.templateMap)

			// Need to sort these due to golang's unordered map handling
			for _, site := range got.getNewSites() {
				sort.Slice(site.DataElements, func(i, j int) bool { return site.DataElements[i].ElementId < site.DataElements[j].ElementId })
			}

			for _, site := range tt.expected.getNewSites() {
				sort.Slice(site.DataElements, func(i, j int) bool { return site.DataElements[i].ElementId < site.DataElements[j].ElementId })
			}

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("buildNewSites() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestConvertDataRowsToMap(t *testing.T) {
	logging, _ = rpcHelp.NewLoggingClient("test", "", "")
	var blankNewSitesElements NewSitesElements

	tests := []struct {
		name     string
		records  [][]string
		columns  []DataColumn
		expected NewSitesElements
		err      error
	}{
		{"Successful test", setupDataRecords(), setupColumnsStructs(), setupNewSitesElements(), nil},
		{"Error test", setupErrorDataRecords(), setupColumnsStructs(), blankNewSitesElements, Invalid_Element_Positions},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertDataRowsToElements(tt.records, tt.columns, true)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("convertDataRowsToElements() = %v, want %v", got, tt.expected)
			}
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("convertDataRowsToElements() error response = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetElementIdFromPosition(t *testing.T) {
	logging, _ = rpcHelp.NewLoggingClient("test", "", "")

	tests := []struct {
		name     string
		input    []DataColumn
		position int
		expected int
	}{
		{"First position", setupColumnsStructs(), 0, 888},
		{"Last position", setupColumnsStructs(), 3, 55},
		{"Invalid position", setupColumnsStructs(), 999, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getElementIdFromPosition(tt.input, tt.position)
			if got != tt.expected {
				t.Errorf("getElementIdFromPosition() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestConvertRecordsToColumns(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []DataColumn
	}{
		{"Site Name and data elements", []string{"store.merchantNo", "store.name", "store.addressLine1"}, recordsToColumnsTestOne},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertRecordsToColumns(tt.input, false)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("convertRecordsToColumns() = \n%v, \nwant \n%v", got, tt.expected)
			}
		})
	}
}

func TestCheckColumn(t *testing.T) {
	tests := []struct {
		name               string
		records            [][]string
		column             int
		is_allowed_empty   bool
		expectedBadRecords []string
		expectedBoolean    bool
	}{
		{"Successful test", setupColumnDataRecords(), 0, false, nil, false},
		{"Error test Not Allowed Empty", setupColumnErrorRecords(), 0, false, setupColumnExpectedOne(), true},
		{"Error test Allowed Empty", setupColumnErrorRecordsEmpty(), 0, true, nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBoolean, recordsOut := checkColumn(tt.records, tt.column, tt.is_allowed_empty)
			if gotBoolean != tt.expectedBoolean {
				t.Errorf("checkColumn() = %v, want %v", gotBoolean, tt.expectedBoolean)
			}
			if !reflect.DeepEqual(recordsOut, tt.expectedBadRecords) {
				t.Errorf("checkColumn() error response = %v, want %v", recordsOut, tt.expectedBadRecords)
			}
		})
	}
}
