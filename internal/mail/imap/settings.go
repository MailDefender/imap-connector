package imap

type GetMessageIn struct {
	FetchHeaders bool
	Unread       bool
	From         int
	To           int
	Sender       *string
}

type GetMailboxesIn struct {
	Patterns MailboxPattern
}
