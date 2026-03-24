package main

import (
	"fmt"
	"github.com/versenilvis/go-vitk/search"
)

func main() {
	// 1. Tạo Pipeline tìm kiếm (Chuẩn hoá + Tách từ + Stopwords)
	pipeline, _ := search.NewPipeline()

	// 2. Chuẩn bị dữ liệu cào về (trước khi bắn lên Meilisearch)
	rawText := "TP.HCM là trung tâm kinh tế sầm uất nhất VN!!!"
	indexedText := search.PrepareDocument(pipeline, rawText)
	fmt.Printf("Index Text: %s\n", indexedText)
	// > "thành phố hồ chí minh trung tâm kinh tế sầm uất"

	// 3. Xử lý Search Query của người dùng
	userQuery := "tại sài gòn có gì vui?"
	searchQuery := search.PrepareQuery(pipeline, userQuery)
	fmt.Printf("Search Query: %s\n", searchQuery)
	// > "sài gòn vui"

	// 4. Sinh JSON cấu hình cho Meilisearch dùng Sonic
	settings := search.DefaultMeilisearchSettings()
	settingsJSON, _ := settings.ToJSON()
	fmt.Printf("Meili Settings JSON:\n%s\n", string(settingsJSON))
}
