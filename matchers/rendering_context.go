package matchers

type RenderingContext struct {
	templates []string
	data      map[string]interface{}
}

func (r RenderingContext) WithData(data map[string]interface{}) RenderingContext {
	r.data = data
	return r
}

func NewRenderingContext(templates ...string) RenderingContext {
	return RenderingContext{templates, nil}
}
