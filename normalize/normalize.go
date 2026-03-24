// https://github.com/AlasdairF/Tokenize
package normalize

import (
	"maps"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

type Option func(*Normalizer)

func WithAbbreviations(abbr map[string]string) Option {
	return func(n *Normalizer) {
		for k, v := range abbr {
			n.abbreviations[strings.ToLower(k)] = strings.ToLower(v)
		}
	}
}

func WithoutAbbreviations() Option {
	return func(n *Normalizer) { n.expandAbbr = false }
}

func WithoutPunctuation() Option {
	return func(n *Normalizer) { n.removePunct = false }
}

func WithoutLowercase() Option {
	return func(n *Normalizer) { n.lowercase = false }
}

func WithoutUnicodeNFC() Option {
	return func(n *Normalizer) { n.unicodeNFC = false }
}

func WithStripAccents() Option {
	return func(n *Normalizer) { n.stripAccents = true }
}

type Normalizer struct {
	abbreviations map[string]string
	removePunct   bool
	lowercase     bool
	unicodeNFC    bool
	expandAbbr    bool
	stripAccents  bool
}

func New(opts ...Option) *Normalizer {
	n := &Normalizer{
		abbreviations: make(map[string]string),
		removePunct:   true,
		lowercase:     true,
		unicodeNFC:    true,
		expandAbbr:    true,
		stripAccents:  false,
	}

	maps.Copy(n.abbreviations, defaultAbbreviations)
	for _, opt := range opts {
		opt(n)
	}
	return n
}


func (n *Normalizer) Normalize(text string) string {
	if text == "" {
		return ""
	}

	if n.unicodeNFC && !norm.NFC.IsNormalString(text) {
		text = norm.NFC.String(text)
	}

	var b strings.Builder
	b.Grow(len(text))

	lastWasSpace := true
	
	for _, r := range text {
		if n.stripAccents {
			r = getBaseChar(r)
		} else if n.lowercase {
			r = unicode.ToLower(r)
		}

		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			lastWasSpace = false
		} else {
			if !lastWasSpace {
				b.WriteRune(' ')
				lastWasSpace = true
			}
		}
	}

	res := strings.TrimSpace(b.String())

	if n.expandAbbr {
		res = n.expandAbbreviations(res)
		if n.stripAccents {
			res = StripAccents(res)
		}
	}

	return res
}

func getBaseChar(r rune) rune {
	r = unicode.ToLower(r)
	switch r {
	case 'á', 'à', 'ả', 'ã', 'ạ', 'ă', 'ắ', 'ằ', 'ẳ', 'ẵ', 'ặ', 'â', 'ấ', 'ầ', 'ẩ', 'ẫ', 'ậ':
		return 'a'
	case 'đ':
		return 'd'
	case 'é', 'è', 'ẻ', 'ẽ', 'ẹ', 'ê', 'ế', 'ề', 'ể', 'ễ', 'ệ':
		return 'e'
	case 'í', 'ì', 'ỉ', 'ĩ', 'ị':
		return 'i'
	case 'ó', 'ò', 'ỏ', 'õ', 'ọ', 'ô', 'ố', 'ồ', 'ổ', 'ỗ', 'ộ', 'ơ', 'ớ', 'ờ', 'ở', 'ỡ', 'ợ':
		return 'o'
	case 'ú', 'ù', 'ủ', 'ũ', 'ụ', 'ư', 'ứ', 'ừ', 'ử', 'ữ', 'ự':
		return 'u'
	case 'ý', 'ỳ', 'ỷ', 'ỹ', 'ỵ':
		return 'y'
	}
	return r
}

func (n *Normalizer) expandAbbreviations(text string) string {
	if text == "" {
		return ""
	}
	words := strings.Fields(text)
	changed := false
	for i, w := range words {
		if expanded, ok := n.abbreviations[w]; ok {
			words[i] = expanded
			changed = true
		}
	}
	if !changed {
		return text
	}
	return strings.Join(words, " ")
}

func NormalizeUnicode(s string) string {
	if norm.NFC.IsNormalString(s) {
		return s
	}
	return norm.NFC.String(s)
}

func StripAccents(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		b.WriteRune(getBaseChar(r))
	}
	return b.String()
}

func (n *Normalizer) AddAbbreviation(abbr, expansion string) {
	n.abbreviations[strings.ToLower(abbr)] = strings.ToLower(expansion)
}
