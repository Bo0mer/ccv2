package ccv2_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCcv2(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CCv2 Suite")
}
