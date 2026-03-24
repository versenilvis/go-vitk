package main

import (
	"fmt"
	"github.com/versenilvis/go-vitk/tokenizer"
)

func main() {
	// 1. Tạo Tokenizer mới với những từ riêng biệt của bạn (VD: Mã sản phẩm, tên hãng)
	// vitk nạp 75k từ chuẩn, nhưng bạn có thể thêm bất cứ điều gì
	customTokenizer, _ := tokenizer.New(
		tokenizer.WithExtraWords([]string{"iphone 16 pro max", "macbook m4", "omniseek"}),
	)

	text := "Hôm nay tôi mua iPhone 16 Pro Max và tải Omniseek về dùng!!!"

	// 2. Tách từ bằng bộ não mới của bạn
	tokens := customTokenizer.TokenizeToStrings(text)
	fmt.Printf("Custom Tokens: %v\n", tokens)
	// > ["iphone 16 pro max", "omniseek", ...]

	// Thấy đấy, vitk sẽ nhận diện "iphone 16 pro max" là 1 token duy nhất thay vì 4 từ đơn!
}
