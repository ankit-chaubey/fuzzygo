// Package fuzzy provides fast, Unicode-safe fuzzy string matching and ranking.
//
// Author: Ankit Chaubey // @ankit-chaubey
// Repository: https://github.com/ankit-chaubey/fuzzygo
// Description: A high-performance fuzzy search and ranking engine for Go,
//              optimized for symbol, API, and identifier search.

package fuzzy

import (
	"sort"
	"strings"
	"unicode"
)

type Match struct {
	Item  string
	Score int
}

func isBoundary(r rune) bool {
	return r == '_' || r == '.' || r == '-' || r == '/'
}

func Score(query, target string) int {
	q := []rune(strings.ToLower(query))
	t := []rune(strings.ToLower(target))

	qi := 0
	score := 0
	consecutive := 0
	maxConsecutive := 0
	lastMatch := -1

	for ti := 0; ti < len(t) && qi < len(q); ti++ {
		if t[ti] == q[qi] {
			base := 10

			if ti == 0 || isBoundary(t[ti-1]) {
				base += 15
			}
			if lastMatch+1 == ti {
				consecutive++
				base += consecutive * 5
			} else {
				consecutive = 1
			}
			if unicode.IsUpper(t[ti]) {
				base += 3
			}

			score += base
			lastMatch = ti
			qi++

			if consecutive > maxConsecutive {
				maxConsecutive = consecutive
			}
		} else {
			score -= 1
			consecutive = 0
		}
	}

	if qi != len(q) {
		return -1
	}

	score += maxConsecutive * 10
	score -= len(t) - len(q)

	return score
}

func Rank(query string, data []string, limit int) ([]Match, int) {
	results := make([]Match, 0)
	total := 0

	for _, item := range data {
		s := Score(query, item)
		if s > 0 {
			total++
			results = append(results, Match{
				Item:  item,
				Score: s,
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}

	return results, total
}
