package ui

import (
	"strings"
	"testing"

	"github.com/nathanmazzapica/goto/internal/marker"
	"github.com/rivo/uniseg"
)

func TestUnicodeTruncation(t *testing.T) {
	tests := []struct {
		name, value string
		width       int
		head, tail  string
	}{
		{"ascii", "abcdef", 4, "abc…", "…def"},
		{"wide", "项目目录", 4, "项…", "…录"},
		{"combining", "e\u0301cole", 3, "e\u0301c…", "…le"},
		{"joined emoji", "👩🏽‍💻abc👩🏽‍💻", 4, "👩🏽‍💻a…", "…c👩🏽‍💻"},
		{"flag", "🇺🇸abc🇺🇸", 3, "🇺🇸…", "…🇺🇸"},
		{"exact fit", "项目", 4, "项目", "项目"},
		{"ellipsis only", "项目", 1, "…", "…"},
		{"no space", "项目", 0, "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := truncateHead(tt.value, tt.width); got != tt.head {
				t.Errorf("head = %q, want %q", got, tt.head)
			}
			if got := truncateTail(tt.value, tt.width); got != tt.tail {
				t.Errorf("tail = %q, want %q", got, tt.tail)
			}
		})
	}
}

func TestListingUnicodeAlignment(t *testing.T) {
	markers := marker.MarkerMap{
		"项目":          {Path: "/项目项目项目项目项目项目", Usage: 1},
		"e\u0301cole": {Path: "/cafe\u0301", Usage: 12},
		"👩🏽‍💻":        {Path: "/work/🇺🇸/👩🏽‍💻", Usage: 3},
	}
	keys := []string{"项目", "e\u0301cole", "👩🏽‍💻"}
	for _, width := range []int{9, 16, 26, 80} {
		out := FormatListing(markers, keys, ListingOptions{Format: ListingFormatBox, MaxWidth: width})
		lines := strings.Split(strings.TrimSuffix(out, "\n"), "\n")
		borderWidth := uniseg.StringWidth(lines[0])
		if borderWidth > width {
			t.Fatalf("width %d: border occupies %d columns", width, borderWidth)
		}
		var separators []int
		for _, line := range lines {
			if got := uniseg.StringWidth(line); got != borderWidth {
				t.Errorf("width %d: line occupies %d columns, want %d: %q", width, got, borderWidth, line)
			}
			if !strings.HasPrefix(line, "│") {
				continue
			}
			var positions []int
			for offset, char := range line {
				if char == '│' {
					positions = append(positions, uniseg.StringWidth(line[:offset]))
				}
			}
			if separators == nil {
				separators = positions
			}
			for i, position := range positions {
				if position != separators[i] {
					t.Errorf("width %d: separator at column %d, want %d: %q", width, position, separators[i], line)
				}
			}
		}
	}
}
