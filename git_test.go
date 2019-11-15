package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

func insideWorkTree() (cleanupFunc func()) {
	dir, _ := ioutil.TempDir(os.TempDir(), "git-user_test")
	if _, found := os.LookupEnv("DEBUG"); found {
		fmt.Println("# insideWorkTree", dir)
	}
	os.Chdir(dir)
	exec.Command("git", "init").Run()
	return func() {
		os.RemoveAll(dir)
	}
}

func outsideWorkTree() (cleanupFunc func()) {
	dir, _ := ioutil.TempDir(os.TempDir(), "git-user_test")
	if _, found := os.LookupEnv("DEBUG"); found {
		fmt.Println("# outsideWorkTree", dir)
	}
	os.Chdir(dir)
	return func() {
		os.RemoveAll(dir)
	}
}

func TestGit_GetLocalUserEmail(t *testing.T) {
	tests := []struct {
		name string
		want string
		init func() func()
	}{
		{
			"valid",
			"tsuty@example.com",
			func() func() {
				fn := insideWorkTree()
				exec.Command("git", "config", "--local", "user.email", "tsuty@example.com").Run()
				return fn
			},
		},
		{
			"invalid",
			"",
			func() func() {
				return insideWorkTree()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gi := &Git{}
			fn := tt.init()
			defer fn()
			if got := gi.GetLocalUserEmail(); got != tt.want {
				t.Errorf("GetLocalUserEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGit_GetLocalUserName(t *testing.T) {
	tests := []struct {
		name string
		want string
		init func() func()
	}{
		{
			"valid",
			"Mike Wazowski",
			func() func() {
				fn := insideWorkTree()
				exec.Command("git", "config", "--local", "user.name", "Mike Wazowski").Run()
				return fn
			},
		},
		{
			"invalid",
			"",
			func() func() {
				return insideWorkTree()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gi := &Git{}
			fn := tt.init()
			defer fn()
			if got := gi.GetLocalUserName(); got != tt.want {
				t.Errorf("GetLocalUserName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGit_GetLocalUserSigningKey(t *testing.T) {
	tests := []struct {
		name string
		want string
		init func() func()
	}{
		{
			"valid",
			"AAABBBCCCDDD",
			func() func() {
				fn := insideWorkTree()
				exec.Command("git", "config", "--local", "user.signingkey", "AAABBBCCCDDD").Run()
				return fn
			},
		},
		{
			"invalid",
			"",
			func() func() {
				return insideWorkTree()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gi := &Git{}
			fn := tt.init()
			defer fn()
			if got := gi.GetLocalUserSigningKey(); got != tt.want {
				t.Errorf("GetLocalUserSigningKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGit_GetRemoteOriginURL(t *testing.T) {
	tests := []struct {
		name string
		want string
		init func() func()
	}{
		{
			"valid",
			"git@example.com:foo/bar",
			func() func() {
				fn := insideWorkTree()
				exec.Command("git", "remote", "add", "origin", "git@example.com:foo/bar").Run()
				return fn
			},
		},
		{
			"invalid",
			"",
			func() func() {
				return insideWorkTree()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gi := &Git{}
			fn := tt.init()
			defer fn()
			if got := gi.GetRemoteOriginURL(); got != tt.want {
				t.Errorf("GetRemoteOriginURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGit_IsInsideWorkTree(t *testing.T) {
	tests := []struct {
		name string
		want bool
		init func() func()
	}{
		{
			"valid",
			true,
			func() func() {
				return insideWorkTree()
			},
		},
		{
			"invalid",
			false,
			func() func() {
				return outsideWorkTree()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gi := &Git{}
			fn := tt.init()
			defer fn()
			if got := gi.IsInsideWorkTree(); got != tt.want {
				t.Errorf("IsInsideWorkTree() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGit_SetLocalUserEmail(t *testing.T) {
	type args struct {
		email string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    string
		init    func() func()
	}{
		{
			"valid",
			args{email: "tsuty@example.com"},
			false,
			"tsuty@example.com",
			func() func() {
				return insideWorkTree()
			},
		},
		{
			"invalid",
			args{email: "tsuty@example.com"},
			true,
			"",
			func() func() {
				return outsideWorkTree()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gi := &Git{}
			fn := tt.init()
			defer fn()
			err := gi.SetLocalUserEmail(tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetLocalUserEmail() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				if got := gi.GetLocalUserEmail(); got != tt.want {
					t.Errorf("GetLocalUserEmail() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestGit_SetLocalUserName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    string
		init    func() func()
	}{
		{
			"valid",
			args{name: "Mike Wazowski"},
			false,
			"Mike Wazowski",
			func() func() {
				return insideWorkTree()
			},
		},
		{
			"invalid",
			args{name: "Mike Wazowski"},
			true,
			"",
			func() func() {
				return outsideWorkTree()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gi := &Git{}
			fn := tt.init()
			defer fn()
			err := gi.SetLocalUserName(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetLocalUserName() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				if got := gi.GetLocalUserName(); got != tt.want {
					t.Errorf("GetLocalUserName() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestGit_SetLocalUserSigningKey(t *testing.T) {
	type args struct {
		signingkey string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    string
		init    func() func()
	}{
		{
			"valid",
			args{signingkey: "aaabbbcccddeeefff"},
			false,
			"aaabbbcccddeeefff",
			func() func() {
				return insideWorkTree()
			},
		},
		{
			"invalid",
			args{signingkey: "aaabbbcccddeeefff"},
			true,
			"",
			func() func() {
				return outsideWorkTree()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gi := &Git{}
			fn := tt.init()
			defer fn()

			err := gi.SetLocalUserSigningKey(tt.args.signingkey)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetLocalUserSigningKey() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				if got := gi.GetLocalUserSigningKey(); got != tt.want {
					t.Errorf("GetLocalUserSigningKey() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestGit_UnsetLocalUserEmail(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
		want    string
		init    func() func()
	}{
		{
			"valid",
			false,
			"",
			func() func() {
				fn := insideWorkTree()
				exec.Command("git", "config", "--local", "user.email", "tsuty@example.com").Run()
				return fn
			},
		},
		{
			"invalid",
			true,
			"",
			func() func() {
				return insideWorkTree()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gi := &Git{}
			fn := tt.init()
			defer fn()
			if err := gi.UnsetLocalUserEmail(); (err != nil) != tt.wantErr {
				t.Errorf("UnsetLocalUserEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got := gi.GetLocalUserEmail(); got != tt.want {
				t.Errorf("GetLocalUserEmail() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGit_UnsetLocalUserName(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
		want    string
		init    func() func()
	}{
		{
			"valid",
			false,
			"",
			func() func() {
				fn := insideWorkTree()
				exec.Command("git", "config", "--local", "user.name", "Mike Wazowski").Run()
				return fn
			},
		},
		{
			"invalid",
			true,
			"",
			func() func() {
				return insideWorkTree()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gi := &Git{}
			fn := tt.init()
			defer fn()
			if err := gi.UnsetLocalUserName(); (err != nil) != tt.wantErr {
				t.Errorf("UnsetLocalUserName() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got := gi.GetLocalUserName(); got != tt.want {
				t.Errorf("UnsetLocalUserEmail() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGit_UnsetLocalUserSigningKey(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
		want    string
		init    func() func()
	}{
		{
			"valid",
			false,
			"",
			func() func() {
				fn := insideWorkTree()
				exec.Command("git", "config", "--local", "user.signingkey", "aaabbbcccddeeefff").Run()
				return fn
			},
		},
		{
			"invalid",
			true,
			"",
			func() func() {
				return insideWorkTree()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gi := &Git{}
			fn := tt.init()
			defer fn()
			if err := gi.UnsetLocalUserSigningKey(); (err != nil) != tt.wantErr {
				t.Errorf("UnsetLocalUserSigningKey() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got := gi.GetLocalUserSigningKey(); got != tt.want {
				t.Errorf("GetLocalUserSigningKey() got = %v, want %v", got, tt.want)
			}
		})
	}
}
