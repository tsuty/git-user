package main

import (
	"errors"
	"fmt"
	"github.com/urfave/cli"
	"os"
	"strings"
	"time"
)

func main() {
	gitUserApp().Run(os.Args)
}

func gitUserApp() *cli.App {
	app := cli.NewApp()
	app.Name = "git-user"
	app.Version = "0.0.1"
	app.Compiled = time.Now()
	app.Usage = "git client multi-user support"
	app.ExitErrHandler = func(_ *cli.Context, err error) {
		if err == nil {
			return
		}
		msg := err.Error()
		if msg != "" {
			fmt.Fprintln(os.Stderr, msg)
		}
		cli.OsExiter(1)
		return
	}
	app.EnableBashCompletion = true
	app.BashComplete = func(c *cli.Context) {
		fmt.Fprintf(c.App.Writer, "show\nset\ndelete\nlocal\nlist\nsync\n")
	}
	app.UsageText = "git-user command [options] [arguments...]"
	app.Commands = []cli.Command{
		{
			Name:      "show",
			Usage:     "Show user configuration",
			UsageText: "git-user show [options]",
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "name, n", Usage: "show user name"},
				cli.BoolFlag{Name: "email, e", Usage: "show user email"},
				cli.BoolFlag{Name: "signingkey, s", Usage: "show user signingkey"},
				cli.BoolFlag{Name: "url, u", Usage: "show url"},
			},
			Action: func(c *cli.Context) error {
				return showAction(c)
			},
		},
		{
			Name:      "prompt",
			Usage:     "For shell prompt",
			UsageText: "git-user prompt [options]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "format, f",
					Usage: "format (name: $n, email: $e, signingkey: $s, url: $u)",
					Value: "[$e]",
				},
				cli.BoolFlag{
					Name:  "quiet, q",
					Usage: "Hide message",
				},
			},
			Action: func(c *cli.Context) error {
				promptAction(c)
				return nil
			},
		},
		{
			Name:      "set",
			Usage:     "Set git user configuration",
			UsageText: "git-user set [options] name email [signingkey]",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "url, u", Usage: "Repository url (default: current repository url)"},
			},
			Action: func(c *cli.Context) error {
				return setAction(c)
			},
		},
		{
			Name:      "delete",
			Usage:     "Delete git user configuration",
			UsageText: "git-user delete id",
			ArgsUsage: "ID",
			Action: func(c *cli.Context) error {
				return deleteAction(c)
			},
		},
		{
			Name:      "local",
			Usage:     "Show local git user.* configuration",
			UsageText: "git-user local [options]",
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "name, n", Usage: "show user name"},
				cli.BoolFlag{Name: "email, e", Usage: "show user email"},
				cli.BoolFlag{Name: "signingkey, s", Usage: "show user signingkey"},
				cli.BoolFlag{Name: "url, u", Usage: "show url"},
			},
			Action: func(c *cli.Context) error {
				return localAction(c)
			},
		},
		{
			Name:      "list",
			Usage:     "Show all git user configuration",
			UsageText: "git-user list",
			Action: func(c *cli.Context) error {
				return listAction(c)
			},
		},
		{
			Name:      "sync",
			Usage:     "Sync to local git configuration",
			UsageText: "git-user sync [options]",
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "quiet, q", Usage: "Hide any message"},
			},
			Action: func(c *cli.Context) error {
				return syncAction(c)
			},
		},
	}
	return app
}

func showAction(c *cli.Context) error {
	if !isInsideWorkTree() {
		return errors.New("outside work tree")
	}

	config := loadConfig()
	url, err := getRemoteOriginUrl()
	if err != nil {
		return err
	}
	u, err := config.findUserByUrl(url)
	if err != nil {
		return err
	}

	o := output{
		Url:        c.Bool("url"),
		Name:       c.Bool("name"),
		Email:      c.Bool("email"),
		SigningKey: c.Bool("signingkey"),
	}

	fmt.Fprintln(c.App.Writer, u.toString(o))

	return nil
}

