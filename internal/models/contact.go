package models

import (
	"fmt"
	"strings"
)

type Contact struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Host  string `json:"host"`
}

type Contacts []Contact

func (c *Contact) String() string {
	if c.Name == "" {
		return c.Email
	}

	return fmt.Sprintf("%s <%s>", c.Name, c.Email)
}

func (cs *Contacts) String() string {
	var items []string
	for _, c := range *cs {
		items = append(items, c.String())
	}

	return strings.Join(items, ",")
}
