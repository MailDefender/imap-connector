package imap

type MailboxPattern []string

var (
	AllMaibloxesPattern MailboxPattern = []string{
		"%",
		"%.*",
	}
)
