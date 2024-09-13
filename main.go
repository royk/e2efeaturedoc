package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Feature struct {
	Name        string
	Subsets     []string
	Subfeatures []*Feature
}

// Extract test cases from the content
func extractTestCases(content string) Feature {
	root := &Feature{
		Name:        "",
		Subsets:     []string{},
		Subfeatures: []*Feature{},
	}
	describeRegex := regexp.MustCompile(`(?ms)describe\(["'](.*?)["'](.*?)`)
	testRegex := regexp.MustCompile(`(?m)test\(["'](.*?)["'],`)
	contentLines := strings.Split(content, "\n")
	currentContextDepth := 0
	currentFeatureDepth := 0
	var currentFeature = root
	for _, line := range contentLines {
		if strings.Contains(line, "{") {
			currentContextDepth++
		}
		if strings.Contains(line, "}") {
			currentContextDepth--
		}
		if strings.Contains(line, "describe") {
			describeMatches := describeRegex.FindAllStringSubmatch(line, -1)
			if len(describeMatches) > 0 {
				describe := describeMatches[0][1]
				nextFeature := &Feature{
					Name:        describe,
					Subsets:     []string{},
					Subfeatures: []*Feature{},
				}
				// find the current feature from root based on currentContextDepth
				currentFeature = root
				for i := 1; i < currentContextDepth; i++ {
					currentFeature = currentFeature.Subfeatures[0]
				}
				currentFeature.Subfeatures = append(currentFeature.Subfeatures, nextFeature)
				currentFeature = nextFeature
				currentFeatureDepth++
			}
		}
		if strings.Contains(line, "test") {
			testMatches := testRegex.FindAllStringSubmatch(line, -1)
			for _, test := range testMatches {
				currentFeature.Subsets = append(currentFeature.Subsets, test[1])
			}
		}
	}

	return *root
}

func generateFeatureDocFromFeatures(features Feature) (string, error) {
	var docBuilder strings.Builder

	docBuilder.WriteString("# App Features Document\n\n")
	for _, feature := range features.Subfeatures {
		docBuilder.WriteString(fmt.Sprintf("## Feature: %s\n\n", feature.Name))
		for _, test := range feature.Subsets {
			docBuilder.WriteString(fmt.Sprintf("- %s\n", test))
		}
		docBuilder.WriteString("\n")
	}
	return docBuilder.String(), nil
}

// Generate the markdown document from the test files
func generateFeatureDoc(testDirectory string) (string, error) {
	features := Feature{}
	err := filepath.Walk(testDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".js") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			extractedFeatures := extractTestCases(string(content))
			features.Subfeatures = append(features.Subfeatures, &extractedFeatures)
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return generateFeatureDocFromFeatures(features)
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
