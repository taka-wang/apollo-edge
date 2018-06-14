package boltstore

import (
	"os"

	"github.com/koding/multiconfig"
	log "github.com/sirupsen/logrus"
)

// envCONF keeps the config file path
const envCONF = "BOLTSTORE_CONF"

func loadConfig() {
	// try to read config file path from the EV
	confPath := os.Getenv(envCONF)
	confLoader := multiconfig.NewWithPath(confPath)
	defaultConfig = new(ConfigType)
	confLoader.Load(defaultConfig)
	//loader.MustLoad(defaultConfig)
}

func setLogger() {
	if defaultConfig == nil {
		return
	}

	// set logger level: [Panic : 0, Fatal : 1, Error : 2, Warn  : 3, Info  : 4, Debug : 5]
	if defaultConfig.BoltStore.Logger.Level >= 0 || defaultConfig.BoltStore.Logger.Level <= 5 {
		log.SetLevel(log.Level(defaultConfig.BoltStore.Logger.Level))
	} else {
		// default logger level
		log.SetLevel(log.InfoLevel)
	}

	// set logger formatter
	if defaultConfig.BoltStore.Logger.JSON {
		// Log as JSON instead of the default ASCII formatter.
		log.SetFormatter(&log.JSONFormatter{})
	}

	// set logger output
	if defaultConfig.BoltStore.Logger.ToFile && len(defaultConfig.BoltStore.Logger.Filename) > 0 {
		// try to create log file
		logfile, err := os.OpenFile(defaultConfig.BoltStore.Logger.Filename, os.O_WRONLY|os.O_CREATE, 0755)
		if err == nil {
			log.SetOutput(logfile)
			return
		}
	}
	// Output to stdout instead of the default stderr
	log.SetOutput(os.Stdout)
}

func init() {
	loadConfig()
	setLogger()
}
