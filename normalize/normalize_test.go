package normalize

import "testing"

func TestNormalize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hà Nội", "ha noi"},
		{"Tiếng Việt", "tieng viet"},
		{"Thành phố Hồ Chí Minh", "thanh pho ho chi minh"},
		{"đường xá", "duong xa"},
		{"TP.HCM", "thanh pho ho chi minh"},
	}

	n := New(WithAbbreviations(defaultAbbreviations), WithStripAccents())
	for _, tc := range tests {
		got := n.Normalize(tc.input)
		if got != tc.expected {
			t.Errorf("Normalize(%q) = %q; want %q", tc.input, got, tc.expected)
		}
	}
}

func TestStripAccents(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hà Nội", "ha noi"},
		{"Tiếng Việt", "tieng viet"},
		{"đây là tiếng Việt có dấu", "day la tieng viet co dau"},
	}

	for _, tc := range tests {
		got := StripAccents(tc.input)
		if got != tc.expected {
			t.Errorf("StripAccents(%q) = %q; want %q", tc.input, got, tc.expected)
		}
	}
}
