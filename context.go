package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"

	"github.com/mitchellh/go-homedir"
)

type Context struct {
	Option Option
	Users  Users
}

type nullIO struct{}

// Write io.Writer
func (nullIO) Write(p []byte) (n int, err error) {
	return 0, nil
}

// NewContext init context
func NewContext() *Context {
	return &Context{
		Option: Option{
			Show: ShowOption{},
			Set: SetOption{
				Args: SetArgs{},
			},
			Delete: DeleteOption{
				Args: DeleteArgs{},
			},
			Local: LocalOption{},
			List:  ListOption{},
			Sync:  SyncOption{},
			Print: PrintOption{},
		},
	}
}

// LoadConfig load config from json
func (c *Context) LoadConfig() error {
	path, err := c.configPath()
	if err != nil {
		return err
	}
	if _, statErr := os.Stat(path); os.IsNotExist(statErr) {
		return nil
	} else {
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		return json.Unmarshal(bytes, &c.Users)
	}
}

// SaveConfig save config to json
func (c *Context) SaveConfig() error {
	path, err := c.configPath()
	if err != nil {
		return err
	}
	bytes, err := json.Marshal(&c.Users)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, bytes, 0644)
}

// Execute execute action
func (c *Context) Execute(command string) error {
	if err := c.LoadConfig(); err != nil {
		return err
	}

	switch command {
	case "show":
		a := &Action{
			printer: NewPrinter(c.Option.Show.printFlag(), os.Stdout),
		}
		return a.ShowUser(c)
	case "set":
		a := &Action{
			printer: NewPrinter(PrintDefault, os.Stdout),
		}
		return a.SetUser(c)
	case "delete":
		a := &Action{
			printer: NewPrinter(PrintDefault, os.Stdout),
		}
		return a.DeleteUser(c)
	case "local":
		a := &Action{
			printer: NewPrinter(c.Option.Local.printFlag(), os.Stdout),
		}
		return a.ShowLocalUser(c)
	case "list":
		a := &Action{
			printer: NewPrinter(c.Option.List.printFlag(), os.Stdout),
		}
		return a.ListUsers(c)
	case "sync":
		var writer io.Writer
		if c.Option.Sync.Quiet {
			writer = &nullIO{}
		} else {
			writer = os.Stdout
		}
		a := &Action{
			printer: NewPrinter(PrintDefault, writer),
		}
		return a.SyncGitUserToLocal(c)
	case "print":
		a := &Action{
			printer: NewPrinter(PrintDefault, os.Stdout),
		}
		return a.Print(c)
	}

	return nil
}

func (c Context) configPath() (string, error) {
	return homedir.Expand(c.Option.Config)
}
