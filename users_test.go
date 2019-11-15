package main

import (
	"reflect"
	"sort"
	"testing"
)

func TestUser_sorting(t *testing.T) {
	tests := []struct {
		name             string
		users            Users
		wantFirstUserURL string
		wantLastUserURL  string
	}{
		{
			"same url",
			Users{
				&User{URL: "git@example.com:foo/bar"},
				&User{URL: "git@example.com:foo/bar"},
			},
			"git@example.com:foo/bar",
			"git@example.com:foo/bar",
		},
		{
			"different url",
			Users{
				&User{URL: "git@example.com:*"},
				&User{URL: "git@example.com:foo/bar"},
				&User{URL: "git@example.com:foo/*"},
			},
			"git@example.com:foo/bar",
			"git@example.com:*",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users := tt.users
			sort.Sort(users)
			first := users[0]
			last := users[len(users)-1]
			if !reflect.DeepEqual(first.URL, tt.wantFirstUserURL) {
				t.Errorf("sorting first user URL is %v, want %v", first.URL, tt.wantFirstUserURL)
			}
			if !reflect.DeepEqual(last.URL, tt.wantLastUserURL) {
				t.Errorf("sorting last user URL is %v, want %v", last.URL, tt.wantLastUserURL)
			}
		})
	}
}

func TestUser_Hash(t *testing.T) {
	type fields struct {
		URL        string
		Name       string
		Email      string
		SigningKey string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
		right  bool
	}{
		{
			"same field same hash",
			fields{Name: "Mike Wazowski"},
			(&User{Name: "Mike Wazowski"}).Hash(),
			true,
		},
		{
			"different field other hash",
			fields{Name: "Mike Wazowski"},
			(&User{Name: "James Phil. Sullivan"}).Hash(),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				URL:        tt.fields.URL,
				Name:       tt.fields.Name,
				Email:      tt.fields.Email,
				SigningKey: tt.fields.SigningKey,
			}
			if got := u.Hash(); (got != tt.want) == tt.right {
				t.Errorf("Hash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsers_FindByURL(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		us   Users
		args args
		want *User
	}{
		{
			"completely match",
			Users{
				&User{URL: "git@github.com:tsuty/git-user.git"},
			},
			args{"git@github.com:tsuty/git-user.git"},
			&User{URL: "git@github.com:tsuty/git-user.git"},
		},
		{
			"wild card match",
			Users{
				&User{URL: "git@github.com:tsuty/*"},
			},
			args{"git@github.com:tsuty/git-user.git"},
			&User{URL: "git@github.com:tsuty/*"},
		},
		{
			"wild card but completely match",
			Users{
				&User{URL: "git@github.com:tsuty/git-user.git", Name: "Completely"},
				&User{URL: "git@github.com:tsuty/*", Name: "Wild"},
			},
			args{"git@github.com:tsuty/git-user.git"},
			&User{URL: "git@github.com:tsuty/git-user.git", Name: "Completely"},
		},
		{
			"not found",
			Users{
				&User{URL: "git@github.com:tsuty/git-user.git", Name: "Completely"},
				&User{URL: "git@github.com:tsuty/*", Name: "Wild"},
			},
			args{"git@github.com:fooo/git-user.git"},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.us.TakeByURL(tt.args.url); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TakeByURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsers_Set(t *testing.T) {
	type args struct {
		url        string
		name       string
		email      string
		signingkey string
	}
	tests := []struct {
		name string
		us   Users
		args args
		want *User
		len  int
	}{
		{
			"append user",
			Users{},
			args{
				url:   "git@github.com:tsuty/git-user.git",
				name:  "tsuty",
				email: "tsuty@example.com",
			},
			&User{
				Name:  "tsuty",
				Email: "tsuty@example.com",
				URL:   "git@github.com:tsuty/git-user.git",
			},
			1,
		},
		{
			"update user",
			Users{
				&User{
					Name:  "tsuty",
					Email: "tsuty@example.com",
					URL:   "git@github.com:tsuty/git-user.git",
				},
			},
			args{
				url:   "git@github.com:tsuty/git-user.git",
				name:  "tsuty",
				email: "tsuty@subdomain.example.com",
			},
			&User{
				Name:  "tsuty",
				Email: "tsuty@subdomain.example.com",
				URL:   "git@github.com:tsuty/git-user.git",
			},
			1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.us.Set(tt.args.url, tt.args.name, tt.args.email, tt.args.signingkey); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Set() = %v, want %v", got, tt.want)
			}
			if got := len(tt.us); got != tt.len {
				t.Errorf("len(Users) = %v, want %v", got, tt.len)
			}
		})
	}
}

func TestUsers_TakeByHash(t *testing.T) {
	exampleUser := &User{
		Name:  "tsuty",
		Email: "tsuty@example.com",
		URL:   "git@github.com:tsuty/git-user.git",
	}
	exampleHash := exampleUser.Hash()

	type args struct {
		hash string
	}
	tests := []struct {
		name string
		us   Users
		args args
		want *User
	}{
		{
			"match",
			Users{
				exampleUser,
				&User{
					Name:  "Mike Wazowski",
					Email: "Mike@example.com",
					URL:   "git@example.com:monsters/inc.git",
				},
			},
			args{hash: exampleHash},
			exampleUser,
		},
		{
			"no match",
			Users{
				exampleUser,
				&User{
					Name:  "Mike Wazowski",
					Email: "Mike@example.com",
					URL:   "git@example.com:monsters/inc.git",
				},
			},
			args{hash: "abcedfg"},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.us.TakeByHash(tt.args.hash); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TakeByHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsers_Delete(t *testing.T) {
	exampleUser := &User{
		Name:  "tsuty",
		Email: "tsuty@example.com",
		URL:   "git@github.com:tsuty/git-user.git",
	}
	exampleHash := exampleUser.Hash()

	type args struct {
		hash string
	}
	tests := []struct {
		name       string
		us         Users
		args       args
		want       *User
		wantLength int
	}{
		{
			"match",
			Users{
				exampleUser,
				&User{
					Name:  "Mike Wazowski",
					Email: "mike@example.com",
					URL:   "git@example.com:monsters/inc.git",
				},
			},
			args{hash: exampleHash},
			exampleUser,
			1,
		},
		{
			"no match",
			Users{
				exampleUser,
				&User{
					Name:  "Mike Wazowski",
					Email: "mike@example.com",
					URL:   "git@example.com:monsters/inc.git",
				},
			},
			args{hash: "abcedfg"},
			nil,
			2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.us.Delete(tt.args.hash); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Delete() = %v, want %v", got, tt.want)
			}
			if got := len(tt.us); got != tt.wantLength {
				t.Errorf("Len(Users) = %v, want %v", got, tt.want)
			}
		})
	}
}
