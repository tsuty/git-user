package main

import (
	"testing"
	"bytes"
)

func TestBashCompletion(t *testing.T) {
	output := new(bytes.Buffer)

	app := gitUserApp()
	app.Writer = output
	app.Run([]string{"foo", "--generate-bash-completion"})

	expected := "show\nset\ndelete\nlocal\nlist\nsync\n"
	actual := output.String()

	if actual != expected {
		t.Errorf("expected: %s, actual: %s", expected, actual)
	}
}


