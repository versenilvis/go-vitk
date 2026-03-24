package tokenizer

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/versenilvis/go-vitk/data"
)

type Dictionary struct {
	trie  Tree[string, string]
	words map[string]struct{}
}

func NewDictionary() (*Dictionary, error) {
	d := &Dictionary{
		words: make(map[string]struct{}),
	}

	f, err := data.DictionaryFS.Open("dictionary.txt")
	if err != nil {
		return nil, fmt.Errorf("open embedded dictionary: %w", err)
	}
	defer f.Close()

	if err := d.loadFromReader(f); err != nil {
		return nil, fmt.Errorf("load embedded dictionary: %w", err)
	}

	d.rebuild()
	return d, nil
}

func NewDictionaryFromReader(r io.Reader) (*Dictionary, error) {
	d := &Dictionary{
		words: make(map[string]struct{}),
	}
	if err := d.loadFromReader(r); err != nil {
		return nil, err
	}
	d.rebuild()
	return d, nil
}

func NewEmptyDictionary() *Dictionary {
	return &Dictionary{
		words: make(map[string]struct{}),
	}
}

func (d *Dictionary) loadFromReader(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word == "" {
			continue
		}
		d.words[strings.ToLower(word)] = struct{}{}
	}
	return scanner.Err()
}

func (d *Dictionary) rebuild() {
	if len(d.words) == 0 {
		return
	}

	keys := make([][]string, 0, len(d.words))
	values := make([]string, 0, len(d.words))

	seen := make(map[string]struct{})

	for w := range d.words {
		syllables := splitSyllables(w)
		sig := strings.Join(syllables, "|")
		if _, ok := seen[sig]; ok {
			continue
		}
		seen[sig] = struct{}{}

		keys = append(keys, syllables)
		values = append(values, w)
	}

	d.trie = NewTrie(keys, values)
}

func (d *Dictionary) AddWord(word string) {
	d.words[strings.ToLower(word)] = struct{}{}
	d.rebuild()
}

func (d *Dictionary) AddWords(words ...string) {
	for _, w := range words {
		d.words[strings.ToLower(w)] = struct{}{}
	}
	d.rebuild()
}

func (d *Dictionary) Contains(word string) bool {
	_, ok := d.words[strings.ToLower(word)]
	return ok
}

func (d *Dictionary) SearchLongestMatch(syllables []string, start int) (string, int) {
	if d.trie == nil {
		return "", 0
	}
	return d.trie.SearchLongestMatch(syllables, start)
}

func (d *Dictionary) Size() int {
	return len(d.words)
}

func splitSyllables(text string) []string {
	var result []string
	for _, field := range strings.Fields(text) {
		cleaned := strings.TrimFunc(field, func(r rune) bool {
			return unicode.IsPunct(r) || unicode.IsSymbol(r)
		})
		if cleaned != "" {
			result = append(result, cleaned)
		}
	}
	return result
}
