package api

import (
	"testing"

	"github.com/NavroO/go-key-value-store/pkg/storage"
)

func BenchmarkSetKey(b *testing.B) {
	store := storage.NewInMemoryStorage()
	for i := 0; i < b.N; i++ {
		store.Set("key", "value", 10)
	}
}

func BenchmarkGetKey(b *testing.B) {
	store := storage.NewInMemoryStorage()
	store.Set("key", "value", 10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.Get("key")
	}
}

func BenchmarkDeleteKey(b *testing.B) {
	store := storage.NewInMemoryStorage()
	store.Set("key", "value", 10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.Delete("key")
	}
}
