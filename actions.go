package main

import (
	"os"

	"github.com/valyala/fasttemplate"
)

// Action actions
type Action struct {
	printer *Printer
}

// ShowUser
func (a *Action) ShowUser(c *Context) error {
	git := &Git{}
	if !git.IsInsideWorkTree() {
		current, err := os.Getwd()
		a.printer.Printf("outside work tree. %s %v\n", current, err)
		return nil
	}

	url := git.GetRemoteOriginURL()
	if url == "" {
		a.printer.Println("no remote origin url. set your remote origin url!")
		return nil
	}

	user := c.Users.TakeByURL(url)
	if user == nil {
		a.printer.Println("no git-user config. `git-user set name email`")
		return nil
	}

	a.printer.PrintUser(user)

	return nil
}

func (a *Action) SetUser(c *Context) error {
	option := c.Option.Set
	if err := option.Args.Valid(); err != nil {
		return err
	}

	url := option.URL
	if url == "" {
		git := &Git{}
		if !git.IsInsideWorkTree() {
			current, err := os.Getwd()
			a.printer.Printf("outside work tree. %s %v\n", current, err)
			a.printer.Println("required `--url` option or inside work tree")
			return nil
		}

		url = git.GetRemoteOriginURL()
		if url == "" {
			a.printer.Println("no remote origin url")
			a.printer.Println("required `--url` option or set your remote origin url!")
			return nil
		}
	}

	user := c.Users.Set(
		url,
		option.Args.Name,
		option.Args.Email,
		option.Args.SigningKey,
	)

	if err := c.SaveConfig(); err != nil {
		return err
	}

	a.printer.PrintUser(user)

	return nil
}

func (a *Action) DeleteUser(c *Context) error {
	hash := c.Option.Delete.Args.Hash

	user := c.Users.Delete(hash)
	if user == nil {
		a.printer.Printf("not found user by %s\n", hash)
		return nil
	}

	a.printer.PrintUser(user)

	return nil
}

func (a *Action) ShowLocalUser(c *Context) error {
	git := &Git{}
	if !git.IsInsideWorkTree() {
		current, err := os.Getwd()
		a.printer.Printf("outside work tree. %s %v\n", current, err)
		return nil
	}

	user := &User{
		URL:        git.GetRemoteOriginURL(),
		Name:       git.GetLocalUserName(),
		Email:      git.GetLocalUserEmail(),
		SigningKey: git.GetLocalUserSigningKey(),
	}

	a.printer.PrintUser(user)
	return nil
}

func (a *Action) ListUsers(c *Context) error {
	a.printer.PrintUsers(c.Users)
	return nil
}

func (a *Action) SyncGitUserToLocal(c *Context) error {
	git := &Git{}
	if !git.IsInsideWorkTree() {
		current, err := os.Getwd()
		a.printer.Printf("outside work tree. %s %v\n", current, err)
		return nil
	}

	url := git.GetRemoteOriginURL()
	if url == "" {
		a.printer.Println("no remote origin url. set your remote origin url!")
		return nil
	}

	user := c.Users.TakeByURL(url)

	if user != nil && user.Name != "" {
		n := git.GetLocalUserName()
		if user.Name != n {
			if err := git.SetLocalUserName(user.Name); err != nil {
				a.printer.Println(err.Error())
			}
		}
	} else {
		if err := git.UnsetLocalUserName(); err != nil {
			a.printer.Println(err.Error())
		}
	}

	if user != nil && user.Email != "" {
		e := git.GetLocalUserEmail()
		if user.Email != e {
			if err := git.SetLocalUserEmail(user.Email); err != nil {
				a.printer.Println(err.Error())
			}
		}
	} else {
		if err := git.UnsetLocalUserEmail(); err != nil {
			a.printer.Println(err.Error())
		}
	}

	if user != nil && user.SigningKey != "" {
		s := git.GetLocalUserSigningKey()
		if user.SigningKey != s {
			if err := git.SetLocalUserSigningKey(user.SigningKey); err != nil {
				a.printer.Println(err.Error())
			}
		}
	} else {
		if err := git.UnsetLocalUserSigningKey(); err != nil {
			a.printer.Println(err.Error())
		}
	}

	return nil
}

func (a *Action) Print(c *Context) error {
	syncAction := &Action{
		printer: NewPrinter(PrintDefault, &nullIO{}),
	}

	if err := syncAction.SyncGitUserToLocal(c); err != nil {
		return nil
	}

	git := &Git{}
	user := c.Users.TakeByURL(git.GetRemoteOriginURL())

	temp := fasttemplate.New(c.Option.Print.Format, "{", "}")
	temp.Execute(
		a.printer.writer,
		map[string]interface{}{
			"n": user.Name,
			"e": user.Email,
			"u": user.URL,
			"s": user.SigningKey,
		},
	)
	return nil
}
