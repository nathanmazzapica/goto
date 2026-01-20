package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nathanmazzapica/goto/internal/marker"
)

// FormatListing returns a table of markers and destinations using box-drawing characters.
func FormatListing(markers marker.MarkerMap, keys []string) string {
	markerWidth := len("MARKER")
	destinationWidth := len("DESTINATION")
	usageWidth := len("USAGE")

	for _, key := range keys {
		if len(key) > markerWidth {
			markerWidth = len(key)
		}
		if len(markers[key].Path) > destinationWidth {
			destinationWidth = len(markers[key].Path)
		}

		usageLength := len(strconv.Itoa(markers[key].Usage))
		if usageLength > usageWidth {
			usageWidth = usageLength
		}
	}

	var builder strings.Builder
	fmt.Fprintf(&builder, "┌─%s─┬─%s─┬─%s─┐\n", strings.Repeat("─", markerWidth), strings.Repeat("─", usageWidth), strings.Repeat("─", destinationWidth))
	fmt.Fprintf(&builder, "│ %-*s │ %*s │ %-*s │\n", markerWidth, "MARKER", usageWidth, "USAGE", destinationWidth, "DESTINATION")
	fmt.Fprintf(&builder, "├─%s─┼─%s─┼─%s─┤\n", strings.Repeat("─", markerWidth), strings.Repeat("─", usageWidth), strings.Repeat("─", destinationWidth))

	for _, key := range keys {
		fmt.Fprintf(&builder, "│ %-*s │ %*d │ %-*s │\n", markerWidth, key, usageWidth, markers[key].Usage, destinationWidth, markers[key].Path)
	}

	fmt.Fprintf(&builder, "└─%s─┴─%s─┴─%s─┘\n", strings.Repeat("─", markerWidth), strings.Repeat("─", usageWidth), strings.Repeat("─", destinationWidth))
	return builder.String()
}
