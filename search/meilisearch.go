package search

import (
	"github.com/bytedance/sonic"
)

type MeilisearchSettings struct {
	SearchableAttributes []string

	SeparatorTokens []string
}

func DefaultMeilisearchSettings() *MeilisearchSettings {
	return &MeilisearchSettings{
		SeparatorTokens: []string{"|"},
	}
}

func (s *MeilisearchSettings) ToJSON() ([]byte, error) {
	settings := map[string]interface{}{
		"separatorTokens": s.SeparatorTokens,
	}
	if len(s.SearchableAttributes) > 0 {
		settings["searchableAttributes"] = s.SearchableAttributes
	}
	return sonic.ConfigDefault.MarshalIndent(settings, "", "  ")
}

func PrepareDocument(pipeline *Pipeline, text string) string {
	return pipeline.ProcessDocument(text)
}

func PrepareQuery(pipeline *Pipeline, query string) string {
	return pipeline.ProcessQuery(query)
}
