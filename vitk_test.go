package vitk_test

import (
	"reflect"
	"testing"

	"github.com/versenilvis/go-vitk"
	"github.com/versenilvis/go-vitk/tokenizer"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"Hà Nội là thủ đô của Việt Nam", []string{"hà nội", "là", "thủ đô", "của", "việt nam"}},
		{"Học sinh đi học tại trường", []string{"học sinh", "đi học", "tại", "trường"}},
		{"Tốc độ truyền thông tin ngày càng cao", []string{"tốc độ", "truyền thông", "tin", "ngày càng", "cao"}},
	}

	for _, tc := range tests {
		got := vitk.Tokenize(tc.input)
		if !reflect.DeepEqual(got, tc.expected) {
			t.Errorf("Tokenize(%q) = %#v; want %#v", tc.input, got, tc.expected)
		}
	}
}

func TestForSearch(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"Hà Nội là thủ đô của Việt Nam", []string{"ha noi", "thu do", "viet nam"}},
		{"TP.HCM là trung tâm kinh tế", []string{"thanh pho", "ho chi minh", "trung tam", "kinh te"}},
	}

	for _, tc := range tests {
		got := vitk.ForSearch(tc.input)

		if len(got) == 0 {
			t.Errorf("ForSearch(%q) returned no tokens", tc.input)
		} else {
			t.Logf("ForSearch(%q) = %v", tc.input, got)
		}
	}
}
func setupTokenizer(t *testing.T) *tokenizer.Tokenizer {
	t.Helper()
	tok, err := tokenizer.New()
	if err != nil {
		t.Fatalf("failed to create tokenizer: %v", err)
	}
	return tok
}

func TestTokenize_BasicVietnamese(t *testing.T) {
	tok := setupTokenizer(t)

	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "thủ đô",
			input:    "Hà Nội là thủ đô của Việt Nam",
			expected: []string{"hà nội", "là", "thủ đô", "của", "việt nam"},
		},
		{
			name:     "từ ghép",
			input:    "học sinh giỏi toán",
			expected: []string{"học sinh", "giỏi", "toán"},
		},
		{
			name:     "câu đơn giản",
			input:    "tôi đi học",
			expected: []string{"tôi", "đi học"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tok.TokenizeToStrings(tt.input)
			assertTokens(t, got, tt.expected)
		})
	}
}

func TestTokenize_EdgeCases(t *testing.T) {
	tok := setupTokenizer(t)

	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: nil,
		},
		{
			name:     "chỉ có dấu câu",
			input:    "!!!",
			expected: []string{},
		},
		{
			name:     "khoảng trắng thừa",
			input:    "học   sinh",
			expected: []string{"học sinh"},
		},
		{
			name:     "mix tiếng Anh và tiếng Việt",
			input:    "học AI rất thú vị",
			expected: []string{"học", "ai", "rất", "thú vị"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tok.TokenizeToStrings(tt.input)
			assertTokens(t, got, tt.expected)
		})
	}
}

func TestTokenize_PreserveCase(t *testing.T) {
	tok, err := tokenizer.New(tokenizer.WithoutNormalize())
	if err != nil {
		t.Fatal(err)
	}

	got := tok.TokenizeToStrings("Hà Nội")
	if len(got) == 0 {
		t.Fatal("expected tokens, got none")
	}
	if got[0] != "Hà Nội" {
		t.Errorf("expected 'Hà Nội', got %q", got[0])
	}
}

func TestTokenize_Position(t *testing.T) {
	tok := setupTokenizer(t)

	tokens := tok.Tokenize("học sinh")
	if len(tokens) == 0 {
		t.Fatal("expected tokens")
	}

	first := tokens[0]
	if first.Start != 0 {
		t.Errorf("expected Start=0, got %d", first.Start)
	}
	if !first.IsWord {
		t.Errorf("expected IsWord=true for 'học sinh'")
	}
}

func TestTokenize_UnknownWord(t *testing.T) {
	tok := setupTokenizer(t)

	tokens := tok.Tokenize("blockchain")
	if len(tokens) == 0 {
		t.Fatal("expected tokens")
	}
	if tokens[0].IsWord {
		t.Errorf("expected IsWord=false for unknown word 'blockchain'")
	}
}

func assertTokens(t *testing.T, got, expected []string) {
	t.Helper()

	filtered := got[:0]
	for _, s := range got {
		if s != "" {
			filtered = append(filtered, s)
		}
	}

	if expected == nil && filtered == nil {
		return
	}

	if len(filtered) != len(expected) {
		t.Errorf("length mismatch\ngot:      %v\nexpected: %v", filtered, expected)
		return
	}

	for i := range expected {
		if filtered[i] != expected[i] {
			t.Errorf("token[%d]: got %q, expected %q", i, filtered[i], expected[i])
		}
	}
}
