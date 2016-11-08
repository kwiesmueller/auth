package model

import (
	"fmt"
)

type Config struct {
	HttpPort            Port
	HttpPrefix          Prefix
	ApplicationName     ApplicationName
	ApplicationPassword ApplicationPassword
	LedisdbAddress      LedisdbAddress
	LedisdbPassword     LedisdbPassword
}

func (c *Config) Validate() error {
	if c.HttpPort <= 0 {
		return fmt.Errorf("parameter Port invalid")
	}
	if len(c.ApplicationName) == 0 {
		return fmt.Errorf("parameter ApplicationName invalid")
	}
	if len(c.ApplicationPassword) == 0 {
		return fmt.Errorf("parameter ApplicationPassword invalid")
	}
	if len(c.LedisdbAddress) == 0 {
		return fmt.Errorf("parameter LedisdbAddress invalid")
	}
	return nil
}
