package matchers

import (
	"fmt"
	"github.com/benjamintf1/unmarshalledmatchers"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/yaml"
)

type WithDeploymentMatcher struct {
	name, namespace, errorMsg, errorMsgNotted string
	matcher                                   types.GomegaMatcher
	metas                                     []types.GomegaMatcher
	specYaml                                  string
	failedMatcher                             types.GomegaMatcher
	failedMatcherActual                       interface{}
}

func WithDeployment(name, ns string) *WithDeploymentMatcher {
	meta := NewObjectMetaMatcher()
	meta.WithNamespace(ns)
	var metas []types.GomegaMatcher
	metas = append(metas, meta)
	return &WithDeploymentMatcher{name: name, metas: metas}
}

func (matcher *WithDeploymentMatcher) WithSpecYaml(yaml string) *WithDeploymentMatcher {
	matcher.specYaml = yaml
	return matcher
}

func (matcher *WithDeploymentMatcher) Match(actual interface{}) (bool, error) {
	matcher.failedMatcherActual = actual
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

	if matcher.specYaml != "" {
		deploymentSpecYaml, err := yaml.Marshal(typedDeployment.Spec)
		if err != nil {
			return false, err
		}

		yamlMatcher := unmarshalledmatchers.ContainUnorderedYAML(matcher.specYaml)
		ok, err = yamlMatcher.Match(deploymentSpecYaml)
		if !ok || err != nil {
			matcher.failedMatcher = yamlMatcher
			matcher.failedMatcherActual = deploymentSpecYaml
			return ok, err
		}
	}

	return true, nil
}

func (matcher *WithDeploymentMatcher) FailureMessage(actual interface{}) string {
	if matcher.failedMatcher == nil {
		msg := fmt.Sprintf(
			"FailureMessage: A deployment with name %q doesnt exist",
			matcher.name,
		)
		return msg
	}
	return matcher.failedMatcher.FailureMessage(matcher.failedMatcherActual)
}

func (matcher *WithDeploymentMatcher) NegatedFailureMessage(actual interface{}) string {
	if matcher.failedMatcher == nil {
		msg := fmt.Sprintf(
			"FailureMessage: A deployment with name %q exists and/or the spec yaml was not correct %v",
			matcher.name,
			matcher.specYaml)
		return msg
	}
	return matcher.failedMatcher.NegatedFailureMessage(matcher.failedMatcherActual)
}
