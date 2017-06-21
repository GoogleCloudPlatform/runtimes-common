package utils

import (
	"github.com/pmezard/go-difflib/difflib"
)

// Modification of difflib's unified differ
func GetAdditions(a, b []string) []string {
	matcher := difflib.NewMatcher(a, b)
	differences := matcher.GetGroupedOpCodes(0)

	var adds []string
	for _, g := range differences {
		for _, c := range g {
			j1, j2 := c.J1, c.J2
			if c.Tag == 'r' || c.Tag == 'i' {
				for _, line := range b[j1:j2] {
					adds = append(adds, line)
				}
			}
		}
	}
	return adds
}

func GetDeletions(a, b []string) []string {
	matcher := difflib.NewMatcher(a, b)
	differences := matcher.GetGroupedOpCodes(0)

	var dels []string
	for _, g := range differences {
		for _, c := range g {
			i1, i2 := c.I1, c.I2
			if c.Tag == 'r' || c.Tag == 'd' {
				for _, line := range a[i1:i2] {
					dels = append(dels, line)
				}
			}
		}
	}
	return dels
}

func GetMatches(a, b []string) []string {
	matcher := difflib.NewMatcher(a, b)
	matchindexes := matcher.GetMatchingBlocks()

	var matches []string
	for i, m := range matchindexes {
		if i != len(matches) - 1 {
			start := m.A
			end := m.A + m.Size
			for _, line := range a[start:end] {
				matches = append(matches, line)
			}
		}
	}
	return matches
}

