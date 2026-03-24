/*
Package: go-vitk - Vietnamese tokenizer
Author: @versenilvis (Github)

vitk cung cấp khả năng xử lý văn bản tiếng Việt được tối ưu hóa cho công cụ tìm kiếm

vitk kết hợp tokenizer, loại bỏ stopwords và chuẩn hóa văn bản
vào một thư viện duy nhất được thiết kế cho các ứng dụng tìm kiếm tiếng Việt

Cách dùng:

token := vitk.Tokenize("Hà Nội là thủ đô của Việt Nam")
> ["hà nội", "là", "thủ đô", "của", "việt nam"]

clean := vitk.TokenizeAndClean("Hà Nội là thủ đô của Việt Nam")
> ["hà nội", "thủ đô", "Việt Nam"]

text := vitk.Normalize("TP.HCM ko co gi dc")
> "thành phố hồ chí minh không có gì được"

result := vitk.ForSearch("Tốc độ truyền thông tin ngày càng cao!!!")
> chuẩn hóa -> phân tách từ -> loại bỏ từ dừng -> làm sạch từ
*/
package vitk

import (
	"sync"

	"github.com/versenilvis/go-vitk/normalize"
	"github.com/versenilvis/go-vitk/search"
	"github.com/versenilvis/go-vitk/stopwords"
	"github.com/versenilvis/go-vitk/tokenizer"
)

var (
	defaultTokenizer  *tokenizer.Tokenizer
	defaultNormalizer *normalize.Normalizer
	defaultStopwords  *stopwords.Filter
	defaultPipeline   *search.Pipeline
	initOnce          sync.Once
	initErr           error
)

func ensureInit() error {
	initOnce.Do(func() {
		defaultTokenizer, initErr = tokenizer.New()
		if initErr != nil {
			return
		}

		defaultStopwords, initErr = stopwords.New()
		if initErr != nil {
			return
		}

		defaultNormalizer = normalize.New()

		defaultPipeline, initErr = search.NewPipeline()
	})
	return initErr
}

func Tokenize(text string) []string {
	if err := ensureInit(); err != nil {
		return nil
	}
	return defaultTokenizer.TokenizeToStrings(text)
}

func TokenizeDetailed(text string) []tokenizer.Token {
	if err := ensureInit(); err != nil {
		return nil
	}
	return defaultTokenizer.Tokenize(text)
}

func TokenizeAndClean(text string) []string {
	if err := ensureInit(); err != nil {
		return nil
	}
	tokens := defaultTokenizer.TokenizeToStrings(text)
	return defaultStopwords.Remove(tokens)
}

func Normalize(text string) string {
	if err := ensureInit(); err != nil {
		return text
	}
	return defaultNormalizer.Normalize(text)
}

func IsStopWord(word string) bool {
	if err := ensureInit(); err != nil {
		return false
	}
	return defaultStopwords.IsStopWord(word)
}

func ForSearch(text string) []string {
	if err := ensureInit(); err != nil {
		return nil
	}
	return defaultPipeline.ProcessTokens(text)
}

func ForSearchString(text string) string {
	if err := ensureInit(); err != nil {
		return ""
	}
	return defaultPipeline.ProcessDocument(text)
}

func NewTokenizer(opts ...tokenizer.Option) (*tokenizer.Tokenizer, error) {
	return tokenizer.New(opts...)
}

func NewNormalizer(opts ...normalize.Option) *normalize.Normalizer {
	return normalize.New(opts...)
}

func NewStopwords() (*stopwords.Filter, error) {
	return stopwords.New()
}

func NewPipeline() (*search.Pipeline, error) {
	return search.NewPipeline()
}
