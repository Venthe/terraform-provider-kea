package kea_test

import (
	"bytes"
	"testing"

	"github.com/venthe/terraform-provider-kea/kea"
)

func TestParseHexID(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    kea.HexID
		wantErr bool
	}{
		{"empty string", "", nil, false},

		{"plain even", "aabbcc", kea.HexID{0xaa, 0xbb, 0xcc}, false},
		{"plain odd (padding)", "abc", kea.HexID{0x0a, 0xbc}, false},
		{"uppercase", "AABBCC", kea.HexID{0xaa, 0xbb, 0xcc}, false},

		{"prefix even", "0xaabbcc", kea.HexID{0xaa, 0xbb, 0xcc}, false},
		{"prefix odd (padding)", "0xabc", kea.HexID{0x0a, 0xbc}, false},

		{"colons", "aa:bb:cc", kea.HexID{0xaa, 0xbb, 0xcc}, false},
		{"colons padding", "a:b:c", kea.HexID{0x0a, 0x0b, 0x0c}, false},
		{"spaces", "aa bb cc", kea.HexID{0xaa, 0xbb, 0xcc}, false},
		{"spaces padding", "a b c", kea.HexID{0x0a, 0x0b, 0x0c}, false},

		{"invalid hex chars", "zz:yy", nil, true},
		{"invalid hex chars prefix", "0xzz", nil, true},
		{"invalid prefix", "0XAABBCC", nil, true},

		{"forbidden: prefix + colons", "0xaa:bb:cc", nil, true},
		{"forbidden: prefix + spaces", "0xaa bb cc", nil, true},
		{"forbidden: mixed separators", "aa:bb cc", nil, true},

		{"invalid segment length", "aaa:bbb", nil, true},
		{"trailing separator", "aa:bb:", nil, true},
		{"leading separator", ":aa:bb", nil, true},
		{"double separator", "aa::bb", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := kea.ParseHexID(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseHexID(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if !bytes.Equal(got, tt.want) {
					t.Errorf("ParseHexID(%q) = %x, want %x", tt.input, got, tt.want)
				}
			}
		})
	}
}
