package main

import "fmt"

func parseReadme(readmeList *map[string]string) {

	for _, readme := range *readmeList {
		fmt.Printf(readme)
	}
}
