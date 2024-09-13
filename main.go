package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Feature struct {
	Name    string
	subsets []string
}

// Extract test cases from the content
func extractTestCases(content string) []Feature {
	features := []Feature{}
	describeRegex := regexp.MustCompile(`(?ms)describe\(["'](.*?)["'](.*?)`)
	testRegex := regexp.MustCompile(`(?m)test\(["'](.*?)["'],`)
	contentLines := strings.Split(content, "\n")
	var currentFeature *Feature
	for _, line := range contentLines {
		if strings.Contains(line, "describe") {
			describeMatches := describeRegex.FindAllStringSubmatch(line, -1)
			for _, describe := range describeMatches {
				currentFeature = &Feature{
					Name:    describe[1],
					subsets: []string{},
				}
				currentFeature.Name = describe[1]
				currentFeature.subsets = []string{}
				features = append(features, *currentFeature)
			}
		}
		if strings.Contains(line, "test") {
			testMatches := testRegex.FindAllStringSubmatch(line, -1)
			for _, test := range testMatches {
				features[len(features)-1].subsets = append(features[len(features)-1].subsets, test[1])
			}
		}
	}

	return features
}

// Generate the markdown document from the test files
func generateFeatureDoc(testDirectory string) (string, error) {
	var docBuilder strings.Builder

	docBuilder.WriteString("# App Features Document\n\n")

	err := filepath.Walk(testDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".js") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			features := extractTestCases(string(content))
			for _, feature := range features {
				docBuilder.WriteString(fmt.Sprintf("## Feature: %s\n\n", feature.Name))
				for _, test := range feature.subsets {
					docBuilder.WriteString(fmt.Sprintf("- %s\n", test))
				}
				docBuilder.WriteString("\n")
			}
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return docBuilder.String(), nil
}

func main() {
	testDirectory := "./tests" // Adjust this to your test folder
	doc, err := generateFeatureDoc(testDirectory)
	if err != nil {
		fmt.Println("Error generating document:", err)
		return
	}

	err = os.WriteFile("features.md", []byte(doc), 0644)
	if err != nil {
		fmt.Println("Error writing document:", err)
		return
	}

	fmt.Println("Features documentation generated: features.md")
}
