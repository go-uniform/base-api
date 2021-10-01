package service

import "strings"

// get the index of a string inside of a string array
func indexOf(haystack []string, needle string, caseSensitive bool) int {
	if haystack == nil {
		panic("specify an array to search through")
	}
	if needle == "" {
		panic("specify a string to search for")
	}
	if caseSensitive {
		for i, item := range haystack {
			if item == needle {
				return i
			}
		}
	} else {
		lowNeedle := strings.ToLower(needle)
		for i, item := range haystack {
			if strings.ToLower(item) == lowNeedle {
				return i
			}
		}
	}
	return -1
}

// see if a string array contains a given string
func contains(haystack []string, needle string, caseSensitive bool) bool {
	if haystack == nil {
		panic("specify an array to search through")
	}
	if needle == "" {
		panic("specify a string to search for")
	}
	return indexOf(haystack, needle, caseSensitive) != -1
}

// trim the filterItems from the items array
func filter(items []string, filterItems []string) []string {
	if filterItems == nil || len(filterItems) <= 0 {
		return items
	}
	newItems := make([]string, 0)
	for _, item := range items {
		if contains(filterItems, item, false) {
			continue
		}
		newItems = append(newItems, item)
	}
	return newItems
}
