package schoolsout

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSchoolsout(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Schools Out Suite")
}
