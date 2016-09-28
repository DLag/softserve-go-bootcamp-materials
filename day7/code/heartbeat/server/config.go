package main

import "os"

type cfgModel interface {
	Get(string) string
}

type cfgModelEnv struct {
}

func (c *cfgModelEnv) Get(name string) string {
	return os.Getenv(name)
}

func newCfgModelEnv() *cfgModelEnv {
	return new(cfgModelEnv)
}
