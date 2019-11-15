package main

import (
	"crypto/sha1"
	"fmt"
	"sort"

	"github.com/ryanuber/go-glob"
)

const userHashSie = 7

// User is git user conf
type User struct {
	URL        string
	Name       string
	Email      string
	SigningKey string
}

// Hash is identity of user
func (u *User) Hash() string {
	h := fmt.Sprintf("%x", sha1.Sum([]byte(fmt.Sprintf("%+v", *u))))
	return h[0:userHashSie]
}

// PrintUsers is slice of user
type Users []*User

// Len sort.Interface
func (us Users) Len() int {
	return len(us)
}

// Less sort.Interface
func (us Users) Less(i, j int) bool {
	return us[i].URL > us[j].URL
}

// Swap sort.Interface
func (us Users) Swap(i, j int) {
	us[i], us[j] = us[j], us[i]
}

// TakeByURL find user by repository URL
func (us Users) TakeByURL(url string) *User {
	sort.Sort(us)
	for _, user := range us {
		if glob.Glob(user.URL, url) {
			return user
		}
	}
	return nil
}

// TakeByURL find user by user hash
func (us Users) TakeByHash(hash string) *User {
	for _, user := range us {
		if user.Hash() == hash {
			return user
		}
	}
	return nil
}

// Set append or update
func (us *Users) Set(url, name, email, signingkey string) *User {
	nu := &User{
		URL:        url,
		Name:       name,
		Email:      email,
		SigningKey: signingkey,
	}

	updated := false
	for _, user := range *us {
		if user.URL == url {
			updated = true
			user = nu
		}
	}
	if !updated {
		*us = append(*us, nu)
	}
	return nu
}

// Delete delete user if exists
func (us *Users) Delete(hash string) *User {
	var du *User
	var nus Users
	for _, user := range *us {
		if user.Hash() == hash {
			du = user
		} else {
			nus = append(nus, user)
		}
	}
	*us = nus

	return du
}
