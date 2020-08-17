package matchers

import (
	"fmt"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
	appsv1 "k8s.io/api/apps/v1"
)

type WithDeploymentMatcher struct {
	name, namespace, errorMsg, errorMsgNotted string
	matcher                                   types.GomegaMatcher
	metas                                     []types.GomegaMatcher
	failedMatcher                             types.GomegaMatcher
}

func WithDeployment(name, ns string) *WithDeploymentMatcher {
	meta := NewObjectMetaMatcher()
	meta.WithNamespace(ns)
	var metas []types.GomegaMatcher
	metas = append(metas, meta)
	return &WithDeploymentMatcher{name: name, metas: metas}
}

func (matcher *WithDeploymentMatcher) Match(actual interface{}) (bool, error) {
	docsMap, ok := actual.(map[string]interface{})
	if !ok {
		return false, fmt.Errorf("YAMLDocument must be passed a map[string]interface{}. Got\n%s", format.Object(actual, 1))
	}

	deployment, ok := docsMap["Deployment/"+matcher.name]
	if !ok {
		return false, nil
	}

	typedDeployment, _ := deployment.(*appsv1.Deployment)

	for _, meta := range matcher.metas {
		ok, err := meta.Match(typedDeployment.ObjectMeta)
		if !ok || err != nil {
			matcher.failedMatcher = meta
			return ok, err
		}
	}
	return true, nil
}

func (matcher *WithDeploymentMatcher) FailureMessage(actual interface{}) string {
	if matcher.failedMatcher == nil {
		msg := fmt.Sprintf(
			"FailureMessage: A statefulset with name %q doesnt exist",
			matcher.name,
		)
		return msg
	}
	return matcher.failedMatcher.FailureMessage(actual)
}

func (matcher *WithDeploymentMatcher) NegatedFailureMessage(actual interface{}) string {
	if matcher.failedMatcher == nil {
		msg := fmt.Sprintf(
			"FailureMessage: A statefulset with name %q exists",
			matcher.name,
		)
		return msg
	}
	return matcher.failedMatcher.NegatedFailureMessage(actual)
}
