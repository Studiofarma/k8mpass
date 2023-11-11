package main

import (
	"github.com/charmbracelet/bubbles/list"
)

func FindItem(items []list.Item, search list.Item) int {
	var idx = -1
	for i, item := range items {
		if item.FilterValue() == search.FilterValue() {
			idx = i
		}
	}
	return idx
}
