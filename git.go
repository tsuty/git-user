package main

import (
	"os/exec"
	"strings"
)

// Git execution of git command
type Git struct{}

// IsInsideWorkTree `git rev-parse --is-inside-work-tree`
func (*Git) IsInsideWorkTree() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Trim(string(out), "\n") == "true"
}

// GetRemoteOriginURL `git config --get remote.origin.url`
func (*Git) GetRemoteOriginURL() string {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	out, _ := cmd.Output()
	return strings.Trim(string(out), "\n")
}

// GetLocalUserName `git config --local --get user.name`
func (*Git) GetLocalUserName() string {
	cmd := exec.Command("git", "config", "--local", "--get", "user.name")
	out, _ := cmd.Output()
	return strings.Trim(string(out), "\n")
}

// SetLocalUserName `git config --local user.name $name`
func (*Git) SetLocalUserName(name string) error {
	cmd := exec.Command("git", "config", "--local", "user.name", name)
	return cmd.Run()
}

// UnsetLocalUserName git config --local --unset-all user.name
func (*Git) UnsetLocalUserName() error {
	cmd := exec.Command("git", "config", "--local", "--unset-all", "user.name")
	return cmd.Run()
}

// GetLocalUserEmail `git config --local --get user.email`
func (*Git) GetLocalUserEmail() string {
	cmd := exec.Command("git", "config", "--local", "--get", "user.email")
	out, _ := cmd.Output()
	return strings.Trim(string(out), "\n")
}

// SetLocalUserEmail `git config --local user.email $email`
func (*Git) SetLocalUserEmail(email string) error {
	cmd := exec.Command("git", "config", "--local", "user.email", email)
	return cmd.Run()
}

// UnsetLocalUserEmail `git config --local --unset-all user.email`
func (*Git) UnsetLocalUserEmail() error {
	cmd := exec.Command("git", "config", "--local", "--unset-all", "user.email")
	return cmd.Run()
}

// GetLocalUserSigningKey `git config --local --get user.signingkey`
func (*Git) GetLocalUserSigningKey() string {
	cmd := exec.Command("git", "config", "--local", "--get", "user.signingkey")
	out, _ := cmd.Output()
	return strings.Trim(string(out), "\n")
}

// SetLocalUserSigningKey `git config --local user.signingkey $signingkey`
func (*Git) SetLocalUserSigningKey(signingkey string) error {
	cmd := exec.Command("git", "config", "--local", "user.signingkey", signingkey)
	return cmd.Run()
}

// UnsetLocalUserSigningKey `git config --local --unset-all user.signingkey`
func (*Git) UnsetLocalUserSigningKey() error {
	cmd := exec.Command("git", "config", "--local", "--unset-all", "user.signingkey")
	return cmd.Run()
}
