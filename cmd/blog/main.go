package main

import (
	"flag"

	"github.com/Sirupsen/logrus"
	"github.com/ThatsMrTalbot/blog"
)

var config blog.Config

func init() {
	flag.StringVar(&config.Listen, "http", ":8080", "Port to listen on")
	flag.StringVar(&config.Path, "path", "example.git", "Path to git repository")
}

func main() {
	flag.Parse()

	logrus.Info("Initializing application")

	app, err := blog.NewApp(&config)
	if err != nil {
		logrus.WithError(err).Fatal("Could not start application")
	}

	logrus.Info("Listening for requests")

	app.Serve()
}
