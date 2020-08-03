package utils

import (
	"strings"
)

// StringArray type define
type StringArray []string

// Set appends s to string array
func (a *StringArray) Set(s string) error {
	*a = append(*a, s)
	return nil
}

// String formats string array
func (a *StringArray) String() string {
	return strings.Join(*a, ",")
}
