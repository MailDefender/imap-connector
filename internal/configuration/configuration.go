package configuration

import (
	"maildefender/imap-connector/internal/utils"
)

var imapConfiguration MailServerConfiguration = MailServerConfiguration{
	Server: ServerConfiguration{
		Host: utils.GetEnvString("IMAP_HOST", ""),
		Port: utils.GetEnvInt("IMAP_PORT", 0),
		Tls:  utils.GetEnvBool("IMAP_TLS", true),
	}, Authentication: AuthenticationConfiguration{
		Username: utils.GetEnvString("IMAP_USERNAME", ""),
		Password: utils.GetEnvString("IMAP_PASSWORD", ""),
	},
}

func ImapConfiguration() MailServerConfiguration {
	return imapConfiguration
}
