package matchers

import (
	"fmt"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
	networkingv1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
)

type WithGatewayMatcher struct {
	name, namespace, errorMsg, errorMsgNotted string
	matcher                                   types.GomegaMatcher
	metas                                     []types.GomegaMatcher
	failedMatcher                             types.GomegaMatcher
}

func WithGateway(name, ns string) *WithGatewayMatcher {
	meta := NewObjectMetaMatcher()
	meta.WithNamespace(ns)
	var metas []types.GomegaMatcher
	metas = append(metas, meta)
	return &WithGatewayMatcher{name: name, metas: metas}
}

func (matcher *WithGatewayMatcher) Match(actual interface{}) (bool, error) {
	docsMap, ok := actual.(map[string]interface{})
	if !ok {
		return false, fmt.Errorf("YAMLDocument must be passed a map[string]interface{}. Got\n%s", format.Object(actual, 1))
	}

	gateway, ok := docsMap["Gateway/"+matcher.name]
	if !ok {
		return false, nil
	}

	typedGateway, _ := gateway.(*networkingv1alpha3.Gateway)
	for _, meta := range matcher.metas {
		ok, err := meta.Match(typedGateway.ObjectMeta)
		if !ok || err != nil {
			matcher.failedMatcher = meta
			return ok, err
		}
	}
	return true, nil
}

func (matcher *WithGatewayMatcher) FailureMessage(actual interface{}) string {
	if matcher.failedMatcher == nil {
		msg := fmt.Sprintf(
			"FailureMessage: A gateway with name %q doesnt exist",
			matcher.name,
		)
		return msg
	}
	return matcher.failedMatcher.FailureMessage(actual)
}

func (matcher *WithGatewayMatcher) NegatedFailureMessage(actual interface{}) string {
	if matcher.failedMatcher == nil {
		msg := fmt.Sprintf(
			"FailureMessage: A gateway with name %q exists",
			matcher.name,
		)
		return msg
	}
	return matcher.failedMatcher.NegatedFailureMessage(actual)
}
