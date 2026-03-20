package api

import (
	"net/http"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/helpers/tdhttp"
	"github.com/maxatome/go-testdeep/td"
)

func parseDateTime(in string) time.Time {
	t, _ := time.ParseInLocation(time.RFC1123Z, in, time.Local)
	return t
}

func TestFetchMessage(t *testing.T) {
	tests := []struct {
		name                string
		endpoint            string
		expectedStatusCode  int
		expectedJsonPayload string
	}{
		{
			name:               "list all messages",
			endpoint:           "/v1/mailbox/recovered-lost-folder-d8de382e758459691200000052c9e01c/message",
			expectedStatusCode: http.StatusOK,
			expectedJsonPayload: `{
			"count": 2,
			"messages": [
				{
					"messageId": "html1234567890@example.org",
					"from": [
						{
							"host":  "example.org",
							"name":  "Test Sender",
							"email": "test.sender@example.org",
						}
					],
					"headers": {
						"MIME-Version": ["1.0"],
						"Content-Type": ["text/html; charset=\"utf-8\""],
						"Content-Transfer-Encoding": ["quoted-printable"],
						"Received-SPF": ["pass (example.org: domain of test.sender@example.org designates 192.0.2.1 as permitted sender) client-ip=192.0.2.1;"],
						"Authentication-Results": ["example.org; spf=pass smtp.mailfrom=test.sender@example.org; dmarc=pass header.from=example.org"],
						"Message-ID": ["<html1234567890@example.org>"],
						"Date":       ["Sat, 3 Jan 2026 15:00:00 +0100"],
						"Subject":    ["Test Email - HTML"],
						"From":       ["\"Test Sender\" <test.sender@example.org>"],
						"To":         ["\"Recipient Test\" <recipient.test@example.com>"]
					},
					"to": [
						{
							"host":  "example.com",
							"name":  "Recipient Test",
							"email": "recipient.test@example.com",
						}
					],
					"subject": "Test Email - HTML",
					"date": NotEmpty()
				},
				{
					"messageId": "plain1234567890@example.org",
					"from": [
						{
							"host":  "example.org",
							"name":  "Test Sender",
							"email": "test.sender@example.org",
						}
					],
					"headers": {
						"MIME-Version": ["1.0"],
						"Content-Type": ["text/plain; charset=\"utf-8\""],
						"Content-Transfer-Encoding": ["quoted-printable"],
						"Received-SPF": ["pass (example.org: domain of test.sender@example.org designates 192.0.2.1 as permitted sender) client-ip=192.0.2.1;"],
						"Authentication-Results": ["example.org; spf=pass smtp.mailfrom=test.sender@example.org; dmarc=pass header.from=example.org"],
						"Message-ID": ["<plain1234567890@example.org>"],
						"Date":       ["Sat, 3 Jan 2026 15:00:00 +0100"],
						"Subject":    ["Test Email - Plain Text"],
						"From":       ["\"Test Sender\" <test.sender@example.org>"],
						"To":         ["\"Recipient Test\" <recipient.test@example.com>"]
					},
					"to": [
						{
							"host":  "example.com",
							"name":  "Recipient Test",
							"email": "recipient.test@example.com",
						}
					],
					"subject": "Test Email - Plain Text",
					"date": NotEmpty()
				}
			]
		}`,
		},
		{
			name:               "unknown mailbox",
			endpoint:           "/v1/mailbox/unknown/message",
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:               "offset too high",
			endpoint:           "/v1/mailbox/recovered-lost-folder-d8de382e758459691200000052c9e01c/message?offset=11&limit=5",
			expectedStatusCode: http.StatusOK,
			expectedJsonPayload: `{
			"count": 0,
			"messages": null
		}`,
		},
	}

	tester := tdhttp.NewTestAPI(t, api.Handle())
	for _, test := range tests {
		tt := tester.
			Name(test.name).
			Get(test.endpoint).
			CmpStatus(test.expectedStatusCode)

		if test.expectedJsonPayload != "" {
			tt.CmpJSONBody(td.JSON(test.expectedJsonPayload))
		}
	}
}
