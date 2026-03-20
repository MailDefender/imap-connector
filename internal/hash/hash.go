package hash

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"maildefender/imap-connector/internal/models"
)

const LatestVersion = -1

type hash struct {
	version int
	content string
}

func (h hash) String() string {
	return fmt.Sprintf("%d:%s", h.version, h.content)
}

func (h hash) Match(w hash) bool {
	return h.version == w.version && h.content == w.content
}

func (h hash) Version() int {
	return h.version
}

func Parse(input string) (hash, error) {
	elts := strings.Split(input, ":")

	errFormatError := errors.New("wrong hash format")

	if len(elts) != 2 {
		return hash{}, errFormatError
	}

	v, c := elts[0], elts[1]
	if v == "" || c == "" {
		return hash{}, errFormatError
	}

	h := hash{content: c}
	var err error
	if h.version, err = strconv.Atoi(v); err != nil {
		return hash{}, err
	}

	return h, nil
}

func HashMessage(msg models.Message, version int) (hash, error) {
	switch version {
	case 1, LatestVersion:
		return hashV1(msg)
	}

	return hash{}, errors.New("unknown hash version")
}
