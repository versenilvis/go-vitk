package stopwords

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/versenilvis/go-vitk/data"
)

type Filter struct {
	words map[string]struct{}
	re    *regexp.Regexp
}

type Option func(*Filter)

func WithExtraWords(words []string) Option {
	return func(f *Filter) {
		f.Add(words...)
	}
}

func New(opts ...Option) (*Filter, error) {
	f := &Filter{words: make(map[string]struct{})}

	file, err := data.StopwordsFS.Open("stopwords.txt")
	if err != nil {
		return nil, fmt.Errorf("open embedded stopwords: %w", err)
	}
	defer file.Close()

	if err := f.loadFromReader(file); err != nil {
		return nil, fmt.Errorf("load embedded stopwords: %w", err)
	}

	for _, opt := range opts {
		opt(f)
	}

	f.rebuildRegexp()
	return f, nil
}

func NewFromReader(r io.Reader) (*Filter, error) {
	f := &Filter{words: make(map[string]struct{})}
	if err := f.loadFromReader(r); err != nil {
		return nil, err
	}
	f.rebuildRegexp()
	return f, nil
}

func NewEmpty() *Filter {
	return &Filter{words: make(map[string]struct{})}
}

func (f *Filter) loadFromReader(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word == "" {
			continue
		}
		f.words[strings.ToLower(word)] = struct{}{}
	}
	return scanner.Err()
}

func (f *Filter) rebuildRegexp() {
	if len(f.words) == 0 {
		f.re = nil
		return
	}
	var sb strings.Builder
	i := 0
	for word := range f.words {
		if i > 0 {
			sb.WriteByte('|')
		}
		sb.WriteString(`\Q`)
		sb.WriteString(word)
		sb.WriteString(`\E`)
		i++
	}
	// using case-insensitive (?i) and word boundary \b
	f.re = regexp.MustCompile(`(?i)\b(` + sb.String() + `)\b`)
}

func (f *Filter) IsStopWord(word string) bool {
	_, ok := f.words[strings.ToLower(word)]
	return ok
}

func (f *Filter) Remove(tokens []string) []string {
	result := make([]string, 0, len(tokens))
	for _, t := range tokens {
		if !f.IsStopWord(t) {
			result = append(result, t)
		}
	}
	return result
}

func (f *Filter) RemoveFromText(text string) string {
	if f.re == nil {
		return text
	}
	res := f.re.ReplaceAllString(text, "")
	return strings.Join(strings.Fields(res), " ")
}

func (f *Filter) Add(words ...string) {
	changed := false
	for _, w := range words {
		lower := strings.ToLower(w)
		if _, ok := f.words[lower]; !ok {
			f.words[lower] = struct{}{}
			changed = true
		}
	}
	if changed {
		f.rebuildRegexp()
	}
}

func (f *Filter) RemoveWord(words ...string) {
	changed := false
	for _, w := range words {
		lower := strings.ToLower(w)
		if _, ok := f.words[lower]; ok {
			delete(f.words, lower)
			changed = true
		}
	}
	if changed {
		f.rebuildRegexp()
	}
}

func (f *Filter) Words() []string {
	result := make([]string, 0, len(f.words))
	for w := range f.words {
		result = append(result, w)
	}
	return result
}

func (f *Filter) Size() int {
	return len(f.words)
}
