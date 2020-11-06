package matchers

import (
	"fmt"

	qsv1a1 "code.cloudfoundry.org/quarks-secret/pkg/kube/apis/quarkssecret/v1alpha1"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

type WithQuarksSecretMatcher struct {
	name, namespace, errorMsg, errorMsgNotted string
	matcher                                   types.GomegaMatcher
	metas                                     []types.GomegaMatcher
	failedMatcher                             types.GomegaMatcher
}

func WithQuarksSecret(name, ns string) *WithQuarksSecretMatcher {
	meta := NewObjectMetaMatcher()
	meta.WithNamespace(ns)
	var metas []types.GomegaMatcher
	metas = append(metas, meta)
	return &WithQuarksSecretMatcher{name: name, metas: metas}
}

func (matcher *WithQuarksSecretMatcher) Match(actual interface{}) (bool, error) {
	docsMap, ok := actual.(map[string]interface{})
	if !ok {
		return false, fmt.Errorf("YAMLDocument must be passed a map[string]interface{}. Got\n%s", format.Object(actual, 1))
	}

	gateway, ok := docsMap["QuarksSecret/"+matcher.name]
	if !ok {
		return false, nil
	}

	typedQuarksSecret, _ := gateway.(*qsv1a1.QuarksSecret)
	for _, meta := range matcher.metas {
		ok, err := meta.Match(typedQuarksSecret.ObjectMeta)
		if !ok || err != nil {
			matcher.failedMatcher = meta
			return ok, err
		}
	}
	return true, nil
}

func (matcher *WithQuarksSecretMatcher) FailureMessage(actual interface{}) string {
	if matcher.failedMatcher == nil {
		msg := fmt.Sprintf(
			"FailureMessage: A quarks secret with name %q doesnt exist",
			matcher.name,
		)
		return msg
	}
	return matcher.failedMatcher.FailureMessage(actual)
}

func (matcher *WithQuarksSecretMatcher) NegatedFailureMessage(actual interface{}) string {
	if matcher.failedMatcher == nil {
		msg := fmt.Sprintf(
			"FailureMessage: A quarks secret with name %q exists",
			matcher.name,
		)
		return msg
	}
	return matcher.failedMatcher.NegatedFailureMessage(actual)
}
