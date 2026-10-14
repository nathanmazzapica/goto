package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nathanmazzapica/goto/internal/marker"
	"github.com/rivo/uniseg"
)

type ListingFormat string

const (
	ListingFormatBox   ListingFormat = "box"
	ListingFormatPlain ListingFormat = "plain"
)

type ListingOptions struct {
	Format   ListingFormat
	MaxWidth int
}

// FormatListing renders a list of markers based on the provided options.
func FormatListing(markers marker.MarkerMap, keys []string, opts ListingOptions) string {
	switch opts.Format {
	case ListingFormatPlain:
		return FormatListingPlain(markers, keys)
	case ListingFormatBox:
		return formatListingBox(markers, keys, opts.MaxWidth)
	default:
		return FormatListingPlain(markers, keys)
	}
}

// FormatListingPlain renders markers as tab-separated rows with no headers.
func FormatListingPlain(markers marker.MarkerMap, keys []string) string {
	var builder strings.Builder
	for _, key := range keys {
		fmt.Fprintf(&builder, "%s\t%d\t%s\n", key, markers[key].Usage, markers[key].Path)
	}
	return builder.String()
}

// formatListingBox returns a table of markers and destinations using box-drawing characters.
func formatListingBox(markers marker.MarkerMap, keys []string, maxWidth int) string {
	if len(keys) == 0 {
		return ""
	}

	markerWidth, usageWidth, destinationWidth := computeColumnWidths(markers, keys)

	if maxWidth <= 0 || tableWidth([]int{markerWidth, usageWidth, destinationWidth}) <= maxWidth {
		return renderBox(markers, keys, markerWidth, usageWidth, destinationWidth, true)
	}

	// Drop usage column before truncating.
	available := maxWidth - tableWidth([]int{0, 0})
	if available < 2 {
		return FormatListingPlain(markers, keys)
	}

	markerFit := markerWidth
	destinationFit := destinationWidth

	if markerFit+destinationFit > available {
		destinationFit = available - markerFit
		if destinationFit < 1 {
			destinationFit = 1
			markerFit = available - destinationFit
		}
	}

	if markerFit < 1 || destinationFit < 1 {
		return FormatListingPlain(markers, keys)
	}

	return renderBox(markers, keys, markerFit, 0, destinationFit, false)
}

func computeColumnWidths(markers marker.MarkerMap, keys []string) (int, int, int) {
	markerWidth := cellWidth("MARKER")
	destinationWidth := cellWidth("DESTINATION")
	usageWidth := cellWidth("USAGE")

	for _, key := range keys {
		if width := cellWidth(key); width > markerWidth {
			markerWidth = width
		}
		if width := cellWidth(markers[key].Path); width > destinationWidth {
			destinationWidth = width
		}

		usageLength := cellWidth(strconv.Itoa(markers[key].Usage))
		if usageLength > usageWidth {
			usageWidth = usageLength
		}
	}

	return markerWidth, usageWidth, destinationWidth
}

func renderBox(markers marker.MarkerMap, keys []string, markerWidth, usageWidth, destinationWidth int, includeUsage bool) string {
	var builder strings.Builder

	markerHeader := truncateCell("MARKER", markerWidth, truncateHead)
	destinationHeader := truncateCell("DESTINATION", destinationWidth, truncateTail)

	if includeUsage {
		usageHeader := truncateCell("USAGE", usageWidth, truncateHead)
		fmt.Fprintf(&builder, "┌─%s─┬─%s─┬─%s─┐\n", strings.Repeat("─", markerWidth), strings.Repeat("─", usageWidth), strings.Repeat("─", destinationWidth))
		fmt.Fprintf(&builder, "│ %s │ %*s │ %s │\n", padCell(markerHeader, markerWidth), usageWidth, usageHeader, padCell(destinationHeader, destinationWidth))
		fmt.Fprintf(&builder, "├─%s─┼─%s─┼─%s─┤\n", strings.Repeat("─", markerWidth), strings.Repeat("─", usageWidth), strings.Repeat("─", destinationWidth))

		for _, key := range keys {
			markerCell := truncateCell(key, markerWidth, truncateHead)
			destCell := truncateCell(markers[key].Path, destinationWidth, truncateTail)
			fmt.Fprintf(&builder, "│ %s │ %*d │ %s │\n", padCell(markerCell, markerWidth), usageWidth, markers[key].Usage, padCell(destCell, destinationWidth))
		}

		fmt.Fprintf(&builder, "└─%s─┴─%s─┴─%s─┘\n", strings.Repeat("─", markerWidth), strings.Repeat("─", usageWidth), strings.Repeat("─", destinationWidth))
		return builder.String()
	}

	fmt.Fprintf(&builder, "┌─%s─┬─%s─┐\n", strings.Repeat("─", markerWidth), strings.Repeat("─", destinationWidth))
	fmt.Fprintf(&builder, "│ %s │ %s │\n", padCell(markerHeader, markerWidth), padCell(destinationHeader, destinationWidth))
	fmt.Fprintf(&builder, "├─%s─┼─%s─┤\n", strings.Repeat("─", markerWidth), strings.Repeat("─", destinationWidth))

	for _, key := range keys {
		markerCell := truncateCell(key, markerWidth, truncateHead)
		destCell := truncateCell(markers[key].Path, destinationWidth, truncateTail)
		fmt.Fprintf(&builder, "│ %s │ %s │\n", padCell(markerCell, markerWidth), padCell(destCell, destinationWidth))
	}

	fmt.Fprintf(&builder, "└─%s─┴─%s─┘\n", strings.Repeat("─", markerWidth), strings.Repeat("─", destinationWidth))
	return builder.String()
}

func tableWidth(widths []int) int {
	if len(widths) == 0 {
		return 0
	}

	total := 0
	for _, width := range widths {
		total += width
	}

	return total + (3 * (len(widths) - 1)) + 4
}

func cellWidth(value string) int {
	return uniseg.StringWidth(value)
}

// padCell pads an already truncated cell by terminal columns, not rune count.
func padCell(value string, width int) string {
	return value + strings.Repeat(" ", max(0, width-cellWidth(value)))
}

func truncateCell(value string, width int, truncate func(string, int) string) string {
	if width <= 0 {
		return ""
	}
	if cellWidth(value) <= width {
		return value
	}
	return truncate(value, width)
}

func truncateHead(value string, width int) string {
	if width <= 0 {
		return ""
	}
	if cellWidth(value) <= width {
		return value
	}
	if width == 1 {
		return "…"
	}

	graphemes := uniseg.NewGraphemes(value)
	used, end := 0, 0
	for graphemes.Next() {
		if used+graphemes.Width() > width-1 {
			break
		}
		used += graphemes.Width()
		_, end = graphemes.Positions()
	}
	return value[:end] + "…"
}

func truncateTail(value string, width int) string {
	if width <= 0 {
		return ""
	}
	if cellWidth(value) <= width {
		return value
	}
	if width == 1 {
		return "…"
	}

	remaining := cellWidth(value)
	graphemes := uniseg.NewGraphemes(value)
	start := 0
	for graphemes.Next() {
		remaining -= graphemes.Width()
		_, start = graphemes.Positions()
		if remaining <= width-1 {
			break
		}
	}
	return "…" + value[start:]
}
