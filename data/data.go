package data

import "embed"

//go:embed dictionary.txt
var DictionaryFS embed.FS

//go:embed stopwords.txt
var StopwordsFS embed.FS
