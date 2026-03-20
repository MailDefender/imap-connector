package configuration

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerConfiguration(t *testing.T) {
	tests := []struct {
		input         MailServerConfiguration
		expectedError error
	}{
		{
			input: MailServerConfiguration{
				Server: ServerConfiguration{
					Host: "ssl.example.com",
					Port: 52,
				},
				Authentication: AuthenticationConfiguration{
					Username: "foo@example.com",
					Password: "P4ssw0rd",
				},
			},
			expectedError: nil,
		},
		{
			input: MailServerConfiguration{
				Server: ServerConfiguration{
					Port: 52,
				},
				Authentication: AuthenticationConfiguration{
					Username: "foo@example.com",
					Password: "P4ssw0rd",
				},
			},
			expectedError: errors.New("host not set"),
		},
		{
			input: MailServerConfiguration{
				Server: ServerConfiguration{
					Port: 52,
				},
				Authentication: AuthenticationConfiguration{
					Username: "foo@example.com",
				},
			},
			expectedError: errors.New("host not set"),
		},
	}

	for _, test := range tests {
		res := test.input.Check()
		assert.Equal(t, test.expectedError, res)
	}
}
