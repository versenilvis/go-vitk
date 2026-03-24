
# go-vitk

**go-vitk** là thư viện xử lý tiếng Việt được tối ưu cho bài toán tìm kiếm (Meilisearch, Elasticsearch, file finder, ...), sử dụng cấu trúc Immutable Array-based Trie và Single-pass Normalization

> [!WARNING]
> **Đây không phải thư viện tokenizer cho LLM/AI**

## Tính năng nổi bật
- **Siêu nhanh:** Tốc độ tách từ cấp độ Nano giây (~8.6 µs cho chuỗi dài)
- **Tiết kiệm RAM:** Sử dụng cấu trúc mảng flat thay vì Map, giảm thiểu GC pressure
- **Dữ liệu thực tế:** Từ điển hơn 75000 từ (bao gồm 63 tỉnh thành), và hơn 2000 stopwords
- **Chuẩn hóa thông minh:** Unicode NFC, chuyển chữ thường, xử lý teencode/slang chỉ trong 1 lần duyệt (single-pass)
- **Sẵn sàng cho Search:** Tích hợp bộ sinh cấu hình Meilisearch/Elasticsearch với Sonic JSON, ngoài ra còn phù hợp cho file finder

## Cài đặt
```bash
go get github.com/versenilvis/go-vitk
```

## Cách dùng
```go
import "github.com/versenilvis/go-vitk"

func main() {
    // 1. Tách từ ghép (compound words)
    tokens := vitk.Tokenize("Hà Nội là thủ đô của Việt Nam")
    // > ["hà nội", "là", "thủ đô", "của", "việt nam"]

    // 2. Chế độ lọc sạch cho Search (Chuẩn hoá -> Tách từ -> Bỏ stopwords)
    clean := vitk.ForSearch("TP.HCM mùa này ko co gi vui hết!!!")
    // > ["thành phố hồ chí minh", "mùa", "vui"]

    // 3. Chuẩn hóa chuỗi raw (Unicode NFC, Teencode, Lowercase)
    norm := vitk.Normalize("Ko bit lam gi lun ak")
    // > "không biết làm gì luôn à"
}
```
---
> [!IMPORTANT]
> **Nếu bạn có thắc mắc, xin hãy liên hệ qua email versedev.store@proton.me**  
> **Thư viện sử dụng [0BSD LICENSE](./LICENSE), đồng nghĩa với việc bạn có thể sử dụng tự do cũng như có thể xoá, sửa, hay làm bất cứ điều gì bạn muốn**  
> **Tôi sẽ không chịu trách nhiệm trước bất cứ vấn đề nào xảy ra trong hệ thống của bạn**
