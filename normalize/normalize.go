package normalize

import (
	"maps"
	"strings"
	"unicode"
	"unicode/utf8"

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

	if isASCII(text) {
		if n.lowercase {
			text = strings.ToLower(text)
		}
		if n.removePunct {
			text = n.cleanPunctuation(text)
		}
		if n.expandAbbr {
			text = n.expandAbbreviations(text)
		}
		if n.stripAccents {
			text = stripAccents(text)
		}
		return strings.TrimSpace(collapseWhitespace(text))
	}

	if n.unicodeNFC && !norm.NFC.IsNormalString(text) {
		text = norm.NFC.String(text)
	}

	if n.lowercase {
		text = strings.ToLower(text)
	}

	if n.removePunct {
		text = n.cleanPunctuation(text)
	}

	if n.expandAbbr {
		text = n.expandAbbreviations(text)
	}

	if n.stripAccents {
		text = stripAccents(text)
	}

	return strings.TrimSpace(collapseWhitespace(text))
}

func NormalizeUnicode(s string) string {
	if norm.NFC.IsNormalString(s) {
		return s
	}
	return norm.NFC.String(s)
}

func StripAccents(s string) string {
	if isASCII(s) {
		return strings.ToLower(s)
	}
	s = NormalizeUnicode(s)
	return stripAccents(s)
}

func stripAccents(s string) string {
	buf := make([]byte, 0, len(s))
	for _, r := range s {
		r = unicode.ToLower(r)
		switch r {
		case 'á', 'à', 'ả', 'ã', 'ạ', 'ă', 'ắ', 'ằ', 'ẳ', 'ẵ', 'ặ', 'â', 'ấ', 'ầ', 'ẩ', 'ẫ', 'ậ':
			buf = append(buf, 'a')
		case 'đ':
			buf = append(buf, 'd')
		case 'é', 'è', 'ẻ', 'ẽ', 'ẹ', 'ê', 'ế', 'ề', 'ể', 'ễ', 'ệ':
			buf = append(buf, 'e')
		case 'í', 'ì', 'ỉ', 'ĩ', 'ị':
			buf = append(buf, 'i')
		case 'ó', 'ò', 'ỏ', 'õ', 'ọ', 'ô', 'ố', 'ồ', 'ổ', 'ỗ', 'ộ', 'ơ', 'ớ', 'ờ', 'ở', 'ỡ', 'ợ':
			buf = append(buf, 'o')
		case 'ú', 'ù', 'ủ', 'ũ', 'ụ', 'ư', 'ứ', 'ừ', 'ử', 'ữ', 'ự':
			buf = append(buf, 'u')
		case 'ý', 'ỳ', 'ỷ', 'ỹ', 'ỵ':
			buf = append(buf, 'y')
		default:
			if r < 128 {
				buf = append(buf, byte(r))
			} else if unicode.IsLetter(r) || unicode.IsDigit(r) {
				var tmp [4]byte
				n := utf8.EncodeRune(tmp[:], r)
				buf = append(buf, tmp[:n]...)
			}
		}
	}
	return string(buf)
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > 127 {
			return false
		}
	}
	return true
}

func (n *Normalizer) expandAbbreviations(text string) string {
	words := strings.Fields(text)
	for i, w := range words {
		cleaned := strings.TrimFunc(w, func(r rune) bool {
			return unicode.IsPunct(r) || unicode.IsSymbol(r)
		})
		if expanded, ok := n.abbreviations[strings.ToLower(cleaned)]; ok {
			words[i] = expanded
		}
	}
	return strings.Join(words, " ")
}

func (n *Normalizer) cleanPunctuation(text string) string {
	var b strings.Builder
	b.Grow(len(text))
	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			b.WriteRune(r)
		} else {
			b.WriteRune(' ')
		}
	}
	return b.String()
}

func collapseWhitespace(text string) string {
	var b strings.Builder
	b.Grow(len(text))
	prevSpace := false
	for _, r := range text {
		if unicode.IsSpace(r) {
			if !prevSpace {
				b.WriteRune(' ')
			}
			prevSpace = true
		} else {
			b.WriteRune(r)
			prevSpace = false
		}
	}
	return b.String()
}

func (n *Normalizer) AddAbbreviation(abbr, expansion string) {
	n.abbreviations[strings.ToLower(abbr)] = strings.ToLower(expansion)
}
