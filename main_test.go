package main

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("E2E Feature Documentation", func() {
	Describe("Feature Doc Generation", func() {
		It("Generates an MD doc file for a given feature set", func() {
			features := []Feature{
				{
					Name:    "Feature 1",
					Subsets: []string{"Test 1"},
				},
				{
					Name:    "Feature 2",
					Subsets: []string{"Test 2"},
				},
			}

			doc, err := generateFeatureDocFromFeatures(features)
			Expect(err).To(BeNil())
			// strip the header by popping it
			docLines := strings.Split(doc, "\n")
			docLines = docLines[2:]
			doc = strings.Join(docLines, "\n")
			Expect(doc).To(Equal("## Feature: Feature 1\n\n- Test 1\n\n## Feature: Feature 2\n\n- Test 2\n\n"))
		})
	})
	Describe("Feature listing", func() {
		It("Extracts feature list from Describe blocks", func() {
			content := `
			describe("Feature 1", func() {
				test("Test 1", func() {
					Expect(1).To(Equal(1))
				})
			})
			describe("Feature 2", func() {
				test("Test 1", func() {
					Expect(1).To(Equal(1))
				})
			})
		`
			feature := extractTestCases(content)
			Expect(feature).To(Equal([]Feature{
				{
					Name:    "Feature 1",
					Subsets: []string{"Test 1"},
				},
				{
					Name:    "Feature 2",
					Subsets: []string{"Test 1"},
				},
			}))
		})
		It("Extracts nested features", func() {
			Skip("Not implemented")
			content := `
			describe("Feature 1", func() {
				describe("Feature 1.1", func() {
					test("Test 1", func() {
						Expect(1).To(Equal(1))
					})
				})
			})`
			feature := extractTestCases(content)
			Expect(feature).To(Equal(Feature{
				Name:    "Feature 1",
				Subsets: []string{},
				Subfeatures: []Feature{
					{
						Name:        "Feature 1.1",
						Subsets:     []string{"Test 1"},
						Subfeatures: []Feature{},
					},
				},
			}))
		})
	})
})
