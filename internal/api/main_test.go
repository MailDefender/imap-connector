package api

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"

	"maildefender/imap-connector/internal/configuration"
	"maildefender/imap-connector/internal/mail/imap"
)

func TestMain(m *testing.M) {
	logrus.Info("Crafting integration test environment...")
	conf := configuration.MailServerConfiguration{
		Server: configuration.ServerConfiguration{
			Host: "localhost",
			Port: 31143,
			Tls:  false,
		},
		Authentication: configuration.AuthenticationConfiguration{
			Username: "test",
			Password: "test",
		},
	}

	logrus.WithFields(logrus.Fields{
		"host":     conf.Server.Host,
		"port":     conf.Server.Port,
		"tls":      conf.Server.Tls,
		"username": conf.Authentication.Username,
		"password": conf.Authentication.Password,
	}).Info("Starting integration test")

	if err := conf.Check(); err != nil {
		logrus.WithError(err).Fatal("imap configuration is invalid")
		os.Exit(1)
	}

	if err := imap.Dial(conf.Server); err != nil {
		logrus.WithError(err).Fatal("cannot dial imap server")
		os.Exit(1)
	}

	if err := imap.Login(conf.Authentication); err != nil {
		logrus.WithError(err).Fatal("cannot login to imap")
		os.Exit(1)
	}

	os.Exit(m.Run())
}
