package main

type cfgModelMock struct {
	data map[string]string
}

func (c *cfgModelMock) Get(name string) string {
	if v, ok := c.data[name]; ok {
		return v
	}
	return ""
}
