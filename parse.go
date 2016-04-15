package main

import (
	"fmt"
	"regexp"
	"strings"
)

// Readme shows the required areas of the README.md file that can easily
// be parsed into a report.
type Readme struct {
	Title              bool   // Follow the convention
	Description        string // Anything
	GoVersion          bool
	SystemRequirements string // Does it list privileges? Link to tools?
	// Installation       string
	// Example            string
	// ExampleJSON        string
	// MetricCatalog      string
	// Roadmap            string
	// License            string
	// Acknowledgement    string // including Authors
	// BuildStatus        bool   // Should find it somewhere on the file
	// LicenseEmbedded    bool   // Should find it somewhere on the file
}

func parseReadme(readmeList *map[string]string) (review Readme) {
	for _, readme := range *readmeList {
		switch {
		case strings.ContainsAny(readme, "1.5 & 1.6"):
			review.GoVersion = true
		default:
			fmt.Println("Didn't find any supported Go versions")
		}

		re, _ := regexp.MatchString("#*plugin*", readme)
		review.Title = re
	}
	return review
}
