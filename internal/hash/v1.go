package hash

import (
	"errors"
	"hash/fnv"
	"strconv"

	"maildefender/imap-connector/internal/models"
)

type hashField []byte

func hashV1(msg models.Message) (hash, error) {
	hasher := fnv.New64a()

	if len(msg.From) == 0 {
		return hash{}, errors.New("this message has no sender, cannot compute hash")
	}
	from := msg.From[0]

	fieldsToHash := []hashField{
		hashField(from.Email),
		hashField(from.Name),
		hashField(msg.Subject),
		hashField(msg.Date.String()),
	}

	for _, c := range msg.To {
		fieldsToHash = append(fieldsToHash, hashField(c.Email))
		fieldsToHash = append(fieldsToHash, hashField(c.Name))
	}

	for _, f := range fieldsToHash {
		if _, err := hasher.Write(f); err != nil {
			return hash{}, err
		}
	}

	return hash{
		version: 1,
		content: strconv.FormatUint(hasher.Sum64(), 10),
	}, nil
}
