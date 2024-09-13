package main

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("E2E Feature Documentation", func() {
	Describe("Feature Doc Generation", func() {
		It("Generates an MD doc file for a given feature set", func() {
			features := []Feature{
				{
					Name:    "Feature 1",
					subsets: []string{"Test 1"},
				},
				{
					Name:    "Feature 2",
					subsets: []string{"Test 1"},
				},
			}

			doc, err := generateFeatureDocFromFeatures(features)
			Expect(err).To(BeNil())
			Expect(doc).To(ContainSubstring("Feature 1"))
			Expect(doc).To(ContainSubstring("Feature 2"))
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
					subsets: []string{"Test 1"},
				},
				{
					Name:    "Feature 2",
					subsets: []string{"Test 1"},
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
			Expect(feature).To(Equal(map[string][]string{
				"Feature 1":   {"Feature 1.1"},
				"Feature 1.1": {"Test 1"},
			}))
		})
	})
})
