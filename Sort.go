package main

import (
	"sort"
	"nextgen-tms-website/dal"
)

func SortGroups(groups []*dal.DataGroup) []*dal.DataGroup {
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].DataGroup < groups[j].DataGroup
	})

	for _, g := range groups {
		g.DataElements = SortElements(g.DataElements)
	}
	return groups
}

func SortElements(elements []dal.DataElement) []dal.DataElement {
	sort.Slice(elements, func(i, j int) bool {
		return elements[i].SortOrderInGroup < elements[j].SortOrderInGroup
	})

	return elements
}

func SortHistory(history []*dal.ProfileChangeHistory) []*dal.ProfileChangeHistory {
	sort.Slice(history, func(i, j int) bool {
		return history[i].ChangedAt > history[j].ChangedAt
	})
	return history
}
