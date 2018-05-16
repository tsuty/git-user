package main

import (
	"testing"
	"io/ioutil"
	"os"
	"os/exec"
)

func TestOutsideWorkTree(t *testing.T) {
	outsideWorkTree()

	if isInsideWorkTree() {
		t.Errorf("expected: false, actual: true")
	}
}

func TestInsideWorkTree(t *testing.T) {
	insideWorkTree()

	if !isInsideWorkTree() {
		t.Errorf("expected: true, actual: false")
	}
}

func TestRemoteOriginUrlOutsideWorkTree(t *testing.T) {
	outsideWorkTree()

	_, er := getRemoteOriginUrl()
	if er == nil {
		t.Errorf("expected: error, actual: nil")
	}
}

func TestRemoteOriginUrlInsideWorkTreeNonSetUrl(t *testing.T) {
	insideWorkTree()

	_, er := getRemoteOriginUrl()
	if er == nil {
		t.Errorf("expected: error, actual: nil")
	}
}

func TestRemoteOriginUrlInsideWorkTreeSetUrl(t *testing.T) {
	insideWorkTree()

	expected := "git@example.com:foo/bar"
	exec.Command("git", "remote", "add", "origin", expected).Run()

	actual, er := getRemoteOriginUrl()
	if er != nil {
		t.Errorf("expected: nil, actual: error")
		return
	}

	if actual != expected {
		t.Errorf("expected: %+v, actual: %+v", expected, actual)
	}
}

func TestLocalUserName(t *testing.T) {
	insideWorkTree()

	expected := "someone"

	name1 := getLocalUserName()
	if name1 != "" {
		t.Errorf("expected: \"\", actual: %+v", name1)
		return
	}

	setLocalUserName(expected)
	name2 := getLocalUserName()
	if name2 != expected {
		t.Errorf("expected: %+v, actual: %+v", expected, name2)
		return
	}

	unsetLocalUserName()
	name3 := getLocalUserName()
	if name3 != "" {
		t.Errorf("expected: \"\", actual: %+v", name3)
	}
}

func TestLocalUserEmail(t *testing.T) {
	insideWorkTree()

	expected := "someone@example.com"

	mail1 := getLocalUserEmail()
	if mail1 != "" {
		t.Errorf("expected: \"\", actual: %+v", mail1)
		return
	}

	setLocalUserEmail(expected)
	mail2 := getLocalUserEmail()
	if mail2 != expected {
		t.Errorf("expected: %+v, actual: %+v", expected, mail2)
		return
	}

	unsetLocalUserEmail()
	mail3 := getLocalUserEmail()
	if mail3 != "" {
		t.Errorf("expected: \"\", actual: %+v", mail3)
	}
}

func TestLocalUserSigningKey(t *testing.T) {
	insideWorkTree()

	expected := "AAABBBCCCDDD"

	key1 := getLocalUserSigningKey()
	if key1 != "" {
		t.Errorf("expected: \"\", actual: %+v", key1)
		return
	}

	setLocalUserSigningKey(expected)
	key2 := getLocalUserSigningKey()
	if key2 != expected {
		t.Errorf("expected: %+v, actual: %+v", expected, key2)
		return
	}

	unsetLocalUserSigningKey()
	key3 := getLocalUserSigningKey()
	if key3 != "" {
		t.Errorf("expected: \"\", actual: %+v", key3)
	}
}

func insideWorkTree() error {
	dir, _ := ioutil.TempDir(os.TempDir(), "git-user_command_test")
	os.Chdir(dir)
	return exec.Command("git", "init").Run()
}

func outsideWorkTree() {
	dir, _ := ioutil.TempDir(os.TempDir(), "git-user_command_test")
	os.Chdir(dir)
}
