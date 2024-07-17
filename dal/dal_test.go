package dal

import "testing"

func TestDataGroupsContainsDataElement(t *testing.T) {
	tests := []struct {
		name                     string
		groups                   []*DataGroup
		elementId                int
		expectedFound            bool
		expectedDgIndex          int
		expectedDataElementIndex int
	}{
		{
			name: "Nil slice - element not found", groups: nil, elementId: 1, expectedFound: false, expectedDgIndex: -1, expectedDataElementIndex: -1,
		},
		{
			name: "Empty slice - element not found", groups: nil, elementId: 1, expectedFound: false, expectedDgIndex: -1, expectedDataElementIndex: -1,
		},
		{
			name: "Slice with single data group & element - given element ID exists in slice - returns true",
			groups: []*DataGroup{{DataElements: []DataElement{{ElementId: 1}}}},
			elementId: 1, expectedFound: true, expectedDgIndex: 0, expectedDataElementIndex: 0,
		},
		{
			name: "Slice with single data group & element - given element ID does not exist in slice - returns false",
			groups: []*DataGroup{{DataElements: []DataElement{{ElementId: 1}}}},
			elementId: 2, expectedFound: false, expectedDgIndex: -1, expectedDataElementIndex: -1,
		},
		{
			name: "Single data group with 4 elements - required element present",
			groups: []*DataGroup{
				{DataElements: []DataElement{{ElementId: 1}, {ElementId: 2}, {ElementId: 3}, {ElementId: 4}}},
			},
			elementId: 4, expectedFound: true, expectedDgIndex: 0, expectedDataElementIndex: 3,
		},
		{
			name: "2 data groups - required element present",
			groups: []*DataGroup{
				{DataElements: []DataElement{{ElementId: 1}}},
				{DataElements: []DataElement{{ElementId: 2}}},
			},
			elementId: 2, expectedFound: true, expectedDgIndex: 1, expectedDataElementIndex: 0,
		},
		{
			name: "3 data groups - required element present",
			groups: []*DataGroup{
				{DataElements: []DataElement{{ElementId: 1}}},
				{DataElements: []DataElement{{ElementId: 2}}},
				{DataElements: []DataElement{{ElementId: 3}}},
			},
			elementId: 3, expectedFound: true, expectedDgIndex: 2, expectedDataElementIndex: 0,
		},
		{
			name: "4 data groups - required element present",
			groups: []*DataGroup{
				{DataElements: []DataElement{{ElementId: 1}}},
				{DataElements: []DataElement{{ElementId: 2}}},
				{DataElements: []DataElement{{ElementId: 3}}},
				{DataElements: []DataElement{{ElementId: 4}}},
			},
			elementId: 4, expectedFound: true, expectedDgIndex: 3, expectedDataElementIndex: 0,
		},
		{
			name: "Multiple data groups and elements",
			groups: []*DataGroup{
				{DataElements: []DataElement{{ElementId: 1}, {ElementId: 2}, {ElementId: 3}, {ElementId: 4}}},
				{DataElements: []DataElement{{ElementId: 5}, {ElementId: 6}, {ElementId: 7}, {ElementId: 8}}},
				{DataElements: []DataElement{{ElementId: 9}, {ElementId: 10}, {ElementId: 11}, {ElementId: 12}}},
				{DataElements: []DataElement{{ElementId: 13}, {ElementId: 14}, {ElementId: 15}, {ElementId: 16}}},
			},
			elementId: 10, expectedFound: true, expectedDgIndex: 2, expectedDataElementIndex: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := DataGroupsContainsDataElement(tt.groups, tt.elementId)
			if got != tt.expectedFound {
				t.Errorf("DataGroupsContainsDataElement() got = %v, expectedFound %v", got, tt.expectedFound)
			}

			if got1 != tt.expectedDgIndex {
				t.Errorf("DataGroupsContainsDataElement() got1 = %v, expectedFound %v", got1, tt.expectedDgIndex)
			}

			if got2 != tt.expectedDataElementIndex {
				t.Errorf("DataGroupsContainsDataElement() got2 = %v, expectedFound %v", got2, tt.expectedDataElementIndex)
			}
		})
	}
}
