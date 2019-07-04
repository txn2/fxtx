package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/txn2/fxtx"

	"go.uber.org/zap"
)

var (
	configFileEnv = getEnv("CONFIG", "./cfg/example.yml")
	debugEnv      = getEnv("DEBUG", "false")
	destHostEnv   = getEnv("DEST", "127.0.0.1:30000")
	timeoutEnv    = getEnv("TIMEOUT", "10")
)

func main() {

	debugEnvBool := false
	if debugEnv == "true" {
		debugEnvBool = true
	}

	debugEnvInt, err := strconv.Atoi(timeoutEnv)
	if err != nil {
		fmt.Printf("TIMEOUT must be an integer: %s\n", err.Error())
	}

	var (
		configFile = flag.String("config", configFileEnv, "Config file")
		debug      = flag.Bool("debug", debugEnvBool, "Debug logging mode")
		dest       = flag.String("dest", destHostEnv, "Destination host")
		timeout    = flag.Int("tcpTimeout", debugEnvInt, "TCP Timeout")
	)

	flag.Parse()

	zapCfg := zap.NewProductionConfig()
	zapCfg.DisableCaller = true
	zapCfg.DisableStacktrace = true

	if *debug == true {
		zapCfg = zap.NewDevelopmentConfig()
	}

	logger, err := zapCfg.Build()
	if err != nil {
		fmt.Printf("Can not build logger: %s\n", err.Error())
		os.Exit(1)
	}

	logger.Info("Loading configuration...")
	genCfg, err := fxtx.GenCfgFromFile(*configFile)
	if err != nil {
		logger.Fatal("Config file error", zap.Error(err))
	}

	fxtxApi, err := fxtx.NewFxtx(&fxtx.Cfg{
		GenCfg:      genCfg,
		Destination: *dest,
		Timeout:     time.Duration(*timeout) * time.Second,
		Logger:      logger,
	})
	if err != nil {
		logger.Fatal("Error instantiating Fxtx", zap.Error(err))
	}

	logger.Info("Run Fxtx...")
	fxtxApi.Run()

}

// getEnv gets an environment variable or sets a default if
// one does not exist.
func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}

	return value
}
