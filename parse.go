package main

import "strings"

func parseReadme(readmeList *map[string]string) (review string) {
	all := []string{""}
	for _, readme := range *readmeList {
		strings.Join(all, readme)
	}
	review = string(all)
	return review
}
