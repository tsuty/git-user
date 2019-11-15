package main

import "testing"

func TestSetArgs_Valid(t *testing.T) {
	type fields struct {
		Name       string
		Email      string
		SigningKey string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"fulfill field",
			fields{Name: "Mike Wazowski", Email: "mike@example.com", SigningKey: "aaabbbccc"},
			false,
		},
		{
			"fill Name and Email",
			fields{Name: "Mike Wazowski", Email: "mike@example.com"},
			false,
		},
		{
			"fill Email and SigningKey",
			fields{Email: "mike@example.com", SigningKey: "aaabbbccc"},
			true,
		},
		{
			"fill Name and SigningKey",
			fields{Name: "Mike Wazowski", SigningKey: "aaabbbccc"},
			true,
		},
		{
			"only Name",
			fields{Name: "Mike Wazowski"},
			true,
		},
		{
			"only Email",
			fields{Email: "mike@example.com"},
			true,
		},
		{
			"only SigningKey",
			fields{SigningKey: "aaabbbccc"},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			as := SetArgs{
				Name:       tt.fields.Name,
				Email:      tt.fields.Email,
				SigningKey: tt.fields.SigningKey,
			}
			if err := as.Valid(); (err != nil) != tt.wantErr {
				t.Errorf("Valid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_printOption_printFlag(t *testing.T) {
	type fields struct {
		URL        bool
		Name       bool
		Email      bool
		SigningKey bool
	}
	tests := []struct {
		name   string
		fields fields
		want   uint
	}{
		{
			"all false",
			fields{false, false, false, false},
			PrintDefault,
		},
		{
			"all true",
			fields{true, true, true, true},
			PrintALL,
		},
		{
			"URL true",
			fields{true, false, false, false},
			PrintURL,
		},
		{
			"Name false",
			fields{false, true, false, false},
			PrintName,
		},
		{
			"Email false",
			fields{false, false, true, false},
			PrintEmail,
		},
		{
			"SigningKey false",
			fields{false, false, false, true},
			PrintSigningKey,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := printOption{
				URL:        tt.fields.URL,
				Name:       tt.fields.Name,
				Email:      tt.fields.Email,
				SigningKey: tt.fields.SigningKey,
			}
			if got := o.printFlag(); got|tt.want != tt.want {
				t.Errorf("printFlag() = %v, want %v", got, tt.want)
			}
		})
	}
}
