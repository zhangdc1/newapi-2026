package service

import (
	"embed"
	"encoding/json"
	"math"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/dlclark/regexp2"
)

//go:embed tokenizer_assets/deepseek/tokenizer.json
var deepSeekTokenizerFS embed.FS

const (
	deepSeekTokenizerPath = "tokenizer_assets/deepseek/tokenizer.json"
	deepSeekSplitNumbers  = `\p{N}{1,3}`
	deepSeekSplitCJK      = `[\u4E00-\u9FFF\u3040-\u309F\u30A0-\u30FF]+`
	deepSeekSplitText     = `[!"#$%&'()*+,\-./:;<=>?@\[\]\\^_` + "`" + `{|}~][A-Za-z]+|[^\r\n\p{L}\p{P}\p{S}]?[\p{L}\p{M}]+| ?[\p{P}\p{S}]+[\r\n]*|\s*[\r\n]+|\s+(?!\S)|\s+`
)

var (
	deepSeekTokenizerOnce sync.Once
	deepSeekTokenizer     *deepSeekBPETokenizer
	deepSeekTokenizerErr  error
)

type deepSeekBPETokenizer struct {
	vocab     map[string]uint
	splitters []*regexp2.Regexp
	byteChars [256]string
}

type deepSeekTokenizerConfig struct {
	Model struct {
		Vocab map[string]uint `json:"vocab"`
	} `json:"model"`
}

func IsDeepSeekModel(model string) bool {
	model = strings.ToLower(strings.TrimSpace(model))
	model = strings.TrimPrefix(model, "deepseek/")
	return strings.HasPrefix(model, "deepseek") || strings.Contains(model, "/deepseek-")
}

func CountDeepSeekTextToken(text string) (int, bool) {
	if text == "" {
		return 0, true
	}
	tk, err := getDeepSeekTokenizer()
	if err != nil {
		return 0, false
	}
	return tk.Count(text), true
}

func getDeepSeekTokenizer() (*deepSeekBPETokenizer, error) {
	deepSeekTokenizerOnce.Do(func() {
		deepSeekTokenizer, deepSeekTokenizerErr = newDeepSeekBPETokenizer()
	})
	return deepSeekTokenizer, deepSeekTokenizerErr
}

func newDeepSeekBPETokenizer() (*deepSeekBPETokenizer, error) {
	data, err := deepSeekTokenizerFS.ReadFile(deepSeekTokenizerPath)
	if err != nil {
		return nil, err
	}
	var cfg deepSeekTokenizerConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	tk := &deepSeekBPETokenizer{
		vocab: cfg.Model.Vocab,
		splitters: []*regexp2.Regexp{
			regexp2.MustCompile(deepSeekSplitNumbers, regexp2.RE2),
			regexp2.MustCompile(deepSeekSplitCJK, regexp2.RE2),
			regexp2.MustCompile(deepSeekSplitText, regexp2.RE2),
		},
		byteChars: deepSeekByteChars(),
	}
	return tk, nil
}

func (t *deepSeekBPETokenizer) Count(text string) int {
	count := 0
	for _, piece := range t.preTokenize(text) {
		byteLevel := t.byteLevel(piece)
		if byteLevel == "" {
			continue
		}
		if _, ok := t.vocab[byteLevel]; ok {
			count++
			continue
		}
		count += t.countBPE(byteLevel)
	}
	return count
}

func (t *deepSeekBPETokenizer) preTokenize(text string) []string {
	pieces := []string{text}
	for _, splitter := range t.splitters {
		var next []string
		for _, piece := range pieces {
			next = append(next, splitRegexp2(piece, splitter)...)
		}
		pieces = next
	}
	return pieces
}

func splitRegexp2(text string, re *regexp2.Regexp) []string {
	if text == "" {
		return nil
	}
	runes := []rune(text)
	var pieces []string
	cursor := 0
	match, err := re.FindStringMatch(text)
	for err == nil && match != nil {
		start := match.Index
		end := match.Index + match.Length
		if start > cursor {
			pieces = append(pieces, string(runes[cursor:start]))
		}
		if end > start {
			pieces = append(pieces, string(runes[start:end]))
		}
		cursor = end
		match, err = re.FindNextMatch(match)
	}
	if cursor < len(runes) {
		pieces = append(pieces, string(runes[cursor:]))
	}
	return pieces
}

func (t *deepSeekBPETokenizer) byteLevel(text string) string {
	var b strings.Builder
	for i := 0; i < len(text); i++ {
		b.WriteString(t.byteChars[text[i]])
	}
	return b.String()
}

func (t *deepSeekBPETokenizer) countBPE(piece string) int {
	offsets := tokenOffsets(piece)
	if len(offsets) <= 2 {
		return 1
	}

	type part struct {
		offset int
		rank   uint
	}
	parts := make([]part, len(offsets))
	for i, offset := range offsets {
		parts[i] = part{offset: offset, rank: math.MaxUint}
	}

	getRank := func(index, skip int) uint {
		if index+skip+2 < len(parts) {
			token := piece[parts[index].offset:parts[index+skip+2].offset]
			if rank, ok := t.vocab[token]; ok {
				return rank
			}
		}
		return math.MaxUint
	}

	for i := 0; i < len(parts)-2; i++ {
		parts[i].rank = getRank(i, 0)
	}

	for {
		minRank := uint(math.MaxUint)
		minIndex := -1
		for i, p := range parts[:len(parts)-1] {
			if p.rank < minRank {
				minRank = p.rank
				minIndex = i
			}
		}
		if minIndex < 0 || minRank == math.MaxUint {
			break
		}

		parts[minIndex].rank = getRank(minIndex, 1)
		if minIndex > 0 {
			parts[minIndex-1].rank = getRank(minIndex-1, 1)
		}
		parts = append(parts[:minIndex+1], parts[minIndex+2:]...)
	}

	if len(parts) > 0 {
		return len(parts) - 1
	}
	return 0
}

func tokenOffsets(s string) []int {
	offsets := make([]int, 0, utf8.RuneCountInString(s)+1)
	for i := range s {
		offsets = append(offsets, i)
	}
	offsets = append(offsets, len(s))
	return offsets
}

func deepSeekByteChars() [256]string {
	var out [256]string
	n := 0
	for i := 256; i <= 288; i++ {
		out[n] = string(rune(i))
		n++
	}
	for i := 33; i <= 126; i++ {
		out[n] = string(rune(i))
		n++
	}
	for i := 289; i <= 322; i++ {
		out[n] = string(rune(i))
		n++
	}
	for i := 161; i <= 172; i++ {
		out[n] = string(rune(i))
		n++
	}
	if n == 173 {
		out[n] = string(rune(323))
		n++
	}
	for i := 174; i <= 255; i++ {
		out[n] = string(rune(i))
		n++
	}
	return out
}
