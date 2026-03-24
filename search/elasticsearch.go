package search

import (
	"github.com/bytedance/sonic"
)

type ElasticsearchAnalyzer struct {
	Name string
}

func DefaultElasticsearchAnalyzer() *ElasticsearchAnalyzer {
	return &ElasticsearchAnalyzer{
		Name: "vitk_vietnamese",
	}
}

func (a *ElasticsearchAnalyzer) IndexSettings() ([]byte, error) {
	settings := map[string]interface{}{
		"settings": map[string]interface{}{
			"analysis": map[string]interface{}{
				"analyzer": map[string]interface{}{
					a.Name: map[string]interface{}{
						"type":      "custom",
						"tokenizer": "whitespace",
						"filter":    []string{"lowercase"},
					},
				},
			},
		},
	}
	return sonic.ConfigDefault.MarshalIndent(settings, "", "  ")
}

func (a *ElasticsearchAnalyzer) FieldMapping(fieldName string) ([]byte, error) {
	mapping := map[string]interface{}{
		"properties": map[string]interface{}{
			fieldName: map[string]interface{}{
				"type":     "text",
				"analyzer": a.Name,
			},
		},
	}
	return sonic.ConfigDefault.MarshalIndent(mapping, "", "  ")
}
