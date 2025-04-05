package main

import (
	"sync"
	"testing"
	"time"
)

func TestSetAndGet(t *testing.T) {
	cache := NewTtlCache()
	defer cache.Stop()

	cache.Set("key1", "value1", 0)
	val, ok := cache.Get("key1")
	if !ok {
		t.Error("Failed to get value that was just set")
	}
	if val != "value1" {
		t.Errorf("Expected 'value1', got %s", val)
	}

	val, ok = cache.Get("non-existent")
	if ok {
		t.Error("Get should return false for non-existent key")
	}
	if val != "" {
		t.Errorf("Expected empty string for non-existent key, got %s", val)
	}

	cache.Set("key1", "updated", 0)
	val, ok = cache.Get("key1")
	if !ok {
		t.Error("Failed to get updated value")
	}
	if val != "updated" {
		t.Errorf("Expected 'updated', got %s", val)
	}
}

func TestDelete(t *testing.T) {
	cache := NewTtlCache()
	defer cache.Stop()

	cache.Set("key1", "value1", 0)
	_, ok := cache.Get("key1")
	if !ok {
		t.Error("Key should exist before deletion")
	}

	cache.Delete("key1")
	_, ok = cache.Get("key1")
	if ok {
		t.Error("Key should have been deleted")
	}

	cache.Delete("non-existent")
}

func TestExpiration(t *testing.T) {
	cache := NewTtlCache()
	defer cache.Stop()

	cache.Set("short", "shortvalue", 50*time.Millisecond)
	cache.Set("forever", "eternalvalue", 0)

	_, ok1 := cache.Get("short")
	_, ok2 := cache.Get("forever")
	if !ok1 || !ok2 {
		t.Error("Both keys should exist initially")
	}

	time.Sleep(100 * time.Millisecond)

	_, ok1 = cache.Get("short")
	_, ok2 = cache.Get("forever")
	if ok1 {
		t.Error("Key with short TTL should have expired")
	}
	if !ok2 {
		t.Error("Key with no TTL should not expire")
	}
}

func TestAutomaticCleanup(t *testing.T) {
	cache := NewTtlCache()
	defer cache.Stop()

	cache.Set("expiring", "value", 1*time.Second)

	val, ok := cache.Get("expiring")
	if !ok || val != "value" {
		t.Error("Key should exist initially")
	}

	time.Sleep(6 * time.Second)

	_, ok = cache.Get("expiring")
	if ok {
		t.Error("Expired key should have been automatically cleaned up")
	}
}

func TestConcurrentAccess(t *testing.T) {
	cache := NewTtlCache()
	defer cache.Stop()

	var wg sync.WaitGroup
	numOperations := 100

	for i := range 5 {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := range numOperations {
				key := "key" + string(rune('A'+workerID))
				cache.Set(key, "value"+string(rune('0'+j%10)), 500*time.Millisecond)
				time.Sleep(5 * time.Millisecond)
			}
		}(i)
	}

	for i := range 5 {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := range numOperations {
				key := "key" + string(rune('A'+j%5))
				cache.Get(key)
				time.Sleep(2 * time.Millisecond)
			}
		}(i)
	}

	for i := range 2 {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := range numOperations / 5 {
				key := "key" + string(rune('A'+j%5))
				cache.Delete(key)
				time.Sleep(10 * time.Millisecond)
			}
		}(i)
	}

	wg.Wait()
}

func TestStopSafety(t *testing.T) {
	cache := NewTtlCache()

	cache.Stop()

	cache.Set("key", "value", 0)
	val, ok := cache.Get("key")
	if !ok || val != "value" {
		t.Error("Cache operations should still work after Stop")
	}

	defer func() {
		if r := recover(); r != nil {
			t.Error("Multiple calls to Stop should not panic")
		}
	}()
	cache.Stop()
}

func TestEmptyStrings(t *testing.T) {
	cache := NewTtlCache()
	defer cache.Stop()

	cache.Set("", "empty key", 0)
	val, ok := cache.Get("")
	if !ok {
		t.Error("Failed to get value with empty string key")
	}
	if val != "empty key" {
		t.Errorf("Expected 'empty key', got '%s'", val)
	}

	cache.Set("empty-value", "", 0)
	val, ok = cache.Get("empty-value")
	if !ok {
		t.Error("Failed to get empty string value")
	}
	if val != "" {
		t.Errorf("Expected empty string, got '%s'", val)
	}

	cache.Delete("")
	_, ok = cache.Get("")
	if ok {
		t.Error("Empty string key should be deleted")
	}
}

func TestTtlUpdates(t *testing.T) {
	cache := NewTtlCache()
	defer cache.Stop()

	cache.Set("key", "value", 100*time.Millisecond)

	time.Sleep(50 * time.Millisecond)
	cache.Set("key", "value", 500*time.Millisecond)

	time.Sleep(100 * time.Millisecond)
	val, ok := cache.Get("key")
	if !ok {
		t.Error("Key should not have expired after TTL update")
	}
	if val != "value" {
		t.Errorf("Expected 'value', got '%s'", val)
	}

	time.Sleep(400 * time.Millisecond)
	_, ok = cache.Get("key")
	if ok {
		t.Error("Key should have expired after the updated TTL")
	}

	cache.Set("convert", "value", 100*time.Millisecond)
	time.Sleep(50 * time.Millisecond)
	cache.Set("convert", "permanent", 0)

	time.Sleep(100 * time.Millisecond)
	val, ok = cache.Get("convert")
	if !ok {
		t.Error("Key should not expire after conversion to non-expiring")
	}
	if val != "permanent" {
		t.Errorf("Expected 'permanent', got '%s'", val)
	}
}
