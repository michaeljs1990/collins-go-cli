package main

import (
	"fmt"
	"os"
)

// This is kinda dumb but go has pretty limited data structure types
// this means no Sets. We only support adding values via the Add method
// creating a UniqueOrderedSet with duplicate values will not cause it to
// be filtered. This is not efficient.
type UniqueOrderedSet []string

func (u UniqueOrderedSet) Contains(s string) bool {

	for _, val := range u {
		if val == s {
			return true
		}
	}

	return false
}

func (u UniqueOrderedSet) Add(s string) UniqueOrderedSet {
	if u.Contains(s) {
		return u
	}

	return append(u, s)
}

func logAndDie(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
