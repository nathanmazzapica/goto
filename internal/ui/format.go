package ui

import (
	"fmt"
	"sort"
	"strings"
)

// FormatListing returns a table of markers and destinations using box-drawing characters.
func FormatListing(markers map[string]string) string {
	markerWidth := len("MARKER")
	destinationWidth := len("DESTINATION")

	keys := make([]string, 0, len(markers))
	for key := range markers {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		if len(key) > markerWidth {
			markerWidth = len(key)
		}
		if len(markers[key]) > destinationWidth {
			destinationWidth = len(markers[key])
		}
	}

	var builder strings.Builder
	fmt.Fprintf(&builder, "┌─%s─┬─%s─┐\n", strings.Repeat("─", markerWidth), strings.Repeat("─", destinationWidth))
	fmt.Fprintf(&builder, "│ %-*s │ %-*s │\n", markerWidth, "MARKER", destinationWidth, "DESTINATION")
	fmt.Fprintf(&builder, "├─%s─┼─%s─┤\n", strings.Repeat("─", markerWidth), strings.Repeat("─", destinationWidth))

	for _, key := range keys {
		fmt.Fprintf(&builder, "│ %-*s │ %-*s │\n", markerWidth, key, destinationWidth, markers[key])
	}

	fmt.Fprintf(&builder, "└─%s─┴─%s─┘\n", strings.Repeat("─", markerWidth), strings.Repeat("─", destinationWidth))
	return builder.String()
}
