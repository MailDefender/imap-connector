package hash

import (
	"log"
	"testing"
	"time"

	"maildefender/imap-connector/internal/models"
)

func TestHashComparison(t *testing.T) {
	tests := []struct {
		messageA       models.Message
		versionHashA   int
		messageB       models.Message
		versionHashB   int
		shouldMatch    bool
		shouldGetError bool
	}{
		{
			messageA: models.Message{
				From:    models.Contacts{models.Contact{Name: "toto", Email: "emailToto", Host: "example.com"}},
				Subject: "Test mail 1",
				Date:    time.Date(2025, 10, 10, 05, 00, 00, 00, time.Local),
				To:      models.Contacts{models.Contact{Name: "To", Email: "emailTo", Host: "example.com"}},
			},
			versionHashA: 1,
			messageB: models.Message{
				From:    models.Contacts{models.Contact{Name: "toto", Email: "emailToto", Host: "example.com"}},
				Subject: "Test mail 1",
				Date:    time.Date(2025, 10, 10, 05, 00, 00, 00, time.Local),
				To:      models.Contacts{models.Contact{Name: "To", Email: "emailTo", Host: "example.com"}},
			},
			versionHashB:   1,
			shouldMatch:    true,
			shouldGetError: false,
		},
		{
			messageA: models.Message{
				From:    models.Contacts{models.Contact{Name: "toto", Email: "emailToto", Host: "example.com"}},
				Subject: "Test mail 1",
				Date:    time.Date(2025, 10, 10, 05, 00, 00, 00, time.Local),
				To:      models.Contacts{models.Contact{Name: "To", Email: "emailTo", Host: "example.com"}},
			},
			versionHashA: 1,
			messageB: models.Message{
				From:    models.Contacts{models.Contact{Name: "toto", Email: "emailToto", Host: "example.com"}},
				Subject: "Test mail 2",
				Date:    time.Date(2025, 10, 10, 05, 00, 00, 00, time.Local),
				To:      models.Contacts{models.Contact{Name: "To", Email: "emailTo", Host: "example.com"}},
			},
			versionHashB:   1,
			shouldMatch:    false,
			shouldGetError: false,
		},
		{
			messageA: models.Message{
				From:    models.Contacts{models.Contact{Name: "toto", Email: "emailToto", Host: "example.com"}},
				Subject: "Test mail 1",
				Date:    time.Date(2025, 10, 10, 05, 00, 00, 00, time.Local),
				To:      models.Contacts{models.Contact{Name: "To", Email: "emailTo", Host: "example.com"}},
			},
			versionHashA: 1,
			messageB: models.Message{
				From:    models.Contacts{models.Contact{Name: "toto", Email: "emailToto", Host: "example.com"}},
				Subject: "Test mail 1",
				Date:    time.Date(2025, 10, 10, 05, 00, 00, 00, time.Local),
				To:      models.Contacts{models.Contact{Name: "To", Email: "emailTo", Host: "example.com"}},
			},
			versionHashB:   2,
			shouldMatch:    false,
			shouldGetError: true,
		},
		{
			messageA: models.Message{
				From:    models.Contacts{models.Contact{Name: "toto", Email: "emailToto", Host: "example.com"}},
				Subject: "Test mail 1",
				Date:    time.Date(2025, 10, 10, 05, 00, 00, 00, time.Local),
				To: models.Contacts{
					models.Contact{Name: "To", Email: "emailTo", Host: "example.com"},
					models.Contact{Name: "ToT", Email: "emailToT", Host: "example2.com"},
				},
			},
			versionHashA: 1,
			messageB: models.Message{
				From:    models.Contacts{models.Contact{Name: "toto", Email: "emailToto", Host: "example.com"}},
				Subject: "Test mail 1",
				Date:    time.Date(2025, 10, 10, 05, 00, 00, 00, time.Local),
				To: models.Contacts{
					models.Contact{Name: "To", Email: "emailTo", Host: "example.com"},
					models.Contact{Name: "ToT", Email: "emailToT", Host: "example2.com"},
				}},
			versionHashB:   1,
			shouldMatch:    true,
			shouldGetError: false,
		},
		{
			messageA: models.Message{
				From:    models.Contacts{models.Contact{Name: "tot", Email: "emailToto", Host: "example.com"}},
				Subject: "Test mail 1",
				Date:    time.Date(2025, 10, 10, 05, 00, 00, 00, time.Local),
				To: models.Contacts{
					models.Contact{Name: "To", Email: "emailTo", Host: "example.com"},
					models.Contact{Name: "ToT", Email: "emailToT", Host: "example2.com"},
				},
			},
			versionHashA: 1,
			messageB: models.Message{
				From:    models.Contacts{models.Contact{Name: "toto", Email: "emailToto", Host: "example.com"}},
				Subject: "Test mail 1",
				Date:    time.Date(2025, 10, 10, 05, 00, 00, 00, time.Local),
				To: models.Contacts{
					models.Contact{Name: "To", Email: "emailTo", Host: "example.com"},
					models.Contact{Name: "ToT", Email: "emailToT", Host: "example2.com"},
				}},
			versionHashB:   1,
			shouldMatch:    false,
			shouldGetError: false,
		},
	}

	for _, t := range tests {
		hashMsgA, err := HashMessage(t.messageA, t.versionHashA)
		if err != nil && !t.shouldGetError {
			log.Fatal("HashMessage on message A returns an error and no error was expected")
		}

		if err != nil {
			continue
		}

		if hashMsgA.version != t.versionHashA {
			log.Fatal("hash vesion of message A does not match with the wanted one")
		}

		hashMsgB, err := HashMessage(t.messageB, t.versionHashB)
		if err != nil && !t.shouldGetError {
			log.Fatal("HashMessage on message B returns an error and no error was expected")
		}

		if err != nil {
			continue
		}

		if hashMsgB.version != t.versionHashB {
			log.Fatal("hash vesion of message B does not match with the wanted one")
		}

		if t.shouldGetError {
			log.Fatal("an error was expected")
		}

		if hashMsgA.Match(hashMsgB) != t.shouldMatch || hashMsgB.Match(hashMsgA) != t.shouldMatch {
			log.Fatal("Match() result does not match with 'shouldMatch'")
		}
	}
}
