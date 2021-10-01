package service

import (
	"fmt"
	"testing"
)

func TestIndexOf(t *testing.T) {
	tests := []struct {
		Haystacks      []string
		Needle         string
		CaseSensitive  bool
		ExpectedOutput int
		ExpectedError  string
	}{
		//Defensive tests
		{nil, "Needle1", false, -1, "specify an array to search through"},
		{[]string{"needle1", "needle2"}, "", false, -1, "specify a string to search for"},

		//Positive tests
		{[]string{"needle1", "needle2"}, "needle1", false, 0, ""},
		{[]string{"needle1", "needle1", "needle2"}, "needle1", false, 0, ""},
		{[]string{"needle1", "needle1", "needle2", "Needle2"}, "Needle2", true, 3, ""},

		//Negative tests
		{[]string{"needle1", "needle2"}, "needle3", false, -1, ""},
		{[]string{"needle1", "needle2"}, "Needle2", true, -1, ""},
	}

	for i, test := range tests {
		func() {
			// recovery
			defer func() {
				if err := recover(); err != nil {
					if test.ExpectedError != fmt.Sprintf("%v", err) {
						t.Errorf("[%d] | haystack %s | needle %s | case-sensitive %v | error = %v; want %v", i, test.Haystacks, test.Needle, test.CaseSensitive, err, test.ExpectedError)
					}
				}
			}()

			// execution
			output := indexOf(test.Haystacks, test.Needle, test.CaseSensitive)

			// assertion
			if test.ExpectedError != "" {
				t.Errorf("[%d] | haystack %s | needle %s | case-sensitive %v  error = <empty>; want %v", i, test.Haystacks, test.Needle, test.CaseSensitive, test.ExpectedError)
			} else if output != test.ExpectedOutput {
				t.Errorf("[%d] | haystack %s | needle %s | case-sensitive %v  output = %v; want %v", i, test.Haystacks, test.Needle, test.CaseSensitive, output, test.ExpectedOutput)
			}
		}()
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		Haystacks      []string
		Needle         string
		CaseSensitive  bool
		ExpectedOutput bool
		ExpectedError  string
	}{

		//Defensive tests
		{nil, "Needle1", true, false, "specify an array to search through"},
		{[]string{"needle1", "needle2"}, "", true, false, "specify a string to search for"},

		//Positive tests
		{[]string{"needle1", "needle2"}, "needle2", false, true, ""},

		//Negative tests
		{[]string{"needle1", "needle2"}, "Needle1", true, false, ""},
		{[]string{"needle1", "needle1"}, "needle1", true, true, ""},
	}

	for i, test := range tests {
		func() {
			// recovery
			defer func() {
				if err := recover(); err != nil {
					if test.ExpectedError != fmt.Sprintf("%v", err) {
						t.Errorf("[%d] | haystack %s | needle %s | case-sensitive %v | error = %v; want %v", i, test.Haystacks, test.Needle, test.CaseSensitive, err, test.ExpectedError)
					}
				}
			}()

			// execution
			output := contains(test.Haystacks, test.Needle, test.CaseSensitive)

			// assertion
			if test.ExpectedError != "" {
				t.Errorf("[%d] | haystack %s | needle %s | case-sensitive %v  error = <empty>; want %v", i, test.Haystacks, test.Needle, test.CaseSensitive, test.ExpectedError)
			} else if output != test.ExpectedOutput {
				t.Errorf("[%d] | haystack %s | needle %s | case-sensitive %v  output = %v; want %v", i, test.Haystacks, test.Needle, test.CaseSensitive, output, test.ExpectedOutput)
			}
		}()
	}
}

func TestFilter(t *testing.T) {
	tests := []struct {
		Items          []string
		FilterItems    []string
		ExpectedOutput []string
	}{
		//Defensive tests
		{nil, nil, nil},
		{nil, []string{}, nil},
		{[]string{}, nil, []string{}},
		{[]string{}, []string{}, []string{}},

		//Positive tests
		{[]string{"needle1", "needle2"}, []string{"needle1"}, []string{"needle2"}},

		//Negative tests
		{[]string{"needle1", "needle2"}, []string{"needle3"}, []string{"needle1", "needle2"}},
	}

	for i, test := range tests {
		func() {
			// execution
			output := filter(test.Items, test.FilterItems)

			hasAll := true
			if len(output) != len(test.ExpectedOutput) {
				hasAll = false
			} else {
				for _, item := range test.ExpectedOutput {
					if !contains(output, item, true) {
						hasAll = false
						break
					}
				}
			}

			// assertion
			if !hasAll {
				t.Errorf("[%d] | items %v | filter-items %v | output = %v; want %v", i, test.Items, test.FilterItems, output, test.ExpectedOutput)
			}
		}()
	}
}
