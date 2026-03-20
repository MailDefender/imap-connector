package imap

import (
	"errors"
	"fmt"
	"slices"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/sirupsen/logrus"

	"maildefender/imap-connector/internal/configuration"
	"maildefender/imap-connector/internal/hash"
	"maildefender/imap-connector/internal/mail/parser"
	"maildefender/imap-connector/internal/models"
	"maildefender/imap-connector/internal/utils"
)

var client *imapclient.Client

func Dial(conf configuration.ServerConfiguration) error {
	addr := fmt.Sprintf("%s:%d", conf.Host, conf.Port)

	var err error
	if conf.Tls {
		client, err = imapclient.DialTLS(addr, &imapclient.Options{})
	} else {
		client, err = imapclient.DialInsecure(addr, &imapclient.Options{})
	}
	return err
}

func Login(auth configuration.AuthenticationConfiguration) error {
	return client.Login(auth.Username, auth.Password).Wait()
}

func Logout() error {
	return client.Logout().Wait()
}

func mailboxes(pattern string) (models.Mailboxes, error) {
	listCmd := client.List("", pattern, &imap.ListOptions{
		ReturnStatus: &imap.StatusOptions{
			NumMessages: true,
			NumUnseen:   true,
		},
	})

	var mailboxes models.Mailboxes

	for {
		mbox := listCmd.Next()
		if mbox == nil {
			break
		}

		messagesCount, unseenCount := 0, 0
		if mbox.Status != nil {
			if mbox.Status.NumMessages != nil {
				messagesCount = int(*mbox.Status.NumMessages)
			}
			if mbox.Status.NumUnseen != nil {
				unseenCount = int(*mbox.Status.NumUnseen)
			}
		}

		mailboxes = append(mailboxes, models.Mailbox{
			Name:         mbox.Mailbox,
			MessageCount: messagesCount,
			UnseenCount:  unseenCount,
		})
	}

	if err := listCmd.Close(); err != nil {
		return mailboxes, err
	}

	return mailboxes, nil
}

func ListMailboxes(in GetMailboxesIn) (models.Mailboxes, error) {
	var mboxes models.Mailboxes

	for _, pattern := range in.Patterns {
		mbx, err := mailboxes(pattern)
		if err != nil {
			return nil, err
		}
		for _, m := range mbx {
			mboxes = append(mboxes, m)
		}
	}

	return mboxes, nil
}

func MailboxStatus(mailbox string) (models.MailboxStatus, error) {
	options := imap.StatusOptions{NumMessages: true, NumUnseen: true}
	data, err := client.Status(mailbox, &options).Wait()
	if err != nil {
		return models.MailboxStatus{}, err
	}
	return models.MailboxStatus{
		MessageCount: uint(*data.NumMessages),
	}, nil
}

func SelectMailbox(mailbox string) error {
	if _, err := client.Select(mailbox, nil).Wait(); err != nil {
		return err
	}

	return nil
}

func FetchMessages(in GetMessageIn) ([]models.Message, error) {
	if err := checkParams(in); err != nil {
		return nil, err
	}

	var messages []models.Message

	var seqSet imap.SeqSet
	seqSet.AddRange(uint32(in.From), uint32(in.To))

	msgs, err := client.Fetch(
		seqSet,
		&imap.FetchOptions{
			UID:        true,
			Envelope:   in.FetchHeaders,
			Flags:      in.FetchHeaders,
			RFC822Size: in.FetchHeaders,
			ModSeq:     false,
			BodySection: []*imap.FetchItemBodySection{
				{
					Specifier: imap.PartSpecifierHeader,
					Peek:      true,
				},
			},
		},
	).Collect()

	if err != nil {
		return nil, err
	}

	for _, msg := range msgs {
		message := models.Message{}

		for _, from := range msg.Envelope.From {
			message.From = append(message.From, models.Contact{
				Name:  from.Name,
				Email: fmt.Sprintf("%s@%s", from.Mailbox, from.Host),
				Host:  from.Host,
			})
		}
		for _, from := range msg.Envelope.To {
			message.To = append(message.To, models.Contact{
				Name:  from.Name,
				Email: fmt.Sprintf("%s@%s", from.Mailbox, from.Host),
				Host:  from.Host,
			})
		}
		message.Date = msg.Envelope.Date
		if msg.Envelope.MessageID != "" {
			message.MessageID = msg.Envelope.MessageID
		} else {
			h, err := hash.HashMessage(message, hash.LatestVersion)
			if err != nil {
				logrus.WithError(err).Error("this message has an empty Message-ID and cannot generate associated hash, skipping this one...")
				continue
			}
			message.MessageID = h.String()
		}
		message.Subject = msg.Envelope.Subject
		message.SequenceIndex = msg.SeqNum

		headers := msg.FindBodySection(&imap.FetchItemBodySection{Specifier: imap.PartSpecifierHeader})
		parser := parser.NewHeaderParser(string(headers))

		if err := parser.Parse(); err != nil {
			logrus.WithError(err).Error("cannot parse message's headers")
		}

		message.Headers = parser.GetAll()

		if in.Sender == nil || (in.Sender != nil && slices.ContainsFunc(message.From, func(contact models.Contact) bool {
			return contact.Email == *in.Sender
		})) {
			messages = append(messages, message)
		}
	}

	return messages, nil
}

func CreateMailbox(mailbox string) error {
	return client.Create(mailbox, nil).Wait()
}

func checkParams(in GetMessageIn) error {
	if in.From < 0 || in.To < 0 || in.From > in.To {
		return errors.New("invalid parameters")
	}
	return nil
}

func MoveMessage(message models.Message, mailboxName string) error {
	_, err := client.Move(
		imap.SeqSetNum(message.SequenceIndex),
		mailboxName).
		Wait()
	return err
}

func SearchMessageFromID(mailbox string, messageID string) (models.Message, error) {
	status, err := MailboxStatus(mailbox)
	if err != nil {
		return models.Message{}, err
	}

	for index := 1; index < int(status.MessageCount)+1; index += 10 {
		from, to := index, utils.MinOf(index+10, int(status.MessageCount))
		messages, err := FetchMessages(GetMessageIn{
			FetchHeaders: true,
			From:         from,
			To:           to,
		})

		if err != nil {
			return models.Message{}, err
		}

		for _, msg := range messages {
			hashFromID, err := hash.Parse(msg.MessageID)

			// Message-ID seems to be an hash
			// Trying to determine if this hash and the message's hash match
			if err == nil {
				hashFromMsg, err := hash.HashMessage(msg, hashFromID.Version())
				if err == nil && hashFromMsg.Match(hashFromID) {
					return msg, nil
				}
			}

			if msg.MessageID == messageID {
				return msg, nil
			}
		}
	}

	return models.Message{}, errors.New("message not found")

}
