package utils

import "os"

type ConfigEntry struct {
	Value       *string
	EnvVar      string
	Description string
	Default     string
}

func NewConfigEntry(c *Config, value *string, env string, desc string, defval string) (e ConfigEntry) {
	e = ConfigEntry{}
	c.entries[env] = &e
	e.Value = value
	e.EnvVar = env
	e.Description = desc
	e.Default = defval
	if v, ok := os.LookupEnv(e.EnvVar); ok {
		*e.Value = v
	} else {
		*e.Value = e.Default
	}
	return
}

type Config struct {
	entries       map[string]*ConfigEntry
	ListenAddress string
}

func (c *Config) init() {
	c.entries = make(map[string]*ConfigEntry)
	NewConfigEntry(c, &c.ListenAddress, "GO_LISTEN_ADDRESS", "Service listen address", ":3000")
}

func GetConfig() (c *Config) {
	c = &Config{}
	c.init()
	return
}
