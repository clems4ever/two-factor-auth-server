package suites

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type KubernetesSuite struct {
	*SeleniumSuite
}

func NewKubernetesSuite() *KubernetesSuite {
	return &KubernetesSuite{SeleniumSuite: new(SeleniumSuite)}
}

func (s *KubernetesSuite) TestOneFactorScenario() {
	suite.Run(s.T(), NewOneFactorScenario())
}

func (s *KubernetesSuite) TestTwoFactorScenario() {
	suite.Run(s.T(), NewTwoFactorScenario())
}

func (s *KubernetesSuite) TestRedirectionURLScenario() {
	suite.Run(s.T(), NewRedirectionURLScenario())
}

func TestKubernetesSuite(t *testing.T) {
	suite.Run(t, NewKubernetesSuite())
}
