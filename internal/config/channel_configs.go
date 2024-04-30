package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/kamva/tracer"
	"gopkg.in/yaml.v3"
)

// ChannelsEnvSuffix is the suffix that we use to specify our channel config should be read from env.
// e.g., "__ENV__USERNAME" means we should read the field's value from the "USERNAME" environment variable.
const ChannelsEnvSuffix = "__ENV__"

type ChannelsConfig struct {
	ChannelHomes map[string]ChannelHome `yaml:"channel_homes"` // map's key is the client name.
	Channels     map[string]Channel     `yaml:"channels"`
	Policies     []Policy               `yaml:"policies"`
}

func (c *ChannelsConfig) Validate() error {
	var homeNames []string
	var channelNames []string
	for k, _ := range c.ChannelHomes {
		homeNames = append(homeNames, k)
	}
	for k, ch := range c.Channels {
		channelNames = append(channelNames, k)

		// Check if the home is in the homes list.
		if _, ok := c.ChannelHomes[ch.Home]; !ok {
			return fmt.Errorf("home %s not found in the channel homes config", ch.Home)
		}
	}

	for _, p := range c.Policies {
		for _, ch := range p.Channels {
			if _, ok := c.Channels[ch]; !ok {
				return fmt.Errorf("channel %s not found in the config", ch)
			}
		}
	}

	return nil
}

type MatrixHomeConfig struct {
	// CommandPrefix is the prefix of each command. e.g., "!t"
	CommandPrefix  string `json:"command_prefix" yaml:"command_prefix"`
	OKEmoji        string `json:"ok_emoji" yaml:"ok_emoji"` // The emoji we use to set a command as done.
	HomeServerAddr string `json:"home_server_addr" yaml:"home_server_addr"`
	IdentifierType string `json:"identifier_type" yaml:"identifier_type"`
	Medium         string `json:"medium" yaml:"medium"`
	Username       string `json:"username" yaml:"username"`
	Address        string `json:"address" yaml:"address"` // The email address if medium field is email.
	Password       string `json:"password" yaml:"password"`
}

func (c *MatrixHomeConfig) ResolveEnvs(envSuffix string) {
	// Resolve env variables
	if strings.HasPrefix(c.HomeServerAddr, envSuffix) {
		c.HomeServerAddr = os.Getenv(strings.TrimPrefix(c.HomeServerAddr, envSuffix))
	}

	if strings.HasPrefix(c.Username, envSuffix) {
		c.Username = os.Getenv(strings.TrimPrefix(c.Username, envSuffix))
	}

	if strings.HasPrefix(c.Address, envSuffix) {
		c.Address = os.Getenv(strings.TrimPrefix(c.Address, envSuffix))
	}

	if strings.HasPrefix(c.Password, envSuffix) {
		c.Password = os.Getenv(strings.TrimPrefix(c.Password, envSuffix))
	}
}

type ChannelHome struct {
	Type   string    `yaml:"type"`
	Config yaml.Node `yaml:"config"` // The raw config of the client.
}

type Channel struct {
	Home   string    `yaml:"home"`
	Config yaml.Node `yaml:"config"`
}

type Policy struct {
	Channels []string          `yaml:"channels"`
	Labels   map[string]string `yaml:"labels"`
}

func LoadChannelsConfig(fname string) (*ChannelsConfig, error) {
	b, err := os.ReadFile(fname)
	if err != nil {
		return nil, tracer.Trace(err)
	}
	var cfg ChannelsConfig
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, tracer.Trace(err)
	}
	return &cfg, nil
}
