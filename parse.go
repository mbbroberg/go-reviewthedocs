package main

import (
	"fmt"
	"strings"
)

type Readme struct {
	Title       string // Follow the convention
	Description string // Anything
	// 	SystemRequirements string // Does it list privileges? Link to tools?
	// 	Installation       string
	// 	Example            string
	// 	ExampleJSON        string
	// 	MetricCatalog      string
	// 	Roadmap            string
	// 	License            string
	// 	Acknowledgement    string // including Authors
	// 	BuildStatus        bool   // Should find it somewhere on the file
	// 	LicenseEmbedded    bool   // Should find it somewhere on the file
}

func parseReadme(readmeList *map[string]string) (review string) {
	for _, readme := range *readmeList {
		results := Readme{}
		splitReadme := strings.Split(readme, "#")
		results.Title = splitReadme[1]
		results.Description = splitReadme[0]
		fmt.Println(results.Title)
		fmt.Println("Title ^ and now, description:")
		fmt.Println(results.Description)
	}
	return review
}
