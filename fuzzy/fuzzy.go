// Package fuzzy provides fast, Unicode-safe fuzzy string matching and ranking.
//
// Author: Ankit Chaubey // @ankit-chaubey
// Repository: https://github.com/ankit-chaubey/fuzzygo
// Description: A high-performance fuzzy search and ranking engine for Go,
//              optimized for symbol, API, and identifier search.

package fuzzy

import (
	"container/heap"
	"strings"
	"unicode"
)

type Match struct {
	Item  string
	Score int
}

type Entry struct {
	Text      string
	LowerText string
	Runes     []rune
	Lower     []rune
}

type Config struct {
	BaseScore        int
	BoundaryBonus    int
	ConsecutiveBonus int
	MaxSeqBonus      int
	LengthPenalty    int
	ExactBonus       int
	PrefixBonus      int
}

var DefaultConfig = Config{
	BaseScore:        10,
	BoundaryBonus:    15,
	ConsecutiveBonus: 5,
	MaxSeqBonus:      10,
	LengthPenalty:    1,
	ExactBonus:       100000,
	PrefixBonus:      50000,
}

func isBoundary(r rune) bool {
	return r == '_' || r == '.' || r == '-' || r == '/'
}

/* ---------- Preprocessing ---------- */

func Preprocess(data []string) []Entry {
	out := make([]Entry, len(data))
	for i, s := range data {
		ls := strings.ToLower(s)
		out[i] = Entry{
			Text:      s,
			LowerText: ls,
			Runes:     []rune(s),
			Lower:     []rune(ls),
		}
	}
	return out
}

/* ---------- Core Scoring ---------- */

func scoreRunes(queryLower []rune, e Entry, cfg Config) int {
	qi := 0
	score := 0
	consecutive := 0
	maxConsecutive := 0
	lastMatch := -1

	for ti := 0; ti < len(e.Lower); ti++ {
		if len(e.Lower)-ti < len(queryLower)-qi {
			break
		}

		if e.Lower[ti] == queryLower[qi] {
			current := cfg.BaseScore

			if ti == 0 || isBoundary(e.Lower[ti-1]) {
				current += cfg.BoundaryBonus
			}
			if lastMatch+1 == ti {
				consecutive++
				current += consecutive * cfg.ConsecutiveBonus
			} else {
				consecutive = 1
			}
			if unicode.IsUpper(e.Runes[ti]) {
				current += 3
			}

			score += current
			lastMatch = ti
			qi++

			if consecutive > maxConsecutive {
				maxConsecutive = consecutive
			}

			if qi == len(queryLower) {
				score += maxConsecutive * cfg.MaxSeqBonus
				score -= (len(e.Lower) - len(queryLower)) * cfg.LengthPenalty
				return score
			}
		} else {
			consecutive = 0
			score--
		}
	}

	return -1
}

/* ---------- Heap ---------- */

type matchHeap []Match

func (h matchHeap) Len() int           { return len(h) }
func (h matchHeap) Less(i, j int) bool { return h[i].Score < h[j].Score }
func (h matchHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *matchHeap) Push(x any) {
	*h = append(*h, x.(Match))
}

func (h *matchHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

/* ---------- Public API ---------- */

func Rank(query string, data []string, limit int) ([]Match, int) {
	entries := Preprocess(data)
	return RankEntries(query, entries, limit, DefaultConfig)
}

func RankEntries(query string, entries []Entry, limit int, cfg Config) ([]Match, int) {
	if query == "" || limit <= 0 {
		return nil, 0
	}

	lq := strings.ToLower(query)
	qRunes := []rune(lq)

	h := &matchHeap{}
	heap.Init(h)

	total := 0

	for _, e := range entries {
		if e.LowerText == lq {
			total++
			heap.Push(h, Match{Item: e.Text, Score: cfg.ExactBonus})
			continue
		}

		if strings.HasPrefix(e.LowerText, lq) {
			total++
			heap.Push(h, Match{Item: e.Text, Score: cfg.PrefixBonus})
			continue
		}

		s := scoreRunes(qRunes, e, cfg)
		if s > 0 {
			total++
			if h.Len() < limit {
				heap.Push(h, Match{Item: e.Text, Score: s})
			} else if (*h)[0].Score < s {
				(*h)[0] = Match{Item: e.Text, Score: s}
				heap.Fix(h, 0)
			}
		}
	}

	results := make([]Match, h.Len())
	for i := len(results) - 1; i >= 0; i-- {
		results[i] = heap.Pop(h).(Match)
	}

	return results, total
}
