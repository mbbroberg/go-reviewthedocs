package main

import "strings"

// Readme shows the required areas of the README.md file that can easily
// be parsed into a report.
type Readme struct {
	Title            string // Follow the convention
	Description      string // Anything there?
	GoVersionValue   string
	CollectedMetrics string // Should find a few but not all if > 20
}

func parseReadme(readmeList *map[string]string) (review Readme) {
	for _, readme := range *readmeList {
		review.Title = checkTitle(readme)
		review.Description = checkDescription(readme)
		review.GoVersionValue = checkGoVersionValue(readme)
		review.CollectedMetrics = checkCellectedMetrics(readme)

	}
	return review
}

func checkTitle(readme string) string {
	return "This is a title"
	// re, _ := regexp.MatchString("#*plugin*", readme) << re = bool
	// review.Title = re
}

func checkDescription(readme string) string {
	return "Good for now"
}

func checkGoVersionValue(readme string) string {
	versions := []string{}
	if strings.Contains(readme, "1.4") {
		versions = append(versions, "1.4")
	}
	if strings.Contains(readme, "1.5") {
		versions = append(versions, "1.5")
	}
	if strings.Contains(readme, "1.6") {
		versions = append(versions, "1.6")
	}
	value := strings.Join(versions, "")
	if value == "" {
		return "No Go version found"
	}
	return value

}

func checkCellectedMetrics(readme string) string {
	return "Found one."
}
