package main

import (
	"github.com/charmbracelet/bubbles/list"
	"sort"
	"strings"
)

func FindItem(items []list.Item, search list.Item) int {
	i, found := sort.Find(len(items), func(i int) int {
		return strings.Compare(items[i].FilterValue(), search.FilterValue())
	})
	if !found {
		return -1
	}
	return i
}
