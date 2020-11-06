package matchers

import (
	"fmt"

	qsv1a1 "code.cloudfoundry.org/quarks-secret/pkg/kube/apis/quarkssecret/v1alpha1"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	. "github.com/onsi/gomega/gstruct"
	"github.com/onsi/gomega/types"
)

type QuarksSecretMatcher struct {
	stringData types.GomegaMatcher
	data       types.GomegaMatcher
	meta       *ObjectMetaMatcher

	executed types.GomegaMatcher
}

func RepresentingAQuarksSecret() *QuarksSecretMatcher {
	return &QuarksSecretMatcher{
		nil,
		nil,
		NewObjectMetaMatcher(),
		nil,
	}
}

func (matcher *QuarksSecretMatcher) WithStringData(name string, value string) *QuarksSecretMatcher {
	matcher.stringData = MatchKeys(IgnoreExtras, Keys{
		name: Equal(value),
	})
	return matcher
}

func (matcher *QuarksSecretMatcher) WithData(name string, value []byte) *QuarksSecretMatcher {
	matcher.data = MatchKeys(IgnoreExtras, Keys{
		name: Equal(value),
	})
	return matcher
}

func (matcher *QuarksSecretMatcher) WithName(name string) *QuarksSecretMatcher {
	matcher.meta.WithName(name)

	return matcher
}

func (matcher *QuarksSecretMatcher) Match(actual interface{}) (success bool, err error) {
	secret, ok := actual.(*qsv1a1.QuarksSecret)
	if !ok {
		return false, fmt.Errorf("Expected a secret. Got\n%s", format.Object(actual, 1))
	}
	if secret == nil {
		return false, fmt.Errorf("secret matcher: secret is nil")
	}

	// if matcher.stringData != nil {
	// 	matcher.executed = matcher.stringData
	// 	if pass, err := matcher.stringData.Match(secret.StringData); !pass || err != nil {
	// 		return pass, err
	// 	}
	// }

	// if matcher.data != nil {
	// 	matcher.executed = matcher.data
	// 	if pass, err := matcher.data.Match(secret.Data); !pass || err != nil {
	// 		return pass, err
	// 	}
	// }

	matcher.executed = matcher.meta
	if pass, err := matcher.meta.Match(secret.ObjectMeta); !pass || err != nil {
		return pass, err
	}

	return true, nil
}

func (matcher *QuarksSecretMatcher) FailureMessage(actual interface{}) string {
	if matcher.executed != nil {
		return matcher.executed.FailureMessage(actual)
	}
	return "FailureMessage: Possibly missing .WithDataValue or .WithStringDataValue method in the test."
}

func (matcher *QuarksSecretMatcher) NegatedFailureMessage(actual interface{}) string {
	return matcher.executed.NegatedFailureMessage(actual)
}
