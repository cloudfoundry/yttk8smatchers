package matchers

import (
	"fmt"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/gstruct"
	"github.com/onsi/gomega/types"
	coreV1 "k8s.io/api/core/v1"
)

type WithConfigMapMatcher struct {
	name, namespace, errorMsg, errorMsgNotted string

	matcher types.GomegaMatcher
	metas   []types.GomegaMatcher

	data gstruct.Keys

	failedMatcher types.GomegaMatcher
}

func WithConfigMap(name, ns string) *WithConfigMapMatcher {
	meta := NewObjectMetaMatcher()
	meta.WithNamespace(ns)
	var metas []types.GomegaMatcher
	metas = append(metas, meta)
	return &WithConfigMapMatcher{name: name, metas: metas}
}

func (matcher *WithConfigMapMatcher) WithData(dm gstruct.Keys) *WithConfigMapMatcher {
	matcher.data = dm
	return matcher
}

func (matcher *WithConfigMapMatcher) Match(actual interface{}) (bool, error) {
	docsMap, ok := actual.(map[string]interface{})
	if !ok {
		return false, fmt.Errorf("YAMLDocument must be passed a map[string]interface{}. Got\n%s", format.Object(actual, 1))
	}

	configMap, ok := docsMap["ConfigMap/"+matcher.name]
	if !ok {
		return false, nil
	}

	typedConfigMap, _ := configMap.(*coreV1.ConfigMap)

	for _, meta := range matcher.metas {
		ok, err := meta.Match(typedConfigMap.ObjectMeta)
		if !ok || err != nil {
			matcher.failedMatcher = meta
			return ok, err
		}
	}

	dataMatcher := gstruct.MatchKeys(gstruct.IgnoreExtras, matcher.data)
	ok, err := dataMatcher.Match(typedConfigMap.Data)
	if !ok || err != nil {
		fmt.Println("here")
		matcher.failedMatcher = dataMatcher
		return ok, err
	}

	return true, nil
}

func (matcher *WithConfigMapMatcher) FailureMessage(actual interface{}) string {
	if matcher.failedMatcher == nil {
		msg := fmt.Sprintf(
			"FailureMessage: A statefulset with name %q doesnt exist",
			matcher.name,
		)
		return msg
	}
	return matcher.failedMatcher.FailureMessage(actual)
}

func (matcher *WithConfigMapMatcher) NegatedFailureMessage(actual interface{}) string {
	if matcher.failedMatcher == nil {
		msg := fmt.Sprintf(
			"FailureMessage: A statefulset with name %q exists",
			matcher.name,
		)
		return msg
	}
	return matcher.failedMatcher.NegatedFailureMessage(actual)
}
