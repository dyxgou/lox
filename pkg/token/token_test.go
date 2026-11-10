package token

import "testing"

func TestIsIdentifier(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		isIdent bool
	}{
		{"Empty", "", false},
		{"Space", " ", false},
		{"SpaceSuffix", "foo ", false},
		{"Number", "123", false},
		{"Keyword", "func", false},
		{"Struct", "struct", false},
		{"Range", "range", false},

		{"LetterASCII", "hello", true},
		{"MixedASCII", "hello123", true},
		{"UppercaseKeyword", "Func", true},
		{"LettersUnicode", "páez", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsIdentifier(tt.in); got != tt.isIdent {
				t.Fatalf("IsIdentifier(%q) = %t, want %v", tt.in, got, tt.isIdent)
			}
		})
	}
}
