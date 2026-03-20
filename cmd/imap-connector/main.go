package main

import (
	"github.com/sirupsen/logrus"

	"maildefender/imap-connector/internal/api"
	"maildefender/imap-connector/internal/configuration"
	"maildefender/imap-connector/internal/mail/imap"
)

func main() {
	logrus.Info("Starting imap-connector...")

	logrus.Info("Retrieving configuration")
	conf := configuration.ImapConfiguration()

	if err := conf.Check(); err != nil {
		logrus.WithError(err).Error("imap configuration is invalid")
		return
	}

	if err := imap.Dial(conf.Server); err != nil {
		logrus.WithError(err).Error("cannot dial imap server")
		return
	}

	if err := imap.Login(conf.Authentication); err != nil {
		logrus.WithError(err).Error("cannot login to imap")
		return
	}

	if err := api.Run(); err != nil {
		logrus.WithError(err).Error("API stopped")
	}
}
