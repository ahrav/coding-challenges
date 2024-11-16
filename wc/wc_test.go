package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCountBytes(t *testing.T) {
	want := 342190
	got := countBytes("test.txt")
	assert.Equal(t, want, got)
}

func BenchmarkCountBytes(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		countBytes("test.txt")
	}
}

func TestCountWords(t *testing.T) {
	want := 58164
	got, err := countWords("test.txt")
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func BenchmarkCountWords(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = countWords("test.txt")
	}
}

func TestCountLines(t *testing.T) {
	want := 7145
	got := countLines("test.txt")
	assert.Equal(t, want, got)
}

func BenchmarkCountLines(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		countLines("test.txt")
	}
}
