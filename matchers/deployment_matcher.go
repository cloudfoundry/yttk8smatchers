package matchers

import (
	"fmt"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	. "github.com/onsi/gomega/gstruct"
	"github.com/onsi/gomega/types"
	coreV1 "k8s.io/api/core/v1"
)

type DeploymentMatcher struct {
	replicas       types.GomegaMatcher
	meta       *ObjectMetaMatcher

	executed types.GomegaMatcher
}

func RepresentingADeployment() *DeploymentMatcher {
	return &DeploymentMatcher{
		nil,
		NewObjectMetaMatcher(),
		nil,
	}
}

func (matcher *DeploymentMatcher) WithReplicas(count int) *DeploymentMatcher {
	matcher.replicas = MatchKeys(IgnoreExtras, Keys{
		"spec": MatchKeys(IgnoreExtras, Keys{
			"replicas": Equal(count),
		}),
	})
	return matcher
}

func (matcher *DeploymentMatcher) Match(actual interface{}) (success bool, err error) {
	deployment, ok := actual.(*coreV1.Secret)
	if !ok {
		return false, fmt.Errorf("Expected a deployment. Got\n%s", format.Object(actual, 1))
	}
	if deployment == nil {
		return false, fmt.Errorf("deployment matcher: deployment is nil")
	}

	if matcher.replicas != nil {
		matcher.executed = matcher.stringData
		if pass, err := matcher.stringData.Match(deployment.StringData); !pass || err != nil {
			return pass, err
		}
	}

	if matcher.data != nil {
		matcher.executed = matcher.data
		if pass, err := matcher.data.Match(deployment.Data); !pass || err != nil {
			return pass, err
		}
	}

	matcher.executed = matcher.meta
	if pass, err := matcher.meta.Match(deployment.ObjectMeta); !pass || err != nil {
		return pass, err
	}

	return true, nil
}

func (matcher *DeploymentMatcher) FailureMessage(actual interface{}) string {
	if matcher.executed != nil {
		return matcher.executed.FailureMessage(actual)
	}
	return "FailureMessage: Possibly missing .WithDataValue or .WithStringDataValue method in the test."
}

func (matcher *DeploymentMatcher) NegatedFailureMessage(actual interface{}) string {
	return matcher.executed.NegatedFailureMessage(actual)
}
