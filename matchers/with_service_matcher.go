package matchers

import (
	"fmt"

	"github.com/benjamintf1/unmarshalledmatchers"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/gstruct"
	"github.com/onsi/gomega/types"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"
)

type WithServiceMatcher struct {
	name                string
	serviceMatcher      *ServiceMatcher
	specYaml            string
	data                gstruct.Keys
	failedMatcher       types.GomegaMatcher
	failedMatcherActual interface{}
}

func WithService(name string) *WithServiceMatcher {
	return &WithServiceMatcher{name: name, serviceMatcher: RepresentingAService()}
}

func (matcher *WithServiceMatcher) WithSpecYaml(yaml string) *WithServiceMatcher {
	matcher.specYaml = yaml
	return matcher
}

func (matcher *WithServiceMatcher) Match(actual interface{}) (bool, error) {
	matcher.failedMatcherActual = actual
	docsMap, ok := actual.(map[string]interface{})
	if !ok {
		return false, fmt.Errorf("YAMLDocument must be passed a map[string]interface{}. Got\n%s", format.Object(actual, 1))
	}

	value, ok := docsMap["Service/"+matcher.name]
	if !ok {
		return false, nil
	}

	typedService := value.(*v1.Service)
	ok, err := matcher.serviceMatcher.Match(value.(*v1.Service))
	if !ok || err != nil {
		matcher.failedMatcher = matcher.serviceMatcher
		return ok, err
	}

	if matcher.specYaml != "" {
		serviceSpecYaml, err := yaml.Marshal(typedService.Spec)
		if err != nil {
			return false, err
		}

		yamlMatcher := unmarshalledmatchers.ContainUnorderedYAML(matcher.specYaml)
		ok, err = yamlMatcher.Match(serviceSpecYaml)
		if !ok || err != nil {
			matcher.failedMatcher = yamlMatcher
			matcher.failedMatcherActual = serviceSpecYaml
			return ok, err
		}
	}

	return true, nil
}

func (matcher *WithServiceMatcher) FailureMessage(actual interface{}) string {
	if matcher.failedMatcher == nil {
		msg := fmt.Sprintf(
			"FailureMessage: A Service with name %q doesnt exist",
			matcher.name,
		)
		return msg
	}
	return matcher.failedMatcher.FailureMessage(matcher.failedMatcherActual)
}

func (matcher *WithServiceMatcher) NegatedFailureMessage(actual interface{}) string {
	serviceMatcherFailureMessage := matcher.serviceMatcher.NegatedFailureMessage(actual)
	if serviceMatcherFailureMessage != "" {
		return serviceMatcherFailureMessage
	}

	msg := fmt.Sprintf(
		"FailureMessage: A Service with name %q exists",
		matcher.name,
	)
	return msg
}

func (matcher *WithServiceMatcher) WithType(value string) types.GomegaMatcher {
	matcher.serviceMatcher.WithType(value)
	return matcher
}

func (matcher *WithServiceMatcher) WithData(dm gstruct.Keys) *WithServiceMatcher {
	matcher.data = dm
	return matcher
}

type WithoutServiceMatcher struct {
	name               string
	withServiceMatcher *WithServiceMatcher
}

func WithoutService(name string) *WithoutServiceMatcher {
	return &WithoutServiceMatcher{name, &WithServiceMatcher{}}
}

func (matcher *WithoutServiceMatcher) Match(actual interface{}) (bool, error) {
	result, err := matcher.withServiceMatcher.Match(actual)
	return !result, err
}

func (matcher *WithoutServiceMatcher) FailureMessage(actual interface{}) string {
	msg := fmt.Sprintf(
		"FailureMessage: A Service with name %q does exist",
		matcher.name,
	)
	return msg
}

func (matcher *WithoutServiceMatcher) NegatedFailureMessage(actual interface{}) string {
	msg := fmt.Sprintf(
		"FailureMessage: A Service with name %q does not exist",
		matcher.name,
	)
	return msg
}
