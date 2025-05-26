package service

import (
	. "github.com/onsi/gomega"
	"testing"
)

func TestShouldReturnValue(t *testing.T) {
	g := NewWithT(t)
	g.Expect(GetName()).To(Equal("Test"))
}
