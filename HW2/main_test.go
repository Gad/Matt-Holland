package main

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func BenchmarkCounterTree(b *testing.B) {

	raw = strings.ReplaceAll(raw, "\n", "")
	s := strings.NewReader(raw)
	tree, _ := html.Parse(s)
	for range 1000000 {
		_, _ = countTree(tree)
	}
}

func BenchmarkCounterTreeRecur(b *testing.B) {

	raw = strings.ReplaceAll(raw, "\n", "")
	s := strings.NewReader(raw)
	tree, _ := html.Parse(s)
	for range 1000000 {
		var tc, ic int
		countTreeRecur(tree, &tc, &ic)
	}

}
