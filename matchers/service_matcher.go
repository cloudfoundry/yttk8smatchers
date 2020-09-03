package matchers

import (
	"fmt"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	. "github.com/onsi/gomega/gstruct"
	"github.com/onsi/gomega/types"
	"k8s.io/api/core/v1"
)

type ServiceMatcher struct {
	spec       *ServiceSpecMatcher

	executed types.GomegaMatcher
}

func RepresentingAService() *ServiceMatcher {
	return &ServiceMatcher{
		NewServiceSpecMatcher(),
		nil,
	}
}

func (matcher *ServiceMatcher) WithType(value string) *ServiceMatcher {
	matcher.spec.WithType(value)

	return matcher
}

func (matcher *ServiceMatcher) Match(actual interface{}) (success bool, err error) {
	service, ok := actual.(*v1.Service)
	if !ok {
		return false, fmt.Errorf("Expected a service. Got\n%s", format.Object(actual, 1))
	}

	matcher.executed = matcher.spec
	if pass, err := matcher.spec.Match(service.Spec); !pass || err != nil {
		return pass, err
	}

	return true, nil
}

func (matcher *ServiceMatcher) FailureMessage(actual interface{}) string {
	return matcher.executed.FailureMessage(actual)
}

func (matcher *ServiceMatcher) NegatedFailureMessage(actual interface{}) string {
	return matcher.executed.NegatedFailureMessage(actual)
}

type ServiceSpecMatcher struct {
	fields map[string]types.GomegaMatcher

	spec     *v1.ServiceSpec
	executed types.GomegaMatcher
}

func NewServiceSpecMatcher() *ServiceSpecMatcher {
	return &ServiceSpecMatcher{map[string]types.GomegaMatcher{}, nil, nil}
}

func (matcher *ServiceSpecMatcher) WithType(value string) *ServiceSpecMatcher {
	matcher.fields["Type"] = Equal(v1.ServiceType (value))

	return matcher
}

func (matcher *ServiceSpecMatcher) Match(actual interface{}) (bool, error) {
	spec, ok := actual.(v1.ServiceSpec)
	if !ok {
		return false, fmt.Errorf("Expecting ServiceSpec. Got\n%s", format.Object(actual, 1))
	}

	matcher.spec = &spec
	matcher.executed = MatchFields(IgnoreExtras, matcher.fields)
	return matcher.executed.Match(spec)
}

func (matcher *ServiceSpecMatcher) FailureMessage(actual interface{}) string {
	return fmt.Sprintf(
		"ServiceSpec should match: \n%s",
		matcher.executed.FailureMessage(&matcher.spec),
	)
}

func (matcher *ServiceSpecMatcher) NegatedFailureMessage(actual interface{}) string {
	return fmt.Sprintf(
		"ServiceSpec should not match: \n%s",
		matcher.executed.FailureMessage(&matcher.spec),
	)
}
