package main

import (
	"fmt"
	"io"
	"math/bits"
	"strings"
)

// print item
const (
	PrintURL uint = 1 << iota
	PrintName
	PrintEmail
	PrintSigningKey
)

// print item
const (
	PrintDefault uint = 0
	PrintALL          = PrintURL | PrintName | PrintEmail | PrintSigningKey
)

// Printer print formatter
type Printer struct {
	flag   uint
	writer io.Writer
}

// NewPrinter inti Printer
func NewPrinter(flag uint, writer io.Writer) *Printer {
	if flag == PrintDefault {
		flag = PrintALL
	}
	return &Printer{
		flag:   flag,
		writer: writer,
	}
}

func (p Printer) buf(user *User) []string {
	var buf []string
	if p.flag&PrintURL == PrintURL {
		buf = append(buf, fmt.Sprintf("URL: %s", user.URL))
	}
	if p.flag&PrintName == PrintName {
		buf = append(buf, fmt.Sprintf("Name: %s", user.Name))
	}
	if p.flag&PrintEmail == PrintEmail {
		buf = append(buf, fmt.Sprintf("Email: %s", user.Email))
	}
	if p.flag&PrintSigningKey == PrintSigningKey {
		buf = append(buf, strings.TrimSpace(fmt.Sprintf("SigningKey: %s", user.SigningKey)))
	}
	if p.flag == PrintALL {
		buf = append(buf, fmt.Sprintf("Hash: %s", user.Hash()))
	}
	return buf
}

// PrintUser print user
func (p Printer) PrintUser(user *User) Printer {
	fmt.Fprintln(p.writer, strings.Join(p.buf(user), "  "))
	return p
}

// PrintUsers print users
func (p Printer) PrintUsers(users []*User) Printer {
	lines := make([][]string, len(users))
	colMaxLen := make([]int, bits.OnesCount(p.flag)+1)
	for i, u := range users {
		lines[i] = p.buf(u)
		for j, s := range lines[i] {
			l := len(s)
			if colMaxLen[j] < l {
				colMaxLen[j] = l
			}
		}
	}

	for _, line := range lines {
		buf := make([]string, len(line))
		for i, col := range line {
			buf = append(
				buf,
				fmt.Sprintf("%-0"+fmt.Sprint(colMaxLen[i]+2)+"s", col),
			)
		}
		fmt.Fprintln(p.writer, strings.TrimSpace(strings.Join(buf, "")))
	}

	return p
}

// Println print message with line feed
func (p Printer) Println(message string) Printer {
	fmt.Fprintln(p.writer, message)
	return p
}

// Printf print formatted message
func (p Printer) Printf(format string, a ...interface{}) Printer {
	fmt.Fprintf(p.writer, format, a...)
	return p
}
