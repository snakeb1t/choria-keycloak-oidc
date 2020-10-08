package dictionary

import (
	"strings"
)

// Store is a set of words that exist in the dictionary
type Store struct {
	dict map[string]struct{}
}

// New returns a new, empty Store
func New() *Store {
	return &Store{
		dict: make(map[string]struct{}),
	}
}

// Add inserts a new word into the Store
func (s *Store) Add(word string) {
	s.dict[strings.ToLower(word)] = struct{}{}
}

// Remove removes a word from the store, and returns true if the word existed, false if it did not exist
func (s *Store) Remove(word string) bool {
	if !s.Lookup(word) {
		return false
	}
	delete(s.dict, strings.ToLower(word))
	return true
}

// Lookup checks if a word exists in the Store and returns true if it does, false if it does not
func (s *Store) Lookup(word string) bool {
	_, ok := s.dict[strings.ToLower(word)]
	return ok
}
