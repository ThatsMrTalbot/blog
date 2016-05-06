package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/ThatsMrTalbot/blog"
)

func main() {
	logrus.Info("Opening config")

	config, err := blog.OpenConfig("config.json")
	if err != nil {
		logrus.WithError(err).Fatal("Could not open config")
	}

	logrus.Info("Initializing application")

	app, err := blog.NewApp(config)
	if err != nil {
		logrus.WithError(err).Fatal("Could not start application")
	}

	logrus.Info("Listening for requests")

	app.Serve()
}
