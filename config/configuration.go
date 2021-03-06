package config

import (
	"time"

	logging "github.com/op/go-logging"
	"github.com/spf13/viper"
)

var logger = logging.MustGetLogger("config")

// Config is a configuration structure
type Config struct {
	Common       Common
	Construction Construction
	Optimization Optimization
}

type Construction struct {
	Strategy string
	LevelMax int
	Penalty  Penalty
}
type Penalty struct {
	TimeWindows    int
	PickupDelivery int
	Capacity       int
}

type Common struct {
	IterMax int
	MaxTime time.Duration
}

type Optimization struct {
	Objective string
	Asymetric bool
	VNS       VNS
	SA        SA
}

type VNS struct {
	IterMax     int
	LevelMax    int
	LocalSearch LocalSearch
}

type SA struct {
	IterMax     float64
	LocalSearch LocalSearch
}

type LocalSearch string

const (
	Const2Opt LocalSearch = "2opt"
	Shifting  LocalSearch = "shifting"
	VND       LocalSearch = "vnd"
)

// LoadConfig loads configuration file
func (c *Config) LoadConfig(conf string) (err error) {
	viper.SetConfigFile(conf)

	if err = viper.ReadInConfig(); err != nil {
		logger.Errorf("Couldn't load file %s\n", conf)
		return
	}

	if err = viper.UnmarshalKey("common", &c.Common); err != nil {
		logger.Error(err.Error())
		return
	}

	if err = viper.UnmarshalKey("construction", &c.Construction); err != nil {
		logger.Error(err.Error())
		return
	}

	if err = viper.UnmarshalKey("optimization", &c.Optimization); err != nil {
		logger.Error(err.Error())
		return
	}

	logger.Info("Configuration loaded successfully")
	return
}
