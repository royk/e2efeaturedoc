package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestE2efeaturedoc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2efeaturedoc Suite")
}
