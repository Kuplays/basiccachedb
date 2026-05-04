package storage

import (
	"errors"
	"fmt"
	"sync"
	"testing"
)

func TestStorageSetAndGet(t *testing.T) {
	s := NewStorage()

	err := s.Set("name", "Aleksandr")

	if err != nil {
		t.Fatalf("Expecting no error at adding key value, got %v", err)
	}

	value, err := s.Get("name")
	if err != nil {
		t.Fatalf("Expecting no error, got %v", err)
	}

	if value != "Aleksandr" {
		t.Fatalf("Expecting Aleksandr, got %s", value)
	}
}

func TestStorageGetMissingKey(t *testing.T) {
	s := NewStorage()

	_, err := s.Get("no_key")

	if !errors.Is(err, ErrorKeyNotFound) {
		t.Fatalf("Expecting ErrorKeyNotFound, got %v", err)
	}
}

func TestStorageDeleteExistingKey(t *testing.T) {
	s := NewStorage()

	err := s.Set("name", "Aleskandr")

	if err != nil {
		t.Fatalf("Expecting no error at adding key value, got %v", err)
	}

	deleted, err := s.Delete("name")

	if err != nil {
		t.Fatalf("Expecting no error, got %v", err)
	}

	if !deleted {
		t.Fatal("Expecting deletion to return True on success")
	}

	_, err = s.Get("name")

	if !errors.Is(err, ErrorKeyNotFound) {
		t.Fatalf("Expecting ErrorKeyNotFound after deleting key, got %v", err)
	}
}

func TestStorageDeleteMissingKey(t *testing.T) {
	s := NewStorage()

	deleted, err := s.Delete("no_key")

	if err != nil {
		t.Fatalf("Expecting no error, got %v", err)
	}

	if deleted {
		t.Fatal("Expecting false since there is no key")
	}
}

func TestStorageExists(t *testing.T) {
	s := NewStorage()

	exists, err := s.Exists("name")

	if err != nil {
		t.Fatalf("Expecting no error, got %v", err)
	}

	if exists {
		t.Fatal("Expecting false since there is no key at creation of storage")
	}

	err = s.Set("name", "Aleksandr")

	if err != nil {
		t.Fatalf("Expecting no error, got %v", err)
	}

	exists, err = s.Exists("name")

	if err != nil {
		t.Fatalf("Expecting no error, got %v", err)
	}

	if !exists {
		t.Fatal("Expecting true since key is existing after inserting")
	}
}

func TestStorageLength(t *testing.T) {
	s := NewStorage()

	if s.Length() != 0 {
		t.Fatalf("Expecting length 0 at creation of an instance, got %d", s.Length())
	}

	_ = s.Set("a", "10")
	_ = s.Set("b", "20")

	if s.Length() != 2 {
		t.Fatalf("Expecting length to be equal to 2, got %d", s.Length())
	}

	_, _ = s.Delete("b")

	if s.Length() != 1 {
		t.Fatalf("Expecting length to be 1 after one deletion, got %d", s.Length())
	}
}

func TestStorageOverwriteValueAtKey(t *testing.T) {
	s := NewStorage()

	_ = s.Set("name", "Aleksandr")
	_ = s.Set("name", "John")

	value, err := s.Get("name")

	if err != nil {
		t.Fatalf("Expecting no error, got %v", err)
	}

	if value != "John" {
		t.Fatalf("Expecting overwrite key with John, got %s", value)
	}

	if s.Length() != 1 {
		t.Fatalf("Expecting Length to be equal to 1 after overwrite, got Length %d", s.Length())
	}
}

func TestStorageConcurrentAccess(t *testing.T) {
	s := NewStorage()

	_ = s.Set("shared", "initial")

	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(2)

		go func(i int) {
			defer wg.Done()

			_ = s.Set("shared", fmt.Sprintf("value:%d", i))
		}(i)

		go func() {
			defer wg.Done()

			_, err := s.Get("shared")
			if err != nil && !errors.Is(err, ErrorKeyNotFound) {
				t.Errorf("unexpected error: %v", err)
			}
		}()
	}

	wg.Wait()

	exists, err := s.Exists("shared")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !exists {
		t.Fatal("expected shared key to exist")
	}
}
