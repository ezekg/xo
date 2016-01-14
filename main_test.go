package main

import (
	"os/exec"
	"testing"
)

func TestMain(t *testing.T) {
	testCommand(t, `echo 'Hello there!' | xo 'i/hello(.*)/Hi$1/'`,
		`Hi there!
`)
	testCommand(t, `echo 'Hello! - Bob' | xo 'i/(hello).*?-.*?(\w+)/Why $1, $2!/'`,
		`Why Hello, Bob!
`)
	testCommand(t, `cat fixtures/servers.yml | xo 'mis/.*?(production):\s*server:\s+([^:\n]+):?(\d+)?.*?user:\s+([^\n]+).*/$4@$2 -p $3?:22/'`,
		`user-1@192.168.1.1 -p 1234
`)
	testCommand(t, `cat fixtures/starwars.txt | xo 'mi/^(\w+):(\s*\[(.*?)\]\s*)?\s*([^\n]+)/$1 said, "$4" in a $3?:normal voice./'`,
		`Vader said, "If only you knew the power of the Dark Side. Obi-Wan never told you what happened to your father." in a normal voice.
Luke said, "He told me enough! He told me you killed him!" in a normal voice.
Vader said, "No, I am your father." in a normal voice.
Luke said, "No. No! That's not true! That's impossible!" in a shocked voice.
`)
	testCommand(t, `echo '123' | xo 'i/(\d)(\d)(\d)(\d)?(\d)?/$1, $2, $3, 4?:FOUR $5?:FIVE/'`,
		`1, 2, 3, 4?:FOUR FIVE
`)
}

func testCommand(t *testing.T, cmd string, expected string) {
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		t.Fatalf(err.Error())
	}
	result := string(out)
	if result != expected {
		t.Fatalf("'%s' should be '%s'", result, expected)
	}
}
