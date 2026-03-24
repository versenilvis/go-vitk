package main

import (
	"fmt"
	"github.com/versenilvis/go-vitk"
)

func main() {
	text := "Hà Nội là thủ đô của Việt Nam, ko co gi vui hết!!!"

	// 1. Tách từ ghép (Tokenize)
	tokens := vitk.Tokenize(text)
	fmt.Printf("Tokenize: %v\n", tokens)

	// 2. Chuẩn hóa chuỗi (Normalize) - Tự động sửa teencode
	norm := vitk.Normalize(text)
	fmt.Printf("Normalize: %s\n", norm)

	// 3. Quy trình cho Search (Pipeline ForSearch)
	// Tự động Chuẩn hoá -> Tách từ -> Bỏ từ dừng
	searchTokens := vitk.ForSearch(text)
	fmt.Printf("ForSearch: %v\n", searchTokens)
}
