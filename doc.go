// Package fracindex generates lexicographically sortable fractional index keys.
//
// An Indexer represents one key space. New without options creates a
// Rocicorp-compatible base62 key space. Empty bounds passed to KeyBetween and
// NKeysBetween mean the beginning or end of that key space.
package fracindex
