package models

type Mailbox struct {
	Name         string `json:"name"`
	MessageCount int    `json:"messageCount"`
	UnseenCount  int    `json:"unseenCount"`
}

type Mailboxes []Mailbox

func (mboxes Mailboxes) Has(mailbox string) bool {
	for _, mbx := range mboxes {
		if mbx.Name == mailbox {
			return true
		}
	}

	return false
}
