package matchers

import (
	"fmt"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
	v1 "k8s.io/api/core/v1"
)

type WithSecretMatcher struct {
	name string
	secretMatcher *SecretMatcher
}

func WithSecret(name string) *WithSecretMatcher {
	return &WithSecretMatcher{name, RepresentingASecret()}
}

func (matcher *WithSecretMatcher) Match(actual interface{}) (bool, error) {
	docsMap, ok := actual.(map[string]interface{})
	if !ok {
		return false, fmt.Errorf("YAMLDocument must be passed a map[string]interface{}. Got\n%s", format.Object(actual, 1))
	}

	value, ok := docsMap["Secret/"+matcher.name]
	if !ok {
		return false, nil
	}

	return matcher.secretMatcher.Match(value.(*v1.Secret))
}

func (matcher *WithSecretMatcher) FailureMessage(actual interface{}) string {
	secretMatcherFailureMessage := matcher.secretMatcher.FailureMessage(actual)
	if secretMatcherFailureMessage != "" {
		return secretMatcherFailureMessage
	}

	msg := fmt.Sprintf(
		"FailureMessage: A Secret with name %q doesnt exist",
		matcher.name,
	)
	return msg
}

func (matcher *WithSecretMatcher) NegatedFailureMessage(actual interface{}) string {
	secretMatcherFailureMessage := matcher.secretMatcher.NegatedFailureMessage(actual)
	if secretMatcherFailureMessage != "" {
		return secretMatcherFailureMessage
	}

	msg := fmt.Sprintf(
		"FailureMessage: A Secret with name %q exists",
		matcher.name,
	)
	return msg
}

func (matcher *WithSecretMatcher) WithDataValue(key string, expectedBase64DecodedValue []byte) types.GomegaMatcher {
	matcher.secretMatcher.WithData(key, expectedBase64DecodedValue)
	return matcher
}

type WithoutSecretMatcher struct {
	name              string
	withSecretMatcher *WithSecretMatcher
}

func WithoutSecret(name string) *WithoutSecretMatcher {
	return &WithoutSecretMatcher{name, &WithSecretMatcher{}}
}

func (matcher *WithoutSecretMatcher) Match(actual interface{}) (bool, error) {
	result, err := matcher.withSecretMatcher.Match(actual)
	return !result, err
}

func (matcher *WithoutSecretMatcher) FailureMessage(actual interface{}) string {
	msg := fmt.Sprintf(
		"FailureMessage: A Secret with name %q does exist",
		matcher.name,
	)
	return msg
}

func (matcher *WithoutSecretMatcher) NegatedFailureMessage(actual interface{}) string {
	msg := fmt.Sprintf(
		"FailureMessage: A Secret with name %q does not exist",
		matcher.name,
	)
	return msg
}
