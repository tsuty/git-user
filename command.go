package main

import (
	"errors"
	"os/exec"
	"strings"
)

func isInsideWorkTree() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	out, _ := cmd.Output()
	return strings.Trim(string(out), "\n") == "true"
}

func getRemoteOriginUrl() (string, error) {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	out, _ := cmd.Output()
	url := strings.Trim(string(out), "\n")
	if url == "" {
		return "", errors.New("remote origin url is not Set")
	} else {
		return url, nil
	}
}

func getLocalUserName() string {
	cmd := exec.Command("git", "config", "--local", "--get", "user.name")
	out, _ := cmd.Output()
	return strings.Trim(string(out), "\n")
}

func setLocalUserName(name string) error {
	cmd := exec.Command("git", "config", "--local", "user.name", name)
	return cmd.Run()
}

func unsetLocalUserName() error {
	cmd := exec.Command("git", "config", "--local", "--unset-all", "user.name")
	return cmd.Run()
}

func getLocalUserEmail() string {
	cmd := exec.Command("git", "config", "--local", "--get", "user.email")
	out, _ := cmd.Output()
	return strings.Trim(string(out), "\n")
}

func setLocalUserEmail(email string) error {
	cmd := exec.Command("git", "config", "--local", "user.email", email)
	return cmd.Run()
}

func unsetLocalUserEmail() error {
	cmd := exec.Command("git", "config", "--local", "--unset-all", "user.email")
	return cmd.Run()
}

func getLocalUserSigningKey() string {
	cmd := exec.Command("git", "config", "--local", "--get", "user.signingkey")
	out, _ := cmd.Output()
	return strings.Trim(string(out), "\n")
}

func setLocalUserSigningKey(signingkey string) error {
	cmd := exec.Command("git", "config", "--local", "user.signingkey", signingkey)
	return cmd.Run()
}

func unsetLocalUserSigningKey() error {
	cmd := exec.Command("git", "config", "--local", "--unset-all", "user.signingkey")
	return cmd.Run()
}
