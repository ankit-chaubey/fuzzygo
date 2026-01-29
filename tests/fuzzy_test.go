package fuzzy

import (
	"fmt"
	"testing"
)

// github.com/ankit-chaubey/fuzzygo
// stress & speed tests // @ankit-chaubey

func TestScoreRunes_Basic(t *testing.T) {
	cfg := DefaultConfig
	target := "messages.sendMessage"

	entry := Entry{
		Text:      target,
		LowerText: "messages.sendmessage",
		Runes:     []rune(target),
		Lower:     []rune("messages.sendmessage"),
	}

	tests := []struct {
		query string
		ok    bool
	}{
		{"msg", true},
		{"sndmsg", true},
		{"sendmsg", true},
		{"sendMessage", false},
		{"xyz", false},
	}

	for _, tt := range tests {
		score := scoreRunes([]rune(tt.query), entry, cfg)
		if tt.ok && score <= 0 {
			t.Errorf("expected match for %q, got %d", tt.query, score)
		}
		if !tt.ok && score > 0 {
			t.Errorf("expected no match for %q, got %d", tt.query, score)
		}
	}
}

func TestRankEntries_PrefixAndExact(t *testing.T) {
	data := []string{
		"messages.sendMessage",
		"messages.sendReaction",
		"users.sendMessage",
	}

	entries := Preprocess(data)
	results, _ := RankEntries("messages.sendMessage", entries, 3, DefaultConfig)

	if len(results) == 0 {
		t.Fatal("no results")
	}
	if results[0].Item != "messages.sendMessage" {
		t.Fatalf("unexpected top result: %s", results[0].Item)
	}
}

func TestLargeDataset_100k(t *testing.T) {
	const N = 100_000

	data := make([]string, N)
	for i := 0; i < N; i++ {
		data[i] = fmt.Sprintf("api.v1.service.method_%d", i)
	}

	entries := Preprocess(data)
	results, total := RankEntries("service.method_99", entries, 10, DefaultConfig)

	if total == 0 {
		t.Fatal("no matches")
	}
	if len(results) > 10 {
		t.Fatalf("limit exceeded: %d", len(results))
	}
}

func TestLargeDataset_1Million(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	const N = 1_000_000

	data := make([]string, N)
	for i := 0; i < N; i++ {
		data[i] = fmt.Sprintf("api.v2.user.profile.fetch_%d", i)
	}

	entries := Preprocess(data)
	results, total := RankEntries("user.profile.fetch_9999", entries, 10, DefaultConfig)

	if total == 0 {
		t.Fatal("no matches")
	}
	if len(results) > 10 {
		t.Fatalf("limit exceeded: %d", len(results))
	}
}

func BenchmarkRankEntries_100k(b *testing.B) {
	const N = 100_000

	data := make([]string, N)
	for i := 0; i < N; i++ {
		data[i] = fmt.Sprintf("api.v1.resource.fetch_data_%d", i)
	}

	entries := Preprocess(data)
	query := "res.fetch_999"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = RankEntries(query, entries, 10, DefaultConfig)
	}
}

func BenchmarkRankEntries_1Million(b *testing.B) {
	if testing.Short() {
		b.Skip()
	}

	const N = 1_000_000

	data := make([]string, N)
	for i := 0; i < N; i++ {
		data[i] = fmt.Sprintf("api.v3.notification.send_bulk_%d", i)
	}

	entries := Preprocess(data)
	query := "notif.send_bulk_9999"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = RankEntries(query, entries, 10, DefaultConfig)
	}
}
