package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	"notifier/internal"
)

var configFilePath = flag.String("config", "./config.yaml", "configuration file path")

func main() {
	logLvl := os.Getenv("APP_LOG_LEVEL")
	if logLvl == "" {
		logLvl = "info"
	}
	parsedLvl, parsedLvlErr := logrus.ParseLevel(logLvl)
	if parsedLvlErr != nil {
		logrus.WithField("error", parsedLvl).Infoln("unknown log level, set default log level - info")
		parsedLvl = logrus.InfoLevel
	}

	logrus.SetLevel(parsedLvl)

	ctx, stop := context.WithCancel(context.Background())
	signal.NotifyContext(ctx, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGABRT)
	defer func() {
		if err := recover(); err != nil {
			logrus.WithField("error", err).Errorln("application critical error")
			stop()
		}
		logrus.Infoln("application has been stopped")
	}()

	logrus.Infoln("application init")
	flag.Parse()
	app, initAppErr := internal.NewApp(ctx, *configFilePath)
	if initAppErr != nil {
		logrus.WithField("error", initAppErr).Errorln("application init error")
		return
	}
	logrus.Infoln("application is starting")

	if appRunErr := app.Run(); appRunErr != nil {
		logrus.WithField("error", appRunErr).Errorln("application run error")
	}
}
