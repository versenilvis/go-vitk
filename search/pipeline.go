package search

import (
	"strings"

	"github.com/versenilvis/go-vitk/normalize"
	"github.com/versenilvis/go-vitk/stopwords"
	"github.com/versenilvis/go-vitk/tokenizer"
)

type Pipeline struct {
	tokenizer  *tokenizer.Tokenizer
	normalizer *normalize.Normalizer
	stopwords  *stopwords.Filter
}

func NewPipeline() (*Pipeline, error) {
	tok, err := tokenizer.New()
	if err != nil {
		return nil, err
	}

	sw, err := stopwords.New()
	if err != nil {
		return nil, err
	}

	return &Pipeline{
		tokenizer:  tok,
		normalizer: normalize.New(),
		stopwords:  sw,
	}, nil
}

func (p *Pipeline) ProcessDocument(text string) string {
	normalized := p.normalizer.Normalize(text)
	tokens := p.tokenizer.TokenizeToStrings(normalized)
	cleaned := p.stopwords.Remove(tokens)
	return strings.Join(cleaned, " ")
}

func (p *Pipeline) ProcessQuery(query string) string {
	normalized := p.normalizer.Normalize(query)
	tokens := p.tokenizer.TokenizeToStrings(normalized)
	cleaned := p.stopwords.Remove(tokens)
	return strings.Join(cleaned, " ")
}

func (p *Pipeline) ProcessTokens(text string) []string {
	normalized := p.normalizer.Normalize(text)
	tokens := p.tokenizer.TokenizeToStrings(normalized)
	return p.stopwords.Remove(tokens)
}

func (p *Pipeline) Tokenizer() *tokenizer.Tokenizer {
	return p.tokenizer
}

func (p *Pipeline) Normalizer() *normalize.Normalizer {
	return p.normalizer
}

func (p *Pipeline) Stopwords() *stopwords.Filter {
	return p.stopwords
}
