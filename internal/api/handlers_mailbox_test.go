package api

import (
	"net/http"
	"testing"

	"github.com/maxatome/go-testdeep/helpers/tdhttp"

	"maildefender/imap-connector/internal/models"
)

func TestListMailbox(t *testing.T) {
	tests := []struct {
		name                string
		endpoint            string
		expectedStatusCode  int
		expectedJsonPayload listMailboxesOut
	}{
		{
			name:               "list all mailboxes",
			endpoint:           "/v1/mailbox",
			expectedStatusCode: http.StatusOK,
			expectedJsonPayload: listMailboxesOut{
				Count: 2,
				Mailboxes: models.Mailboxes{
					models.Mailbox{
						Name:         "recovered-lost-folder-d8de382e758459691200000052c9e01c",
						MessageCount: 2,
						UnseenCount:  2,
					},
					models.Mailbox{
						Name:         "INBOX",
						MessageCount: 0,
						UnseenCount:  0,
					},
				},
			},
		},
	}

	tester := tdhttp.NewTestAPI(t, api.Handle())
	for _, test := range tests {
		tt := tester.Name(test.name).Get(test.endpoint).CmpStatus(test.expectedStatusCode)

		tt.CmpJSONBody(test.expectedJsonPayload)

	}
}
