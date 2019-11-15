package main

import "errors"

// Option command option
type Option struct {
	Show   ShowOption   `command:"show" description:"Show git-user"`
	Set    SetOption    `command:"set" description:"Set git-user"`
	Delete DeleteOption `command:"delete" description:"Delete git-user"`
	Local  LocalOption  `command:"local" description:"Show local git user.*"`
	List   ListOption   `command:"list" description:"Show all git-user"`
	Sync   SyncOption   `command:"sync" description:"Sync to local git"`
	Print  PrintOption  `command:"print" description:"Print with sync"`

	Config string `long:"config" value-name:"file" description:"configuration file name" default:"~/git-user.json" env:"GIT_USER_CONFIG"`
}

// ShowOption show command option
type ShowOption struct {
	printOption
}

// SetOption set command option
type SetOption struct {
	URL  string  `long:"url" value-name:"url" short:"u" description:"Repository url (default: current repository url)"`
	Args SetArgs `positional-args:"yes"`
}

// SetArgs set command args
type SetArgs struct {
	Name       string `positional-arg-name:"name" description:"name (required)"`
	Email      string `positional-arg-name:"email" description:"email address (required)"`
	SigningKey string `positional-arg-name:"signingkey" description:"signing key (optional)"`
}

// Valid validate set command args
func (as SetArgs) Valid() error {
	if as.Name == "" {
		return errors.New("required name argument")
	}
	if as.Email == "" {
		return errors.New("required email argument")
	}
	return nil
}

// DeleteOption delete command option
type DeleteOption struct {
	Args DeleteArgs `positional-args:"yes" required:"yes"`
}

// DeleteArgs delete command args
type DeleteArgs struct {
	Hash string `positional-arg-name:"hash"`
}

// LocalOption local command option
type LocalOption struct {
	printOption
}

// ListOption list command option
type ListOption struct {
	printOption
}

// SyncOption sync command option
type SyncOption struct {
	Quiet bool `long:"quiet" short:"q" description:"Hide any message"`
}

// PrintOption print command option
type PrintOption struct {
	Format string `long:"format" short:"f" description:"Print format. name:{n}, email:{e}, signingkey:{s}, url:{u}" default:"[{e}]" env:"GIT_USER_PROMPT"`
}

type printOption struct {
	URL        bool `long:"url" short:"u" description:"show url"`
	Name       bool `long:"name" short:"n" description:"show user name"`
	Email      bool `long:"email" short:"e" description:"show user email"`
	SigningKey bool `long:"signingkey" short:"s" description:"show user signingkey"`
}

// printFlag return view flag bit
func (o printOption) printFlag() uint {
	flag := PrintDefault
	if o.URL {
		flag = flag | PrintURL
	}
	if o.Name {
		flag = flag | PrintName
	}
	if o.Email {
		flag = flag | PrintEmail
	}
	if o.SigningKey {
		flag = flag | PrintSigningKey
	}
	return flag
}
