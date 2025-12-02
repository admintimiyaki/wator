package main

import "testing"

func BenchmarkInitWorld(b *testing.B) {
	for i := 0; i < b.N; i++ {
		initWorld()
	}
}

func BenchmarkUpdate(b *testing.B) {
	initWorld()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		update()
	}
}
