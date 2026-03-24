package main

import (
	"fmt"
	"strings"

	"github.com/versenilvis/go-vitk/normalize"
)

func main() {
	// Danh sách file thực tế trong hệ thống của bạn (viết hoa, viết thường, có dấu, không dấu xen kẽ)
	files := []string{
		"Báo cáo tài chính 2024.pdf",
		"Danh sách khách hàng VIP.xlsx",
		"Hợp đồng du lịch Hà Nội.docx",
		"ds_khach_hang_cu.csv",
		"hanoi-pho-co-dem-mua.jpg",
	}

	// 1. Bình thường hoá: Chuyển hết về chữ thường, bỏ dấu để so sánh "mềm" (fuzzy)
	// Đây chính là sức mạnh của Normalize trong vitk
	normalizer := normalize.New(normalize.WithStripAccents())

	// 2. Chế độ tìm kiếm không dấu "thông minh"
	// Người dùng gõ: "danh sach khach hn"
	query := "danh sach khach"
	queryNorm := normalizer.Normalize(query)

	fmt.Printf("Tìm kiếm file với từ khoá: '%s'\n", query)
	fmt.Println("----------------------------------------")

	for _, f := range files {
		// Chuẩn hoá tên file trước khi check
		fileNameNorm := normalizer.Normalize(f)

		// Simple search: Kiểm tra xem query có nằm trong tên file đã chuẩn hoá ko
		if strings.Contains(fileNameNorm, queryNorm) {
			fmt.Printf("Found: [ %s ] -> (%s)\n", f, fileNameNorm)
		}
	}
}
