package api

import (
	"fmt"
	"net/http"

	"github.com/kessaro/shaker"
	"github.com/sirupsen/logrus"

	"maildefender/imap-connector/internal/locker"
	"maildefender/imap-connector/internal/mail/imap"
	"maildefender/imap-connector/internal/models"
)

type listMailboxesOut struct {
	Count     int              `json:"count"`
	Mailboxes models.Mailboxes `json:"mailboxes"`
}

func listMailboxes(c *shaker.Context) (listMailboxesOut, error) {
	if !locker.TryLock() {
		return listMailboxesOut{}, errFailedToLockImapMutex
	}
	defer locker.Unlock()

	mbox, err := imap.ListMailboxes(imap.GetMailboxesIn{
		Patterns: imap.AllMaibloxesPattern,
	})

	if err != nil {
		return listMailboxesOut{}, err
	}

	return listMailboxesOut{
		Count:     len(mbox),
		Mailboxes: mbox,
	}, nil
}

type getMailboxStatusIn struct {
	Mailbox string `form:"mailbox" binding:"required"`
}

func getMailboxStatus(c *shaker.Context, in *getMailboxStatusIn) (*models.MailboxStatus, error) {
	if !locker.TryLock() {
		return nil, errFailedToLockImapMutex
	}

	defer locker.Unlock()

	mboxes, err := imap.ListMailboxes(imap.GetMailboxesIn{Patterns: imap.AllMaibloxesPattern})
	if err != nil {
		return nil, err
	}

	if !mboxes.Has(in.Mailbox) {
		return nil, shaker.ErrRessourceNotFound
	}

	status, err := imap.MailboxStatus(in.Mailbox)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

type createMailboxStatusIn struct {
	Mailbox string `json:"mailbox" binding:"required"`
}

func createMailbox(c *shaker.Context, in *createMailboxStatusIn) error {

	if !locker.TryLock() {
		logrus.Error("cannot lock the imap client, maybe another operation is in progress...")
		return errFailedToLockImapMutex
	}

	defer locker.Unlock()

	mboxes, err := imap.ListMailboxes(imap.GetMailboxesIn{Patterns: imap.AllMaibloxesPattern})
	if err != nil {
		logrus.WithError(err).WithField("pattern", imap.AllMaibloxesPattern).Error("cannot list mailboxes")
		return fmt.Errorf("cannot list mailbox : %v", err)
	}

	if mboxes.Has(in.Mailbox) {
		logrus.WithField("mailbox", in.Mailbox).Error("this mailbox already exist")
		return errAlreadyExist
	}

	err = imap.CreateMailbox(in.Mailbox)
	locker.Unlock()
	if err != nil {
		logrus.WithError(err).WithField("mailbox", in.Mailbox).Error("cannot create mailbox")
		c.JSON(http.StatusInternalServerError, apiError{Error: err.Error()})
		return fmt.Errorf("cannot create mailbox : %v", err)
	}

	return nil
}
