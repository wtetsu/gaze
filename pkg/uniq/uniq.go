/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package uniq

// Uniq can deal with unique list.
type Uniq struct {
	list []string
	keys map[string]struct{}
}

// New returns a new Uniq.
func New() *Uniq {
	return &Uniq{
		list: []string{},
		keys: map[string]struct{}{},
	}
}

// Add adds a new entry.
func (u *Uniq) Add(newEntry string) {
	_, ok := u.keys[newEntry]
	if ok {
		return
	}
	u.keys[newEntry] = struct{}{}
	u.list = (append(u.list, newEntry))
}

// AddAll adds multiple new entries.
func (u *Uniq) AddAll(newEntries []string) {
	for _, e := range newEntries {
		u.Add(e)
	}
}

// List returns a internal unique list.
func (u *Uniq) List() []string {
	return u.list
}

// Len returns the length
func (u *Uniq) Len() int {
	return len(u.list)
}
