package config

import "os"

type Environ struct {
	configs map[string]string
}

func Env() *Environ {
	configs := make(map[string]string)
	for key, value := range defaults() {
		configs[key] = value
		v, found := os.LookupEnv(key)
		if found {
			configs[key] = v
		}
	}
	return &Environ{
		configs: configs,
	}
}

func (e *Environ) Get(key string) string {
	return e.configs[key]
}