package models

import (
	"time"
)

const LATEST_HASH_VERSION = 1

type Message struct {
	MessageID     string    `json:"messageId"`
	SequenceIndex uint32    `json:"-"`
	Headers       Headers   `json:"headers"`
	From          Contacts  `json:"from,omitempty"`
	To            Contacts  `json:"to,omitempty"`
	Cc            Contacts  `json:"cc,omitempty"`
	Subject       string    `json:"subject,omitempty"`
	Date          time.Time `json:"date"`
}
