package matchers

import (
	"errors"
	"fmt"
	"os/exec"
	"reflect"
	"regexp"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/types"
)

type ProduceYAMLMatcher struct {
	matcher  types.GomegaMatcher
	rendered string
}

func ProduceYAML(matcher types.GomegaMatcher) *ProduceYAMLMatcher {
	return &ProduceYAMLMatcher{matcher, ""}
}

func (matcher *ProduceYAMLMatcher) Match(actual interface{}) (bool, error) {
	rendering, ok := actual.(RenderingContext)
	if !ok {
		return false, fmt.Errorf("ProduceYAML must be passed a RenderingContext. Got\n%s", format.Object(actual, 1))
	}

	session, err := renderWithData(rendering.templates, rendering.data)
	if err != nil {
		return false, fmt.Errorf("render error, exit status={%v}, command={%s}, error={%v}", session.ExitCode(), session.Command, err)
	}

	matcher.rendered = string(session.Out.Contents())

	if session.ExitCode() != 0 {
		return matcher.matcher.Match(errors.New(string(session.Err.Contents())))
	}

	docsMap, err := parseYAML(session.Out)
	if err != nil {
		return false, err
	}

	return matcher.matcher.Match(docsMap)
}

func (matcher *ProduceYAMLMatcher) FailureMessage(actual interface{}) string {
	msg := fmt.Sprintf(
		"FailureMessage: There is a problem with this YAML:\n\n%s\n\n%s",
		matcher.rendered,
		matcher.matcher.FailureMessage(actual),
	)
	return msg
}

func (matcher *ProduceYAMLMatcher) NegatedFailureMessage(actual interface{}) string {
	msg := fmt.Sprintf(
		"NegatedFailureMessage: There is a problem with this YAML:\n\n%s\n\n%s",
		matcher.rendered,
		matcher.matcher.NegatedFailureMessage(actual),
	)
	return msg
}

func renderWithData(templates []string, data map[string]interface{}) (*gexec.Session, error) {
	var args []string
	for _, template := range templates {
		args = append(args, "-f", template)
	}

	for k, i := range data {
		switch v := i.(type) {
		case bool:
				args = append(args, "--data-value-yaml", fmt.Sprintf("%s=%t", k, v))
		case int:
				args = append(args, "--data-value-yaml", fmt.Sprintf("%s=%d", k, v))
		case string:
				args = append(args, "--data-value-yaml", fmt.Sprintf("%s=%q", k, v))
		default:
				return nil, fmt.Errorf("Unsupported data value type for key %q: %T", k, v)
		}
}

	command := exec.Command("ytt", args...)
	session, err := gexec.Start(command, nil, GinkgoWriter)
	if err != nil {
		return session, err
	}

	return session.Wait(), nil
}

func parseYAML(yaml *gbytes.Buffer) (interface{}, error) {
	apiObjects := map[string]interface{}{}
	decode := GetDecoder().Decode

	// for each document
	ptn := regexp.MustCompile(`(?m)^---\s*$`)
	docStrings := ptn.Split(string(yaml.Contents()), -1)
	for _, docString := range docStrings {
		// Checks for empty documents
		if docString == "" {
			continue
		}

		obj, gk, err := decode([]byte(docString), nil, nil)
		if err != nil {
			continue
		}

		apiObjects[gk.Kind+"/"+reflect.ValueOf(obj).Elem().FieldByName("Name").String()] = obj
	}

	return apiObjects, nil
}
