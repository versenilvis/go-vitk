// https://github.com/AlasdairF/Tokenize
package tokenizer

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

type Token struct {
	Text   string
	Start  int
	End    int
	IsWord bool
}

type Option func(*Tokenizer)

func WithDictionary(dict *Dictionary) Option {
	return func(t *Tokenizer) {
		t.dict = dict
	}
}

func WithoutNormalize() Option {
	return func(t *Tokenizer) {
		t.normalize = false
	}
}

type Tokenizer struct {
	dict      *Dictionary
	normalize bool
}

func New(opts ...Option) (*Tokenizer, error) {
	t := &Tokenizer{
		normalize: true,
	}

	for _, opt := range opts {
		opt(t)
	}

	if t.dict == nil {
		dict, err := NewDictionary()
		if err != nil {
			return nil, err
		}
		t.dict = dict
	}

	return t, nil
}

func (t *Tokenizer) Tokenize(text string) []Token {
	if text == "" {
		return nil
	}

	segments := t.extractSegments(text)
	if len(segments) == 0 {
		return nil
	}

	lowerSyllables := make([]string, len(segments))
	for i, seg := range segments {
		if seg.isPunct {
			lowerSyllables[i] = seg.text
		} else {
			lowerSyllables[i] = strings.ToLower(seg.text)
		}
	}

	var tokens []Token
	i := 0
	numSegments := len(segments)

	for i < numSegments {

		if segments[i].isPunct {
			i++
			continue
		}

		matched, length := t.dict.SearchLongestMatch(lowerSyllables, i)

		if length > 0 {
			tokenText := ""
			if t.normalize {
				tokenText = matched
			} else {
				if length == 1 {
					tokenText = segments[i].text
				} else {
					var b strings.Builder
					b.Grow(segments[i+length-1].end - segments[i].start)
					for j := 0; j < length; j++ {
						if j > 0 {
							b.WriteByte(' ')
						}
						b.WriteString(segments[i+j].text)
					}
					tokenText = b.String()
				}
			}

			tokens = append(tokens, Token{
				Text:   tokenText,
				Start:  segments[i].start,
				End:    segments[i+length-1].end,
				IsWord: true,
			})
			i += length
		} else {
			text := segments[i].text
			if t.normalize {
				text = lowerSyllables[i]
			}
			tokens = append(tokens, Token{
				Text:   text,
				Start:  segments[i].start,
				End:    segments[i].end,
				IsWord: false,
			})
			i++
		}
	}

	return tokens
}

func (t *Tokenizer) TokenizeToStrings(text string) []string {
	tokens := t.Tokenize(text)
	if len(tokens) == 0 {
		return nil
	}
	result := make([]string, len(tokens))
	for i, tok := range tokens {
		result[i] = tok.Text
	}
	return result
}

type segment struct {
	text    string
	start   int
	end     int
	isPunct bool
}

func (t *Tokenizer) extractSegments(text string) []segment {
	var segments []segment
	n := len(text)
	bytePos := 0

	for bytePos < n {
		r, width := utf8.DecodeRuneInString(text[bytePos:])

		if unicode.IsSpace(r) {
			bytePos += width
			continue
		}

		startByte := bytePos
		isPunct := unicode.IsPunct(r) || unicode.IsSymbol(r)

		if isPunct {

			for bytePos < n {
				r, width = utf8.DecodeRuneInString(text[bytePos:])
				if !unicode.IsPunct(r) && !unicode.IsSymbol(r) {
					break
				}
				bytePos += width
			}
			segments = append(segments, segment{
				text:    text[startByte:bytePos],
				start:   startByte,
				end:     bytePos,
				isPunct: true,
			})
		} else {

			for bytePos < n {
				r, width = utf8.DecodeRuneInString(text[bytePos:])
				if unicode.IsSpace(r) || unicode.IsPunct(r) || unicode.IsSymbol(r) {
					break
				}
				bytePos += width
			}
			segments = append(segments, segment{
				text:    text[startByte:bytePos],
				start:   startByte,
				end:     bytePos,
				isPunct: false,
			})
		}
	}

	return segments
}

func (t *Tokenizer) Dict() *Dictionary {
	return t.dict
}
