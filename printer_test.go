package main

import (
	"bytes"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestNewPrinter(t *testing.T) {
	type args struct {
		flag   uint
		writer io.Writer
	}
	tests := []struct {
		name string
		args args
		want *Printer
	}{
		{
			"with PrintALL flag",
			args{flag: PrintALL, writer: &bytes.Buffer{}},
			&Printer{flag: PrintALL, writer: &bytes.Buffer{}},
		},
		{
			"with PrintDefault flag",
			args{flag: PrintDefault, writer: &bytes.Buffer{}},
			&Printer{flag: PrintALL, writer: &bytes.Buffer{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			got := NewPrinter(tt.args.flag, writer)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPrinter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrinter_PrintUser(t *testing.T) {
	type fields struct {
		flag uint
	}
	type args struct {
		user *User
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantWriter string
	}{
		{
			"with PrintALL",
			fields{flag: PrintALL},
			args{
				&User{
					URL:        "git@example.com:c/d",
					Name:       "Mike Wazowski",
					Email:      "mike@example.com",
					SigningKey: "",
				},
			},
			"URL: git@example.com:c/d  Name: Mike Wazowski  Email: mike@example.com  SigningKey:  Hash: ",
		},
		{
			"with PrintURL",
			fields{flag: PrintURL},
			args{
				&User{
					URL:        "git@example.com:e/f",
					Name:       "Mike Wazowski",
					Email:      "mike@example.com",
					SigningKey: "",
				},
			},
			"URL: git@example.com:e/f",
		},
		{
			"with PrintEmail",
			fields{flag: PrintEmail},
			args{
				&User{
					URL:        "git@example.com:e/f",
					Name:       "Mike Wazowski",
					Email:      "mike@example.com",
					SigningKey: "",
				},
			},
			"Email: mike@example.com",
		},
		{
			"with PrintName | PrintEmail",
			fields{flag: PrintName | PrintEmail},
			args{
				&User{
					URL:        "git@example.com:e/f",
					Name:       "Mike Wazowski",
					Email:      "mike@example.com",
					SigningKey: "",
				},
			},
			"Name: Mike Wazowski  Email: mike@example.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			p := Printer{
				flag:   tt.fields.flag,
				writer: writer,
			}
			p.PrintUser(tt.args.user)
			if got := writer.String(); !strings.HasPrefix(got, tt.wantWriter) {
				t.Errorf("PrintUser() write `%v`, want `%v`", got, tt.wantWriter)
			}
		})
	}
}

func TestPrinter_PrintUsers(t *testing.T) {
	type fields struct {
		flag uint
	}
	type args struct {
		users []*User
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantWriter string
	}{
		{
			"with PrintURL | PrintName | PrintEmail",
			fields{flag: PrintURL | PrintName | PrintEmail},
			args{
				Users{
					&User{
						URL:        "git@example.com:a/b",
						Name:       "Mike Wazowski",
						Email:      "mike@example.com",
						SigningKey: "",
					},
					&User{
						URL:        "git@example.com:c/d",
						Name:       "James Phil. Sullivan",
						Email:      "sulley@example.com",
						SigningKey: "",
					},
				},
			},
			`URL: git@example.com:a/b  Name: Mike Wazowski         Email: mike@example.com
URL: git@example.com:c/d  Name: James Phil. Sullivan  Email: sulley@example.com
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			p := Printer{
				flag:   tt.fields.flag,
				writer: writer,
			}
			p.PrintUsers(tt.args.users)
			if got := writer.String(); got != tt.wantWriter {
				t.Errorf("PrintUsers() write `%v`, want `%v`", got, tt.wantWriter)
			}
		})
	}
}

func TestPrinter_Printf(t *testing.T) {
	type args struct {
		format string
		a      []interface{}
	}
	tests := []struct {
		name       string
		args       args
		wantWriter string
	}{
		{
			"with format",
			args{"%d %d %d", []interface{}{1, 2, 3}},
			"1 2 3",
		},
		{
			"outwith format",
			args{"a b c", []interface{}{}},
			"a b c",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			p := Printer{
				flag:   PrintALL,
				writer: writer,
			}
			p.Printf(tt.args.format, tt.args.a...)
			if got := writer.String(); got != tt.wantWriter {
				t.Errorf("Printf() write %v, want %v", got, tt.wantWriter)
			}
		})
	}
}

func TestPrinter_Println(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name       string
		args       args
		wantWriter string
	}{
		{
			"with some message",
			args{"This is a pen."},
			"This is a pen.\n",
		},
		{
			"without any message",
			args{""},
			"\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			p := Printer{
				flag:   PrintALL,
				writer: writer,
			}
			p.Println(tt.args.message)
			if got := writer.String(); got != tt.wantWriter {
				t.Errorf("Println() write `%v`, want `%v`", got, tt.wantWriter)
			}
		})
	}
}
