package namespace

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/sahilm/fuzzy"
)

func NamespaceFilter(term string, targets []string) []list.Rank {
	var ranks = fuzzy.FindNoSort(term, targets)
	result := make([]list.Rank, len(ranks))
	for i, r := range ranks {
		result[i] = list.Rank{
			Index:          r.Index,
			MatchedIndexes: r.MatchedIndexes,
		}
	}
	return result
}

// func NamespaceFilter(term string, targets []string) []list.Rank {
// 	var matches []int
// 	for i, target := range targets {
// 		if strings.Contains(target, term) {
// 			matches = append(matches, i)
// 		}
// 	}
// 	result := make([]list.Rank, len(matches))
// 	for i, r := range matches {
// 		result[i] = list.Rank{
// 			Index:          r,
// 			MatchedIndexes: []int{},
// 		}
// 	}
// 	return result
// }
