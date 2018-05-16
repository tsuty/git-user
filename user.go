package main

import (
	"crypto/sha1"
	"fmt"
)

type output struct {
	Url        bool
	Name       bool
	Email      bool
	SigningKey bool
}

type user struct {
	Id         string `json:"id"`
	Url        string `json:"url"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	SigningKey string `json:"signingkey"`
}

func (u *user) init() {
	bytes := fmt.Sprintf("%s$%s$%s$%s", u.Url, u.Name, u.Email, u.SigningKey)
	data := []byte(bytes)
	u.Id = fmt.Sprintf("%x", sha1.Sum(data))
}

func (u user) toString(o output) string {
	if o.Url {
		return u.Url
	}
	if o.Name {
		return u.Name
	}
	if o.Email {
		return u.Email
	}
	if o.SigningKey {
		return u.SigningKey
	}

	if u.Id != "" {
		return fmt.Sprintf(
			"Url: %s\tName: %s\tEmail: %s\tSigningKey: %s\tID: %s",
			u.Url,
			u.Name,
			u.Email,
			u.SigningKey,
			u.Id,
		)
	} else {
		return fmt.Sprintf(
			"Url: %s\tName: %s\tEmail: %s\tSigningKey: %s",
			u.Url,
			u.Name,
			u.Email,
			u.SigningKey,
		)
	}
}
