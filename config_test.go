package main

import (
	"testing"
	"io/ioutil"
	"os"
)

func TestLoadingConfigWhenNotExist(t *testing.T) {
	configNotExist()

	config := loadConfig()
	if len(config) != 0 {
		t.Errorf("expected: no data, actual: %d data", len(config))
	}
}

func TestLoadingConfigWhenExist(t *testing.T) {
	configExist()

	config := loadConfig()
	if len(config) == 0 {
		t.Errorf("expected: data exist, actual: no data")
	}
}

func TestFindingByUrlWhenExist(t *testing.T) {
	configExist()

	url := "git@example.com:foo/bar"

	config := loadConfig()
	_, err := config.findUserByUrl(url)
	if err != nil {
		t.Errorf("expected: nil, actual: error")
	}
}

func TestFindingByUrlWhenNotExist(t *testing.T) {
	configExist()

	url := "git@example.org:hoge/fuga"

	config := loadConfig()
	_, err := config.findUserByUrl(url)
	if err == nil {
		t.Errorf("expected: error, actual: nil")
	}
}

func TestFindingByIdWhenExist(t *testing.T) {
	configExist()

	id := "a1cd4ef7fa69afc33cff820e70a447e4f83ff9a1"

	config := loadConfig()
	_, err := config.findUserById(id)
	if err != nil {
		t.Errorf("expected: nil, actual: error")
	}
}

func TestFindingByIdWhenNotExist(t *testing.T) {
	configExist()

	id := "some_id"

	config := loadConfig()
	_, err := config.findUserById(id)
	if err == nil {
		t.Errorf("expected: error, actual: nil")
	}
}

func TestSaveConfig(t *testing.T) {
	configNotExist()

	config := loadConfig()
	config.save()

	if _, statErr := os.Stat(config.path()); os.IsNotExist(statErr) {
		t.Errorf("expected: success, actual: failure")
	}
}

func TestSetUser(t *testing.T) {
	configNotExist()

	config1 := loadConfig()
	config1.set(
		"git@other.example.com:foo/bar",
		"Foo",
		"foo@example.com",
		"",
	)

	config2 := loadConfig()
	if len(config2) < 1 {
		t.Errorf("expected: data increase, actual: data no change")
	}
}

func TestSetUserWhenSameUrl(t *testing.T) {
	configNotExist()
	url := "git@example.com:foo/bar"

	config1 := loadConfig()
	user1, _ := config1.set(
		url,
		"Foo",
		"foo@example.com",
		"",
	)

	config2 := loadConfig()
	user2, _ := config2.set(
		url,
		"Bar",
		"bar@example.com",
		"",
	)

	config3 := loadConfig()

	if len(config3) != 1 {
		t.Errorf("expected: data no change, actual: data change")
		return
	}

	if user1.Id != user2.Id {
		t.Errorf("expected: data no change, actual: data change")
	}
}

func TestSetUserWhenWildcardMatchedUrl(t *testing.T) {
	configNotExist()

	wildcardUrl := "git@example.com:*"
	matchedUrl := "git@example.com:foo/bar"

	config1 := loadConfig()
	config1.set(
		wildcardUrl,
		"Foo",
		"foo@example.com",
		"",
	)

	config2 := loadConfig()
	user, _ := config2.set(
		matchedUrl,
		"Bar",
		"bar@example.com",
		"",
	)

	config3 := loadConfig()

	if len(config3) != 1 {
		t.Errorf("expected: data no change, actual: data change")
		return
	}

	if user.Url != wildcardUrl {
		t.Errorf("expected: data no change, actual: data change")
	}
}

func TestSetUserWhenWildcardUrl(t *testing.T) {
	configNotExist()

	baseUrl := "git@example.com:foo/bar"
	wildcardUrl := "git@example.com:*"

	config1 := loadConfig()
	config1.set(
		baseUrl,
		"Foo",
		"foo@example.com",
		"",
	)

	config2 := loadConfig()
	config2.set(
		wildcardUrl,
		"Bar",
		"bar@example.com",
		"",
	)

	config3 := loadConfig()

	if len(config3) != 1 {
		t.Errorf("expected: data no change, actual: data change")
		return
	}

	user := config3[0]
	if user.Url != wildcardUrl {
		t.Errorf("expected: data override, actual: data no override")
	}
}

func TestDeleteWhenNotExist(t *testing.T) {
	configExist()

	id := "dummy_id"

	config := loadConfig()
	_, err := config.delete(id)
	if err == nil {
		t.Errorf("expected: error, actual: nil")
	}
}

func TestDeleteWhenExist(t *testing.T) {
	configExist()

	id := "a1cd4ef7fa69afc33cff820e70a447e4f83ff9a1"

	config1 := loadConfig()
	_, err := config1.delete(id)
	if err != nil {
		t.Errorf("expected: nil, actual: error")
		return
	}

	config2 := loadConfig()
	if len(config2) != 0 {
		t.Errorf("expected: data delete, actual: data no delete")
	}
}

func configExist() error {
	dir, _ := ioutil.TempDir(os.TempDir(), "git-user_config_test")
	path := dir + "/" + defaultConfigFile
	os.Setenv("GIT_USER_CONFIG", path)

	bytes := []byte(`
[
  {
    "id":"a1cd4ef7fa69afc33cff820e70a447e4f83ff9a1",
    "url":"git@example.com:*",
    "name":"someone",
    "email":"someone@example.com",
    "signingkey":"AAABBBCCC"
  }
]
`)
	return ioutil.WriteFile(path, bytes, 0644)
}

func configNotExist() {
	dir, _ := ioutil.TempDir(os.TempDir(), "git-user_config_test")
	path := dir + "/" + defaultConfigFile
	os.Setenv("GIT_USER_CONFIG", path)
}
