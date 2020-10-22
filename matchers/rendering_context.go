package matchers

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/otiai10/copy"
)

type RenderingContext struct {
	targetDir string
	templates []string
	data      map[string]interface{}
}

type RenderingContextOption func(RenderingContext) (RenderingContext, error)

func NewRenderingContext(opts ...RenderingContextOption) (RenderingContext, error) {
	r := RenderingContext{}
	var err error

	for _, opt := range opts {
		r, err = opt(r)

		if err != nil {
			return RenderingContext{}, nil
		}
	}

	return r, nil
}

func WithData(data map[string]interface{}) RenderingContextOption {
	return func(r RenderingContext) (RenderingContext, error) {
		r.data = data
		return r, nil
	}
}

func WithTargetDir(targetDir string) RenderingContextOption {
	return func(r RenderingContext) (RenderingContext, error) {
		r.targetDir = targetDir
		r.templates = []string{targetDir}
		return r, nil
	}
}

func WithTemplateFiles(templateFiles ...string) RenderingContextOption {
	return func(r RenderingContext) (RenderingContext, error) {
		var err error
		var targetPath string

		for _, path := range templateFiles {
			start := strings.LastIndex(path, "/config/")

			if start != -1 {
				targetPath = filepath.Join(r.targetDir, path[start+8:])
			} else {
				start := strings.LastIndex(path, "/config")

				if start != -1 {
					targetPath = ""

					files, err := ioutil.ReadDir(path)
					if err != nil {
						return RenderingContext{}, err
					}

					for _, file := range files {
						err = copy.Copy(filepath.Join(path, file.Name()), filepath.Join(r.targetDir, file.Name()))
						if err != nil {
							return RenderingContext{}, err
						}
					}
				} else {
					start := strings.LastIndex(path, "/cf-for-k8s/")
					targetPath = filepath.Join(r.targetDir, path[start+12:])
				}
			}

			if targetPath != "" {
				err = copy.Copy(path, targetPath)
				if err != nil {
					return RenderingContext{}, err
				}
			}
		}

		return r, nil
	}
}

func WithValueFiles(valueFiles ...string) RenderingContextOption {
	return func(r RenderingContext) (RenderingContext, error) {
		for _, valueFile := range valueFiles {
			r.templates = append(r.templates, valueFile)
		}

		return r, nil
	}
}
