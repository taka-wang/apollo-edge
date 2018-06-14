package boltstore

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/takawang/sugar"
)

func reset() {
	os.Clearenv()
	defaultConfig = nil
	files, _ := filepath.Glob("*.log")
	for _, f := range files {
		os.Remove(f)
	}
}

// Test Cases =========================
func TestInit(t *testing.T) {
	s := sugar.New(t)

	s.Assert("Load default config", func(logf sugar.Log) bool {
		defer reset()

		loadConfig()
		logf("current config: %+v", defaultConfig)
		if defaultConfig.BoltStore.Logger.Filename == "boltstore.log" {
			return true
		}
		return false
	})

	s.Assert("Load config from environment variable", func(logf sugar.Log) bool {
		defer reset()

		os.Setenv("BOLTSTORE_CONF", "testdata/config.json")
		loadConfig()
		logf("current config: %+v", defaultConfig)
		if defaultConfig.BoltStore.Logger.Filename == "boltstore.test.log" {
			return true
		}
		return false
	})

	s.Assert("Set default logger", func(logf sugar.Log) bool {
		defer reset()

		loadConfig()
		logf("current config: %+v", defaultConfig)
		setLogger()
		logf("current logger: %+v", log.StandardLogger())
		if log.GetLevel() == log.InfoLevel {
			return true
		}
		return false
	})

	s.Assert("Set test logger from environment variable", func(logf sugar.Log) bool {
		defer reset()

		os.Setenv("BOLTSTORE_CONF", "testdata/config.json")
		loadConfig()
		logf("current config: %+v", defaultConfig)
		setLogger()
		logf("current logger: %+v", log.StandardLogger())

		if _, err := os.Stat("boltstore.test.log"); os.IsNotExist(err) {
			logf("check if log file exists: %v", err)
			return false
		}
		return true
	})

	if s.IsFailed() {
		fmt.Println("the tests failed :/")
	}
}
