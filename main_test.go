package main

import (
	"testing"
)

func TestDeduplicate(t *testing.T) {
	s := []string{"a", "b", "b", "c"}

	// Remove the extra "b" from the slice.
	s = deduplicate(s)

	// Check that the slice has three elements.
	if len(s) != 3 {
		t.Errorf("Length of de-duped slice is %v, expected 3", len(s))
	}
}

func TestCompare(t *testing.T) {
	a := []string{"a", "b", "c"}
	b := []string{"b", "c", "d"}

	// Subtract slice b from slice a, leaving []string{"a"}.
	s := compare(a, b)

	// Check that the slice has only one element.
	if len(s) != 1 {
		t.Errorf("Length of comparison slice is %v, expected 1", len(s))
	}
	// Check that the slice doesn't contain anything other than "a".
	for _, v := range s {
		if v != "a" {
			t.Errorf("Found string \"%s\" in comparison slice, expected \"a\"", v)
		}
	}
	// Check that we didn't mix up the order and include the unique "d" from slice b.
	for _, v := range s {
		if v == "d" {
			t.Errorf("Found string \"d\" in comparison slice, expected \"a\"")
		}
	}
}
