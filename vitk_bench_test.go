package vitk_test

import (
	"testing"

	"github.com/versenilvis/go-vitk"
)

var (
	shortText = "Hà Nội là thủ đô của Việt Nam"
	longText  = "Tốc độ truyền thông tin ngày càng cao, đặc biệt là trong bối cảnh cuộc cách mạng công nghiệp lần thứ tư đang diễn ra mạnh mẽ trên toàn thế giới. Việt Nam đang nỗ lực không ngừng để bắt kịp xu hướng này."
)

func BenchmarkTokenize_Short(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = vitk.Tokenize(shortText)
	}
}

func BenchmarkTokenize_Long(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = vitk.Tokenize(longText)
	}
}

func BenchmarkNormalize_Short(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = vitk.Normalize(shortText)
	}
}

func BenchmarkNormalize_Long(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = vitk.Normalize(longText)
	}
}

func BenchmarkForSearch_Short(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = vitk.ForSearch(shortText)
	}
}

func BenchmarkForSearch_Long(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = vitk.ForSearch(longText)
	}
}
