package cmd

import (
	"testing"
)

type testpair struct {
	values   []string
	expected bool
}

var tests = []testpair{
	{[]string{}, false},
	{[]string{"one"}, false},
	{[]string{"one", "two"}, true},
	{[]string{"one", "two", "three"}, true},
	{[]string{"one", "two", "three", "four"}, false},
}

func TestArgNum(t *testing.T) {
	for _, test := range tests {
		valid, err := checkArgNum(test.values)
		if valid != test.expected {
			if test.expected {
				t.Errorf("Got unexpected error: %s", err)
			} else {
				t.Errorf("Expected error but got none")
			}
		}
	}
}
