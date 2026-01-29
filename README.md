# fuzzygo

**fuzzygo** is a small, fast Go library for fuzzy string matching and ranking.
It is Unicode-safe, typo-tolerant, and designed for searching identifiers like APIs, symbols, and commands.

---

## Install

```bash
go get github.com/ankit-chaubey/fuzzygo
```

---

## Basic Usage

```go
import "github.com/ankit-chaubey/fuzzygo/fuzzy"

results, total := fuzzy.Rank("mesages.getMesage", data, 5)
```

* `results` → top ranked matches (best first)
* `total` → total number of matched items

---

## Scoring

```go
score := fuzzy.Score(query, target)
```

* Returns a positive score for a match
* Returns `-1` if the target does not match the query

---

## TL / JSON Helper (optional)

```go
methods, _ := fuzzy.LoadTLMethods("output.json")
results, total := fuzzy.Rank("sndmsg", methods, 10)
```

---

## License

MIT © Ankit Chaubey
