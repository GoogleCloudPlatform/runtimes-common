package cmd

import (
	"testing"
)

type testpair struct {
	input           []string
	expected_output bool
}

var tests = []testpair{
	{[]string{}, false},
	{[]string{"one"}, false},
	{[]string{"one", "two"}, false},
	{[]string{"one", "two", "three"}, true},
	{[]string{"one", "two", "three", "four"}, false},
}

func TestArgNum(t *testing.T) {
	for _, test := range tests {
		valid, err := checkArgNum(test.input)
		if valid != test.expected_output {
			if test.expected_output {
				t.Errorf("Got unexpected error: %s", err)
			} else {
				t.Errorf("Expected error but got none")
			}
		}
	}
}
