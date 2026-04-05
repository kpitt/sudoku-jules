package solver

import (
	"testing"
)

func BenchmarkFormatHouses(b *testing.B) {
	houses := []*House{
		{Kind: kindRow, Index: 0},
		{Kind: kindRow, Index: 1},
		{Kind: kindRow, Index: 2},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = formatHouses(houses)
	}
}
