package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/nathanmazzapica/goto/internal/marker"
)

// FormatListing returns a table of markers and destinations using box-drawing characters.
func FormatListing(markers marker.MarkerMap) string {
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
		if len(markers[key].Path) > destinationWidth {
			destinationWidth = len(markers[key].Path)
		}
	}

	var builder strings.Builder
	fmt.Fprintf(&builder, "┌─%s─┬─%s─┐\n", strings.Repeat("─", markerWidth), strings.Repeat("─", destinationWidth))
	fmt.Fprintf(&builder, "│ %-*s │ %-*s │\n", markerWidth, "MARKER", destinationWidth, "DESTINATION")
	fmt.Fprintf(&builder, "├─%s─┼─%s─┤\n", strings.Repeat("─", markerWidth), strings.Repeat("─", destinationWidth))

	for _, key := range keys {
		fmt.Fprintf(&builder, "│ %-*s │ %-*s │\n", markerWidth, key, destinationWidth, markers[key].Path)
	}

	fmt.Fprintf(&builder, "└─%s─┴─%s─┘\n", strings.Repeat("─", markerWidth), strings.Repeat("─", destinationWidth))
	return builder.String()
}
