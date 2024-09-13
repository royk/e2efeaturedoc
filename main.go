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
	Subfeatures []Feature
}

// Extract test cases from the content
func extractTestCases(content string) Feature {
	features := Feature{}
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
					Subsets: []string{},
				}
				features.Subfeatures = append(features.Subfeatures, *currentFeature)
			}
		}
		if strings.Contains(line, "test") {
			testMatches := testRegex.FindAllStringSubmatch(line, -1)
			for _, test := range testMatches {
				features.Subfeatures[len(features.Subfeatures)-1].Subsets = append(features.Subfeatures[len(features.Subfeatures)-1].Subsets, test[1])
			}
		}
	}

	return features
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
			features.Subfeatures = append(features.Subfeatures, extractedFeatures)
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
