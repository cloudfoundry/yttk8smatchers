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

func NewRenderingContext(targetDir string, valuesFiles ...string) RenderingContext {
	templates := []string{targetDir}

	for _, valuesFile := range valuesFiles {
		templates = append(templates, valuesFile)
	}

	return RenderingContext{targetDir, templates, nil}
}

func (r RenderingContext) WithData(data map[string]interface{}) RenderingContext {
	r.data = data
	return r
}

func (r RenderingContext) CopyTemplatesToTargetDir(templates ...string) error {
	var err error
	var targetPath string

	for _, path := range templates {
		start := strings.LastIndex(path, "/config/")

		if start != -1 {
			targetPath = filepath.Join(r.targetDir, path[start+8:])
		} else {
			start := strings.LastIndex(path, "/config")

			if start != -1 {
				targetPath = ""

				files, err := ioutil.ReadDir(path)
				if err != nil {
					return err
				}

				for _, file := range files {
					err = copy.Copy(filepath.Join(path, file.Name()), filepath.Join(r.targetDir, file.Name()))
					if err != nil {
						return err
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
				return err
			}
		}
	}

	return nil
}
