package env

import (
	"os"

	"github.com/caarlos0/env/v11"
)

type Options struct {
	Environment     map[string]string
	TagName         string
	RequiredIfNoDef bool
	OnSet           env.OnSetFn
	Prefix          string
}

func Parse(v interface{}, opts ...Options) error {
	altOpts := []env.Options{}

	for _, opt := range opts {
		altOpts = append(altOpts, env.Options{
			Environment:     opt.Environment,
			TagName:         opt.TagName,
			RequiredIfNoDef: opt.RequiredIfNoDef,
			OnSet:           opt.OnSet,
			Prefix:          opt.Prefix,
		})
	}

	return env.Parse(v)
}

func GetEnvVar(key, fall string) string {
	val := os.Getenv(key)

	if val == "" {
		return fall
	}

	return val
}
