package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestGort(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gort Suite")
}