func promptAction(c *cli.Context) {
	if !isInsideWorkTree() {
		return
	}

	quiet := c.Bool("quiet")
	config := loadConfig()

	url, err := getRemoteOriginUrl()
	if err != nil {
		if !quiet {
			fmt.Fprintf(c.App.Writer, "[git-user:%s]", err.Error())
		}
		return
	}

	u, err := config.findUserByUrl(url)
	if err != nil {
		if !quiet {
			fmt.Fprintf(c.App.Writer, "[git-user:%s]", err.Error())
		}
		return
	}

	syncToLocal(url, config)

	format := c.String("format")
	format = strings.Replace(format, "$u", u.Url, -1)
	format = strings.Replace(format, "$n", u.Name, -1)
	format = strings.Replace(format, "$e", u.Email, -1)
	format = strings.Replace(format, "$s", u.SigningKey, -1)
	fmt.Fprint(c.App.Writer, strings.TrimSpace(format))
}

func localAction(c *cli.Context) error {
	if !isInsideWorkTree() {
		return errors.New("outside work tree")
	}

	url, _ := getRemoteOriginUrl()

	u := user{
		Url:        url,
		Name:       getLocalUserName(),
		Email:      getLocalUserEmail(),
		SigningKey: getLocalUserSigningKey(),
	}

	o := output{
		Url:        c.Bool("url"),
		Name:       c.Bool("name"),
		Email:      c.Bool("email"),
		SigningKey: c.Bool("signingkey"),
	}

	fmt.Fprintln(c.App.Writer, u.toString(o))
	return nil
}

func listAction(c *cli.Context) error {
	conf := loadConfig()
	for _, u := range conf {
		fmt.Fprintln(c.App.Writer, u.toString(output{}))
	}

	return nil
}

func setAction(c *cli.Context) error {
	if len(c.Args()) < 2 {
		cli.ShowCommandHelp(c, "set")
		return errors.New("not enough arguments")
	}
	url := c.String("url")
	if url == "" {
		u, err := getRemoteOriginUrl()
		if err != nil {
			return err
		}
		url = u
	}

	name := ""
	email := ""
	signingkey := ""
	for i, v := range c.Args() {
		switch i {
		case 0:
			name = v
		case 1:
			email = v
		case 2:
			signingkey = v
		}
	}

	config := loadConfig()
	user, err := config.set(url, name, email, signingkey)

	fmt.Fprintln(c.App.Writer, user.toString(output{}))

	return err
}

func deleteAction(c *cli.Context) error {
	if len(c.Args()) != 1 {
		cli.ShowCommandHelp(c, "delete")
		return errors.New("not enough arguments")
	}

	id := c.Args()[0]
	config := loadConfig()
	u, err := config.delete(id)
	if err != nil {
		return errors.New(fmt.Sprintf("no such user. ID: %s", id))
	} else {
		fmt.Fprintln(c.App.Writer, u.toString(output{}))
	}

	return nil
}

func syncAction(c *cli.Context) error {
	quiet := c.Bool("quiet")

	if !isInsideWorkTree() {
		if quiet {
			return errors.New("")
		} else {
			return errors.New("outside work tree")
		}
	}

	url, err1 := getRemoteOriginUrl()
	if err1 != nil {
		if quiet {
			return errors.New("")
		} else {
			return err1
		}
	}

	conf := loadConfig()
	msg, e := syncToLocal(url, conf)

	if e != nil {
		if quiet {
			return errors.New("")
		} else {
			return e
		}
	} else if !quiet {
		fmt.Fprintln(c.App.Writer, "sync succeeded")
		fmt.Fprintln(c.App.Writer, msg)
	}

	return nil
}

func syncToLocal(url string, conf config) (string, error) {
	u, err := conf.findUserByUrl(url)
	if err != nil {
		unsetLocalUserName()
		unsetLocalUserEmail()
		unsetLocalUserSigningKey()
		return "", nil
	}

	if u.Name != "" {
		if getLocalUserName() != u.Name {
			e := setLocalUserName(u.Name)
			if e != nil {
				return "", e
			}
		}
	} else {
		unsetLocalUserName()
	}

	if u.Email != "" {
		if getLocalUserEmail() != u.Email {
			e := setLocalUserEmail(u.Email)
			if e != nil {
				return "", e
			}
		}
	} else {
		unsetLocalUserEmail()
	}

	if u.SigningKey != "" {
		if getLocalUserSigningKey() != u.SigningKey {
			e := setLocalUserSigningKey(u.SigningKey)
			if e != nil {
				return "", e
			}
		}
	} else {
		unsetLocalUserSigningKey()
	}

	return u.toString(output{}), nil
}
