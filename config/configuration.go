package config

import (
	logging "github.com/op/go-logging"
	"github.com/spf13/viper"
)

var logger = logging.MustGetLogger("config")

// Config is a configuration structure
type Config struct {
	VNS          VNS
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

type VNS struct {
	IterMax int
}

type Optimization struct {
	Strategy string
	IterMax  int
	LevelMax int
}

// LoadConfig loads configuration file
func (c *Config) LoadConfig(conf string) (err error) {
	viper.SetConfigFile(conf)

	if err = viper.ReadInConfig(); err != nil {
		logger.Errorf("Couldn't load file %s\n", conf)
		return
	}

	if err = viper.UnmarshalKey("vns", &c.VNS); err != nil {
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
