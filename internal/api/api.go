package api

import (
	"net/http"

	"github.com/kessaro/shaker"
)

var api = shaker.NewShaker(&shaker.MappedErrors{
	errFailedToLockImapMutex: http.StatusTooEarly,
	errAlreadyExist:          http.StatusConflict,
})

func init() {
	api.Get("/v1/mailbox", listMailboxes, http.StatusOK)
	api.Get("/v1/mailbox/status", getMailboxStatus, http.StatusOK)
	api.Get("/v1/mailbox/:mailbox/message", fetchMessages, http.StatusOK)
	api.Post("/v1/mailbox/:mailbox/message/:messageId/move", moveMessage, http.StatusOK)

	api.Post("/v1/mailbox", createMailbox, http.StatusCreated)
}

func Run() error {
	return api.Shake()
}
