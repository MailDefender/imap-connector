package api

import (
	"errors"
	"fmt"
	"math"

	"github.com/kessaro/shaker"
	"github.com/sirupsen/logrus"

	"maildefender/imap-connector/internal/locker"
	"maildefender/imap-connector/internal/mail/imap"
	types "maildefender/imap-connector/internal/models"
)

type fetchEmailIn struct {
	Mailbox string  `uri:"mailbox"`
	Limit   int     `form:"limit"`
	Offset  int     `form:"offset"`
	Sender  *string `form:"sender"`
}

type fetchEmailsOut struct {
	Count    int             `json:"count"`
	Messages []types.Message `json:"messages"`
}

func fetchMessages(c *shaker.Context, in *fetchEmailIn) (fetchEmailsOut, error) {
	if !locker.TryLock() {
		logrus.Warn("internal imap client is already in use, try later")
		return fetchEmailsOut{}, errFailedToLockImapMutex
	}

	defer locker.Unlock()

	if err := imap.SelectMailbox(in.Mailbox); err != nil {
		logrus.WithFields(logrus.Fields{"mailbox": in.Mailbox}).WithError(err).Error("cannot select mailbox")
		return fetchEmailsOut{}, fmt.Errorf("cannot select mailbox : %v", err)
	}

	// 0 offset = position 1
	in.Offset++

	if in.Limit == 0 {
		in.Limit = in.Offset + 11
	} else {
		in.Limit += in.Offset + int(math.Max(float64(in.Limit), 10))
	}
	status, err := imap.MailboxStatus(in.Mailbox)
	if err != nil {
		logrus.WithFields(logrus.Fields{"mailbox": in.Mailbox}).WithError(err).Error("cannot get mailbox status")
		return fetchEmailsOut{}, fmt.Errorf("cannot get mailblox status : %v", err)
	}

	if status.MessageCount == 0 || in.Offset > int(status.MessageCount) {
		return fetchEmailsOut{Count: 0}, nil
	}

	if in.Limit > int(status.MessageCount) {
		in.Limit = int(status.MessageCount)
	}

	msgs, err := imap.FetchMessages(imap.GetMessageIn{
		FetchHeaders: true,
		From:         in.Offset,
		To:           in.Limit,
		Sender:       in.Sender,
	})

	if err != nil {
		logrus.WithFields(logrus.Fields{"mailbox": in.Mailbox, "offset": in.Offset, "limit": in.Limit}).WithError(err).Error("cannot fetch messages")
		return fetchEmailsOut{}, fmt.Errorf("cannot fetch messages : %v", err)
	}

	logrus.WithFields(logrus.Fields{"message_count": len(msgs)}).Info("messages succesfully fetched")
	return fetchEmailsOut{
		Count:    len(msgs),
		Messages: msgs,
	}, nil
}

type MoveEmailIn struct {
	Source                string `uri:"mailbox"`
	MessageID             string `uri:"messageId" `
	Destination           string `json:"destination" binding:"required"`
	CreateMailboxIfNeeded bool   `json:"createMailbox"`
}

func moveMessage(c *shaker.Context, in *MoveEmailIn) error {
	if !locker.TryLock() {
		logrus.Error("cannot lock the imap client, maybe another operation is in progress...")
		return errFailedToLockImapMutex
	}

	defer locker.Unlock()

	{
		mboxes, err := imap.ListMailboxes(imap.GetMailboxesIn{
			Patterns: imap.AllMaibloxesPattern,
		})

		if err != nil {
			logrus.WithError(err).Error("cannot list mailboxes")
			return err
		}

		if !mboxes.Has(in.Source) {
			logrus.WithField("mailbox", in.Source).Error("this mailbox does not exist")
			return shaker.ErrRessourceNotFound
		}

		if !mboxes.Has(in.Destination) {
			if in.CreateMailboxIfNeeded {
				if err := imap.CreateMailbox(in.Destination); err != nil {
					logrus.WithError(err).Error("cannot create mailbox")
					return errors.New("cannot create destination mailbox")
				}
			} else {
				logrus.WithField("mailbox", in.Destination).Error("this mailbox does not exist")
				return shaker.ErrRessourceNotFound
			}
		}
	}

	if err := imap.SelectMailbox(in.Source); err != nil {
		logrus.WithField("mailbox", in.Source).Error("cannot select this mailbox")
		return errors.New("cannot select this mailbox")
	}

	msg, err := imap.SearchMessageFromID(in.Source, in.MessageID)
	if err != nil {
		logrus.WithField("message_id", in.MessageID).Error("cannot find this message")
		return shaker.ErrRessourceNotFound
	}

	err = imap.MoveMessage(msg, in.Destination)

	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"message_id":  in.MessageID,
			"source":      in.Source,
			"destination": in.Destination,
		}).Error("cannot move message")
		return fmt.Errorf("cannot move message : %v", err)
	}

	return nil
}
