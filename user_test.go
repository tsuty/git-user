package main

import "testing"

func TestInitUser(t *testing.T) {
	expected := "a1cd4ef7fa69afc33cff820e70a447e4f83ff9a1"
	u := user{
		Url: "git@example.com:*",
		Name: "someone",
		Email: "someone@example.com",
		SigningKey: "AAABBBCCC",
	}
	u.init()

	if u.Id != expected {
		t.Errorf("expected: %s, actual: %s", expected, u.Id)
	}
}

func TestToStringWithUrlOption(t *testing.T) {
	expected := "git@example.com:*"

	u := user{
		Url: expected,
		Name: "someone",
		Email: "someone@example.com",
		SigningKey: "AAABBBCCC",
	}
	u.init()

	actual := u.toString(output{Url: true})
	if actual != expected {
		t.Errorf("expected: %s, actual: %s", expected, actual)
	}
}

func TestToStringWithNameOption(t *testing.T) {
	expected := "someone"

	u := user{
		Url: "git@example.com:*",
		Name: expected,
		Email: "someone@example.com",
		SigningKey: "AAABBBCCC",
	}
	u.init()

	actual := u.toString(output{Name: true})
	if actual != expected {
		t.Errorf("expected: %s, actual: %s", expected, actual)
	}
}

func TestToStringWithEmailOption(t *testing.T) {
	expected := "someone@example.com"

	u := user{
		Url: "git@example.com:*",
		Name: "someone",
		Email: expected,
		SigningKey: "AAABBBCCC",
	}
	u.init()

	actual := u.toString(output{Email: true})
	if actual != expected {
		t.Errorf("expected: %s, actual: %s", expected, actual)
	}
}

func TestToStringWithSigningKeyOption(t *testing.T) {
	expected := "AAABBBCCC"

	u := user{
		Url: "git@example.com:*",
		Name: "someone",
		Email: "someone@example.com",
		SigningKey: expected,
	}
	u.init()

	actual := u.toString(output{SigningKey: true})
	if actual != expected {
		t.Errorf("expected: %s, actual: %s", expected, actual)
	}
}

func TestToStringWithoutOption(t *testing.T) {
	expected1 := "Url: git@example.com:*\tName: someone\tEmail: someone@example.com\tSigningKey: AAABBBCCC\tID: a1cd4ef7fa69afc33cff820e70a447e4f83ff9a1"
	u1 := user{
		Url: "git@example.com:*",
		Name: "someone",
		Email: "someone@example.com",
		SigningKey: "AAABBBCCC",
	}
	u1.init()

	actual1 := u1.toString(output{})
	if actual1 != expected1 {
		t.Errorf("expected: %s, actual: %s", expected1, actual1)
		return
	}

	expected2 := "Url: git@example.com:*\tName: someone\tEmail: someone@example.com\tSigningKey: AAABBBCCC"
	u2 := user{
		Url: "git@example.com:*",
		Name: "someone",
		Email: "someone@example.com",
		SigningKey: "AAABBBCCC",
	}

	actual2 := u2.toString(output{})
	if actual2 != expected2 {
		t.Errorf("expected: %s, actual: %s", expected2, actual2)
	}

}
