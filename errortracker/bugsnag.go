package errortracker

import (
	"os"

	logrus_bugsnag "github.com/Shopify/logrus-bugsnag"
	"github.com/Sirupsen/logrus"
	bugsnag "github.com/bugsnag/bugsnag-go"
)

func init() {
	apiKey := os.Getenv("BUGSNAG_API_KEY")
	if apiKey == "" {
		return
	}

	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "development"
	}

	bugsnag.Configure(bugsnag.Configuration{
		APIKey:       apiKey,
		ReleaseStage: environment,
	})

	hook, err := logrus_bugsnag.NewBugsnagHook()
	if err != nil {
		logrus.WithField("err", err).Fatal("Error creating new Bugsnag hook for logrus")
	}
	logrus.StandardLogger().Hooks.Add(hook)
}
