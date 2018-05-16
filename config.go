package main

import (
	"encoding/json"
	"errors"
	"github.com/ryanuber/go-glob"
	"io/ioutil"
	"os"
)

const (
	defaultConfigFile = ".git-user.json"
)

type config []user

func (c config) findUserByUrl(url string) (user, error) {
	var u user
	for _, u := range c {
		if glob.Glob(u.Url, url) {
			return u, nil
		}
	}
	return u, errors.New("git user not found")
}

func (c config) findUserById(id string) (user, error) {
	var u user
	for _, u := range c {
		if u.Id == id {
			return u, nil
		}
	}
	return u, errors.New("git user not found")
}

func (c config) save() error {
	bytes, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(c.path(), bytes, 0644)
}

func (c config) set(url string, name string, email string, signingkey string) (user, error) {
	u := user{
		Url:        url,
		Name:       name,
		Email:      email,
		SigningKey: signingkey,
	}
	u.init()

	nc := config{}
	for _, t := range c {
		if glob.Glob(t.Url, url) {
			return t, errors.New("already exist or match wildcard")
		}

		if !glob.Glob(url, t.Url) {
			nc = append(nc, t)
		}
	}
	nc = append(nc, u)
	return u, nc.save()
}

func (c config) delete(id string) (user, error) {
	u, e := c.findUserById(id)
	if e != nil {
		return u, e
	}

	nc := config{}
	for _, t := range c {
		if t.Id != id {
			nc = append(nc, t)
		}
	}
	nc.save()

	return u, nil
}

func (_ config) path() string {
	path := os.Getenv("GIT_USER_CONFIG")
	if path == "" {
		return os.Getenv("HOME") + string(os.PathSeparator) + defaultConfigFile
	}
	return path
}

func loadConfig() config {
	var c config
	path := c.path()
	if _, statErr := os.Stat(path); os.IsNotExist(statErr) {
		return c
	} else {
		bytes, readErr := ioutil.ReadFile(path)
		if readErr == nil {
			json.Unmarshal(bytes, &c)
		}
		return c
	}
}
