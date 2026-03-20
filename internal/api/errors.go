package api

import "errors"

type apiError struct {
	Error string `json:"error"`
}

var (
	errFailedToLockImapMutex error = errors.New("imap connection not available, retry later")
	errAlreadyExist          error = errors.New("already exist")
	errMailboxNotFound       error = errors.New("mailbox not found")
)
