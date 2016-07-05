package main

// Add type declarations for *Result structs.
// These will be pre-generated based on the go:generate templates.
// Declare both a db field and json entry like so:
type SampleResult struct {
	Id string `db:"sample.id" json:"sample.id"`
}
