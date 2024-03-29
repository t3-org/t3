package config

import (
	"fmt"
	"os"

	"github.com/kamva/tracer"
	"gopkg.in/yaml.v3"
)

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
	CommandPrefix  string `yaml:"command_prefix" yaml:"command_prefix"`
	OKEmoji        string `yaml:"ok_emoji" yaml:"ok_emoji"` // The emoji we use to set a command as done.
	PickleKey      string `yaml:"pickle_key" yaml:"pickle_key"`
	HomeServerAddr string `yaml:"home_server_addr" yaml:"home_server_addr"`
	Username       string `yaml:"username" yaml:"username"`
	Password       string `yaml:"password" yaml:"password"`
	// If you want to use multiple clients with the same DB,
	// you should set a distinct database account ID for each one.
	DBAccountID string `yaml:"db_account_id" yaml:"db_account_id"`
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
