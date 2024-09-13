package main

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Feature listing", func() {
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
		Expect(feature).To(Equal(map[string][]string{
			"Feature 1": {"Test 1"},
			"Feature 2": {"Test 1"},
		}))
	})
})
