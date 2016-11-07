package model

import (
	"fmt"
)

type Config struct {
	Port                Port
	Prefix              Prefix
	ApplicationPassword ApplicationPassword
	LedisdbAddress      LedisdbAddress
	LedisdbPassword     LedisdbPassword
}

func (c *Config) Validate() error {
	if c.Port <= 0 {
		return fmt.Errorf("parameter Port invalid")
	}
	if len(c.ApplicationPassword) == 0 {
		return fmt.Errorf("parameter ApplicationPassword invalid")
	}
	if len(c.LedisdbAddress) == 0 {
		return fmt.Errorf("parameter LedisdbAddress invalid")
	}
	return nil
}
