package main

import (
	"testing"
)

func createTestMarkers(n int) map[string]string {
	markers := make(map[string]string, n)
	names := []string{"zulu", "alpha", "mike", "bravo", "yankee", "charlie",
		"xray", "delta", "whiskey", "echo", "victor", "foxtrot", "uniform",
		"golf", "tango", "hotel", "sierra", "india", "romeo", "juliet",
		"quebec", "kilo", "papa", "lima", "oscar", "november"}

	for i := 0; i < n && i < len(names); i++ {
		markers[names[i]] = "/some/path/" + names[i]
	}
	return markers
}

func BenchmarkSortKeys(b *testing.B) {
	markers := createTestMarkers(10)
	for i := 0; i < b.N; i++ {
		sortKeys(markers)
	}
}

func BenchmarkIterateUnsorted(b *testing.B) {
	markers := createTestMarkers(10)
	for i := 0; i < b.N; i++ {
		for key := range markers {
			_ = key
			_ = markers[key]
		}
	}
}

func BenchmarkSortKeys_25Markers(b *testing.B) {
	markers := createTestMarkers(25)
	for i := 0; i < b.N; i++ {
		sortKeys(markers)
	}
}

func BenchmarkIterateUnsorted_25Markers(b *testing.B) {
	markers := createTestMarkers(25)
	for i := 0; i < b.N; i++ {
		for key := range markers {
			_ = key
			_ = markers[key]
		}
	}
}
