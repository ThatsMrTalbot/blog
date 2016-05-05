package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/ThatsMrTalbot/blog"
)

func main() {
	config, err := blog.OpenConfig("config.json")
	if err != nil {
		logrus.WithError(err).Fatal("Could not start application")
	}

	app, err := blog.NewApp(config)
	if err != nil {
		logrus.WithError(err).Fatal("Could not start application")
	}

	app.Serve()
}
