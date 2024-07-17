package main

import (
	"nextgen-tms-website/dal"
	"reflect"
	"testing"
)

var (
	profOne = dal.ProfileChangeHistory{ChangedAt: "1",}
	profTwo = dal.ProfileChangeHistory{ChangedAt: "2",}
	profThree = dal.ProfileChangeHistory{ChangedAt: "3",}
	elemOne = dal.DataElement{Name: "Aname", SortOrderInGroup: 1}
	elemTwo = dal.DataElement{Name: "Bname", SortOrderInGroup: 2}
	elemThree = dal.DataElement{Name: "Cname", SortOrderInGroup: 3}
	dataOne = dal.DataGroup{DataGroup: "GroupA", DataElements: nil,}
	dataTwo = dal.DataGroup{DataGroup: "GroupB", DataElements: []dal.DataElement{ elemOne }}
	dataThree = dal.DataGroup{DataGroup: "GroupB", DataElements: []dal.DataElement{ elemTwo, elemThree, elemOne }}
	dataThreeSorted = dal.DataGroup{DataGroup: "GroupB", DataElements: []dal.DataElement{ elemOne, elemTwo, elemThree }}
)

func TestSortGroups(t *testing.T) {
	tests := []struct {
		name string
		groups []*dal.DataGroup
		expected []*dal.DataGroup
	}{
		{"Nil in Nil out", nil, nil},
		{"Single in, no elements", []*dal.DataGroup{&dataOne}, []*dal.DataGroup{&dataOne}},
		{"Single in, one element", []*dal.DataGroup{&dataTwo}, []*dal.DataGroup{&dataTwo}},
		{"Single in, unsorted elements", []*dal.DataGroup{&dataThree}, []*dal.DataGroup{&dataThreeSorted}},
		{"Multiple in, pre-sorted", []*dal.DataGroup{&dataOne, &dataTwo, &dataThree}, []*dal.DataGroup{&dataOne, &dataTwo, &dataThreeSorted}},
		/*{"Multiple in, unsorted", }*/
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SortGroups(tt.groups)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("SortGroups() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSortElements(t *testing.T) {
	tests := []struct {
		name     string
		elements []dal.DataElement
		expected []dal.DataElement
	}{
		{"Nil in Nil out", nil, nil},
		{"Single in", []dal.DataElement{elemOne}, []dal.DataElement{elemOne}},
		{"Two in, already ordered", []dal.DataElement{elemOne, elemTwo}, []dal.DataElement{elemOne, elemTwo}},
		{"Two in, unordered", []dal.DataElement{elemTwo, elemOne}, []dal.DataElement{elemOne, elemTwo}},
		{"Three in, already ordered", []dal.DataElement{elemOne, elemTwo, elemThree}, []dal.DataElement{elemOne, elemTwo, elemThree}},
		{"Three in, unordered", []dal.DataElement{elemTwo, elemOne, elemThree}, []dal.DataElement{elemOne, elemTwo, elemThree}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SortElements(tt.elements)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("SortElements() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSortHistory(t *testing.T) {
	tests := []struct {
		name     string
		history  []*dal.ProfileChangeHistory
		expected []*dal.ProfileChangeHistory
	}{
		{"Nil in Nil out", nil, nil},
		{"Single in", []*dal.ProfileChangeHistory{&profOne}, []*dal.ProfileChangeHistory{&profOne}},
		{"Two in, already ordered", []*dal.ProfileChangeHistory{&profTwo, &profOne}, []*dal.ProfileChangeHistory{&profTwo, &profOne}},
		{"Two in, unordered", []*dal.ProfileChangeHistory{&profOne, &profTwo}, []*dal.ProfileChangeHistory{&profTwo, &profOne}},
		{"Three in, already ordered", []*dal.ProfileChangeHistory{&profThree, &profTwo, &profOne}, []*dal.ProfileChangeHistory{&profThree, &profTwo, &profOne}},
		{"Three in, unordered", []*dal.ProfileChangeHistory{&profTwo, &profOne, &profThree}, []*dal.ProfileChangeHistory{&profThree, &profTwo, &profOne}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SortHistory(tt.history)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("SortHistory() = %v, want %v", got, tt.expected)
			}
		})
	}
}
