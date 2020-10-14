package matchers

import (
	"fmt"
	"github.com/benjamintf1/unmarshalledmatchers"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/yaml"
)

type WithDaemonSetMatcher struct {
	name, namespace, errorMsg, errorMsgNotted string
	matcher                                   types.GomegaMatcher
	metas                                     []types.GomegaMatcher
	specYaml                                  string
	failedMatcher                             types.GomegaMatcher
	failedMatcherActual                       interface{}
}

func WithDaemonSet(name, ns string) *WithDaemonSetMatcher {
	meta := NewObjectMetaMatcher()
	meta.WithNamespace(ns)
	var metas []types.GomegaMatcher
	metas = append(metas, meta)
	return &WithDaemonSetMatcher{name: name, metas: metas}
}

func (matcher *WithDaemonSetMatcher) WithSpecYaml(yaml string) *WithDaemonSetMatcher {
	matcher.specYaml = yaml
	return matcher
}

func (matcher *WithDaemonSetMatcher) Match(actual interface{}) (bool, error) {
	matcher.failedMatcherActual = actual
	docsMap, ok := actual.(map[string]interface{})
	if !ok {
		return false, fmt.Errorf("YAMLDocument must be passed a map[string]interface{}. Got\n%s", format.Object(actual, 1))
	}

	daemonSet, ok := docsMap["DaemonSet/"+matcher.name]
	if !ok {
		return false, nil
	}

	typedDaemonSet, _ := daemonSet.(*appsv1.DaemonSet)

	for _, meta := range matcher.metas {
		ok, err := meta.Match(typedDaemonSet.ObjectMeta)
		if !ok || err != nil {
			matcher.failedMatcher = meta
			return ok, err
		}
	}

	if matcher.specYaml != "" {
		daemonSetSpecYaml, err := yaml.Marshal(typedDaemonSet.Spec)
		if err != nil {
			return false, err
		}

		yamlMatcher := unmarshalledmatchers.ContainUnorderedYAML(matcher.specYaml)
		ok, err = yamlMatcher.Match(daemonSetSpecYaml)
		if !ok || err != nil {
			matcher.failedMatcher = yamlMatcher
			matcher.failedMatcherActual = daemonSetSpecYaml
			return ok, err
		}
	}

	return true, nil
}

func (matcher *WithDaemonSetMatcher) FailureMessage(actual interface{}) string {
	if matcher.failedMatcher == nil {
		msg := fmt.Sprintf(
			"FailureMessage: A daemon set with name %q doesnt exist",
			matcher.name,
		)
		return msg
	}
	return matcher.failedMatcher.FailureMessage(matcher.failedMatcherActual)
}

func (matcher *WithDaemonSetMatcher) NegatedFailureMessage(actual interface{}) string {
	if matcher.failedMatcher == nil {
		msg := fmt.Sprintf(
			"FailureMessage: A daemon set with name %q exists and/or the spec yaml was not correct %v",
			matcher.name,
			matcher.specYaml)
		return msg
	}
	return matcher.failedMatcher.NegatedFailureMessage(matcher.failedMatcherActual)
}
